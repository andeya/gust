// Package gust provides Rust-inspired error handling, optional values, and iteration utilities for Go.
// This file contains the Errable type for error handling.
package gust

import (
	"fmt"
	"reflect"
)

// Errable is the type that indicates whether there is an error.
type Errable[E any] struct {
	errVal *E
}

// NonErrable returns no error object.
//
//go:inline
func NonErrable[E any]() Errable[E] {
	return Errable[E]{}
}

// ToErrable converts an error value (E) to `Errable[T]`.
func ToErrable[E any](errVal E) Errable[E] {
	switch t := any(errVal).(type) {
	case nil:
		return Errable[E]{}
	case error:
		if t == nil {
			return Errable[E]{}
		}
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8, float32, float64, complex64, complex128, string, bool:
	case *int:
		if t == (*int)(nil) {
			return Errable[E]{}
		}
	case *int64:
		if t == (*int64)(nil) {
			return Errable[E]{}
		}
	case *int32:
		if t == (*int32)(nil) {
			return Errable[E]{}
		}
	case *int16:
		if t == (*int16)(nil) {
			return Errable[E]{}
		}
	case *int8:
		if t == (*int8)(nil) {
			return Errable[E]{}
		}
	case *uint:
		if t == (*uint)(nil) {
			return Errable[E]{}
		}
	case *uint64:
		if t == (*uint64)(nil) {
			return Errable[E]{}
		}
	case *uint32:
		if t == (*uint32)(nil) {
			return Errable[E]{}
		}
	case *uint16:
		if t == (*uint16)(nil) {
			return Errable[E]{}
		}
	case *uint8:
		if t == (*uint8)(nil) {
			return Errable[E]{}
		}
	case *float32:
		if t == (*float32)(nil) {
			return Errable[E]{}
		}
	case *float64:
		if t == (*float64)(nil) {
			return Errable[E]{}
		}
	case *complex64:
		if t == (*complex64)(nil) {
			return Errable[E]{}
		}
	case *complex128:
		if t == (*complex128)(nil) {
			return Errable[E]{}
		}
	case *string:
		if t == (*string)(nil) {
			return Errable[E]{}
		}
	case *bool:
		if t == (*bool)(nil) {
			return Errable[E]{}
		}
	default:
		v := reflect.ValueOf(errVal)
		if v.Kind() == reflect.Ptr && v.IsNil() {
			return Errable[E]{}
		}
	}
	return Errable[E]{errVal: &errVal}
}

// FmtErrable wraps an Errable with a formatted error.
//
//go:inline
func FmtErrable(format string, args ...any) Errable[error] {
	return ToErrable(fmt.Errorf(format, args...))
}

// TryPanic panics if the errVal is not nil.
//
//go:inline
func TryPanic[E any](errVal E) {
	ToErrable(errVal).TryPanic()
}

// TryThrow panic returns E (panicValue[*E]) if the errVal is not nil.
// NOTE:
//
//	If there is an E, that panic should be caught with CatchErrable[E] or CatchEnumResult[U, E].
//
//go:inline
func TryThrow[E any](errVal E) {
	ToErrable(errVal).TryThrow()
}

// IsErr returns true if the Errable contains an error.
//
//go:inline
func (e Errable[E]) IsErr() bool {
	return e.errVal != nil
}

// IsOk returns true if the Errable does not contain an error.
//
//go:inline
func (e Errable[E]) IsOk() bool {
	return e.errVal == nil
}

// ToError converts the Errable to a standard Go error.
// Returns nil if IsOk() is true.
func (e Errable[E]) ToError() error {
	if e.IsOk() {
		return nil
	}
	return toError(e.UnwrapErr())
}

// UnwrapErr returns the contained error value.
// Panics if IsOk() is true.
//
//go:inline
func (e Errable[E]) UnwrapErr() E {
	return *e.errVal
}

// UnwrapErrOr returns the contained error value or a provided default.
//
//go:inline
func (e Errable[E]) UnwrapErrOr(def E) E {
	if e.IsErr() {
		return e.UnwrapErr()
	}
	return def
}

// EnumResult converts from Errable[E] to EnumResult[Void, E].
func (e Errable[E]) EnumResult() EnumResult[Void, E] {
	if e.IsErr() {
		return EnumErr[Void, E](e.UnwrapErr())
	}
	return EnumOk[Void, E](nil)
}

// Result converts from Errable[E] to Result[Void].
func (e Errable[E]) Result() Result[Void] {
	if e.IsErr() {
		return Err[Void](e.UnwrapErr())
	}
	return Ok[Void](nil)
}

// Option converts from Errable[E] to Option[E].
func (e Errable[E]) Option() Option[E] {
	if e.IsErr() {
		return Some[E](e.UnwrapErr())
	}
	return None[E]()
}

// CtrlFlow returns the `CtrlFlow[E, Void]`.
func (e Errable[E]) CtrlFlow() CtrlFlow[E, Void] {
	if e.IsErr() {
		return Break[E, Void](e.UnwrapErr())
	}
	return Continue[E, Void](nil)
}

// InspectErr calls the provided closure with a reference to the contained error (if error).
func (e Errable[E]) InspectErr(f func(err E)) Errable[E] {
	if e.IsErr() {
		f(e.UnwrapErr())
	}
	return e
}

// Inspect calls the provided closure if the Errable is Ok.
func (e Errable[E]) Inspect(f func()) Errable[E] {
	if e.IsOk() {
		f()
	}
	return e
}

// TryPanic panics if the errVal is not nil.
func (e Errable[E]) TryPanic() {
	if e.IsErr() {
		panic(e.UnwrapErr())
	}
}

// TryThrow panic returns E (panicValue[*E]) if the errVal is not nil.
// NOTE:
//
//	If there is an E, that panic should be caught with CatchErrable[E] or CatchEnumResult[U, E].
func (e Errable[E]) TryThrow() {
	if e.errVal != nil {
		panic(panicValue[E]{value: e.errVal})
	}
}

// CatchErrable catches panic caused by Errable[E].TryThrow() and sets E to *Errable[E]
// Example:
//
//	```go
//	func example() (errable Errable[string]) {
//		defer errable.Catch()
//		ToErrable("panic error").TryThrow()
//		return ToErrable("return error")
//	}
//	```
func (e *Errable[E]) Catch() {
	switch p := recover().(type) {
	case nil:
	case panicValue[E]:
		if e == nil {
			panic(p.ValueOrDefault())
		}
		e.errVal = p.value
	default:
		panic(p)
	}
}

// CatchErrable catches panic caused by `Errable[E].TryThrow()` and sets E to `*Errable[E]`
// Example:
//
//	```go
//	func example() (errable Errable[string]) {
//		defer CatchErrable[E](&errable)
//		ToErrable("panic error").TryThrow()
//		return ToErrable("return error")
//	}
//	```
func CatchErrable[E any](errable *Errable[E]) {
	switch p := recover().(type) {
	case nil:
	case panicValue[E]:
		if errable == nil {
			panic(p.ValueOrDefault())
		}
		errable.errVal = p.value
	default:
		panic(p)
	}
}
