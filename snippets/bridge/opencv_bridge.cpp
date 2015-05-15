#include "opencv_bridge.h"

#include <string.h>
#include <opencv2/opencv.hpp>

int VideoCapture_Open(char* uri, VideoCapture vcap) {
  cv::VideoCapture tempVcap; // TODO default constructor
  //if (!tempVcap.open(uri)) {
  //  return 0;
  //}
  free(vcap);
  vcap = &tempVcap;
  return 1;
}

int VideoCapture_IsOpened(VideoCapture vcap) {
  cv::VideoCapture* vc = (cv::VideoCapture*) vcap;
  return vc->isOpened();
}

int VideoCapture_Read(VideoCapture vcap, MatVec3b buf) {
  cv::VideoCapture *vc = (cv::VideoCapture*) vcap;
  cv::Mat_<cv::Vec3b> *result = (cv::Mat_<cv::Vec3b>*) buf;
  if (!vc->read(*result)) {
    return 0;
  }
  buf = &result;
  return 1;
}

void VideoCapture_Grab(VideoCapture vcap) {
  cv::VideoCapture* vc = (cv::VideoCapture*) vcap;
  vc->grab();
}

void MatVec3b_Clone(MatVec3b buf, MatVec3b cloneBuf) {
  cv::Mat_<cv::Vec3b> *mat = (cv::Mat_<cv::Vec3b>*) buf;
  cv::Mat_<cv::Vec3b> result;
  result = mat->clone();
  free(cloneBuf);
  cloneBuf = &result;
}

int MatVec3b_Empty(MatVec3b buf) {
  cv::Mat_<cv::Vec3b> *mat = (cv::Mat_<cv::Vec3b>*) buf;
  return mat->empty();
}
