package iterator

import (
	"github.com/andeya/gust/constraints"
	"github.com/andeya/gust/option"
	"github.com/andeya/gust/pair"
	"github.com/andeya/gust/result"
	"github.com/andeya/gust/void"
)

//go:inline
func countImpl[T any](iter Iterable[T]) uint {
	var count uint
	it := Iterator[T]{iterable: iter}
	for it.Next().IsSome() {
		count++
	}
	return count
}

//go:inline
func lastImpl[T any](iter Iterable[T]) option.Option[T] {
	var last option.Option[T] = option.None[T]()
	it := Iterator[T]{iterable: iter}
	for {
		item := it.Next()
		if item.IsNone() {
			break
		}
		last = item
	}
	return last
}

//go:inline
func advanceByImpl[T any](iter Iterable[T], n uint) result.VoidResult {
	it := Iterator[T]{iterable: iter}
	for i := uint(0); i < n; i++ {
		if it.Next().IsNone() {
			return result.TryErr[void.Void](n - i)
		}
	}
	return result.Ok[void.Void](nil)
}

//go:inline
func nthImpl[T any](iter Iterable[T], n uint) option.Option[T] {
	it := Iterator[T]{iterable: iter}
	if advanceByImpl(iter, n).IsErr() {
		return option.None[T]()
	}
	return it.Next()
}

//go:inline
func nextChunkImpl[T any](iter Iterable[T], n uint) ChunkResult[[]T] {
	if n == 0 {
		return chunkOk[[]T]([]T{})
	}
	result := make([]T, 0, n)
	it := Iterator[T]{iterable: iter}
	for i := uint(0); i < n; i++ {
		item := it.Next()
		if item.IsNone() {
			// Return error with remaining elements
			return chunkErr[[]T](result)
		}
		result = append(result, item.Unwrap())
	}
	return chunkOk[[]T](result)
}

// Collect collects all items into a slice.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var result = iterator.Collect()
//	assert.Equal(t, []int{1, 2, 3}, result)
//
//go:inline
func (it Iterator[T]) Collect() []T {
	return collectImpl(it.iterable)
}

// Count consumes the iterator, counting the number of iterations.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3, 4, 5})
//	assert.Equal(t, uint(5), iterator.Count())
//
//go:inline
func (it Iterator[T]) Count() uint {
	return countImpl(it.iterable)
}

// Last returns the last element of the iterator.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	assert.Equal(t, option.Some(3), iterator.Last())
//
//go:inline
func (it Iterator[T]) Last() option.Option[T] {
	return lastImpl(it.iterable)
}

// Reduce reduces the iterator to a single value.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var sum = iterator.Reduce(func(acc int, x int) int { return acc + x })
//	assert.Equal(t, option.Some(6), sum)
//
//go:inline
func (it Iterator[T]) Reduce(f func(T, T) T) option.Option[T] {
	return reduceImpl(it.iterable, f)
}

// TryReduce reduces the elements to a single one by repeatedly applying a reducing operation.
//
// # Examples
//
//	var numbers = []int{10, 20, 5, 23, 0}
//	var sum = iterator.TryReduce(func(x, y int) result.Result[int] {
//		if x+y > 100 {
//			return result.TryErr[int](errors.New("overflow"))
//		}
//		return result.Ok(x + y)
//	})
//	assert.True(t, sum.IsOk())
//
//go:inline
func (it Iterator[T]) TryReduce(f func(T, T) result.Result[T]) result.Result[option.Option[T]] {
	return tryReduceImpl(it, f)
}

// ForEach calls a closure on each element.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	iterator.ForEach(func(x int) { fmt.Println(x) })
//
//go:inline
func (it Iterator[T]) ForEach(f func(T)) {
	forEachImpl(it.iterable, f)
}

// XTryForEach applies a fallible function to each item in the iterator,
// stopping at the first error and returning that error.
//
// # Examples
//
//	var data = []string{"no_tea.txt", "stale_bread.json", "torrential_rain.png"}
//	var res = iterator.TryForEach(func(x string) result.Result[any] {
//		fmt.Println(x)
//		return result.Ok[any](nil)
//	})
//	assert.True(t, res.IsOk())
//
//go:inline
func (it Iterator[T]) XTryForEach(f func(T) result.Result[any]) result.Result[any] {
	return TryForEach(it, f)
}

// TryForEach applies a fallible function to each item in the iterator,
// stopping at the first error and returning that error.
//
// # Examples
//
//	var data = []string{"no_tea.txt", "stale_bread.json", "torrential_rain.png"}
//	var res = iterator.TryForEach(func(x string) result.Result[string] {
//		fmt.Println(x)
//		return result.Ok[string](x+"_processed")
//	})
//	assert.True(t, res.IsOk())
//	assert.Equal(t, "no_tea.txt_processed", res.Unwrap())
//
//go:inline
func (it Iterator[T]) TryForEach(f func(T) result.Result[T]) result.Result[T] {
	return TryForEach(it, f)
}

// Partition partitions the iterator into two slices.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3, 4, 5})
//	evens, odds := iterator.Partition(func(x int) bool { return x%2 == 0 })
//	assert.Equal(t, []int{2, 4}, evens)
//	assert.Equal(t, []int{1, 3, 5}, odds)
//
//go:inline
func (it Iterator[T]) Partition(f func(T) bool) (truePart []T, falsePart []T) {
	return partitionImpl(it.iterable, f)
}

// AdvanceBy advances the iterator by n elements.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3, 4})
//	assert.True(t, iterator.AdvanceBy(2).IsOk())
//
//go:inline
func (it Iterator[T]) AdvanceBy(n uint) result.VoidResult {
	return advanceByImpl(it.iterable, n)
}

// Nth returns the nth element of the iterator.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	assert.Equal(t, option.Some(2), iterator.Nth(1))
//
//go:inline
func (it Iterator[T]) Nth(n uint) option.Option[T] {
	return nthImpl(it.iterable, n)
}

// NextChunk advances the iterator and returns an array containing the next N values.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3, 4, 5})
//	chunk := iterator.NextChunk(2)
//	assert.True(t, chunk.IsOk())
//
//go:inline
func (it Iterator[T]) NextChunk(n uint) ChunkResult[[]T] {
	return nextChunkImpl(it.iterable, n)
}

// Unzip converts an iterator of pairs into a pair of containers.
//
// Unzip() consumes an entire iterator of pairs, producing two
// collections: one from the left elements of the pairs, and one
// from the right elements.
//
// This function is, in some sense, the opposite of Zip.
//
// # Examples
//
//	var a = []pair.Pair[int, string]{
//		{A: 1, B: "a"},
//		{A: 2, B: "b"},
//		{A: 3, B: "c"},
//	}
//	var (left, right) = Unzip(FromSlice(a))
//	assert.Equal(t, []int{1, 2, 3}, left)
//	assert.Equal(t, []string{"a", "b", "c"}, right)
func Unzip[T any, U any](iter Iterator[pair.Pair[T, U]]) ([]T, []U) {
	var left []T
	var right []U
	for {
		item := iter.Next()
		if item.IsNone() {
			break
		}
		pair := item.Unwrap()
		left = append(left, pair.A)
		right = append(right, pair.B)
	}
	return left, right
}

// Sum sums the elements of an iterator.
//
// Takes each element, adds them together, and returns the result.
//
// An empty iterator returns the *additive identity* ("zero") of the type,
// which is 0 for integers and -0.0 for floats.
//
// # Panics
//
// When calling Sum() and a primitive integer type is being returned, this
// method will panic if the computation overflows and overflow checks are
// enabled.
//
// # Examples
//
//	var a = []int{1, 2, 3}
//	var sum = Sum(FromSlice(a))
//	assert.Equal(t, 6, sum)
//
//	var b = []float64{}
//	var sumFloat = Sum(FromSlice(b))
//	assert.Equal(t, -0.0, sumFloat)
//
//go:inline
func Sum[T constraints.Digit](iter Iterator[T]) T {
	var zero T
	return Fold(iter, zero, func(acc T, x T) T { return acc + x })
}

// Product iterates over the entire iterator, multiplying all the elements
//
// An empty iterator returns the one value of the type.
//
// # Panics
//
// When calling Product() and a primitive integer type is being returned,
// method will panic if the computation overflows and overflow checks are
// enabled.
//
// # Examples
//
//	var factorial = func(n int) int {
//		return Product(FromRange(1, n+1))
//	}
//	assert.Equal(t, 1, factorial(0))
//	assert.Equal(t, 1, factorial(1))
//	assert.Equal(t, 120, factorial(5))
func Product[T constraints.Digit](iter Iterator[T]) T {
	// For numeric types, we need to handle the identity element
	// For integers, it's 1, for floats it's 1.0
	// We'll use a type switch or helper function
	first := iter.Next()
	if first.IsNone() {
		return getProductIdentity[T]()
	}
	return Fold(iter, first.Unwrap(), func(acc T, x T) T { return acc * x })
}

// getProductIdentity returns the multiplicative identity for type T
func getProductIdentity[T constraints.Digit]() T {
	var zero T
	switch any(zero).(type) {
	case int:
		return T(1)
	case int8:
		return T(int8(1))
	case int16:
		return T(int16(1))
	case int32:
		return T(int32(1))
	case int64:
		return T(int64(1))
	case uint:
		return T(uint(1))
	case uint8:
		return T(uint8(1))
	case uint16:
		return T(uint16(1))
	case uint32:
		return T(uint32(1))
	case uint64:
		return T(uint64(1))
	case float32:
		return T(float32(1.0))
	case float64:
		return T(1.0)
	default:
		// For unknown types, return zero (this shouldn't happen for constraints.Digit)
		return zero
	}
}
