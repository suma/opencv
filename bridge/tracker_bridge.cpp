#include "tracker_bridge.h"
#include <scouter-core/mvom.hpp>

void ResolveCandidatez(struct Candidatez candidatez, Candidates* obj) {
  for (size_t i = 0; i < candidatez.candidatesVec->size(); ++i) {
    Candidates c = {
      new std::vector<scouter::ObjectCandidate>((*candidatez.candidatesVec)[i]),
      (*candidatez.candidatesVec)[i].size()};
    obj[i] = c;
  }
  return;
}

void Candidatez_Delete(struct Candidatez candidatesVec) {
  delete candidatesVec.candidatesVec;
}

struct Candidatez MVOM_GetMatching(Candidate** candidatez, int* lengths, int length, float kThreshold) {
  std::vector<std::vector<scouter::ObjectCandidate> > frames;
  for (int i = 0; i < length; ++i) {
    std::vector<scouter::ObjectCandidate> candidates;
    for (int j = 0; j < lengths[i]; ++j) {
      candidates.push_back(*candidatez[i][j]);
    }
    frames.push_back(candidates);
  }
  std::vector<std::vector<scouter::ObjectCandidate> >& views =
    scouter::mvom::get_matching(frames, kThreshold);
  std::vector<std::vector<scouter::ObjectCandidate> >* ret =
    new std::vector<std::vector<scouter::ObjectCandidate> >();
  for (size_t i = 0; i < views.size(); ++i) {
    ret->push_back(views[i]);
  }
  Candidatez cv = {ret, (int)views.size()};
  return cv;
}