package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]       = (*chainIterator[any])(nil)
	_ iRealLast[any]      = (*chainIterator[any])(nil)
	_ iRealFind[any]      = (*chainIterator[any])(nil)
	_ iRealNext[any]      = (*chainIterator[any])(nil)
	_ iRealSizeHint       = (*chainIterator[any])(nil)
	_ iRealCount          = (*chainIterator[any])(nil)
	_ iRealTryFold[any]   = (*chainIterator[any])(nil)
	_ iRealFold[any]      = (*chainIterator[any])(nil)
	_ iRealAdvanceBy[any] = (*chainIterator[any])(nil)
	_ iRealNth[any]       = (*chainIterator[any])(nil)
)

func newChainIterator[T any](inner Iterator[T], other Iterator[T]) Iterator[T] {
	iter := &chainIterator[T]{inner: inner, other: other}
	iter.setFacade(iter)
	return iter
}

type chainIterator[T any] struct {
	iterBackground[T]
	inner Iterator[T]
	other Iterator[T]
}

func (s *chainIterator[T]) realLast() gust.Option[T] {
	// Must exhaust a before b.
	var aLast gust.Option[T]
	var bLast gust.Option[T]
	if s.inner != nil {
		aLast = s.inner.Last()
	}
	if s.other != nil {
		bLast = s.other.Last()
	}
	if bLast.IsSome() {
		return bLast
	}
	return aLast
}

func (s *chainIterator[T]) realFind(predicate func(T) bool) gust.Option[T] {
	if s.inner != nil {
		item := s.inner.Find(predicate)
		if item.IsSome() {
			return item
		}
		s.inner = nil
	}
	if s.other != nil {
		return s.other.Find(predicate)
	}
	return gust.None[T]()
}

func (s *chainIterator[T]) realNth(n uint) gust.Option[T] {
	if s.inner != nil {
		r := s.inner.AdvanceBy(n)
		if r.IsErr() {
			n -= r.UnwrapErr()
		} else {
			item := s.inner.Next()
			if item.IsSome() {
				return item
			}
			n = 0
		}
		s.inner = nil
	}
	if s.other != nil {
		return s.other.Nth(n)
	}
	return gust.None[T]()
}

func (s *chainIterator[T]) realAdvanceBy(n uint) gust.Errable[uint] {
	var rem = n
	if s.inner != nil {
		r := s.inner.AdvanceBy(rem)
		if !r.IsErr() {
			return r
		}
		rem -= r.UnwrapErr()
		s.inner = nil
	}
	if s.other != nil {
		r := s.other.AdvanceBy(rem)
		if !r.IsErr() {
			return r
		}
		rem -= r.UnwrapErr()
		// we don't fuse the second iterator
	}
	if rem == 0 {
		return gust.NonErrable[uint]()
	}
	return gust.ToErrable(n - rem)
}

func (s *chainIterator[T]) realFold(acc any, f func(any, T) any) any {
	if s.inner != nil {
		acc = s.inner.Fold(acc, f)
	}
	if s.other != nil {
		acc = s.other.Fold(acc, f)
	}
	return acc
}

func (s *chainIterator[T]) realTryFold(acc any, f func(any, T) gust.Result[any]) gust.Result[any] {
	if s.inner != nil {
		r := s.inner.TryFold(acc, f)
		if r.IsErr() {
			return r
		}
		acc = r.Unwrap()
		s.inner = nil
	}
	if s.other != nil {
		r := s.other.TryFold(acc, f)
		if r.IsErr() {
			return r
		}
		acc = r.Unwrap()
		// we don't fuse the second iterator
	}
	return gust.Ok(acc)
}

func (s *chainIterator[T]) realNext() gust.Option[T] {
	if s.inner != nil {
		item := s.inner.Next()
		if item.IsSome() {
			return item
		}
		s.inner = nil
	}
	if s.other != nil {
		return s.other.Next()
	}
	return gust.None[T]()
}

func (s *chainIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	if s.inner != nil && s.other != nil {
		var aLower, aUpper = s.inner.SizeHint()
		var bLower, bUpper = s.other.SizeHint()
		var lower = saturatingAdd(aLower, bLower)
		var upper gust.Option[uint]
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
	return 0, gust.Some[uint](0)
}

func (s *chainIterator[T]) realCount() uint {
	var aCount uint
	if s.inner != nil {
		aCount = s.inner.Count()
	}
	var bCount uint
	if s.other != nil {
		bCount = s.other.Count()
	}
	return aCount + bCount
}
