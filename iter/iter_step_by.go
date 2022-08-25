package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]  = (*StepByIterator[any])(nil)
	_ iRealNext[any] = (*StepByIterator[any])(nil)
	// _ iRealSizeHint  = (*StepByIterator[any])(nil)
	// _ iRealCount     = (*StepByIterator[any])(nil)
)

func newStepByIterator[T any](inner Iterator[T], step uint) *StepByIterator[T] {
	if step == 0 {
		panic("step must be non-zero")
	}
	iter := &StepByIterator[T]{inner: inner, step: step - 1, firstTake: true}
	iter.setFacade(iter)
	return iter
}

type StepByIterator[T any] struct {
	iterTrait[T]
	inner     Iterator[T]
	step      uint
	firstTake bool
}

func (s StepByIterator[T]) realNext() gust.Option[T] {
	// TODO implement me
	panic("implement me")
}
