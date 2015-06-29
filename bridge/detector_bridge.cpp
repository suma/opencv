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

void ResolveCondidates(struct Candidates candidates, Candidate* obj) {
  return;
}

struct Candidates Detector_ACFDetect(Detector detector, MatVec3b image, int offsetX, int offssetY) {
  Candidates c = {};
  return c;
}

struct Candidates Detector_FilterCndidateByMask(struct Candidates candidates) {
  return candidates;
}

void Detector_EstimateCandidateHeight(struct Candidates candidates, int offsetX, int offsetY) {
   return;
}