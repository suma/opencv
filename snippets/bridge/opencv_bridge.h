#ifndef _OPENCV_BRIDGE_H_
#define _OPENCV_BRIDGE_H_

#ifdef __cplusplus
extern "C" {
#endif

// VideoCapture
typedef void* VideoCapture;
typedef void* MatVec3b;

VideoCapture VideoCapture_Open(char* uri);
int VideoCapture_IsOpened(VideoCapture vcap);
void VideoCapture_Read(VideoCapture vcap, MatVec3b buf);
void VideoCapture_Grab(VideoCapture vcap);

MatVec3b MatVec3b_Clone(MatVec3b buf);
int MatVec3b_Empty(MatVec3b buf);

#ifdef __cplusplus
}
#endif

#endif //_OPENCV_BRIDGE_H_