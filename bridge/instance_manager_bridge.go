package bridge

/*
#cgo pkg-config: scouter-core
#include <stdlib.h>
#include "instance_manager_bridge.h"
*/
import "C"
import (
	"unsafe"
)

type InstanceState struct {
	p C.InstanceState
}

func (s InstanceState) Serialize() []byte {
	b := C.InstanceState_Serialize(s.p)
	defer C.ByteArray_Release(b)
	return ToGoBytes(b)
}

func DeserializeInstanceState(s []byte) InstanceState {
	b := toByteArray(s)
	return InstanceState{p: C.InstanceState_Deserialize(b)}
}

type InstanceManager struct {
	p C.InstanceManager
}

func NewInstanceManager(config string) InstanceManager {
	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	return InstanceManager{p: C.InstanceManager_New(cConfig)}
}

func (m *InstanceManager) Delete() {
	C.InstanceManager_Delete(m.p)
	m.p = nil
}

func (m *InstanceManager) Updaate(tr TrackingResult) {
	C.InstanceManager_Update(m.p, tr.p)
}

func (m *InstanceManager) GetCurrentStates() []InstanceState {
	currentStates := C.InstanceManager_GetCurrentStates(m.p)
	defer C.InstanceStates_Delete(currentStates)
	l := int(currentStates.length)
	statesPtr := make([]C.InstanceState, l)
	C.ResolveInstanceStates(currentStates, (*C.InstanceState)(&statesPtr[0]))

	states := make([]InstanceState, l)
	for i := 0; i < l; i++ {
		states[i] = InstanceState{p: statesPtr[i]}
	}
	return states
}
