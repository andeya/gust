package examples_test

import (
	"fmt"

	"github.com/andeya/gust"
)

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

// ExampleOption_safeDivision demonstrates safe handling of division by zero.
func ExampleOption_safeDivision() {
	divide := func(a, b float64) gust.Option[float64] {
		if b == 0 {
			return gust.None[float64]()
		}
		return gust.Some(a / b)
	}

	// Safe division - no panic
	result := divide(10, 0)
	if result.IsNone() {
		fmt.Println("Cannot divide by zero")
	} else {
		fmt.Println("Result:", result.Unwrap())
	}
	// Output: Cannot divide by zero
}
