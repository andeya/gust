package iter

import (
	"github.com/andeya/gust"
)

//go:inline
func skipWhileImpl[T any](iter Iterator[T], predicate func(T) bool) Iterator[T] {
	return Iterator[T]{iterable: &skipWhileIterable[T]{iter: iter.iterable, predicate: predicate, done: false}}
}

type skipWhileIterable[T any] struct {
	iter      Iterable[T]
	predicate func(T) bool
	done      bool
}

func (s *skipWhileIterable[T]) Next() gust.Option[T] {
	if !s.done {
		for {
			item := s.iter.Next()
			if item.IsNone() {
				s.done = true
				return gust.None[T]()
			}
			if !s.predicate(item.Unwrap()) {
				s.done = true
				return item
			}
		}
	}
	return s.iter.Next()
}

func (s *skipWhileIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	// SkipWhile can reduce the size, but we don't know by how much
	_, upper := s.iter.SizeHint()
	return 0, upper
}

//go:inline
func takeWhileImpl[T any](iter Iterator[T], predicate func(T) bool) Iterator[T] {
	return Iterator[T]{iterable: &takeWhileIterable[T]{iter: iter.iterable, predicate: predicate}}
}

type takeWhileIterable[T any] struct {
	iter      Iterable[T]
	predicate func(T) bool
}

func (t *takeWhileIterable[T]) Next() gust.Option[T] {
	item := t.iter.Next()
	if item.IsNone() {
		return gust.None[T]()
	}
	if t.predicate(item.Unwrap()) {
		return item
	}
	return gust.None[T]()
}

func (t *takeWhileIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	_, upper := t.iter.SizeHint()
	// TakeWhile can reduce the size, but we don't know by how much
	return 0, upper
}

// MapWhile creates an iterator that both yields elements based on a predicate and maps.
//
// MapWhile() takes a closure as an argument. It will call this
// closure on each element of the iterator, and yield elements
// while it returns gust.Some(_).
//
// # Examples
//
// Basic usage:
//
//	var a = []int{-1, 4, 0, 1}
//	var iter = MapWhile(FromSlice(a), func(x int) gust.Option[int] {
//		if x != 0 {
//			return gust.Some(16 / x)
//		}
//		return gust.None[int]()
//	})
//
//	assert.Equal(t, gust.Some(-16), iter.Next())
//	assert.Equal(t, gust.Some(4), iter.Next())
//	assert.Equal(t, gust.None[int](), iter.Next())
func MapWhile[T any, U any](iter Iterator[T], predicate func(T) gust.Option[U]) Iterator[U] {
	return Iterator[U]{iterable: &mapWhileIterable[T, U]{iter: iter.iterable, predicate: predicate}}
}

type mapWhileIterable[T any, U any] struct {
	iter      Iterable[T]
	predicate func(T) gust.Option[U]
}

func (m *mapWhileIterable[T, U]) Next() gust.Option[U] {
	item := m.iter.Next()
	if item.IsNone() {
		return gust.None[U]()
	}
	return m.predicate(item.Unwrap())
}

func (m *mapWhileIterable[T, U]) SizeHint() (uint, gust.Option[uint]) {
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
// iterator and the return value from the closure, an gust.Option, is
// returned by the Next() method. Thus the closure can return
// gust.Some(value) to yield value, or gust.None[T]() to end the iteration.
//
// # Examples
//
//	var a = []int{1, 2, 3, 4}
//	var iter = Scan(FromSlice(a), 1, func(state *int, x int) gust.Option[int] {
//		*state = *state * x
//		if *state > 6 {
//			return gust.None[int]()
//		}
//		return gust.Some(-*state)
//	})
//
//	assert.Equal(t, gust.Some(-1), iter.Next())
//	assert.Equal(t, gust.Some(-2), iter.Next())
//	assert.Equal(t, gust.Some(-6), iter.Next())
//	assert.Equal(t, gust.None[int](), iter.Next())
func Scan[T any, U any, St any](iter Iterator[T], initialState St, f func(*St, T) gust.Option[U]) Iterator[U] {
	return Iterator[U]{iterable: &scanIterable[T, U, St]{iter: iter.iterable, state: initialState, f: f}}
}

type scanIterable[T any, U any, St any] struct {
	iter  Iterable[T]
	state St
	f     func(*St, T) gust.Option[U]
}

func (s *scanIterable[T, U, St]) Next() gust.Option[U] {
	item := s.iter.Next()
	if item.IsNone() {
		return gust.None[U]()
	}
	return s.f(&s.state, item.Unwrap())
}

func (s *scanIterable[T, U, St]) SizeHint() (uint, gust.Option[uint]) {
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

func (f *flatMapIterable[T, U]) Next() gust.Option[U] {
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
			return gust.None[U]()
		}
		iter := f.f(item.Unwrap())
		f.current = iter.iterable
	}
}

func (f *flatMapIterable[T, U]) SizeHint() (uint, gust.Option[uint]) {
	// FlatMap can expand or contract the size, so we can't provide accurate size hint
	return 0, gust.None[uint]()
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
//
// flattenImpl is the internal implementation of Flatten.
//
//go:inline
func flattenImpl[T any](iter Iterable[Iterator[T]]) Iterable[T] {
	return &flattenIterable[T]{iter: iter, current: nil}
}

func Flatten[T any](iter Iterator[Iterator[T]]) Iterator[T] {
	return Iterator[T]{iterable: flattenImpl(iter.iterable)}
}

type flattenIterable[T any] struct {
	iter    Iterable[Iterator[T]]
	current Iterable[T]
}

func (f *flattenIterable[T]) Next() gust.Option[T] {
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
			return gust.None[T]()
		}
		iter := item.Unwrap()
		f.current = iter.iterable
	}
}

func (f *flattenIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	// Flatten can expand or contract the size, so we can't provide accurate size hint
	return 0, gust.None[uint]()
}

//go:inline
func fuseImpl[T any](iter Iterator[T]) Iterator[T] {
	return Iterator[T]{iterable: &fuseIterable[T]{iter: iter.iterable, done: false}}
}

type fuseIterable[T any] struct {
	iter Iterable[T]
	done bool
}

func (f *fuseIterable[T]) Next() gust.Option[T] {
	if f.done {
		return gust.None[T]()
	}
	item := f.iter.Next()
	if item.IsNone() {
		f.done = true
	}
	return item
}

func (f *fuseIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	return f.iter.SizeHint()
}

//go:inline
func inspectImpl[T any](iter Iterator[T], f func(T)) Iterator[T] {
	return Iterator[T]{iterable: &inspectIterable[T]{iter: iter.iterable, f: f}}
}

type inspectIterable[T any] struct {
	iter Iterable[T]
	f    func(T)
}

func (i *inspectIterable[T]) Next() gust.Option[T] {
	item := i.iter.Next()
	if item.IsSome() {
		i.f(item.Unwrap())
	}
	return item
}

func (i *inspectIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	return i.iter.SizeHint()
}

//go:inline
func intersperseImpl[T any](iter Iterable[T], separator T) Iterator[T] {
	return intersperseWithImpl(iter, func() T { return separator })
}

//go:inline
func intersperseWithImpl[T any](iter Iterable[T], separator func() T) Iterator[T] {
	return Iterator[T]{iterable: &intersperseIterable[T]{iter: iter, separator: separator, peeked: gust.None[T](), state: intersperseStateFirst}}
}

type intersperseState int

const (
	intersperseStateFirst intersperseState = iota
	intersperseStateItem
	intersperseStateSeparator
)

type intersperseIterable[T any] struct {
	iter      Iterable[T]
	separator func() T
	peeked    gust.Option[T]
	state     intersperseState
}

func (i *intersperseIterable[T]) Next() gust.Option[T] {
	switch i.state {
	case intersperseStateFirst:
		item := i.iter.Next()
		if item.IsNone() {
			return gust.None[T]()
		}
		i.peeked = i.iter.Next()
		if i.peeked.IsNone() {
			i.state = intersperseStateItem
			return item
		}
		i.state = intersperseStateSeparator
		return item
	case intersperseStateSeparator:
		i.state = intersperseStateItem
		return gust.Some(i.separator())
	case intersperseStateItem:
		item := i.peeked
		i.peeked = i.iter.Next()
		if i.peeked.IsNone() {
			return item
		}
		i.state = intersperseStateSeparator
		return item
	}
	return gust.None[T]()
}

func (i *intersperseIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	lower, upper := i.iter.SizeHint()
	if lower > 0 {
		lower = lower*2 - 1
	}
	if upper.IsSome() {
		upperVal := upper.Unwrap()
		if upperVal > 0 {
			upper = gust.Some(upperVal*2 - 1)
		}
	}
	return lower, upper
}

// ArrayChunks creates an iterator that yields arrays of a fixed size.
//
// The chunks are arrays and do not overlap. If the iterator length is not
// divisible by the chunk size, the last chunk will be shorter.
//
// # Panics
//
// Panics if chunk_size is 0.
//
// # Examples
//
//	var iter = ArrayChunks(FromSlice([]int{1, 2, 3, 4, 5, 6}), 2)
//	chunk1 := iter.Next()
//	assert.True(t, chunk1.IsSome())
//	assert.Equal(t, []int{1, 2}, chunk1.Unwrap())
//
//	chunk2 := iter.Next()
//	assert.True(t, chunk2.IsSome())
//	assert.Equal(t, []int{3, 4}, chunk2.Unwrap())
//
//	chunk3 := iter.Next()
//	assert.True(t, chunk3.IsSome())
//	assert.Equal(t, []int{5, 6}, chunk3.Unwrap())
//
//	assert.True(t, iter.Next().IsNone())
//
// arrayChunksImpl is the internal implementation of ArrayChunks.
//
//go:inline
func arrayChunksImpl[T any](iter Iterable[T], chunkSize uint) Iterable[[]T] {
	return &arrayChunksIterable[T]{iter: iter, chunkSize: chunkSize, buffer: make([]T, 0, chunkSize)}
}

func ArrayChunks[T any](iter Iterator[T], chunkSize uint) Iterator[[]T] {
	if chunkSize == 0 {
		panic("ArrayChunks: chunk_size must be non-zero")
	}
	return Iterator[[]T]{iterable: arrayChunksImpl(iter.iterable, chunkSize)}
}

type arrayChunksIterable[T any] struct {
	iter      Iterable[T]
	chunkSize uint
	buffer    []T
}

func (a *arrayChunksIterable[T]) Next() gust.Option[[]T] {
	// Clear buffer but keep capacity
	a.buffer = a.buffer[:0]

	for i := uint(0); i < a.chunkSize; i++ {
		item := a.iter.Next()
		if item.IsNone() {
			if len(a.buffer) == 0 {
				return gust.None[[]T]()
			}
			// Return partial chunk
			result := make([]T, len(a.buffer))
			copy(result, a.buffer)
			return gust.Some(result)
		}
		a.buffer = append(a.buffer, item.Unwrap())
	}

	// Full chunk
	result := make([]T, len(a.buffer))
	copy(result, a.buffer)
	return gust.Some(result)
}

func (a *arrayChunksIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	lower, upper := a.iter.SizeHint()
	if lower > 0 {
		lower = (lower + a.chunkSize - 1) / a.chunkSize
	}
	if upper.IsSome() {
		upperVal := upper.Unwrap()
		if upperVal > 0 {
			upper = gust.Some((upperVal + a.chunkSize - 1) / a.chunkSize)
		}
	}
	return lower, upper
}

// ChunkBy creates an iterator that groups consecutive elements that are equal
// according to the predicate function.
//
// The predicate function receives two consecutive elements and should return
// true if they should be in the same group.
//
// # Examples
//
//	var iter = ChunkBy(FromSlice([]int{1, 1, 2, 2, 2, 3, 3}), func(a, b int) bool { return a == b })
//	chunk1 := iter.Next()
//	assert.True(t, chunk1.IsSome())
//	assert.Equal(t, []int{1, 1}, chunk1.Unwrap())
//
//	chunk2 := iter.Next()
//	assert.True(t, chunk2.IsSome())
//	assert.Equal(t, []int{2, 2, 2}, chunk2.Unwrap())
//
//	chunk3 := iter.Next()
//	assert.True(t, chunk3.IsSome())
//	assert.Equal(t, []int{3, 3}, chunk3.Unwrap())
//
//	assert.True(t, iter.Next().IsNone())
//
// chunkByImpl is the internal implementation of ChunkBy.
//
//go:inline
func chunkByImpl[T any](iter Iterable[T], predicate func(T, T) bool) Iterable[[]T] {
	return &chunkByIterable[T]{iter: iter, predicate: predicate, first: true, prev: gust.None[T](), current: []T{}}
}

func ChunkBy[T any](iter Iterator[T], predicate func(T, T) bool) Iterator[[]T] {
	return Iterator[[]T]{iterable: chunkByImpl(iter.iterable, predicate)}
}

type chunkByIterable[T any] struct {
	iter      Iterable[T]
	predicate func(T, T) bool
	first     bool
	prev      gust.Option[T]
	current   []T
}

func (c *chunkByIterable[T]) Next() gust.Option[[]T] {
	if c.first {
		item := c.iter.Next()
		if item.IsNone() {
			return gust.None[[]T]()
		}
		c.prev = item
		c.current = []T{item.Unwrap()}
		c.first = false
	}

	for {
		item := c.iter.Next()
		if item.IsNone() {
			if len(c.current) == 0 {
				return gust.None[[]T]()
			}
			result := make([]T, len(c.current))
			copy(result, c.current)
			c.current = []T{}
			return gust.Some(result)
		}

		itemVal := item.Unwrap()
		prevVal := c.prev.Unwrap()

		if c.predicate(prevVal, itemVal) {
			c.current = append(c.current, itemVal)
			c.prev = item
		} else {
			result := make([]T, len(c.current))
			copy(result, c.current)
			c.current = []T{itemVal}
			c.prev = item
			return gust.Some(result)
		}
	}
}

func (c *chunkByIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	// We can't know the exact size without iterating
	return 0, gust.None[uint]()
}

// MapWindows creates an iterator that applies a function to overlapping windows
// of a given size.
//
// The windows are arrays of a fixed size that overlap. The first window contains
// the first N elements, the second window contains elements [1..N+1], etc.
//
// # Panics
//
// Panics if window_size is 0 or if the iterator has fewer than window_size elements.
//
// # Examples
//
//	var iter = MapWindows(FromSlice([]int{1, 2, 3, 4, 5}), 3, func(window []int) int {
//		return window[0] + window[1] + window[2]
//	})
//	assert.Equal(t, gust.Some(6), iter.Next())  // 1+2+3
//	assert.Equal(t, gust.Some(9), iter.Next())  // 2+3+4
//	assert.Equal(t, gust.Some(12), iter.Next()) // 3+4+5
//	assert.Equal(t, gust.None[int](), iter.Next())
func MapWindows[T any, U any](iter Iterator[T], windowSize uint, f func([]T) U) Iterator[U] {
	if windowSize == 0 {
		panic("MapWindows: window_size must be non-zero")
	}
	return Iterator[U]{iterable: &mapWindowsIterable[T, U]{iter: iter.iterable, windowSize: windowSize, f: f, buffer: make([]T, 0, windowSize)}}
}

type mapWindowsIterable[T any, U any] struct {
	iter       Iterable[T]
	windowSize uint
	f          func([]T) U
	buffer     []T
}

func (m *mapWindowsIterable[T, U]) Next() gust.Option[U] {
	// Fill buffer to windowSize
	for uint(len(m.buffer)) < m.windowSize {
		item := m.iter.Next()
		if item.IsNone() {
			return gust.None[U]()
		}
		m.buffer = append(m.buffer, item.Unwrap())
	}

	// Apply function to current window
	result := m.f(m.buffer)

	// Shift window: remove first element
	m.buffer = m.buffer[1:]

	return gust.Some(result)
}

func (m *mapWindowsIterable[T, U]) SizeHint() (uint, gust.Option[uint]) {
	lower, upper := m.iter.SizeHint()
	if lower >= m.windowSize {
		lower = lower - m.windowSize + 1
	} else {
		lower = 0
	}
	if upper.IsSome() {
		upperVal := upper.Unwrap()
		if upperVal >= m.windowSize {
			upper = gust.Some(upperVal - m.windowSize + 1)
		} else {
			upper = gust.Some(uint(0))
		}
	}
	return lower, upper
}
