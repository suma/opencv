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

Candidate ImageTaggerCaffe_PredictTags(ImageTaggerCaffe tagger, Candidate candidate,
    MatVec3b croppedImg) {
  std::vector<scouter::ObjectCandidate> candidateVec;
  std::vector<cv::Mat_<cv::Vec3b> > croppedVec;
  candidateVec.push_back(*candidate);
  croppedVec.push_back(*croppedImg);

  tagger->predict_tags_batch(candidateVec, croppedVec);
  scouter::ObjectCandidate* ret = new scouter::ObjectCandidate(candidateVec[0]);
  return ret;
}

struct Candidates ImageTaggerCaffe_PredictTagsBatch(ImageTaggerCaffe tagger,
    Candidate* candidates, MatVec3b* croppedImages, int length) {
  std::vector<scouter::ObjectCandidate> candidateVec;
  std::vector<cv::Mat_<cv::Vec3b> > croppedVec;
  for (int i = 0; i < length; ++i) {
    candidateVec.push_back(*candidates[i]);
    croppedVec.push_back(*croppedImages[i]);
  }
  tagger->predict_tags_batch(candidateVec, croppedVec);
  int l = (int)candidateVec.size();
  scouter::ObjectCandidate** ret = new scouter::ObjectCandidate*[l];
  for (size_t i = 0; i < l; ++i) {
    ret[i] = new scouter::ObjectCandidate(candidateVec[i]);
  }
  Candidates cs = {ret, l};
  return cs;
}

Candidate ImageTaggerCaffe_CropAndPredictTags(ImageTaggerCaffe tagger, Candidate candidate,
    MatVec3b image) {
  cv::Mat_<cv::Vec3b>* cropped = ImageTaggerCaffe_Crop(tagger, candidate, image);
  return ImageTaggerCaffe_PredictTags(tagger, candidate, cropped);
}
