package bridge

/*
#cgo linux pkg-config: opencv
#cgo darwin pkg-config: opencv
#include "opencv_bridge.h"
*/
import "C"

type VideoCapture C.VideoCapture
type MatVec3b C.MatVec3b

func VideoCapture_Open(uri string) VideoCapture {
	return VideoCapture(C.VideoCapture_Open(C.CString(uri)))
}

func VideoCapture_IsOpened(vcap VideoCapture) bool {
	isOpened := C.VideoCapture_IsOpened(C.VideoCapture(vcap))
	return isOpened != 0
}

func VideoCapture_Read(vcap VideoCapture, buf MatVec3b) {
	C.VideoCapture_Read(C.VideoCapture(vcap), C.MatVec3b(buf))
}

func VideoCapture_Grab(vcap VideoCapture) {
	C.VideoCapture_Grab(C.VideoCapture(vcap))
}

func MatVec3b_Clone(buf MatVec3b) MatVec3b {
	return MatVec3b(C.MatVec3b_Clone(C.MatVec3b(buf)))
}

func MatVec3b_Empty(buf MatVec3b) bool {
	isEmpty := C.MatVec3b_Empty(C.MatVec3b(buf))
	return isEmpty != 0
}
