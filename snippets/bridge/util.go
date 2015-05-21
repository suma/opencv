package bridge

/*
#cgo LDFLAGS: -lstdc++ -lmsgpack
#include "util.h"
*/
import "C"
import (
	"unsafe"
)

func toByteArray(b []byte) C.struct_ByteArray {
	return C.struct_ByteArray{
		data:   (*C.char)(unsafe.Pointer(&b[0])),
		length: C.int(len(b)),
	}
}

func ToGoBytes(b C.struct_ByteArray) []byte {
	return C.GoBytes(unsafe.Pointer(b.data), b.length)
}
