#include "opencv_bridge.h"
#include "util.hpp"

#include <string.h>
#include <opencv2/opencv.hpp>

MatVec3b MatVec3b_New() {
  return new cv::Mat_<cv::Vec3b>();
}

struct ByteArray MatVec3b_ToJpegData(MatVec3b m, int quality){
  cv::Mat_<cv::Vec3b>& mat = *static_cast<cv::Mat_<cv::Vec3b>*>(m);
  std::vector<int> param(2);
  param[0] = CV_IMWRITE_JPEG_QUALITY;
  param[1] = quality;
  std::vector<uchar> data;
  cv::imencode(".jpg", mat, data, param);
  return toByteArray(reinterpret_cast<const char*>(&data[0]), data.size());
}

void MatVec3b_Delete(MatVec3b m) {
  delete static_cast<cv::Mat_<cv::Vec3b>*>(m);
}

void MatVec3b_CopyTo(MatVec3b src, MatVec3b dst) {
  static_cast<cv::Mat_<cv::Vec3b>*>(src)->copyTo(*static_cast<cv::Mat_<cv::Vec3b>*>(dst));
}

int MatVec3b_Empty(MatVec3b m) {
  return static_cast<cv::Mat_<cv::Vec3b>*>(m)->empty();
}

VideoCapture VideoCapture_New() {
  return new cv::VideoCapture();
}

void VideoCapture_Delete(VideoCapture v) {
  delete static_cast<cv::VideoCapture*>(v);
}

int VideoCapture_Open(VideoCapture v, const char* uri) {
  return static_cast<cv::VideoCapture*>(v)->open(uri);
}

int VideoCapture_OpenDevice(VideoCapture v, int device) {
  return static_cast<cv::VideoCapture*>(v)->open(device);
}

void VideoCapture_Set(VideoCapture v, int prop, int param) {
  static_cast<cv::VideoCapture*>(v)->set(prop, param);
}

int VideoCapture_IsOpened(VideoCapture v) {
  return static_cast<cv::VideoCapture*>(v)->isOpened();
}

int VideoCapture_Read(VideoCapture v, MatVec3b buf) {
  return static_cast<cv::VideoCapture*>(v)->read(*static_cast<cv::Mat_<cv::Vec3b>*>(buf));
}

void VideoCapture_Grab(VideoCapture v, int skip) {
  cv::VideoCapture vcap = *static_cast<cv::VideoCapture*>(v);
  for (int i =0; i < skip; i++) {
    vcap.grab();
  }
}
