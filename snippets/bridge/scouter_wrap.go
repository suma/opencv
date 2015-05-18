package bridge

/*
#include "scouter_bridge.h"
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

type FrameProcessorConfig C.FrameProcessorConfig
type FrameProcessor C.FrameProcessor
type Frame C.Frame
type DetectorConfig C.DetectorConfig
type Detector C.Detector
type DetectionResult C.DetectionResult

func FrameProcessor_SetUp(fp FrameProcessor, config FrameProcessorConfig) {
	C.FrameProcessor_SetUp(C.FrameProcessor(fp), C.FrameProcessorConfig(config))
}

// must C.free(unsafe.Pointer(b)) ???
func FrameProcessor_Apply(fp FrameProcessor, buf MatVec3b,
	timestamp int64, cameraID int) (Frame, []byte) {
	var fr Frame
	b := make([]byte, 1)
	var l int

	//data := (*reflect.SliceHeader)(unsafe.Pointer(&b)).Data
	C.FrameProcessor_Apply(
		C.FrameProcessor(fp), C.MatVec3b(buf),
		C.longlong(timestamp), C.int(cameraID), C.Frame(fr),
		(*C.char)(unsafe.Pointer(&b)), (*C.int)(unsafe.Pointer(&l)))
	frByte := C.GoBytes(unsafe.Pointer(&b), C.int(l))
	return fr, frByte
}

func Detector_SetUp(detector Detector, config DetectorConfig) {
	C.Detector_SetUp(C.Detector(detector), C.DetectorConfig(config))
}

func Detector_Detect(detector Detector, frame Frame) (DetectionResult, []byte) {
	var dr DetectionResult
	b := make([]byte, 1)
	var l int

	C.Detector_Detect(C.Detector(detector), C.Frame(frame), C.DetectionResult(dr),
		(*C.char)(unsafe.Pointer(&b)), (*C.int)(unsafe.Pointer(&l)))
	drByte := C.GoBytes(unsafe.Pointer(&b), C.int(l))
	return dr, drByte
}

func ConvertToFramePointer(fr []byte) {
	var f Frame
	C.ConvertToFramePointer((*C.char)(unsafe.Pointer(&fr)), C.Frame(f))
	return f
}

func Scouter_GetEpochms() uint64 {
	return uint64(C.Scouter_GetEpochms())
}

func DetectDrawResult(frame Frame, dr DetectionResult, ms uint64) (MatVec3b, []byte) {
	var draw MatVec3b
	b := make([]byte, 1)
	var l int

	C.DetectDrawResult(
		C.Frame(frame), C.DetectionResult(dr), C.ulonglong(ms),
		C.MatVec3b(draw), (*C.char)(unsafe.Pointer(b)), (*C.int)(unsafe.Pointer(&l)))
	drwByte := C.GoBytes(unsafe.Pointer(&b), C.int(l))
	return draw, drwByte
}
