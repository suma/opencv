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
func (t *Tracker) Push(frames []ScouterFrame, mvRegions []MVCandidate) {
	t.mu.Lock()
	defer t.mu.Unlock()

	fLength := len(frames)
	fs := []C.struct_ScouterFrame2{}
	for _, v := range frames {
		sf := C.struct_ScouterFrame2{
			image:     v.Image.p,
			offset_x:  C.int(v.OffsetX),
			offset_y:  C.int(v.OffsetY),
			timestamp: C.ulonglong(v.Timestamp),
			camera_id: C.int(v.CameraID),
		}
		fs = append(fs, sf)
	}

	mvLength := len(mvRegions)
	mvos := []C.MVCandidate{}
	for _, r := range mvRegions {
		mvos = append(mvos, r.p)
	}

	C.Tracker_Push(t.p, (*C.struct_ScouterFrame2)(&fs[0]), C.int(fLength),
		(*C.MVCandidate)(&mvos[0]), C.int(mvLength))
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
