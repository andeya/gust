package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]  = (*Filter[any])(nil)
	_ iRealNext[any] = (*Filter[any])(nil)
	// _ iRealSizeHint  = (*Filter[any])(nil)
	// _ iRealCount     = (*Filter[any])(nil)
)

func newFilter[T any](inner Iterator[T], filter func(T) bool) *Filter[T] {
	iter := &Filter[T]{inner: inner, filter: filter}
	iter.setFacade(iter)
	return iter
}

type Filter[T any] struct {
	iterTrait[T]
	inner  Iterator[T]
	filter func(T) bool
}

func (f Filter[T]) realNext() gust.Option[T] {
	// TODO implement me
	panic("implement me")
}
