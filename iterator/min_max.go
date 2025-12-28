package iterator

import (
	"github.com/andeya/gust/constraints"
	"github.com/andeya/gust/option"
)

// Max returns the maximum element of an iterator.
//
// If several elements are equally maximum, the last element is
// returned. If the iterator is empty, option.None[T]() is returned.
//
// Note that f32/f64 doesn't implement Ord due to NaN being
// incomparable. You can work around this by using Reduce:
//
//	var max = Reduce(FromSlice([]float32{2.4, float32(math.NaN()), 1.3}), func(a, b float32) float32 {
//		if a > b {
//			return a
//		}
//		return b
//	})
//
// # Examples
//
//	var a = []int{1, 2, 3}
//	var b = []int{}
//	assert.Equal(t, option.Some(3), Max(FromSlice(a)))
//	assert.Equal(t, option.None[int](), Max(FromSlice(b)))
//
//go:inline
func Max[T constraints.Ord](iter Iterator[T]) option.Option[T] {
	return maxByImpl(iter, func(a, b T) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})
}

// Min returns the minimum element of an iterator.
//
// If several elements are equally minimum, the first element is returned.
// If the iterator is empty, option.None[T]() is returned.
//
// Note that f32/f64 doesn't implement Ord due to NaN being
// incomparable. You can work around this by using Reduce:
//
//	var min = Reduce(FromSlice([]float32{2.4, float32(math.NaN()), 1.3}), func(a, b float32) float32 {
//		if a < b {
//			return a
//		}
//		return b
//	})
//
// # Examples
//
//	var a = []int{1, 2, 3}
//	var b = []int{}
//	assert.Equal(t, option.Some(1), Min(FromSlice(a)))
//	assert.Equal(t, option.None[int](), Min(FromSlice(b)))
//
//go:inline
func Min[T constraints.Ord](iter Iterator[T]) option.Option[T] {
	return minImpl(iter.iterable)
}

//go:inline
func minImpl[T constraints.Ord](iter Iterable[T]) option.Option[T] {
	return minByImpl(Iterator[T]{iterable: iter}, func(a, b T) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})
}

// MaxByKey returns the element that gives the maximum value from the
// specified function.
//
// If several elements are equally maximum, the last element is
// returned. If the iterator is empty, option.None[T]() is returned.
//
// # Examples
//
//	var a = []int{-3, 0, 1, 5, -10}
//	var max = MaxByKey(FromSlice(a), func(x int) int {
//		if x < 0 {
//			return -x
//		}
//		return x
//	})
//	assert.Equal(t, option.Some(-10), max)
func MaxByKey[T any, K constraints.Ord](iter Iterator[T], f func(T) K) option.Option[T] {
	return maxByImpl(iter, func(a, b T) int {
		keyA := f(a)
		keyB := f(b)
		if keyA < keyB {
			return -1
		}
		if keyA > keyB {
			return 1
		}
		return 0
	})
}

func maxByImpl[T any](iter Iterator[T], compare func(T, T) int) option.Option[T] {
	first := iter.Next()
	if first.IsNone() {
		return option.None[T]()
	}
	result := Fold(iter, first.Unwrap(), func(acc T, x T) T {
		if compare(acc, x) <= 0 {
			return x
		}
		return acc
	})
	return option.Some(result)
}

// MinByKey returns the element that gives the minimum value from the
// specified function.
//
// If several elements are equally minimum, the first element is
// returned. If the iterator is empty, option.None[T]() is returned.
//
// # Examples
//
//	var a = []int{-3, 0, 1, 5, -10}
//	var min = MinByKey(FromSlice(a), func(x int) int {
//		if x < 0 {
//			return -x
//		}
//		return x
//	})
//	assert.Equal(t, option.Some(0), min)
func MinByKey[T any, K constraints.Ord](iter Iterator[T], f func(T) K) option.Option[T] {
	return minByImpl(iter, func(a, b T) int {
		keyA := f(a)
		keyB := f(b)
		if keyA < keyB {
			return -1
		}
		if keyA > keyB {
			return 1
		}
		return 0
	})
}

func minByImpl[T any](iter Iterator[T], compare func(T, T) int) option.Option[T] {
	first := iter.Next()
	if first.IsNone() {
		return option.None[T]()
	}
	result := Fold(iter, first.Unwrap(), func(acc T, x T) T {
		if compare(acc, x) > 0 {
			return x
		}
		return acc
	})
	return option.Some(result)
}

// MaxBy returns the element that gives the maximum value with respect to the
// specified comparison function.
//
// # Examples
//
//	var iter = FromSlice([]int{-3, 0, 1, 5, -10})
//	var max = iterator.MaxBy(func(x, y int) int {
//		if x < y {
//			return -1
//		}
//		if x > y {
//			return 1
//		}
//		return 0
//	})
//	assert.Equal(t, option.Some(5), max)
//
//go:inline
func (it Iterator[T]) MaxBy(compare func(T, T) int) option.Option[T] {
	return maxByImpl(it, compare)
}

// MinBy returns the element that gives the minimum value with respect to the
// specified comparison function.
//
// # Examples
//
//	var iter = FromSlice([]int{-3, 0, 1, 5, -10})
//	var min = iterator.MinBy(func(x, y int) int {
//		if x < y {
//			return -1
//		}
//		if x > y {
//			return 1
//		}
//		return 0
//	})
//	assert.Equal(t, option.Some(-10), min)
//
//go:inline
func (it Iterator[T]) MinBy(compare func(T, T) int) option.Option[T] {
	return minByImpl(it, compare)
}
