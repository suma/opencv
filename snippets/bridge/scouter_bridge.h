#ifndef _SCOUTER_CORE_BRIDGE_H_
#define _SCOUTER_CORE_BRIDGE_H_

#include "opencv_bridge.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef void* FrameProcessorConfig;
typedef void* FrameProcessor;
typedef void* Frame; // = pointer, want to use []uint8_t

FrameProcessor FrameProcessor_SetUp(FrameProcessorConfig config);
void FrameProcessor_Apply(FrameProcessor frameProcessor, MatVec3b buf,
                             long long timestamp, int cameraID,
                             Frame frame, int frameLength);

#ifdef __cplusplus
}
#endif

#endif //_SCOUTER_CORE_BRIDGE_H_