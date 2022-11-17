package iter

import (
	"unicode/utf8"
	"unsafe"

	"github.com/andeya/gust"
)

// FromIterable creates an iterator from an Iterable.
func FromIterable[T any](data gust.Iterable[T]) Iterator[T] {
	iter, _ := data.(Iterator[T])
	if iter != nil {
		return iter
	}
	return fromIterable[T](data)
}

// EnumIterable creates an iterator with index from an Iterable.
func EnumIterable[T any](data gust.Iterable[T]) Iterator[gust.KV[T]] {
	return ToEnumerate[T](FromIterable[T](data))
}

// FromDeIterable creates a double ended iterator from an Iterable.
func FromDeIterable[T any](data gust.DeIterable[T]) DeIterator[T] {
	iter, _ := data.(DeIterator[T])
	if iter != nil {
		return iter
	}
	return fromDeIterable[T](data)
}

// EnumDeIterable creates a double ended iterator with index from an Iterable.
func EnumDeIterable[T any](data gust.DeIterable[T]) DeIterator[gust.KV[T]] {
	return ToDeEnumerate[T](FromDeIterable[T](data))
}

// FromVec creates a double ended iterator from a slice.
func FromVec[T any](slice []T) DeIterator[T] {
	return NewIterableVec(slice).ToDeIterator()
}

// EnumVec creates a double ended iterator with index from a slice.
func EnumVec[T any](slice []T) DeIterator[gust.KV[T]] {
	return ToDeEnumerate[T](FromVec[T](slice))
}

// FromElements creates a double ended iterator from a set of elements.
func FromElements[T any](elems ...T) DeIterator[T] {
	return NewIterableVec(elems).ToDeIterator()
}

// EnumElements creates a double ended iterator with index from a set of elements.
func EnumElements[T any](elems ...T) DeIterator[gust.KV[T]] {
	return ToDeEnumerate[T](FromVec[T](elems))
}

// FromRange creates a double ended iterator from a range.
func FromRange[T gust.Integer](start T, end T, rightClosed ...bool) DeIterator[T] {
	return NewIterableRange[T](start, end, rightClosed...).ToDeIterator()
}

// EnumRange creates a double ended iterator with index from a range.
func EnumRange[T gust.Integer](start T, end T, rightClosed ...bool) DeIterator[gust.KV[T]] {
	return ToDeEnumerate[T](FromRange[T](start, end, rightClosed...))
}

// FromChan creates an iterator from a channel.
func FromChan[T any](c chan T) Iterator[T] {
	return NewIterableChan[T](c).ToIterator()
}

// EnumChan creates an iterator with index from a channel.
func EnumChan[T any](c chan T) Iterator[gust.KV[T]] {
	return ToEnumerate[T](FromChan(c))
}

// FromResult creates a double ended iterator from a result.
func FromResult[T any](ret gust.Result[T]) DeIterator[T] {
	return FromDeIterable[T](ret)
}

// EnumResult creates a double ended iterator with index from a result.
func EnumResult[T any](ret gust.Result[T]) DeIterator[gust.KV[T]] {
	return EnumDeIterable[T](ret)
}

// FromOption creates a double ended iterator from an option.
func FromOption[T any](opt gust.Option[T]) DeIterator[T] {
	return FromDeIterable[T](opt)
}

// EnumOption creates a double ended iterator with index from an option.
func EnumOption[T any](opt gust.Option[T]) DeIterator[gust.KV[T]] {
	return EnumDeIterable[T](opt)
}

// FromString creates a double ended iterator from a string.
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

// EnumString creates a double ended iterator with index from a string.
func EnumString[T ~byte | ~rune](s string) DeIterator[gust.KV[T]] {
	return ToDeEnumerate[T](FromString[T](s))
}

// ToUnique return an iterator adaptor that filters out elements that have
// already been produced once during the iteration. Duplicates
// are detected using hash and equality.
//
// Clones of visited elements are stored in a hash set in the
// iterator.
//
// The iterator is stable, returning the non-duplicate items in the order
// in which they occur in the adapted iterator. In a set of duplicate
// items, the first item encountered is the item retained.
//
// ```
// var data = FromElements(10, 20, 30, 20, 40, 10, 50);
// ToUnique(data).Collect() // [10, 20, 30, 40, 50]
// ```
func ToUnique[T comparable](iter Iterator[T]) Iterator[T] {
	min, _ := iter.SizeHint()
	var set = make(map[T]struct{}, min)
	return iter.ToFilter(func(x T) bool {
		if _, ok := set[x]; ok {
			return false
		} else {
			set[x] = struct{}{}
			return true
		}
	})
}

// ToDeUnique return a double ended iterator adaptor that filters out elements that have
// already been produced once during the iteration. Duplicates
// are detected using hash and equality.
//
// Clones of visited elements are stored in a hash set in the
// iterator.
//
// The iterator is stable, returning the non-duplicate items in the order
// in which they occur in the adapted iterator. In a set of duplicate
// items, the first item encountered is the item retained.
//
// ```
// var data = FromElements(10, 20, 30, 20, 40, 10, 50);
// ToDeUnique(data).Collect() // [10, 20, 30, 40, 50]
// ```
func ToDeUnique[T comparable](iter DeIterator[T]) DeIterator[T] {
	var set = make(map[T]struct{}, iter.Remaining())
	return iter.ToDeFilter(func(x T) bool {
		if _, ok := set[x]; ok {
			return false
		} else {
			set[x] = struct{}{}
			return true
		}
	})
}

// ToUniqueBy return an iterator adaptor that filters out elements that have
// already been produced once during the iteration.
//
// Duplicates are detected by comparing the key they map to
// with the keying function `f` by hash and equality.
// The keys are stored in a hash set in the iterator.
//
// The iterator is stable, returning the non-duplicate items in the order
// in which they occur in the adapted iterator. In a set of duplicate
// items, the first item encountered is the item retained.
//
// ```
// var data = FromElements("a", "bb", "aa", "c", "ccc");
// ToUniqueBy(data, func(s string)int {return len(s)}).Collect() // "a", "bb", "ccc"
// ```
func ToUniqueBy[T any, K comparable](iter Iterator[T], f func(T) K) Iterator[T] {
	min, _ := iter.SizeHint()
	var set = make(map[K]struct{}, min)
	return iter.ToFilter(func(x T) bool {
		k := f(x)
		if _, ok := set[k]; ok {
			return false
		} else {
			set[k] = struct{}{}
			return true
		}
	})
}

// ToDeUniqueBy return an iterator adaptor that filters out elements that have
// already been produced once during the iteration.
//
// Duplicates are detected by comparing the key they map to
// with the keying function `f` by hash and equality.
// The keys are stored in a hash set in the iterator.
//
// The iterator is stable, returning the non-duplicate items in the order
// in which they occur in the adapted iterator. In a set of duplicate
// items, the first item encountered is the item retained.
//
// ```
// var data = FromElements("a", "bb", "aa", "c", "ccc");
// ToDeUniqueBy(data, func(s string)int {return len(s)}).Collect() // "a", "bb", "ccc"
// ```
func ToDeUniqueBy[T any, K comparable](iter DeIterator[T], f func(T) K) DeIterator[T] {
	var set = make(map[K]struct{}, iter.Remaining())
	return iter.ToDeFilter(func(x T) bool {
		k := f(x)
		if _, ok := set[k]; ok {
			return false
		} else {
			set[k] = struct{}{}
			return true
		}
	})
}

// ToMap takes a closure and creates an iterator which calls that closure on each
// element.
//
// If you are good at thinking in types, you can think of `ToMap()` like this:
// If you have an iterator that gives you elements of some type `A`, and
// you want an iterator of some other type `B`, you can use `ToMap()`,
// passing a closure that takes an `A` and returns a `B`.
//
// `ToMap()` is conceptually similar to a [`for`] loop. However, as `ToMap()` is
// lazy, it is best used when you're already working with other iterators.
// If you're doing some sort of looping for a side effect, it's considered
// more idiomatic to use [`for`] than `ToMap()`.
//
// # Examples
//
// Basic usage:
//
// ```
// var a = []int{1, 2, 3};
//
// var iter = FromVec(a).ToMap(func(x)int{ return 2 * x});
//
// assert.Equal(iter.Next(), gust.Some(2));
// assert.Equal(iter.Next(), gust.Some(4));
// assert.Equal(iter.Next(), gust.Some(6));
// assert.Equal(iter.Next(), gust.None[int]());
// ```
func ToMap[T any, B any](iter Iterator[T], f func(T) B) Iterator[B] {
	return newMapIterator(iter, f)
}

// ToDeMap takes a closure and creates a double ended iterator which calls that closure on each
// element.
//
// If you are good at thinking in types, you can think of `ToDeMap()` like this:
// If you have an iterator that gives you elements of some type `A`, and
// you want an iterator of some other type `B`, you can use `ToDeMap()`,
// passing a closure that takes an `A` and returns a `B`.
//
// `ToDeMap()` is conceptually similar to a [`for`] loop. However, as `ToDeMap()` is
// lazy, it is best used when you're already working with other iterators.
// If you're doing some sort of looping for a side effect, it's considered
// more idiomatic to use [`for`] than `ToDeMap()`.
func ToDeMap[T any, B any](iter DeIterator[T], f func(T) B) DeIterator[B] {
	return newDeMapIterator(iter, f)
}

// ToFilterMap creates an iterator that both filters and maps.
//
// The returned iterator yields only the `value`s for which the supplied
// closure returns `gust.Some(value)`.
func ToFilterMap[T any, B any](iter Iterator[T], f func(T) gust.Option[B]) Iterator[B] {
	return newFilterMapIterator[T, B](iter, f)
}

// ToDeFilterMap creates a double ended iterator that both filters and maps.
//
// The returned iterator yields only the `value`s for which the supplied
// closure returns `gust.Some(value)`.
func ToDeFilterMap[T any, B any](iter DeIterator[T], f func(T) gust.Option[B]) DeIterator[B] {
	return newDeFilterMapIterator[T, B](iter, f)
}

// ToZip 'Zips up' two iterators into a single iterator of pairs.
//
// `ToZip()` returns a new iterator that will iterate over two other
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
func ToZip[A any, B any](a Iterator[A], b Iterator[B]) Iterator[gust.Pair[A, B]] {
	return newZipIterator[A, B](a, b)
}

// ToDeZip is similar to `ToZip`, but it supports take elements starting from the back of the iterator.
func ToDeZip[A any, B any](a DeIterator[A], b DeIterator[B]) DeIterator[gust.Pair[A, B]] {
	return newDeZipIterator[A, B](a, b)
}

// ToEnumerate creates an iterator that yields pairs of the index and the value.
func ToEnumerate[T any](iter Iterator[T]) Iterator[gust.KV[T]] {
	return newEnumerateIterator(iter)
}

// ToDeEnumerate creates a double ended iterator that yields pairs of the index and the value.
func ToDeEnumerate[T any](iter DeIterator[T]) DeIterator[gust.KV[T]] {
	return newDeEnumerateIterator(iter)
}

// ToMapWhile creates an iterator that both yields elements based on a predicate and maps.
//
// `ToMapWhile()` takes a closure as an argument. It will call this
// closure on each element of the iterator, and yield elements
// while it returns [`Some`].
func ToMapWhile[T any, B any](iter Iterator[T], predicate func(T) gust.Option[B]) Iterator[B] {
	return newMapWhileIterator[T, B](iter, predicate)
}

// ToScan is an iterator adapter similar to [`Fold`] that holds internal state and
// produces a new iterator.
//
// [`Fold`]: Iterator.Fold
//
// `ToScan()` takes two arguments: an initial value which seeds the internal
// state, and a closure with two arguments, the first being a mutable
// reference to the internal state and the second an iterator element.
// The closure can assign to the internal state to share state between
// iterations.
//
// On iteration, the closure will be applied to each element of the
// iterator and the return value from the closure, an [`Option`], is
// yielded by the iterator.
func ToScan[T any, St any, B any](iter Iterator[T], initialState St, f func(state *St, item T) gust.Option[B]) Iterator[B] {
	return newScanIterator[T, St, B](iter, initialState, f)
}

// ToFlatten creates an iterator that flattens nested structure.
func ToFlatten[I gust.Iterable[T], T any](iter Iterator[I]) Iterator[T] {
	return newFlattenIterator[I, T](iter)
}

// ToDeFlatten creates a double ended iterator that flattens nested structure.
func ToDeFlatten[I gust.DeIterable[T], T any](iter DeIterator[I]) DeIterator[T] {
	return newDeFlattenIterator[I, T](iter)
}

// ToFlatMap creates an iterator that works like map, but flattens nested structure.
//
// The [`ToMap`] adapter is very useful, but only when the closure
// argument produces values. If it produces an iterator instead, there's
// an extra layer of indirection. `ToFlatMap()` will remove this extra layer
// on its own.
//
// You can think of `ToFlatMap(f)` as the semantic equivalent
// of [`ToMap`]ping, and then [`ToFlatten`]ing as in `ToMap(f).ToFlatten()`.
//
// Another way of thinking about `ToFlatMap()`: [`ToMap`]'s closure returns
// one item for each element, and `ToFlatMap()`'s closure returns an
// iterator for each element.
func ToFlatMap[T any, B any](iter Iterator[T], f func(T) Iterator[B]) Iterator[B] {
	return newFlatMapIterator[T, B](iter, f)
}

// ToDeFlatMap creates a double ended iterator that works like map, but flattens nested structure.
//
// The [`ToDeMap`] adapter is very useful, but only when the closure
// argument produces values. If it produces an iterator instead, there's
// an extra layer of indirection. `ToDeFlatMap()` will remove this extra layer
// on its own.
//
// You can think of `ToDeFlatMap(f)` as the semantic equivalent
// of [`ToDeMap`]ping, and then [`ToDeFlatten`]ing as in `ToDeFlatten(ToDeMap(f))`.
//
// Another way of thinking about `ToDeFlatMap()`: [`ToDeMap`]'s closure returns
// one item for each element, and `ToDeFlatMap()`'s closure returns an
// iterator for each element.
func ToDeFlatMap[T any, B any](iter DeIterator[T], f func(T) DeIterator[B]) DeIterator[B] {
	return newDeFlatMapIterator[T, B](iter, f)
}

// FindMap applies function to the elements of data and returns
// the first non-none
//
// `FindMap(iter, f)` is equivalent to `ToFilterMap(iter, f).Next()`.
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

// SigTryFold a data method that applies a function as long as it returns
// successfully, producing a single, final value.
//
// assert.Equal(t, sum, Ok(6));
func SigTryFold[T any, CB any](iter Iterator[T], init CB, f func(CB, T) gust.SigCtrlFlow[CB]) gust.SigCtrlFlow[CB] {
	var accum = gust.SigContinue[CB](init)
	for {
		x := iter.Next()
		if x.IsNone() {
			return accum
		}
		accum = f(accum.UnwrapContinue(), x.Unwrap())
		if accum.IsBreak() {
			return accum
		}
	}
}

// TryFold a data method that applies a function as long as it returns
// successfully, producing a single, final value.
func TryFold[T any, B any, C any](iter Iterator[T], init C, f func(C, T) gust.CtrlFlow[B, C]) gust.CtrlFlow[B, C] {
	var accum = gust.Continue[B, C](init)
	for {
		x := iter.Next()
		if x.IsNone() {
			return accum
		}
		accum = f(accum.UnwrapContinue(), x.Unwrap())
		if accum.IsBreak() {
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
// This operation is sometimes called 'Reduce' or 'inject'.
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

// TryRfold is the reverse version of [`Iterator[T].TryFold()`]: it takes
// elements starting from the back of the iterator.
func TryRfold[T any, CB any](iter DeIterator[T], init CB, f func(CB, T) gust.SigCtrlFlow[CB]) gust.SigCtrlFlow[CB] {
	var accum = gust.SigContinue[CB](init)
	for {
		x := iter.NextBack()
		if x.IsNone() {
			return accum
		}
		accum = f(accum.UnwrapContinue(), x.Unwrap())
		if accum.IsBreak() {
			return accum
		}
	}
}

// Rfold is an iterator method that reduces the iterator elements to a single,
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
