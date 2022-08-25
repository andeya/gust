package iter

import (
	"github.com/andeya/gust"
	"github.com/andeya/gust/digit"
)

var (
	_ NextForIter[uint64] = (*RangeNext[uint64])(nil)
	_ SizeHintForIter     = (*RangeNext[uint64])(nil)
	_ CountForIter        = (*RangeNext[uint64])(nil)
)

type RangeNext[T digit.Integer] struct {
	nextValue T
	max       T
	ended     bool
}

func NewRangeNext[T digit.Integer](start T, end T, rightClosed ...bool) *RangeNext[T] {
	max := end
	if len(rightClosed) == 0 || !rightClosed[0] {
		if end <= start {
			return &RangeNext[T]{nextValue: start, max: max, ended: true}
		}
		max = end - 1
	}
	ended := false
	if max < start {
		ended = true
	}
	return &RangeNext[T]{
		max:       max,
		nextValue: start,
		ended:     ended,
	}
}

func (r *RangeNext[T]) ToIter() *Iter[T] {
	return newIter[T](r)
}

func (r *RangeNext[T]) NextForIter() gust.Option[T] {
	if r.ended {
		return gust.None[T]()
	}
	value := r.nextValue
	if r.nextValue == r.max {
		r.ended = true
	} else {
		r.nextValue++
	}
	return gust.Some(value)
}

func (r *RangeNext[T]) SizeHintForIter() (uint64, gust.Option[uint64]) {
	size := uint64(r.max - r.nextValue + 1)
	return size, gust.Some(size)
}

func (r *RangeNext[T]) CountForIter() uint64 {
	if !r.ended {
		return 0
	}
	r.ended = true
	return uint64(r.max - r.nextValue + 1)
}
