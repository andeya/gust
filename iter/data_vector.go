package iter

import (
	"github.com/andeya/gust"
)

var (
	_ DataForIter[any]      = (*DataVec[any])(nil)
	_ SizeHintForIter       = (*DataVec[any])(nil)
	_ CountForIter          = (*DataVec[any])(nil)
	_ NextBackForIter[uint] = (*DataVec[uint])(nil)
	_ RemainingLenForIter   = (*DataVec[uint])(nil)
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

func (v *DataVec[T]) ToIterator() Iterator[T] {
	return newIter[T](v)
}

func (v *DataVec[T]) ToDoubleEndedIterator() DoubleEndedIterator[T] {
	return newDoubleEndedIter[T](v)
}

func (v *DataVec[T]) NextForIter() gust.Option[T] {
	if v.nextIndex <= v.backNextIndex {
		opt := gust.Some(v.slice[v.nextIndex])
		v.nextIndex++
		return opt
	}
	return gust.None[T]()
}

func (v *DataVec[T]) NextBackForIter() gust.Option[T] {
	if v.backNextIndex >= 0 {
		opt := gust.Some(v.slice[v.backNextIndex])
		v.backNextIndex--
		return opt
	}
	return gust.None[T]()
}

func (v *DataVec[T]) SizeHintForIter() (uint, gust.Option[uint]) {
	n := uint(v.backNextIndex - v.nextIndex + 1)
	return n, gust.Some(n)
}

func (v *DataVec[T]) CountForIter() uint {
	v.nextIndex = v.backNextIndex
	return uint(v.backNextIndex - v.nextIndex + 1)
}

func (v *DataVec[T]) RemainingLenForIter() uint {
	return uint(v.backNextIndex - v.nextIndex + 1)
}
