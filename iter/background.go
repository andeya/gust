package iter

import (
	"github.com/andeya/gust"
	"github.com/andeya/gust/opt"
)

var _ Iterator[any] = iterBackground[any]{}

type iterBackground[T any] struct {
	facade iRealNext[T]
}

//goland:noinspection GoMixedReceiverTypes
func (iter *iterBackground[T]) setFacade(facade iRealNext[T]) {
	iter.facade = facade
}

func (iter iterBackground[T]) Collect() []T {
	lower, _ := iter.SizeHint()
	return Fold[T, []T](iter, make([]T, 0, lower), func(slice []T, x T) []T {
		return append(slice, x)
	})
}

func (iter iterBackground[T]) Next() gust.Option[T] {
	return iter.facade.realNext()
}

func (iter iterBackground[T]) NextChunk(n uint) gust.EnumResult[[]T, []T] {
	var chunk = make([]T, 0, n)
	for i := uint(0); i < n; i++ {
		item := iter.Next()
		if item.IsSome() {
			chunk = append(chunk, item.Unwrap())
		} else {
			return gust.EnumErr[[]T, []T](chunk)
		}
	}
	return gust.EnumOk[[]T, []T](chunk)
}

func (iter iterBackground[T]) SizeHint() (uint, gust.Option[uint]) {
	if cover, ok := iter.facade.(iRealSizeHint); ok {
		return cover.realSizeHint()
	}
	return 0, gust.None[uint]()
}

func (iter iterBackground[T]) Count() uint {
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

func (iter iterBackground[T]) Fold(init any, f func(any, T) any) any {
	if cover, ok := iter.facade.(iRealFold[T]); ok {
		return cover.realFold(init, f)
	}
	return Fold[T, any](iter, init, f)
}

func (iter iterBackground[T]) TryFold(init any, f func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	if cover, ok := iter.facade.(iRealTryFold[T]); ok {
		return cover.realTryFold(init, f)
	}
	return TryFold[T, any](iter, init, f)
}

func (iter iterBackground[T]) Last() gust.Option[T] {
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

func (iter iterBackground[T]) AdvanceBy(n uint) gust.Errable[uint] {
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

func (iter iterBackground[T]) Nth(n uint) gust.Option[T] {
	if cover, ok := iter.facade.(iRealNth[T]); ok {
		return cover.realNth(n)
	}
	var res = iter.AdvanceBy(n)
	if res.IsErr() {
		return gust.None[T]()
	}
	return iter.Next()
}

func (iter iterBackground[T]) ForEach(f func(T)) {
	var call = func(f func(T)) func(any, T) any {
		return func(_ any, item T) any {
			f(item)
			return nil
		}
	}
	_ = iter.Fold(nil, call(f))
}

func (iter iterBackground[T]) TryForEach(f func(T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return iter.TryFold(nil, func(_ any, x T) gust.AnyCtrlFlow {
		return f(x)
	})
}

func (iter iterBackground[T]) Reduce(f func(accum T, item T) T) gust.Option[T] {
	var first = iter.Next()
	if first.IsNone() {
		return first
	}
	return gust.Some(Fold[T, T](iter, first.Unwrap(), func(accum T, item T) T {
		return f(accum, item)
	}))
}

func (iter iterBackground[T]) All(predicate func(T) bool) bool {
	var check = func(f func(T) bool) func(any, T) gust.AnyCtrlFlow {
		return func(_ any, x T) gust.AnyCtrlFlow {
			if f(x) {
				return gust.AnyContinue(nil)
			} else {
				return gust.AnyBreak(nil)
			}
		}
	}
	return iter.TryFold(nil, check(predicate)).IsContinue()
}

func (iter iterBackground[T]) Any(predicate func(T) bool) bool {
	var check = func(f func(T) bool) func(any, T) gust.AnyCtrlFlow {
		return func(_ any, x T) gust.AnyCtrlFlow {
			if f(x) {
				return gust.AnyBreak(nil)
			} else {
				return gust.AnyContinue(nil)
			}
		}
	}
	return iter.TryFold(nil, check(predicate)).IsBreak()
}

func (iter iterBackground[T]) Find(predicate func(T) bool) gust.Option[T] {
	if cover, ok := iter.facade.(iRealFind[T]); ok {
		return cover.realFind(predicate)
	}
	var check = func(f func(T) bool) func(any, T) gust.AnyCtrlFlow {
		return func(_ any, x T) gust.AnyCtrlFlow {
			if f(x) {
				return gust.AnyBreak(x)
			} else {
				return gust.AnyContinue(nil)
			}
		}
	}
	r := iter.TryFold(nil, check(predicate))
	if r.IsBreak() {
		return gust.Some[T](r.UnwrapBreak().(T))
	}
	return gust.None[T]()
}

func (iter iterBackground[T]) FindMap(f func(T) gust.Option[T]) gust.Option[T] {
	return FindMap[T, T](iter, f)
}

func (iter iterBackground[T]) XFindMap(f func(T) gust.Option[any]) gust.Option[any] {
	return FindMap[T, any](iter, f)
}

func (iter iterBackground[T]) Partition(f func(T) bool) (truePart []T, falsePart []T) {
	var left []T
	var right []T
	iter.Fold(nil, func(_ any, x T) any {
		if f(x) {
			left = append(left, x)
		} else {
			right = append(right, x)
		}
		return nil
	})
	return left, right
}

func (iter iterBackground[T]) IsPartitioned(predicate func(T) bool) bool {
	// Either all items test `true`, or the first clause stops at `false`
	// and we check that there are no more `true` items after that.
	return iter.All(predicate) || !iter.Any(predicate)
}

func (iter iterBackground[T]) TryFind(predicate func(T) gust.Result[bool]) gust.Result[gust.Option[T]] {
	var check = func(f func(T) gust.Result[bool]) func(any, T) gust.AnyCtrlFlow {
		return func(_ any, x T) gust.AnyCtrlFlow {
			r := f(x)
			if r.IsOk() {
				if r.Unwrap() {
					return gust.AnyBreak(gust.Ok[gust.Option[T]](gust.Some(x)))
				} else {
					return gust.AnyContinue(nil)
				}
			} else {
				return gust.AnyBreak(gust.Err[gust.Option[T]](r.Err()))
			}
		}
	}
	r := iter.TryFold(nil, check(predicate))
	if r.IsBreak() {
		return r.UnwrapBreak().(gust.Result[gust.Option[T]])
	}
	return gust.Ok[gust.Option[T]](gust.None[T]())
}

func (iter iterBackground[T]) Position(predicate func(T) bool) gust.Option[int] {
	var check = func(f func(T) bool) func(int, T) gust.SigCtrlFlow[int] {
		return func(i int, x T) gust.SigCtrlFlow[int] {
			if f(x) {
				return gust.SigBreak[int](i)
			} else {
				return gust.SigContinue[int](i + 1)
			}
		}
	}
	r := TryFold[T, int](iter, 0, check(predicate))
	if r.IsBreak() {
		return gust.Some[int](r.UnwrapBreak())
	}
	return gust.None[int]()
}

func (iter iterBackground[T]) ToStepBy(step uint) Iterator[T] {
	return newStepByIterator[T](iter, step)
}

func (iter iterBackground[T]) ToFilter(f func(T) bool) Iterator[T] {
	return newFilterIterator[T](iter, f)
}

func (iter iterBackground[T]) ToFilterMap(f func(T) gust.Option[T]) Iterator[T] {
	return newFilterMapIterator[T, T](iter, f)
}

func (iter iterBackground[T]) ToXFilterMap(f func(T) gust.Option[any]) Iterator[any] {
	return newFilterMapIterator[T, any](iter, f)
}

func (iter iterBackground[T]) ToChain(other Iterator[T]) Iterator[T] {
	return newChainIterator[T](iter, other)
}

func (iter iterBackground[T]) ToMap(f func(T) T) Iterator[T] {
	return newMapIterator[T, T](iter, f)
}

func (iter iterBackground[T]) ToXMap(f func(T) any) Iterator[any] {
	return newMapIterator[T, any](iter, f)
}

func (iter iterBackground[T]) ToInspect(f func(T)) Iterator[T] {
	return newInspectIterator[T](iter, f)
}

func (iter iterBackground[T]) ToFuse() Iterator[T] {
	return newFuseIterator[T](iter)
}

func (iter iterBackground[T]) ToPeekable() PeekableIterator[T] {
	return newPeekableIterator[T](iter)
}

func (iter iterBackground[T]) ToIntersperse(separator T) Iterator[T] {
	return newIntersperseIterator[T](iter.ToPeekable(), separator)
}

func (iter iterBackground[T]) ToIntersperseWith(separator func() T) Iterator[T] {
	return newIntersperseWithIterator[T](iter.ToPeekable(), separator)
}

func (iter iterBackground[T]) ToSkipWhile(predicate func(T) bool) Iterator[T] {
	return newSkipWhileIterator[T](iter, predicate)
}

func (iter iterBackground[T]) ToTakeWhile(predicate func(T) bool) Iterator[T] {
	return newTakeWhileIterator[T](iter, predicate)
}

func (iter iterBackground[T]) ToMapWhile(predicate func(T) gust.Option[T]) Iterator[T] {
	return newMapWhileIterator[T, T](iter, predicate)
}

func (iter iterBackground[T]) ToXMapWhile(predicate func(T) gust.Option[any]) Iterator[any] {
	return newMapWhileIterator[T, any](iter, predicate)
}

func (iter iterBackground[T]) ToSkip(n uint) Iterator[T] {
	return newSkipIterator[T](iter, n)
}

func (iter iterBackground[T]) ToTake(n uint) Iterator[T] {
	return newTakeIterator[T](iter, n)
}

func (iter iterBackground[T]) ToScan(initialState any, f func(state *any, item T) gust.Option[any]) Iterator[any] {
	return newScanIterator[T, any, any](iter, initialState, f)
}

var _ DeIterator[any] = deIterBackground[any]{}

type deIterBackground[T any] struct {
	iterBackground[T]
}

//goland:noinspection GoMixedReceiverTypes
func (iter *deIterBackground[T]) setFacade(facade iRealDeIterable[T]) {
	iter.iterBackground.facade = facade
}

func (iter deIterBackground[T]) Remaining() uint {
	if size, ok := iter.facade.(iRealRemaining); ok {
		return size.realRemaining()
	}
	return defaultRemaining[T](iter)
}

func defaultRemaining[T any](iter Iterator[T]) uint {
	lo, hi := iter.SizeHint()
	if opt.MapOr[uint, bool](hi, false, func(x uint) bool {
		return x == lo
	}) {
		return lo
	}
	return lo
}

func (iter deIterBackground[T]) NextBack() gust.Option[T] {
	return iter.facade.(iRealDeIterable[T]).realNextBack()
}

func (iter deIterBackground[T]) AdvanceBackBy(n uint) gust.Errable[uint] {
	if cover, ok := iter.facade.(iRealAdvanceBackBy[T]); ok {
		return cover.realAdvanceBackBy(n)
	}
	for i := uint(0); i < n; i++ {
		if iter.NextBack().IsNone() {
			return gust.ToErrable[uint](i)
		}
	}
	return gust.NonErrable[uint]()
}

func (iter deIterBackground[T]) NthBack(n uint) gust.Option[T] {
	if cover, ok := iter.facade.(iRealNthBack[T]); ok {
		return cover.realNthBack(n)
	}
	if iter.AdvanceBackBy(n).IsErr() {
		return gust.None[T]()
	}
	return iter.NextBack()
}

func (iter deIterBackground[T]) TryRfold(init any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	if cover, ok := iter.facade.(iRealTryRfold[T]); ok {
		return cover.realTryRfold(init, fold)
	}
	return TryRfold[T, any](iter, init, fold)
}

func (iter deIterBackground[T]) Rfold(init any, fold func(any, T) any) any {
	if cover, ok := iter.facade.(iRealRfold[T]); ok {
		return cover.realRfold(init, fold)
	}
	return Rfold[T](iter, init, fold)
}

func (iter deIterBackground[T]) Rfind(predicate func(T) bool) gust.Option[T] {
	if cover, ok := iter.facade.(iRealRfind[T]); ok {
		return cover.realRfind(predicate)
	}
	var check = func(f func(T) bool) func(any, T) gust.AnyCtrlFlow {
		return func(_ any, x T) gust.AnyCtrlFlow {
			if f(x) {
				return gust.AnyBreak(x)
			} else {
				return gust.AnyContinue(nil)
			}
		}
	}
	r := iter.TryRfold(nil, check(predicate))
	if r.IsBreak() {
		return gust.Some[T](r.UnwrapBreak().(T))
	}
	return gust.None[T]()
}

func (iter deIterBackground[T]) ToDeFuse() DeIterator[T] {
	return newDeFuseIterator[T](iter)
}

func (iter deIterBackground[T]) ToDePeekable() DePeekableIterator[T] {
	return newDePeekableIterator[T](iter)
}

func (iter deIterBackground[T]) ToDeSkip(n uint) DeIterator[T] {
	return newDeSkipIterator[T](iter, n)
}

func (iter deIterBackground[T]) ToDeTake(n uint) DeIterator[T] {
	return newDeTakeIterator[T](iter, n)
}

func (iter deIterBackground[T]) ToDeChain(other DeIterator[T]) DeIterator[T] {
	return newDeChainIterator[T](iter, other)
}

func (iter deIterBackground[T]) ToDeFilter(f func(T) bool) DeIterator[T] {
	return newDeFilterIterator[T](iter, f)
}

func (iter deIterBackground[T]) ToDeFilterMap(f func(T) gust.Option[T]) DeIterator[T] {
	return newDeFilterMapIterator[T, T](iter, f)
}

func (iter deIterBackground[T]) ToXDeFilterMap(f func(T) gust.Option[any]) DeIterator[any] {
	return newDeFilterMapIterator[T, any](iter, f)
}
func (iter deIterBackground[T]) ToDeInspect(f func(T)) DeIterator[T] {
	return newDeInspectIterator[T](iter, f)
}

func (iter deIterBackground[T]) ToDeMap(f func(T) T) DeIterator[T] {
	return newDeMapIterator[T, T](iter, f)
}

func (iter deIterBackground[T]) ToXDeMap(f func(T) any) DeIterator[any] {
	return newDeMapIterator[T, any](iter, f)
}
