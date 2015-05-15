package bridge

/*
#include "scouter_bridge.h"
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type FrameProcessorConfig C.FrameProcessorConfig
type FrameProcessor C.FrameProcessor
type DetectorConfig C.DetectorConfig
type Detector C.Detector
type DetectionResult C.DetectionResult

func FrameProcessor_SetUp(fp FrameProcessor, config FrameProcessorConfig) {
	C.FrameProcessor_SetUp(C.FrameProcessor(fp), C.FrameProcessorConfig(config))
}

// must C.free(unsafe.Pointer(b)) ???
func FrameProcessor_Apply(fp FrameProcessor, buf MatVec3b,
	timestamp int64, cameraID int) ([]byte, bool) {
	var b []byte

	data := (*reflect.SliceHeader)(unsafe.Pointer(&b)).Data
	err := C.FrameProcessor_Apply(
		C.FrameProcessor(fp), C.MatVec3b(buf),
		C.longlong(timestamp), C.int(cameraID), (*C.char)(unsafe.Pointer(data)))
	ok := err != 0
	return b, ok
}

func Detector_SetUp(detector Detector, config DetectorConfig) {
	C.Detector_SetUp(C.Detector(detector), C.DetectorConfig(config))
}

func Detector_Detect(detector Detector, frame []byte) ([]byte, bool) {
	var dr []byte

	frameData := (*reflect.SliceHeader)(unsafe.Pointer(&frame)).Data
	drData := (*reflect.SliceHeader)(unsafe.Pointer(&dr)).Data
	err := C.Detector_Detect(C.Detector(detector), (*C.char)(unsafe.Pointer(frameData)),
		(*C.char)(unsafe.Pointer(drData)))
	ok := err != 0
	return dr, ok
}

func Scouter_GetEpochms() uint64 {
	return uint64(C.Scouter_GetEpochms())
}

func DetectDrawResult(frame []byte, dr []byte, ms uint64) ([]byte, bool) {
	var resultFrame []byte

	frameData := (*reflect.SliceHeader)(unsafe.Pointer(&frame)).Data
	drData := (*reflect.SliceHeader)(unsafe.Pointer(&dr)).Data
	rfData := (*reflect.SliceHeader)(unsafe.Pointer(&resultFrame)).Data
	err := C.DetectDrawResult(
		(*C.char)(unsafe.Pointer(frameData)), (*C.char)(unsafe.Pointer(drData)), C.ulonglong(ms),
		(*C.char)(unsafe.Pointer(rfData)))
	ok := err != 0
	return resultFrame, ok
}
