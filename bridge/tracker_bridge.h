#ifndef _TRACKER_BRIDGE_H_
#define _TRACKER_BRIDGE_H_

#include "opencv_bridge.h"
#include "detector_bridge.h"

#ifdef __cplusplus
extern "C" {
#endif

#ifdef __cplusplus
typedef struct Candidatez {
  std::vector<std::vector<scouter::ObjectCandidate> >* candidatesVec;
  int length;
} Candidatez;
#else
typedef struct Candidatez {
  void* candidatesVec;
  int length;
} Candidatez;
#endif

void ResolveCandidatez(struct Candidatez candidatez, Candidates* obj);
void Candidatez_Delete(struct Candidatez candidatesVec);
struct Candidatez MVOM_GetMatching(Candidate** candidatez, int* lengths, int length, float kThreshold);

#ifdef __cplusplus
}
#endif

#endif //_TRACKER_BRIDGE_H_