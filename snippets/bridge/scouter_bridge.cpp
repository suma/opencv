#include "scouter_bridge.h"
#include "util.hpp"

#include <string>
#include <sstream>
#include <pficommon/text/json.h>
#include <jsonconfig.hpp>
#include <opencv2/opencv.hpp>
#include <scouter-core/mv_detection_result.hpp>

template <class Type>
Type load_json(const char *config) {
  std::stringstream ss(config);
  pfi::text::json::json config_raw;
  ss >> config_raw;
  return jsonconfig::config_cast<Type>(jsonconfig::config_root(config_raw));
}

struct ByteArray Frame_Serialize(Frame f) {
  return serializeObject(*f);
}

Frame Freme_Deserialize(struct ByteArray src) {
  return deserializeObject<scouter::Frame>(src);
}

void Frame_Delete(Frame f) {
  delete f;
}

struct ByteArray DetectionResult_Serialize(DetectionResult dr) {
  return serializeObject(*dr);
}

DetectionResult DetectionResult_Deserialize(struct ByteArray src) {
  return deserializeObject<scouter::DetectionResult>(src);
}

void DetectionResult_Delete(DetectionResult dr) {
  delete dr;
}

FrameProcessor FrameProcessor_New(const char *config) {
  scouter::FrameProcessor::Config fpc =
      load_json<scouter::FrameProcessor::Config>(config);
  return new scouter::FrameProcessor(fpc);
}

void FrameProcessor_Delete(FrameProcessor fp) {
  delete fp;
}

Frame FrameProcessor_Apply(FrameProcessor fp, MatVec3b buf,
                           long long timestamp, int cameraID) {
  scouter::FrameMeta meta(timestamp, cameraID);
  return new scouter::Frame(fp->apply(*buf, meta));
}

Detector Detector_New(const char *config) {
  scouter::Detector::Config dc = load_json<scouter::Detector::Config>(config);
  return new scouter::Detector(dc);
}

void Detector_Delete(Detector detector) {
  delete detector;
}

DetectionResult Detector_Detect(Detector detector, Frame frame) {
  return new scouter::DetectionResult(detector->detect(*frame));
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
  cv::Mat_<cv::Vec3b>* target = new cv::Mat_<cv::Vec3b>();
  draw_result(*frame, *dr, ms, *target);
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
  for (size_t i = 0; i < taggers->size(); ++i) {
    scouter::ImageTaggerCaffe tagger = (*taggers)[i];
    tagger.predict_tags_batch(dr->object_candidates, *frame);
  }
  return dr;
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
  std::map<std::string, cv::Mat_<cv::Vec3b> >* target = draw_result(*frame, *dr);
  return target;
}

Integrator Integrator_New(const char *config) {
  scouter::Integrator::Config ic = load_json<scouter::Integrator::Config>(config);
  return new scouter::Integrator(ic);
}

void Integrator_Delete(Integrator integrator) {
  delete integrator;
}

void Integrator_Push(Integrator integrator, Frame frame, DetectionResult dr) {
  std::vector<scouter::Frame> frames;
  frames.push_back(*frame);
  std::vector<scouter::DetectionResult> drs;
  drs.push_back(*dr);
  integrator->push(scouter::make_frames(frames), drs);
}

int Integrator_TrackerReady(Integrator integrator) {
  return integrator->tracker_ready();
}

TrackingResult Integrator_Track(Integrator integrator) {
  return new scouter::TrackingResult(integrator->track());
}

InstanceManager InstanceManager_New(const char *config) {
  scouter::InstanceManager::Config ic = load_json<scouter::InstanceManager::Config>(config);
  return new scouter::InstanceManager(ic);
}

void InstanceManager_Delete(InstanceManager instanceManager) {
  delete instanceManager;
}

InstanceStates InstanceManager_GetCurrentStates(InstanceManager instanceManager,
                                                TrackingResult result) {
  instanceManager->update(*result);
  std::vector<scouter::InstanceState> states = instanceManager->get_current_states();
  return new std::vector<scouter::InstanceState>(states);
}

void InstanceStates_Delete(InstanceStates states) {
  delete states;
}

const char* ConvertStatesToJson(InstanceStates instanceStates, int floorID) {
  std::string dummy = "states";
  return dummy.c_str();
}


