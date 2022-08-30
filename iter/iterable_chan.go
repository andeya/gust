package iter

import (
	"github.com/andeya/gust"
)

var (
	_ gust.Iterable[any] = (*DataChan[any])(nil)
)

type DataChan[T any] struct {
	c <-chan T
}

func NewDataChan[T any](c <-chan T) DataChan[T] {
	return DataChan[T]{c: c}
}

func (c DataChan[T]) ToIterator() Iterator[T] {
	return FromIterable[T](c)
}

func (c DataChan[T]) Next() gust.Option[T] {
	var x, ok = <-c.c
	if ok {
		return gust.Some(x)
	}
	return gust.None[T]()
}
