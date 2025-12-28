package iterator_test

import (
	"iter"
	"testing"

	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/option"
	"github.com/andeya/gust/pair"
	"github.com/stretchr/testify/assert"
)

func TestFromSlice(t *testing.T) {
	a := []int{1, 2, 3}
	iter := iterator.FromSlice(a)

	assert.Equal(t, option.Some(1), iter.Next())
	assert.Equal(t, option.Some(2), iter.Next())
	assert.Equal(t, option.Some(3), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestFromRange(t *testing.T) {
	iter := iterator.FromRange(0, 5)
	var result []int
	for {
		opt := iter.Next()
		if opt.IsNone() {
			break
		}
		result = append(result, opt.Unwrap())
	}
	assert.Equal(t, []int{0, 1, 2, 3, 4}, result)
}

func TestFromElements(t *testing.T) {
	iter := iterator.FromElements(1, 2, 3)
	assert.Equal(t, option.Some(1), iter.Next())
	assert.Equal(t, option.Some(2), iter.Next())
	assert.Equal(t, option.Some(3), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestOnce(t *testing.T) {
	iter := iterator.Once(42)
	assert.Equal(t, option.Some(42), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestEmpty(t *testing.T) {
	iter := iterator.Empty[int]()
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestFromFunc(t *testing.T) {
	count := 0
	iter := iterator.FromFunc(func() option.Option[int] {
		if count < 3 {
			count++
			return option.Some(count)
		}
		return option.None[int]()
	})
	assert.Equal(t, option.Some(1), iter.Next())
	assert.Equal(t, option.Some(2), iter.Next())
	assert.Equal(t, option.Some(3), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestRepeat(t *testing.T) {
	iter := iterator.Repeat(42)
	assert.Equal(t, option.Some(42), iter.Next())
	assert.Equal(t, option.Some(42), iter.Next())
	assert.Equal(t, option.Some(42), iter.Next())
	// Can continue forever
}

func TestFromBitSet(t *testing.T) {
	bitset := &mockBitSet{bits: []byte{0b10101010}}
	iter := iterator.FromBitSet(bitset)
	p := iter.Next()
	assert.True(t, p.IsSome())
	assert.Equal(t, 0, p.Unwrap().A)    // offset
	assert.Equal(t, true, p.Unwrap().B) // bit value

	// Filter only set bits
	setBits := iterator.Map(
		iterator.FromBitSet(bitset).
			Filter(func(p pair.Pair[int, bool]) bool { return p.B }),
		func(p pair.Pair[int, bool]) int { return p.A },
	).Collect()
	assert.Contains(t, setBits, 0)
	assert.Contains(t, setBits, 2)
}

func TestFromBitSetOnes(t *testing.T) {
	bitset := &mockBitSet{bits: []byte{0b10101010}}
	iter := iterator.FromBitSetOnes(bitset)
	assert.Equal(t, option.Some(0), iter.Next()) // first set bit
	assert.Equal(t, option.Some(2), iter.Next()) // second set bit
	assert.Equal(t, option.Some(4), iter.Next()) // third set bit
	assert.Equal(t, option.Some(6), iter.Next()) // fourth set bit
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestFromBitSetZeros(t *testing.T) {
	bitset := &mockBitSet{bits: []byte{0b10101010}}
	iter := iterator.FromBitSetZeros(bitset)
	assert.Equal(t, option.Some(1), iter.Next()) // first unset bit
	assert.Equal(t, option.Some(3), iter.Next()) // second unset bit
	assert.Equal(t, option.Some(5), iter.Next()) // third unset bit
	assert.Equal(t, option.Some(7), iter.Next()) // fourth unset bit
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestFromBitSetBytes(t *testing.T) {
	bytes := []byte{0b10101010, 0b11001100}
	iter := iterator.FromBitSetBytes(bytes)
	p := iter.Next()
	assert.True(t, p.IsSome())
	assert.Equal(t, 0, p.Unwrap().A)    // offset
	assert.Equal(t, true, p.Unwrap().B) // bit value (MSB of first byte)

	// Get all set bit offsets
	setBits := iterator.Map(
		iterator.FromBitSetBytes(bytes).
			Filter(func(p pair.Pair[int, bool]) bool { return p.B }),
		func(p pair.Pair[int, bool]) int { return p.A },
	).Collect()
	assert.Contains(t, setBits, 0)
	assert.Contains(t, setBits, 2)
}

func TestFromBitSetBytesOnes(t *testing.T) {
	bytes := []byte{0b10101010, 0b11001100}
	iter := iterator.FromBitSetBytesOnes(bytes)
	assert.Equal(t, option.Some(0), iter.Next()) // first set bit
	assert.Equal(t, option.Some(2), iter.Next()) // second set bit
	// ... continues with all set bits
}

func TestFromBitSetBytesZeros(t *testing.T) {
	bytes := []byte{0b10101010, 0b11001100}
	iter := iterator.FromBitSetBytesZeros(bytes)
	assert.Equal(t, option.Some(1), iter.Next()) // first unset bit
	assert.Equal(t, option.Some(3), iter.Next()) // second unset bit
	// ... continues with all unset bits
}

// TestRangeIteratorSizeHint tests rangeIterator SizeHint when exhausted
func TestRangeIteratorSizeHint(t *testing.T) {
	iter := iterator.FromRange(0, 3)
	iter.Next()
	iter.Next()
	iter.Next() // Exhaust iterator
	lower, upper := iter.SizeHint()
	assert.Equal(t, uint(0), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(0), upper.Unwrap())
}

// TestOnceIteratorSizeHint tests onceIterator SizeHint
func TestOnceIteratorSizeHint(t *testing.T) {
	iter := iterator.Once(42)
	lower, upper := iter.SizeHint()
	assert.Equal(t, uint(1), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(1), upper.Unwrap())

	iter.Next() // Consume the value
	lower2, upper2 := iter.SizeHint()
	assert.Equal(t, uint(0), lower2)
	assert.True(t, upper2.IsSome())
	assert.Equal(t, uint(0), upper2.Unwrap())
}

// TestRepeatIteratorSizeHint tests repeatIterator SizeHint
func TestRepeatIteratorSizeHint(t *testing.T) {
	iter := iterator.Repeat(42)
	lower, upper := iter.SizeHint()
	assert.Equal(t, uint(0), lower)
	assert.True(t, upper.IsNone()) // Infinite iterator
}

// TestEmptyIteratorSizeHint tests emptyIterator SizeHint
func TestEmptyIteratorSizeHint(t *testing.T) {
	iter := iterator.Empty[int]()
	lower, upper := iter.SizeHint()
	assert.Equal(t, uint(0), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(0), upper.Unwrap())
}

// TestRangeIteratorSizeHintNotExhausted tests rangeIterator SizeHint when not exhausted
func TestRangeIteratorSizeHintNotExhausted(t *testing.T) {
	iter := iterator.FromRange(0, 5)
	lower, upper := iter.SizeHint()
	assert.Equal(t, uint(5), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(5), upper.Unwrap())
}

// TestOnceIteratorSizeHintDone tests onceIterator SizeHint when done
func TestOnceIteratorSizeHintDone(t *testing.T) {
	iter := iterator.Once(42)
	iter.Next() // Consume the value
	lower, upper := iter.SizeHint()
	assert.Equal(t, uint(0), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(0), upper.Unwrap())
}

// mockBitSet is a simple implementation of BitSetLike for testing
type mockBitSet struct {
	bits []byte
}

func (m *mockBitSet) Size() int {
	return len(m.bits) * 8
}

func (m *mockBitSet) Get(offset int) bool {
	if offset < 0 || offset >= m.Size() {
		return false
	}
	byteIdx := offset / 8
	bitIdx := offset % 8
	return (m.bits[byteIdx] & (1 << (7 - bitIdx))) != 0
}

func TestFromSeq(t *testing.T) {
	// Test with custom sequence function
	goSeq := func(yield func(int) bool) {
		for i := 0; i < 5; i++ {
			if !yield(i) {
				return
			}
		}
	}
	gustIter, deferStop := iterator.FromSeq(goSeq)
	defer deferStop()

	var result []int
	for i := 0; i < 5; i++ {
		opt := gustIter.Next()
		assert.True(t, opt.IsSome())
		result = append(result, opt.Unwrap())
	}
	assert.Equal(t, []int{0, 1, 2, 3, 4}, result)
	assert.True(t, gustIter.Next().IsNone())

	// Test with custom sequence
	customSeq := func(yield func(int) bool) {
		for i := 0; i < 3; i++ {
			if !yield(i * 2) {
				return
			}
		}
	}
	gustIter, deferStop = iterator.FromSeq(customSeq)
	defer deferStop()
	result = nil
	for {
		opt := gustIter.Next()
		if opt.IsNone() {
			break
		}
		result = append(result, opt.Unwrap())
	}
	assert.Equal(t, []int{0, 2, 4}, result)

	// Test with empty sequence
	emptySeq := func(yield func(int) bool) {}
	gustIter, deferStop = iterator.FromSeq(emptySeq)
	defer deferStop()
	assert.True(t, gustIter.Next().IsNone())

	// Test chaining gust methods after FromSeq
	goSeq2 := func(yield func(int) bool) {
		for i := 0; i < 5; i++ {
			if !yield(i) {
				return
			}
		}
	}
	gustIter, deferStop = iterator.FromSeq(goSeq2)
	defer deferStop()
	filtered := gustIter.Filter(func(x int) bool { return x > 2 })
	result = filtered.Collect()
	assert.Equal(t, []int{3, 4}, result)
}

func TestFromSeq_WithCustomSequence(t *testing.T) {
	// Test with custom sequence function
	seq := func(yield func(int) bool) {
		for i := 0; i < 3; i++ {
			if !yield(i) {
				return
			}
		}
	}
	it, deferStop := iterator.FromSeq(seq)
	defer deferStop()
	result := it.Collect()
	assert.Equal(t, []int{0, 1, 2}, result)

	// Test with slice-based sequence
	slice := []string{"a", "b", "c"}
	sliceSeq := func(yield func(string) bool) {
		for _, v := range slice {
			if !yield(v) {
				return
			}
		}
	}
	it2, deferStop2 := iterator.FromSeq(sliceSeq)
	defer deferStop2()
	result2 := it2.Collect()
	assert.Equal(t, []string{"a", "b", "c"}, result2)
}

func TestFromSeq2(t *testing.T) {
	// Test with custom map-like sequence
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	seq2 := func(yield func(string, int) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
	gustIter, deferStop := iterator.FromSeq2(seq2)
	defer deferStop()

	var result []pair.Pair[string, int]
	for {
		opt := gustIter.Next()
		if opt.IsNone() {
			break
		}
		result = append(result, opt.Unwrap())
	}
	assert.Len(t, result, 3)
	// Check that all pairs are present (order may vary)
	keys := make(map[string]int)
	for _, p := range result {
		keys[p.A] = p.B
	}
	assert.Equal(t, 1, keys["a"])
	assert.Equal(t, 2, keys["b"])
	assert.Equal(t, 3, keys["c"])

	// Test with empty sequence
	emptySeq2 := func(yield func(string, int) bool) {}
	gustIter, deferStop = iterator.FromSeq2(emptySeq2)
	defer deferStop()
	assert.True(t, gustIter.Next().IsNone())

	// Test with custom key-value sequence
	customSeq2 := func(yield func(int, string) bool) {
		for i := 0; i < 3; i++ {
			if !yield(i, string(rune('a'+i))) {
				return
			}
		}
	}
	var customIter iterator.Iterator[pair.Pair[int, string]]
	customIter, deferStop = iterator.FromSeq2(customSeq2)
	defer deferStop()
	var customResult []pair.Pair[int, string]
	for {
		opt := customIter.Next()
		if opt.IsNone() {
			break
		}
		customResult = append(customResult, opt.Unwrap())
	}
	assert.Equal(t, []pair.Pair[int, string]{
		{A: 0, B: "a"},
		{A: 1, B: "b"},
		{A: 2, B: "c"},
	}, customResult)

	// Test chaining gust methods after FromSeq2
	seq2Chain := func(yield func(int, int) bool) {
		for i := 0; i < 5; i++ {
			if !yield(i, i*2) {
				return
			}
		}
	}
	var chainIter iterator.Iterator[pair.Pair[int, int]]
	chainIter, deferStop = iterator.FromSeq2(seq2Chain)
	defer deferStop()
	// Filter pairs where value > 3
	// Sequence: (0,0), (1,2), (2,4), (3,6), (4,8)
	// Filter p.B > 3: (2,4), (3,6), (4,8) = 3 items
	filtered := chainIter.Filter(func(p pair.Pair[int, int]) bool {
		return p.B > 3
	})
	chainResult := filtered.Collect()
	assert.Len(t, chainResult, 3)
	assert.Equal(t, []pair.Pair[int, int]{
		{A: 2, B: 4},
		{A: 3, B: 6},
		{A: 4, B: 8},
	}, chainResult)
}

func TestFromPull(t *testing.T) {
	// Test with iterator.Pull result
	seq := func(yield func(int) bool) {
		for i := 0; i < 5; i++ {
			if !yield(i) {
				return
			}
		}
	}
	next, stop := iter.Pull(seq)
	defer stop()

	gustIter, _ := iterator.FromPull(next, stop)
	var result []int
	for {
		opt := gustIter.Next()
		if opt.IsNone() {
			break
		}
		result = append(result, opt.Unwrap())
	}
	assert.Equal(t, []int{0, 1, 2, 3, 4}, result)

	// Test with custom pull-style iterator
	count := 0
	customNext := func() (int, bool) {
		if count >= 3 {
			return 0, false
		}
		val := count * 2
		count++
		return val, true
	}
	customStop := func() {}

	gustIter, _ = iterator.FromPull(customNext, customStop)
	result = nil
	for {
		opt := gustIter.Next()
		if opt.IsNone() {
			break
		}
		result = append(result, opt.Unwrap())
	}
	assert.Equal(t, []int{0, 2, 4}, result)

	// Test with empty pull iterator
	emptyNext := func() (int, bool) {
		return 0, false
	}
	emptyStop := func() {}

	gustIter, _ = iterator.FromPull(emptyNext, emptyStop)
	assert.True(t, gustIter.Next().IsNone())

	// Test chaining gust methods after FromPull
	seq2 := func(yield func(int) bool) {
		for i := 0; i < 5; i++ {
			if !yield(i) {
				return
			}
		}
	}
	next2, stop2 := iter.Pull(seq2)
	defer stop2()

	gustIter, _ = iterator.FromPull(next2, stop2)
	filtered := gustIter.Filter(func(x int) bool { return x > 2 })
	result = filtered.Collect()
	assert.Equal(t, []int{3, 4}, result)
}

func TestFromPull2(t *testing.T) {
	// Test with iterator.Pull2 result
	seq2 := func(yield func(int, string) bool) {
		pairs := []pair.Pair[int, string]{
			{A: 1, B: "a"},
			{A: 2, B: "b"},
			{A: 3, B: "c"},
		}
		for _, p := range pairs {
			if !yield(p.A, p.B) {
				return
			}
		}
	}
	next, stop := iter.Pull2(seq2)
	defer stop()

	gustIter, deferStop := iterator.FromPull2(next, stop)
	defer deferStop()
	var result []pair.Pair[int, string]
	for {
		opt := gustIter.Next()
		if opt.IsNone() {
			break
		}
		result = append(result, opt.Unwrap())
	}
	assert.Equal(t, []pair.Pair[int, string]{
		{A: 1, B: "a"},
		{A: 2, B: "b"},
		{A: 3, B: "c"},
	}, result)

	// Test with custom pull-style iterator
	count := 0
	customNext := func() (int, string, bool) {
		if count >= 3 {
			return 0, "", false
		}
		k := count
		v := string(rune('a' + count))
		count++
		return k, v, true
	}
	customStop := func() {}

	gustIter, deferStop = iterator.FromPull2(customNext, customStop)
	defer deferStop()
	result = nil
	for {
		opt := gustIter.Next()
		if opt.IsNone() {
			break
		}
		result = append(result, opt.Unwrap())
	}
	assert.Equal(t, []pair.Pair[int, string]{
		{A: 0, B: "a"},
		{A: 1, B: "b"},
		{A: 2, B: "c"},
	}, result)

	// Test with empty pull iterator
	emptyNext := func() (int, string, bool) {
		return 0, "", false
	}
	emptyStop := func() {}

	gustIter, deferStop = iterator.FromPull2(emptyNext, emptyStop)
	defer deferStop()
	assert.True(t, gustIter.Next().IsNone())

	// Test chaining gust methods after FromPull2
	seq2Chain := func(yield func(int, int) bool) {
		for i := 0; i < 5; i++ {
			if !yield(i, i*2) {
				return
			}
		}
	}
	next2, stop2 := iter.Pull2(seq2Chain)
	defer stop2()

	var chainIter iterator.Iterator[pair.Pair[int, int]]
	chainIter, deferStop = iterator.FromPull2(next2, stop2)
	defer deferStop()
	// Filter pairs where value > 3
	// Sequence: (0,0), (1,2), (2,4), (3,6), (4,8)
	// Filter p.B > 3: (2,4), (3,6), (4,8) = 3 items
	filtered := chainIter.Filter(func(p pair.Pair[int, int]) bool {
		return p.B > 3
	})
	chainResult := filtered.Collect()
	assert.Len(t, chainResult, 3)
	assert.Equal(t, []pair.Pair[int, int]{
		{A: 2, B: 4},
		{A: 3, B: 6},
		{A: 4, B: 8},
	}, chainResult)
}

func TestBitSetIteratorChaining(t *testing.T) {
	// Test chaining with Filter, Map, Take, etc.
	bitset := &mockBitSet{bits: []byte{0b10101010, 0b11001100}}

	// Get offsets of set bits that are greater than 5
	result := iterator.FromBitSetOnes(bitset).
		Filter(func(offset int) bool { return offset > 5 }).
		Collect()
	assert.Equal(t, []int{6, 8, 9, 12, 13}, result)

	// Sum of offsets of set bits
	sum := iterator.FromBitSetOnes(bitset).
		Fold(0, func(acc, offset int) int { return acc + offset })
	assert.Equal(t, 54, sum) // 0+2+4+6+8+9+12+13 = 54

	// Count set bits
	count := iterator.FromBitSetOnes(bitset).Count()
	assert.Equal(t, uint(8), count)

	// Take first 3 set bits
	firstThree := iterator.FromBitSetOnes(bitset).
		Take(3).
		Collect()
	assert.Equal(t, []int{0, 2, 4}, firstThree)
}

func TestBytesIteratorChaining(t *testing.T) {
	// Test chaining with Filter, Map, Take, etc.
	bytes := []byte{0b10101010, 0b11001100}

	// Get offsets of set bits that are greater than 5
	result := iterator.FromBitSetBytesOnes(bytes).
		Filter(func(offset int) bool { return offset > 5 }).
		Collect()
	assert.Equal(t, []int{6, 8, 9, 12, 13}, result)

	// Sum of offsets of set bits
	sum := iterator.FromBitSetBytesOnes(bytes).
		Fold(0, func(acc, offset int) int { return acc + offset })
	assert.Equal(t, 54, sum) // 0+2+4+6+8+9+12+13 = 54

	// Count set bits
	count := iterator.FromBitSetBytesOnes(bytes).Count()
	assert.Equal(t, uint(8), count)

	// Take first 3 set bits
	firstThree := iterator.FromBitSetBytesOnes(bytes).
		Take(3).
		Collect()
	assert.Equal(t, []int{0, 2, 4}, firstThree)
}
