package helper

import (
	"unsafe"
)

func BytesToString(v []byte) string {
	return *(*string)(unsafe.Pointer(&v))
}

func StringToBytes(v string) []byte {
	return unsafe.Slice(unsafe.StringData(v), len(v))
}
