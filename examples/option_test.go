package examples_test

import (
	"fmt"
	"os"
	"strconv"

	"github.com/andeya/gust/option"
)

// ExampleOption demonstrates how Option eliminates nil pointer panics.
func ExampleOption() {
	// Safe division without nil checks
	divide := func(a, b float64) option.Option[float64] {
		if b == 0 {
			return option.None[float64]()
		}
		return option.Some(a / b)
	}

	quotient := divide(10, 2).
		Map(func(x float64) float64 { return x * 2 }).
		UnwrapOr(0)

	fmt.Println("Result:", quotient)
	// Output: Result: 10
}

// ExampleOption_Map demonstrates chaining Option operations.
func ExampleOption_Map() {
	// Chain operations on optional values
	value := option.Some(5).
		Map(func(x int) int { return x * 2 }).
		Filter(func(x int) bool { return x > 8 }).
		XMap(func(x int) any {
			return fmt.Sprintf("Value: %d", x)
		}).
		UnwrapOr("No value")

	fmt.Println(value)
	// Output: Value: 10
}

// ExampleOption_safeDivision demonstrates safe handling of division by zero.
func ExampleOption_safeDivision() {
	divide := func(a, b float64) option.Option[float64] {
		if b == 0 {
			return option.None[float64]()
		}
		return option.Some(a / b)
	}

	// Safe division - no panic
	quotient := divide(10, 0)
	if quotient.IsNone() {
		fmt.Println("Cannot divide by zero")
	} else {
		fmt.Println("Result:", quotient.Unwrap())
	}
	// Output: Cannot divide by zero
}

// ExampleOption_configManagement demonstrates using Option for configuration management.
func ExampleOption_configManagement() {
	type Config struct {
		APIKey option.Option[string]
		Port   option.Option[int]
	}

	loadConfig := func() Config {
		apiKeyEnv := os.Getenv("API_KEY")
		var apiKeyPtr *string
		if apiKeyEnv != "" {
			apiKeyPtr = &apiKeyEnv
		}
		return Config{
			APIKey: option.ElemOpt(apiKeyPtr),
			Port:   option.RetOpt(strconv.Atoi(os.Getenv("PORT"))),
		}
	}

	config := loadConfig()
	port := config.Port.UnwrapOr(8080)   // Default to 8080 if not set
	apiKey := config.APIKey.UnwrapOr("") // Default to empty string

	fmt.Printf("Port: %d, APIKey set: %v\n", port, config.APIKey.IsSome())
	_ = apiKey // Use apiKey to avoid unused variable
	// Output: Port: 8080, APIKey set: false
}
