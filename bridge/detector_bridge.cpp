#include "detector_bridge.h"
#include "util.hpp"

struct ByteArray Candidate_Serialize(Candidate c) {
  return serializeObject(*c);
}

Candidate Candidate_Deserialize(struct ByteArray src) {
  return deserializeObject<scouter::ObjectCandidate>(src);
}

void Candidate_Delete(Candidate c) {
  delete c;
}

Detector Detector_New(const char *config) {
  scouter::Detector::Config dc = load_json<scouter::Detector::Config>(config);
  return new scouter::Detector(dc);
}

void Detector_Delete(Detector detector) {
  delete detector;
}

void ResolveCandidates(struct Candidates candidates, Candidate* obj) {
  for (size_t i = 0; i < candidates.candidateVec->size(); ++i) {
    obj[i] = new scouter::ObjectCandidate((*candidates.candidateVec)[i]);
  }
  return;
}

void Candidates_Delete(struct Candidates candidates) {
  delete candidates.candidateVec;
}

struct Candidates Detector_ACFDetect(Detector detector, MatVec3b image, int offsetX, int offsetY) {
  std::vector<scouter::ObjectCandidate> candidates = detector->acf_detect(*image, offsetX, offsetY);
  std::vector<scouter::ObjectCandidate>* ret = new std::vector<scouter::ObjectCandidate>();
  for (size_t i = 0; i < candidates.size(); ++i) {
    ret->push_back(candidates[i]);
  }
  Candidates c = {ret, (int)candidates.size()};
  return c;
}

struct Candidates Detector_FilterCandidateByMask(Detector detector, Candidate* candidates, int length) {
  std::vector<scouter::ObjectCandidate> candidateVec;
  for (int i = 0; i < length; ++i) {
    candidateVec.push_back(*candidates[i]);
  }
  std::vector<scouter::ObjectCandidate> filtered = detector->filter_candidate_by_mask(candidateVec);
  std::vector<scouter::ObjectCandidate>* ret = new std::vector<scouter::ObjectCandidate>();
  for (size_t i = 0; i < filtered.size(); ++i) {
    ret->push_back(filtered[i]);
  }
  Candidates c = {ret, (int)filtered.size()};
  return c;
}

struct Candidates Detector_EstimateCandidateHeight(Detector detector, Candidate* candidates, int length,
    int offsetX, int offsetY) {
  std::vector<scouter::ObjectCandidate>* candidateVec = new std::vector<scouter::ObjectCandidate>();
  for (int i = 0; i < length; ++i) {
    candidateVec->push_back(*candidates[i]);
  }
  detector->estimate_candidate_height(*candidateVec, offsetX, offsetY);
  Candidates c = {candidateVec, length};
  return c;
}
