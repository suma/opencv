package bridge

/*
#cgo darwin CXXFLAGS: -I/System/Library/Frameworks/Accelerate.framework/Versions/Current/Frameworks/vecLib.framework/Headers/ -DCPU_ONLY
#cgo pkg-config: scouter-core
#cgo pkg-config: pficv
#cgo pkg-config: pficommon
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

type TrackingResult struct {
	p C.TrackingResult
}

type FrameProcessor struct {
	p C.FrameProcessor
}

type Detector struct {
	p C.Detector
}

type ImageTaggerCaffe struct {
	p C.ImageTaggerCaffe
}

type Taggers struct {
	p C.Taggers
}

type Integrator struct {
	p C.Integrator
}

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

func (t TrackingResult) Serialize() []byte {
	b := C.TrackingResult_Serialize(t.p)
	defer C.ByteArray_Release(b)
	return ToGoBytes(b)
}

func DeserializeTrackingResult(t []byte) TrackingResult {
	b := toByteArray(t)
	defer C.ByteArray_Release(b)
	return TrackingResult{p: C.TrackingResult_Deserialize(b)}
}

func (t TrackingResult) Delete() {
	C.TrackingResult_Delete(t.p)
	t.p = nil
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

func DetectDrawResult(f Frame, dr DetectionResult, ms int64) MatVec3b {
	return MatVec3b{p: C.DetectDrawResult(f.p, dr.p, C.longlong(ms))}
}

func ImageTaggerCaffe_New(configTaggers RecognizeConfigTaggers) ImageTaggerCaffe {
	return ImageTaggerCaffe{
		p: C.ImageTaggerCaffe_New(C.RecognizeConfigTaggers(configTaggers.p)),
	}
}

func (itc *ImageTaggerCaffe) Delete() {
	C.ImageTaggerCaffe_Delete(itc.p)
	itc.p = nil
}

func (itc *ImageTaggerCaffe) Recognize(
	f Frame, dr DetectionResult) DetectionResult {
	return DetectionResult{p: C.Recognize(itc.p, f.p, dr.p)}
}

func RecognizeDrawResult(f Frame, dr DetectionResult) Taggers {
	return Taggers{p: C.RecognizeDrawResult(f.p, dr.p)}
}

func Integrator_New(config IntegratorConfig) Integrator {
	return Integrator{p: C.Integrator_New(config.p)}
}

func (itr *Integrator) Delete() {
	C.Integrator_Delete(itr.p)
	itr.p = nil
}

func (itr *Integrator) Integrator_Push(f Frame, dr DetectionResult) {
	C.Integrator_Push(itr.p, f.p, dr.p)
}

func (itr *Integrator) Integrator_TrackerReady() bool {
	return C.Integrator_TrackerReady(itr.p) != 0
}

func (itr *Integrator) Integrator_Track() TrackingResult {
	return TrackingResult{C.Integrator_Track(itr.p)}
}
