#ifndef _DETECTOR_BRIDGE_H_
#define _DETECTOR_BRIDGE_H_

#include "opencv_bridge.h"

#ifdef __cplusplus
#include <scouter-core/detector.hpp>
#include <scouter-core/mm_detector.hpp>
extern "C" {
#endif

#ifdef __cplusplus
typedef scouter::Detector* Detector;
typedef scouter::MultiModelDetector* MMDetector;
typedef scouter::ObjectCandidate* Candidate;
typedef struct Candidates {
  std::vector<scouter::ObjectCandidate>* candidateVec;
  int length;
} Candidates;
#else
typedef void* Detector;
typedef void* MMDetector;
typedef void* Candidate;
typedef struct Candidates {
  void* candidateVec;
  int length;
} Candidates;
#endif

struct ByteArray Candidate_Serialize(Candidate c);
Candidate Candidate_Deserialize(struct ByteArray src);
void Candidate_Delete(Candidate c);

void ResolveCandidates(struct Candidates candidates, Candidate* obj);
void Candidates_Delete(struct Candidates candidates);

Detector Detector_New(const char *config);
void Detector_Delete(Detector detector);

struct Candidates Detector_ACFDetect(Detector detector, MatVec3b image, int offsetX, int offsetY);
int Detector_FilterByMask(Detector detector, Candidate candidate);
void Detector_EstimateHeight(Detector detector, Candidate candidate, int offsetX, int offsetY);
void Detector_PutFeature(Detector detector, Candidate candidate, MatVec3b image);

MMDetector MMDetector_New(const char *config);
void MMDetector_Delete(MMDetector detector);

struct Candidates MMDetector_MMDetect(MMDetector detector, MatVec3b image, int offsetX, int offsetY);
int MMDetector_FilterByMask(MMDetector detector, Candidate candidate);
void MMDetector_EstimateHeight(MMDetector detector, Candidate candidate, int offsetX, int offsetY);

MatVec3b Candidates_Draw(MatVec3b image, Candidate* candidates, int length);

#ifdef __cplusplus
}
#endif

#endif //_DETECTOR_BRIDGE_H_
