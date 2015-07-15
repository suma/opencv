#ifndef _TRACKER_BRIDGE_H_
#define _TRACKER_BRIDGE_H_

#include "opencv_bridge.h"
#include "moving_matcher_bridge.h"

#ifdef __cplusplus
#include <scouter-core/tracker_sp.hpp>
extern "C" {
#endif

#ifdef __cplusplus
typedef scouter::TrackerSP* Tracker;
#else
typedef void* Tracker;
#endif
typedef struct MatWithCameraID {
  int cameraID;
  MatVec3b mat;
} MatWithCameraID;

Tracker Tracker_New(const char *config);
void Tracker_Delete(Tracker tracker);

void Tracker_Push(Tracker tracker, struct MatWithCameraID* frames, int length,
  struct MVCandidates mvCandidates, unsigned long long timestamp);
void Tracker_track(Tracker tracker, unsigned long long timestamp);
int ready(Tracker tracker);

#ifdef __cplusplus
}
#endif

#endif // _TRACKER_BRIDGE_H_