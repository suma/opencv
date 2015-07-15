#ifndef _MOVING_MATCHER_BRIDGE_H_
#define _MOVING_MATCHER_BRIDGE_H_

#include "opencv_bridge.h"
#include "detector_bridge.h"

#ifdef __cplusplus
#include <scouter-core/mv_detection_result.hpp>
extern "C" {
#endif

#ifdef __cplusplus
typedef scouter::MVObjectCandidate* MVCandidate;
typedef struct MVCandidates {
  std::vector<scouter::MVObjectCandidate>* candidateVec;
  int length;
} MVCandidates;
#else
typedef void* MVCandidate;
typedef struct MVCandidates {
  void* candidateVec;
  int length;
} MVCandidates;
#endif
typedef struct RegionsWithCameraID {
  Candidates candidates;
  int cameraID;
} RegionsWithCameraID;

struct ByteArray MVCandidate_Serialize(MVCandidate c);
MVCandidate MVCandidate_Deserialize(struct ByteArray src);
void MVCandidate_Delete(MVCandidate c);

void ResolveMVCandidates(struct MVCandidates mvCandidates, MVCandidate* obj);
struct MVCandidates InvertMVCandidates(MVCandidate* obj, int length);
void MVCandidates_Delete(struct MVCandidates mcCandidates);
struct MVCandidates MVOM_GetMatching(RegionsWithCameraID* regions, int length, float kThreshold);

#ifdef __cplusplus
}
#endif

#endif //_MOVING_MATCHER_BRIDGE_H_