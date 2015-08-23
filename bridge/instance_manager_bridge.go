package bridge

/*
#cgo pkg-config: scouter-core
#include <stdlib.h>
#include "instance_manager_bridge.h"
*/
import "C"
import (
	"reflect"
	"sync"
	"unsafe"
)

// InstanceState is a bind ob `scouter::InstanceState`.
type InstanceState struct {
	p C.InstanceState
}

// Serialize object.
func (s InstanceState) Serialize() []byte {
	b := C.InstanceState_Serialize(s.p)
	defer C.ByteArray_Release(b)
	return ToGoBytes(b)
}

// DeserializeInstanceState deserializes object.
func DeserializeInstanceState(s []byte) InstanceState {
	b := toByteArray(s)
	return InstanceState{p: C.InstanceState_Deserialize(b)}
}

// Delete object.
func (s InstanceState) Delete() {
	C.InstanceState_Delete(s.p)
	s.p = nil
}

// ConvertInstanceStatesToJSON converts instance state object to JSON style text.
func ConvertInstanceStatesToJSON(iss []InstanceState, floorID int,
	timestamp uint64) string {

	issPtr := []C.InstanceState{}
	for _, is := range iss {
		issPtr = append(issPtr, is.p)
	}
	cIssPtr := C.InvertInstanceStates((*C.InstanceState)(&issPtr[0]),
		C.int(len(iss)))
	defer C.InstanceStates_Delete(cIssPtr)

	cJSON := C.InstanceState_ToJSON(cIssPtr, C.int(floorID),
		C.ulonglong(timestamp))
	return C.GoStringN(cJSON.str, cJSON.length)
}

// InstanceManager is a bind of `scouter::InstanceManager`.
type InstanceManager struct {
	mu sync.RWMutex
	p  C.InstanceManager
}

// NewInstanceManager return a new instance manager.
func NewInstanceManager(config string) InstanceManager {
	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	return InstanceManager{p: C.InstanceManager_New(cConfig)}
}

// Delete object.
func (m *InstanceManager) Delete() {
	C.InstanceManager_Delete(m.p)
	m.p = nil
}

// TrackAndGetStates returns current states. First, tracker returns
// `scouter::TrackingResult`. Second, instance manager returns instance states
// using the tracking result.
func (m *InstanceManager) TrackAndGetStates(tr Tracker) []InstanceState {
	m.mu.Lock()
	defer m.mu.Unlock()
	tr.mu.Lock()
	defer tr.mu.Unlock()

	iss := C.TrackAndGetStates(tr.p, m.p)
	defer C.InstanceStates_Delete(iss)

	cArray := iss.instanceStates
	length := int(iss.length)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cArray)),
		Len:  length,
		Cap:  length,
	}
	goSlice := *(*[]C.InstanceState)(unsafe.Pointer(&hdr))

	states := make([]InstanceState, length)
	for i, s := range goSlice {
		states[i] = InstanceState{p: s}
	}
	return states
}

// Update cached frames managed by instance manager.
func (m *InstanceManager) Update(frames map[int]MatVec3b, trackees []Trackee,
	timestamp uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	fLength := len(frames)
	framesPtr := []C.struct_MatWithCameraID{}
	for k, v := range frames {
		matWithID := C.struct_MatWithCameraID{
			cameraID: C.int(k),
			mat:      v.p,
		}
		framesPtr = append(framesPtr, matWithID)
	}

	tLength := len(trackees)
	trs := []C.struct_Trackee{}
	for _, t := range trackees {
		var interpo int
		if t.Interpolated {
			interpo = 1
		} else {
			interpo = 0
		}
		trackee := C.struct_Trackee{
			colorID:      C.ulonglong(t.ColorID),
			mvCandidate:  t.MVCandidate.p,
			interpolated: C.int(interpo),
		}

		trs = append(trs, trackee)
	}
	C.InstanceManager_Update(m.p, (*C.MatWithCameraID)(&framesPtr[0]),
		C.int(fLength), (*C.Trackee)(&trs[0]), C.int(tLength),
		C.ulonglong(timestamp))
}

// GetCurrentStates returns a current state managed and cached in instance
// manager.
func (m *InstanceManager) GetCurrentStates() []InstanceState {
	m.mu.RLock()
	defer m.mu.RUnlock()

	currentStates := C.InstanceManager_GetCurrentStates(m.p)
	defer C.InstanceStates_Delete(currentStates)

	cArray := currentStates.instanceStates
	length := int(currentStates.length)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cArray)),
		Len:  length,
		Cap:  length,
	}
	goSlice := *(*[]C.InstanceState)(unsafe.Pointer(&hdr))

	states := make([]InstanceState, length)
	for i, s := range goSlice {
		states[i] = InstanceState{p: s}
	}
	return states
}
