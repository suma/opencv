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
  for (int i = 0; i < candidates.length; ++i) {
    obj[i] = candidates.candidates[i];
  }
}

struct Candidates InvertCandidates(Candidate* obj, int length) {
  scouter::ObjectCandidate** os = new scouter::ObjectCandidate*[length];
  for (int i = 0; i < length; ++i) {
    os[i] = new scouter::ObjectCandidate(*obj[i]);
  }
  Candidates cs = {os, length};
  return cs;
}

void Candidates_Delete(struct Candidates candidates) {
  delete candidates.candidates;
}

Detector Detector_New(const char *config) {
  scouter::Detector::Config dc = load_json<scouter::Detector::Config>(config);
  return new scouter::Detector(dc);
}

void Detector_Delete(Detector detector) {
  delete detector;
}

struct Candidates Detector_ACFDetect(Detector detector, MatVec3b image, int offsetX, int offsetY) {
  const std::vector<scouter::ObjectCandidate>& candidates = detector->acf_detect(*image, offsetX, offsetY);
  scouter::ObjectCandidate** ret = new scouter::ObjectCandidate*[candidates.size()];
  for (size_t i = 0; i < candidates.size(); ++i) {
    ret[i] = new scouter::ObjectCandidate(candidates[i]);
  }
  Candidates cs = {ret, (int)candidates.size()};
  return cs;
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
  const std::vector<scouter::ObjectCandidate>& candidates = detector->mm_detect(*image, offsetX, offsetY);
  scouter::ObjectCandidate** ret = new scouter::ObjectCandidate*[candidates.size()];
  for (size_t i = 0; i < candidates.size(); ++i) {
    ret[i] = new scouter::ObjectCandidate(candidates[i]);
  }
  Candidates cs = {ret, (int)candidates.size()};
  return cs;
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
