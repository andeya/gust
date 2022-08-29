package iter

import (
	"github.com/andeya/gust"
	"github.com/andeya/gust/digit"
)

var (
	_ gust.DataForIter[uint]            = (*DataRange[uint])(nil)
	_ gust.SizeHintForIter              = (*DataRange[uint])(nil)
	_ gust.CountForIter                 = (*DataRange[uint])(nil)
	_ gust.DataForDoubleEndedIter[uint] = (*DataRange[uint])(nil)
)

type DataRange[T digit.Integer] struct {
	nextValue     T
	backNextValue T
	ended         bool
}

func NewDataRange[T digit.Integer](start T, end T, rightClosed ...bool) *DataRange[T] {
	max := end
	if len(rightClosed) == 0 || !rightClosed[0] {
		if end <= start {
			return &DataRange[T]{
				nextValue:     start,
				ended:         true,
				backNextValue: max,
			}
		}
		max = end - 1
	}
	ended := false
	if max < start {
		ended = true
	}
	return &DataRange[T]{
		nextValue:     start,
		ended:         ended,
		backNextValue: max,
	}
}

func (r *DataRange[T]) ToIterator() Iterator[T] {
	return newIter[T](r)
}

func (r *DataRange[T]) ToDoubleEndedIterator() DoubleEndedIterator[T] {
	return newDoubleEndedIter[T](r)
}

func (r *DataRange[T]) NextForIter() gust.Option[T] {
	if r.ended {
		return gust.None[T]()
	}
	value := r.nextValue
	if r.nextValue == r.backNextValue {
		r.ended = true
	} else {
		r.nextValue++
	}
	return gust.Some(value)
}

func (r *DataRange[T]) NextBackForIter() gust.Option[T] {
	if r.ended {
		return gust.None[T]()
	}
	value := r.backNextValue
	if r.backNextValue == r.nextValue {
		r.ended = true
	} else {
		r.backNextValue++
	}
	return gust.Some(value)
}

func (r *DataRange[T]) SizeHintForIter() (uint, gust.Option[uint]) {
	size := uint(r.backNextValue - r.nextValue + 1)
	return size, gust.Some(size)
}

func (r *DataRange[T]) CountForIter() uint {
	if !r.ended {
		return 0
	}
	r.ended = true
	return uint(r.backNextValue - r.nextValue + 1)
}

func (r *DataRange[T]) RemainingLenForIter() uint {
	if !r.ended {
		return 0
	}
	return uint(r.backNextValue - r.nextValue + 1)
}
