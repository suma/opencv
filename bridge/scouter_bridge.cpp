#include "scouter_bridge.h"
#include "util.hpp"

#include <iostream>
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

struct ScouterFrame FrameProcessor_Projection(FrameProcessor fp, MatVec3b buf) {
  scouter::FrameMeta meta = scouter::FrameMeta();
  scouter::Frame* frame = new scouter::Frame(fp->apply(*buf, meta));
  ScouterFrame result = {
    &(frame->image),
    frame->meta.offset_x,
    frame->meta.offset_y,
  };
  return result;
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

MultiModelDetector MultiModelDetector_New(const char *config) {
  scouter::MultiModelDetector::Config dc
    = load_json<scouter::MultiModelDetector::Config>(config);
  return new scouter::MultiModelDetector(dc);
}

void MultiModelDetector_Delete(MultiModelDetector detector) {
  delete detector;
}

DetectionResult MultiModelDetector_Detect(MultiModelDetector detector, Frame frame) {
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
  delete taggers;
}

DetectionResult Recognize(ImageTaggerCaffe taggers, Frame frame, DetectionResult dr) {
  for (size_t i = 0; i < taggers->size(); ++i) {
    scouter::ImageTaggerCaffe tagger = (*taggers)[i];
    tagger.predict_tags_batch(dr->object_candidates, *frame);
  }
  return dr;
}

std::map<std::string, cv::Mat_<cv::Vec3b> > draw_result(
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

  std::map<std::string, cv::Mat_<cv::Vec3b> > ret;
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
    ret.insert(std::make_pair(*it, c));
  }
  return ret;
}

Taggers RecognizeDrawResult(Frame frame, DetectionResult dr) {
  std::map<std::string, cv::Mat_<cv::Vec3b> > target = draw_result(*frame, *dr);
  std::map<std::string, cv::Mat_<cv::Vec3b> >* taggers =
      new std::map<std::string, cv::Mat_<cv::Vec3b> >();
  for (std::map<std::string, cv::Mat_<cv::Vec3b> >::iterator it = target.begin();
      it != target.end(); it++) {
    taggers->insert(std::make_pair(it->first, it->second));
  }
  Taggers result = {taggers, (int)taggers->size()};
  return result;
}

void Taggers_Delete(Taggers taggers) {
  delete taggers.drawResultsMap;
}

void ResolveDrawResult(struct Taggers taggers, const char** keys, MatVec3b* drawResults) {
  std::map<std::string, cv::Mat_<cv::Vec3b> > resultMap = *taggers.drawResultsMap;
  int i = 0;
  for (std::map<std::string, cv::Mat_<cv::Vec3b> >::iterator it = resultMap.begin();
      it != resultMap.end(); it++) {
    const char* key = it->first.c_str();
    keys[i] = key;
    drawResults[i] = new cv::Mat_<cv::Vec3b>(it->second);
    i++;
  }
}

Integrator Integrator_New(const char *config) {
  scouter::Integrator::Config ic = load_json<scouter::Integrator::Config>(config);
  return new scouter::Integrator(ic);
}

void Integrator_Delete(Integrator integrator) {
  delete integrator;
}

void Integrator_Push(Integrator integrator, Frame* frame, DetectionResult* dr, int size) {
  std::vector<scouter::Frame> frames;
  std::vector<scouter::DetectionResult> drs;
  for (int i = 0; i < size; i++) {
    frames.push_back(*(frame[i]));
    drs.push_back(*(dr[i]));
  }
  integrator->push(scouter::make_frames(frames), drs);
}

int Integrator_TrackerReady(Integrator integrator) {
  return integrator->tracker_ready();
}

TrackingResult Integrator_Track(Integrator integrator) {
  return new scouter::TrackingResult(integrator->track());
}

void TrackingResult_Delete(TrackingResult tr) {
  delete tr;
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

String ConvertStatesToJson(InstanceStates instanceStates,
                                int floorID, long long timestamp) {
  using pfi::text::json::json;
  using pfi::text::json::json_object;
  using pfi::text::json::json_array;
  using pfi::text::json::json_integer;
  using pfi::text::json::json_string;

  json instances_json(new json_array);
  for (size_t i = 0; i < instanceStates->size(); ++i) {
    const scouter::InstanceState s = (*instanceStates)[i];
    json instance_json(new json_object);
    instance_json["id"] = json(new json_integer(s.id));
    instance_json["location"] = json(new json_object);
    instance_json["location"]["x"] = json(new json_integer(s.position.x));
    instance_json["location"]["y"] = json(new json_integer(s.position.y));
    instance_json["location"]["floor_id"] = json(new json_integer(floorID));
    instance_json["labels"] = json(new json_array);
    for (size_t j = 0; j < s.tags.size(); ++j) {
      const std::pair<std::string, std::string>& tag = s.tags[j];
      instance_json["labels"].add(json(new json_string(tag.first + "=" + tag.second)));
    }
    instances_json.add(instance_json);
  }

  json ret(new json_object);
  ret["time"] = json(new json_integer(timestamp));
  ret["instances"] = instances_json;

  // json to string
  std::stringstream ss;
  ss << ret;
  String str = {ss.str().c_str(), (int)ss.str().size()};
  return str;
}

Visualizer Visualizer_New(const char *config, InstanceManager instanceManager) {
  scouter::InstancesVisualizer::Config vc = load_json<scouter::InstancesVisualizer::Config>(config);
  return new scouter::InstancesVisualizer(vc, *instanceManager);
}

void Visualizer_Delete(Visualizer visualizer) {
  delete visualizer;
}

PlotTrajectories Visualizer_PlotTrajectories(Visualizer visualizer) {
  std::vector<cv::Mat_<cv::Vec3b> > traj_plots =
      visualizer->plot_trajectories();
  std::vector<cv::Mat_<cv::Vec3b> >* results = new std::vector<cv::Mat_<cv::Vec3b> >();
  for (size_t i = 0; i < traj_plots.size(); ++i) {
    results->push_back(traj_plots[i]);
  }
  PlotTrajectories pt = {results, (int)traj_plots.size()};
  return pt;
}

void PlotTrajectories_Delete(PlotTrajectories plotTrajectories) {
  delete plotTrajectories.trajectories;
}

void ResolvePlotTrajectories(struct PlotTrajectories plotTrajectories, MatVec3b* trajectories) {
  for (size_t i = 0; i < plotTrajectories.trajectories->size(); ++i) {
    trajectories[i] = new cv::Mat_<cv::Vec3b>((*plotTrajectories.trajectories)[i]);
  }
}
