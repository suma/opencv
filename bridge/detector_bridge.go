package bridge

/*
#cgo pkg-config: scouter-core
#include <stdlib.h>
#include "detector_bridge.h"
*/
import "C"
import (
	"reflect"
	"sync"
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
	candidatePointers := make([]C.Candidate, len(candidates))
	for i, c := range candidates {
		candidatePointers[i] = c.p
	}
	return candidatePointers
}

func convertCandidatesToSlice(candidates C.struct_Candidates) []Candidate {
	var cArray *C.Candidate = candidates.candidates
	length := int(candidates.length)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cArray)),
		Len:  length,
		Cap:  length,
	}
	goSlice := *(*[]C.Candidate)(unsafe.Pointer(&hdr))
	cs := make([]Candidate, length)
	for i, c := range goSlice {
		cs[i] = Candidate{p: c}
	}
	return cs
}

func (c Candidate) Delete() {
	C.Candidate_Delete(c.p)
	c.p = nil
}

type Detector struct {
	mu sync.RWMutex
	p  C.Detector
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

func (d *Detector) UpdateCameraParameter(config string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	C.Detector_UpdateCameraParameter(d.p, cConfig)
}

func (d *Detector) ACFDetect(img MatVec3b, offsetX int, offsetY int) []Candidate {
	d.mu.RLock()
	defer d.mu.RUnlock()

	candidates := C.Detector_ACFDetect(d.p, img.p, C.int(offsetX), C.int(offsetY))
	defer C.Candidates_Delete(candidates)

	return convertCandidatesToSlice(candidates)
}

func (d *Detector) FilterByMask(candidate Candidate) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	masked := C.Detector_FilterByMask(d.p, candidate.p)
	return int(masked) == 0
}

func (d *Detector) EstimateHeight(candidate *Candidate, offsetX int, offsetY int) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	C.Detector_EstimateHeight(d.p, candidate.p, C.int(offsetX), C.int(offsetY))
}

func (d *Detector) PutFeature(candidate *Candidate, img MatVec3b) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	C.Detector_PutFeature(d.p, candidate.p, img.p)
}

type MMDetector struct {
	mu sync.RWMutex
	p  C.MMDetector
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

func (d *MMDetector) UpdateCameraParameter(config string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	C.MMDetector_UpdateCameraParameter(d.p, cConfig)
}

func (d *MMDetector) MMDetect(img MatVec3b, offsetX int, offsetY int) []Candidate {
	d.mu.RLock()
	defer d.mu.RUnlock()

	candidates := C.MMDetector_MMDetect(d.p, img.p, C.int(offsetX), C.int(offsetY))
	defer C.Candidates_Delete(candidates)

	return convertCandidatesToSlice(candidates)
}

func (d *MMDetector) FilterByMask(candidate Candidate) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	masked := C.MMDetector_FilterByMask(d.p, candidate.p)
	return int(masked) == 0
}

func (d *MMDetector) EstimateHeight(candidate *Candidate, offsetX int, offsetY int) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	C.MMDetector_EstimateHeight(d.p, candidate.p, C.int(offsetX), C.int(offsetY))
}

// draw result should be called by each candidate,
// but think the cost of copying MatVec3b, called by []C.Candidate
func DrawDetectionResult(img MatVec3b, candidates []Candidate) MatVec3b {
	l := len(candidates)
	candidatePointer := convertCandidatesToPointer(candidates)
	ret := C.Candidates_Draw(img.p, (*C.Candidate)(&candidatePointer[0]), C.int(l))
	return MatVec3b{p: ret}
}

func DrawDetectionResultWithTags(img MatVec3b, candidates []Candidate) MatVec3b {
	l := len(candidates)
	candidatePointer := convertCandidatesToPointer(candidates)
	ret := C.Candidates_DrawTags(img.p, (*C.Candidate)(&candidatePointer[0]), C.int(l))
	return MatVec3b{p: ret}
}
