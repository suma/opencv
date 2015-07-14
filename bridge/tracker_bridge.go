package bridge

/*
#cgo darwin CXXFLAGS: -I/System/Library/Frameworks/Accelerate.framework/Versions/Current/Frameworks/vecLib.framework/Headers/ -DCPU_ONLY
#cgo LDFLAGS: -ljsonconfig
#cgo pkg-config: scouter-core
#cgo pkg-config: pficv
#cgo pkg-config: pficommon
#include <stdlib.h>
#include "tracker_bridge.h"
*/
import "C"

type RegionsWithCamerID struct {
	cameraID   int
	candidates []Candidate
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

func convertCandidatezToPointer(regions []RegionsWithCamerID) []C.struct_RegionsWithCamerID {
	regionsPointers := []C.struct_RegionsWithCamerID{}
	for _, f := range regions {
		candidatePointers := convertCandidatesToPointer(f.candidates) // -> []C.Candidate
		candidateVec := C.InvertCandidates((*C.Candidate)(&candidatePointers[0]),
			C.int(len(candidatePointers))) // -> C.struct_Candidates
		defer C.Candidates_Delete(candidateVec)
		f := C.struct_RegionsWithCamerID{
			candidates: C.struct_Candidates{
				candidateVec: candidateVec.candidateVec,
				length:       C.int(len(candidatePointers)),
			},
			cameraID: C.int(f.cameraID),
		}
		regionsPointers = append(regionsPointers, f)
	}
	return regionsPointers
}

func GetMatching(kThreashold float32, regions ...RegionsWithCamerID) []MVCandidate {
	regionsPointers := convertCandidatezToPointer(regions) // -> []C.struct_RegionsWithCamerID
	mvCandidatePointers := C.MVOM_GetMatching((*C.struct_RegionsWithCamerID)(&regionsPointers[0]),
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
