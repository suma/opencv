#include "opencv_bridge.h"

#include <string.h>
#include <msgpack.hpp>

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

struct RawData MatVec3b_ToRawData(MatVec3b m) {
  int width = m->cols;
  int height = m->rows;
  int size = width * height * 3;
  char* data = reinterpret_cast<char*>(m->data);
  ByteArray byteData = toByteArray(data, size);
  RawData raw = {width, height, byteData};
  return raw;
}

MatVec3b RawData_ToMatVec3b(struct RawData r) {
  int rows = r.height;
  int cols = r.width;
  cv::Mat_<cv::Vec3b>* mat = new cv::Mat_<cv::Vec3b>(rows, cols);
  memcpy(mat->data, r.data.data, r.data.length);
  return mat;
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

void VideoCapture_Release(VideoCapture v) {
  v->release();
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

VideoWriter VideoWriter_New() {
  return new cv::VideoWriter();
}

void VideoWriter_Delete(VideoWriter vw) {
  delete vw;
}

void VideoWriter_Open(VideoWriter vw, const char* name, double fps, int width,
    int height) {
  vw->open(name, CV_FOURCC('M', 'J', 'P', 'G'), fps, cv::Size(width, height), true);
}

void VideoWriter_OpenWithMat(VideoWriter vw, const char*name, double fps,
    MatVec3b img) {
  vw->open(name, CV_FOURCC('M', 'J', 'P', 'G'), fps, img->size(), true);
}

int VideoWriter_IsOpened(VideoWriter vw) {
  return vw->isOpened();
}

void VideoWriter_Write(VideoWriter vw, MatVec3b img) {
  *vw << *img;
}
