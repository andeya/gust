// nolint:unused
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
	deIterBackground[B]
	iter Iterator[T]
	f    func(T) B
}

func (*mapIterator[T, B]) realNextBack() gust.Option[B] {
	panic("unreachable")
}

func (s *mapIterator[T, B]) realSizeHint() (uint, gust.Option[uint]) {
	return s.iter.SizeHint()
}

func (s *mapIterator[T, B]) realFold(init any, g func(any, B) any) any {
	return Fold[T, any](s.iter, init, func(acc any, elt T) any { return g(acc, s.f(elt)) })
}

func (s *mapIterator[T, B]) realTryFold(init any, g func(any, B) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return SigTryFold[T, any](s.iter, init, func(acc any, elt T) gust.AnyCtrlFlow { return g(acc, s.f(elt)) })
}

func (s *mapIterator[T, B]) realNext() gust.Option[B] {
	r := s.iter.Next()
	if r.IsSome() {
		return gust.Some(s.f(r.Unwrap()))
	}
	return gust.None[B]()
}

var (
	_ DeIterator[any]    = (*deMapIterator[any, any])(nil)
	_ iRealRemaining     = (*deMapIterator[any, any])(nil)
	_ iRealNextBack[any] = (*deMapIterator[any, any])(nil)
	_ iRealTryRfold[any] = (*deMapIterator[any, any])(nil)
	_ iRealRfold[any]    = (*deMapIterator[any, any])(nil)
)

func newDeMapIterator[T any, B any](iter DeIterator[T], f func(T) B) DeIterator[B] {
	p := &deMapIterator[T, B]{}
	p.iter = iter
	p.f = f
	p.setFacade(p)
	return p
}

type deMapIterator[T any, B any] struct {
	mapIterator[T, B]
}

func (d *deMapIterator[T, B]) realRemaining() uint {
	return d.iter.(DeIterator[T]).Remaining()
}

func (d *deMapIterator[T, B]) realNextBack() gust.Option[B] {
	r := d.iter.(DeIterator[T]).NextBack()
	if r.IsSome() {
		return gust.Some(d.f(r.Unwrap()))
	}
	return gust.None[B]()
}

func (d *deMapIterator[T, B]) realRfold(init any, g func(any, B) any) any {
	return Rfold[T, any](d.iter.(DeIterator[T]), init, func(acc any, elt T) any { return g(acc, d.f(elt)) })
}

func (d *deMapIterator[T, B]) realTryRfold(init any, g func(any, B) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return TryRfold[T, any](d.iter.(DeIterator[T]), init, func(acc any, elt T) gust.AnyCtrlFlow { return g(acc, d.f(elt)) })
}
