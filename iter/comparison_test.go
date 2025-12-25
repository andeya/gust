package iter

import (
	"testing"

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
}

func TestLe(t *testing.T) {
	assert.True(t, Le(FromSlice([]int{1}), FromSlice([]int{1})))
	assert.True(t, Le(FromSlice([]int{1}), FromSlice([]int{1, 2})))
	assert.False(t, Le(FromSlice([]int{1, 2}), FromSlice([]int{1})))
}

func TestGt(t *testing.T) {
	assert.False(t, Gt(FromSlice([]int{1}), FromSlice([]int{1})))
	assert.False(t, Gt(FromSlice([]int{1}), FromSlice([]int{1, 2})))
	assert.True(t, Gt(FromSlice([]int{1, 2}), FromSlice([]int{1})))
}

func TestGe(t *testing.T) {
	assert.True(t, Ge(FromSlice([]int{1}), FromSlice([]int{1})))
	assert.False(t, Ge(FromSlice([]int{1}), FromSlice([]int{1, 2})))
	assert.True(t, Ge(FromSlice([]int{1, 2}), FromSlice([]int{1})))
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
}

func TestIsSortedByKey(t *testing.T) {
	assert.True(t, IsSortedByKey(FromSlice([]string{"c", "bb", "aaa"}), func(s string) int { return len(s) }))
	assert.False(t, IsSortedByKey(FromSlice([]int{-2, -1, 0, 3}), func(n int) int {
		if n < 0 {
			return -n
		}
		return n
	}))
}
