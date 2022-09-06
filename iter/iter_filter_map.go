package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]     = (*FilterMapIterator[any, any])(nil)
	_ iRealNext[any]    = (*FilterMapIterator[any, any])(nil)
	_ iRealSizeHint     = (*FilterMapIterator[any, any])(nil)
	_ iRealFold[any]    = (*FilterMapIterator[any, any])(nil)
	_ iRealTryFold[any] = (*FilterMapIterator[any, any])(nil)
)

func newFilterMapIterator[T any, B any](iter Iterator[T], filterMap func(T) gust.Option[B]) *FilterMapIterator[T, B] {
	p := &FilterMapIterator[T, B]{iter: iter, f: filterMap}
	p.facade = p
	return p
}

type FilterMapIterator[T any, B any] struct {
	iterTrait[B]
	iter Iterator[T]
	f    func(T) gust.Option[B]
}

func (f FilterMapIterator[T, B]) realSizeHint() (uint, gust.Option[uint]) {
	var _, upper = f.iter.SizeHint()
	return 0, upper // can't know a lower bound, due to the f
}

func (f FilterMapIterator[T, B]) realNext() gust.Option[B] {
	return FindMap(f.iter, f.f)
}

func (f FilterMapIterator[T, B]) realFold(init any, fold func(any, B) any) any {
	return f.iter.Fold(init, func(acc any, item T) any {
		r := f.f(item)
		if r.IsSome() {
			return fold(acc, r.Unwrap())
		}
		return acc
	})
}

func (f FilterMapIterator[T, B]) realTryFold(init any, fold func(any, B) gust.Result[any]) gust.Result[any] {
	return f.iter.TryFold(init, func(acc any, item T) gust.Result[any] {
		r := f.f(item)
		if r.IsSome() {
			return fold(acc, r.Unwrap())
		}
		return gust.Ok(acc)
	})
}
