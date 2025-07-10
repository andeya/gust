# gust

[![tag](https://img.shields.io/github/tag/andeya/gust.svg)](https://github.com/andeya/gust/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.19-%23007d9c)
[![GoDoc](https://godoc.org/github.com/andeya/gust?status.svg)](https://pkg.go.dev/github.com/andeya/gust)
![Build Status](https://github.com/andeya/gust/actions/workflows/go-ci.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/andeya/gust)](https://goreportcard.com/report/github.com/andeya/gust)
[![Coverage](https://img.shields.io/codecov/c/github/andeya/gust)](https://codecov.io/gh/andeya/gust)
[![License](https://img.shields.io/github/license/andeya/gust)](./LICENSE)

A Rust-inspired **declarative-programming and generic-type module** for Golang that helps avoid bugs and improve development efficiency. 

![declarative_vs_imperative.jpg](doc/declarative_vs_imperative.jpg)

After using this package, your code style will be like this:

```go
import (
	"github.com/andeya/gust"
	"github.com/andeya/gust/dict"
	"github.com/andeya/gust/valconv"
)

func Convert() (result gust.Result[map[string]int64]) {
	defer result.Catch()
	var a = map[string]any{"a": 1, "b": 2, "c": 3}
	return gust.Ok(dict.MapValue(
		valconv.SafeAssertMap[string, int](a).UnwrapOrThrow(),
		func(k string, v int) int64 {
			return int64(v) + 1
		}),
	)
}
```

## Go Version

go≥1.19

## Features

- `gust.Result` is a type that represents either a success or an error.
- `gust.Option` is a type that represents either a value or nothing.
- `gust.Mutex` is a better generic-type wrapper for `sync.Mutex` that holds a value.
- `gust.RWMutex` is a better generic-type wrapper for `sync.RWMutex` that holds a value.
- `gust.SyncMap` is a better generic-type wrapper for `sync.Map`.
- `gust.AtomicValue` is a better generic-type wrapper for `atomic.Value`.
- `vec` is a package of generic-type functions for slices.
- `valconv` is a package that provides a generic-type value converter.
- `digit` is a package of generic-type functions for digit.
- and more...

### Result

Improve `func() (T,error)`, handle result with chain methods.

- Result Example

```go
func TestResult(t *testing.T) {
	var goodResult1 = gust.Ok(10)
	var badResult1 = gust.Err[int](10)

	// The `IsOk` and `IsErr` methods do what they say.
	assert.True(t, goodResult1.IsOk() && !goodResult1.IsErr())
	assert.True(t, badResult1.IsErr() && !badResult1.IsOk())

	// `map` consumes the `Result` and produces another.
	var goodResult2 = goodResult1.Map(func(i int) int { return i + 1 })
	var badResult2 = badResult1.Map(func(i int) int { return i - 1 })

	// Use `AndThen` to continue the computation.
	var goodResult3 = ret.AndThen(goodResult2, func(i int) gust.Result[bool] { return gust.Ok(i == 11) })

	// Use `OrElse` to handle the error.
	var _ = badResult2.OrElse(func(err error) gust.Result[int] {
		fmt.Println(err)
		return gust.Ok(20)
	})

	// Consume the result and return the contents with `Unwrap`.
	var _ = goodResult3.Unwrap()
}
```

### Option

Improve `func()(T, bool)` and `if *U != nil`, handle value with `Option` type.

Type [`Option`] represents an optional value, and has a number of uses:
* Initial values
* Return values for functions that are not defined
  over their entire input range (partial functions)
* Return value for otherwise reporting simple errors, where [`None`] is
  returned on error
* Optional struct fields
* Optional function arguments
* Nil-able pointers

- Option Example

```go
func TestOption(t *testing.T) {
	var divide = func(numerator, denominator float64) gust.Option[float64] {
		if denominator == 0.0 {
			return gust.None[float64]()
		}
		return gust.Some(numerator / denominator)
	}
	// The return value of the function is an option
	divide(2.0, 3.0).
		Inspect(func(x float64) {
			// Pattern match to retrieve the value
			t.Log("Result:", x)
		}).
		InspectNone(func() {
			t.Log("Cannot divide by 0")
		})
}
```

### Errable

Improve `func() error`, handle error with chain methods.

- Errable Example

```go
func ExampleErrable() {
	var hasErr = true
	var f = func() gust.Errable[int] {
		if hasErr {
			return gust.ToErrable(1)
		}
		return gust.NonErrable[int]()
	}
	var r = f()
	fmt.Println(r.IsErr())
	fmt.Println(r.UnwrapErr())
	fmt.Printf("%#v", r.ToError())
	// Output:
	// true
	// 1
	// &gust.errorWithVal{val:1}
}
```
