package examples_test

import (
	"fmt"

	"github.com/andeya/gust/vec"
)

// ExampleMapAlone demonstrates mapping slice elements.
func ExampleMapAlone() {
	// Map slice elements
	numbers := []int{1, 2, 3, 4, 5}
	doubled := vec.MapAlone(numbers, func(x int) int {
		return x * 2
	})
	fmt.Println(doubled)
	// Output: [2 4 6 8 10]
}

// ExampleMapAlone_typeConversion demonstrates converting []any to specific type.
func ExampleMapAlone_typeConversion() {
	// Convert []any to specific type
	anySlice := []any{1, 2, 3, 4, 5}
	intSlice := vec.MapAlone(anySlice, func(v any) int {
		return v.(int)
	})
	fmt.Println(intSlice)
	// Output: [1 2 3 4 5]
}

