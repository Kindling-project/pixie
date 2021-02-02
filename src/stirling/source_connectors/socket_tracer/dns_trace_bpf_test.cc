#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <string>

#include "src/common/base/base.h"
#include "src/common/exec/exec.h"
#include "src/common/testing/test_utils/container_runner.h"
#include "src/stirling/core/output.h"
#include "src/stirling/source_connectors/socket_tracer/testing/socket_trace_bpf_test_fixture.h"
#include "src/stirling/testing/common.h"

namespace pl {
namespace stirling {

using ::pl::stirling::testing::ColWrapperSizeIs;
using ::pl::stirling::testing::FindRecordsMatchingPID;
using ::pl::stirling::testing::SocketTraceBPFTest;
using ::pl::testing::BazelBinTestFilePath;

using ::testing::Each;
using ::testing::MatchesRegex;

// A DNS server using the bind9 DNS server image.
class DNSServerContainer : public ContainerRunner {
 public:
  DNSServerContainer()
      : ContainerRunner(BazelBinTestFilePath(kBazelImageTar), kInstanceNamePrefix, kReadyMessage) {}

 private:
  static constexpr std::string_view kBazelImageTar =
      "src/stirling/testing/dns2/dns_server_image.tar";
  static constexpr std::string_view kInstanceNamePrefix = "dns_server";
  static constexpr std::string_view kReadyMessage = "all zones loaded";
};

class DNSTraceTest : public SocketTraceBPFTest</* TClientSideTracing */ true> {
 protected:
  DNSTraceTest() {
    // Run the bind DNS server.
    // The container runner will make sure it is in the ready state before unblocking.
    // Stirling will run after this unblocks, as part of SocketTraceBPFTest SetUp().
    // Note that this step will make an access to docker hub to download the bind image.
    PL_CHECK_OK(container_.Run(150, {}));
  }

  DNSServerContainer container_;
};

//-----------------------------------------------------------------------------
// Test Scenarios
//-----------------------------------------------------------------------------

TEST_F(DNSTraceTest, Capture) {
  // Sleep an additional second, just to be safe.
  sleep(1);

  // Uncomment to enable tracing:
  FLAGS_stirling_conn_trace_pid = container_.process_pid();

  // Run dig to generate a DNS request.
  // Run it through bash, and return the PID, so we can use it to filter captured results.
  std::string cmd =
      absl::StrFormat("docker exec %s sh -c 'dig @127.0.0.1 server.dnstest.com & echo $! && wait'",
                      container_.container_name());
  ASSERT_OK_AND_ASSIGN(std::string out, pl::Exec(cmd));
  LOG(INFO) << out;

  // Grab the data from Stirling.
  DataTable data_table(kDNSTable);
  source_->TransferData(ctx_.get(), SocketTraceConnector::kDNSTableNum, &data_table);
  std::vector<TaggedRecordBatch> tablets = data_table.ConsumeRecords();
  ASSERT_FALSE(tablets.empty());

  types::ColumnWrapperRecordBatch rb = tablets[0].records;
  PrintRecordBatch("dns", kDNSTable.ToProto(), rb);

  // Check server-side.
  {
    types::ColumnWrapperRecordBatch records =
        FindRecordsMatchingPID(tablets[0].records, kDNSUPIDIdx, container_.process_pid());

    ASSERT_THAT(records, Each(ColWrapperSizeIs(1)));

    const std::string& req_hdr = records[kDNSReqHdrIdx]->Get<types::StringValue>(0);
    const std::string& req_body = records[kDNSReqBodyIdx]->Get<types::StringValue>(0);
    const std::string& resp_hdr = records[kDNSRespHdrIdx]->Get<types::StringValue>(0);
    const std::string& resp_body = records[kDNSRespBodyIdx]->Get<types::StringValue>(0);

    EXPECT_THAT(
        req_hdr,
        MatchesRegex(
            R"(\{"txid":[0-9]+,"qr":0,"opcode":0,"aa":0,"tc":0,"rd":1,"ra":0,"ad":1,"cd":0,"rcode":0,)"
            R"("num_queries":1,"num_answers":0,"num_auth":0,"num_addl":1\})"));
    EXPECT_EQ(req_body, R"({"queries":[{"name":"server.dnstest.com","type":"A"}]})");
    EXPECT_THAT(
        resp_hdr,
        MatchesRegex(
            R"(\{"txid":[0-9]+,"qr":1,"opcode":0,"aa":1,"tc":0,"rd":1,"ra":1,"ad":0,"cd":0,"rcode":0,)"
            R"("num_queries":1,"num_answers":1,"num_auth":0,"num_addl":1\})"));
    EXPECT_EQ(resp_body,
              R"({"answers":[{"name":"server.dnstest.com","type":"A","addr":"192.168.32.200"}]})");
  }
}

}  // namespace stirling
}  // namespace pl