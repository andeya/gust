package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]  = (*FilterMap[any])(nil)
	_ iRealNext[any] = (*FilterMap[any])(nil)
	// _ iRealSizeHint  = (*FilterMap[any])(nil)
	// _ iRealCount     = (*FilterMap[any])(nil)
)

func newFilterMap[T any](inner Iterator[T], filterMap func(T) gust.Option[T]) *FilterMap[T] {
	iter := &FilterMap[T]{inner: inner, filterMap: filterMap}
	iter.setFacade(iter)
	return iter
}

type FilterMap[T any] struct {
	iterTrait[T]
	inner     Iterator[T]
	filterMap func(T) gust.Option[T]
}

func (f FilterMap[T]) realNext() gust.Option[T] {
	// TODO implement me
	panic("implement me")
}
