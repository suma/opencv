#include "scouter_bridge.h"
#include "util.hpp"

#include <string>
#include <opencv2/opencv.hpp>
#include <scouter-core/frame.hpp>
#include <scouter-core/frame_processor.hpp>
#include <scouter-core/detection_result.hpp>
#include <scouter-core/detector.hpp>
#include <scouter-core/epochms.hpp>

struct ByteArray Frame_Serialize(Frame f) {
  return serializeObject(static_cast<scouter::Frame*>(f));
}

Frame Freme_Deserialize(struct ByteArray src) {
  return deserializeObject<scouter::Frame>(src);
}

void Frame_Delete(Frame f) {
  delete static_cast<scouter::Frame*>(f);
}

struct ByteArray DetectionResult_Serialize(DetectionResult dr) {
  return serializeObject(static_cast<scouter::DetectionResult*>(dr));
}

DetectionResult DetectionResult_Deserialize(struct ByteArray src) {
  return deserializeObject<scouter::DetectionResult>(src);
}

void DetectionResult_Delete(DetectionResult dr) {
  delete static_cast<scouter::DetectionResult*>(dr);
}

FrameProcessor FrameProcessor_New(FrameProcessorConfig config) {
  return new scouter::FrameProcessor(
    *static_cast<scouter::FrameProcessor::Config*>(config));
}

void FrameProcessor_Delete(FrameProcessor fp) {
  delete static_cast<scouter::FrameProcessor*>(fp);
}

Frame FrameProcessor_Apply(FrameProcessor fp, MatVec3b buf,
                           long long timestamp, int cameraID) {
  scouter::FrameProcessor* processor = static_cast<scouter::FrameProcessor*>(fp);
  cv::Mat_<cv::Vec3b>* mat = static_cast<cv::Mat_<cv::Vec3b>*>(buf);

  scouter::FrameMeta meta(timestamp, cameraID);
  return new scouter::Frame(processor->apply(*mat, meta));
}

Detector Detector_New(DetectorConfig config) {
  return new scouter::Detector(
    *static_cast<scouter::Detector::Config*>(config));
}

void Detector_Delete(Detector detector) {
  delete static_cast<scouter::Detector*>(detector);
}

DetectionResult Detector_Detect(Detector detector, Frame frame) {
  scouter::Frame& fr = *static_cast<scouter::Frame*>(frame);
  scouter::Detector& d = *static_cast<scouter::Detector*>(detector);
  return new scouter::DetectionResult(d.detect(fr));
}

void draw_result(
    const scouter::Frame& frame,
    const scouter::DetectionResult& dr,
    uint64_t ms,
    cv::Mat_<cv::Vec3b>& target) {
  frame.image.copyTo(target);
  for (size_t i = 0; i < dr.object_candidates.size(); ++i) {
    const scouter::ObjectCandidate& o = dr.object_candidates[i];
    o.draw(target, cv::Scalar(0, 0, 255), 2);
  }
  std::stringstream ss;
  ss << ms << "ms";
  cv::putText(target, ss.str(), cv::Point(48, 48),
              cv::FONT_HERSHEY_SIMPLEX, 1.5,
              cv::Scalar(255, 0, 0));
}

MatVec3b DetectDrawResult(Frame frame, DetectionResult dr, long long ms) {
  scouter::Frame& fr = *static_cast<scouter::Frame*>(frame);
  scouter::DetectionResult& detected = *static_cast<scouter::DetectionResult*>(dr);
  cv::Mat_<cv::Vec3b>* target = new cv::Mat_<cv::Vec3b>();
  draw_result(fr, detected, ms, *target);
  return target;
}

unsigned long long Scouter_GetEpochms() {
  return scouter::get_epochms();
}

void ConvertToFramePointer(char* frByte, Frame frame) {
  // msgpack::unpacked frMsg;
  // msgpack::unpack(&frMsg, frByte, sizeof(frByte));
  // msgpack::object frObj = frMsg.get();
  // scouter::Frame fr;
  // frObj.convert(&fr);
  // frame = *fr;
}

void ImageTaggerCaffe_SetUp(ImageTaggerCaffes taggers, RecognizeConfig config) {

}
void ImageTaggerCaffe_PredictTagsBatch(ImageTaggerCaffes taggers, Frame frame, DetectionResult dr,
                                       DetectionResult resultDr, char** retByte, int* retLength) {

}
void RecognizeDrawResult(Frame frame, DetectionResult dr,
                         char** drwByte, int* drwLength) {
}

void ConvertToDetectionResultPointer(char* drByte, DetectionResult dr) {

}

void IntegratorSetUp(Integrator integrator, IntegratorConfig config) {

}
void Integrator_Push(Integrator integrator, Frame frame, DetectionResult dr) {

}
int Integrator_TrackerReady(Integrator integrator) {
  return 1;
}
void  Integrator_Track(Integrator integrator, TrackingResult tr, char** trByte, int* trLength) {

}

