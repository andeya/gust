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

func (s ZipIterator[A, B]) realSizeHint() (uint64, gust.Option[uint64]) {
	var aLower, aUpper = s.a.SizeHint()
	var bLower, bUpper = s.b.SizeHint()

	var lower = aLower
	if lower > bLower {
		lower = bLower
	}

	var upper gust.Option[uint64]
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

//
// func (s ZipIterator[A, B]) realSizeHint() (uint64, gust.Option[uint64]) {
// 	return s.iter.SizeHint()
// }
//
// func (s ZipIterator[A, B]) realFold(init any, g func(any, B) any) any {
// 	return Fold[A, any](s.iter, init, func(acc any, elt A) any { return g(acc, s.f(elt)) })
// }
//
// func (s ZipIterator[A, B]) realTryFold(init any, g func(any, B) gust.Result[any]) gust.Result[any] {
// 	return TryFold[A, any](s.iter, init, func(acc any, elt A) gust.Result[any] { return g(acc, s.f(elt)) })
// }
//
