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
  delete trackingResult.trackees;
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

struct TrackingResult Tracker_Track(Tracker tracker) {
  const scouter::TrackingResult& ret = tracker->track();

  int length = int(ret.trackees.size());
  struct Trackee* trackees = new Trackee[length];
  for (int i = 0; i < length; ++i) {
    const scouter::Trackee& t = ret.trackees[i];
    scouter::MVObjectCandidate* o = new scouter::MVObjectCandidate(t.object);
    int interpolated = t.interpolated ? 1 : 0;
    Trackee trackee = {t.id, o, interpolated};
    trackees[i] = trackee;
  }

  struct TrackingResult tr = {trackees, length, ret.timestamp};
  return tr;
}

int Tracker_Ready(Tracker tracker) {
  return tracker->ready() ? 1 : 0;
}
