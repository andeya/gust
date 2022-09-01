package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[gust.Pair[any, any]]  = (*ZipIterator[any, any])(nil)
	_ iRealNext[gust.Pair[any, any]] = (*ZipIterator[any, any])(nil)
	_ iRealSizeHint                  = (*ZipIterator[any, any])(nil)
	_ iRealNth[gust.Pair[any, any]]  = (*ZipIterator[any, any])(nil)
)

func newZipIterator[A any, B any](a Iterator[A], b Iterator[B]) *ZipIterator[A, B] {
	p := &ZipIterator[A, B]{a: a, b: b}
	p.facade = p
	return p
}

type (
	ZipIterator[A any, B any] struct {
		iterTrait[gust.Pair[A, B]]
		a Iterator[A]
		b Iterator[B]
	}
)

func (s ZipIterator[A, B]) SuperNth(n uint) gust.Option[gust.Pair[A, B]] {
	for {
		p := s.Next()
		if p.IsNone() {
			return gust.None[gust.Pair[A, B]]()
		}
		if n == 0 {
			return p
		}
		n -= 1
	}
}

func (s ZipIterator[A, B]) realNext() gust.Option[gust.Pair[A, B]] {
	var x = s.a.Next()
	if x.IsNone() {
		return gust.None[gust.Pair[A, B]]()
	}
	var y = s.b.Next()
	if y.IsNone() {
		return gust.None[gust.Pair[A, B]]()
	}
	return gust.Some(gust.Pair[A, B]{A: x.Unwrap(), B: y.Unwrap()})
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

func (s ZipIterator[A, B]) realNth(n uint) gust.Option[gust.Pair[A, B]] {
	return s.SuperNth(n)
}

var (
	_ SizeDeIterator[gust.Pair[any, any]]  = (*ZipSizeDeIterator[any, any])(nil)
	_ iRealSizeHint                        = (*ZipSizeDeIterator[any, any])(nil)
	_ iRealNth[gust.Pair[any, any]]        = (*ZipSizeDeIterator[any, any])(nil)
	_ iRealDeIterable[gust.Pair[any, any]] = (*ZipSizeDeIterator[any, any])(nil)
)

func newZipSizeDeIterator[A any, B any](a SizeDeIterator[A], b SizeDeIterator[B]) *ZipSizeDeIterator[A, B] {
	p := &ZipSizeDeIterator[A, B]{a: a, b: b}
	p.setFacade(p)
	return p
}

// ZipSizeDeIterator is a double-ended 2-in-1 iterator with explicit size
type ZipSizeDeIterator[A any, B any] struct {
	sizeDeIterTrait[gust.Pair[A, B]]
	a SizeDeIterator[A]
	b SizeDeIterator[B]
}

func (s ZipSizeDeIterator[A, B]) realRemaining() uint {
	aLen := s.a.Remaining()
	bLen := s.b.Remaining()
	if aLen < bLen {
		return aLen
	}
	return bLen
}

func (s ZipSizeDeIterator[A, B]) realNextBack() gust.Option[gust.Pair[A, B]] {
	var aLen = s.a.Remaining()
	var bLen = s.b.Remaining()
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
		return gust.Some(gust.Pair[A, B]{A: x.Unwrap(), B: y.Unwrap()})
	}
	return gust.None[gust.Pair[A, B]]()
}

func (s ZipSizeDeIterator[A, B]) SuperNth(n uint) gust.Option[gust.Pair[A, B]] {
	for {
		p := s.Next()
		if p.IsNone() {
			return gust.None[gust.Pair[A, B]]()
		}
		if n == 0 {
			return p
		}
		n -= 1
	}
}

func (s ZipSizeDeIterator[A, B]) realNext() gust.Option[gust.Pair[A, B]] {
	var x = s.a.Next()
	if x.IsNone() {
		return gust.None[gust.Pair[A, B]]()
	}
	var y = s.b.Next()
	if y.IsNone() {
		return gust.None[gust.Pair[A, B]]()
	}
	return gust.Some(gust.Pair[A, B]{A: x.Unwrap(), B: y.Unwrap()})
}

func (s ZipSizeDeIterator[A, B]) realSizeHint() (uint, gust.Option[uint]) {
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

func (s ZipSizeDeIterator[A, B]) realNth(n uint) gust.Option[gust.Pair[A, B]] {
	return s.SuperNth(n)
}
