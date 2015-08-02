#ifndef _INSTANCES_VISUALIZER_BRIDGE_H_
#define _INSTANCES_VISUALIZER_BRIDGE_H_

#include "instance_manager_bridge.h"

#ifdef __cplusplus
#include <scouter-core/instances_visualizer.hpp>
extern "C" {
#endif

#ifdef __cplusplus
typedef scouter::InstancesVisualizer* InstancesVisualizer;
#else
typedef void* InstancesVisualizer;
#endif

InstancesVisualizer InstancesVisualizer_New(InstanceManager im,
  const char *config);
void InstancesVisualizer_Delete(InstancesVisualizer iv);

MatVec3b InstancesVisualizer_Draw(InstancesVisualizer iv,
    struct MatWithCameraID* frames, int fLength, struct InstanceStates states,
    struct Trackee* trackees, int tLength);

#ifdef __cplusplus
}
#endif

#endif // _INSTANCES_VISUALIZER_BRIDGE_H_