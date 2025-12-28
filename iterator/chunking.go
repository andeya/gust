package iterator

import (
	"github.com/andeya/gust/option"
)

//go:inline
func arrayChunksImpl[T any](iter Iterable[T], chunkSize uint) Iterable[[]T] {
	return &arrayChunksIterable[T]{iter: iter, chunkSize: chunkSize, buffer: make([]T, 0, chunkSize)}
}

// ArrayChunks creates an iterator that yields chunks of a given size.
//
// # Panics
//
// Panics if chunk_size is 0.
//
// # Examples
//
//	var iter = ArrayChunks(FromSlice([]int{1, 2, 3, 4, 5}), 2)
//	chunk1 := iterator.Next()
//	assert.True(t, chunk1.IsSome())
//	assert.Equal(t, []int{1, 2}, chunk1.Unwrap())
//
//	chunk2 := iterator.Next()
//	assert.True(t, chunk2.IsSome())
//	assert.Equal(t, []int{3, 4}, chunk2.Unwrap())
//
//	chunk3 := iterator.Next()
//	assert.True(t, chunk3.IsSome())
//	assert.Equal(t, []int{5}, chunk3.Unwrap()) // Last partial chunk
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

func (a *arrayChunksIterable[T]) Next() option.Option[[]T] {
	// Clear buffer but keep capacity
	a.buffer = a.buffer[:0]

	for i := uint(0); i < a.chunkSize; i++ {
		item := a.iter.Next()
		if item.IsNone() {
			if len(a.buffer) == 0 {
				return option.None[[]T]()
			}
			// Return partial chunk
			result := make([]T, len(a.buffer))
			copy(result, a.buffer)
			return option.Some(result)
		}
		a.buffer = append(a.buffer, item.Unwrap())
	}

	// Full chunk
	result := make([]T, len(a.buffer))
	copy(result, a.buffer)
	return option.Some(result)
}

func (a *arrayChunksIterable[T]) SizeHint() (uint, option.Option[uint]) {
	lower, upper := a.iter.SizeHint()
	if lower > 0 {
		lower = (lower + a.chunkSize - 1) / a.chunkSize
	}
	if upper.IsSome() {
		upperVal := upper.Unwrap()
		if upperVal > 0 {
			upper = option.Some((upperVal + a.chunkSize - 1) / a.chunkSize)
		}
	}
	return lower, upper
}

//go:inline
func chunkByImpl[T any](iter Iterable[T], predicate func(T, T) bool) Iterable[[]T] {
	return &chunkByIterable[T]{iter: iter, predicate: predicate, first: true, prev: option.None[T](), current: []T{}}
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
//	chunk1 := iterator.Next()
//	assert.True(t, chunk1.IsSome())
//	assert.Equal(t, []int{1, 1}, chunk1.Unwrap())
//
//	chunk2 := iterator.Next()
//	assert.True(t, chunk2.IsSome())
//	assert.Equal(t, []int{2, 2, 2}, chunk2.Unwrap())
//
//	chunk3 := iterator.Next()
//	assert.True(t, chunk3.IsSome())
//	assert.Equal(t, []int{3, 3}, chunk3.Unwrap())
//
//	assert.True(t, iterator.Next().IsNone())
func ChunkBy[T any](iter Iterator[T], predicate func(T, T) bool) Iterator[[]T] {
	return Iterator[[]T]{iterable: chunkByImpl(iter.iterable, predicate)}
}

type chunkByIterable[T any] struct {
	iter      Iterable[T]
	predicate func(T, T) bool
	first     bool
	prev      option.Option[T]
	current   []T
}

func (c *chunkByIterable[T]) Next() option.Option[[]T] {
	if c.first {
		item := c.iter.Next()
		if item.IsNone() {
			return option.None[[]T]()
		}
		c.prev = item
		c.current = []T{item.Unwrap()}
		c.first = false
	}

	for {
		item := c.iter.Next()
		if item.IsNone() {
			if len(c.current) == 0 {
				return option.None[[]T]()
			}
			result := make([]T, len(c.current))
			copy(result, c.current)
			c.current = []T{}
			return option.Some(result)
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
			return option.Some(result)
		}
	}
}

func (c *chunkByIterable[T]) SizeHint() (uint, option.Option[uint]) {
	// We can't know the exact size without iterating
	return 0, option.None[uint]()
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
//	assert.Equal(t, option.Some(6), iterator.Next())  // 1+2+3
//	assert.Equal(t, option.Some(9), iterator.Next())  // 2+3+4
//	assert.Equal(t, option.Some(12), iterator.Next()) // 3+4+5
//	assert.Equal(t, option.None[int](), iterator.Next())
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

func (m *mapWindowsIterable[T, U]) Next() option.Option[U] {
	// Fill buffer to windowSize
	for uint(len(m.buffer)) < m.windowSize {
		item := m.iter.Next()
		if item.IsNone() {
			return option.None[U]()
		}
		m.buffer = append(m.buffer, item.Unwrap())
	}

	// Apply function to current window
	result := m.f(m.buffer)

	// Shift window: remove first element
	m.buffer = m.buffer[1:]

	return option.Some(result)
}

func (m *mapWindowsIterable[T, U]) SizeHint() (uint, option.Option[uint]) {
	lower, upper := m.iter.SizeHint()
	if lower >= m.windowSize {
		lower = lower - m.windowSize + 1
	} else {
		lower = 0
	}
	if upper.IsSome() {
		upperVal := upper.Unwrap()
		if upperVal >= m.windowSize {
			upper = option.Some(upperVal - m.windowSize + 1)
		} else {
			upper = option.Some(uint(0))
		}
	}
	return lower, upper
}

// XMapWindows creates an iterator that applies a function to overlapping windows (any version).
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3, 4, 5})
//	var windows = iterator.XMapWindows(3, func(window []int) any {
//		return window[0] + window[1] + window[2]
//	})
//	// Can chain: windows.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) XMapWindows(windowSize uint, f func([]T) any) Iterator[any] {
	return MapWindows(it, windowSize, f)
}

// MapWindows creates an iterator that applies a function to overlapping windows.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3, 4, 5})
//	var windows = iterator.MapWindows(3, func(window []int) int {
//		return window[0] + window[1] + window[2]
//	})
//	// Can chain: windows.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) MapWindows(windowSize uint, f func([]T) T) Iterator[T] {
	return MapWindows(it, windowSize, f)
}
