package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]     = (*FilterMapIterator[any])(nil)
	_ iRealNext[any]    = (*FilterMapIterator[any])(nil)
	_ iRealSizeHint     = (*FilterMapIterator[any])(nil)
	_ iRealFold[any]    = (*FilterMapIterator[any])(nil)
	_ iRealTryFold[any] = (*FilterMapIterator[any])(nil)
)

func newFilterMapIterator[T any](iter Iterator[T], filterMap func(T) gust.Option[T]) *FilterMapIterator[T] {
	p := &FilterMapIterator[T]{iter: iter, f: filterMap}
	p.setFacade(p)
	return p
}

type FilterMapIterator[T any] struct {
	iterTrait[T]
	iter Iterator[T]
	f    func(T) gust.Option[T]
}

func (f FilterMapIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	var _, upper = f.iter.SizeHint()
	return 0, upper // can't know a lower bound, due to the f
}

func (f FilterMapIterator[T]) realNext() gust.Option[T] {
	return FindMap(f.iter, f.f)
}

func (f FilterMapIterator[T]) realFold(init any, fold func(any, T) any) any {
	return f.iter.Fold(init, func(acc any, item T) any {
		r := f.f(item)
		if r.IsSome() {
			return fold(acc, item)
		}
		return acc
	})
}

func (f FilterMapIterator[T]) realTryFold(init any, fold func(any, T) gust.Result[any]) gust.Result[any] {
	return f.iter.TryFold(init, func(acc any, item T) gust.Result[any] {
		r := f.f(item)
		if r.IsSome() {
			return fold(acc, item)
		}
		return gust.Ok(acc)
	})
}
