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

void TrackingResult_Delete(TrackingResult trackingResult) {
  delete trackingResult;
}

void Tracker_Push(Tracker tracker, struct MatWithCameraID* frames, int length,
  struct MVCandidates mvCandidates, unsigned long long timestamp) {
  scouter::MatMapPtr ret(new scouter::MatMap);
  for (int i = 0; i < length; ++i) {
    ret->insert(std::make_pair(frames[i].cameraID, *(frames[i].mat)));
  }

  std::vector<scouter::MVObjectCandidate> mvCans;
  for (int i = 0; i < mvCandidates.length; ++i) {
    mvCans.push_back(*(mvCandidates.mvCandidates[i]));
  }

  tracker->push(ret, mvCans, timestamp);
}

TrackingResult Tracker_Track(Tracker tracker, unsigned long long timestamp) {
  scouter::TrackingResult ret = tracker->track(timestamp);
  return new scouter::TrackingResult(ret);
}

int Tracker_Ready(Tracker tracker) {
  return tracker->ready() ? 1 : 0;
}
