package examples_test

import (
	"strconv"
	"testing"

	"github.com/andeya/gust/dict"
	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/option"
	"github.com/andeya/gust/result"
	"github.com/andeya/gust/vec"
	"github.com/stretchr/testify/assert"
)

// TestExamples runs all examples to ensure they work correctly.
func TestExamples(t *testing.T) {
	// Test Result example
	numbers := []string{"1", "2", "three", "4"}
	results := iterator.FilterMap(
		iterator.RetMap(iterator.FromSlice(numbers), strconv.Atoi),
		result.Result[int].Ok,
	).
		Collect()

	assert.Equal(t, []int{1, 2, 4}, results)

	// Test Option example
	divide := func(a, b float64) option.Option[float64] {
		if b == 0 {
			return option.None[float64]()
		}
		return option.Some(a / b)
	}
	result := divide(10, 2).UnwrapOr(0)
	assert.Equal(t, 5.0, result)

	// Test Iterator example
	sum := iterator.FromSlice([]int{1, 2, 3, 4, 5, 6}).
		Filter(func(x int) bool { return x%2 == 0 }).
		Map(func(x int) int { return x * x }).
		Take(3).
		Fold(0, func(acc int, x int) int {
			return acc + x
		})
	assert.Equal(t, 56, sum)

	// Test FlatMap
	words := []string{"ab", "cd"}
	chars := iterator.FromSlice(words).
		XFlatMap(func(s string) iterator.Iterator[any] {
			return iterator.FromSlice([]rune(s)).XMap(func(r rune) any { return r })
		}).
		Collect()

	// Convert []any to []rune using vec.MapAlone
	runeSlice := vec.MapAlone(chars, func(v any) rune {
		return v.(rune)
	})
	assert.Equal(t, []rune{'a', 'b', 'c', 'd'}, runeSlice)

	// Test Partition
	numbers2 := []int{1, 2, 3, 4, 5}
	evens, odds := iterator.FromSlice(numbers2).
		Partition(func(x int) bool {
			return x%2 == 0
		})
	assert.Equal(t, []int{2, 4}, evens)
	assert.Equal(t, []int{1, 3, 5}, odds)

	// Test Dict
	m := map[string]int{"a": 1, "b": 2}
	value := dict.Get(m, "b")
	assert.True(t, value.IsSome())
	assert.Equal(t, 2, value.Unwrap())
}
