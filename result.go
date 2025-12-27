package gust

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type (
	// Result can be used to improve `func()(T,error)`,
	// represents either success (T) or failure (error).
	Result[T any] struct {
		t Option[T]
		e ErrBox
	}

	// VoidResult is an alias for Result[Void], representing a result that only indicates success or failure.
	// This is equivalent to Rust's Result<(), E> and provides a simpler API than Result[Void].
	//
	// Example:
	//
	//	```go
	//	var result gust.VoidResult = gust.RetVoid(err)
	//	if result.IsErr() {
	//		fmt.Println(result.Err())
	//	}
	//	```
	VoidResult = Result[Void]
)

// Ret wraps a result.
//
//go:inline
func Ret[T any](some T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(some)
}

// RetVoid wraps an error as VoidResult (Result[Void]).
// Returns Ok[Void](nil) if err is nil, otherwise returns Err[Void](err).
//
// Example:
//
//	```go
//	var result gust.VoidResult = gust.RetVoid(err)
//	```
//
//go:inline
func RetVoid(err any) VoidResult {
	if err == nil {
		return Ok[Void](nil)
	}
	return Err[Void](err)
}

// Ok wraps a successful result.
//
//go:inline
func Ok[T any](ok T) Result[T] {
	return Result[T]{t: Some(ok)}
}

// Err wraps a failure result.
// NOTE: Even if err is nil, Err(nil) still represents an error state.
// This follows declarative programming principles: calling Err() explicitly declares an error result,
// regardless of the error value. This is consistent with Rust's Result::Err(()) semantics.
//
// Example:
//
//	```go
//	result := gust.Err[string](nil)  // This is an error state, even though err is nil
//	if result.IsErr() {
//		fmt.Println("This will be printed")
//	}
//	```
//
//go:inline
func Err[T any](err any) Result[T] {
	eb := BoxErr(err)
	if eb == nil {
		return Result[T]{e: ErrBox{}}
	}
	return Result[T]{e: *eb}
}

// FmtErr wraps a failure result with a formatted error.
//
//go:inline
func FmtErr[T any](format string, a ...any) Result[T] {
	return Err[T](fmt.Errorf(format, a...))
}

// NonResult returns Ok[Void](nil).
//
// Example:
//
//	```go
//	var result gust.VoidResult = gust.NonResult()
//	```
//
//go:inline
func NonResult() VoidResult {
	return Ok[Void](nil)
}

// AssertRet returns the Result[T] of asserting `i` to type `T`
func AssertRet[T any](i any) Result[T] {
	value, ok := i.(T)
	if !ok {
		return FmtErr[T]("type assert error, got %T, want %T", i, value)
	}
	return Ok(value)
}

// Ref returns the pointer of the object.
//
//go:inline
func (r Result[T]) Ref() *Result[T] {
	return &r
}

// safeGetT safely gets the T value.
//
//go:inline
func (r Result[T]) safeGetT() T {
	if r.t.IsSome() {
		return r.t.UnwrapUnchecked()
	}
	var t T
	return t
}

// safeGetE safely gets the error value.
//
//go:inline
func (r Result[T]) safeGetE() error {
	return (&r.e).ToError()
}

// IsValid returns true if the object is initialized.
//
//go:inline
func (r *Result[T]) IsValid() bool {
	return r != nil && (!r.e.IsEmpty() || r.t.IsSome())
}

// IsErr returns true if the result is error.
// NOTE: This is determined by whether t.IsSome() is false, not by e.IsEmpty().
// This ensures that Err(nil) is correctly identified as an error state,
// following declarative programming principles where Err() explicitly declares an error result.
//
//go:inline
func (r Result[T]) IsErr() bool {
	return !r.t.IsSome()
}

// IsOk returns true if the result is Ok.
// NOTE: This is determined by whether t.IsSome() is true, not by e.IsEmpty().
// This ensures that Err(nil) is correctly identified as an error state,
// following declarative programming principles where Err() explicitly declares an error result.
//
//go:inline
func (r Result[T]) IsOk() bool {
	return r.t.IsSome()
}

// String returns the string representation.
func (r Result[T]) String() string {
	if r.IsErr() {
		return fmt.Sprintf("Err(%v)", r.safeGetE())
	}
	return fmt.Sprintf("Ok(%v)", r.safeGetT())
}

// Split returns the tuple (T, error).
//
//go:inline
func (r Result[T]) Split() (T, error) {
	return r.safeGetT(), r.safeGetE()
}

// IsOkAnd returns true if the result is Ok and the value inside it matches a predicate.
//
//go:inline
func (r Result[T]) IsOkAnd(f func(T) bool) bool {
	if r.IsOk() {
		return f(r.safeGetT())
	}
	return false
}

// IsErrAnd returns true if the result is error and the value inside it matches a predicate.
//
//go:inline
func (r Result[T]) IsErrAnd(f func(error) bool) bool {
	if r.IsErr() {
		return f((&r.e).ToError())
	}
	return false
}

// Ok converts from `Result[T]` to `Option[T]`.
//
//go:inline
func (r Result[T]) Ok() Option[T] {
	return r.t
}

// XOk converts from `Result[T]` to `Option[any]`.
//
//go:inline
func (r Result[T]) XOk() Option[any] {
	return r.t.ToX()
}

// Err returns error.
//
//go:inline
func (r Result[T]) Err() error {
	return r.safeGetE()
}

// ToError converts VoidResult to a standard Go error.
// Returns nil if IsOk() is true, otherwise returns the error.
//
// Example:
//
//	```go
//	var result gust.VoidResult = gust.RetVoid(err)
//	if err := gust.ToError(result); err != nil {
//		return err
//	}
//	```
//
//go:inline
func ToError(r VoidResult) error {
	return r.Err()
}

// UnwrapErrOr returns the contained error value or a provided default for VoidResult.
//
// Example:
//
//	```go
//	var result gust.VoidResult = gust.RetVoid(err)
//	err := gust.UnwrapErrOr(result, errors.New("default error"))
//	```
//
//go:inline
func UnwrapErrOr(r VoidResult, def error) error {
	if r.IsErr() {
		return r.Err()
	}
	return def
}

// ErrVal returns error inner value.
//
//go:inline
func (r Result[T]) ErrVal() any {
	if r.IsErr() {
		return (&r.e).Value()
	}
	return nil
}

// ToX converts from `Result[T]` to Result[any].
//
//go:inline
func (r Result[T]) ToX() Result[any] {
	if r.IsErr() {
		return Result[any]{t: None[any](), e: r.e}
	}
	return Ok[any](r.safeGetT())
}

// Map maps a Result[T] to Result[T] by applying a function to a contained Ok value, leaving an error untouched.
// This function can be used to compose the results of two functions.
//
//go:inline
func (r Result[T]) Map(f func(T) T) Result[T] {
	if r.IsOk() {
		return Ok[T](f(r.safeGetT()))
	}
	return Result[T]{t: None[T](), e: r.e}
}

// XMap maps a Result[T] to Result[any] by applying a function to a contained Ok value, leaving an error untouched.
// This function can be used to compose the results of two functions.
//
//go:inline
func (r Result[T]) XMap(f func(T) any) Result[any] {
	if r.IsOk() {
		return Ok[any](f(r.safeGetT()))
	}
	return Result[any]{t: None[any](), e: r.e}
}

// MapOr returns the provided default (if error), or applies a function to the contained value (if no error),
// Arguments passed to map_or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use MapOrElse, which is lazily evaluated.
func (r Result[T]) MapOr(defaultOk T, f func(T) T) T {
	if r.IsOk() {
		return f(r.safeGetT())
	}
	return defaultOk
}

// XMapOr returns the provided default (if error), or applies a function to the contained value (if no error),
// Arguments passed to map_or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use MapOrElse, which is lazily evaluated.
func (r Result[T]) XMapOr(defaultOk any, f func(T) any) any {
	if r.IsOk() {
		return f(r.safeGetT())
	}
	return defaultOk
}

// MapOrElse maps a Result[T] to T by applying fallback function default to a contained error, or function f to a contained Ok value.
// This function can be used to unpack a successful result while handling an error.
func (r Result[T]) MapOrElse(defaultFn func(error) T, f func(T) T) T {
	if r.IsOk() {
		return f(r.safeGetT())
	}
	return defaultFn(r.safeGetE())
}

// XMapOrElse maps a Result[T] to `any` by applying fallback function default to a contained error, or function f to a contained Ok value.
// This function can be used to unpack a successful result while handling an error.
func (r Result[T]) XMapOrElse(defaultFn func(error) any, f func(T) any) any {
	if r.IsOk() {
		return f(r.safeGetT())
	}
	return defaultFn(r.safeGetE())
}

// MapErr maps a Result[T] to Result[T] by applying a function to a contained error, leaving an Ok value untouched.
// This function can be used to pass through a successful result while handling an error.
//
//go:inline
func (r Result[T]) MapErr(op func(error) (newErr any)) Result[T] {
	if r.IsErr() {
		return Err[T](op(r.safeGetE()))
	}
	return r
}

// Inspect calls the provided closure with a reference to the contained value (if no error).
//
//go:inline
func (r Result[T]) Inspect(f func(T)) Result[T] {
	if r.IsOk() {
		f(r.safeGetT())
	}
	return r
}

// InspectErr calls the provided closure with a reference to the contained error (if error).
//
//go:inline
func (r Result[T]) InspectErr(f func(error)) Result[T] {
	if r.IsErr() {
		f(r.safeGetE())
	}
	return r
}

// wrapError wraps an error with a message.
func (r Result[T]) wrapError(msg string) error {
	if r.IsErr() {
		val := (&r.e).Value()
		if val == nil {
			// Err(nil) case: provide a descriptive error message
			return BoxErr(fmt.Errorf("%s: gust.Err(nil)", msg)).ToError()
		}
		if err, ok := val.(error); ok {
			return BoxErr(fmt.Errorf("%s: %w", msg, err)).ToError()
		}
		return BoxErr(fmt.Errorf("%s: %v", msg, val)).ToError()
	}
	return BoxErr(msg).ToError()
}

// Expect returns the contained Ok value.
// Panics if the value is an error, with a panic message including the
// passed message, and the content of the error.
func (r Result[T]) Expect(msg string) T {
	if r.IsErr() {
		panic(r.wrapError(msg))
	}
	return r.safeGetT()
}

// Unwrap returns the contained Ok value.
// Because this function may panic, its use is generally discouraged.
// Instead, prefer to use pattern matching and handle the error case explicitly, or call UnwrapOr or UnwrapOrElse.
// NOTE: This panics *ErrBox (not error) to be consistent with UnwrapOrThrow() and allow Catch() to properly handle it.
func (r Result[T]) Unwrap() T {
	if r.IsErr() {
		panic(&r.e)
	}
	return r.safeGetT()
}

// UnwrapOrDefault returns the contained T or a non-nil-pointer zero T.
func (r Result[T]) UnwrapOrDefault() T {
	if r.IsOk() {
		return r.safeGetT()
	}
	return defaultValue[T]()
}

// UnwrapUnchecked returns the contained T.
//
//go:inline
func (r Result[T]) UnwrapUnchecked() T {
	return r.Ok().UnwrapUnchecked()
}

// ExpectErr returns the contained error.
// Panics if the value is not an error, with a panic message including the
// passed message, and the content of the [`Ok`].
//
//go:inline
func (r Result[T]) ExpectErr(msg string) error {
	if r.IsErr() {
		return r.safeGetE()
	}
	panic(BoxErr(fmt.Sprintf("%s: %v", msg, r.safeGetT())))
}

// UnwrapErr returns the contained error.
// Panics if the value is not an error, with a custom panic message provided
// by the [`Ok`]'s value.
//
//go:inline
func (r Result[T]) UnwrapErr() error {
	if r.IsErr() {
		return r.safeGetE()
	}
	panic(BoxErr(fmt.Sprintf("called `Result.UnwrapErr()` on an `ok` value: %v", r.safeGetT())))
}

// And returns `r` if `r` is `Err`, otherwise returns `r2`.
//
//go:inline
func (r Result[T]) And(r2 Result[T]) Result[T] {
	if r.IsErr() {
		return r
	}
	return r2
}

// And2 returns `r` if `r` is `Err`, otherwise returns `Ret(v2, err2)`.
//
//go:inline
func (r Result[T]) And2(v2 T, err2 error) Result[T] {
	if r.IsErr() {
		return r
	}
	return Ret(v2, err2)
}

// XAnd returns res if the result is Ok, otherwise returns the error of self.
//
//go:inline
func (r Result[T]) XAnd(res Result[any]) Result[any] {
	if r.IsErr() {
		return Result[any]{t: None[any](), e: r.e}
	}
	return res
}

// XAnd2 returns `r` if `r` is `Err`, otherwise returns `Ret(v2, err2)`.
//
//go:inline
func (r Result[T]) XAnd2(v2 any, err2 error) Result[any] {
	if r.IsErr() {
		return Err[any](r.Err())
	}
	return Ret[any](v2, err2)
}

// AndThen calls op if the result is Ok, otherwise returns the error of self.
// This function can be used for control flow based on Result values.
func (r Result[T]) AndThen(op func(T) Result[T]) Result[T] {
	if r.IsErr() {
		return r
	}
	return op(r.safeGetT())
}

// AndThen2 calls op if the result is Ok, otherwise returns the error of self.
// This function can be used for control flow based on Result values.
func (r Result[T]) AndThen2(op func(T) (T, error)) Result[T] {
	if r.IsErr() {
		return r
	}
	return Ret[T](op(r.safeGetT()))
}

// XAndThen calls op if the result is Ok, otherwise returns the error of self.
// This function can be used for control flow based on Result values.
func (r Result[T]) XAndThen(op func(T) Result[any]) Result[any] {
	if r.IsErr() {
		return Err[any](r.Err())
	}
	return op(r.safeGetT())
}

// XAndThen2 calls op if the result is Ok, otherwise returns the error of self.
// This function can be used for control flow based on Result values.
func (r Result[T]) XAndThen2(op func(T) (any, error)) Result[any] {
	if r.IsErr() {
		return Err[any](r.Err())
	}
	return Ret[any](op(r.safeGetT()))
}

// Or returns `r2` if `r` is `Err`, otherwise returns `r`.
// Arguments passed to or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use OrElse, which is lazily evaluated.
//
//go:inline
func (r Result[T]) Or(r2 Result[T]) Result[T] {
	if r.IsErr() {
		return r2
	}
	return r
}

// Or2 returns `Ret(v2, err2)` if `r` is `Err`, otherwise returns `r`.
//
//go:inline
func (r Result[T]) Or2(v2 T, err2 error) Result[T] {
	if r.IsErr() {
		return Ret[T](v2, err2)
	}
	return r
}

// OrElse calls op if the result is Err, otherwise returns the Ok value of self.
// This function can be used for control flow based on result values.
func (r Result[T]) OrElse(op func(error) Result[T]) Result[T] {
	if r.IsErr() {
		return op(r.safeGetE())
	}
	return r
}

// OrElse2 calls op if the result is Err, otherwise returns the Ok value of self.
// This function can be used for control flow based on result values.
func (r Result[T]) OrElse2(op func(error) (T, error)) Result[T] {
	if r.IsErr() {
		return Ret[T](op(r.safeGetE()))
	}
	return r
}

// UnwrapOr returns the contained Ok value or a provided default.
// Arguments passed to UnwrapOr are eagerly evaluated; if you are passing the result of a function call, it is recommended to use UnwrapOrElse, which is lazily evaluated.
func (r Result[T]) UnwrapOr(defaultOk T) T {
	if r.IsErr() {
		return defaultOk
	}
	return r.safeGetT()
}

// UnwrapOrElse returns the contained Ok value or computes it from a closure.
func (r Result[T]) UnwrapOrElse(defaultFn func(error) T) T {
	if r.IsErr() {
		return defaultFn(r.safeGetE())
	}
	return r.safeGetT()
}

// ContainsErr returns true if the result is an error containing the given value.
func (r Result[T]) ContainsErr(err any) bool {
	if r.IsOk() {
		return false
	}
	if r.IsErr() {
		return errors.Is((&r.e).ToError(), BoxErr(err).ToError())
	}
	return false
}

// Flatten converts from `(gust.Result[T], error)` to gust.Result[T].
//
// # Examples
//
//	var r1 = gust.Ok(1)
//	var err1 = nil
//	var result1 = r1.Flatten(err1)
//	assert.Equal(t, gust.Ok[int](1), result1)
//
//	var r2 = gust.Ok(1)
//	var err2 = errors.New("error")
//	var result2 = r2.Flatten(err2)
//	assert.Equal(t, "error", result2.Err().Error())
//
// var r3 = gust.Err(errors.New("error"))
// var err3 = nil
// var result3 = r3.Flatten(err3)
// assert.Equal(t, "error", result3.Err().Error())
//
//go:inline
func (r Result[T]) Flatten(err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return r
}

// AsPtr returns its pointer or nil.
//
//go:inline
func (r Result[T]) AsPtr() *T {
	return r.t.AsPtr()
}

// MarshalJSON implements the json.Marshaler interface.
func (r Result[T]) MarshalJSON() ([]byte, error) {
	if r.IsErr() {
		return nil, r.safeGetE()
	}
	return json.Marshal(r.safeGetT())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (r *Result[T]) UnmarshalJSON(b []byte) error {
	var t T
	if r == nil {
		return &json.InvalidUnmarshalError{Type: reflect.TypeOf(t)}
	}
	err := json.Unmarshal(b, &t)
	if err != nil {
		r.t = None[T]()
		eb := BoxErr(err)
		if eb == nil {
			r.e = ErrBox{}
		} else {
			r.e = *eb
		}
	} else {
		r.t = Some(t)
		r.e = ErrBox{}
	}
	return err
}

var (
	_ Iterable[any]            = new(Result[any])
	_ DoubleEndedIterable[any] = new(Result[any])
)

// Next returns the next element of the iterator.
func (r *Result[T]) Next() Option[T] {
	if r == nil {
		return None[T]()
	}
	return r.t.Next()
}

// NextBack returns the next element from the back of the iterator.
//
//go:inline
func (r *Result[T]) NextBack() Option[T] {
	return r.Next()
}

// Remaining returns the number of remaining elements in the iterator.
func (r *Result[T]) Remaining() uint {
	if r == nil {
		return 0
	}
	return r.t.Remaining()
}

// UnwrapOrThrow returns the contained T or panic returns error (*ErrBox).
// NOTE:
//
//	If there is an error, that panic should be caught with `result.Catch()`
func (r Result[T]) UnwrapOrThrow() T {
	if r.IsErr() {
		panic(r.e)
	}
	return r.safeGetT()
}

// Catch catches any panic and converts it to a Result error.
// It catches:
//   - *ErrBox (gust's own error type)
//   - error (regular Go errors, wrapped in ErrBox)
//   - any other type (wrapped in ErrBox)
//
// Example:
//
//	```go
//	func example() (result Result[string]) {
//	   defer result.Catch()
//	   Err[int]("int error").UnwrapOrThrow()
//	   return Ok[string]("ok")
//	}
//	```
func (r *Result[T]) Catch() {
	if r == nil {
		// If receiver is nil, let panic propagate
		return
	}
	switch p := recover().(type) {
	case nil:
		// No panic occurred
	case *ErrBox:
		// Gust's own ErrBox type (pointer)
		if r.t.IsSome() {
			r.t = None[T]()
		}
		if p != nil {
			r.e = *p
		} else {
			r.e = ErrBox{}
		}
	case ErrBox:
		// Gust's own ErrBox type (value)
		if r.t.IsSome() {
			r.t = None[T]()
		}
		r.e = p
	case error:
		// Regular error panic - wrap in ErrBox
		if r.t.IsSome() {
			r.t = None[T]()
		}
		eb := BoxErr(p)
		if eb == nil {
			r.e = ErrBox{}
		} else {
			r.e = *eb
		}
	default:
		// Other types - wrap in ErrBox
		// This ensures all panics are caught and converted to errors
		if r.t.IsSome() {
			r.t = None[T]()
		}
		eb := BoxErr(p)
		if eb == nil {
			r.e = ErrBox{}
		} else {
			r.e = *eb
		}
	}
}
