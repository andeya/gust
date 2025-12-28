package iterator_test

import (
	"iter"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iterator"
	"github.com/stretchr/testify/assert"
)

func TestIterator_Seq(t *testing.T) {
	// Test basic conversion
	iter := iterator.FromSlice([]int{1, 2, 3})
	var result []int
	for v := range iter.Seq() {
		result = append(result, v)
	}
	assert.Equal(t, []int{1, 2, 3}, result)

	// Test with filtered iterator
	filtered := iterator.FromSlice([]int{1, 2, 3, 4, 5}).Filter(func(x int) bool { return x%2 == 0 })
	result = nil
	for v := range filtered.Seq() {
		result = append(result, v)
	}
	assert.Equal(t, []int{2, 4}, result)

	// Test with empty iterator
	empty := iterator.Empty[int]()
	result = nil
	for v := range empty.Seq() {
		result = append(result, v)
	}
	assert.Nil(t, result)
	assert.Len(t, result, 0)

	// Test early termination
	iter = iterator.FromSlice([]int{1, 2, 3, 4, 5})
	result = nil
	count := 0
	for v := range iter.Seq() {
		result = append(result, v)
		count++
		if count >= 3 {
			break
		}
	}
	assert.Equal(t, []int{1, 2, 3}, result)
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

func TestSeq_RoundTrip(t *testing.T) {
	// Test round trip: gust Iterator -> Seq -> gust Iterator
	// Create two independent iterators: one for seq, one for expected result
	originalForSeq := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	originalForExpected := iterator.FromSlice([]int{1, 2, 3, 4, 5})

	seq := originalForSeq.Seq()
	converted, deferStop := iterator.FromSeq(seq)
	defer deferStop()

	// Get expected result from independent iterator
	var expectedResult []int
	for {
		opt := originalForExpected.Next()
		if opt.IsNone() {
			break
		}
		expectedResult = append(expectedResult, opt.Unwrap())
	}

	// Get actual result from converted iterator
	var convertedResult []int
	for {
		opt := converted.Next()
		if opt.IsNone() {
			break
		}
		convertedResult = append(convertedResult, opt.Unwrap())
	}

	assert.Equal(t, expectedResult, convertedResult)
}

func TestSeq_WithMap(t *testing.T) {
	// Test Seq with mapped iterator
	iter := iterator.FromSlice([]int{1, 2, 3}).Map(func(x int) int { return x * 2 })
	var result []int
	for v := range iter.Seq() {
		result = append(result, v)
	}
	assert.Equal(t, []int{2, 4, 6}, result)
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

func TestSeq2(t *testing.T) {
	// Test with Zip iterator
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	iter2 := iterator.FromSlice([]string{"a", "b", "c"})
	zipped := iterator.Zip(iter1, iter2)

	var result []gust.Pair[int, string]
	for k, v := range iterator.Seq2(zipped) {
		result = append(result, gust.Pair[int, string]{A: k, B: v})
	}
	assert.Equal(t, []gust.Pair[int, string]{
		{A: 1, B: "a"},
		{A: 2, B: "b"},
		{A: 3, B: "c"},
	}, result)

	// Test with Enumerate iterator
	enumerated := iterator.Enumerate(iterator.FromSlice([]string{"x", "y", "z"}))
	var enumResult []gust.Pair[uint, string]
	for idx, val := range iterator.Seq2(enumerated) {
		enumResult = append(enumResult, gust.Pair[uint, string]{A: idx, B: val})
	}
	assert.Equal(t, []gust.Pair[uint, string]{
		{A: 0, B: "x"},
		{A: 1, B: "y"},
		{A: 2, B: "z"},
	}, enumResult)

	// Test with empty iterator
	empty := iterator.Zip(iterator.Empty[int](), iterator.Empty[string]())
	var emptyResult []gust.Pair[int, string]
	for k, v := range iterator.Seq2(empty) {
		emptyResult = append(emptyResult, gust.Pair[int, string]{A: k, B: v})
	}
	assert.Nil(t, emptyResult)
	assert.Len(t, emptyResult, 0)

	// Test early termination
	iter1 = iterator.FromSlice([]int{1, 2, 3, 4, 5})
	iter2 = iterator.FromSlice([]string{"a", "b", "c", "d", "e"})
	zipped = iterator.Zip(iter1, iter2)
	var earlyResult []gust.Pair[int, string]
	count := 0
	for k, v := range iterator.Seq2(zipped) {
		earlyResult = append(earlyResult, gust.Pair[int, string]{A: k, B: v})
		count++
		if count >= 2 {
			break
		}
	}
	assert.Equal(t, []gust.Pair[int, string]{
		{A: 1, B: "a"},
		{A: 2, B: "b"},
	}, earlyResult)
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

	var result []gust.Pair[string, int]
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
	var customIter iterator.Iterator[gust.Pair[int, string]]
	customIter, deferStop = iterator.FromSeq2(customSeq2)
	defer deferStop()
	var customResult []gust.Pair[int, string]
	for {
		opt := customIter.Next()
		if opt.IsNone() {
			break
		}
		customResult = append(customResult, opt.Unwrap())
	}
	assert.Equal(t, []gust.Pair[int, string]{
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
	var chainIter iterator.Iterator[gust.Pair[int, int]]
	chainIter, deferStop = iterator.FromSeq2(seq2Chain)
	defer deferStop()
	// Filter pairs where value > 3
	// Sequence: (0,0), (1,2), (2,4), (3,6), (4,8)
	// Filter p.B > 3: (2,4), (3,6), (4,8) = 3 items
	filtered := chainIter.Filter(func(p gust.Pair[int, int]) bool {
		return p.B > 3
	})
	chainResult := filtered.Collect()
	assert.Len(t, chainResult, 3)
	assert.Equal(t, []gust.Pair[int, int]{
		{A: 2, B: 4},
		{A: 3, B: 6},
		{A: 4, B: 8},
	}, chainResult)
}

func TestSeq2_RoundTrip(t *testing.T) {
	// Test round trip: gust Iterator -> Seq2 -> gust Iterator
	// Create two independent iterators: one for seq2, one for expected result
	iter1ForSeq := iterator.FromSlice([]int{1, 2, 3})
	iter2ForSeq := iterator.FromSlice([]string{"a", "b", "c"})
	originalForSeq := iterator.Zip(iter1ForSeq, iter2ForSeq)

	iter1ForExpected := iterator.FromSlice([]int{1, 2, 3})
	iter2ForExpected := iterator.FromSlice([]string{"a", "b", "c"})
	originalForExpected := iterator.Zip(iter1ForExpected, iter2ForExpected)

	seq2 := iterator.Seq2(originalForSeq)
	converted, deferStop := iterator.FromSeq2(seq2)
	defer deferStop()

	// Get expected result from independent iterator
	var expectedResult []gust.Pair[int, string]
	for {
		opt := originalForExpected.Next()
		if opt.IsNone() {
			break
		}
		expectedResult = append(expectedResult, opt.Unwrap())
	}

	// Get actual result from converted iterator
	var convertedResult []gust.Pair[int, string]
	for {
		opt := converted.Next()
		if opt.IsNone() {
			break
		}
		convertedResult = append(convertedResult, opt.Unwrap())
	}

	assert.Equal(t, expectedResult, convertedResult)
}

func TestSeq2_WithGoStandardLibrary(t *testing.T) {
	// Test that Seq2 works with Go's standard library iterator.Pull2
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	iter2 := iterator.FromSlice([]string{"a", "b", "c"})
	zipped := iterator.Zip(iter1, iter2)
	seq2 := iterator.Seq2(zipped)

	// Use iter.Pull2 to pull values manually
	next, stop := iter.Pull2(seq2)
	defer stop()

	var result []gust.Pair[int, string]
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		result = append(result, gust.Pair[int, string]{A: k, B: v})
	}

	assert.Equal(t, []gust.Pair[int, string]{
		{A: 1, B: "a"},
		{A: 2, B: "b"},
		{A: 3, B: "c"},
	}, result)
}

func TestIterator_Pull(t *testing.T) {
	// Test basic Pull functionality
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	next, stop := iter.Pull()
	defer stop()

	var result []int
	for {
		v, ok := next()
		if !ok {
			break
		}
		result = append(result, v)
	}
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)

	// Test early termination
	iter = iterator.FromSlice([]int{1, 2, 3, 4, 5})
	next, stop = iter.Pull()
	defer stop()

	result = nil
	for {
		v, ok := next()
		if !ok {
			break
		}
		result = append(result, v)
		if v == 3 {
			break // Early termination
		}
	}
	assert.Equal(t, []int{1, 2, 3}, result)

	// Test with empty iterator
	empty := iterator.Empty[int]()
	next, stop = empty.Pull()
	defer stop()

	result = nil
	for {
		v, ok := next()
		if !ok {
			break
		}
		result = append(result, v)
	}
	assert.Nil(t, result)
	assert.Len(t, result, 0)

	// Test with filtered iterator
	filtered := iterator.FromSlice([]int{1, 2, 3, 4, 5}).Filter(func(x int) bool { return x%2 == 0 })
	next, stop = filtered.Pull()
	defer stop()

	result = nil
	for {
		v, ok := next()
		if !ok {
			break
		}
		result = append(result, v)
	}
	assert.Equal(t, []int{2, 4}, result)
}

func TestPull2(t *testing.T) {
	// Test basic Pull2 functionality
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	iter2 := iterator.FromSlice([]string{"a", "b", "c"})
	zipped := iterator.Zip(iter1, iter2)

	next, stop := iterator.Pull2(zipped)
	defer stop()

	var result []gust.Pair[int, string]
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		result = append(result, gust.Pair[int, string]{A: k, B: v})
	}

	assert.Equal(t, []gust.Pair[int, string]{
		{A: 1, B: "a"},
		{A: 2, B: "b"},
		{A: 3, B: "c"},
	}, result)

	// Test early termination
	iter1 = iterator.FromSlice([]int{1, 2, 3, 4, 5})
	iter2 = iterator.FromSlice([]string{"a", "b", "c", "d", "e"})
	zipped = iterator.Zip(iter1, iter2)

	next, stop = iterator.Pull2(zipped)
	defer stop()

	result = nil
	count := 0
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		result = append(result, gust.Pair[int, string]{A: k, B: v})
		count++
		if count >= 2 {
			break // Early termination
		}
	}
	assert.Equal(t, []gust.Pair[int, string]{
		{A: 1, B: "a"},
		{A: 2, B: "b"},
	}, result)

	// Test with empty iterator
	empty1 := iterator.Empty[int]()
	empty2 := iterator.Empty[string]()
	zipped = iterator.Zip(empty1, empty2)

	next, stop = iterator.Pull2(zipped)
	defer stop()

	result = nil
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		result = append(result, gust.Pair[int, string]{A: k, B: v})
	}
	assert.Nil(t, result)
	assert.Len(t, result, 0)

	// Test with Enumerate
	enumerated := iterator.Enumerate(iterator.FromSlice([]string{"x", "y", "z"}))
	nextEnum, stopEnum := iterator.Pull2(enumerated)
	defer stopEnum()

	var enumResult []gust.Pair[uint, string]
	for {
		idx, val, ok := nextEnum()
		if !ok {
			break
		}
		enumResult = append(enumResult, gust.Pair[uint, string]{A: idx, B: val})
	}

	assert.Equal(t, []gust.Pair[uint, string]{
		{A: 0, B: "x"},
		{A: 1, B: "y"},
		{A: 2, B: "z"},
	}, enumResult)
}

func TestIterator_Seq2(t *testing.T) {
	// Test basic Seq2 conversion - converts Iterator[T] to iterator.Seq2[uint, T]
	iter := iterator.FromSlice([]int{1, 2, 3})
	var result []gust.Pair[uint, int]
	for k, v := range iter.Seq2() {
		result = append(result, gust.Pair[uint, int]{A: k, B: v})
	}
	assert.Equal(t, []gust.Pair[uint, int]{
		{A: 0, B: 1},
		{A: 1, B: 2},
		{A: 2, B: 3},
	}, result)

	// Test with filtered iterator
	filtered := iterator.FromSlice([]int{10, 20, 30, 40, 50}).Filter(func(x int) bool { return x > 20 })
	var filteredResult []gust.Pair[uint, int]
	for k, v := range filtered.Seq2() {
		filteredResult = append(filteredResult, gust.Pair[uint, int]{A: k, B: v})
	}
	assert.Equal(t, []gust.Pair[uint, int]{
		{A: 0, B: 30},
		{A: 1, B: 40},
		{A: 2, B: 50},
	}, filteredResult)

	// Test with empty iterator
	empty := iterator.Empty[int]()
	var emptyResult []gust.Pair[uint, int]
	for k, v := range empty.Seq2() {
		emptyResult = append(emptyResult, gust.Pair[uint, int]{A: k, B: v})
	}
	assert.Nil(t, emptyResult)
	assert.Len(t, emptyResult, 0)

	// Test early termination
	iter = iterator.FromSlice([]int{1, 2, 3, 4, 5})
	var earlyResult []gust.Pair[uint, int]
	count := 0
	for k, v := range iter.Seq2() {
		earlyResult = append(earlyResult, gust.Pair[uint, int]{A: k, B: v})
		count++
		if count >= 3 {
			break
		}
	}
	assert.Equal(t, []gust.Pair[uint, int]{
		{A: 0, B: 1},
		{A: 1, B: 2},
		{A: 2, B: 3},
	}, earlyResult)

	// Test with string iterator
	strIter := iterator.FromSlice([]string{"hello", "world", "rust"})
	var strResult []gust.Pair[uint, string]
	for k, v := range strIter.Seq2() {
		strResult = append(strResult, gust.Pair[uint, string]{A: k, B: v})
	}
	assert.Equal(t, []gust.Pair[uint, string]{
		{A: 0, B: "hello"},
		{A: 1, B: "world"},
		{A: 2, B: "rust"},
	}, strResult)
}

func TestIterator_Pull2(t *testing.T) {
	// Test basic Pull2 functionality - converts Iterator[T] to pull-style iterator with index-value pairs
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	next, stop := iter.Pull2()
	defer stop()

	var result []gust.Pair[uint, int]
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		result = append(result, gust.Pair[uint, int]{A: k, B: v})
	}
	assert.Equal(t, []gust.Pair[uint, int]{
		{A: 0, B: 1},
		{A: 1, B: 2},
		{A: 2, B: 3},
		{A: 3, B: 4},
		{A: 4, B: 5},
	}, result)

	// Test early termination
	iter = iterator.FromSlice([]int{1, 2, 3, 4, 5})
	next, stop = iter.Pull2()
	defer stop()

	result = nil
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		result = append(result, gust.Pair[uint, int]{A: k, B: v})
		if v == 3 {
			break // Early termination
		}
	}
	assert.Equal(t, []gust.Pair[uint, int]{
		{A: 0, B: 1},
		{A: 1, B: 2},
		{A: 2, B: 3},
	}, result)

	// Test with empty iterator
	empty := iterator.Empty[int]()
	next, stop = empty.Pull2()
	defer stop()

	result = nil
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		result = append(result, gust.Pair[uint, int]{A: k, B: v})
	}
	assert.Nil(t, result)
	assert.Len(t, result, 0)

	// Test with filtered iterator
	filtered := iterator.FromSlice([]int{10, 20, 30, 40, 50}).Filter(func(x int) bool { return x%20 == 0 })
	next, stop = filtered.Pull2()
	defer stop()

	result = nil
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		result = append(result, gust.Pair[uint, int]{A: k, B: v})
	}
	assert.Equal(t, []gust.Pair[uint, int]{
		{A: 0, B: 20},
		{A: 1, B: 40},
	}, result)

	// Test with string iterator
	strIter := iterator.FromSlice([]string{"a", "b", "c"})
	nextStr, stopStr := strIter.Pull2()
	defer stopStr()

	var strResult []gust.Pair[uint, string]
	for {
		k, v, ok := nextStr()
		if !ok {
			break
		}
		strResult = append(strResult, gust.Pair[uint, string]{A: k, B: v})
	}
	assert.Equal(t, []gust.Pair[uint, string]{
		{A: 0, B: "a"},
		{A: 1, B: "b"},
		{A: 2, B: "c"},
	}, strResult)
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
		pairs := []gust.Pair[int, string]{
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
	var result []gust.Pair[int, string]
	for {
		opt := gustIter.Next()
		if opt.IsNone() {
			break
		}
		result = append(result, opt.Unwrap())
	}
	assert.Equal(t, []gust.Pair[int, string]{
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
	assert.Equal(t, []gust.Pair[int, string]{
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

	var chainIter iterator.Iterator[gust.Pair[int, int]]
	chainIter, deferStop = iterator.FromPull2(next2, stop2)
	defer deferStop()
	// Filter pairs where value > 3
	// Sequence: (0,0), (1,2), (2,4), (3,6), (4,8)
	// Filter p.B > 3: (2,4), (3,6), (4,8) = 3 items
	filtered := chainIter.Filter(func(p gust.Pair[int, int]) bool {
		return p.B > 3
	})
	chainResult := filtered.Collect()
	assert.Len(t, chainResult, 3)
	assert.Equal(t, []gust.Pair[int, int]{
		{A: 2, B: 4},
		{A: 3, B: 6},
		{A: 4, B: 8},
	}, chainResult)
}
