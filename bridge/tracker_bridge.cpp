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

void Tracker_Push(Tracker tracker, struct ScouterFrame2* frames, int fLength,
  MVCandidate* mvCandidates, int mvLength) {

  std::vector<scouter::Frame> frameVec;
  for (int i = 0; i < fLength; ++i) {
    ScouterFrame2 fs = frames[i];
    scouter::FrameMeta fm = scouter::FrameMeta(fs.timestamp, fs.offset_x, fs.offset_y);
    scouter::Frame f = scouter::Frame(fm, *(fs.image));

    frameVec.push_back(f);
  }

  std::vector<scouter::MVObjectCandidate> mvos;
  for (int i = 0; i < mvLength; ++i) {
    mvos.push_back(*(mvCandidates[i]));
  }

  tracker->push(scouter::make_frames(frameVec), mvos);
}

struct TrackingResult Tracker_Track(Tracker tracker) {
  struct TrackingResult tr = {};
  return tr;
}

int Tracker_Ready(Tracker tracker) {
  return tracker->ready() ? 1 : 0;
}
