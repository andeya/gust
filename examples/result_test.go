package examples_test

import (
	"fmt"
	"strconv"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/andeya/gust/ret"
)

// ExampleResult demonstrates elegant error handling with Result.
func ExampleResult() {
	// Parse numbers with automatic error handling
	numbers := []string{"1", "2", "three", "4", "five"}

	results := iter.FilterMap(
		iter.RetMap(iter.FromSlice(numbers), strconv.Atoi),
		gust.Result[int].Ok,
	).
		Collect()

	fmt.Println("Parsed numbers:", results)
	// Output: Parsed numbers: [1 2 4]
}

// ExampleResult_AndThen demonstrates chaining Result operations elegantly.
func ExampleResult_AndThen() {
	// Chain multiple operations that can fail
	result := gust.Ok(10).
		Map(func(x int) int { return x * 2 }).
		AndThen(func(x int) gust.Result[int] {
			if x > 15 {
				return gust.Err[int]("too large")
			}
			return gust.Ok(x + 5)
		}).
		OrElse(func(err error) gust.Result[int] {
			fmt.Println("Error handled:", err)
			return gust.Ok(0)
		})

	fmt.Println("Final value:", result.Unwrap())
	// Output: Error handled: too large
	// Final value: 0
}

// ExampleAndThen demonstrates elegant error handling patterns.
func ExampleAndThen() {
	// Handle multiple operations with automatic error propagation
	result := ret.AndThen(
		gust.Ret(strconv.Atoi("42")),
		func(n int) gust.Result[string] {
			return gust.Ok(fmt.Sprintf("Number: %d", n))
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

	multiplied := gust.Ok(42).
		Map(func(x int) int { return x * 2 })

	result := ret.AndThen(multiplied, func(x int) gust.Result[string] {
		if x > 100 {
			return gust.Err[string]("value too large")
		}
		return gust.Ok(fmt.Sprintf("Result: %d", x))
	})

	if result.IsOk() {
		fmt.Println(result.Unwrap())
	} else {
		fmt.Println("Error:", result.UnwrapErr())
	}
	// Output: Result: 84
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

	// Using gust.Result for elegant error handling
	fetchUserData := func(userID int) gust.Result[string] {
		return ret.AndThen(gust.Ret(getUser(userID)), func(user *User) gust.Result[string] {
			if user == nil || user.Email == "" {
				return gust.Err[string]("invalid user")
			}
			return ret.Map(gust.Ret(getProfile(user.Email)), func(profile *Profile) string {
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
