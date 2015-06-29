package bridge

/*
#cgo darwin CXXFLAGS: -I/System/Library/Frameworks/Accelerate.framework/Versions/Current/Frameworks/vecLib.framework/Headers/ -DCPU_ONLY
#cgo LDFLAGS: -ljsonconfig
#cgo pkg-config: scouter-core
#cgo pkg-config: pficv
#cgo pkg-config: pficommon
#include <stdlib.h>
#include "detector_bridge.h"
*/
import "C"
import (
	"unsafe"
)

const (
	// CvCapPropFrameWidth is OpenCV parameter of Frame Width
	CvCapPropFrameWidth = 3
	// CvCapPropFrameHeight is OpenCV parameter of Frame Height
	CvCapPropFrameHeight = 4
	// CvCapPropFps is OpenCV parameter of FPS
	CvCapPropFps = 5
)

type Candidate struct {
	p C.Candidate
}

func (c Candidate) Serialize() []byte {
	b := C.Candidate_Serialize(c.p)
	defer C.ByteArray_Release(b)
	return ToGoBytes(b)
}

func DeserializeCandidate(c []byte) Candidate {
	b := toByteArray(c)
	return Candidate{p: C.Candidate_Deserialize(b)}
}

func (c Candidate) Delete() {
	C.Candidate_Delete(c.p)
	c.p = nil
}

type Detector struct {
	p C.Detector
}

func NewDetector(config string) Detector {
	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	return Detector{p: C.Detector_New(cConfig)}
}

func (d *Detector) Delete() {
	C.Detector_Delete(d.p)
	d.p = nil
}

func (d *Detector) ACFDetect(img MatVec3b, offsetX int, offsetY int) []Candidate {
	candidateVecPointer := C.Detector_ACFDetect(d.p, img.p, C.int(offsetX), C.int(offsetY))
	defer C.Candidates_Delete(candidateVecPointer)
	l := int(candidateVecPointer.length)
	candidates := make([]C.Candidate, l)
	C.ResolveCandidates(candidateVecPointer, (*C.Candidate)(&candidates[0]))

	ret := make([]Candidate, l)
	for i := 0; i < l; i++ {
		ret[i] = Candidate{p: candidates[i]}
	}
	return ret
}

func (d *Detector) FilterByMask(candidates [][]byte) []Candidate {
	l := len(candidates)
	candidatePointer := make([]C.Candidate, l)
	for _, candidate := range candidates {
		candidatePointer = append(candidatePointer, DeserializeCandidate(candidate).p)
	}
	filteredVec := C.Detector_FilterCandidateByMask(d.p, (*C.Candidate)(&candidatePointer[0]), C.int(l))
	defer C.Candidates_Delete(filteredVec)
	filteredLength := int(filteredVec.length)
	filtered := make([]C.Candidate, filteredLength)
	C.ResolveCandidates(filteredVec, (*C.Candidate)(&filtered[0]))

	ret := make([]Candidate, filteredLength)
	for i := 0; i < filteredLength; i++ {
		ret[i] = Candidate{p: filtered[i]}
	}
	return ret
}

func (d *Detector) EstimateHeight(candidates [][]byte, offsetX int, offsetY int) []Candidate {
	l := len(candidates)
	candidatePointer := make([]C.Candidate, l)
	for _, candidate := range candidates {
		candidatePointer = append(candidatePointer, DeserializeCandidate(candidate).p)
	}
	estimatedVec := C.Detector_EstimateCandidateHeight(d.p, (*C.Candidate)(&candidatePointer[0]),
		C.int(l), C.int(offsetX), C.int(offsetY))
	defer C.Candidates_Delete(estimatedVec)
	estimatedLength := int(estimatedVec.length)
	estimated := make([]C.Candidate, estimatedLength)
	C.ResolveCandidates(estimatedVec, (*C.Candidate)(&estimated[0]))

	ret := make([]Candidate, estimatedLength)
	for i := 0; i < estimatedLength; i++ {
		ret[i] = Candidate{p: estimated[i]}
	}
	return ret
}
