package iter

import (
	"math"

	"github.com/andeya/gust"
)

var (
	_ Iterator[any]  = (*Chain[any])(nil)
	_ iRealNext[any] = (*Chain[any])(nil)
	_ iRealSizeHint  = (*Chain[any])(nil)
	_ iRealCount     = (*Chain[any])(nil)
)

func newChain[T any](inner Iterator[T], other Iterator[T]) *Chain[T] {
	iter := &Chain[T]{inner: inner, other: other}
	iter.setFacade(iter)
	return iter
}

type Chain[T any] struct {
	iterTrait[T]
	inner Iterator[T]
	other Iterator[T]
}

func (s *Chain[T]) realNext() gust.Option[T] {
	if s.inner != nil {
		item := s.inner.Next()
		if item.IsSome() {
			return item
		}
		s.inner = nil
	}
	if s.other != nil {
		item := s.other.Next()
		if item.IsSome() {
			return item
		}
		s.other = nil
	}
	return gust.None[T]()
}

func saturatingAdd(a, b uint64) uint64 {
	if a < math.MaxUint64-b {
		return a + b
	}
	return math.MaxUint64
}

func checkedAdd(a, b uint64) gust.Option[uint64] {
	if a <= math.MaxUint64-b {
		return gust.Some(a + b)
	}
	return gust.None[uint64]()
}

func (s *Chain[T]) realSizeHint() (uint64, gust.Option[uint64]) {
	if s.inner != nil && s.other != nil {
		var aLower, aUpper = s.inner.SizeHint()
		var bLower, bUpper = s.other.SizeHint()
		var lower = saturatingAdd(aLower, bLower)
		var upper gust.Option[uint64]
		if aUpper.IsSome() && bUpper.IsSome() {
			upper = checkedAdd(aUpper.Unwrap(), bUpper.Unwrap())
		}
		return lower, upper
	}
	if s.inner != nil && s.other == nil {
		return s.inner.SizeHint()
	}
	if s.inner == nil && s.other != nil {
		return s.other.SizeHint()
	}
	return 0, gust.Some[uint64](0)
}

func (s *Chain[T]) realCount() uint64 {
	var aCount uint64
	if s.inner != nil {
		aCount = s.inner.Count()
	}
	var bCount uint64
	if s.other != nil {
		bCount = s.other.Count()
	}
	return aCount + bCount
}
