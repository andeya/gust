package iter

import "github.com/andeya/gust"

type Iterator[T any] interface {
	iNext[T]
	iSizeHint
	iCount
	iFold[T]
	iTryFold[T]
	iLast[T]
	iAdvanceBy[T]
	iNth[T]
	iForEach[T]
	iReduce[T]
	iAll[T]
	iAny[T]
	iFind[T]
	iFindMap[T]
	iTryFind[T]
	iPosition[T]
	iStepBy[T]
	iFilter[T]
	iFilterMap[T]
	iChain[T]
}

type (
	iNext[T any] interface {
		// Next advances the next and returns the next value.
		//
		// Returns [`gust.None[T]()`] when iteration is finished. Individual next
		// implementations may choose to resume iteration, and so calling `next()`
		// again may or may not eventually min returning [`gust.Some(T)`] again at some
		// point.
		//
		// # Examples
		//
		// Basic usage:
		//
		// var a = []int{1, 2, 3};
		//
		// var iter = FromVec(a);
		//
		// // A call to next() returns the next value...
		// assert.Equal(t, gust.Some(1), iter.iNext());
		// assert.Equal(t, gust.Some(2), iter.iNext());
		// assert.Equal(t, gust.Some(3), iter.iNext());
		//
		// // ... and then None once it's over.
		// assert.Equal(t, gust.None[int](), iter.iNext());
		//
		// // More calls may or may not return `gust.None[T]()`. Here, they always will.
		// assert.Equal(t, gust.None[int](), iter.iNext());
		// assert.Equal(t, gust.None[int](), iter.iNext());
		Next() gust.Option[T]
	}
	iRealNext[T any] interface {
		realNext() gust.Option[T]
	}
	NextForIter[T any] interface {
		NextForIter() gust.Option[T]
	}
)
type (
	iCount interface {
		// Count consumes the next, counting the number of iterations and returning it.
		//
		// This method will call [`iNext`] repeatedly until [`gust.None[T]()`] is encountered,
		// returning the number of times it saw [`gust.Some`]. Note that [`iNext`] has to be
		// called at least once even if the next does not have any elements.
		//
		// # Overflow Behavior
		//
		// The method does no guarding against overflows, so counting elements of
		// a next with more than [`math.MaxInt`] elements either produces the
		// wrong result or panics. If debug assertions are enabled, a panic is
		// guaranteed.
		//
		// # Panics
		//
		// This function might panic if the next has more than [`math.MaxInt`]
		// elements.
		//
		// # Examples
		//
		// Basic usage:
		//
		// var a = []int{1, 2, 3};
		// assert.Equal(t, FromVec(a).iCount(), 3);
		//
		// var a = []int{1, 2, 3, 4, 5};
		// assert.Equal(t, FromVec(a).iCount(), 5);
		Count() uint64
	}
	iRealCount interface {
		realCount() uint64
	}
	CountForIter interface {
		CountForIter() uint64
	}
)
type (
	iSizeHint interface {
		// SizeHint returns the bounds on the remaining length of the next.
		//
		// Specifically, `iSizeHint()` returns a tuple where the first element
		// is the lower bound, and the second element is the upper bound.
		//
		// The second half of the tuple that is returned is an <code>Option[T]</code>.
		// A [`gust.None[T]()`] here means that either there is no known upper bound, or the
		// upper bound is larger than [`int`].
		//
		// # Implementation notes
		//
		// It is not enforced that a next implementation yields the declared
		// number of elements. A buggy next may yield less than the lower bound
		// or more than the upper bound of elements.
		//
		// `iSizeHint()` is primarily intended to be used for optimizations such as
		// reserving space for the elements of the next, but must not be
		// trusted to e.g., omit bounds checks in unsafe code. An incorrect
		// implementation of `iSizeHint()` should not lead to memory safety
		// violations.
		//
		// That said, the implementation should provide a correct estimation,
		// because otherwise it would be a violation of the interface's protocol.
		//
		// The default implementation returns <code>(0, [None[int]()])</code> which is correct for any
		// next.
		//
		// # Examples
		//
		// Basic usage:
		//
		// var a = []int{1, 2, 3};
		// var iter = FromVec(a);
		//
		// assert.Equal(t, (3, gust.Some(3)), iter.iSizeHint());
		//
		// A more complex example:
		//
		// // The even numbers in the range of zero to nine.
		// var iter = FromRange(0..10).Filter(func(x T) {return x % 2 == 0});
		//
		// // We might iterate from zero to ten times. Knowing that it's five
		// // exactly wouldn't be possible without executing filter().
		// assert.Equal(t, (0, gust.Some(10)), iter.iSizeHint());
		//
		// // Let's add five more numbers with chain()
		// var iter = FromRange(0, 10).Filter(func(x T) {return x % 2 == 0}).Chain(FromRange(15, 20));
		//
		// // now both bounds are increased by five
		// assert.Equal(t, (5, gust.Some(15)), iter.iSizeHint());
		//
		// Returning `gust.None[int]()` for an upper bound:
		//
		// // an infinite next has no upper bound
		// // and the maximum possible lower bound
		// var iter = FromRange(0, math.MaxInt);
		//
		// assert.Equal(t, (math.MaxInt, gust.None[int]()), iter.iSizeHint());
		SizeHint() (uint64, gust.Option[uint64])
	}
	iRealSizeHint interface {
		realSizeHint() (uint64, gust.Option[uint64])
	}
	SizeHintForIter interface {
		SizeHintForIter() (uint64, gust.Option[uint64])
	}
)
type (
	iFold[T any] interface {
		// Fold folds every element into an accumulator by applying an operation,
		// returning the final
		//
		// `iFold()` takes two arguments: an initial value, and a closure with two
		// arguments: an 'accumulator', and an element. The closure returns the value that
		// the accumulator should have for the next iteration.
		//
		// The initial value is the value the accumulator will have on the first
		// call.
		//
		// After applying this closure to every element of the next, `iFold()`
		// returns the accumulator.
		//
		// This operation is sometimes called 'iReduce' or 'inject'.
		//
		// Folding is useful whenever you have a collection of something, and want
		// to produce a single value from it.
		//
		// Note: `iFold()`, and similar methods that traverse the entire next,
		// might not terminate for infinite iterators, even on interfaces for which a
		// result is determinable in finite time.
		//
		// Note: [`Reduce()`] can be used to use the first element as the initial
		// value, if the accumulator type and item type is the same.
		//
		// Note: `iFold()` combines elements in a *left-associative* fashion. For associative
		// operators like `+`, the order the elements are combined in is not important, but for non-associative
		// operators like `-` the order will affect the final
		//
		// # Note to Implementors
		//
		// Several of the other (forward) methods have default implementations in
		// terms of this one, so try to implement this explicitly if it can
		// do something better than the default `for` loop implementation.
		//
		// In particular, try to have this call `iFold()` on the internal parts
		// from which this next is composed.
		//
		// # Examples
		//
		// Basic usage:
		//
		// var a = []int{1, 2, 3};
		//
		// the sum of iAll the elements of the array
		// var sum = FromVec(a).iFold((0, func(acc any, x T) any { return acc.(int) + x });
		//
		// assert.Equal(t, sum, 6);
		//
		// Let's walk through each step of the iteration here:
		//
		// | element | acc | x | result |
		// |---------|-----|---|--------|
		// |         | 0   |   |        |
		// | 1       | 0   | 1 | 1      |
		// | 2       | 1   | 2 | 3      |
		// | 3       | 3   | 3 | 6      |
		//
		// And so, our final result, `6`.
		Fold(init any, f func(any, T) any) any
	}
	iRealFold[T any] interface {
		realFold(init any, f func(any, T) any) any
	}
)
type (
	iTryFold[T any] interface {
		// TryFold a next method that applies a function as long as it returns
		// successfully, producing a single, final value.
		//
		// # Examples
		//
		// Basic usage:
		//
		// var a = []int{1, 2, 3};
		//
		// the checked sum of iAll the elements of the array
		// var sum = FromVec(a).TryFold(0, func(acc any, x T) { return Ok(acc.(int)+x) });
		//
		// assert.Equal(t, sum, Ok(6));
		TryFold(init any, f func(any, T) gust.Result[any]) gust.Result[any]
	}
	iRealTryFold[T any] interface {
		realTryFold(init any, f func(any, T) gust.Result[any]) gust.Result[any]
	}
)
type (
	iLast[T any] interface {
		// Last consumes the next, returning the iLast element.
		//
		// This method will evaluate the next until it returns [`None[T]()`]. While
		// doing so, it keeps track of the current element. After [`None[T]()`] is
		// returned, `Last()` will then return the iLast element it saw.
		//
		// # Examples
		//
		// Basic usage:
		//
		// var a = []int{1, 2, 3};
		// assert.Equal(t, FromVec(a).Last(), Some(3));
		//
		// var a = [1, 2, 3, 4, 5];
		// assert.Equal(t, FromVec(a).Last(), Some(5));
		Last() gust.Option[T]
	}
	iRealLast[T any] interface {
		realLast() gust.Option[T]
	}
)
type (
	iAdvanceBy[T any] interface {
		// AdvanceBy advances the next by `n` elements.
		//
		// This method will eagerly skip `n` elements by calling [`iNext`] up to `n`
		// times until [`None[T]()`] is encountered.
		//
		// `AdvanceBy(n)` will return [`Ok[struct{}](struct{}{})`] if the next successfully advances by
		// `n` elements, or [`Err[struct{}](err)`] if [`None[T]()`] is encountered, where `k` is the number
		// of elements the next is advanced by before running out of elements (i.e. the
		// length of the next). Note that `k` is always less than `n`.
		//
		// Calling `AdvanceBy(0)` can do meaningful work, for example [`Flatten`]
		// can advance its outer next until it finds an facade next that is not empty, which
		// then often allows it to return a more accurate `iSizeHint()` than in its initial state.
		// `AdvanceBy(0)` may either return `T()` or `Err(0)`. The former conveys no information
		// whether the next is or is not exhausted, the latter can be treated as if [`iNext`]
		// had returned `None[T]()`. Replacing a `Err(0)` with `T` is only correct for `n = 0`.
		//
		// # Examples
		//
		// Basic usage:
		//
		// var a = []int{1, 2, 3, 4};
		// var iter = FromVec(a);
		//
		// assert.Equal(t, iter.AdvanceBy(2), Ok[struct{}](struct{}{}));
		// assert.Equal(t, iter.iNext(), Some(3));
		// assert.Equal(t, iter.AdvanceBy(0), Ok[struct{}](struct{}{}));
		// assert.Equal(t, iter.AdvanceBy(100), Err[struct{}](fmt.Errorf("%d", 1))); // only `4` was skipped
		AdvanceBy(n uint) gust.Result[struct{}]
	}
	iRealAdvanceBy[T any] interface {
		realAdvanceBy(n uint) gust.Result[struct{}]
	}
)
type (
	iNth[T any] interface {
		// Nth returns the `n`th element of the next.
		//
		// Like most indexing operations, the iCount starts from zero, so `Nth(0)`
		// returns the first value, `Nth(1)` the second, and so on.
		//
		// Note that iAll preceding elements, as well as the returned element, will be
		// consumed from the next. That means that the preceding elements will be
		// discarded, and also that calling `iNth(0)` multiple times on the same next
		// will return different elements.
		//
		// `Nth()` will return [`None[T]()`] if `n` is greater than or equal to the length of the
		// next.
		//
		// # Examples
		//
		// Basic usage:
		//
		// var a = []int{1, 2, 3};
		// assert.Equal(t, FromVec(a).Nth(1), Some(2));
		//
		// Calling `Nth()` multiple times doesn't rewind the next:
		//
		// var a = []int{1, 2, 3};
		//
		// var iter = FromVec(a);
		//
		// assert.Equal(t, iter.Nth(1), Some(2));
		// assert.Equal(t, iter.Nth(1), None[int]());
		//
		// Returning `None[T]()` if there are less than `n + 1` elements:
		//
		// var a = []int{1, 2, 3};
		// assert.Equal(t, FromVec(a).Nth(10), None[int]());
		Nth(n uint) gust.Option[T]
	}
	iRealNth[T any] interface {
		realNth(n uint) gust.Option[T]
	}
)
type (
	iForEach[T any] interface {
		// ForEach calls a closure on each element of a next.
		//
		// This is equivalent to using a [`for`] loop on the next, although
		// `break` and `continue` are not possible from a closure. It's generally
		// more idiomatic to use a `for` loop, but `ForEach` may be more legible
		// when processing items at the end of longer next chains. In some
		// cases `ForEach` may also be faster than a loop, because it will use
		// internal iteration on adapters like `Chain`.
		//
		// # Examples
		//
		// Basic usage:
		//
		// var c = make(chan int, 1000)
		// FromRange(0, 5).Map(func(x T)any{return x * 2 + 1})
		//
		//	.ForEach(func(x any){ c<-x });
		//
		// var v = FromChan(c).Collect();
		// assert.Equal(t, v, []int{1, 3, 5, 7, 9});
		ForEach(f func(T))
	}
	iRealForEach[T any] interface {
		realForEach(f func(T))
	}
)
type (
	iReduce[T any] interface {
		// Reduce reduces the elements to a single one, by repeatedly applying a reducing
		// operation.
		//
		// If the next is empty, returns [`gust.None[T]()`]; otherwise, returns the
		// result of the reduction.
		//
		// The reducing function is a closure with two arguments: an 'accumulator', and an element.
		// For iterators with at least one element, this is the same as [`iFold()`]
		// with the first element of the next as the initial accumulator value, folding
		// every subsequent element into it.
		//
		// # Example
		//
		// Find the maximum value:
		//
		//	func findMax[T any](iter: Iterator[T])  Option[T] {
		//	    iter.Reduce(func(accum T, item T) T {
		//	        if accum >= item { accum } else { item }
		//	    })
		//	}
		//
		// var a = []int{10, 20, 5, -23, 0};
		// var b = []int{};
		//
		// assert.Equal(t, findMax(FromVec(a)), gust.Some(20));
		// assert.Equal(t, findMax(FromVec(b)), gust.None[int]());
		Reduce(f func(accum T, item T) T) gust.Option[T]
	}
	iRealReduce[T any] interface {
		realReduce(f func(accum T, item T) T) gust.Option[T]
	}
)
type (
	iAll[T any] interface {
		// All tests if every element of the next matches a predicate.
		//
		// `All()` takes a closure that returns `true` or `false`. It applies
		// this closure to each element of the next, and if they iAll return
		// `true`, then so does `All()`. If any of them return `false`, it
		// returns `false`.
		//
		// `All()` is short-circuiting; in other words, it will stop processing
		// as soon as it finds a `false`, given that no matter what else happens,
		// the result will also be `false`.
		//
		// An empty next returns `true`.
		//
		// # Examples
		//
		// Basic usage:
		//
		// var a = []int{1, 2, 3};
		//
		// assert.True(t, FromVec(a).All(func(x T) bool { return x > 0}));
		//
		// assert.True(t, !FromVec(a).All(func(x T) bool { return x > 2}));
		//
		// Stopping at the first `false`:
		//
		// var a = []int{1, 2, 3};
		//
		// var iter = FromVec(a);
		//
		// assert.True(t, !iter.All(func(x T) bool { return x != 2}));
		//
		// we can still use `iter`, as there are more elements.
		// assert.Equal(t, iter.iNext(), gust.Some(3));
		All(predicate func(T) bool) bool
	}
	iRealAll[T any] interface {
		realAll(predicate func(T) bool) bool
	}
)
type (
	iAny[T any] interface {
		// Any tests if any element of the next matches a predicate.
		//
		// `Any()` takes a closure that returns `true` or `false`. It applies
		// this closure to each element of the next, and if any of them return
		// `true`, then so does `Any()`. If they iAll return `false`, it
		// returns `false`.
		//
		// `Any()` is short-circuiting; in other words, it will stop processing
		// as soon as it finds a `true`, given that no matter what else happens,
		// the result will also be `true`.
		//
		// An empty next returns `false`.
		//
		// # Examples
		//
		// Basic usage:
		//
		// var a = []int{1, 2, 3};
		//
		// assert.True(t, FromVec(a).Any(func(x T) bool{return x>0}));
		//
		// assert.True(t, !FromVec(a).Any(func(x T) bool{return x>5}));
		//
		// Stopping at the first `true`:
		//
		// var a = []int{1, 2, 3};
		//
		// var iter = FromVec(a);
		//
		// assert.True(t, iter.Any(func(x T) bool { return x != 2}));
		//
		// we can still use `iter`, as there are more elements.
		// assert.Equal(t, iter.iNext(), gust.Some(2));
		Any(predicate func(T) bool) bool
	}
	iRealAny[T any] interface {
		realAny(predicate func(T) bool) bool
	}
)
type (
	iFind[T any] interface {
		// Find searches for an element of a next that satisfies a predicate.
		//
		// `Find()` takes a closure that returns `true` or `false`. It applies
		// this closure to each element of the next, and if any of them return
		// `true`, then `Find()` returns [`gust.Some(element)`]. If they iAll return
		// `false`, it returns [`gust.None[T]()`].
		//
		// `Find()` is short-circuiting; in other words, it will stop processing
		// as soon as the closure returns `true`.
		//
		// Because `Find()` takes a reference, and many iterators iterate over
		// references, this leads to a possibly confusing situation where the
		// argument is a double reference. You can see this effect in the
		// examples below, with `&&x`.
		//
		// [`gust.Some(element)`]: gust.Some
		//
		// # Examples
		//
		// Basic usage:
		//
		// var a = []int{1, 2, 3};
		//
		// assert.Equal(t, FromVec(a).Find(func(x T) bool{return x==2}), gust.Some(2));
		//
		// assert.Equal(t, FromVec(a).Find(func(x T) bool{return x==5}), gust.None[int]());
		//
		// Stopping at the first `true`:
		//
		// var a = []int{1, 2, 3};
		//
		// var iter = FromVec(a);
		//
		// assert.Equal(t, iter.Find(func(x T) bool{return x==2}), gust.Some(2));
		//
		// we can still use `iter`, as there are more elements.
		// assert.Equal(t, iter.iNext(), gust.Some(3));
		//
		// Note that `iter.Find(f)` is equivalent to `iter.Filter(f).iNext()`.
		Find(predicate func(T) bool) gust.Option[T]
	}
	iRealFind[T any] interface {
		realFind(predicate func(T) bool) gust.Option[T]
	}
)
type (
	iFindMap[T any] interface {
		// FindMap applies function to the elements of next and returns
		// the first non-none
		//
		// `iter.FindMap(f)` is equivalent to `iter.FilterMap(f).iNext()`.
		//
		// # Examples
		//
		// var a = []string{"lol", "NaN", "2", "5"};
		//
		// var first_number = FromVec(a).FindMap(func(s T) Option[any]{ return Wrap[any](strconv.Atoi(s))});
		//
		// assert.Equal(t, first_number, gust.Some(2));
		FindMap(f func(T) gust.Option[any]) gust.Option[any]
	}
	iRealFindMap[T any] interface {
		realFindMap(f func(T) gust.Option[any]) gust.Option[any]
	}
)
type (
	iTryFind[T any] interface {
		// TryFind applies function to the elements of next and returns
		// the first true result or the first error.
		//
		// # Examples
		//
		// var a = []string{"1", "2", "lol", "NaN", "5"};
		//
		//	var is_my_num = func(s string, search int) Result[bool] {
		//	    return ret.Map(gust.Ret(strconv.Atoi(s)), func(x int) bool { return x == search })
		//	}
		//
		// var result = FromVec(a).TryFind(func(s string)bool{return is_my_num(s, 2)});
		// assert.Equal(t, result, T(Some("2")));
		//
		// var result = FromVec(a).TryFind(func(s string)bool{return is_my_num(s, 5)});
		// assert.True(t, IsErr());
		TryFind(predicate func(T) gust.Result[bool]) gust.Result[gust.Option[T]]
	}
	iRealTryFind[T any] interface {
		realTryFind(predicate func(T) gust.Result[bool]) gust.Result[gust.Option[T]]
	}
)
type (
	iPosition[T any] interface {
		// Position searches for an element in a next, returning its index.
		//
		// `Position()` takes a closure that returns `true` or `false`. It applies
		// this closure to each element of the next, and if one of them
		// returns `true`, then `Position()` returns [`gust.Some(index)`]. If iAll of
		// them return `false`, it returns [`gust.None[T]()`].
		//
		// `Position()` is short-circuiting; in other words, it will stop
		// processing as soon as it finds a `true`.
		//
		// # Overflow Behavior
		//
		// The method does no guarding against overflows, so if there are more
		// than [`math.MaxInt`] non-matching elements, it either produces the wrong
		// result or panics. If debug assertions are enabled, a panic is
		// guaranteed.
		//
		// # Panics
		//
		// This function might panic if the next has more than `math.MaxInt`
		// non-matching elements.
		//
		// [`gust.Some(index)`]: gust.Some
		//
		// # Examples
		//
		// Basic usage:
		//
		// var a = []int{1, 2, 3};
		//
		// assert.Equal(t, FromVec(a).Position(func(x int)bool{return x==2}), gust.Some(1));
		//
		// assert.Equal(t, FromVec(a).Position(func(x int)bool{return x==5}), gust.None[int]());
		//
		// Stopping at the first `true`:
		//
		// var a = []int{1, 2, 3, 4};
		//
		// var iter = FromVec(a);
		//
		// assert.Equal(t, iter.Position(func(x int)bool{return x >= 2}), gust.Some(1));
		//
		// we can still use `iter`, as there are more elements.
		// assert.Equal(t, iter.iNext(), gust.Some(3));
		//
		// The returned index depends on next state
		// assert.Equal(t, iter.Position(func(x int)bool{return x == 4}), gust.Some(0));
		Position(predicate func(T) bool) gust.Option[int]
	}
	iRealPosition[T any] interface {
		realPosition(predicate func(T) bool) gust.Option[int]
	}
)
type (
	iStepBy[T any] interface {
		// StepBy creates a next starting at the same point, but stepping by
		// the given amount at each iteration.
		//
		// Note 1: The first element of the next will always be returned,
		// regardless of the step given.
		//
		// Note 2: The time at which ignored elements are pulled is not fixed.
		// `StepBy` behaves like the sequence `iter.iNext()`, `iter.Nth(step-1)`,
		// `iter.Nth(step-1)`, …, but is also free to behave like the sequence.
		//
		// # Examples
		//
		// Basic usage:
		//
		// var a = []int{0, 1, 2, 3, 4, 5};
		// var iter = FromVec(a).StepBy(2);
		//
		// assert.Equal(t, iter.iNext(), Some(0));
		// assert.Equal(t, iter.iNext(), Some(2));
		// assert.Equal(t, iter.iNext(), Some(4));
		// assert.Equal(t, iter.iNext(), None[T]());
		StepBy(step uint) *StepBy[T]
	}
	iRealStepBy[T any] interface {
		realStepBy(step uint) *StepBy[T]
	}
)
type (
	iFilter[T any] interface {
		// Filter creates an iterator which uses a closure to determine if an element
		// should be yielded.
		//
		// Given an element the closure must return `true` or `false`. The returned
		// iterator will yield only the elements for which the closure returns
		// true.
		//
		// # Examples
		//
		// Basic usage:
		//
		// ```
		// var a = []int{0, 1, 2};
		//
		// var iter = FromVec(a).Filter(func(x int)bool{return x>0});
		//
		// assert_eq!(iter.iNext(), gust.Some(&1));
		// assert_eq!(iter.iNext(), gust.Some(&2));
		// assert_eq!(iter.iNext(), gust.None[int]());
		// ```
		//
		// Note that `iter.Filter(f).iNext()` is equivalent to `iter.Find(f)`.
		Filter(f func(T) bool) *Filter[T]
	}
	iRealFilter[T any] interface {
		realFilter(f func(T) bool) *Filter[T]
	}
)
type (
	iFilterMap[T any] interface {
		// FilterMap creates an iterator that both filters and maps.
		//
		// The returned iterator yields only the `value`s for which the supplied
		// closure returns `Some(value)`.
		//
		// `FilterMap` can be used to make chains of [`filter`] and [`map`] more
		// concise. The example below shows how a `map().filter().map()` can be
		// shortened to a single call to `FilterMap`.
		//
		// # Examples
		//
		// Basic usage:
		//
		// ```
		// var a = []string{"1", "two", "NaN", "four", "5"}
		//
		// var iter = FromVec(a).FilterMap(|s| s.parse().ok());
		//
		// assert_eq!(iter.next(), gust.Some(1));
		// assert_eq!(iter.next(), gust.Some(5));
		// assert_eq!(iter.next(), gust.None[string]());
		// ```
		FilterMap(f func(T) gust.Option[T]) *FilterMap[T]
	}
	iRealFilterMap[T any] interface {
		realFilterMap(f func(T) gust.Option[T]) *FilterMap[T]
	}
)
type (
	iChain[T any] interface {
		// Chain takes two iterators and creates a new next over both in sequence.
		//
		// `Chain()` will return a new next which will first iterate over
		// values from the first next and then over values from the second
		// next.
		//
		// In other words, it links two iterators together, in a chain. 🔗
		//
		// [`once`] is commonly used to adapt a single value into a chain of
		// other kinds of iteration.
		//
		// # Examples
		//
		// Basic usage:
		//
		//
		// var a1 = []int{1, 2, 3};
		// var a2 = []int{4, 5, 6};
		//
		// var iter = FromVec(a1).Chain(FromVec(a2));
		//
		// assert.Equal(t, iter.iNext(), Some(1));
		// assert.Equal(t, iter.iNext(), Some(2));
		// assert.Equal(t, iter.iNext(), Some(3));
		// assert.Equal(t, iter.iNext(), Some(4));
		// assert.Equal(t, iter.iNext(), Some(5));
		// assert.Equal(t, iter.iNext(), Some(6));
		// assert.Equal(t, iter.iNext(), None[int]());
		//
		Chain(other Iterator[T]) *Chain[T]
	}
	iRealChain[T any] interface {
		realChain(other Iterator[T]) *Chain[T]
	}
)
