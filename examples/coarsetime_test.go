package examples

import (
	"fmt"
	"time"

	"github.com/andeya/gust/coarsetime"
)

// ExampleCoarseTime_basic demonstrates basic usage of CoarseTime
func ExampleCoarseTime_basic() {
	// Use default instance (100ms precision)
	now := coarsetime.Now()
	fmt.Printf("Current time (coarse): %v\n", now)

	// Use custom precision
	ct := coarsetime.New(10 * time.Millisecond)
	defer ct.Stop()

	customNow := ct.Now()
	fmt.Printf("Current time (10ms precision): %v\n", customNow)
}

// ExampleCoarseTime_monotonic demonstrates using monotonic time for performance measurement
func ExampleCoarseTime_monotonic() {
	// Start timing using monotonic time
	start := coarsetime.Monotonic()

	// Simulate some work
	time.Sleep(50 * time.Millisecond)

	// Calculate elapsed time
	elapsed := coarsetime.Since(start)
	fmt.Printf("Work took: %v\n", elapsed)
}

// ExampleCoarseTime_presets demonstrates using preset precision instances
func ExampleCoarseTime_presets() {
	// Use fast precision (10ms)
	fastTime := coarsetime.Fast10ms.Now()
	fmt.Printf("Fast precision time: %v\n", fastTime)

	// Use standard precision (100ms, default)
	standardTime := coarsetime.Standard100ms.Now()
	fmt.Printf("Standard precision time: %v\n", standardTime)

	// Use coarse precision (1s)
	coarseTime := coarsetime.Coarse1s.Now()
	fmt.Printf("Coarse precision time: %v\n", coarseTime)
}

// ExampleCoarseTime_timeout demonstrates using coarse time for timeout checks
func ExampleCoarseTime_timeout() {
	// Set a deadline
	deadline := coarsetime.Monotonic().Add(100 * time.Millisecond)

	// Check timeout in a loop
	for i := 0; i < 10; i++ {
		if coarsetime.Monotonic().After(deadline) {
			fmt.Println("Timeout!")
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Completed before timeout")
}

// ExampleCoarseTime_logging demonstrates using coarse time for logging
func ExampleCoarseTime_logging() {
	// Use coarse time for log timestamps (reduces system calls)
	timestamp := coarsetime.Now()

	fmt.Printf("[%s] Log message\n", timestamp.Format(time.RFC3339))
}

// ExampleCoarseTime_highConcurrency demonstrates performance in high concurrency scenarios
func ExampleCoarseTime_highConcurrency() {
	// In high concurrency scenarios, coarse time significantly reduces system calls
	const iterations = 10000

	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = coarsetime.Now()
	}
	coarseDuration := time.Since(start)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = time.Now()
	}
	standardDuration := time.Since(start)

	fmt.Printf("Coarse time: %v for %d calls\n", coarseDuration, iterations)
	fmt.Printf("Standard time: %v for %d calls\n", standardDuration, iterations)
}

// ExampleCoarseTime_customInstance demonstrates creating and using a custom instance
func ExampleCoarseTime_customInstance() {
	// Create a custom instance with 50ms precision
	ct := coarsetime.New(50 * time.Millisecond)
	defer ct.Stop()

	// Use the instance
	now := ct.Now()
	monotonic := ct.Monotonic()

	fmt.Printf("Wall clock time: %v\n", now)
	fmt.Printf("Monotonic time: %v\n", monotonic)

	// Calculate time difference
	elapsed := ct.Since(monotonic)
	fmt.Printf("Elapsed: %v\n", elapsed)
}

// ExampleCoarseTime_floorCeiling demonstrates using Floor and Ceiling methods
func ExampleCoarseTime_floorCeiling() {
	ct := coarsetime.New(100 * time.Millisecond)
	defer ct.Stop()

	floor := ct.Floor()
	ceiling := ct.Ceiling()

	fmt.Printf("Floor (rounded down): %v\n", floor)
	fmt.Printf("Ceiling (rounded up): %v\n", ceiling)
	fmt.Printf("Difference: %v\n", ceiling.Sub(floor))
}
