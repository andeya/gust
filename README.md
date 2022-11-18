# gust [![Docs](https://img.shields.io/badge/Docs-pkg.go.dev-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/andeya/gust)

A Rust-inspired **declarative generic programming module** for Golang that helps reduce bugs and improve development efficiency. 

![declarative_vs_imperative.jpg](doc/declarative_vs_imperative.jpg)

After using this package, your code style will be like this:

```go
package examples_test

import (
	"errors"
	"fmt"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/andeya/gust/ret"
)

type Version int8

const (
	Version1 Version = iota + 1
	Version2
)

func ParseVersion(header iter.Iterator[byte]) gust.Result[Version] {
	return ret.AndThen(
		header.Next().
			OkOr("invalid header length"),
		func(b byte) gust.Result[Version] {
			switch b {
			case 1:
				return gust.Ok(Version1)
			case 2:
				return gust.Ok(Version2)
			}
			return gust.Err[Version]("invalid version")
		},
	)
}

func ExampleVersion() {
	ParseVersion(iter.FromElements[byte](1, 2, 3, 4)).
		Inspect(func(v Version) {
			fmt.Printf("working with version: %v\n", v)
		}).
		InspectErr(func(err error) {
			fmt.Printf("error parsing header: %v\n", err)
		})
	// Output:
	// working with version: 1
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
- `iter` is a package that provides a generic-type iterator type.
- `vec` is a toolkit for efficient handling of generic-type slices.
- `valconv` is a package that provides a generic-type value converter.
- `digit` is a package that provides generic-type digit operations.
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

### Iterator

Feature-rich iterators.

- Iterator Example

```go
func TestAny(t *testing.T) {
	var iter = FromVec([]int{1, 2, 3})
	if !iter.Any(func(x int) bool {
		return x > 1
	}) {
		t.Error("Any failed")
	}
}
```
