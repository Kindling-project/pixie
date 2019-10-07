#pragma once

#include <rapidjson/stringbuffer.h>
#include <rapidjson/writer.h>
#include <string>
#include <utility>
#include <vector>

#include "src/carnot/udf/registry.h"
#include "src/shared/metadata/metadata_state.h"
#include "src/shared/types/types.h"

namespace pl {
namespace carnot {
namespace funcs {
namespace metadata {

using ScalarUDF = pl::carnot::udf::ScalarUDF;
using FunctionContext = pl::carnot::udf::FunctionContext;

inline const pl::md::AgentMetadataState* GetMetadataState(FunctionContext* ctx) {
  DCHECK(ctx != nullptr);
  auto md = ctx->metadata_state();
  DCHECK(md != nullptr);
  return md;
}

class ASIDUDF : public ScalarUDF {
 public:
  types::Int64Value Exec(FunctionContext* ctx) {
    auto md = GetMetadataState(ctx);
    return md->asid();
  }
};

class UPIDToASIDUDF : public ScalarUDF {
 public:
  types::Int64Value Exec(FunctionContext*, types::UInt128Value upid_value) {
    return upid_value.High64() >> 32;
  }
};

class PodIDToPodNameUDF : public ScalarUDF {
 public:
  types::StringValue Exec(FunctionContext* ctx, types::StringValue pod_id) {
    auto md = GetMetadataState(ctx);

    const auto* pod_info = md->k8s_metadata_state().PodInfoByID(pod_id);
    if (pod_info != nullptr) {
      return absl::Substitute("$0/$1", pod_info->ns(), pod_info->name());
    }

    return "";
  }
};

class PodNameToPodIDUDF : public ScalarUDF {
 public:
  types::StringValue Exec(FunctionContext* ctx, types::StringValue pod_name) {
    auto md = GetMetadataState(ctx);
    return GetPodID(md, pod_name);
  }

  static types::StringValue GetPodID(const pl::md::AgentMetadataState* md,
                                     types::StringValue pod_name) {
    // This UDF expects the pod name to be in the format of "<ns>/<pod-name>".
    std::vector<std::string_view> name_parts = absl::StrSplit(pod_name, "/");
    if (name_parts.size() != 2) {
      return "";
    }

    auto pod_name_view = std::make_pair(name_parts[0], name_parts[1]);
    auto pod_id = md->k8s_metadata_state().PodIDByName(pod_name_view);

    return pod_id;
  }
};

class UPIDToContainerIDUDF : public ScalarUDF {
 public:
  types::StringValue Exec(FunctionContext* ctx, types::UInt128Value upid_value) {
    auto md = GetMetadataState(ctx);

    auto upid_uint128 = absl::MakeUint128(upid_value.High64(), upid_value.Low64());
    auto upid = md::UPID(upid_uint128);
    auto pid = md->GetPIDByUPID(upid);
    if (pid == nullptr) {
      return "";
    }
    return pid->cid();
  }
};

class UPIDToPodIDUDF : public ScalarUDF {
 public:
  types::StringValue Exec(FunctionContext* ctx, types::UInt128Value upid_value) {
    auto md = GetMetadataState(ctx);

    auto upid_uint128 = absl::MakeUint128(upid_value.High64(), upid_value.Low64());
    auto upid = md::UPID(upid_uint128);
    auto pid = md->GetPIDByUPID(upid);
    if (pid == nullptr) {
      return "";
    }
    auto container_info = md->k8s_metadata_state().ContainerInfoByID(pid->cid());
    if (container_info == nullptr) {
      return "";
    }
    return std::string(container_info->pod_id());
  }
};

class UPIDToPodNameUDF : public ScalarUDF {
 public:
  types::StringValue Exec(FunctionContext* ctx, types::UInt128Value upid_value) {
    auto md = GetMetadataState(ctx);

    auto upid_uint128 = absl::MakeUint128(upid_value.High64(), upid_value.Low64());
    auto upid = md::UPID(upid_uint128);
    auto pid = md->GetPIDByUPID(upid);
    if (pid == nullptr) {
      return "";
    }
    auto container_info = md->k8s_metadata_state().ContainerInfoByID(pid->cid());
    if (container_info == nullptr) {
      return "";
    }
    const auto* pod_info = md->k8s_metadata_state().PodInfoByID(container_info->pod_id());
    if (pod_info == nullptr) {
      return "";
    }
    return absl::Substitute("$0/$1", pod_info->ns(), pod_info->name());
  }
};

class ServiceIDToServiceNameUDF : public ScalarUDF {
 public:
  types::StringValue Exec(FunctionContext* ctx, types::StringValue service_id) {
    auto md = GetMetadataState(ctx);

    const auto* service_info = md->k8s_metadata_state().ServiceInfoByID(service_id);
    if (service_info != nullptr) {
      return absl::Substitute("$0/$1", service_info->ns(), service_info->name());
    }

    return "";
  }
};

class ServiceNameToServiceIDUDF : public ScalarUDF {
 public:
  types::StringValue Exec(FunctionContext* ctx, types::StringValue service_name) {
    auto md = GetMetadataState(ctx);
    // This UDF expects the service name to be in the format of "<ns>/<service-name>".
    std::vector<std::string_view> name_parts = absl::StrSplit(service_name, "/");
    if (name_parts.size() != 2) {
      return "";
    }

    auto service_name_view = std::make_pair(name_parts[0], name_parts[1]);
    auto service_id = md->k8s_metadata_state().ServiceIDByName(service_name_view);

    return service_id;
  }
};

class UPIDToK8S {
 public:
  static const pl::md::PodInfo* UPIDtoPod(const pl::md::AgentMetadataState* md,
                                          types::UInt128Value upid_value) {
    auto upid_uint128 = absl::MakeUint128(upid_value.High64(), upid_value.Low64());
    auto upid = md::UPID(upid_uint128);
    auto pid = md->GetPIDByUPID(upid);
    if (pid == nullptr) {
      return nullptr;
    }
    auto container_info = md->k8s_metadata_state().ContainerInfoByID(pid->cid());
    if (container_info == nullptr) {
      return nullptr;
    }
    auto pod_info = md->k8s_metadata_state().PodInfoByID(container_info->pod_id());
    return pod_info;
  }
  static types::StringValue StringifyVector(const std::vector<std::string>& vec) {
    if (vec.size() == 1) {
      return std::string(vec[0]);
    } else if (vec.size() > 1) {
      rapidjson::StringBuffer s;
      rapidjson::Writer<rapidjson::StringBuffer> writer(s);

      writer.StartArray();
      for (const auto& str : vec) {
        writer.String(str.c_str());
      }
      writer.EndArray();
      return s.GetString();
    }
    return "";
  }
};

/**
 * @brief Returns the service ids for services that are currently running.
 */
class UPIDToServiceIDUDF : public ScalarUDF {
 public:
  types::StringValue Exec(FunctionContext* ctx, types::UInt128Value upid_value) {
    auto md = GetMetadataState(ctx);
    auto pod_info = UPIDToK8S::UPIDtoPod(md, upid_value);
    if (pod_info == nullptr || pod_info->services().size() == 0) {
      return "";
    }
    std::vector<std::string> running_service_ids;
    for (const auto& service_id : pod_info->services()) {
      auto service_info = md->k8s_metadata_state().ServiceInfoByID(service_id);
      if (service_info == nullptr) {
        continue;
      }
      if (service_info->stop_time_ns() == 0) {
        running_service_ids.push_back(service_id);
      }
    }

    return UPIDToK8S::StringifyVector(running_service_ids);
  }
};

/**
 * @brief Returns the service names for services that are currently running.
 */
class UPIDToServiceNameUDF : public ScalarUDF {
 public:
  types::StringValue Exec(FunctionContext* ctx, types::UInt128Value upid_value) {
    auto md = GetMetadataState(ctx);
    auto pod_info = UPIDToK8S::UPIDtoPod(md, upid_value);
    if (pod_info == nullptr || pod_info->services().size() == 0) {
      return "";
    }
    std::vector<std::string> running_service_names;
    for (const auto& service_id : pod_info->services()) {
      auto service_info = md->k8s_metadata_state().ServiceInfoByID(service_id);
      if (service_info == nullptr) {
        continue;
      }
      if (service_info->stop_time_ns() == 0) {
        running_service_names.push_back(
            absl::Substitute("$0/$1", service_info->ns(), service_info->name()));
      }
    }
    return UPIDToK8S::StringifyVector(running_service_names);
  }
};

/**
 * @brief Returns the service names for the given pod ID.
 */
class PodIDToServiceNameUDF : public ScalarUDF {
 public:
  types::StringValue Exec(FunctionContext* ctx, types::StringValue pod_id) {
    auto md = GetMetadataState(ctx);

    const auto* pod_info = md->k8s_metadata_state().PodInfoByID(pod_id);
    if (pod_info == nullptr) {
      return "";
    }

    std::vector<std::string> running_service_names;
    for (const auto& service_id : pod_info->services()) {
      auto service_info = md->k8s_metadata_state().ServiceInfoByID(service_id);
      if (service_info == nullptr) {
        continue;
      }
      if (service_info->stop_time_ns() == 0) {
        running_service_names.push_back(
            absl::Substitute("$0/$1", service_info->ns(), service_info->name()));
      }
    }
    return UPIDToK8S::StringifyVector(running_service_names);
  }
};

/**
 * @brief Returns the service ids for the given pod ID.
 */
class PodIDToServiceIDUDF : public ScalarUDF {
 public:
  types::StringValue Exec(FunctionContext* ctx, types::StringValue pod_id) {
    auto md = GetMetadataState(ctx);

    const auto* pod_info = md->k8s_metadata_state().PodInfoByID(pod_id);
    if (pod_info == nullptr) {
      return "";
    }

    std::vector<std::string> running_service_ids;
    for (const auto& service_id : pod_info->services()) {
      auto service_info = md->k8s_metadata_state().ServiceInfoByID(service_id);
      if (service_info == nullptr) {
        continue;
      }
      if (service_info->stop_time_ns() == 0) {
        running_service_ids.push_back(service_id);
      }
    }
    return UPIDToK8S::StringifyVector(running_service_ids);
  }
};

/**
 * @brief Returns the service names for the given pod name.
 */
class PodNameToServiceNameUDF : public ScalarUDF {
 public:
  types::StringValue Exec(FunctionContext* ctx, types::StringValue pod_name) {
    auto md = GetMetadataState(ctx);

    // This UDF expects the pod name to be in the format of "<ns>/<pod-name>".
    std::vector<std::string_view> name_parts = absl::StrSplit(pod_name, "/");
    if (name_parts.size() != 2) {
      return "";
    }

    auto pod_name_view = std::make_pair(name_parts[0], name_parts[1]);
    auto pod_id = md->k8s_metadata_state().PodIDByName(pod_name_view);

    const auto* pod_info = md->k8s_metadata_state().PodInfoByID(pod_id);
    if (pod_info == nullptr) {
      return "";
    }

    std::vector<std::string> running_service_names;
    for (const auto& service_id : pod_info->services()) {
      auto service_info = md->k8s_metadata_state().ServiceInfoByID(service_id);
      if (service_info == nullptr) {
        continue;
      }
      if (service_info->stop_time_ns() == 0) {
        running_service_names.push_back(
            absl::Substitute("$0/$1", service_info->ns(), service_info->name()));
      }
    }
    return UPIDToK8S::StringifyVector(running_service_names);
  }
};

/**
 * @brief Returns the service ids for the given pod name.
 */
class PodNameToServiceIDUDF : public ScalarUDF {
 public:
  types::StringValue Exec(FunctionContext* ctx, types::StringValue pod_name) {
    auto md = GetMetadataState(ctx);

    // This UDF expects the pod name to be in the format of "<ns>/<pod-name>".
    std::vector<std::string_view> name_parts = absl::StrSplit(pod_name, "/");
    if (name_parts.size() != 2) {
      return "";
    }

    auto pod_name_view = std::make_pair(name_parts[0], name_parts[1]);
    auto pod_id = md->k8s_metadata_state().PodIDByName(pod_name_view);

    const auto* pod_info = md->k8s_metadata_state().PodInfoByID(pod_id);
    if (pod_info == nullptr) {
      return "";
    }

    std::vector<std::string> running_service_ids;
    for (const auto& service_id : pod_info->services()) {
      auto service_info = md->k8s_metadata_state().ServiceInfoByID(service_id);
      if (service_info == nullptr) {
        continue;
      }
      if (service_info->stop_time_ns() == 0) {
        running_service_ids.push_back(service_id);
      }
    }
    return UPIDToK8S::StringifyVector(running_service_ids);
  }
};

void RegisterMetadataOpsOrDie(pl::carnot::udf::ScalarUDFRegistry* registry);

class UPIDToStringUDF : public ScalarUDF {
 public:
  types::StringValue Exec(FunctionContext*, types::UInt128Value upid_value) {
    auto upid_uint128 = absl::MakeUint128(upid_value.High64(), upid_value.Low64());
    auto upid = md::UPID(upid_uint128);
    return upid.String();
  }
};

class UPIDToPIDUDF : public ScalarUDF {
 public:
  types::Int64Value Exec(FunctionContext*, types::UInt128Value upid_value) {
    auto upid_uint128 = absl::MakeUint128(upid_value.High64(), upid_value.Low64());
    auto upid = md::UPID(upid_uint128);
    return static_cast<int64_t>(upid.pid());
  }
};

class PodIDToPodStartTimeUDF : public ScalarUDF {
 public:
  types::Time64NSValue Exec(FunctionContext* ctx, types::StringValue pod_id) {
    auto md = GetMetadataState(ctx);
    const pl::md::PodInfo* pod_info = md->k8s_metadata_state().PodInfoByID(pod_id);
    if (pod_info == nullptr) {
      return 0;
    }
    return pod_info->start_time_ns();
  }
};

class PodNameToPodStartTimeUDF : public ScalarUDF {
 public:
  types::Time64NSValue Exec(FunctionContext* ctx, types::StringValue pod_name) {
    auto md = GetMetadataState(ctx);
    types::StringValue pod_id = PodNameToPodIDUDF::GetPodID(md, pod_name);
    const pl::md::PodInfo* pod_info = md->k8s_metadata_state().PodInfoByID(pod_id);
    if (pod_info == nullptr) {
      return 0;
    }
    return pod_info->start_time_ns();
  }
};

class PodNameToPodStatusUDF : public ScalarUDF {
 public:
  /**
   * @brief Gets the Pod status for a passed in pod.
   *
   * @param ctx: the function context
   * @param pod_name: the Value containing a pod name.
   * @return types::StringValue: the status of the pod.
   */
  types::StringValue Exec(FunctionContext* ctx, types::StringValue pod_name) {
    auto md = GetMetadataState(ctx);
    types::StringValue pod_id = PodNameToPodIDUDF::GetPodID(md, pod_name);
    const pl::md::PodInfo* pod_info = md->k8s_metadata_state().PodInfoByID(pod_id);
    if (pod_info == nullptr) {
      return "";
    }
    if (pod_info->stop_time_ns() != 0) {
      return "Terminated";
    }
    return "Running";
  }
};

}  // namespace metadata
}  // namespace funcs
}  // namespace carnot
}  // namespace pl
