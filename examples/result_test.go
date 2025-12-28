package examples_test

import (
	"fmt"
	"strconv"

	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/result"
)

// ExampleResult demonstrates elegant error handling with Result.
func ExampleResult() {
	// Parse numbers with automatic error handling
	numbers := []string{"1", "2", "three", "4", "five"}

	resList := iterator.FilterMap(
		iterator.RetMap(iterator.FromSlice(numbers), strconv.Atoi),
		result.Result[int].Ok,
	).
		Collect()

	fmt.Println("Parsed numbers:", resList)
	// Output: Parsed numbers: [1 2 4]
}

// ExampleResult_AndThen demonstrates chaining Result operations elegantly.
func ExampleResult_AndThen() {
	// Chain multiple operations that can fail
	result := result.Ok(10).
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

	fmt.Println("Final value:", result.Unwrap())
	// Output: Error handled: too large
	// Final value: 0
}

// ExampleAndThen demonstrates elegant error handling patterns.
func ExampleAndThen() {
	// Handle multiple operations with automatic error propagation
	result := result.AndThen(
		result.Ret(strconv.Atoi("42")),
		func(n int) result.Result[string] {
			return result.Ok(fmt.Sprintf("Number: %d", n))
		},
	)

	fmt.Println(result.Unwrap())
	// Output: Number: 42
}

// ExampleResult_beforeAfter demonstrates the power of Result in real-world scenarios.
func ExampleResult_beforeAfter() {
	// This example shows how Result simplifies error handling
	// compared to traditional Go error handling patterns

	// Traditional approach would require multiple if err != nil checks
	// With Result, errors flow naturally through the chain

	multiplied := result.Ok(42).
		Map(func(x int) int { return x * 2 })

	result := result.AndThen(multiplied, func(x int) result.Result[string] {
		if x > 100 {
			return result.TryErr[string]("value too large")
		}
		return result.Ok(fmt.Sprintf("Result: %d", x))
	})

	if result.IsOk() {
		fmt.Println(result.Unwrap())
	} else {
		fmt.Println("Error:", result.UnwrapErr())
	}
	// Output: Result: 84
}

// Example_quickStart demonstrates the first gust program from README.
func Example_quickStart() {
	// Chain operations elegantly
	res := result.Ok(10).
		Map(func(x int) int { return x * 2 }).
		AndThen(func(x int) result.Result[int] {
			if x > 15 {
				return result.TryErr[int]("too large")
			}
			return result.Ok(x + 5)
		})

	if res.IsOk() {
		fmt.Println("Success:", res.Unwrap())
	} else {
		fmt.Println("Error:", res.UnwrapErr())
	}
	// Output: Success: 25
}

// ExampleResult_fetchUserData demonstrates a real-world error handling scenario
// similar to fetching user data from a database and API.
func ExampleResult_fetchUserData() {
	// Simulate database and API calls
	type User struct {
		ID    int
		Name  string
		Email string
	}
	type Profile struct {
		Bio string
	}

	// Simulate database call
	getUser := func(userID int) (*User, error) {
		if userID <= 0 {
			return nil, fmt.Errorf("invalid user ID")
		}
		return &User{ID: userID, Name: "Alice", Email: "alice@example.com"}, nil
	}

	// Simulate API call
	getProfile := func(email string) (*Profile, error) {
		if email == "" {
			return nil, fmt.Errorf("email required")
		}
		return &Profile{Bio: "Software developer"}, nil
	}

	// Using result.Result for elegant error handling
	fetchUserData := func(userID int) result.Result[string] {
		return result.AndThen(result.Ret(getUser(userID)), func(user *User) result.Result[string] {
			if user == nil || user.Email == "" {
				return result.TryErr[string]("invalid user")
			}
			return result.Map(result.Ret(getProfile(user.Email)), func(profile *Profile) string {
				return fmt.Sprintf("%s: %s", user.Name, profile.Bio)
			})
		})
	}

	result := fetchUserData(1)
	if result.IsOk() {
		fmt.Println(result.Unwrap())
	} else {
		fmt.Println("Error:", result.UnwrapErr())
	}
	// Output: Alice: Software developer
}
