package iterator

import (
	"github.com/andeya/gust"
)

// TryFold is an iterator method that applies a function as long as it returns
// successfully, producing a single, final value.
//
// TryFold() takes two arguments: an initial value, and a closure with
// two arguments: an 'accumulator', and an element. The closure either
// returns successfully, with the value that the accumulator should have
// for the next iteration, or it returns failure, with an error value that
// is propagated back to the caller immediately (short-circuiting).
//
// The initial value is the value the accumulator will have on the first
// call. If applying the closure succeeded against every element of the
// iterator, TryFold() returns the final accumulator as success.
//
// Folding is useful whenever you have a collection of something, and want
// to produce a single value from it.
//
// # Examples
//
// Basic usage:
//
//	var a = []int{1, 2, 3}
//	// the checked sum of all of the elements of the array
//	var sum = TryFold(FromSlice(a), 0, func(acc int, x int) gust.Result[int] {
//		// Simulate checked addition
//		if acc > 100 {
//			return gust.Err[int](errors.New("overflow"))
//		}
//		return gust.Ok(acc + x)
//	})
//	assert.True(t, sum.IsOk())
//	assert.Equal(t, 6, sum.Unwrap())
//
// Short-circuiting:
//
//	var a = []int{10, 20, 30, 100, 40, 50}
//	var iter = FromSlice(a)
//	// This sum overflows when adding the 100 element
//	var sum = TryFold(iter, 0, func(acc int, x int) gust.Result[int] {
//		if acc+x > 50 {
//			return gust.Err[int](errors.New("overflow"))
//		}
//		return gust.Ok(acc + x)
//	})
//	assert.True(t, sum.IsErr())
//
//go:inline
func TryFold[T any, B any](iter Iterator[T], init B, f func(B, T) gust.Result[B]) gust.Result[B] {
	accum := init
	for {
		item := iter.Next()
		if item.IsNone() {
			return gust.Ok(accum)
		}
		result := f(accum, item.Unwrap())
		if result.IsErr() {
			return result
		}
		accum = result.Unwrap()
	}
}

// TryForEach is an iterator method that applies a fallible function to each item in the
// iterator, stopping at the first error and returning that error.
//
// This can also be thought of as the fallible form of ForEach()
// or as the stateless version of TryFold().
//
// # Examples
//
//	var data = []string{"no_tea.txt", "stale_bread.json", "torrential_rain.png"}
//	var res = TryForEach(FromSlice(data), func(x string) gust.Result[any] {
//		fmt.Println(x)
//		return gust.Ok[any](nil)
//	})
//	assert.True(t, res.IsOk())
//
//go:inline
func TryForEach[T any, B any](iter Iterator[T], f func(T) gust.Result[B]) gust.Result[B] {
	var zero B
	return TryFold(iter, zero, func(_ B, x T) gust.Result[B] {
		return f(x)
	})
}

//
//go:inline
func tryReduceImpl[T any](iter Iterator[T], f func(T, T) gust.Result[T]) gust.Result[gust.Option[T]] {
	first := iter.Next()
	if first.IsNone() {
		return gust.Ok(gust.None[T]())
	}

	result := TryFold(iter, first.Unwrap(), f)
	if result.IsErr() {
		return gust.TryErr[gust.Option[T]](result.UnwrapErr())
	}
	return gust.Ok(gust.Some(result.Unwrap()))
}

//
//go:inline
func tryFindImpl[T any](iter Iterator[T], f func(T) gust.Result[bool]) gust.Result[gust.Option[T]] {
	for {
		item := iter.Next()
		if item.IsNone() {
			return gust.Ok(gust.None[T]())
		}
		result := f(item.Unwrap())
		if result.IsErr() {
			return gust.TryErr[gust.Option[T]](result.UnwrapErr())
		}
		if result.Unwrap() {
			return gust.Ok(item)
		}
	}
}
