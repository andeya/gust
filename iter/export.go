package iter

import (
	"unicode/utf8"
	"unsafe"

	"github.com/andeya/gust"
	"github.com/andeya/gust/digit"
)

// FromIterable creates an iterator from an Iterable.
func FromIterable[T any](data gust.Iterable[T]) Iterator[T] {
	iter, _ := data.(Iterator[T])
	if iter != nil {
		return iter
	}
	return fromIterable[T](data)
}

// FromDeIterable creates an iterator from an Iterable.
func FromDeIterable[T any](data gust.DeIterable[T]) DeIterator[T] {
	iter, _ := data.(DeIterator[T])
	if iter != nil {
		return iter
	}
	return fromDeIterable[T](data)
}

// FromVec creates an iterator from a slice.
func FromVec[T any](slice []T) DeIterator[T] {
	return NewIterableVec(slice).ToDeIterator()
}

// FromElements creates an iterator from a set of elements.
func FromElements[T any](elem ...T) DeIterator[T] {
	return NewIterableVec(elem).ToDeIterator()
}

// FromRange creates an iterator from a range.
func FromRange[T digit.Integer](start T, end T, rightClosed ...bool) DeIterator[T] {
	return NewIterableRange[T](start, end, rightClosed...).ToDeIterator()
}

// FromChan creates an iterator from a channel.
func FromChan[T any](c chan T) Iterator[T] {
	return NewIterableChan[T](c).ToIterator()
}

// FromResult creates an iterator from a result.
func FromResult[T any](ret gust.Result[T]) DeIterator[T] {
	return FromDeIterable[T](ret)
}

// FromOption creates an iterator from an option.
func FromOption[T any](opt gust.Option[T]) DeIterator[T] {
	return FromDeIterable[T](opt)
}

// FromString creates an iterator from a string.
func FromString[T ~byte | ~rune](s string) DeIterator[T] {
	if len(s) == 0 {
		return NewIterableVec[T]([]T{}).ToDeIterator()
	}
	const bn = rune(^byte(0))
	if bn == rune(^T(0)) {
		var rs = *(*[]T)(unsafe.Pointer(
			&struct {
				string
				Cap int
			}{s, len(s)},
		))
		return NewIterableVec[T](rs).ToDeIterator()
	}
	var rs = make([]T, 0, len(s))
	var b = *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		rs = append(rs, T(r))
		b = b[size:]
	}
	return NewIterableVec[T](rs).ToDeIterator()
}

// TryFold a data method that applies a function as long as it returns
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
func TryFold[T any, B any](iter Iterator[T], init B, f func(B, T) gust.Result[B]) gust.Result[B] {
	var accum = gust.Ok(init)
	for {
		x := iter.Next()
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
// the accumulator should have for the data iteration.
//
// The initial value is the value the accumulator will have on the first
// call.
//
// After applying this closure to every element of the data, `Fold()`
// returns the accumulator.
//
// This operation is sometimes called 'iReduce' or 'inject'.
//
// Folding is useful whenever you have a collection of something, and want
// to produce a single value from it.
//
// Note: `Fold()`, and similar methods that traverse the entire data,
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
// from which this data is composed.
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
func Fold[T any, B any](iter Iterator[T], init B, f func(B, T) B) B {
	var accum = init
	for {
		x := iter.Next()
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
func Map[T any, B any](iter Iterator[T], f func(T) B) Iterator[B] {
	return newMapIterator(iter, f)
}

// FilterMap creates an iterator that both filters and maps.
//
// The returned iterator yields only the `value`s for which the supplied
// closure returns `gust.Some(value)`.
//
// `FilterMap` can be used to make chains of [`Filter`] and [`Map`] more
// concise. The example below shows how a `Map().Filter().Map()` can be
// shortened to a single call to `FilterMap`.
func FilterMap[T any, B any](iter Iterator[T], f func(T) gust.Option[B]) Iterator[B] {
	return newFilterMapIterator[T, B](iter, f)
}

// FindMap applies function to the elements of data and returns
// the first non-none
//
// `FindMap(iter, f)` is equivalent to `FilterMap(iter, f).Next()`.
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
func Zip[A any, B any](a Iterator[A], b Iterator[B]) *ZipIterator[A, B] {
	return newZipIterator[A, B](a, b)
}

// DeZip is similar to `Zip`, but it supports take elements starting from the back of the iterator.
func DeZip[A any, B any](a DeIterator[A], b DeIterator[B]) *ZipDeIterator[A, B] {
	return newZipDeIterator[A, B](a, b)
}

// TryRfold is the reverse version of [`Iterator[T].TryFold()`]: it takes
// elements starting from the back of the iterator.
func TryRfold[T any, B any](iter DeIterator[T], init B, f func(B, T) gust.Result[B]) gust.Result[B] {
	var accum = gust.Ok(init)
	for {
		x := iter.NextBack()
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
func Rfold[T any, B any](iter DeIterator[T], init B, f func(B, T) B) B {
	var accum = init
	for {
		x := iter.NextBack()
		if x.IsNone() {
			return accum
		}
		accum = f(accum, x.Unwrap())
	}
}

// Flatten creates an iterator that flattens nested structure.
func Flatten[T any, D gust.Iterable[T]](iter Iterator[D]) *FlattenIterator[T, D] {
	return newFlattenIterator[T, D](iter)
}

// DeFlatten creates a double ended iterator that flattens nested structure.
func DeFlatten[T any, D gust.DeIterable[T]](iter DeIterator[D]) *FlattenDeIterator[T, D] {
	return newFlattenDeIterator[T, D](iter)
}

// Enumerate creates an iterator that yields pairs of the index and the value.
func Enumerate[T any](iter Iterator[T]) Iterator[KV[T]] {
	return newEnumerateIterator(iter)
}

// DeEnumerate creates a double ended iterator that yields pairs of the index and the value.
func DeEnumerate[T any](iter DeIterator[T]) DeIterator[KV[T]] {
	return newDeEnumerateIterator(iter)
}
