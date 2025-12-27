package gust

import (
	"errors"
	"fmt"
	"io"
	"reflect"
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
	// panicError wraps an error with its panic stack trace.
	// This is used by Catch() to preserve stack trace information when converting panics to errors.
	panicError struct {
		err   error
		stack StackTrace
	}
)

var (
	_ fmt.Stringer      = (*ErrBox)(nil)
	_ fmt.GoStringer    = (*ErrBox)(nil)
	_ fmt.Formatter     = (*panicError)(nil)
	_ error             = (*innerErrBox)(nil)
	_ error             = (*panicError)(nil)
	_ StackTraceCarrier = (*panicError)(nil)
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

// Error returns the error message without stack trace information.
// To include stack trace, use fmt.Sprintf("%+v", e) or access StackTrace() directly.
func (e *panicError) Error() string {
	if e.err == nil {
		return "<nil>"
	}
	return e.err.Error()
}

// Format formats the panicError according to the fmt.Formatter interface.
// The formatting behavior is consistent with Frame.Format() and StackTrace.Format():
//
//	%v    error message only (same as Error())
//	%+v   error message with detailed stack trace (each frame shows function name, file path, and line)
//	%s    error message only (same as Error())
//	%+s   error message with stack trace (each frame shows file name only)
//
// Note: Unlike Frame.Format() which supports 's', 'd', 'n', 'v', panicError only supports 'v' and 's'.
// The 'd' (line number) and 'n' (function name) verbs are not applicable to panicError.
// For unsupported verbs, the behavior defaults to 's' (error message only).
//
// Similar to Frame.Format(), %v recursively calls Format(s, 's') for the error message part,
// then conditionally appends stack trace based on flags.
//
// Example:
//
//	err := result.Err()
//	fmt.Printf("%v", err)    // "test error"
//	fmt.Printf("%+v", err)   // "test error\n\nfunction_name\n\tfile_path:line\n..."
//	fmt.Printf("%s", err)    // "test error"
//	fmt.Printf("%+s", err)   // "test error\n\n[file:line file:line ...]"
func (e *panicError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		// Similar to Frame.Format(): recursively call Format(s, 's') for error message
		e.Format(s, 's')
		// If + flag is set and stack trace exists, append detailed stack trace
		// Pass the verb 'v' and flags to stack trace formatting
		// This ensures %+v on panicError results in %+v on StackTrace
		if s.Flag('+') && len(e.stack) > 0 {
			io.WriteString(s, "\n")
			e.stack.Format(s, verb)
		}
	case 's':
		// %s displays error message (base formatting)
		if e.err == nil {
			io.WriteString(s, "<nil>")
			return
		}
		io.WriteString(s, e.err.Error())
		// If + flag is set and stack trace exists, append stack trace with %s formatting
		if s.Flag('+') && len(e.stack) > 0 {
			io.WriteString(s, "\n")
			// Use %s formatting for stack trace (shows file names only, via formatSlice)
			e.stack.Format(s, verb)
		}
	default:
		// For unsupported verbs ('d', 'n', etc.), default to 's' behavior
		// This provides graceful fallback for unknown format verbs
		e.Format(s, 's')
	}
}

// Unwrap returns the wrapped error.
func (e *panicError) Unwrap() error {
	return e.err
}

// StackTrace returns the panic stack trace.
func (e *panicError) StackTrace() StackTrace {
	return e.stack
}

// newPanicError creates a new panicError with the given error and stack trace.
func newPanicError(err any, stack StackTrace) *panicError {
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
	return &panicError{
		err:   wrappedErr,
		stack: stack,
	}
}
