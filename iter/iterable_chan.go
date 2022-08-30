package iter

import (
	"github.com/andeya/gust"
)

var (
	_ gust.Iterable[any] = (*IterableChan[any])(nil)
)

type IterableChan[T any] struct {
	c <-chan T
}

func NewIterableChan[T any](c <-chan T) IterableChan[T] {
	return IterableChan[T]{c: c}
}

func (c IterableChan[T]) ToIterator() Iterator[T] {
	return FromIterable[T](c)
}

func (c IterableChan[T]) Next() gust.Option[T] {
	var x, ok = <-c.c
	if ok {
		return gust.Some(x)
	}
	return gust.None[T]()
}
