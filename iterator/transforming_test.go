package iterator_test

import (
	"testing"

	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/option"
	"github.com/stretchr/testify/assert"
)

// TestMapWhile tests iterator.MapWhile functionality
func TestMapWhile(t *testing.T) {
	a := []int{-1, 4, 0, 1}
	iter := iterator.MapWhile(iterator.FromSlice(a), func(x int) option.Option[int] {
		if x != 0 {
			return option.Some(16 / x)
		}
		return option.None[int]()
	})

	assert.Equal(t, option.Some(-16), iter.Next())
	assert.Equal(t, option.Some(4), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestScan tests iterator.Scan functionality
func TestScan(t *testing.T) {
	a := []int{1, 2, 3, 4}
	iter := iterator.Scan(iterator.FromSlice(a), 1, func(state *int, x int) option.Option[int] {
		*state = *state * x
		if *state > 6 {
			return option.None[int]()
		}
		return option.Some(-*state)
	})

	assert.Equal(t, option.Some(-1), iter.Next())
	assert.Equal(t, option.Some(-2), iter.Next())
	assert.Equal(t, option.Some(-6), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestScanNoEarlyTermination tests iterator.Scan without early termination
func TestScanNoEarlyTermination(t *testing.T) {
	a := []int{1, 2, 3}
	iter := iterator.Scan(iterator.FromSlice(a), 0, func(state *int, x int) option.Option[int] {
		*state = *state + x
		return option.Some(*state)
	})

	assert.Equal(t, option.Some(1), iter.Next())
	assert.Equal(t, option.Some(3), iter.Next())
	assert.Equal(t, option.Some(6), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestFlatMap tests iterator.FlatMap functionality
func TestFlatMap(t *testing.T) {
	words := []string{"alpha", "beta"}
	iter := iterator.FlatMap(iterator.FromSlice(words), func(s string) iterator.Iterator[rune] {
		return iterator.FromSlice([]rune(s))
	})

	result := iter.Collect()
	expected := []rune{'a', 'l', 'p', 'h', 'a', 'b', 'e', 't', 'a'}
	assert.Equal(t, expected, result)
}

// TestFlatMapEmptyInner tests iterator.FlatMap with empty inner iterators
func TestFlatMapEmptyInner(t *testing.T) {
	words := []string{"", "a", ""}
	iter := iterator.FlatMap(iterator.FromSlice(words), func(s string) iterator.Iterator[rune] {
		return iterator.FromSlice([]rune(s))
	})

	result := iter.Collect()
	assert.Equal(t, []rune{'a'}, result)
}

// TestFlatten tests iterator.Flatten functionality
func TestFlatten(t *testing.T) {
	data := [][]int{{1, 2, 3, 4}, {5, 6}}
	iters := make([]iterator.Iterator[int], len(data))
	for i, slice := range data {
		iters[i] = iterator.FromSlice(slice)
	}
	iter := iterator.Flatten(iterator.FromSlice(iters))
	result := iter.Collect()
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
}

// TestFlattenEmptyInner tests iterator.Flatten with empty inner iterators
func TestFlattenEmptyInner(t *testing.T) {
	data := [][]int{{}, {1, 2}, {}}
	iters := make([]iterator.Iterator[int], len(data))
	for i, slice := range data {
		iters[i] = iterator.FromSlice(slice)
	}
	iter := iterator.Flatten(iterator.FromSlice(iters))
	result := iter.Collect()
	assert.Equal(t, []int{1, 2}, result)
}

// TestFlatMap_CurrentNil tests flatMapIterable when current becomes nil
func TestFlatMap_CurrentNil(t *testing.T) {
	iter := iterator.FlatMap(iterator.FromSlice([]int{1, 2}), func(x int) iterator.Iterator[int] {
		if x == 1 {
			return iterator.FromSlice([]int{10, 20}) // Non-empty iterator
		}
		return iterator.Empty[int]() // iterator.Empty iterator - current will become nil
	})

	// Should yield elements from first iterator
	assert.Equal(t, option.Some(10), iter.Next())
	assert.Equal(t, option.Some(20), iter.Next())

	// After first iterator is exhausted, current becomes nil, should move to next
	// Second iterator is empty, so should return None
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestFlatten_CurrentNil tests flattenIterable when current becomes nil
func TestFlatten_CurrentNil(t *testing.T) {
	iter := iterator.Flatten(iterator.FromSlice([]iterator.Iterator[int]{
		iterator.FromSlice([]int{1, 2}),
		iterator.Empty[int](), // iterator.Empty iterator - current will become nil
		iterator.FromSlice([]int{3, 4}),
	}))

	// Should yield elements from first iterator
	assert.Equal(t, option.Some(1), iter.Next())
	assert.Equal(t, option.Some(2), iter.Next())

	// After first iterator is exhausted, current becomes nil, should move to next
	// Second iterator is empty, so should skip to third
	assert.Equal(t, option.Some(3), iter.Next())
	assert.Equal(t, option.Some(4), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestIterator_XMapWhile(t *testing.T) {
	iter := iterator.FromSlice([]int{-1, 4, 0, 1})
	mapped := iter.XMapWhile(func(x int) option.Option[any] {
		if x != 0 {
			return option.Some(any(16 / x))
		}
		return option.None[any]()
	})
	result := mapped.Collect()
	assert.Equal(t, []any{-16, 4}, result)
}

func TestIterator_XScan(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	scanned := iter.XScan(0, func(state *any, x int) option.Option[any] {
		s := (*state).(int) + x
		*state = s
		return option.Some(any(s))
	})
	result := scanned.Collect()
	assert.Equal(t, []any{1, 3, 6}, result)
}

// TestIterator_WrapperMethods_Transforming tests transforming wrapper methods from TestIterator_WrapperMethods
func TestIterator_WrapperMethods_Transforming(t *testing.T) {
	// Test MapWhile (covers iterator_methods.go:713-715)
	iter7 := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	mapped := iter7.MapWhile(func(x int) option.Option[int] {
		if x < 4 {
			return option.Some(x * 2)
		}
		return option.None[int]()
	})
	result5 := mapped.Collect()
	assert.Equal(t, []int{2, 4, 6}, result5)

	// Test Scan (covers iterator_methods.go:746-748)
	iter8 := iterator.FromSlice([]int{1, 2, 3})
	scanned := iter8.Scan(0, func(acc *int, x int) option.Option[int] {
		*acc = *acc + x
		return option.Some(*acc)
	})
	result6 := scanned.Collect()
	assert.Equal(t, []int{1, 3, 6}, result6)
}
