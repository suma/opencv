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
	imgPointers := []C.MatVec3b{}
	for _, img := range mats {
		imgPointers = append(imgPointers, img.p)
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
