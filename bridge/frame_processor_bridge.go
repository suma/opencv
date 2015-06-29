package bridge

/*
#cgo darwin CXXFLAGS: -I/System/Library/Frameworks/Accelerate.framework/Versions/Current/Frameworks/vecLib.framework/Headers/ -DCPU_ONLY
#cgo LDFLAGS: -ljsonconfig
#cgo pkg-config: scouter-core
#cgo pkg-config: pficv
#cgo pkg-config: pficommon
#include <stdlib.h>
#include "frame_processor_bridge.h"
*/
import "C"
import (
	"unsafe"
)

type FrameProcessor struct {
	p C.FrameProcessor
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

func (fp *FrameProcessor) Projection(buf MatVec3b) (MatVec3b, int, int) {
	frame := C.FrameProcessor_Projection(fp.p, buf.p)
	img := MatVec3b{p: frame.image}
	offsetX := int(frame.offset_x)
	offsetY := int(frame.offset_y)
	return img, offsetX, offsetY
}
