#include "config_bridge.h"

#include <scouter-core/frame_processor.hpp>

FrameProcessorConfig FrameProcessorConfig_New(const char *config) {
  return 0;
}

void FrameProcessorConfig_Delete(FrameProcessorConfig config) {
  delete static_cast<scouter::FrameProcessor::Config*>(config);
}