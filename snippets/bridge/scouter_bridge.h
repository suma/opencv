#ifndef _SCOUTER_CORE_BRIDGE_H_
#define _SCOUTER_CORE_BRIDGE_H_

#include "opencv_bridge.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef void* Frame;
typedef void* DetectionResult;
typedef void* FrameProcessorConfig;
typedef void* FrameProcessor;
typedef void* DetectorConfig;
typedef void* Detector;
typedef void* RecognizeConfig;
typedef void* ImageTaggerCaffes;
typedef void* Integrator;
typedef void* IntegratorConfig;
typedef void* TrackingResult;

struct ByteArray Frame_Serialize(Frame f);
Frame Freme_Deserialize(struct ByteArray src);
void Frame_Delete(Frame f);

struct ByteArray DetectionResult_Serialize(DetectionResult dr);
DetectionResult DetectionResult_Deserialize(struct ByteArray src);
void DetectionResult_Delete(DetectionResult dr);

FrameProcessor FrameProcessor_New(FrameProcessorConfig config);
void FrameProcessor_Delete(FrameProcessor fp);
Frame FrameProcessor_Apply(FrameProcessor fp, MatVec3b buf,
                           long long timestamp, int cameraID);

Detector Detector_New(DetectorConfig config);
void Detector_Delete(Detector detector);
DetectionResult Detector_Detect(Detector detector, Frame frame);
MatVec3b DetectDrawResult(Frame frame, DetectionResult dr, long long ms);
unsigned long long Scouter_GetEpochms();

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