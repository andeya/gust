package iterator_test

import (
	"testing"

	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/option"
	"github.com/stretchr/testify/assert"
)

func TestMax(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{}
	assert.Equal(t, option.Some(3), iterator.Max(iterator.FromSlice(a)))
	assert.Equal(t, option.None[int](), iterator.Max(iterator.FromSlice(b)))
}

func TestMin(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{}
	assert.Equal(t, option.Some(1), iterator.Min(iterator.FromSlice(a)))
	assert.Equal(t, option.None[int](), iterator.Min(iterator.FromSlice(b)))
}

func TestMaxByKey(t *testing.T) {
	a := []int{-3, 0, 1, 5, -10}
	max := iterator.MaxByKey(iterator.FromSlice(a), func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	})
	assert.Equal(t, option.Some(-10), max)
}

func TestMinByKey(t *testing.T) {
	a := []int{-3, 0, 1, 5, -10}
	min := iterator.MinByKey(iterator.FromSlice(a), func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	})
	assert.Equal(t, option.Some(0), min)
}

func TestMaxBy(t *testing.T) {
	a := []int{-3, 0, 1, 5, -10}
	max := iterator.FromSlice(a).MaxBy(func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
	assert.Equal(t, option.Some(5), max)
}

func TestMinBy(t *testing.T) {
	a := []int{-3, 0, 1, 5, -10}
	min := iterator.FromSlice(a).MinBy(func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
	assert.Equal(t, option.Some(-10), min)
}

func TestMax_EdgeCases(t *testing.T) {
	// Test Max with single element
	{
		result := iterator.Max(iterator.FromSlice([]int{42}))
		assert.True(t, result.IsSome())
		assert.Equal(t, 42, result.Unwrap())
	}

	// Test Max with equal elements (should return last)
	{
		result := iterator.Max(iterator.FromSlice([]int{5, 5, 5}))
		assert.True(t, result.IsSome())
		assert.Equal(t, 5, result.Unwrap())
	}

	// Test Max with string
	{
		result := iterator.Max(iterator.FromSlice([]string{"z", "a", "m"}))
		assert.True(t, result.IsSome())
		assert.Equal(t, "z", result.Unwrap())
	}

	// Test Max with uint
	{
		result := iterator.Max(iterator.FromSlice([]uint{100, 50, 200}))
		assert.True(t, result.IsSome())
		assert.Equal(t, uint(200), result.Unwrap())
	}
}

func TestMin_EdgeCases(t *testing.T) {
	// Test Min with single element
	{
		result := iterator.Min(iterator.FromSlice([]int{42}))
		assert.True(t, result.IsSome())
		assert.Equal(t, 42, result.Unwrap())
	}

	// Test Min with equal elements (should return first)
	{
		result := iterator.Min(iterator.FromSlice([]int{5, 5, 5}))
		assert.True(t, result.IsSome())
		assert.Equal(t, 5, result.Unwrap())
	}

	// Test Min with string
	{
		result := iterator.Min(iterator.FromSlice([]string{"z", "a", "m"}))
		assert.True(t, result.IsSome())
		assert.Equal(t, "a", result.Unwrap())
	}

	// Test Min with uint
	{
		result := iterator.Min(iterator.FromSlice([]uint{100, 50, 200}))
		assert.True(t, result.IsSome())
		assert.Equal(t, uint(50), result.Unwrap())
	}
}

func TestMaxByKey_EdgeCases(t *testing.T) {
	// Test MaxByKey with single element
	{
		result := iterator.MaxByKey(iterator.FromSlice([]int{42}), func(x int) int {
			return x * 2
		})
		assert.True(t, result.IsSome())
		assert.Equal(t, 42, result.Unwrap())
	}

	// Test MaxByKey with equal keys (should return last)
	{
		result := iterator.MaxByKey(iterator.FromSlice([]int{1, 2, 3}), func(x int) int {
			return 10 // All have same key
		})
		assert.True(t, result.IsSome())
		assert.Equal(t, 3, result.Unwrap()) // Last element
	}

	// Test MaxByKey with negative values
	{
		result := iterator.MaxByKey(iterator.FromSlice([]int{-3, -1, -5}), func(x int) int {
			return -x // Negate to find max absolute value
		})
		assert.True(t, result.IsSome())
		assert.Equal(t, -5, result.Unwrap())
	}
}

func TestMinByKey_EdgeCases(t *testing.T) {
	// Test MinByKey with single element
	{
		result := iterator.MinByKey(iterator.FromSlice([]int{42}), func(x int) int {
			return x * 2
		})
		assert.True(t, result.IsSome())
		assert.Equal(t, 42, result.Unwrap())
	}

	// Test MinByKey with equal keys (should return first)
	{
		result := iterator.MinByKey(iterator.FromSlice([]int{1, 2, 3}), func(x int) int {
			return 10 // All have same key
		})
		assert.True(t, result.IsSome())
		assert.Equal(t, 1, result.Unwrap()) // First element
	}

	// Test MinByKey with negative values
	{
		result := iterator.MinByKey(iterator.FromSlice([]int{-3, -1, -5}), func(x int) int {
			return -x // Negate to find min absolute value
		})
		assert.True(t, result.IsSome())
		assert.Equal(t, -1, result.Unwrap())
	}
}

// TestIteratorExtConsumerMethods_MinMax tests min/max methods from TestIteratorExtConsumerMethods
func TestIteratorExtConsumerMethods_MinMax(t *testing.T) {
	// Test Max
	iter6 := iterator.FromSlice([]int{1, 3, 2})
	assert.Equal(t, option.Some(3), iterator.Max(iter6))

	// Test Min
	iter7 := iterator.FromSlice([]int{3, 1, 2})
	assert.Equal(t, option.Some(1), iterator.Min(iter7))

	// Test MaxBy
	iter8 := iterator.FromSlice([]int{-3, 0, 1, 5, -10})
	maxBy := iter8.MaxBy(func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
	assert.Equal(t, option.Some(5), maxBy)

	// Test MinBy
	iter9 := iterator.FromSlice([]int{-3, 0, 1, 5, -10})
	minBy := iter9.MinBy(func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
	assert.Equal(t, option.Some(-10), minBy)

	// Test MaxByKey (function version)
	iter10 := iterator.FromSlice([]int{-3, 0, 1, 5, -10})
	maxByKey := iterator.MaxByKey(iter10, func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	})
	assert.Equal(t, option.Some(-10), maxByKey)

	// Test MinByKey (function version)
	iter11 := iterator.FromSlice([]int{-3, 0, 1, 5, -10})
	minByKey := iterator.MinByKey(iter11, func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	})
	assert.Equal(t, option.Some(0), minByKey)
}
