#include "image_tagger_bridge.h"
#include "util.hpp"

ImageTaggerCaffe ImageTaggerCaffe_New(const char* config) {
  scouter::ImageTaggerCaffe::Config conf = load_json<scouter::ImageTaggerCaffe::Config>(config);
  return new scouter::ImageTaggerCaffe(conf);
}

void ImageTaggerCaffe_Delete(ImageTaggerCaffe tagger) {
  delete tagger;
}

Candidates ImageTaggerCaffe_PredictTagsBatch(ImageTaggerCaffe tagger,
    Candidate* candidates, int length, MatVec3b image) {
  std::vector<scouter::ObjectCandidate>* candidateVec = new std::vector<scouter::ObjectCandidate>();
  for (int i = 0; i < length; ++i) {
    candidateVec->push_back(*candidates[i]);
  }
  tagger->predict_tags_batch(*candidateVec, *image);

  Candidates c = {candidateVec, (int)candidateVec->size()};
  return c;
}
