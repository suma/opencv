#include "instance_manager_bridge.h"
#include "util.hpp"
#include <scouter-core/tracking_result.hpp>

struct ByteArray InstanceState_Serialize(InstanceState s) {
  return serializeObject(*s);
}

InstanceState InstanceState_Deserialize(struct ByteArray src) {
  return deserializeObject<scouter::InstanceState>(src);
}

void InstanceState_Delete(InstanceState s) {
  delete s;
}

struct InstanceStates InvertInstanceStates(InstanceState* obj, int length) {
  scouter::InstanceState** iss = new scouter::InstanceState*[length];
  for (int i = 0; i < length; ++i) {
    iss[i] = obj[i];
  }
  InstanceStates is = {iss, length};
  return is;
}

void InstanceStates_Delete(struct InstanceStates instanceStates) {
  delete instanceStates.instanceStates;
}

String InstanceState_ToJSON(struct InstanceStates instanceStates, int floorID,
    unsigned long long timestamp) {
  using pfi::text::json::json;
  using pfi::text::json::json_object;
  using pfi::text::json::json_array;
  using pfi::text::json::json_integer;
  using pfi::text::json::json_string;

  json instances_json(new json_array);
  for (int i = 0; i < instanceStates.length; ++i) {
    const scouter::InstanceState& s = *(instanceStates.instanceStates[i]);
    json instance_json(new json_object);
    instance_json["id"] = json(new json_integer(s.id));
    instance_json["location"] = json(new json_object);
    instance_json["location"]["x"] = json(new json_integer(s.position.x));
    instance_json["location"]["y"] = json(new json_integer(s.position.y));
    instance_json["location"]["floor_id"] = json(new json_integer(floorID));
    instance_json["labels"] = json(new json_array);
    for (size_t j = 0; j < s.tags.size(); ++j) {
      const scouter::Tag& tag = s.tags[j];
      instance_json["labels"].add(json(new json_string(tag.key + "=" + tag.value)));
    }
    instances_json.add(instance_json);
  }

  json ret(new json_object);
  ret["time"] = json(new json_integer(timestamp));
  ret["instances"] = instances_json;

  // json to string
  std::stringstream ss;
  ss << ret;
  String str = {ss.str().c_str(), (int)ss.str().size()};
  return str;
}

InstanceManager InstanceManager_New(const char *config) {
  scouter::InstanceManager::Config ic =
    load_json<scouter::InstanceManager::Config>(config);
  return new scouter::InstanceManager(ic);
}

void InstanceManager_Delete(InstanceManager instanceManager) {
  delete instanceManager;
}

struct InstanceStates TrackAndGetStates(Tracker tracker, InstanceManager im) {
  scouter::TrackingResult result = tracker->track();
  im->update(result);
  std::vector<scouter::InstanceState> states = im->get_current_states();

  scouter::InstanceState** ret = new scouter::InstanceState*[states.size()];
  for (size_t i = 0; i < states.size(); ++i) {
    ret[i] = new scouter::InstanceState(states[i]);
  }

  InstanceStates iss = {ret, (int)states.size()};
  return iss;
}
