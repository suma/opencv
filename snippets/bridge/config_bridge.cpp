#include "config_bridge.h"

#include <sstream>
#include <pficommon/text/json.h>
#include <jsonconfig.hpp>
#include <scouter-core/frame_processor.hpp>
#include <scouter-core/detector.hpp>
#include <scouter-core/image_tagger_caffe.hpp>
#include <scouter-core/integrator.hpp>

template <class Type>
void load_json(const char *config, Type *configPointer) {
  std::stringstream ss(config);
  pfi::text::json::json config_raw;
  ss >> config_raw;
  *configPointer = jsonconfig::config_cast<Type>(jsonconfig::config_root(config_raw));
}

FrameProcessorConfig FrameProcessorConfig_New(const char *config) {
  scouter::FrameProcessor::Config *fpc = new scouter::FrameProcessor::Config();
  load_json<scouter::FrameProcessor::Config>(config, fpc);
  return fpc;
}

void FrameProcessorConfig_Delete(FrameProcessorConfig config) {
  delete static_cast<scouter::FrameProcessor::Config*>(config);
}

DetectorConfig DetectorConfig_New(const char *config) {
  scouter::Detector::Config *dc = new scouter::Detector::Config();
  load_json<scouter::Detector::Config>(config, dc);
  return dc;
}

void DetectorConfig_Delete(DetectorConfig config) {
  delete static_cast<scouter::Detector::Config*>(config);
}

RecognizeConfigTaggers RecognizeConfigTaggers_New(const char *config) {
  std::vector<scouter::ImageTaggerCaffe::Config> *taggers =
      new std::vector<scouter::ImageTaggerCaffe::Config>();
  load_json<std:: vector<scouter::ImageTaggerCaffe::Config> >(config, taggers);
  return taggers;
}

void RecognizeConfigTaggers_Delete(RecognizeConfigTaggers taggers) {
  delete static_cast<std::vector<scouter::ImageTaggerCaffe::Config>* >(taggers);
}

IntegratorConfig IntegratorConfig_New(const char *config) {
  scouter::Integrator::Config *ic = new scouter::Integrator::Config();
  load_json<scouter::Integrator::Config>(config, ic);
  return ic;
}

void IntegratorConfig_Delete(IntegratorConfig config) {
  delete static_cast<scouter::Integrator::Config*>(config);
}
