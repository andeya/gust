// Package iter provides a complete implementation of Rust's Iterator trait in Go.
// This package translates the Rust Iterator trait and its methods to Go, maintaining
// semantic equivalence while adapting to Go's type system and idioms.
package iterator

import (
	"github.com/andeya/gust/option"
	"github.com/andeya/gust/result"
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
