// Package conv provides generic functions for type conversion and value transformation.
//
// This package offers safe and unsafe conversion utilities between compatible types,
// including byte/string conversions, reference conversions, reflection-based operations,
// and string manipulation utilities.
//
// # Examples
//
//	// Convert bytes to string (zero-copy)
//	bytes := []byte{'h', 'e', 'l', 'l', 'o'}
//	str := conv.BytesToString[string](bytes)
//	fmt.Println(str) // Output: hello
//
//	// Convert string to readonly bytes
//	readonlyBytes := conv.StringToReadonlyBytes("hello")
//	// Note: modifying readonlyBytes will panic
//
//	// Convert to snake_case
//	snake := conv.SnakeString("UserID")
//	fmt.Println(snake) // "user_id"
//
//	// Convert to CamelCase
//	camel := conv.CamelString("user_id")
//	fmt.Println(camel) // "UserId"
package conv

import (
	"encoding/binary"
	"unsafe"

	"github.com/andeya/gust/result"
)

// BytesToString convert []byte type to ~string type.
func BytesToString[STRING ~string](b []byte) STRING {
	return *(*STRING)(unsafe.Pointer(&b))
}

// SystemEndian returns the byte order of the current system.
//
//go:inline
func SystemEndian() binary.ByteOrder {
	return systemEndian
}

var systemEndian = func() binary.ByteOrder {
	var i int32 = 0x01020304
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb
	if b == 0x04 {
		return binary.LittleEndian
	}
	return binary.BigEndian
}()

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

// SafeAssert asserts any value up to (zero)T.
func SafeAssert[T any](v any) T {
	t, _ := v.(T)
	return t
}

// SafeAssertSlice convert []any to []T.
func SafeAssertSlice[T any](a []any) result.Result[[]T] {
	if a == nil {
		return result.Ok[[]T](nil)
	}
	var ok bool
	r := make([]T, len(a))
	for k, v := range a {
		r[k], ok = v.(T)
		if !ok {
			return result.FmtErr[[]T]("assert slice[%v] type failed, got %T want %T", k, v, *new(T))
		}
	}
	return result.Ok(r)
}

// SafeAssertMap convert map[K]any to map[K]V.
func SafeAssertMap[K comparable, V any](a map[K]any) result.Result[map[K]V] {
	if a == nil {
		return result.Ok[map[K]V](nil)
	}
	var ok bool
	r := make(map[K]V, len(a))
	for k, v := range a {
		r[k], ok = v.(V)
		if !ok {
			return result.FmtErr[map[K]V]("assert map[%v] type failed, got %T want %T", k, v, *new(V))
		}
	}
	return result.Ok(r)
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
