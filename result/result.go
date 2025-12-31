// Package result provides helper functions for working with Result types.
//
// Result represents a value that can be either Ok(T) or Err(error).
// It provides a Rust-inspired way to handle errors in Go without using
// the traditional (T, error) tuple pattern, enabling chainable error handling.
//
// # Examples
//
//	// Create a Result
//	success := result.Ok(42)
//	failure := result.TryErr[int](errors.New("error"))
//
//	// Chain operations
//	value := success.
//		Map(func(x int) int { return x * 2 }).
//		UnwrapOr(0) // Output: 84
//
//	// Handle errors
//	err := failure.
//		MapErr(func(e error) error { return fmt.Errorf("wrapped: %w", e) }).
//		Err() // Returns the error
package result

import (
	"github.com/andeya/gust/internal/core"
)

// Result is an alias for core.Result[T].
// This allows using result.Result[T] instead of core.Result[T].
//
// Result represents a value that can be either Ok(T) or Err(error).
// It provides a Rust-inspired way to handle errors in Go without using
// the traditional (T, error) tuple pattern.
type Result[T any] = core.Result[T]

// VoidResult is an alias for core.VoidResult.
// VoidResult is a type alias for Result[Void], used for operations that
// only return success or failure without a value.
type VoidResult = core.VoidResult

// Ok wraps a successful result.
//
//go:inline
func Ok[T any](ok T) Result[T] {
	return core.Ok(ok)
}

// OkVoid returns Ok[Void](nil).
//
// Example:
//
//	```go
//	var result result.VoidResult = result.OkVoid()
//	```
//
//go:inline
func OkVoid() VoidResult {
	return core.OkVoid()
}

// TryErr wraps a failure result.
// NOTE: If err is nil, TryErr(nil) returns Ok with zero value.
// This follows the principle that nil represents "no error", so TryErr(nil) should be Ok.
//
// Example:
//
//	```go
//	result := result.TryErr[string](nil)  // This is an ok state
//	if result.IsOk() {
//		fmt.Println("This will be printed")
//	}
//	```
//
//go:inline
func TryErr[T any](err any) Result[T] {
	return core.TryErr[T](err)
}

// TryErrVoid wraps a failure result as VoidResult.
// NOTE: If err is nil, TryErrVoid(nil) returns OkVoid().
//
//go:inline
func TryErrVoid(err any) VoidResult {
	return core.TryErrVoid(err)
}

// FmtErr wraps a failure result with a formatted error.
//
//go:inline
func FmtErr[T any](format string, args ...any) Result[T] {
	return core.FmtErr[T](format, args...)
}

// FmtErrVoid wraps a failure result with a formatted error as VoidResult.
//
// Example:
//
//	```go
//	var res result.VoidResult = result.FmtErrVoid("operation failed: %s", "file not found")
//	if res.IsErr() {
//		fmt.Println(res.Err())
//	}
//	```
//
//go:inline
func FmtErrVoid(format string, args ...any) VoidResult {
	return core.FmtErrVoid(format, args...)
}

// AssertRet returns the Result[T] of asserting `i` to type `T`
//
//go:inline
func AssertRet[T any](i any) Result[T] {
	return core.AssertRet[T](i)
}

// Ret wraps a result.
//
//go:inline
func Ret[T any](v T, err error) Result[T] {
	return core.Ret(v, err)
}

// RetVoid wraps an error as VoidResult (Result[Void]).
// Returns Ok[Void](nil) if maybeError is nil, otherwise returns Err[Void](maybeError).
//
// Example:
//
//	```go
//	var result result.VoidResult = result.RetVoid(maybeError)
//	```
//
//go:inline
func RetVoid(err error) VoidResult {
	return core.RetVoid(err)
}

// ToError converts VoidResult to a standard Go error.
// Returns nil if IsOk() is true, otherwise returns the error.
//
// Example:
//
//	```go
//	var result result.VoidResult = result.RetVoid(err)
//	if err := result.ToError(result); err != nil {
//		return err
//	}
//	```
//
//go:inline
func ToError(r VoidResult) error {
	return core.ToError(r)
}

// UnwrapErrOr returns the contained error value or a provided default for VoidResult.
//
// Example:
//
//	```go
//	var result result.VoidResult = result.RetVoid(err)
//	err := result.UnwrapErrOr(result, errors.New("default error"))
//	```
//
//go:inline
func UnwrapErrOr(r VoidResult, defaultErr error) error {
	return core.UnwrapErrOr(r, defaultErr)
}

// Assert asserts Result[T] as Result[U].
//
//go:inline
func Assert[T any, U any](o Result[T]) Result[U] {
	if o.IsOk() {
		u, ok := any(o.Unwrap()).(U)
		if ok {
			return Ok[U](u)
		}
		return FmtErr[U]("type assert error, got %T, want %T", o.Unwrap(), u)
	}
	return TryErr[U](o.UnwrapErr())
}

// Assert2 asserts a value and an error as a Result[U].
//
// # Examples
//
//	var v = 1
//	var err = nil
//	var result = Assert2(v, err)
//	assert.Equal(t, core.Ok[int](1), result)
//
//go:inline
func Assert2[T any, U any](v T, err error) Result[U] {
	if err != nil {
		return TryErr[U](err)
	}
	u, ok := any(v).(U)
	if ok {
		return Ok[U](u)
	}
	return FmtErr[U]("type assert error, got %T, want %T", v, u)
}

// XAssert asserts Result[any] as Result[U].
//
//go:inline
func XAssert[U any](o Result[any]) Result[U] {
	if o.IsOk() {
		u, ok := o.Unwrap().(U)
		if ok {
			return Ok[U](u)
		}
		return FmtErr[U]("type assert error, got %T, want %T", o.Unwrap(), u)
	}
	return TryErr[U](o.UnwrapErr())
}

// XAssert2 asserts a value and an error as a Result[U].
//
// # Examples
//
//	var v = 1
//	var err = nil
//	var result = XAssert2(v, err)
//	assert.Equal(t, core.Ok[int](1), result)
//
//go:inline
func XAssert2[U any](v any, err error) Result[U] {
	if err != nil {
		return core.TryErr[U](err)
	}
	u, ok := v.(U)
	if ok {
		return core.Ok[U](u)
	}
	return core.FmtErr[U]("type assert error, got %T, want %T", v, u)
}

// Map maps a Result[T] to Result[U] by applying a function to a contained Ok value, leaving an error untouched.
// This function can be used to compose the results of two functions.
//
//go:inline
func Map[T any, U any](r Result[T], f func(T) U) Result[U] {
	if r.IsOk() {
		return Ok[U](f(r.Unwrap()))
	}
	return TryErr[U](r.Err())
}

// Map2 maps a value and an error as a Result[U] by applying a function to the value.
//
// # Examples
//
//	var v = 1
//	var err = nil
//	var result = Map2(v, err, func(v int) int { return v * 2 })
//	assert.Equal(t, core.Ok[int](2), result)
//
//go:inline
func Map2[T any, U any](v T, err error, f func(T) U) Result[U] {
	if err != nil {
		return TryErr[U](err)
	}
	return Ok[U](f(v))
}

// MapOr returns the provided default (if error), or applies a function to the contained value (if no error),
// Arguments passed to map_or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use MapOrElse, which is lazily evaluated.
//
//go:inline
func MapOr[T any, U any](r Result[T], defaultOk U, f func(T) U) U {
	if r.IsOk() {
		return f(r.Unwrap())
	}
	return defaultOk
}

// MapOr2 maps a value and an error to U by applying a function to the value, or returns the default if error.
//
// # Examples
//
//	var v = 1
//	var err = nil
//	var result = MapOr2(v, err, func(v int) int { return v * 2 })
//	assert.Equal(t, core.Ok[int](2), result)
//
//go:inline
func MapOr2[T any, U any](v T, err error, defaultOk U, f func(T) U) U {
	if err != nil {
		return defaultOk
	}
	return f(v)
}

// MapOrElse maps a Result[T] to U by applying fallback function default to a contained error, or function f to a contained Ok value.
// This function can be used to unpack a successful result while handling an error.
//
//go:inline
func MapOrElse[T any, U any](r Result[T], defaultFn func(error) U, f func(T) U) U {
	if r.IsOk() {
		return f(r.Unwrap())
	}
	return defaultFn(r.Err())
}

// MapOrElse2 maps a value and an error to U by applying a function to the value, or applies the default function if error.
//
// # Examples
//
//	var v = 1
//	var err = nil
//	var result = MapOrElse2(v, err, func(err error) int { return 0 }, func(v int) int { return v * 2 })
//	assert.Equal(t, core.Ok[int](2), result)
//
//go:inline
func MapOrElse2[T any, U any](v T, err error, defaultFn func(error) U, f func(T) U) U {
	if err != nil {
		return defaultFn(err)
	}
	return f(v)
}

// And returns `r1` if `r1` is Err, otherwise returns `r2`.
//
//go:inline
func And[T any, U any](r1 Result[T], r2 Result[U]) Result[U] {
	if r1.IsErr() {
		return TryErr[U](r1.Err())
	}
	return r2
}

// And2 returns `Ret(v1, err1)` if `r1` is `Err`, otherwise returns `Ret(v2, err2)`.
//
// # Examples
//
//	var v1 = 1
//	var err1 = nil
//	var v2 = 2
//	var err2 = nil
//	var result = And2(v1, err1, v2, err2)
//	assert.Equal(t, core.Ok[int](2), result)
//
//	var v1 = 1
//	var err1 = errors.New("error1")
//	var v2 = 2
//	var err2 = nil
//	var result = And2(v1, err1, v2, err2)
//	assert.Equal(t, "error1", result.Err().Error())
//
//	var v1 = 1
//	var err1 = nil
//	var v2 = 2
//	var err2 = errors.New("error2")
//	var result = And2(v1, err1, v2, err2)
//	assert.Equal(t, "error2", result.Err().Error())
//
//	var v1 = 1
//	var err1 = errors.New("error1")
//	var v2 = 2
//	var err2 = errors.New("error2")
//	var result = And2(v1, err1, v2, err2)
//	assert.Equal(t, "error1", result.Err().Error())
//
//go:inline
func And2[T any, U any](v1 T, err1 error, v2 U, err2 error) Result[U] {
	if err1 != nil {
		return TryErr[U](err1)
	}
	return Ret(v2, err2)
}

// AndThen calls op if the result is Ok, otherwise returns the error of self.
// This function can be used for control flow based on Result values.
//
//go:inline
func AndThen[T any, U any](r Result[T], op func(T) Result[U]) Result[U] {
	if r.IsErr() {
		return TryErr[U](r.Err())
	}
	return op(r.Unwrap())
}

// AndThen2 calls op if the result is Ok, otherwise returns the error of self.
// This function can be used for control flow based on Result values.
//
//go:inline
func AndThen2[T any, U any](r Result[T], op func(T) (U, error)) Result[U] {
	if r.IsErr() {
		return TryErr[U](r.Err())
	}
	return Ret[U](op(r.Unwrap()))
}

// AndThen3 calls op if the result is Ok, otherwise returns the error of self.
// This function can be used for control flow based on Result values.
//
//go:inline
func AndThen3[T any, U any](v T, err error, op func(T) (U, error)) Result[U] {
	if err != nil {
		return TryErr[U](err)
	}
	return Ret[U](op(v))
}

// Contains returns true if the result is an Ok value containing the given value.
//
//go:inline
func Contains[T comparable](r Result[T], x T) bool {
	if r.IsErr() {
		return false
	}
	return r.Unwrap() == x
}

// Contains2 returns true if the result is an Ok value containing the given value.
//
//go:inline
func Contains2[T comparable](v T, err error, x T) bool {
	if err != nil {
		return false
	}
	return v == x
}

// Flatten converts from Result[Result[T]] to Result[T].
//
// # Examples
//
//	var r1 = core.Ok(core.Ok(1))
//	var result1 = Flatten(r1)
//	assert.Equal(t, core.Ok[int](1), result1)
//	var r2 = core.Ok(core.TryErr[int](errors.New("error")))
//	var result2 = Flatten(r2)
//	assert.Equal(t, "error", result2.Err().Error())
//	var r3 = core.TryErr[Result[int]](errors.New("error"))
//	var result3 = Flatten(r3)
//	assert.Equal(t, "error", result3.Err().Error())
//
//go:inline
func Flatten[T any](r Result[Result[T]]) Result[T] {
	if r.IsErr() {
		return TryErr[T](r.Err())
	}
	return r.Unwrap()
}

// Flatten2 converts from `(Result[T], error)` to Result[T].
//
// # Examples
//
//	var r1 = core.Ok(1)
//	var err1 = nil
//	var result1 = Flatten2(r1, err1)
//	assert.Equal(t, core.Ok[int](1), result1)
//	var r2 = core.Ok(1)
//	var err2 = errors.New("error")
//	var result2 = Flatten2(r2, err2)
//	assert.Equal(t, "error", result2.Err().Error())
//	var r3 = core.TryErr[int](errors.New("error"))
//	var err3 = nil
//	var result3 = Flatten2(r3, err3)
//	assert.Equal(t, "error", result3.Err().Error())
//
//go:inline
func Flatten2[T any](r Result[T], err error) Result[T] {
	if err != nil {
		return TryErr[T](err)
	}
	return r
}
