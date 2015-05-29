#ifndef _SCOUTER_CORE_BRIDGE_H_
#define _SCOUTER_CORE_BRIDGE_H_

#include "opencv_bridge.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef void* Frame;
typedef void* DetectionResult;
typedef void* TrackingResult;
typedef void* FrameProcessor;
typedef void* Detector;
typedef void* ImageTaggerCaffe;
typedef void* Taggers;
typedef void* Integrator;
typedef void* InstanceManager;
typedef void* InstanceStates;

struct ByteArray Frame_Serialize(Frame f);
Frame Freme_Deserialize(struct ByteArray src);
void Frame_Delete(Frame f);

struct ByteArray DetectionResult_Serialize(DetectionResult dr);
DetectionResult DetectionResult_Deserialize(struct ByteArray src);
void DetectionResult_Delete(DetectionResult dr);

struct ByteArray TrackingResult_Serialize(TrackingResult tr);
TrackingResult TrackingResult_Deserialize(struct ByteArray src);
void TrackingResult_Delete(TrackingResult tr);

FrameProcessor FrameProcessor_New(const char *config);
void FrameProcessor_Delete(FrameProcessor fp);
Frame FrameProcessor_Apply(FrameProcessor fp, MatVec3b buf,
                           long long timestamp, int cameraID);

Detector Detector_New(const char *config);
void Detector_Delete(Detector detector);
DetectionResult Detector_Detect(Detector detector, Frame frame);
MatVec3b DetectDrawResult(Frame frame, DetectionResult dr, long long ms);

ImageTaggerCaffe ImageTaggerCaffe_New(const char *config);
void ImageTaggerCaffe_Delete(ImageTaggerCaffe taggers);
DetectionResult Recognize(ImageTaggerCaffe taggers, Frame frame, DetectionResult dr);
Taggers RecognizeDrawResult(Frame frame, DetectionResult dr);

Integrator Integrator_New(const char *config);
void Integrator_Delete(Integrator integrator);
void Integrator_Push(Integrator integrator, Frame frame, DetectionResult dr);
int Integrator_TrackerReady(Integrator integrator);
TrackingResult Integrator_Track(Integrator integrator);

InstanceManager InstanceManager_New(const char *config);
void InstanceManager_Delete(InstanceManager instanceManager);
InstanceStates InstanceManager_GetCurrentStates(InstanceManager instanceManager,
                                                TrackingResult result);
void InstanceStates_Delete(InstanceStates states);
const char* ConvertStatesToJson(InstanceStates instanceStates, int floorID);


#ifdef __cplusplus
}
#endif

#endif //_SCOUTER_CORE_BRIDGE_H_
