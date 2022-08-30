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
	_ DeIterator[any]    = (*FuseDeIterator[any])(nil)
	_ iRealNext[any]     = (*FuseDeIterator[any])(nil)
	_ iRealNth[any]      = (*FuseDeIterator[any])(nil)
	_ iRealTryFold[any]  = (*FuseDeIterator[any])(nil)
	_ iRealFind[any]     = (*FuseDeIterator[any])(nil)
	_ iRealNextBack[any] = (*FuseDeIterator[any])(nil)
	_ iRealNthBack[any]  = (*FuseDeIterator[any])(nil)
	_ iRealTryRfold[any] = (*FuseDeIterator[any])(nil)
	_ iRealRfind[any]    = (*FuseDeIterator[any])(nil)
)

func newFuseDeIterator[T any](iter DeIterator[T]) *FuseDeIterator[T] {
	p := &FuseDeIterator[T]{iter: iter}
	p.setFacade(p)
	return p
}

type FuseDeIterator[T any] struct {
	sizeDeIterTrait[T]
	iter   DeIterator[T]
	isNone bool
}

func (d FuseDeIterator[T]) realNext() gust.Option[T] {
	if d.isNone {
		return gust.None[T]()
	}
	return d.iter.Next().InspectNone(func() {
		d.isNone = true
	})
}

func (d FuseDeIterator[T]) realNextBack() gust.Option[T] {
	if d.isNone {
		return gust.None[T]()
	}
	return d.iter.NextBack().InspectNone(func() {
		d.isNone = true
	})
}

func (d FuseDeIterator[T]) realNth(n uint) gust.Option[T] {
	if d.isNone {
		return gust.None[T]()
	}
	return d.iter.Nth(n).InspectNone(func() {
		d.isNone = true
	})
}

func (d FuseDeIterator[T]) realNthBack(n uint) gust.Option[T] {
	if d.isNone {
		return gust.None[T]()
	}
	return d.iter.NthBack(n).InspectNone(func() {
		d.isNone = true
	})
}

func (d FuseDeIterator[T]) realTryFold(acc any, fold func(any, T) gust.Result[any]) gust.Result[any] {
	if d.isNone {
		return gust.Err[any]("fuse iterator is empty")
	}
	r := d.iter.TryFold(acc, fold)
	d.isNone = true
	return r
}

func (d FuseDeIterator[T]) realTryRfold(acc any, fold func(any, T) gust.Result[any]) gust.Result[any] {
	if d.isNone {
		return gust.Err[any]("fuse iterator is empty")
	}
	r := d.iter.TryRfold(acc, fold)
	d.isNone = true
	return r
}

func (d FuseDeIterator[T]) realFind(predicate func(T) bool) gust.Option[T] {
	if d.isNone {
		return gust.None[T]()
	}
	return d.iter.Find(predicate).InspectNone(func() {
		d.isNone = true
	})
}

func (d FuseDeIterator[T]) realRfind(predicate func(T) bool) gust.Option[T] {
	if d.isNone {
		return gust.None[T]()
	}
	return d.iter.Rfind(predicate).InspectNone(func() {
		d.isNone = true
	})
}
