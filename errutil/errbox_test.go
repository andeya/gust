package errutil_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/andeya/gust/errutil"
	"github.com/andeya/gust/result"
	"github.com/stretchr/testify/assert"
)

func TestBoxErr(t *testing.T) {
	// Test with string
	eb1 := errutil.BoxErr("test error")
	assert.NotNil(t, eb1)
	assert.Equal(t, "test error", eb1.Value())

	// Test with error
	err := errors.New("test error")
	eb2 := errutil.BoxErr(err)
	assert.NotNil(t, eb2)
	assert.Equal(t, err, eb2.Value())

	// Test with ErrBox (should return same)
	eb3 := errutil.BoxErr(eb2)
	assert.Equal(t, eb2, eb3)

	// Test with ErrBox value type (not pointer) (covers errbox.go:39-40)
	ebValue := errutil.BoxErr("test")
	var ebValueType errutil.ErrBox = *ebValue
	eb4 := errutil.BoxErr(ebValueType)
	assert.NotNil(t, eb4)
	assert.Equal(t, "test", eb4.Value())

	// Test with int
	eb6 := errutil.BoxErr(42)
	assert.NotNil(t, eb6)
	assert.Equal(t, 42, eb6.Value())

	// Test with nil
	eb7 := errutil.BoxErr(nil)
	assert.Nil(t, eb7)
}

func TestErrBox_Value(t *testing.T) {
	// Test with value
	eb := errutil.BoxErr("test")
	assert.Equal(t, "test", eb.Value())

	// Test with nil ErrBox
	var nilEb *errutil.ErrBox
	assert.Nil(t, nilEb.Value())
}

// TestErrBox_Is_NilReceiver tests Is method with nil receiver
func TestErrBox_Is_NilReceiver(t *testing.T) {
	var nilEb *errutil.ErrBox
	target := errors.New("test")
	assert.False(t, nilEb.Is(target))
}

func TestErrBox_ToError(t *testing.T) {
	// Test with string
	eb1 := errutil.BoxErr("test error")
	err1 := eb1.ToError()
	assert.NotNil(t, err1)
	assert.Contains(t, err1.Error(), "test error")

	// Test with error
	err := errors.New("test error")
	eb2 := errutil.BoxErr(err)
	err2 := eb2.ToError()
	assert.Equal(t, err, err2)

	// Test with nil value
	eb3 := errutil.BoxErr(nil)
	if eb3 == nil {
		assert.Nil(t, eb3)
	} else {
		err3 := eb3.ToError()
		assert.Nil(t, err3)
	}

	// Test with nil ErrBox pointer
	var nilEb *errutil.ErrBox
	err4 := nilEb.ToError()
	assert.Nil(t, err4)

	// Test with int
	eb4 := errutil.BoxErr(42)
	err5 := eb4.ToError()
	assert.NotNil(t, err5)
	assert.Contains(t, err5.Error(), "42")
}

func TestErrBox_Unwrap(t *testing.T) {
	// Test with error
	err := errors.New("test error")
	eb1 := errutil.BoxErr(err)
	unwrapped := eb1.Unwrap()
	assert.Equal(t, err, unwrapped)

	// Test with wrapped error
	wrappedErr := &wrappedError{err: err}
	eb2 := errutil.BoxErr(wrappedErr)
	unwrapped2 := eb2.Unwrap()
	assert.Equal(t, err, unwrapped2)

	// Test with nil value
	eb3 := errutil.BoxErr(nil)
	if eb3 == nil {
		assert.Nil(t, eb3)
	} else {
		assert.Nil(t, eb3.Unwrap())
	}

	// Test with nil ErrBox
	var nilEb *errutil.ErrBox
	assert.Nil(t, nilEb.Unwrap())

	// Test with non-error value
	eb4 := errutil.BoxErr(42)
	assert.Nil(t, eb4.Unwrap())

	// Test with error that doesn't implement Unwrap
	se := simpleError("simple")
	eb5 := errutil.BoxErr(se)
	unwrapped5 := eb5.Unwrap()
	assert.Equal(t, se, unwrapped5)
}

type wrappedError struct {
	err error
}

func (w *wrappedError) Error() string {
	return w.err.Error()
}

func (w *wrappedError) Unwrap() error {
	return w.err
}

type simpleError string

func (e simpleError) Error() string {
	return string(e)
}

// errorWithNilUnwrap2 is an error type that implements Unwrap() but returns nil
type errorWithNilUnwrap2 struct {
	msg string
}

func (e *errorWithNilUnwrap2) Error() string {
	return e.msg
}

func (e *errorWithNilUnwrap2) Unwrap() error {
	return nil // Explicitly return nil
}

func TestErrBox_Is(t *testing.T) {
	// Test with wrapped error
	err := errors.New("test error")
	eb1 := errutil.BoxErr(err)
	assert.True(t, eb1.Is(err))

	// Test with different error
	err2 := errors.New("different error")
	assert.False(t, eb1.Is(err2))

	// Test with nil ErrBox and nil target
	var nilEb *errutil.ErrBox
	assert.True(t, nilEb.Is(nil))

	// Test with non-nil ErrBox and nil target
	assert.False(t, eb1.Is(nil))

	// Test with wrapped error that has Unwrap
	wrappedErr := &wrappedError{err: err}
	eb2 := errutil.BoxErr(wrappedErr)
	assert.True(t, eb2.Is(err))

	// Test with wrapped error chain
	wrappedErr1 := errors.New("wrapped")
	wrappedErr2 := fmt.Errorf("outer: %w", wrappedErr1)
	eb3 := errutil.BoxErr(wrappedErr2)
	assert.True(t, eb3.Is(wrappedErr1)) // wrappedErr2 wraps wrappedErr1
}

func TestErrBox_As(t *testing.T) {
	// Test with wrapped error
	err := errors.New("test error")
	eb1 := errutil.BoxErr(err)
	var targetErr error
	assert.True(t, eb1.As(&targetErr))
	assert.Equal(t, err, targetErr)

	// Test with nil target (should panic)
	assert.Panics(t, func() {
		eb1.As(nil)
	})

	// Test with nil ErrBox
	var nilEb *errutil.ErrBox
	var nilTarget error
	assert.False(t, nilEb.As(&nilTarget))

	// Test with non-error value
	eb2 := errutil.BoxErr(42)
	var targetErr2 error
	assert.False(t, eb2.As(&targetErr2))
}

func TestErrBox_String(t *testing.T) {
	// Test with string
	eb1 := errutil.BoxErr("test error")
	assert.Equal(t, "test error", eb1.String())

	// Test with error
	err := errors.New("test error")
	eb2 := errutil.BoxErr(err)
	assert.Equal(t, "test error", eb2.String())

	// Test with nil value
	eb3 := errutil.BoxErr(nil)
	if eb3 == nil {
		assert.Nil(t, eb3)
	} else {
		assert.Equal(t, "<nil>", eb3.String())
	}

	// Test with nil ErrBox
	var nilEb *errutil.ErrBox
	assert.Equal(t, "<nil>", nilEb.String())

	// Test with int
	eb4 := errutil.BoxErr(42)
	assert.Equal(t, "42", eb4.String())
}

func TestErrBox_GoString(t *testing.T) {
	// Test with string
	eb1 := errutil.BoxErr("test")
	assert.Contains(t, eb1.GoString(), "ErrBox")
	assert.Contains(t, eb1.GoString(), "test")

	// Test with nil ErrBox
	var nilEb *errutil.ErrBox
	assert.Equal(t, "(*errutil.ErrBox)(nil)", nilEb.GoString())

	// Test with nil value
	eb2 := errutil.BoxErr(nil)
	if eb2 == nil {
		assert.Nil(t, eb2)
	} else {
		assert.Contains(t, eb2.GoString(), "ErrBox")
		assert.Contains(t, eb2.GoString(), "nil")
	}

	// Test with int
	eb3 := errutil.BoxErr(42)
	assert.Contains(t, eb3.GoString(), "ErrBox")
	assert.Contains(t, eb3.GoString(), "42")
}

// TestToError tests the ToError function
func TestToError(t *testing.T) {
	// Test toError with error type (through VoidResult.ToError)
	err := errors.New("test error")
	resultErr := result.TryErrVoid(err).Err()
	assert.NotNil(t, resultErr)
	assert.Equal(t, "test error", resultErr.Error())

	// Test toError with non-error type (through VoidResult.ToError)
	resultErr2 := result.TryErrVoid("test error").Err()
	assert.NotNil(t, resultErr2)
	assert.Equal(t, "test error", resultErr2.Error())

	// Test toError with int (through VoidResult.ToError)
	resultErr3 := result.TryErrVoid(42).Err()
	assert.NotNil(t, resultErr3)
	assert.Equal(t, "42", resultErr3.Error())
}

// TestInnerErrBox_Error_NilValue tests innerErrBox.Error() with nil value
func TestInnerErrBox_Error_NilValue(t *testing.T) {
	// Create an ErrBox with nil value
	eb := errutil.BoxErr(nil)
	// When val is nil, ToError returns nil, so we need to test Error() directly
	// We can't directly call Error() on innerErrBox, but we can test through String()
	assert.Equal(t, "<nil>", eb.String())

	// Test innerErrBox.Error() with nil val directly
	// We need to create an ErrBox with nil val and call ToError() to get innerErrBox
	// But ToError() returns nil when val is nil, so we need a different approach
	// Actually, we can create an ErrBox with a nil pointer value
	var nilPtr *int
	eb2 := errutil.BoxErr(nilPtr)
	err2 := eb2.ToError()
	assert.NotNil(t, err2)
	// When val is a nil pointer, Error() should return "<nil>"
	assert.Equal(t, "<nil>", err2.Error())
}

// TestInnerErrBox_Unwrap_NilValue tests innerErrBox.Unwrap() with nil value
func TestInnerErrBox_Unwrap_NilValue(t *testing.T) {
	// Create an ErrBox with nil value
	eb := errutil.BoxErr(nil)
	unwrapped := eb.Unwrap()
	assert.Nil(t, unwrapped)

	// Test innerErrBox.Unwrap() with nil val directly (case nil:)
	// We need to create an ErrBox with nil val and call ToError() to get innerErrBox
	// But ToError() returns nil when val is nil, so we need a different approach
	// Actually, we can create an ErrBox with a nil pointer value
	var nilPtr *int
	eb2 := errutil.BoxErr(nilPtr)
	err2 := eb2.ToError()
	assert.NotNil(t, err2)
	// When val is a nil pointer, Unwrap() should return nil (case nil:)
	unwrapped2 := eb2.Unwrap()
	assert.Nil(t, unwrapped2)
}

// TestInnerErrBox_Is_UnwrapNil tests innerErrBox.Is() when Unwrap returns nil
func TestInnerErrBox_Is_UnwrapNil(t *testing.T) {
	// Create an ErrBox with a non-error value (so Unwrap returns nil)
	eb := errutil.BoxErr(42)
	// When Unwrap returns nil, Is checks if the wrapped value matches target
	// Since 42 is not an error, Is should return false
	target := errors.New("test")
	assert.False(t, eb.Is(target))

	// Test with nil target
	assert.False(t, eb.Is(nil))

	// Test with an error value that doesn't implement Unwrap
	se := simpleError("simple")
	eb2 := errutil.BoxErr(se)
	// When wrapped value is an error but Unwrap returns nil (no Unwrap method),
	// Is should check the wrapped value directly
	assert.True(t, eb2.Is(se))
	assert.False(t, eb2.Is(errors.New("different")))

	// Test innerErrBox.Is() when Unwrap returns nil but val is an error
	// This covers the case where e.val.(error) is true but Unwrap() returns nil
	// We need an error that implements error but Unwrap() returns nil
	// Actually, simpleError already does this - it implements error but doesn't have Unwrap()
	// But we already tested that above. Let's test with a nil error pointer
	var nilErr *simpleError
	eb3 := errutil.BoxErr(nilErr)
	// When val is a nil error pointer, Unwrap() returns nil (case nil:)
	// Then Is() checks if val is an error - it is (*simpleError), but it's nil
	// So errors.Is(nil, target) should return false
	assert.False(t, eb3.Is(errors.New("test")))

	// Test innerErrBox.Is() when Unwrap returns nil but val is an error (non-nil)
	// This covers the case at line 1907: if err, ok := e.val.(error); ok
	// Use errorWithNilUnwrap2 which implements Unwrap() but returns nil
	errWithNilUnwrap2 := &errorWithNilUnwrap2{msg: "test"}
	eb4 := errutil.BoxErr(errWithNilUnwrap2)
	// When Unwrap() returns nil, Is() should check if e.val is an error
	// errWithNilUnwrap2 is an error, so it should check errors.Is(errWithNilUnwrap2, target)
	targetErr2 := errors.New("target")
	assert.False(t, eb4.Is(targetErr2))
	// But it should match itself
	assert.True(t, eb4.Is(errWithNilUnwrap2))
}

// TestInnerErrBox_Error_String tests innerErrBox.Error() with string value (covers errbox.go:151-152)
func TestInnerErrBox_Error_String(t *testing.T) {
	eb := errutil.BoxErr("test string")
	err := eb.ToError()
	assert.NotNil(t, err)
	// When val is string, Error() should return the string directly
	assert.Equal(t, "test string", err.Error())
}

// TestInnerErrBox_Unwrap_WithUnwrapInterface tests innerErrBox.Unwrap() with error that implements Unwrap (covers errbox.go:166-170)
func TestInnerErrBox_Unwrap_WithUnwrapInterface(t *testing.T) {
	baseErr := errors.New("base error")
	wrappedErr := &wrappedError{err: baseErr}
	eb := errutil.BoxErr(wrappedErr)

	unwrapped := eb.Unwrap()
	assert.NotNil(t, unwrapped)
	assert.Equal(t, baseErr, unwrapped)
}

// TestInnerErrBox_As_WithUnwrap tests innerErrBox.As() when Unwrap returns non-nil (covers errbox.go:198-199)
func TestInnerErrBox_As_WithUnwrap(t *testing.T) {
	baseErr := errors.New("base error")
	wrappedErr := &wrappedError{err: baseErr}
	eb := errutil.BoxErr(wrappedErr)

	var target *wrappedError
	assert.True(t, eb.As(&target))
	assert.NotNil(t, target)
	assert.Equal(t, wrappedErr, target)

	// Test with different target type
	var targetErr error
	assert.True(t, eb.As(&targetErr))
	assert.NotNil(t, targetErr)
}

// TestInnerErrBox_As_UnwrapNil tests innerErrBox.As() when Unwrap returns nil (covers errbox.go:201)
func TestInnerErrBox_As_UnwrapNil(t *testing.T) {
	// Test with non-error value (Unwrap returns nil)
	eb := errutil.BoxErr(42)
	var target int
	assert.False(t, eb.As(&target))

	// Test with error that doesn't implement Unwrap
	// Note: simpleError implements error interface, so errors.As should return true
	se := simpleError("simple")
	eb2 := errutil.BoxErr(se)
	var targetErr error
	// simpleError implements error, so As should return true and set targetErr to se
	assert.True(t, eb2.As(&targetErr))
	assert.NotNil(t, targetErr)
	assert.Equal(t, se, targetErr)
}

// TestInnerErrBox_As_ErrorDirectMatch tests innerErrBox.As() when error directly matches target type
func TestInnerErrBox_As_ErrorDirectMatch(t *testing.T) {
	// Test with wrappedError that directly matches target type
	baseErr := errors.New("base error")
	wrappedErr := &wrappedError{err: baseErr}
	eb := errutil.BoxErr(wrappedErr)

	var target *wrappedError
	assert.True(t, eb.As(&target))
	assert.NotNil(t, target)
	assert.Equal(t, wrappedErr, target)
}

// TestInnerErrBox_As_ErrorChain tests innerErrBox.As() when error chain needs to be traversed
func TestInnerErrBox_As_ErrorChain(t *testing.T) {
	// Create a chain of wrapped errors
	baseErr := errors.New("base error")
	wrappedErr1 := &wrappedError{err: baseErr}
	wrappedErr2 := fmt.Errorf("outer: %w", wrappedErr1)
	eb := errutil.BoxErr(wrappedErr2)

	// Target should match wrappedErr1 through the chain
	var target *wrappedError
	assert.True(t, eb.As(&target))
	assert.NotNil(t, target)
	assert.Equal(t, wrappedErr1, target)
}

// TestInnerErrBox_As_NonMatchingType tests innerErrBox.As() with non-matching target type
func TestInnerErrBox_As_NonMatchingType(t *testing.T) {
	eb := errutil.BoxErr(errors.New("test"))
	var target int
	// Target type doesn't match error type, should return false
	assert.False(t, eb.As(&target))
}

// TestInnerErrBox_As_NilTargetPanic tests innerErrBox.As() panics when target is nil
func TestInnerErrBox_As_NilTargetPanic(t *testing.T) {
	eb := errutil.BoxErr(errors.New("test"))
	assert.Panics(t, func() {
		eb.As(nil)
	})
}

// TestInnerErrBox_As_NonPointerTarget tests innerErrBox.As() with non-pointer target
func TestInnerErrBox_As_NonPointerTarget(t *testing.T) {
	eb := errutil.BoxErr(errors.New("test"))
	// Test with non-pointer target (should return false)
	var target error = errors.New("target")
	assert.False(t, eb.As(target))

	// Test with int (not a pointer)
	var targetInt int
	assert.False(t, eb.As(targetInt))

	// Test with string (not a pointer)
	var targetStr string
	assert.False(t, eb.As(targetStr))
}

// TestPanicError_Error tests panicError.Error() method
func TestPanicError_Error(t *testing.T) {
	// Test with nil error
	{
		pe := errutil.NewPanicError(nil, errutil.StackTrace{})
		assert.Equal(t, "<nil>", pe.Error())
	}

	// Test with error but no stack trace
	{
		err := errors.New("test error")
		pe := errutil.NewPanicError(err, errutil.StackTrace{})
		assert.Equal(t, "test error", pe.Error())
	}

	// Test with error and stack trace
	{
		err := errors.New("test error")
		stack := errutil.GetStackTrace(0)
		pe := errutil.NewPanicError(err, stack)
		// Error() should only return error message, not stack trace
		errorMsg := pe.Error()
		assert.Equal(t, "test error", errorMsg)
		assert.NotContains(t, errorMsg, "\n")
		// Verify stack trace exists but is not in Error()
		assert.True(t, len(stack) > 0)
	}
}

// TestPanicError_Unwrap tests panicError.Unwrap() method
func TestPanicError_Unwrap(t *testing.T) {
	// Test with nil error
	{
		pe := errutil.NewPanicError(nil, errutil.StackTrace{})
		assert.Nil(t, pe.Unwrap())
	}

	// Test with error
	{
		err := errors.New("test error")
		pe := errutil.NewPanicError(err, errutil.StackTrace{})
		assert.Equal(t, err, pe.Unwrap())
	}

	// Test with wrapped error
	{
		originalErr := errors.New("original error")
		wrappedErr := fmt.Errorf("wrapped: %w", originalErr)
		pe := errutil.NewPanicError(wrappedErr, errutil.StackTrace{})
		assert.Equal(t, wrappedErr, pe.Unwrap())
		// Verify error chain
		assert.True(t, errors.Is(pe.Unwrap(), originalErr))
	}
}

// TestPanicError_StackTrace tests panicError.StackTrace() method
func TestPanicError_StackTrace(t *testing.T) {
	// Test with empty stack trace
	{
		pe := errutil.NewPanicError(errors.New("test"), errutil.StackTrace{})
		stack := pe.StackTrace()
		assert.Equal(t, 0, len(stack))
	}

	// Test with stack trace
	{
		stack := errutil.GetStackTrace(0)
		pe := errutil.NewPanicError(errors.New("test"), stack)
		returnedStack := pe.StackTrace()
		assert.Equal(t, stack, returnedStack)
		assert.True(t, len(returnedStack) > 0)
	}
}

// TestPanicError_StackTraceCarrier tests that panicError implements StackTraceCarrier interface
func TestPanicError_StackTraceCarrier(t *testing.T) {
	var _ errutil.StackTraceCarrier = (*errutil.PanicError)(nil)
	stack := errutil.GetStackTrace(0)
	pe := errutil.NewPanicError(errors.New("test"), stack)
	// Verify it implements the interface
	var carrier errutil.StackTraceCarrier = pe
	assert.Equal(t, stack, carrier.StackTrace())
}

// TestNewPanicError tests NewPanicError function with various input types
func TestNewPanicError(t *testing.T) {
	stack := errutil.GetStackTrace(0)

	// Test with nil
	{
		pe := errutil.NewPanicError(nil, stack)
		assert.NotNil(t, pe)
		assert.Nil(t, pe.Unwrap())
		assert.Equal(t, stack, pe.StackTrace())
	}

	// Test with error
	{
		err := errors.New("test error")
		pe := errutil.NewPanicError(err, stack)
		assert.NotNil(t, pe)
		assert.Equal(t, err, pe.Unwrap())
		assert.Equal(t, stack, pe.StackTrace())
	}

	// Test with *ErrBox
	{
		eb := errutil.BoxErr(errors.New("errbox error"))
		pe := errutil.NewPanicError(eb, stack)
		assert.NotNil(t, pe)
		assert.NotNil(t, pe.Unwrap())
		assert.Equal(t, stack, pe.StackTrace())
		// Verify error message
		assert.Contains(t, pe.Error(), "errbox error")
	}

	// Test with ErrBox (value)
	{
		eb := errutil.BoxErr(errors.New("errbox value error"))
		pe := errutil.NewPanicError(*eb, stack)
		assert.NotNil(t, pe)
		assert.NotNil(t, pe.Unwrap())
		assert.Equal(t, stack, pe.StackTrace())
	}

	// Test with string
	{
		pe := errutil.NewPanicError("string panic", stack)
		assert.NotNil(t, pe)
		assert.NotNil(t, pe.Unwrap())
		assert.Equal(t, stack, pe.StackTrace())
		assert.Contains(t, pe.Error(), "string panic")
	}

	// Test with int
	{
		pe := errutil.NewPanicError(42, stack)
		assert.NotNil(t, pe)
		assert.NotNil(t, pe.Unwrap())
		assert.Equal(t, stack, pe.StackTrace())
		assert.Contains(t, pe.Error(), "42")
	}

	// Test with custom type
	{
		type CustomType struct {
			Value string
		}
		ct := CustomType{Value: "custom"}
		pe := errutil.NewPanicError(ct, stack)
		assert.NotNil(t, pe)
		assert.NotNil(t, pe.Unwrap())
		assert.Equal(t, stack, pe.StackTrace())
		assert.Contains(t, pe.Error(), "custom")
	}
}

// TestPanicError_Format tests panicError.Format() method
func TestPanicError_Format(t *testing.T) {
	// Test %v (error message only)
	{
		err := errors.New("test error")
		stack := errutil.GetStackTrace(0)
		pe := errutil.NewPanicError(err, stack)
		output := fmt.Sprintf("%v", pe)
		assert.Equal(t, "test error", output)
		assert.NotContains(t, output, "\n")
	}

	// Test %+v (error message with stack trace)
	{
		err := errors.New("test error")
		stack := errutil.GetStackTrace(0)
		pe := errutil.NewPanicError(err, stack)
		output := fmt.Sprintf("%+v", pe)
		assert.Contains(t, output, "test error")
		assert.Contains(t, output, "\n")
		// Should contain stack trace
		assert.True(t, len(stack) > 0)
	}

	// Test %s (error message only)
	{
		err := errors.New("test error")
		stack := errutil.GetStackTrace(0)
		pe := errutil.NewPanicError(err, stack)
		output := fmt.Sprintf("%s", pe)
		assert.Equal(t, "test error", output)
		assert.NotContains(t, output, "\n")
	}

	// Test with nil error
	{
		stack := errutil.GetStackTrace(0)
		pe := errutil.NewPanicError(nil, stack)
		assert.Equal(t, "<nil>", fmt.Sprintf("%v", pe))
		// %+v with nil error but stack trace should still show stack
		output := fmt.Sprintf("%+v", pe)
		assert.Contains(t, output, "<nil>")
		// If stack exists, %+v should include it
		if len(pe.StackTrace()) > 0 {
			assert.Contains(t, output, "\n")
		}
		assert.Equal(t, "<nil>", fmt.Sprintf("%s", pe))
	}

	// Test with empty stack trace
	{
		err := errors.New("test error")
		pe := errutil.NewPanicError(err, errutil.StackTrace{})
		output := fmt.Sprintf("%+v", pe)
		assert.Equal(t, "test error", output)
		assert.NotContains(t, output, "\n")
	}
}

// TestPanicError_Format_AllVerbs tests all format verbs for panicError
func TestPanicError_Format_AllVerbs(t *testing.T) {
	err := errors.New("test error")
	stack := errutil.GetStackTrace(0)
	pe := errutil.NewPanicError(err, stack)

	// Test %v (error message only)
	outputV := fmt.Sprintf("%v", pe)
	assert.Equal(t, "test error", outputV)
	assert.NotContains(t, outputV, "\n")

	// Test %+v (error message with detailed stack trace)
	outputPlusV := fmt.Sprintf("%+v", pe)
	assert.Contains(t, outputPlusV, "test error")
	assert.Contains(t, outputPlusV, "\n")

	// Test %s (error message only)
	outputS := fmt.Sprintf("%s", pe)
	assert.Equal(t, "test error", outputS)
	assert.NotContains(t, outputS, "\n")

	// Test %+s (error message with stack trace in %s format)
	outputPlusS := fmt.Sprintf("%+s", pe)
	assert.Contains(t, outputPlusS, "test error")
	assert.Contains(t, outputPlusS, "\n")

	// Test unsupported verbs (should default to %s behavior)
	outputD := fmt.Sprintf("%d", pe)
	assert.Equal(t, "test error", outputD)
	outputN := fmt.Sprintf("%n", pe)
	assert.Equal(t, "test error", outputN)
}

// TestPanicError_Format_EmptyStack tests Format with empty stack trace
func TestPanicError_Format_EmptyStack(t *testing.T) {
	err := errors.New("test error")
	pe := errutil.NewPanicError(err, errutil.StackTrace{})

	// %v should only show error message
	outputV := fmt.Sprintf("%v", pe)
	assert.Equal(t, "test error", outputV)

	// %+v with empty stack should only show error message
	outputPlusV := fmt.Sprintf("%+v", pe)
	assert.Equal(t, "test error", outputPlusV)
	assert.NotContains(t, outputPlusV, "\n")

	// %s should only show error message
	outputS := fmt.Sprintf("%s", pe)
	assert.Equal(t, "test error", outputS)

	// %+s with empty stack should only show error message
	outputPlusS := fmt.Sprintf("%+s", pe)
	assert.Equal(t, "test error", outputPlusS)
	assert.NotContains(t, outputPlusS, "\n")
}

// TestPanicError_Unwrap_Chain tests error unwrapping chain
func TestPanicError_Unwrap_Chain(t *testing.T) {
	baseErr := errors.New("base error")
	wrappedErr := fmt.Errorf("wrapped: %w", baseErr)
	pe := errutil.NewPanicError(wrappedErr, errutil.GetStackTrace(0))

	// Unwrap should return the wrapped error
	unwrapped := pe.Unwrap()
	assert.Equal(t, wrappedErr, unwrapped)

	// Error chain should be preserved
	assert.True(t, errors.Is(pe, baseErr))
	assert.True(t, errors.Is(pe, wrappedErr))
}
