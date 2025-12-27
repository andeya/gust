package iter

import (
	"strconv"
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

func TestRetMap(t *testing.T) {
	iter := RetMap(FromSlice([]string{"1", "2", "3", "NaN"}), strconv.Atoi)

	assert.Equal(t, gust.Some(gust.Ok(1)), iter.Next())
	assert.Equal(t, gust.Some(gust.Ok(2)), iter.Next())
	assert.Equal(t, gust.Some(gust.Ok(3)), iter.Next())
	assert.Equal(t, true, iter.Next().Unwrap().IsErr())
	assert.Equal(t, gust.None[gust.Result[int]](), iter.Next())
}

func TestOptMap(t *testing.T) {
	iter := OptMap(FromSlice([]string{"1", "2", "3", "NaN"}), func(s string) *int {
		if v, err := strconv.Atoi(s); err == nil {
			return &v
		} else {
			return nil
		}
	})
	var newInt = func(v int) *int {
		return &v
	}
	assert.Equal(t, gust.Some(gust.Some(newInt(1))), iter.Next())
	assert.Equal(t, gust.Some(gust.Some(newInt(2))), iter.Next())
	assert.Equal(t, gust.Some(gust.Some(newInt(3))), iter.Next())
	assert.Equal(t, gust.Some(gust.None[*int]()), iter.Next())
	assert.Equal(t, gust.None[gust.Option[*int]](), iter.Next())
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

// TestSkipWhile_DoneBranch tests skipWhileIterable when done == true (covers adapters_extended.go:32)
func TestSkipWhile_DoneBranch(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3, 4, 5}).SkipWhile(func(x int) bool { return x < 3 })

	// First call should skip 1, 2 and return 3
	assert.Equal(t, gust.Some(3), iter.Next())

	// After done is set to true, subsequent calls should use s.iter.Next() directly
	assert.Equal(t, gust.Some(4), iter.Next())
	assert.Equal(t, gust.Some(5), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

// TestTakeWhile_PredicateFalse tests takeWhileIterable when predicate returns false (covers adapters_extended.go:59)
func TestTakeWhile_PredicateFalse(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3, 4, 5}).TakeWhile(func(x int) bool { return x < 3 })

	// Should return elements while predicate is true
	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.Some(2), iter.Next())

	// When predicate returns false, should return None
	assert.Equal(t, gust.None[int](), iter.Next())

	// Subsequent calls should also return None
	assert.Equal(t, gust.None[int](), iter.Next())
}

// TestFlatMap_CurrentNil tests flatMapIterable when current becomes nil (covers adapters_extended.go:204)
func TestFlatMap_CurrentNil(t *testing.T) {
	iter := FlatMap(FromSlice([]int{1, 2}), func(x int) Iterator[int] {
		if x == 1 {
			return FromSlice([]int{10, 20}) // Non-empty iterator
		}
		return Empty[int]() // Empty iterator - current will become nil
	})

	// Should yield elements from first iterator
	assert.Equal(t, gust.Some(10), iter.Next())
	assert.Equal(t, gust.Some(20), iter.Next())

	// After first iterator is exhausted, current becomes nil, should move to next
	// Second iterator is empty, so should return None
	assert.Equal(t, gust.None[int](), iter.Next())
}

// TestFlatten_CurrentNil tests flattenIterable when current becomes nil (covers adapters_extended.go:259)
func TestFlatten_CurrentNil(t *testing.T) {
	iter := Flatten(FromSlice([]Iterator[int]{
		FromSlice([]int{1, 2}),
		Empty[int](), // Empty iterator - current will become nil
		FromSlice([]int{3, 4}),
	}))

	// Should yield elements from first iterator
	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.Some(2), iter.Next())

	// After first iterator is exhausted, current becomes nil, should move to next
	// Second iterator is empty, so should skip to third
	assert.Equal(t, gust.Some(3), iter.Next())
	assert.Equal(t, gust.Some(4), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

// TestFuse_DoneBranch tests fuseIterable when done == true (covers adapters_extended.go:287)
func TestFuse_DoneBranch(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3}).Fuse()

	// Should yield elements normally
	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.Some(2), iter.Next())
	assert.Equal(t, gust.Some(3), iter.Next())

	// After None is encountered, done is set to true
	assert.Equal(t, gust.None[int](), iter.Next())

	// Subsequent calls should return None immediately (done == true branch)
	assert.Equal(t, gust.None[int](), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

// TestArrayChunks_EmptyBuffer tests arrayChunksIterable when buffer is empty (covers adapters_extended.go:444)
func TestArrayChunks_EmptyBuffer(t *testing.T) {
	// Test with empty iterator
	iter := ArrayChunks(Empty[int](), 2)
	assert.Equal(t, gust.None[[]int](), iter.Next())

	// Test with iterator that becomes empty during chunking
	iter2 := ArrayChunks(FromSlice([]int{1}), 2)
	chunk1 := iter2.Next()
	assert.True(t, chunk1.IsSome())
	assert.Equal(t, []int{1}, chunk1.Unwrap())

	// Next call should return None (buffer is empty)
	assert.Equal(t, gust.None[[]int](), iter2.Next())
}

// TestChunkBy_CurrentEmpty tests chunkByIterable when current is empty (covers adapters_extended.go:531)
func TestChunkBy_CurrentEmpty(t *testing.T) {
	// This is a tricky case - current should never be empty in normal operation
	// But we can test the edge case where iter becomes None and current is empty
	iter := ChunkBy(FromSlice([]int{1, 2}), func(a, b int) bool { return a == b })

	chunk1 := iter.Next()
	assert.True(t, chunk1.IsSome())
	assert.Equal(t, []int{1}, chunk1.Unwrap())

	chunk2 := iter.Next()
	assert.True(t, chunk2.IsSome())
	assert.Equal(t, []int{2}, chunk2.Unwrap())

	// After all chunks are consumed, should return None
	assert.Equal(t, gust.None[[]int](), iter.Next())
}

// TestMapWindows_SizeHint_UpperLessThanWindowSize tests MapWindows SizeHint when upper < windowSize (covers adapters_extended.go:625)
func TestMapWindows_SizeHint_UpperLessThanWindowSize(t *testing.T) {
	// Create an iterator with SizeHint upper < windowSize
	iter := MapWindows(FromSlice([]int{1, 2}), 3, func(window []int) int {
		return window[0] + window[1] + window[2]
	})

	lower, upper := iter.SizeHint()
	// lower should be 0 (since 2 < 3)
	assert.Equal(t, uint(0), lower)
	// upper should be Some(0) (since 2 < 3)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(0), upper.Unwrap())
}
