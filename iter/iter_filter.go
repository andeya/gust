package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]  = (*FilterIterator[any])(nil)
	_ iRealNext[any] = (*FilterIterator[any])(nil)
	// _ iRealSizeHint  = (*FilterIterator[any])(nil)
	// _ iRealCount     = (*FilterIterator[any])(nil)
)

func newFilterIterator[T any](inner Iterator[T], filter func(T) bool) *FilterIterator[T] {
	iter := &FilterIterator[T]{inner: inner, filter: filter}
	iter.setFacade(iter)
	return iter
}

type FilterIterator[T any] struct {
	iterTrait[T]
	inner  Iterator[T]
	filter func(T) bool
}

func (f FilterIterator[T]) realNext() gust.Option[T] {
	// TODO implement me
	panic("implement me")
}
