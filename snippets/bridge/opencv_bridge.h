#ifndef _OPENCV_BRIDGE_H_
#define _OPENCV_BRIDGE_H_

#include "util.h"

#ifdef __cplusplus
#include <opencv2/opencv.hpp>
extern "C" {
#endif

#ifdef __cplusplus
typedef cv::Mat_<cv::Vec3b>* MatVec3b;
typedef cv::VideoCapture* VideoCapture;
#else
typedef void* MatVec3b;
typedef void* VideoCapture;
#endif
// MatVec3b is golang type wrapper for cv::Mat_<cv::Vec3b>
MatVec3b MatVec3b_New();
struct ByteArray MatVec3b_ToJpegData(MatVec3b m, int quality);
void MatVec3b_Delete(MatVec3b m);
void MatVec3b_CopyTo(MatVec3b src, MatVec3b dst);
int MatVec3b_Empty(MatVec3b m);

// VideoCapture is gloang type wrapper for cv::VideoCapture
VideoCapture VideoCapture_New();
void VideoCapture_Delete(VideoCapture v);
int VideoCapture_Open(VideoCapture v, const char* uri);
int VideoCapture_OpenDevice(VideoCapture v, int device);
void VideoCapture_Set(VideoCapture v, int prop, int param);
int VideoCapture_IsOpened(VideoCapture v);
int VideoCapture_Read(VideoCapture v, MatVec3b buf);
void VideoCapture_Grab(VideoCapture v, int skip);

#ifdef __cplusplus
}
#endif

#endif //_OPENCV_BRIDGE_H_
