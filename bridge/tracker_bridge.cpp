#include "tracker_bridge.h"
#include "util.hpp"
#include <scouter-core/mvom.hpp>

struct ByteArray MVCandidate_Serialize(MVCandidate c) {
  return serializeObject(*c);
}

MVCandidate MVCandidate_Deserialize(struct ByteArray src) {
  return deserializeObject<scouter::MVObjectCandidate>(src);
}

void Candidate_Delete(MVCandidate c) {
  delete c;
}

void ResolveMVCandidates(struct MVCandidates mvCandidates, MVCandidate* obj) {
  for (size_t i = 0; i < mvCandidates.candidateVec->size(); ++i) {
    obj[i] = new scouter::MVObjectCandidate((*mvCandidates.candidateVec)[i]);
  }
  return;
}

void MVCandidates_Delete(struct MVCandidates mvCandidates) {
  delete mvCandidates.candidateVec;
}

struct MVCandidates MVOM_GetMatching(Frame* frames, int length, float kThreshold) {
  std::vector<std::vector<scouter::ObjectCandidate> > candidatez;
  for (int i = 0; i < length; ++i) {
    std::vector<scouter::ObjectCandidate> candidates;
    for (int j = 0; j < frames[i].candidates.length; ++j) {
      scouter::ObjectCandidate& o = (*(frames[i].candidates.candidateVec))[j];
      o.camera_id = frames[i].cameraID;
      candidates.push_back(o);
    }
    candidatez.push_back(candidates);
  }
  std::vector<scouter::MVObjectCandidate> views =
    scouter::mvom::get_matching(candidatez, kThreshold);
  std::vector<scouter::MVObjectCandidate>* ret =
    new std::vector<scouter::MVObjectCandidate>();
  for (size_t i = 0; i < views.size(); ++i) {
    ret->push_back(views[i]);
  }
  MVCandidates mc = {ret, (int)views.size()};
  return mc;
}