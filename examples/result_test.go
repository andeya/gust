package examples_test

import (
	"fmt"
	"os"
	"strconv"

	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/result"
)

// ExampleResult_catchPattern demonstrates the powerful Catch pattern that eliminates all error checks.
func ExampleResult_catchPattern() {
	// Before: Traditional Go (15+ lines, 4 error checks)
	// func fetchUserData(userID int) (string, error) {
	//     user, err := db.GetUser(userID)
	//     if err != nil {
	//         return "", fmt.Errorf("db error: %w", err)
	//     }
	//     if user == nil {
	//         return "", fmt.Errorf("user not found")
	//     }
	//     if user.Email == "" {
	//         return "", fmt.Errorf("invalid user: no email")
	//     }
	//     profile, err := api.GetProfile(user.Email)
	//     if err != nil {
	//         return "", fmt.Errorf("api error: %w", err)
	//     }
	//     if profile == nil {
	//         return "", fmt.Errorf("profile not found")
	//     }
	//     return fmt.Sprintf("%s: %s", user.Name, profile.Bio), nil
	// }

	// After: gust Catch pattern (8 lines, 0 error checks)
	type User struct {
		Name  string
		Email string
	}
	type Profile struct {
		Bio string
	}

	getUser := func(userID int) (*User, error) {
		if userID <= 0 {
			return nil, fmt.Errorf("invalid user ID")
		}
		return &User{Name: "Alice", Email: "alice@example.com"}, nil
	}

	getProfile := func(email string) (*Profile, error) {
		if email == "" {
			return nil, fmt.Errorf("email required")
		}
		return &Profile{Bio: "Software developer"}, nil
	}

	fetchUserData := func(userID int) (r result.Result[string]) {
		defer r.Catch()
		user := result.Ret(getUser(userID)).Unwrap()
		if user == nil || user.Email == "" {
			return result.TryErr[string]("invalid user")
		}
		profile := result.Ret(getProfile(user.Email)).Unwrap()
		if profile == nil {
			return result.TryErr[string]("profile not found")
		}
		return result.Ok(fmt.Sprintf("%s: %s", user.Name, profile.Bio))
	}

	res := fetchUserData(1)
	if res.IsOk() {
		fmt.Println(res.Unwrap())
	} else {
		fmt.Println("Error:", res.UnwrapErr())
	}
	// Output: Alice: Software developer
}

// ExampleResult_fileIO demonstrates Catch pattern for file operations.
func ExampleResult_fileIO() {
	// Before: Traditional Go (multiple error checks)
	// func readConfigFile(filename string) (string, error) {
	//     f, err := os.Open(filename)
	//     if err != nil {
	//         return "", err
	//     }
	//     defer f.Close()
	//     data, err := io.ReadAll(f)
	//     if err != nil {
	//         return "", err
	//     }
	//     return string(data), nil
	// }

	// After: gust Catch pattern (linear flow, no error checks)
	readConfigFile := func(filename string) (r result.Result[string]) {
		defer r.Catch()
		f := result.Ret(os.Open(filename)).Unwrap()
		defer f.Close()
		data := result.Ret(os.ReadFile(filename)).Unwrap()
		return result.Ok(string(data))
	}

	// Create a temporary file for demonstration
	tmpfile, _ := os.CreateTemp("", "example")
	tmpfile.WriteString("config data")
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	res := readConfigFile(tmpfile.Name())
	if res.IsOk() {
		fmt.Println("Config:", res.Unwrap())
	} else {
		fmt.Println("Error:", res.UnwrapErr())
	}
	// Output: Config: config data
}

// ExampleResult_validationChain demonstrates chaining validations with Catch pattern.
func ExampleResult_validationChain() {
	// Before: Traditional Go (nested validations)
	// func validateAndProcess(input string) (int, error) {
	//     n, err := strconv.Atoi(input)
	//     if err != nil {
	//         return 0, err
	//     }
	//     if n < 0 {
	//         return 0, fmt.Errorf("negative not allowed")
	//     }
	//     if n > 100 {
	//         return 0, fmt.Errorf("too large")
	//     }
	//     return n * 2, nil
	// }

	// After: gust Catch pattern (linear validations)
	validateAndProcess := func(input string) (r result.Result[int]) {
		defer r.Catch()
		n := result.Ret(strconv.Atoi(input)).Unwrap()
		if n < 0 {
			return result.TryErr[int]("negative not allowed")
		}
		if n > 100 {
			return result.TryErr[int]("too large")
		}
		return result.Ok(n * 2)
	}

	res := validateAndProcess("42")
	if res.IsOk() {
		fmt.Println("Result:", res.Unwrap())
	} else {
		fmt.Println("Error:", res.UnwrapErr())
	}
	// Output: Result: 84
}

// ExampleResult_iteratorIntegration demonstrates Result with Iterator for data processing.
func ExampleResult_iteratorIntegration() {
	// Before: Traditional Go (nested loops + error handling)
	// func parseNumbers(input []string) ([]int, error) {
	//     var results []int
	//     for _, s := range input {
	//         n, err := strconv.Atoi(s)
	//         if err != nil {
	//             continue // Skip invalid
	//         }
	//         if n > 0 {
	//             results = append(results, n*2)
	//         }
	//     }
	//     if len(results) == 0 {
	//         return nil, fmt.Errorf("no valid numbers")
	//     }
	//     return results, nil
	// }

	// After: gust Iterator + Result (declarative, type-safe)
	parseNumbers := func(input []string) result.Result[[]int] {
		resList := iterator.FilterMap(
			iterator.RetMap(iterator.FromSlice(input), strconv.Atoi),
			result.Result[int].Ok,
		).Filter(func(x int) bool { return x > 0 }).
			Map(func(x int) int { return x * 2 }).
			Collect()

		if len(resList) == 0 {
			return result.TryErr[[]int]("no valid numbers")
		}
		return result.Ok(resList)
	}

	res := parseNumbers([]string{"1", "2", "three", "4", "five", "6"})
	if res.IsOk() {
		fmt.Println("Parsed:", res.Unwrap())
	} else {
		fmt.Println("Error:", res.UnwrapErr())
	}
	// Output: Parsed: [2 4 8 12]
}

// ExampleResult_chainOperations demonstrates chaining Result operations elegantly.
func ExampleResult_chainOperations() {
	// Chain multiple operations that can fail
	res := result.Ok(10).
		Map(func(x int) int { return x * 2 }).
		AndThen(func(x int) result.Result[int] {
			if x > 15 {
				return result.TryErr[int]("too large")
			}
			return result.Ok(x + 5)
		}).
		OrElse(func(err error) result.Result[int] {
			fmt.Println("Error handled:", err)
			return result.Ok(0)
		})

	fmt.Println("Final value:", res.Unwrap())
	// Output: Error handled: too large
	// Final value: 0
}

// ExampleResult_quickStart demonstrates the first gust program with Catch pattern.
func ExampleResult_quickStart() {
	// Using Catch pattern for simple operations
	processValue := func(value int) (r result.Result[int]) {
		defer r.Catch()
		doubled := value * 2
		if doubled > 20 {
			return result.TryErr[int]("too large")
		}
		return result.Ok(doubled + 5)
	}

	res := processValue(10)
	if res.IsOk() {
		fmt.Println("Success:", res.Unwrap())
	} else {
		fmt.Println("Error:", res.UnwrapErr())
	}
	// Output: Success: 25
}
