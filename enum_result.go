package gust

import (
	"encoding/json"
	"fmt"
)

// EnumOk wraps a successful result enumeration.
func EnumOk[T any, E any](ok T) EnumResult[T, E] {
	v := any(ok)
	return EnumResult[T, E]{value: &v, isErr: false}
}

// EnumErr wraps a failure result enumeration.
func EnumErr[T any, E any](err E) EnumResult[T, E] {
	v := any(err)
	return EnumResult[T, E]{value: &v, isErr: true}
}

// EnumResult represents a success (T) or failure (E) enumeration.
type EnumResult[T any, E any] struct {
	value *any
	isErr bool
}

func (r EnumResult[T, E]) safeGetT() T {
	if r.isErr || r.value == nil {
		var t T
		return t
	}
	v, _ := (*r.value).(T)
	return v
}

func (r EnumResult[T, E]) safeGetE() E {
	if !r.isErr || r.value == nil {
		var e E
		return e
	}
	v, _ := (*r.value).(E)
	return v
}

// IsErr returns true if the result is E.
func (r EnumResult[T, E]) IsErr() bool {
	return r.isErr
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

// Ok converts from `Result[T,E]` to `Option[T,E]`.
func (r EnumResult[T, E]) Ok() Option[T] {
	if r.IsOk() {
		return Some(r.safeGetT())
	}
	return None[T]()
}

// Err returns E value.
func (r EnumResult[T, E]) Err() Option[E] {
	if r.IsErr() {
		return Some(r.safeGetE())
	}
	return None[E]()
}

// ToXOk converts from `EnumResult[T,E]` to EnumResult[any,E].
func (r EnumResult[T, E]) ToXOk() EnumResult[any, E] {
	return EnumResult[any, E]{
		value: r.value,
		isErr: r.isErr,
	}
}

// ToXErr converts from `EnumResult[T,E]` to EnumResult[T,any].
func (r EnumResult[T, E]) ToXErr() EnumResult[T, any] {
	return EnumResult[T, any]{
		value: r.value,
		isErr: r.isErr,
	}
}

// ToX converts from `EnumResult[T,E]` to EnumResult[any,any].
func (r EnumResult[T, E]) ToX() EnumResult[any, any] {
	return EnumResult[any, any]{
		value: r.value,
		isErr: r.isErr,
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
		panic(fmt.Sprintf("%s: %v", msg, r.safeGetE()))
	}
	return r.safeGetT()
}

// Unwrap returns the contained T value.
// Because this function may panic, its use is generally discouraged.
// Instead, prefer to use pattern matching and handle the E case explicitly, or call UnwrapOr or UnwrapOrElse.
func (r EnumResult[T, E]) Unwrap() T {
	if r.IsErr() {
		panic(fmt.Sprintf("called `EnumResult.Unwrap()` on an `err` value: %v", r.safeGetE()))
	}
	return r.safeGetT()
}

// ExpectErr returns the contained E.
// Panics if the value is not an E, with a panic message including the
// passed message, and the content of the T.
func (r EnumResult[T, E]) ExpectErr(msg string) E {
	if r.IsErr() {
		return r.safeGetE()
	}
	panic(fmt.Sprintf("%s: %v", msg, r.safeGetT()))
}

// UnwrapErr returns the contained E.
// Panics if the value is not an E, with a custom panic message provided
// by the T's value.
func (r EnumResult[T, E]) UnwrapErr() E {
	if r.IsErr() {
		return r.safeGetE()
	}
	panic(fmt.Sprintf("called `EnumResult.UnwrapErr()` on an `ok` value: %v", r.safeGetT()))
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

func (r EnumResult[T, E]) MarshalJSON() ([]byte, error) {
	if r.IsErr() {
		return nil, newAnyError(r.safeGetE())
	}
	return json.Marshal(r.safeGetT())
}

func (r *EnumResult[T, E]) UnmarshalJSON(b []byte) error {
	var t T
	err := json.Unmarshal(b, &t)
	if err != nil {
		r.isErr = true
		e := any(fromError[E](err))
		r.value = &e
	} else {
		r.isErr = false
		v := any(t)
		r.value = &v
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
	_ Iterable[any]       = EnumResult[any, any]{}
	_ SizeDeIterable[any] = EnumResult[any, any]{}
)

func (r EnumResult[T, E]) Next() Option[T] {
	if r.isErr || r.value == nil || *r.value == nil {
		return None[T]()
	}
	v := *r.value
	*r.value = nil
	return Some[T](v.(T))
}

func (r EnumResult[T, E]) NextBack() Option[T] {
	return r.Next()
}

func (r EnumResult[T, E]) Remaining() uint {
	if r.isErr || r.value == nil || *r.value == nil {
		return 0
	}
	return 1
}
