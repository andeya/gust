package iterator

import (
	"github.com/andeya/gust/option"
	"github.com/andeya/gust/result"
)

// Fold folds every element into an accumulator by applying an operation,
// returning the final result.
//
// Fold() takes two arguments: an initial value, and a closure with two
// arguments: an 'accumulator', and an element. The closure returns the value that
// the accumulator should have for the next iteration.
//
// The initial value is the value the accumulator will have on the first
// call.
//
// After applying this closure to every element of the iterator, Fold()
// returns the accumulator.
//
// This operation is sometimes called 'reduce' or 'inject'.
//
// Folding is useful whenever you have a collection of something, and want
// to produce a single value from it.
//
// Note: Fold(), and similar methods that traverse the entire iterator,
// might not terminate for infinite iterators, even on traits for which a
// result is determinable in finite time.
//
// Note: Reduce() can be used to use the first element as the initial
// value, if the accumulator type and item type is the same.
//
// Note: Fold() combines elements in a *left-associative* fashion. For associative
// operators like +, the order the elements are combined in is not important, but for non-associative
// operators like - the order will affect the final result.
//
// # Examples
//
// Basic usage:
//
//	var a = []int{1, 2, 3}
//	// the sum of all of the elements of the array
//	var sum = Fold(FromSlice(a), 0, func(acc int, x int) int { return acc + x })
//	assert.Equal(t, 6, sum)
//
// Let's walk through each step of the iteration here:
//
// | element | acc | x | result |
// |---------|-----|---|--------|
// |         | 0   |   |        |
// | 1       | 0   | 1 | 1      |
// | 2       | 1   | 2 | 3      |
// | 3       | 3   | 3 | 6      |
//
// And so, our final result, 6.
//
//go:inline
func Fold[T any, B any](iter Iterator[T], init B, f func(B, T) B) B {
	return foldImpl(iter.iterable, init, f)
}

//go:inline
func foldImpl[T any, B any](iter Iterable[T], init B, f func(B, T) B) B {
	accum := init
	it := Iterator[T]{iterable: iter}
	for {
		item := it.Next()
		if item.IsNone() {
			break
		}
		accum = f(accum, item.Unwrap())
	}
	return accum
}

//go:inline
func reduceImpl[T any](iter Iterable[T], f func(T, T) T) option.Option[T] {
	it := Iterator[T]{iterable: iter}
	first := it.Next()
	if first.IsNone() {
		return option.None[T]()
	}
	result := Fold(Iterator[T]{iterable: iter}, first.Unwrap(), f)
	return option.Some(result)
}

//go:inline
func forEachImpl[T any](iter Iterable[T], f func(T)) {
	it := Iterator[T]{iterable: iter}
	for {
		item := it.Next()
		if item.IsNone() {
			break
		}
		f(item.Unwrap())
	}
}

//go:inline
func collectImpl[T any](iter Iterable[T]) []T {
	it := Iterator[T]{iterable: iter}
	lower, upper := it.SizeHint()
	var capacity = lower
	if upper.IsSome() && upper.Unwrap() > lower {
		capacity = upper.Unwrap()
	}
	var result = make([]T, 0, capacity)
	for {
		item := it.Next()
		if item.IsNone() {
			break
		}
		result = append(result, item.Unwrap())
	}
	return result
}

//go:inline
func partitionImpl[T any](iter Iterable[T], f func(T) bool) (truePart []T, falsePart []T) {
	var left []T
	var right []T
	it := Iterator[T]{iterable: iter}
	for {
		item := it.Next()
		if item.IsNone() {
			break
		}
		if f(item.Unwrap()) {
			left = append(left, item.Unwrap())
		} else {
			right = append(right, item.Unwrap())
		}
	}
	// Return empty slices instead of nil for consistency with design intent
	if left == nil {
		left = []T{}
	}
	if right == nil {
		right = []T{}
	}
	return left, right
}

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
//	var sum = TryFold(FromSlice(a), 0, func(acc int, x int) result.Result[int] {
//		// Simulate checked addition
//		if acc > 100 {
//			return result.TryErr[int](errors.New("overflow"))
//		}
//		return result.Ok(acc + x)
//	})
//	assert.True(t, sum.IsOk())
//	assert.Equal(t, 6, sum.Unwrap())
//
// Short-circuiting:
//
//	var a = []int{10, 20, 30, 100, 40, 50}
//	var iter = FromSlice(a)
//	// This sum overflows when adding the 100 element
//	var sum = TryFold(iter, 0, func(acc int, x int) result.Result[int] {
//		if acc+x > 50 {
//			return result.TryErr[int](errors.New("overflow"))
//		}
//		return result.Ok(acc + x)
//	})
//	assert.True(t, sum.IsErr())
//
//go:inline
func TryFold[T any, B any](iter Iterator[T], init B, f func(B, T) result.Result[B]) result.Result[B] {
	accum := init
	for {
		item := iter.Next()
		if item.IsNone() {
			return result.Ok(accum)
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
//	var res = TryForEach(FromSlice(data), func(x string) result.Result[any] {
//		fmt.Println(x)
//		return result.Ok[any](nil)
//	})
//	assert.True(t, res.IsOk())
//
//go:inline
func TryForEach[T any, B any](iter Iterator[T], f func(T) result.Result[B]) result.Result[B] {
	var zero B
	return TryFold(iter, zero, func(_ B, x T) result.Result[B] {
		return f(x)
	})
}

//go:inline
func tryReduceImpl[T any](iter Iterator[T], f func(T, T) result.Result[T]) result.Result[option.Option[T]] {
	first := iter.Next()
	if first.IsNone() {
		return result.Ok(option.None[T]())
	}

	res := TryFold(iter, first.Unwrap(), f)
	if res.IsErr() {
		return result.TryErr[option.Option[T]](res.UnwrapErr())
	}
	return result.Ok(option.Some(res.Unwrap()))
}

// XFold folds every element into an accumulator.
// This wrapper method allows XFold to be called as a method.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var sum = iterator.XFold(0, func(acc any, x int) any { return acc.(int) + x })
//	assert.Equal(t, 6, sum)
//
//go:inline
func (it Iterator[T]) XFold(init any, f func(any, T) any) any {
	return Fold(it, init, f)
}

// Fold folds every element into an accumulator.
// This wrapper method allows Fold to be called as a method.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var sum = iterator.Fold(0, func(acc int, x int) int { return acc + x })
//	assert.Equal(t, 6, sum)
//
//go:inline
func (it Iterator[T]) Fold(init T, f func(T, T) T) T {
	return Fold(it, init, f)
}

// XTryFold applies a function as long as it returns successfully, producing a single, final value (any version).
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var sum = iterator.XTryFold(0, func(acc any, x int) result.Result[any] {
//		return result.Ok[any](any(acc.(int) + x))
//	})
//	assert.True(t, sum.IsOk())
//	assert.Equal(t, 6, sum.Unwrap().(int))
//
//go:inline
func (it Iterator[T]) XTryFold(init any, f func(any, T) result.Result[any]) result.Result[any] {
	return TryFold(it, init, f)
}

// TryFold applies a function as long as it returns successfully, producing a single, final value.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var sum = iterator.TryFold(0, func(acc int, x int) result.Result[int] {
//		return result.Ok(acc + x)
//	})
//	assert.True(t, sum.IsOk())
//	assert.Equal(t, 6, sum.Unwrap())
//
//go:inline
func (it Iterator[T]) TryFold(init T, f func(T, T) result.Result[T]) result.Result[T] {
	return TryFold(it, init, f)
}

// XRfold folds every element into an accumulator by applying an operation,
// starting from the back (any version).
//
// This is the reverse version of Fold(): it takes elements starting from
// the back of the iterator.
//
// # Examples
//
//	var a = []int{1, 2, 3}
//	var deIter = FromSlice(a).MustToDoubleEnded()
//	var deIter = iterator.MustToDoubleEnded()
//	var sum = deIter.XRfold(0, func(acc any, x int) any { return acc.(int) + x })
//	assert.Equal(t, 6, sum)
//
//go:inline
func (de DoubleEndedIterator[T]) XRfold(init any, f func(any, T) any) any {
	return rfoldImpl(de.iterable, init, f)
}

// Rfold folds every element into an accumulator by applying an operation,
// starting from the back.
//
// This is the reverse version of Fold(): it takes elements starting from
// the back of the iterator.
//
// Rfold() takes two arguments: an initial value, and a closure with two
// arguments: an 'accumulator', and an element. The closure returns the value that
// the accumulator should have for the next iteration.
//
// The initial value is the value the accumulator will have on the first
// call.
//
// After applying this closure to every element of the iterator, Rfold()
// returns the accumulator.
//
// Note: Rfold() combines elements in a *right-associative* fashion. For associative
// operators like +, the order the elements are combined in is not important, but for non-associative
// operators like - the order will affect the final result.
//
// # Examples
//
//	var a = []int{1, 2, 3}
//	var deIter = FromSlice(a).MustToDoubleEnded()
//	var sum = deIter.Rfold(0, func(acc int, x int) int { return acc + x })
//	assert.Equal(t, 6, sum)
//
//go:inline
func (de DoubleEndedIterator[T]) Rfold(init T, f func(T, T) T) T {
	return rfoldImpl(de.iterable, init, f)
}

// Rfold folds every element into an accumulator by applying an operation,
// starting from the back (generic version).
//
// # Examples
//
//	var a = []int{1, 2, 3}
//	var deIter = FromSlice(a).MustToDoubleEnded()
//	var sum = Rfold(deIter, 0, func(acc int, x int) int { return acc + x })
//	assert.Equal(t, 6, sum)
func Rfold[T any, B any](de DoubleEndedIterator[T], init B, f func(B, T) B) B {
	return rfoldImpl(de.iterable, init, f)
}

// rfoldImpl is the internal implementation of Rfold.
//
//go:inline
func rfoldImpl[T any, B any](iter DoubleEndedIterable[T], init B, f func(B, T) B) B {
	accum := init
	de := DoubleEndedIterator[T]{iterable: iter}
	for {
		item := de.NextBack()
		if item.IsNone() {
			break
		}
		accum = f(accum, item.Unwrap())
	}
	return accum
}

// XTryRfold is the reverse version of TryFold: it takes elements starting from
// the back of the iterator (any version).
//
// # Examples
//
//	var a = []string{"1", "2", "3"}
//	var deIter = FromSlice(a).MustToDoubleEnded()
//	var sum = deIter.XTryRfold(0, func(acc any, s string) result.Result[any] {
//		if v, err := strconv.Atoi(s); err == nil {
//			return result.Ok(any(acc.(int) + v))
//		}
//		return result.TryErr[any](err)
//	})
//	assert.True(t, sum.IsOk())
//	assert.Equal(t, 6, sum.Unwrap().(int))
//
//go:inline
func (de DoubleEndedIterator[T]) XTryRfold(init any, f func(any, T) result.Result[any]) result.Result[any] {
	return tryRfoldImpl(de.iterable, init, f)
}

// TryRfold is the reverse version of TryFold: it takes elements starting from
// the back of the iterator.
//
// # Examples
//
//	var a = []string{"1", "2", "3"}
//	var deIter = FromSlice(a).MustToDoubleEnded()
//	var sum = deIter.TryRfold(0, func(acc int, s string) result.Result[int] {
//		if v, err := strconv.Atoi(s); err == nil {
//			return result.Ok(acc + v)
//		}
//		return result.TryErr[int](err)
//	})
//	assert.True(t, sum.IsOk())
//	assert.Equal(t, 6, sum.Unwrap())
//
//go:inline
func (de DoubleEndedIterator[T]) TryRfold(init T, f func(T, T) result.Result[T]) result.Result[T] {
	return tryRfoldImpl(de.iterable, init, f)
}

// TryRfold is the reverse version of TryFold: it takes elements starting from
// the back of the iterator (generic version).
//
// # Examples
//
//	var a = []string{"1", "2", "3"}
//	var deIter = FromSlice(a).MustToDoubleEnded()
//	var sum = TryRfold(deIter, 0, func(acc int, s string) result.Result[int] {
//		if v, err := strconv.Atoi(s); err == nil {
//			return result.Ok(acc + v)
//		}
//		return result.TryErr[int](err)
//	})
//	assert.True(t, sum.IsOk())
//	assert.Equal(t, 6, sum.Unwrap())
func TryRfold[T any, B any](de DoubleEndedIterator[T], init B, f func(B, T) result.Result[B]) result.Result[B] {
	return tryRfoldImpl(de.iterable, init, f)
}

// tryRfoldImpl is the internal implementation of TryRfold.
//
//go:inline
func tryRfoldImpl[T any, B any](iter DoubleEndedIterable[T], init B, f func(B, T) result.Result[B]) result.Result[B] {
	accum := init
	de := DoubleEndedIterator[T]{iterable: iter}
	for {
		item := de.NextBack()
		if item.IsNone() {
			break
		}
		result := f(accum, item.Unwrap())
		if result.IsErr() {
			return result
		}
		accum = result.Unwrap()
	}
	return result.Ok(accum)
}
