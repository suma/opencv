#include "scouter_bridge.h"

#include <stdlib.h>
#include <string>
#include <msgpack.hpp>
#include <opencv2/opencv.hpp>
// #include <scouter-core/frame.hpp>
// #include <scouter-core/frame_processor.hpp>
// #include <scouter-core/detection_result.hpp>
// #include <scouter-core/detector.hpp>
// #include <scouter-core/epochms.hpp>

void FrameProcessor_SetUp(FrameProcessor fp, FrameProcessorConfig config) {
  // scouter::FrameProcessor::Config *c = (scouter::FrameProcessor::Config*) config;
  // scouter::FrameProcessor tempFp(*c);
  // fp = &tempFp;
}

void FrameProcessor_Apply(FrameProcessor frameProcessor, MatVec3b buf,
                          long long timestamp, int cameraID,
                          Frame frame, char** frByte, int* frLength) {
  // scouter::FrameProcessor *fp = (scouter::FrameProcessor*) frameProcessor;
  // cv::Mat_<cv::Vec3b> *mat = (cv::Mat_<cv::Vec3b>*) buf;

  // scouter::FrameMeta meta(timestamp, cameraID);
  // scouter::Frame f = fp->apply(*mat, meta);
  // *frame = new scouter::Frame(f.meta, f.image);

  // msgpack::sbuffer buffer;
  // msgpack::packer<msgpack::sbuffer> pk(&buffer);
  // pk.pack(f);

  // char *tmp = buffer.data();
  // *frByte = (char *)malloc(buffer.size());
  // memcopy(frByte, tmp, buffer.size());
  // *frLength = buffer.size();
}

void Detector_SetUp(Detector detector, DetectorConfig config) {
  // Detector::Config *c = (Detector::Config*) config;
  // Detector tempDetector(*c)
  // detector = &tempDetector;
}

void Detector_Detect(Detector detector, Frame frame,
                     DetectionResult dr, char** drByte, int* drLength) {
  // scouter::Frame* fr = (scouter::Frame*) frame;

  // scouter::Detector *d = (scouter::Detector*) detector;
  // scouter::DetectionResult detected = d->detect(*fr);
  // dr = *detected;

  // msgpack::sbuffer buffer;
  // msgpack::packer<msgpack::sbuffer> pk(&buffer);
  // pk.pack(detected);

  // char *tmp = buffer.data();
  // *drByte = *tmp;
  // *drLength = buffer.size();
}


unsigned long long Scouter_GetEpochms() {
  return 0; //scouter::get_epochms();
}

void DetectDrawResult(Frame frame, DetectionResult dr, unsigned long long ms,
                      char** drwByte, int* drwLength) {
  // scouter::Frame* fr = (scouter::Frame*) frame;
  // scouter::DetectionResult detected = (scouter::DetectionResult*) dr;

  // cv::Mat_<cv::Vec3b> c = draw_result(*fr, *detected, ms);
  // draw = *c

  // msgpack::sbuffer buffer;
  // msgpack::packer<msgpack::sbuffer> pk(&buffer);
  // pk.pack(c);

  // char *tmp = buffer.data();
  // *drwByte = *tmp;
  // *drwLength = buffer.size();
}

// cv::Mat_<cv::Vec3b> draw_result(
//     const scouter::Frame& frame,
//     const scouter::DetectionResult& dr,
//     uint64_t ms) const {
//   cv::Mat_<cv::Vec3b> c = frame.image.clone();
//   for (size_t i = 0; i < dr.object_candidates.size(); ++i) {
//     const scouter::ObjectCandidate& o = dr.object_candidates[i];
//     o.draw(c, cv::Scalar(0, 0, 255), 2);
//   }
//   std::stringstream ss;
//   ss << ms << "ms";
//   cv::putText(c, ss.str(), cv::Point(48, 48),
//               cv::FONT_HERSHEY_SIMPLEX, 1.5,
//               cv::Scalar(255, 0, 0));
//   return c;
// }

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

