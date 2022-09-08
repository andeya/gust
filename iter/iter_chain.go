package iter

import (
	"github.com/andeya/gust"
	"github.com/andeya/gust/digit"
)

var (
	_ Iterator[any]       = (*chainIterator[any])(nil)
	_ iRealNext[any]      = (*chainIterator[any])(nil)
	_ iRealCount          = (*chainIterator[any])(nil)
	_ iRealTryFold[any]   = (*chainIterator[any])(nil)
	_ iRealFold[any]      = (*chainIterator[any])(nil)
	_ iRealAdvanceBy[any] = (*chainIterator[any])(nil)
	_ iRealNth[any]       = (*chainIterator[any])(nil)
	_ iRealFind[any]      = (*chainIterator[any])(nil)
	_ iRealLast[any]      = (*chainIterator[any])(nil)
	_ iRealSizeHint       = (*chainIterator[any])(nil)
)

func newChainIterator[T any](a Iterator[T], b Iterator[T]) Iterator[T] {
	iter := &chainIterator[T]{a: a, b: b}
	iter.setFacade(iter)
	return iter
}

type chainIterator[T any] struct {
	deIterBackground[T]
	a Iterator[T]
	b Iterator[T]
}

func (s *chainIterator[T]) realNextBack() gust.Option[T] {
	panic("unreachable")
}

func (s *chainIterator[T]) realLast() gust.Option[T] {
	// Must exhaust a before b.
	var aLast gust.Option[T]
	var bLast gust.Option[T]
	if s.a != nil {
		aLast = s.a.Last()
	}
	if s.b != nil {
		bLast = s.b.Last()
	}
	if bLast.IsSome() {
		return bLast
	}
	return aLast
}

func (s *chainIterator[T]) realFind(predicate func(T) bool) gust.Option[T] {
	if s.a != nil {
		item := s.a.Find(predicate)
		if item.IsSome() {
			return item
		}
		s.a = nil
	}
	if s.b != nil {
		return s.b.Find(predicate)
	}
	return gust.None[T]()
}

func (s *chainIterator[T]) realNth(n uint) gust.Option[T] {
	if s.a != nil {
		r := s.a.AdvanceBy(n)
		if r.IsErr() {
			n -= r.UnwrapErr()
		} else {
			item := s.a.Next()
			if item.IsSome() {
				return item
			}
			n = 0
		}
		s.a = nil
	}
	if s.b != nil {
		return s.b.Nth(n)
	}
	return gust.None[T]()
}

func (s *chainIterator[T]) realAdvanceBy(n uint) gust.Errable[uint] {
	var rem = n
	if s.a != nil {
		r := s.a.AdvanceBy(rem)
		if !r.IsErr() {
			return r
		}
		rem -= r.UnwrapErr()
		s.a = nil
	}
	if s.b != nil {
		r := s.b.AdvanceBy(rem)
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
	if s.a != nil {
		acc = s.a.Fold(acc, f)
	}
	if s.b != nil {
		acc = s.b.Fold(acc, f)
	}
	return acc
}

func (s *chainIterator[T]) realTryFold(acc any, f func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	if s.a != nil {
		r := s.a.TryFold(acc, f)
		if r.IsBreak() {
			return r
		}
		acc = r.UnwrapContinue()
		s.a = nil
	}
	if s.b != nil {
		r := s.b.TryFold(acc, f)
		if r.IsBreak() {
			return r
		}
		acc = r.UnwrapContinue()
		// we don't fuse the second iterator
	}
	return gust.AnyContinue(acc)
}

func (s *chainIterator[T]) realNext() gust.Option[T] {
	if s.a != nil {
		item := s.a.Next()
		if item.IsSome() {
			return item
		}
		s.a = nil
	}
	if s.b != nil {
		return s.b.Next()
	}
	return gust.None[T]()
}

func (s *chainIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	if s.a != nil && s.b != nil {
		var aLower, aUpper = s.a.SizeHint()
		var bLower, bUpper = s.b.SizeHint()
		var lower = digit.SaturatingAdd(aLower, bLower)
		var upper gust.Option[uint]
		if aUpper.IsSome() && bUpper.IsSome() {
			upper = digit.CheckedAdd(aUpper.Unwrap(), bUpper.Unwrap())
		}
		return lower, upper
	}
	if s.a != nil && s.b == nil {
		return s.a.SizeHint()
	}
	if s.a == nil && s.b != nil {
		return s.b.SizeHint()
	}
	return 0, gust.Some[uint](0)
}

func (s *chainIterator[T]) realCount() uint {
	var aCount uint
	if s.a != nil {
		aCount = s.a.Count()
	}
	var bCount uint
	if s.b != nil {
		bCount = s.b.Count()
	}
	return aCount + bCount
}

var (
	_ DeIterator[any]         = (*deChainIterator[any])(nil)
	_ iRealRemaining          = (*deChainIterator[any])(nil)
	_ iRealNextBack[any]      = (*deChainIterator[any])(nil)
	_ iRealAdvanceBackBy[any] = (*deChainIterator[any])(nil)
	_ iRealNthBack[any]       = (*deChainIterator[any])(nil)
	_ iRealRfind[any]         = (*deChainIterator[any])(nil)
	_ iRealTryRfold[any]      = (*deChainIterator[any])(nil)
	_ iRealRfold[any]         = (*deChainIterator[any])(nil)
)

func newDeChainIterator[T any](a DeIterator[T], b DeIterator[T]) DeIterator[T] {
	iter := &deChainIterator[T]{chainIterator: chainIterator[T]{a: a, b: b}}
	iter.setFacade(iter)
	return iter
}

type deChainIterator[T any] struct {
	chainIterator[T]
}

func (d *deChainIterator[T]) realRemaining() uint {
	return d.a.(DeIterator[T]).Remaining() + d.b.(DeIterator[T]).Remaining()
}

func andThenOrClear[T any, U any](opt *Iterator[T], f func(Iterator[T]) gust.Option[U]) gust.Option[U] {
	if *opt == nil {
		return gust.None[U]()
	}
	var x = f(*opt)
	if x.IsNone() {
		*opt = nil
	}
	return x
}

func (d *deChainIterator[T]) realNextBack() gust.Option[T] {
	return andThenOrClear[T, T](&d.b, func(iter Iterator[T]) gust.Option[T] {
		return iter.(DeIterator[T]).NextBack()
	}).OrElse(func() gust.Option[T] {
		if d.a == nil {
			return gust.None[T]()
		}
		return d.a.(DeIterator[T]).NextBack()
	})
}

func (d *deChainIterator[T]) realAdvanceBackBy(n uint) gust.Errable[uint] {
	var rem = n
	if d.b != nil {
		var r = d.b.(DeIterator[T]).AdvanceBackBy(rem)
		if r.IsOk() {
			return gust.NonErrable[uint]()
		}
		rem -= r.UnwrapErr()
		d.b = nil
	}
	if d.a != nil {
		var r = d.a.(DeIterator[T]).AdvanceBackBy(rem)
		if r.IsOk() {
			return gust.NonErrable[uint]()
		}
		rem -= r.UnwrapErr()
		// we don't fuse the second iterator
	}
	if rem == 0 {
		return gust.NonErrable[uint]()
	}
	return gust.ToErrable(n - rem)
}

func (d *deChainIterator[T]) realNthBack(n uint) gust.Option[T] {
	if d.b != nil {
		var b = d.b.(DeIterator[T])
		var r = b.AdvanceBackBy(n)
		if r.IsOk() {
			x := b.NextBack()
			if x.IsSome() {
				return x
			}
			n = 0
		} else {
			n -= r.UnwrapErr()
		}
		d.b = nil
	}
	if d.a != nil {
		return d.a.(DeIterator[T]).NthBack(n)
	}
	return gust.None[T]()
}

func (d *deChainIterator[T]) realRfind(f func(T) bool) gust.Option[T] {
	return andThenOrClear[T, T](&d.b, func(iter Iterator[T]) gust.Option[T] {
		return iter.(DeIterator[T]).Rfind(f)
	}).OrElse(func() gust.Option[T] {
		if d.a == nil {
			return gust.None[T]()
		}
		return d.a.(DeIterator[T]).Rfind(f)
	})
}

func (d *deChainIterator[T]) realTryRfold(acc any, f func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	if d.b != nil {
		b := d.b.(DeIterator[T])
		r := b.TryRfold(acc, f)
		if r.IsBreak() {
			return r
		}
		acc = r.UnwrapContinue()
		d.b = nil
	}
	if d.a != nil {
		a := d.a.(DeIterator[T])
		r := a.TryRfold(acc, f)
		if r.IsBreak() {
			return r
		}
		acc = r.UnwrapContinue()
		// we don't fuse the second iterator
	}
	return gust.AnyContinue(acc)
}

func (d *deChainIterator[T]) realRfold(acc any, f func(any, T) any) any {
	if d.b != nil {
		b := d.b.(DeIterator[T])
		acc = b.Rfold(acc, f)
	}
	if d.a != nil {
		a := d.a.(DeIterator[T])
		acc = a.Rfold(acc, f)
	}
	return acc
}
