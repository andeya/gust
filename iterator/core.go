// Package iter provides a complete implementation of Rust's Iterator trait in Go.
// This package translates the Rust Iterator trait and its methods to Go, maintaining
// semantic equivalence while adapting to Go's type system and idioms.
package iterator

import (
	"github.com/andeya/gust"
)

// Iterable is the interface that extends gust.Iterable[T] and gust.IterableSizeHint.
// This is used internally for concrete iterator implementations.
type Iterable[T any] interface {
	gust.Iterable[T]
	gust.IterableSizeHint
}

// DefaultSizeHint provides a default implementation of SizeHint that returns (0, None).
// This can be used by iterator implementations that don't have size information.
//
//go:inline
func DefaultSizeHint[T any]() (uint, gust.Option[uint]) {
	return 0, gust.None[uint]()
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

var (
	// Verify that Iterator[T] implements Iterable[T]
	_ Iterable[any] = Iterator[any]{}
	// Verify that DoubleEndedIterator[T] implements Iterable[T]
	_ Iterable[any] = DoubleEndedIterator[any]{}
	// Verify that DoubleEndedIterator[T] implements DoubleEndedIterable[T]
	_ DoubleEndedIterable[any] = DoubleEndedIterator[any]{}
)

// DoubleEndedIterable is the interface for double-ended iterator implementations.
type DoubleEndedIterable[T any] interface {
	Iterable[T]                 // Includes Next() and SizeHint()
	gust.DoubleEndedIterable[T] // Includes NextBack() and Remaining()
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
//	assert.Equal(t, gust.Some(1), deIter.Next())
//	assert.Equal(t, gust.Some(6), deIter.NextBack())
//	assert.Equal(t, gust.Some(5), deIter.NextBack())
//	assert.Equal(t, gust.Some(2), deIter.Next())
//	assert.Equal(t, gust.Some(3), deIter.Next())
//	assert.Equal(t, gust.Some(4), deIter.Next())
//	assert.Equal(t, gust.None[int](), deIter.Next())
//	assert.Equal(t, gust.None[int](), deIter.NextBack())
//
//	// Can use all Iterator methods:
//	var filtered = deIter.Filter(func(x int) bool { return x > 2 })
//	var collected = filtered.Collect()
type DoubleEndedIterator[T any] struct {
	Iterator[T] // Embed Iterator to inherit all its methods
	iterable    DoubleEndedIterable[T]
}
