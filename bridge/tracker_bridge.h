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

Tracker Tracker_New(const char *config);
void Tracker_Delete(Tracker tracker);

void Tracker_Push(Tracker tracker, struct ScouterFrame2* frames, int fLength,
  MVCandidate* mvCandidates, int mvLength);
int Tracker_Ready(Tracker tracker);

#ifdef __cplusplus
}
#endif

#endif // _TRACKER_BRIDGE_H_