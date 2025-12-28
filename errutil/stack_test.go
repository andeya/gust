package errutil_test

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/andeya/gust/errutil"
	"github.com/stretchr/testify/assert"
)

// TestGetStackTrace tests GetStackTrace function
func TestGetStackTrace(t *testing.T) {
	// Test basic functionality
	stack := errutil.GetStackTrace(0)
	assert.NotNil(t, stack)
	assert.Greater(t, len(stack), 0)

	// Test with different skip values
	stack1 := errutil.GetStackTrace(0)
	stack2 := errutil.GetStackTrace(1)
	assert.NotNil(t, stack1)
	assert.NotNil(t, stack2)
	// stack2 should have fewer frames than stack1
	assert.LessOrEqual(t, len(stack2), len(stack1))

	// Test that stack trace contains valid frames
	for _, frame := range stack {
		assert.NotEqual(t, errutil.Frame(0), frame)
	}
}

// TestPanicStackTrace tests PanicStackTrace function
func TestPanicStackTrace(t *testing.T) {
	// Test without panic (should return empty or minimal stack)
	stack := errutil.PanicStackTrace()
	assert.NotNil(t, stack)
	// Without an actual panic, findPanicStack may return empty stack
	// This is expected behavior

	// Test with actual panic
	var panicStack errutil.StackTrace
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicStack = errutil.PanicStackTrace()
			}
		}()
		panic("test panic")
	}()
	// After panic, stack should be captured
	assert.NotNil(t, panicStack)
}

// TestFrame_pc tests Frame.pc() method
// Note: pc() is not exported, so we can't test it directly
func TestFrame_pc(t *testing.T) {
	stack := errutil.GetStackTrace(0)
	if len(stack) == 0 {
		t.Skip("No frames in stack trace")
		return
	}

	frame := stack[0]
	// Note: pc() is not exported, so we can't test it directly
	// We'll test through other exported methods
	// Verify frame is valid by checking it's not zero
	assert.NotEqual(t, errutil.Frame(0), frame)
}

// TestFrame_file tests Frame.file() method
// Note: file() is not exported, so we test through Format()
func TestFrame_file(t *testing.T) {
	stack := errutil.GetStackTrace(0)
	if len(stack) == 0 {
		t.Skip("No frames in stack trace")
		return
	}

	frame := stack[0]
	// Test through Format() method
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s", frame)
	output := buf.String()
	assert.NotEmpty(t, output)
	assert.NotEqual(t, "unknown", output)

	// Test with invalid frame (should return "unknown")
	invalidFrame := errutil.Frame(0)
	buf.Reset()
	fmt.Fprintf(&buf, "%s", invalidFrame)
	assert.Equal(t, "unknown", buf.String())
}

// TestFrame_line tests Frame.line() method
// Note: line() is not exported, so we test through Format()
func TestFrame_line(t *testing.T) {
	stack := errutil.GetStackTrace(0)
	if len(stack) == 0 {
		t.Skip("No frames in stack trace")
		return
	}

	frame := stack[0]
	// Test through Format() method
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d", frame)
	output := buf.String()
	assert.NotEmpty(t, output)
	// Should be a number
	assert.True(t, strings.TrimSpace(output) != "")

	// Test with invalid frame (should return 0)
	invalidFrame := errutil.Frame(0)
	buf.Reset()
	fmt.Fprintf(&buf, "%d", invalidFrame)
	assert.Equal(t, "0", buf.String())
}

// TestFrame_name tests Frame.name() method
// Note: name() is not exported, so we test through Format()
func TestFrame_name(t *testing.T) {
	stack := errutil.GetStackTrace(0)
	if len(stack) == 0 {
		t.Skip("No frames in stack trace")
		return
	}

	frame := stack[0]
	// Test through Format() method
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%n", frame)
	output := buf.String()
	assert.NotEmpty(t, output)
	assert.NotEqual(t, "unknown", output)
	// Should not contain package path
	assert.False(t, strings.Contains(output, "/"))

	// Test with invalid frame (should return "unknown")
	invalidFrame := errutil.Frame(0)
	buf.Reset()
	fmt.Fprintf(&buf, "%n", invalidFrame)
	assert.Equal(t, "unknown", buf.String())
}

// TestFrame_Format tests Frame.Format() method with various verbs and flags
func TestFrame_Format(t *testing.T) {
	stack := errutil.GetStackTrace(0)
	if len(stack) == 0 {
		t.Skip("No frames in stack trace")
		return
	}

	frame := stack[0]
	var buf bytes.Buffer

	// Test %s (without + flag)
	buf.Reset()
	fmt.Fprintf(&buf, "%s", frame)
	output := buf.String()
	assert.NotEmpty(t, output)
	// Should be just the base filename
	assert.False(t, strings.Contains(output, "\n"))

	// Test %+s (with + flag)
	buf.Reset()
	fmt.Fprintf(&buf, "%+s", frame)
	output = buf.String()
	assert.NotEmpty(t, output)
	// Should contain function name and file path with newline
	assert.Contains(t, output, "\n")
	assert.Contains(t, output, "\t")

	// Test %d (line number)
	buf.Reset()
	fmt.Fprintf(&buf, "%d", frame)
	output = buf.String()
	assert.NotEmpty(t, output)
	// Should be a number
	assert.True(t, strings.TrimSpace(output) != "")

	// Test %n (function name)
	buf.Reset()
	fmt.Fprintf(&buf, "%n", frame)
	output = buf.String()
	assert.NotEmpty(t, output)
	// Should not contain package path
	assert.False(t, strings.Contains(output, "/"))

	// Test %v (equivalent to %s:%d)
	buf.Reset()
	fmt.Fprintf(&buf, "%v", frame)
	output = buf.String()
	assert.NotEmpty(t, output)
	assert.Contains(t, output, ":")

	// Test %+v (equivalent to %+s:%d)
	buf.Reset()
	fmt.Fprintf(&buf, "%+v", frame)
	output = buf.String()
	assert.NotEmpty(t, output)
	assert.Contains(t, output, "\n")
	assert.Contains(t, output, ":")
}

// TestFrame_MarshalText tests Frame.MarshalText() method
func TestFrame_MarshalText(t *testing.T) {
	stack := errutil.GetStackTrace(0)
	if len(stack) == 0 {
		t.Skip("No frames in stack trace")
		return
	}

	frame := stack[0]
	text, err := frame.MarshalText()
	assert.NoError(t, err)
	assert.NotNil(t, text)
	assert.Greater(t, len(text), 0)
	// Should contain function name, file, and line number
	textStr := string(text)
	assert.Contains(t, textStr, " ")
	assert.Contains(t, textStr, ":")

	// Test with invalid frame (should return "unknown")
	invalidFrame := errutil.Frame(0)
	text, err = invalidFrame.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("unknown"), text)
}

// TestStackTrace_Format tests StackTrace.Format() method
func TestStackTrace_Format(t *testing.T) {
	stack := errutil.GetStackTrace(0)
	if len(stack) == 0 {
		t.Skip("No frames in stack trace")
		return
	}

	var buf bytes.Buffer

	// Test %v (without flags)
	buf.Reset()
	fmt.Fprintf(&buf, "%v", stack)
	output := buf.String()
	assert.NotEmpty(t, output)
	assert.True(t, strings.HasPrefix(output, "["))
	assert.True(t, strings.HasSuffix(output, "]"))

	// Test %+v (with + flag)
	buf.Reset()
	fmt.Fprintf(&buf, "%+v", stack)
	output = buf.String()
	assert.NotEmpty(t, output)
	// Should contain newlines for each frame
	assert.Contains(t, output, "\n")

	// Test %#v (with # flag)
	buf.Reset()
	fmt.Fprintf(&buf, "%#v", stack)
	output = buf.String()
	assert.NotEmpty(t, output)
	// Should be Go-syntax representation
	assert.True(t, strings.HasPrefix(output, "[]errutil.Frame"))

	// Test %s
	buf.Reset()
	fmt.Fprintf(&buf, "%s", stack)
	output = buf.String()
	assert.NotEmpty(t, output)
	assert.True(t, strings.HasPrefix(output, "["))

	// Test empty stack trace
	emptyStack := errutil.StackTrace{}
	buf.Reset()
	fmt.Fprintf(&buf, "%v", emptyStack)
	output = buf.String()
	assert.Equal(t, "[]", output)

	// Test %d (line numbers) - covers default case in StackTrace.Format()
	buf.Reset()
	fmt.Fprintf(&buf, "%d", stack)
	output = buf.String()
	assert.NotEmpty(t, output)
	assert.True(t, strings.HasPrefix(output, "["))
	assert.True(t, strings.HasSuffix(output, "]"))

	// Test %n (function names) - covers default case in StackTrace.Format()
	buf.Reset()
	fmt.Fprintf(&buf, "%n", stack)
	output = buf.String()
	assert.NotEmpty(t, output)
	assert.True(t, strings.HasPrefix(output, "["))
	assert.True(t, strings.HasSuffix(output, "]"))
}

// TestStackTrace_formatSlice tests StackTrace.formatSlice() method
func TestStackTrace_formatSlice(t *testing.T) {
	st := errutil.GetStackTrace(0)
	if len(st) == 0 {
		t.Skip("No frames in stack trace")
		return
	}

	// Test with %s (using fmt.Fprintf which provides fmt.State)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s", st)
	output := buf.String()
	assert.True(t, strings.HasPrefix(output, "["))
	assert.True(t, strings.HasSuffix(output, "]"))

	// Test with %v
	buf.Reset()
	fmt.Fprintf(&buf, "%v", st)
	output = buf.String()
	assert.True(t, strings.HasPrefix(output, "["))
	assert.True(t, strings.HasSuffix(output, "]"))

	// Test with single frame
	singleStack := errutil.StackTrace{st[0]}
	buf.Reset()
	fmt.Fprintf(&buf, "%s", singleStack)
	output = buf.String()
	assert.True(t, strings.HasPrefix(output, "["))
	assert.True(t, strings.HasSuffix(output, "]"))

	// Test with empty stack
	emptyStack := errutil.StackTrace{}
	buf.Reset()
	fmt.Fprintf(&buf, "%s", emptyStack)
	output = buf.String()
	assert.Equal(t, "[]", output)
}

// TestStackTraceCarrier tests StackTraceCarrier interface
func TestStackTraceCarrier(t *testing.T) {
	// Test with PanicError (which implements StackTraceCarrier)
	st := errutil.GetStackTrace(0)
	err := errors.New("test error")
	pe := errutil.NewPanicError(err, st)
	var carrier errutil.StackTraceCarrier = pe
	assert.NotNil(t, carrier)
	result := carrier.StackTrace()
	assert.Equal(t, st, result)
}

// TestFrame_EdgeCases tests edge cases for Frame methods
func TestFrame_EdgeCases(t *testing.T) {
	// Test with zero Frame
	zeroFrame := errutil.Frame(0)
	// Note: pc(), file(), line(), name() are not exported, so we test through Format()
	assert.Equal(t, "unknown", fmt.Sprintf("%s", zeroFrame))
	assert.Equal(t, "0", fmt.Sprintf("%d", zeroFrame))
	assert.Equal(t, "unknown", fmt.Sprintf("%n", zeroFrame))

	// Test MarshalText with zero frame
	text, err := zeroFrame.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("unknown"), text)

	// Test Format with zero frame
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s", zeroFrame)
	assert.Equal(t, "unknown", buf.String())
}

// TestStackTrace_EdgeCases tests edge cases for StackTrace
func TestStackTrace_EdgeCases(t *testing.T) {
	// Test nil StackTrace (should not panic)
	var nilStack errutil.StackTrace
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%v", nilStack)
	assert.Equal(t, "[]", buf.String())

	// Test single frame StackTrace
	stack := errutil.GetStackTrace(0)
	if len(stack) > 0 {
		singleStack := errutil.StackTrace{stack[0]}
		buf.Reset()
		fmt.Fprintf(&buf, "%v", singleStack)
		output := buf.String()
		assert.True(t, strings.HasPrefix(output, "["))
		assert.True(t, strings.HasSuffix(output, "]"))
	}
}

// TestStack_Format tests stack.Format() method
// Since stack is unexported, we test it indirectly through GetStackTrace
// and verify that stack.Format() behavior is correct
func TestStack_Format(t *testing.T) {
	// Create a stack using runtime.Callers (similar to callers() function)
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(0, pcs[:])

	// Verify that we have valid program counters
	assert.Greater(t, n, 0, "Should have at least one program counter")

	// Test through GetStackTrace which uses stack internally
	// and verify the output format matches what stack.Format() would produce
	st := errutil.GetStackTrace(0)
	if len(st) == 0 {
		t.Skip("No frames in stack trace")
		return
	}

	// Test that stack.Format() would be called when formatting with %+v
	// We can't directly test it, but we can verify the behavior is correct
	// by checking that GetStackTrace returns valid frames

	// Since we can't directly test unexported stack.Format(),
	// we'll verify that the stack trace functionality works correctly
	// which indirectly tests stack.Format()
	assert.Greater(t, len(st), 0, "Stack trace should contain frames")

	// Test formatting the stack trace with %+v to verify Format() behavior
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%+v", st)
	output := buf.String()
	assert.NotEmpty(t, output)
	assert.Contains(t, output, "\n")
}
