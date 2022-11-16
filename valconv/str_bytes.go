package valconv

import (
	"unsafe"
)

// BytesToString convert []byte type to ~string type.
func BytesToString[STRING ~string](b []byte) STRING {
	return *(*STRING)(unsafe.Pointer(&b))
}

type ReadonlyBytes = []byte

// StringToReadonlyBytes convert ~string to unsafe read-only []byte.
// NOTE:
//
//	panic if modify the member value of the ReadonlyBytes.
func StringToReadonlyBytes[STRING ~string](s STRING) ReadonlyBytes {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{*(*string)(unsafe.Pointer(&s)), len(s)},
	))
}
