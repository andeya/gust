package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]     = (*MapIterator[any, any])(nil)
	_ iRealNext[any]    = (*MapIterator[any, any])(nil)
	_ iRealTryFold[any] = (*MapIterator[any, any])(nil)
	_ iRealFold[any]    = (*MapIterator[any, any])(nil)
	_ iRealSizeHint     = (*MapIterator[any, any])(nil)
)

func newMapIterator[T any, B any](iter Iterator[T], f func(T) B) *MapIterator[T, B] {
	p := &MapIterator[T, B]{iter: iter, f: f}
	p.facade = p
	return p
}

type MapIterator[T any, B any] struct {
	iterTrait[B]
	iter Iterator[T]
	f    func(T) B
}

func (s MapIterator[T, B]) realSizeHint() (uint64, gust.Option[uint64]) {
	return s.iter.SizeHint()
}

func (s MapIterator[T, B]) realFold(init any, g func(any, B) any) any {
	return Fold[T, any](s.iter, init, func(acc any, elt T) any { return g(acc, s.f(elt)) })
}

func (s MapIterator[T, B]) realTryFold(init any, g func(any, B) gust.Result[any]) gust.Result[any] {
	return TryFold[T, any](s.iter, init, func(acc any, elt T) gust.Result[any] { return g(acc, s.f(elt)) })
}

func (s MapIterator[T, B]) realNext() gust.Option[B] {
	r := s.iter.Next()
	if r.IsSome() {
		return gust.Some(s.f(r.Unwrap()))
	}
	return gust.None[B]()
}
