package iterator

import (
	"github.com/andeya/gust/option"
)

// XMapWhile creates an iterator that both yields elements based on a predicate and maps (any version).
//
//go:inline
func (it Iterator[T]) XMapWhile(predicate func(T) option.Option[any]) Iterator[any] {
	return MapWhile(it, predicate)
}

// MapWhile creates an iterator that both yields elements based on a predicate and maps.
//
//go:inline
func (it Iterator[T]) MapWhile(predicate func(T) option.Option[T]) Iterator[T] {
	return MapWhile(it, predicate)
}

// XScan creates an iterator that scans over the iterator with a state (any version).
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var scanned = iterator.XScan(0, func(state *any, x int) option.Option[any] {
//		s := (*state).(int) + x
//		*state = s
//		return option.Some(any(s))
//	})
//	// Can chain: scanned.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) XScan(initialState any, f func(*any, T) option.Option[any]) Iterator[any] {
	return Scan(it, initialState, f)
}

// Scan creates an iterator that scans over the iterator with a state.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var scanned = iterator.Scan(0, func(state *int, x int) option.Option[int] {
//		*state = *state + x
//		return option.Some(*state)
//	})
//	// Can chain: scanned.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) Scan(initialState T, f func(*T, T) option.Option[T]) Iterator[T] {
	return Scan(it, initialState, f)
}

// MapWhile creates an iterator that both yields elements based on a predicate and maps.
//
// MapWhile() takes a closure as an argument. It will call this
// closure on each element of the iterator, and yield elements
// while it returns option.Some(_).
//
// # Examples
//
// Basic usage:
//
//	var a = []int{-1, 4, 0, 1}
//	var iter = MapWhile(FromSlice(a), func(x int) option.Option[int] {
//		if x != 0 {
//			return option.Some(16 / x)
//		}
//		return option.None[int]()
//	})
//
//	assert.Equal(t, option.Some(-16), iterator.Next())
//	assert.Equal(t, option.Some(4), iterator.Next())
//	assert.Equal(t, option.None[int](), iterator.Next())
func MapWhile[T any, U any](iter Iterator[T], predicate func(T) option.Option[U]) Iterator[U] {
	return Iterator[U]{iterable: &mapWhileIterable[T, U]{iter: iter.iterable, predicate: predicate}}
}

type mapWhileIterable[T any, U any] struct {
	iter      Iterable[T]
	predicate func(T) option.Option[U]
}

func (m *mapWhileIterable[T, U]) Next() option.Option[U] {
	item := m.iter.Next()
	if item.IsNone() {
		return option.None[U]()
	}
	return m.predicate(item.Unwrap())
}

func (m *mapWhileIterable[T, U]) SizeHint() (uint, option.Option[uint]) {
	_, upper := m.iter.SizeHint()
	// MapWhile can reduce the size, but we don't know by how much
	return 0, upper
}

// Scan is an iterator adapter which, like Fold, holds internal state, but
// unlike Fold, produces a new iterator.
//
// Scan() takes two arguments: an initial value which seeds the internal
// state, and a closure with two arguments, the first being a mutable
// reference to the internal state and the second an iterator element.
// The closure can assign to the internal state to share state between
// iterations.
//
// On iteration, the closure will be applied to each element of the
// iterator and the return value from the closure, an option.Option, is
// returned by the Next() method. Thus the closure can return
// option.Some(value) to yield value, or option.None[T]() to end the iteration.
//
// # Examples
//
//	var a = []int{1, 2, 3, 4}
//	var iter = Scan(FromSlice(a), 1, func(state *int, x int) option.Option[int] {
//		*state = *state * x
//		if *state > 6 {
//			return option.None[int]()
//		}
//		return option.Some(-*state)
//	})
//
//	assert.Equal(t, option.Some(-1), iterator.Next())
//	assert.Equal(t, option.Some(-2), iterator.Next())
//	assert.Equal(t, option.Some(-6), iterator.Next())
//	assert.Equal(t, option.None[int](), iterator.Next())
func Scan[T any, U any, St any](iter Iterator[T], initialState St, f func(*St, T) option.Option[U]) Iterator[U] {
	return Iterator[U]{iterable: &scanIterable[T, U, St]{iter: iter.iterable, state: initialState, f: f}}
}

type scanIterable[T any, U any, St any] struct {
	iter  Iterable[T]
	state St
	f     func(*St, T) option.Option[U]
}

func (s *scanIterable[T, U, St]) Next() option.Option[U] {
	item := s.iter.Next()
	if item.IsNone() {
		return option.None[U]()
	}
	return s.f(&s.state, item.Unwrap())
}

func (s *scanIterable[T, U, St]) SizeHint() (uint, option.Option[uint]) {
	// Scan can terminate early, so we can't provide accurate size hint
	_, upper := s.iter.SizeHint()
	return 0, upper
}

// FlatMap creates an iterator that works like map, but flattens nested structure.
//
// The Map adapter is very useful, but only when the closure
// argument produces values. If it produces an iterator instead, there's
// an extra layer of indirection. FlatMap() will remove this extra layer
// on its own.
//
// You can think of FlatMap(f) as the semantic equivalent
// of Mapping, and then Flattening as in Map(f).Flatten().
//
// Another way of thinking about FlatMap(): Map's closure returns
// one item for each element, and FlatMap()'s closure returns an
// iterator for each element.
//
// # Examples
//
//	var words = []string{"alpha", "beta", "gamma"}
//	var iter = FlatMap(FromSlice(words), func(s string) Iterator[rune] {
//		return FromSlice([]rune(s))
//	})
//	var result = Collect[rune](iter)
//	// result contains all characters from all words
func FlatMap[T any, U any](iter Iterator[T], f func(T) Iterator[U]) Iterator[U] {
	return Iterator[U]{iterable: &flatMapIterable[T, U]{iter: iter.iterable, f: f, current: nil}}
}

type flatMapIterable[T any, U any] struct {
	iter    Iterable[T]
	f       func(T) Iterator[U]
	current Iterable[U]
}

func (f *flatMapIterable[T, U]) Next() option.Option[U] {
	for {
		if f.current != nil {
			item := f.current.Next()
			if item.IsSome() {
				return item
			}
			f.current = nil
		}

		item := f.iter.Next()
		if item.IsNone() {
			return option.None[U]()
		}
		it := f.f(item.Unwrap())
		f.current = it.iterable
	}
}

func (f *flatMapIterable[T, U]) SizeHint() (uint, option.Option[uint]) {
	// FlatMap can expand or contract the size, so we can't provide accurate size hint
	return 0, option.None[uint]()
}

// flattenImpl is the internal implementation of Flatten.
//
//go:inline
func flattenImpl[T any](iter Iterable[Iterator[T]]) Iterable[T] {
	return &flattenIterable[T]{iter: iter, current: nil}
}

// Flatten creates an iterator that flattens nested structure.
//
// This is useful when you have an iterator of iterators or an iterator of
// things that can be turned into iterators and you want to remove one
// level of indirection.
//
// # Examples
//
// Basic usage:
//
//	var data = [][]int{{1, 2, 3, 4}, {5, 6}}
//	var iter = Flatten(FromSlice(data))
//	var result = Collect[int](iter)
//	// result is []int{1, 2, 3, 4, 5, 6}
func Flatten[T any](iter Iterator[Iterator[T]]) Iterator[T] {
	return Iterator[T]{iterable: flattenImpl(iter.iterable)}
}

type flattenIterable[T any] struct {
	iter    Iterable[Iterator[T]]
	current Iterable[T]
}

func (f *flattenIterable[T]) Next() option.Option[T] {
	for {
		if f.current != nil {
			item := f.current.Next()
			if item.IsSome() {
				return item
			}
			f.current = nil
		}

		item := f.iter.Next()
		if item.IsNone() {
			return option.None[T]()
		}
		it := item.Unwrap()
		f.current = it.iterable
	}
}

func (f *flattenIterable[T]) SizeHint() (uint, option.Option[uint]) {
	// Flatten can expand or contract the size, so we can't provide accurate size hint
	return 0, option.None[uint]()
}
