package iterator_test

import (
	"errors"
	"testing"

	"github.com/andeya/gust/errutil"
	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/option"
	"github.com/andeya/gust/pair"
	"github.com/andeya/gust/result"
	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	a := []int{1, 2, 3}
	assert.Equal(t, uint(3), iterator.FromSlice(a).Count())

	b := []int{1, 2, 3, 4, 5}
	assert.Equal(t, uint(5), iterator.FromSlice(b).Count())
}

func TestLast(t *testing.T) {
	a := []int{1, 2, 3}
	assert.Equal(t, option.Some(3), iterator.FromSlice(a).Last())

	b := []int{1, 2, 3, 4, 5}
	assert.Equal(t, option.Some(5), iterator.FromSlice(b).Last())
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

func TestAdvanceBy(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3, 4})
	result := iter.AdvanceBy(2)
	assert.True(t, result.IsOk())
	assert.Equal(t, option.Some(3), iter.Next())
}

func TestNth(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	assert.Equal(t, option.Some(2), iter.Nth(1))
	assert.Equal(t, option.None[int](), iter.Nth(10))
}

func TestNextChunk(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	chunk := iter.NextChunk(2)
	assert.True(t, chunk.IsOk())
	assert.Equal(t, []int{1, 2}, chunk.Unwrap())
}

func TestForEach(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	var result []int
	iter.ForEach(func(x int) {
		result = append(result, x)
	})
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestForEachImplDirectly(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	var result []int
	iter.ForEach(func(x int) {
		result = append(result, x)
	})
	assert.Equal(t, []int{1, 2, 3}, result)
}

// TestTryForEach tests TryForEach functionality
func TestTryForEach(t *testing.T) {
	iter := iterator.FromSlice([]string{"no_tea.txt", "stale_bread.json", "torrential_rain.png"})
	res := iterator.TryForEach(iter, func(x string) result.Result[any] {
		return result.Ok[any](nil)
	})
	assert.True(t, res.IsOk())
}

// TestIterator_XTryForEach tests XTryForEach method
func TestIterator_XTryForEach(t *testing.T) {
	iter := iterator.FromSlice([]string{"no_tea.txt", "stale_bread.json", "torrential_rain.png"})
	res := iter.XTryForEach(func(x string) result.Result[any] {
		return result.Ok[any](nil)
	})
	assert.True(t, res.IsOk())
}

// TestTryReduce tests TryReduce functionality
func TestTryReduce(t *testing.T) {
	numbers := []int{10, 20, 5, 23, 0}
	iter := iterator.FromSlice(numbers)
	sum := iter.TryReduce(func(x, y int) result.Result[int] {
		if x+y > 100 {
			return result.TryErr[int](errors.New("overflow"))
		}
		return result.Ok(x + y)
	})
	assert.True(t, sum.IsOk())
}

// TestNextChunkZero tests NextChunk with zero size
func TestNextChunkZero(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	chunk := iter.NextChunk(0)
	assert.True(t, chunk.IsOk())
	assert.Equal(t, []int{}, chunk.Unwrap())
}

// TestAdvanceByZero tests AdvanceBy with zero
func TestAdvanceByZero(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	result := iter.AdvanceBy(0)
	assert.True(t, result.IsOk())
	assert.Equal(t, option.Some(1), iter.Next())
}

// TestAdvanceByTooMany tests AdvanceBy with too many steps
func TestAdvanceByTooMany(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	result := iter.AdvanceBy(100)
	assert.True(t, result.IsErr())
	assert.Equal(t, uint(97), result.ErrVal()) // 100 - 3 = 97
}

// TestNthEmpty tests Nth with empty iterator
func TestNthEmpty(t *testing.T) {
	iter := iterator.Empty[int]()
	result := iter.Nth(0)
	assert.True(t, result.IsNone())
}

// TestNthTooLarge tests Nth with index too large
func TestNthTooLarge(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	result := iter.Nth(10)
	assert.True(t, result.IsNone())
}

// TestPartitionEmpty tests Partition with empty iterator
func TestPartitionEmpty(t *testing.T) {
	iter := iterator.Empty[int]()
	truePart, falsePart := iter.Partition(func(x int) bool { return x%2 == 0 })
	// Partition should return empty slices (not nil) for consistency with design intent
	// iterator.Empty slice []int{} is semantically different from nil slice
	assert.Equal(t, []int{}, truePart)
	assert.Equal(t, []int{}, falsePart)
}

// TestPartitionAllTrue tests Partition where all elements match
func TestPartitionAllTrue(t *testing.T) {
	iter := iterator.FromSlice([]int{2, 4, 6})
	truePart, falsePart := iter.Partition(func(x int) bool { return x%2 == 0 })
	assert.Equal(t, []int{2, 4, 6}, truePart)
	// When no elements match predicate, should return empty slice (not nil)
	assert.Equal(t, []int{}, falsePart)
}

// TestPartitionAllFalse tests Partition where no elements match
func TestPartitionAllFalse(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 3, 5})
	truePart, falsePart := iter.Partition(func(x int) bool { return x%2 == 0 })
	// When no elements match predicate, should return empty slice (not nil)
	assert.Equal(t, []int{}, truePart)
	assert.Equal(t, []int{1, 3, 5}, falsePart)
}

// TestLastEmpty tests Last with empty iterator
func TestLastEmpty(t *testing.T) {
	iter := iterator.Empty[int]()
	result := iter.Last()
	assert.True(t, result.IsNone())
}

// TestCollectSizeHintEdgeCases tests Collect SizeHint edge cases
func TestCollectSizeHintEdgeCases(t *testing.T) {
	// Test with upper.IsSome() && upper.Unwrap() > lower
	iter := iterator.FromSlice([]int{1, 2, 3})
	collected := iter.Collect()
	assert.Equal(t, []int{1, 2, 3}, collected)

	// Test with upper.IsNone()
	iter2 := iterator.Repeat(1)
	collected2 := iter2.Take(3).Collect()
	assert.Equal(t, []int{1, 1, 1}, collected2)
}

// TestCollectWithUpperNone tests Collect when upper is None
func TestCollectWithUpperNone(t *testing.T) {
	iter := iterator.Repeat(1)
	collected := iter.Take(3).Collect()
	assert.Equal(t, []int{1, 1, 1}, collected)
}

// TestCollectWithUpperGreaterThanLower tests Collect when upper > lower
func TestCollectWithUpperGreaterThanLower(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	collected := iter.Collect()
	assert.Equal(t, []int{1, 2, 3}, collected)
}

// TestCollectWithUpperEqualLower tests Collect when upper == lower
func TestCollectWithUpperEqualLower(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	collected := iter.Collect()
	assert.Equal(t, []int{1, 2, 3}, collected)
}

// TestCollectWithUpperLessThanLower tests Collect when upper < lower (shouldn't happen, but test for coverage)
func TestCollectWithUpperLessThanLower(t *testing.T) {
	// This case shouldn't normally happen, but we test for completeness
	iter := iterator.FromSlice([]int{1, 2, 3})
	collected := iter.Collect()
	assert.Equal(t, []int{1, 2, 3}, collected)
}

// TestUnzip tests Unzip functionality
func TestUnzip(t *testing.T) {
	a := []pair.Pair[int, string]{
		{A: 1, B: "a"},
		{A: 2, B: "b"},
		{A: 3, B: "c"},
	}
	left, right := iterator.Unzip(iterator.FromSlice(a))
	assert.Equal(t, []int{1, 2, 3}, left)
	assert.Equal(t, []string{"a", "b", "c"}, right)
}

// TestSum tests Sum functionality
func TestSum(t *testing.T) {
	a := []int{1, 2, 3}
	sum := iterator.Sum(iterator.FromSlice(a))
	assert.Equal(t, 6, sum)

	b := []float64{}
	sumFloat := iterator.Sum(iterator.FromSlice(b))
	assert.Equal(t, 0.0, sumFloat)
}

// TestIteratorExtConsumerMethods_Consumers tests consumer methods from TestIteratorExtConsumerMethods
func TestIteratorExtConsumerMethods_Consumers(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})

	// Test Count
	assert.Equal(t, uint(5), iter.Count())

	// Test Last
	iter2 := iterator.FromSlice([]int{1, 2, 3})
	assert.Equal(t, option.Some(3), iter2.Last())

	// Test ForEach
	iter12 := iterator.FromSlice([]int{1, 2, 3})
	var forEachResult []int
	iter12.ForEach(func(x int) {
		forEachResult = append(forEachResult, x)
	})
	assert.Equal(t, []int{1, 2, 3}, forEachResult)

	// Test TryForEach
	iter13 := iterator.FromSlice([]int{1, 2, 3})
	res := iter13.TryForEach(func(x int) result.Result[int] {
		return result.Ok[int](x)
	})
	assert.True(t, res.IsOk())

	// Test TryReduce
	iter15 := iterator.FromSlice([]int{10, 20, 5})
	sumResult := iter15.TryReduce(func(x, y int) result.Result[int] {
		if x+y > 100 {
			return result.TryErr[int](errors.New("overflow"))
		}
		return result.Ok[int](x + y)
	})
	assert.True(t, sumResult.IsOk())
	assert.True(t, sumResult.Unwrap().IsSome())
	assert.Equal(t, 35, sumResult.Unwrap().Unwrap())
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
		eb, ok := p.(*errutil.ErrBox)
		assert.True(t, ok)
		assert.NotNil(t, eb)
	}()
	errChunk.Unwrap()
}

// TestProduct tests Product functionality
func TestProduct(t *testing.T) {
	factorial := func(n int) int {
		return iterator.Product(iterator.FromRange(1, n+1))
	}
	assert.Equal(t, 1, factorial(0))
	assert.Equal(t, 1, factorial(1))
	assert.Equal(t, 120, factorial(5))
}

// TestProduct_AllTypes tests Product with all numeric types to cover getProductIdentity
func TestProduct_AllTypes(t *testing.T) {
	// Test Product with all numeric types to cover getProductIdentity
	// int
	assert.Equal(t, 1, iterator.Product(iterator.FromSlice([]int{})))
	assert.Equal(t, 6, iterator.Product(iterator.FromSlice([]int{1, 2, 3})))

	// int8
	assert.Equal(t, int8(1), iterator.Product(iterator.FromSlice([]int8{})))
	assert.Equal(t, int8(6), iterator.Product(iterator.FromSlice([]int8{1, 2, 3})))

	// int16
	assert.Equal(t, int16(1), iterator.Product(iterator.FromSlice([]int16{})))
	assert.Equal(t, int16(6), iterator.Product(iterator.FromSlice([]int16{1, 2, 3})))

	// int32
	assert.Equal(t, int32(1), iterator.Product(iterator.FromSlice([]int32{})))
	assert.Equal(t, int32(6), iterator.Product(iterator.FromSlice([]int32{1, 2, 3})))

	// int64
	assert.Equal(t, int64(1), iterator.Product(iterator.FromSlice([]int64{})))
	assert.Equal(t, int64(6), iterator.Product(iterator.FromSlice([]int64{1, 2, 3})))

	// uint
	assert.Equal(t, uint(1), iterator.Product(iterator.FromSlice([]uint{})))
	assert.Equal(t, uint(6), iterator.Product(iterator.FromSlice([]uint{1, 2, 3})))

	// uint8
	assert.Equal(t, uint8(1), iterator.Product(iterator.FromSlice([]uint8{})))
	assert.Equal(t, uint8(6), iterator.Product(iterator.FromSlice([]uint8{1, 2, 3})))

	// uint16
	assert.Equal(t, uint16(1), iterator.Product(iterator.FromSlice([]uint16{})))
	assert.Equal(t, uint16(6), iterator.Product(iterator.FromSlice([]uint16{1, 2, 3})))

	// uint32
	assert.Equal(t, uint32(1), iterator.Product(iterator.FromSlice([]uint32{})))
	assert.Equal(t, uint32(6), iterator.Product(iterator.FromSlice([]uint32{1, 2, 3})))

	// uint64
	assert.Equal(t, uint64(1), iterator.Product(iterator.FromSlice([]uint64{})))
	assert.Equal(t, uint64(6), iterator.Product(iterator.FromSlice([]uint64{1, 2, 3})))

	// float32
	assert.Equal(t, float32(1.0), iterator.Product(iterator.FromSlice([]float32{})))
	assert.Equal(t, float32(6.0), iterator.Product(iterator.FromSlice([]float32{1.0, 2.0, 3.0})))

	// float64
	assert.Equal(t, 1.0, iterator.Product(iterator.FromSlice([]float64{})))
	assert.Equal(t, 6.0, iterator.Product(iterator.FromSlice([]float64{1.0, 2.0, 3.0})))
}

// TestSum_AllTypes tests Sum with all numeric types
func TestSum_AllTypes(t *testing.T) {
	// Test Sum with all numeric types
	// int
	assert.Equal(t, 0, iterator.Sum(iterator.FromSlice([]int{})))
	assert.Equal(t, 6, iterator.Sum(iterator.FromSlice([]int{1, 2, 3})))

	// float64
	assert.Equal(t, 0.0, iterator.Sum(iterator.FromSlice([]float64{})))
	assert.Equal(t, 6.0, iterator.Sum(iterator.FromSlice([]float64{1.0, 2.0, 3.0})))
}
