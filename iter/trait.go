package iter

import (
	"github.com/andeya/gust"
	"github.com/andeya/gust/opt"
)

var _ Iterator[any] = iterTrait[any]{}

type iterTrait[T any] struct {
	facade iRealNext[T]
}

//goland:noinspection GoMixedReceiverTypes
func (iter *iterTrait[T]) setFacade(facade iRealNext[T]) {
	iter.facade = facade
}

func (iter iterTrait[T]) Next() gust.Option[T] {
	return iter.facade.realNext()
}

func (iter iterTrait[T]) NextChunk(n uint) ([]T, bool) {
	var chunk = make([]T, 0, n)
	for i := uint(0); i < n; i++ {
		item := iter.Next()
		if item.IsSome() {
			chunk = append(chunk, item.Unwrap())
		} else {
			return chunk, false
		}
	}
	return chunk, true
}

func (iter iterTrait[T]) SizeHint() (uint, gust.Option[uint]) {
	if cover, ok := iter.facade.(iRealSizeHint); ok {
		return cover.realSizeHint()
	}
	return 0, gust.None[uint]()
}

func (iter iterTrait[T]) Count() uint {
	if cover, ok := iter.facade.(iRealCount); ok {
		return cover.realCount()
	}
	return Fold[T, uint](iter,
		uint(0),
		func(count uint, _ T) uint {
			return count + 1
		},
	)
}

func (iter iterTrait[T]) Fold(init any, f func(any, T) any) any {
	if cover, ok := iter.facade.(iRealFold[T]); ok {
		return cover.realFold(init, f)
	}
	return Fold[T, any](iter, init, f)
}

func (iter iterTrait[T]) TryFold(init any, f func(any, T) gust.Result[any]) gust.Result[any] {
	if cover, ok := iter.facade.(iRealTryFold[T]); ok {
		return cover.realTryFold(init, f)
	}
	return TryFold[T, any](iter, init, f)
}

func (iter iterTrait[T]) Last() gust.Option[T] {
	if cover, ok := iter.facade.(iRealLast[T]); ok {
		return cover.realLast()
	}
	return Fold[T, gust.Option[T]](
		iter,
		gust.None[T](),
		func(_ gust.Option[T], x T) gust.Option[T] {
			return gust.Some(x)
		})
}

func (iter iterTrait[T]) AdvanceBy(n uint) gust.Errable[uint] {
	if cover, ok := iter.facade.(iRealAdvanceBy[T]); ok {
		return cover.realAdvanceBy(n)
	}
	for i := uint(0); i < n; i++ {
		if iter.Next().IsNone() {
			return gust.ToErrable[uint](i)
		}
	}
	return gust.NonErrable[uint]()
}

func (iter iterTrait[T]) Nth(n uint) gust.Option[T] {
	if cover, ok := iter.facade.(iRealNth[T]); ok {
		return cover.realNth(n)
	}
	var res = iter.AdvanceBy(n)
	if res.AsError() {
		return gust.None[T]()
	}
	return iter.Next()
}

func (iter iterTrait[T]) ForEach(f func(T)) {
	if cover, ok := iter.facade.(iRealForEach[T]); ok {
		cover.realForEach(f)
		return
	}
	var call = func(f func(T)) func(any, T) any {
		return func(_ any, item T) any {
			f(item)
			return nil
		}
	}
	_ = iter.Fold(nil, call(f))
}

func (iter iterTrait[T]) Reduce(f func(accum T, item T) T) gust.Option[T] {
	if cover, ok := iter.facade.(iRealReduce[T]); ok {
		return cover.realReduce(f)
	}
	var first = iter.Next()
	if first.IsNone() {
		return first
	}
	return gust.Some(Fold[T, T](iter, first.Unwrap(), func(accum T, item T) T {
		return f(accum, item)
	}))
}

func (iter iterTrait[T]) All(predicate func(T) bool) bool {
	if cover, ok := iter.facade.(iRealAll[T]); ok {
		return cover.realAll(predicate)
	}
	var check = func(f func(T) bool) func(any, T) gust.Result[any] {
		return func(_ any, x T) gust.Result[any] {
			if f(x) {
				return gust.Ok[any](nil)
			} else {
				return gust.Err[any](nil)
			}
		}
	}
	return iter.TryFold(nil, check(predicate)).IsOk()
}

func (iter iterTrait[T]) Any(predicate func(T) bool) bool {
	if cover, ok := iter.facade.(iRealAny[T]); ok {
		return cover.realAny(predicate)
	}
	var check = func(f func(T) bool) func(any, T) gust.Result[any] {
		return func(_ any, x T) gust.Result[any] {
			if f(x) {
				return gust.Err[any](nil)
			} else {
				return gust.Ok[any](nil)
			}
		}
	}
	return iter.TryFold(nil, check(predicate)).IsErr()
}

func (iter iterTrait[T]) Find(predicate func(T) bool) gust.Option[T] {
	if cover, ok := iter.facade.(iRealFind[T]); ok {
		return cover.realFind(predicate)
	}
	var check = func(f func(T) bool) func(any, T) gust.Result[any] {
		return func(_ any, x T) gust.Result[any] {
			if f(x) {
				return gust.Err[any](x)
			} else {
				return gust.Ok[any](nil)
			}
		}
	}
	r := iter.TryFold(nil, check(predicate))
	if r.IsErr() {
		return gust.Some[T](r.ErrVal().(T))
	}
	return gust.None[T]()
}

func (iter iterTrait[T]) FindMap(f func(T) gust.Option[any]) gust.Option[any] {
	if cover, ok := iter.facade.(iRealFindMap[T]); ok {
		return cover.realFindMap(f)
	}
	return FindMap[T, any](iter, f)
}

func (iter iterTrait[T]) TryFind(predicate func(T) gust.Result[bool]) gust.Result[gust.Option[T]] {
	if cover, ok := iter.facade.(iRealTryFind[T]); ok {
		return cover.realTryFind(predicate)
	}
	var check = func(f func(T) gust.Result[bool]) func(any, T) gust.Result[any] {
		return func(_ any, x T) gust.Result[any] {
			r := f(x)
			if r.IsOk() {
				if r.Unwrap() {
					return gust.Err[any](gust.Ok[gust.Option[T]](gust.Some(x)))
				} else {
					return gust.Ok[any](nil)
				}
			} else {
				return gust.Err[any](gust.Err[gust.Option[T]](r.Err()))
			}
		}
	}
	r := iter.TryFold(nil, check(predicate))
	if r.IsErr() {
		return r.ErrVal().(gust.Result[gust.Option[T]])
	}
	return gust.Ok[gust.Option[T]](gust.None[T]())
}

func (iter iterTrait[T]) Position(predicate func(T) bool) gust.Option[int] {
	if cover, ok := iter.facade.(iRealPosition[T]); ok {
		return cover.realPosition(predicate)
	}
	var check = func(f func(T) bool) func(int, T) gust.Result[int] {
		return func(i int, x T) gust.Result[int] {
			if f(x) {
				return gust.Err[int](i)
			} else {
				return gust.Ok[int](i + 1)
			}
		}
	}
	r := TryFold[T, int](iter, 0, check(predicate))
	if r.IsErr() {
		return gust.Some[int](r.ErrVal().(int))
	}
	return gust.None[int]()
}

func (iter iterTrait[T]) StepBy(step uint) *StepByIterator[T] {
	if cover, ok := iter.facade.(iRealStepBy[T]); ok {
		return cover.realStepBy(step)
	}
	return newStepByIterator[T](iter, step)
}

func (iter iterTrait[T]) Filter(f func(T) bool) *FilterIterator[T] {
	if cover, ok := iter.facade.(iRealFilter[T]); ok {
		return cover.realFilter(f)
	}
	return newFilterIterator[T](iter, f)
}

func (iter iterTrait[T]) FilterMap(f func(T) gust.Option[T]) *FilterMapIterator[T] {
	if cover, ok := iter.facade.(iRealFilterMap[T]); ok {
		return cover.realFilterMap(f)
	}
	return newFilterMapIterator[T](iter, f)
}

func (iter iterTrait[T]) Chain(other Iterator[T]) *ChainIterator[T] {
	if cover, ok := iter.facade.(iRealChain[T]); ok {
		return cover.realChain(other)
	}
	return newChainIterator[T](iter, other)
}

func (iter iterTrait[T]) Map(f func(T) any) *MapIterator[T, any] {
	if cover, ok := iter.facade.(iRealMap[T]); ok {
		return cover.realMap(f)
	}
	return newMapIterator[T, any](iter, f)
}

func (iter iterTrait[T]) Inspect(f func(T)) *InspectIterator[T] {
	return newInspectIterator[T](iter, f)
}

func (iter iterTrait[T]) Fuse() *FuseIterator[T] {
	return newFuseIterator[T](iter)
}

func (iter iterTrait[T]) Collect() []T {
	lower, _ := iter.SizeHint()
	return Fold[T, []T](iter, make([]T, 0, lower), func(slice []T, x T) []T {
		return append(slice, x)
	})
}

var _ SizeDeIterator[any] = sizeDeIterTrait[any]{}

type sizeDeIterTrait[T any] struct {
	iterTrait[T]
	facade iRealDeIterable[T]
}

//goland:noinspection GoMixedReceiverTypes
func (d *sizeDeIterTrait[T]) setFacade(facade iRealDeIterable[T]) {
	d.iterTrait.facade = facade
	d.facade = facade
}

func (d sizeDeIterTrait[T]) Remaining() uint {
	if size, ok := d.facade.(iRealSizeDeIterable[T]); ok {
		return size.realRemaining()
	}
	lo, hi := d.SizeHint()
	if opt.MapOr[uint, bool](hi, false, func(x uint) bool {
		return x == lo
	}) {
		return lo
	}
	return lo
}

func (d sizeDeIterTrait[T]) NextBack() gust.Option[T] {
	return d.facade.realNextBack()
}

func (d sizeDeIterTrait[T]) AdvanceBackBy(n uint) gust.Errable[uint] {
	if cover, ok := d.facade.(iRealAdvanceBackBy[T]); ok {
		return cover.realAdvanceBackBy(n)
	}
	for i := uint(0); i < n; i++ {
		if d.NextBack().IsNone() {
			return gust.ToErrable[uint](i)
		}
	}
	return gust.NonErrable[uint]()
}

func (d sizeDeIterTrait[T]) NthBack(n uint) gust.Option[T] {
	if cover, ok := d.facade.(iRealNthBack[T]); ok {
		return cover.realNthBack(n)
	}
	if d.AdvanceBackBy(n).AsError() {
		return gust.None[T]()
	}
	return d.NextBack()
}

func (d sizeDeIterTrait[T]) TryRfold(init any, fold func(any, T) gust.Result[any]) gust.Result[any] {
	if cover, ok := d.facade.(iRealTryRfold[T]); ok {
		return cover.realTryRfold(init, fold)
	}
	return TryRfold[T](d, init, fold)
}

func (d sizeDeIterTrait[T]) Rfold(init any, fold func(any, T) any) any {
	if cover, ok := d.facade.(iRealRfold[T]); ok {
		return cover.realRfold(init, fold)
	}
	return Rfold[T](d, init, fold)
}

func (d sizeDeIterTrait[T]) Rfind(predicate func(T) bool) gust.Option[T] {
	if cover, ok := d.facade.(iRealRfind[T]); ok {
		return cover.realRfind(predicate)
	}
	var check = func(f func(T) bool) func(any, T) gust.Result[any] {
		return func(_ any, x T) gust.Result[any] {
			if f(x) {
				return gust.Err[any](x)
			} else {
				return gust.Ok[any](nil)
			}
		}
	}
	r := d.TryRfold(nil, check(predicate))
	if r.IsErr() {
		return gust.Some[T](r.ErrVal().(T))
	}
	return gust.None[T]()
}

// DeFuse creates an iterator which ends after the first [`gust.None[T]()`].
//
// After an iterator returns [`gust.None[T]()`], future calls may or may not yield
// [`gust.Some(T)`] again. `Fuse()` adapts an iterator, ensuring that after a
// [`gust.None[T]()`] is given, it will always return [`gust.None[T]()`] forever.
func (d sizeDeIterTrait[T]) DeFuse() *FuseDeIterator[T] {
	return newFuseDeIterator[T](d)
}
