#include "instance_manager_bridge.h"
#include "util.hpp"

struct ByteArray InstanceState_Serialize(InstanceState s) {
  return serializeObject(*s);
}

InstanceState InstanceState_Deserialize(struct ByteArray src) {
  return deserializeObject<scouter::InstanceState>(src);
}

void InstanceState_Delete(InstanceState s) {
  delete s;
}

void ResolveInstanceStates(struct InstanceStates instanceStates, InstanceState* obj) {
  for (size_t i = 0; i < instanceStates.instanceStateVec->size(); ++i) {
    obj[i] = new scouter::InstanceState((*instanceStates.instanceStateVec)[i]);
  }
  return;
}

void InstanceStates_Delete(struct InstanceStates instanceStates) {
  delete instanceStates.instanceStateVec;
}

InstanceManager InstanceManager_New(const char *config) {
  scouter::InstanceManager::Config ic = load_json<scouter::InstanceManager::Config>(config);
  return new scouter::InstanceManager(ic);
}

void InstanceManager_Delete(InstanceManager instanceManager) {
  delete instanceManager;
}

void InstanceManager_Update(InstanceManager instanceManager, TrackingResult tr) {
  instanceManager->update(*tr);
}

struct InstanceStates InstanceManager_GetCurrentStates(InstanceManager instanceManager) {
  std::vector<scouter::InstanceState> currentStates = instanceManager->get_current_states();
  std::vector<scouter::InstanceState>* states = new std::vector<scouter::InstanceState>(currentStates);
  InstanceStates ret = {states->size(), states};
  return ret;
}