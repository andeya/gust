package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[gust.VecEntry[any]]       = (*enumerateIterator[any])(nil)
	_ iRealNext[gust.VecEntry[any]]      = (*enumerateIterator[any])(nil)
	_ iRealSizeHint                      = (*enumerateIterator[any])(nil)
	_ iRealNth[gust.VecEntry[any]]       = (*enumerateIterator[any])(nil)
	_ iRealCount                         = (*enumerateIterator[any])(nil)
	_ iRealTryFold[gust.VecEntry[any]]   = (*enumerateIterator[any])(nil)
	_ iRealFold[gust.VecEntry[any]]      = (*enumerateIterator[any])(nil)
	_ iRealAdvanceBy[gust.VecEntry[any]] = (*enumerateIterator[any])(nil)
)

func newEnumerateIterator[T any](iter Iterator[T]) Iterator[gust.VecEntry[T]] {
	p := &enumerateIterator[T]{iter: iter}
	p.setFacade(p)
	return p
}

type (
	// enumerateIterator is an iterator that yields the current count and the element during iteration.
	enumerateIterator[T any] struct {
		deIterBackground[gust.VecEntry[T]]
		iter  Iterator[T]
		count uint
	}
)

func (f *enumerateIterator[T]) realNextBack() gust.Option[gust.VecEntry[T]] {
	panic("unreachable")
}

func (f *enumerateIterator[T]) realNext() gust.Option[gust.VecEntry[T]] {
	var a = f.iter.Next()
	if a.IsNone() {
		return gust.None[gust.VecEntry[T]]()
	}
	var i = f.count
	f.count += 1
	return gust.Some(gust.VecEntry[T]{Index: int(i), Elem: a.Unwrap()})
}

func (f *enumerateIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	return f.iter.SizeHint()
}

func (f *enumerateIterator[T]) realNth(n uint) gust.Option[gust.VecEntry[T]] {
	var a = f.iter.Nth(n)
	if a.IsNone() {
		return gust.None[gust.VecEntry[T]]()
	}
	var i = f.count + n
	f.count = i + 1
	return gust.Some(gust.VecEntry[T]{Index: int(i), Elem: a.Unwrap()})
}

func (f *enumerateIterator[T]) realCount() uint {
	return f.iter.Count()
}

func (f *enumerateIterator[T]) realTryFold(acc any, fold func(any, gust.VecEntry[T]) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return f.iter.TryFold(acc, func(acc any, item T) gust.AnyCtrlFlow {
		var r = fold(acc, gust.VecEntry[T]{Index: int(f.count), Elem: item})
		f.count += 1
		return r
	})
}

func (f *enumerateIterator[T]) realFold(acc any, fold func(any, gust.VecEntry[T]) any) any {
	return f.iter.Fold(acc, func(acc any, item T) any {
		var r = fold(acc, gust.VecEntry[T]{Index: int(f.count), Elem: item})
		f.count += 1
		return r
	})
}

func (f *enumerateIterator[T]) realAdvanceBy(n uint) gust.Errable[uint] {
	var ret = f.iter.AdvanceBy(n)
	if !ret.IsErr() {
		f.count += n
		return ret
	}
	f.count += ret.UnwrapErr()
	return ret
}

var (
	_ DeIterator[gust.VecEntry[any]]         = (*deEnumerateIterator[any])(nil)
	_ iRealNext[gust.VecEntry[any]]          = (*deEnumerateIterator[any])(nil)
	_ iRealSizeHint                          = (*deEnumerateIterator[any])(nil)
	_ iRealNth[gust.VecEntry[any]]           = (*deEnumerateIterator[any])(nil)
	_ iRealCount                             = (*deEnumerateIterator[any])(nil)
	_ iRealTryFold[gust.VecEntry[any]]       = (*deEnumerateIterator[any])(nil)
	_ iRealFold[gust.VecEntry[any]]          = (*deEnumerateIterator[any])(nil)
	_ iRealAdvanceBy[gust.VecEntry[any]]     = (*deEnumerateIterator[any])(nil)
	_ iRealNextBack[gust.VecEntry[any]]      = (*deEnumerateIterator[any])(nil)
	_ iRealNthBack[gust.VecEntry[any]]       = (*deEnumerateIterator[any])(nil)
	_ iRealTryRfold[gust.VecEntry[any]]      = (*deEnumerateIterator[any])(nil)
	_ iRealRfold[gust.VecEntry[any]]         = (*deEnumerateIterator[any])(nil)
	_ iRealAdvanceBackBy[gust.VecEntry[any]] = (*deEnumerateIterator[any])(nil)
	_ iRealRemaining                         = (*deEnumerateIterator[any])(nil)
)

func newDeEnumerateIterator[T any](iter DeIterator[T]) DeIterator[gust.VecEntry[T]] {
	p := &deEnumerateIterator[T]{}
	p.iter = iter
	p.setFacade(p)
	return p
}

// deEnumerateIterator is an iterator that yields the current count and the element during iteration.
type deEnumerateIterator[T any] struct {
	enumerateIterator[T]
}

func (d *deEnumerateIterator[T]) realNextBack() gust.Option[gust.VecEntry[T]] {
	var sizeDeIter = d.iter.(DeIterator[T])
	var a = sizeDeIter.NextBack()
	if a.IsNone() {
		return gust.None[gust.VecEntry[T]]()
	}
	return gust.Some(gust.VecEntry[T]{Index: int(d.count + sizeDeIter.Remaining()), Elem: a.Unwrap()})
}

func (d *deEnumerateIterator[T]) realNthBack(n uint) gust.Option[gust.VecEntry[T]] {
	var sizeDeIter = d.iter.(DeIterator[T])
	var a = sizeDeIter.NthBack(n)
	if a.IsNone() {
		return gust.None[gust.VecEntry[T]]()
	}
	return gust.Some(gust.VecEntry[T]{Index: int(d.count + sizeDeIter.Remaining()), Elem: a.Unwrap()})
}

func (d *deEnumerateIterator[T]) realTryRfold(init any, fold func(any, gust.VecEntry[T]) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	var sizeDeIter = d.iter.(DeIterator[T])
	var count = d.count + sizeDeIter.Remaining()
	return sizeDeIter.TryRfold(init, func(acc any, item T) gust.AnyCtrlFlow {
		count -= 1
		return fold(acc, gust.VecEntry[T]{Index: int(count), Elem: item})
	})
}

func (d *deEnumerateIterator[T]) realRfold(acc any, fold func(any, gust.VecEntry[T]) any) any {
	var sizeDeIter = d.iter.(DeIterator[T])
	var count = d.count + sizeDeIter.Remaining()
	return sizeDeIter.Rfold(acc, func(acc any, item T) any {
		count -= 1
		return fold(acc, gust.VecEntry[T]{Index: int(count), Elem: item})
	})
}

func (d *deEnumerateIterator[T]) realAdvanceBackBy(n uint) gust.Errable[uint] {
	return d.iter.(DeIterator[T]).AdvanceBackBy(n)
}

func (d *deEnumerateIterator[T]) realRemaining() uint {
	return d.iter.(DeIterator[T]).Remaining()
}
