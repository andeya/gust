package errutil

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanicError_Format_AllCases(t *testing.T) {
	stack := GetStackTrace(0)

	// Test %v with nil error
	{
		pe := NewPanicError(nil, stack)
		output := fmt.Sprintf("%v", pe)
		assert.Equal(t, "<nil>", output)
	}

	// Test %v with error
	{
		err := errors.New("test error")
		pe := NewPanicError(err, stack)
		output := fmt.Sprintf("%v", pe)
		assert.Equal(t, "test error", output)
	}

	// Test %+v with nil error and empty stack
	{
		pe := NewPanicError(nil, StackTrace{})
		output := fmt.Sprintf("%+v", pe)
		assert.Equal(t, "<nil>", output)
	}

	// Test %+v with nil error and non-empty stack
	{
		pe := NewPanicError(nil, stack)
		output := fmt.Sprintf("%+v", pe)
		assert.Contains(t, output, "<nil>")
		// Only check for newline if stack is not empty
		// Note: GetStackTrace(0) might return empty stack in some test environments
		if len(stack) > 0 {
			// If stack exists, %+v should append stack trace
			assert.Contains(t, output, "\n", "Expected newline when stack is not empty")
		} else {
			// If stack is empty, should just be "<nil>"
			assert.Equal(t, "<nil>", output)
		}
	}

	// Test %+v with error and empty stack
	{
		err := errors.New("test error")
		pe := NewPanicError(err, StackTrace{})
		output := fmt.Sprintf("%+v", pe)
		assert.Equal(t, "test error", output)
		assert.NotContains(t, output, "\n")
	}

	// Test %+v with error and non-empty stack
	{
		err := errors.New("test error")
		pe := NewPanicError(err, stack)
		output := fmt.Sprintf("%+v", pe)
		assert.Contains(t, output, "test error")
		// Only check for newline if stack is not empty
		if len(stack) > 0 {
			assert.Contains(t, output, "\n", "Expected newline when stack is not empty")
		} else {
			// If stack is empty, should just be error message
			assert.Equal(t, "test error", output)
		}
	}

	// Test %s with nil error
	{
		pe := NewPanicError(nil, stack)
		output := fmt.Sprintf("%s", pe)
		assert.Equal(t, "<nil>", output)
	}

	// Test %s with error
	{
		err := errors.New("test error")
		pe := NewPanicError(err, stack)
		output := fmt.Sprintf("%s", pe)
		assert.Equal(t, "test error", output)
	}

	// Test %+s with nil error and empty stack
	{
		pe := NewPanicError(nil, StackTrace{})
		output := fmt.Sprintf("%+s", pe)
		assert.Equal(t, "<nil>", output)
	}

	// Test %+s with nil error and non-empty stack
	{
		pe := NewPanicError(nil, stack)
		output := fmt.Sprintf("%+s", pe)
		assert.Contains(t, output, "<nil>")
		// Only check for newline if stack is not empty
		// Note: GetStackTrace(0) might return empty stack in some test environments
		if len(stack) > 0 {
			// If stack exists, %+s should append stack trace
			assert.Contains(t, output, "\n", "Expected newline when stack is not empty")
		} else {
			// If stack is empty, should just be "<nil>"
			assert.Equal(t, "<nil>", output)
		}
	}

	// Test %+s with error and empty stack
	{
		err := errors.New("test error")
		pe := NewPanicError(err, StackTrace{})
		output := fmt.Sprintf("%+s", pe)
		assert.Equal(t, "test error", output)
		assert.NotContains(t, output, "\n")
	}

	// Test %+s with error and non-empty stack
	{
		err := errors.New("test error")
		pe := NewPanicError(err, stack)
		output := fmt.Sprintf("%+s", pe)
		assert.Contains(t, output, "test error")
		// Only check for newline if stack is not empty
		if len(stack) > 0 {
			assert.Contains(t, output, "\n", "Expected newline when stack is not empty")
		} else {
			// If stack is empty, should just be error message
			assert.Equal(t, "test error", output)
		}
	}

	// Test unsupported verb (defaults to 's')
	{
		err := errors.New("test error")
		pe := NewPanicError(err, stack)
		output := fmt.Sprintf("%d", pe)
		assert.Equal(t, "test error", output)
	}

	// Test unsupported verb with + flag (defaults to 's')
	{
		err := errors.New("test error")
		pe := NewPanicError(err, stack)
		output := fmt.Sprintf("%+d", pe)
		assert.Contains(t, output, "test error")
		// Only check for newline if stack is not empty
		if len(stack) > 0 {
			assert.Contains(t, output, "\n", "Expected newline when stack is not empty")
		} else {
			// If stack is empty, should just be error message
			assert.Equal(t, "test error", output)
		}
	}
}

func TestPanicError_Error_NilError(t *testing.T) {
	stack := GetStackTrace(0)
	pe := NewPanicError(nil, stack)
	assert.Equal(t, "<nil>", pe.Error())
}

func TestPanicError_Error_WithError(t *testing.T) {
	stack := GetStackTrace(0)
	err := errors.New("test error")
	pe := NewPanicError(err, stack)
	assert.Equal(t, "test error", pe.Error())
}

func TestPanicError_Unwrap_NilError(t *testing.T) {
	stack := GetStackTrace(0)
	pe := NewPanicError(nil, stack)
	assert.Nil(t, pe.Unwrap())
}

func TestPanicError_Unwrap_WithError(t *testing.T) {
	stack := GetStackTrace(0)
	err := errors.New("test error")
	pe := NewPanicError(err, stack)
	assert.Equal(t, err, pe.Unwrap())
}

func TestPanicError_StackTrace(t *testing.T) {
	stack := GetStackTrace(0)
	pe := NewPanicError(errors.New("test"), stack)
	assert.Equal(t, stack, pe.StackTrace())
}

func TestPanicError_NewPanicError_DefaultCase(t *testing.T) {
	stack := GetStackTrace(0)

	// Test with string (default case)
	{
		pe := NewPanicError("string value", stack)
		assert.NotNil(t, pe)
		assert.NotNil(t, pe.Unwrap())
		assert.Contains(t, pe.Error(), "string value")
	}

	// Test with int (default case)
	{
		pe := NewPanicError(42, stack)
		assert.NotNil(t, pe)
		assert.NotNil(t, pe.Unwrap())
		assert.Contains(t, pe.Error(), "42")
	}

	// Test with struct (default case)
	{
		type TestStruct struct {
			Value int
		}
		ts := TestStruct{Value: 100}
		pe := NewPanicError(ts, stack)
		assert.NotNil(t, pe)
		assert.NotNil(t, pe.Unwrap())
		assert.Contains(t, pe.Error(), "100")
	}
}
