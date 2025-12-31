package examples_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/andeya/gust/fileutil"
	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/result"
)

// Example_realWorld_dataProcessing demonstrates Iterator + Result for data processing pipelines.
func Example_realWorld_dataProcessing() {
	// Before: Traditional Go (nested loops + error handling)
	// func processUserInput(input []string) ([]int, error) {
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

	// After: gust Iterator + Result (declarative, type-safe, 70% less code)
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

// Example_realWorld_batchFileOperations demonstrates Iterator + Catch for batch file operations.
func Example_realWorld_batchFileOperations() {
	// Before: Traditional Go (for loop + error handling)
	// func normalizePaths(paths []string) ([]string, error) {
	//     var results []string
	//     for _, p := range paths {
	//         abs, err := filepath.Abs(p)
	//         if err != nil {
	//             return nil, err
	//         }
	//         results = append(results, abs)
	//     }
	//     return results, nil
	// }

	// After: gust Iterator + Catch (linear flow, automatic error propagation)
	normalizePaths := func(paths []string) (r result.Result[[]string]) {
		defer r.Catch()
		res := iterator.FromSlice(paths).Map(func(p string) string {
			return result.Ret(filepath.Abs(p)).Unwrap()
		}).Collect()
		return result.Ok(res)
	}

	// Use current directory for demonstration
	paths := []string{".", ".."}

	res := normalizePaths(paths)
	if res.IsOk() {
		fmt.Println("Normalized paths:", len(res.Unwrap()), "items")
	} else {
		fmt.Println("Error:", res.UnwrapErr())
	}
	// Output: Normalized paths: 2 items
}

// Example_realWorld_dataTransformation demonstrates complex data transformation with Iterator.
func Example_realWorld_dataTransformation() {
	// Before: Traditional Go (nested loops, manual aggregation)
	// func analyzeNumbers(input []string) (int, error) {
	//     sum := 0
	//     count := 0
	//     for _, s := range input {
	//         n, err := strconv.Atoi(s)
	//         if err != nil {
	//             continue
	//         }
	//         if n > 0 && n < 100 {
	//             sum += n * n
	//             count++
	//         }
	//     }
	//     if count == 0 {
	//         return 0, fmt.Errorf("no valid numbers")
	//     }
	//     return sum, nil
	// }

	// After: gust Iterator (declarative pipeline, type-safe)
	input := []string{"1", "2", "three", "4", "five", "6", "150"}

	sum := iterator.FilterMap(
		iterator.RetMap(iterator.FromSlice(input), strconv.Atoi),
		result.Result[int].Ok,
	).
		Filter(func(x int) bool { return x > 0 && x < 100 }).
		Map(func(x int) int { return x * x }).
		Fold(0, func(acc, x int) int { return acc + x })

	fmt.Println("Sum of squares:", sum)
	// Output: Sum of squares: 57
}

// Example_realWorld_fileSystemOperations demonstrates file system operations with Catch pattern.
func Example_realWorld_fileSystemOperations() {
	// Before: Traditional Go (multiple error checks, nested conditions)
	// func copyDirectory(src, dst string) error {
	//     info, err := os.Stat(src)
	//     if err != nil {
	//         return err
	//     }
	//     if err = os.MkdirAll(dst, info.Mode()); err != nil {
	//         return err
	//     }
	//     entries, err := os.ReadDir(src)
	//     if err != nil {
	//         return err
	//     }
	//     for _, entry := range entries {
	//         // ... more error checks
	//     }
	//     return nil
	// }

	// After: gust Catch pattern (linear flow, single error handler)
	// This demonstrates the pattern used in fileutil.CopyDir
	var copyDirExample func(string, string) result.VoidResult
	copyDirExample = func(src, dst string) (r result.VoidResult) {
		defer r.Catch()
		info := result.Ret(os.Stat(src)).Unwrap()
		result.RetVoid(os.MkdirAll(dst, info.Mode())).Unwrap()
		entries := result.Ret(os.ReadDir(src)).Unwrap()
		for _, entry := range entries {
			srcPath := filepath.Join(src, entry.Name())
			dstPath := filepath.Join(dst, entry.Name())
			if entry.IsDir() {
				// Recursive call - errors automatically propagate via Catch
				copyDirExample(srcPath, dstPath).Unwrap()
			} else {
				// Copy file - errors automatically propagate
				fileutil.CopyFile(srcPath, dstPath).Unwrap()
			}
		}
		return result.OkVoid()
	}

	// Create temporary directories for demonstration
	tmpDir, _ := os.MkdirTemp("", "example")
	defer os.RemoveAll(tmpDir)
	srcDir := filepath.Join(tmpDir, "src")
	dstDir := filepath.Join(tmpDir, "dst")
	os.MkdirAll(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "test.txt"), []byte("test"), 0644)

	res := copyDirExample(srcDir, dstDir)
	if res.IsOk() {
		fmt.Println("Directory copied successfully")
	} else {
		fmt.Println("Error:", res.Err())
	}
	// Output: Directory copied successfully
}

// Example_realWorld_dataValidation demonstrates data validation pipeline with Iterator + Result.
func Example_realWorld_dataValidation() {
	// Before: Traditional Go (manual validation, error accumulation)
	// func validateAndProcess(input []string) ([]int, error) {
	//     var results []int
	//     var errors []error
	//     for _, s := range input {
	//         n, err := strconv.Atoi(s)
	//         if err != nil {
	//             errors = append(errors, err)
	//             continue
	//         }
	//         if n < 0 || n > 100 {
	//             continue
	//         }
	//         results = append(results, n)
	//     }
	//     if len(results) == 0 {
	//         return nil, fmt.Errorf("no valid numbers: %v", errors)
	//     }
	//     return results, nil
	// }

	// After: gust Iterator + Result (automatic error filtering, declarative)
	input := []string{"1", "2", "three", "4", "five", "6", "150", "-5"}

	results := iterator.FilterMap(
		iterator.RetMap(iterator.FromSlice(input), strconv.Atoi),
		result.Result[int].Ok,
	).
		Filter(func(x int) bool { return x > 0 && x <= 100 }).
		Collect()

	fmt.Println("Valid numbers:", results)
	// Output: Valid numbers: [1 2 4 6]
}

// Example_realWorld_apiCallChain demonstrates API call chain with Catch pattern.
func Example_realWorld_apiCallChain() {
	// Before: Traditional Go (nested error handling, 15+ lines)
	// func fetchUserProfile(userID int) (string, error) {
	//     user, err := db.GetUser(userID)
	//     if err != nil {
	//         return "", fmt.Errorf("db error: %w", err)
	//     }
	//     if user == nil || user.Email == "" {
	//         return "", fmt.Errorf("invalid user")
	//     }
	//     profile, err := api.GetProfile(user.Email)
	//     if err != nil {
	//         return "", fmt.Errorf("api error: %w", err)
	//     }
	//     return fmt.Sprintf("%s: %s", user.Name, profile.Bio), nil
	// }

	// After: gust Catch pattern (8 lines, 0 error checks, linear flow)
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

	fetchUserProfile := func(userID int) (r result.Result[string]) {
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

	res := fetchUserProfile(1)
	if res.IsOk() {
		fmt.Println(res.Unwrap())
	} else {
		fmt.Println("Error:", res.UnwrapErr())
	}
	// Output: Alice: Software developer
}
