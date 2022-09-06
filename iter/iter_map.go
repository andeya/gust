package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]     = (*mapIterator[any, any])(nil)
	_ iRealNext[any]    = (*mapIterator[any, any])(nil)
	_ iRealTryFold[any] = (*mapIterator[any, any])(nil)
	_ iRealFold[any]    = (*mapIterator[any, any])(nil)
	_ iRealSizeHint     = (*mapIterator[any, any])(nil)
)

func newMapIterator[T any, B any](iter Iterator[T], f func(T) B) Iterator[B] {
	p := &mapIterator[T, B]{iter: iter, f: f}
	p.setFacade(p)
	return p
}

type mapIterator[T any, B any] struct {
	iterBackground[B]
	iter Iterator[T]
	f    func(T) B
}

func (s mapIterator[T, B]) realSizeHint() (uint, gust.Option[uint]) {
	return s.iter.SizeHint()
}

func (s mapIterator[T, B]) realFold(init any, g func(any, B) any) any {
	return Fold[T, any](s.iter, init, func(acc any, elt T) any { return g(acc, s.f(elt)) })
}

func (s mapIterator[T, B]) realTryFold(init any, g func(any, B) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return TryFold[T, any](s.iter, init, func(acc any, elt T) gust.AnyCtrlFlow { return g(acc, s.f(elt)) })
}

func (s mapIterator[T, B]) realNext() gust.Option[B] {
	r := s.iter.Next()
	if r.IsSome() {
		return gust.Some(s.f(r.Unwrap()))
	}
	return gust.None[B]()
}
