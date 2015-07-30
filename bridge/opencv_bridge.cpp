#include "opencv_bridge.h"
#include "util.hpp"

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

struct ByteArray MatVec3b_Serialize(MatVec3b m) {
  msgpack::sbuffer buf;
  msgpack::packer<msgpack::sbuffer> pk(&buf);
  pk.pack_array(3);
  pk.pack(m->rows);
  pk.pack(m->cols);
  int size = m->rows * m->cols * 3;
  pk.pack_raw(size);
  assert(m->isContinuous());
  pk.pack_raw_body(reinterpret_cast<char*>(m->data), size);
  return toByteArray(buf.data(), buf.size());
}

MatVec3b MatVec3b_Deserialize(struct ByteArray src) {
  msgpack::unpacked msg;
  msgpack::unpack(&msg, src.data, src.length);
  msgpack::object obj = msg.get();

  msgpack::object_array obj_array = obj.via.array;
  assert(obj_array.size == 3);
  int rows, cols;
  obj_array.ptr[0] >> rows;
  obj_array.ptr[1] >> cols;
  cv::Mat_<cv::Vec3b>* mat = new cv::Mat_<cv::Vec3b>(rows, cols);
  memcpy(mat->data, obj_array.ptr[2].via.raw.ptr, obj_array.ptr[2].via.raw.size);
  return mat;
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
  std::string path = "";
  path += name;
  path += ".avi";
  vw->open(path, CV_FOURCC('M', 'J', 'P', 'G'), fps, cv::Size(width, height), true);
}

int VideoWriter_IsOpened(VideoWriter vw) {
  return vw->isOpened();
}

void VideoWriter_Write(VideoWriter vw, MatVec3b img) {
  *vw << *img;
}
