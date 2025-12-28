package iterator_test

import (
	"errors"
	"testing"

	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/result"
	"github.com/stretchr/testify/assert"
)

func TestTryFold(t *testing.T) {
	a := []int{1, 2, 3}
	sum := iterator.TryFold(iterator.FromSlice(a), 0, func(acc int, x int) result.Result[int] {
		// Simulate checked addition
		if acc > 100 {
			return result.TryErr[int](errors.New("overflow"))
		}
		return result.Ok(acc + x)
	})
	assert.True(t, sum.IsOk())
	assert.Equal(t, 6, sum.Unwrap())

	// Test short-circuiting
	a2 := []int{10, 20, 30, 100, 40, 50}
	iter := iterator.FromSlice(a2)
	sum2 := iterator.TryFold(iter, 0, func(acc int, x int) result.Result[int] {
		if acc+x > 50 {
			return result.TryErr[int](errors.New("overflow"))
		}
		return result.Ok(acc + x)
	})
	assert.True(t, sum2.IsErr())
}

func TestTryForEach(t *testing.T) {
	data := []string{"no_tea.txt", "stale_bread.json", "torrential_rain.png"}
	res := iterator.TryForEach(iterator.FromSlice(data), func(x string) result.Result[any] {
		// Simulate processing
		return result.Ok[any](nil)
	})
	assert.True(t, res.IsOk())
}

func TestTryReduce(t *testing.T) {
	numbers := []int{10, 20, 5, 23, 0}
	sum := iterator.FromSlice(numbers).TryReduce(func(x, y int) result.Result[int] {
		// Simulate checked addition
		if x+y > 100 {
			return result.TryErr[int](errors.New("overflow"))
		}
		return result.Ok(x + y)
	})
	assert.True(t, sum.IsOk())
	assert.True(t, sum.Unwrap().IsSome())
	assert.Equal(t, 58, sum.Unwrap().Unwrap())

	// Test empty iterator
	numbers2 := []int{}
	sum2 := iterator.FromSlice(numbers2).TryReduce(func(x, y int) result.Result[int] {
		return result.Ok(x + y)
	})
	assert.True(t, sum2.IsOk())
	assert.True(t, sum2.Unwrap().IsNone())
}

func TestTryFind(t *testing.T) {
	a := []string{"1", "2", "lol", "NaN", "5"}
	res := iterator.FromSlice(a).TryFind(func(s string) result.Result[bool] {
		if s == "lol" {
			return result.TryErr[bool](errors.New("invalid"))
		}
		if s == "2" {
			return result.Ok[bool](true)
		}
		return result.Ok[bool](false)
	})
	assert.True(t, res.IsOk())
	assert.True(t, res.Unwrap().IsSome())
	assert.Equal(t, "2", res.Unwrap().Unwrap())

	// Test TryFind with error before finding
	a2 := []string{"lol", "2", "3"}
	res2 := iterator.FromSlice(a2).TryFind(func(s string) result.Result[bool] {
		if s == "lol" {
			return result.TryErr[bool](errors.New("invalid"))
		}
		return result.Ok[bool](false)
	})
	assert.True(t, res2.IsErr())

	// Test TryFind with no match
	a3 := []string{"1", "3", "5"}
	res3 := iterator.FromSlice(a3).TryFind(func(s string) result.Result[bool] {
		return result.Ok[bool](false)
	})
	assert.True(t, res3.IsOk())
	assert.True(t, res3.Unwrap().IsNone())

	// Test TryFind with empty iterator
	a4 := []string{}
	res4 := iterator.FromSlice(a4).TryFind(func(s string) result.Result[bool] {
		return result.Ok[bool](true)
	})
	assert.True(t, res4.IsOk())
	assert.True(t, res4.Unwrap().IsNone())
}

func TestTryReduce_ErrorPath(t *testing.T) {
	// Test TryReduce with error
	numbers := []int{10, 20, 100}
	sum := iterator.FromSlice(numbers).TryReduce(func(x, y int) result.Result[int] {
		if x+y > 100 {
			return result.TryErr[int](errors.New("overflow"))
		}
		return result.Ok(x + y)
	})
	assert.True(t, sum.IsErr())
}

func TestTryForEach_ErrorPath(t *testing.T) {
	// Test TryForEach with error
	data := []string{"a", "error", "c"}
	res := iterator.TryForEach(iterator.FromSlice(data), func(x string) result.Result[any] {
		if x == "error" {
			return result.TryErr[any](errors.New("processing error"))
		}
		return result.Ok[any](x)
	})
	assert.True(t, res.IsErr())
}
