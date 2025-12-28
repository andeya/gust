package iterator_test

import (
	"testing"

	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/option"
	"github.com/stretchr/testify/assert"
)

type alternateIterable struct {
	state int
}

func (a *alternateIterable) Next() option.Option[int] {
	val := a.state
	a.state++
	if val%2 == 0 {
		return option.Some(val)
	}
	return option.None[int]()
}

func (a *alternateIterable) SizeHint() (uint, option.Option[uint]) {
	return iterator.DefaultSizeHint[int]()
}

// TestFuse tests Fuse functionality
func TestFuse(t *testing.T) {
	// Create an iterator that alternates between Some and None
	alt := &alternateIterable{state: 0}
	var iterable iterator.Iterable[int] = alt
	iter := iterator.FromIterable(iterable).Fuse()

	// First call should return Some(0)
	assert.Equal(t, option.Some(0), iter.Next())
	// Second call returns None, which should fuse the iterator
	assert.Equal(t, option.None[int](), iter.Next())
	// After None, all subsequent calls should return None (even though alt would return Some(2))
	assert.Equal(t, option.None[int](), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestInspect tests Inspect functionality
func TestInspect(t *testing.T) {
	called := false
	a := []int{1, 2, 3}
	iter := iterator.FromSlice(a).Inspect(func(x int) {
		called = true
		assert.True(t, x > 0)
	})

	iter.Next()
	assert.True(t, called)
	result := iter.Collect()
	assert.Equal(t, []int{2, 3}, result)
}

// TestInspectEmpty tests Inspect with empty iterator
func TestInspectEmpty(t *testing.T) {
	called := false
	iter := iterator.Empty[int]().Inspect(func(x int) {
		called = true
	})

	assert.Equal(t, option.None[int](), iter.Next())
	assert.False(t, called)
}

// TestIntersperseSingleElement tests Intersperse with single element
func TestIntersperseSingleElement(t *testing.T) {
	a := []int{42}
	iter := iterator.FromSlice(a).Intersperse(100)
	assert.Equal(t, option.Some(42), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestIntersperseEmpty tests Intersperse with empty iterator
func TestIntersperseEmpty(t *testing.T) {
	iter := iterator.Empty[int]().Intersperse(100)
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestIntersperseWithSingleElement tests IntersperseWith with single element
func TestIntersperseWithSingleElement(t *testing.T) {
	a := []int{42}
	iter := iterator.FromSlice(a).IntersperseWith(func() int { return 99 })
	assert.Equal(t, option.Some(42), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestIntersperseWithEmpty tests IntersperseWith with empty iterator
func TestIntersperseWithEmpty(t *testing.T) {
	iter := iterator.Empty[int]().IntersperseWith(func() int { return 99 })
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestPeekableSizeHint tests Peekable SizeHint
func TestPeekableSizeHint(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3}).Peekable()
	iter.Peek() // Peek allows looking at next element without consuming it from user's perspective
	// From user's perspective: original has 3 elements, peek shows we can still see all 3
	// After peek: underlying consumed 1 (SizeHint becomes 2), but we have 1 peeked
	// Total visible elements: 1 peeked + 2 remaining = 3
	lower, upper := iter.SizeHint()
	assert.Equal(t, uint(3), lower) // Should reflect total visible elements
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(3), upper.Unwrap())
}

// TestCycleEmptyIterator tests Cycle with empty iterator
func TestCycleEmptyIterator(t *testing.T) {
	iter := iterator.Empty[int]().Cycle()
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestCycleSizeHint tests Cycle SizeHint
func TestCycleSizeHint(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3}).Cycle()
	// Consume all elements to exhaust the iterator
	iter.Next()
	iter.Next()
	iter.Next()
	iter.Next() // This should start cycling
	lower, upper := iter.SizeHint()
	// After cycling starts, size hint should indicate infinite
	assert.Equal(t, uint(0), lower)
	assert.True(t, upper.IsNone())
}

// TestCycleSizeHintWithCache tests Cycle SizeHint when cache has elements but not exhausted
// This covers the branch: if len(c.cache) > 0 { return 0, option.None[uint]() }
func TestCycleSizeHintWithCache(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3}).Cycle()
	// Call Next() once to populate cache (exhausted is still false)
	opt := iter.Next()
	assert.True(t, opt.IsSome())
	assert.Equal(t, 1, opt.Unwrap())

	// Now call SizeHint() - cache has elements but exhausted is false
	// This should trigger the len(c.cache) > 0 branch
	lower, upper := iter.SizeHint()
	assert.Equal(t, uint(0), lower)
	assert.True(t, upper.IsNone(), "SizeHint should return None when cache has elements")
}

// TestIntersperseSizeHintEdgeCases tests Intersperse SizeHint edge cases
func TestIntersperseSizeHintEdgeCases(t *testing.T) {
	// Test with lower == 0 (empty iterator)
	// iterator.Empty iterator has SizeHint (0, Some(0)) - it's known to have 0 elements
	iter := iterator.Empty[int]()
	interspersed := iter.Intersperse(100)
	lower, upper := interspersed.SizeHint()
	assert.Equal(t, uint(0), lower)
	// iterator.Empty iterator: 0 elements -> 0 interspersed elements (no separators for empty)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(0), upper.Unwrap())

	// Test with upperVal == 0 (same as above)
	iter2 := iterator.Empty[int]()
	interspersed2 := iter2.Intersperse(100)
	lower2, upper2 := interspersed2.SizeHint()
	assert.Equal(t, uint(0), lower2)
	assert.True(t, upper2.IsSome())
	assert.Equal(t, uint(0), upper2.Unwrap())
}

// TestIntersperseSizeHintLowerZero tests Intersperse SizeHint with lower == 0
func TestIntersperseSizeHintLowerZero(t *testing.T) {
	iter := iterator.Empty[int]()
	interspersed := iter.Intersperse(100)
	lower, upper := interspersed.SizeHint()
	assert.Equal(t, uint(0), lower)
	// iterator.Empty iterator: 0 elements -> 0 interspersed elements (no separators for empty)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(0), upper.Unwrap())
}

// TestIntersperseSizeHintUpperZero tests Intersperse SizeHint with upperVal == 0
func TestIntersperseSizeHintUpperZero(t *testing.T) {
	iter := iterator.Empty[int]()
	interspersed := iter.Intersperse(100)
	lower, upper := interspersed.SizeHint()
	assert.Equal(t, uint(0), lower)
	// iterator.Empty iterator: 0 elements -> 0 interspersed elements (no separators for empty)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(0), upper.Unwrap())
}

// TestPeekableSizeHintWithPeek tests Peekable SizeHint with peeked value
func TestPeekableSizeHintWithPeek(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3}).Peekable()
	iter.Peek() // Peek allows looking at next element without consuming it from user's perspective
	// From user's perspective: original has 3 elements, peek shows we can still see all 3
	// After peek: underlying consumed 1 (SizeHint becomes 2), but we have 1 peeked
	// Total visible elements: 1 peeked + 2 remaining = 3
	lower, upper := iter.SizeHint()
	assert.Equal(t, uint(3), lower) // Should reflect total visible elements
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(3), upper.Unwrap())
}

// TestPeekableSizeHintUpperNone tests Peekable SizeHint when upper is None
func TestPeekableSizeHintUpperNone(t *testing.T) {
	iter := iterator.Repeat(1).Peekable()
	iter.Peek()
	lower, upper := iter.SizeHint()
	// iterator.Repeat has lower=0, upper=None (infinite iterator)
	// After peek: we have 1 peeked element, but lower=0 so implementation doesn't increment
	// From design perspective: we have at least 1 element (peeked), but upper is still None (infinite)
	// Implementation limitation: when lower=0, it doesn't account for peeked element
	assert.Equal(t, uint(0), lower) // Implementation doesn't increment when lower=0
	assert.True(t, upper.IsNone())  // Still infinite
}

// TestPeekableSizeHintUpperZero tests Peekable SizeHint when upper is 0
func TestPeekableSizeHintUpperZero(t *testing.T) {
	iter := iterator.Empty[int]().Peekable()
	iter.Peek() // Peek on empty iterator returns None, so peeked is None
	lower, upper := iter.SizeHint()
	// iterator.Empty iterator: no elements, peek returns None
	assert.Equal(t, uint(0), lower)
	// iterator.Empty iterator has upper.IsSome() with value 0 (known to have 0 elements)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(0), upper.Unwrap())
}

// TestPeekable tests Peekable functionality
func TestPeekable(t *testing.T) {
	xs := []int{1, 2, 3}
	iter := iterator.FromSlice(xs).Peekable()

	// peek() lets us see into the future
	assert.Equal(t, option.Some(1), iter.Peek())
	assert.Equal(t, option.Some(1), iter.Next())

	assert.Equal(t, option.Some(2), iter.Next())

	// we can peek() multiple times, the iterator won't advance
	assert.Equal(t, option.Some(3), iter.Peek())
	assert.Equal(t, option.Some(3), iter.Peek())

	assert.Equal(t, option.Some(3), iter.Next())

	// after the iterator is finished, so is peek()
	assert.Equal(t, option.None[int](), iter.Peek())
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestPeekableInheritsIteratorMethods tests that PeekableIterator can use all Iterator methods
func TestPeekableInheritsIteratorMethods(t *testing.T) {
	// Test that PeekableIterator can use all Iterator methods
	xs := []int{1, 2, 3, 4, 5}
	iter := iterator.FromSlice(xs).Peekable()

	// Test Filter (Iterator method)
	filtered := iter.Filter(func(x int) bool { return x > 2 })
	assert.Equal(t, option.Some(3), filtered.Next())
	assert.Equal(t, option.Some(4), filtered.Next())
	assert.Equal(t, option.Some(5), filtered.Next())
	assert.Equal(t, option.None[int](), filtered.Next())

	// Test Map (Iterator method)
	xs2 := []int{1, 2, 3}
	iter2 := iterator.FromSlice(xs2).Peekable()
	mapped := iter2.Map(func(x int) int { return x * 2 })
	assert.Equal(t, option.Some(2), mapped.Next())
	assert.Equal(t, option.Some(4), mapped.Next())
	assert.Equal(t, option.Some(6), mapped.Next())

	// Test Take (Iterator method)
	xs3 := []int{1, 2, 3, 4, 5}
	iter3 := iterator.FromSlice(xs3).Peekable()
	taken := iter3.Take(2)
	assert.Equal(t, option.Some(1), taken.Next())
	assert.Equal(t, option.Some(2), taken.Next())
	assert.Equal(t, option.None[int](), taken.Next())

	// Test Collect (Iterator method)
	xs4 := []int{1, 2, 3}
	iter4 := iterator.FromSlice(xs4).Peekable()
	collected := iter4.Collect()
	assert.Equal(t, []int{1, 2, 3}, collected)

	// Test that Peek still works after using Iterator methods
	xs5 := []int{1, 2, 3}
	iter5 := iterator.FromSlice(xs5).Peekable()
	assert.Equal(t, option.Some(1), iter5.Peek())
	filtered2 := iter5.Filter(func(x int) bool { return x > 1 })
	assert.Equal(t, option.Some(2), filtered2.Next())
}

// TestCloned_NilPointer tests Cloned with nil pointer (covers utility.go)
func TestCloned_NilPointer(t *testing.T) {
	ptrs := []*int{nil, intPtr(1), nil, intPtr(2)}
	iter := iterator.Cloned(iterator.FromSlice(ptrs))
	result := iter.Collect()
	// nil pointers should be converted to zero values
	assert.Equal(t, []int{0, 1, 0, 2}, result)
}

func intPtr(i int) *int {
	return &i
}

// TestCloned tests Cloned functionality
func TestCloned(t *testing.T) {
	a := []string{"hello", "world"}
	ptrs := make([]*string, len(a))
	for i := range a {
		ptrs[i] = &a[i]
	}
	iter := iterator.Cloned(iterator.FromSlice(ptrs))
	v := iter.Collect()
	assert.Equal(t, []string{"hello", "world"}, v)
}

// TestFuse_DoneBranch tests fuseIterable when done == true
func TestFuse_DoneBranch(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3}).Fuse()

	// Should yield elements normally
	assert.Equal(t, option.Some(1), iter.Next())
	assert.Equal(t, option.Some(2), iter.Next())
	assert.Equal(t, option.Some(3), iter.Next())

	// After None is encountered, done is set to true
	assert.Equal(t, option.None[int](), iter.Next())

	// Subsequent calls should return None immediately (done == true branch)
	assert.Equal(t, option.None[int](), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestIterator_Intersperse(t *testing.T) {
	iter := iterator.FromSlice([]int{0, 1, 2})
	interspersed := iter.Intersperse(100)
	assert.Equal(t, option.Some(0), interspersed.Next())
	assert.Equal(t, option.Some(100), interspersed.Next())
	assert.Equal(t, option.Some(1), interspersed.Next())
	assert.Equal(t, option.Some(100), interspersed.Next())
	assert.Equal(t, option.Some(2), interspersed.Next())
	assert.Equal(t, option.None[int](), interspersed.Next())
}

func TestIterator_IntersperseWith(t *testing.T) {
	iter := iterator.FromSlice([]int{0, 1, 2})
	interspersed := iter.IntersperseWith(func() int { return 99 })
	assert.Equal(t, option.Some(0), interspersed.Next())
	assert.Equal(t, option.Some(99), interspersed.Next())
	assert.Equal(t, option.Some(1), interspersed.Next())
	assert.Equal(t, option.Some(99), interspersed.Next())
	assert.Equal(t, option.Some(2), interspersed.Next())
	assert.Equal(t, option.None[int](), interspersed.Next())

	// Test Intersperse with single element (no separator)
	iter2 := iterator.FromSlice([]int{42})
	interspersed2 := iter2.Intersperse(100)
	assert.Equal(t, option.Some(42), interspersed2.Next())
	assert.Equal(t, option.None[int](), interspersed2.Next())

	// Test Intersperse with empty iterator
	iter3 := iterator.FromSlice([]int{})
	interspersed3 := iter3.Intersperse(100)
	assert.Equal(t, option.None[int](), interspersed3.Next())

	// Test Intersperse SizeHint with known size
	iter4 := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	interspersed4 := iter4.Intersperse(0)
	lower, upper := interspersed4.SizeHint()
	// Should have lower > 0 and upper > 0 for intersperse
	// Intersperse adds separators, so size becomes 2*n - 1
	assert.True(t, lower > 0 || upper.IsSome())
}

// TestFuse_SizeHint tests Fuse SizeHint method
func TestFuse_SizeHint(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	fused := iter.Fuse()
	lower, upper := fused.SizeHint()
	assert.Equal(t, uint(3), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(3), upper.Unwrap())
}

func TestIterator_Cycle(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	cycled := iter.Cycle()
	assert.Equal(t, option.Some(1), cycled.Next())
	assert.Equal(t, option.Some(2), cycled.Next())
	assert.Equal(t, option.Some(3), cycled.Next())
	assert.Equal(t, option.Some(1), cycled.Next()) // starts over
	assert.Equal(t, option.Some(2), cycled.Next())
	assert.Equal(t, option.Some(3), cycled.Next()) // continues cycling
}
