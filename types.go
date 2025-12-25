// Package gust provides Rust-inspired error handling, optional values, and iteration utilities for Go.
// This file contains common types and interfaces used throughout the package.
package gust

// For implements iterators, the following methods are available:

type (
	// Iterable represents a sequence of values that can be iterated over.
	Iterable[T any] interface {
		Next() Option[T]
	}
	// SizeIterable represents an iterable that knows its remaining size.
	SizeIterable[T any] interface {
		Remaining() uint
	}
	// DoubleEndedIterable represents an iterable that can be iterated from both ends.
	DoubleEndedIterable[T any] interface {
		Iterable[T]
		SizeIterable[T]
		NextBack() Option[T]
	}
	// IterableCount represents an iterable that can count its elements.
	IterableCount interface {
		Count() uint
	}
	// IterableSizeHint represents an iterable that can provide size hints.
	IterableSizeHint interface {
		SizeHint() (uint, Option[uint])
	}
)

// Pair is a pair of values.
type Pair[A any, B any] struct {
	A A
	B B
}

// VecEntry is an index-element entry of slice or array.
type VecEntry[T any] struct {
	Index int
	Elem  T
}

// DictEntry is a key-value entry of map.
type DictEntry[K comparable, V any] struct {
	Key   K
	Value V
}
