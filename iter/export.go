package iter

import (
	"github.com/andeya/gust"
	"github.com/andeya/gust/digit"
)

func FromVec[T any](slice []T) *Iter[T] {
	return NewVecNext(slice).ToIter()
}

func FromRange[T digit.Integer](start T, end T, rightClosed ...bool) *Iter[T] {
	return NewRangeNext[T](start, end, rightClosed...).ToIter()
}

func FromChan[T any](c <-chan T) *Iter[T] {
	return NewChanNext[T](c).ToIter()
}

// FromIterator conversion from an [`Iterator`].
//
// By implementing `FromIterator` for a type, you define how it will be
// created from an iterator. This is common for types which describe a
// collection of some kind.
type FromIterator[T any, R any] interface {
	FromIter(Iterator[T]) R
}

// Collect collects all the items in the iterator into a slice.
func Collect[T any, R any](iter Iterator[T], fromIter FromIterator[T, R]) R {
	return fromIter.FromIter(iter)
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
// the checked sum of iAll the elements of the array
// var sum = FromVec(a).TryFold(0, func(acc int, x int) { return Ok(acc+x) });
//
// assert.Equal(t, sum, Ok(6));
func TryFold[T any, B any](next iNext[T], init B, f func(B, T) gust.Result[B]) gust.Result[B] {
	var accum = gust.Ok(init)
	for {
		x := next.Next()
		if x.IsNone() {
			return accum
		}
		accum = f(accum.Unwrap(), x.Unwrap())
		if accum.IsErr() {
			return accum
		}
	}
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
// This operation is sometimes called 'iReduce' or 'inject'.
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
// the sum of iAll the elements of the array
// var sum = FromVec(a).Fold((0, func(acc int, x int) any { return acc + x });
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
func Fold[T any, B any](next iNext[T], init B, f func(B, T) B) B {
	var accum = init
	for {
		x := next.Next()
		if x.IsNone() {
			return accum
		}
		accum = f(accum, x.Unwrap())
	}
}

// Map takes a closure and creates an iterator which calls that closure on each
// element.
//
// If you are good at thinking in types, you can think of `Map()` like this:
// If you have an iterator that gives you elements of some type `A`, and
// you want an iterator of some other type `B`, you can use `Map()`,
// passing a closure that takes an `A` and returns a `B`.
//
// `Map()` is conceptually similar to a [`for`] loop. However, as `Map()` is
// lazy, it is best used when you're already working with other iterators.
// If you're doing some sort of looping for a side effect, it's considered
// more idiomatic to use [`for`] than `Map()`.
//
// # Examples
//
// Basic usage:
//
// ```
// var a = []int{1, 2, 3};
//
// var iter = FromVec(a).Map(func(x)int{ return 2 * x});
//
// assert.Equal(iter.Next(), gust.Some(2));
// assert.Equal(iter.Next(), gust.Some(4));
// assert.Equal(iter.Next(), gust.Some(6));
// assert.Equal(iter.Next(), gust.None[int]());
// ```
func Map[T any, B any](iter Iterator[T], f func(T) B) *MapIterator[T, B] {
	return newMapIterator(iter, f)
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
// var first_number = FromVec(a).FindMap(func(s A) Option[any]{ return Wrap[any](strconv.Atoi(s))});
//
// assert.Equal(t, first_number, gust.Some(2));
func FindMap[T any, B any](iter Iterator[T], f func(T) gust.Option[B]) gust.Option[B] {
	for {
		x := iter.Next()
		if x.IsNone() {
			break
		}
		y := f(x.Unwrap())
		if y.IsSome() {
			return y
		}
	}
	return gust.None[B]()
}

// Zip 'Zips up' two iterators into a single iterator of pairs.
//
// `Zip()` returns a new iterator that will iterate over two other
// iterators, returning a tuple where the first element comes from the
// first iterator, and the second element comes from the second iterator.
//
// In other words, it zips two iterators together, into a single one.
//
// If either iterator returns [`gust.None[A]()`], [`Next`] from the zipped iterator
// will return [gust.None[A]()].
// If the zipped iterator has no more elements to return then each further attempt to advance
// it will first try to advance the first iterator at most one time and if it still yielded an item
// try to advance the second iterator at most one time.
func Zip[T any, B any](a Iterator[T], b Iterator[B]) *ZipIterator[T, B] {
	return newZipIterator[T, B](a, b)
}

// TryRfold is the reverse version of [`Iterator[T].TryFold()`]: it takes
// elements starting from the back of the iterator.
func TryRfold[T any, B any](next iNextBack[T], init B, f func(B, T) gust.Result[B]) gust.Result[B] {
	var accum = gust.Ok(init)
	for {
		x := next.NextBack()
		if x.IsNone() {
			return accum
		}
		accum = f(accum.Unwrap(), x.Unwrap())
		if accum.IsErr() {
			return accum
		}
	}
}

// Rfold is an iterator method that reduces the iterator's elements to a single,
// final value, starting from the back.
func Rfold[T any, B any](next iNextBack[T], init B, f func(B, T) B) B {
	var accum = init
	for {
		x := next.NextBack()
		if x.IsNone() {
			return accum
		}
		accum = f(accum, x.Unwrap())
	}
}
