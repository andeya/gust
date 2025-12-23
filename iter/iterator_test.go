package iter

import (
	"errors"
	"strconv"
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

type customIter struct {
	values []int
	index  int
}

func (c *customIter) Next() gust.Option[int] {
	if c.index >= len(c.values) {
		return gust.None[int]()
	}
	val := c.values[c.index]
	c.index++
	return gust.Some(val)
}

func TestFromIterable(t *testing.T) {
	// Test with Iterator[T] - should return the same iterator
	iter1 := FromSlice([]int{1, 2, 3})
	var gustIter1 gust.Iterable[int] = iter1
	iter2 := FromIterable(gustIter1)
	assert.Equal(t, gust.Some(1), iter2.Next())
	assert.Equal(t, gust.Some(2), iter2.Next())
	assert.Equal(t, gust.Some(3), iter2.Next())
	assert.Equal(t, gust.None[int](), iter2.Next())

	// Test with gust.Iterable[T] that is not Iterator[T]
	custom := &customIter{values: []int{10, 20, 30}, index: 0}
	var gustIter2 gust.Iterable[int] = custom
	iter3 := FromIterable(gustIter2)
	assert.Equal(t, gust.Some(10), iter3.Next())
	assert.Equal(t, gust.Some(20), iter3.Next())
	assert.Equal(t, gust.Some(30), iter3.Next())
	assert.Equal(t, gust.None[int](), iter3.Next())
}

func TestIteratorIterable(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3})
	iterable := iter.Iterable()
	assert.Equal(t, gust.Some(1), iterable.Next())
}

func TestTryToDoubleEnded(t *testing.T) {
	// Test with double-ended iterator
	iter1 := FromSlice([]int{1, 2, 3})
	deOpt := iter1.TryToDoubleEnded()
	assert.True(t, deOpt.IsSome())
	deIter := deOpt.Unwrap()
	assert.Equal(t, gust.Some(3), deIter.NextBack())

	// Test with non-double-ended iterator (would need a custom iterator)
	// For now, sliceIterator supports double-ended, so this will succeed
}

func TestFromSlice(t *testing.T) {
	a := []int{1, 2, 3}
	iter := FromSlice(a)

	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.Some(2), iter.Next())
	assert.Equal(t, gust.Some(3), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestCount(t *testing.T) {
	a := []int{1, 2, 3}
	assert.Equal(t, uint(3), FromSlice(a).Count())

	b := []int{1, 2, 3, 4, 5}
	assert.Equal(t, uint(5), FromSlice(b).Count())
}

func TestLast(t *testing.T) {
	a := []int{1, 2, 3}
	assert.Equal(t, gust.Some(3), FromSlice(a).Last())

	b := []int{1, 2, 3, 4, 5}
	assert.Equal(t, gust.Some(5), FromSlice(b).Last())
}

func TestMap(t *testing.T) {
	a := []int{1, 2, 3}
	iter := Map(FromSlice(a), func(x int) int { return 2 * x })

	assert.Equal(t, gust.Some(2), iter.Next())
	assert.Equal(t, gust.Some(4), iter.Next())
	assert.Equal(t, gust.Some(6), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestFilter(t *testing.T) {
	a := []int{0, 1, 2}
	iter := FromSlice(a).Filter(func(x int) bool { return x > 0 })

	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.Some(2), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestChain(t *testing.T) {
	s1 := FromSlice([]int{1, 2, 3})
	s2 := FromSlice([]int{4, 5, 6})
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
	s1 := FromSlice([]int{1, 2, 3})
	s2 := FromSlice([]int{4, 5, 6})
	iter := Zip(s1, s2)

	assert.Equal(t, gust.Some(gust.Pair[int, int]{A: 1, B: 4}), iter.Next())
	assert.Equal(t, gust.Some(gust.Pair[int, int]{A: 2, B: 5}), iter.Next())
	assert.Equal(t, gust.Some(gust.Pair[int, int]{A: 3, B: 6}), iter.Next())
	assert.Equal(t, gust.None[gust.Pair[int, int]](), iter.Next())
}

func TestEnumerate(t *testing.T) {
	a := FromSlice([]int{10, 20, 30})
	iter := Enumerate(a)

	assert.Equal(t, gust.Some(gust.Pair[uint, int]{A: 0, B: 10}), iter.Next())
	assert.Equal(t, gust.Some(gust.Pair[uint, int]{A: 1, B: 20}), iter.Next())
	assert.Equal(t, gust.Some(gust.Pair[uint, int]{A: 2, B: 30}), iter.Next())
	assert.Equal(t, gust.None[gust.Pair[uint, int]](), iter.Next())
}

func TestSkip(t *testing.T) {
	a := []int{1, 2, 3}
	iter := FromSlice(a).Skip(2)

	assert.Equal(t, gust.Some(3), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestTake(t *testing.T) {
	a := []int{1, 2, 3}
	iter := FromSlice(a).Take(2)

	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.Some(2), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestFold(t *testing.T) {
	a := []int{1, 2, 3}
	sum := Fold(FromSlice(a), 0, func(acc int, x int) int { return acc + x })
	assert.Equal(t, 6, sum)
}

func TestReduce(t *testing.T) {
	reduced := FromRange(1, 10).Reduce(func(acc int, e int) int { return acc + e })
	assert.Equal(t, gust.Some(45), reduced)
}

func TestCollect(t *testing.T) {
	a := []int{1, 2, 3}
	doubled := Map(FromSlice(a), func(x int) int { return x * 2 }).Collect()
	assert.Equal(t, []int{2, 4, 6}, doubled)
}

func TestPartition(t *testing.T) {
	a := []int{1, 2, 3}
	even, odd := FromSlice(a).Partition(func(n int) bool { return n%2 == 0 })
	assert.Equal(t, []int{2}, even)
	assert.Equal(t, []int{1, 3}, odd)
}

func TestAll(t *testing.T) {
	a := []int{1, 2, 3}
	assert.True(t, FromSlice(a).All(func(x int) bool { return x > 0 }))
	assert.False(t, FromSlice(a).All(func(x int) bool { return x > 2 }))
}

func TestAny(t *testing.T) {
	a := []int{1, 2, 3}
	assert.True(t, FromSlice(a).Any(func(x int) bool { return x > 0 }))
	assert.False(t, FromSlice(a).Any(func(x int) bool { return x > 5 }))
}

func TestFind(t *testing.T) {
	a := []int{1, 2, 3}
	assert.Equal(t, gust.Some(2), FromSlice(a).Find(func(x int) bool { return x == 2 }))
	assert.Equal(t, gust.None[int](), FromSlice(a).Find(func(x int) bool { return x == 5 }))
}

func TestPosition(t *testing.T) {
	a := []int{1, 2, 3}
	assert.Equal(t, gust.Some(uint(1)), FromSlice(a).Position(func(x int) bool { return x == 2 }))
	assert.Equal(t, gust.None[uint](), FromSlice(a).Position(func(x int) bool { return x == 5 }))
}

func TestMax(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{}
	assert.Equal(t, gust.Some(3), Max(FromSlice(a)))
	assert.Equal(t, gust.None[int](), Max(FromSlice(b)))
}

func TestMin(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{}
	assert.Equal(t, gust.Some(1), Min(FromSlice(a)))
	assert.Equal(t, gust.None[int](), Min(FromSlice(b)))
}

func TestMaxByKey(t *testing.T) {
	a := []int{-3, 0, 1, 5, -10}
	max := MaxByKey(FromSlice(a), func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	})
	assert.Equal(t, gust.Some(-10), max)

	// Test with empty iterator
	empty := []int{}
	assert.True(t, MaxByKey(FromSlice(empty), func(x int) int { return x }).IsNone())
}

func TestMaxBy(t *testing.T) {
	a := []int{-3, 0, 1, 5, -10}
	max := FromSlice(a).MaxBy(func(x, y int) int {
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
	assert.True(t, FromSlice(empty).MaxBy(func(x, y int) int {
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
	max2 := FromSlice(equal).MaxBy(func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
	assert.Equal(t, gust.Some(2), max2)
}

func TestMinByKey(t *testing.T) {
	a := []int{-3, 0, 1, 5, -10}
	min := MinByKey(FromSlice(a), func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	})
	assert.Equal(t, gust.Some(0), min)

	// Test with empty iterator
	empty := []int{}
	assert.True(t, MinByKey(FromSlice(empty), func(x int) int { return x }).IsNone())
}

func TestMinBy(t *testing.T) {
	a := []int{-3, 0, 1, 5, -10}
	min := FromSlice(a).MinBy(func(x, y int) int {
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
	assert.True(t, FromSlice(empty).MinBy(func(x, y int) int {
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
	min2 := FromSlice(equal).MinBy(func(x, y int) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
	assert.Equal(t, gust.Some(2), min2)
}

func TestFromRange(t *testing.T) {
	iter := FromRange(0, 5)
	assert.Equal(t, gust.Some(0), iter.Next())
	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.Some(2), iter.Next())
	assert.Equal(t, gust.Some(3), iter.Next())
	assert.Equal(t, gust.Some(4), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestFromElements(t *testing.T) {
	iter := FromElements(1, 2, 3)
	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.Some(2), iter.Next())
	assert.Equal(t, gust.Some(3), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestOnce(t *testing.T) {
	iter := Once(42)
	assert.Equal(t, gust.Some(42), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestEmpty(t *testing.T) {
	iter := Empty[int]()
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestFromFunc(t *testing.T) {
	count := 0
	iter := FromFunc(func() gust.Option[int] {
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
	iter := Repeat(42)
	assert.Equal(t, gust.Some(42), iter.Next())
	assert.Equal(t, gust.Some(42), iter.Next())
	assert.Equal(t, gust.Some(42), iter.Next())
	// Should repeat forever
}

func TestSizeHint(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3})
	lower, upper := iter.SizeHint()
	assert.Equal(t, uint(3), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(3), upper.Unwrap())
}

func TestSkipWhile(t *testing.T) {
	a := []int{-1, 0, 1}
	iter := FromSlice(a).SkipWhile(func(x int) bool { return x < 0 })

	assert.Equal(t, gust.Some(0), iter.Next())
	assert.Equal(t, gust.Some(1), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestTakeWhile(t *testing.T) {
	a := []int{-1, 0, 1}
	iter := FromSlice(a).TakeWhile(func(x int) bool { return x < 0 })

	assert.Equal(t, gust.Some(-1), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestAdvanceBy(t *testing.T) {
	a := []int{1, 2, 3, 4}
	iter := FromSlice(a)

	assert.Equal(t, gust.NonErrable[uint](), iter.AdvanceBy(2))
	assert.Equal(t, gust.Some(3), iter.Next())
	assert.Equal(t, gust.NonErrable[uint](), iter.AdvanceBy(0))
	assert.Equal(t, gust.ToErrable[uint](99), iter.AdvanceBy(100))
}

func TestNth(t *testing.T) {
	a := []int{1, 2, 3}
	assert.Equal(t, gust.Some(2), FromSlice(a).Nth(1))

	b := []int{1, 2, 3}
	iter := FromSlice(b)
	assert.Equal(t, gust.Some(2), iter.Nth(1))
	assert.Equal(t, gust.Some(3), iter.Nth(0))

	c := []int{1, 2, 3}
	assert.Equal(t, gust.None[int](), FromSlice(c).Nth(10))
}

func TestFilterMap(t *testing.T) {
	a := []string{"1", "two", "NaN", "four", "5"}
	iter := FilterMap(FromSlice(a), func(s string) gust.Option[int] {
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
	result := Map(
		FromRange(1, 10).Filter(func(x int) bool { return x%2 == 0 }),
		func(x int) int { return x * 2 },
	).Take(3).Collect()
	assert.Equal(t, []int{4, 8, 12}, result)
}

func TestIteratorExtChaining(t *testing.T) {
	// Test method chaining with Iterator
	iter := FromSlice([]int{1, 2, 3, 4, 5})

	// Chain Filter and Take
	filtered := iter.Filter(func(x int) bool { return x > 2 })
	taken := filtered.Take(2)
	result := taken.Collect()

	assert.Equal(t, []int{3, 4}, result)
}

func TestIteratorExtMap(t *testing.T) {
	// Test Map with Iterator (using function-style API)
	iter := FromSlice([]int{1, 2, 3})
	doubled := Map(iter, func(x int) int { return x * 2 })
	result := doubled.Collect()

	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestIteratorExtComplexChain(t *testing.T) {
	// Test complex chaining
	iter := FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	// Filter -> Map -> Take -> Collect
	mapped := Map(iter, func(x int) int { return x * 2 })
	filtered := mapped.Filter(func(x int) bool { return x > 5 })
	taken := filtered.Take(3)
	result := taken.Collect()

	assert.Equal(t, []int{6, 8, 10}, result)
}

func TestIteratorExtEnumerate(t *testing.T) {
	iter := FromSlice([]int{10, 20, 30})
	enumerated := Enumerate(iter)

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
	iter := FromSlice([]int{1, 2, 3, 4, 5})

	// Test Count
	assert.Equal(t, uint(5), iter.Count())

	// Test Last
	iter2 := FromSlice([]int{1, 2, 3})
	assert.Equal(t, gust.Some(3), iter2.Last())

	// Test All
	iter3 := FromSlice([]int{2, 4, 6})
	assert.True(t, iter3.All(func(x int) bool { return x%2 == 0 }))

	// Test Any
	iter4 := FromSlice([]int{1, 2, 3})
	assert.True(t, iter4.Any(func(x int) bool { return x > 2 }))

	// Test Find
	iter5 := FromSlice([]int{1, 2, 3})
	assert.Equal(t, gust.Some(2), iter5.Find(func(x int) bool { return x > 1 }))

	// Test Max
	iter6 := FromSlice([]int{1, 3, 2})
	assert.Equal(t, gust.Some(3), Max(iter6))

	// Test Min
	iter7 := FromSlice([]int{3, 1, 2})
	assert.Equal(t, gust.Some(1), Min(iter7))

	// Test MaxBy
	iter8 := FromSlice([]int{-3, 0, 1, 5, -10})
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
	iter9 := FromSlice([]int{-3, 0, 1, 5, -10})
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
	iter10 := FromSlice([]int{-3, 0, 1, 5, -10})
	maxByKey := MaxByKey(iter10, func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	})
	assert.Equal(t, gust.Some(-10), maxByKey)

	// Test MinByKey (function version)
	iter11 := FromSlice([]int{-3, 0, 1, 5, -10})
	minByKey := MinByKey(iter11, func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	})
	assert.Equal(t, gust.Some(0), minByKey)

	// Test TryForEach
	iter12 := FromSlice([]int{1, 2, 3})
	result := iter12.TryForEach(func(x int) gust.Result[int] {
		return gust.Ok(x)
	})
	assert.True(t, result.IsOk())

	// Test TryReduce
	iter13 := FromSlice([]int{10, 20, 5})
	sumResult := iter13.TryReduce(func(x, y int) gust.Result[int] {
		if x+y > 100 {
			return gust.Err[int](errors.New("overflow"))
		}
		return gust.Ok(x + y)
	})
	assert.True(t, sumResult.IsOk())
	assert.True(t, sumResult.Unwrap().IsSome())
	assert.Equal(t, 35, sumResult.Unwrap().Unwrap())

	// Test TryFind
	iter14 := FromSlice([]string{"1", "2", "lol", "NaN", "5"})
	findResult := iter14.TryFind(func(s string) gust.Result[bool] {
		if s == "lol" {
			return gust.Err[bool](errors.New("invalid"))
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
	iter := FromSlice([]int{1, 2, 3, 4, 5})

	// Skip first 2, then take 2
	skipped := iter.Skip(2)
	taken := skipped.Take(2)
	result := taken.Collect()

	assert.Equal(t, []int{3, 4}, result)
}

func TestIteratorExtZip(t *testing.T) {
	iter1 := FromSlice([]int{1, 2, 3})
	iter2 := FromSlice([]string{"a", "b", "c"})

	zipped := Zip(iter1, iter2)
	pair := zipped.Next()

	assert.True(t, pair.IsSome())
	assert.Equal(t, 1, pair.Unwrap().A)
	assert.Equal(t, "a", pair.Unwrap().B)
}
