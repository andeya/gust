package iterator_test

import (
	"testing"

	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/option"
	"github.com/stretchr/testify/assert"
)

// TestArrayChunksPanic tests iterator.ArrayChunks panic on zero chunk size
func TestArrayChunksPanic(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	assert.Panics(t, func() {
		iterator.ArrayChunks(iter, 0)
	})
}

// TestArrayChunksEmptyBuffer tests iterator.ArrayChunks with empty buffer
func TestArrayChunksEmptyBuffer(t *testing.T) {
	iter := iterator.ArrayChunks(iterator.Empty[int](), 2)
	assert.Equal(t, option.None[[]int](), iter.Next())
}

// TestChunkBySingleElement tests iterator.ChunkBy with single element
func TestChunkBySingleElement(t *testing.T) {
	iter := iterator.ChunkBy(iterator.FromSlice([]int{1}), func(a, b int) bool { return a == b })
	chunk := iter.Next()
	assert.True(t, chunk.IsSome())
	assert.Equal(t, []int{1}, chunk.Unwrap())
	assert.Equal(t, option.None[[]int](), iter.Next())
}

// TestMapWindowsEmpty tests iterator.MapWindows with empty iterator
func TestMapWindowsEmpty(t *testing.T) {
	iter := iterator.MapWindows(iterator.Empty[int](), 3, func(window []int) int {
		return len(window)
	})
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestMapWindowsPanic tests iterator.MapWindows panic on zero window size
func TestMapWindowsPanic(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	assert.Panics(t, func() {
		iterator.MapWindows(iter, 0, func(window []int) int { return len(window) })
	})
}

// TestArrayChunksSizeHintEdgeCases tests iterator.ArrayChunks SizeHint edge cases
func TestArrayChunksSizeHintEdgeCases(t *testing.T) {
	// Test with lower == 0
	iter := iterator.Empty[int]()
	chunks := iterator.ArrayChunks(iter, 2)
	lower, upper := chunks.SizeHint()
	assert.Equal(t, uint(0), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(0), upper.Unwrap())

	// Test with upperVal == 0
	iter2 := iterator.Empty[int]()
	chunks2 := iterator.ArrayChunks(iter2, 2)
	lower2, upper2 := chunks2.SizeHint()
	assert.Equal(t, uint(0), lower2)
	assert.True(t, upper2.IsSome())
	assert.Equal(t, uint(0), upper2.Unwrap())
}

// TestMapWindowsSizeHintEdgeCases tests iterator.MapWindows SizeHint edge cases
func TestMapWindowsSizeHintEdgeCases(t *testing.T) {
	// Test with lower < windowSize
	iter := iterator.FromSlice([]int{1, 2})
	windows := iterator.MapWindows(iter, 3, func(window []int) int { return len(window) })
	lower, upper := windows.SizeHint()
	assert.Equal(t, uint(0), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(0), upper.Unwrap())

	// Test with upperVal < windowSize
	iter2 := iterator.FromSlice([]int{1, 2})
	windows2 := iterator.MapWindows(iter2, 3, func(window []int) int { return len(window) })
	lower2, upper2 := windows2.SizeHint()
	assert.Equal(t, uint(0), lower2)
	assert.True(t, upper2.IsSome())
	assert.Equal(t, uint(0), upper2.Unwrap())
}

// TestMapWindowsSizeHintLower tests iterator.MapWindows SizeHint with lower >= windowSize
func TestMapWindowsSizeHintLower(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	windows := iterator.MapWindows(iter, 3, func(window []int) int { return len(window) })
	lower, upper := windows.SizeHint()
	assert.Equal(t, uint(3), lower) // 5 - 3 + 1 = 3
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(3), upper.Unwrap())
}

// TestChunkByFirstEmpty tests iterator.ChunkBy when first element is None
func TestChunkByFirstEmpty(t *testing.T) {
	iter := iterator.Empty[int]()
	chunked := iterator.ChunkBy(iter, func(a, b int) bool { return a == b })
	assert.Equal(t, option.None[[]int](), chunked.Next())
}

// TestChunkByCurrentEmpty tests iterator.ChunkBy when current is empty after None
func TestChunkByCurrentEmpty(t *testing.T) {
	// This tests the len(c.current) == 0 case
	// This should not happen in normal usage, but we test it for coverage
	iter := iterator.FromSlice([]int{1})
	chunked := iterator.ChunkBy(iter, func(a, b int) bool { return a == b })
	chunk1 := chunked.Next()
	assert.True(t, chunk1.IsSome())
	assert.Equal(t, []int{1}, chunk1.Unwrap())
	// After consuming, next should be None
	assert.Equal(t, option.None[[]int](), chunked.Next())
}

func TestArrayChunks(t *testing.T) {
	iter := iterator.ArrayChunks(iterator.FromSlice([]int{1, 2, 3, 4, 5, 6}), 2)

	chunk1 := iter.Next()
	assert.True(t, chunk1.IsSome())
	assert.Equal(t, []int{1, 2}, chunk1.Unwrap())

	chunk2 := iter.Next()
	assert.True(t, chunk2.IsSome())
	assert.Equal(t, []int{3, 4}, chunk2.Unwrap())

	chunk3 := iter.Next()
	assert.True(t, chunk3.IsSome())
	assert.Equal(t, []int{5, 6}, chunk3.Unwrap())

	assert.True(t, iter.Next().IsNone())
}

func TestArrayChunksPartial(t *testing.T) {
	iter := iterator.ArrayChunks(iterator.FromSlice([]int{1, 2, 3, 4, 5}), 2)

	chunk1 := iter.Next()
	assert.True(t, chunk1.IsSome())
	assert.Equal(t, []int{1, 2}, chunk1.Unwrap())

	chunk2 := iter.Next()
	assert.True(t, chunk2.IsSome())
	assert.Equal(t, []int{3, 4}, chunk2.Unwrap())

	chunk3 := iter.Next()
	assert.True(t, chunk3.IsSome())
	assert.Equal(t, []int{5}, chunk3.Unwrap()) // Partial chunk

	assert.True(t, iter.Next().IsNone())
}

func TestArrayChunksEmpty(t *testing.T) {
	iter := iterator.ArrayChunks(iterator.Empty[int](), 2)
	assert.True(t, iter.Next().IsNone())
}

func TestChunkBy(t *testing.T) {
	iter := iterator.ChunkBy(iterator.FromSlice([]int{1, 1, 2, 2, 2, 3, 3}), func(a, b int) bool { return a == b })

	chunk1 := iter.Next()
	assert.True(t, chunk1.IsSome())
	assert.Equal(t, []int{1, 1}, chunk1.Unwrap())

	chunk2 := iter.Next()
	assert.True(t, chunk2.IsSome())
	assert.Equal(t, []int{2, 2, 2}, chunk2.Unwrap())

	chunk3 := iter.Next()
	assert.True(t, chunk3.IsSome())
	assert.Equal(t, []int{3, 3}, chunk3.Unwrap())

	assert.True(t, iter.Next().IsNone())
}

func TestChunkByEmpty(t *testing.T) {
	iter := iterator.ChunkBy(iterator.Empty[int](), func(a, b int) bool { return a == b })
	assert.True(t, iter.Next().IsNone())
}

func TestChunkBySingle(t *testing.T) {
	iter := iterator.ChunkBy(iterator.FromSlice([]int{1}), func(a, b int) bool { return a == b })

	chunk1 := iter.Next()
	assert.True(t, chunk1.IsSome())
	assert.Equal(t, []int{1}, chunk1.Unwrap())

	assert.True(t, iter.Next().IsNone())
}

func TestMapWindows(t *testing.T) {
	iter := iterator.MapWindows(iterator.FromSlice([]int{1, 2, 3, 4, 5}), 3, func(window []int) int {
		return window[0] + window[1] + window[2]
	})

	assert.Equal(t, option.Some(6), iter.Next())  // 1+2+3
	assert.Equal(t, option.Some(9), iter.Next())  // 2+3+4
	assert.Equal(t, option.Some(12), iter.Next()) // 3+4+5
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestMapWindowsSmall(t *testing.T) {
	iter := iterator.MapWindows(iterator.FromSlice([]int{1, 2}), 2, func(window []int) int {
		return window[0] + window[1]
	})

	assert.Equal(t, option.Some(3), iter.Next()) // 1+2
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestMapWindowsTooSmall(t *testing.T) {
	iter := iterator.MapWindows(iterator.FromSlice([]int{1}), 2, func(window []int) int {
		return window[0] + window[1]
	})

	// Should return None immediately since we don't have enough elements
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestArrayChunks_EmptyBuffer tests iterator.ArrayChunksIterable when buffer is empty
func TestArrayChunks_EmptyBuffer(t *testing.T) {
	// Test with empty iterator
	iter := iterator.ArrayChunks(iterator.Empty[int](), 2)
	assert.Equal(t, option.None[[]int](), iter.Next())

	// Test with iterator that becomes empty during chunking
	iter2 := iterator.ArrayChunks(iterator.FromSlice([]int{1}), 2)
	chunk1 := iter2.Next()
	assert.True(t, chunk1.IsSome())
	assert.Equal(t, []int{1}, chunk1.Unwrap())

	// Next call should return None (buffer is empty)
	assert.Equal(t, option.None[[]int](), iter2.Next())
}

// TestChunkBy_CurrentEmpty tests iterator.ChunkByIterable when current is empty
func TestChunkBy_CurrentEmpty(t *testing.T) {
	// This is a tricky case - current should never be empty in normal operation
	// But we can test the edge case where iter becomes None and current is empty
	iter := iterator.ChunkBy(iterator.FromSlice([]int{1, 2}), func(a, b int) bool { return a == b })

	chunk1 := iter.Next()
	assert.True(t, chunk1.IsSome())
	assert.Equal(t, []int{1}, chunk1.Unwrap())

	chunk2 := iter.Next()
	assert.True(t, chunk2.IsSome())
	assert.Equal(t, []int{2}, chunk2.Unwrap())

	// After all chunks are consumed, should return None
	assert.Equal(t, option.None[[]int](), iter.Next())
}

// TestMapWindows_SizeHint_UpperLessThanWindowSize tests iterator.MapWindows SizeHint when upper < windowSize
func TestMapWindows_SizeHint_UpperLessThanWindowSize(t *testing.T) {
	// Create an iterator with SizeHint upper < windowSize
	iter := iterator.MapWindows(iterator.FromSlice([]int{1, 2}), 3, func(window []int) int {
		return window[0] + window[1] + window[2]
	})

	lower, upper := iter.SizeHint()
	// lower should be 0 (since 2 < 3)
	assert.Equal(t, uint(0), lower)
	// upper should be Some(0) (since 2 < 3)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(0), upper.Unwrap())
}

// TestArrayChunks_SizeHint tests ArrayChunks SizeHint method
func TestArrayChunks_SizeHint(t *testing.T) {
	// Test with known size
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	chunks := iterator.ArrayChunks(iter, 3)
	lower, upper := chunks.SizeHint()
	// Should have lower > 0 and upper > 0
	assert.True(t, lower > 0 || upper.IsSome())
	if upper.IsSome() {
		assert.True(t, upper.Unwrap() > 0)
	}
}

// TestChunkBy_SizeHint tests ChunkBy SizeHint method
func TestChunkBy_SizeHint(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 1, 2, 2, 2, 3, 3})
	chunks := iterator.ChunkBy(iter, func(a, b int) bool { return a == b })
	lower, upper := chunks.SizeHint()
	// ChunkBy can't provide accurate size hint
	assert.Equal(t, uint(0), lower)
	assert.True(t, upper.IsNone())
}

func TestIterator_XMapWindows(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	windows := iter.XMapWindows(3, func(window []int) any {
		return window[0] + window[1] + window[2]
	})
	result := windows.Collect()
	assert.Equal(t, []any{6, 9, 12}, result)
}

// TestIterator_WrapperMethods_Chunking tests chunking wrapper methods from TestIterator_WrapperMethods
func TestIterator_WrapperMethods_Chunking(t *testing.T) {
	// Test MapWindows (covers iterator_methods.go:776-778)
	iter9 := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	windows := iter9.MapWindows(3, func(window []int) int {
		return window[0] + window[1] + window[2]
	})
	result7 := windows.Collect()
	assert.Equal(t, []int{6, 9, 12}, result7)
}
