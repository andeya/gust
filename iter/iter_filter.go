package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]     = (*filterIterator[any])(nil)
	_ iRealFold[any]    = (*filterIterator[any])(nil)
	_ iRealTryFold[any] = (*filterIterator[any])(nil)
	_ iRealNext[any]    = (*filterIterator[any])(nil)
	_ iRealSizeHint     = (*filterIterator[any])(nil)
	_ iRealCount        = (*filterIterator[any])(nil)
)

func newFilterIterator[T any](iter Iterator[T], predicate func(T) bool) Iterator[T] {
	p := &filterIterator[T]{iter: iter, predicate: predicate}
	p.setFacade(p)
	return p
}

type filterIterator[T any] struct {
	iterBackground[T]
	iter      Iterator[T]
	predicate func(T) bool
}

func (f filterIterator[T]) realFold(init any, fold func(any, T) any) any {
	return f.iter.Fold(init, func(acc any, item T) any {
		if f.predicate(item) {
			return fold(acc, item)
		}
		return acc
	})
}

func (f filterIterator[T]) realTryFold(init any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return f.iter.TryFold(init, func(acc any, item T) gust.AnyCtrlFlow {
		if f.predicate(item) {
			return fold(acc, item)
		}
		return gust.AnyContinue(acc)
	})
}

func (f filterIterator[T]) realCount() uint {
	return Fold[uint, uint](Map[T, uint](f.iter, func(x T) uint {
		if f.predicate(x) {
			return 1
		}
		return 0
	}), uint(0), func(count uint, x uint) uint {
		return count + x
	})
}

func (f filterIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	var _, upper = f.iter.SizeHint()
	return 0, upper // can't know a lower bound, due to the predicate
}

func (f filterIterator[T]) realNext() gust.Option[T] {
	return f.iter.Find(f.predicate)
}
