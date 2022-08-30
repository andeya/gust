package iter

import (
	"github.com/andeya/gust"
)

var (
	_ gust.Iterable[any]    = (*IterableVec[any])(nil)
	_ gust.IterableSizeHint = (*IterableVec[any])(nil)
	_ gust.IterableCount    = (*IterableVec[any])(nil)
	_ gust.DeIterable[uint] = (*IterableVec[uint])(nil)
)

type IterableVec[T any] struct {
	slice         []T
	nextIndex     int
	backNextIndex int
}

func NewIterableVec[T any](slice []T) *IterableVec[T] {
	return &IterableVec[T]{
		slice:         slice,
		nextIndex:     0,
		backNextIndex: len(slice) - 1,
	}
}

func (v *IterableVec[T]) ToSizeDeIterator() SizeDeIterator[T] {
	return FromSizeDeIterable[T](v)
}

func (v *IterableVec[T]) Next() gust.Option[T] {
	if v.nextIndex <= v.backNextIndex {
		opt := gust.Some(v.slice[v.nextIndex])
		v.nextIndex++
		return opt
	}
	return gust.None[T]()
}

func (v *IterableVec[T]) NextBack() gust.Option[T] {
	if v.backNextIndex >= 0 {
		opt := gust.Some(v.slice[v.backNextIndex])
		v.backNextIndex--
		return opt
	}
	return gust.None[T]()
}

func (v *IterableVec[T]) SizeHint() (uint, gust.Option[uint]) {
	n := uint(v.backNextIndex - v.nextIndex + 1)
	return n, gust.Some(n)
}

func (v *IterableVec[T]) Count() uint {
	v.nextIndex = v.backNextIndex
	return uint(v.backNextIndex - v.nextIndex + 1)
}

func (v *IterableVec[T]) Remaining() uint {
	return uint(v.backNextIndex - v.nextIndex + 1)
}
