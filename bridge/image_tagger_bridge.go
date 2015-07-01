package bridge

/*
#cgo darwin CXXFLAGS: -I/System/Library/Frameworks/Accelerate.framework/Versions/Current/Frameworks/vecLib.framework/Headers/ -DCPU_ONLY
#cgo LDFLAGS: -ljsonconfig
#cgo pkg-config: scouter-core
#cgo pkg-config: pficv
#cgo pkg-config: pficommon
#include <stdlib.h>
#include "detector_bridge.h"
#include "image_tagger_bridge.h"
*/
import "C"
import (
	"unsafe"
)

type ImageTaggerCaffe struct {
	p C.ImageTaggerCaffe
}

func NewImageTaggerCaffe(config string) ImageTaggerCaffe {
	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	return ImageTaggerCaffe{p: C.ImageTaggerCaffe_New(cConfig)}
}

func (t *ImageTaggerCaffe) Delete() {
	C.ImageTaggerCaffe_Delete(t.p)
	t.p = nil
}

func (t *ImageTaggerCaffe) PredictTagsBatch(candidates []Candidate, img MatVec3b) []Candidate {
	l := len(candidates)
	candidatePointer := ConvertCandidatesToPointer(candidates)
	recognizedVec := C.ImageTaggerCaffe_PredictTagsBatch(t.p,
		(*C.Candidate)(&candidatePointer[0]), C.int(l), img.p)
	defer C.Candidates_Delete(recognizedVec)
	recognizedLength := int(recognizedVec.length)
	recognized := make([]C.Candidate, recognizedLength)
	C.ResolveCandidates(recognizedVec, (*C.Candidate)(&recognized[0]))

	ret := make([]Candidate, recognizedLength)
	for i := 0; i < recognizedLength; i++ {
		ret[i] = Candidate{p: recognized[i]}
	}
	return ret
}
