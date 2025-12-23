package examples_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/dict"
	"github.com/andeya/gust/iter"
	"github.com/andeya/gust/ret"
	"github.com/stretchr/testify/assert"
)

// ExampleResult demonstrates elegant error handling with Result.
func ExampleResult() {
	// Parse numbers with automatic error handling
	numbers := []string{"1", "2", "three", "4", "five"}

	results := iter.FilterMap(
		iter.Map(iter.FromSlice(numbers), func(s string) gust.Result[int] {
			return gust.Ret(strconv.Atoi(s))
		}),
		gust.Result[int].Ok).
		Collect()

	fmt.Println("Parsed numbers:", results)
	// Output: Parsed numbers: [1 2 4]
}

// ExampleResult_AndThen demonstrates chaining Result operations elegantly.
func ExampleResult_AndThen() {
	// Chain multiple operations that can fail
	result := gust.Ok(10).
		Map(func(x int) int { return x * 2 }).
		AndThen(func(x int) gust.Result[int] {
			if x > 15 {
				return gust.Err[int]("too large")
			}
			return gust.Ok(x + 5)
		}).
		OrElse(func(err error) gust.Result[int] {
			fmt.Println("Error handled:", err)
			return gust.Ok(0)
		})

	fmt.Println("Final value:", result.Unwrap())
	// Output: Error handled: too large
	// Final value: 0
}

// ExampleOption demonstrates how Option eliminates nil pointer panics.
func ExampleOption() {
	// Safe division without nil checks
	divide := func(a, b float64) gust.Option[float64] {
		if b == 0 {
			return gust.None[float64]()
		}
		return gust.Some(a / b)
	}

	result := divide(10, 2).
		Map(func(x float64) float64 { return x * 2 }).
		UnwrapOr(0)

	fmt.Println("Result:", result)
	// Output: Result: 10
}

// ExampleOption_Map demonstrates chaining Option operations.
func ExampleOption_Map() {
	// Chain operations on optional values
	result := gust.Some(5).
		Map(func(x int) int { return x * 2 }).
		Filter(func(x int) bool { return x > 8 }).
		XMap(func(x int) any {
			return fmt.Sprintf("Value: %d", x)
		}).
		UnwrapOr("No value")

	fmt.Println(result)
	// Output: Value: 10
}

// ExampleIterator demonstrates Rust-like iterator chains in Go.
func ExampleIterator() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	sum := iter.FromSlice(numbers).
		Filter(func(x int) bool { return x%2 == 0 }).
		Map(func(x int) int { return x * x }).
		Take(3).
		Fold(0, func(acc int, x int) int {
			return acc + x
		})

	fmt.Println("Sum of first 3 even squares:", sum)
	// Output: Sum of first 3 even squares: 56
}

// ExampleEnumerate demonstrates advanced iterator operations.
func ExampleEnumerate() {
	data := []string{"hello", "world", "rust", "go", "iterator"}

	// Chain operations elegantly - use function API for Enumerate, then continue chaining
	enumerated := iter.Enumerate(
		iter.FromSlice(data).
			Filter(func(s string) bool { return len(s) > 2 }).
			XMap(func(s string) any { return len(s) }),
	)
	// Enumerate returns Iterator[gust.Pair[uint, any]], so we need to use Map with proper type
	result := iter.Map(enumerated, func(p gust.Pair[uint, any]) string {
		return fmt.Sprintf("%d: %d", p.A, p.B)
	}).
		Collect()

	fmt.Println(result)
	// Output: [0: 5 1: 5 2: 4 3: 8]
}

// ExampleDoubleEndedIterator demonstrates bidirectional iteration.
func ExampleDoubleEndedIterator() {
	numbers := []int{1, 2, 3, 4, 5}
	deIter := iter.FromSlice(numbers).MustToDoubleEnded()

	// Iterate from front
	fmt.Print("Front: ")
	for i := 0; i < 2; i++ {
		if val := deIter.Next(); val.IsSome() {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(val.Unwrap())
		}
	}

	// Iterate from back
	fmt.Print("\nBack: ")
	for i := 0; i < 2; i++ {
		if val := deIter.NextBack(); val.IsSome() {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(val.Unwrap())
		}
	}
	fmt.Println()
	// Output:
	// Front: 1 2
	// Back: 5 4
}

// ExampleFlatMap demonstrates flattening nested structures.
func ExampleFlatMap() {
	words := []string{"hello", "world"}

	chars := iter.FromSlice(words).
		XFlatMap(func(s string) iter.Iterator[any] {
			return iter.FromSlice([]rune(s)).XMap(func(r rune) any { return r })
		}).
		Collect()

	// Convert []any to []rune
	runeSlice := make([]rune, 0, len(chars))
	for _, v := range chars {
		runeSlice = append(runeSlice, v.(rune))
	}

	fmt.Println("Characters:", string(runeSlice))
	// Output: Characters: helloworld
}

// ExampleFindMap demonstrates finding and mapping in one operation.
func ExampleFindMap() {
	numbers := []string{"lol", "NaN", "2", "5"}

	result := iter.FromSlice(numbers).
		XFilterMap(func(s string) gust.Option[any] {
			if v, err := strconv.Atoi(s); err == nil {
				return gust.Some[any](v)
			}
			return gust.None[any]()
		}).
		Take(1).
		Collect()

	if len(result) > 0 {
		fmt.Println("First number:", result[0].(int))
	} else {
		fmt.Println("First number: 0")
	}
	// Output: First number: 2
}

// Example_iteratorPartition demonstrates splitting iterators.
func Example_iteratorPartition() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	evens, odds := iter.FromSlice(numbers).
		Partition(func(x int) bool {
			return x%2 == 0
		})

	fmt.Println("Evens:", evens)
	fmt.Println("Odds:", odds)
	// Output:
	// Evens: [2 4 6 8 10]
	// Odds: [1 3 5 7 9]
}

// ExampleGet demonstrates map utilities.
func ExampleGet() {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// Get with Option
	value := dict.Get(m, "b")
	fmt.Println("Value for 'b':", value.UnwrapOr(0))

	// Filter map
	filtered := dict.Filter(m, func(k string, v int) bool {
		return v > 1
	})
	fmt.Println("Filtered:", filtered)

	// Map values
	mapped := dict.MapValue(m, func(k string, v int) int {
		return v * 2
	})
	fmt.Println("Mapped:", mapped)
	// Output:
	// Value for 'b': 2
	// Filtered: map[b:2 c:3]
	// Mapped: map[a:2 b:4 c:6]
}

// ExampleAndThen demonstrates elegant error handling patterns.
func ExampleAndThen() {
	parseInt := func(s string) gust.Result[int] {
		return gust.Ret(strconv.Atoi(s))
	}

	// Handle multiple operations with automatic error propagation
	result := ret.AndThen(
		parseInt("42"),
		func(n int) gust.Result[string] {
			return gust.Ok(fmt.Sprintf("Number: %d", n))
		},
	)

	fmt.Println(result.Unwrap())
	// Output: Number: 42
}

// Example_realWorld demonstrates a real-world data processing scenario.
func Example_realWorld() {
	// Process user input: parse, validate, transform
	input := []string{"10", "20", "invalid", "30", "0", "40"}

	results := iter.FilterMap(
		iter.Map(iter.FromSlice(input), func(s string) gust.Result[int] {
			return gust.Ret(strconv.Atoi(s))
		}),
		gust.Result[int].Ok).
		Filter(func(x int) bool { return x > 0 }).
		Map(func(x int) int { return x * 2 }).
		Take(3).
		Collect()

	fmt.Println(results)
	// Output:
	// [20 40 60]
}

// TestExamples runs all examples to ensure they work correctly.
func TestExamples(t *testing.T) {
	// Test Result example
	numbers := []string{"1", "2", "three", "4"}
	results := iter.FromSlice(numbers).
		XMap(func(s string) any {
			return gust.Ret(strconv.Atoi(s))
		}).
		XFilterMap(func(r any) gust.Option[any] {
			return r.(gust.Result[int]).XOk()
		}).
		XMap(func(r any) any {
			return r.(int)
		}).
		Collect()

	// Convert []any to []int
	intSlice := make([]int, 0, len(results))
	for _, v := range results {
		intSlice = append(intSlice, v.(int))
	}
	assert.Equal(t, []int{1, 2, 4}, intSlice)

	// Test Option example
	divide := func(a, b float64) gust.Option[float64] {
		if b == 0 {
			return gust.None[float64]()
		}
		return gust.Some(a / b)
	}
	result := divide(10, 2).UnwrapOr(0)
	assert.Equal(t, 5.0, result)

	// Test Iterator example
	sum := iter.FromSlice([]int{1, 2, 3, 4, 5, 6}).
		Filter(func(x int) bool { return x%2 == 0 }).
		Map(func(x int) int { return x * x }).
		Take(3).
		Fold(0, func(acc int, x int) int {
			return acc + x
		})
	assert.Equal(t, 56, sum)

	// Test FlatMap
	words := []string{"ab", "cd"}
	chars := iter.FromSlice(words).
		XFlatMap(func(s string) iter.Iterator[any] {
			return iter.FromSlice([]rune(s)).XMap(func(r rune) any { return r })
		}).
		Collect()

	// Convert []any to []rune
	runeSlice := make([]rune, 0, len(chars))
	for _, v := range chars {
		runeSlice = append(runeSlice, v.(rune))
	}
	assert.Equal(t, []rune{'a', 'b', 'c', 'd'}, runeSlice)

	// Test Partition
	numbers2 := []int{1, 2, 3, 4, 5}
	evens, odds := iter.FromSlice(numbers2).
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
