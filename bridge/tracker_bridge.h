#ifndef _TRACKER_BRIDGE_H_
#define _TRACKER_BRIDGE_H_

#include "opencv_bridge.h"
#include "frame_processor_bridge.h"
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
typedef struct Trackee {
  unsigned long long colorID;
  MVCandidate mvCandidate;
  int interpolated;
} Trackee;
typedef struct TrackingResult {
  Trackee* trackees;
  int length;
  unsigned long long timestamp;
} TrackingResult;

Tracker Tracker_New(const char *config);
void Tracker_Delete(Tracker tracker);

void TrackingResult_Delete(TrackingResult trackingResult);

void Tracker_Push(Tracker tracker, struct ScouterFrame2* frames, int fLength,
  MVCandidate* mvCandidates, int mvLength);
struct TrackingResult Tracker_Track(Tracker tracker);
int Tracker_Ready(Tracker tracker);

#ifdef __cplusplus
}
#endif

#endif // _TRACKER_BRIDGE_H_