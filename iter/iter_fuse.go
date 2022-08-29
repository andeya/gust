package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]     = (*FuseIterator[any])(nil)
	_ iRealNext[any]    = (*FuseIterator[any])(nil)
	_ iRealNth[any]     = (*FuseIterator[any])(nil)
	_ iRealTryFold[any] = (*FuseIterator[any])(nil)
	_ iRealFind[any]    = (*FuseIterator[any])(nil)
)

func newFuseIterator[T any](iter Iterator[T]) *FuseIterator[T] {
	p := &FuseIterator[T]{iter: iter}
	p.setFacade(p)
	return p
}

type FuseIterator[T any] struct {
	iterTrait[T]
	iter   Iterator[T]
	isNone bool
}

func (f FuseIterator[T]) realNext() gust.Option[T] {
	if f.isNone {
		return gust.None[T]()
	}
	return f.iter.Next().InspectNone(func() {
		f.isNone = true
	})
}

func (f FuseIterator[T]) realNth(n uint) gust.Option[T] {
	if f.isNone {
		return gust.None[T]()
	}
	return f.iter.Nth(n).InspectNone(func() {
		f.isNone = true
	})
}

func (f FuseIterator[T]) realTryFold(acc any, fold func(any, T) gust.Result[any]) gust.Result[any] {
	if f.isNone {
		return gust.Err[any]("fuse iterator is empty")
	}
	r := f.iter.TryFold(acc, fold)
	f.isNone = true
	return r
}

func (f FuseIterator[T]) realFind(predicate func(T) bool) gust.Option[T] {
	if f.isNone {
		return gust.None[T]()
	}
	return f.iter.Find(predicate).InspectNone(func() {
		f.isNone = true
	})
}

var (
	_ DoubleEndedIterator[any] = (*DoubleEndedFuseIterator[any])(nil)
	_ iRealNext[any]           = (*DoubleEndedFuseIterator[any])(nil)
	_ iRealNth[any]            = (*DoubleEndedFuseIterator[any])(nil)
	_ iRealTryFold[any]        = (*DoubleEndedFuseIterator[any])(nil)
	_ iRealFind[any]           = (*DoubleEndedFuseIterator[any])(nil)
	_ iRealNextBack[any]       = (*DoubleEndedFuseIterator[any])(nil)
	_ iRealNthBack[any]        = (*DoubleEndedFuseIterator[any])(nil)
	_ iRealTryRfold[any]       = (*DoubleEndedFuseIterator[any])(nil)
	_ iRealRfind[any]          = (*DoubleEndedFuseIterator[any])(nil)
)

func newDoubleEndedFuseIterator[T any](iter DoubleEndedIterator[T]) *DoubleEndedFuseIterator[T] {
	p := &DoubleEndedFuseIterator[T]{iter: iter}
	p.setFacade(p)
	return p
}

type DoubleEndedFuseIterator[T any] struct {
	doubleEndedIterTrait[T]
	iter   DoubleEndedIterator[T]
	isNone bool
}

func (d DoubleEndedFuseIterator[T]) realRemainingLen() uint {
	return d.iter.RemainingLen()
}

func (d DoubleEndedFuseIterator[T]) realNext() gust.Option[T] {
	if d.isNone {
		return gust.None[T]()
	}
	return d.iter.Next().InspectNone(func() {
		d.isNone = true
	})
}

func (d DoubleEndedFuseIterator[T]) realNextBack() gust.Option[T] {
	if d.isNone {
		return gust.None[T]()
	}
	return d.iter.NextBack().InspectNone(func() {
		d.isNone = true
	})
}

func (d DoubleEndedFuseIterator[T]) realNth(n uint) gust.Option[T] {
	if d.isNone {
		return gust.None[T]()
	}
	return d.iter.Nth(n).InspectNone(func() {
		d.isNone = true
	})
}

func (d DoubleEndedFuseIterator[T]) realNthBack(n uint) gust.Option[T] {
	if d.isNone {
		return gust.None[T]()
	}
	return d.iter.NthBack(n).InspectNone(func() {
		d.isNone = true
	})
}

func (d DoubleEndedFuseIterator[T]) realTryFold(acc any, fold func(any, T) gust.Result[any]) gust.Result[any] {
	if d.isNone {
		return gust.Err[any]("fuse iterator is empty")
	}
	r := d.iter.TryFold(acc, fold)
	d.isNone = true
	return r
}

func (d DoubleEndedFuseIterator[T]) realTryRfold(acc any, fold func(any, T) gust.Result[any]) gust.Result[any] {
	if d.isNone {
		return gust.Err[any]("fuse iterator is empty")
	}
	r := d.iter.TryRfold(acc, fold)
	d.isNone = true
	return r
}

func (d DoubleEndedFuseIterator[T]) realFind(predicate func(T) bool) gust.Option[T] {
	if d.isNone {
		return gust.None[T]()
	}
	return d.iter.Find(predicate).InspectNone(func() {
		d.isNone = true
	})
}

func (d DoubleEndedFuseIterator[T]) realRfind(predicate func(T) bool) gust.Option[T] {
	if d.isNone {
		return gust.None[T]()
	}
	return d.iter.Rfind(predicate).InspectNone(func() {
		d.isNone = true
	})
}
