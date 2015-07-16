package bridge

/*
#cgo pkg-config: scouter-core
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

func convertCandidatesToPointer(candidates []Candidate) []C.Candidate {
	candidatePointers := []C.Candidate{}
	for _, c := range candidates {
		candidatePointers = append(candidatePointers, c.p)
	}
	return candidatePointers
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

func (d *Detector) FilterByMask(candidate Candidate) bool {
	masked := C.Detector_FilterByMask(d.p, candidate.p)
	return int(masked) == 0
}

func (d *Detector) EstimateHeight(candidate *Candidate, offsetX int, offsetY int) {
	C.Detector_EstimateHeight(d.p, candidate.p, C.int(offsetX), C.int(offsetY))
}

func (d *Detector) PutFeature(candidate *Candidate, img MatVec3b) {
	C.Detector_PutFeature(d.p, candidate.p, img.p)
}

type MMDetector struct {
	p C.MMDetector
}

func NewMMDetector(config string) MMDetector {
	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	return MMDetector{p: C.MMDetector_New(cConfig)}
}

func (d *MMDetector) Delete() {
	C.MMDetector_Delete(d.p)
	d.p = nil
}

func (d *MMDetector) MMDetect(img MatVec3b, offsetX int, offsetY int) []Candidate {
	candidateVecPointer := C.MMDetector_MMDetect(d.p, img.p, C.int(offsetX), C.int(offsetY))
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

func (d *MMDetector) FilterByMask(candidate Candidate) bool {
	masked := C.MMDetector_FilterByMask(d.p, candidate.p)
	return int(masked) == 0
}

func (d *MMDetector) EstimateHeight(candidate *Candidate, offsetX int, offsetY int) {
	C.MMDetector_EstimateHeight(d.p, candidate.p, C.int(offsetX), C.int(offsetY))
}

func DrawDetectionResult(img MatVec3b, candidates []Candidate) MatVec3b {
	// draw result should be called by each candidate,
	// but think the cost of copying MatVec3b, called by []C.Candidate
	l := len(candidates)
	candidatePointer := convertCandidatesToPointer(candidates)
	ret := C.Candidates_Draw(img.p, (*C.Candidate)(&candidatePointer[0]), C.int(l))
	return MatVec3b{p: ret}
}
