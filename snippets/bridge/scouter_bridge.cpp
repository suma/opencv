#include "scouter_bridge.h"

#include <stdlib.h>
#include <msgpack.hpp>
#include <opencv2/opencv.hpp>
//#include <scouter-core/frame.hpp>
//#include <scouter-core/frame_processor.hpp>
//#include <scouter-core/detection_result.hpp>
//#include <scouter-core/detector.hpp>
//#include <scouter-core/epochms.hpp>

void FrameProcessor_SetUp(FrameProcessor fp, FrameProcessorConfig config) {
  /*
  scouter::FrameProcessor::Config *c = (scouter::FrameProcessor::Config*) config;
  scouter::FrameProcessor tempFp(*c);
  free(fp);
  fp = &tempFp;
  */
}

int FrameProcessor_Apply(FrameProcessor frameProcessor, MatVec3b buf,
                          long long timestamp, int cameraID, char* frame) {
  /*
  scouter::FrameProcessor *fp = (scouter::FrameProcessor*) frameProcessor;
  cv::Mat_<cv::Vec3b> *mat = (cv::Mat_<cv::Vec3b>*) buf;

  scouter::FrameMeta meta(timestamp, cameraID);
  scouter::Frame f = fp->apply(*mat, meta);
  */

  msgpack::sbuffer buffer;
  msgpack::packer<msgpack::sbuffer> pk(&buffer);
  pk.pack(999); // TOBE replace f

  free(frame);
  int size = buffer.size() * sizeof(char);
  frame = (char *)malloc(size);
  if (frame == NULL) {
    return 0;
  }
  frame = buffer.data();
  return 1;
}

void Detector_SetUp(Detector detector, DetectorConfig config) {
  /*
  Detector::Config *c = (Detector::Config*) config;
  Detector tempDetector(*c)
  free(detector)
  detector = &tempDetector;
  */
}

int Detector_Detect(Detector detector, char* frame, char* dr) {
  /*
  msgpack::unpacked frameMsg;
  msgpack::unpack(&frameMsg, frame, sizeof(frame)); // sizeof is ok?
  msgpack::object frameObj = frameMsg.get();
  scouter::Frame frame;
  frameObj.convert(&frame);

  scouter::Detector *d = (scouter::Detector*) detector;
  scouter::DetectionResult result = d->detect(frame);

  msgpack::sbuffer buffer;
  msgpack::packer<msgpack::sbuffer> pk(&buffer);
  pk.pack(result);

  free(dr);
  int size = buffer.size * sizeof(char);
  dr = (char *)malloc(size);
  if (dr == NULL) {
    return 0;
  }

  dr = buffer.data();
  */
  return 1;
}

unsigned long long Scouter_GetEpochms() {
  return 0; //scouter::get_epochms();
}

int DetectDrawResult(char* frame, char* dr, unsigned long long ms, char* resultFrame) {
  /*
  msgpack::unpacked frameMsg;
  msgpack::unpack(&frameMsg, frame, sizeof(frame)); // sizeof is ok?
  msgpack::object frameObj = frameMsg.get();
  scouter::Frame scouterFrame;
  frameObj.convert(&scouterFrame);

  msgpack::unpacked drMsg;
  msgpack::unpack(&drMsg, dr, sizeof(dr)); // sizeof is ok?
  msgpack::object drObj = drMsg.get();
  scouter::DetectionResult detectionResult;
  drObj.convert(&detectionResult);

  cv::Mat_<cv::Vec3b> c = draw_result(scouterFrame, detectonResult, ms);

  msgpack::sbuffer buffer;
  msgpack::packer<msgpack::sbuffer> pk(&buffer);
  pk.pack(c);

  free(resultFrame);
  int size = c.size * sizeof(char);
  c = (char *)malloc(size);
  if (c == NULL) {
    return 0;
  }

  resultFrame = buffer.data();
  */
  return 1;
}

/*
cv::Mat_<cv::Vec3b> draw_result(
    const scouter::Frame& frame,
    const scouter::DetectionResult& dr,
    uint64_t ms) const {
  cv::Mat_<cv::Vec3b> c = frame.image.clone();
  for (size_t i = 0; i < dr.object_candidates.size(); ++i) {
    const scouter::ObjectCandidate& o = dr.object_candidates[i];
    o.draw(c, cv::Scalar(0, 0, 255), 2);
  }
  std::stringstream ss;
  ss << ms << "ms";
  cv::putText(c, ss.str(), cv::Point(48, 48),
              cv::FONT_HERSHEY_SIMPLEX, 1.5,
              cv::Scalar(255, 0, 0));
  return c;
}
*/

