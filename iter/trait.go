package iter

import (
	"github.com/andeya/gust"
)

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

func (iter iterTrait[T]) SizeHint() (uint64, gust.Option[uint64]) {
	if cover, ok := iter.facade.(iRealSizeHint); ok {
		return cover.realSizeHint()
	}
	return 0, gust.None[uint64]()
}

func (iter iterTrait[T]) Count() uint64 {
	if cover, ok := iter.facade.(iRealCount); ok {
		return cover.realCount()
	}
	return iter.Fold(
		uint64(0),
		func(count any, _ T) any {
			return count.(uint64) + 1
		},
	).(uint64)
}

func (iter iterTrait[T]) Fold(init any, f func(any, T) any) any {
	if cover, ok := iter.facade.(iRealFold[T]); ok {
		return cover.realFold(init, f)
	}
	var accum = init
	for {
		x := iter.Next()
		if x.IsNone() {
			break
		}
		accum = f(accum, x.Unwrap())
	}
	return accum
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
	return iter.Fold(gust.None[T](), func(_ any, x T) any { return gust.Some(x) }).(gust.Option[T])
}

func (iter iterTrait[T]) AdvanceBy(n uint) gust.Result[struct{}] {
	if cover, ok := iter.facade.(iRealAdvanceBy[T]); ok {
		return cover.realAdvanceBy(n)
	}
	for i := uint(0); i < n; i++ {
		if iter.Next().IsNone() {
			return gust.Err[struct{}](i)
		}
	}
	return gust.Ok(struct{}{})
}

func (iter iterTrait[T]) Nth(n uint) gust.Option[T] {
	if cover, ok := iter.facade.(iRealNth[T]); ok {
		return cover.realNth(n)
	}
	var res = iter.AdvanceBy(n)
	if res.IsErr() {
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
	return gust.Some(iter.Fold(first, func(accum any, item T) any {
		return f(accum.(T), item)
	}).(T))
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
	var check = func(f func(T) gust.Option[any]) func(any, T) gust.Result[any] {
		return func(_ any, x T) gust.Result[any] {
			r := f(x)
			if r.IsSome() {
				return gust.Err[any](x)
			} else {
				return gust.Ok[any](nil)
			}
		}
	}
	r := iter.TryFold(nil, check(f))
	if r.IsErr() {
		return gust.Some(r.ErrVal())
	}
	return gust.None[any]()
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
