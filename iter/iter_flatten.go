package iter

import (
	"github.com/andeya/gust"
	"github.com/andeya/gust/digit"
	"github.com/andeya/gust/opt"
)

func newFlattenIterator[I gust.Iterable[T], T any](iter Iterator[I]) Iterator[T] {
	p := &flattenIterator[I, T]{iter: iter.ToFuse()}
	p.setFacade(p)
	return p
}

func newDeFlattenIterator[I gust.DeIterable[T], T any](iter DeIterator[I]) DeIterator[T] {
	p := &deFlattenIterator[I, T]{iter: iter.ToDeFuse()}
	p.setFacade(p)
	return p
}

func newFlatMapIterator[T any, B any](iter Iterator[T], f func(T) Iterator[B]) Iterator[B] {
	var m = newMapIterator[T, Iterator[B]](iter, f)
	return newFlattenIterator[Iterator[B], B](m)
}

func newDeFlatMapIterator[T any, B any](iter DeIterator[T], f func(T) DeIterator[B]) DeIterator[B] {
	var m = newDeMapIterator[T, DeIterator[B]](iter, f)
	return newDeFlattenIterator[DeIterator[B], B](m)
}

var (
	_ Iterator[any]       = (*flattenIterator[gust.Iterable[any], any])(nil)
	_ iRealNext[any]      = (*flattenIterator[gust.Iterable[any], any])(nil)
	_ iRealSizeHint       = (*flattenIterator[gust.Iterable[any], any])(nil)
	_ iRealTryFold[any]   = (*flattenIterator[gust.Iterable[any], any])(nil)
	_ iRealFold[any]      = (*flattenIterator[gust.Iterable[any], any])(nil)
	_ iRealAdvanceBy[any] = (*flattenIterator[gust.Iterable[any], any])(nil)
	_ iRealCount          = (*flattenIterator[gust.Iterable[any], any])(nil)
	_ iRealLast[any]      = (*flattenIterator[gust.Iterable[any], any])(nil)
)

// flattenIterator is an iterator that flattens one level of nesting in an iterator of things
// that can be turned into iterators.
type flattenIterator[I gust.Iterable[T], T any] struct {
	iterBackground[T]
	iter      Iterator[I]
	frontiter gust.Option[Iterator[T]]
	backiter  gust.Option[Iterator[T]]
}

func (f *flattenIterator[I, T]) realNext() gust.Option[T] {
	for {
		if f.frontiter.IsSome() {
			x := f.frontiter.Unwrap().Next().InspectNone(func() {
				f.frontiter = gust.None[Iterator[T]]()
			})
			if x.IsSome() {
				return x
			}
		} else {
			x := f.iter.Next().Inspect(func(t I) {
				f.frontiter = gust.Some(FromIterable[T](t))
			})
			if x.IsNone() {
				if f.backiter.IsSome() {
					return f.backiter.Unwrap().Next().InspectNone(func() {
						f.backiter = gust.None[Iterator[T]]()
					})
				}
				return gust.None[T]()
			}
		}
	}
}

func (f *flattenIterator[I, T]) realSizeHint() (uint, gust.Option[uint]) {
	var fl = opt.MapOr(f.frontiter, gust.Pair[uint, gust.Option[uint]]{0, gust.Some[uint](0)}, func(i Iterator[T]) (x gust.Pair[uint, gust.Option[uint]]) {
		x.A, x.B = i.SizeHint()
		return x
	})
	var bl = opt.MapOr(f.backiter, gust.Pair[uint, gust.Option[uint]]{0, gust.Some[uint](0)}, func(i Iterator[T]) (x gust.Pair[uint, gust.Option[uint]]) {
		x.A, x.B = i.SizeHint()
		return x
	})
	var lo = digit.SaturatingAdd(fl.A, bl.A)
	// TODO: check fixed size
	var a, b = f.iter.SizeHint()
	if a == 0 && b.IsSome() && b.Unwrap() == 0 && fl.B.IsSome() && bl.B.IsSome() {
		return lo, digit.CheckedAdd(fl.B.Unwrap(), bl.B.Unwrap())
	}
	return lo, gust.None[uint]()
}

func (f *flattenIterator[I, T]) realTryFold(init any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return f.iterTryFold(init, func(acc any, item Iterator[T]) gust.AnyCtrlFlow {
		return item.TryFold(acc, fold)
	})
}

func (f *flattenIterator[I, T]) realFold(init any, fold func(any, T) any) any {
	return f.iterFold(init, func(acc any, item Iterator[T]) any {
		return item.Fold(acc, fold)
	})
}

func (f *flattenIterator[I, T]) iterTryFold(acc any, fold func(any, Iterator[T]) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	var flatten = func(frontiter *gust.Option[Iterator[T]], fold func(any, Iterator[T]) gust.AnyCtrlFlow) func(any, I) gust.AnyCtrlFlow {
		return func(acc any, iter I) gust.AnyCtrlFlow {
			return fold(acc, *frontiter.Insert(FromIterable[T](iter)))
		}
	}
	if f.frontiter.IsSome() {
		x := fold(acc, f.frontiter.Unwrap())
		if x.IsBreak() {
			return x
		}
		acc = x
		f.frontiter = gust.None[Iterator[T]]()
	}

	acc = f.iter.TryFold(acc, flatten(&f.frontiter, fold))
	f.frontiter = gust.None[Iterator[T]]()

	if f.backiter.IsSome() {
		x := fold(acc, f.backiter.Unwrap())
		if x.IsBreak() {
			return x
		}
		acc = x
		f.backiter = gust.None[Iterator[T]]()
	}
	return gust.AnyContinue(acc)
}

func (f *flattenIterator[I, T]) iterFold(acc any, fold func(any, Iterator[T]) any) any {
	f.frontiter.Inspect(func(i Iterator[T]) {
		acc = fold(acc, i)
	})
	acc = f.iter.Fold(acc, func(acc any, i I) any {
		return fold(acc, FromIterable[T](i))
	})
	f.backiter.Inspect(func(i Iterator[T]) {
		acc = fold(acc, i)
	})
	return acc
}

func (f *flattenIterator[I, T]) realAdvanceBy(n uint) gust.Errable[uint] {
	var advance = func(n any, iter Iterator[T]) gust.AnyCtrlFlow {
		x := iter.AdvanceBy(n.(uint))
		if x.IsErr() {
			return gust.AnyBreak(n.(uint) - x.UnwrapErr())
		}
		return gust.AnyContinue(nil)
	}
	x := f.iterTryFold(n, advance)
	if x.IsBreak() {
		remaining := x.UnwrapBreak().(uint)
		if remaining > 0 {
			return gust.ToErrable(n - remaining)
		}
	}
	return gust.NonErrable[uint]()
}

func (f *flattenIterator[I, T]) realCount() uint {
	var count = func(acc any, iter Iterator[T]) any {
		return acc.(uint) + iter.Count()
	}
	return f.iterFold(uint(0), count).(uint)
}

func (f *flattenIterator[I, T]) realLast() gust.Option[T] {
	var last = func(last any, iter Iterator[T]) any {
		return iter.Last().Or(last.(gust.Option[T]))
	}
	return f.iterFold(gust.None[T](), last).(gust.Option[T])
}

var (
	_ DeIterator[any]         = (*deFlattenIterator[gust.DeIterable[any], any])(nil)
	_ iRealNext[any]          = (*deFlattenIterator[gust.DeIterable[any], any])(nil)
	_ iRealNextBack[any]      = (*deFlattenIterator[gust.DeIterable[any], any])(nil)
	_ iRealSizeHint           = (*deFlattenIterator[gust.DeIterable[any], any])(nil)
	_ iRealTryFold[any]       = (*deFlattenIterator[gust.DeIterable[any], any])(nil)
	_ iRealTryRfold[any]      = (*deFlattenIterator[gust.DeIterable[any], any])(nil)
	_ iRealFold[any]          = (*deFlattenIterator[gust.DeIterable[any], any])(nil)
	_ iRealRfold[any]         = (*deFlattenIterator[gust.DeIterable[any], any])(nil)
	_ iRealAdvanceBy[any]     = (*deFlattenIterator[gust.DeIterable[any], any])(nil)
	_ iRealAdvanceBackBy[any] = (*deFlattenIterator[gust.DeIterable[any], any])(nil)
	_ iRealCount              = (*deFlattenIterator[gust.DeIterable[any], any])(nil)
	_ iRealLast[any]          = (*deFlattenIterator[gust.DeIterable[any], any])(nil)
	_ iRealRemaining          = (*deFlattenIterator[gust.DeIterable[any], any])(nil)
)

type deFlattenIterator[I gust.DeIterable[T], T any] struct {
	deIterBackground[T]
	iter      DeIterator[I]
	frontiter gust.Option[DeIterator[T]]
	backiter  gust.Option[DeIterator[T]]
}

func (f *deFlattenIterator[I, T]) realRemaining() uint {
	return f.iter.Remaining()
}

func (f *deFlattenIterator[I, T]) realNextBack() gust.Option[T] {
	for {
		if f.backiter.IsSome() {
			x := f.backiter.Unwrap().NextBack().InspectNone(func() {
				f.backiter = gust.None[DeIterator[T]]()
			})
			if x.IsSome() {
				return x
			}
		} else {
			x := f.iter.NextBack().Inspect(func(t I) {
				f.backiter = gust.Some(FromDeIterable[T](t))
			})
			if x.IsNone() {
				if f.frontiter.IsSome() {
					return f.frontiter.Unwrap().NextBack().InspectNone(func() {
						f.frontiter = gust.None[DeIterator[T]]()
					})
				}
				return gust.None[T]()
			}
		}
	}
}

func (f *deFlattenIterator[I, T]) realNext() gust.Option[T] {
	for {
		if f.frontiter.IsSome() {
			x := f.frontiter.Unwrap().Next().InspectNone(func() {
				f.frontiter = gust.None[DeIterator[T]]()
			})
			if x.IsSome() {
				return x
			}
		} else {
			x := f.iter.Next().Inspect(func(t I) {
				f.frontiter = gust.Some(FromDeIterable[T](t))
			})
			if x.IsNone() {
				if f.backiter.IsSome() {
					return f.backiter.Unwrap().Next().InspectNone(func() {
						f.backiter = gust.None[DeIterator[T]]()
					})
				}
				return gust.None[T]()
			}
		}
	}
}

func (f *deFlattenIterator[I, T]) realSizeHint() (uint, gust.Option[uint]) {
	var fl = opt.MapOr(f.frontiter, gust.Pair[uint, gust.Option[uint]]{0, gust.Some[uint](0)}, func(i DeIterator[T]) (x gust.Pair[uint, gust.Option[uint]]) {
		x.A, x.B = i.SizeHint()
		return x
	})
	var bl = opt.MapOr(f.backiter, gust.Pair[uint, gust.Option[uint]]{0, gust.Some[uint](0)}, func(i DeIterator[T]) (x gust.Pair[uint, gust.Option[uint]]) {
		x.A, x.B = i.SizeHint()
		return x
	})
	var lo = digit.SaturatingAdd(fl.A, bl.A)
	// TODO: check fixed size
	var a, b = f.iter.SizeHint()
	if a == 0 && b.IsSome() && b.Unwrap() == 0 && fl.B.IsSome() && bl.B.IsSome() {
		return lo, digit.CheckedAdd(fl.B.Unwrap(), bl.B.Unwrap())
	}
	return lo, gust.None[uint]()
}

func (f *deFlattenIterator[I, T]) realTryFold(init any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return f.iterTryFold(init, func(acc any, item DeIterator[T]) gust.AnyCtrlFlow {
		return item.TryFold(acc, fold)
	})
}

func (f *deFlattenIterator[I, T]) realTryRfold(init any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return f.iterTryRfold(init, func(acc any, item DeIterator[T]) gust.AnyCtrlFlow {
		return item.TryRfold(acc, fold)
	})
}

func (f *deFlattenIterator[I, T]) realFold(init any, fold func(any, T) any) any {
	return f.iterFold(init, func(acc any, item DeIterator[T]) any {
		return item.Fold(acc, fold)
	})
}

func (f *deFlattenIterator[I, T]) realRfold(init any, fold func(any, T) any) any {
	return f.iterRfold(init, func(acc any, item DeIterator[T]) any {
		return item.Rfold(acc, fold)
	})
}

func (f *deFlattenIterator[I, T]) realAdvanceBy(n uint) gust.Errable[uint] {
	var advance = func(n any, iter DeIterator[T]) gust.AnyCtrlFlow {
		x := iter.AdvanceBy(n.(uint))
		if x.IsErr() {
			return gust.AnyBreak(n.(uint) - x.UnwrapErr())
		}
		return gust.AnyContinue(nil)
	}
	x := f.iterTryFold(n, advance)
	if x.IsBreak() {
		remaining := x.UnwrapBreak().(uint)
		if remaining > 0 {
			return gust.ToErrable(n - remaining)
		}
	}
	return gust.NonErrable[uint]()
}

func (f *deFlattenIterator[I, T]) realAdvanceBackBy(n uint) gust.Errable[uint] {
	var advance = func(n any, iter DeIterator[T]) gust.AnyCtrlFlow {
		x := iter.AdvanceBackBy(n.(uint))
		if x.IsErr() {
			return gust.AnyBreak(n.(uint) - x.UnwrapErr())
		}
		return gust.AnyContinue(nil)
	}
	x := f.iterTryRfold(n, advance)
	if x.IsBreak() {
		remaining := x.UnwrapBreak().(uint)
		if remaining > 0 {
			return gust.ToErrable(n - remaining)
		}
	}
	return gust.NonErrable[uint]()
}

func (f *deFlattenIterator[I, T]) realCount() uint {
	var count = func(acc any, iter DeIterator[T]) any {
		return acc.(uint) + iter.Count()
	}
	return f.iterFold(uint(0), count).(uint)
}

func (f *deFlattenIterator[I, T]) realLast() gust.Option[T] {
	var last = func(last any, iter DeIterator[T]) any {
		return iter.Last().Or(last.(gust.Option[T]))
	}
	return f.iterFold(gust.None[T](), last).(gust.Option[T])
}

func (f *deFlattenIterator[I, T]) iterTryFold(acc any, fold func(any, DeIterator[T]) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	var flatten = func(frontiter *gust.Option[DeIterator[T]], fold func(any, DeIterator[T]) gust.AnyCtrlFlow) func(any, I) gust.AnyCtrlFlow {
		return func(acc any, iter I) gust.AnyCtrlFlow {
			return fold(acc, *frontiter.Insert(FromDeIterable[T](iter)))
		}
	}
	if f.frontiter.IsSome() {
		x := fold(acc, f.frontiter.Unwrap())
		if x.IsBreak() {
			return x
		}
		acc = x
		f.frontiter = gust.None[DeIterator[T]]()
	}

	acc = f.iter.TryFold(acc, flatten(&f.frontiter, fold))
	f.frontiter = gust.None[DeIterator[T]]()

	if f.backiter.IsSome() {
		x := fold(acc, f.backiter.Unwrap())
		if x.IsBreak() {
			return x
		}
		acc = x
		f.backiter = gust.None[DeIterator[T]]()
	}
	return gust.AnyContinue(acc)
}

func (f *deFlattenIterator[I, T]) iterFold(acc any, fold func(any, DeIterator[T]) any) any {
	f.frontiter.Inspect(func(i DeIterator[T]) {
		acc = fold(acc, i)
	})
	acc = f.iter.Fold(acc, func(acc any, i I) any {
		return fold(acc, FromDeIterable[T](i))
	})
	f.backiter.Inspect(func(i DeIterator[T]) {
		acc = fold(acc, i)
	})
	return acc
}

func (f *deFlattenIterator[I, T]) iterTryRfold(acc any, fold func(any, DeIterator[T]) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	var flatten = func(backiter *gust.Option[DeIterator[T]], fold func(any, DeIterator[T]) gust.AnyCtrlFlow) func(any, I) gust.AnyCtrlFlow {
		return func(acc any, iter I) gust.AnyCtrlFlow {
			return fold(acc, *backiter.Insert(FromDeIterable[T](iter)))
		}
	}
	if f.backiter.IsSome() {
		x := fold(acc, f.backiter.Unwrap())
		if x.IsBreak() {
			return x
		}
		acc = x
		f.backiter = gust.None[DeIterator[T]]()
	}

	acc = f.iter.TryRfold(acc, flatten(&f.backiter, fold))
	f.backiter = gust.None[DeIterator[T]]()

	if f.frontiter.IsSome() {
		x := fold(acc, f.frontiter.Unwrap())
		if x.IsBreak() {
			return x
		}
		acc = x
		f.frontiter = gust.None[DeIterator[T]]()
	}
	return gust.AnyContinue(acc)
}

func (f *deFlattenIterator[I, T]) iterRfold(acc any, fold func(any, DeIterator[T]) any) any {
	f.backiter.Inspect(func(i DeIterator[T]) {
		acc = fold(acc, i)
	})
	acc = f.iter.Rfold(acc, func(acc any, i I) any {
		return fold(acc, FromDeIterable[T](i))
	})
	f.frontiter.Inspect(func(i DeIterator[T]) {
		acc = fold(acc, i)
	})
	return acc
}
