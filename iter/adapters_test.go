package iter

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestStepBy(t *testing.T) {
	a := []int{0, 1, 2, 3, 4, 5}
	iter := FromSlice(a).StepBy(2)

	assert.Equal(t, gust.Some(0), iter.Next())
	assert.Equal(t, gust.Some(2), iter.Next())
	assert.Equal(t, gust.Some(4), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestIntersperse(t *testing.T) {
	a := []int{0, 1, 2}
	iter := FromSlice(a).Intersperse(100)

	assert.Equal(t, gust.Some(0), iter.Next())
	assert.Equal(t, gust.Some(100), iter.Next())
	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.Some(100), iter.Next())
	assert.Equal(t, gust.Some(2), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestIntersperseWith(t *testing.T) {
	v := []int{0, 1, 2}
	iter := FromSlice(v).IntersperseWith(func() int { return 99 })

	assert.Equal(t, gust.Some(0), iter.Next())
	assert.Equal(t, gust.Some(99), iter.Next())
	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.Some(99), iter.Next())
	assert.Equal(t, gust.Some(2), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestArrayChunks(t *testing.T) {
	iter := ArrayChunks(FromSlice([]int{1, 2, 3, 4, 5, 6}), 2)

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
	iter := ArrayChunks(FromSlice([]int{1, 2, 3, 4, 5}), 2)

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
	iter := ArrayChunks(Empty[int](), 2)
	assert.True(t, iter.Next().IsNone())
}

func TestChunkBy(t *testing.T) {
	iter := ChunkBy(FromSlice([]int{1, 1, 2, 2, 2, 3, 3}), func(a, b int) bool { return a == b })

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
	iter := ChunkBy(Empty[int](), func(a, b int) bool { return a == b })
	assert.True(t, iter.Next().IsNone())
}

func TestChunkBySingle(t *testing.T) {
	iter := ChunkBy(FromSlice([]int{1}), func(a, b int) bool { return a == b })

	chunk1 := iter.Next()
	assert.True(t, chunk1.IsSome())
	assert.Equal(t, []int{1}, chunk1.Unwrap())

	assert.True(t, iter.Next().IsNone())
}

func TestMapWindows(t *testing.T) {
	iter := MapWindows(FromSlice([]int{1, 2, 3, 4, 5}), 3, func(window []int) int {
		return window[0] + window[1] + window[2]
	})

	assert.Equal(t, gust.Some(6), iter.Next())  // 1+2+3
	assert.Equal(t, gust.Some(9), iter.Next())  // 2+3+4
	assert.Equal(t, gust.Some(12), iter.Next()) // 3+4+5
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestMapWindowsSmall(t *testing.T) {
	iter := MapWindows(FromSlice([]int{1, 2}), 2, func(window []int) int {
		return window[0] + window[1]
	})

	assert.Equal(t, gust.Some(3), iter.Next()) // 1+2
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestMapWindowsTooSmall(t *testing.T) {
	iter := MapWindows(FromSlice([]int{1}), 2, func(window []int) int {
		return window[0] + window[1]
	})

	// Should return None immediately since we don't have enough elements
	assert.Equal(t, gust.None[int](), iter.Next())
}
