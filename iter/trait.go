package iter

import (
	"fmt"

	"github.com/andeya/gust"
)

type (
	iterTrait[T any, C iterCore[T]] struct {
		core C
	}
	iterCore[T any] interface {
		Nextor[T]
		SizeHint
		Count() uint64
	}

	SizeHint interface {
		// SizeHint returns the bounds on the remaining length of the next.
		//
		// Specifically, `SizeHint()` returns a tuple where the first element
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
		// `SizeHint()` is primarily intended to be used for optimizations such as
		// reserving space for the elements of the next, but must not be
		// trusted to e.g., omit bounds checks in unsafe code. An incorrect
		// implementation of `SizeHint()` should not lead to memory safety
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
		//
		// var a = []int{1, 2, 3};
		// var iter = IterAnyFromVec(a);
		//
		// assert.Equal(t, (3, gust.Some(3)), iter.SizeHint());
		//
		//
		// A more complex example:
		//
		//
		// // The even numbers in the range of zero to nine.
		// var iter = IterAnyFromRange(0..10).Filter(func(x T) {return x % 2 == 0});
		//
		// // We might iterate from zero to ten times. Knowing that it's five
		// // exactly wouldn't be possible without executing filter().
		// assert.Equal(t, (0, gust.Some(10)), iter.SizeHint());
		//
		// // Let's add five more numbers with chain()
		// var iter = IterAnyFromRange(0, 10).Filter(func(x T) {return x % 2 == 0}).Chain(IterAnyFromRange(15, 20));
		//
		// // now both bounds are increased by five
		// assert.Equal(t, (5, gust.Some(15)), iter.SizeHint());
		//
		//
		// Returning `gust.None[int]()` for an upper bound:
		//
		//
		// // an infinite next has no upper bound
		// // and the maximum possible lower bound
		// var iter = IterAnyFromRange(0, math.MaxInt);
		//
		// assert.Equal(t, (math.MaxInt, gust.None[int]()), iter.SizeHint());
		//
		SizeHint() (uint64, gust.Option[uint64])
	}
)

func (iter iterTrait[T, N]) Next() gust.Option[T] {
	return iter.core.Next()
}

func (iter iterTrait[T, N]) SizeHint() (uint64, gust.Option[uint64]) {
	return iter.core.SizeHint()
}

func (iter iterTrait[T, N]) Count() uint64 {
	return iter.core.Count()
}

// Fold folds every element into an accumulator by applying an operation,
// returning the final
//
// `Fold()` takes two arguments: an initial value, and a closure with two
// arguments: an 'accumulator', and an element. The closure returns the value that
// the accumulator should have for the next iteration.
//
// The initial value is the value the accumulator will have on the first
// call.
//
// After applying this closure to every element of the next, `Fold()`
// returns the accumulator.
//
// This operation is sometimes called 'reduce' or 'inject'.
//
// Folding is useful whenever you have a collection of something, and want
// to produce a single value from it.
//
// Note: `Fold()`, and similar methods that traverse the entire next,
// might not terminate for infinite iterators, even on interfaces for which a
// result is determinable in finite time.
//
// Note: [`Reduce()`] can be used to use the first element as the initial
// value, if the accumulator type and item type is the same.
//
// Note: `Fold()` combines elements in a *left-associative* fashion. For associative
// operators like `+`, the order the elements are combined in is not important, but for non-associative
// operators like `-` the order will affect the final
//
// # Note to Implementors
//
// Several of the other (forward) methods have default implementations in
// terms of this one, so try to implement this explicitly if it can
// do something better than the default `for` loop implementation.
//
// In particular, try to have this call `Fold()` on the internal parts
// from which this next is composed.
//
// # Examples
//
// Basic usage:
//
// var a = []int{1, 2, 3};
//
// the sum of all the elements of the array
// var sum = IterAnyFromVec(a).Fold((0, func(acc any, x T) any { return acc.(int) + x });
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
func (iter iterTrait[T, N]) Fold(init any, f func(any, T) any) any {
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

// TryFold a next method that applies a function as long as it returns
// successfully, producing a single, final value.
//
// # Examples
//
// Basic usage:
//
// var a = []int{1, 2, 3};
//
// the checked sum of all the elements of the array
// var sum = IterAnyFromVec(a).TryFold(0, func(acc any, x T) { return Ok(acc.(int)+x) });
//
// assert.Equal(t, sum, Ok(6));
func (iter iterTrait[T, N]) TryFold(init any, f func(any, T) gust.Result[any]) gust.Result[any] {
	return TryFold[T, any](iter, init, f)
}

// Last consumes the next, returning the last element.
//
// This method will evaluate the next until it returns [`None[T]()`]. While
// doing so, it keeps track of the current element. After [`None[T]()`] is
// returned, `Last()` will then return the last element it saw.
//
// # Examples
//
// Basic usage:
//
// var a = []int{1, 2, 3};
// assert.Equal(t, IterAnyVec(a).Last(), Some(3));
//
// var a = [1, 2, 3, 4, 5];
// assert.Equal(t, IterAnyVec(a).Last(), Some(5));
func (iter iterTrait[T, N]) Last() gust.Option[T] {
	return iter.Fold(gust.None[T](), func(_ any, x T) any { return gust.Some(x) }).(gust.Option[T])
}

// AdvanceBy advances the next by `n` elements.
//
// This method will eagerly skip `n` elements by calling [`Nextor`] up to `n`
// times until [`None[T]()`] is encountered.
//
// `AdvanceBy(n)` will return [`Ok[struct{}](struct{}{})`] if the next successfully advances by
// `n` elements, or [`Err[struct{}](err)`] if [`None[T]()`] is encountered, where `k` is the number
// of elements the next is advanced by before running out of elements (i.e. the
// length of the next). Note that `k` is always less than `n`.
//
// Calling `AdvanceBy(0)` can do meaningful work, for example [`Flatten`]
// can advance its outer next until it finds an core next that is not empty, which
// then often allows it to return a more accurate `SizeHint()` than in its initial state.
// `AdvanceBy(0)` may either return `T()` or `Err(0)`. The former conveys no information
// whether the next is or is not exhausted, the latter can be treated as if [`Nextor`]
// had returned `None[T]()`. Replacing a `Err(0)` with `T` is only correct for `n = 0`.
//
// # Examples
//
// Basic usage:
//
// var a = []int{1, 2, 3, 4};
// var iter = IterAnyVec(a);
//
// assert.Equal(t, iter.AdvanceBy(2), Ok[struct{}](struct{}{}));
// assert.Equal(t, iter.Next(), Some(3));
// assert.Equal(t, iter.AdvanceBy(0), Ok[struct{}](struct{}{}));
// assert.Equal(t, iter.AdvanceBy(100), Err[struct{}](fmt.Errorf("%d", 1))); // only `4` was skipped
func (iter iterTrait[T, N]) AdvanceBy(n uint) gust.Result[struct{}] {
	for i := uint(0); i < n; i++ {
		if iter.Next().IsNone() {
			return gust.Err[struct{}](fmt.Errorf("%d", i))
		}
	}
	return gust.Ok(struct{}{})
}

// Nth returns the `n`th element of the next.
//
// Like most indexing operations, the count starts from zero, so `Nth(0)`
// returns the first value, `Nth(1)` the second, and so on.
//
// Note that all preceding elements, as well as the returned element, will be
// consumed from the next. That means that the preceding elements will be
// discarded, and also that calling `nth(0)` multiple times on the same next
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
// assert.Equal(t, IterAnyVec(a).Nth(1), Some(2));
//
// Calling `Nth()` multiple times doesn't rewind the next:
//
// var a = []int{1, 2, 3};
//
// var iter = IterAnyVec(a);
//
// assert.Equal(t, iter.Nth(1), Some(2));
// assert.Equal(t, iter.Nth(1), None[int]());
//
// Returning `None[T]()` if there are less than `n + 1` elements:
//
// var a = []int{1, 2, 3};
// assert.Equal(t, IterAnyVec(a).Nth(10), None[int]());
func (iter iterTrait[T, N]) Nth(n uint) gust.Option[T] {
	var res = iter.AdvanceBy(n)
	if res.IsErr() {
		return gust.None[T]()
	}
	return iter.Next()
}

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
// IterAnyFromRange(0, 5).Map(func(x T)any{return x * 2 + 1})
//
//	.ForEach(func(x any){ c<-x });
//
// var v = IterAnyFromChan(c).Collect();
// assert.Equal(t, v, []int{1, 3, 5, 7, 9});
func (iter iterTrait[T, N]) ForEach(f func(T)) {
	var call = func(f func(T)) func(any, T) any {
		return func(_ any, item T) any {
			f(item)
			return nil
		}
	}
	_ = iter.Fold(nil, call(f))
}

// Reduce reduces the elements to a single one, by repeatedly applying a reducing
// operation.
//
// If the next is empty, returns [`gust.None[T]()`]; otherwise, returns the
// result of the reduction.
//
// The reducing function is a closure with two arguments: an 'accumulator', and an element.
// For iterators with at least one element, this is the same as [`Fold()`]
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
// assert.Equal(t, findMax(IterAnyFromVec(a)), gust.Some(20));
// assert.Equal(t, findMax(IterAnyFromVec(b)), gust.None[int]());
func (iter iterTrait[T, N]) Reduce(f func(accum T, item T) T) gust.Option[T] {
	var first = iter.Next()
	if first.IsNone() {
		return first
	}
	return gust.Some(iter.Fold(first, func(accum any, item T) any {
		return f(accum.(T), item)
	}).(T))
}

// All tests if every element of the next matches a predicate.
//
// `All()` takes a closure that returns `true` or `false`. It applies
// this closure to each element of the next, and if they all return
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
// assert.True(t, IterAnyFromVec(a).All(func(x T) bool { return x > 0}));
//
// assert.True(t, !IterAnyFromVec(a).All(func(x T) bool { return x > 2}));
//
// Stopping at the first `false`:
//
// var a = []int{1, 2, 3};
//
// var iter = IterAnyFromVec(a);
//
// assert.True(t, !iter.All(func(x T) bool { return x != 2}));
//
// we can still use `iter`, as there are more elements.
// assert.Equal(t, iter.Next(), gust.Some(3));
func (iter iterTrait[T, N]) All(predicate func(T) bool) bool {
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

// Any tests if any element of the next matches a predicate.
//
// `Any()` takes a closure that returns `true` or `false`. It applies
// this closure to each element of the next, and if any of them return
// `true`, then so does `Any()`. If they all return `false`, it
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
// assert.True(t, IterAnyFromVec(a).Any(func(x T) bool{return x>0}));
//
// assert.True(t, !IterAnyFromVec(a).Any(func(x T) bool{return x>5}));
//
// Stopping at the first `true`:
//
// var a = []int{1, 2, 3};
//
// var iter = IterAnyFromVec(a);
//
// assert.True(t, iter.Any(func(x T) bool { return x != 2}));
//
// we can still use `iter`, as there are more elements.
// assert.Equal(t, iter.Next(), gust.Some(2));
func (iter iterTrait[T, N]) Any(predicate func(T) bool) bool {
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

// Find searches for an element of a next that satisfies a predicate.
//
// `Find()` takes a closure that returns `true` or `false`. It applies
// this closure to each element of the next, and if any of them return
// `true`, then `Find()` returns [`gust.Some(element)`]. If they all return
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
// assert.Equal(t, IterAnyFromVec(a).Find(func(x T) bool{return x==2}), gust.Some(2));
//
// assert.Equal(t, IterAnyFromVec(a).Find(func(x T) bool{return x==5}), gust.None[int]());
//
// Stopping at the first `true`:
//
// var a = []int{1, 2, 3};
//
// var iter = IterAnyFromVec(a);
//
// assert.Equal(t, iter.Find(func(x T) bool{return x==2}), gust.Some(2));
//
// we can still use `iter`, as there are more elements.
// assert.Equal(t, iter.Next(), gust.Some(3));
//
// Note that `iter.Find(f)` is equivalent to `iter.Filter(f).Next()`.
func (iter iterTrait[T, N]) Find(predicate func(T) bool) gust.Option[T] {
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

// FindMap applies function to the elements of next and returns
// the first non-none
//
// `iter.FindMap(f)` is equivalent to `iter.FilterMap(f).Next()`.
//
// # Examples
//
// var a = []string{"lol", "NaN", "2", "5"};
//
// var first_number = IterAnyFromVec(a).FindMap(func(s T) Option[any]{ return Wrap[any](strconv.Atoi(s))});
//
// assert.Equal(t, first_number, gust.Some(2));
func (iter iterTrait[T, N]) FindMap(f func(T) gust.Option[any]) gust.Option[any] {
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
// var result = IterAnyFromVec(a).TryFind(func(s string)bool{return is_my_num(s, 2)});
// assert.Equal(t, result, T(Some("2")));
//
// var result = IterAnyFromVec(a).TryFind(func(s string)bool{return is_my_num(s, 5)});
// assert.True(t, IsErr());
func (iter iterTrait[T, N]) TryFind(predicate func(T) gust.Result[bool]) gust.Result[gust.Option[T]] {
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

// Position searches for an element in a next, returning its index.
//
// `Position()` takes a closure that returns `true` or `false`. It applies
// this closure to each element of the next, and if one of them
// returns `true`, then `Position()` returns [`gust.Some(index)`]. If all of
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
// assert.Equal(t, IterAnyFromVec(a).Position(func(x int)bool{return x==2}), gust.Some(1));
//
// assert.Equal(t, IterAnyFromVec(a).Position(func(x int)bool{return x==5}), gust.None[int]());
//
// Stopping at the first `true`:
//
// var a = []int{1, 2, 3, 4};
//
// var iter = IterAnyFromVec(a);
//
// assert.Equal(t, iter.Position(func(x int)bool{return x >= 2}), gust.Some(1));
//
// we can still use `iter`, as there are more elements.
// assert.Equal(t, iter.Next(), gust.Some(3));
//
// The returned index depends on next state
// assert.Equal(t, iter.Position(func(x int)bool{return x == 4}), gust.Some(0));
func (iter iterTrait[T, N]) Position(predicate func(T) bool) gust.Option[int] {
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
