package examples_test

import (
	"fmt"

	"github.com/andeya/gust/dict"
)

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

// ExampleFilter demonstrates filtering maps.
func ExampleFilter() {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

	// Filter to keep only values greater than 2
	filtered := dict.Filter(m, func(k string, v int) bool {
		return v > 2
	})

	fmt.Println("Filtered map:", filtered)
	// Output: Filtered map: map[c:3 d:4]
}

// ExampleMapValue demonstrates mapping map values.
func ExampleMapValue() {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// Double all values
	mapped := dict.MapValue(m, func(k string, v int) int {
		return v * 2
	})

	fmt.Println("Mapped values:", mapped)
	// Output: Mapped values: map[a:2 b:4 c:6]
}
