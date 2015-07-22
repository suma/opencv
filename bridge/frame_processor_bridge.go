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

type FrameProcessor struct {
	mu sync.RWMutex
	p  C.FrameProcessor
}

func NewFrameProcessor(config string) FrameProcessor {
	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	return FrameProcessor{p: C.FrameProcessor_New(cConfig)}
}

func (fp *FrameProcessor) Delete() {
	C.FrameProcessor_Delete(fp.p)
	fp.p = nil
}

func (fp *FrameProcessor) UpdateConfig(config string) {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	C.FrameProcessor_UpdateConfig(fp.p, cConfig)
}

func (fp *FrameProcessor) Projection(buf MatVec3b) (MatVec3b, int, int) {
	fp.mu.RLock()
	defer fp.mu.RUnlock()

	frame := C.FrameProcessor_Projection(fp.p, buf.p)
	img := MatVec3b{p: frame.image}
	offsetX := int(frame.offset_x)
	offsetY := int(frame.offset_y)
	return img, offsetX, offsetY
}
