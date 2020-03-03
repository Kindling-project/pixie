package bridge

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/nats-io/nats.go"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc/metadata"
	"pixielabs.ai/pixielabs/src/cloud/vzconn/vzconnpb"
	"pixielabs.ai/pixielabs/src/shared/cvmsgspb"
	utils2 "pixielabs.ai/pixielabs/src/shared/services/utils"
	"pixielabs.ai/pixielabs/src/utils"
	certmgrpb "pixielabs.ai/pixielabs/src/vizier/services/certmgr/certmgrpb"
)

const heartbeatIntervalS = 5 * time.Second

// HeartbeatTopic is the topic that heartbeats are written to.
const HeartbeatTopic = "heartbeat"
const natsCloudUpdateTopic = "v2c"
const registrationTimeout = 30 * time.Second

// ErrRegistrationTimeout is the registration timeout error.
var ErrRegistrationTimeout = errors.New("Registration timeout")

// VizierInfo fetches information about Vizier.
type VizierInfo interface {
	GetAddress() (string, int32, error)
}

// Bridge is the NATS<->GRPC bridge.
type Bridge struct {
	vizierID      uuid.UUID
	jwtSigningKey string
	sessionID     int64

	vzConnClient  vzconnpb.VZConnServiceClient
	certMgrClient certmgrpb.CertMgrServiceClient

	vzInfo VizierInfo

	hbSeqNum int64

	nc         *nats.Conn
	natsCh     chan *nats.Msg
	registered bool
	// There are a two sets of streams that we manage for the GRPC side. The incoming
	// data and the outgoing data. GRPC does not natively provide a channel based interface
	// so we wrap the Send/Recv calls with goroutines that are responsible for
	// performing the read/write operations.
	//
	// If the GRPC connection gets disrupted, we close all the readers and writer routines, but we leave the
	// channels in place so that data does not get lost. The data will simply be resent
	// once the connection comes back alive. If data is lost due to a crash, the rest of the system
	// is resilient to this loss, but reducing it is optimal to prevent a lot of replay traffic.
	grpcOutCh chan *vzconnpb.V2CBridgeMessage
	grpcInCh  chan *vzconnpb.C2VBridgeMessage
	// This tracks the message we are trying to send, but has not been sent yet.
	pendingGRPCOutMsg *vzconnpb.V2CBridgeMessage

	quitCh chan bool      // Channel is used to signal that things should shutdown.
	wg     sync.WaitGroup // Tracks all the active goroutines.
}

// New creates a cloud connector to cloud bridge.
func New(vizierID uuid.UUID, jwtSigningKey string, sessionID int64, vzClient vzconnpb.VZConnServiceClient, certMgrClient certmgrpb.CertMgrServiceClient, vzInfo VizierInfo, nc *nats.Conn) *Bridge {
	return &Bridge{
		vizierID:          vizierID,
		jwtSigningKey:     jwtSigningKey,
		sessionID:         sessionID,
		vzConnClient:      vzClient,
		certMgrClient:     certMgrClient,
		vzInfo:            vzInfo,
		hbSeqNum:          0,
		nc:                nc,
		natsCh:            make(chan *nats.Msg),
		registered:        false,
		grpcOutCh:         make(chan *vzconnpb.V2CBridgeMessage),
		grpcInCh:          make(chan *vzconnpb.C2VBridgeMessage),
		pendingGRPCOutMsg: nil,
		quitCh:            make(chan bool),
		wg:                sync.WaitGroup{},
	}
}

// RunStream manages starting and restarting the stream to VZConn.
func (s *Bridge) RunStream() {
	natsTopic := fmt.Sprintf("%s.*", natsCloudUpdateTopic)
	log.WithField("topic", natsTopic).Trace("Subscribing to NATS")
	natsSub, err := s.nc.ChanSubscribe(natsTopic, s.natsCh)
	if err != nil {
		log.WithError(err).Fatal("Failed to subscribe to NATS.")
	}
	defer natsSub.Unsubscribe()

	for {
		s.registered = false
		select {
		case <-s.quitCh:
			return
		default:
			log.Trace("Starting stream")
			errCh := make(chan error)
			err := s.StartStream(errCh)
			if err == nil {
				log.Trace("Stream ending")
			} else {
				log.WithError(err).Error("Stream errored. Restarting stream")
			}
			close(errCh)
		}
	}
}

func (s *Bridge) doRegistrationHandshake(stream vzconnpb.VZConnService_NATSBridgeClient) error {
	addr, _, err := s.vzInfo.GetAddress()
	if err != nil {
		log.WithError(err).Error("Unable to get vizier proxy address")
	}

	// Send over a registration request and wait for ACK.
	regReq := &cvmsgspb.RegisterVizierRequest{
		VizierID: utils.ProtoFromUUID(&s.vizierID),
		JwtKey:   s.jwtSigningKey,
		Address:  addr,
	}

	err = s.publishBridgeSync(stream, "register", regReq)
	if err != nil {
		return err
	}

	for {
		select {
		case <-time.After(registrationTimeout):
			log.Error("Timeout with registration terminating stream")
			stream.CloseSend()
			return ErrRegistrationTimeout
		case resp := <-s.grpcInCh:
			// Try to receive the registerAck.
			if resp.Topic != "registerAck" {
				log.Error("Unexpected message type while waiting for ACK")
			}
			registerAck := &cvmsgspb.RegisterVizierAck{}
			err = types.UnmarshalAny(resp.Msg, registerAck)
			if err != nil {
				return err
			}
			switch registerAck.Status {
			case cvmsgspb.ST_FAILED_NOT_FOUND:
				return errors.New("registration not found, cluster unknown in pixie-cloud")
			case cvmsgspb.ST_OK:
				s.registered = true
				return nil
			default:
				return errors.New("registration unsuccessful: " + err.Error())
			}
		}
	}
}

// StartStream starts the stream between the cloud connector and Vizier connector.
func (s *Bridge) StartStream(errCh chan error) error {
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := s.vzConnClient.NATSBridge(ctx)
	if err != nil {
		log.WithError(err).Error("Error starting stream")
		return err
	}
	// Wait for  all goroutines to terminate.
	defer func() {
		s.wg.Wait()
	}()

	// Setup the stream reader go routine.
	done := make(chan bool)
	defer close(done)
	// Cancel the stream to make sure everything get shutdown properly.
	defer func() {
		cancel()
	}()

	s.wg.Add(1)
	go s.startStreamGRPCReader(stream, done, errCh)
	s.wg.Add(1)
	go s.startStreamGRPCWriter(stream, done, errCh)

	if !s.registered {
		// Need to do registration handshake before we allow any cvmsgs.
		err := s.doRegistrationHandshake(stream)
		if err != nil {
			return err
		}
	}

	log.Trace("Registration Complete.")
	s.wg.Add(1)
	err = s.HandleNATSBridging(stream, done, errCh)

	return err
}

func (s *Bridge) startStreamGRPCReader(stream vzconnpb.VZConnService_NATSBridgeClient, done chan bool, errCh chan<- error) {
	defer s.wg.Done()
	log.Trace("Starting GRPC reader stream")
	defer log.Trace("Closing GRPC read stream")
	for {
		select {
		case <-stream.Context().Done():
			return
		case <-done:
			log.Info("Closing GRPC reader because of <-done")
			return
		default:
			log.Trace("Waiting for next message")
			msg, err := stream.Recv()
			if err != nil && err == io.EOF {
				log.Trace("Stream has closed(Read)")
				// stream closed.
				return
			}
			if err != nil {
				log.WithError(err).Error("Got a stream read error")
				return
			}
			s.grpcInCh <- msg
		}
	}
}

func (s *Bridge) startStreamGRPCWriter(stream vzconnpb.VZConnService_NATSBridgeClient, done chan bool, errCh chan<- error) {
	defer s.wg.Done()
	log.Trace("Starting GRPC writer stream")
	defer log.Trace("Closing GRPC writer stream")

	for {
		// Pending message try to send it first.
		if s.pendingGRPCOutMsg != nil {
			err := stream.Send(s.pendingGRPCOutMsg)
			if err != nil {
				// Error sending message. Second error aborts.
				errCh <- err
				return
			}
			s.pendingGRPCOutMsg = nil
		}
		select {
		case <-stream.Context().Done():
			log.Trace("Write stream has closed")
			return
		case <-done:
			log.Trace("Closing GRPC writer because of <-done")
			// Quit called.
			return
		case m := <-s.grpcOutCh:
			s.pendingGRPCOutMsg = m
			// Write message to GRPC if it exists.
			err := stream.Send(s.pendingGRPCOutMsg)
			if err != nil {
				// Need to resend this message.
				continue
			}
			s.pendingGRPCOutMsg = nil
		}
	}
}

// HandleNATSBridging routes message to and from cloud NATS.
func (s *Bridge) HandleNATSBridging(stream vzconnpb.VZConnService_NATSBridgeClient, done chan bool, errCh chan error) error {
	defer s.wg.Done()
	defer log.Info("Closing NATS Bridge")
	// Vizier -> Cloud side:
	// 1. Listen to NATS on cloudupdate.<topic>.
	// 2. Extract Topic from the stream name above.
	// 3. Wrap the message and throw it over the wire.

	// Cloud -> Vizier side:
	// 1. Read the stream.
	// 2. For cvmsgs of type: C2VBridgeMessage, read the topic
	//    and throw it onto nats under cloudresponse.topic

	log.Info("Starting NATS bridge.")
	hbChan := s.generateHeartbeats(done)
	s.wg.Add(1)
	go s.requestSSLCerts(stream.Context(), done)
	for {
		select {
		case <-s.quitCh:
			return nil
		case <-done:
			return nil
		case e := <-errCh:
			log.WithError(e).Error("GRPC error, terminating stream")
			return e
		case data := <-s.natsCh:
			if !strings.HasPrefix(data.Subject, natsCloudUpdateTopic+".") {
				return errors.New("invalid subject: " + data.Subject)
			}
			topic := strings.TrimPrefix(data.Subject, natsCloudUpdateTopic+".")
			// Message over nats should be wrapped in a V2CMessage.
			v2cMsg := &cvmsgspb.V2CMessage{}
			err := proto.Unmarshal(data.Data, v2cMsg)
			if err != nil {
				return err
			}
			err = s.publishBridgeCh(topic, v2cMsg.Msg)
			if err != nil {
				return err
			}
		case bridgeMsg := <-s.grpcInCh:
			log.Info("Got Message on GRPC channel")
			if bridgeMsg == nil {
				return nil
			}

			log.
				WithField("msg", bridgeMsg.String()).
				Trace("Got Message on GRPC channel")
			topic := fmt.Sprintf("c2v.%s", bridgeMsg.Topic)

			natsMsg := &cvmsgspb.C2VMessage{
				VizierID: s.vizierID.String(),
				Msg:      bridgeMsg.Msg,
			}
			b, err := natsMsg.Marshal()
			if err != nil {
				log.WithError(err).Error("Failed to marshal")
				return err
			}

			log.WithField("topic", topic).
				WithField("msg", natsMsg.String()).
				Info("Publishing to NATS")
			err = s.nc.Publish(topic, b)
			if err != nil {
				log.WithError(err).Error("Failed to publish")
				return err
			}
		case hbMsg := <-hbChan:
			err := s.publishProtoToBridgeCh(HeartbeatTopic, hbMsg)
			if err != nil {
				return err
			}
		case <-stream.Context().Done():
			log.Info("Stream has been closed, shutting down grpc readers")
			return nil
		}
	}
	return nil
}

// Stop terminates the server. Don't reuse this server object after stop has been called.
func (s *Bridge) Stop() {
	close(s.quitCh)
	// Wait fo all goroutines to stop.
	s.wg.Wait()
}
func (s *Bridge) publishBridgeCh(topic string, msg *types.Any) error {
	wrappedReq := &vzconnpb.V2CBridgeMessage{
		Topic:     topic,
		SessionId: s.sessionID,
		Msg:       msg,
	}
	s.grpcOutCh <- wrappedReq
	return nil
}

func (s *Bridge) publishProtoToBridgeCh(topic string, msg proto.Message) error {
	anyMsg, err := types.MarshalAny(msg)
	if err != nil {
		return err
	}

	return s.publishBridgeCh(topic, anyMsg)
}

func (s *Bridge) publishBridgeSync(stream vzconnpb.VZConnService_NATSBridgeClient, topic string, msg proto.Message) error {
	anyMsg, err := types.MarshalAny(msg)
	if err != nil {
		return err
	}

	wrappedReq := &vzconnpb.V2CBridgeMessage{
		Topic:     topic,
		SessionId: s.sessionID,
		Msg:       anyMsg,
	}

	if err := stream.Send(wrappedReq); err != nil {
		return err
	}
	return nil
}

func (s *Bridge) generateHeartbeats(done <-chan bool) (hbCh chan *cvmsgspb.VizierHeartbeat) {
	hbCh = make(chan *cvmsgspb.VizierHeartbeat)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(heartbeatIntervalS)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				log.Info("Stopping heartbeat routine")
				return
			case <-ticker.C:
				addr, port, err := s.vzInfo.GetAddress()
				if err != nil {
					log.WithError(err).Error("Failed to get vizier address")
				}
				hbCh <- &cvmsgspb.VizierHeartbeat{
					VizierID:       utils.ProtoFromUUID(&s.vizierID),
					Time:           time.Now().UnixNano(),
					SequenceNumber: s.hbSeqNum,
					Address:        addr,
					Port:           port,
				}
				s.hbSeqNum++
			}
		}
	}()
	return
}

// TODO(zasgar/michelle): Move this out of here.
// RequestAndHandleSSLCerts registers the cluster with VZConn.
func (s *Bridge) requestSSLCerts(ctx context.Context, done chan bool) error {
	defer s.wg.Done()
	log.Info("Requesting SSL certs")
	sslCh := make(chan *nats.Msg)
	sub, err := s.nc.ChanSubscribe("c2v.sslResp", sslCh)
	defer sub.Unsubscribe()

	// Send over a request for SSL certs.
	regReq := &cvmsgspb.VizierSSLCertRequest{
		VizierID: utils.ProtoFromUUID(&s.vizierID),
	}
	s.publishProtoToBridgeCh("ssl", regReq)
	sslResp := cvmsgspb.VizierSSLCertResponse{}

loop:
	for {
		select {
		case <-done:
			return nil
		case <-time.After(30 * time.Second):
			log.Error("Timeout waiting for SSL certs. Re-requesting")
			s.publishProtoToBridgeCh("ssl", regReq)
		case msg := <-sslCh:
			log.Trace("Got SSL message")
			envelope := &cvmsgspb.C2VMessage{}
			err := envelope.Unmarshal(msg.Data)
			if err != nil {
				// jump out and wait for timeout.
				log.WithError(err).Error("Got bad SSL response.")
				break
			}

			err = types.UnmarshalAny(envelope.Msg, &sslResp)
			if err != nil {
				log.WithError(err).Error("Got bad SSL response.")
				break
			}
			break loop
		}

	}

	certMgrReq := &certmgrpb.UpdateCertsRequest{
		Key:  sslResp.Key,
		Cert: sslResp.Cert,
	}

	serviceAuthToken, err := getServiceCredentials(viper.GetString("jwt_signing_key"))
	if err != nil {
		return err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization",
		fmt.Sprintf("bearer %s", serviceAuthToken))

	certMgrResp, err := s.certMgrClient.UpdateCerts(ctx, certMgrReq)
	if err != nil {
		return err
	}

	if !certMgrResp.OK {
		log.Error("Failed to update certs")
		return errors.New("Failed to update certs")
	}
	log.WithField("reply", certMgrResp.String()).Info("Certs updated")
	return nil
}

// TODO(zasgar/michelle): Remove this, we need to make this into cluster credentials.
// getServiceCredentials returns JWT credentials for inter-service requests.
func getServiceCredentials(signingKey string) (string, error) {
	claims := utils2.GenerateJWTForService("cloud connector Service")
	return utils2.SignJWTClaims(claims, signingKey)
}
