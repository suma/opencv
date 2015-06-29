#ifndef _FRAME_PROCESSOR_BRIDGE_H_
#define _FRAME_PROCESSOR_BRIDGE_H_

#include "opencv_bridge.h"

#ifdef __cplusplus
#include <scouter-core/frame_processor.hpp>
extern "C" {
#endif

#ifdef __cplusplus
typedef struct ScouterFrame {
  cv::Mat_<cv::Vec3b>* image;
  int offset_x;
  int offset_y;
} ScouterFrame;
typedef scouter::FrameProcessor* FrameProcessor;
#else
typedef struct ScouterFrame {
  MatVec3b image;
  int offset_x;
  int offset_y;
} ScouterFrame;
typedef void* FrameProcessor;
#endif

FrameProcessor FrameProcessor_New(const char *config);
void FrameProcessor_Delete(FrameProcessor fp);
struct ScouterFrame FrameProcessor_Projection(FrameProcessor pf, MatVec3b buf);

#ifdef __cplusplus
}
#endif

#endif //_FRAME_PROCESSOR_BRIDGE_H_
