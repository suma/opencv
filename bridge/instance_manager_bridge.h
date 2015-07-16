#ifndef _INSTANCE_MANAGER_BRIDGE_H_
#define _INSTANCE_MANAGER_BRIDGE_H_

#include "tracker_bridge.h"
#include "util.h"

#ifdef __cplusplus
#include <scouter-core/instance_manager.hpp>
extern "C" {
#endif

#ifdef __cplusplus
typedef scouter::InstanceManager* InstanceManager;
typedef scouter::InstanceState* InstanceState;
typedef struct InstanceStates {
  int length;
  std::vector<scouter::InstanceState>* instanceStateVec;
} InstanceStates;
#else
typedef void* InstanceManager;
typedef void* InstanceState;
typedef struct InstanceStates {
  int length;
  void* instanceStateVec;
} InstanceStates;
#endif

struct ByteArray InstanceState_Serialize(InstanceState s);
InstanceState InstanceState_Deserialize(struct ByteArray src);
void InstanceState_Delete(InstanceState s);

void ResolveInstanceStates(struct InstanceStates instanceStates, InstanceState* obj);
void InstanceStates_Delete(struct InstanceStates instanceStates);

InstanceManager InstanceManager_New(const char *config);
void InstanceManager_Delete(InstanceManager instanceManager);

void InstanceManager_Update(InstanceManager instanceManager, TrackingResult tr);
struct InstanceStates InstanceManager_GetCurrentStates(InstanceManager instanceManager);

#ifdef __cplusplus
}
#endif

#endif // _INSTANCE_MANAGER_BRIDGE_H_