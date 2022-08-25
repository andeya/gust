package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]  = (*Iter[any])(nil)
	_ iRealNext[any] = (*Iter[any])(nil)
	_ iRealCount     = (*Iter[any])(nil)
	_ iRealSizeHint  = (*Iter[any])(nil)
)

func newIter[T any](next NextForIter[T]) *Iter[T] {
	iter := &Iter[T]{next: next}
	iter.setFacade(iter)
	return iter
}

type Iter[T any] struct {
	iterTrait[T]
	next NextForIter[T]
}

func (iter *Iter[T]) realNext() gust.Option[T] {
	return iter.next.NextForIter()
}

func (iter *Iter[T]) realCount() uint64 {
	if c, ok := iter.next.(CountForIter); ok {
		return c.CountForIter()
	}
	var a uint64
	for iter.next.NextForIter().IsSome() {
		a++
	}
	return a
}

func (iter *Iter[T]) realSizeHint() (uint64, gust.Option[uint64]) {
	if cover, ok := iter.next.(SizeHintForIter); ok {
		return cover.SizeHintForIter()
	}
	return 0, gust.None[uint64]()
}
