package iter

import (
	"errors"
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestTryFold(t *testing.T) {
	a := []int{1, 2, 3}
	sum := TryFold(FromSlice(a), 0, func(acc int, x int) gust.Result[int] {
		// Simulate checked addition
		if acc > 100 {
			return gust.Err[int](errors.New("overflow"))
		}
		return gust.Ok(acc + x)
	})
	assert.True(t, sum.IsOk())
	assert.Equal(t, 6, sum.Unwrap())

	// Test short-circuiting
	a2 := []int{10, 20, 30, 100, 40, 50}
	iter := FromSlice(a2)
	sum2 := TryFold(iter, 0, func(acc int, x int) gust.Result[int] {
		if acc+x > 50 {
			return gust.Err[int](errors.New("overflow"))
		}
		return gust.Ok(acc + x)
	})
	assert.True(t, sum2.IsErr())
}

func TestTryForEach(t *testing.T) {
	data := []string{"no_tea.txt", "stale_bread.json", "torrential_rain.png"}
	res := TryForEach(FromSlice(data), func(x string) gust.Result[any] {
		// Simulate processing
		return gust.Ok[any](nil)
	})
	assert.True(t, res.IsOk())
}

func TestTryReduce(t *testing.T) {
	numbers := []int{10, 20, 5, 23, 0}
	sum := FromSlice(numbers).TryReduce(func(x, y int) gust.Result[int] {
		// Simulate checked addition
		if x+y > 100 {
			return gust.Err[int](errors.New("overflow"))
		}
		return gust.Ok(x + y)
	})
	assert.True(t, sum.IsOk())
	assert.True(t, sum.Unwrap().IsSome())
	assert.Equal(t, 58, sum.Unwrap().Unwrap())

	// Test empty iterator
	numbers2 := []int{}
	sum2 := FromSlice(numbers2).TryReduce(func(x, y int) gust.Result[int] {
		return gust.Ok(x + y)
	})
	assert.True(t, sum2.IsOk())
	assert.True(t, sum2.Unwrap().IsNone())
}

func TestTryFind(t *testing.T) {
	a := []string{"1", "2", "lol", "NaN", "5"}
	result := FromSlice(a).TryFind(func(s string) gust.Result[bool] {
		if s == "lol" {
			return gust.Err[bool](errors.New("invalid"))
		}
		if s == "2" {
			return gust.Ok(true)
		}
		return gust.Ok(false)
	})
	assert.True(t, result.IsOk())
	assert.True(t, result.Unwrap().IsSome())
	assert.Equal(t, "2", result.Unwrap().Unwrap())
}

