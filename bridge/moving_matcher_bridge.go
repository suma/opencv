package bridge

/*
#cgo pkg-config: scouter-core
#include "moving_matcher_bridge.h"
*/
import "C"
import (
	"reflect"
	"unsafe"
)

// RegionsWithCameraID is a utility structure to manage ObjectCandidates with
// camera ID.
type RegionsWithCameraID struct {
	CameraID   int
	Candidates []Candidate
}

// MVCandidate is a bind of `scouter::MVObjectCandidate`
type MVCandidate struct {
	p C.MVCandidate
}

// Serialize object.
func (c MVCandidate) Serialize() []byte {
	b := C.MVCandidate_Serialize(c.p)
	defer C.ByteArray_Release(b)
	return ToGoBytes(b)
}

// DeserializeMVCandidate deserializes object.
func DeserializeMVCandidate(c []byte) MVCandidate {
	b := toByteArray(c)
	return MVCandidate{p: C.MVCandidate_Deserialize(b)}
}

// Delete object.
func (c MVCandidate) Delete() {
	C.MVCandidate_Delete(c.p)
	c.p = nil
}

func convertCandidatezToPointer(
	regions []RegionsWithCameraID) []C.struct_RegionsWithCameraID {
	regionsPointers := []C.struct_RegionsWithCameraID{}
	for _, r := range regions {
		var candidates C.struct_Candidates
		// []Candidate -> []C.Candidate
		if len(r.Candidates) == 0 {
			candidates = C.struct_Candidates{
				length: C.int(0),
			}
		} else {
			candidatePointers := convertCandidatesToPointer(r.Candidates)
			// -> C.struct_Candidates
			candidates = C.struct_Candidates{
				candidates: (*C.Candidate)(&candidatePointers[0]),
				length:     C.int(len(candidatePointers)),
			}
		}
		f := C.struct_RegionsWithCameraID{
			candidates: candidates,
			cameraID:   C.int(r.CameraID),
		}
		regionsPointers = append(regionsPointers, f)
	}
	return regionsPointers
}

// GetMatching find matched ObjectCandidate and aggregation with
// MVObjectCandidate. "kthreshold" is used top-k algorithm.
func GetMatching(kthreshold float32, regions []RegionsWithCameraID) []MVCandidate {
	//  -> []C.struct_RegionsWithCameraID
	regionsPointers := convertCandidatezToPointer(regions)
	// -> *C.MVCandidate
	mvCandidatePointers := C.MVOM_GetMatching(
		(*C.struct_RegionsWithCameraID)(&regionsPointers[0]),
		C.int(len(regions)), C.float(kthreshold))
	defer C.MVCandidates_Delete(mvCandidatePointers)

	cArray := mvCandidatePointers.mvCandidates
	length := int(mvCandidatePointers.length)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cArray)),
		Len:  length,
		Cap:  length,
	}
	goSlice := *(*[]C.MVCandidate)(unsafe.Pointer(&hdr))

	ret := make([]MVCandidate, length)
	for i, c := range goSlice {
		ret[i] = MVCandidate{p: c}
	}
	return ret
}
