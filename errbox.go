package gust

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/andeya/gust/errutil"
)

type (
	// ErrBox is a wrapper for any error type.
	// Use ToError() method to convert to error interface, or access ErrBox through Result.Err().
	ErrBox struct {
		inner innerErrBox
	}
	// innerErrBox is an internal type that implements the error interface.
	// It wraps a value and provides error representation.
	innerErrBox struct {
		val any
	}
)

var (
	_ fmt.Stringer   = (*ErrBox)(nil)
	_ fmt.GoStringer = (*ErrBox)(nil)
	_ error          = (*innerErrBox)(nil)
	// errorInterfaceType is cached to avoid repeated reflection calls
	errorInterfaceType = reflect.TypeOf((*error)(nil)).Elem()
)

// BoxErr wraps any error type into ErrBox.
//
//go:inline
func BoxErr(val any) *ErrBox {
	switch val := val.(type) {
	case nil:
		return nil
	case ErrBox:
		return &val
	case *ErrBox:
		return val
	case innerErrBox:
		return &ErrBox{inner: val}
	default:
		return &ErrBox{inner: innerErrBox{val: val}}
	}
}

// IsEmpty returns true if ErrBox is empty (nil receiver or nil val).
//
//go:inline
func (e *ErrBox) IsEmpty() bool {
	return e == nil || e.inner.val == nil
}

// Value returns the inner value.
//
//go:inline
func (e *ErrBox) Value() any {
	if e == nil {
		return nil
	}
	return e.inner.val
}

// ValueOrDefault returns the inner value, or nil if ErrBox is nil or val is nil.
// This is useful when you want to safely extract a typed value from ErrBox.
//
//go:inline
func (e *ErrBox) ValueOrDefault() any {
	if e == nil || e.inner.val == nil {
		return nil
	}
	return e.inner.val
}

// String returns the string representation.
// This implements the fmt.Stringer interface.
func (e *ErrBox) String() string {
	if e == nil {
		return "<nil>"
	}
	return e.inner.Error()
}

// GoString returns the Go-syntax representation.
// This implements the fmt.GoStringer interface.
func (e *ErrBox) GoString() string {
	if e == nil {
		return "(*gust.ErrBox)(nil)"
	}
	return fmt.Sprintf("&gust.ErrBox{inner: %#v}", e.inner.val)
}

// ToError converts ErrBox to error interface.
// Returns nil if the receiver is nil or val is nil.
// Returns a pointer to innerErrBox which implements the error interface.
//
// Example:
//
//	```go
//	var eb gust.ErrBox = gust.BoxErr(errors.New("test"))
//	var err error = eb.ToError() // err is *innerErrBox implementing error
//	```
//
//go:inline
func (e *ErrBox) ToError() error {
	if e == nil || e.inner.val == nil {
		return nil
	}
	// If the wrapped value is already an error, return it directly
	if err, ok := e.inner.val.(error); ok {
		return err
	}
	// Otherwise return the innerErrBox wrapper
	return &e.inner
}

// Unwrap returns the inner error.
//
//go:inline
func (e *ErrBox) Unwrap() error {
	if e == nil || e.inner.val == nil {
		return nil
	}
	return e.inner.Unwrap()
}

// Is reports whether any error in err's chain matches target.
//
//go:inline
func (e *ErrBox) Is(target error) bool {
	if e == nil || e.inner.val == nil {
		return target == nil
	}
	return e.inner.Is(target)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false.
//
//go:inline
func (e *ErrBox) As(target any) bool {
	if e == nil || e.inner.val == nil {
		return false
	}
	return e.inner.As(target)
}

// Error returns the string representation.
// This implements the error interface.
func (e innerErrBox) Error() string {
	if e.val == nil {
		return "<nil>"
	}
	switch val := e.val.(type) {
	case string:
		return val
	case error:
		return val.Error()
	default:
		return fmt.Sprintf("%v", val)
	}
}

// Unwrap returns the inner error.
func (e innerErrBox) Unwrap() error {
	switch val := e.val.(type) {
	case nil:
		return nil
	case error:
		u, ok := val.(interface {
			Unwrap() error
		})
		if ok {
			return u.Unwrap()
		}
		return val
	default:
		return nil
	}
}

// Is reports whether any error in err's chain matches target.
func (e innerErrBox) Is(target error) bool {
	b := e.Unwrap()
	if b != nil {
		return errors.Is(b, target)
	}
	// If Unwrap returns nil, check if target matches the wrapped value
	if err, ok := e.val.(error); ok {
		return errors.Is(err, target)
	}
	return false
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false.
func (e innerErrBox) As(target any) bool {
	if target == nil {
		panic("errors: target cannot be nil")
	}
	targetVal := reflect.ValueOf(target)
	if targetVal.Kind() != reflect.Ptr {
		return false
	}
	targetType := targetVal.Type().Elem()

	// Check if target type implements error interface
	// errors.As requires *target to be interface or implement error
	isErrorTarget := targetType.Implements(errorInterfaceType)

	// If e.val is an error, handle error matching
	if err, ok := e.val.(error); ok {
		// First check if the error itself matches the target type directly
		errType := reflect.TypeOf(err)
		if errType.AssignableTo(targetType) {
			targetVal.Elem().Set(reflect.ValueOf(err))
			return true
		}
		// Only use errors.As if target implements error interface
		// errors.As requires target to be interface or implement error
		if isErrorTarget {
			return errors.As(err, target)
		}
		return false
	}

	// For non-error types, As only works if target is error interface
	// and we can convert the value to error (which we can't for non-error types)
	return false
}

// newPanicError creates a new PanicError with the given error and stack trace.
// This is a helper function that wraps the errutil.NewPanicError with ErrBox handling.
func newPanicError(err any, stack errutil.StackTrace) *errutil.PanicError {
	var wrappedErr error
	switch e := err.(type) {
	case nil:
		wrappedErr = nil
	case error:
		wrappedErr = e
	case *ErrBox:
		if e != nil {
			wrappedErr = e.ToError()
		}
	case ErrBox:
		wrappedErr = e.ToError()
	default:
		wrappedErr = BoxErr(err).ToError()
	}
	return errutil.NewPanicError(wrappedErr, stack)
}
