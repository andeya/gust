package examples_test

import (
	"errors"
	"fmt"

	"github.com/andeya/gust/errutil"
)

// Example_errutil_errorBoxing demonstrates error boxing.
func Example_errutil_errorBoxing() {
	// Box any error type
	err := errors.New("something went wrong")
	errBox := errutil.BoxErr(err)
	fmt.Println("Boxed error:", errBox.ToError().Error())

	// Box nil (returns empty ErrBox)
	emptyBox := errutil.BoxErr(nil)
	fmt.Println("Empty box:", emptyBox.IsEmpty())

	// Convert to error interface
	boxedErr := errBox.ToError()
	fmt.Println("As error:", boxedErr.Error())

	// Output:
	// Boxed error: something went wrong
	// Empty box: true
	// As error: something went wrong
}

// Example_errutil_stackTraces demonstrates stack trace utilities.
func Example_errutil_stackTraces() {
	// Get current stack trace
	stack := errutil.GetStackTrace(0)
	fmt.Println("Stack trace available:", len(stack) > 0)

	// Format stack trace
	if len(stack) > 0 {
		firstFrame := stack[0]
		formatted := fmt.Sprintf("%v", firstFrame)
		fmt.Println("First frame formatted:", len(formatted) > 0)
	}

	// Output:
	// Stack trace available: true
	// First frame formatted: true
}

// Example_errutil_panicRecovery demonstrates panic recovery with stack traces.
func Example_errutil_panicRecovery() {
	// Recover from panic with stack trace
	defer func() {
		if r := recover(); r != nil {
			// Get panic stack trace
			stack := errutil.PanicStackTrace()
			panicErr := errutil.NewPanicError(r, stack)
			fmt.Println("Panic recovered:", panicErr.Error())
			fmt.Println("Stack trace available:", len(stack) > 0)
		}
	}()

	// Simulate a panic
	panic("test panic")
	// Output:
	// Panic recovered: test panic
	// Stack trace available: true
}

// Example_errutil_errorChaining demonstrates error chaining with ErrBox.
func Example_errutil_errorChaining() {
	// Create a chain of errors
	err1 := errors.New("first error")
	err2 := fmt.Errorf("second error: %w", err1)
	err3 := fmt.Errorf("third error: %w", err2)

	// Box the error chain
	boxed := errutil.BoxErr(err3)
	boxedErr := boxed.ToError()

	// Check if errors.Is works
	if errors.Is(boxedErr, err1) {
		fmt.Println("Error chain preserved")
	}

	// Unwrap error
	unwrapped := errors.Unwrap(boxedErr)
	if unwrapped != nil {
		fmt.Println("Unwrapped error:", unwrapped.Error())
	}

	// Output:
	// Error chain preserved
	// Unwrapped error: second error: first error
}

// Example_errutil_panicError demonstrates PanicError with stack traces.
func Example_errutil_panicError() {
	// Create a panic error with stack trace
	panicErr := errors.New("panic occurred")
	stack := errutil.GetStackTrace(0)
	panicError := errutil.NewPanicError(panicErr, stack)

	// Format error message
	fmt.Println("Error:", panicError.Error())

	// Check if stack trace is available
	fmt.Println("Stack trace available:", len(panicError.StackTrace()) > 0)

	// Output:
	// Error: panic occurred
	// Stack trace available: true
}

// Example_errutil_errorWrapping demonstrates wrapping errors with ErrBox.
func Example_errutil_errorWrapping() {
	// Wrap a string as error
	strErr := errutil.BoxErr("string error")
	fmt.Println("String error:", strErr.ToError().Error())

	// Wrap an integer as error
	intErr := errutil.BoxErr(42)
	fmt.Println("Int error:", intErr.ToError().Error())

	// Wrap a custom type
	type CustomError struct {
		Code    int
		Message string
	}
	customErr := errutil.BoxErr(CustomError{Code: 404, Message: "Not Found"})
	fmt.Println("Custom error:", customErr.ToError().Error())

	// Output:
	// String error: string error
	// Int error: 42
	// Custom error: {404 Not Found}
}
