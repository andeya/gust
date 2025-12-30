package examples_test

import (
	"fmt"
	"time"

	"github.com/andeya/gust/random"
)

// ExampleGenerator_RandomString demonstrates basic random string generation.
func ExampleGenerator_RandomString() {
	// Create a generator with case-insensitive encoding (base62: 0-9, a-z, A-Z)
	gen := random.NewGenerator(false)

	// Generate a random string of length 16
	str := gen.RandomString(16).Unwrap()
	fmt.Printf("Random string length: %d\n", len(str))
	// Output: Random string length: 16
}

// ExampleGenerator_RandomString_caseSensitive demonstrates case-sensitive random string generation.
func ExampleGenerator_RandomString_caseSensitive() {
	// Create a generator with case-sensitive encoding (base36: 0-9, a-z)
	gen := random.NewGenerator(true)

	// Generate a random string of length 12
	str := gen.RandomString(12).Unwrap()
	fmt.Printf("Case-sensitive random string length: %d\n", len(str))
	// Output: Case-sensitive random string length: 12
}

// ExampleGenerator_StringWithNow demonstrates generating random strings with embedded current timestamp.
func ExampleGenerator_StringWithNow() {
	gen := random.NewGenerator(false)

	// Use a future timestamp to ensure it's after epoch (2026-01-01)
	futureTime := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	str := gen.StringWithTimestamp(20, futureTime).Unwrap()
	fmt.Printf("String length: %d\n", len(str))

	// Parse the timestamp back
	timestamp := gen.ParseTimestamp(str).Unwrap()
	fmt.Println("Parsed timestamp matches:", timestamp == futureTime)
	// Output:
	// String length: 20
	// Parsed timestamp matches: true
}

// ExampleGenerator_StringWithTimestamp demonstrates generating random strings with a specific timestamp.
func ExampleGenerator_StringWithTimestamp() {
	gen := random.NewGenerator(false)

	// Generate a random string with a specific timestamp
	timestamp := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	str := gen.StringWithTimestamp(18, timestamp).Unwrap()
	fmt.Printf("String length: %d\n", len(str))

	// Parse and verify the timestamp
	parsedTs := gen.ParseTimestamp(str).Unwrap()
	fmt.Println("Parsed timestamp matches:", parsedTs == timestamp)
	// Output:
	// String length: 18
	// Parsed timestamp matches: true
}

// ExampleGenerator_ParseTimestamp demonstrates parsing timestamps from random strings.
func ExampleGenerator_ParseTimestamp() {
	gen := random.NewGenerator(false)

	// Generate a string with a specific timestamp
	testTime := time.Date(2030, 6, 15, 12, 30, 45, 0, time.UTC).Unix()
	str := gen.StringWithTimestamp(16, testTime).Unwrap()

	// Parse the timestamp
	timestampRes := gen.ParseTimestamp(str)
	if timestampRes.IsErr() {
		fmt.Println("Error:", timestampRes.Err())
		return
	}

	timestamp := timestampRes.Unwrap()
	fmt.Println("Parsed timestamp matches:", timestamp == testTime)
	fmt.Println("Time:", time.Unix(timestamp, 0).UTC().Format(time.RFC3339))
	// Output:
	// Parsed timestamp matches: true
	// Time: 2030-06-15T12:30:45Z
}

// ExampleRandomBytes demonstrates generating random bytes.
func ExampleRandomBytes() {
	// Generate 32 random bytes (e.g., for encryption keys)
	bytes := random.RandomBytes(32).Unwrap()
	fmt.Printf("Random bytes length: %d\n", len(bytes))
	// Output: Random bytes length: 32
}

// ExampleGenerator_StringWithTimestamp_errorHandling demonstrates error handling with Result types.
func ExampleGenerator_StringWithTimestamp_errorHandling() {
	gen := random.NewGenerator(false)

	// Try to generate a string with invalid length (too short)
	strRes := gen.StringWithTimestamp(5, time.Now().Unix())
	if strRes.IsErr() {
		fmt.Println("Error:", strRes.Err())
	}

	// Try to generate with timestamp before epoch
	oldStrRes := gen.StringWithTimestamp(16, 1000000000) // Before 2026-01-01
	if oldStrRes.IsErr() {
		fmt.Println("Error:", oldStrRes.Err())
	}

	// Output:
	// Error: length must be greater than 6
	// Error: timestamp must be >= timestampEpoch (1767225600)
}

// ExampleGenerator_differentBases demonstrates the timestamp range for different bases.
func ExampleGenerator_differentBases() {
	// Base62 generator (case-insensitive)
	gen62 := random.NewGenerator(false)
	fmt.Println("Base62 timestamp range:")
	fmt.Println("  From: 2026-01-01 00:00:00 UTC")
	fmt.Println("  To:   3825-12-06 03:13:03 UTC")

	// Base36 generator (case-sensitive)
	gen36 := random.NewGenerator(true)
	fmt.Println("\nBase36 timestamp range:")
	fmt.Println("  From: 2026-01-01 00:00:00 UTC")
	fmt.Println("  To:   2094-12-24 05:45:35 UTC")

	// Generate strings at epoch time
	epoch := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	str62 := gen62.StringWithTimestamp(16, epoch).Unwrap()
	str36 := gen36.StringWithTimestamp(16, epoch).Unwrap()

	fmt.Println("\nEpoch time strings:")
	fmt.Printf("  Base62 length: %d, ends with: %s\n", len(str62), str62[len(str62)-6:])
	fmt.Printf("  Base36 length: %d, ends with: %s\n", len(str36), str36[len(str36)-6:])
	// Output:
	// Base62 timestamp range:
	//   From: 2026-01-01 00:00:00 UTC
	//   To:   3825-12-06 03:13:03 UTC
	//
	// Base36 timestamp range:
	//   From: 2026-01-01 00:00:00 UTC
	//   To:   2094-12-24 05:45:35 UTC
	//
	// Epoch time strings:
	//   Base62 length: 16, ends with: 000000
	//   Base36 length: 16, ends with: 000000
}
