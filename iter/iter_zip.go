package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[Pair[any, any]]  = (*ZipIterator[any, any])(nil)
	_ iRealNext[Pair[any, any]] = (*ZipIterator[any, any])(nil)
	_ iRealSizeHint             = (*ZipIterator[any, any])(nil)
	_ iRealNth[Pair[any, any]]  = (*ZipIterator[any, any])(nil)
)

func newZipIterator[A any, B any](a Iterator[A], b Iterator[B]) *ZipIterator[A, B] {
	p := &ZipIterator[A, B]{a: a, b: b}
	p.facade = p
	return p
}

type (
	ZipIterator[A any, B any] struct {
		iterTrait[Pair[A, B]]
		a Iterator[A]
		b Iterator[B]
	}
	Pair[A any, B any] struct {
		A A
		B B
	}
)

func (s ZipIterator[A, B]) SuperNth(n uint) gust.Option[Pair[A, B]] {
	for {
		p := s.Next()
		if p.IsNone() {
			return gust.None[Pair[A, B]]()
		}
		if n == 0 {
			return p
		}
		n -= 1
	}
}

func (s ZipIterator[A, B]) realNext() gust.Option[Pair[A, B]] {
	var x = s.a.Next()
	if x.IsNone() {
		return gust.None[Pair[A, B]]()
	}
	var y = s.b.Next()
	if y.IsNone() {
		return gust.None[Pair[A, B]]()
	}
	return gust.Some(Pair[A, B]{A: x.Unwrap(), B: y.Unwrap()})
}

func (s ZipIterator[A, B]) realSizeHint() (uint, gust.Option[uint]) {
	var aLower, aUpper = s.a.SizeHint()
	var bLower, bUpper = s.b.SizeHint()

	var lower = aLower
	if lower > bLower {
		lower = bLower
	}

	var upper gust.Option[uint]
	if aUpper.IsSome() && bUpper.IsSome() {
		if aUpper.Unwrap() <= bUpper.Unwrap() {
			upper = aUpper
		} else {
			upper = bUpper
		}
	} else if aUpper.IsSome() {
		upper = aUpper
	} else if bUpper.IsSome() {
		upper = bUpper
	}
	return lower, upper
}

func (s ZipIterator[A, B]) realNth(n uint) gust.Option[Pair[A, B]] {
	return s.SuperNth(n)
}

var (
	_ Iterator[Pair[any, any]]             = (*DoubleEndedZipIterator[any, any])(nil)
	_ DoubleEndedIterator[Pair[any, any]]  = (*DoubleEndedZipIterator[any, any])(nil)
	_ iRealSizeHint                        = (*DoubleEndedZipIterator[any, any])(nil)
	_ iRealNth[Pair[any, any]]             = (*DoubleEndedZipIterator[any, any])(nil)
	_ iRealDoubleEndedNext[Pair[any, any]] = (*DoubleEndedZipIterator[any, any])(nil)
)

func newDoubleEndedZipIterator[A any, B any](a DoubleEndedIterator[A], b DoubleEndedIterator[B]) *DoubleEndedZipIterator[A, B] {
	p := &DoubleEndedZipIterator[A, B]{a: a, b: b}
	p.setFacade(p)
	return p
}

type DoubleEndedZipIterator[A any, B any] struct {
	doubleEndedIterTrait[Pair[A, B]]
	a DoubleEndedIterator[A]
	b DoubleEndedIterator[B]
}

func (s DoubleEndedZipIterator[A, B]) realRemainingLen() uint {
	aLen := s.a.RemainingLen()
	bLen := s.b.RemainingLen()
	if aLen < bLen {
		return aLen
	}
	return bLen
}

func (s DoubleEndedZipIterator[A, B]) realNextBack() gust.Option[Pair[A, B]] {
	var aLen = s.a.RemainingLen()
	var bLen = s.b.RemainingLen()
	if aLen != bLen {
		// Adjust a, b to equal length
		if aLen > bLen {
			u := aLen - bLen
			for i := uint(0); i < u; i++ {
				s.a.NextBack()
			}
		} else {
			u := bLen - aLen
			for i := uint(0); i < u; i++ {
				s.b.NextBack()
			}
		}
	}
	var x = s.a.NextBack()
	var y = s.b.NextBack()
	if x.IsSome() && y.IsSome() {
		return gust.Some(Pair[A, B]{A: x.Unwrap(), B: y.Unwrap()})
	}
	return gust.None[Pair[A, B]]()
}

func (s DoubleEndedZipIterator[A, B]) SuperNth(n uint) gust.Option[Pair[A, B]] {
	for {
		p := s.Next()
		if p.IsNone() {
			return gust.None[Pair[A, B]]()
		}
		if n == 0 {
			return p
		}
		n -= 1
	}
}

func (s DoubleEndedZipIterator[A, B]) realNext() gust.Option[Pair[A, B]] {
	var x = s.a.Next()
	if x.IsNone() {
		return gust.None[Pair[A, B]]()
	}
	var y = s.b.Next()
	if y.IsNone() {
		return gust.None[Pair[A, B]]()
	}
	return gust.Some(Pair[A, B]{A: x.Unwrap(), B: y.Unwrap()})
}

func (s DoubleEndedZipIterator[A, B]) realSizeHint() (uint, gust.Option[uint]) {
	var aLower, aUpper = s.a.SizeHint()
	var bLower, bUpper = s.b.SizeHint()

	var lower = aLower
	if lower > bLower {
		lower = bLower
	}

	var upper gust.Option[uint]
	if aUpper.IsSome() && bUpper.IsSome() {
		if aUpper.Unwrap() <= bUpper.Unwrap() {
			upper = aUpper
		} else {
			upper = bUpper
		}
	} else if aUpper.IsSome() {
		upper = aUpper
	} else if bUpper.IsSome() {
		upper = bUpper
	}
	return lower, upper
}

func (s DoubleEndedZipIterator[A, B]) realNth(n uint) gust.Option[Pair[A, B]] {
	return s.SuperNth(n)
}
