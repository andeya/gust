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

// UnsafeConvert convert a value to another type.
func UnsafeConvert[T any, U any](t T) U {
	return *(*U)(unsafe.Pointer(&t))
}

// ToAnySlice convert []T to []any.
func ToAnySlice[T any](a []T) []any {
	if a == nil {
		return nil
	}
	r := make([]any, len(a))
	for k, v := range a {
		r[k] = v
	}
	return r
}

// ToAnyMap convert map[K]V to map[K]any.
func ToAnyMap[K comparable, V any](a map[K]V) map[K]any {
	if a == nil {
		return nil
	}
	r := make(map[K]any, len(a))
	for k, v := range a {
		r[k] = v
	}
	return r
}

// UnsafeAssertSlice convert []any to []T.
func UnsafeAssertSlice[T any](a []any) []T {
	if a == nil {
		return nil
	}
	r := make([]T, len(a))
	for k, v := range a {
		r[k] = v.(T)
	}
	return r
}

// UnsafeAssertMap convert map[K]any to map[K]V.
func UnsafeAssertMap[K comparable, V any](a map[K]any) map[K]V {
	if a == nil {
		return nil
	}
	r := make(map[K]V, len(a))
	for k, v := range a {
		r[k] = v.(V)
	}
	return r
}
