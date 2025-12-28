package iterator

import (
	"github.com/andeya/gust/option"
	"github.com/andeya/gust/result"
	"github.com/andeya/gust/void"
)

// Remaining returns the number of elements remaining in the iterator.
//
// # Examples
//
// var numbers = []int{1, 2, 3, 4, 5, 6}
// var deIter = FromSlice(numbers).MustToDoubleEnded()
// assert.Equal(t, uint(6), deIter.Remaining())
// deIter.Next()
// assert.Equal(t, uint(5), deIter.Remaining())
// deIter.NextBack()
// assert.Equal(t, uint(4), deIter.Remaining())
// deIter.NextBack()
// assert.Equal(t, uint(3), deIter.Remaining())
// deIter.NextBack()
// assert.Equal(t, uint(2), deIter.Remaining())
// deIter.NextBack()
// assert.Equal(t, uint(1), deIter.Remaining())
// deIter.NextBack()
// assert.Equal(t, uint(0), deIter.Remaining())
// deIter.NextBack()
// assert.Equal(t, uint(0), deIter.Remaining())
func (de DoubleEndedIterator[T]) Remaining() uint {
	return de.iterable.Remaining()
}

// NextBack removes and returns an element from the end of the iterator.
//
// Returns None when there are no more elements.
//
// # Examples
//
//	var numbers = []int{1, 2, 3, 4, 5, 6}
//	var deIter = FromSlice(numbers).MustToDoubleEnded()
//	assert.Equal(t, option.Some(6), deIter.NextBack())
//	assert.Equal(t, option.Some(5), deIter.NextBack())
//	assert.Equal(t, option.Some(4), deIter.NextBack())
//	assert.Equal(t, option.Some(3), deIter.NextBack())
//	assert.Equal(t, option.Some(2), deIter.NextBack())
//	assert.Equal(t, option.Some(1), deIter.NextBack())
//	assert.Equal(t, option.None[int](), deIter.NextBack())
//
//go:inline
func (de DoubleEndedIterator[T]) NextBack() option.Option[T] {
	return de.iterable.NextBack()
}

// AdvanceBackBy advances the iterator from the back by n elements.
//
// AdvanceBackBy is the reverse version of AdvanceBy. This method will
// eagerly skip n elements starting from the back by calling NextBack up
// to n times until None is encountered.
//
// AdvanceBackBy(n) will return Ok[Void](nil) if the iterator successfully advances by
// n elements, or Err[Void](k) with value k if None is encountered, where k
// is remaining number of steps that could not be advanced because the iterator ran out.
// If iter is empty and n is non-zero, then this returns Err[Void](n).
// Otherwise, k is always less than n.
//
// Calling AdvanceBackBy(0) can do meaningful work.
//
// # Examples
//
//	var a = []int{3, 4, 5, 6}
//	var deIter = FromSlice(a).MustToDoubleEnded()
//	assert.True(t, deIter.AdvanceBackBy(2).IsOk())
//	assert.Equal(t, option.Some(4), deIter.NextBack())
//	assert.True(t, deIter.AdvanceBackBy(0).IsOk())
//	assert.True(t, deIter.AdvanceBackBy(100).IsErr())
func (de DoubleEndedIterator[T]) AdvanceBackBy(n uint) result.VoidResult {
	for i := uint(0); i < n; i++ {
		if de.iterable.NextBack().IsNone() {
			return result.TryErr[void.Void](n - i)
		}
	}
	return result.Ok[void.Void](nil)
}

// NthBack returns the nth element from the end of the iterator.
//
// This is essentially the reversed version of Nth().
// Although like most indexing operations, the count starts from zero, so
// NthBack(0) returns the first value from the end, NthBack(1) the
// second, and so on.
//
// Note that all elements between the end and the returned element will be
// consumed, including the returned element. This also means that calling
// NthBack(0) multiple times on the same iterator will return different
// elements.
//
// NthBack() will return None if n is greater than or equal to the length of the
// iterator.
//
// # Examples
//
//	var a = []int{1, 2, 3}
//	var deIter = FromSlice(a).MustToDoubleEnded()
//	assert.Equal(t, option.Some(1), deIter.NthBack(2))
//	assert.Equal(t, option.Some(2), deIter.NthBack(1))
//	assert.Equal(t, option.Some(3), deIter.NthBack(0))
//	assert.Equal(t, option.None[int](), deIter.NthBack(10))
func (de DoubleEndedIterator[T]) NthBack(n uint) option.Option[T] {
	if de.AdvanceBackBy(n).IsErr() {
		return option.None[T]()
	}
	return de.iterable.NextBack()
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
//	var deIter = FromSlice(a).MustToDoubleEnded()
//	assert.Equal(t, option.Some(2), deIter.Rfind(func(x int) bool { return x == 2 }))
//
//	var b = []int{1, 2, 3}
//	var deIter2 = FromSlice(b).MustToDoubleEnded()
//	assert.Equal(t, option.None[int](), deIter2.Rfind(func(x int) bool { return x == 5 }))
func (de DoubleEndedIterator[T]) Rfind(predicate func(T) bool) option.Option[T] {
	return rfindImpl(de.iterable, predicate)
}

// Rfind searches for an element of an iterator from the back that satisfies a predicate (function version).
//
// # Examples
//
//	var a = []int{1, 2, 3}
//	var deIter = FromSlice(a).MustToDoubleEnded()
//	assert.Equal(t, option.Some(2), Rfind(deIter, func(x int) bool { return x == 2 }))
func Rfind[T any](de DoubleEndedIterator[T], predicate func(T) bool) option.Option[T] {
	return de.Rfind(predicate)
}

// rfindImpl is the internal implementation of Rfind.
//
//go:inline
func rfindImpl[T any](iter DoubleEndedIterable[T], predicate func(T) bool) option.Option[T] {
	de := DoubleEndedIterator[T]{iterable: iter}
	for {
		item := de.NextBack()
		if item.IsNone() {
			return option.None[T]()
		}
		if predicate(item.Unwrap()) {
			return item
		}
	}
}
