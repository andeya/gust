package iter

import (
	"github.com/andeya/gust"
	"github.com/andeya/gust/digit"
)

var (
	_ gust.Iterable[uint]   = (*IterableRange[uint])(nil)
	_ gust.IterableSizeHint = (*IterableRange[uint])(nil)
	_ gust.IterableCount    = (*IterableRange[uint])(nil)
	_ gust.DeIterable[uint] = (*IterableRange[uint])(nil)
)

type IterableRange[T digit.Integer] struct {
	nextValue     T
	backNextValue T
	ended         bool
}

func NewIterableRange[T digit.Integer](start T, end T, rightClosed ...bool) *IterableRange[T] {
	max := end
	if len(rightClosed) == 0 || !rightClosed[0] {
		if end <= start {
			return &IterableRange[T]{
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
	return &IterableRange[T]{
		nextValue:     start,
		ended:         ended,
		backNextValue: max,
	}
}

func (r *IterableRange[T]) ToDeIterator() DeIterator[T] {
	return FromDeIterable[T](r)
}

func (r *IterableRange[T]) Next() gust.Option[T] {
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

func (r *IterableRange[T]) NextBack() gust.Option[T] {
	if r.ended {
		return gust.None[T]()
	}
	value := r.backNextValue
	if r.backNextValue == r.nextValue {
		r.ended = true
	} else {
		r.backNextValue--
	}
	return gust.Some(value)
}

func (r *IterableRange[T]) SizeHint() (uint, gust.Option[uint]) {
	size := uint(r.backNextValue - r.nextValue + 1)
	return size, gust.Some(size)
}

func (r *IterableRange[T]) Count() uint {
	if r.ended {
		return 0
	}
	r.ended = true
	return uint(r.backNextValue - r.nextValue + 1)
}

func (r *IterableRange[T]) Remaining() uint {
	if r.ended {
		return 0
	}
	return uint(r.backNextValue - r.nextValue + 1)
}
