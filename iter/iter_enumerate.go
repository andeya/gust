package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[gust.KV[any]]       = (*enumerateIterator[any])(nil)
	_ iRealNext[gust.KV[any]]      = (*enumerateIterator[any])(nil)
	_ iRealSizeHint                = (*enumerateIterator[any])(nil)
	_ iRealNth[gust.KV[any]]       = (*enumerateIterator[any])(nil)
	_ iRealCount                   = (*enumerateIterator[any])(nil)
	_ iRealTryFold[gust.KV[any]]   = (*enumerateIterator[any])(nil)
	_ iRealFold[gust.KV[any]]      = (*enumerateIterator[any])(nil)
	_ iRealAdvanceBy[gust.KV[any]] = (*enumerateIterator[any])(nil)
)

func newEnumerateIterator[T any](iter Iterator[T]) Iterator[gust.KV[T]] {
	p := &enumerateIterator[T]{iter: iter}
	p.setFacade(p)
	return p
}

type (
	// enumerateIterator is an iterator that yields the current count and the element during iteration.
	enumerateIterator[T any] struct {
		deIterBackground[gust.KV[T]]
		iter  Iterator[T]
		count uint
	}
)

func (f *enumerateIterator[T]) realNextBack() gust.Option[gust.KV[T]] {
	panic("unreachable")
}

func (f *enumerateIterator[T]) realNext() gust.Option[gust.KV[T]] {
	var a = f.iter.Next()
	if a.IsNone() {
		return gust.None[gust.KV[T]]()
	}
	var i = f.count
	f.count += 1
	return gust.Some(gust.KV[T]{Index: i, Value: a.Unwrap()})
}

func (f *enumerateIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	return f.iter.SizeHint()
}

func (f *enumerateIterator[T]) realNth(n uint) gust.Option[gust.KV[T]] {
	var a = f.iter.Nth(n)
	if a.IsNone() {
		return gust.None[gust.KV[T]]()
	}
	var i = f.count + n
	f.count = i + 1
	return gust.Some(gust.KV[T]{Index: i, Value: a.Unwrap()})
}

func (f *enumerateIterator[T]) realCount() uint {
	return f.iter.Count()
}

func (f *enumerateIterator[T]) realTryFold(acc any, fold func(any, gust.KV[T]) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return f.iter.TryFold(acc, func(acc any, item T) gust.AnyCtrlFlow {
		var r = fold(acc, gust.KV[T]{Index: f.count, Value: item})
		f.count += 1
		return r
	})
}

func (f *enumerateIterator[T]) realFold(acc any, fold func(any, gust.KV[T]) any) any {
	return f.iter.Fold(acc, func(acc any, item T) any {
		var r = fold(acc, gust.KV[T]{Index: f.count, Value: item})
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
	_ DeIterator[gust.KV[any]]         = (*deEnumerateIterator[any])(nil)
	_ iRealNext[gust.KV[any]]          = (*deEnumerateIterator[any])(nil)
	_ iRealSizeHint                    = (*deEnumerateIterator[any])(nil)
	_ iRealNth[gust.KV[any]]           = (*deEnumerateIterator[any])(nil)
	_ iRealCount                       = (*deEnumerateIterator[any])(nil)
	_ iRealTryFold[gust.KV[any]]       = (*deEnumerateIterator[any])(nil)
	_ iRealFold[gust.KV[any]]          = (*deEnumerateIterator[any])(nil)
	_ iRealAdvanceBy[gust.KV[any]]     = (*deEnumerateIterator[any])(nil)
	_ iRealNextBack[gust.KV[any]]      = (*deEnumerateIterator[any])(nil)
	_ iRealNthBack[gust.KV[any]]       = (*deEnumerateIterator[any])(nil)
	_ iRealTryRfold[gust.KV[any]]      = (*deEnumerateIterator[any])(nil)
	_ iRealRfold[gust.KV[any]]         = (*deEnumerateIterator[any])(nil)
	_ iRealAdvanceBackBy[gust.KV[any]] = (*deEnumerateIterator[any])(nil)
	_ iRealRemaining                   = (*deEnumerateIterator[any])(nil)
)

func newDeEnumerateIterator[T any](iter DeIterator[T]) DeIterator[gust.KV[T]] {
	p := &deEnumerateIterator[T]{}
	p.iter = iter
	p.setFacade(p)
	return p
}

// deEnumerateIterator is an iterator that yields the current count and the element during iteration.
type deEnumerateIterator[T any] struct {
	enumerateIterator[T]
}

func (d *deEnumerateIterator[T]) realNextBack() gust.Option[gust.KV[T]] {
	var sizeDeIter = d.iter.(DeIterator[T])
	var a = sizeDeIter.NextBack()
	if a.IsNone() {
		return gust.None[gust.KV[T]]()
	}
	return gust.Some(gust.KV[T]{Index: d.count + sizeDeIter.Remaining(), Value: a.Unwrap()})
}

func (d *deEnumerateIterator[T]) realNthBack(n uint) gust.Option[gust.KV[T]] {
	var sizeDeIter = d.iter.(DeIterator[T])
	var a = sizeDeIter.NthBack(n)
	if a.IsNone() {
		return gust.None[gust.KV[T]]()
	}
	return gust.Some(gust.KV[T]{Index: d.count + sizeDeIter.Remaining(), Value: a.Unwrap()})
}

func (d *deEnumerateIterator[T]) realTryRfold(init any, fold func(any, gust.KV[T]) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	var sizeDeIter = d.iter.(DeIterator[T])
	var count = d.count + sizeDeIter.Remaining()
	return sizeDeIter.TryRfold(init, func(acc any, item T) gust.AnyCtrlFlow {
		count -= 1
		return fold(acc, gust.KV[T]{Index: count, Value: item})
	})
}

func (d *deEnumerateIterator[T]) realRfold(acc any, fold func(any, gust.KV[T]) any) any {
	var sizeDeIter = d.iter.(DeIterator[T])
	var count = d.count + sizeDeIter.Remaining()
	return sizeDeIter.Rfold(acc, func(acc any, item T) any {
		count -= 1
		return fold(acc, gust.KV[T]{Index: count, Value: item})
	})
}

func (d *deEnumerateIterator[T]) realAdvanceBackBy(n uint) gust.Errable[uint] {
	return d.iter.(DeIterator[T]).AdvanceBackBy(n)
}

func (d *deEnumerateIterator[T]) realRemaining() uint {
	return d.iter.(DeIterator[T]).Remaining()
}
