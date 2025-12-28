package examples_test

import (
	"fmt"
	"strconv"

	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/result"
)

// Example_realWorld demonstrates a real-world data processing scenario.
func Example_realWorld() {
	// Process user input: parse, validate, transform
	input := []string{"10", "20", "invalid", "30", "0", "40"}

	results := iterator.FilterMap(
		iterator.RetMap(iterator.FromSlice(input), strconv.Atoi),
		result.Result[int].Ok,
	).
		Filter(func(x int) bool { return x > 0 }).
		Map(func(x int) int { return x * 2 }).
		Take(3).
		Collect()

	fmt.Println(results)
	// Output:
	// [20 40 60]
}

// Example_dataProcessing demonstrates processing data with error handling.
func Example_dataProcessing() {
	// Parse and validate user input
	input := []string{"1", "2", "three", "4", "five", "6"}

	// Parse strings to integers, filter out errors, validate, and transform
	results := iterator.FilterMap(
		iterator.RetMap(iterator.FromSlice(input), strconv.Atoi),
		result.Result[int].Ok,
	).
		Filter(func(x int) bool { return x > 0 }).
		Map(func(x int) int { return x * x }).
		Collect()

	fmt.Println("Processed numbers:", results)
	// Output: Processed numbers: [1 4 16 36]
}

// Example_errorHandling demonstrates elegant error handling in data pipelines.
func Example_errorHandling() {
	// Simulate processing data that might fail at various stages
	processData := func(input []string) result.Result[[]int] {
		resList := iterator.FilterMap(
			iterator.RetMap(iterator.FromSlice(input), strconv.Atoi),
			result.Result[int].Ok,
		).
			Collect()

		if len(resList) == 0 {
			return result.TryErr[[]int]("no valid numbers found")
		}

		return result.Ok(resList)
	}

	result := processData([]string{"1", "2", "3"})
	if result.IsOk() {
		fmt.Println("Success:", result.Unwrap())
	} else {
		fmt.Println("Error:", result.UnwrapErr())
	}
	// Output: Success: [1 2 3]
}
