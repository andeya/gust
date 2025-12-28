package errutil

import (
	"fmt"
	"io"
)

// PanicError wraps an error with its panic stack trace.
// This is used by Catch() to preserve stack trace information when converting panics to errors.
type PanicError struct {
	err   error
	stack StackTrace
}

var (
	_ fmt.Formatter     = (*PanicError)(nil)
	_ error             = (*PanicError)(nil)
	_ StackTraceCarrier = (*PanicError)(nil)
)

// NewPanicError creates a new PanicError with the given error and stack trace.
func NewPanicError(err any, stack StackTrace) *PanicError {
	var wrappedErr error
	switch e := err.(type) {
	case nil:
		wrappedErr = nil
	case error:
		wrappedErr = e
	default:
		wrappedErr = fmt.Errorf("%v", e)
	}
	return &PanicError{
		err:   wrappedErr,
		stack: stack,
	}
}

// Error returns the error message without stack trace information.
// To include stack trace, use fmt.Sprintf("%+v", e) or access StackTrace() directly.
func (e *PanicError) Error() string {
	if e.err == nil {
		return "<nil>"
	}
	return e.err.Error()
}

// Format formats the PanicError according to the fmt.Formatter interface.
// The formatting behavior is consistent with Frame.Format() and StackTrace.Format():
//
//	%v    error message only (same as Error())
//	%+v   error message with detailed stack trace (each frame shows function name, file path, and line)
//	%s    error message only (same as Error())
//	%+s   error message with stack trace (each frame shows file name only)
//
// Note: Unlike Frame.Format() which supports 's', 'd', 'n', 'v', PanicError only supports 'v' and 's'.
// The 'd' (line number) and 'n' (function name) verbs are not applicable to PanicError.
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
func (e *PanicError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if e.err == nil {
			io.WriteString(s, "<nil>")
		} else {
			io.WriteString(s, e.err.Error())
		}
		// If + flag is set and stack trace exists, append detailed stack trace
		// Pass the verb 'v' and flags to stack trace formatting
		// This ensures %+v on PanicError results in %+v on StackTrace
		if s.Flag('+') && len(e.stack) > 0 {
			io.WriteString(s, "\n")
			e.stack.Format(s, verb)
		}
	case 's':
		// %s displays error message (base formatting)
		if e.err == nil {
			io.WriteString(s, "<nil>")
		} else {
			io.WriteString(s, e.err.Error())
		}
		// If + flag is set and stack trace exists, append stack trace with %s formatting
		if s.Flag('+') && len(e.stack) > 0 {
			io.WriteString(s, "\n")
			// Use %s formatting for stack trace (shows file names only, via formatSlice)
			e.stack.Format(s, verb)
		}
	default:
		// For unsupported verbs, default to 's' behavior (error message only)
		e.Format(s, 's')
	}
}

// Unwrap returns the wrapped error.
func (e *PanicError) Unwrap() error {
	return e.err
}

// StackTrace returns the stack trace associated with this panic error.
func (e *PanicError) StackTrace() StackTrace {
	return e.stack
}
