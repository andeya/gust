package iter

import (
	"github.com/andeya/gust"
)

// FromIterable creates an iterator from a gust.Iterable[T].
// If the data is already an Iterator[T], it returns the same iterator.
// If the data is an Iterable[T], it returns an Iterator[T] with the core.
// If the data is a gust.Iterable[T], it returns an Iterator[T] with the iterable wrapper.
//
// # Examples
//
//	var iter = FromIterable(FromSlice([]int{1, 2, 3}))
//	assert.Equal(t, gust.Some(1), iter.Next())
//	assert.Equal(t, gust.Some(2), iter.Next())
//	assert.Equal(t, gust.Some(3), iter.Next())
//	assert.Equal(t, gust.None[int](), iter.Next())
//
//go:inline
func FromIterable[T any](data gust.Iterable[T]) Iterator[T] {
	switch iter := data.(type) {
	case Iterator[T]:
		return iter
	case Iterable[T]:
		return Iterator[T]{iterable: iter}
	default:
		return Iterator[T]{iterable: iterableWrapper[T]{Iterable: iter}}
	}
}

type iterableWrapper[T any] struct {
	gust.Iterable[T]
}

//go:inline
func (iter iterableWrapper[T]) SizeHint() (uint, gust.Option[uint]) {
	return 0, gust.None[uint]()
}

// FromSlice creates an iterator from a slice.
//
// The returned iterator supports double-ended iteration, allowing iteration
// from both ends. Use AsDoubleEnded() to convert to DoubleEndedIterator.
//
// # Examples
//
//	var a = []int{1, 2, 3}
//	var iter = FromSlice(a)
//	assert.Equal(t, gust.Some(1), iter.Next())
//	assert.Equal(t, gust.Some(2), iter.Next())
//	assert.Equal(t, gust.Some(3), iter.Next())
//	assert.Equal(t, gust.None[int](), iter.Next())
//
//	// As DoubleEndedIterator:
//	var deIter = AsDoubleEnded(FromSlice([]int{1, 2, 3, 4, 5, 6}))
//	assert.Equal(t, gust.Some(1), deIter.Next())
//	assert.Equal(t, gust.Some(6), deIter.NextBack())
//	assert.Equal(t, gust.Some(5), deIter.NextBack())
//
//go:inline
func FromSlice[T any](slice []T) Iterator[T] {
	return Iterator[T]{iterable: &sliceIterable[T]{slice: slice, front: 0, back: len(slice)}}
}

type sliceIterable[T any] struct {
	slice []T
	front int // front index (inclusive)
	back  int // back index (exclusive)
}

func (s *sliceIterable[T]) Next() gust.Option[T] {
	if s.front >= s.back {
		return gust.None[T]()
	}
	item := s.slice[s.front]
	s.front++
	return gust.Some(item)
}

func (s *sliceIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	remaining := uint(s.back - s.front)
	return remaining, gust.Some(remaining)
}

func (s *sliceIterable[T]) Remaining() uint {
	return uint(s.back - s.front)
}

// NextBack removes and returns an element from the end of the iterator.
func (s *sliceIterable[T]) NextBack() gust.Option[T] {
	if s.front >= s.back {
		return gust.None[T]()
	}
	s.back--
	item := s.slice[s.back]
	return gust.Some(item)
}

// FromElements creates an iterator from a set of elements.
//
// # Examples
//
//	var iter = FromElements(1, 2, 3)
//	assert.Equal(t, gust.Some(1), iter.Next())
//	assert.Equal(t, gust.Some(2), iter.Next())
//	assert.Equal(t, gust.Some(3), iter.Next())
//	assert.Equal(t, gust.None[int](), iter.Next())
//
//go:inline
func FromElements[T any](elems ...T) Iterator[T] {
	return FromSlice(elems)
}

// FromRange creates an iterator from a range of integers.
//
// The range is [start, end), meaning start is inclusive and end is exclusive.
//
// # Examples
//
//	var iter = FromRange(0, 5)
//	assert.Equal(t, gust.Some(0), iter.Next())
//	assert.Equal(t, gust.Some(1), iter.Next())
//	assert.Equal(t, gust.Some(2), iter.Next())
//	assert.Equal(t, gust.Some(3), iter.Next())
//	assert.Equal(t, gust.Some(4), iter.Next())
//	assert.Equal(t, gust.None[int](), iter.Next())
//
//go:inline
func FromRange[T gust.Integer](start T, end T) Iterator[T] {
	return Iterator[T]{iterable: &rangeIterable[T]{start: start, end: end, current: start}}
}

type rangeIterable[T gust.Integer] struct {
	start   T
	end     T
	current T
}

func (r *rangeIterable[T]) Next() gust.Option[T] {
	if r.current >= r.end {
		return gust.None[T]()
	}
	item := r.current
	r.current++
	return gust.Some(item)
}

func (r *rangeIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	if r.current >= r.end {
		return 0, gust.Some(uint(0))
	}
	remaining := uint(r.end - r.current)
	return remaining, gust.Some(remaining)
}

// FromFunc creates an iterator from a function that generates values.
//
// The function is called repeatedly until it returns gust.None[T]().
//
// # Examples
//
//	var count = 0
//	var iter = FromFunc(func() gust.Option[int] {
//		if count < 3 {
//			count++
//			return gust.Some(count)
//		}
//		return gust.None[int]()
//	})
//	assert.Equal(t, gust.Some(1), iter.Next())
//	assert.Equal(t, gust.Some(2), iter.Next())
//	assert.Equal(t, gust.Some(3), iter.Next())
//	assert.Equal(t, gust.None[int](), iter.Next())
func FromFunc[T any](f func() gust.Option[T]) Iterator[T] {
	return Iterator[T]{iterable: &funcIterable[T]{f: f}}
}

type funcIterable[T any] struct {
	f func() gust.Option[T]
}

func (f *funcIterable[T]) Next() gust.Option[T] {
	return f.f()
}

func (f *funcIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	return DefaultSizeHint[T]()
}

// Once creates an iterator that yields a single value.
//
// # Examples
//
// Once creates an iterator that yields a value exactly once.
//
// # Examples
//
//	var iter = Once(42)
//	assert.Equal(t, gust.Some(42), iter.Next())
//	assert.Equal(t, gust.None[int](), iter.Next())
func Once[T any](value T) Iterator[T] {
	return Iterator[T]{iterable: &onceIterable[T]{value: value, done: false}}
}

type onceIterable[T any] struct {
	value T
	done  bool
}

func (o *onceIterable[T]) Next() gust.Option[T] {
	if o.done {
		return gust.None[T]()
	}
	o.done = true
	return gust.Some(o.value)
}

func (o *onceIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	if o.done {
		return 0, gust.Some(uint(0))
	}
	return 1, gust.Some(uint(1))
}

// Repeat creates an iterator that repeats a value endlessly.
//
// # Examples
//
//	var iter = Repeat(42)
//	assert.Equal(t, gust.Some(42), iter.Next())
//	assert.Equal(t, gust.Some(42), iter.Next())
//	assert.Equal(t, gust.Some(42), iter.Next())
//	// ... continues forever
func Repeat[T any](value T) Iterator[T] {
	return Iterator[T]{iterable: &repeatIterable[T]{value: value}}
}

type repeatIterable[T any] struct {
	value T
}

func (r *repeatIterable[T]) Next() gust.Option[T] {
	return gust.Some(r.value)
}

func (r *repeatIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	// Infinite iterator
	return 0, gust.None[uint]()
}

// Empty creates an iterator that yields no values.
//
// # Examples
//
//	var iter = Empty[int]()
//	assert.Equal(t, gust.None[int](), iter.Next())
func Empty[T any]() Iterator[T] {
	return Iterator[T]{iterable: &emptyIterable[T]{}}
}

type emptyIterable[T any] struct{}

func (e *emptyIterable[T]) Next() gust.Option[T] {
	return gust.None[T]()
}

func (e *emptyIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	return 0, gust.Some(uint(0))
}
