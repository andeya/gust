package iter

import (
	"github.com/andeya/gust"
)

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
//	assert.Equal(t, gust.Some(2), iter.Next())
//	assert.Equal(t, gust.Some(4), iter.Next())
//	assert.Equal(t, gust.Some(6), iter.Next())
//	assert.Equal(t, gust.None[int](), iter.Next())
//
// Map creates an iterator which calls a closure on each element.
// This function accepts Iterator[T] and returns Iterator[U] for chainable calls.
func Map[T any, U any](iter Iterator[T], f func(T) U) Iterator[U] {
	return Iterator[U]{iterable: &mapIterable[T, U]{iter: iter.iterable, f: f}}
}

// RetMap creates an iterator which calls a closure on each element and returns a Result[U].
//
// # Examples
//
// iter := RetMap(FromSlice([]string{"1", "2", "3", "NaN"}), strconv.Atoi)

// assert.Equal(t, gust.Some(gust.Ok(1)), iter.Next())
// assert.Equal(t, gust.Some(gust.Ok(2)), iter.Next())
// assert.Equal(t, gust.Some(gust.Ok(3)), iter.Next())
// assert.Equal(t, true, iter.Next().Unwrap().IsErr())
// assert.Equal(t, gust.None[gust.Result[int]](), iter.Next())
//
//go:inline
func RetMap[T any, U any](iter Iterator[T], f func(T) (U, error)) Iterator[gust.Result[U]] {
	return Map(iter, func(t T) gust.Result[U] {
		return gust.Ret(f(t))
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
// assert.Equal(t, gust.Some(gust.Some(newInt(1))), iter.Next())
// assert.Equal(t, gust.Some(gust.Some(newInt(2))), iter.Next())
// assert.Equal(t, gust.Some(gust.Some(newInt(3))), iter.Next())
// assert.Equal(t, gust.Some(gust.None[*int]()), iter.Next())
// assert.Equal(t, gust.None[gust.Option[*int]](), iter.Next())
//
//go:inline
func OptMap[T any, U any](iter Iterator[T], f func(T) *U) Iterator[gust.Option[*U]] {
	return Map(iter, func(t T) gust.Option[*U] {
		return gust.PtrOpt(f(t))
	})
}

type mapIterable[T any, U any] struct {
	iter Iterable[T]
	f    func(T) U
}

func (m *mapIterable[T, U]) Next() gust.Option[U] {
	item := m.iter.Next()
	if item.IsNone() {
		return gust.None[U]()
	}
	return gust.Some(m.f(item.Unwrap()))
}

func (m *mapIterable[T, U]) SizeHint() (uint, gust.Option[uint]) {
	return m.iter.SizeHint()
}

//go:inline
func filterImpl[T any](iter Iterator[T], predicate func(T) bool) Iterator[T] {
	return Iterator[T]{iterable: &filterIterable[T]{iter: iter.iterable, predicate: predicate}}
}

type filterIterable[T any] struct {
	iter      Iterable[T]
	predicate func(T) bool
}

func (f *filterIterable[T]) Next() gust.Option[T] {
	for {
		item := f.iter.Next()
		if item.IsNone() {
			return gust.None[T]()
		}
		if f.predicate(item.Unwrap()) {
			return item
		}
	}
}

func (f *filterIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	_, upper := f.iter.SizeHint()
	// Filter can reduce the size, but we don't know by how much
	return 0, upper
}

// FilterMap creates an iterator that both filters and maps.
//
// The returned iterator yields only the values for which the supplied
// closure returns gust.Some(value).
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
//	var iter = FilterMap(FromSlice(a), func(s string) gust.Option[int] {
//		if v, err := strconv.Atoi(s); err == nil {
//			return gust.Some(v)
//		}
//		return gust.None[int]()
//	})
//
//	assert.Equal(t, gust.Some(1), iter.Next())
//	assert.Equal(t, gust.Some(5), iter.Next())
//	assert.Equal(t, gust.None[int](), iter.Next())
//
// FilterMap creates an iterator that both filters and maps.
// This function accepts Iterator[T] and returns Iterator[U] for chainable calls.
func FilterMap[T any, U any](iter Iterator[T], f func(T) gust.Option[U]) Iterator[U] {
	return Iterator[U]{iterable: &filterMapIterable[T, U]{iter: iter.iterable, f: f}}
}

type filterMapIterable[T any, U any] struct {
	iter Iterable[T]
	f    func(T) gust.Option[U]
}

func (f *filterMapIterable[T, U]) Next() gust.Option[U] {
	for {
		item := f.iter.Next()
		if item.IsNone() {
			return gust.None[U]()
		}
		if result := f.f(item.Unwrap()); result.IsSome() {
			return result
		}
	}
}

func (f *filterMapIterable[T, U]) SizeHint() (uint, gust.Option[uint]) {
	_, upper := f.iter.SizeHint()
	// FilterMap can reduce the size, but we don't know by how much
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

func (c *chainIterable[T]) Next() gust.Option[T] {
	if c.useA {
		item := c.a.Next()
		if item.IsSome() {
			return item
		}
		c.useA = false
	}
	return c.b.Next()
}

func (c *chainIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	lowerA, upperA := c.a.SizeHint()
	lowerB, upperB := c.b.SizeHint()

	lower := lowerA + lowerB

	var upper gust.Option[uint]
	if upperA.IsSome() && upperB.IsSome() {
		upper = gust.Some(upperA.Unwrap() + upperB.Unwrap())
	} else {
		upper = gust.None[uint]()
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
// If either iterator returns gust.None[T](), Next() from the zipped iterator
// will return gust.None[T]().
//
// # Examples
//
// Basic usage:
//
//	var s1 = FromSlice([]rune{'a', 'b', 'c'})
//	var s2 = FromSlice([]rune{'d', 'e', 'f'})
//	var iter = Zip(s1, s2)
//
//	assert.Equal(t, gust.Some(gust.Pair[rune, rune]{A: 'a', B: 'd'}), iter.Next())
//	assert.Equal(t, gust.Some(gust.Pair[rune, rune]{A: 'b', B: 'e'}), iter.Next())
//	assert.Equal(t, gust.Some(gust.Pair[rune, rune]{A: 'c', B: 'f'}), iter.Next())
//	assert.Equal(t, gust.None[gust.Pair[rune, rune]](), iter.Next())
//
// Zip creates an iterator that zips two iterators together.
// This function accepts Iterator[T] and Iterator[U] and returns Iterator[gust.Pair[T, U]] for chainable calls.
func Zip[T any, U any](a Iterator[T], b Iterator[U]) Iterator[gust.Pair[T, U]] {
	return Iterator[gust.Pair[T, U]]{iterable: &zipIterable[T, U]{a: a.iterable, b: b.iterable}}
}

type zipIterable[T any, U any] struct {
	a Iterable[T]
	b Iterable[U]
}

func (z *zipIterable[T, U]) Next() gust.Option[gust.Pair[T, U]] {
	itemA := z.a.Next()
	if itemA.IsNone() {
		return gust.None[gust.Pair[T, U]]()
	}

	itemB := z.b.Next()
	if itemB.IsNone() {
		return gust.None[gust.Pair[T, U]]()
	}

	return gust.Some(gust.Pair[T, U]{A: itemA.Unwrap(), B: itemB.Unwrap()})
}

func (z *zipIterable[T, U]) SizeHint() (uint, gust.Option[uint]) {
	lowerA, upperA := z.a.SizeHint()
	lowerB, upperB := z.b.SizeHint()

	lower := lowerA
	if lowerB < lower {
		lower = lowerB
	}

	var upper gust.Option[uint]
	if upperA.IsSome() && upperB.IsSome() {
		upperAVal := upperA.Unwrap()
		upperBVal := upperB.Unwrap()
		if upperAVal < upperBVal {
			upper = gust.Some(upperAVal)
		} else {
			upper = gust.Some(upperBVal)
		}
	} else if upperA.IsSome() {
		upper = upperA
	} else if upperB.IsSome() {
		upper = upperB
	} else {
		upper = gust.None[uint]()
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
//	assert.Equal(t, gust.Some(gust.Pair[uint, rune]{A: 0, B: 'a'}), iter.Next())
//	assert.Equal(t, gust.Some(gust.Pair[uint, rune]{A: 1, B: 'b'}), iter.Next())
//	assert.Equal(t, gust.Some(gust.Pair[uint, rune]{A: 2, B: 'c'}), iter.Next())
//	assert.Equal(t, gust.None[gust.Pair[uint, rune]](), iter.Next())
//
// enumerateImpl is the internal implementation of Enumerate.
//
//go:inline
func enumerateImpl[T any](iter Iterable[T]) Iterable[gust.Pair[uint, T]] {
	return &enumerateIterable[T]{iter: iter, count: 0}
}

// Enumerate creates an iterator that yields pairs of (index, value).
// This function accepts Iterator[T] and returns Iterator[gust.Pair[uint, T]] for chainable calls.
func Enumerate[T any](iter Iterator[T]) Iterator[gust.Pair[uint, T]] {
	return Iterator[gust.Pair[uint, T]]{iterable: enumerateImpl(iter.iterable)}
}

type enumerateIterable[T any] struct {
	iter  Iterable[T]
	count uint
}

func (e *enumerateIterable[T]) Next() gust.Option[gust.Pair[uint, T]] {
	item := e.iter.Next()
	if item.IsNone() {
		return gust.None[gust.Pair[uint, T]]()
	}
	pair := gust.Pair[uint, T]{A: e.count, B: item.Unwrap()}
	e.count++
	return gust.Some(pair)
}

func (e *enumerateIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	return e.iter.SizeHint()
}

//go:inline
func skipImpl[T any](iter Iterator[T], n uint) Iterator[T] {
	return Iterator[T]{iterable: &skipIterable[T]{iter: iter.iterable, n: n, done: false}}
}

type skipIterable[T any] struct {
	iter Iterable[T]
	n    uint
	done bool
}

func (s *skipIterable[T]) Next() gust.Option[T] {
	if !s.done {
		advanceByImpl(s.iter, s.n)
		s.done = true
	}
	return s.iter.Next()
}

func (s *skipIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	lower, upper := s.iter.SizeHint()
	if lower >= s.n {
		lower -= s.n
	} else {
		lower = 0
	}
	if upper.IsSome() {
		upperVal := upper.Unwrap()
		if upperVal >= s.n {
			upper = gust.Some(upperVal - s.n)
		} else {
			upper = gust.Some(uint(0))
		}
	}
	return lower, upper
}

//go:inline
func takeImpl[T any](iter Iterator[T], n uint) Iterator[T] {
	return Iterator[T]{iterable: &takeIterable[T]{iter: iter.iterable, n: n, taken: 0}}
}

type takeIterable[T any] struct {
	iter  Iterable[T]
	n     uint
	taken uint
}

func (t *takeIterable[T]) Next() gust.Option[T] {
	if t.taken >= t.n {
		return gust.None[T]()
	}
	item := t.iter.Next()
	if item.IsSome() {
		t.taken++
	}
	return item
}

func (t *takeIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	lower, upper := t.iter.SizeHint()
	if lower > t.n {
		lower = t.n
	}
	if upper.IsSome() && upper.Unwrap() > t.n {
		upper = gust.Some(t.n)
	}
	return lower, upper
}

//go:inline
func stepByImpl[T any](iter Iterator[T], step uint) Iterator[T] {
	if step == 0 {
		panic("StepBy: step must be non-zero")
	}
	return Iterator[T]{iterable: &stepByIterable[T]{iter: iter.iterable, step: step, first: true}}
}

type stepByIterable[T any] struct {
	iter  Iterable[T]
	step  uint
	first bool
}

func (s *stepByIterable[T]) Next() gust.Option[T] {
	if s.first {
		s.first = false
		return s.iter.Next()
	}
	if advanceByImpl(s.iter, s.step-1).IsErr() {
		return gust.None[T]()
	}
	return s.iter.Next()
}

func (s *stepByIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	lower, upper := s.iter.SizeHint()
	if lower > 0 {
		lower = (lower + s.step - 1) / s.step
	}
	if upper.IsSome() {
		upperVal := upper.Unwrap()
		if upperVal > 0 {
			upper = gust.Some((upperVal + s.step - 1) / s.step)
		} else {
			upper = gust.Some(uint(0))
		}
	}
	return lower, upper
}
