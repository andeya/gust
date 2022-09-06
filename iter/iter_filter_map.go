package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]     = (*filterMapIterator[any, any])(nil)
	_ iRealNext[any]    = (*filterMapIterator[any, any])(nil)
	_ iRealSizeHint     = (*filterMapIterator[any, any])(nil)
	_ iRealFold[any]    = (*filterMapIterator[any, any])(nil)
	_ iRealTryFold[any] = (*filterMapIterator[any, any])(nil)
)

func newFilterMapIterator[T any, B any](iter Iterator[T], filterMap func(T) gust.Option[B]) Iterator[B] {
	p := &filterMapIterator[T, B]{iter: iter, f: filterMap}
	p.setFacade(p)
	return p
}

type filterMapIterator[T any, B any] struct {
	iterBackground[B]
	iter Iterator[T]
	f    func(T) gust.Option[B]
}

func (f filterMapIterator[T, B]) realSizeHint() (uint, gust.Option[uint]) {
	var _, upper = f.iter.SizeHint()
	return 0, upper // can't know a lower bound, due to the f
}

func (f filterMapIterator[T, B]) realNext() gust.Option[B] {
	return FindMap(f.iter, f.f)
}

func (f filterMapIterator[T, B]) realFold(init any, fold func(any, B) any) any {
	return f.iter.Fold(init, func(acc any, item T) any {
		r := f.f(item)
		if r.IsSome() {
			return fold(acc, r.Unwrap())
		}
		return acc
	})
}

func (f filterMapIterator[T, B]) realTryFold(init any, fold func(any, B) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return f.iter.TryFold(init, func(acc any, item T) gust.AnyCtrlFlow {
		r := f.f(item)
		if r.IsSome() {
			return fold(acc, r.Unwrap())
		}
		return gust.AnyContinue(acc)
	})
}
