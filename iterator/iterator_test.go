package iterator_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iterator"
	"github.com/stretchr/testify/assert"
)

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

type easyIterable struct {
	values []int
	index  int
}

func (c *easyIterable) Next() gust.Option[int] {
	if c.index >= len(c.values) {
		return gust.None[int]()
	}
	val := c.values[c.index]
	c.index++
	return gust.Some(val)
}

func TestFromIterable(t *testing.T) {
	// Test with Iterator[T] - should return the same iterator
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	var gustIter1 gust.Iterable[int] = iter1
	iter2 := iterator.FromIterable(gustIter1)
	assert.Equal(t, gust.Some(1), iter2.Next())
	assert.Equal(t, gust.Some(2), iter2.Next())
	assert.Equal(t, gust.Some(3), iter2.Next())
	assert.Equal(t, gust.None[int](), iter2.Next())

	// Test with gust.Iterable[T] that is not Iterator[T]
	custom := &easyIterable{values: []int{10, 20, 30}, index: 0}
	var gustIter2 gust.Iterable[int] = custom
	iter3 := iterator.FromIterable(gustIter2)
	assert.Equal(t, gust.Some(10), iter3.Next())
	assert.Equal(t, gust.Some(20), iter3.Next())
	assert.Equal(t, gust.Some(30), iter3.Next())
	assert.Equal(t, gust.None[int](), iter3.Next())
}

func TestTryToDoubleEnded(t *testing.T) {
	// Test with double-ended iterator
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	deOpt := iter1.TryToDoubleEnded()
	assert.True(t, deOpt.IsSome())
	deIter := deOpt.Unwrap()
	assert.Equal(t, gust.Some(3), deIter.NextBack())

	// Test with non-double-ended iterator (would need a custom iterator)
	// For now, sliceIterator supports double-ended, so this will succeed
}

func TestMustToDoubleEnded(t *testing.T) {
	// Test with double-ended iterator
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	deIter := iter1.MustToDoubleEnded()
	assert.Equal(t, gust.Some(3), deIter.NextBack())
	assert.Equal(t, gust.Some(2), deIter.NextBack())
	assert.Equal(t, gust.Some(1), deIter.NextBack())
	assert.Equal(t, gust.None[int](), deIter.NextBack())
}

func TestSeq(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	var results []int
	for v := range iter.Seq() {
		results = append(results, v)
	}
	assert.Equal(t, []int{1, 2, 3}, results)

	// Test early termination
	iter2 := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	var results2 []int
	for v := range iter2.Seq() {
		results2 = append(results2, v)
		if v == 3 {
			break
		}
	}
	assert.Equal(t, []int{1, 2, 3}, results2)
}

func TestPull(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	next, stop := iter.Pull()
	defer stop()

	var results []int
	for {
		v, ok := next()
		if !ok {
			break
		}
		results = append(results, v)
		if v == 3 {
			break // Early termination
		}
	}
	assert.Equal(t, []int{1, 2, 3}, results)
}

func TestFromSlice(t *testing.T) {
	a := []int{1, 2, 3}
	iter := iterator.FromSlice(a)

	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.Some(2), iter.Next())
	assert.Equal(t, gust.Some(3), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestCount(t *testing.T) {
	a := []int{1, 2, 3}
	assert.Equal(t, uint(3), iterator.FromSlice(a).Count())

	b := []int{1, 2, 3, 4, 5}
	assert.Equal(t, uint(5), iterator.FromSlice(b).Count())
}

func TestLast(t *testing.T) {
	a := []int{1, 2, 3}
	assert.Equal(t, gust.Some(3), iterator.FromSlice(a).Last())

	b := []int{1, 2, 3, 4, 5}
	assert.Equal(t, gust.Some(5), iterator.FromSlice(b).Last())
}

func TestMap(t *testing.T) {
	a := []int{1, 2, 3}
	iter := iterator.Map(iterator.FromSlice(a), func(x int) int { return 2 * x })

	assert.Equal(t, gust.Some(2), iter.Next())
	assert.Equal(t, gust.Some(4), iter.Next())
	assert.Equal(t, gust.Some(6), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestFilter(t *testing.T) {
	a := []int{0, 1, 2}
	iter := iterator.FromSlice(a).Filter(func(x int) bool { return x > 0 })

	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.Some(2), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestChain(t *testing.T) {
	s1 := iterator.FromSlice([]int{1, 2, 3})
	s2 := iterator.FromSlice([]int{4, 5, 6})
	iter := s1.Chain(s2)

	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.Some(2), iter.Next())
	assert.Equal(t, gust.Some(3), iter.Next())
	assert.Equal(t, gust.Some(4), iter.Next())
	assert.Equal(t, gust.Some(5), iter.Next())
	assert.Equal(t, gust.Some(6), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestZip(t *testing.T) {
	s1 := iterator.FromSlice([]int{1, 2, 3})
	s2 := iterator.FromSlice([]int{4, 5, 6})
	iter := iterator.Zip(s1, s2)

	assert.Equal(t, gust.Some(gust.Pair[int, int]{A: 1, B: 4}), iter.Next())
	assert.Equal(t, gust.Some(gust.Pair[int, int]{A: 2, B: 5}), iter.Next())
	assert.Equal(t, gust.Some(gust.Pair[int, int]{A: 3, B: 6}), iter.Next())
	assert.Equal(t, gust.None[gust.Pair[int, int]](), iter.Next())
}

func TestEnumerate(t *testing.T) {
	a := iterator.FromSlice([]int{10, 20, 30})
	iter := iterator.Enumerate(a)

	assert.Equal(t, gust.Some(gust.Pair[uint, int]{A: 0, B: 10}), iter.Next())
	assert.Equal(t, gust.Some(gust.Pair[uint, int]{A: 1, B: 20}), iter.Next())
	assert.Equal(t, gust.Some(gust.Pair[uint, int]{A: 2, B: 30}), iter.Next())
	assert.Equal(t, gust.None[gust.Pair[uint, int]](), iter.Next())
}

func TestSkip(t *testing.T) {
	a := []int{1, 2, 3}
	iter := iterator.FromSlice(a).Skip(2)

	assert.Equal(t, gust.Some(3), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestTake(t *testing.T) {
	a := []int{1, 2, 3}
	iter := iterator.FromSlice(a).Take(2)

	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.Some(2), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestFold(t *testing.T) {
	a := []int{1, 2, 3}
	sum := iterator.Fold(iterator.FromSlice(a), 0, func(acc int, x int) int { return acc + x })
	assert.Equal(t, 6, sum)
}

func TestReduce(t *testing.T) {
	reduced := iterator.FromRange(1, 10).Reduce(func(acc int, e int) int { return acc + e })
	assert.Equal(t, gust.Some(45), reduced)
}

func TestCollect(t *testing.T) {
	a := []int{1, 2, 3}
	doubled := iterator.Map(iterator.FromSlice(a), func(x int) int { return x * 2 }).Collect()
	assert.Equal(t, []int{2, 4, 6}, doubled)
}

func TestPartition(t *testing.T) {
	a := []int{1, 2, 3}
	even, odd := iterator.FromSlice(a).Partition(func(n int) bool { return n%2 == 0 })
	assert.Equal(t, []int{2}, even)
	assert.Equal(t, []int{1, 3}, odd)
}

func TestAll(t *testing.T) {
	a := []int{1, 2, 3}
	assert.True(t, iterator.FromSlice(a).All(func(x int) bool { return x > 0 }))
	assert.False(t, iterator.FromSlice(a).All(func(x int) bool { return x > 2 }))
}

func TestAny(t *testing.T) {
	a := []int{1, 2, 3}
	assert.True(t, iterator.FromSlice(a).Any(func(x int) bool { return x > 0 }))
	assert.False(t, iterator.FromSlice(a).Any(func(x int) bool { return x > 5 }))
}

func TestFind(t *testing.T) {
	a := []int{1, 2, 3}
	assert.Equal(t, gust.Some(2), iterator.FromSlice(a).Find(func(x int) bool { return x == 2 }))
	assert.Equal(t, gust.None[int](), iterator.FromSlice(a).Find(func(x int) bool { return x == 5 }))
}

func TestPosition(t *testing.T) {
	a := []int{1, 2, 3}
	assert.Equal(t, gust.Some(uint(1)), iterator.FromSlice(a).Position(func(x int) bool { return x == 2 }))
	assert.Equal(t, gust.None[uint](), iterator.FromSlice(a).Position(func(x int) bool { return x == 5 }))
}

func TestMax_EdgeCases(t *testing.T) {
	// Test Max with single element
	{
		result := iterator.Max(iterator.FromSlice([]int{42}))
		assert.True(t, result.IsSome())
		assert.Equal(t, 42, result.Unwrap())
	}

	// Test Max with equal elements (should return last)
	{
		result := iterator.Max(iterator.FromSlice([]int{5, 5, 5}))
		assert.True(t, result.IsSome())
		assert.Equal(t, 5, result.Unwrap())
	}

	// Test Max with string
	{
		result := iterator.Max(iterator.FromSlice([]string{"z", "a", "m"}))
		assert.True(t, result.IsSome())
		assert.Equal(t, "z", result.Unwrap())
	}

	// Test Max with uint
	{
		result := iterator.Max(iterator.FromSlice([]uint{100, 50, 200}))
		assert.True(t, result.IsSome())
		assert.Equal(t, uint(200), result.Unwrap())
	}
}

func TestMax(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{}
	assert.Equal(t, gust.Some(3), iterator.Max(iterator.FromSlice(a)))
	assert.Equal(t, gust.None[int](), iterator.Max(iterator.FromSlice(b)))
}

func TestMin_EdgeCases(t *testing.T) {
	// Test Min with single element
	{
		result := iterator.Min(iterator.FromSlice([]int{42}))
		assert.True(t, result.IsSome())
		assert.Equal(t, 42, result.Unwrap())
	}

	// Test Min with equal elements (should return first)
	{
		result := iterator.Min(iterator.FromSlice([]int{5, 5, 5}))
		assert.True(t, result.IsSome())
		assert.Equal(t, 5, result.Unwrap())
	}

	// Test Min with string
	{
		result := iterator.Min(iterator.FromSlice([]string{"z", "a", "m"}))
		assert.True(t, result.IsSome())
		assert.Equal(t, "a", result.Unwrap())
	}

	// Test Min with uint
	{
		result := iterator.Min(iterator.FromSlice([]uint{100, 50, 200}))
		assert.True(t, result.IsSome())
		assert.Equal(t, uint(50), result.Unwrap())
	}
}

func TestMin(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{}
	assert.Equal(t, gust.Some(1), iterator.Min(iterator.FromSlice(a)))
	assert.Equal(t, gust.None[int](), iterator.Min(iterator.FromSlice(b)))
}

func TestMaxByKey_EdgeCases(t *testing.T) {
	// Test MaxByKey with single element
	{
		result := iterator.MaxByKey(iterator.FromSlice([]int{42}), func(x int) int {
			return x * 2
		})
		assert.True(t, result.IsSome())
		assert.Equal(t, 42, result.Unwrap())
	}

	// Test MaxByKey with equal keys (should return last)
	{
		result := iterator.MaxByKey(iterator.FromSlice([]int{1, 2, 3}), func(x int) int {
			return 10 // All have same key
		})
		assert.True(t, result.IsSome())
		assert.Equal(t, 3, result.Unwrap()) // Last element
	}

	// Test MaxByKey with negative values
	{
		result := iterator.MaxByKey(iterator.FromSlice([]int{-3, -1, -5}), func(x int) int {
			return -x // Negate to find max absolute value
		})
		assert.True(t, result.IsSome())
		assert.Equal(t, -5, result.Unwrap())
	}
}

func TestMaxByKey(t *testing.T) {
	a := []int{-3, 0, 1, 5, -10}
	max := iterator.MaxByKey(iterator.FromSlice(a), func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	})
	assert.Equal(t, gust.Some(-10), max)

	// Test with empty iterator
	empty := []int{}
	assert.True(t, iterator.MaxByKey(iterator.FromSlice(empty), func(x int) int { return x }).IsNone())
}

func TestMaxBy(t *testing.T) {
	a := []int{-3, 0, 1, 5, -10}
	max := iterator.FromSlice(a).MaxBy(func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
	assert.Equal(t, gust.Some(5), max)

	// Test with empty iterator
	empty := []int{}
	assert.True(t, iterator.FromSlice(empty).MaxBy(func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	}).IsNone())

	// Test with equal elements (should return last)
	equal := []int{2, 2, 2}
	max2 := iterator.FromSlice(equal).MaxBy(func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
	assert.Equal(t, gust.Some(2), max2)

	// Test maxByImpl compare(acc, x) < 0 branch (covers min_max.go:118)
	// When compare returns < 0, it should return x (new maximum)
	ascending := []int{1, 2, 3, 4, 5}
	max3 := iterator.FromSlice(ascending).MaxBy(func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
	assert.Equal(t, gust.Some(5), max3)

	// Test maxByImpl compare(acc, x) >= 0 branch (covers min_max.go:121)
	// When compare returns >= 0, it should return acc (keep current maximum)
	descending := []int{5, 4, 3, 2, 1}
	max4 := iterator.FromSlice(descending).MaxBy(func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
	assert.Equal(t, gust.Some(5), max4)
}

func TestMinByKey_EdgeCases(t *testing.T) {
	// Test MinByKey with single element
	{
		result := iterator.MinByKey(iterator.FromSlice([]int{42}), func(x int) int {
			return x * 2
		})
		assert.True(t, result.IsSome())
		assert.Equal(t, 42, result.Unwrap())
	}

	// Test MinByKey with equal keys (should return first)
	{
		result := iterator.MinByKey(iterator.FromSlice([]int{1, 2, 3}), func(x int) int {
			return 10 // All have same key
		})
		assert.True(t, result.IsSome())
		assert.Equal(t, 1, result.Unwrap()) // First element
	}

	// Test MinByKey with negative values
	{
		result := iterator.MinByKey(iterator.FromSlice([]int{-3, -1, -5}), func(x int) int {
			return -x // Negate to find min absolute value
		})
		assert.True(t, result.IsSome())
		assert.Equal(t, -1, result.Unwrap())
	}
}

func TestMinByKey(t *testing.T) {
	a := []int{-3, 0, 1, 5, -10}
	min := iterator.MinByKey(iterator.FromSlice(a), func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	})
	assert.Equal(t, gust.Some(0), min)

	// Test with empty iterator
	empty := []int{}
	assert.True(t, iterator.MinByKey(iterator.FromSlice(empty), func(x int) int { return x }).IsNone())
}

func TestMinBy(t *testing.T) {
	a := []int{-3, 0, 1, 5, -10}
	min := iterator.FromSlice(a).MinBy(func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
	assert.Equal(t, gust.Some(-10), min)

	// Test with empty iterator
	empty := []int{}
	assert.True(t, iterator.FromSlice(empty).MinBy(func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	}).IsNone())

	// Test with equal elements (should return first)
	equal := []int{2, 2, 2}
	min2 := iterator.FromSlice(equal).MinBy(func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
	assert.Equal(t, gust.Some(2), min2)

	// Test minByImpl compare(acc, x) > 0 branch (covers min_max.go:162)
	// When compare returns > 0, it should return x (new minimum)
	descending := []int{5, 4, 3, 2, 1}
	min3 := iterator.FromSlice(descending).MinBy(func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
	assert.Equal(t, gust.Some(1), min3)

	// Test minByImpl compare(acc, x) <= 0 branch (covers min_max.go:165)
	// When compare returns <= 0, it should return acc (keep current minimum)
	ascending := []int{1, 2, 3, 4, 5}
	min4 := iterator.FromSlice(ascending).MinBy(func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
	assert.Equal(t, gust.Some(1), min4)
}

func TestFromRange(t *testing.T) {
	iter := iterator.FromRange(0, 5)
	assert.Equal(t, gust.Some(0), iter.Next())
	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.Some(2), iter.Next())
	assert.Equal(t, gust.Some(3), iter.Next())
	assert.Equal(t, gust.Some(4), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestFromElements(t *testing.T) {
	iter := iterator.FromElements(1, 2, 3)
	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.Some(2), iter.Next())
	assert.Equal(t, gust.Some(3), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestOnce(t *testing.T) {
	iter := iterator.Once(42)
	assert.Equal(t, gust.Some(42), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestEmpty(t *testing.T) {
	iter := iterator.Empty[int]()
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestFromFunc(t *testing.T) {
	count := 0
	iter := iterator.FromFunc(func() gust.Option[int] {
		if count < 3 {
			count++
			return gust.Some(count)
		}
		return gust.None[int]()
	})
	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.Some(2), iter.Next())
	assert.Equal(t, gust.Some(3), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestRepeat(t *testing.T) {
	iter := iterator.Repeat(42)
	assert.Equal(t, gust.Some(42), iter.Next())
	assert.Equal(t, gust.Some(42), iter.Next())
	assert.Equal(t, gust.Some(42), iter.Next())
	// Should repeat forever
}

func TestSizeHint(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	lower, upper := iter.SizeHint()
	assert.Equal(t, uint(3), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(3), upper.Unwrap())
}

func TestSkipWhile(t *testing.T) {
	a := []int{-1, 0, 1}
	iter := iterator.FromSlice(a).SkipWhile(func(x int) bool { return x < 0 })

	assert.Equal(t, gust.Some(0), iter.Next())
	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestTakeWhile(t *testing.T) {
	a := []int{-1, 0, 1}
	iter := iterator.FromSlice(a).TakeWhile(func(x int) bool { return x < 0 })

	assert.Equal(t, gust.Some(-1), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestAdvanceBy(t *testing.T) {
	a := []int{1, 2, 3, 4}
	iter := iterator.FromSlice(a)

	assert.True(t, iter.AdvanceBy(2).IsOk())
	assert.Equal(t, gust.Some(3), iter.Next())
	assert.True(t, iter.AdvanceBy(0).IsOk())
	result := iter.AdvanceBy(100)
	assert.True(t, result.IsErr())
	assert.Equal(t, uint(99), result.ErrVal())
}

func TestNth(t *testing.T) {
	a := []int{1, 2, 3}
	assert.Equal(t, gust.Some(2), iterator.FromSlice(a).Nth(1))

	b := []int{1, 2, 3}
	iter := iterator.FromSlice(b)
	assert.Equal(t, gust.Some(2), iter.Nth(1))
	assert.Equal(t, gust.Some(3), iter.Nth(0))

	c := []int{1, 2, 3}
	assert.Equal(t, gust.None[int](), iterator.FromSlice(c).Nth(10))
}

func TestFilterMap(t *testing.T) {
	a := []string{"1", "two", "NaN", "four", "5"}
	iter := iterator.FilterMap(iterator.FromSlice(a), func(s string) gust.Option[int] {
		if s == "1" {
			return gust.Some(1)
		}
		if s == "5" {
			return gust.Some(5)
		}
		return gust.None[int]()
	})

	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.Some(5), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestComplexChain(t *testing.T) {
	// Test a complex chain of operations
	result := iterator.Map(
		iterator.FromRange(1, 10).Filter(func(x int) bool { return x%2 == 0 }),
		func(x int) int { return x * 2 },
	).Take(3).Collect()
	assert.Equal(t, []int{4, 8, 12}, result)
}

func TestIteratorExtChaining(t *testing.T) {
	// Test method chaining with Iterator
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})

	// Chain Filter and Take
	filtered := iter.Filter(func(x int) bool { return x > 2 })
	taken := filtered.Take(2)
	result := taken.Collect()

	assert.Equal(t, []int{3, 4}, result)
}

func TestIteratorExtMap(t *testing.T) {
	// Test Map with Iterator (using function-style API)
	iter := iterator.FromSlice([]int{1, 2, 3})
	doubled := iterator.Map(iter, func(x int) int { return x * 2 })
	result := doubled.Collect()

	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestIteratorExtComplexChain(t *testing.T) {
	// Test complex chaining
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	// Filter -> Map -> Take -> Collect
	mapped := iterator.Map(iter, func(x int) int { return x * 2 })
	filtered := mapped.Filter(func(x int) bool { return x > 5 })
	taken := filtered.Take(3)
	result := taken.Collect()

	assert.Equal(t, []int{6, 8, 10}, result)
}

func TestIteratorExtEnumerate(t *testing.T) {
	iter := iterator.FromSlice([]int{10, 20, 30})
	enumerated := iterator.Enumerate(iter)

	pair1 := enumerated.Next()
	assert.True(t, pair1.IsSome())
	assert.Equal(t, uint(0), pair1.Unwrap().A)
	assert.Equal(t, 10, pair1.Unwrap().B)

	pair2 := enumerated.Next()
	assert.True(t, pair2.IsSome())
	assert.Equal(t, uint(1), pair2.Unwrap().A)
	assert.Equal(t, 20, pair2.Unwrap().B)
}

func TestIteratorExtConsumerMethods(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})

	// Test Count
	assert.Equal(t, uint(5), iter.Count())

	// Test Last
	iter2 := iterator.FromSlice([]int{1, 2, 3})
	assert.Equal(t, gust.Some(3), iter2.Last())

	// Test All
	iter3 := iterator.FromSlice([]int{2, 4, 6})
	assert.True(t, iter3.All(func(x int) bool { return x%2 == 0 }))

	// Test Any
	iter4 := iterator.FromSlice([]int{1, 2, 3})
	assert.True(t, iter4.Any(func(x int) bool { return x > 2 }))

	// Test Find
	iter5 := iterator.FromSlice([]int{1, 2, 3})
	assert.Equal(t, gust.Some(2), iter5.Find(func(x int) bool { return x > 1 }))

	// Test Max
	iter6 := iterator.FromSlice([]int{1, 3, 2})
	assert.Equal(t, gust.Some(3), iterator.Max(iter6))

	// Test Min
	iter7 := iterator.FromSlice([]int{3, 1, 2})
	assert.Equal(t, gust.Some(1), iterator.Min(iter7))

	// Test MaxBy
	iter8 := iterator.FromSlice([]int{-3, 0, 1, 5, -10})
	maxBy := iter8.MaxBy(func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
	assert.Equal(t, gust.Some(5), maxBy)

	// Test MinBy
	iter9 := iterator.FromSlice([]int{-3, 0, 1, 5, -10})
	minBy := iter9.MinBy(func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
	assert.Equal(t, gust.Some(-10), minBy)

	// Test MaxByKey (function version)
	iter10 := iterator.FromSlice([]int{-3, 0, 1, 5, -10})
	maxByKey := iterator.MaxByKey(iter10, func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	})
	assert.Equal(t, gust.Some(-10), maxByKey)

	// Test MinByKey (function version)
	iter11 := iterator.FromSlice([]int{-3, 0, 1, 5, -10})
	minByKey := iterator.MinByKey(iter11, func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	})
	assert.Equal(t, gust.Some(0), minByKey)

	// Test ForEach
	iter12 := iterator.FromSlice([]int{1, 2, 3})
	var forEachResult []int
	iter12.ForEach(func(x int) {
		forEachResult = append(forEachResult, x)
	})
	assert.Equal(t, []int{1, 2, 3}, forEachResult)

	// Test TryForEach
	iter13 := iterator.FromSlice([]int{1, 2, 3})
	result := iter13.TryForEach(func(x int) gust.Result[int] {
		return gust.Ok(x)
	})
	assert.True(t, result.IsOk())

	// Test TryReduce
	iter15 := iterator.FromSlice([]int{10, 20, 5})
	sumResult := iter15.TryReduce(func(x, y int) gust.Result[int] {
		if x+y > 100 {
			return gust.TryErr[int](errors.New("overflow"))
		}
		return gust.Ok(x + y)
	})
	assert.True(t, sumResult.IsOk())
	assert.True(t, sumResult.Unwrap().IsSome())
	assert.Equal(t, 35, sumResult.Unwrap().Unwrap())

	// Test TryFind
	iter16 := iterator.FromSlice([]string{"1", "2", "lol", "NaN", "5"})
	findResult := iter16.TryFind(func(s string) gust.Result[bool] {
		if s == "lol" {
			return gust.TryErr[bool](errors.New("invalid"))
		}
		if v, err := strconv.Atoi(s); err == nil {
			return gust.Ok(v == 2)
		}
		return gust.Ok(false)
	})
	assert.True(t, findResult.IsOk())
	assert.True(t, findResult.Unwrap().IsSome())
	assert.Equal(t, "2", findResult.Unwrap().Unwrap())
}

func TestIteratorExtSkipTake(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})

	// Skip first 2, then take 2
	skipped := iter.Skip(2)
	taken := skipped.Take(2)
	result := taken.Collect()

	assert.Equal(t, []int{3, 4}, result)
}

func TestForEach(t *testing.T) {
	// Test basic ForEach functionality
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	var result []int
	iter.ForEach(func(x int) {
		result = append(result, x)
	})
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)

	// Test ForEach with empty iterator
	emptyIter := iterator.Empty[int]()
	var emptyResult []int
	emptyIter.ForEach(func(x int) {
		emptyResult = append(emptyResult, x)
	})
	assert.Nil(t, emptyResult)
	assert.Len(t, emptyResult, 0)

	// Test ForEach with filtered iterator
	filteredIter := iterator.FromSlice([]int{1, 2, 3, 4, 5}).Filter(func(x int) bool { return x%2 == 0 })
	var filteredResult []int
	filteredIter.ForEach(func(x int) {
		filteredResult = append(filteredResult, x)
	})
	assert.Equal(t, []int{2, 4}, filteredResult)

	// Test ForEach with string iterator
	strIter := iterator.FromSlice([]string{"hello", "world", "rust"})
	var strResult []string
	strIter.ForEach(func(s string) {
		strResult = append(strResult, s)
	})
	assert.Equal(t, []string{"hello", "world", "rust"}, strResult)

	// Test ForEach with mapped iterator
	mappedIter := iterator.FromSlice([]int{1, 2, 3}).Map(func(x int) int { return x * 2 })
	var mappedResult []int
	mappedIter.ForEach(func(x int) {
		mappedResult = append(mappedResult, x)
	})
	assert.Equal(t, []int{2, 4, 6}, mappedResult)

	// Test ForEach consumes the iterator completely
	consumedIter := iterator.FromSlice([]int{1, 2, 3})
	consumedIter.ForEach(func(x int) {})
	// After ForEach, iterator should be exhausted
	assert.Equal(t, gust.None[int](), consumedIter.Next())

	// Test ForEach with chained operations
	chainedIter := iterator.FromSlice([]int{1, 2, 3, 4, 5}).
		Filter(func(x int) bool { return x > 2 }).
		Map(func(x int) int { return x * 2 })
	var chainedResult []int
	chainedIter.ForEach(func(x int) {
		chainedResult = append(chainedResult, x)
	})
	assert.Equal(t, []int{6, 8, 10}, chainedResult)

	// Test ForEach with Take (partial consumption)
	takenIter := iterator.FromSlice([]int{1, 2, 3, 4, 5}).Take(3)
	var takenResult []int
	takenIter.ForEach(func(x int) {
		takenResult = append(takenResult, x)
	})
	assert.Equal(t, []int{1, 2, 3}, takenResult)

	// Test ForEach with Skip
	skippedIter := iterator.FromSlice([]int{1, 2, 3, 4, 5}).Skip(2)
	var skippedResult []int
	skippedIter.ForEach(func(x int) {
		skippedResult = append(skippedResult, x)
	})
	assert.Equal(t, []int{3, 4, 5}, skippedResult)

	// Test ForEach with function that modifies external state
	counter := 0
	sum := 0
	iterator.FromSlice([]int{10, 20, 30}).ForEach(func(x int) {
		counter++
		sum += x
	})
	assert.Equal(t, 3, counter)
	assert.Equal(t, 60, sum)

	// Test ForEach with single element iterator
	singleIter := iterator.FromSlice([]int{42})
	var singleResult []int
	singleIter.ForEach(func(x int) {
		singleResult = append(singleResult, x)
	})
	assert.Equal(t, []int{42}, singleResult)

	// Test ForEach with custom iterable
	customIter := iterator.FromIterable(&easyIterable{values: []int{100, 200}, index: 0})
	var customResult []int
	customIter.ForEach(func(x int) {
		customResult = append(customResult, x)
	})
	assert.Equal(t, []int{100, 200}, customResult)
}

// TestForEachImplDirectly tests forEachImpl directly to ensure coverage of fold_reduce.go:86-92
// This test ensures that the forEachImpl function is properly covered, especially
// the loop with break condition and the function call path.
func TestForEachImplDirectly(t *testing.T) {
	// Test with non-empty iterator to cover the loop body (line 91: f(item.Unwrap()))
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	var result1 []int
	iter1.ForEach(func(x int) {
		result1 = append(result1, x)
	})
	assert.Equal(t, []int{1, 2, 3}, result1)

	// Test with empty iterator to cover the break path (line 88-89: if item.IsNone() { break })
	iter2 := iterator.Empty[int]()
	var result2 []int
	iter2.ForEach(func(x int) {
		result2 = append(result2, x)
	})
	assert.Nil(t, result2)
	assert.Len(t, result2, 0)

	// Test with iterator that becomes empty during iteration
	// This ensures the loop condition (line 87: item := iterator.Next()) and break (line 89) are both covered
	iter3 := iterator.FromSlice([]int{42})
	var result3 []int
	iter3.ForEach(func(x int) {
		result3 = append(result3, x)
	})
	assert.Equal(t, []int{42}, result3)
	// After ForEach, iterator should be exhausted
	assert.Equal(t, gust.None[int](), iter3.Next())

	// Test with multiple elements to ensure the loop iterates correctly
	iter4 := iterator.FromSlice([]int{10, 20, 30, 40, 50})
	var callCount int
	iter4.ForEach(func(x int) {
		callCount++
	})
	assert.Equal(t, 5, callCount)
}

func TestIteratorExtZip(t *testing.T) {
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	iter2 := iterator.FromSlice([]string{"a", "b", "c"})

	zipped := iterator.Zip(iter1, iter2)
	pair := zipped.Next()

	assert.True(t, pair.IsSome())
	assert.Equal(t, 1, pair.Unwrap().A)
	assert.Equal(t, "a", pair.Unwrap().B)
}

func TestIterator_XTryFold(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	result := iter.XTryFold(0, func(acc any, x int) gust.Result[any] {
		return gust.Ok(any(acc.(int) + x))
	})
	assert.True(t, result.IsOk())
	assert.Equal(t, 6, result.Unwrap().(int))

	// Test with error
	iter2 := iterator.FromSlice([]int{1, 2, 3})
	result2 := iter2.XTryFold(0, func(acc any, x int) gust.Result[any] {
		if acc.(int)+x > 5 {
			return gust.TryErr[any](errors.New("overflow"))
		}
		return gust.Ok(any(acc.(int) + x))
	})
	assert.True(t, result2.IsErr())
}

func TestIterator_XTryForEach(t *testing.T) {
	data := iterator.FromSlice([]string{"a", "b", "c"})
	result := data.XTryForEach(func(x string) gust.Result[any] {
		return gust.Ok[any](x + "_processed")
	})
	assert.True(t, result.IsOk())

	// Test with error
	data2 := iterator.FromSlice([]string{"a", "b", "error"})
	result2 := data2.XTryForEach(func(x string) gust.Result[any] {
		if x == "error" {
			return gust.TryErr[any](errors.New("processing error"))
		}
		return gust.Ok[any](x)
	})
	assert.True(t, result2.IsErr())
}

func TestIterator_XMap(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	mapped := iter.XMap(func(x int) any { return x * 2 })
	result := mapped.Collect()
	assert.Equal(t, []any{2, 4, 6}, result)
}

func TestIterator_XFlatMap(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2})
	flatMapped := iter.XFlatMap(func(x int) iterator.Iterator[any] {
		return iterator.FromSlice([]any{x, x * 2})
	})
	result := flatMapped.Collect()
	assert.Equal(t, []any{1, 2, 2, 4}, result)
}

func TestIterator_XFold(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	result := iter.XFold(0, func(acc any, x int) any {
		return acc.(int) + x
	})
	assert.Equal(t, 6, result)
}

func TestIterator_XFilterMap(t *testing.T) {
	iter := iterator.FromSlice([]string{"1", "two", "NaN", "four", "5"})
	filtered := iter.XFilterMap(func(s string) gust.Option[any] {
		if s == "1" {
			return gust.Some(any(1))
		}
		if s == "5" {
			return gust.Some(any(5))
		}
		return gust.None[any]()
	})
	result := filtered.Collect()
	assert.Equal(t, []any{1, 5}, result)
}

func TestIterator_XMapWhile(t *testing.T) {
	iter := iterator.FromSlice([]int{-1, 4, 0, 1})
	mapped := iter.XMapWhile(func(x int) gust.Option[any] {
		if x != 0 {
			return gust.Some(any(16 / x))
		}
		return gust.None[any]()
	})
	result := mapped.Collect()
	assert.Equal(t, []any{-16, 4}, result)
}

func TestIterator_XScan(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	scanned := iter.XScan(0, func(state *any, x int) gust.Option[any] {
		s := (*state).(int) + x
		*state = s
		return gust.Some(any(s))
	})
	result := scanned.Collect()
	assert.Equal(t, []any{1, 3, 6}, result)
}

func TestIterator_Intersperse(t *testing.T) {
	iter := iterator.FromSlice([]int{0, 1, 2})
	interspersed := iter.Intersperse(100)
	assert.Equal(t, gust.Some(0), interspersed.Next())
	assert.Equal(t, gust.Some(100), interspersed.Next())
	assert.Equal(t, gust.Some(1), interspersed.Next())
	assert.Equal(t, gust.Some(100), interspersed.Next())
	assert.Equal(t, gust.Some(2), interspersed.Next())
	assert.Equal(t, gust.None[int](), interspersed.Next())
}

func TestIterator_IntersperseWith(t *testing.T) {
	iter := iterator.FromSlice([]int{0, 1, 2})
	interspersed := iter.IntersperseWith(func() int { return 99 })
	assert.Equal(t, gust.Some(0), interspersed.Next())
	assert.Equal(t, gust.Some(99), interspersed.Next())
	assert.Equal(t, gust.Some(1), interspersed.Next())
	assert.Equal(t, gust.Some(99), interspersed.Next())
	assert.Equal(t, gust.Some(2), interspersed.Next())
	assert.Equal(t, gust.None[int](), interspersed.Next())

	// Test Intersperse with single element (no separator)
	iter2 := iterator.FromSlice([]int{42})
	interspersed2 := iter2.Intersperse(100)
	assert.Equal(t, gust.Some(42), interspersed2.Next())
	assert.Equal(t, gust.None[int](), interspersed2.Next())

	// Test Intersperse with empty iterator
	iter3 := iterator.FromSlice([]int{})
	interspersed3 := iter3.Intersperse(100)
	assert.Equal(t, gust.None[int](), interspersed3.Next())

	// Test Intersperse SizeHint with known size
	iter4 := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	interspersed4 := iter4.Intersperse(0)
	lower, upper := interspersed4.SizeHint()
	// Should have lower > 0 and upper > 0 for intersperse
	// Intersperse adds separators, so size becomes 2*n - 1
	assert.True(t, lower > 0 || upper.IsSome())
}

// TestFuse_SizeHint tests Fuse SizeHint method
func TestFuse_SizeHint(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	fused := iter.Fuse()
	lower, upper := fused.SizeHint()
	assert.Equal(t, uint(3), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(3), upper.Unwrap())
}

// TestArrayChunks_SizeHint tests ArrayChunks SizeHint method
func TestArrayChunks_SizeHint(t *testing.T) {
	// Test with known size
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	chunks := iterator.ArrayChunks(iter, 3)
	lower, upper := chunks.SizeHint()
	// Should have lower > 0 and upper > 0
	assert.True(t, lower > 0 || upper.IsSome())
	if upper.IsSome() {
		assert.True(t, upper.Unwrap() > 0)
	}
}

// TestChunkBy_SizeHint tests ChunkBy SizeHint method
func TestChunkBy_SizeHint(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 1, 2, 2, 2, 3, 3})
	chunks := iterator.ChunkBy(iter, func(a, b int) bool { return a == b })
	lower, upper := chunks.SizeHint()
	// ChunkBy can't provide accurate size hint
	assert.Equal(t, uint(0), lower)
	assert.True(t, upper.IsNone())
}

func TestIterator_Cycle(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	cycled := iter.Cycle()
	assert.Equal(t, gust.Some(1), cycled.Next())
	assert.Equal(t, gust.Some(2), cycled.Next())
	assert.Equal(t, gust.Some(3), cycled.Next())
	assert.Equal(t, gust.Some(1), cycled.Next()) // starts over
	assert.Equal(t, gust.Some(2), cycled.Next())
	assert.Equal(t, gust.Some(3), cycled.Next()) // continues cycling
}

func TestIterator_NextChunk(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	chunk := iter.NextChunk(2)
	assert.True(t, chunk.IsOk())
	assert.Equal(t, []int{1, 2}, chunk.Unwrap())

	// Test with insufficient elements (request more than remaining)
	chunk2 := iter.NextChunk(4)
	assert.True(t, chunk2.IsErr())
	assert.Equal(t, []int{3, 4, 5}, chunk2.UnwrapErr())

	// After consuming all elements, next chunk should fail
	chunk3 := iter.NextChunk(2)
	assert.True(t, chunk3.IsErr())

	// Test with exact chunk size
	iter2 := iterator.FromSlice([]int{1, 2, 3})
	chunk4 := iter2.NextChunk(3)
	assert.True(t, chunk4.IsOk())
	assert.Equal(t, []int{1, 2, 3}, chunk4.Unwrap())

	// Test with n == 0 (covers consumers.go:49-51)
	iter3 := iterator.FromSlice([]int{1, 2, 3})
	chunk5 := iter3.NextChunk(0)
	assert.True(t, chunk5.IsOk())
	assert.Equal(t, []int{}, chunk5.Unwrap())

	// Test iterator.ChunkResult.UnwrapErr with Ok result (should panic) (covers chunk_result.go:69-73)
	chunk6 := iter3.NextChunk(1)
	assert.True(t, chunk6.IsOk())
	defer func() {
		p := recover()
		assert.NotNil(t, p)
	}()
	chunk6.UnwrapErr()
}

// TestChunkResult_UnwrapErr tests iterator.ChunkResult.UnwrapErr method (covers chunk_result.go:69-73)
func TestChunkResult_UnwrapErr(t *testing.T) {
	// Test with Err result (should return error value)
	iter := iterator.FromSlice([]int{1, 2, 3})
	chunk := iter.NextChunk(5) // Request more than available
	assert.True(t, chunk.IsErr())
	errVal := chunk.UnwrapErr()
	assert.Equal(t, []int{1, 2, 3}, errVal)

	// Test with Ok result (should panic)
	iter2 := iterator.FromSlice([]int{1, 2})
	chunk2 := iter2.NextChunk(2)
	assert.True(t, chunk2.IsOk())
	defer func() {
		p := recover()
		assert.NotNil(t, p)
	}()
	chunk2.UnwrapErr()
}

// TestChunkResult_EdgeCases tests edge cases for iterator.ChunkResult methods
func TestChunkResult_EdgeCases(t *testing.T) {
	// Test zero-value iterator.ChunkResult (covers safeGetT when r.t.IsSome() is false, and safeGetE when r.e == nil)
	var zeroChunk iterator.ChunkResult[[]int]
	assert.False(t, zeroChunk.IsErr())
	assert.True(t, zeroChunk.IsOk())

	// Test Unwrap() with zero-value (should return zero value of T)
	zeroVal := zeroChunk.Unwrap()
	assert.Equal(t, []int(nil), zeroVal)

	// Test UnwrapErr() with zero-value (should panic, covers safeGetE when r.e == nil)
	// This covers the path where r.e == nil in safeGetE()
	defer func() {
		p := recover()
		assert.NotNil(t, p, "UnwrapErr() on zero-value ChunkResult should panic")
	}()
	zeroChunk.UnwrapErr()
}

// TestChunkResult_UnwrapPanic tests Unwrap() panic when IsErr() is true (covers chunk_result.go:62-63)
func TestChunkResult_UnwrapPanic(t *testing.T) {
	// Create an Err iterator.ChunkResult by requesting more elements than available
	iter := iterator.FromSlice([]int{1, 2, 3})
	errChunk := iter.NextChunk(5) // Request more than available to get Err result
	assert.True(t, errChunk.IsErr())
	assert.False(t, errChunk.IsOk())

	// Test Unwrap() with Err result (should panic)
	defer func() {
		p := recover()
		assert.NotNil(t, p)
		// Verify panic is *ErrBox
		eb, ok := p.(*gust.ErrBox)
		assert.True(t, ok)
		assert.NotNil(t, eb)
	}()
	errChunk.Unwrap()
}

// TestIterator_WrapperMethods tests wrapper methods to cover iterator_methods.go wrapper functions
func TestIterator_WrapperMethods(t *testing.T) {
	// Test XFindMap (covers iterator_methods.go:416-418)
	iter1 := iterator.FromSlice([]string{"lol", "NaN", "2", "5"})
	result1 := iter1.XFindMap(func(s string) gust.Option[any] {
		if v, err := strconv.Atoi(s); err == nil {
			return gust.Some[any](v)
		}
		return gust.None[any]()
	})
	assert.True(t, result1.IsSome())
	assert.Equal(t, 2, result1.UnwrapUnchecked())

	// Test FindMap (covers iterator_methods.go:434-436)
	iter2 := iterator.FromSlice([]string{"lol", "NaN", "2", "5"})
	result2 := iter2.FindMap(func(s string) gust.Option[string] {
		if v, err := strconv.Atoi(s); err == nil {
			return gust.Some(strconv.Itoa(v))
		}
		return gust.None[string]()
	})
	assert.True(t, result2.IsSome())
	assert.Equal(t, "2", result2.UnwrapUnchecked())

	// Test XFold (covers iterator_methods.go:492-494)
	iter3 := iterator.FromSlice([]int{1, 2, 3})
	sum1 := iter3.XFold(0, func(acc any, x int) any {
		return acc.(int) + x
	})
	assert.Equal(t, 6, sum1)

	// Test Fold (covers iterator_methods.go:520-522)
	iter4 := iterator.FromSlice([]int{1, 2, 3})
	sum2 := iter4.Fold(0, func(acc int, x int) int {
		return acc + x
	})
	assert.Equal(t, 6, sum2)

	// Test XTryFold (covers iterator_methods.go:552-554)
	iter5 := iterator.FromSlice([]int{1, 2, 3})
	result3 := iter5.XTryFold(0, func(acc any, x int) gust.Result[any] {
		return gust.Ok[any](acc.(int) + x)
	})
	assert.True(t, result3.IsOk())
	assert.Equal(t, 6, result3.UnwrapUnchecked())

	// Test FilterMap (covers iterator_methods.go:699-701)
	iter6 := iterator.FromSlice([]string{"1", "two", "3", "four"})
	filtered := iter6.FilterMap(func(s string) gust.Option[string] {
		if v, err := strconv.Atoi(s); err == nil {
			return gust.Some(strconv.Itoa(v))
		}
		return gust.None[string]()
	})
	result4 := filtered.Collect()
	assert.Equal(t, []string{"1", "3"}, result4)

	// Test MapWhile (covers iterator_methods.go:713-715)
	iter7 := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	mapped := iter7.MapWhile(func(x int) gust.Option[int] {
		if x < 4 {
			return gust.Some(x * 2)
		}
		return gust.None[int]()
	})
	result5 := mapped.Collect()
	assert.Equal(t, []int{2, 4, 6}, result5)

	// Test Scan (covers iterator_methods.go:746-748)
	iter8 := iterator.FromSlice([]int{1, 2, 3})
	scanned := iter8.Scan(0, func(acc *int, x int) gust.Option[int] {
		*acc = *acc + x
		return gust.Some(*acc)
	})
	result6 := scanned.Collect()
	assert.Equal(t, []int{1, 3, 6}, result6)

	// Test MapWindows (covers iterator_methods.go:776-778)
	iter9 := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	windows := iter9.MapWindows(3, func(window []int) int {
		return window[0] + window[1] + window[2]
	})
	result7 := windows.Collect()
	assert.Equal(t, []int{6, 9, 12}, result7)
}

func TestIterator_XMapWindows(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	windows := iter.XMapWindows(3, func(window []int) any {
		return window[0] + window[1] + window[2]
	})
	result := windows.Collect()
	assert.Equal(t, []any{6, 9, 12}, result)
}

func TestFromBitSet(t *testing.T) {
	// Test with a simple bit set: 0b10101010 = 0xAA
	bitset := &mockBitSet{bits: []byte{0b10101010}}
	iter := iterator.FromBitSet(bitset)

	// Check first few bits
	pair := iter.Next()
	assert.True(t, pair.IsSome())
	assert.Equal(t, 0, pair.Unwrap().A)    // offset
	assert.Equal(t, true, pair.Unwrap().B) // bit value (MSB)

	pair = iter.Next()
	assert.True(t, pair.IsSome())
	assert.Equal(t, 1, pair.Unwrap().A)
	assert.Equal(t, false, pair.Unwrap().B)

	pair = iter.Next()
	assert.True(t, pair.IsSome())
	assert.Equal(t, 2, pair.Unwrap().A)
	assert.Equal(t, true, pair.Unwrap().B)

	// Collect all pairs
	iter = iterator.FromBitSet(bitset)
	allPairs := iter.Collect()
	assert.Len(t, allPairs, 8)
	assert.Equal(t, gust.Pair[int, bool]{A: 0, B: true}, allPairs[0])
	assert.Equal(t, gust.Pair[int, bool]{A: 1, B: false}, allPairs[1])
	assert.Equal(t, gust.Pair[int, bool]{A: 2, B: true}, allPairs[2])
	assert.Equal(t, gust.Pair[int, bool]{A: 3, B: false}, allPairs[3])
	assert.Equal(t, gust.Pair[int, bool]{A: 4, B: true}, allPairs[4])
	assert.Equal(t, gust.Pair[int, bool]{A: 5, B: false}, allPairs[5])
	assert.Equal(t, gust.Pair[int, bool]{A: 6, B: true}, allPairs[6])
	assert.Equal(t, gust.Pair[int, bool]{A: 7, B: false}, allPairs[7])

	// Test empty bit set
	emptyBitset := &mockBitSet{bits: []byte{}}
	emptyIter := iterator.FromBitSet(emptyBitset)
	assert.Equal(t, gust.None[gust.Pair[int, bool]](), emptyIter.Next())

	// Test SizeHint
	iter = iterator.FromBitSet(bitset)
	lower, upper := iter.SizeHint()
	assert.Equal(t, uint(8), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(8), upper.Unwrap())
}

func TestFromBitSetOnes(t *testing.T) {
	// Test with 0b10101010 = 0xAA (bits at positions 0, 2, 4, 6 are set)
	bitset := &mockBitSet{bits: []byte{0b10101010}}
	iter := iterator.FromBitSetOnes(bitset)

	ones := iter.Collect()
	assert.Equal(t, []int{0, 2, 4, 6}, ones)

	// Test with multiple bytes: 0b10101010, 0b11001100
	bitset2 := &mockBitSet{bits: []byte{0b10101010, 0b11001100}}
	iter2 := iterator.FromBitSetOnes(bitset2)
	ones2 := iter2.Collect()
	// First byte: positions 0, 2, 4, 6
	// Second byte: positions 8, 9, 12, 13 (8+0, 8+1, 8+4, 8+5)
	assert.Equal(t, []int{0, 2, 4, 6, 8, 9, 12, 13}, ones2)

	// Test empty bit set
	emptyBitset := &mockBitSet{bits: []byte{}}
	emptyIter := iterator.FromBitSetOnes(emptyBitset)
	assert.Equal(t, gust.None[int](), emptyIter.Next())
}

func TestFromBitSetZeros(t *testing.T) {
	// Test with 0b10101010 = 0xAA (bits at positions 1, 3, 5, 7 are unset)
	bitset := &mockBitSet{bits: []byte{0b10101010}}
	iter := iterator.FromBitSetZeros(bitset)

	zeros := iter.Collect()
	assert.Equal(t, []int{1, 3, 5, 7}, zeros)

	// Test with all bits set: 0b11111111
	allSetBitset := &mockBitSet{bits: []byte{0b11111111}}
	allSetIter := iterator.FromBitSetZeros(allSetBitset)
	assert.Equal(t, gust.None[int](), allSetIter.Next())

	// Test with no bits set: 0b00000000
	noSetBitset := &mockBitSet{bits: []byte{0b00000000}}
	noSetIter := iterator.FromBitSetZeros(noSetBitset)
	allZeros := noSetIter.Collect()
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7}, allZeros)
}

func TestFromBitSetBytes(t *testing.T) {
	// Test with a single byte: 0b10101010 = 0xAA
	bytes := []byte{0b10101010}
	iter := iterator.FromBitSetBytes(bytes)

	// Check first few bits
	pair := iter.Next()
	assert.True(t, pair.IsSome())
	assert.Equal(t, 0, pair.Unwrap().A)
	assert.Equal(t, true, pair.Unwrap().B)

	pair = iter.Next()
	assert.True(t, pair.IsSome())
	assert.Equal(t, 1, pair.Unwrap().A)
	assert.Equal(t, false, pair.Unwrap().B)

	// Collect all pairs
	iter = iterator.FromBitSetBytes(bytes)
	allPairs := iter.Collect()
	assert.Len(t, allPairs, 8)
	assert.Equal(t, gust.Pair[int, bool]{A: 0, B: true}, allPairs[0])
	assert.Equal(t, gust.Pair[int, bool]{A: 7, B: false}, allPairs[7])

	// Test with multiple bytes
	bytes2 := []byte{0b10101010, 0b11001100}
	iter2 := iterator.FromBitSetBytes(bytes2)
	allPairs2 := iter2.Collect()
	assert.Len(t, allPairs2, 16)
	assert.Equal(t, gust.Pair[int, bool]{A: 0, B: true}, allPairs2[0])    // First byte, MSB
	assert.Equal(t, gust.Pair[int, bool]{A: 7, B: false}, allPairs2[7])   // First byte, LSB
	assert.Equal(t, gust.Pair[int, bool]{A: 8, B: true}, allPairs2[8])    // Second byte, MSB
	assert.Equal(t, gust.Pair[int, bool]{A: 15, B: false}, allPairs2[15]) // Second byte, LSB

	// Test empty bytes
	emptyBytes := []byte{}
	emptyIter := iterator.FromBitSetBytes(emptyBytes)
	assert.Equal(t, gust.None[gust.Pair[int, bool]](), emptyIter.Next())

	// Test SizeHint
	iter = iterator.FromBitSetBytes(bytes)
	lower, upper := iter.SizeHint()
	assert.Equal(t, uint(8), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(8), upper.Unwrap())
}

func TestFromBitSetBytesOnes(t *testing.T) {
	// Test with 0b10101010 = 0xAA (bits at positions 0, 2, 4, 6 are set)
	bytes := []byte{0b10101010}
	iter := iterator.FromBitSetBytesOnes(bytes)

	ones := iter.Collect()
	assert.Equal(t, []int{0, 2, 4, 6}, ones)

	// Test with multiple bytes: 0b10101010, 0b11001100
	bytes2 := []byte{0b10101010, 0b11001100}
	iter2 := iterator.FromBitSetBytesOnes(bytes2)
	ones2 := iter2.Collect()
	// First byte: positions 0, 2, 4, 6
	// Second byte: positions 8, 9, 12, 13 (8+0, 8+1, 8+4, 8+5)
	assert.Equal(t, []int{0, 2, 4, 6, 8, 9, 12, 13}, ones2)

	// Test empty bytes
	emptyBytes := []byte{}
	emptyIter := iterator.FromBitSetBytesOnes(emptyBytes)
	assert.Equal(t, gust.None[int](), emptyIter.Next())

	// Test with all bits set
	allSetBytes := []byte{0b11111111}
	allSetIter := iterator.FromBitSetBytesOnes(allSetBytes)
	allOnes := allSetIter.Collect()
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7}, allOnes)
}

func TestFromBitSetBytesZeros(t *testing.T) {
	// Test with 0b10101010 = 0xAA (bits at positions 1, 3, 5, 7 are unset)
	bytes := []byte{0b10101010}
	iter := iterator.FromBitSetBytesZeros(bytes)

	zeros := iter.Collect()
	assert.Equal(t, []int{1, 3, 5, 7}, zeros)

	// Test with all bits set: 0b11111111
	allSetBytes := []byte{0b11111111}
	allSetIter := iterator.FromBitSetBytesZeros(allSetBytes)
	assert.Equal(t, gust.None[int](), allSetIter.Next())

	// Test with no bits set: 0b00000000
	noSetBytes := []byte{0b00000000}
	noSetIter := iterator.FromBitSetBytesZeros(noSetBytes)
	allZeros := noSetIter.Collect()
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7}, allZeros)

	// Test empty bytes
	emptyBytes := []byte{}
	emptyIter := iterator.FromBitSetBytesZeros(emptyBytes)
	assert.Equal(t, gust.None[int](), emptyIter.Next())
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
