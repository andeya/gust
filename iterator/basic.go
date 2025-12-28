package iterator

import (
	"github.com/andeya/gust/option"
	"github.com/andeya/gust/pair"
	"github.com/andeya/gust/result"
)

// XMap creates an iterator which calls a closure on each element (any version).
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var doubled = iterator.XMap(func(x int) any { return x * 2 })
//	assert.Equal(t, option.Some(2), doubled.Next())
//	// Can chain: doubled.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) XMap(f func(T) any) Iterator[any] {
	return Map(it, f)
}

// Map creates an iterator which calls a closure on each element.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var doubled = iterator.Map(func(x int) int { return x * 2 })
//	assert.Equal(t, option.Some(2), doubled.Next())
//	// Can chain: doubled.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) Map(f func(T) T) Iterator[T] {
	return Map(it, f)
}

// Filter creates an iterator which uses a closure to determine if an element should be yielded.
//
// # Examples
//
//	var iter = FromSlice([]int{0, 1, 2})
//	var filtered = iterator.Filter(func(x int) bool { return x > 0 })
//	assert.Equal(t, option.Some(1), filtered.Next())
//
//go:inline
func (it Iterator[T]) Filter(predicate func(T) bool) Iterator[T] {
	return filterImpl(it, predicate)
}

// Chain takes two iterators and creates a new iterator over both in sequence.
//
// # Examples
//
//	var iter1 = FromSlice([]int{1, 2, 3})
//	var iter2 = FromSlice([]int{4, 5, 6})
//	var chained = iter1.Chain(iter2)
//	assert.Equal(t, option.Some(1), chained.Next())
//
//go:inline
func (it Iterator[T]) Chain(other Iterator[T]) Iterator[T] {
	return chainImpl(it, other)
}

// XFilterMap creates an iterator that both filters and maps (any version).
//
// # Examples
//
//	var iter = FromSlice([]string{"1", "two", "NaN", "four", "5"})
//	var filtered = iterator.XFilterMap(func(s string) option.Option[any] {
//		if s == "1" {
//			return option.Some(any(1))
//		}
//		if s == "5" {
//			return option.Some(any(5))
//		}
//		return option.None[any]()
//	})
//	// Can chain: filtered.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) XFilterMap(f func(T) option.Option[any]) Iterator[any] {
	return FilterMap(it, f)
}

// FilterMap creates an iterator that both filters and maps.
//
// # Examples
//
//	var iter = FromSlice([]string{"1", "two", "NaN", "four", "5"})
//	var filtered = iterator.FilterMap(func(s string) option.Option[string] {
//		if s == "1" || s == "5" {
//			return option.Some(s)
//		}
//		return option.None[string]()
//	})
//	// Can chain: filtered.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) FilterMap(f func(T) option.Option[T]) Iterator[T] {
	return FilterMap(it, f)
}

// XFlatMap creates an iterator that maps and flattens nested iterators (any version).
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var flatMapped = iterator.XFlatMap(func(x int) Iterator[any] {
//		return FromSlice([]any{x, x * 2})
//	})
//	// Can chain: flatMapped.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) XFlatMap(f func(T) Iterator[any]) Iterator[any] {
	return FlatMap(it, f)
}

// FlatMap creates an iterator that maps and flattens nested iterators.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var flatMapped = iterator.FlatMap(func(x int) Iterator[int] {
//		return FromSlice([]int{x, x * 2})
//	})
//	// Can chain: flatMapped.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) FlatMap(f func(T) Iterator[T]) Iterator[T] {
	return FlatMap(it, f)
}

// Map creates an iterator which calls a closure on each element.
//
// Map() transforms one iterator into another, by means of its argument:
// something that implements a function. It produces a new iterator which
// calls this closure on each element of the original iterator.
//
// If you are good at thinking in types, you can think of Map() like this:
// If you have an iterator that gives you elements of some type T, and
// you want an iterator of some other type U, you can use Map(),
// passing a closure that takes a T and returns a U.
//
// Map() is conceptually similar to a for loop. However, as Map() is
// lazy, it is best used when you're already working with other iterators.
// If you're doing some sort of looping for a side effect, it's considered
// more idiomatic to use for than Map().
//
// # Examples
//
// Basic usage:
//
//	var a = []int{1, 2, 3}
//	var iter = Map(FromSlice(a), func(x int) int { return 2 * x })
//
//	assert.Equal(t, option.Some(2), iterator.Next())
//	assert.Equal(t, option.Some(4), iterator.Next())
//	assert.Equal(t, option.Some(6), iterator.Next())
//	assert.Equal(t, option.None[int](), iterator.Next())
//
// Map creates an iterator which calls a closure on each element.
// This function accepts Iterator[T] and returns Iterator[U] for chainable calls.
func Map[T any, U any](iter Iterator[T], f func(T) U) Iterator[U] {
	return Iterator[U]{iterable: &mapIterable[T, U]{iter: iter.iterable, f: f}}
}

type mapIterable[T any, U any] struct {
	iter Iterable[T]
	f    func(T) U
}

func (m *mapIterable[T, U]) Next() option.Option[U] {
	item := m.iter.Next()
	if item.IsNone() {
		return option.None[U]()
	}
	return option.Some(m.f(item.Unwrap()))
}

func (m *mapIterable[T, U]) SizeHint() (uint, option.Option[uint]) {
	return m.iter.SizeHint()
}

// RetMap creates an iterator which calls a closure on each element and returns a Result[U].
//
// # Examples
//
// iter := RetMap(FromSlice([]string{"1", "2", "3", "NaN"}), strconv.Atoi)
//
// assert.Equal(t, option.Some(result.Ok(1)), iterator.Next())
// assert.Equal(t, option.Some(result.Ok(2)), iterator.Next())
// assert.Equal(t, option.Some(result.Ok(3)), iterator.Next())
// assert.Equal(t, true, iterator.Next().Unwrap().IsErr())
// assert.Equal(t, option.None[result.Result[int]](), iterator.Next())
//
//go:inline
func RetMap[T any, U any](iter Iterator[T], f func(T) (U, error)) Iterator[result.Result[U]] {
	return Map(iter, func(t T) result.Result[U] {
		return result.Ret(f(t))
	})
}

// OptMap creates an iterator which calls a closure on each element and returns a Option[*U].
// NOTE:
//
//	`non-nil pointer` is wrapped as Some,
//	and `nil pointer` is wrapped as None.
//
// # Examples
//
//	iter := OptMap(FromSlice([]string{"1", "2", "3", "NaN"}), func(s string) *int {
//		if v, err := strconv.Atoi(s); err == nil {
//			return &v
//		} else {
//			return nil
//		}
//	})
//
//	var newInt = func(v int) *int {
//		return &v
//	}
//
// assert.Equal(t, option.Some(option.Some(newInt(1))), iterator.Next())
// assert.Equal(t, option.Some(option.Some(newInt(2))), iterator.Next())
// assert.Equal(t, option.Some(option.Some(newInt(3))), iterator.Next())
// assert.Equal(t, option.Some(option.None[*int]()), iterator.Next())
// assert.Equal(t, option.None[option.Option[*int]](), iterator.Next())
//
//go:inline
func OptMap[T any, U any](iter Iterator[T], f func(T) *U) Iterator[option.Option[*U]] {
	return Map(iter, func(t T) option.Option[*U] {
		return option.PtrOpt(f(t))
	})
}

// FilterMap creates an iterator that both filters and maps.
//
// The returned iterator yields only the values for which the supplied
// closure returns option.Some(value).
//
// FilterMap can be used to make chains of Filter and Map more
// concise. The example below shows how a Map().Filter().Map() can be
// shortened to a single call to FilterMap.
//
// # Examples
//
// Basic usage:
//
//	var a = []string{"1", "two", "NaN", "four", "5"}
//	var iter = FilterMap(FromSlice(a), func(s string) option.Option[int] {
//		if v, err := strconv.Atoi(s); err == nil {
//			return option.Some(v)
//		}
//		return option.None[int]()
//	})
//
//	assert.Equal(t, option.Some(1), iterator.Next())
//	assert.Equal(t, option.Some(5), iterator.Next())
//	assert.Equal(t, option.None[int](), iterator.Next())
//
// FilterMap creates an iterator that both filters and maps.
// This function accepts Iterator[T] and returns Iterator[U] for chainable calls.
func FilterMap[T any, U any](iter Iterator[T], f func(T) option.Option[U]) Iterator[U] {
	return Iterator[U]{iterable: &filterMapIterable[T, U]{iter: iter.iterable, f: f}}
}

type filterMapIterable[T any, U any] struct {
	iter Iterable[T]
	f    func(T) option.Option[U]
}

func (f *filterMapIterable[T, U]) Next() option.Option[U] {
	for {
		item := f.iter.Next()
		if item.IsNone() {
			return option.None[U]()
		}
		if result := f.f(item.Unwrap()); result.IsSome() {
			return result
		}
	}
}

func (f *filterMapIterable[T, U]) SizeHint() (uint, option.Option[uint]) {
	_, upper := f.iter.SizeHint()
	// FilterMap can reduce the size, but we don't know by how much
	return 0, upper
}

//go:inline
func filterImpl[T any](iter Iterator[T], predicate func(T) bool) Iterator[T] {
	return Iterator[T]{iterable: &filterIterable[T]{iter: iter.iterable, predicate: predicate}}
}

type filterIterable[T any] struct {
	iter      Iterable[T]
	predicate func(T) bool
}

func (f *filterIterable[T]) Next() option.Option[T] {
	for {
		item := f.iter.Next()
		if item.IsNone() {
			return option.None[T]()
		}
		if f.predicate(item.Unwrap()) {
			return item
		}
	}
}

func (f *filterIterable[T]) SizeHint() (uint, option.Option[uint]) {
	_, upper := f.iter.SizeHint()
	// Filter can reduce the size, but we don't know by how much
	return 0, upper
}

//go:inline
func chainImpl[T any](a Iterator[T], b Iterator[T]) Iterator[T] {
	return Iterator[T]{iterable: &chainIterable[T]{a: a.iterable, b: b.iterable, useA: true}}
}

type chainIterable[T any] struct {
	a    Iterable[T]
	b    Iterable[T]
	useA bool
}

func (c *chainIterable[T]) Next() option.Option[T] {
	if c.useA {
		item := c.a.Next()
		if item.IsSome() {
			return item
		}
		c.useA = false
	}
	return c.b.Next()
}

func (c *chainIterable[T]) SizeHint() (uint, option.Option[uint]) {
	lowerA, upperA := c.a.SizeHint()
	lowerB, upperB := c.b.SizeHint()

	lower := lowerA + lowerB

	var upper option.Option[uint]
	if upperA.IsSome() && upperB.IsSome() {
		upper = option.Some(upperA.Unwrap() + upperB.Unwrap())
	} else {
		upper = option.None[uint]()
	}

	return lower, upper
}

// Zip 'zips up' two iterators into a single iterator of pairs.
//
// Zip() returns a new iterator that will iterate over two other
// iterators, returning a tuple where the first element comes from the
// first iterator, and the second element comes from the second iterator.
//
// In other words, it zips two iterators together, into a single one.
//
// If either iterator returns option.None[T](), Next() from the zipped iterator
// will return option.None[T]().
//
// # Examples
//
// Basic usage:
//
//	var s1 = FromSlice([]rune{'a', 'b', 'c'})
//	var s2 = FromSlice([]rune{'d', 'e', 'f'})
//	var iter = Zip(s1, s2)
//
//	assert.Equal(t, option.Some(pair.Pair[rune, rune]{A: 'a', B: 'd'}), iterator.Next())
//	assert.Equal(t, option.Some(pair.Pair[rune, rune]{A: 'b', B: 'e'}), iterator.Next())
//	assert.Equal(t, option.Some(pair.Pair[rune, rune]{A: 'c', B: 'f'}), iterator.Next())
//	assert.Equal(t, option.None[pair.Pair[rune, rune]](), iterator.Next())
//
// Zip creates an iterator that zips two iterators together.
// This function accepts Iterator[T] and Iterator[U] and returns Iterator[pair.Pair[T, U]] for chainable calls.
func Zip[T any, U any](a Iterator[T], b Iterator[U]) Iterator[pair.Pair[T, U]] {
	return Iterator[pair.Pair[T, U]]{iterable: &zipIterable[T, U]{a: a.iterable, b: b.iterable}}
}

type zipIterable[T any, U any] struct {
	a Iterable[T]
	b Iterable[U]
}

func (z *zipIterable[T, U]) Next() option.Option[pair.Pair[T, U]] {
	itemA := z.a.Next()
	if itemA.IsNone() {
		return option.None[pair.Pair[T, U]]()
	}

	itemB := z.b.Next()
	if itemB.IsNone() {
		return option.None[pair.Pair[T, U]]()
	}

	return option.Some(pair.Pair[T, U]{A: itemA.Unwrap(), B: itemB.Unwrap()})
}

func (z *zipIterable[T, U]) SizeHint() (uint, option.Option[uint]) {
	lowerA, upperA := z.a.SizeHint()
	lowerB, upperB := z.b.SizeHint()

	lower := lowerA
	if lowerB < lower {
		lower = lowerB
	}

	var upper option.Option[uint]
	if upperA.IsSome() && upperB.IsSome() {
		upperAVal := upperA.Unwrap()
		upperBVal := upperB.Unwrap()
		if upperAVal < upperBVal {
			upper = option.Some(upperAVal)
		} else {
			upper = option.Some(upperBVal)
		}
	} else if upperA.IsSome() {
		upper = upperA
	} else if upperB.IsSome() {
		upper = upperB
	} else {
		upper = option.None[uint]()
	}

	return lower, upper
}

// Enumerate creates an iterator which gives the current iteration count as well as
// the next value.
//
// The iterator returned yields pairs (i, val), where i is the
// current index of iteration and val is the value returned by the
// iterator.
//
// Enumerate() keeps its count as a uint. If you want to count by a
// different sized integer, the Zip function provides similar
// functionality.
//
// # Overflow Behavior
//
// The method does no guarding against overflows, so enumerating more than
// uint elements either produces the wrong result or panics. If
// overflow checks are enabled, a panic is guaranteed.
//
// # Panics
//
// The returned iterator might panic if the to-be-returned index would
// overflow a uint.
//
// # Examples
//
//	var a = []rune{'a', 'b', 'c'}
//	var iter = Enumerate(FromSlice(a))
//
//	assert.Equal(t, option.Some(pair.Pair[uint, rune]{A: 0, B: 'a'}), iterator.Next())
//	assert.Equal(t, option.Some(pair.Pair[uint, rune]{A: 1, B: 'b'}), iterator.Next())
//	assert.Equal(t, option.Some(pair.Pair[uint, rune]{A: 2, B: 'c'}), iterator.Next())
//	assert.Equal(t, option.None[pair.Pair[uint, rune]](), iterator.Next())
//
// enumerateImpl is the internal implementation of Enumerate.
//
//go:inline
func enumerateImpl[T any](iter Iterable[T]) Iterable[pair.Pair[uint, T]] {
	return &enumerateIterable[T]{iter: iter, count: 0}
}

// Enumerate creates an iterator that yields pairs of (index, value).
// This function accepts Iterator[T] and returns Iterator[pair.Pair[uint, T]] for chainable calls.
func Enumerate[T any](iter Iterator[T]) Iterator[pair.Pair[uint, T]] {
	return Iterator[pair.Pair[uint, T]]{iterable: enumerateImpl(iter.iterable)}
}

type enumerateIterable[T any] struct {
	iter  Iterable[T]
	count uint
}

func (e *enumerateIterable[T]) Next() option.Option[pair.Pair[uint, T]] {
	item := e.iter.Next()
	if item.IsNone() {
		return option.None[pair.Pair[uint, T]]()
	}
	pair := pair.Pair[uint, T]{A: e.count, B: item.Unwrap()}
	e.count++
	return option.Some(pair)
}

func (e *enumerateIterable[T]) SizeHint() (uint, option.Option[uint]) {
	return e.iter.SizeHint()
}
