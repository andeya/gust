package iter

import (
	"github.com/andeya/gust"
)

func newIter[T any](data gust.DataForIter[T]) Iterator[T] {
	iter := &implIter[T]{data: data}
	iter.setFacade(iter)
	return iter
}

var (
	_ Iterator[any]  = (*implIter[any])(nil)
	_ iRealNext[any] = (*implIter[any])(nil)
	_ iRealCount     = (*implIter[any])(nil)
	_ iRealSizeHint  = (*implIter[any])(nil)
)

type implIter[T any] struct {
	iterTrait[T]
	data gust.DataForIter[T]
}

func (iter *implIter[T]) realNext() gust.Option[T] {
	return iter.data.NextForIter()
}

func (iter *implIter[T]) realCount() uint {
	if c, ok := iter.data.(gust.CountForIter); ok {
		return c.CountForIter()
	}
	var a uint
	for iter.data.NextForIter().IsSome() {
		a++
	}
	return a
}

func (iter *implIter[T]) realSizeHint() (uint, gust.Option[uint]) {
	if cover, ok := iter.data.(gust.SizeHintForIter); ok {
		return cover.SizeHintForIter()
	}
	return 0, gust.None[uint]()
}

func newDoubleEndedIter[T any](data gust.DataForDoubleEndedIter[T]) DoubleEndedIterator[T] {
	iter := &implDoubleEndedIter[T]{data: data}
	iter.setFacade(iter)
	return iter
}

var (
	_ iIterator[any]     = (*implDoubleEndedIter[any])(nil)
	_ iRealNext[any]     = (*implDoubleEndedIter[any])(nil)
	_ iRealNextBack[any] = (*implDoubleEndedIter[any])(nil)
	_ iRealRemainingLen  = (*implDoubleEndedIter[any])(nil)
	_ iRealCount         = (*implDoubleEndedIter[any])(nil)
	_ iRealSizeHint      = (*implDoubleEndedIter[any])(nil)
)

type implDoubleEndedIter[T any] struct {
	doubleEndedIterTrait[T]
	data gust.DataForDoubleEndedIter[T]
}

func (iter *implDoubleEndedIter[T]) realNext() gust.Option[T] {
	return iter.data.NextForIter()
}

func (iter *implDoubleEndedIter[T]) realNextBack() gust.Option[T] {
	return iter.data.NextBackForIter()
}

func (iter *implDoubleEndedIter[T]) realRemainingLen() uint {
	return iter.data.RemainingLenForIter()
}

func (iter *implDoubleEndedIter[T]) realCount() uint {
	if c, ok := iter.data.(gust.CountForIter); ok {
		return c.CountForIter()
	}
	var a uint
	for iter.data.NextForIter().IsSome() {
		a++
	}
	return a
}

func (iter *implDoubleEndedIter[T]) realSizeHint() (uint, gust.Option[uint]) {
	if cover, ok := iter.data.(gust.SizeHintForIter); ok {
		return cover.SizeHintForIter()
	}
	return 0, gust.None[uint]()
}
