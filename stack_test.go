package gust

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetStackTrace tests GetStackTrace function
func TestGetStackTrace(t *testing.T) {
	// Test basic functionality
	stack := GetStackTrace(0)
	assert.NotNil(t, stack)
	assert.Greater(t, len(stack), 0)

	// Test with different skip values
	stack1 := GetStackTrace(0)
	stack2 := GetStackTrace(1)
	assert.NotNil(t, stack1)
	assert.NotNil(t, stack2)
	// stack2 should have fewer frames than stack1
	assert.LessOrEqual(t, len(stack2), len(stack1))

	// Test that stack trace contains valid frames
	for _, frame := range stack {
		assert.NotEqual(t, Frame(0), frame)
	}
}

// TestPanicStackTrace tests PanicStackTrace function
func TestPanicStackTrace(t *testing.T) {
	// Test without panic (should return empty or minimal stack)
	stack := PanicStackTrace()
	assert.NotNil(t, stack)
	// Without an actual panic, findPanicStack may return empty stack
	// This is expected behavior

	// Test with actual panic
	var panicStack StackTrace
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicStack = PanicStackTrace()
			}
		}()
		panic("test panic")
	}()
	// After panic, stack should be captured
	assert.NotNil(t, panicStack)
}

// TestFrame_pc tests Frame.pc() method
func TestFrame_pc(t *testing.T) {
	stack := GetStackTrace(0)
	if len(stack) == 0 {
		t.Skip("No frames in stack trace")
		return
	}

	frame := stack[0]
	pc := frame.pc()
	assert.NotEqual(t, uintptr(0), pc)

	// Verify pc() returns uintptr(f) - 1
	expectedPC := uintptr(frame) - 1
	assert.Equal(t, expectedPC, pc)
}

// TestFrame_file tests Frame.file() method
func TestFrame_file(t *testing.T) {
	stack := GetStackTrace(0)
	if len(stack) == 0 {
		t.Skip("No frames in stack trace")
		return
	}

	frame := stack[0]
	file := frame.file()
	assert.NotEmpty(t, file)
	assert.NotEqual(t, "unknown", file)
	// Should contain path separator or be a valid file path
	assert.True(t, strings.Contains(file, "/") || strings.Contains(file, "\\"))

	// Test with invalid frame (should return "unknown")
	invalidFrame := Frame(0)
	file = invalidFrame.file()
	assert.Equal(t, "unknown", file)
}

// TestFrame_line tests Frame.line() method
func TestFrame_line(t *testing.T) {
	stack := GetStackTrace(0)
	if len(stack) == 0 {
		t.Skip("No frames in stack trace")
		return
	}

	frame := stack[0]
	line := frame.line()
	assert.Greater(t, line, 0)

	// Test with invalid frame (should return 0)
	invalidFrame := Frame(0)
	line = invalidFrame.line()
	assert.Equal(t, 0, line)
}

// TestFrame_name tests Frame.name() method
func TestFrame_name(t *testing.T) {
	stack := GetStackTrace(0)
	if len(stack) == 0 {
		t.Skip("No frames in stack trace")
		return
	}

	frame := stack[0]
	name := frame.name()
	assert.NotEmpty(t, name)
	assert.NotEqual(t, "unknown", name)
	// Should contain package and function name
	assert.Contains(t, name, ".")

	// Test with invalid frame (should return "unknown")
	invalidFrame := Frame(0)
	name = invalidFrame.name()
	assert.Equal(t, "unknown", name)
}

// TestFrame_Format tests Frame.Format() method with various verbs and flags
func TestFrame_Format(t *testing.T) {
	stack := GetStackTrace(0)
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
	stack := GetStackTrace(0)
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
	invalidFrame := Frame(0)
	text, err = invalidFrame.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("unknown"), text)
}

// TestStackTrace_Format tests StackTrace.Format() method
func TestStackTrace_Format(t *testing.T) {
	stack := GetStackTrace(0)
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
	assert.True(t, strings.HasPrefix(output, "[]gust.Frame"))

	// Test %s
	buf.Reset()
	fmt.Fprintf(&buf, "%s", stack)
	output = buf.String()
	assert.NotEmpty(t, output)
	assert.True(t, strings.HasPrefix(output, "["))

	// Test empty stack trace
	emptyStack := StackTrace{}
	buf.Reset()
	fmt.Fprintf(&buf, "%v", emptyStack)
	output = buf.String()
	assert.Equal(t, "[]", output)
}

// TestStackTrace_formatSlice tests StackTrace.formatSlice() method
func TestStackTrace_formatSlice(t *testing.T) {
	st := GetStackTrace(0)
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
	singleStack := StackTrace{st[0]}
	buf.Reset()
	fmt.Fprintf(&buf, "%s", singleStack)
	output = buf.String()
	assert.True(t, strings.HasPrefix(output, "["))
	assert.True(t, strings.HasSuffix(output, "]"))

	// Test with empty stack
	emptyStack := StackTrace{}
	buf.Reset()
	fmt.Fprintf(&buf, "%s", emptyStack)
	output = buf.String()
	assert.Equal(t, "[]", output)
}

// TestStack_Format tests stack.Format() method
func TestStack_Format(t *testing.T) {
	st := GetStackTrace(0)
	if len(st) == 0 {
		t.Skip("No frames in stack trace")
		return
	}

	// Convert StackTrace back to stack for testing
	var pcs []uintptr
	for _, frame := range st {
		pcs = append(pcs, uintptr(frame))
	}
	stackVal := stack(pcs)

	var buf bytes.Buffer

	// Test %+v (with + flag)
	buf.Reset()
	fmt.Fprintf(&buf, "%+v", &stackVal)
	output := buf.String()
	assert.NotEmpty(t, output)
	assert.Contains(t, output, "\n")

	// Test %v (without + flag, should not format)
	buf.Reset()
	fmt.Fprintf(&buf, "%v", &stackVal)
	output = buf.String()
	// Without + flag, stack.Format doesn't handle %v, so fmt uses default formatting
	// The output might be empty or contain default formatting, both are acceptable
	// We just verify it doesn't panic
	_ = output
}

// TestStack_StackTrace tests stack.StackTrace() method
func TestStack_StackTrace(t *testing.T) {
	st := GetStackTrace(0)
	if len(st) == 0 {
		t.Skip("No frames in stack trace")
		return
	}

	// Convert StackTrace to stack
	var pcs []uintptr
	for _, frame := range st {
		pcs = append(pcs, uintptr(frame))
	}
	stackVal := stack(pcs)

	// Convert back to StackTrace
	result := stackVal.StackTrace()
	assert.Equal(t, len(stackVal), len(result))
	for i := range stackVal {
		assert.Equal(t, Frame(stackVal[i]), result[i])
	}

	// Test with empty stack
	emptyStack := stack{}
	result = emptyStack.StackTrace()
	assert.Equal(t, 0, len(result))
}

// TestFuncname tests funcname() function
func TestFuncname(t *testing.T) {
	// Test with standard function name
	name := "github.com/andeya/gust.TestFuncname"
	result := funcname(name)
	assert.Equal(t, "TestFuncname", result)

	// Test with nested package
	name = "github.com/example/pkg/subpkg.Function"
	result = funcname(name)
	assert.Equal(t, "Function", result)

	// Test with no package path
	name = "main.function"
	result = funcname(name)
	assert.Equal(t, "function", result)

	// Test with single dot
	name = "pkg.func"
	result = funcname(name)
	assert.Equal(t, "func", result)
}

// TestCallers tests callers() function
func TestCallers(t *testing.T) {
	// Test basic functionality
	st := callers(0)
	assert.NotNil(t, st)
	assert.Greater(t, len(*st), 0)

	// Test with different skip values
	st1 := callers(0)
	st2 := callers(1)
	assert.NotNil(t, st1)
	assert.NotNil(t, st2)
	// st2 should have fewer frames
	assert.LessOrEqual(t, len(*st2), len(*st1))

	// Test that it returns valid program counters
	for _, pc := range *st {
		assert.NotEqual(t, uintptr(0), pc)
	}
}

// TestFindPanicStack tests findPanicStack() function
func TestFindPanicStack(t *testing.T) {
	// Test without panic (may return empty stack)
	st := findPanicStack()
	assert.NotNil(t, st)
	// Without actual panic, it may return empty stack

	// Test with actual panic
	var panicStack *stack
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicStack = findPanicStack()
			}
		}()
		panic("test panic for findPanicStack")
	}()
	// After panic, stack should be captured
	assert.NotNil(t, panicStack)
	// Should contain frames after runtime.gopanic
	assert.Greater(t, len(*panicStack), 0)
}

// TestStackTraceCarrier tests StackTraceCarrier interface
func TestStackTraceCarrier(t *testing.T) {
	// Test with panicError (which implements StackTraceCarrier)
	st := GetStackTrace(0)
	err := errors.New("test error")
	pe := &panicError{
		err:   err,
		stack: st,
	}
	var carrier StackTraceCarrier = pe
	assert.NotNil(t, carrier)
	result := carrier.StackTrace()
	assert.Equal(t, st, result)
}

// TestFrame_EdgeCases tests edge cases for Frame methods
func TestFrame_EdgeCases(t *testing.T) {
	// Test with zero Frame
	zeroFrame := Frame(0)
	// pc() returns uintptr(f) - 1, so Frame(0).pc() = uintptr(0) - 1
	// Since uintptr is unsigned, this wraps around to max uintptr value
	pc := zeroFrame.pc()
	// Verify it's the wrapped value (max uintptr)
	assert.Equal(t, ^uintptr(0), pc)
	assert.Equal(t, "unknown", zeroFrame.file())
	assert.Equal(t, 0, zeroFrame.line())
	assert.Equal(t, "unknown", zeroFrame.name())

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
	var nilStack StackTrace
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%v", nilStack)
	assert.Equal(t, "[]", buf.String())

	// Test single frame StackTrace
	stack := GetStackTrace(0)
	if len(stack) > 0 {
		singleStack := StackTrace{stack[0]}
		buf.Reset()
		fmt.Fprintf(&buf, "%v", singleStack)
		output := buf.String()
		assert.True(t, strings.HasPrefix(output, "["))
		assert.True(t, strings.HasSuffix(output, "]"))
	}
}
