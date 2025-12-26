package iter

import (
	stditer "iter"

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

// FromSeq creates an Iterator[T] from Go's standard iter.Seq[T].
// This allows converting Go standard iterators to gust iterators.
// Returns the iterator and a deferStop function that should be deferred
// to ensure proper cleanup.
//
// Note: While the sequence will automatically clean up when it ends naturally
// (when next() returns false), it is recommended to always use "defer deferStop()"
// to ensure proper cleanup in all cases, including early termination.
//
// # Examples
//
//	// Convert a Go range iterator to gust Iterator
//	seq := func(yield func(int) bool) {
//		for i := 0; i < 5; i++ {
//			if !yield(i) {
//				return
//			}
//		}
//	}
//	iter, deferStop := FromSeq(seq)
//	defer deferStop() // Recommended: ensures cleanup even if iteration ends naturally
//	assert.Equal(t, gust.Some(0), iter.Next())
//	assert.Equal(t, gust.Some(1), iter.Next())
//
//	// Works with Go's standard library iterators
//	iter, deferStop = FromSeq(iter.N(5)) // iter.N(5) returns iter.Seq[int]
//	defer deferStop()
//	assert.Equal(t, gust.Some(0), iter.Next())
func FromSeq[T any](seq stditer.Seq[T]) (Iterator[T], func()) {
	next, stop := stditer.Pull(seq)
	return FromPull(next, stop)
}

// FromSeq2 creates an Iterator[gust.Pair[K, V]] from Go's standard iter.Seq2[K, V].
// This allows converting Go standard key-value iterators to gust pair iterators.
// Returns the iterator and a deferStop function that should be deferred
// to ensure proper cleanup.
//
// Note: While the sequence will automatically clean up when it ends naturally
// (when next() returns false), it is recommended to always use "defer deferStop()"
// to ensure proper cleanup in all cases, including early termination.
//
// # Examples
//
//	// Convert a Go map iterator to gust Iterator
//	m := map[string]int{"a": 1, "b": 2, "c": 3}
//	seq2 := func(yield func(string, int) bool) {
//		for k, v := range m {
//			if !yield(k, v) {
//				return
//			}
//		}
//	}
//	iter, deferStop := FromSeq2(seq2)
//	defer deferStop()
//	pair := iter.Next()
//	assert.True(t, pair.IsSome())
//	assert.Contains(t, []string{"a", "b", "c"}, pair.Unwrap().A)
//
//	// Works with Go's standard library iterators
//	iter, deferStop = FromSeq2(maps.All(myMap)) // maps.All returns iter.Seq2[K, V]
//	defer deferStop()
func FromSeq2[K any, V any](seq stditer.Seq2[K, V]) (Iterator[gust.Pair[K, V]], func()) {
	return FromPull2(stditer.Pull2(seq))
}

// FromPull creates an Iterator[T] from a pull-style iterator (next and stop functions).
// This allows converting pull-style iterators to gust iterators.
// Returns the iterator and a deferStop function that should be deferred
// to ensure proper cleanup.
//
// Note: While the sequence will automatically clean up when it ends naturally
// (when next() returns false), it is recommended to always use "defer deferStop()"
// to ensure proper cleanup in all cases, including early termination.
//
// # Examples
//
//	// Convert a pull-style iterator to gust Iterator
//	next, stop := iter.Pull(someSeq)
//	defer stop()
//	gustIter, deferStop := FromPull(next, stop)
//	defer deferStop()
//	result := gustIter.Filter(func(x int) bool { return x > 2 }).Collect()
//
//	// Works with custom pull-style iterators
//	customNext := func() (int, bool) {
//		// custom implementation
//		return 0, false
//	}
//	customStop := func() {}
//	gustIter, deferStop = FromPull(customNext, customStop)
//	defer deferStop()
func FromPull[T any](next func() (T, bool), stop func()) (Iterator[T], func()) {
	return Iterator[T]{iterable: &pullIterable[T]{next: next}}, stop
}

// FromPull2 creates an Iterator[gust.Pair[K, V]] from a pull-style iterator (next and stop functions).
// This allows converting pull-style key-value iterators to gust pair iterators.
// Returns the iterator and a deferStop function that should be deferred
// to ensure proper cleanup.
//
// Note: While the sequence will automatically clean up when it ends naturally
// (when next() returns false), it is recommended to always use "defer deferStop()"
// to ensure proper cleanup in all cases, including early termination.
//
// # Examples
//
//	// Convert a pull-style iterator to gust Iterator
//	next, stop := iter.Pull2(someSeq2)
//	defer stop()
//	gustIter, deferStop := FromPull2(next, stop)
//	defer deferStop()
//	result := gustIter.Filter(func(p gust.Pair[int, string]) bool {
//		return p.B != ""
//	}).Collect()
//
//	// Works with custom pull-style iterators
//	customNext := func() (int, string, bool) {
//		// custom implementation
//		return 0, "", false
//	}
//	customStop := func() {}
//	gustIter, deferStop = FromPull2(customNext, customStop)
//	defer deferStop()
func FromPull2[K any, V any](next func() (K, V, bool), stop func()) (Iterator[gust.Pair[K, V]], func()) {
	return Iterator[gust.Pair[K, V]]{iterable: &pull2Iterable[K, V]{next: next}}, stop
}

// Seq2 converts the Iterator[gust.Pair[K, V]] to Go's standard iter.Seq2[K, V].
// This allows using gust pair iterators with Go's built-in key-value iteration support.
//
// # Examples
//
//	// Convert Zip iterator to Go Seq2
//	iter1 := FromSlice([]int{1, 2, 3})
//	iter2 := FromSlice([]string{"a", "b", "c"})
//	zipped := Zip(iter1, iter2)
//	for k, v := range Seq2(zipped) {
//		fmt.Println(k, v) // prints 1 a, 2 b, 3 c
//	}
//
//	// Works with Go's standard library functions
//	enumerated := Enumerate(FromSlice([]string{"a", "b", "c"}))
//	for idx, val := range Seq2(enumerated) {
//		fmt.Println(idx, val) // prints 0 a, 1 b, 2 c
//	}
func Seq2[K any, V any](it Iterator[gust.Pair[K, V]]) stditer.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for {
			opt := it.Next()
			if opt.IsNone() {
				return
			}
			pair := opt.Unwrap()
			if !yield(pair.A, pair.B) {
				return
			}
		}
	}
}

// Pull2 converts the Iterator[gust.Pair[K, V]] to a pull-style iterator using Go's standard iter.Pull2.
// This returns two functions: next (to pull key-value pairs) and stop (to stop iteration).
// The caller should defer stop() to ensure proper cleanup.
//
// # Examples
//
//	iter1 := FromSlice([]int{1, 2, 3})
//	iter2 := FromSlice([]string{"a", "b", "c"})
//	zipped := Zip(iter1, iter2)
//	next, stop := Pull2(zipped)
//	defer stop()
//
//	// Pull key-value pairs manually
//	for {
//		k, v, ok := next()
//		if !ok {
//			break
//		}
//		fmt.Println(k, v)
//	}
func Pull2[K any, V any](it Iterator[gust.Pair[K, V]]) (next func() (K, V, bool), stop func()) {
	return stditer.Pull2(Seq2(it))
}

type pullIterable[T any] struct {
	next func() (T, bool)
	done bool
}

func (p *pullIterable[T]) Next() gust.Option[T] {
	if p.done {
		return gust.None[T]()
	}

	v, ok := p.next()
	if !ok {
		p.done = true
		return gust.None[T]()
	}
	return gust.Some(v)
}

func (p *pullIterable[T]) SizeHint() (uint, gust.Option[uint]) {
	return DefaultSizeHint[T]()
}

type pull2Iterable[K any, V any] struct {
	next func() (K, V, bool)
	done bool
}

func (p *pull2Iterable[K, V]) Next() gust.Option[gust.Pair[K, V]] {
	if p.done {
		return gust.None[gust.Pair[K, V]]()
	}

	k, v, ok := p.next()
	if !ok {
		p.done = true
		return gust.None[gust.Pair[K, V]]()
	}
	return gust.Some(gust.Pair[K, V]{A: k, B: v})
}

func (p *pull2Iterable[K, V]) SizeHint() (uint, gust.Option[uint]) {
	return DefaultSizeHint[gust.Pair[K, V]]()
}
