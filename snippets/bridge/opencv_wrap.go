package bridge

/*
#cgo linux pkg-config: opencv
#cgo darwin pkg-config: opencv
#include "opencv_bridge.h"
*/
import "C"

type VideoCapture C.VideoCapture
type MatVec3b C.MatVec3b

func VideoCapture_Open(uri string, vcap VideoCapture) bool {
	ok := C.VideoCapture_Open(C.CString(uri), C.VideoCapture(vcap))
	return ok != 0
}

func VideoCapture_IsOpened(vcap VideoCapture) bool {
	isOpened := C.VideoCapture_IsOpened(C.VideoCapture(vcap))
	return isOpened != 0
}

func VideoCapture_Read(vcap VideoCapture, buf MatVec3b) bool {
	ok := C.VideoCapture_Read(C.VideoCapture(vcap), C.MatVec3b(buf))
	return ok != 0
}

func VideoCapture_Grab(vcap VideoCapture) {
	C.VideoCapture_Grab(C.VideoCapture(vcap))
}

func MatVec3b_Clone(buf MatVec3b, cloneBuf MatVec3b) {
	C.MatVec3b_Clone(C.MatVec3b(buf), C.MatVec3b(cloneBuf))
}

func MatVec3b_Empty(buf MatVec3b) bool {
	isEmpty := C.MatVec3b_Empty(C.MatVec3b(buf))
	return isEmpty != 0
}
