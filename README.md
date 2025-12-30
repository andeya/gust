<div align="center">

# gust 🌬️

**Bring Rust's Elegance to Go**

*A production-ready library that makes error handling, optional values, and iteration as beautiful and safe as in Rust.*

[![GitHub release](https://img.shields.io/github/release/andeya/gust.svg)](https://github.com/andeya/gust/releases)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.24-00ADD8?style=flat&logo=go)](https://golang.org)
[![GoDoc](https://pkg.go.dev/badge/github.com/andeya/gust.svg)](https://pkg.go.dev/github.com/andeya/gust)
[![CI Status](https://github.com/andeya/gust/actions/workflows/go-ci.yml/badge.svg)](https://github.com/andeya/gust/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/andeya/gust)](https://goreportcard.com/report/github.com/andeya/gust)
[![Code Coverage](https://codecov.io/gh/andeya/gust/branch/main/graph/badge.svg)](https://codecov.io/gh/andeya/gust)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)

**[English](./README.md)** | **[中文](./README_ZH.md)**

</div>

---

## 🎯 What is gust?

**gust** is a comprehensive Go library that brings Rust's most powerful patterns to Go, enabling you to write **safer, cleaner, and more expressive code**. With **zero dependencies** and **production-ready** quality, gust transforms how you handle errors, optional values, and data iteration in Go.

### ✨ Why gust?

| Traditional Go | With gust |
|----------------|-----------|
| ❌ Verbose error handling | ✅ Chainable `Result[T]` |
| ❌ Nil pointer panics | ✅ Safe `Option[T]` |
| ❌ Imperative loops | ✅ Declarative iterators |
| ❌ Boilerplate code | ✅ Elegant composition |

---

## 🚀 Quick Start

```bash
go get github.com/andeya/gust
```

### 30-Second Example

```go
package main

import (
    "fmt"
    "github.com/andeya/gust/result"
)

func main() {
    // Chain operations elegantly - no error boilerplate!
    res := result.Ok(10).
        Map(func(x int) int { return x * 2 }).
        AndThen(func(x int) result.Result[int] {
            if x > 20 {
                return result.TryErr[int]("too large")
            }
            return result.Ok(x + 5)
        })

    if res.IsOk() {
        fmt.Println("Success:", res.Unwrap()) // Success: 25 (⚠️ Unwrap may panic if not checked)
    }
}
```

---

## 💡 The Problem gust Solves

### Before: Traditional Go Code

```go
func fetchUserData(userID int) (string, error) {
    user, err := db.GetUser(userID)
    if err != nil {
        return "", fmt.Errorf("db error: %w", err)
    }
    if user == nil {
        return "", fmt.Errorf("user not found")
    }
    if user.Email == "" {
        return "", fmt.Errorf("invalid user: no email")
    }
    profile, err := api.GetProfile(user.Email)
    if err != nil {
        return "", fmt.Errorf("api error: %w", err)
    }
    return fmt.Sprintf("%s: %s", user.Name, profile.Bio), nil
}
```

**Problems:**
- ❌ Repetitive `if err != nil` checks
- ❌ Nested conditionals
- ❌ Hard to compose and test
- ❌ Easy to forget error handling

### After: With gust

```go
import "github.com/andeya/gust/result"

func fetchUserData(userID int) result.Result[string] {
    return result.Ret(db.GetUser(userID)).
        AndThen(func(user *User) result.Result[string] {
            if user == nil || user.Email == "" {
                return result.TryErr[string]("invalid user")
            }
            return result.Ret(api.GetProfile(user.Email)).
                Map(func(profile *Profile) string {
                    return fmt.Sprintf("%s: %s", user.Name, profile.Bio)
                })
        })
}
```

**Benefits:**
- ✅ **No error boilerplate** - Errors flow naturally
- ✅ **Linear flow** - Easy to read and understand
- ✅ **Automatic propagation** - Errors stop the chain automatically
- ✅ **Composable** - Each step is independent and testable
- ✅ **Type-safe** - Compiler enforces correct error handling

---

## 📚 Core Features

### 1. Result<T> - Type-Safe Error Handling

Replace `(T, error)` with chainable `Result[T]`:

```go
import "github.com/andeya/gust/result"

res := result.Ok(10).
    Map(func(x int) int { return x * 2 }).
    AndThen(func(x int) result.Result[int] {
        if x > 15 {
            return result.TryErr[int]("too large")
        }
        return result.Ok(x + 5)
    }).
    OrElse(func(err error) result.Result[int] {
        return result.Ok(0) // Fallback
    })

fmt.Println(res.UnwrapOr(0)) // 25 (safe, returns 0 if error)
// Or check first (Unwrap may panic if not checked):
if res.IsOk() {
    fmt.Println(res.Unwrap()) // 25 (panics if error, only use after IsOk() check)
}
```

**Key Methods:**
- `Map` - Transform value if Ok
- `AndThen` - Chain operations returning Result
- `OrElse` - Handle errors with fallback
- `UnwrapOr` - Extract values safely (with default, **never panics**)
- `Unwrap` - Extract value (⚠️ **panics if error** - use only after `IsOk()` check, prefer `UnwrapOr` for safety)

### 2. Option<T> - No More Nil Panics

Replace `*T` and `(T, bool)` with safe `Option[T]`:

```go
import "github.com/andeya/gust/option"

divide := func(a, b float64) option.Option[float64] {
    if b == 0 {
        return option.None[float64]()
    }
    return option.Some(a / b)
}

res := divide(10, 2).
    Map(func(x float64) float64 { return x * 2 }).
    UnwrapOr(0)

fmt.Println(res) // 10
```

**Key Methods:**
- `Map` - Transform value if Some
- `AndThen` - Chain operations returning Option
- `Filter` - Conditionally filter values
- `UnwrapOr` - Extract values safely (with default, **never panics**)
- `Unwrap` - Extract value (⚠️ **panics if None** - use only after `IsSome()` check, prefer `UnwrapOr` for safety)

### 3. Iterator - Rust-like Iteration

Full Rust Iterator trait implementation with **60+ methods**:

```go
import "github.com/andeya/gust/iterator"

numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

sum := iterator.FromSlice(numbers).
    Filter(func(x int) bool { return x%2 == 0 }).
    Map(func(x int) int { return x * x }).
    Take(3).
    Fold(0, func(acc, x int) int { return acc + x })

fmt.Println(sum) // 56 (4 + 16 + 36)
```

**Highlights:**
- 🚀 **60+ methods** from Rust's Iterator trait
- 🔄 **Lazy evaluation** - Computations happen on-demand
- 🔗 **Method chaining** - Compose complex operations elegantly
- 🔌 **Go 1.24+ integration** - Works with standard `iter.Seq[T]`
- 🎯 **Type-safe** - Compile-time guarantees

**Method Categories:**
- **Constructors**: `FromSlice`, `FromRange`, `FromFunc`, `Empty`, `Once`, `Repeat`
- **BitSet Iterators**: `FromBitSet`, `FromBitSetOnes`, `FromBitSetZeros`
- **Go Integration**: `FromSeq`, `Seq`, `Pull` (Go 1.24+ standard iterators)
- **Basic Adapters**: `Map`, `Filter`, `Chain`, `Zip`, `Enumerate`
- **Filtering**: `Skip`, `Take`, `StepBy`, `SkipWhile`, `TakeWhile`
- **Transforming**: `MapWhile`, `Scan`, `FlatMap`, `Flatten`
- **Chunking**: `MapWindows`, `ArrayChunks`, `ChunkBy`
- **Consumers**: `Collect`, `Fold`, `Reduce`, `Count`, `Sum`, `Product`, `Partition`
- **Search**: `Find`, `FindMap`, `Position`, `All`, `Any`
- **Min/Max**: `Max`, `Min`, `MaxBy`, `MinBy`, `MaxByKey`, `MinByKey`
- **Double-Ended**: `NextBack`, `Rfold`, `Rfind`, `NthBack`

---

## 🌟 Real-World Examples

### Data Processing Pipeline

```go
import (
    "github.com/andeya/gust/iterator"
    "github.com/andeya/gust/result"
    "strconv"
)

// Parse, validate, transform, and limit user input
input := []string{"10", "20", "invalid", "30", "0", "40"}

results := iterator.FilterMap(
    iterator.RetMap(iterator.FromSlice(input), strconv.Atoi),
    result.Result[int].Ok,
).
    Filter(func(x int) bool { return x > 0 }).
    Map(func(x int) int { return x * 2 }).
    Take(3).
    Collect()

fmt.Println(results) // [20 40 60]
```

### Option Chain Operations

```go
import "github.com/andeya/gust/option"

res := option.Some(5).
    Map(func(x int) int { return x * 2 }).
    Filter(func(x int) bool { return x > 8 }).
    UnwrapOr("No value")

fmt.Println(res) // 10
```

### BitSet with Iterators

```go
import (
    "github.com/andeya/gust/bitset"
    "github.com/andeya/gust/iterator"
)

bs := bitset.New()
bs.Set(0, true).Unwrap()
bs.Set(5, true).Unwrap()

// Get all set bits using iterator
setBits := iterator.FromBitSetOnes(bs).Collect() // [0 5]

// Bitwise operations
bs1 := bitset.NewFromString("c0", bitset.EncodingHex).Unwrap()
bs2 := bitset.NewFromString("30", bitset.EncodingHex).Unwrap()
or := bs1.Or(bs2)

// Encoding/decoding (Base64URL by default)
encoded := bs.String()
decoded := bitset.NewFromBase64URL(encoded).Unwrap()
```

---

## 📦 Additional Packages

gust provides comprehensive utility packages:

| Package | Description | Key Features |
|---------|-------------|--------------|
| **`gust/dict`** | Generic map utilities | `Filter`, `Map`, `Keys`, `Values`, `Get` |
| **`gust/vec`** | Generic slice utilities | `MapAlone`, `Get`, `Copy`, `Dict` |
| **`gust/conv`** | Type-safe conversions | `BytesToString`, `StringToReadonlyBytes`, reflection utils |
| **`gust/digit`** | Number conversions | Base 2-62 conversion, `FormatByDict`, `ParseByDict` |
| **`gust/random`** | Secure random strings | Base36/Base62 encoding, timestamp embedding |
| **`gust/encrypt`** | Cryptographic functions | MD5, SHA series, FNV, CRC, Adler-32, AES encryption |
| **`gust/bitset`** | Thread-safe bit sets | Bitwise ops, iterator integration, multiple encodings |
| **`gust/syncutil`** | Concurrent utilities | `SyncMap`, `Lazy`, mutex wrappers |
| **`gust/errutil`** | Error utilities | Stack traces, panic recovery, `ErrBox` |
| **`gust/constraints`** | Type constraints | `Ordering`, `Numeric`, `Digit` |

---

## 🔗 Resources

- 📖 **[Full Documentation](https://pkg.go.dev/github.com/andeya/gust)** - Complete API reference
- 💡 **[Examples](./examples/)** - Comprehensive examples by feature
- 🌐 **[中文文档](./README_ZH.md)** - Chinese documentation
- 🐛 **[Issue Tracker](https://github.com/andeya/gust/issues)** - Report bugs or request features
- 💬 **[Discussions](https://github.com/andeya/gust/discussions)** - Ask questions and share ideas

---

## 📋 Requirements

- **Go 1.24+** (required for generics and standard iterator support)

---

## 🤝 Contributing

We welcome contributions! Whether you're:

- 🐛 **Reporting bugs** - Help us improve
- 💡 **Suggesting features** - Share your ideas
- 📝 **Improving docs** - Make documentation better
- 🔧 **Submitting PRs** - Contribute code improvements

Every contribution makes gust better!

### Development Setup

```bash
# Clone the repository
git clone https://github.com/andeya/gust.git
cd gust

# Run tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](./LICENSE) file for details.

---

<div align="center">

**Made with ❤️ for the Go community**

*Inspired by Rust's `Result`, `Option`, and `Iterator` traits*

[⭐ Star us on GitHub](https://github.com/andeya/gust) • [📖 Documentation](https://pkg.go.dev/github.com/andeya/gust) • [🐛 Report Bug](https://github.com/andeya/gust/issues) • [💡 Request Feature](https://github.com/andeya/gust/issues/new)

</div>
