package gust

import (
	"errors"
	"fmt"
)

// Ret wraps a result.
func Ret[T any](some T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(some)
}

// Ok wraps an Ok result.
func Ok[T any](ok T) Result[T] {
	return Result[T]{ok: ok}
}

// Err wraps an error result.
func Err[T any](err any) Result[T] {
	return Result[T]{err: newAnyError(err)}
}

// Result is a type that represents either success (T) or failure (error).
type Result[T any] struct {
	ok  T
	err error
}

// String returns the string representation.
func (r Result[T]) String() string {
	if r.IsErr() {
		return fmt.Sprintf("Err(%s)", r.err.Error())
	}
	return fmt.Sprintf("Ok(%v)", r.ok)
}

// IsOk returns true if the result is Ok.
func (r Result[T]) IsOk() bool {
	return !r.IsErr()
}

// IsOkAnd returns true if the result is Ok and the value inside of it matches a predicate.
func (r Result[T]) IsOkAnd(f func(T) bool) bool {
	if r.IsOk() {
		return f(r.ok)
	}
	return false
}

// IsErr returns true if the result is error.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// IsErrAnd returns true if the result is error and the value inside of it matches a predicate.
func (r Result[T]) IsErrAnd(f func(error) bool) bool {
	if r.IsErr() {
		return f(r.err)
	}
	return false
}

// Ok converts from `Result[T]` to `Option[T]`.
func (r Result[T]) Ok() Option[T] {
	if r.IsOk() {
		return Some(r.ok)
	}
	return None[T]()
}

// Err returns error.
func (r Result[T]) Err() error {
	return r.err
}

// ErrVal returns error inner value.
func (r Result[T]) ErrVal() any {
	if r.IsErr() {
		return nil
	}
	if ev, _ := r.err.(*errorWithVal); ev != nil {
		return ev.val
	}
	return r.err
}

// Map maps a Result[T] to Result[T] by applying a function to a contained Ok value, leaving an error untouched.
// This function can be used to compose the results of two functions.
func (r Result[T]) Map(f func(T) T) Result[T] {
	if r.IsOk() {
		return Ok[T](f(r.ok))
	}
	return Err[T](r.err)
}

// MapOr returns the provided default (if error), or applies a function to the contained value (if no error),
// Arguments passed to map_or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use MapOrElse, which is lazily evaluated.
func (r Result[T]) MapOr(defaultOk T, f func(T) T) T {
	if r.IsOk() {
		return f(r.ok)
	}
	return defaultOk
}

// MapOrElse maps a Result[T] to T by applying fallback function default to a contained error, or function f to a contained Ok value.
// This function can be used to unpack a successful result while handling an error.
func (r Result[T]) MapOrElse(defaultFn func(error) T, f func(T) T) T {
	if r.IsOk() {
		return f(r.ok)
	}
	return defaultFn(r.err)
}

// MapErr maps a Result[T] to Result[T] by applying a function to a contained error, leaving an Ok value untouched.
// This function can be used to pass through a successful result while handling an error.
func (r *Result[T]) MapErr(op func(error) error) Result[T] {
	if r.IsErr() {
		r.err = op(r.err)
	}
	return *r
}

// Inspect calls the provided closure with a reference to the contained value (if no error).
func (r Result[T]) Inspect(f func(T)) Result[T] {
	if r.IsOk() {
		f(r.ok)
	}
	return r
}

// InspectErr calls the provided closure with a reference to the contained error (if error).
func (r Result[T]) InspectErr(f func(error)) Result[T] {
	if r.IsErr() {
		f(r.err)
	}
	return r
}

// Expect returns the contained Ok value.
// Panics if the value is an error, with a panic message including the
// passed message, and the content of the error.
func (r Result[T]) Expect(msg string) T {
	if r.IsErr() {
		panic(fmt.Errorf("%s: %w", msg, r.err))
	}
	return r.ok
}

// Unwrap returns the contained Ok value.
// Because this function may panic, its use is generally discouraged.
// Instead, prefer to use pattern matching and handle the error case explicitly, or call UnwrapOr or UnwrapOrElse.
func (r Result[T]) Unwrap() T {
	if r.IsErr() {
		panic(fmt.Errorf("called `Result.Unwrap()` on an `err` value: %w", r.err))
	}
	return r.ok
}

// ExpectErr returns the contained error.
// Panics if the value is not an error, with a panic message including the
// passed message, and the content of the [`Ok`].
func (r Result[T]) ExpectErr(msg string) error {
	if r.IsErr() {
		return r.err
	}
	panic(fmt.Errorf("%s: %v", msg, r.ok))
}

// UnwrapErr returns the contained error.
// Panics if the value is not an error, with a custom panic message provided
// by the [`Ok`]'s value.
func (r Result[T]) UnwrapErr() error {
	if r.IsErr() {
		return r.err
	}
	panic(fmt.Errorf("called `Result.UnwrapErr()` on an `ok` value: %v", r.ok))
}

// And returns res if the result is Ok, otherwise returns the error of self.
func (r Result[T]) And(res Result[T]) Result[T] {
	if r.IsErr() {
		return r
	}
	return res
}

// AndThen calls op if the result is Ok, otherwise returns the error of self.
// This function can be used for control flow based on Result values.
func (r Result[T]) AndThen(op func(T) Result[T]) Result[T] {
	if r.IsErr() {
		return r
	}
	return op(r.ok)
}

// Or returns res if the result is Err, otherwise returns the Ok value of r.
// Arguments passed to or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use OrElse, which is lazily evaluated.
func (r Result[T]) Or(res Result[T]) Result[T] {
	if r.IsErr() {
		return res
	}
	return r
}

// OrElse calls op if the result is Err, otherwise returns the Ok value of self.
// This function can be used for control flow based on result values.
func (r Result[T]) OrElse(op func(error) Result[T]) Result[T] {
	if r.IsErr() {
		return op(r.err)
	}
	return r
}

// UnwrapOr returns the contained Ok value or a provided default.
// Arguments passed to UnwrapOr are eagerly evaluated; if you are passing the result of a function call, it is recommended to use UnwrapOrElse, which is lazily evaluated.
func (r Result[T]) UnwrapOr(defaultOk T) T {
	if r.IsErr() {
		return defaultOk
	}
	return r.ok
}

// UnwrapOrElse returns the contained Ok value or computes it from a closure.
func (r Result[T]) UnwrapOrElse(defaultFn func(error) T) T {
	if r.IsErr() {
		return defaultFn(r.err)
	}
	return r.ok
}

// UnwrapUnchecked returns the contained Ok value, without checking that the value is not an error.
func (r Result[T]) UnwrapUnchecked() T {
	return r.ok
}

// ContainsErr returns true if the result is an error containing the given value.
func (r Result[T]) ContainsErr(err error) bool {
	if r.IsOk() {
		return false
	}
	if err == r.err {
		return true
	}
	return errors.Is(r.err, err)
}
