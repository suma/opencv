#include "opencv_bridge.h"
#include "util.hpp"

#include <string.h>

MatVec3b MatVec3b_New() {
  return new cv::Mat_<cv::Vec3b>();
}

struct ByteArray MatVec3b_ToJpegData(MatVec3b m, int quality){
  std::vector<int> param(2);
  param[0] = CV_IMWRITE_JPEG_QUALITY;
  param[1] = quality;
  std::vector<uchar> data;
  cv::imencode(".jpg", *m, data, param);
  return toByteArray(reinterpret_cast<const char*>(&data[0]), data.size());
}

void MatVec3b_Delete(MatVec3b m) {
  delete m;
}

void MatVec3b_CopyTo(MatVec3b src, MatVec3b dst) {
  src->copyTo(*dst);
}

int MatVec3b_Empty(MatVec3b m) {
  return m->empty();
}

VideoCapture VideoCapture_New() {
  return new cv::VideoCapture();
}

void VideoCapture_Delete(VideoCapture v) {
  delete v;
}

int VideoCapture_Open(VideoCapture v, const char* uri) {
  return v->open(uri);
}

int VideoCapture_OpenDevice(VideoCapture v, int device) {
  return v->open(device);
}

void VideoCapture_Set(VideoCapture v, int prop, int param) {
  v->set(prop, param);
}

int VideoCapture_IsOpened(VideoCapture v) {
  return v->isOpened();
}

int VideoCapture_Read(VideoCapture v, MatVec3b buf) {
  return v->read(*buf);
}

void VideoCapture_Grab(VideoCapture v, int skip) {
  for (int i =0; i < skip; i++) {
    v->grab();
  }
}
