package iterator_test

import (
	"math"
	"testing"

	"github.com/andeya/gust/constraints"
	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/option"
	"github.com/stretchr/testify/assert"
)

func TestCmp(t *testing.T) {
	assert.True(t, iterator.Cmp(iterator.FromSlice([]int{1}), iterator.FromSlice([]int{1})).IsEqual())
	assert.True(t, iterator.Cmp(iterator.FromSlice([]int{1}), iterator.FromSlice([]int{1, 2})).IsLess())
	assert.True(t, iterator.Cmp(iterator.FromSlice([]int{1, 2}), iterator.FromSlice([]int{1})).IsGreater())

	// Test with equal iterators (covers iter/comparison.go:17-26)
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	iter2 := iterator.FromSlice([]int{1, 2, 3})
	ordering := iterator.Cmp(iter1, iter2)
	assert.True(t, ordering.IsEqual())

	// Test with less than case
	iter3 := iterator.FromSlice([]int{1, 2})
	iter4 := iterator.FromSlice([]int{1, 2, 3})
	ordering2 := iterator.Cmp(iter3, iter4)
	assert.True(t, ordering2.IsLess())

	// Test with greater than case (need to recreate iterators as they are consumed)
	iter5 := iterator.FromSlice([]int{1, 2, 3})
	iter6 := iterator.FromSlice([]int{1, 2})
	ordering3 := iterator.Cmp(iter5, iter6)
	assert.True(t, ordering3.IsGreater())
}

func TestCmpBy(t *testing.T) {
	xs := []int{1, 2, 3, 4}
	ys := []int{1, 4, 9, 16}
	result := iterator.CmpBy(iterator.FromSlice(xs), iterator.FromSlice(ys), func(x, y int) int {
		if x*x < y {
			return -1
		}
		if x*x > y {
			return 1
		}
		return 0
	})
	assert.True(t, result.IsEqual())

	// Test itemA.IsNone() && itemB.IsNone() returns Equal
	result2 := iterator.CmpBy(iterator.Empty[int](), iterator.Empty[int](), func(x, y int) int {
		return x - y
	})
	assert.True(t, result2.IsEqual())

	// Test itemA.IsNone() returns Less
	result3 := iterator.CmpBy(iterator.Empty[int](), iterator.FromSlice([]int{1}), func(x, y int) int {
		return x - y
	})
	assert.True(t, result3.IsLess())

	// Test itemB.IsNone() returns Greater
	result4 := iterator.CmpBy(iterator.FromSlice([]int{1}), iterator.Empty[int](), func(x, y int) int {
		return x - y
	})
	assert.True(t, result4.IsGreater())

	// Test cmp result < 0 returns Less
	result5 := iterator.CmpBy(iterator.FromSlice([]int{1}), iterator.FromSlice([]int{2}), func(x, y int) int {
		return x - y
	})
	assert.True(t, result5.IsLess())

	// Test cmp result > 0 returns Greater
	result6 := iterator.CmpBy(iterator.FromSlice([]int{2}), iterator.FromSlice([]int{1}), func(x, y int) int {
		return x - y
	})
	assert.True(t, result6.IsGreater())

	// Test Cmp with x < y branch (covers comparison.go:19-21)
	result7 := iterator.Cmp(iterator.FromSlice([]int{1}), iterator.FromSlice([]int{2}))
	assert.True(t, result7.IsLess())

	// Test Cmp with x > y branch (covers comparison.go:22-24)
	result8 := iterator.Cmp(iterator.FromSlice([]int{2}), iterator.FromSlice([]int{1}))
	assert.True(t, result8.IsGreater())

	// Test PartialCmp with x < y branch (covers comparison.go:91-93)
	result9 := iterator.PartialCmp(iterator.FromSlice([]int{1}), iterator.FromSlice([]int{2}))
	assert.True(t, result9.IsSome())
	assert.True(t, result9.UnwrapUnchecked().IsLess())

	// Test PartialCmp with x > y branch (covers comparison.go:94-96)
	result10 := iterator.PartialCmp(iterator.FromSlice([]int{2}), iterator.FromSlice([]int{1}))
	assert.True(t, result10.IsSome())
	assert.True(t, result10.UnwrapUnchecked().IsGreater())

	// Test Cmp with x < y branch (covers comparison.go:19-21)
	result11 := iterator.Cmp(iterator.FromSlice([]int{1}), iterator.FromSlice([]int{2}))
	assert.True(t, result11.IsLess())

	// Test Cmp with x > y branch (covers comparison.go:22-24)
	result12 := iterator.Cmp(iterator.FromSlice([]int{2}), iterator.FromSlice([]int{1}))
	assert.True(t, result12.IsGreater())
}

func TestEq(t *testing.T) {
	assert.True(t, iterator.Eq(iterator.FromSlice([]int{1}), iterator.FromSlice([]int{1})))
	assert.False(t, iterator.Eq(iterator.FromSlice([]int{1}), iterator.FromSlice([]int{1, 2})))
}

func TestEqBy(t *testing.T) {
	xs := []int{1, 2, 3, 4}
	ys := []int{1, 4, 9, 16}
	assert.True(t, iterator.EqBy(iterator.FromSlice(xs), iterator.FromSlice(ys), func(x, y int) bool { return x*x == y }))

	// Test !eq returns false
	assert.False(t, iterator.EqBy(iterator.FromSlice([]int{1, 2}), iterator.FromSlice([]int{1, 3}), func(x, y int) bool { return x == y }))

	// Test itemA.IsNone() || itemB.IsNone() returns false
	assert.False(t, iterator.EqBy(iterator.Empty[int](), iterator.FromSlice([]int{1}), func(x, y int) bool { return x == y }))
	assert.False(t, iterator.EqBy(iterator.FromSlice([]int{1}), iterator.Empty[int](), func(x, y int) bool { return x == y }))

	// Test itemA.IsNone() && itemB.IsNone() returns true
	assert.True(t, iterator.EqBy(iterator.Empty[int](), iterator.Empty[int](), func(x, y int) bool { return x == y }))
}

func TestNe(t *testing.T) {
	assert.False(t, iterator.Ne(iterator.FromSlice([]int{1}), iterator.FromSlice([]int{1})))
	assert.True(t, iterator.Ne(iterator.FromSlice([]int{1}), iterator.FromSlice([]int{1, 2})))
}

func TestLt(t *testing.T) {
	assert.False(t, iterator.Lt(iterator.FromSlice([]int{1}), iterator.FromSlice([]int{1})))
	assert.True(t, iterator.Lt(iterator.FromSlice([]int{1}), iterator.FromSlice([]int{1, 2})))
	assert.False(t, iterator.Lt(iterator.FromSlice([]int{1, 2}), iterator.FromSlice([]int{1})))

	// Test Lt with NaN (should return false)
	nan := []float64{math.NaN()}
	assert.False(t, iterator.Lt(iterator.FromSlice(nan), iterator.FromSlice([]float64{1.0})))
	assert.False(t, iterator.Lt(iterator.FromSlice([]float64{1.0}), iterator.FromSlice(nan)))
}

func TestLe(t *testing.T) {
	assert.True(t, iterator.Le(iterator.FromSlice([]int{1}), iterator.FromSlice([]int{1})))
	assert.True(t, iterator.Le(iterator.FromSlice([]int{1}), iterator.FromSlice([]int{1, 2})))
	assert.False(t, iterator.Le(iterator.FromSlice([]int{1, 2}), iterator.FromSlice([]int{1})))

	// Test Le with NaN (should return false)
	nan := []float64{math.NaN()}
	assert.False(t, iterator.Le(iterator.FromSlice(nan), iterator.FromSlice([]float64{1.0})))
	assert.False(t, iterator.Le(iterator.FromSlice([]float64{1.0}), iterator.FromSlice(nan)))
}

func TestGt(t *testing.T) {
	assert.False(t, iterator.Gt(iterator.FromSlice([]int{1}), iterator.FromSlice([]int{1})))
	assert.False(t, iterator.Gt(iterator.FromSlice([]int{1}), iterator.FromSlice([]int{1, 2})))
	assert.True(t, iterator.Gt(iterator.FromSlice([]int{1, 2}), iterator.FromSlice([]int{1})))

	// Test Gt with NaN (should return false)
	nan := []float64{math.NaN()}
	assert.False(t, iterator.Gt(iterator.FromSlice(nan), iterator.FromSlice([]float64{1.0})))
	assert.False(t, iterator.Gt(iterator.FromSlice([]float64{1.0}), iterator.FromSlice(nan)))
}

func TestGe(t *testing.T) {
	assert.True(t, iterator.Ge(iterator.FromSlice([]int{1}), iterator.FromSlice([]int{1})))
	assert.False(t, iterator.Ge(iterator.FromSlice([]int{1}), iterator.FromSlice([]int{1, 2})))
	assert.True(t, iterator.Ge(iterator.FromSlice([]int{1, 2}), iterator.FromSlice([]int{1})))

	// Test Ge with NaN (should return false)
	nan := []float64{math.NaN()}
	assert.False(t, iterator.Ge(iterator.FromSlice(nan), iterator.FromSlice([]float64{1.0})))
	assert.False(t, iterator.Ge(iterator.FromSlice([]float64{1.0}), iterator.FromSlice(nan)))
}

func TestIsSorted(t *testing.T) {
	assert.True(t, iterator.IsSorted(iterator.FromSlice([]int{1, 2, 2, 9})))
	assert.False(t, iterator.IsSorted(iterator.FromSlice([]int{1, 3, 2, 4})))
	assert.True(t, iterator.IsSorted(iterator.FromSlice([]int{0})))
	assert.True(t, iterator.IsSorted(iterator.Empty[int]()))
}

func TestIsSortedBy(t *testing.T) {
	assert.True(t, iterator.IsSortedBy(iterator.FromSlice([]int{1, 2, 2, 9}), func(a, b int) bool { return a <= b }))
	assert.False(t, iterator.IsSortedBy(iterator.FromSlice([]int{1, 2, 2, 9}), func(a, b int) bool { return a < b }))

	// Test with multiple elements to cover the loop
	assert.True(t, iterator.IsSortedBy(iterator.FromSlice([]int{1, 2, 3, 4, 5}), func(a, b int) bool { return a <= b }))
	assert.False(t, iterator.IsSortedBy(iterator.FromSlice([]int{1, 2, 3, 2, 5}), func(a, b int) bool { return a <= b }))

	// Test with single element
	assert.True(t, iterator.IsSortedBy(iterator.FromSlice([]int{1}), func(a, b int) bool { return a <= b }))

	// Test with empty slice
	assert.True(t, iterator.IsSortedBy(iterator.Empty[int](), func(a, b int) bool { return a <= b }))
}

func TestIsSortedByKey(t *testing.T) {
	assert.True(t, iterator.IsSortedByKey(iterator.FromSlice([]string{"c", "bb", "aaa"}), func(s string) int { return len(s) }))
	assert.False(t, iterator.IsSortedByKey(iterator.FromSlice([]int{-2, -1, 0, 3}), func(n int) int {
		if n < 0 {
			return -n
		}
		return n
	}))

	// Test with multiple elements to cover the loop
	assert.True(t, iterator.IsSortedByKey(iterator.FromSlice([]string{"a", "bb", "ccc", "dddd"}), func(s string) int { return len(s) }))
	assert.False(t, iterator.IsSortedByKey(iterator.FromSlice([]string{"a", "bb", "ccc", "dd"}), func(s string) int { return len(s) }))

	// Test with single element
	assert.True(t, iterator.IsSortedByKey(iterator.FromSlice([]string{"a"}), func(s string) int { return len(s) }))

	// Test with empty slice
	assert.True(t, iterator.IsSortedByKey(iterator.Empty[string](), func(s string) int { return len(s) }))
}

func TestPartialCmp(t *testing.T) {
	// Test equal iterators
	result := iterator.PartialCmp(iterator.FromSlice([]float64{1.0}), iterator.FromSlice([]float64{1.0}))
	assert.True(t, result.IsSome())
	assert.True(t, result.Unwrap().IsEqual())

	// Test less than
	result2 := iterator.PartialCmp(iterator.FromSlice([]float64{1.0}), iterator.FromSlice([]float64{1.0, 2.0}))
	assert.True(t, result2.IsSome())
	assert.True(t, result2.Unwrap().IsLess())

	// Test greater than
	result3 := iterator.PartialCmp(iterator.FromSlice([]float64{1.0, 2.0}), iterator.FromSlice([]float64{1.0}))
	assert.True(t, result3.IsSome())
	assert.True(t, result3.Unwrap().IsGreater())

	// Test NaN case (should return None)
	nan := []float64{math.NaN()}
	result4 := iterator.PartialCmp(iterator.FromSlice(nan), iterator.FromSlice([]float64{1.0}))
	assert.True(t, result4.IsNone())
}

func TestPartialCmpBy(t *testing.T) {
	xs := []float64{1.0, 2.0, 3.0, 4.0}
	ys := []float64{1.0, 4.0, 9.0, 16.0}
	result := iterator.PartialCmpBy(iterator.FromSlice(xs), iterator.FromSlice(ys), func(x, y float64) option.Option[constraints.Ordering] {
		if x*x < y {
			return option.Some(constraints.Less())
		}
		if x*x > y {
			return option.Some(constraints.Greater())
		}
		return option.Some(constraints.Equal())
	})
	assert.True(t, result.IsSome())
	assert.True(t, result.Unwrap().IsEqual())

	// Test with None result (NaN case)
	result2 := iterator.PartialCmpBy(iterator.FromSlice([]float64{1.0}), iterator.FromSlice([]float64{1.0}), func(x, y float64) option.Option[constraints.Ordering] {
		if x != x || y != y { // NaN check
			return option.None[constraints.Ordering]()
		}
		return option.Some(constraints.Equal())
	})
	assert.True(t, result2.IsSome())

	// Test with different lengths
	result3 := iterator.PartialCmpBy(iterator.FromSlice([]float64{1.0, 2.0}), iterator.FromSlice([]float64{1.0}), func(x, y float64) option.Option[constraints.Ordering] {
		if x < y {
			return option.Some(constraints.Less())
		}
		if x > y {
			return option.Some(constraints.Greater())
		}
		return option.Some(constraints.Equal())
	})
	assert.True(t, result3.IsSome())
	assert.True(t, result3.Unwrap().IsGreater())

	// Test with actual NaN
	nan := math.NaN()
	result4 := iterator.PartialCmpBy(iterator.FromSlice([]float64{nan}), iterator.FromSlice([]float64{1.0}), func(x, y float64) option.Option[constraints.Ordering] {
		if x != x || y != y { // NaN check
			return option.None[constraints.Ordering]()
		}
		if x < y {
			return option.Some(constraints.Less())
		}
		if x > y {
			return option.Some(constraints.Greater())
		}
		return option.Some(constraints.Equal())
	})
	assert.True(t, result4.IsNone())

	// Test with equal elements (should continue)
	result5 := iterator.PartialCmpBy(iterator.FromSlice([]float64{1.0, 2.0}), iterator.FromSlice([]float64{1.0, 2.0}), func(x, y float64) option.Option[constraints.Ordering] {
		if x < y {
			return option.Some(constraints.Less())
		}
		if x > y {
			return option.Some(constraints.Greater())
		}
		return option.Some(constraints.Equal())
	})
	assert.True(t, result5.IsSome())
	assert.True(t, result5.Unwrap().IsEqual())

	// Test PartialCmpBy with less than and greater than cases
	result6 := iterator.PartialCmpBy(iterator.FromSlice([]float64{1.0}), iterator.FromSlice([]float64{2.0}), func(x, y float64) option.Option[constraints.Ordering] {
		if x < y {
			return option.Some(constraints.Less())
		}
		if x > y {
			return option.Some(constraints.Greater())
		}
		return option.Some(constraints.Equal())
	})
	assert.True(t, result6.IsSome())
	assert.True(t, result6.Unwrap().IsLess())

	result7 := iterator.PartialCmpBy(iterator.FromSlice([]float64{2.0}), iterator.FromSlice([]float64{1.0}), func(x, y float64) option.Option[constraints.Ordering] {
		if x < y {
			return option.Some(constraints.Less())
		}
		if x > y {
			return option.Some(constraints.Greater())
		}
		return option.Some(constraints.Equal())
	})
	assert.True(t, result7.IsSome())
	assert.True(t, result7.Unwrap().IsGreater())
}
