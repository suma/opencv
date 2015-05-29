#include "scouter_bridge.h"
#include "util.hpp"

#include <string>
#include <sstream>
#include <vector>
#include <map>
#include <set>
#include <pficommon/text/json.h>
#include <jsonconfig.hpp>
#include <opencv2/opencv.hpp>
#include <scouter-core/frame.hpp>
#include <scouter-core/frame_processor.hpp>
#include <scouter-core/detection_result.hpp>
#include <scouter-core/mv_detection_result.hpp>
#include <scouter-core/tracking_result.hpp>
#include <scouter-core/detector.hpp>
#include <scouter-core/epochms.hpp>
#include <scouter-core/image_tagger.hpp>
#include <scouter-core/image_tagger_caffe.hpp>
#include <scouter-core/integrator.hpp>
#include <scouter-core/instance_manager.hpp>

template <class Type>
Type load_json(const char *config) {
  std::stringstream ss(config);
  pfi::text::json::json config_raw;
  ss >> config_raw;
  return jsonconfig::config_cast<Type>(jsonconfig::config_root(config_raw));
}

struct ByteArray Frame_Serialize(Frame f) {
  return serializeObject(*static_cast<scouter::Frame*>(f));
}

Frame Freme_Deserialize(struct ByteArray src) {
  return deserializeObject<scouter::Frame>(src);
}

void Frame_Delete(Frame f) {
  delete static_cast<scouter::Frame*>(f);
}

struct ByteArray DetectionResult_Serialize(DetectionResult dr) {
  return serializeObject(*static_cast<scouter::DetectionResult*>(dr));
}

DetectionResult DetectionResult_Deserialize(struct ByteArray src) {
  return deserializeObject<scouter::DetectionResult>(src);
}

void DetectionResult_Delete(DetectionResult dr) {
  delete static_cast<scouter::DetectionResult*>(dr);
}

//TODO need to convert scouter::TrackingResult
struct ByteArray TrackingResult_Serialize(TrackingResult tr) {
  return serializeObject(*static_cast<scouter::TrackingResult*>(tr));
}

TrackingResult TrackingResult_Deserialize(struct ByteArray src) {
  return deserializeObject<scouter::DetectionResult>(src);
}

void TrackingResult_Delete(TrackingResult tr) {
  delete static_cast<scouter::DetectionResult*>(tr);
}

FrameProcessor FrameProcessor_New(const char *config) {
  scouter::FrameProcessor::Config fpc =
      load_json<scouter::FrameProcessor::Config>(config);
  return new scouter::FrameProcessor(fpc);
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

Detector Detector_New(const char *config) {
  scouter::Detector::Config dc = load_json<scouter::Detector::Config>(config);
  return new scouter::Detector(dc);
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

ImageTaggerCaffe ImageTaggerCaffe_New(const char *config) {
  std::vector<scouter::ImageTaggerCaffe::Config> taggers =
      load_json<std:: vector<scouter::ImageTaggerCaffe::Config> >(config);
  std::vector<scouter::ImageTaggerCaffe>* target = new std::vector<scouter::ImageTaggerCaffe>();
  for (size_t i = 0; i < taggers.size(); ++i) {
    target->push_back(scouter::ImageTaggerCaffe(taggers[i]));
  }
  return target;
}

void ImageTaggerCaffe_Delete(ImageTaggerCaffe taggers) {
  delete static_cast<std::vector<scouter::ImageTaggerCaffe>*>(taggers);
}

DetectionResult Recognize(ImageTaggerCaffe taggers, Frame frame, DetectionResult dr) {
  std::vector<scouter::ImageTaggerCaffe>& tags = *static_cast<
    std::vector<scouter::ImageTaggerCaffe>*>(taggers);
  scouter::Frame& fr = *static_cast<scouter::Frame*>(frame);
  scouter::DetectionResult& detected = *static_cast<scouter::DetectionResult*>(dr);
  for (size_t i = 0; i < tags.size(); ++i) {
    scouter::ImageTaggerCaffe tagger = tags[i];
    tagger.predict_tags_batch(detected.object_candidates, fr);
  }
  return &detected;
}

std::map<std::string, cv::Mat_<cv::Vec3b> >* draw_result(
    scouter::Frame& frame,
    scouter::DetectionResult& dr) {
  typedef std::map<std::string, cv::Scalar> ColorMap;
  ColorMap color_map;   // TODO(tabe): make it configurable
  color_map.insert(std::make_pair("Yes", cv::Scalar(0, 0, 255)));
  color_map.insert(std::make_pair("No", cv::Scalar(96, 96, 96)));
  color_map.insert(std::make_pair("Male", cv::Scalar(255, 0, 0)));
  color_map.insert(std::make_pair("Female", cv::Scalar(0, 0, 255)));

  std::set<std::string> tags;
  for (size_t i = 0; i < dr.object_candidates.size(); ++i) {
    const scouter::ObjectCandidate& o = dr.object_candidates[i];
    for (size_t j = 0; j < o.tags.size(); ++j) {
      tags.insert(o.tags[j].first);
    }
  }

  std::map<std::string, cv::Mat_<cv::Vec3b> >* ret;
  for (std::set<std::string>::const_iterator it = tags.begin();
       it != tags.end(); ++it) {
    cv::Mat_<cv::Vec3b> c = frame.image.clone();
    for (size_t i = 0; i < dr.object_candidates.size(); ++i) {
      const scouter::ObjectCandidate& o = dr.object_candidates[i];
      if (o.tags.size() == 0) {
        o.draw(c, cv::Scalar(96, 96, 96), 2);
      }
      for (size_t j = 0; j < o.tags.size(); ++j) {
        if (o.tags[j].first != *it) {
          continue;
        }
        o.draw(c, color_map[o.tags[j].second], 2);
        break;
      }
    }
    ret->insert(std::make_pair(*it, c));
  }
  return ret;
}

Taggers RecognizeDrawResult(Frame frame, DetectionResult dr) {
  scouter::Frame& fr = *static_cast<scouter::Frame*>(frame);
  scouter::DetectionResult& detected = *static_cast<scouter::DetectionResult*>(dr);
  std::map<std::string, cv::Mat_<cv::Vec3b> >* target = draw_result(fr, detected);
  return target;
}

Integrator Integrator_New(const char *config) {
  scouter::Integrator::Config ic = load_json<scouter::Integrator::Config>(config);
  return new scouter::Integrator(ic);
}

void Integrator_Delete(Integrator integrator) {
  delete static_cast<scouter::Integrator*>(integrator);
}

void Integrator_Push(Integrator integrator, Frame frame, DetectionResult dr) {
  scouter::Frame& fr = *static_cast<scouter::Frame*>(frame);
  scouter::DetectionResult& detected = *static_cast<scouter::DetectionResult*>(dr);
  scouter::Integrator& itr = *static_cast<scouter::Integrator*>(integrator);

  std::vector<scouter::Frame> frames;
  frames.push_back(fr);
  std::vector<scouter::DetectionResult> drs;
  drs.push_back(detected);
  itr.push(scouter::make_frames(frames), drs);
}

int Integrator_TrackerReady(Integrator integrator) {
  return static_cast<scouter::Integrator*>(integrator)->tracker_ready();
}

TrackingResult Integrator_Track(Integrator integrator) {
  scouter::Integrator& itr = *static_cast<scouter::Integrator*>(integrator);
  return new scouter::TrackingResult(itr.track());
}

InstanceManager InstanceManager_New(const char *config) {
  scouter::InstanceManager::Config ic = load_json<scouter::InstanceManager::Config>(config);
  return new scouter::InstanceManager(ic);
}

void InstanceManager_Delete(InstanceManager instanceManager) {
  delete static_cast<scouter::InstanceManager*>(instanceManager);
}

InstanceStates InstanceManager_GetCurrentStates(InstanceManager instanceManager,
                                                TrackingResult result) {
  scouter::InstanceManager im = *static_cast<scouter::InstanceManager*>(instanceManager);
  scouter::TrackingResult tr = *static_cast<scouter::TrackingResult*>(result);
  im.update(tr);
  std::vector<scouter::InstanceState> states = im.get_current_states();
  return new std::vector<scouter::InstanceState>(states);
}

void InstanceStates_Delete(InstanceStates states) {
  delete static_cast<std::vector<scouter::InstanceState>*>(states);
}

const char* ConvertStatesToJson(InstanceStates instanceStates, int floorID) {
  std::string dummy = "states";
  return dummy.c_str();
}


