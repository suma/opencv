#include "moving_matcher_bridge.h"
#include "util.hpp"
#include <scouter-core/mvom.hpp>

struct ByteArray MVCandidate_Serialize(MVCandidate c) {
  return serializeObject(*c);
}

MVCandidate MVCandidate_Deserialize(struct ByteArray src) {
  return deserializeObject<scouter::MVObjectCandidate>(src);
}

void MVCandidate_Delete(MVCandidate c) {
  delete c;
}

struct MVCandidates InvertMVCandidates(MVCandidate* obj, int length) {
  scouter::MVObjectCandidate** os = new scouter::MVObjectCandidate*[length];
  for (int i = 0; i < length; ++i) {
    os[i] = obj[i];
  }
  MVCandidates cs = {os, length};
  return cs;
}

void MVCandidates_Delete(struct MVCandidates mvCandidates) {
  delete mvCandidates.mvCandidates;
}

struct MVCandidates MVOM_GetMatching(RegionsWithCameraID* regions, int length,
    float kThreshold) {
  std::vector<std::vector<scouter::ObjectCandidate> > candidatez;
  for (int i = 0; i < length; ++i) {
    std::vector<scouter::ObjectCandidate> candidates;
    const RegionsWithCameraID& r = regions[i];
    for (int j = 0; j < r.candidates.length; ++j) {
      scouter::ObjectCandidate& o = *(r.candidates.candidates[j]);
      o.camera_id = r.cameraID;
      candidates.push_back(o);
    }
    candidatez.push_back(candidates);
  }
  const std::vector<scouter::MVObjectCandidate>& views =
    scouter::mvom::get_matching(candidatez, kThreshold);
  scouter::MVObjectCandidate** ret = new scouter::MVObjectCandidate*[views.size()];
  for (size_t i = 0; i < views.size(); ++i) {
    ret[i] = new scouter::MVObjectCandidate(views[i]);
  }
  MVCandidates mc = {ret, (int)views.size()};
  return mc;
}