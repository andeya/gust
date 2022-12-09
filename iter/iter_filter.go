// nolint:unused
package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]     = (*filterIterator[any])(nil)
	_ iRealNext[any]    = (*filterIterator[any])(nil)
	_ iRealSizeHint     = (*filterIterator[any])(nil)
	_ iRealCount        = (*filterIterator[any])(nil)
	_ iRealTryFold[any] = (*filterIterator[any])(nil)
	_ iRealFold[any]    = (*filterIterator[any])(nil)
)

func newFilterIterator[T any](iter Iterator[T], predicate func(T) bool) Iterator[T] {
	p := &filterIterator[T]{iter: iter, predicate: predicate}
	p.setFacade(p)
	return p
}

type filterIterator[T any] struct {
	deIterBackground[T]
	iter      Iterator[T]
	predicate func(T) bool
}

func (f *filterIterator[T]) realNextBack() gust.Option[T] {
	panic("unreachable")
}

func (f *filterIterator[T]) realNext() gust.Option[T] {
	return f.iter.Find(f.predicate)
}

func (f *filterIterator[T]) realCount() uint {
	return Fold[uint, uint](ToMap[T, uint](f.iter, func(x T) uint {
		if f.predicate(x) {
			return 1
		}
		return 0
	}), uint(0), func(count uint, x uint) uint {
		return count + x
	})
}

func (f *filterIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	var _, upper = f.iter.SizeHint()
	return 0, upper // can't know a lower bound, due to the predicate
}

func (f *filterIterator[T]) realTryFold(init any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return f.iter.TryFold(init, func(acc any, item T) gust.AnyCtrlFlow {
		if f.predicate(item) {
			return fold(acc, item)
		}
		return gust.AnyContinue(acc)
	})
}

func (f *filterIterator[T]) realFold(init any, fold func(any, T) any) any {
	return f.iter.Fold(init, func(acc any, item T) any {
		if f.predicate(item) {
			return fold(acc, item)
		}
		return acc
	})
}

var (
	_ DeIterator[any]    = (*deFilterIterator[any])(nil)
	_ iRealRemaining     = (*deFilterIterator[any])(nil)
	_ iRealNextBack[any] = (*deFilterIterator[any])(nil)
	_ iRealTryRfold[any] = (*deFilterIterator[any])(nil)
	_ iRealRfold[any]    = (*deFilterIterator[any])(nil)
)

func newDeFilterIterator[T any](iter DeIterator[T], predicate func(T) bool) DeIterator[T] {
	p := &deFilterIterator[T]{filterIterator: filterIterator[T]{iter: iter, predicate: predicate}}
	p.setFacade(p)
	return p
}

type deFilterIterator[T any] struct {
	filterIterator[T]
}

func (d *deFilterIterator[T]) realRemaining() uint {
	return d.iter.(DeIterator[T]).Remaining()
}

func (d *deFilterIterator[T]) realNextBack() gust.Option[T] {
	return d.iter.(DeIterator[T]).Rfind(d.predicate)
}

func (d *deFilterIterator[T]) realTryRfold(init any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return d.iter.(DeIterator[T]).TryRfold(init, func(acc any, item T) gust.AnyCtrlFlow {
		if d.predicate(item) {
			return fold(acc, item)
		}
		return gust.AnyContinue(acc)
	})
}

func (d *deFilterIterator[T]) realRfold(init any, fold func(any, T) any) any {
	return d.iter.(DeIterator[T]).Rfold(init, func(acc any, item T) any {
		if d.predicate(item) {
			return fold(acc, item)
		}
		return acc
	})
}
