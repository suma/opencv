#ifndef _SCOUTER_CORE_BRIDGE_H_
#define _SCOUTER_CORE_BRIDGE_H_

#include "opencv_bridge.h"

#ifdef __cplusplus
#include <scouter-core/frame.hpp>
#include <scouter-core/frame_processor.hpp>
#include <scouter-core/detection_result.hpp>
#include <scouter-core/tracking_result.hpp>
#include <scouter-core/detector.hpp>
#include <scouter-core/mm_detector.hpp>
#include <scouter-core/image_tagger.hpp>
#include <scouter-core/image_tagger_caffe.hpp>
#include <scouter-core/integrator.hpp>
#include <scouter-core/instance_manager.hpp>
#include <scouter-core/instances_visualizer.hpp>
extern "C" {
#endif

typedef struct String {
  const char* str;
  int length;
} String;

#ifdef __cplusplus
typedef scouter::Frame* Frame;
typedef struct ScouterFrame {
  cv::Mat_<cv::Vec3b>* image;
  int offset_x;
  int offset_y;
} ScouterFrame;
typedef scouter::DetectionResult* DetectionResult;
typedef scouter::FrameProcessor* FrameProcessor;
typedef scouter::Detector* Detector;
typedef scouter::MultiModelDetector* MultiModelDetector;
typedef std::vector<scouter::ImageTaggerCaffe>* ImageTaggerCaffe;
typedef struct Taggers {
  std::map<std::string, cv::Mat_<cv::Vec3b> >* drawResultsMap;
  int length;
} Taggers;
typedef scouter::TrackingResult* TrackingResult;
typedef scouter::Integrator* Integrator;
typedef scouter::InstanceManager* InstanceManager;
typedef std::vector<scouter::InstanceState>* InstanceStates;
typedef scouter::InstancesVisualizer* Visualizer;
typedef struct PlotTrajectories {
  std::vector<cv::Mat_<cv::Vec3b> >* trajectories;
  int length;
} PlotTrajectories;
#else
typedef void* Frame;
typedef struct ScouterFrame {
  MatVec3b image;
  int offset_x;
  int offset_y;
} ScouterFrame;
typedef void* DetectionResult;
typedef void* FrameProcessor;
typedef void* Detector;
typedef void* MultiModelDetector;
typedef void* ImageTaggerCaffe;
typedef struct Taggers {
  void* drawResultsMap;
  int length;
} Taggers;
typedef void* TrackingResult;
typedef void* Integrator;
typedef void* InstanceManager;
typedef void* InstanceStates;
typedef void* Visualizer;
typedef struct PlotTrajectories {
  void* trajectories;
  int length;
} PlotTrajectories;
#endif

struct ByteArray Frame_Serialize(Frame f);
Frame Freme_Deserialize(struct ByteArray src);
void Frame_Delete(Frame f);

struct ByteArray DetectionResult_Serialize(DetectionResult dr);
DetectionResult DetectionResult_Deserialize(struct ByteArray src);
void DetectionResult_Delete(DetectionResult dr);

FrameProcessor FrameProcessor_New(const char *config);
void FrameProcessor_Delete(FrameProcessor fp);
struct ScouterFrame FrameProcessor_Projection(FrameProcessor pf, MatVec3b buf);

Detector Detector_New(const char *config);
void Detector_Delete(Detector detector);
DetectionResult Detector_Detect(Detector detector, Frame frame);

MultiModelDetector MultiModelDetector_New(const char *config);
void MultiModelDetector_Delete(MultiModelDetector detector);
DetectionResult MultiModelDetector_Detect(MultiModelDetector detector, Frame frame);

MatVec3b DetectDrawResult(Frame frame, DetectionResult dr, long long ms);

ImageTaggerCaffe ImageTaggerCaffe_New(const char *config);
void ImageTaggerCaffe_Delete(ImageTaggerCaffe taggers);
DetectionResult Recognize(ImageTaggerCaffe taggers, Frame frame, DetectionResult dr);
Taggers RecognizeDrawResult(Frame frame, DetectionResult dr);
void Taggers_Delete(Taggers taggers);
void ResolveDrawResult(struct Taggers taggers, const char** keys, MatVec3b* drawResults);

Integrator Integrator_New(const char *config);
void Integrator_Delete(Integrator integrator);
void Integrator_Push(Integrator integrator, Frame* frame, DetectionResult* dr, int size);
int Integrator_TrackerReady(Integrator integrator);
TrackingResult Integrator_Track(Integrator integrator);
void TrackingResult_Delete(TrackingResult tr);

InstanceManager InstanceManager_New(const char *config);
void InstanceManager_Delete(InstanceManager instanceManager);
InstanceStates InstanceManager_GetCurrentStates(InstanceManager instanceManager,
                                                TrackingResult result);
void InstanceStates_Delete(InstanceStates states);
String ConvertStatesToJson(InstanceStates instanceStates,
                                int floorID, long long timestamp);
Visualizer Visualizer_New(const char *config, InstanceManager instanceManager);
void Visualizer_Delete(Visualizer visualizer);
PlotTrajectories Visualizer_PlotTrajectories(Visualizer visualizer);
void PlotTrajectories_Delete(PlotTrajectories plotTrajectories);
void ResolvePlotTrajectories(struct PlotTrajectories plotTrajectories, MatVec3b* trajectories);

#ifdef __cplusplus
}
#endif

#endif //_SCOUTER_CORE_BRIDGE_H_
