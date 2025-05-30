package gust

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// EnumOk wraps a successful result enumeration.
func EnumOk[T any, E any](ok T) EnumResult[T, E] {
	return EnumResult[T, E]{t: Some(ok)}
}

// EnumErr wraps a failure result enumeration.
func EnumErr[T any, E any](err E) EnumResult[T, E] {
	return EnumResult[T, E]{e: &err}
}

// EnumResult represents a success (T) or failure (E) enumeration.
type EnumResult[T any, E any] struct {
	t Option[T]
	e *E
}

// Ref returns the pointer of the object.
//
//go:inline
func (r EnumResult[T, E]) Ref() *EnumResult[T, E] {
	return &r
}

// IsValid returns true if the object is initialized.
func (r *EnumResult[T, E]) IsValid() bool {
	return r != nil && (r.e != nil || r.t.IsSome())
}

func (r EnumResult[T, E]) safeGetT() T {
	if r.t.IsSome() {
		return r.t.UnwrapUnchecked()
	}
	var t T
	return t
}

func (r EnumResult[T, E]) safeGetE() E {
	if r.e != nil {
		return *r.e
	}
	var e E
	return e
}

// Split returns the tuple (T, E).
func (r EnumResult[T, E]) Split() (T, E) {
	return r.safeGetT(), r.safeGetE()
}

// IsErr returns true if the result is E.
func (r EnumResult[T, E]) IsErr() bool {
	return r.e != nil
}

// IsOk returns true if the result is ok.
func (r EnumResult[T, E]) IsOk() bool {
	return !r.IsErr()
}

// String returns the string representation.
func (r EnumResult[T, E]) String() string {
	if r.IsErr() {
		return fmt.Sprintf("Err(%v)", r.safeGetE())
	}
	return fmt.Sprintf("Ok(%v)", r.safeGetT())
}

// Result converts from `EnumResult[T,E]` to `Result[T]`.
func (r EnumResult[T, E]) Result() Result[T] {
	if r.IsErr() {
		return Err[T](r.safeGetE())
	}
	return Ok[T](r.safeGetT())
}

// Errable converts from `EnumResult[T,E]` to `Errable[E]`.
func (r EnumResult[T, E]) Errable() Errable[E] {
	if r.IsErr() {
		return ToErrable[E](r.safeGetE())
	}
	return NonErrable[E]()
}

// IsOkAnd returns true if the result is Ok and the value inside it matches a predicate.
func (r EnumResult[T, E]) IsOkAnd(f func(T) bool) bool {
	if r.IsOk() {
		return f(r.safeGetT())
	}
	return false
}

// IsErrAnd returns true if the result is E and the value inside it matches a predicate.
func (r EnumResult[T, E]) IsErrAnd(f func(E) bool) bool {
	if r.IsErr() {
		return f(r.safeGetE())
	}
	return false
}

// Ok converts from `Result[T,E]` to `Option[T]`.
func (r EnumResult[T, E]) Ok() Option[T] {
	return r.t
}

// XOk converts from `Result[T,E]` to `Option[any]`.
func (r EnumResult[T, E]) XOk() Option[any] {
	return r.t.ToX()
}

// Err returns E value `Option[E]`.
func (r EnumResult[T, E]) Err() Option[E] {
	if r.IsErr() {
		return Some(r.safeGetE())
	}
	return None[E]()
}

// XErr returns E value `Option[any]`.
func (r EnumResult[T, E]) XErr() Option[any] {
	if r.IsErr() {
		return Some[any](r.safeGetE())
	}
	return None[any]()
}

// ToXOk converts from `EnumResult[T,E]` to EnumResult[any,E].
// nolint:gosimple
func (r EnumResult[T, E]) ToXOk() EnumResult[any, E] {
	if r.IsErr() {
		return EnumResult[any, E]{
			e: r.e,
		}
	}
	return EnumResult[any, E]{
		t: r.t.ToX(),
	}
}

// ToXErr converts from `EnumResult[T,E]` to EnumResult[T,any].
// nolint:gosimple
func (r EnumResult[T, E]) ToXErr() EnumResult[T, any] {
	if r.IsErr() {
		e := any(*r.e)
		return EnumResult[T, any]{
			e: &e,
		}
	}
	return EnumResult[T, any]{
		t: r.t,
	}
}

// ToX converts from `EnumResult[T,E]` to EnumResult[any,any].
// nolint:gosimple
func (r EnumResult[T, E]) ToX() EnumResult[any, any] {
	if r.IsErr() {
		e := any(*r.e)
		return EnumResult[any, any]{
			e: &e,
		}
	}
	return EnumResult[any, any]{
		t: r.t.ToX(),
	}
}

// Map maps a EnumResult[T,E] to EnumResult[T,E] by applying a function to a contained T value, leaving an E untouched.
// This function can be used to compose the results of two functions.
func (r EnumResult[T, E]) Map(f func(T) T) EnumResult[T, E] {
	if r.IsOk() {
		return EnumOk[T, E](f(r.safeGetT()))
	}
	return EnumErr[T, E](r.safeGetE())
}

// XMap maps a EnumResult[T,E] to EnumResult[any,E] by applying a function to a contained `any` value, leaving an E untouched.
// This function can be used to compose the results of two functions.
func (r EnumResult[T, E]) XMap(f func(T) any) EnumResult[any, E] {
	if r.IsOk() {
		return EnumOk[any, E](f(r.safeGetT()))
	}
	return EnumErr[any, E](r.safeGetE())
}

// MapOr returns the provided default (if E), or applies a function to the contained value (if no E),
// Arguments passed to map_or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use MapOrElse, which is lazily evaluated.
func (r EnumResult[T, E]) MapOr(defaultOk T, f func(T) T) T {
	if r.IsOk() {
		return f(r.safeGetT())
	}
	return defaultOk
}

// XMapOr returns the provided default (if E), or applies a function to the contained value (if no E),
// Arguments passed to map_or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use MapOrElse, which is lazily evaluated.
func (r EnumResult[T, E]) XMapOr(defaultOk any, f func(T) any) any {
	if r.IsOk() {
		return f(r.safeGetT())
	}
	return defaultOk
}

// MapOrElse maps a EnumResult[T,E] to T by applying fallback function default to a contained E, or function f to a contained T value.
// This function can be used to unpack a successful result while handling an E.
func (r EnumResult[T, E]) MapOrElse(defaultFn func(E) T, f func(T) T) T {
	if r.IsOk() {
		return f(r.safeGetT())
	}
	return defaultFn(r.safeGetE())
}

// XMapOrElse maps a EnumResult[T,E] to `any` type by applying fallback function default to a contained E, or function f to a contained T value.
// This function can be used to unpack a successful result while handling an E.
func (r EnumResult[T, E]) XMapOrElse(defaultFn func(E) any, f func(T) any) any {
	if r.IsOk() {
		return f(r.safeGetT())
	}
	return defaultFn(r.safeGetE())
}

// MapErr maps a EnumResult[T,E] to EnumResult[T,E] by applying a function to a contained E, leaving an T value untouched.
// This function can be used to pass through a successful result while handling an error.
func (r EnumResult[T, E]) MapErr(op func(E) E) EnumResult[T, E] {
	if r.IsErr() {
		return EnumErr[T, E](op(r.safeGetE()))
	}
	return r
}

// XMapErr maps a EnumResult[T,E] to EnumResult[T,any] by applying a function to a contained `any`, leaving an T value untouched.
// This function can be used to pass through a successful result while handling an error.
func (r EnumResult[T, E]) XMapErr(op func(E) any) EnumResult[T, any] {
	if r.IsErr() {
		return EnumErr[T, any](op(r.safeGetE()))
	}
	return EnumOk[T, any](r.safeGetT())
}

// Inspect calls the provided closure with a reference to the contained value (if no E).
func (r EnumResult[T, E]) Inspect(f func(T)) EnumResult[T, E] {
	if r.IsOk() {
		f(r.safeGetT())
	}
	return r
}

// InspectErr calls the provided closure with a reference to the contained E (if E).
func (r EnumResult[T, E]) InspectErr(f func(E)) EnumResult[T, E] {
	if r.IsErr() {
		f(r.safeGetE())
	}
	return r
}

// Expect returns the contained T value.
// Panics if the value is an E, with a panic message including the
// passed message, and the content of the E.
func (r EnumResult[T, E]) Expect(msg string) T {
	if r.IsErr() {
		panic(r.wrapError(msg))
	}
	return r.safeGetT()
}

func (r EnumResult[T, E]) wrapError(msg string) error {
	e := any(r.safeGetE())
	if err, ok := e.(error); ok {
		return ToErrBox(fmt.Errorf("%s: %w", msg, err))
	}
	return ToErrBox(fmt.Errorf("%s: %v", msg, e))
}

// Unwrap returns the contained T value.
// Because this function may panic, its use is generally discouraged.
// Instead, prefer to use pattern matching and handle the E case explicitly, or call UnwrapOr or UnwrapOrElse.
func (r EnumResult[T, E]) Unwrap() T {
	if r.IsErr() {
		panic(ToErrBox(r.safeGetE()))
	}
	return r.safeGetT()
}

// UnwrapOrDefault returns the contained T or a non-nil-pointer zero T.
func (r EnumResult[T, E]) UnwrapOrDefault() T {
	if r.IsOk() {
		return r.safeGetT()
	}
	return defaultValue[T]()
}

// ExpectErr returns the contained E.
// Panics if the value is not an E, with a panic message including the
// passed message, and the content of the T.
func (r EnumResult[T, E]) ExpectErr(msg string) E {
	if r.IsErr() {
		return r.safeGetE()
	}
	panic(ToErrBox(fmt.Sprintf("%s: %v", msg, r.safeGetT())))
}

// UnwrapErr returns the contained E.
// Panics if the value is not an E, with a custom panic message provided
// by the T's value.
func (r EnumResult[T, E]) UnwrapErr() E {
	if r.IsErr() {
		return r.safeGetE()
	}
	panic(ToErrBox(fmt.Sprintf("called `EnumResult.UnwrapErr()` on an `ok` value: %v", r.safeGetT())))
}

// And returns res if the result is T, otherwise returns the E of self.
func (r EnumResult[T, E]) And(res EnumResult[T, E]) EnumResult[T, E] {
	if r.IsErr() {
		return r
	}
	return res
}

// XAnd returns res if the result is T, otherwise returns the E of self.
func (r EnumResult[T, E]) XAnd(res EnumResult[any, E]) EnumResult[any, E] {
	if r.IsErr() {
		return EnumErr[any](r.safeGetE())
	}
	return res
}

// AndThen calls op if the result is T, otherwise returns the E of self.
// This function can be used for control flow based on EnumResult values.
func (r EnumResult[T, E]) AndThen(op func(T) EnumResult[T, E]) EnumResult[T, E] {
	if r.IsErr() {
		return r
	}
	return op(r.safeGetT())
}

// XAndThen calls op if the result is ok, otherwise returns the E of self.
// This function can be used for control flow based on EnumResult values.
func (r EnumResult[T, E]) XAndThen(op func(T) EnumResult[any, E]) EnumResult[any, E] {
	if r.IsErr() {
		return EnumErr[any, E](r.safeGetE())
	}
	return op(r.safeGetT())
}

// Or returns res if the result is E, otherwise returns the T value of r.
// Arguments passed to or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use OrElse, which is lazily evaluated.
func (r EnumResult[T, E]) Or(res EnumResult[T, E]) EnumResult[T, E] {
	if r.IsErr() {
		return res
	}
	return r
}

// XOr returns res if the result is E, otherwise returns the T value of r.
// Arguments passed to or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use XOrElse, which is lazily evaluated.
func (r EnumResult[T, E]) XOr(res EnumResult[T, any]) EnumResult[T, any] {
	if r.IsErr() {
		return res
	}
	return EnumOk[T, any](r.safeGetT())
}

// OrElse calls op if the result is E, otherwise returns the T value of self.
// This function can be used for control flow based on result values.
func (r EnumResult[T, E]) OrElse(op func(E) EnumResult[T, E]) EnumResult[T, E] {
	if r.IsErr() {
		return op(r.safeGetE())
	}
	return r
}

// XOrElse calls op if the result is E, otherwise returns the T value of self.
// This function can be used for control flow based on result values.
func (r EnumResult[T, E]) XOrElse(op func(E) EnumResult[T, any]) EnumResult[T, any] {
	if r.IsErr() {
		return op(r.safeGetE())
	}
	return EnumOk[T, any](r.safeGetT())
}

// UnwrapOr returns the contained T value or a provided default.
// Arguments passed to UnwrapOr are eagerly evaluated; if you are passing the result of a function call, it is recommended to use UnwrapOrElse, which is lazily evaluated.
func (r EnumResult[T, E]) UnwrapOr(defaultT T) T {
	if r.IsErr() {
		return defaultT
	}
	return r.safeGetT()
}

// UnwrapOrElse returns the contained T value or computes it from a closure.
func (r EnumResult[T, E]) UnwrapOrElse(defaultFn func(E) T) T {
	if r.IsErr() {
		return defaultFn(r.safeGetE())
	}
	return r.safeGetT()
}

// AsPtr returns its pointer or nil.
func (r EnumResult[T, E]) AsPtr() *T {
	return r.t.AsPtr()
}

func (r EnumResult[T, E]) MarshalJSON() ([]byte, error) {
	if r.IsErr() {
		return nil, toError(r.safeGetE())
	}
	return json.Marshal(r.safeGetT())
}

func (r *EnumResult[T, E]) UnmarshalJSON(b []byte) error {
	var t T
	if r == nil {
		return &json.InvalidUnmarshalError{Type: reflect.TypeOf(t)}
	}
	err := json.Unmarshal(b, &t)
	if err != nil {
		r.t = None[T]()
		e := fromError[E](err)
		r.e = &e
	} else {
		r.t = Some(t)
		r.e = nil
	}
	return err
}

func fromError[E any](e error) E {
	if x, is := e.(E); is {
		return x
	}
	var x E
	return x
}

var (
	_ Iterable[any]   = new(EnumResult[any, any])
	_ DeIterable[any] = new(EnumResult[any, any])
)

func (r *EnumResult[T, E]) Next() Option[T] {
	if r == nil {
		return None[T]()
	}
	return r.t.Next()
}

//go:inline
func (r *EnumResult[T, E]) NextBack() Option[T] {
	return r.Next()
}

func (r *EnumResult[T, E]) Remaining() uint {
	if r == nil {
		return 0
	}
	return r.t.Remaining()
}

// CtrlFlow returns the `CtrlFlow[E, T]`.
func (r EnumResult[T, E]) CtrlFlow() CtrlFlow[E, T] {
	if r.IsErr() {
		return Break[E, T](r.safeGetE())
	}
	return Continue[E, T](r.safeGetT())
}

// UnwrapOrThrow returns the contained T or panic returns E (panicValue[*E]).
// NOTE:
//
//	If there is an E, that panic should be caught with CatchEnumResult[U, E]
func (r EnumResult[T, E]) UnwrapOrThrow() T {
	if r.IsErr() {
		panic(panicValue[E]{r.e})
	}
	return r.safeGetT()
}

// CatchEnumResult catches panic caused by EnumResult[U, E].UnwrapOrThrow() or Errable[E].TryThrow(), and sets E to *EnumResult[T,E]
// Example:
//
//	```go
//	func example() (result EnumResult[int, string]) {
//		defer result.Catch()
//		EnumErr[string, string]("error").UnwrapOrThrow()
//		return
//	}
//	```
func (r *EnumResult[T, E]) Catch() {
	switch p := recover().(type) {
	case nil:
	case panicValue[E]:
		if r == nil {
			panic(p.ValueOrDefault())
		}
		if r.t.IsSome() {
			r.t = None[T]()
		}
		r.e = p.value
	default:
		panic(p)
	}
}

// CatchEnumResult catches panic caused by EnumResult[T, E].UnwrapOrThrow() or Errable[E].TryThrow(), and sets E to *EnumResult[U,E]
// Example:
//
//	```go
//	func example() (result EnumResult[int, string]) {
//		defer CatchEnumResult[int, string](&result)
//		EnumErr[string, string]("error").UnwrapOrThrow()
//		return
//	}
//	```
func CatchEnumResult[U any, E any](result *EnumResult[U, E]) {
	switch p := recover().(type) {
	case nil:
	case panicValue[E]:
		if result == nil {
			panic(p.ValueOrDefault())
		}
		if result.t.IsSome() {
			result.t = None[U]()
		}
		result.e = p.value
	default:
		panic(p)
	}
}
