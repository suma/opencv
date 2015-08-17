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

// Crop the image around the candidate (=region) information, and returns
// image data.
func (t *ImageTaggerCaffe) Crop(candidate Candidate, image MatVec3b) MatVec3b {
	cropped := C.ImageTaggerCaffe_Crop(t.p, candidate.p, image.p)
	return MatVec3b{p: cropped}
}

// PredictTags predicts attributes with caffe models at candidate information
// and set tags in the candidate.
func (t *ImageTaggerCaffe) PredictTags(candidate Candidate,
	croppedImg MatVec3b) Candidate {

	recognized := C.ImageTaggerCaffe_PredictTags(t.p, candidate.p, croppedImg.p)
	return Candidate{p: recognized}
}

// PredictTagsBatch predicts attributes with caffe models at candidate
// information and set tags int the candidates. The function executes to predict
// attributes on batch.
func (t *ImageTaggerCaffe) PredictTagsBatch(candidates []Candidate,
	croppedImg []MatVec3b) []Candidate {

	l := len(candidates)
	candidatePointer := convertCandidatesToPointer(candidates)
	imgPointer := convertMatVec3bsToPointer(croppedImg)
	recognized := C.ImageTaggerCaffe_PredictTagsBatch(t.p,
		(*C.Candidate)(&candidatePointer[0]), (*C.MatVec3b)(&imgPointer[0]), C.int(l))
	defer C.Candidates_Delete(recognized)

	return convertCandidatesToSlice(recognized)
}

// CroppingAndPredictTags crops image around the candidate and predict
// attributes with caffe models. The function executes these two tasks.
func (t *ImageTaggerCaffe) CroppingAndPredictTags(candidate Candidate,
	image MatVec3b) Candidate {

	recognized := C.ImageTaggerCaffe_CropAndPredictTags(t.p, candidate.p, image.p)
	return Candidate{p: recognized}
}
