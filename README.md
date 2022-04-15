# gust [![Docs](https://img.shields.io/badge/Docs-pkg.go.dev-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/henrylee2cn/gust)

Golang basic module inspired by rust.

## Go Version

go≥1.18

## Features

### Result

Go-generics result module inspired by rust.

Avoid ifelse, handle result with chain methods, will you choose her?

- Result Example

```go
func ExampleResult_AndThen() {
	var divide = func(i, j float32) gust.Result[float32] {
		if j == 0 {
			return gust.Err[float32]("j can not be 0")
		}
		return gust.Ok(i / j)
	}
	var ret float32 = divide(1, 2).AndThen(func(i float32) gust.Result[float32] {
		return gust.Ok(i * 10)
	}).Unwrap()
	fmt.Println(ret)
	// Output:
	// 5
}
```

```go
func ExampleResult_UnwrapOr() {
	const def int = 10

	// before
	i, err := strconv.Atoi("1")
	if err != nil {
		i = def
	}
	fmt.Println(i * 2)

	// now
	fmt.Println(gust.Ret(strconv.Atoi("1")).UnwrapOr(def) * 2)

	// Output:
	// 2
	// 2
}
```

### Option

Go-generics option module inspired by rust.

Avoid `(T, bool)`, handle value with `Option` type, will you choose her?

- Option Example

```go
func ExampleOption() {
	type A struct {
		X int
	}
	var a = gust.Some(A{X: 1})
	fmt.Println(a.IsSome(), a.IsNone())

	var b = gust.None[A]()
	fmt.Println(b.IsSome(), b.IsNone())

	var x = b.UnwrapOr(A{X: 2})
	fmt.Println(x)

	type B struct {
		Y string
	}
	var c = opt.Map(a, func(t A) B {
		return B{
			Y: strconv.Itoa(t.X),
		}
	})
	fmt.Println(c)

	// Output:
	// true false
	// false true
	// {2}
	// Some({1})
}
```

### OptNil

Go-generics option module inspired by rust.

Avoid `(*T)(nil)`, handle value with `OptNil` type, will you choose her?

- OptNil Example

```go
func ExampleOptNil() {
	type A struct {
		X int
	}
	var a = gust.Ptr(&A{X: 1})
	fmt.Println(a.NotNil(), a.IsNil())

	var b = gust.Nil[A]()
	fmt.Println(b.NotNil(), b.IsNil())

	var x = b.UnwrapOr(&A{X: 2})
	fmt.Println(x)

	type B struct {
		Y string
	}
	var c = optnil.Map(a, func(t *A) *B {
		return &B{
			Y: strconv.Itoa(t.X),
		}
	})
	fmt.Println(c)

	// Output:
	// true false
	// false true
	// &{2}
	// NonNil(&{1})
}
```
