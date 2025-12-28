# gust 🌬️

[![tag](https://img.shields.io/github/tag/andeya/gust.svg)](https://github.com/andeya/gust/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.23-%23007d9c)
[![GoDoc](https://godoc.org/github.com/andeya/gust?status.svg)](https://pkg.go.dev/github.com/andeya/gust)
![Build Status](https://github.com/andeya/gust/actions/workflows/go-ci.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/andeya/gust)](https://goreportcard.com/report/github.com/andeya/gust)
[![Coverage](https://img.shields.io/codecov/c/github/andeya/gust)](https://codecov.io/gh/andeya/gust)
[![License](https://img.shields.io/github/license/andeya/gust)](./LICENSE)

**Bring Rust's elegance to Go** - A powerful library that makes error handling, optional values, and iteration as beautiful and safe as in Rust.

> 🎯 **Zero dependencies** • 🚀 **Production ready** • 📚 **Well documented** • ✨ **Type-safe**

**Languages:** [English](./README.md) | [中文](./README_ZH.md)

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
import "github.com/andeya/gust/result"

func fetchUserData(userID int) gust.Result[string] {
    return result.AndThen(gust.Ret(getUser(userID)), func(user *User) gust.Result[string] {
        if user == nil || user.Email == "" {
            return gust.TryErr[string]("invalid user")
        }
        return result.Map(gust.Ret(getProfile(user.Email)), func(profile *Profile) string {
            return fmt.Sprintf("%s: %s", user.Name, profile.Bio)
        })
    })
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

// Chain operations that can fail
result := gust.Ok(10).
    Map(func(x int) int { return x * 2 }).
    AndThen(func(x int) gust.Result[int] {
        if x > 15 {
            return gust.TryErr[int]("too large")
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
import "github.com/andeya/gust/iterator"

numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

sum := iterator.FromSlice(numbers).
    Filter(func(x int) bool { return x%2 == 0 }).
    Map(func(x int) int { return x * x }).
    Take(3).
    Fold(0, func(acc int, x int) int {
        return acc + x
    })

fmt.Println(sum) // 56 (4 + 16 + 36)
```

**Available Methods:**
- **Constructors**: `FromSlice()`, `FromElements()`, `FromRange()`, `FromFunc()`, `FromIterable()`, `Empty()`, `Once()`, `Repeat()` - Create iterators from various sources
- **BitSet Iterators**: `FromBitSet()`, `FromBitSetOnes()`, `FromBitSetZeros()`, `FromBitSetBytes()`, `FromBitSetBytesOnes()`, `FromBitSetBytesZeros()` - Iterate over bits in bit sets or byte slices
- **Go Integration**: `FromSeq()`, `FromSeq2()`, `FromPull()`, `FromPull2()` - Convert from Go's standard iterators; `Seq()`, `Seq2()`, `Pull()`, `Pull2()` - Convert to Go's standard iterators
- **Adapters**: `Map`, `Filter`, `Chain`, `Zip`, `Enumerate`, `Skip`, `Take`, `StepBy`, `FlatMap`, `Flatten`
- **Consumers**: `Fold`, `Reduce`, `Collect`, `Count`, `All`, `Any`, `Find`, `Sum`, `Product`, `Partition`
- **Advanced**: `Scan`, `Intersperse`, `Peekable`, `ArrayChunks`, `FindMap`, `MapWhile`
- **Double-Ended**: `NextBack`, `Rfold`, `TryRfold`, `Rfind`
- And 60+ more methods from Rust's Iterator trait!

**Note:** For type-changing operations (e.g., `Map` from `string` to `int`), use the function-style API:
```go
iterator.Map(iterator.FromSlice(strings), func(s string) int { return len(s) })
```

For same-type operations, you can use method chaining:
```go
iterator.FromSlice(numbers).Filter(func(x int) bool { return x > 0 }).Map(func(x int) int { return x * 2 })
```

**Key Benefits:**
- ✅ Rust-like method chaining
- ✅ Lazy evaluation
- ✅ Type-safe transformations
- ✅ Zero-copy where possible

#### Iterator Constructors

Create iterators from various sources:

```go
import "github.com/andeya/gust/iterator"

// From slice
iter1 := iterator.FromSlice([]int{1, 2, 3})

// From individual elements
iter2 := iterator.FromElements(1, 2, 3)

// From range [start, end)
iter3 := iterator.FromRange(0, 5) // 0, 1, 2, 3, 4

// From function
count := 0
iter4 := iterator.FromFunc(func() gust.Option[int] {
    if count < 3 {
        count++
        return gust.Some(count)
    }
    return gust.None[int]()
})

// Empty iterator
iter5 := iterator.Empty[int]()

// Single value
iter6 := iterator.Once(42)

// Infinite repeat
iter7 := iterator.Repeat("hello") // "hello", "hello", "hello", ...
```

#### Go Standard Iterator Integration

gust iterators seamlessly integrate with Go 1.23+ standard iterators:

**Convert gust Iterator to Go's `iter.Seq[T]`:**
```go
import "github.com/andeya/gust/iterator"

numbers := []int{1, 2, 3, 4, 5}
gustIter := iterator.FromSlice(numbers).Filter(func(x int) bool { return x%2 == 0 })

// Use in Go's standard for-range loop
for v := range gustIter.Seq() {
    fmt.Println(v) // prints 2, 4
}
```

**Convert Go's `iter.Seq[T]` to gust Iterator:**
```go
import "github.com/andeya/gust/iterator"

// Create a Go standard iterator sequence
goSeq := func(yield func(int) bool) {
    for i := 0; i < 5; i++ {
        if !yield(i) {
            return
        }
    }
}

// Convert to gust Iterator and use gust methods
gustIter, deferStop := iterator.FromSeq(goSeq)
defer deferStop()
result := gustIter.Map(func(x int) int { return x * 2 }).Collect()
fmt.Println(result) // [0 2 4 6 8]
```

### 4. Double-Ended Iterator

Iterate from both ends:

```go
import "github.com/andeya/gust/iterator"

numbers := []int{1, 2, 3, 4, 5}
deIter := iterator.FromSlice(numbers).MustToDoubleEnded()

// Iterate from front
if val := deIter.Next(); val.IsSome() {
    fmt.Println("Front:", val.Unwrap()) // Front: 1
}

// Iterate from back
if val := deIter.NextBack(); val.IsSome() {
    fmt.Println("Back:", val.Unwrap()) // Back: 5
}
```

## 📖 Examples

### Parse and Filter with Error Handling

```go
import "github.com/andeya/gust"
import "github.com/andeya/gust/iterator"
import "strconv"

// Parse strings to integers, automatically filtering out errors
numbers := []string{"1", "2", "three", "4", "five"}

results := iterator.FilterMap(
    iterator.RetMap(iterator.FromSlice(numbers), strconv.Atoi),
    gust.Result[int].Ok,
).
    Collect()

fmt.Println("Parsed numbers:", results)
// Output: Parsed numbers: [1 2 4]
```

### Real-World Data Pipeline

```go
import "github.com/andeya/gust"
import "github.com/andeya/gust/iterator"
import "strconv"

// Process user input: parse, validate, transform, limit
input := []string{"10", "20", "invalid", "30", "0", "40"}

results := iterator.FilterMap(
    iterator.RetMap(iterator.FromSlice(input), strconv.Atoi),
    gust.Result[int].Ok,
).
    Filter(func(x int) bool { return x > 0 }).
    Map(func(x int) int { return x * 2 }).
    Take(3).
    Collect()

fmt.Println(results) // [20 40 60]
```

### Option Chain Operations

```go
import "github.com/andeya/gust"
import "fmt"

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
import "github.com/andeya/gust/iterator"
import "fmt"

// Split numbers into evens and odds
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

evens, odds := iterator.FromSlice(numbers).
    Partition(func(x int) bool { return x%2 == 0 })

fmt.Println("Evens:", evens) // [2 4 6 8 10]
fmt.Println("Odds:", odds)   // [1 3 5 7 9]
```

### BitSet Iteration

Iterate over bits in bit sets or byte slices with full iterator support:

```go
import "github.com/andeya/gust/iterator"
import "fmt"

// Iterate over bits in a byte slice
bytes := []byte{0b10101010, 0b11001100}

// Get all set bit offsets
setBits := iterator.FromBitSetBytesOnes(bytes).
    Filter(func(offset int) bool { return offset > 5 }).
    Collect()
fmt.Println(setBits) // [6 8 9 12 13]

// Count set bits
count := iterator.FromBitSetBytesOnes(bytes).Count()
fmt.Println(count) // 8

// Sum of offsets of set bits
sum := iterator.FromBitSetBytesOnes(bytes).
    Fold(0, func(acc, offset int) int { return acc + offset })
fmt.Println(sum) // 54 (0+2+4+6+8+9+12+13)

// Works with any BitSetLike implementation
type MyBitSet struct {
    bits []byte
}

func (b *MyBitSet) Size() int { return len(b.bits) * 8 }
func (b *MyBitSet) Get(offset int) bool {
    if offset < 0 || offset >= b.Size() {
        return false
    }
    byteIdx := offset / 8
    bitIdx := offset % 8
    return (b.bits[byteIdx] & (1 << (7 - bitIdx))) != 0
}

bitset := &MyBitSet{bits: []byte{0b10101010}}
ones := iterator.FromBitSetOnes(bitset).Collect()
fmt.Println(ones) // [0 2 4 6]
```

## 📦 Additional Packages

gust provides several utility packages to extend its functionality:

- **`gust/dict`** - Generic map utilities (Filter, Map, Keys, Values, etc.)
- **`gust/vec`** - Generic slice utilities  
- **`gust/conv`** - Type-safe value conversion and reflection utilities
- **`gust/digit`** - Number conversion utilities (base conversion, etc.)
- **`gust/opt`** - Helper functions for `Option[T]` (Map, AndThen, Zip, Unzip, Assert, etc.)
- **`gust/result`** - Helper functions for `Result[T]` (Map, AndThen, Assert, Flatten, etc.)
- **`gust/iterator`** - Rust-like iterator implementation (see [Iterator section](#3-iterator---rust-like-iteration-in-go) above)
- **`gust/syncutil`** - Concurrent utilities (SyncMap, Mutex wrappers, Lazy initialization, etc.)
- **`gust/errutil`** - Error utilities (Stack traces, Panic recovery, etc.)
- **`gust/constraints`** - Type constraints (Ordering, Numeric, etc.)

### Quick Examples

**Dict utilities:**
```go
import "github.com/andeya/gust/dict"

m := map[string]int{"a": 1, "b": 2, "c": 3}
value := dict.Get(m, "b").UnwrapOr(0) // 2
filtered := dict.Filter(m, func(k string, v int) bool { return v > 1 })
```

**Vec utilities:**
```go
import "github.com/andeya/gust/vec"

numbers := []int{1, 2, 3, 4, 5}
doubled := vec.MapAlone(numbers, func(x int) int { return x * 2 })
```

**Opt utilities:**
```go
import "github.com/andeya/gust/opt"

some := gust.Some(5)
doubled := opt.Map(some, func(x int) int { return x * 2 })
```

**Result utilities:**
```go
import "github.com/andeya/gust"
import "github.com/andeya/gust/result"

result := gust.Ok(10)
doubled := result.Map(result, func(x int) int { return x * 2 })
```

For more details, see the [full documentation](https://pkg.go.dev/github.com/andeya/gust) and [examples](./examples/).

### Detailed Examples

#### Dict Utilities

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

#### Vec Utilities
```go
import "github.com/andeya/gust/vec"

numbers := []int{1, 2, 3, 4, 5}
doubled := vec.MapAlone(numbers, func(x int) int { return x * 2 })
fmt.Println(doubled) // [2 4 6 8 10]
```

#### Opt Utilities
```go
import "github.com/andeya/gust/opt"

some := gust.Some(5)
doubled := opt.Map(some, func(x int) int { return x * 2 })
zipped := opt.Zip(gust.Some(1), gust.Some("hello"))
```

#### Result Utilities
```go
import "github.com/andeya/gust"
import "github.com/andeya/gust/result"

result := gust.Ok(10)
doubled := result.Map(result, func(x int) int { return x * 2 })
chained := result.AndThen(gust.Ok(5), func(x int) gust.Result[int] {
    return gust.Ok(x * 2)
})
```

#### Conv Utilities
```go
import "github.com/andeya/gust/conv"

// Type-safe conversions
value := conv.To[int]("42") // Returns Option[int]

// Reflection utilities
if conv.IsNil(someValue) {
    // Handle nil case
}
```

#### SyncUtil Utilities
```go
import "github.com/andeya/gust/syncutil"

// Thread-safe map
var m syncutil.SyncMap[string, int]
m.Store("key", 42)
value := m.Load("key") // Returns Option[int]

// Lazy initialization
lazy := syncutil.NewLazy(func() int {
    return expensiveComputation()
})
value := lazy.Get() // Computed only once
```

#### Digit Utilities
```go
import "github.com/andeya/gust/digit"

// Base conversion (e.g., base62)
encoded := digit.Itoa62(12345) // Convert to base62 string
decoded := digit.Atoi62(encoded) // Convert back, returns Result[int]
```

## 🔗 Resources

- 📖 [Full Documentation](https://pkg.go.dev/github.com/andeya/gust) - Complete API reference
- 💡 [Examples](./examples/) - Comprehensive examples organized by feature
- 🌐 [中文文档](./README_ZH.md) - Chinese documentation
- 🐛 [Issue Tracker](https://github.com/andeya/gust/issues) - Report bugs or request features
- 💬 [Discussions](https://github.com/andeya/gust/discussions) - Ask questions and share ideas

## 📋 Requirements

Requires **Go 1.23+** (for generics and standard iterator support)

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
