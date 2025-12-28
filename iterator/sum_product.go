package iterator

import (
	"github.com/andeya/gust/constraints"
)

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
