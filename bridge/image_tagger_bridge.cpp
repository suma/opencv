#include "image_tagger_bridge.h"
#include "util.hpp"

ImageTaggerCaffe ImageTaggerCaffe_New(const char* config) {
  scouter::ImageTaggerCaffe::Config conf =
    load_json<scouter::ImageTaggerCaffe::Config>(config);
  return new scouter::ImageTaggerCaffe(conf);
}

void ImageTaggerCaffe_Delete(ImageTaggerCaffe tagger) {
  delete tagger;
}

Candidate ImageTaggerCaffe_CropAndPredictTags(ImageTaggerCaffe tagger, Candidate candidate,
    MatVec3b image) {
  std::vector<scouter::ObjectCandidate> candidateVec;
  candidateVec.push_back(*candidate);

  tagger->predict_tags_batch(candidateVec, *image);
  scouter::ObjectCandidate* ret = new scouter::ObjectCandidate(candidateVec[0]);
  return ret;
}

struct Candidates ImageTaggerCaffe_CropAndPredictTagsBatch(ImageTaggerCaffe tagger,
    Candidate* candidates, int length, MatVec3b image) {
  std::vector<scouter::ObjectCandidate> candidateVec;
  for (int i = 0; i < length; ++i) {
    candidateVec.push_back(*candidates[i]);
  }
  tagger->predict_tags_batch(candidateVec, *image);
  int l = (int)candidateVec.size();
  scouter::ObjectCandidate** ret = new scouter::ObjectCandidate*[l];
  for (size_t i = 0; i < l; ++i) {
    ret[i] = new scouter::ObjectCandidate(candidateVec[i]);
  }
  Candidates cs = {ret, l};
  return cs;
}
