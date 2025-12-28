// Package iter provides a complete implementation of Rust's Iterator trait in Go.
// This package translates the Rust Iterator trait and its methods to Go, maintaining
// semantic equivalence while adapting to Go's type system and idioms.
package iterator

import (
	"iter"

	"github.com/andeya/gust/option"
	"github.com/andeya/gust/result"
	"github.com/andeya/gust/void"
)

// Iterable represents a sequence of values that can be iterated over.
type Iterable[T any] interface {
	// Next returns the next value in the sequence.
	// Returns Some(value) if there is a next value, None otherwise.
	Next() option.Option[T]
	// SizeHint returns a hint about the remaining size of the sequence.
	// Returns (lower_bound, Some(upper_bound)) if the size is known,
	// or (lower_bound, None) if the upper bound is unknown.
	SizeHint() (uint, option.Option[uint])
}

// SizeIterable represents an iterable that knows its remaining size.
type SizeIterable[T any] interface {
	Iterable[T]
	// Remaining returns the number of remaining elements.
	Remaining() uint
}

// DoubleEndedIterable represents an iterable that can be iterated from both ends.
type DoubleEndedIterable[T any] interface {
	Iterable[T]
	SizeIterable[T]
	// NextBack returns the next value from the back of the sequence.
	// Returns Some(value) if there is a next value, None otherwise.
	NextBack() option.Option[T]
}

// IterableCount represents an iterable that can count its elements.
type IterableCount interface {
	// Count returns the total number of elements.
	Count() uint
}

// IterableSizeHint represents an iterable that can provide size hints.
type IterableSizeHint interface {
	// SizeHint returns a hint about the remaining size of the sequence.
	// Returns (lower_bound, Some(upper_bound)) if the size is known,
	// or (lower_bound, None) if the upper bound is unknown.
	SizeHint() (uint, option.Option[uint])
}

// DefaultSizeHint provides a default implementation of SizeHint that returns (0, None).
// This can be used by iterator implementations that don't have size information.
//
//go:inline
func DefaultSizeHint[T any]() (uint, option.Option[uint]) {
	return 0, option.None[uint]()
}

// Iterator is the main iterator type that provides method chaining.
// It wraps an Iterable[T] interface and provides all adapter methods as struct methods,
// enabling Rust-like method chaining: iterator.Map(...).Filter(...).Collect()
//
// Iterator implements gust.Iterable[T] and gust.IterableSizeHint, so it can be used
// anywhere these interfaces are expected.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3, 4, 5})
//	var result = iterator.Filter(func(x int) bool { return x > 2 }).
//	              Take(2).
//	              Collect()
type Iterator[T any] struct {
	iterable Iterable[T]
}

// DoubleEndedIterator is the main double-ended iterator type that provides method chaining.
// It embeds Iterator[T] to inherit all Iterator methods, and adds double-ended specific methods.
// This enables Rust-like method chaining: deIter.Filter(...).NextBack().Rfold(...)
//
// DoubleEndedIterator implements gust.Iterable[T] and gust.IterableSizeHint, so it can be used
// anywhere these interfaces are expected. It also inherits all Iterator[T] methods.
//
// # Examples
//
// Basic usage:
//
//	var numbers = []int{1, 2, 3, 4, 5, 6}
//	var deIter = FromSlice(numbers).MustToDoubleEnded()
//
//	import "github.com/andeya/gust/option"
//	assert.Equal(t, option.Some(1), deIter.Next())
//	assert.Equal(t, option.Some(6), deIter.NextBack())
//	assert.Equal(t, option.Some(5), deIter.NextBack())
//	assert.Equal(t, option.Some(2), deIter.Next())
//	assert.Equal(t, option.Some(3), deIter.Next())
//	assert.Equal(t, option.Some(4), deIter.Next())
//	assert.Equal(t, option.None[int](), deIter.Next())
//	assert.Equal(t, option.None[int](), deIter.NextBack())
//
//	// Can use all Iterator methods:
//	var filtered = deIter.Filter(func(x int) bool { return x > 2 })
//	var collected = filtered.Collect()
type DoubleEndedIterator[T any] struct {
	Iterator[T] // Embed Iterator to inherit all its methods
	iterable    DoubleEndedIterable[T]
}

var (
	// Verify that Iterator[T] implements Iterable[T]
	_ Iterable[any] = Iterator[any]{}
	// Verify that DoubleEndedIterator[T] implements Iterable[T]
	_ Iterable[any] = DoubleEndedIterator[any]{}
	// Verify that DoubleEndedIterator[T] implements DoubleEndedIterable[T]
	_ DoubleEndedIterable[any] = DoubleEndedIterator[any]{}
	// Verify that Result[T] implements Iterable[T]
	_ Iterable[any] = new(result.Result[any])
	// Verify that Result[T] implements DoubleEndedIterable[T]
	_ DoubleEndedIterable[any] = new(result.Result[any])
	// Verify that Option[T] implements Iterable[T]
	_ Iterable[any] = new(option.Option[any])
	// Verify that Option[T] implements DoubleEndedIterable[T]
	_ DoubleEndedIterable[any] = new(option.Option[any])
)

// Next advances the iterator and returns the next value.
// This implements gust.Iterable[T] interface.
//
//go:inline
func (it Iterator[T]) Next() option.Option[T] {
	return it.iterable.Next()
}

// SizeHint returns the bounds on the remaining length of the iterator.
// This implements gust.IterableSizeHint interface.
//
//go:inline
func (it Iterator[T]) SizeHint() (uint, option.Option[uint]) {
	return it.iterable.SizeHint()
}

// MustToDoubleEnded converts to a DoubleEndedIterator[T] if the underlying
// iterator supports double-ended iteration. Otherwise, it panics.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var deIter = iterator.MustToDoubleEnded()
//	assert.Equal(t, option.Some(3), deIter.NextBack())
//	// Can use Iterator methods:
//	var doubled = deIter.Map(func(x int) any { return x * 2 })
func (it Iterator[T]) MustToDoubleEnded() DoubleEndedIterator[T] {
	if deCore, ok := it.iterable.(DoubleEndedIterable[T]); ok {
		return DoubleEndedIterator[T]{
			Iterator: Iterator[T]{iterable: deCore}, // Embed Iterator with the same core
			iterable: deCore,
		}
	}
	panic("iterator does not support double-ended iteration")
}

// TryToDoubleEnded converts to a DoubleEndedIterator[T] if the underlying
// iterator supports double-ended iteration. Otherwise, it returns None.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var deIter = iterator.TryToDoubleEnded()
//	assert.Equal(t, option.Some(3), deIter.NextBack())
//	// Can use Iterator methods:
//	var doubled = deIter.Map(func(x int) any { return x * 2 })
func (it Iterator[T]) TryToDoubleEnded() option.Option[DoubleEndedIterator[T]] {
	if deCore, ok := it.iterable.(DoubleEndedIterable[T]); ok {
		return option.Some(DoubleEndedIterator[T]{
			Iterator: Iterator[T]{iterable: deCore}, // Embed Iterator with the same core
			iterable: deCore,
		})
	}
	return option.None[DoubleEndedIterator[T]]()
}

// Seq converts the Iterator[T] to Go's standard iterator.Seq[T].
// This allows using gust iterators with Go's built-in iteration support (for loops).
//
// # Examples
//
//	iter := FromSlice([]int{1, 2, 3})
//	for v := range iterator.Seq() {
//		fmt.Println(v) // prints 1, 2, 3
//	}
func (it Iterator[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			opt := it.Next()
			if opt.IsNone() {
				return
			}
			if !yield(opt.Unwrap()) {
				return
			}
		}
	}
}

// Seq2 converts the Iterator[T] to Go's standard iterator.Seq2[T].
// This allows using gust iterators with Go's built-in iteration support (for loops).
//
// # Examples
//
// iter := FromSlice([]int{1, 2, 3})
//
//	for k, v := range iterator.Seq2() {
//		fmt.Println(k, v) // prints 0 1, 1 2, 2 3
//	}
func (it Iterator[T]) Seq2() iter.Seq2[uint, T] {
	pairIter := enumerateImpl(it.iterable)
	return func(yield func(uint, T) bool) {
		for {
			opt := pairIter.Next()
			if opt.IsNone() {
				return
			}
			pair := opt.Unwrap()
			if !yield(pair.A, pair.B) {
				return
			}
		}
	}
}

// Pull converts the Iterator[T] to a pull-style iterator using Go's standard iterator.Pull.
// This returns two functions: next (to pull values) and stop (to stop iteration).
// The caller should defer stop() to ensure proper cleanup.
//
// # Examples
//
//	iter := FromSlice([]int{1, 2, 3, 4, 5})
//	next, stop := iterator.Pull()
//	defer stop()
//
//	// Pull values manually
//	for {
//		v, ok := next()
//		if !ok {
//			break
//		}
//		fmt.Println(v)
//		if v == 3 {
//			break // Early termination
//		}
//	}
func (it Iterator[T]) Pull() (next func() (T, bool), stop func()) {
	return func() (T, bool) {
			return it.Next().Split()
		}, func() {
			// No need to stop here since the iterator will clean up automatically
		}
}

// Pull2 converts the Iterator[T] to a pull-style iterator using Go's standard iterator.Pull2.
// This returns two functions: next (to pull key-value pairs) and stop (to stop iteration).
// The caller should defer stop() to ensure proper cleanup.
//
// # Examples
//
//	iter := FromSlice([]int{1, 2, 3, 4, 5})
//	next, stop := iterator.Pull2()
//	defer stop()
//
//	// Pull key-value pairs manually
//	for {
//		k, v, ok := next()
//		if !ok {
//			break
//		}
//		fmt.Println(k, v)
//		if v == 3 {
//			break // Early termination
//		}
//	}
func (it Iterator[T]) Pull2() (next func() (uint, T, bool), stop func()) {
	pairIter := enumerateImpl(it.iterable)
	return func() (uint, T, bool) {
			pair, ok := pairIter.Next().Split()
			return pair.A, pair.B, ok
		}, func() {
			// No need to stop here since the iterator will clean up automatically
		}
}

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
