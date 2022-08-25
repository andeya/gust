package iter

import (
	"github.com/andeya/gust"
)

var (
	_ NextForIter[any] = (*ChanNext[any])(nil)
)

type ChanNext[T any] struct {
	c <-chan T
}

func NewChanNext[T any](c <-chan T) ChanNext[T] {
	return ChanNext[T]{c: c}
}

func (c ChanNext[T]) ToIter() *Iter[T] {
	return newIter[T](c)
}

func (c ChanNext[T]) NextForIter() gust.Option[T] {
	var x, ok = <-c.c
	if ok {
		return gust.Some(x)
	}
	return gust.None[T]()
}
