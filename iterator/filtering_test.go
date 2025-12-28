package iterator_test

import (
	"testing"

	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/option"
	"github.com/stretchr/testify/assert"
)

func TestSkip(t *testing.T) {
	a := []int{1, 2, 3}
	iter := iterator.FromSlice(a).Skip(2)

	assert.Equal(t, option.Some(3), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestTake(t *testing.T) {
	a := []int{1, 2, 3}
	iter := iterator.FromSlice(a).Take(2)

	assert.Equal(t, option.Some(1), iter.Next())
	assert.Equal(t, option.Some(2), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestSkipWhile(t *testing.T) {
	a := []int{-1, 0, 1}
	iter := iterator.FromSlice(a).SkipWhile(func(x int) bool { return x < 0 })

	assert.Equal(t, option.Some(0), iter.Next())
	assert.Equal(t, option.Some(1), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestTakeWhile(t *testing.T) {
	a := []int{-1, 0, 1}
	iter := iterator.FromSlice(a).TakeWhile(func(x int) bool { return x < 0 })

	assert.Equal(t, option.Some(-1), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestSkipWhileAllSkipped tests SkipWhile when all elements are skipped
func TestSkipWhileAllSkipped(t *testing.T) {
	a := []int{1, 2, 3}
	iter := iterator.FromSlice(a).SkipWhile(func(x int) bool { return x > 0 })
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestSkipWhileNoneSkipped tests SkipWhile when no elements are skipped
func TestSkipWhileNoneSkipped(t *testing.T) {
	a := []int{-1, -2, -3}
	iter := iterator.FromSlice(a).SkipWhile(func(x int) bool { return x > 0 })
	assert.Equal(t, option.Some(-1), iter.Next())
	assert.Equal(t, option.Some(-2), iter.Next())
	assert.Equal(t, option.Some(-3), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestTakeWhileAllTaken tests TakeWhile when all elements are taken
func TestTakeWhileAllTaken(t *testing.T) {
	a := []int{-1, -2, -3}
	iter := iterator.FromSlice(a).TakeWhile(func(x int) bool { return x < 0 })
	assert.Equal(t, option.Some(-1), iter.Next())
	assert.Equal(t, option.Some(-2), iter.Next())
	assert.Equal(t, option.Some(-3), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestTakeWhileNoneTaken tests TakeWhile when no elements are taken
func TestTakeWhileNoneTaken(t *testing.T) {
	a := []int{1, 2, 3}
	iter := iterator.FromSlice(a).TakeWhile(func(x int) bool { return x < 0 })
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestSkipEmpty tests Skip with empty iterator
func TestSkipEmpty(t *testing.T) {
	iter := iterator.Empty[int]()
	skipped := iter.Skip(5)
	assert.Equal(t, option.None[int](), skipped.Next())
}

// TestTakeEmpty tests Take with empty iterator
func TestTakeEmpty(t *testing.T) {
	iter := iterator.Empty[int]()
	taken := iter.Take(5)
	assert.Equal(t, option.None[int](), taken.Next())
}

// TestTakeZero tests Take with zero count
func TestTakeZero(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	taken := iter.Take(0)
	assert.Equal(t, option.None[int](), taken.Next())
}

// TestSkipZero tests Skip with zero count
func TestSkipZero(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	skipped := iter.Skip(0)
	assert.Equal(t, option.Some(1), skipped.Next())
	assert.Equal(t, option.Some(2), skipped.Next())
	assert.Equal(t, option.Some(3), skipped.Next())
	assert.Equal(t, option.None[int](), skipped.Next())
}

// TestStepByPanic tests StepBy panic on zero step
func TestStepByPanic(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	assert.Panics(t, func() {
		iter.StepBy(0)
	})
}

// TestStepBySizeHint tests StepBy SizeHint edge cases
func TestStepBySizeHint(t *testing.T) {
	// Test with upper.IsSome() but upperVal == 0
	iter := iterator.FromSlice([]int{})
	stepIter := iter.StepBy(2)
	lower, upper := stepIter.SizeHint()
	assert.Equal(t, uint(0), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(0), upper.Unwrap())
}

// TestStepByFirstElement tests StepBy first element behavior
func TestStepByFirstElement(t *testing.T) {
	iter := iterator.FromSlice([]int{0, 1, 2, 3, 4, 5})
	stepIter := iter.StepBy(2)
	// First element should always be returned
	assert.Equal(t, option.Some(0), stepIter.Next())
	assert.Equal(t, option.Some(2), stepIter.Next())
	assert.Equal(t, option.Some(4), stepIter.Next())
	assert.Equal(t, option.None[int](), stepIter.Next())
}

// TestStepByAdvanceByError tests StepBy when AdvanceBy returns error
func TestStepByAdvanceByError(t *testing.T) {
	iter := iterator.FromSlice([]int{0, 1})
	stepIter := iter.StepBy(3)
	// First element
	assert.Equal(t, option.Some(0), stepIter.Next())
	// AdvanceBy(2) will fail, so Next() should return None
	assert.Equal(t, option.None[int](), stepIter.Next())
}

// TestSkipSizeHintEdgeCases tests Skip SizeHint edge cases
func TestSkipSizeHintEdgeCases(t *testing.T) {
	// Test with lower < n
	iter := iterator.FromSlice([]int{1, 2})
	skipped := iter.Skip(5)
	lower, upper := skipped.SizeHint()
	assert.Equal(t, uint(0), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(0), upper.Unwrap())

	// Test with upperVal < n
	iter2 := iterator.FromSlice([]int{1, 2})
	skipped2 := iter2.Skip(5)
	lower2, upper2 := skipped2.SizeHint()
	assert.Equal(t, uint(0), lower2)
	assert.True(t, upper2.IsSome())
	assert.Equal(t, uint(0), upper2.Unwrap())
}

// TestTakeSizeHintEdgeCases tests Take SizeHint edge cases
func TestTakeSizeHintEdgeCases(t *testing.T) {
	// Test with upper.IsSome() && upper.Unwrap() > n
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	taken := iter.Take(3)
	lower, upper := taken.SizeHint()
	assert.Equal(t, uint(3), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(3), upper.Unwrap())

	// Test with upper.IsSome() && upper.Unwrap() <= n
	iter2 := iterator.FromSlice([]int{1, 2})
	taken2 := iter2.Take(5)
	lower2, upper2 := taken2.SizeHint()
	assert.Equal(t, uint(2), lower2)
	assert.True(t, upper2.IsSome())
	assert.Equal(t, uint(2), upper2.Unwrap())
}

func TestStepBy(t *testing.T) {
	a := []int{0, 1, 2, 3, 4, 5}
	iter := iterator.FromSlice(a).StepBy(2)

	assert.Equal(t, option.Some(0), iter.Next())
	assert.Equal(t, option.Some(2), iter.Next())
	assert.Equal(t, option.Some(4), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestSkipWhile_DoneBranch tests skipWhileIterable when done == true
func TestSkipWhile_DoneBranch(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5}).SkipWhile(func(x int) bool { return x < 3 })

	// First call should skip 1, 2 and return 3
	assert.Equal(t, option.Some(3), iter.Next())

	// After done is set to true, subsequent calls should use s.iter.Next() directly
	assert.Equal(t, option.Some(4), iter.Next())
	assert.Equal(t, option.Some(5), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestTakeWhile_PredicateFalse tests takeWhileIterable when predicate returns false
func TestTakeWhile_PredicateFalse(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5}).TakeWhile(func(x int) bool { return x < 3 })

	// Should return elements while predicate is true
	assert.Equal(t, option.Some(1), iter.Next())
	assert.Equal(t, option.Some(2), iter.Next())

	// When predicate returns false, should return None
	assert.Equal(t, option.None[int](), iter.Next())

	// Subsequent calls should also return None
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestIteratorExtSkipTake(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})

	// Skip first 2, then take 2
	skipped := iter.Skip(2)
	taken := skipped.Take(2)
	result := taken.Collect()

	assert.Equal(t, []int{3, 4}, result)
}
