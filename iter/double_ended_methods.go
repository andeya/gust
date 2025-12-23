package iter

import (
	"github.com/andeya/gust"
)

// XRfold folds every element into an accumulator by applying an operation,
// starting from the back (any version).
//
// This is the reverse version of Fold(): it takes elements starting from
// the back of the iterator.
//
// # Examples
//
//	var a = []int{1, 2, 3}
//	var iter = FromSlice(a)
//	var deIter = iter.MustToDoubleEnded()
//	var sum = deIter.XRfold(0, func(acc any, x int) any { return acc.(int) + x })
//	assert.Equal(t, 6, sum)
//
//go:inline
func (de DoubleEndedIterator[T]) XRfold(init any, f func(any, T) any) any {
	return rfoldImpl(de.iter, init, f)
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
//	var iter = FromSlice(a)
//	var deIter = iter.MustToDoubleEnded()
//	var sum = deIter.Rfold(0, func(acc int, x int) int { return acc + x })
//	assert.Equal(t, 6, sum)
//
//go:inline
func (de DoubleEndedIterator[T]) Rfold(init T, f func(T, T) T) T {
	return rfoldImpl(de.iter, init, f)
}

// Rfold folds every element into an accumulator by applying an operation,
// starting from the back (generic version).
//
// # Examples
//
//	var a = []int{1, 2, 3}
//	var iter = FromSlice(a)
//	var deIter = iter.MustToDoubleEnded()
//	var sum = Rfold(deIter, 0, func(acc int, x int) int { return acc + x })
//	assert.Equal(t, 6, sum)
func Rfold[T any, B any](de DoubleEndedIterator[T], init B, f func(B, T) B) B {
	return rfoldImpl(de.iter, init, f)
}

// rfoldImpl is the internal implementation of Rfold.
//
//go:inline
func rfoldImpl[T any, B any](iter DoubleEndedIterable[T], init B, f func(B, T) B) B {
	accum := init
	for {
		item := iter.NextBack()
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
//	var iter = FromSlice(a)
//	var deIter = iter.MustToDoubleEnded()
//	var sum = deIter.XTryRfold(0, func(acc any, s string) gust.Result[any] {
//		if v, err := strconv.Atoi(s); err == nil {
//			return gust.Ok(any(acc.(int) + v))
//		}
//		return gust.Err[any](err)
//	})
//	assert.True(t, sum.IsOk())
//	assert.Equal(t, 6, sum.Unwrap().(int))
//
//go:inline
func (de DoubleEndedIterator[T]) XTryRfold(init any, f func(any, T) gust.Result[any]) gust.Result[any] {
	return tryRfoldImpl(de.iter, init, f)
}

// TryRfold is the reverse version of TryFold: it takes elements starting from
// the back of the iterator.
//
// # Examples
//
//	var a = []string{"1", "2", "3"}
//	var iter = FromSlice(a)
//	var deIter = iter.MustToDoubleEnded()
//	var sum = deIter.TryRfold(0, func(acc int, s string) gust.Result[int] {
//		if v, err := strconv.Atoi(s); err == nil {
//			return gust.Ok(acc + v)
//		}
//		return gust.Err[int](err)
//	})
//	assert.True(t, sum.IsOk())
//	assert.Equal(t, 6, sum.Unwrap())
//
//go:inline
func (de DoubleEndedIterator[T]) TryRfold(init T, f func(T, T) gust.Result[T]) gust.Result[T] {
	return tryRfoldImpl(de.iter, init, f)
}

// TryRfold is the reverse version of TryFold: it takes elements starting from
// the back of the iterator (generic version).
//
// # Examples
//
//	var a = []string{"1", "2", "3"}
//	var iter = FromSlice(a)
//	var deIter = iter.MustToDoubleEnded()
//	var sum = TryRfold(deIter, 0, func(acc int, s string) gust.Result[int] {
//		if v, err := strconv.Atoi(s); err == nil {
//			return gust.Ok(acc + v)
//		}
//		return gust.Err[int](err)
//	})
//	assert.True(t, sum.IsOk())
//	assert.Equal(t, 6, sum.Unwrap())
func TryRfold[T any, B any](de DoubleEndedIterator[T], init B, f func(B, T) gust.Result[B]) gust.Result[B] {
	return tryRfoldImpl(de.iter, init, f)
}

// tryRfoldImpl is the internal implementation of TryRfold.
//
//go:inline
func tryRfoldImpl[T any, B any](iter DoubleEndedIterable[T], init B, f func(B, T) gust.Result[B]) gust.Result[B] {
	accum := init
	for {
		item := iter.NextBack()
		if item.IsNone() {
			break
		}
		result := f(accum, item.Unwrap())
		if result.IsErr() {
			return result
		}
		accum = result.Unwrap()
	}
	return gust.Ok(accum)
}

// Rfind searches for an element of an iterator from the back that satisfies a predicate.
//
// Rfind() takes a closure that returns true or false. It applies
// this closure to each element of the iterator, starting at the end, and if any
// of them return true, then Rfind() returns Some(element). If they all return
// false, it returns None.
//
// Rfind() is short-circuiting; in other words, it will stop processing
// as soon as the closure returns true.
//
// # Examples
//
//	var a = []int{1, 2, 3}
//	var iter = FromSlice(a)
//	var deIter = AsDoubleEnded(iter)
//	assert.Equal(t, gust.Some(2), deIter.Rfind(func(x int) bool { return x == 2 }))
//
//	var b = []int{1, 2, 3}
//	var iter2 = FromSlice(b)
//	var deIter2 = AsDoubleEnded(iter2)
//	assert.Equal(t, gust.None[int](), deIter2.Rfind(func(x int) bool { return x == 5 }))
func (de DoubleEndedIterator[T]) Rfind(predicate func(T) bool) gust.Option[T] {
	return rfindImpl(de.iter, predicate)
}

// Rfind searches for an element of an iterator from the back that satisfies a predicate (function version).
//
// # Examples
//
//	var a = []int{1, 2, 3}
//	var iter = FromSlice(a)
//	var deIter = AsDoubleEnded(iter)
//	assert.Equal(t, gust.Some(2), Rfind(deIter, func(x int) bool { return x == 2 }))
func Rfind[T any](de DoubleEndedIterator[T], predicate func(T) bool) gust.Option[T] {
	return de.Rfind(predicate)
}

// rfindImpl is the internal implementation of Rfind.
//
//go:inline
func rfindImpl[T any](iter DoubleEndedIterable[T], predicate func(T) bool) gust.Option[T] {
	for {
		item := iter.NextBack()
		if item.IsNone() {
			return gust.None[T]()
		}
		if predicate(item.Unwrap()) {
			return item
		}
	}
}
