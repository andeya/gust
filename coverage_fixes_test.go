package gust_test

import (
	"testing"

	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

// TestChunkResult_SafeGetMethods tests ChunkResult's safeGetT and safeGetE methods
// (covers iter/chunk_result.go:43-49, 52-58)
func TestChunkResult_SafeGetMethods(t *testing.T) {
	// Test ChunkOk's safeGetT (created indirectly via NextChunk)
	iter1 := iter.FromSlice([]int{1, 2, 3})
	ok := iter1.NextChunk(3)
	// safeGetT is internal, we test it indirectly via Unwrap
	assert.Equal(t, []int{1, 2, 3}, ok.Unwrap())

	// Test ChunkErr's safeGetE (partial chunk via NextChunk)
	iter2 := iter.FromSlice([]int{4, 5})
	err := iter2.NextChunk(3) // Only 2 elements available, should return error
	assert.True(t, err.IsErr())
	// safeGetE is internal, we test it indirectly via UnwrapErr
	assert.Equal(t, []int{4, 5}, err.UnwrapErr())

	// Test empty chunk
	iter3 := iter.FromSlice([]int{})
	emptyOk := iter3.NextChunk(0)
	assert.Equal(t, []int{}, emptyOk.Unwrap())

	// Test Unwrap's IsErr branch (should panic)
	iter4 := iter.FromSlice([]int{1, 2})
	errResult := iter4.NextChunk(3)
	assert.Panics(t, func() {
		errResult.Unwrap()
	})

	// Test UnwrapErr's !IsErr branch (should panic)
	iter5 := iter.FromSlice([]int{1, 2})
	okResult := iter5.NextChunk(2)
	assert.Panics(t, func() {
		okResult.UnwrapErr()
	})
}
