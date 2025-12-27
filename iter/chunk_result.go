package iter

import (
	"github.com/andeya/gust"
)

// ChunkResult represents a success (T) or failure (T) for NextChunk.
// This is used when the error type is the same as the success type (e.g., []T).
type ChunkResult[T any] struct {
	t gust.Option[T]
	e *T
}

// chunkOk wraps a successful chunk result.
//
//go:inline
func chunkOk[T any](ok T) ChunkResult[T] {
	return ChunkResult[T]{t: gust.Some(ok)}
}

// chunkErr wraps a failure chunk result.
//
//go:inline
func chunkErr[T any](err T) ChunkResult[T] {
	return ChunkResult[T]{e: &err}
}

// IsErr returns true if the result is an error.
//
//go:inline
func (r ChunkResult[T]) IsErr() bool {
	return r.e != nil
}

// IsOk returns true if the result is ok.
//
//go:inline
func (r ChunkResult[T]) IsOk() bool {
	return !r.IsErr()
}

// safeGetT safely gets the T value.
func (r ChunkResult[T]) safeGetT() T {
	if r.t.IsSome() {
		return r.t.UnwrapUnchecked()
	}
	var t T
	return t
}

// safeGetE safely gets the error value.
func (r ChunkResult[T]) safeGetE() T {
	if r.e != nil {
		return *r.e
	}
	var t T
	return t
}

// Unwrap returns the contained T value.
func (r ChunkResult[T]) Unwrap() T {
	if r.IsErr() {
		panic(gust.BoxErr(r.safeGetE()))
	}
	return r.safeGetT()
}

// UnwrapErr returns the contained error T value.
func (r ChunkResult[T]) UnwrapErr() T {
	if r.IsErr() {
		return r.safeGetE()
	}
	panic(gust.BoxErr(gust.Some(r.safeGetT())))
}
