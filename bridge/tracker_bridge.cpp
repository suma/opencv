#include "tracker_bridge.h"
#include "util.hpp"
#include <scouter-core/tracking_result.hpp>


Tracker Tracker_New(const char *config) {
  scouter::TrackerSP::Config tc = load_json<scouter::TrackerSP::Config>(config);
  return new scouter::TrackerSP(tc);
}

void Tracker_Delete(Tracker tracker) {
  delete tracker;
}

void Tracker_Push(Tracker tracker, struct MatWithCameraID* frames, int length,
  struct MVCandidates mvCandidates, unsigned long long timestamp) {

}

void Tracker_track(Tracker tracker, unsigned long long timestamp) {

}

int ready(Tracker tracker) {
  return tracker->ready() ? 1 : 0;
}
