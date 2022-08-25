package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]  = (*FilterMapIterator[any])(nil)
	_ iRealNext[any] = (*FilterMapIterator[any])(nil)
	// _ iRealSizeHint  = (*FilterMapIterator[any])(nil)
	// _ iRealCount     = (*FilterMapIterator[any])(nil)
)

func newFilterMapIterator[T any](inner Iterator[T], filterMap func(T) gust.Option[T]) *FilterMapIterator[T] {
	iter := &FilterMapIterator[T]{inner: inner, filterMap: filterMap}
	iter.setFacade(iter)
	return iter
}

type FilterMapIterator[T any] struct {
	iterTrait[T]
	inner     Iterator[T]
	filterMap func(T) gust.Option[T]
}

func (f FilterMapIterator[T]) realNext() gust.Option[T] {
	// TODO implement me
	panic("implement me")
}
