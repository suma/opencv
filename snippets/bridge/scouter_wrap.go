package bridge

/*
#include "scouter_bridge.h"
*/
import "C"
import (
	"unsafe"
)

type FrameProcessorConfig C.FrameProcessorConfig
type FrameProcessor C.FrameProcessor
type Frame C.Frame

func FrameProcessor_SetUp(config FrameProcessorConfig) FrameProcessor {
	return FrameProcessor(
		C.FrameProcessor_SetUp(C.FrameProcessorConfig(config)))
}

func FrameProcessor_Apply(fp FrameProcessor, buf MatVec3b,
	timestamp int64, cameraId int) []byte {
	var frame Frame
	var frameLength int
	C.FrameProcessor_Apply(
		C.FrameProcessor(fp), C.MatVec3b(buf),
		C.longlong(timestamp), C.int(cameraId),
		C.Frame(frame), C.int(frameLength))
	return C.GoBytes(unsafe.Pointer(frame), C.int(frameLength))
}
