package bridge

/*
#cgo pkg-config: scouter-core
#include <stdlib.h>
#include "tracker_bridge.h"
*/
import "C"
import (
	"unsafe"
)

type Tracker struct {
	p C.Tracker
}

func NewTracker(config string) Tracker {
	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	return Tracker{p: C.Tracker_New(cConfig)}
}

func (t *Tracker) Delete() {
	C.Tracker_Delete(t.p)
	t.p = nil
}

type TrackingResult struct {
	p C.TrackingResult
}

func (t *Tracker) Push(frames map[int]MatVec3b, mvRegions []MVCandidate, timestamp uint64) {
	length := len(frames)
	framesPtr := make([]C.struct_MatWithCameraID, length)
	for k, v := range frames {
		matWithID := C.struct_MatWithCameraID{
			cameraID: C.int(k),
			mat:      v.p,
		}
		framesPtr = append(framesPtr, matWithID)
	}

	mvRegionsLen := len(mvRegions)
	mvRegionsPtr := make([]C.MVCandidate, mvRegionsLen)
	for _, r := range mvRegions {
		mvRegionsPtr = append(mvRegionsPtr, r.p)
	}
	mvCandidates := C.InvertMVCandidates((*C.MVCandidate)(&mvRegionsPtr[0]), C.int(mvRegionsLen))

	C.Tracker_Push(t.p, (*C.MatWithCameraID)(&framesPtr[0]), C.int(length),
		mvCandidates, C.ulonglong(timestamp))
}

func (t *Tracker) Track(timestamp uint64) TrackingResult {
	tr := C.Tracker_Track(t.p, C.ulonglong(timestamp))
	return TrackingResult{p: tr}
}

func (t *Tracker) Ready() bool {
	ready := C.Tracker_Ready(t.p)
	if int(ready) == 1 {
		return true
	}
	return false
}