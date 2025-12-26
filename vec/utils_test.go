package vec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFixIndex tests fixIndex function (covers utils.go:26-41)
func TestFixIndex(t *testing.T) {
	// Test idx < 0 && idx < 0 after adjustment (covers utils.go:29-31)
	length := 5
	idx := -10 // length + (-10) = -5, which is < 0
	result := fixIndex(length, idx, false)
	assert.Equal(t, 0, result)

	// Test idx < 0 && idx >= 0 after adjustment
	idx2 := -3 // length + (-3) = 2, which is >= 0
	result2 := fixIndex(length, idx2, false)
	assert.Equal(t, 2, result2)

	// Test idx >= length && canLen = true
	idx3 := 10
	result3 := fixIndex(length, idx3, true)
	assert.Equal(t, length, result3)

	// Test idx >= length && canLen = false
	idx4 := 10
	result4 := fixIndex(length, idx4, false)
	assert.Equal(t, length-1, result4)

	// Test normal case (0 <= idx < length)
	idx5 := 2
	result5 := fixIndex(length, idx5, false)
	assert.Equal(t, 2, result5)
}
