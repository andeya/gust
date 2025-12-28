package syncutil

import (
	"sync/atomic"

	"github.com/andeya/gust"
)

// AtomicValue is a better generic-type wrapper for `atomic.Value`.
// A AtomicValue provides an atomic load and store of a consistently typed value.
// The zero value for a AtomicValue returns nil from Load.
// Once Store has been called, a AtomicValue must not be copied.
//
// A AtomicValue must not be copied after first use.
type AtomicValue[T any] struct {
	inner atomic.Value
}

// Load returns the value set by the most recent Store.
// It returns None if there has been no call to Store for this AtomicValue.
func (v *AtomicValue[T]) Load() (val gust.Option[T]) {
	return gust.AssertOpt[T](v.inner.Load())
}

// Store sets the value of the AtomicValue to x.
// All calls to Store for a given AtomicValue must use values of the same concrete type.
// Store of an inconsistent type panics, as does Store(nil).
func (v *AtomicValue[T]) Store(val T) {
	v.inner.Store(val)
}

// Swap stores new into AtomicValue and returns the previous value. It returns None if
// the AtomicValue is empty.
//
// All calls to Swap for a given AtomicValue must use values of the same concrete
// type. Swap of an inconsistent type panics, as does Swap(nil).
func (v *AtomicValue[T]) Swap(new T) (old gust.Option[T]) {
	return gust.AssertOpt[T](v.inner.Swap(new))
}

// CompareAndSwap executes the compare-and-swap operation for the AtomicValue.
//
// All calls to CompareAndSwap for a given AtomicValue must use values of the same
// concrete type. CompareAndSwap of an inconsistent type panics, as does
// CompareAndSwap(old, nil).
func (v *AtomicValue[T]) CompareAndSwap(old T, new T) (swapped bool) {
	return v.inner.CompareAndSwap(old, new)
}
