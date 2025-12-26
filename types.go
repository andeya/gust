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

// Split splits the pair into its two components.
//
//go:inline
func (p Pair[A, B]) Split() (A, B) {
	return p.A, p.B
}

// VecEntry is an index-element entry of slice or array.
type VecEntry[T any] struct {
	Index int
	Elem  T
}

// Split splits the vector entry into its two components.
//
//go:inline
func (v VecEntry[T]) Split() (int, T) {
	return v.Index, v.Elem
}

// DictEntry is a key-value entry of map.
type DictEntry[K comparable, V any] struct {
	Key   K
	Value V
}

// Split splits the dictionary entry into its two components.
//
//go:inline
func (d DictEntry[K, V]) Split() (K, V) {
	return d.Key, d.Value
}
