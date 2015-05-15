#ifndef _SCOUTER_CORE_BRIDGE_H_
#define _SCOUTER_CORE_BRIDGE_H_

#include "opencv_bridge.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef void* FrameProcessorConfig;
typedef void* FrameProcessor;
typedef void* DetectorConfig;
typedef void* Detector;
typedef void* DetectionResult;

void FrameProcessor_SetUp(FrameProcessor fp, FrameProcessorConfig config);
int FrameProcessor_Apply(FrameProcessor frameProcessor, MatVec3b buf,
                          long long timestamp, int cameraID, char* frame);

void Detector_SetUp(Detector detector, DetectorConfig config);
int Detector_Detect(Detector detector, char* frame, char* dr);
unsigned long long Scouter_GetEpochms();
int DetectDrawResult(char* frame, char* dr, unsigned long long ms, char* resultFrame);

#ifdef __cplusplus
}
#endif

#endif //_SCOUTER_CORE_BRIDGE_H_