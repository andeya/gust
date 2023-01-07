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
	return Result[T]{inner: EnumErr[T, error](toError(err))}
}

// Result can be used to improve `func()(T,error)`,
// represents either success (T) or failure (error).
type Result[T any] struct {
	inner EnumResult[T, error]
}

// IsValid returns true if the object is initialized.
func (r *Result[T]) IsValid() bool {
	return r != nil && r.inner.IsValid()
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

// Split returns the tuple (T, error).
func (r Result[T]) Split() (T, error) {
	return r.inner.Split()
}

// EnumResult returns the inner EnumResult[T, error].
func (r Result[T]) EnumResult() EnumResult[T, error] {
	return r.inner
}

// Errable converts from `Result[T]` to `Errable[error]`.
func (r Result[T]) Errable() Errable[error] {
	if r.IsErr() {
		return ToErrable[error](r.inner.safeGetE())
	}
	return NonErrable[error]()
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

// XOk converts from `Result[T]` to `Option[any]`.
func (r Result[T]) XOk() Option[any] {
	return r.inner.XOk()
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
	if r.IsOk() {
		return nil
	}
	e := r.inner.safeGetE()
	if ev, ok := any(e).(*ErrBox); ok {
		if ev != nil {
			return ev.val
		}
		return nil
	}
	return e
}

// ToX converts from `Result[T]` to Result[any].
func (r Result[T]) ToX() Result[any] {
	return Result[any]{inner: r.inner.ToXOk()}
}

// Map maps a Result[T] to Result[T] by applying a function to a contained Ok value, leaving an error untouched.
// This function can be used to compose the results of two functions.
func (r Result[T]) Map(f func(T) T) Result[T] {
	if r.IsOk() {
		return Ok[T](f(r.inner.safeGetT()))
	}
	return Err[T](r.inner.safeGetE())
}

// XMap maps a Result[T] to Result[any] by applying a function to a contained Ok value, leaving an error untouched.
// This function can be used to compose the results of two functions.
func (r Result[T]) XMap(f func(T) any) Result[any] {
	if r.IsOk() {
		return Ok[any](f(r.inner.safeGetT()))
	}
	return Err[any](r.inner.safeGetE())
}

// MapOr returns the provided default (if error), or applies a function to the contained value (if no error),
// Arguments passed to map_or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use MapOrElse, which is lazily evaluated.
func (r Result[T]) MapOr(defaultOk T, f func(T) T) T {
	if r.IsOk() {
		return f(r.inner.safeGetT())
	}
	return defaultOk
}

// XMapOr returns the provided default (if error), or applies a function to the contained value (if no error),
// Arguments passed to map_or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use MapOrElse, which is lazily evaluated.
func (r Result[T]) XMapOr(defaultOk any, f func(T) any) any {
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

// XMapOrElse maps a Result[T] to `any` by applying fallback function default to a contained error, or function f to a contained Ok value.
// This function can be used to unpack a successful result while handling an error.
func (r Result[T]) XMapOrElse(defaultFn func(error) any, f func(T) any) any {
	if r.IsOk() {
		return f(r.inner.safeGetT())
	}
	return defaultFn(r.inner.safeGetE())
}

// MapErr maps a Result[T] to Result[T] by applying a function to a contained error, leaving an Ok value untouched.
// This function can be used to pass through a successful result while handling an error.
func (r Result[T]) MapErr(op func(error) (newErr any)) Result[T] {
	if r.IsErr() {
		return Err[T](op(r.inner.safeGetE()))
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

// UnwrapOrDefault returns the contained T or a non-nil-pointer zero T.
func (r Result[T]) UnwrapOrDefault() T {
	return r.inner.UnwrapOrDefault()
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
	if r.IsErr() {
		return r.inner.safeGetE()
	}
	panic(ToErrBox(fmt.Sprintf("called `Result.UnwrapErr()` on an `ok` value: %v", r.inner.safeGetT())))
}

// And returns res if the result is Ok, otherwise returns the error of self.
func (r Result[T]) And(res Result[T]) Result[T] {
	if r.IsErr() {
		return r
	}
	return res
}

// XAnd returns res if the result is Ok, otherwise returns the error of self.
func (r Result[T]) XAnd(res Result[any]) Result[any] {
	if r.IsErr() {
		return Err[any](r.inner.safeGetE())
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

// XAndThen calls op if the result is Ok, otherwise returns the error of self.
// This function can be used for control flow based on Result values.
func (r Result[T]) XAndThen(op func(T) Result[any]) Result[any] {
	if r.IsErr() {
		return Err[any](r.Err())
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
func (r Result[T]) ContainsErr(err any) bool {
	if r.IsOk() {
		return false
	}
	return errors.Is(r.inner.safeGetE(), toError(err))
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
	return r.inner.MarshalJSON()
}

func (r *Result[T]) UnmarshalJSON(b []byte) error {
	return r.inner.UnmarshalJSON(b)
}

var (
	_ Iterable[any]   = Result[any]{}
	_ DeIterable[any] = Result[any]{}
)

func (r Result[T]) Next() Option[T] {
	return r.inner.Next()
}

func (r Result[T]) NextBack() Option[T] {
	return r.inner.NextBack()
}

func (r Result[T]) Remaining() uint {
	return r.inner.Remaining()
}

// CtrlFlow returns the `CtrlFlow[error, T]`.
func (r Result[T]) CtrlFlow() CtrlFlow[error, T] {
	if r.IsErr() {
		return Break[error, T](r.inner.safeGetE())
	}
	return Continue[error, T](r.inner.safeGetT())
}

// UnwrapOrThrow returns the contained T or panic returns error (panicValue[*any]).
// NOTE:
//
//	If there is an error, that panic should be caught with CatchResult[U]
func (r Result[T]) UnwrapOrThrow() T {
	if r.inner.isErr {
		if r.inner.value == nil {
			var err error
			v := any(toError(err))
			panic(panicValue[*any]{&v})
		}
		panic(panicValue[*any]{r.inner.value})
	}
	return r.inner.safeGetT()
}

// CatchResult catches panic caused by Result[T].UnwrapOrThrow() and sets error to *Result[U]
func CatchResult[U any](result *Result[U]) {
	switch p := recover().(type) {
	case nil:
	case panicValue[*any]:
		result.inner.value = p.value
		result.inner.isErr = true
	default:
		panic(p)
	}
}
