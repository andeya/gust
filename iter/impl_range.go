package iter

import (
	"github.com/andeya/gust"
	"github.com/andeya/gust/digit"
)

var (
	_ Nextor[uint64] = new(RangeNext[uint64])
	_ SizeHint       = new(RangeNext[uint64])
	_ counter        = new(RangeNext[uint64])
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

func (r *RangeNext[T]) ToIter() *AnyIter[T] {
	return IterAny[T](r)
}

func (r *RangeNext[T]) Next() gust.Option[T] {
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

func (r *RangeNext[T]) SizeHint() (uint64, gust.Option[uint64]) {
	size := uint64(r.max - r.nextValue + 1)
	return size, gust.Some(size)
}

func (r *RangeNext[T]) count() uint64 {
	if !r.ended {
		return 0
	}
	r.ended = true
	return uint64(r.max - r.nextValue + 1)
}
