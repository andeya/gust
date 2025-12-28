package iterator_test

import (
	"strconv"
	"testing"

	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/option"
	"github.com/andeya/gust/pair"
	"github.com/andeya/gust/result"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	a := []int{1, 2, 3}
	iter := iterator.Map(iterator.FromSlice(a), func(x int) int { return 2 * x })

	assert.Equal(t, option.Some(2), iter.Next())
	assert.Equal(t, option.Some(4), iter.Next())
	assert.Equal(t, option.Some(6), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestFilter(t *testing.T) {
	a := []int{0, 1, 2}
	iter := iterator.FromSlice(a).Filter(func(x int) bool { return x > 0 })

	assert.Equal(t, option.Some(1), iter.Next())
	assert.Equal(t, option.Some(2), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestChain(t *testing.T) {
	s1 := iterator.FromSlice([]int{1, 2, 3})
	s2 := iterator.FromSlice([]int{4, 5, 6})
	iter := s1.Chain(s2)

	assert.Equal(t, option.Some(1), iter.Next())
	assert.Equal(t, option.Some(2), iter.Next())
	assert.Equal(t, option.Some(3), iter.Next())
	assert.Equal(t, option.Some(4), iter.Next())
	assert.Equal(t, option.Some(5), iter.Next())
	assert.Equal(t, option.Some(6), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestZip(t *testing.T) {
	s1 := iterator.FromSlice([]int{1, 2, 3})
	s2 := iterator.FromSlice([]int{4, 5, 6})
	iter := iterator.Zip(s1, s2)

	assert.Equal(t, option.Some(pair.Pair[int, int]{A: 1, B: 4}), iter.Next())
	assert.Equal(t, option.Some(pair.Pair[int, int]{A: 2, B: 5}), iter.Next())
	assert.Equal(t, option.Some(pair.Pair[int, int]{A: 3, B: 6}), iter.Next())
	assert.Equal(t, option.None[pair.Pair[int, int]](), iter.Next())
}

func TestEnumerate(t *testing.T) {
	a := []rune{'a', 'b', 'c'}
	iter := iterator.Enumerate(iterator.FromSlice(a))

	assert.Equal(t, option.Some(pair.Pair[uint, rune]{A: 0, B: 'a'}), iter.Next())
	assert.Equal(t, option.Some(pair.Pair[uint, rune]{A: 1, B: 'b'}), iter.Next())
	assert.Equal(t, option.Some(pair.Pair[uint, rune]{A: 2, B: 'c'}), iter.Next())
	assert.Equal(t, option.None[pair.Pair[uint, rune]](), iter.Next())
}

func TestFilterMap(t *testing.T) {
	a := []string{"1", "two", "NaN", "four", "5"}
	iter := iterator.FilterMap(iterator.FromSlice(a), func(s string) option.Option[int] {
		if v, err := strconv.Atoi(s); err == nil {
			return option.Some(v)
		}
		return option.None[int]()
	})

	assert.Equal(t, option.Some(1), iter.Next())
	assert.Equal(t, option.Some(5), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestRetMap(t *testing.T) {
	iter := iterator.RetMap(iterator.FromSlice([]string{"1", "2", "3", "NaN"}), strconv.Atoi)

	// First element: "1" -> Ok(1)
	first := iter.Next()
	assert.True(t, first.IsSome())
	assert.True(t, first.Unwrap().IsOk())
	assert.Equal(t, 1, first.Unwrap().Unwrap())

	// Second element: "2" -> Ok(2)
	second := iter.Next()
	assert.True(t, second.IsSome())
	assert.True(t, second.Unwrap().IsOk())
	assert.Equal(t, 2, second.Unwrap().Unwrap())

	// Third element: "3" -> Ok(3)
	third := iter.Next()
	assert.True(t, third.IsSome())
	assert.True(t, third.Unwrap().IsOk())
	assert.Equal(t, 3, third.Unwrap().Unwrap())

	// Fourth element: "NaN" -> Err
	fourth := iter.Next()
	assert.True(t, fourth.IsSome())
	assert.True(t, fourth.Unwrap().IsErr())

	// Fifth element: None
	assert.Equal(t, option.None[result.Result[int]](), iter.Next())
}

func TestOptMap(t *testing.T) {
	iter := iterator.OptMap(iterator.FromSlice([]string{"1", "2", "3", "NaN"}), func(s string) *int {
		if v, err := strconv.Atoi(s); err == nil {
			return &v
		}
		return nil
	})

	var newInt = func(v int) *int {
		return &v
	}

	assert.Equal(t, option.Some(option.Some(newInt(1))), iter.Next())
	assert.Equal(t, option.Some(option.Some(newInt(2))), iter.Next())
	assert.Equal(t, option.Some(option.Some(newInt(3))), iter.Next())
	assert.Equal(t, option.Some(option.None[*int]()), iter.Next())
	assert.Equal(t, option.None[option.Option[*int]](), iter.Next())
}

func TestIterator_XMap(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	doubled := iter.XMap(func(x int) any { return x * 2 })
	next := doubled.Next()
	assert.True(t, next.IsSome())
	assert.Equal(t, 2, next.Unwrap())
	// Can chain: doubled.Filter(...).Collect()
}

func TestIterator_XFilterMap(t *testing.T) {
	iter := iterator.FromSlice([]string{"1", "two", "NaN", "four", "5"})
	filtered := iter.XFilterMap(func(s string) option.Option[any] {
		if s == "1" {
			return option.Some(any(1))
		}
		if s == "5" {
			return option.Some(any(5))
		}
		return option.None[any]()
	})
	// Can chain: filtered.Filter(...).Collect()
	assert.True(t, filtered.Next().IsSome())
}

func TestIterator_XFlatMap(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	flatMapped := iter.XFlatMap(func(x int) iterator.Iterator[any] {
		return iterator.FromSlice([]any{x, x * 2})
	})
	// Can chain: flatMapped.Filter(...).Collect()
	assert.True(t, flatMapped.Next().IsSome())
}

func TestComplexChain(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	result := iter.
		Filter(func(x int) bool { return x > 2 }).
		Map(func(x int) int { return x * 2 }).
		Take(2).
		Collect()
	assert.Equal(t, []int{6, 8}, result)
}

func TestIteratorExtChaining(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	result := iter.
		Filter(func(x int) bool { return x > 2 }).
		Map(func(x int) int { return x * 2 }).
		Collect()
	assert.Equal(t, []int{6, 8, 10}, result)
}

func TestIteratorExtMap(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	doubled := iter.Map(func(x int) int { return x * 2 })
	assert.Equal(t, option.Some(2), doubled.Next())
	assert.Equal(t, option.Some(4), doubled.Next())
	assert.Equal(t, option.Some(6), doubled.Next())
}

func TestIteratorExtComplexChain(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	result := iter.
		Filter(func(x int) bool { return x%2 == 0 }).
		Map(func(x int) int { return x * x }).
		Collect()
	assert.Equal(t, []int{4, 16}, result)
}

func TestIteratorExtEnumerate(t *testing.T) {
	iter := iterator.FromSlice([]string{"a", "b", "c"})
	enumerated := iterator.Enumerate(iter)
	var result []pair.Pair[uint, string]
	for {
		opt := enumerated.Next()
		if opt.IsNone() {
			break
		}
		result = append(result, opt.Unwrap())
	}
	assert.Equal(t, []pair.Pair[uint, string]{
		{A: 0, B: "a"},
		{A: 1, B: "b"},
		{A: 2, B: "c"},
	}, result)
}

func TestIteratorExtZip(t *testing.T) {
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	iter2 := iterator.FromSlice([]string{"a", "b", "c"})
	zipped := iterator.Zip(iter1, iter2)
	var result []pair.Pair[int, string]
	for {
		opt := zipped.Next()
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
}

// TestZipOneEmpty tests iterator.Zip with one empty iterator
func TestZipOneEmpty(t *testing.T) {
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	iter2 := iterator.Empty[string]()
	zipped := iterator.Zip(iter1, iter2)
	assert.Equal(t, option.None[pair.Pair[int, string]](), zipped.Next())
}

// TestZipBothEmpty tests iterator.Zip with both empty iterators
func TestZipBothEmpty(t *testing.T) {
	iter1 := iterator.Empty[int]()
	iter2 := iterator.Empty[string]()
	zipped := iterator.Zip(iter1, iter2)
	assert.Equal(t, option.None[pair.Pair[int, string]](), zipped.Next())
}

// TestChainFirstEmpty tests Chain with first iterator empty
func TestChainFirstEmpty(t *testing.T) {
	iter1 := iterator.Empty[int]()
	iter2 := iterator.FromSlice([]int{1, 2, 3})
	chained := iter1.Chain(iter2)
	assert.Equal(t, option.Some(1), chained.Next())
	assert.Equal(t, option.Some(2), chained.Next())
	assert.Equal(t, option.Some(3), chained.Next())
	assert.Equal(t, option.None[int](), chained.Next())
}

// TestChainSecondEmpty tests Chain with second iterator empty
func TestChainSecondEmpty(t *testing.T) {
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	iter2 := iterator.Empty[int]()
	chained := iter1.Chain(iter2)
	assert.Equal(t, option.Some(1), chained.Next())
	assert.Equal(t, option.Some(2), chained.Next())
	assert.Equal(t, option.Some(3), chained.Next())
	assert.Equal(t, option.None[int](), chained.Next())
}

// TestChainBothEmpty tests Chain with both iterators empty
func TestChainBothEmpty(t *testing.T) {
	iter1 := iterator.Empty[int]()
	iter2 := iterator.Empty[int]()
	chained := iter1.Chain(iter2)
	assert.Equal(t, option.None[int](), chained.Next())
}

// TestFilterMapAllNone tests iterator.FilterMap where all elements return None
func TestFilterMapAllNone(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	filtered := iterator.FilterMap(iter, func(x int) option.Option[string] {
		return option.None[string]()
	})
	assert.Equal(t, option.None[string](), filtered.Next())
}

// TestFilterMapSomeNone tests iterator.FilterMap with some None results
func TestFilterMapSomeNone(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	filtered := iterator.FilterMap(iter, func(x int) option.Option[int] {
		if x%2 == 0 {
			return option.Some(x * 2)
		}
		return option.None[int]()
	})
	assert.Equal(t, option.Some(4), filtered.Next())
	assert.Equal(t, option.Some(8), filtered.Next())
	assert.Equal(t, option.None[int](), filtered.Next())
}

// TestEnumerateEmpty tests iterator.Enumerate with empty iterator
func TestEnumerateEmpty(t *testing.T) {
	iter := iterator.Empty[int]()
	enumerated := iterator.Enumerate(iter)
	assert.Equal(t, option.None[pair.Pair[uint, int]](), enumerated.Next())
}

// TestEnumerateSizeHint tests iterator.Enumerate SizeHint (covers adapters.go:368-370)
func TestEnumerateSizeHint(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	enumerated := iterator.Enumerate(iter)
	lower, upper := enumerated.SizeHint()
	assert.Equal(t, uint(3), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(3), upper.Unwrap())
}

// TestZipSizeHintEdgeCases tests iterator.Zip SizeHint edge cases
func TestZipSizeHintEdgeCases(t *testing.T) {
	// Test with only upperA.IsSome()
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	iter2 := iterator.Empty[string]()
	zipped := iterator.Zip(iter1, iter2)
	lower, upper := zipped.SizeHint()
	assert.Equal(t, uint(0), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(0), upper.Unwrap())

	// Test with only upperB.IsSome()
	iter3 := iterator.Empty[int]()
	iter4 := iterator.FromSlice([]string{"a", "b"})
	zipped2 := iterator.Zip(iter3, iter4)
	lower2, upper2 := zipped2.SizeHint()
	assert.Equal(t, uint(0), lower2)
	assert.True(t, upper2.IsSome())
	assert.Equal(t, uint(0), upper2.Unwrap())

	// Test with upperAVal < upperBVal
	iter5 := iterator.FromSlice([]int{1, 2})
	iter6 := iterator.FromSlice([]string{"a", "b", "c", "d"})
	zipped3 := iterator.Zip(iter5, iter6)
	lower3, upper3 := zipped3.SizeHint()
	assert.Equal(t, uint(2), lower3)
	assert.True(t, upper3.IsSome())
	assert.Equal(t, uint(2), upper3.Unwrap())

	// Test with upperAVal >= upperBVal
	iter7 := iterator.FromSlice([]int{1, 2, 3, 4})
	iter8 := iterator.FromSlice([]string{"a", "b"})
	zipped4 := iterator.Zip(iter7, iter8)
	lower4, upper4 := zipped4.SizeHint()
	assert.Equal(t, uint(2), lower4)
	assert.True(t, upper4.IsSome())
	assert.Equal(t, uint(2), upper4.Unwrap())

	// Test with neither upperA nor upperB IsSome() (covers adapters.go:301-303)
	// Use iterators that don't provide SizeHint upper bound (infinite iterators)
	iter9 := iterator.Repeat(1)    // iterator.Repeat returns (0, None)
	iter10 := iterator.Repeat("a") // iterator.Repeat returns ("a", None)
	zipped5 := iterator.Zip(iter9, iter10)
	lower5, upper5 := zipped5.SizeHint()
	assert.Equal(t, uint(0), lower5)
	assert.False(t, upper5.IsSome())
}

// TestChainSizeHintEdgeCases tests Chain SizeHint edge cases
func TestChainSizeHintEdgeCases(t *testing.T) {
	// Test with upperA.IsSome() && upperB.IsSome()
	iter1 := iterator.FromSlice([]int{1, 2})
	iter2 := iterator.FromSlice([]int{3, 4})
	chained := iter1.Chain(iter2)
	lower, upper := chained.SizeHint()
	assert.Equal(t, uint(4), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(4), upper.Unwrap())

	// Test with upperA.IsNone() || upperB.IsNone()
	iter3 := iterator.Repeat(1)
	iter4 := iterator.FromSlice([]int{2, 3})
	chained2 := iter3.Chain(iter4)
	lower2, upper2 := chained2.SizeHint()
	assert.Equal(t, uint(2), lower2)
	assert.True(t, upper2.IsNone())
}

// TestIterator_WrapperMethods_Basic tests basic wrapper methods from TestIterator_WrapperMethods
func TestIterator_WrapperMethods_Basic(t *testing.T) {
	// Test XFindMap (covers iterator_methods.go:416-418)
	iter1 := iterator.FromSlice([]string{"lol", "NaN", "2", "5"})
	result1 := iter1.XFindMap(func(s string) option.Option[any] {
		if v, err := strconv.Atoi(s); err == nil {
			return option.Some[any](v)
		}
		return option.None[any]()
	})
	assert.True(t, result1.IsSome())
	assert.Equal(t, 2, result1.UnwrapUnchecked())

	// Test FindMap (covers iterator_methods.go:434-436)
	iter2 := iterator.FromSlice([]string{"lol", "NaN", "2", "5"})
	result2 := iter2.FindMap(func(s string) option.Option[string] {
		if v, err := strconv.Atoi(s); err == nil {
			return option.Some(strconv.Itoa(v))
		}
		return option.None[string]()
	})
	assert.True(t, result2.IsSome())
	assert.Equal(t, "2", result2.UnwrapUnchecked())

	// Test FilterMap (covers iterator_methods.go:699-701)
	iter6 := iterator.FromSlice([]string{"1", "two", "3", "four"})
	filtered := iter6.FilterMap(func(s string) option.Option[string] {
		if v, err := strconv.Atoi(s); err == nil {
			return option.Some(strconv.Itoa(v))
		}
		return option.None[string]()
	})
	result4 := filtered.Collect()
	assert.Equal(t, []string{"1", "3"}, result4)
}
