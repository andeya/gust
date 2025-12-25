package iter

import (
	"github.com/andeya/gust"
)

// Fold folds every element into an accumulator by applying an operation,
// returning the final result.
//
// Fold() takes two arguments: an initial value, and a closure with two
// arguments: an 'accumulator', and an element. The closure returns the value that
// the accumulator should have for the next iteration.
//
// The initial value is the value the accumulator will have on the first
// call.
//
// After applying this closure to every element of the iterator, Fold()
// returns the accumulator.
//
// This operation is sometimes called 'reduce' or 'inject'.
//
// Folding is useful whenever you have a collection of something, and want
// to produce a single value from it.
//
// Note: Fold(), and similar methods that traverse the entire iterator,
// might not terminate for infinite iterators, even on traits for which a
// result is determinable in finite time.
//
// Note: Reduce() can be used to use the first element as the initial
// value, if the accumulator type and item type is the same.
//
// Note: Fold() combines elements in a *left-associative* fashion. For associative
// operators like +, the order the elements are combined in is not important, but for non-associative
// operators like - the order will affect the final result.
//
// # Examples
//
// Basic usage:
//
//	var a = []int{1, 2, 3}
//	// the sum of all of the elements of the array
//	var sum = Fold(FromSlice(a), 0, func(acc int, x int) int { return acc + x })
//	assert.Equal(t, 6, sum)
//
// Let's walk through each step of the iteration here:
//
// | element | acc | x | result |
// |---------|-----|---|--------|
// |         | 0   |   |        |
// | 1       | 0   | 1 | 1      |
// | 2       | 1   | 2 | 3      |
// | 3       | 3   | 3 | 6      |
//
// And so, our final result, 6.
//
//go:inline
func Fold[T any, B any](iter Iterator[T], init B, f func(B, T) B) B {
	return foldImpl(iter.iterable, init, f)
}

//go:inline
func foldImpl[T any, B any](iter Iterable[T], init B, f func(B, T) B) B {
	accum := init
	for {
		item := iter.Next()
		if item.IsNone() {
			break
		}
		accum = f(accum, item.Unwrap())
	}
	return accum
}

//go:inline
func reduceImpl[T any](iter Iterable[T], f func(T, T) T) gust.Option[T] {
	first := iter.Next()
	if first.IsNone() {
		return gust.None[T]()
	}
	result := Fold(Iterator[T]{iterable: iter}, first.Unwrap(), f)
	return gust.Some(result)
}

//go:inline
func forEachImpl[T any](iter Iterable[T], f func(T)) {
	for {
		item := iter.Next()
		if item.IsNone() {
			break
		}
		f(item.Unwrap())
	}
}

//go:inline
func collectImpl[T any](iter Iterable[T]) []T {
	lower, upper := iter.SizeHint()
	var capacity = lower
	if upper.IsSome() && upper.Unwrap() > lower {
		capacity = upper.Unwrap()
	}
	var result = make([]T, 0, capacity)
	for {
		item := iter.Next()
		if item.IsNone() {
			break
		}
		result = append(result, item.Unwrap())
	}
	return result
}

//go:inline
func partitionImpl[T any](iter Iterable[T], f func(T) bool) (truePart []T, falsePart []T) {
	var left []T
	var right []T
	for {
		item := iter.Next()
		if item.IsNone() {
			break
		}
		if f(item.Unwrap()) {
			left = append(left, item.Unwrap())
		} else {
			right = append(right, item.Unwrap())
		}
	}
	// Return empty slices instead of nil for consistency with design intent
	if left == nil {
		left = []T{}
	}
	if right == nil {
		right = []T{}
	}
	return left, right
}
