package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]     = (*inspectIterator[any])(nil)
	_ iRealNext[any]    = (*inspectIterator[any])(nil)
	_ iRealSizeHint     = (*inspectIterator[any])(nil)
	_ iRealTryFold[any] = (*inspectIterator[any])(nil)
	_ iRealFold[any]    = (*inspectIterator[any])(nil)
)

func newInspectIterator[T any](iter Iterator[T], f func(T)) Iterator[T] {
	p := &inspectIterator[T]{iter: iter, f: f}
	p.setFacade(p)
	return p
}

type inspectIterator[T any] struct {
	iterBackground[T]
	iter Iterator[T]
	f    func(T)
}

func (s inspectIterator[T]) doInspect(elt gust.Option[T]) gust.Option[T] {
	if elt.IsSome() {
		s.f(elt.Unwrap())
	}
	return elt
}

func (s inspectIterator[T]) realNext() gust.Option[T] {
	return s.doInspect(s.iter.Next())
}

func (s inspectIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	return s.iter.SizeHint()
}

func (s inspectIterator[T]) realTryFold(init any, g func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return TryFold[T, any](s.iter, init, func(acc any, elt T) gust.AnyCtrlFlow {
		s.f(elt)
		return g(acc, elt)
	})
}

func (s inspectIterator[T]) realFold(init any, g func(any, T) any) any {
	return Fold[T, any](s.iter, init, func(acc any, elt T) any {
		s.f(elt)
		return g(acc, elt)
	})
}
