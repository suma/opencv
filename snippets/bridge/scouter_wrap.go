package bridge

/*
#cgo darwin CXXFLAGS: -I/System/Library/Frameworks/Accelerate.framework/Versions/Current/Frameworks/vecLib.framework/Headers/ -DCPU_ONLY
#cgo LDFLAGS: -ljsonconfig
#cgo pkg-config: scouter-core
#cgo pkg-config: pficv
#cgo pkg-config: pficommon
#include <stdlib.h>
#include "scouter_bridge.h"
#include "util.h"
*/
import "C"
import (
	"unsafe"
)

type Frame struct {
	p C.Frame
}

type DetectionResult struct {
	p C.DetectionResult
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

type TrackingResult struct {
	p C.TrackingResult
}

type Integrator struct {
	p C.Integrator
}

type InstanceManager struct {
	p C.InstanceManager
}

type InstanceStates struct {
	p C.InstanceStates
}

type Visualizer struct {
	p C.Visualizer
}

func (f Frame) Serialize() []byte {
	b := C.Frame_Serialize(f.p)
	defer C.ByteArray_Release(b)
	return ToGoBytes(b)
}

func DeserializeFrame(f []byte) Frame {
	b := toByteArray(f)
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
	return DetectionResult{p: C.DetectionResult_Deserialize(b)}
}

func (d DetectionResult) Delete() {
	C.DetectionResult_Delete(d.p)
	d.p = nil
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

func (fp *FrameProcessor) Apply(buf MatVec3b, timestamp int64,
	cameraID int) Frame {
	return Frame{p: C.FrameProcessor_Apply(fp.p, buf.p, C.longlong(timestamp), C.int(cameraID))}
}

func NewDetector(config string) Detector {
	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	return Detector{p: C.Detector_New(cConfig)}
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

func NewImageTaggerCaffe(configTaggers string) ImageTaggerCaffe {
	cConfig := C.CString(configTaggers)
	defer C.free(unsafe.Pointer(cConfig))
	return ImageTaggerCaffe{
		p: C.ImageTaggerCaffe_New(cConfig),
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

func RecognizeDrawResult(f Frame, dr DetectionResult) map[string]MatVec3b {
	result := C.RecognizeDrawResult(f.p, dr.p)
	defer C.Taggers_Delete(result)
	l := int(result.length)
	keys := make([](*C.char), l)
	drawResults := make([]C.MatVec3b, l)
	C.ResolveDrawResult(result, (**C.char)(&keys[0]), (*C.MatVec3b)(&drawResults[0]))

	resultMap := make(map[string]MatVec3b, l)
	for i := 0; i < l; i++ {
		resultMap[C.GoString(keys[i])] = MatVec3b{p: drawResults[i]}
	}
	return resultMap
}

func NewIntegrator(config string) Integrator {
	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	return Integrator{p: C.Integrator_New(cConfig)}
}

func (itr *Integrator) Delete() {
	C.Integrator_Delete(itr.p)
	itr.p = nil
}

func (itr *Integrator) Integrator_Push(fs []Frame, drs []DetectionResult) {
	size := len(fs)
	frPointers := []C.Frame{}
	drPointers := []C.DetectionResult{}
	for i := 0; i < size; i++ {
		frPointers = append(frPointers, fs[i].p)
		drPointers = append(drPointers, drs[i].p)
	}
	C.Integrator_Push(itr.p, (*C.Frame)(&frPointers[0]),
		(*C.DetectionResult)(&drPointers[0]), C.int(size))
}

func (itr *Integrator) Integrator_TrackerReady() bool {
	return C.Integrator_TrackerReady(itr.p) != 0
}

func (itr *Integrator) Integrator_Track() TrackingResult {
	return TrackingResult{C.Integrator_Track(itr.p)}
}

func (tr *TrackingResult) Delete() {
	C.TrackingResult_Delete(tr.p)
	tr.p = nil
}

func NewInstanceManager(config string) InstanceManager {
	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	return InstanceManager{p: C.InstanceManager_New(cConfig)}
}

func (im *InstanceManager) Delete() {
	C.InstanceManager_Delete(im.p)
	im.p = nil
}

func (im *InstanceManager) GetCurrentStates(tr TrackingResult) InstanceStates {
	return InstanceStates{
		p: C.InstanceManager_GetCurrentStates(im.p, tr.p),
	}
}

func (is *InstanceStates) Delete() {
	C.InstanceStates_Delete(is.p)
	is.p = nil
}

func (is *InstanceStates) ConvertSatesToJson(floorID int, timestamp int64) string {
	c_str := C.ConvertStatesToJson(is.p, C.int(floorID), C.longlong(timestamp))
	return C.GoStringN(c_str.str, c_str.length)
}

func NewVisualizer(config string, instanceManager InstanceManager) Visualizer {
	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	return Visualizer{p: C.Visualizer_New(cConfig, instanceManager.p)}
}

func (v *Visualizer) Delete() {
	C.Visualizer_Delete(v.p)
	v.p = nil
}

func (v *Visualizer) PlotTrajectories() []MatVec3b {
	plots := C.Visualizer_PlotTrajectories(v.p)
	defer C.PlotTrajectories_Delete(plots)
	l := int(plots.length)
	plotsTraj := make([]C.MatVec3b, l)
	C.ResolvePlotTrajectories(plots, (*C.MatVec3b)(&plotsTraj[0]))

	result := make([]MatVec3b, l)
	for i := 0; i < l; i++ {
		result[i] = MatVec3b{p: plotsTraj[i]}
	}
	return result
}
