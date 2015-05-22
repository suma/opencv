#include "config_bridge.h"

#include <sstream>
#include <pficommon/text/json.h>
#include <jsonconfig.hpp>
#include <scouter-core/frame_processor.hpp>

FrameProcessorConfig FrameProcessorConfig_New(const char *config) {
  std::stringstream ss(config);
  pfi::text::json::json config_raw;
  ss >> config_raw;
  scouter::FrameProcessor::Config *fpc = new scouter::FrameProcessor::Config();
  *fpc = jsonconfig::config_cast<scouter::FrameProcessor::Config>(
    jsonconfig::config_root(config_raw));
  return fpc;
}

void FrameProcessorConfig_Delete(FrameProcessorConfig config) {
  delete static_cast<scouter::FrameProcessor::Config*>(config);
}