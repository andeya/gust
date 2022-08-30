package iter

import (
	"github.com/andeya/gust"
)

var (
	_ gust.Iterable[any]    = (*DataVec[any])(nil)
	_ gust.IterableSizeHint = (*DataVec[any])(nil)
	_ gust.IterableCount    = (*DataVec[any])(nil)
	_ gust.DeIterable[uint] = (*DataVec[uint])(nil)
)

type DataVec[T any] struct {
	slice         []T
	nextIndex     int
	backNextIndex int
}

func NewDataVec[T any](slice []T) *DataVec[T] {
	return &DataVec[T]{
		slice:         slice,
		nextIndex:     0,
		backNextIndex: len(slice) - 1,
	}
}

func (v *DataVec[T]) ToSizeDeIterator() SizeDeIterator[T] {
	return FromSizeDeIterable[T](v)
}

func (v *DataVec[T]) Next() gust.Option[T] {
	if v.nextIndex <= v.backNextIndex {
		opt := gust.Some(v.slice[v.nextIndex])
		v.nextIndex++
		return opt
	}
	return gust.None[T]()
}

func (v *DataVec[T]) NextBack() gust.Option[T] {
	if v.backNextIndex >= 0 {
		opt := gust.Some(v.slice[v.backNextIndex])
		v.backNextIndex--
		return opt
	}
	return gust.None[T]()
}

func (v *DataVec[T]) SizeHint() (uint, gust.Option[uint]) {
	n := uint(v.backNextIndex - v.nextIndex + 1)
	return n, gust.Some(n)
}

func (v *DataVec[T]) Count() uint {
	v.nextIndex = v.backNextIndex
	return uint(v.backNextIndex - v.nextIndex + 1)
}

func (v *DataVec[T]) Remaining() uint {
	return uint(v.backNextIndex - v.nextIndex + 1)
}
