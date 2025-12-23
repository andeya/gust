package iter

import (
	"github.com/andeya/gust"
)

// Cmp lexicographically compares the elements of this Iterator with those
// of another.
//
// # Examples
//
//	assert.Equal(t, gust.Less(), Cmp(FromSlice([]int{1}), FromSlice([]int{1})))
//	assert.Equal(t, gust.Less(), Cmp(FromSlice([]int{1}), FromSlice([]int{1, 2})))
//	assert.Equal(t, gust.Greater(), Cmp(FromSlice([]int{1, 2}), FromSlice([]int{1})))
func Cmp[T gust.Ord](a Iterator[T], b Iterator[T]) gust.Ordering {
	return CmpBy(a, b, func(x, y T) int {
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
		return 0
	})
}

// CmpBy lexicographically compares the elements of this Iterator with those
// of another with respect to the specified comparison function.
//
// # Examples
//
//	var xs = []int{1, 2, 3, 4}
//	var ys = []int{1, 4, 9, 16}
//	var result = CmpBy(FromSlice(xs), FromSlice(ys), func(x, y int) int {
//		if x*x < y {
//			return -1
//		}
//		if x*x > y {
//			return 1
//		}
//		return 0
//	})
//	assert.Equal(t, gust.Equal(), result)
func CmpBy[T any, U any](a Iterator[T], b Iterator[U], cmp func(T, U) int) gust.Ordering {
	for {
		itemA := a.Next()
		itemB := b.Next()

		if itemA.IsNone() && itemB.IsNone() {
			return gust.Equal()
		}
		if itemA.IsNone() {
			return gust.Less()
		}
		if itemB.IsNone() {
			return gust.Greater()
		}

		result := cmp(itemA.Unwrap(), itemB.Unwrap())
		if result < 0 {
			return gust.Less()
		}
		if result > 0 {
			return gust.Greater()
		}
		// Continue to next elements
	}
}

// PartialCmp lexicographically compares the PartialOrd elements of
// this Iterator with those of another. The comparison works like short-circuit
// evaluation, returning a result without comparing the remaining elements.
// As soon as an order can be determined, the evaluation stops and a result is returned.
//
// # Examples
//
//	var result = PartialCmp(FromSlice([]float64{1.0}), FromSlice([]float64{1.0}))
//	assert.Equal(t, gust.Some(gust.Equal()), result)
//
//	var result2 = PartialCmp(FromSlice([]float64{1.0}), FromSlice([]float64{1.0, 2.0}))
//	assert.Equal(t, gust.Some(gust.Less()), result2)
//
//	// For floating-point numbers, NaN does not have a total order
//	var nan = []float64{0.0 / 0.0}
//	var result3 = PartialCmp(FromSlice(nan), FromSlice([]float64{1.0}))
//	assert.Equal(t, gust.None[gust.Ordering](), result3)
func PartialCmp[T gust.Digit](a Iterator[T], b Iterator[T]) gust.Option[gust.Ordering] {
	return PartialCmpBy(a, b, func(x, y T) gust.Option[gust.Ordering] {
		if x < y {
			return gust.Some(gust.Less())
		}
		if x > y {
			return gust.Some(gust.Greater())
		}
		if x == y {
			return gust.Some(gust.Equal())
		}
		// NaN case
		return gust.None[gust.Ordering]()
	})
}

// PartialCmpBy lexicographically compares the elements of this Iterator with those
// of another with respect to the specified comparison function.
//
// # Examples
//
//	var xs = []float64{1.0, 2.0, 3.0, 4.0}
//	var ys = []float64{1.0, 4.0, 9.0, 16.0}
//	var result = PartialCmpBy(FromSlice(xs), FromSlice(ys), func(x, y float64) gust.Option[gust.Ordering] {
//		if x*x < y {
//			return gust.Some(gust.Less())
//		}
//		if x*x > y {
//			return gust.Some(gust.Greater())
//		}
//		return gust.Some(gust.Equal())
//	})
//	assert.Equal(t, gust.Some(gust.Equal()), result)
func PartialCmpBy[T any, U any](a Iterator[T], b Iterator[U], partialCmp func(T, U) gust.Option[gust.Ordering]) gust.Option[gust.Ordering] {
	for {
		itemA := a.Next()
		itemB := b.Next()

		if itemA.IsNone() && itemB.IsNone() {
			return gust.Some(gust.Equal())
		}
		if itemA.IsNone() {
			return gust.Some(gust.Less())
		}
		if itemB.IsNone() {
			return gust.Some(gust.Greater())
		}

		result := partialCmp(itemA.Unwrap(), itemB.Unwrap())
		if result.IsNone() {
			return gust.None[gust.Ordering]()
		}
		ord := result.Unwrap()
		if !ord.IsEqual() {
			return result
		}
		// Continue to next elements
	}
}

// Eq determines if the elements of this Iterator are equal to those of
// another.
//
// # Examples
//
//	assert.True(t, Eq(FromSlice([]int{1}), FromSlice([]int{1})))
//	assert.False(t, Eq(FromSlice([]int{1}), FromSlice([]int{1, 2})))
func Eq[T comparable](a Iterator[T], b Iterator[T]) bool {
	return EqBy(a, b, func(x, y T) bool { return x == y })
}

// EqBy determines if the elements of this Iterator are equal to those of
// another with respect to the specified equality function.
//
// # Examples
//
//	var xs = []int{1, 2, 3, 4}
//	var ys = []int{1, 4, 9, 16}
//	assert.True(t, EqBy(FromSlice(xs), FromSlice(ys), func(x, y int) bool { return x*x == y }))
func EqBy[T any, U any](a Iterator[T], b Iterator[U], eq func(T, U) bool) bool {
	for {
		itemA := a.Next()
		itemB := b.Next()

		if itemA.IsNone() && itemB.IsNone() {
			return true
		}
		if itemA.IsNone() || itemB.IsNone() {
			return false
		}

		if !eq(itemA.Unwrap(), itemB.Unwrap()) {
			return false
		}
	}
}

// Ne determines if the elements of this Iterator are not equal to those of
// another.
//
// # Examples
//
//	assert.False(t, Ne(FromSlice([]int{1}), FromSlice([]int{1})))
//	assert.True(t, Ne(FromSlice([]int{1}), FromSlice([]int{1, 2})))
func Ne[T comparable](a Iterator[T], b Iterator[T]) bool {
	return !Eq(a, b)
}

// Lt determines if the elements of this Iterator are lexicographically
// less than those of another.
//
// # Examples
//
//	assert.False(t, Lt(FromSlice([]int{1}), FromSlice([]int{1})))
//	assert.True(t, Lt(FromSlice([]int{1}), FromSlice([]int{1, 2})))
//	assert.False(t, Lt(FromSlice([]int{1, 2}), FromSlice([]int{1})))
func Lt[T gust.Digit](a Iterator[T], b Iterator[T]) bool {
	result := PartialCmp(a, b)
	if result.IsNone() {
		return false
	}
	return result.Unwrap().IsLess()
}

// Le determines if the elements of this Iterator are lexicographically
// less or equal to those of another.
//
// # Examples
//
//	assert.True(t, Le(FromSlice([]int{1}), FromSlice([]int{1})))
//	assert.True(t, Le(FromSlice([]int{1}), FromSlice([]int{1, 2})))
//	assert.False(t, Le(FromSlice([]int{1, 2}), FromSlice([]int{1})))
func Le[T gust.Digit](a Iterator[T], b Iterator[T]) bool {
	result := PartialCmp(a, b)
	if result.IsNone() {
		return false
	}
	ord := result.Unwrap()
	return ord.IsLess() || ord.IsEqual()
}

// Gt determines if the elements of this Iterator are lexicographically
// greater than those of another.
//
// # Examples
//
//	assert.False(t, Gt(FromSlice([]int{1}), FromSlice([]int{1})))
//	assert.False(t, Gt(FromSlice([]int{1}), FromSlice([]int{1, 2})))
//	assert.True(t, Gt(FromSlice([]int{1, 2}), FromSlice([]int{1})))
func Gt[T gust.Digit](a Iterator[T], b Iterator[T]) bool {
	result := PartialCmp(a, b)
	if result.IsNone() {
		return false
	}
	return result.Unwrap().IsGreater()
}

// Ge determines if the elements of this Iterator are lexicographically
// greater than or equal to those of another.
//
// # Examples
//
//	assert.True(t, Ge(FromSlice([]int{1}), FromSlice([]int{1})))
//	assert.False(t, Ge(FromSlice([]int{1}), FromSlice([]int{1, 2})))
//	assert.True(t, Ge(FromSlice([]int{1, 2}), FromSlice([]int{1})))
func Ge[T gust.Digit](a Iterator[T], b Iterator[T]) bool {
	result := PartialCmp(a, b)
	if result.IsNone() {
		return false
	}
	ord := result.Unwrap()
	return ord.IsGreater() || ord.IsEqual()
}

// IsSorted checks if the elements of this iterator are sorted.
//
// That is, for each element a and its following element b, a <= b must hold. If the
// iterator yields exactly zero or one element, true is returned.
//
// Note that if T is only PartialOrd, but not Ord, the above definition
// implies that this function returns false if any two consecutive items are not
// comparable.
//
// # Examples
//
//	assert.True(t, IsSorted(FromSlice([]int{1, 2, 2, 9})))
//	assert.False(t, IsSorted(FromSlice([]int{1, 3, 2, 4})))
//	assert.True(t, IsSorted(FromSlice([]int{0})))
//	assert.True(t, IsSorted(Empty[int]()))
func IsSorted[T gust.Digit](iter Iterator[T]) bool {
	return IsSortedBy(iter, func(a, b T) bool { return a <= b })
}

// IsSortedBy checks if the elements of this iterator are sorted using the given comparator function.
//
// Instead of using PartialOrd::partial_cmp, this function uses the given compare
// function to determine whether two elements are to be considered in sorted order.
//
// # Examples
//
//	assert.True(t, IsSortedBy(FromSlice([]int{1, 2, 2, 9}), func(a, b int) bool { return a <= b }))
//	assert.False(t, IsSortedBy(FromSlice([]int{1, 2, 2, 9}), func(a, b int) bool { return a < b }))
func IsSortedBy[T any](iter Iterator[T], compare func(T, T) bool) bool {
	last := iter.Next()
	if last.IsNone() {
		return true
	}

	for {
		curr := iter.Next()
		if curr.IsNone() {
			return true
		}
		if !compare(last.Unwrap(), curr.Unwrap()) {
			return false
		}
		last = curr
	}
}

// IsSortedByKey checks if the elements of this iterator are sorted using the given key extraction
// function.
//
// Instead of comparing the iterator's elements directly, this function compares the keys of
// the elements, as determined by f. Apart from that, it's equivalent to IsSorted; see
// its documentation for more information.
//
// # Examples
//
//	assert.True(t, IsSortedByKey(FromSlice([]string{"c", "bb", "aaa"}), func(s string) int { return len(s) }))
//	assert.False(t, IsSortedByKey(FromSlice([]int{-2, -1, 0, 3}), func(n int) int {
//		if n < 0 {
//			return -n
//		}
//		return n
//	}))
func IsSortedByKey[T any, K gust.Digit](iter Iterator[T], f func(T) K) bool {
	return IsSorted(Map(iter, f))
}

