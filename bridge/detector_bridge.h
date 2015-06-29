#ifndef _DETECTOR_BRIDGE_H_
#define _DETECTOR_BRIDGE_H_

#include "opencv_bridge.h"

#ifdef __cplusplus
#include <scouter-core/detector.hpp>
extern "C" {
#endif

#ifdef __cplusplus
typedef scouter::Detector* Detector;
typedef scouter::ObjectCandidate* Candidate;
typedef struct Candidates {
  std::vector<scouter::ObjectCandidate>* candidatesVec;
  int length;
} Candidates;
#else
typedef void* Detector;
typedef void* Candidate;
typedef struct Candidates {
  void* candidatesVec;
  int length;
} Candidates;
#endif

struct ByteArray Candidate_Serialize(Candidate c);
Candidate Candidate_Deserialize(struct ByteArray src);
void Candidate_Delete(Candidate c);

Detector Detector_New(const char *config);
void Detector_Delete(Detector detector);
void ResolveCondidates(struct Candidates candidates, Candidate* obj);
struct Candidates Detector_ACFDetect(Detector detector, MatVec3b image, int offsetX, int offssetY);
struct Candidates Detector_FilterCndidateByMask(struct Candidates candidates);
void Detector_EstimateCandidateHeight(struct Candidates candidates, int offsetX, int offsetY);

#ifdef __cplusplus
}
#endif

#endif //_DETECTOR_BRIDGE_H_
