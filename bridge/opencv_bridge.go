package bridge

/*
#cgo linux pkg-config: opencv
#cgo darwin pkg-config: opencv
#include <stdlib.h>
#include "util.h"
#include "opencv_bridge.h"
*/
import "C"
import (
	"sync"
	"unsafe"
)

type MatVec3b struct {
	p C.MatVec3b
}

func NewMatVec3b() MatVec3b {
	return MatVec3b{p: C.MatVec3b_New()}
}

func (m *MatVec3b) ToJpegData(quality int) []byte {
	b := C.MatVec3b_ToJpegData(m.p, C.int(quality))
	defer C.ByteArray_Release(b)
	return ToGoBytes(b)
}

func (m *MatVec3b) Serialize() []byte {
	b := C.MatVec3b_Serialize(m.p)
	defer C.ByteArray_Release(b)
	return ToGoBytes(b)
}

func convertMatVec3bsToPointer(mats []MatVec3b) []C.MatVec3b {
	imgPointers := make([]C.MatVec3b, len(mats))
	for i, img := range mats {
		imgPointers[i] = img.p
	}
	return imgPointers
}

func DeserializeMatVec3b(m []byte) MatVec3b {
	b := toByteArray(m)
	return MatVec3b{p: C.MatVec3b_Deserialize(b)}
}

func (m *MatVec3b) Delete() {
	C.MatVec3b_Delete(m.p)
	m.p = nil
}

func (m *MatVec3b) CopyTo(dst MatVec3b) {
	C.MatVec3b_CopyTo(m.p, dst.p)
}

func (m *MatVec3b) Empty() bool {
	isEmpty := C.MatVec3b_Empty(m.p)
	return isEmpty != 0
}

type VideoCapture struct {
	p C.VideoCapture
}

func NewVideoCapture() VideoCapture {
	return VideoCapture{p: C.VideoCapture_New()}
}

func (v *VideoCapture) Delete() {
	C.VideoCapture_Delete(v.p)
	v.p = nil
}

func (v *VideoCapture) Open(uri string) bool {
	c_uri := C.CString(uri)
	defer C.free(unsafe.Pointer(c_uri))
	return C.VideoCapture_Open(v.p, c_uri) != 0
}

func (v *VideoCapture) OpenDevice(device int) bool {
	return C.VideoCapture_OpenDevice(v.p, C.int(device)) != 0
}

func (v *VideoCapture) Release() {
	C.VideoCapture_Release(v.p)
}

func (v *VideoCapture) Set(prop int, param int) {
	C.VideoCapture_Set(v.p, C.int(prop), C.int(param))
}

func (v *VideoCapture) IsOpened() bool {
	isOpened := C.VideoCapture_IsOpened(v.p)
	return isOpened != 0
}

func (v *VideoCapture) Read(m MatVec3b) bool {
	return C.VideoCapture_Read(v.p, m.p) != 0
}

func (v *VideoCapture) Grab(skip int) {
	C.VideoCapture_Grab(v.p, C.int(skip))
}

type VideoWriter struct {
	mu sync.RWMutex
	p  C.VideoWriter
}

func NewVideoWriter() VideoWriter {
	return VideoWriter{p: C.VideoWriter_New()}
}

func (vw *VideoWriter) Delete() {
	C.VideoWriter_Delete(vw.p)
	vw.p = nil
}

func (vw *VideoWriter) Open(name string, fps float64, width int, height int) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.VideoWriter_Open(vw.p, cName, C.double(fps), C.int(width), C.int(height))
}

func (vw *VideoWriter) IsOpened() bool {
	isOpend := C.VideoWriter_IsOpened(vw.p)
	return isOpend != 0
}

func (vw *VideoWriter) Write(img MatVec3b) {
	vw.mu.Lock()
	defer vw.mu.Unlock()
	C.VideoWriter_Write(vw.p, img.p)
}
