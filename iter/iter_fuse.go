package iter

import (
	"github.com/andeya/gust"
)

var (
	_ Iterator[any]     = (*fuseIterator[any])(nil)
	_ iRealNext[any]    = (*fuseIterator[any])(nil)
	_ iRealNth[any]     = (*fuseIterator[any])(nil)
	_ iRealLast[any]    = (*fuseIterator[any])(nil)
	_ iRealCount        = (*fuseIterator[any])(nil)
	_ iRealSizeHint     = (*fuseIterator[any])(nil)
	_ iRealTryFold[any] = (*fuseIterator[any])(nil)
	_ iRealFold[any]    = (*fuseIterator[any])(nil)
	_ iRealFind[any]    = (*fuseIterator[any])(nil)
)

func newFuseIterator[T any](iter Iterator[T]) Iterator[T] {
	p := &fuseIterator[T]{iter: iter}
	p.setFacade(p)
	return p
}

type fuseIterator[T any] struct {
	deIterBackground[T]
	iter   Iterator[T]
	isNone bool
}

func (f *fuseIterator[T]) realNextBack() gust.Option[T] {
	panic("unreachable")
}

func (f *fuseIterator[T]) realNext() gust.Option[T] {
	if f.isNone {
		return gust.None[T]()
	}
	return f.iter.Next().InspectNone(func() {
		f.isNone = true
	})
}

func (f *fuseIterator[T]) realNth(n uint) gust.Option[T] {
	if f.isNone {
		return gust.None[T]()
	}
	return f.iter.Nth(n).InspectNone(func() {
		f.isNone = true
	})
}

func (f *fuseIterator[T]) realLast() gust.Option[T] {
	if f.isNone {
		return gust.None[T]()
	}
	return f.iter.Last()
}

func (f *fuseIterator[T]) realCount() uint {
	if f.isNone {
		return 0
	}
	f.isNone = true
	return f.iter.Count()
}

func (f *fuseIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	if f.isNone {
		return 0, gust.Some[uint](0)
	}
	return f.iter.SizeHint()
}

func (f *fuseIterator[T]) realTryFold(acc any, fold func(any, T) gust.Result[any]) gust.Result[any] {
	if f.isNone {
		return gust.Err[any]("fuse iterator is empty")
	}
	r := f.iter.TryFold(acc, fold)
	f.isNone = true
	return r
}

func (f *fuseIterator[T]) realFold(acc any, fold func(any, T) any) any {
	if f.isNone {
		return acc
	}
	r := f.iter.Fold(acc, fold)
	f.isNone = true
	return r
}

func (f *fuseIterator[T]) realFind(predicate func(T) bool) gust.Option[T] {
	if f.isNone {
		return gust.None[T]()
	}
	return f.iter.Find(predicate).InspectNone(func() {
		f.isNone = true
	})
}

var (
	_ DeIterator[any]    = (*fuseDeIterator[any])(nil)
	_ iRealNext[any]     = (*fuseDeIterator[any])(nil)
	_ iRealNth[any]      = (*fuseDeIterator[any])(nil)
	_ iRealLast[any]     = (*fuseDeIterator[any])(nil)
	_ iRealCount         = (*fuseDeIterator[any])(nil)
	_ iRealSizeHint      = (*fuseDeIterator[any])(nil)
	_ iRealTryFold[any]  = (*fuseDeIterator[any])(nil)
	_ iRealFold[any]     = (*fuseDeIterator[any])(nil)
	_ iRealFind[any]     = (*fuseDeIterator[any])(nil)
	_ iRealRemaining     = (*fuseDeIterator[any])(nil)
	_ iRealNextBack[any] = (*fuseDeIterator[any])(nil)
	_ iRealNthBack[any]  = (*fuseDeIterator[any])(nil)
	_ iRealTryRfold[any] = (*fuseDeIterator[any])(nil)
	_ iRealRfold[any]    = (*fuseDeIterator[any])(nil)
	_ iRealRfind[any]    = (*fuseDeIterator[any])(nil)
)

func newFuseDeIterator[T any](iter DeIterator[T]) DeIterator[T] {
	p := &fuseDeIterator[T]{}
	p.iter = iter
	p.setFacade(p)
	return p
}

// fuseDeIterator double ended fuse iterator
type fuseDeIterator[T any] struct {
	fuseIterator[T]
}

func (d *fuseDeIterator[T]) realRemaining() uint {
	return d.iter.(DeIterator[T]).Remaining()
}

func (d *fuseDeIterator[T]) realNextBack() gust.Option[T] {
	if d.isNone {
		return gust.None[T]()
	}
	deIter := d.iter.(DeIterator[T])
	return deIter.NextBack().InspectNone(func() {
		d.isNone = true
	})
}

func (d *fuseDeIterator[T]) realNthBack(n uint) gust.Option[T] {
	if d.isNone {
		return gust.None[T]()
	}
	deIter := d.iter.(DeIterator[T])
	return deIter.NthBack(n).InspectNone(func() {
		d.isNone = true
	})
}

func (d *fuseDeIterator[T]) realTryRfold(acc any, fold func(any, T) gust.Result[any]) gust.Result[any] {
	if d.isNone {
		return gust.Err[any]("fuse iterator is empty")
	}
	deIter := d.iter.(DeIterator[T])
	r := deIter.TryRfold(acc, fold)
	d.isNone = true
	return r
}

func (d *fuseDeIterator[T]) realRfold(acc any, fold func(any, T) any) any {
	if d.isNone {
		return acc
	}
	deIter := d.iter.(DeIterator[T])
	r := deIter.Rfold(acc, fold)
	d.isNone = true
	return r
}

func (d *fuseDeIterator[T]) realRfind(predicate func(T) bool) gust.Option[T] {
	if d.isNone {
		return gust.None[T]()
	}
	deIter := d.iter.(DeIterator[T])
	return deIter.Rfind(predicate).InspectNone(func() {
		d.isNone = true
	})
}
