package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[KV[any]]       = (*EnumerateIterator[any])(nil)
	_ iRealNext[KV[any]]      = (*EnumerateIterator[any])(nil)
	_ iRealSizeHint           = (*EnumerateIterator[any])(nil)
	_ iRealNth[KV[any]]       = (*EnumerateIterator[any])(nil)
	_ iRealCount              = (*EnumerateIterator[any])(nil)
	_ iRealTryFold[KV[any]]   = (*EnumerateIterator[any])(nil)
	_ iRealFold[KV[any]]      = (*EnumerateIterator[any])(nil)
	_ iRealAdvanceBy[KV[any]] = (*EnumerateIterator[any])(nil)
)

func newEnumerateIterator[T any](iter Iterator[T]) *EnumerateIterator[T] {
	p := &EnumerateIterator[T]{iter: iter}
	p.setFacade(p)
	return p
}

type (
	EnumerateIterator[T any] struct {
		deIterBackground[KV[T]]
		iter  Iterator[T]
		count uint
	}
	KV[T any] struct {
		Index uint
		Value T
	}
)

func (f *EnumerateIterator[T]) realNextBack() gust.Option[KV[T]] {
	panic("unreachable")
}

func (f *EnumerateIterator[T]) realNext() gust.Option[KV[T]] {
	var a = f.iter.Next()
	if a.IsNone() {
		return gust.None[KV[T]]()
	}
	var i = f.count
	f.count += 1
	return gust.Some(KV[T]{Index: i, Value: a.Unwrap()})
}

func (f *EnumerateIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	return f.iter.SizeHint()
}

func (f *EnumerateIterator[T]) realNth(n uint) gust.Option[KV[T]] {
	var a = f.iter.Nth(n)
	if a.IsNone() {
		return gust.None[KV[T]]()
	}
	var i = f.count + n
	f.count = i + 1
	return gust.Some(KV[T]{Index: i, Value: a.Unwrap()})
}

func (f *EnumerateIterator[T]) realCount() uint {
	return f.iter.Count()
}

func (f *EnumerateIterator[T]) realTryFold(acc any, fold func(any, KV[T]) gust.Result[any]) gust.Result[any] {
	return f.iter.TryFold(acc, func(acc any, item T) gust.Result[any] {
		var r = fold(acc, KV[T]{Index: f.count, Value: item})
		f.count += 1
		return r
	})
}

func (f *EnumerateIterator[T]) realFold(acc any, fold func(any, KV[T]) any) any {
	return f.iter.Fold(acc, func(acc any, item T) any {
		var r = fold(acc, KV[T]{Index: f.count, Value: item})
		f.count += 1
		return r
	})
}

func (f *EnumerateIterator[T]) realAdvanceBy(n uint) gust.Errable[uint] {
	var ret = f.iter.AdvanceBy(n)
	if !ret.IsErr() {
		f.count += n
		return ret
	}
	f.count += ret.UnwrapErr()
	return ret
}

var (
	_ DeIterator[KV[any]]         = (*EnumerateDeIterator[any])(nil)
	_ iRealNext[KV[any]]          = (*EnumerateDeIterator[any])(nil)
	_ iRealSizeHint               = (*EnumerateDeIterator[any])(nil)
	_ iRealNth[KV[any]]           = (*EnumerateDeIterator[any])(nil)
	_ iRealCount                  = (*EnumerateDeIterator[any])(nil)
	_ iRealTryFold[KV[any]]       = (*EnumerateDeIterator[any])(nil)
	_ iRealFold[KV[any]]          = (*EnumerateDeIterator[any])(nil)
	_ iRealAdvanceBy[KV[any]]     = (*EnumerateDeIterator[any])(nil)
	_ iRealNextBack[KV[any]]      = (*EnumerateDeIterator[any])(nil)
	_ iRealNthBack[KV[any]]       = (*EnumerateDeIterator[any])(nil)
	_ iRealTryRfold[KV[any]]      = (*EnumerateDeIterator[any])(nil)
	_ iRealRfold[KV[any]]         = (*EnumerateDeIterator[any])(nil)
	_ iRealAdvanceBackBy[KV[any]] = (*EnumerateDeIterator[any])(nil)
	_ iRealRemaining              = (*EnumerateDeIterator[any])(nil)
)

func newEnumerateDeIterator[T any](iter DeIterator[T]) *EnumerateDeIterator[T] {
	p := &EnumerateDeIterator[T]{}
	p.iter = iter
	p.setFacade(p)
	return p
}

// EnumerateDeIterator double ended fuse iterator
type EnumerateDeIterator[T any] struct {
	EnumerateIterator[T]
}

func (d *EnumerateDeIterator[T]) realNextBack() gust.Option[KV[T]] {
	var sizeDeIter = d.iter.(DeIterator[T])
	var a = sizeDeIter.NextBack()
	if a.IsNone() {
		return gust.None[KV[T]]()
	}
	return gust.Some(KV[T]{Index: d.count + sizeDeIter.Remaining(), Value: a.Unwrap()})
}

func (d *EnumerateDeIterator[T]) realNthBack(n uint) gust.Option[KV[T]] {
	var sizeDeIter = d.iter.(DeIterator[T])
	var a = sizeDeIter.NthBack(n)
	if a.IsNone() {
		return gust.None[KV[T]]()
	}
	return gust.Some(KV[T]{Index: d.count + sizeDeIter.Remaining(), Value: a.Unwrap()})
}

func (d *EnumerateDeIterator[T]) realTryRfold(acc any, fold func(any, KV[T]) gust.Result[any]) gust.Result[any] {
	var sizeDeIter = d.iter.(DeIterator[T])
	var count = d.count + sizeDeIter.Remaining()
	return sizeDeIter.TryRfold(acc, func(acc any, item T) gust.Result[any] {
		count -= 1
		return fold(acc, KV[T]{Index: count, Value: item})
	})
}

func (d *EnumerateDeIterator[T]) realRfold(acc any, fold func(any, KV[T]) any) any {
	var sizeDeIter = d.iter.(DeIterator[T])
	var count = d.count + sizeDeIter.Remaining()
	return sizeDeIter.Rfold(acc, func(acc any, item T) any {
		count -= 1
		return fold(acc, KV[T]{Index: count, Value: item})
	})
}

func (d *EnumerateDeIterator[T]) realAdvanceBackBy(n uint) gust.Errable[uint] {
	return d.iter.(DeIterator[T]).AdvanceBackBy(n)
}

func (d *EnumerateDeIterator[T]) realRemaining() uint {
	return d.iter.(DeIterator[T]).Remaining()
}
