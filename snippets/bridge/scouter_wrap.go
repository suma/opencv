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
type DetectorConfig C.DetectorConfig
type Detector C.Detector
type DetectionResult C.DetectionResult
type RecognizeConfig C.RecognizeConfig
type ImageTaggerCaffes C.ImageTaggerCaffes
type IntegratorConfig C.IntegratorConfig
type Integrator C.Integrator
type TrackingResult C.TrackingResult

func FrameProcessor_SetUp(fp FrameProcessor, config FrameProcessorConfig) {
	C.FrameProcessor_SetUp(C.FrameProcessor(fp), C.FrameProcessorConfig(config))
}

func FrameProcessor_Apply(fp FrameProcessor, buf MatVec3b,
	timestamp int64, cameraID int) (Frame, []byte) {
	var fr Frame
	b := make([]byte, 1)
	var l int

	C.FrameProcessor_Apply(
		C.FrameProcessor(fp), C.MatVec3b(buf),
		C.longlong(timestamp), C.int(cameraID), C.Frame(fr),
		(**C.char)(unsafe.Pointer(&b)), (*C.int)(unsafe.Pointer(&l)))
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
		(**C.char)(unsafe.Pointer(&b)), (*C.int)(unsafe.Pointer(&l)))
	drByte := C.GoBytes(unsafe.Pointer(&b), C.int(l))
	return dr, drByte
}

func Scouter_GetEpochms() uint64 {
	return uint64(C.Scouter_GetEpochms())
}

func DetectDrawResult(frame Frame, dr DetectionResult, ms uint64) []byte {
	b := make([]byte, 1)
	var l int

	C.DetectDrawResult(
		C.Frame(frame), C.DetectionResult(dr), C.ulonglong(ms),
		(**C.char)(unsafe.Pointer(&b)), (*C.int)(unsafe.Pointer(&l)))
	drwByte := C.GoBytes(unsafe.Pointer(&b), C.int(l))
	return drwByte
}

func ConvertToFramePointer(fr []byte) Frame {
	var f Frame
	C.ConvertToFramePointer((*C.char)(unsafe.Pointer(&fr)), C.Frame(f))
	return f
}

func ImageTaggerCaffe_SetUp(taggers ImageTaggerCaffes, config RecognizeConfig) {
	C.ImageTaggerCaffe_SetUp(C.ImageTaggerCaffes(taggers), C.RecognizeConfig(config))
}

func ImageTaggerCaffe_PredictTagsBatch(taggers ImageTaggerCaffes,
	frame Frame, dr DetectionResult) (DetectionResult, []byte) {
	var resultDr DetectionResult
	b := make([]byte, 1)
	var l int

	C.ImageTaggerCaffe_PredictTagsBatch(C.ImageTaggerCaffes(taggers), C.Frame(frame),
		C.DetectionResult(dr), C.DetectionResult(resultDr),
		(**C.char)(unsafe.Pointer(&b)), (*C.int)(unsafe.Pointer(&l)))
	retByte := C.GoBytes(unsafe.Pointer(&resultDr), C.int(l))
	return resultDr, retByte
}

func RecognizeDrawResult(frame Frame, dr DetectionResult) []byte {
	b := make([]byte, 1)
	var l int

	C.RecognizeDrawResult(C.Frame(frame), C.DetectionResult(dr),
		(**C.char)(unsafe.Pointer(&b)), (*C.int)(unsafe.Pointer(&l)))
	drwByte := C.GoBytes(unsafe.Pointer(&b), C.int(l))
	return drwByte
}

func ConvertToDetectionResultPointer(drByte []byte) DetectionResult {
	var dr DetectionResult
	C.ConvertToDetectionResultPointer((*C.char)(unsafe.Pointer(&drByte)), C.DetectionResult(dr))
	return dr
}

func IntegratorSetUp(integrator Integrator, config IntegratorConfig) {
	C.IntegratorSetUp(C.Integrator(integrator), C.IntegratorConfig(config))
}

func Integrator_Push(integrator Integrator, frame Frame, dr DetectionResult) {
	C.Integrator_Push(C.Integrator(integrator), C.Frame(frame), C.DetectionResult(dr))
}

func Integrator_TrackerReady(integrator Integrator) bool {
	i := C.Integrator_TrackerReady(C.Integrator(integrator))
	return i != 0
}

func Integrator_Track(integrator Integrator) (TrackingResult, []byte) {
	var tr TrackingResult
	b := make([]byte, 1)
	var l int

	C.Integrator_Track(C.Integrator(integrator), C.TrackingResult(tr),
		(**C.char)(unsafe.Pointer(&b)), (*C.int)(unsafe.Pointer(&l)))
	trByte := C.GoBytes(unsafe.Pointer(&tr), C.int(l))
	return tr, trByte
}
