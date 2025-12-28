package iterator

import (
	"iter"

	"github.com/andeya/gust"
	"github.com/andeya/gust/constraints"
)

// FromIterable creates an iterator from a gust.Iterable[T].
// If the data is already an Iterator[T], it returns the same iterator.
// If the data is an Iterable[T], it returns an Iterator[T] with the core.
// If the data is a gust.Iterable[T], it returns an Iterator[T] with the iterable wrapper.
//
// # Examples
//
//	var iter = FromIterable(FromSlice([]int{1, 2, 3}))
//	assert.Equal(t, gust.Some(1), iterator.Next())
//	assert.Equal(t, gust.Some(2), iterator.Next())
//	assert.Equal(t, gust.Some(3), iterator.Next())
//	assert.Equal(t, gust.None[int](), iterator.Next())
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
//	assert.Equal(t, gust.Some(1), iterator.Next())
//	assert.Equal(t, gust.Some(2), iterator.Next())
//	assert.Equal(t, gust.Some(3), iterator.Next())
//	assert.Equal(t, gust.None[int](), iterator.Next())
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
//	assert.Equal(t, gust.Some(1), iterator.Next())
//	assert.Equal(t, gust.Some(2), iterator.Next())
//	assert.Equal(t, gust.Some(3), iterator.Next())
//	assert.Equal(t, gust.None[int](), iterator.Next())
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
//	assert.Equal(t, gust.Some(0), iterator.Next())
//	assert.Equal(t, gust.Some(1), iterator.Next())
//	assert.Equal(t, gust.Some(2), iterator.Next())
//	assert.Equal(t, gust.Some(3), iterator.Next())
//	assert.Equal(t, gust.Some(4), iterator.Next())
//	assert.Equal(t, gust.None[int](), iterator.Next())
//
//go:inline
func FromRange[T constraints.Integer](start T, end T) Iterator[T] {
	return Iterator[T]{iterable: &rangeIterable[T]{start: start, end: end, current: start}}
}

type rangeIterable[T constraints.Integer] struct {
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
//	assert.Equal(t, gust.Some(1), iterator.Next())
//	assert.Equal(t, gust.Some(2), iterator.Next())
//	assert.Equal(t, gust.Some(3), iterator.Next())
//	assert.Equal(t, gust.None[int](), iterator.Next())
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
//	assert.Equal(t, gust.Some(42), iterator.Next())
//	assert.Equal(t, gust.None[int](), iterator.Next())
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
//	assert.Equal(t, gust.Some(42), iterator.Next())
//	assert.Equal(t, gust.Some(42), iterator.Next())
//	assert.Equal(t, gust.Some(42), iterator.Next())
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
//	assert.Equal(t, gust.None[int](), iterator.Next())
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

// FromSeq creates an Iterator[T] from Go's standard iterator.Seq[T].
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
//	assert.Equal(t, gust.Some(0), iterator.Next())
//	assert.Equal(t, gust.Some(1), iterator.Next())
//
//	// Works with Go's standard library iterators
//	iter, deferStop = FromSeq(iterator.N(5)) // iterator.N(5) returns iterator.Seq[int]
//	defer deferStop()
//	assert.Equal(t, gust.Some(0), iterator.Next())
func FromSeq[T any](seq iter.Seq[T]) (Iterator[T], func()) {
	next, stop := iter.Pull(seq)
	return FromPull(next, stop)
}

// FromSeq2 creates an Iterator[gust.Pair[K, V]] from Go's standard iterator.Seq2[K, V].
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
//	pair := iterator.Next()
//	assert.True(t, pair.IsSome())
//	assert.Contains(t, []string{"a", "b", "c"}, pair.Unwrap().A)
//
//	// Works with Go's standard library iterators
//	iter, deferStop = FromSeq2(maps.All(myMap)) // maps.All returns iterator.Seq2[K, V]
//	defer deferStop()
func FromSeq2[K any, V any](seq iter.Seq2[K, V]) (Iterator[gust.Pair[K, V]], func()) {
	return FromPull2(iter.Pull2(seq))
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
//	next, stop := iterator.Pull(someSeq)
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
//	next, stop := iterator.Pull2(someSeq2)
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

// Seq2 converts the Iterator[gust.Pair[K, V]] to Go's standard iterator.Seq2[K, V].
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
func Seq2[K any, V any](it Iterator[gust.Pair[K, V]]) iter.Seq2[K, V] {
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

// Pull2 converts the Iterator[gust.Pair[K, V]] to a pull-style iterator using Go's standard iterator.Pull2.
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
	return iter.Pull2(Seq2(it))
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

// BitSetLike is an interface for bit set implementations.
// This allows FromBitSet to work with any bit set implementation
// without depending on a specific package.
type BitSetLike interface {
	// Size returns the number of bits in the bit set.
	Size() int
	// Get returns the value of the bit at the specified offset.
	//
	// Returns:
	//   - true if the bit at offset is set to 1
	//   - false if the bit at offset is set to 0
	//   - false if offset is out of range (offset < 0 or offset >= Size())
	Get(offset int) bool
}

// FromBitSet creates an iterator over all bits in a bit set,
// yielding pairs of (offset, bool) where offset is the bit position
// and bool indicates whether the bit is set.
//
// # Examples
//
//	// Assuming you have a BitSet implementation
//	type MyBitSet struct {
//		bits []byte
//	}
//	func (b *MyBitSet) Size() int { return len(b.bits) * 8 }
//	func (b *MyBitSet) Get(offset int) bool {
//		if offset < 0 || offset >= b.Size() {
//			return false
//		}
//		byteIdx := offset / 8
//		bitIdx := offset % 8
//		return (b.bits[byteIdx] & (1 << (7 - bitIdx))) != 0
//	}
//
//	bitset := &MyBitSet{bits: []byte{0b10101010}}
//	iter := FromBitSet(bitset)
//	pair := iterator.Next()
//	assert.True(t, pair.IsSome())
//	assert.Equal(t, 0, pair.Unwrap().A)  // offset
//	assert.Equal(t, true, pair.Unwrap().B) // bit value
//
//	// Filter only set bits
//	setBits := FromBitSet(bitset).
//		Filter(func(p gust.Pair[int, bool]) bool { return p.B }).
//		Map(func(p gust.Pair[int, bool]) int { return p.A }).
//		Collect()
func FromBitSet(bitset BitSetLike) Iterator[gust.Pair[int, bool]] {
	size := bitset.Size()
	if size <= 0 {
		return Empty[gust.Pair[int, bool]]()
	}
	return Iterator[gust.Pair[int, bool]]{
		iterable: &bitsetIterable{bitset: bitset, size: size, offset: 0},
	}
}

type bitsetIterable struct {
	bitset BitSetLike
	size   int
	offset int
}

func (b *bitsetIterable) Next() gust.Option[gust.Pair[int, bool]] {
	if b.offset >= b.size {
		return gust.None[gust.Pair[int, bool]]()
	}
	offset := b.offset
	value := b.bitset.Get(offset)
	b.offset++
	return gust.Some(gust.Pair[int, bool]{A: offset, B: value})
}

func (b *bitsetIterable) SizeHint() (uint, gust.Option[uint]) {
	remaining := uint(b.size - b.offset)
	return remaining, gust.Some(remaining)
}

// FromBitSetOnes creates an iterator over only the bits that are set to true (1)
// in a bit set, yielding the offset of each set bit.
//
// # Examples
//
//	bitset := &MyBitSet{bits: []byte{0b10101010}}
//	iter := FromBitSetOnes(bitset)
//	assert.Equal(t, gust.Some(0), iterator.Next())  // first set bit
//	assert.Equal(t, gust.Some(2), iterator.Next())  // second set bit
//	assert.Equal(t, gust.Some(4), iterator.Next())  // third set bit
//	assert.Equal(t, gust.Some(6), iterator.Next())  // fourth set bit
//	assert.Equal(t, gust.None[int](), iterator.Next())
//
//	// Chain with other iterator methods
//	sum := FromBitSetOnes(bitset).
//		Filter(func(offset int) bool { return offset > 2 }).
//		Fold(0, func(acc, offset int) int { return acc + offset })
func FromBitSetOnes(bitset BitSetLike) Iterator[int] {
	return Map(
		filterImpl(FromBitSet(bitset), func(p gust.Pair[int, bool]) bool { return p.B }),
		func(p gust.Pair[int, bool]) int { return p.A },
	)
}

// FromBitSetZeros creates an iterator over only the bits that are set to false (0)
// in a bit set, yielding the offset of each unset bit.
//
// # Examples
//
//	bitset := &MyBitSet{bits: []byte{0b10101010}}
//	iter := FromBitSetZeros(bitset)
//	assert.Equal(t, gust.Some(1), iterator.Next())  // first unset bit
//	assert.Equal(t, gust.Some(3), iterator.Next())  // second unset bit
//	assert.Equal(t, gust.Some(5), iterator.Next())  // third unset bit
//	assert.Equal(t, gust.Some(7), iterator.Next())  // fourth unset bit
//	assert.Equal(t, gust.None[int](), iterator.Next())
func FromBitSetZeros(bitset BitSetLike) Iterator[int] {
	return Map(
		filterImpl(FromBitSet(bitset), func(p gust.Pair[int, bool]) bool { return !p.B }),
		func(p gust.Pair[int, bool]) int { return p.A },
	)
}

// FromBitSetBytes creates an iterator over all bits in a byte slice,
// treating the bytes as a bit set and yielding pairs of (offset, bool)
// where offset is the bit position (0-indexed from the first byte, first bit)
// and bool indicates whether the bit is set.
//
// Bits are ordered from the most significant bit (MSB) to the least significant bit (LSB)
// within each byte, and bytes are ordered from first to last.
//
// # Examples
//
//	bytes := []byte{0b10101010, 0b11001100}
//	iter := FromBitSetBytes(bytes)
//	pair := iterator.Next()
//	assert.True(t, pair.IsSome())
//	assert.Equal(t, 0, pair.Unwrap().A)  // offset
//	assert.Equal(t, true, pair.Unwrap().B) // bit value (MSB of first byte)
//
//	// Get all set bit offsets
//	setBits := FromBitSetBytes(bytes).
//		Filter(func(p gust.Pair[int, bool]) bool { return p.B }).
//		Map(func(p gust.Pair[int, bool]) int { return p.A }).
//		Collect()
//	// setBits contains [0, 2, 4, 6, 8, 9, 12, 13]
//
//	// Count set bits
//	count := FromBitSetBytes(bytes).
//		Filter(func(p gust.Pair[int, bool]) bool { return p.B }).
//		Count()
func FromBitSetBytes(bytes []byte) Iterator[gust.Pair[int, bool]] {
	if len(bytes) == 0 {
		return Empty[gust.Pair[int, bool]]()
	}
	return Iterator[gust.Pair[int, bool]]{
		iterable: &bitSetBytesIterable{bytes: bytes, size: len(bytes) * 8, offset: 0},
	}
}

type bitSetBytesIterable struct {
	bytes  []byte
	size   int
	offset int
}

func (b *bitSetBytesIterable) Next() gust.Option[gust.Pair[int, bool]] {
	if b.offset >= b.size {
		return gust.None[gust.Pair[int, bool]]()
	}
	offset := b.offset
	byteIdx := offset / 8
	bitIdx := offset % 8
	value := (b.bytes[byteIdx] & (1 << (7 - bitIdx))) != 0
	b.offset++
	return gust.Some(gust.Pair[int, bool]{A: offset, B: value})
}

func (b *bitSetBytesIterable) SizeHint() (uint, gust.Option[uint]) {
	remaining := uint(b.size - b.offset)
	return remaining, gust.Some(remaining)
}

// FromBitSetBytesOnes creates an iterator over only the bits that are set to true (1)
// in a byte slice (treated as a bit set), yielding the offset of each set bit.
//
// # Examples
//
//	bytes := []byte{0b10101010, 0b11001100}
//	iter := FromBitSetBytesOnes(bytes)
//	assert.Equal(t, gust.Some(0), iterator.Next())  // first set bit
//	assert.Equal(t, gust.Some(2), iterator.Next())  // second set bit
//	assert.Equal(t, gust.Some(4), iterator.Next())  // third set bit
//	// ... continues with all set bits
func FromBitSetBytesOnes(bytes []byte) Iterator[int] {
	return Map(
		filterImpl(FromBitSetBytes(bytes), func(p gust.Pair[int, bool]) bool { return p.B }),
		func(p gust.Pair[int, bool]) int { return p.A },
	)
}

// FromBitSetBytesZeros creates an iterator over only the bits that are set to false (0)
// in a byte slice (treated as a bit set), yielding the offset of each unset bit.
//
// # Examples
//
//	bytes := []byte{0b10101010, 0b11001100}
//	iter := FromBitSetBytesZeros(bytes)
//	assert.Equal(t, gust.Some(1), iterator.Next())  // first unset bit
//	assert.Equal(t, gust.Some(3), iterator.Next())  // second unset bit
//	assert.Equal(t, gust.Some(5), iterator.Next())  // third unset bit
//	// ... continues with all unset bits
func FromBitSetBytesZeros(bytes []byte) Iterator[int] {
	return Map(
		filterImpl(FromBitSetBytes(bytes), func(p gust.Pair[int, bool]) bool { return !p.B }),
		func(p gust.Pair[int, bool]) int { return p.A },
	)
}
