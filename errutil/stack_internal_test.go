package errutil

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStack_Format_Direct tests stack.Format() method directly - covers errbox.go:163-174
func TestStack_Format_Direct(t *testing.T) {
	// Create a stack using runtime.Callers (similar to callers() function)
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(0, pcs[:])

	if n == 0 {
		t.Skip("No program counters available")
		return
	}

	// Create a stack instance
	var st stack = pcs[0:n]

	// Test with %v (without + flag) - should not output anything (no case matches)
	var buf1 bytes.Buffer
	fmt.Fprintf(&buf1, "%v", &st)
	output1 := buf1.String()
	// When verb is 'v' but no '+' flag, the switch case doesn't match, so no output
	assert.Empty(t, output1)

	// Test with %+v (with + flag) - covers errbox.go:163-174
	var buf2 bytes.Buffer
	fmt.Fprintf(&buf2, "%+v", &st)
	output2 := buf2.String()
	// Should contain newlines (one per frame)
	assert.NotEmpty(t, output2)
	assert.Contains(t, output2, "\n")
	// Should contain frame information (function names, file paths)
	assert.True(t, len(output2) > 0)

	// Count newlines to verify each frame is formatted
	// Each frame produces 2 newlines: one from "\n" prefix and one from Frame's "\n\t" in %+v format
	newlineCount := strings.Count(output2, "\n")
	assert.Greater(t, newlineCount, 0, "Should have at least one newline per frame")
	assert.GreaterOrEqual(t, newlineCount, n, "Should have at least one newline per frame")
	assert.LessOrEqual(t, newlineCount, n*2, "Should have at most two newlines per frame")

	// Test with other verbs (should not output anything)
	var buf3 bytes.Buffer
	fmt.Fprintf(&buf3, "%s", &st)
	output3 := buf3.String()
	assert.Empty(t, output3, "Non-'v' verbs should produce no output")

	var buf4 bytes.Buffer
	fmt.Fprintf(&buf4, "%d", &st)
	output4 := buf4.String()
	assert.Empty(t, output4, "Non-'v' verbs should produce no output")
}

// TestStack_Format_EmptyStack tests stack.Format() with empty stack
func TestStack_Format_EmptyStack(t *testing.T) {
	// Create an empty stack
	var st stack = []uintptr{}

	// Test with %+v (with + flag)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%+v", &st)
	output := buf.String()
	// Empty stack should produce no output (loop doesn't execute)
	assert.Empty(t, output)
}

// TestStack_Format_SingleFrame tests stack.Format() with single frame
func TestStack_Format_SingleFrame(t *testing.T) {
	// Create a stack with a single frame
	var pcs [1]uintptr
	n := runtime.Callers(0, pcs[:])

	if n == 0 {
		t.Skip("No program counters available")
		return
	}

	var st stack = pcs[0:1]

	// Test with %+v (with + flag)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%+v", &st)
	output := buf.String()
	// Each frame produces 2 newlines: one from "\n" prefix and one from Frame's "\n\t" in %+v format
	assert.NotEmpty(t, output)
	newlineCount := strings.Count(output, "\n")
	assert.Equal(t, 2, newlineCount, "Single frame should produce two newlines (one from prefix, one from Frame format)")
}

// TestStack_Format_MultipleFrames tests stack.Format() with multiple frames
func TestStack_Format_MultipleFrames(t *testing.T) {
	// Create a stack with multiple frames
	const depth = 10
	var pcs [depth]uintptr
	n := runtime.Callers(0, pcs[:])

	if n == 0 {
		t.Skip("No program counters available")
		return
	}

	// Use first few frames
	frameCount := n
	if frameCount > 5 {
		frameCount = 5
	}
	var st stack = pcs[0:frameCount]

	// Test with %+v (with + flag)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%+v", &st)
	output := buf.String()
	// Each frame produces 2 newlines: one from "\n" prefix and one from Frame's "\n\t" in %+v format
	assert.NotEmpty(t, output)
	newlineCount := strings.Count(output, "\n")
	assert.Equal(t, frameCount*2, newlineCount, "Should have two newlines per frame")

	// Verify each frame is formatted with %+v (should contain function names and file paths)
	// Split by "\n" to get all lines (including empty lines from the prefix)
	lines := strings.Split(output, "\n")
	// Filter out empty lines to get actual frame content lines
	var contentLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			contentLines = append(contentLines, line)
		}
	}
	// Each frame produces: function name, file path (with tab), and line number
	// So we should have at least frameCount content lines (function names)
	assert.GreaterOrEqual(t, len(contentLines), frameCount, "Should have at least one content line per frame")

	// Verify output contains frame information (function names, file paths)
	assert.Contains(t, output, "\t", "Should contain tab characters from Frame format")
}

// TestStack_Format_VerbCases tests stack.Format() with different verbs
func TestStack_Format_VerbCases(t *testing.T) {
	// Create a stack
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(0, pcs[:])

	if n == 0 {
		t.Skip("No program counters available")
		return
	}

	var st stack = pcs[0:n]

	// Test with 'v' verb and '+' flag - should output
	var buf1 bytes.Buffer
	fmt.Fprintf(&buf1, "%+v", &st)
	output1 := buf1.String()
	assert.NotEmpty(t, output1, "%+v should produce output")

	// Test with 'v' verb without '+' flag - should not output (no case matches)
	var buf2 bytes.Buffer
	fmt.Fprintf(&buf2, "%v", &st)
	output2 := buf2.String()
	assert.Empty(t, output2, "%v without + flag should produce no output")

	// Test with other verbs - should not output
	testCases := []string{"%s", "%d", "%x", "%#v", "%q"}
	for _, verb := range testCases {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, verb, &st)
		output := buf.String()
		assert.Empty(t, output, "Verb %s should produce no output", verb)
	}
}
