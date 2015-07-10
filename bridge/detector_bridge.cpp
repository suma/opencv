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

void ResolveCandidates(struct Candidates candidates, Candidate* obj) {
  for (size_t i = 0; i < candidates.candidateVec->size(); ++i) {
    obj[i] = new scouter::ObjectCandidate((*candidates.candidateVec)[i]);
  }
  return;
}

void Candidates_Delete(struct Candidates candidates) {
  delete candidates.candidateVec;
}

Detector Detector_New(const char *config) {
  scouter::Detector::Config dc = load_json<scouter::Detector::Config>(config);
  return new scouter::Detector(dc);
}

void Detector_Delete(Detector detector) {
  delete detector;
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

int Detector_FilterByMask(Detector detector, Candidate candidate) {
  return detector->filter_by_mask(*candidate) ? 0 : 1;
}

void Detector_EstimateHeight(Detector detector, Candidate candidate, int offsetX, int offsetY) {
  detector->estimate_height(*candidate, offsetX, offsetY);
}

void Detector_PutFeature(Detector detector, Candidate candidate, MatVec3b image) {
  detector->put_feature(*candidate, *image);
}

MMDetector MMDetector_New(const char *config) {
  scouter::MultiModelDetector::Config dc = load_json<scouter::MultiModelDetector::Config>(config);
  return new scouter::MultiModelDetector(dc);
}

void MMDetector_Delete(MMDetector detector) {
  delete detector;
}

struct Candidates MMDetector_MMDetect(MMDetector detector, MatVec3b image, int offsetX, int offsetY) {
  std::vector<scouter::ObjectCandidate> candidates = detector->mm_detect(*image, offsetX, offsetY);
  std::vector<scouter::ObjectCandidate>* ret = new std::vector<scouter::ObjectCandidate>();
  for (size_t i = 0; i < candidates.size(); ++i) {
    ret->push_back(candidates[i]);
  }
  Candidates c = {ret, (int)candidates.size()};
  return c;
}

int MMDetector_FilterByMask(MMDetector detector, Candidate candidate) {
  return detector->filter_by_mask(*candidate) ? 0 : 1;
}

void MMDetector_EstimateHeight(MMDetector detector, Candidate candidate, int offsetX, int offsetY) {
  detector->estimate_height(*candidate, offsetX, offsetY);
}

MatVec3b Candidates_Draw(MatVec3b image, Candidate* candidates, int length) {
  cv::Mat_<cv::Vec3b>* c = new cv::Mat_<cv::Vec3b>();
  image->copyTo(*c);
  for (int i = 0; i < length; ++i) {
    const scouter::ObjectCandidate& o = *candidates[i];
    o.draw(*c, cv::Scalar(0, 0, 255), 2);
  }
  return c;
}
