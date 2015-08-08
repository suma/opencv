#ifndef _IMAGE_TAGGER_BRIDGE_H_
#define _IMAGE_TAGGER_BRIDGE_H_

#include "opencv_bridge.h"
#include "detector_bridge.h"

#ifdef __cplusplus
#include <scouter-core/image_tagger_caffe.hpp>
extern "C" {
#endif

#ifdef __cplusplus
typedef scouter::ImageTaggerCaffe* ImageTaggerCaffe;
#else
typedef void* ImageTaggerCaffe;
#endif

ImageTaggerCaffe ImageTaggerCaffe_New(const char* config);
void ImageTaggerCaffe_Delete(ImageTaggerCaffe tagger);
MatVec3b ImageTaggerCaffe_Crop(ImageTaggerCaffe tagger, Candidate candidate,
  MatVec3b image);
Candidate ImageTaggerCaffe_PredictTags(ImageTaggerCaffe tagger, Candidate candidate,
  MatVec3b cropedImg);
struct Candidates ImageTaggerCaffe_PredictTagsBatch(ImageTaggerCaffe tagger,
  Candidate* candidates, MatVec3b* croppedImages, int length);
Candidate ImageTaggerCaffe_CropAndPredictTags(ImageTaggerCaffe tagger,
  Candidate candidate, MatVec3b image);

#ifdef __cplusplus
}
#endif

#endif //_IMAGE_TAGGER_BRIDGE_H_