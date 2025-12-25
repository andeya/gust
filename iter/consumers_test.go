package iter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextChunk(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3, 4, 5})
	chunk := iter.NextChunk(2)
	assert.True(t, chunk.IsOk())
	assert.Equal(t, []int{1, 2}, chunk.Unwrap())

	// Next chunk should fail because only 3 elements remain, but we request 4
	chunk2 := iter.NextChunk(4)
	assert.True(t, chunk2.IsErr())
	// The error contains the remaining elements
	assert.Equal(t, []int{3, 4, 5}, chunk2.UnwrapErr())

	// After consuming all elements, next chunk should fail
	chunk3 := iter.NextChunk(2)
	assert.True(t, chunk3.IsErr())
}
