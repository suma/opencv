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
#else
typedef void* InstanceManager;
typedef void* InstanceState;
#endif
typedef struct InstanceStates {
  InstanceState* instanceStates;
  int length;
} InstanceStates;

struct ByteArray InstanceState_Serialize(InstanceState s);
InstanceState InstanceState_Deserialize(struct ByteArray src);
void InstanceState_Delete(InstanceState s);

struct InstanceStates InvertInstanceStates(InstanceState* obj, int length);
void InstanceStates_Delete(struct InstanceStates instanceStates);
String InstanceState_ToJSON(struct InstanceStates instanceStates, int floorID,
  unsigned long long timestamp);

InstanceManager InstanceManager_New(const char *config);
void InstanceManager_Delete(InstanceManager instanceManager);

void InstanceManager_Update(InstanceManager instanceManager,
  struct MatWithCameraID* frames, int fLength, struct Trackee* trackees,
  int tLength, unsigned long long timestamp);
struct InstanceStates InstanceManager_GetCurrentStates(
  InstanceManager instanceManager);

#ifdef __cplusplus
}
#endif

#endif // _INSTANCE_MANAGER_BRIDGE_H_