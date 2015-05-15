#ifndef _OPENCV_BRIDGE_H_
#define _OPENCV_BRIDGE_H_

#ifdef __cplusplus
extern "C" {
#endif

// VideoCapture
typedef void* VideoCapture;
typedef void* MatVec3b;

int VideoCapture_Open(char* uri, VideoCapture vcap);
int VideoCapture_IsOpened(VideoCapture vcap);
int VideoCapture_Read(VideoCapture vcap, MatVec3b buf);
void VideoCapture_Grab(VideoCapture vcap);

void MatVec3b_Clone(MatVec3b buf, MatVec3b cloneBuf);
int MatVec3b_Empty(MatVec3b buf);

#ifdef __cplusplus
}
#endif

#endif //_OPENCV_BRIDGE_H_