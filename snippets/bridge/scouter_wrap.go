package bridge

/*
#cgo linux pkg-config: scouter-core
#cgo darwin pkg-config: scouter-core
#cgo linux pkg-config: pficv
#cgo darwin pkg-config: pficv
#include "scouter_bridge.h"
#include "util.h"
*/
import "C"

type Frame struct {
	p C.Frame
}

type DetectionResult struct {
	p C.DetectionResult
}

type FrameProcessorConfig struct {
	p C.FrameProcessorConfig
}

type FrameProcessor struct {
	p C.FrameProcessor
}

type DetectorConfig struct {
	p C.DetectorConfig
}

type Detector struct {
	p C.Detector
}

type RecognizeConfig C.RecognizeConfig
type ImageTaggerCaffes C.ImageTaggerCaffes
type IntegratorConfig C.IntegratorConfig
type Integrator C.Integrator
type TrackingResult C.TrackingResult

func (f Frame) Serialize() []byte {
	b := C.Frame_Serialize(f.p)
	defer C.ByteArray_Release(b)
	return ToGoBytes(b)
}

func DeserializeFrame(f []byte) Frame {
	b := toByteArray(f)
	defer C.ByteArray_Release(b)
	return Frame{p: C.Freme_Deserialize(b)}
}

func (f Frame) Delete() {
	C.Frame_Delete(f.p)
	f.p = nil
}

func (d DetectionResult) Serialize() []byte {
	b := C.DetectionResult_Serialize(d.p)
	defer C.ByteArray_Release(b)
	return ToGoBytes(b)
}

func DeserializeDetectionResult(d []byte) DetectionResult {
	b := toByteArray(d)
	defer C.ByteArray_Release(b)
	return DetectionResult{p: C.DetectionResult_Deserialize(b)}
}

func (d DetectionResult) Delete() {
	C.DetectionResult_Delete(d.p)
	d.p = nil
}

func NewFrameProcessor(config FrameProcessorConfig) FrameProcessor {
	return FrameProcessor{p: C.FrameProcessor_New(config.p)}
}

func (fp *FrameProcessor) Delete() {
	C.FrameProcessor_Delete(fp.p)
	fp.p = nil
}

func (fp *FrameProcessor) Apply(buf MatVec3b, timestamp int64,
	cameraID int) Frame {
	return Frame{p: C.FrameProcessor_Apply(fp.p, buf.p, C.longlong(timestamp), C.int(cameraID))}
}

func NewDetector(config DetectorConfig) Detector {
	return Detector{p: C.Detector_New(config.p)}
}

func (d *Detector) Delete() {
	C.Detector_Delete(d.p)
	d.p = nil
}

func (d *Detector) Detect(f Frame) DetectionResult {
	return DetectionResult{p: C.Detector_Detect(d.p, f.p)}
}

func DetectDrawResult(f Frame, dr DetectionResult, ms uint64) MatVec3b {
	return MatVec3b{p: C.DetectDrawResult(f.p, dr.p, C.longlong(ms))}
}

func Scouter_GetEpochms() uint64 {
	return uint64(C.Scouter_GetEpochms())
}

func ImageTaggerCaffe_SetUp(taggers ImageTaggerCaffes, config RecognizeConfig) {
	C.ImageTaggerCaffe_SetUp(C.ImageTaggerCaffes(taggers), C.RecognizeConfig(config))
}

func ImageTaggerCaffe_PredictTagsBatch(taggers ImageTaggerCaffes,
	frame Frame, dr DetectionResult) (DetectionResult, []byte) {
	return DetectionResult{}, []byte{}
}

func RecognizeDrawResult(frame Frame, dr DetectionResult) []byte {
	return []byte{}
}

func IntegratorSetUp(integrator Integrator, config IntegratorConfig) {
}

func Integrator_Push(integrator Integrator, frame Frame, dr DetectionResult) {
}

func Integrator_TrackerReady(integrator Integrator) bool {
	return false
}

func Integrator_Track(integrator Integrator) (TrackingResult, []byte) {
	return nil, []byte{}
}
