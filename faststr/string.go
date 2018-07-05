package faststr

import (
	"unsafe"
)

func StringFromByteSlice(slice []byte) string {
	return *(*string)(unsafe.Pointer(&slice))
}

func ByteSliceFromString(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(&str))
}
