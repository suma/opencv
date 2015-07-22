#include "frame_processor_bridge.h"
#include "util.hpp"

FrameProcessor FrameProcessor_New(const char *config) {
  scouter::FrameProcessor::Config fpc =
      load_json<scouter::FrameProcessor::Config>(config);
  return new scouter::FrameProcessor(fpc);
}

void FrameProcessor_Delete(FrameProcessor fp) {
  delete fp;
}

void FrameProcessor_UpdateConfig(FrameProcessor fp, const char *config) {
  scouter::FrameProcessor::Config fpc =
      load_json<scouter::FrameProcessor::Config>(config);
  fp->update_config(fpc);
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
