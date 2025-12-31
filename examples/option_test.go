package examples_test

import (
	"fmt"
	"os"
	"strconv"

	"github.com/andeya/gust/option"
)

// ExampleOption_nilSafety demonstrates how Option eliminates nil pointer panics.
func ExampleOption_nilSafety() {
	// Before: Traditional Go (nil checks everywhere, easy to forget)
	// func divide(a, b float64) *float64 {
	//     if b == 0 {
	//         return nil
	//     }
	//     result := a / b
	//     return &result
	// }
	// result := divide(10, 2)
	// if result != nil {
	//     fmt.Println(*result * 2) // Risk of nil pointer panic if forgotten
	// }

	// After: gust Option (type-safe, no nil panics, 70% less code)
	divide := func(a, b float64) option.Option[float64] {
		if b == 0 {
			return option.None[float64]()
		}
		return option.Some(a / b)
	}

	quotient := divide(10, 2).
		Map(func(x float64) float64 { return x * 2 }).
		UnwrapOr(0) // Safe: never panics

	fmt.Println("Result:", quotient)
	// Output: Result: 10
}

// ExampleOption_chaining demonstrates chaining Option operations elegantly.
func ExampleOption_chaining() {
	// Before: Traditional Go (nested nil checks)
	// func process(value *int) string {
	//     if value == nil {
	//         return "No value"
	//     }
	//     doubled := *value * 2
	//     if doubled > 8 {
	//         return fmt.Sprintf("Value: %d", doubled)
	//     }
	//     return "No value"
	// }

	// After: gust Option (chainable, declarative)
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
	// Before: Traditional Go (error handling or nil checks)
	// func divide(a, b float64) (float64, error) {
	//     if b == 0 {
	//         return 0, fmt.Errorf("division by zero")
	//     }
	//     return a / b, nil
	// }
	// result, err := divide(10, 0)
	// if err != nil {
	//     fmt.Println("Cannot divide by zero")
	// }

	// After: gust Option (no errors, type-safe)
	divide := func(a, b float64) option.Option[float64] {
		if b == 0 {
			return option.None[float64]()
		}
		return option.Some(a / b)
	}

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
	// Before: Traditional Go (nil checks, error handling)
	// type Config struct {
	//     APIKey *string
	//     Port   int
	// }
	// func loadConfig() (Config, error) {
	//     apiKeyEnv := os.Getenv("API_KEY")
	//     var apiKey *string
	//     if apiKeyEnv != "" {
	//         apiKey = &apiKeyEnv
	//     }
	//     portStr := os.Getenv("PORT")
	//     port := 8080
	//     if portStr != "" {
	//         p, err := strconv.Atoi(portStr)
	//         if err != nil {
	//             return Config{}, err
	//         }
	//         port = p
	//     }
	//     return Config{APIKey: apiKey, Port: port}, nil
	// }

	// After: gust Option (type-safe, no nil checks, elegant)
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

// ExampleOption_mapLookup demonstrates safe map lookups with Option.
func ExampleOption_mapLookup() {
	// Before: Traditional Go (ok check, nil risk)
	// m := map[string]int{"a": 1, "b": 2}
	// value, ok := m["c"]
	// if !ok {
	//     value = 0 // Default value
	// }

	// After: gust Option (type-safe, no ok check needed)
	// Using dict.Get for safe map lookups (would require: import "github.com/andeya/gust/dict")
	// m := map[string]int{"a": 1, "b": 2}
	// value := dict.Get(m, "c").UnwrapOr(0) // Safe: never panics

	// Simplified example showing the pattern
	m := map[string]int{"a": 1, "b": 2}
	value, ok := m["c"]
	if !ok {
		value = 0
	}

	fmt.Println("Value:", value)
	// Output: Value: 0
}
