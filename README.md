# gust 🌬️

[![tag](https://img.shields.io/github/tag/andeya/gust.svg)](https://github.com/andeya/gust/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.19-%23007d9c)
[![GoDoc](https://godoc.org/github.com/andeya/gust?status.svg)](https://pkg.go.dev/github.com/andeya/gust)
![Build Status](https://github.com/andeya/gust/actions/workflows/go-ci.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/andeya/gust)](https://goreportcard.com/report/github.com/andeya/gust)
[![Coverage](https://img.shields.io/codecov/c/github/andeya/gust)](https://codecov.io/gh/andeya/gust)
[![License](https://img.shields.io/github/license/andeya/gust)](./LICENSE)

**Bring Rust's elegance to Go** - A powerful library that makes error handling, optional values, and iteration as beautiful and safe as in Rust.

> 🎯 **Zero dependencies** • 🚀 **Production ready** • 📚 **Well documented** • ✨ **Type-safe**

## ✨ Why gust?

Tired of writing `if err != nil` everywhere? Frustrated with nil pointer panics? Want Rust-like iterator chains in Go?

**gust** brings Rust's best patterns to Go, making your code:
- 🛡️ **Safer** - No more nil pointer panics
- 🎯 **Cleaner** - Chain operations elegantly
- 🚀 **More Expressive** - Write what you mean, not boilerplate

### From Imperative to Declarative

gust helps you shift from **imperative** (focusing on *how*) to **declarative** (focusing on *what*) programming:

![Declarative vs Imperative](./doc/declarative_vs_imperative.jpg)

With gust, you describe **what** you want to achieve, not **how** to achieve it step-by-step. This makes your code more readable, maintainable, and less error-prone.

### Before gust (Traditional Go)
```go
func fetchUserData(userID int) (string, error) {
    // Step 1: Fetch from database
    user, err := db.GetUser(userID)
    if err != nil {
        return "", fmt.Errorf("db error: %w", err)
    }
    
    // Step 2: Validate user
    if user == nil {
        return "", fmt.Errorf("user not found")
    }
    if user.Email == "" {
        return "", fmt.Errorf("invalid user: no email")
    }
    
    // Step 3: Fetch profile
    profile, err := api.GetProfile(user.Email)
    if err != nil {
        return "", fmt.Errorf("api error: %w", err)
    }
    
    // Step 4: Format result
    result := fmt.Sprintf("%s: %s", user.Name, profile.Bio)
    return result, nil
}
```

### After gust (Elegant & Safe)
```go
import "github.com/andeya/gust"
import "github.com/andeya/gust/ret"

func fetchUserData(userID int) gust.Result[string] {
    return ret.AndThen(
        gust.Ret(db.GetUser(userID)),
        func(user *User) gust.Result[string] {
            if user == nil || user.Email == "" {
                return gust.Err[string]("invalid user")
            }
            return ret.AndThen(
                gust.Ret(api.GetProfile(user.Email)),
                func(profile *Profile) gust.Result[string] {
                    return gust.Ok(fmt.Sprintf("%s: %s", user.Name, profile.Bio))
                },
            )
        },
    )
}

// See ExampleResult_fetchUserData in examples/ for a complete runnable example
```

**What changed?**
- ✅ **No error boilerplate** - Errors flow naturally through the chain
- ✅ **No nested if-else** - Linear flow, easy to read
- ✅ **Automatic propagation** - Errors stop the chain automatically
- ✅ **Composable** - Each step is independent and testable
- ✅ **Type-safe** - Compiler enforces correct error handling

## 🚀 Quick Start

```bash
go get github.com/andeya/gust
```

## 📚 Core Features

### 1. Result<T> - Elegant Error Handling

Replace `(T, error)` with chainable `Result[T]`:

```go
import "github.com/andeya/gust"
import "github.com/andeya/gust/ret"

// Chain operations that can fail
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
        return gust.Ok(0) // Fallback
    })

fmt.Println("Final value:", result.Unwrap())
// Output: Error handled: too large
// Final value: 0
```

**Key Benefits:**
- ✅ No more `if err != nil` boilerplate
- ✅ Automatic error propagation
- ✅ Chain multiple operations elegantly
- ✅ Type-safe error handling

### 2. Option<T> - No More Nil Panics

Replace `*T` and `(T, bool)` with safe `Option[T]`:

```go
// Safe division without nil checks
divide := func(a, b float64) gust.Option[float64] {
    if b == 0 {
        return gust.None[float64]()
    }
    return gust.Some(a / b)
}

result := divide(10, 2).
    Map(func(x float64) float64 { return x * 2 }).
    UnwrapOr(0)

fmt.Println(result) // 10
```

**Key Benefits:**
- ✅ Eliminates nil pointer panics
- ✅ Explicit optional values
- ✅ Chain operations safely
- ✅ Compiler-enforced safety

### 3. Iterator - Rust-like Iteration in Go

Full Rust Iterator trait implementation with method chaining:

```go
import "github.com/andeya/gust/iter"

numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

sum := iter.FromSlice(numbers).
    Filter(func(x int) bool { return x%2 == 0 }).
    Map(func(x int) int { return x * x }).
    Take(3).
    Fold(0, func(acc int, x int) int {
        return acc + x
    })

fmt.Println(sum) // 56 (4 + 16 + 36)
```

**Available Methods:**
- **Adapters**: `Map`, `Filter`, `Chain`, `Zip`, `Enumerate`, `Skip`, `Take`, `StepBy`, `FlatMap`, `Flatten`
- **Consumers**: `Fold`, `Reduce`, `Collect`, `Count`, `All`, `Any`, `Find`, `Sum`, `Product`, `Partition`
- **Advanced**: `Scan`, `Intersperse`, `Peekable`, `ArrayChunks`, `FindMap`, `MapWhile`
- **Double-Ended**: `NextBack`, `Rfold`, `TryRfold`, `Rfind`
- And 60+ more methods from Rust's Iterator trait!

**Note:** For type-changing operations (e.g., `Map` from `string` to `int`), use the function-style API:
```go
iter.Map(iter.FromSlice(strings), func(s string) int { return len(s) })
```

For same-type operations, you can use method chaining:
```go
iter.FromSlice(numbers).Filter(func(x int) bool { return x > 0 }).Map(func(x int) int { return x * 2 })
```

**Key Benefits:**
- ✅ Rust-like method chaining
- ✅ Lazy evaluation
- ✅ Type-safe transformations
- ✅ Zero-copy where possible

### 4. Double-Ended Iterator

Iterate from both ends:

```go
import "github.com/andeya/gust/iter"

numbers := []int{1, 2, 3, 4, 5}
deIter := iter.FromSlice(numbers).MustToDoubleEnded()

// Iterate from front
if val := deIter.Next(); val.IsSome() {
    fmt.Println("Front:", val.Unwrap()) // Front: 1
}

// Iterate from back
if val := deIter.NextBack(); val.IsSome() {
    fmt.Println("Back:", val.Unwrap()) // Back: 5
}
```

## 📖 More Examples

### Parse and Filter with Error Handling

```go
import "github.com/andeya/gust"
import "github.com/andeya/gust/iter"
import "strconv"

// Parse strings to integers, automatically filtering out errors
numbers := []string{"1", "2", "three", "4", "five"}

results := iter.FilterMap(
    iter.Map(iter.FromSlice(numbers), func(s string) gust.Result[int] {
        return gust.Ret(strconv.Atoi(s))
    }),
    gust.Result[int].Ok).
    Collect()

fmt.Println("Parsed numbers:", results)
// Output: Parsed numbers: [1 2 4]
```

### Real-World Data Pipeline

```go
// Process user input: parse, validate, transform, limit
input := []string{"10", "20", "invalid", "30", "0", "40"}

results := iter.FilterMap(
    iter.Map(iter.FromSlice(input), func(s string) gust.Result[int] {
        return gust.Ret(strconv.Atoi(s))
    }),
    gust.Result[int].Ok).
    Filter(func(x int) bool { return x > 0 }).
    Map(func(x int) int { return x * 2 }).
    Take(3).
    Collect()

fmt.Println(results) // [20 40 60]
```

### Option Chain Operations

```go
// Chain operations on optional values with filtering
result := gust.Some(5).
    Map(func(x int) int { return x * 2 }).
    Filter(func(x int) bool { return x > 8 }).
    XMap(func(x int) any {
        return fmt.Sprintf("Value: %d", x)
    }).
    UnwrapOr("No value")

fmt.Println(result) // "Value: 10"
```

### Partition Data

```go
// Split numbers into evens and odds
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

evens, odds := iter.FromSlice(numbers).
    Partition(func(x int) bool { return x%2 == 0 })

fmt.Println("Evens:", evens) // [2 4 6 8 10]
fmt.Println("Odds:", odds)   // [1 3 5 7 9]
```

## 🎯 Use Cases

- **Error Handling**: Replace `(T, error)` with `Result[T]` for cleaner, chainable error handling
- **Optional Values**: Use `Option[T]` instead of `*T` for nil safety and explicit optional semantics
- **Data Processing**: Chain iterator operations for elegant, lazy-evaluated data transformations
- **API Responses**: Handle optional/error cases explicitly without nil checks
- **Configuration**: Use `Option` for optional config fields with type safety
- **Data Validation**: Combine `Result` and `Option` for robust input validation pipelines

## 📦 Additional Packages

- **`gust/dict`** - Generic map utilities (Filter, Map, Keys, Values, etc.)
- **`gust/vec`** - Generic slice utilities  
- **`gust/valconv`** - Type-safe value conversion
- **`gust/digit`** - Number conversion utilities
- **`gust/sync`** - Generic sync primitives (Mutex, RWMutex, etc.)

### Dict Utilities Example

```go
import "github.com/andeya/gust/dict"

m := map[string]int{"a": 1, "b": 2, "c": 3}

// Get with Option
value := dict.Get(m, "b")
fmt.Println(value.UnwrapOr(0)) // 2

// Filter map
filtered := dict.Filter(m, func(k string, v int) bool {
    return v > 1
})
fmt.Println(filtered) // map[b:2 c:3]

// Map values
mapped := dict.MapValue(m, func(k string, v int) int {
    return v * 2
})
fmt.Println(mapped) // map[a:2 b:4 c:6]
```

### Vec Utilities Example

```go
import "github.com/andeya/gust/vec"

// Map slice elements
numbers := []int{1, 2, 3, 4, 5}
doubled := vec.MapAlone(numbers, func(x int) int {
    return x * 2
})
fmt.Println(doubled) // [2 4 6 8 10]

// Convert []any to specific type
anySlice := []any{1, 2, 3, 4, 5}
intSlice := vec.MapAlone(anySlice, func(v any) int {
    return v.(int)
})
fmt.Println(intSlice) // [1 2 3 4 5]
```

## 🔗 Resources

- 📖 [Full Documentation](https://pkg.go.dev/github.com/andeya/gust) - Complete API reference
- 💡 [Examples](./examples/) - Comprehensive examples organized by feature
- 🐛 [Issue Tracker](https://github.com/andeya/gust/issues) - Report bugs or request features
- 💬 [Discussions](https://github.com/andeya/gust/discussions) - Ask questions and share ideas

## 📋 Go Version

Requires **Go 1.19+** (for generics support)

## 🤝 Contributing

Contributions are welcome! Whether it's:
- 🐛 Reporting bugs
- 💡 Suggesting new features
- 📝 Improving documentation
- 🔧 Submitting pull requests

Every contribution makes gust better! Please feel free to submit a Pull Request or open an issue.

## 📄 License

This project is licensed under the MIT License.

---

**Made with ❤️ for the Go community**

*Inspired by Rust's `Result`, `Option`, and `Iterator` traits*
