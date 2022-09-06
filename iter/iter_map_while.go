package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]     = (*mapWhileIterator[any, any])(nil)
	_ iRealNext[any]    = (*mapWhileIterator[any, any])(nil)
	_ iRealTryFold[any] = (*mapWhileIterator[any, any])(nil)
	_ iRealFold[any]    = (*mapWhileIterator[any, any])(nil)
	_ iRealSizeHint     = (*mapWhileIterator[any, any])(nil)
)

func newMapWhileIterator[T any, B any](iter Iterator[T], f func(T) gust.Option[B]) Iterator[B] {
	p := &mapWhileIterator[T, B]{iter: iter, f: f}
	p.setFacade(p)
	return p
}

type mapWhileIterator[T any, B any] struct {
	iterBackground[B]
	iter Iterator[T]
	f    func(T) gust.Option[B]
}

func (s mapWhileIterator[T, B]) realNext() gust.Option[B] {
	r := s.iter.Next()
	if r.IsSome() {
		return s.f(r.Unwrap())
	}
	return gust.None[B]()
}

func (s mapWhileIterator[T, B]) realSizeHint() (uint, gust.Option[uint]) {
	var _, upper = s.iter.SizeHint()
	return 0, upper // can't know a lower bound, due to the predicate
}

func (s mapWhileIterator[T, B]) realTryFold(init any, fold func(any, B) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return TryFold[T, any](s.iter, init, func(acc any, x T) gust.AnyCtrlFlow {
		r := s.f(x)
		if r.IsSome() {
			return fold(acc, r.Unwrap())
		}
		return gust.AnyBreak(acc)
	})
}

func (s mapWhileIterator[T, B]) realFold(init any, fold func(any, B) any) any {
	return s.TryFold(init, func(acc any, x B) gust.AnyCtrlFlow {
		return gust.AnyContinue(fold(acc, x))
	}).UnwrapContinue()
}
