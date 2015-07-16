#include "image_tagger_bridge.h"
#include "util.hpp"

ImageTaggerCaffe ImageTaggerCaffe_New(const char* config) {
  scouter::ImageTaggerCaffe::Config conf = load_json<scouter::ImageTaggerCaffe::Config>(config);
  return new scouter::ImageTaggerCaffe(conf);
}

void ImageTaggerCaffe_Delete(ImageTaggerCaffe tagger) {
  delete tagger;
}

MatVec3b ImageTaggerCaffe_Crop(ImageTaggerCaffe tagger, Candidate candidate, MatVec3b image) {
  cv::Mat_<cv::Vec3b> cropped = tagger->crop(*candidate, *image);
  return new cv::Mat_<cv::Vec3b>(cropped);
}

Candidates ImageTaggerCaffe_PredictTagsBatch(ImageTaggerCaffe tagger,
    Candidate* candidates, MatVec3b* croppedImages, int length) {
  std::vector<scouter::ObjectCandidate>* candidateVec = new std::vector<scouter::ObjectCandidate>();
  std::vector<cv::Mat_<cv::Vec3b> > croppedVec;
  for (int i = 0; i < length; ++i) {
    candidateVec->push_back(*candidates[i]);
    croppedVec.push_back(*croppedImages[i]);
  }
  tagger->predict_tags_batch(*candidateVec, croppedVec);
  Candidates c = {candidateVec, (int)candidateVec->size()};
  return c;
}
