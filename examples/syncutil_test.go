package examples_test

import (
	"fmt"

	"github.com/andeya/gust/result"
	"github.com/andeya/gust/syncutil"
)

// Example_syncUtil demonstrates concurrent utilities.
func Example_syncUtil() {
	// Thread-safe map
	var m syncutil.SyncMap[string, int]
	m.Store("key", 42)
	value := m.Load("key") // Returns Option[int]
	if value.IsSome() {
		fmt.Println("Value:", value.Unwrap())
	}

	// Lazy initialization
	callCount := 0
	expensiveComputation := func() int {
		callCount++
		return 42
	}

	lazy := syncutil.NewLazyValueWithFunc(func() result.Result[int] {
		return result.Ok(expensiveComputation())
	})

	// Get value multiple times - should only compute once
	v1 := lazy.TryGetValue()
	v2 := lazy.TryGetValue()
	fmt.Println("First call:", v1.Unwrap())
	fmt.Println("Second call:", v2.Unwrap())
	fmt.Println("Computation called:", callCount, "times")
	// Output:
	// Value: 42
	// First call: 42
	// Second call: 42
	// Computation called: 1 times
}
