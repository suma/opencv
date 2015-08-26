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

void InstanceStates_Delete(struct InstanceStates instanceStates) {
  delete instanceStates.instanceStates;
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
