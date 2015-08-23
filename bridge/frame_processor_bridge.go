package bridge

/*
#cgo pkg-config: scouter-core
#include <stdlib.h>
#include "frame_processor_bridge.h"
*/
import "C"
import (
	"sync"
	"unsafe"
)

// ScouterFrame is a bind of `scouter::Frame`
type ScouterFrame struct {
	Image     MatVec3b
	OffsetX   int
	OffsetY   int
	Timestamp uint64
	CameraID  int
}

// FrameProcessor is a bind of `scouter::FrameProcessor`.
type FrameProcessor struct {
	mu sync.RWMutex
	p  C.FrameProcessor
}

// NewFrameProcessor returns a new frame processor.
func NewFrameProcessor(config string) FrameProcessor {
	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	return FrameProcessor{p: C.FrameProcessor_New(cConfig)}
}

// Delete object.
func (fp *FrameProcessor) Delete() {
	C.FrameProcessor_Delete(fp.p)
	fp.p = nil
}

// UpdateConfig updates camera parameter in the frame processor.
func (fp *FrameProcessor) UpdateConfig(config string) {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	C.FrameProcessor_UpdateConfig(fp.p, cConfig)
}

// Projection the image data with frame processor parameters, and returns with
// offset information.
func (fp *FrameProcessor) Projection(buf MatVec3b) (MatVec3b, int, int) {
	fp.mu.RLock()
	defer fp.mu.RUnlock()

	frame := C.FrameProcessor_Projection(fp.p, buf.p)
	img := MatVec3b{p: frame.image}
	offsetX := int(frame.offset_x)
	offsetY := int(frame.offset_y)
	return img, offsetX, offsetY
}
