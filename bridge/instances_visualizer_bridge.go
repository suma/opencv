package bridge

/*
#cgo pkg-config: scouter-core
#include <stdlib.h>
#include "instances_visualizer_bridge.h"
*/
import "C"
import (
	"sync"
	"unsafe"
)

// InstancesVisualizer is a bind of `scouter::InstancesVisualizer`.
type InstancesVisualizer struct {
	mu sync.RWMutex
	p  C.InstancesVisualizer
}

// NewInstancesVisualizer create InstancesVisualizer. InstanceManager is
// necessary to create a visualizer, but not use the manager in drawing.
func NewInstancesVisualizer(im *InstanceManager, config string) InstancesVisualizer {
	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	return InstancesVisualizer{p: C.InstancesVisualizer_New(im.p, cConfig)}
}

// Delete object.
func (v *InstancesVisualizer) Delete() {
	C.InstancesVisualizer_Delete(v.p)
	v.p = nil
}

// UpdateCameraParameter updates camera parameter with camera ID.
func (v *InstancesVisualizer) UpdateCameraParameter(cameraID int, config string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	C.InstancesVisualizer_UpdateCameraParam(v.p, C.int(cameraID), cConfig)
}

// DrawWithStates draws image with instance states information
func (v *InstancesVisualizer) DrawWithStates() MatVec3b {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return MatVec3b{p: C.InstancesVisualizer_Draw(v.p)}
}
