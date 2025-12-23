package gust

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToErrBox(t *testing.T) {
	// Test with string
	eb1 := ToErrBox("test error")
	assert.NotNil(t, eb1)
	assert.Equal(t, "test error", eb1.Value())

	// Test with error
	err := errors.New("test error")
	eb2 := ToErrBox(err)
	assert.NotNil(t, eb2)
	assert.Equal(t, err, eb2.Value())

	// Test with ErrBox (should return same)
	eb3 := ToErrBox(eb2)
	assert.Equal(t, eb2, eb3)

	// Test with int
	eb4 := ToErrBox(42)
	assert.NotNil(t, eb4)
	assert.Equal(t, 42, eb4.Value())

	// Test with nil
	eb5 := ToErrBox(nil)
	assert.NotNil(t, eb5)
	assert.Nil(t, eb5.Value())
}

func TestErrBox_Value(t *testing.T) {
	// Test with value
	eb := &ErrBox{val: "test"}
	assert.Equal(t, "test", eb.Value())

	// Test with nil ErrBox
	var nilEb *ErrBox
	assert.Nil(t, nilEb.Value())
}

func TestErrBox_Error(t *testing.T) {
	// Test with string
	eb1 := &ErrBox{val: "test error"}
	assert.Equal(t, "test error", eb1.Error())

	// Test with error
	err := errors.New("test error")
	eb2 := &ErrBox{val: err}
	assert.Equal(t, "test error", eb2.Error())

	// Test with nil value
	eb3 := &ErrBox{val: nil}
	assert.Equal(t, "", eb3.Error())

	// Test with nil ErrBox
	var nilEb *ErrBox
	assert.Equal(t, "", nilEb.Error())

	// Test with int
	eb4 := &ErrBox{val: 42}
	assert.Equal(t, "42", eb4.Error())
}

func TestErrBox_Unwrap(t *testing.T) {
	// Test with error
	err := errors.New("test error")
	eb1 := &ErrBox{val: err}
	unwrapped := eb1.Unwrap()
	assert.Equal(t, err, unwrapped)

	// Test with wrapped error
	wrappedErr := &wrappedError{err: err}
	eb2 := &ErrBox{val: wrappedErr}
	unwrapped2 := eb2.Unwrap()
	assert.Equal(t, err, unwrapped2)

	// Test with nil value
	eb3 := &ErrBox{val: nil}
	assert.Nil(t, eb3.Unwrap())

	// Test with nil ErrBox
	var nilEb *ErrBox
	assert.Nil(t, nilEb.Unwrap())

	// Test with non-error value
	eb4 := &ErrBox{val: 42}
	assert.Nil(t, eb4.Unwrap())
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

func TestErrBox_Is(t *testing.T) {
	// Test with same ErrBox
	eb1 := &ErrBox{val: "test"}
	eb2 := &ErrBox{val: "test"}
	assert.True(t, eb1.Is(eb2))

	// Test with different ErrBox
	eb3 := &ErrBox{val: "different"}
	assert.False(t, eb1.Is(eb3))

	// Test with nil target
	var nilEb *ErrBox
	assert.False(t, eb1.Is(nilEb))
	assert.True(t, nilEb.Is(nilEb))

	// Test with wrapped error
	err := errors.New("test error")
	eb4 := &ErrBox{val: err}
	assert.True(t, eb4.Is(err))

	// Test with different error
	err2 := errors.New("different error")
	assert.False(t, eb4.Is(err2))
}

func TestErrBox_As(t *testing.T) {
	// Test with ErrBox target (non-nil target)
	eb1 := &ErrBox{val: "test"}
	target := &ErrBox{} // Create a non-nil ErrBox
	assert.True(t, eb1.As(target))
	assert.Equal(t, "test", target.Value())

	// Test with nil target (should panic)
	assert.Panics(t, func() {
		eb1.As(nil)
	})

	// Test with nil ErrBox target (should return false)
	var nilTarget *ErrBox
	assert.False(t, eb1.As(nilTarget))

	// Test with wrapped error
	err := errors.New("test error")
	eb2 := &ErrBox{val: err}
	var targetErr error
	assert.True(t, eb2.As(&targetErr))
	assert.Equal(t, err, targetErr)

	// Test with non-error value
	eb3 := &ErrBox{val: 42}
	var targetErr2 error
	assert.False(t, eb3.As(&targetErr2))
}
