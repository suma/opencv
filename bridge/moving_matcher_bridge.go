package bridge

/*
#cgo pkg-config: scouter-core
#include "moving_matcher_bridge.h"
*/
import "C"

type RegionsWithCameraID struct {
	CameraID   int
	Candidates []Candidate
}

type MVCandidate struct {
	p C.MVCandidate
}

func (c MVCandidate) Serialize() []byte {
	b := C.MVCandidate_Serialize(c.p)
	defer C.ByteArray_Release(b)
	return ToGoBytes(b)
}

func DeserializeMVCandiate(c []byte) MVCandidate {
	b := toByteArray(c)
	return MVCandidate{p: C.MVCandidate_Deserialize(b)}
}

func convertCandidatezToPointer(regions []RegionsWithCameraID) []C.struct_RegionsWithCameraID {
	regionsPointers := []C.struct_RegionsWithCameraID{}
	for _, f := range regions {
		candidatePointers := convertCandidatesToPointer(f.Candidates) // -> []C.Candidate
		candidateVec := C.InvertCandidates((*C.Candidate)(&candidatePointers[0]),
			C.int(len(candidatePointers))) // -> C.struct_Candidates
		defer C.Candidates_Delete(candidateVec)
		f := C.struct_RegionsWithCameraID{
			candidates: C.struct_Candidates{
				candidateVec: candidateVec.candidateVec,
				length:       C.int(len(candidatePointers)),
			},
			cameraID: C.int(f.CameraID),
		}
		regionsPointers = append(regionsPointers, f)
	}
	return regionsPointers
}

func GetMatching(kThreashold float32, regions []RegionsWithCameraID) []MVCandidate {
	regionsPointers := convertCandidatezToPointer(regions) // -> []C.struct_RegionsWithCameraID
	mvCandidatePointers := C.MVOM_GetMatching((*C.struct_RegionsWithCameraID)(&regionsPointers[0]),
		C.int(len(regions)), C.float(kThreashold)) // -> vector<vector<ObjectCandidate>>
	defer C.MVCandidates_Delete(mvCandidatePointers)

	l := int(mvCandidatePointers.length)
	mvCandidates := make([]C.MVCandidate, l)
	C.ResolveMVCandidates(mvCandidatePointers, (*C.MVCandidate)(&mvCandidates[0])) // -> []C.MVCandidate

	ret := make([]MVCandidate, l)
	for i := 0; i < l; i++ {
		ret[i] = MVCandidate{p: mvCandidates[i]}
	}
	return ret
}
