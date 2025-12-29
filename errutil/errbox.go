// Package errutil provides utilities for error handling and manipulation.
//
// This package offers error boxing, stack traces, and panic recovery utilities
// to enhance error handling capabilities in Go applications.
//
// # Examples
//
//	// Box any error type
//	errBox := errutil.BoxErr(errors.New("something went wrong"))
//	fmt.Println(errBox.Error())
//
//	// Recover from panic with stack trace
//	defer errutil.Recover(func(err error) {
//		fmt.Printf("Panic recovered: %v\n", err)
//	})
package errutil

import (
	"errors"
	"fmt"
	"reflect"
)

type (
	// ErrBox is a wrapper for any error type.
	// Use ToError() method to convert to error interface, or access ErrBox through Result.Err().
	ErrBox struct {
		val any
	}
	// innerErrBox is an internal type that implements the error interface.
	// It wraps a value and provides error representation.
	// Created on-demand in ToError() when the wrapped value is not already an error.
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
		return &ErrBox{val: val.val}
	default:
		return &ErrBox{val: val}
	}
}

// IsEmpty returns true if ErrBox is empty (nil receiver or nil val).
//
//go:inline
func (e *ErrBox) IsEmpty() bool {
	return e == nil || e.val == nil
}

// Value returns the inner value.
//
//go:inline
func (e *ErrBox) Value() any {
	if e == nil {
		return nil
	}
	return e.val
}

// String returns the string representation.
// This implements the fmt.Stringer interface.
//
//go:inline
func (e *ErrBox) String() string {
	if e == nil {
		return "<nil>"
	}
	return errorString(e.val)
}

// errorString returns the string representation of a value for error display.
//
//go:inline
func errorString(val any) string {
	if val == nil {
		return "<nil>"
	}
	switch v := val.(type) {
	case string:
		return v
	case error:
		return v.Error()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// GoString returns the Go-syntax representation.
// This implements the fmt.GoStringer interface.
func (e *ErrBox) GoString() string {
	if e == nil {
		return "(*errutil.ErrBox)(nil)"
	}
	return fmt.Sprintf("&errutil.ErrBox{val: %#v}", e.val)
}

// ToError converts ErrBox to error interface.
// Returns nil if the receiver is nil or val is nil.
// Returns the wrapped error directly if it's already an error, otherwise
// returns a pointer to innerErrBox which implements the error interface.
//
// Example:
//
//	```go
//	var eb errutil.ErrBox = errutil.BoxErr(errors.New("test"))
//	var err error = eb.ToError() // err is the original error
//	```
//
//go:inline
func (e *ErrBox) ToError() error {
	if e == nil || e.val == nil {
		return nil
	}
	// If the wrapped value is already an error, return it directly
	if err, ok := e.val.(error); ok {
		return err
	}
	// Otherwise create and return an innerErrBox wrapper
	return &innerErrBox{val: e.val}
}

// Unwrap returns the inner error.
//
//go:inline
func (e *ErrBox) Unwrap() error {
	if e == nil || e.val == nil {
		return nil
	}
	return unwrapError(e.val)
}

// unwrapError extracts the inner error from a value.
//
//go:inline
func unwrapError(val any) error {
	switch v := val.(type) {
	case nil:
		return nil
	case error:
		u, ok := v.(interface {
			Unwrap() error
		})
		if ok {
			return u.Unwrap()
		}
		return v
	default:
		return nil
	}
}

// Is reports whether any error in err's chain matches target.
//
//go:inline
func (e *ErrBox) Is(target error) bool {
	if e == nil || e.val == nil {
		return target == nil
	}
	return isError(e.val, target)
}

// isError checks if the value or its unwrapped error matches the target.
//
//go:inline
func isError(val any, target error) bool {
	b := unwrapError(val)
	if b != nil {
		return errors.Is(b, target)
	}
	// If Unwrap returns nil, check if target matches the wrapped value
	if err, ok := val.(error); ok {
		return errors.Is(err, target)
	}
	return false
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false.
//
//go:inline
func (e *ErrBox) As(target any) bool {
	if e == nil || e.val == nil {
		return false
	}
	return asError(e.val, target)
}

// asError finds the first error in the value's chain that matches target.
//
//go:inline
func asError(val any, target any) bool {
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

	// If val is an error, handle error matching
	if err, ok := val.(error); ok {
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

// Error returns the string representation.
// This implements the error interface.
//
//go:inline
func (e innerErrBox) Error() string {
	return errorString(e.val)
}

// Unwrap returns the inner error.
//
//go:inline
func (e innerErrBox) Unwrap() error {
	return unwrapError(e.val)
}

// Is reports whether any error in err's chain matches target.
//
//go:inline
func (e innerErrBox) Is(target error) bool {
	return isError(e.val, target)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false.
//
//go:inline
func (e innerErrBox) As(target any) bool {
	return asError(e.val, target)
}
