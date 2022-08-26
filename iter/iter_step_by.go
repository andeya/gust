package iter

import (
	"math"

	"github.com/andeya/gust"
)

var (
	_ Iterator[any]     = (*StepByIterator[any])(nil)
	_ iRealNext[any]    = (*StepByIterator[any])(nil)
	_ iRealSizeHint     = (*StepByIterator[any])(nil)
	_ iRealNth[any]     = (*StepByIterator[any])(nil)
	_ iRealTryFold[any] = (*StepByIterator[any])(nil)
	_ iRealFold[any]    = (*StepByIterator[any])(nil)
)

func newStepByIterator[T any](iter Iterator[T], step uint) *StepByIterator[T] {
	if step == 0 {
		panic("step must be non-zero")
	}
	p := &StepByIterator[T]{iter: iter, step: step - 1, firstTake: true}
	p.setFacade(p)
	return p
}

type StepByIterator[T any] struct {
	iterTrait[T]
	iter      Iterator[T]
	step      uint
	firstTake bool
}

func (s *StepByIterator[T]) realFold(acc any, f func(any, T) any) any {
	if s.firstTake {
		s.firstTake = false
		r := s.iter.Next()
		if r.IsNone() {
			return acc
		}
		acc = f(acc, r.Unwrap())
	}
	r := s.iter.Nth(s.step)
	if r.IsSome() {
		return f(acc, r.Unwrap())
	}
	return acc
}

func (s *StepByIterator[T]) realTryFold(acc any, f func(any, T) gust.Result[any]) gust.Result[any] {
	if s.firstTake {
		s.firstTake = false
		r := s.iter.Next()
		if r.IsNone() {
			return gust.Ok(acc)
		}
		v := f(acc, r.Unwrap())
		if v.IsErr() {
			return v
		}
		acc = v.Unwrap()
	}
	r := s.iter.Nth(s.step)
	if r.IsSome() {
		return f(acc, r.Unwrap())
	}
	return gust.Ok(acc)
}

func (s *StepByIterator[T]) realNth(n uint) gust.Option[T] {
	if s.firstTake {
		s.firstTake = false
		var first = s.iter.Next()
		if n == 0 {
			return first
		}
		n -= 1
	}
	// n and s.step are indices, we need to add 1 to get the amount of elements
	// When calling `.Nth`, we need to subtract 1 again to convert back to an index
	// step + 1 can't overflow because `.step_by` sets `s.step` to `step - 1`
	var step = s.step + 1
	// n + 1 could overflow
	// thus, if n is math.MaxUint, instead of adding one, we call .Nth(step)
	if n == math.MaxUint {
		s.iter.Nth(step - 1)
	} else {
		n += 1
	}

	// overflow handling
	for {
		var mul = uintCheckedMul(n, step)
		if mul.IsSome() {
			return s.iter.Nth(mul.Unwrap() - 1)
		}
		var divN = math.MaxUint / n
		var divStep = math.MaxUint / step
		var nthN = divN * n
		var nthStep = divStep * step
		var nth uint
		if nthN > nthStep {
			step -= divN
			nth = nthN
		} else {
			n -= divStep
			nth = nthStep
		}
		s.iter.Nth(nth - 1)
	}
}

func (s *StepByIterator[T]) realSizeHint() (uint64, gust.Option[uint64]) {
	var firstSize = func(step uint64) func(uint64) uint64 {
		return func(n uint64) uint64 {
			if n == 0 {
				return 0
			}
			return 1 + (n-1)/(step+1)
		}
	}

	var otherSize = func(step uint64) func(uint64) uint64 {
		return func(n uint64) uint64 { return n / (step + 1) }
	}

	var low, high = s.iter.SizeHint()

	if s.firstTake {
		var f = firstSize(uint64(s.step))
		return f(low), high.Map(f)
	}
	var f = otherSize(uint64(s.step))
	return f(low), high.Map(f)
}

func (s *StepByIterator[T]) realNext() gust.Option[T] {
	if s.firstTake {
		s.firstTake = false
		return s.iter.Next()
	}
	return s.iter.Nth(s.step)
}
