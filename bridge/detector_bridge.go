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
	return nil
}
