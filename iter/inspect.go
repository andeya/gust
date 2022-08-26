package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]     = (*InspectIterator[any])(nil)
	_ iRealNext[any]    = (*InspectIterator[any])(nil)
	_ iRealSizeHint     = (*InspectIterator[any])(nil)
	_ iRealTryFold[any] = (*InspectIterator[any])(nil)
	_ iRealFold[any]    = (*InspectIterator[any])(nil)
)

func newInspectIterator[T any](iter Iterator[T], f func(T)) *InspectIterator[T] {
	p := &InspectIterator[T]{iter: iter, f: f}
	p.facade = p
	return p
}

type InspectIterator[T any] struct {
	iterTrait[T]
	iter Iterator[T]
	f    func(T)
}

func (s InspectIterator[T]) doInspect(elt gust.Option[T]) gust.Option[T] {
	if elt.IsSome() {
		s.f(elt.Unwrap())
	}
	return elt
}

func (s InspectIterator[T]) realNext() gust.Option[T] {
	return s.doInspect(s.iter.Next())
}

func (s InspectIterator[T]) realSizeHint() (uint64, gust.Option[uint64]) {
	return s.iter.SizeHint()
}

func (s InspectIterator[T]) realTryFold(init any, g func(any, T) gust.Result[any]) gust.Result[any] {
	return TryFold[T, any](s.iter, init, func(acc any, elt T) gust.Result[any] {
		s.f(elt)
		return g(acc, elt)
	})
}

func (s InspectIterator[T]) realFold(init any, g func(any, T) any) any {
	return Fold[T, any](s.iter, init, func(acc any, elt T) any {
		s.f(elt)
		return g(acc, elt)
	})
}
