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

// ImageTaggerCaffe is a bind of `scouter::ImageTaggerCaffe`.
type ImageTaggerCaffe struct {
	p C.ImageTaggerCaffe
}

// NewImageTaggerCaffe returns a new tagger.
func NewImageTaggerCaffe(config string) ImageTaggerCaffe {
	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	return ImageTaggerCaffe{p: C.ImageTaggerCaffe_New(cConfig)}
}

// Delete object.
func (t *ImageTaggerCaffe) Delete() {
	C.ImageTaggerCaffe_Delete(t.p)
	t.p = nil
}

// CropAndPredictTags crops image around the candidate and predict
// attributes with caffe models. The function executes these two tasks.
func (t *ImageTaggerCaffe) CropAndPredictTags(candidate Candidate,
	image MatVec3b) Candidate {

	recognized := C.ImageTaggerCaffe_CropAndPredictTags(t.p, candidate.p, image.p)
	return Candidate{p: recognized}
}

// CropAndPredictTagsBatch predicts attributes with caffe models at candidate
// information and set tags int the candidates. The function executes to predict
// attributes on batch.
func (t *ImageTaggerCaffe) CropAndPredictTagsBatch(candidates []Candidate,
	image MatVec3b) []Candidate {

	l := len(candidates)
	candidatePointer := convertCandidatesToPointer(candidates)
	recognized := C.ImageTaggerCaffe_CropAndPredictTagsBatch(t.p,
		(*C.Candidate)(&candidatePointer[0]), C.int(l), image.p)
	defer C.Candidates_Delete(recognized)

	return convertCandidatesToSlice(recognized)
}
