package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]     = (*FilterIterator[any])(nil)
	_ iRealFold[any]    = (*FilterIterator[any])(nil)
	_ iRealTryFold[any] = (*FilterIterator[any])(nil)
	_ iRealNext[any]    = (*FilterIterator[any])(nil)
	_ iRealSizeHint     = (*FilterIterator[any])(nil)
	_ iRealCount        = (*FilterIterator[any])(nil)
)

func newFilterIterator[T any](iter Iterator[T], predicate func(T) bool) *FilterIterator[T] {
	p := &FilterIterator[T]{iter: iter, predicate: predicate}
	p.setFacade(p)
	return p
}

type FilterIterator[T any] struct {
	iterTrait[T]
	iter      Iterator[T]
	predicate func(T) bool
}

func (f FilterIterator[T]) realFold(init any, fold func(any, T) any) any {
	return f.iter.Fold(init, func(acc any, item T) any {
		if f.predicate(item) {
			return fold(acc, item)
		}
		return acc
	})
}

func (f FilterIterator[T]) realTryFold(init any, fold func(any, T) gust.Result[any]) gust.Result[any] {
	return f.iter.TryFold(init, func(acc any, item T) gust.Result[any] {
		if f.predicate(item) {
			return fold(acc, item)
		}
		return gust.Ok(acc)
	})
}

func (f FilterIterator[T]) realCount() uint {
	return Map[T, uint](f.iter, func(x T) uint {
		if f.predicate(x) {
			return 1
		}
		return 0
	}).Fold(uint(0), func(count any, x uint) any {
		return count.(uint) + x
	}).(uint)
}

func (f FilterIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	var _, upper = f.iter.SizeHint()
	return 0, upper // can't know a lower bound, due to the predicate
}

func (f FilterIterator[T]) realNext() gust.Option[T] {
	return f.iter.Find(f.predicate)
}
