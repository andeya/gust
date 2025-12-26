package iter

import (
	"math"
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestCmp(t *testing.T) {
	assert.True(t, Cmp(FromSlice([]int{1}), FromSlice([]int{1})).IsEqual())
	assert.True(t, Cmp(FromSlice([]int{1}), FromSlice([]int{1, 2})).IsLess())
	assert.True(t, Cmp(FromSlice([]int{1, 2}), FromSlice([]int{1})).IsGreater())
}

func TestCmpBy(t *testing.T) {
	xs := []int{1, 2, 3, 4}
	ys := []int{1, 4, 9, 16}
	result := CmpBy(FromSlice(xs), FromSlice(ys), func(x, y int) int {
		if x*x < y {
			return -1
		}
		if x*x > y {
			return 1
		}
		return 0
	})
	assert.True(t, result.IsEqual())
}

func TestEq(t *testing.T) {
	assert.True(t, Eq(FromSlice([]int{1}), FromSlice([]int{1})))
	assert.False(t, Eq(FromSlice([]int{1}), FromSlice([]int{1, 2})))
}

func TestEqBy(t *testing.T) {
	xs := []int{1, 2, 3, 4}
	ys := []int{1, 4, 9, 16}
	assert.True(t, EqBy(FromSlice(xs), FromSlice(ys), func(x, y int) bool { return x*x == y }))
}

func TestNe(t *testing.T) {
	assert.False(t, Ne(FromSlice([]int{1}), FromSlice([]int{1})))
	assert.True(t, Ne(FromSlice([]int{1}), FromSlice([]int{1, 2})))
}

func TestLt(t *testing.T) {
	assert.False(t, Lt(FromSlice([]int{1}), FromSlice([]int{1})))
	assert.True(t, Lt(FromSlice([]int{1}), FromSlice([]int{1, 2})))
	assert.False(t, Lt(FromSlice([]int{1, 2}), FromSlice([]int{1})))

	// Test Lt with NaN (should return false)
	nan := []float64{math.NaN()}
	assert.False(t, Lt(FromSlice(nan), FromSlice([]float64{1.0})))
	assert.False(t, Lt(FromSlice([]float64{1.0}), FromSlice(nan)))
}

func TestLe(t *testing.T) {
	assert.True(t, Le(FromSlice([]int{1}), FromSlice([]int{1})))
	assert.True(t, Le(FromSlice([]int{1}), FromSlice([]int{1, 2})))
	assert.False(t, Le(FromSlice([]int{1, 2}), FromSlice([]int{1})))

	// Test Le with NaN (should return false)
	nan := []float64{math.NaN()}
	assert.False(t, Le(FromSlice(nan), FromSlice([]float64{1.0})))
	assert.False(t, Le(FromSlice([]float64{1.0}), FromSlice(nan)))
}

func TestGt(t *testing.T) {
	assert.False(t, Gt(FromSlice([]int{1}), FromSlice([]int{1})))
	assert.False(t, Gt(FromSlice([]int{1}), FromSlice([]int{1, 2})))
	assert.True(t, Gt(FromSlice([]int{1, 2}), FromSlice([]int{1})))

	// Test Gt with NaN (should return false)
	nan := []float64{math.NaN()}
	assert.False(t, Gt(FromSlice(nan), FromSlice([]float64{1.0})))
	assert.False(t, Gt(FromSlice([]float64{1.0}), FromSlice(nan)))
}

func TestGe(t *testing.T) {
	assert.True(t, Ge(FromSlice([]int{1}), FromSlice([]int{1})))
	assert.False(t, Ge(FromSlice([]int{1}), FromSlice([]int{1, 2})))
	assert.True(t, Ge(FromSlice([]int{1, 2}), FromSlice([]int{1})))

	// Test Ge with NaN (should return false)
	nan := []float64{math.NaN()}
	assert.False(t, Ge(FromSlice(nan), FromSlice([]float64{1.0})))
	assert.False(t, Ge(FromSlice([]float64{1.0}), FromSlice(nan)))
}

func TestIsSorted(t *testing.T) {
	assert.True(t, IsSorted(FromSlice([]int{1, 2, 2, 9})))
	assert.False(t, IsSorted(FromSlice([]int{1, 3, 2, 4})))
	assert.True(t, IsSorted(FromSlice([]int{0})))
	assert.True(t, IsSorted(Empty[int]()))
}

func TestIsSortedBy(t *testing.T) {
	assert.True(t, IsSortedBy(FromSlice([]int{1, 2, 2, 9}), func(a, b int) bool { return a <= b }))
	assert.False(t, IsSortedBy(FromSlice([]int{1, 2, 2, 9}), func(a, b int) bool { return a < b }))

	// Test with multiple elements to cover the loop
	assert.True(t, IsSortedBy(FromSlice([]int{1, 2, 3, 4, 5}), func(a, b int) bool { return a <= b }))
	assert.False(t, IsSortedBy(FromSlice([]int{1, 2, 3, 2, 5}), func(a, b int) bool { return a <= b }))

	// Test with single element
	assert.True(t, IsSortedBy(FromSlice([]int{1}), func(a, b int) bool { return a <= b }))

	// Test with empty slice
	assert.True(t, IsSortedBy(Empty[int](), func(a, b int) bool { return a <= b }))
}

func TestIsSortedByKey(t *testing.T) {
	assert.True(t, IsSortedByKey(FromSlice([]string{"c", "bb", "aaa"}), func(s string) int { return len(s) }))
	assert.False(t, IsSortedByKey(FromSlice([]int{-2, -1, 0, 3}), func(n int) int {
		if n < 0 {
			return -n
		}
		return n
	}))

	// Test with multiple elements to cover the loop
	assert.True(t, IsSortedByKey(FromSlice([]string{"a", "bb", "ccc", "dddd"}), func(s string) int { return len(s) }))
	assert.False(t, IsSortedByKey(FromSlice([]string{"a", "bb", "ccc", "dd"}), func(s string) int { return len(s) }))

	// Test with single element
	assert.True(t, IsSortedByKey(FromSlice([]string{"a"}), func(s string) int { return len(s) }))

	// Test with empty slice
	assert.True(t, IsSortedByKey(Empty[string](), func(s string) int { return len(s) }))
}

func TestPartialCmp(t *testing.T) {
	// Test equal iterators
	result := PartialCmp(FromSlice([]float64{1.0}), FromSlice([]float64{1.0}))
	assert.True(t, result.IsSome())
	assert.True(t, result.Unwrap().IsEqual())

	// Test less than
	result2 := PartialCmp(FromSlice([]float64{1.0}), FromSlice([]float64{1.0, 2.0}))
	assert.True(t, result2.IsSome())
	assert.True(t, result2.Unwrap().IsLess())

	// Test greater than
	result3 := PartialCmp(FromSlice([]float64{1.0, 2.0}), FromSlice([]float64{1.0}))
	assert.True(t, result3.IsSome())
	assert.True(t, result3.Unwrap().IsGreater())

	// Test NaN case (should return None)
	nan := []float64{math.NaN()}
	result4 := PartialCmp(FromSlice(nan), FromSlice([]float64{1.0}))
	assert.True(t, result4.IsNone())
}

func TestPartialCmpBy(t *testing.T) {
	xs := []float64{1.0, 2.0, 3.0, 4.0}
	ys := []float64{1.0, 4.0, 9.0, 16.0}
	result := PartialCmpBy(FromSlice(xs), FromSlice(ys), func(x, y float64) gust.Option[gust.Ordering] {
		if x*x < y {
			return gust.Some(gust.Less())
		}
		if x*x > y {
			return gust.Some(gust.Greater())
		}
		return gust.Some(gust.Equal())
	})
	assert.True(t, result.IsSome())
	assert.True(t, result.Unwrap().IsEqual())

	// Test with None result (NaN case)
	result2 := PartialCmpBy(FromSlice([]float64{1.0}), FromSlice([]float64{1.0}), func(x, y float64) gust.Option[gust.Ordering] {
		if x != x || y != y { // NaN check
			return gust.None[gust.Ordering]()
		}
		return gust.Some(gust.Equal())
	})
	assert.True(t, result2.IsSome())

	// Test with different lengths
	result3 := PartialCmpBy(FromSlice([]float64{1.0, 2.0}), FromSlice([]float64{1.0}), func(x, y float64) gust.Option[gust.Ordering] {
		if x < y {
			return gust.Some(gust.Less())
		}
		if x > y {
			return gust.Some(gust.Greater())
		}
		return gust.Some(gust.Equal())
	})
	assert.True(t, result3.IsSome())
	assert.True(t, result3.Unwrap().IsGreater())

	// Test with actual NaN
	nan := math.NaN()
	result4 := PartialCmpBy(FromSlice([]float64{nan}), FromSlice([]float64{1.0}), func(x, y float64) gust.Option[gust.Ordering] {
		if x != x || y != y { // NaN check
			return gust.None[gust.Ordering]()
		}
		if x < y {
			return gust.Some(gust.Less())
		}
		if x > y {
			return gust.Some(gust.Greater())
		}
		return gust.Some(gust.Equal())
	})
	assert.True(t, result4.IsNone())

	// Test with equal elements (should continue)
	result5 := PartialCmpBy(FromSlice([]float64{1.0, 2.0}), FromSlice([]float64{1.0, 2.0}), func(x, y float64) gust.Option[gust.Ordering] {
		if x < y {
			return gust.Some(gust.Less())
		}
		if x > y {
			return gust.Some(gust.Greater())
		}
		return gust.Some(gust.Equal())
	})
	assert.True(t, result5.IsSome())
	assert.True(t, result5.Unwrap().IsEqual())

	// Test PartialCmpBy with less than and greater than cases
	result6 := PartialCmpBy(FromSlice([]float64{1.0}), FromSlice([]float64{2.0}), func(x, y float64) gust.Option[gust.Ordering] {
		if x < y {
			return gust.Some(gust.Less())
		}
		if x > y {
			return gust.Some(gust.Greater())
		}
		return gust.Some(gust.Equal())
	})
	assert.True(t, result6.IsSome())
	assert.True(t, result6.Unwrap().IsLess())

	result7 := PartialCmpBy(FromSlice([]float64{2.0}), FromSlice([]float64{1.0}), func(x, y float64) gust.Option[gust.Ordering] {
		if x < y {
			return gust.Some(gust.Less())
		}
		if x > y {
			return gust.Some(gust.Greater())
		}
		return gust.Some(gust.Equal())
	})
	assert.True(t, result7.IsSome())
	assert.True(t, result7.Unwrap().IsGreater())
}
