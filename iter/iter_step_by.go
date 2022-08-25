package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]  = (*StepBy[any])(nil)
	_ iRealNext[any] = (*StepBy[any])(nil)
	// _ iRealSizeHint  = (*StepBy[any])(nil)
	// _ iRealCount     = (*StepBy[any])(nil)
)

func newStepBy[T any](inner Iterator[T], step uint) *StepBy[T] {
	if step == 0 {
		panic("step must be non-zero")
	}
	iter := &StepBy[T]{inner: inner, step: step - 1, firstTake: true}
	iter.setFacade(iter)
	return iter
}

type StepBy[T any] struct {
	iterTrait[T]
	inner     Iterator[T]
	step      uint
	firstTake bool
}

func (s StepBy[T]) realNext() gust.Option[T] {
	// TODO implement me
	panic("implement me")
}
