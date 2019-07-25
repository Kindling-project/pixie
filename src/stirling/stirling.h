#pragma once

#include <memory>
#include <string>
#include <unordered_map>

#include "src/common/base/base.h"
#include "src/stirling/proto/stirling.pb.h"
#include "src/stirling/source_registry.h"

namespace pl {
namespace stirling {

/**
 * @brief Convenience function to subscribe to all info classes of
 * a published proto message. This should actually be in an agent.
 * TODO(kgandhi): Move to agent or common utils for agent when available.
 *
 * @param publish_proto
 * @return stirlingpb::Subscribe
 */
stirlingpb::Subscribe SubscribeToAllInfoClasses(const stirlingpb::Publish& publish_proto);

/**
 * @brief Convenience function to subscribe to a single info classes of
 * a published proto message. This should actually be in an agent.
 * TODO(kgandhi): Move to agent or common utils for agent when available.
 *
 * @param publish_proto
 * @return stirlingpb::Subscribe
 */
stirlingpb::Subscribe SubscribeToInfoClass(const stirlingpb::Publish& publish_proto,
                                           std::string_view name);

/**
 * @brief Creates a registry of all the available source connectors, including
 *        source connectors used in development only.
 *        They will all be intialized, and consume memory, even if they are disabled.
 *
 * @param unique_ptr to the created registry.
 */
std::unique_ptr<SourceRegistry> CreateAllSourceRegistry();

/**
 * @brief Creates a registry of main production source connectors.
 *        Does not include registries used for development (e.g. SeqGenConnector).
 *
 * @param unique_ptr to the created registry.
 */
std::unique_ptr<SourceRegistry> CreateProdSourceRegistry();

/**
 * The data collector collects data from various different 'sources',
 * and makes them available via a structured API, where the data can then be used and queried as
 * needed (by Pixie or others). Its function is to unify various, disparate sources of data into a
 * common, structured data format.
 */
class Stirling : public NotCopyable {
 public:
  Stirling() = default;
  virtual ~Stirling() = default;

  /**
   * @brief Create a Stirling object
   * Factory method to create Stirling with a default registry containing
   * all sources
   *
   * @return std::unique_ptr<Stirling>
   */
  static std::unique_ptr<Stirling> Create();

  /**
   * @brief Create a Stirling object
   * Factory method to create Stirling with a source registry.
   *
   * @param registry
   * @return std::unique_ptr<Stirling>
   */
  static std::unique_ptr<Stirling> Create(std::unique_ptr<SourceRegistry> registry);

  /**
   * @brief Populate the Publish Proto object. Agent calls this function to get the Publish
   * proto message. The proto publish message contains information (InfoClassSchema) on
   * all the Source Connectors that can be run to gather data and information on the types
   * for the data. The agent can then subscribe to a subset of the published message. The proto
   * is defined in //src/stirling/proto/stirling.proto.
   *
   */
  virtual void GetPublishProto(stirlingpb::Publish* publish_pb) = 0;

  /**
   * @brief Get the Subscription object. Receive a Subscribe proto message from the agent.
   * Update the schemas based on the subscription message. Generate the appropriate tables
   * that conform to subscription information.
   *
   * @param subscribe_proto
   * @return Status
   */
  virtual Status SetSubscription(const stirlingpb::Subscribe& subscribe_proto) = 0;

  /**
   * @brief Register call-back from Agent. Used to periodically send data.
   *
   * Function signature is:
   *   uint64_t table_id
   *   std::unique_ptr<ColumnWrapperRecordBatch> data
   */
  virtual void RegisterCallback(PushDataCallback f) = 0;

  /**
   * @brief Return a map of table ID to InfoClassManager names.
   */
  virtual std::unordered_map<uint64_t, std::string> TableIDToNameMap() const = 0;

  /**
   * @brief Main data collection call. This version blocks, so make sure to wrap a thread around it.
   */
  virtual void Run() = 0;

  /**
   * @brief Main data collection call. This version spawns off as an independent thread.
   */
  virtual Status RunAsThread() = 0;

  /**
   * @brief Wait for the running thread to terminate. Assumes previous call to RunThread().
   */
  virtual void WaitForThreadJoin() = 0;

  /**
   * @brief Stop Stirling data collection, and perform any final clean-up actions.
   * Blocking, so will only return once the main loop has stopped.
   *
   * If Stirling is managing the thread, it will wait for thread to exit.
   * If external agent is managing the thread, it will wait until the main loop has exited.
   *
   * Note: this should be called in the case of a signal (e.g. SIGINT, SIGTERM, etc.)
   * to clean-up BPF deployed resources.
   */
  virtual void Stop() = 0;
};

}  // namespace stirling
}  // namespace pl
