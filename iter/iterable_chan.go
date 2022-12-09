package iter

import (
	"github.com/andeya/gust"
)

var (
	_ gust.Iterable[any]    = (*IterableChan[any])(nil)
	_ gust.IterableSizeHint = (*IterableChan[any])(nil)
	_ gust.IterableCount    = (*IterableChan[any])(nil)
)

func NewIterableChan[T any](c chan T) IterableChan[T] {
	return IterableChan[T]{c: c}
}

type IterableChan[T any] struct {
	c chan T
}

func (c IterableChan[T]) ToIterator() Iterator[T] {
	return FromIterable[T](c)
}

func (c IterableChan[T]) Next() gust.Option[T] {
	if cap(c.c) > 0 && len(c.c) == 0 {
		return gust.None[T]()
	}
	var x, ok = <-c.c
	if ok {
		return gust.Some(x)
	}
	return gust.None[T]()
}

func (c IterableChan[T]) SizeHint() (uint, gust.Option[uint]) {
	return uint(len(c.c)), gust.Some(uint(cap(c.c)))
}

func (c IterableChan[T]) Count() uint {
	defer func() { _ = recover() }()
	var count = uint(len(c.c))
	close(c.c)
	return count
}
