package gust

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoxErr(t *testing.T) {
	// Test with string
	eb1 := BoxErr("test error")
	assert.NotNil(t, eb1)
	assert.Equal(t, "test error", eb1.Value())

	// Test with error
	err := errors.New("test error")
	eb2 := BoxErr(err)
	assert.NotNil(t, eb2)
	assert.Equal(t, err, eb2.Value())

	// Test with ErrBox (should return same)
	eb3 := BoxErr(eb2)
	assert.Equal(t, eb2, eb3)

	// Test with ErrBox value type (not pointer) (covers errbox.go:34-35)
	ebValue := ErrBox{inner: innerErrBox{val: "test"}}
	eb4 := BoxErr(ebValue)
	assert.NotNil(t, eb4)
	assert.Equal(t, "test", eb4.Value())

	// Test with int
	eb5 := BoxErr(42)
	assert.NotNil(t, eb5)
	assert.Equal(t, 42, eb5.Value())

	// Test with nil
	eb6 := BoxErr(nil)
	assert.Nil(t, eb6)
}

func TestErrBox_Value(t *testing.T) {
	// Test with value
	eb := BoxErr("test")
	assert.Equal(t, "test", eb.Value())

	// Test with nil ErrBox
	var nilEb *ErrBox
	assert.Nil(t, nilEb.Value())
}

// TestErrBox_Is_NilReceiver tests Is method with nil receiver
func TestErrBox_Is_NilReceiver(t *testing.T) {
	var nilEb *ErrBox
	target := errors.New("test")
	assert.False(t, nilEb.Is(target))
}

func TestErrBox_ToError(t *testing.T) {
	// Test with string
	eb1 := BoxErr("test error")
	err1 := eb1.ToError()
	assert.NotNil(t, err1)
	assert.Contains(t, err1.Error(), "test error")

	// Test with error
	err := errors.New("test error")
	eb2 := BoxErr(err)
	err2 := eb2.ToError()
	assert.Equal(t, err, err2)

	// Test with nil value
	eb3 := BoxErr(nil)
	if eb3 == nil {
		assert.Nil(t, eb3)
	} else {
		err3 := eb3.ToError()
		assert.Nil(t, err3)
	}

	// Test with nil ErrBox pointer
	var nilEb *ErrBox
	err4 := nilEb.ToError()
	assert.Nil(t, err4)

	// Test with int
	eb4 := BoxErr(42)
	err5 := eb4.ToError()
	assert.NotNil(t, err5)
	assert.Contains(t, err5.Error(), "42")
}

func TestErrBox_Unwrap(t *testing.T) {
	// Test with error
	err := errors.New("test error")
	eb1 := BoxErr(err)
	unwrapped := eb1.Unwrap()
	assert.Equal(t, err, unwrapped)

	// Test with wrapped error
	wrappedErr := &wrappedError{err: err}
	eb2 := BoxErr(wrappedErr)
	unwrapped2 := eb2.Unwrap()
	assert.Equal(t, err, unwrapped2)

	// Test with nil value
	eb3 := BoxErr(nil)
	if eb3 == nil {
		assert.Nil(t, eb3)
	} else {
		assert.Nil(t, eb3.Unwrap())
	}

	// Test with nil ErrBox
	var nilEb *ErrBox
	assert.Nil(t, nilEb.Unwrap())

	// Test with non-error value
	eb4 := BoxErr(42)
	assert.Nil(t, eb4.Unwrap())

	// Test with error that doesn't implement Unwrap
	se := simpleError("simple")
	eb5 := BoxErr(se)
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

func TestErrBox_Is(t *testing.T) {
	// Test with wrapped error
	err := errors.New("test error")
	eb1 := BoxErr(err)
	assert.True(t, eb1.Is(err))

	// Test with different error
	err2 := errors.New("different error")
	assert.False(t, eb1.Is(err2))

	// Test with nil ErrBox and nil target
	var nilEb *ErrBox
	assert.True(t, nilEb.Is(nil))

	// Test with non-nil ErrBox and nil target
	assert.False(t, eb1.Is(nil))

	// Test with wrapped error that has Unwrap
	wrappedErr := &wrappedError{err: err}
	eb2 := BoxErr(wrappedErr)
	assert.True(t, eb2.Is(err))

	// Test with wrapped error chain
	wrappedErr1 := errors.New("wrapped")
	wrappedErr2 := fmt.Errorf("outer: %w", wrappedErr1)
	eb3 := BoxErr(wrappedErr2)
	assert.True(t, eb3.Is(wrappedErr1)) // wrappedErr2 wraps wrappedErr1
}

func TestErrBox_As(t *testing.T) {
	// Test with wrapped error
	err := errors.New("test error")
	eb1 := BoxErr(err)
	var targetErr error
	assert.True(t, eb1.As(&targetErr))
	assert.Equal(t, err, targetErr)

	// Test with nil target (should panic)
	assert.Panics(t, func() {
		eb1.As(nil)
	})

	// Test with nil ErrBox
	var nilEb *ErrBox
	var nilTarget error
	assert.False(t, nilEb.As(&nilTarget))

	// Test with non-error value
	eb2 := BoxErr(42)
	var targetErr2 error
	assert.False(t, eb2.As(&targetErr2))
}

func TestErrBox_ValueOrDefault(t *testing.T) {
	// Test with value
	eb := BoxErr("test")
	assert.Equal(t, "test", eb.ValueOrDefault())

	// Test with nil ErrBox
	var nilEb *ErrBox
	assert.Nil(t, nilEb.ValueOrDefault())

	// Test with nil value
	eb2 := BoxErr(nil)
	if eb2 == nil {
		assert.Nil(t, eb2)
	} else {
		assert.Nil(t, eb2.ValueOrDefault())
	}

	// Test with int value
	eb3 := BoxErr(42)
	assert.Equal(t, 42, eb3.ValueOrDefault())

	// Test with error value
	err := errors.New("test error")
	eb4 := BoxErr(err)
	assert.Equal(t, err, eb4.ValueOrDefault())
}

func TestErrBox_String(t *testing.T) {
	// Test with string
	eb1 := BoxErr("test error")
	assert.Equal(t, "test error", eb1.String())

	// Test with error
	err := errors.New("test error")
	eb2 := BoxErr(err)
	assert.Equal(t, "test error", eb2.String())

	// Test with nil value
	eb3 := BoxErr(nil)
	if eb3 == nil {
		assert.Nil(t, eb3)
	} else {
		assert.Equal(t, "<nil>", eb3.String())
	}

	// Test with nil ErrBox
	var nilEb *ErrBox
	assert.Equal(t, "<nil>", nilEb.String())

	// Test with int
	eb4 := BoxErr(42)
	assert.Equal(t, "42", eb4.String())
}

func TestErrBox_GoString(t *testing.T) {
	// Test with string
	eb1 := BoxErr("test")
	assert.Contains(t, eb1.GoString(), "ErrBox")
	assert.Contains(t, eb1.GoString(), "test")

	// Test with nil ErrBox
	var nilEb *ErrBox
	assert.Equal(t, "(*gust.ErrBox)(nil)", nilEb.GoString())

	// Test with nil value
	eb2 := BoxErr(nil)
	if eb2 == nil {
		assert.Nil(t, eb2)
	} else {
		assert.Contains(t, eb2.GoString(), "ErrBox")
		assert.Contains(t, eb2.GoString(), "nil")
	}

	// Test with int
	eb3 := BoxErr(42)
	assert.Contains(t, eb3.GoString(), "ErrBox")
	assert.Contains(t, eb3.GoString(), "42")
}

// TestToError tests the toError function indirectly through Errable.ToError
func TestToError(t *testing.T) {
	// Test toError with error type (through VoidResult.ToError)
	err := errors.New("test error")
	result := RetVoid(err)
	resultErr := ToError(result)
	assert.NotNil(t, resultErr)
	assert.Equal(t, "test error", resultErr.Error())

	// Test toError with non-error type (through VoidResult.ToError)
	result2 := Err[Void]("test error")
	resultErr2 := ToError(result2)
	assert.NotNil(t, resultErr2)
	assert.Equal(t, "test error", resultErr2.Error())

	// Test toError with int (through VoidResult.ToError)
	result3 := Err[Void](42)
	resultErr3 := ToError(result3)
	assert.NotNil(t, resultErr3)
	assert.Equal(t, "42", resultErr3.Error())
}

// TestInnerErrBox_Error_NilValue tests innerErrBox.Error() with nil value
func TestInnerErrBox_Error_NilValue(t *testing.T) {
	// Create an ErrBox with nil value
	eb := ErrBox{inner: innerErrBox{val: nil}}
	// When val is nil, ToError returns nil, so we need to test Error() directly
	// We can't directly call Error() on innerErrBox, but we can test through String()
	assert.Equal(t, "<nil>", eb.String())
}

// TestInnerErrBox_Unwrap_NilValue tests innerErrBox.Unwrap() with nil value
func TestInnerErrBox_Unwrap_NilValue(t *testing.T) {
	// Create an ErrBox with nil value
	eb := ErrBox{inner: innerErrBox{val: nil}}
	unwrapped := eb.Unwrap()
	assert.Nil(t, unwrapped)
}

// TestInnerErrBox_Is_UnwrapNil tests innerErrBox.Is() when Unwrap returns nil
func TestInnerErrBox_Is_UnwrapNil(t *testing.T) {
	// Create an ErrBox with a non-error value (so Unwrap returns nil)
	eb := BoxErr(42)
	// When Unwrap returns nil, Is checks if the wrapped value matches target
	// Since 42 is not an error, Is should return false
	target := errors.New("test")
	assert.False(t, eb.Is(target))

	// Test with nil target
	assert.False(t, eb.Is(nil))

	// Test with an error value that doesn't implement Unwrap
	se := simpleError("simple")
	eb2 := BoxErr(se)
	// When wrapped value is an error but Unwrap returns nil (no Unwrap method),
	// Is should check the wrapped value directly
	assert.True(t, eb2.Is(se))
	assert.False(t, eb2.Is(errors.New("different")))
}

// TestInnerErrBox_Error_String tests innerErrBox.Error() with string value (covers errbox.go:151-152)
func TestInnerErrBox_Error_String(t *testing.T) {
	eb := BoxErr("test string")
	err := eb.ToError()
	assert.NotNil(t, err)
	// When val is string, Error() should return the string directly
	assert.Equal(t, "test string", err.Error())
}

// TestInnerErrBox_Unwrap_WithUnwrapInterface tests innerErrBox.Unwrap() with error that implements Unwrap (covers errbox.go:166-170)
func TestInnerErrBox_Unwrap_WithUnwrapInterface(t *testing.T) {
	baseErr := errors.New("base error")
	wrappedErr := &wrappedError{err: baseErr}
	eb := BoxErr(wrappedErr)

	unwrapped := eb.Unwrap()
	assert.NotNil(t, unwrapped)
	assert.Equal(t, baseErr, unwrapped)
}

// TestInnerErrBox_As_WithUnwrap tests innerErrBox.As() when Unwrap returns non-nil (covers errbox.go:198-199)
func TestInnerErrBox_As_WithUnwrap(t *testing.T) {
	baseErr := errors.New("base error")
	wrappedErr := &wrappedError{err: baseErr}
	eb := BoxErr(wrappedErr)

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
	eb := BoxErr(42)
	var target int
	assert.False(t, eb.As(&target))

	// Test with error that doesn't implement Unwrap
	// Note: simpleError implements error interface, so errors.As should return true
	se := simpleError("simple")
	eb2 := BoxErr(se)
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
	eb := BoxErr(wrappedErr)

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
	eb := BoxErr(wrappedErr2)

	// Target should match wrappedErr1 through the chain
	var target *wrappedError
	assert.True(t, eb.As(&target))
	assert.NotNil(t, target)
	assert.Equal(t, wrappedErr1, target)
}

// TestInnerErrBox_As_NonMatchingType tests innerErrBox.As() with non-matching target type
func TestInnerErrBox_As_NonMatchingType(t *testing.T) {
	eb := BoxErr(errors.New("test"))
	var target int
	// Target type doesn't match error type, should return false
	assert.False(t, eb.As(&target))
}

// TestInnerErrBox_As_NilTargetPanic tests innerErrBox.As() panics when target is nil
func TestInnerErrBox_As_NilTargetPanic(t *testing.T) {
	eb := BoxErr(errors.New("test"))
	assert.Panics(t, func() {
		eb.As(nil)
	})
}
