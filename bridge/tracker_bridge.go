package bridge

/*
#cgo pkg-config: scouter-core
#include <stdlib.h>
#include "tracker_bridge.h"
*/
import "C"
import (
	"reflect"
	"sync"
	"unsafe"
)

// Tracker is a bind of `scouter::Tracker`.
type Tracker struct {
	mu sync.RWMutex
	p  C.Tracker
}

// NewTracker returns a new tracker.
func NewTracker(config string) Tracker {
	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	return Tracker{p: C.Tracker_New(cConfig)}
}

// Delete object.
func (t *Tracker) Delete() {
	C.Tracker_Delete(t.p)
	t.p = nil
}

// Trackee is a utility structure, as scouter-core Tracke structure.
type Trackee struct {
	ColorID      uint64
	MVCandidate  MVCandidate
	Interpolated bool
	Timestamp    uint64 // should placed in TrackingResult
}

// Push frame and regions with the tracker.
func (t *Tracker) Push(frames map[int]MatVec3b, mvRegions []MVCandidate,
	timestamp uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	length := len(frames)
	framesPtr := []C.struct_MatWithCameraID{}
	for k, v := range frames {
		matWithID := C.struct_MatWithCameraID{
			cameraID: C.int(k),
			mat:      v.p,
		}
		framesPtr = append(framesPtr, matWithID)
	}

	mvRegionsLen := len(mvRegions)
	mvRegionsPtr := []C.MVCandidate{}
	for _, r := range mvRegions {
		mvRegionsPtr = append(mvRegionsPtr, r.p)
	}
	mvCandidates := C.InvertMVCandidates(
		(*C.MVCandidate)(&mvRegionsPtr[0]), C.int(mvRegionsLen))

	C.Tracker_Push(t.p, (*C.MatWithCameraID)(&framesPtr[0]), C.int(length),
		mvCandidates, C.ulonglong(timestamp))
}

// Track returns Trackee array cached in the tracker.
func (t *Tracker) Track() []Trackee {
	t.mu.Lock()
	defer t.mu.Unlock()

	tr := C.Tracker_Track(t.p)
	defer C.TrackingResult_Delete(tr)

	cArray := tr.trackees
	length := int(tr.length)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cArray)),
		Len:  length,
		Cap:  length,
	}
	goSlice := *(*[]C.Trackee)(unsafe.Pointer(&hdr))
	trs := make([]Trackee, length)
	for i, t := range goSlice {
		trs[i] = Trackee{
			ColorID:      uint64(t.colorID),
			MVCandidate:  MVCandidate{p: t.mvCandidate},
			Interpolated: int(t.interpolated) != 0,
			Timestamp:    uint64(tr.timestamp),
		}
	}
	return trs
}

// Ready return the tracker could start tracking by acceptance value of tracking
// configuration.
func (t *Tracker) Ready() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	ready := C.Tracker_Ready(t.p)
	if int(ready) == 1 {
		return true
	}
	return false
}
