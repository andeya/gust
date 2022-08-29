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

// Ok wraps a successful result.
func Ok[T any](ok T) Result[T] {
	return Result[T]{inner: EnumOk[T, error](ok)}
}

// Err wraps a failure result.
func Err[T any](err any) Result[T] {
	return Result[T]{inner: EnumErr[T, error](newAnyError(err))}
}

// Result can be used to improve `func()(T,error)`,
// represents either success (T) or failure (error).
type Result[T any] struct {
	inner EnumResult[T, error]
}

// Ref returns pointer.
func (r Result[T]) Ref() *Result[T] {
	return &r
}

// IsErr returns true if the result is error.
func (r Result[T]) IsErr() bool {
	return r.inner.IsErr()
}

// IsOk returns true if the result is Ok.
func (r Result[T]) IsOk() bool {
	return !r.IsErr()
}

// String returns the string representation.
func (r Result[T]) String() string {
	return r.inner.String()
}

// IsOkAnd returns true if the result is Ok and the value inside it matches a predicate.
func (r Result[T]) IsOkAnd(f func(T) bool) bool {
	return r.inner.IsOkAnd(f)
}

// IsErrAnd returns true if the result is error and the value inside it matches a predicate.
func (r Result[T]) IsErrAnd(f func(error) bool) bool {
	return r.inner.IsErrAnd(f)
}

// Ok converts from `Result[T]` to `Option[T]`.
func (r Result[T]) Ok() Option[T] {
	return r.inner.Ok()
}

// Err returns error.
func (r Result[T]) Err() error {
	if r.IsErr() {
		return r.inner.safeGetE()
	}
	return nil
}

// ErrVal returns error inner value.
func (r Result[T]) ErrVal() any {
	if r.IsErr() {
		return nil
	}
	if ev, _ := r.inner.val.(*errorWithVal); ev != nil {
		return ev.val
	}
	return r.inner.val
}

// Map maps a Result[T] to Result[T] by applying a function to a contained Ok value, leaving an error untouched.
// This function can be used to compose the results of two functions.
func (r Result[T]) Map(f func(T) T) Result[T] {
	if r.IsOk() {
		return Ok[T](f(r.inner.safeGetT()))
	}
	return Err[T](r.inner.val)
}

// MapOr returns the provided default (if error), or applies a function to the contained value (if no error),
// Arguments passed to map_or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use MapOrElse, which is lazily evaluated.
func (r Result[T]) MapOr(defaultOk T, f func(T) T) T {
	if r.IsOk() {
		return f(r.inner.safeGetT())
	}
	return defaultOk
}

// MapOrElse maps a Result[T] to T by applying fallback function default to a contained error, or function f to a contained Ok value.
// This function can be used to unpack a successful result while handling an error.
func (r Result[T]) MapOrElse(defaultFn func(error) T, f func(T) T) T {
	if r.IsOk() {
		return f(r.inner.safeGetT())
	}
	return defaultFn(r.inner.safeGetE())
}

// MapErr maps a Result[T] to Result[T] by applying a function to a contained error, leaving an Ok value untouched.
// This function can be used to pass through a successful result while handling an error.
func (r Result[T]) MapErr(op func(error) error) Result[T] {
	if r.IsErr() {
		r.inner.val = op(r.inner.safeGetE())
	}
	return r
}

// Inspect calls the provided closure with a reference to the contained value (if no error).
func (r Result[T]) Inspect(f func(T)) Result[T] {
	if r.IsOk() {
		f(r.inner.safeGetT())
	}
	return r
}

// InspectErr calls the provided closure with a reference to the contained error (if error).
func (r Result[T]) InspectErr(f func(error)) Result[T] {
	if r.IsErr() {
		f(r.inner.safeGetE())
	}
	return r
}

// Expect returns the contained Ok value.
// Panics if the value is an error, with a panic message including the
// passed message, and the content of the error.
func (r Result[T]) Expect(msg string) T {
	return r.inner.Expect(msg)
}

// Unwrap returns the contained Ok value.
// Because this function may panic, its use is generally discouraged.
// Instead, prefer to use pattern matching and handle the error case explicitly, or call UnwrapOr or UnwrapOrElse.
func (r Result[T]) Unwrap() T {
	return r.inner.Unwrap()
}

// ExpectErr returns the contained error.
// Panics if the value is not an error, with a panic message including the
// passed message, and the content of the [`Ok`].
func (r Result[T]) ExpectErr(msg string) error {
	return r.inner.ExpectErr(msg)
}

// UnwrapErr returns the contained error.
// Panics if the value is not an error, with a custom panic message provided
// by the [`Ok`]'s value.
func (r Result[T]) UnwrapErr() error {
	return r.inner.UnwrapErr()
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
	return op(r.inner.safeGetT())
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
		return op(r.inner.safeGetE())
	}
	return r
}

// UnwrapOr returns the contained Ok value or a provided default.
// Arguments passed to UnwrapOr are eagerly evaluated; if you are passing the result of a function call, it is recommended to use UnwrapOrElse, which is lazily evaluated.
func (r Result[T]) UnwrapOr(defaultOk T) T {
	return r.inner.UnwrapOr(defaultOk)
}

// UnwrapOrElse returns the contained Ok value or computes it from a closure.
func (r Result[T]) UnwrapOrElse(defaultFn func(error) T) T {
	return r.inner.UnwrapOrElse(defaultFn)
}

// ContainsErr returns true if the result is an error containing the given value.
func (r Result[T]) ContainsErr(err error) bool {
	if r.IsOk() {
		return false
	}
	return errors.Is(r.inner.safeGetE(), err)
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	return r.inner.MarshalJSON()
}

func (r *Result[T]) UnmarshalJSON(b []byte) error {
	return r.inner.UnmarshalJSON(b)
}

type errorWithVal struct {
	val any
}

func newAnyError(val any) error {
	if err, _ := val.(error); err != nil {
		return err
	}
	return &errorWithVal{val: val}
}

func (a *errorWithVal) Error() string {
	return fmt.Sprintf("%v", a.val)
}

var (
	_ DataForIter[any]            = (*Result[any])(nil)
	_ DataForDoubleEndedIter[any] = (*Result[any])(nil)
)

func (r *Result[T]) NextForIter() Option[T] {
	return r.inner.NextForIter()
}

func (r *Result[T]) NextBackForIter() Option[T] {
	return r.inner.NextBackForIter()
}

func (r Result[T]) RemainingLenForIter() uint {
	return r.inner.RemainingLenForIter()
}
