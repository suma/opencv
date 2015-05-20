#ifndef _SCOUTER_CORE_BRIDGE_H_
#define _SCOUTER_CORE_BRIDGE_H_

#include "opencv_bridge.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef void* FrameProcessorConfig;
typedef void* FrameProcessor;
typedef void* Frame;
typedef void* DetectorConfig;
typedef void* Detector;
typedef void* DetectionResult;
typedef void* RecognizeConfig;
typedef void* ImageTaggerCaffes;
typedef void* Integrator;
typedef void* IntegratorConfig;
typedef void* TrackingResult;

void FrameProcessor_SetUp(FrameProcessor fp, FrameProcessorConfig config);
void FrameProcessor_Apply(FrameProcessor frameProcessor, MatVec3b buf,
                          long long timestamp, int cameraID,
                          Frame frame, char** frByte, int* frLength);

void Detector_SetUp(Detector detector, DetectorConfig config);
void Detector_Detect(Detector detector, Frame frame,
                    DetectionResult dr, char** drByte, int* drLength);
unsigned long long Scouter_GetEpochms();
void DetectDrawResult(Frame frame, DetectionResult dr, unsigned long long ms,
                      char** drwByte, int* drwLength);
void ConvertToFramePointer(char* frByte, Frame frame);

void ImageTaggerCaffe_SetUp(ImageTaggerCaffes taggers, RecognizeConfig config);
void ImageTaggerCaffe_PredictTagsBatch(ImageTaggerCaffes taggers, Frame frame, DetectionResult dr,
                                       DetectionResult resultDr, char** retByte, int* retLength);
void RecognizeDrawResult(Frame frame, DetectionResult dr,
                         char** drwByte, int* drwLength);
void ConvertToDetectionResultPointer(char* drByte, DetectionResult dr);

void IntegratorSetUp(Integrator integrator, IntegratorConfig config);
void Integrator_Push(Integrator integrator, Frame frame, DetectionResult dr);
int Integrator_TrackerReady(Integrator integrator);
void Integrator_Track(Integrator integrator, TrackingResult tr, char** trByte, int* trLength);

#ifdef __cplusplus
}
#endif

#endif //_SCOUTER_CORE_BRIDGE_H_