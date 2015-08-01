package bridge

/*
#cgo pkg-config: scouter-core
#include <stdlib.h>
#include "instances_visualizer_bridge.h"
*/
import "C"
import (
	"unsafe"
)

type InstancesVisualizer struct {
	im *InstanceManager
	p  C.InstancesVisualizer
}

func NewInstancesVisualizer(im *InstanceManager, config string) InstancesVisualizer {
	cConfig := C.CString(config)
	defer C.free(unsafe.Pointer(cConfig))
	return InstancesVisualizer{
		im: im,
		p:  C.InstancesVisualizer_New(im.p, cConfig),
	}
}

func (v *InstancesVisualizer) Delete() {
	C.InstancesVisualizer_Delete(v.p)
	v.p = nil
}

func (v *InstancesVisualizer) Draw() MatVec3b {
	v.im.mu.RLock()
	v.im.mu.RUnlock()
	img := C.InstancesVisualizer_Draw(v.p)
	return MatVec3b{p: img}
}
