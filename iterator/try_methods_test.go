package iterator_test

import (
	"errors"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iterator"
	"github.com/stretchr/testify/assert"
)

func TestTryFold(t *testing.T) {
	a := []int{1, 2, 3}
	sum := iterator.TryFold(iterator.FromSlice(a), 0, func(acc int, x int) gust.Result[int] {
		// Simulate checked addition
		if acc > 100 {
			return gust.TryErr[int](errors.New("overflow"))
		}
		return gust.Ok(acc + x)
	})
	assert.True(t, sum.IsOk())
	assert.Equal(t, 6, sum.Unwrap())

	// Test short-circuiting
	a2 := []int{10, 20, 30, 100, 40, 50}
	iter := iterator.FromSlice(a2)
	sum2 := iterator.TryFold(iter, 0, func(acc int, x int) gust.Result[int] {
		if acc+x > 50 {
			return gust.TryErr[int](errors.New("overflow"))
		}
		return gust.Ok(acc + x)
	})
	assert.True(t, sum2.IsErr())
}

func TestTryForEach(t *testing.T) {
	data := []string{"no_tea.txt", "stale_bread.json", "torrential_rain.png"}
	res := iterator.TryForEach(iterator.FromSlice(data), func(x string) gust.Result[any] {
		// Simulate processing
		return gust.Ok[any](nil)
	})
	assert.True(t, res.IsOk())
}

func TestTryReduce(t *testing.T) {
	numbers := []int{10, 20, 5, 23, 0}
	sum := iterator.FromSlice(numbers).TryReduce(func(x, y int) gust.Result[int] {
		// Simulate checked addition
		if x+y > 100 {
			return gust.TryErr[int](errors.New("overflow"))
		}
		return gust.Ok(x + y)
	})
	assert.True(t, sum.IsOk())
	assert.True(t, sum.Unwrap().IsSome())
	assert.Equal(t, 58, sum.Unwrap().Unwrap())

	// Test empty iterator
	numbers2 := []int{}
	sum2 := iterator.FromSlice(numbers2).TryReduce(func(x, y int) gust.Result[int] {
		return gust.Ok(x + y)
	})
	assert.True(t, sum2.IsOk())
	assert.True(t, sum2.Unwrap().IsNone())
}

func TestTryFind(t *testing.T) {
	a := []string{"1", "2", "lol", "NaN", "5"}
	result := iterator.FromSlice(a).TryFind(func(s string) gust.Result[bool] {
		if s == "lol" {
			return gust.TryErr[bool](errors.New("invalid"))
		}
		if s == "2" {
			return gust.Ok(true)
		}
		return gust.Ok(false)
	})
	assert.True(t, result.IsOk())
	assert.True(t, result.Unwrap().IsSome())
	assert.Equal(t, "2", result.Unwrap().Unwrap())

	// Test TryFind with error before finding
	a2 := []string{"lol", "2", "3"}
	result2 := iterator.FromSlice(a2).TryFind(func(s string) gust.Result[bool] {
		if s == "lol" {
			return gust.TryErr[bool](errors.New("invalid"))
		}
		return gust.Ok(false)
	})
	assert.True(t, result2.IsErr())

	// Test TryFind with no match
	a3 := []string{"1", "3", "5"}
	result3 := iterator.FromSlice(a3).TryFind(func(s string) gust.Result[bool] {
		return gust.Ok(false)
	})
	assert.True(t, result3.IsOk())
	assert.True(t, result3.Unwrap().IsNone())

	// Test TryFind with empty iterator
	a4 := []string{}
	result4 := iterator.FromSlice(a4).TryFind(func(s string) gust.Result[bool] {
		return gust.Ok(true)
	})
	assert.True(t, result4.IsOk())
	assert.True(t, result4.Unwrap().IsNone())
}

func TestTryReduce_ErrorPath(t *testing.T) {
	// Test TryReduce with error
	numbers := []int{10, 20, 100}
	sum := iterator.FromSlice(numbers).TryReduce(func(x, y int) gust.Result[int] {
		if x+y > 100 {
			return gust.TryErr[int](errors.New("overflow"))
		}
		return gust.Ok(x + y)
	})
	assert.True(t, sum.IsErr())
}

func TestTryForEach_ErrorPath(t *testing.T) {
	// Test TryForEach with error
	data := []string{"a", "error", "c"}
	res := iterator.TryForEach(iterator.FromSlice(data), func(x string) gust.Result[any] {
		if x == "error" {
			return gust.TryErr[any](errors.New("processing error"))
		}
		return gust.Ok[any](x)
	})
	assert.True(t, res.IsErr())
}
