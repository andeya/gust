package iterator_test

import (
	"errors"
	"testing"

	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/option"
	"github.com/andeya/gust/result"
	"github.com/stretchr/testify/assert"
)

func TestFold(t *testing.T) {
	a := []int{1, 2, 3}
	sum := iterator.Fold(iterator.FromSlice(a), 0, func(acc int, x int) int { return acc + x })
	assert.Equal(t, 6, sum)
}

func TestReduce(t *testing.T) {
	reduced := iterator.FromRange(1, 10).Reduce(func(acc int, e int) int { return acc + e })
	assert.Equal(t, option.Some(45), reduced)
}

func TestIterator_XFold(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	sum := iter.XFold(0, func(acc any, x int) any { return acc.(int) + x })
	assert.Equal(t, 6, sum)
}

func TestIterator_XTryFold(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	sum := iter.XTryFold(0, func(acc any, x int) result.Result[any] {
		return result.Ok[any](any(acc.(int) + x))
	})
	assert.True(t, sum.IsOk())
	assert.Equal(t, 6, sum.Unwrap().(int))
}

// TestReduceEmpty tests Reduce with empty iterator
func TestReduceEmpty(t *testing.T) {
	iter := iterator.Empty[int]()
	result := iter.Reduce(func(acc int, x int) int { return acc + x })
	assert.True(t, result.IsNone())
}

// TestTryFold tests TryFold functionality
func TestTryFold(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	sum := iterator.TryFold(iter, 0, func(acc int, x int) result.Result[int] {
		return result.Ok(acc + x)
	})
	assert.True(t, sum.IsOk())
	assert.Equal(t, 6, sum.Unwrap())
}

func TestRfold(t *testing.T) {
	a := []int{1, 2, 3}
	iter := iterator.FromSlice(a)
	deIter := iter.MustToDoubleEnded()

	// Test XRfold (any version)
	sum := deIter.XRfold(0, func(acc any, x int) any { return acc.(int) + x })
	assert.Equal(t, 6, sum.(int))

	// Test Rfold (T version)
	iter2 := iterator.FromSlice(a)
	deIter2 := iter2.MustToDoubleEnded()
	sum2 := deIter2.Rfold(0, func(acc int, x int) int { return acc + x })
	assert.Equal(t, 6, sum2)

	// Test generic function version
	iter3 := iterator.FromSlice(a)
	deIter3 := iter3.MustToDoubleEnded()
	sum3 := iterator.Rfold(deIter3, 0, func(acc int, x int) int { return acc + x })
	assert.Equal(t, 6, sum3)
}

func TestTryRfold(t *testing.T) {
	a := []int{1, 2, 3}
	iter := iterator.FromSlice(a)
	deIter := iter.MustToDoubleEnded()

	// Test XTryRfold (any version)
	sum := deIter.XTryRfold(0, func(acc any, x int) result.Result[any] {
		return result.Ok(any(acc.(int) + x))
	})
	assert.True(t, sum.IsOk())
	assert.Equal(t, 6, sum.Unwrap().(int))

	// Test TryRfold (T version)
	iter2 := iterator.FromSlice(a)
	deIter2 := iter2.MustToDoubleEnded()
	sum2 := deIter2.TryRfold(0, func(acc int, x int) result.Result[int] {
		return result.Ok(acc + x)
	})
	assert.True(t, sum2.IsOk())
	assert.Equal(t, 6, sum2.Unwrap())

	// Test generic function version
	iter3 := iterator.FromSlice(a)
	deIter3 := iter3.MustToDoubleEnded()
	sum3 := iterator.TryRfold(deIter3, 0, func(acc int, x int) result.Result[int] {
		return result.Ok(acc + x)
	})
	assert.True(t, sum3.IsOk())
	assert.Equal(t, 6, sum3.Unwrap())
}

// TestTryRfoldError tests TryRfold with error
func TestTryRfoldError(t *testing.T) {
	a := []int{1, 2, 3}
	iter := iterator.FromSlice(a)
	deIter := iter.MustToDoubleEnded()

	sum := deIter.TryRfold(0, func(acc int, x int) result.Result[int] {
		if acc+x > 5 {
			return result.TryErr[int](errors.New("overflow"))
		}
		return result.Ok(acc + x)
	})
	assert.True(t, sum.IsErr())
}

// TestTryRfoldEmpty tests TryRfold with empty iterator
func TestTryRfoldEmpty(t *testing.T) {
	iter := iterator.FromSlice([]int{})
	deIter := iter.MustToDoubleEnded()

	sum := deIter.TryRfold(0, func(acc int, x int) result.Result[int] {
		return result.Ok(acc + x)
	})
	assert.True(t, sum.IsOk())
	assert.Equal(t, 0, sum.Unwrap())
}

// TestIterator_WrapperMethods_FoldReduce tests fold/reduce wrapper methods from TestIterator_WrapperMethods
func TestIterator_WrapperMethods_FoldReduce(t *testing.T) {
	// Test XFold (covers iterator_methods.go:492-494)
	iter3 := iterator.FromSlice([]int{1, 2, 3})
	sum1 := iter3.XFold(0, func(acc any, x int) any {
		return acc.(int) + x
	})
	assert.Equal(t, 6, sum1)

	// Test Fold (covers iterator_methods.go:520-522)
	iter4 := iterator.FromSlice([]int{1, 2, 3})
	sum2 := iter4.Fold(0, func(acc int, x int) int {
		return acc + x
	})
	assert.Equal(t, 6, sum2)

	// Test XTryFold (covers iterator_methods.go:552-554)
	iter5 := iterator.FromSlice([]int{1, 2, 3})
	result3 := iter5.XTryFold(0, func(acc any, x int) result.Result[any] {
		return result.Ok[any](acc.(int) + x)
	})
	assert.True(t, result3.IsOk())
	assert.Equal(t, 6, result3.UnwrapUnchecked())
}
