package gust

import (
	"encoding/json"
	"fmt"
)

// EnumOk wraps a successful result enumeration.
func EnumOk[T any, E any](ok T) EnumResult[T, E] {
	return EnumResult[T, E]{val: ok, isErr: false}
}

// EnumErr wraps a failure result enumeration.
func EnumErr[T any, E any](err E) EnumResult[T, E] {
	return EnumResult[T, E]{val: err, isErr: true}
}

// EnumResult represents a success (T) or failure (E) enumeration.
type EnumResult[T any, E any] struct {
	val   any
	isErr bool
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
		return fmt.Sprintf("Err(%v)", r.val)
	}
	return fmt.Sprintf("Ok(%v)", r.val)
}

// IsOkAnd returns true if the result is Ok and the value inside it matches a predicate.
func (r EnumResult[T, E]) IsOkAnd(f func(T) bool) bool {
	if r.IsOk() {
		return f(r.val.(T))
	}
	return false
}

// IsErrAnd returns true if the result is E and the value inside it matches a predicate.
func (r EnumResult[T, E]) IsErrAnd(f func(E) bool) bool {
	if r.IsErr() {
		return f(r.val.(E))
	}
	return false
}

// Ok converts from `Result[T,E]` to `Option[T,E]`.
func (r EnumResult[T, E]) Ok() Option[T] {
	if r.IsOk() {
		return Some(r.val.(T))
	}
	return None[T]()
}

// Err returns E value.
func (r EnumResult[T, E]) Err() Option[E] {
	if r.IsErr() {
		return Some(r.val.(E))
	}
	return None[E]()
}

// Map maps a EnumResult[T,E] to EnumResult[T,E] by applying a function to a contained T value, leaving an E untouched.
// This function can be used to compose the results of two functions.
func (r EnumResult[T, E]) Map(f func(T) T) EnumResult[T, E] {
	if r.IsOk() {
		return EnumOk[T, E](f(r.val.(T)))
	}
	return EnumErr[T, E](r.val.(E))
}

// MapOr returns the provided default (if E), or applies a function to the contained value (if no E),
// Arguments passed to map_or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use MapOrElse, which is lazily evaluated.
func (r EnumResult[T, E]) MapOr(defaultOk T, f func(T) T) T {
	if r.IsOk() {
		return f(r.val.(T))
	}
	return defaultOk
}

// MapOrElse maps a EnumResult[T,E] to T by applying fallback function default to a contained E, or function f to a contained T value.
// This function can be used to unpack a successful result while handling an E.
func (r EnumResult[T, E]) MapOrElse(defaultFn func(E) T, f func(T) T) T {
	if r.IsOk() {
		return f(r.val.(T))
	}
	return defaultFn(r.val.(E))
}

// MapErr maps a EnumResult[T,E] to EnumResult[T,E] by applying a function to a contained E, leaving an T value untouched.
// This function can be used to pass through a successful result while handling an error.
func (r EnumResult[T, E]) MapErr(op func(E) E) EnumResult[T, E] {
	if r.IsErr() {
		r.val = op(r.val.(E))
	}
	return r
}

// Inspect calls the provided closure with a reference to the contained value (if no E).
func (r EnumResult[T, E]) Inspect(f func(T)) EnumResult[T, E] {
	if r.IsOk() {
		f(r.val.(T))
	}
	return r
}

// InspectErr calls the provided closure with a reference to the contained E (if E).
func (r EnumResult[T, E]) InspectErr(f func(E)) EnumResult[T, E] {
	if r.IsErr() {
		f(r.val.(E))
	}
	return r
}

// Expect returns the contained T value.
// Panics if the value is an E, with a panic message including the
// passed message, and the content of the E.
func (r EnumResult[T, E]) Expect(msg string) T {
	if r.IsErr() {
		panic(fmt.Errorf("%s: %v", msg, r.val))
	}
	return r.val.(T)
}

// Unwrap returns the contained T value.
// Because this function may panic, its use is generally discouraged.
// Instead, prefer to use pattern matching and handle the E case explicitly, or call UnwrapOr or UnwrapOrElse.
func (r EnumResult[T, E]) Unwrap() T {
	if r.IsErr() {
		panic(fmt.Errorf("called `Result.Unwrap()` on an `err` value: %s", r.val))
	}
	return r.val.(T)
}

// ExpectErr returns the contained E.
// Panics if the value is not an E, with a panic message including the
// passed message, and the content of the T.
func (r EnumResult[T, E]) ExpectErr(msg string) E {
	if r.IsErr() {
		return r.val.(E)
	}
	panic(fmt.Errorf("%s: %v", msg, r.val))
}

// UnwrapErr returns the contained E.
// Panics if the value is not an E, with a custom panic message provided
// by the T's value.
func (r EnumResult[T, E]) UnwrapErr() E {
	if r.IsErr() {
		return r.val.(E)
	}
	panic(fmt.Errorf("called `Result.UnwrapErr()` on an `ok` value: %v", r.val))
}

// And returns res if the result is T, otherwise returns the E of self.
func (r EnumResult[T, E]) And(res EnumResult[T, E]) EnumResult[T, E] {
	if r.IsErr() {
		return r
	}
	return res
}

// AndThen calls op if the result is T, otherwise returns the E of self.
// This function can be used for control flow based on EnumResult values.
func (r EnumResult[T, E]) AndThen(op func(T) EnumResult[T, E]) EnumResult[T, E] {
	if r.IsErr() {
		return r
	}
	return op(r.val.(T))
}

// Or returns res if the result is E, otherwise returns the T value of r.
// Arguments passed to or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use OrElse, which is lazily evaluated.
func (r EnumResult[T, E]) Or(res EnumResult[T, E]) EnumResult[T, E] {
	if r.IsErr() {
		return res
	}
	return r
}

// OrElse calls op if the result is E, otherwise returns the T value of self.
// This function can be used for control flow based on result values.
func (r EnumResult[T, E]) OrElse(op func(E) EnumResult[T, E]) EnumResult[T, E] {
	if r.IsErr() {
		return op(r.val.(E))
	}
	return r
}

// UnwrapOr returns the contained T value or a provided default.
// Arguments passed to UnwrapOr are eagerly evaluated; if you are passing the result of a function call, it is recommended to use UnwrapOrElse, which is lazily evaluated.
func (r EnumResult[T, E]) UnwrapOr(defaultT T) T {
	if r.IsErr() {
		return defaultT
	}
	return r.val.(T)
}

// UnwrapOrElse returns the contained T value or computes it from a closure.
func (r EnumResult[T, E]) UnwrapOrElse(defaultFn func(E) T) T {
	if r.IsErr() {
		return defaultFn(r.val.(E))
	}
	return r.val.(T)
}

func (r EnumResult[T, E]) MarshalJSON() ([]byte, error) {
	if r.IsErr() {
		return nil, toError(r.val)
	}
	return json.Marshal(r.val)
}

func (r *EnumResult[T, E]) UnmarshalJSON(b []byte) error {
	var t T
	err := json.Unmarshal(b, &t)
	if err != nil {
		r.isErr = true
		r.val = fromError[E](err)
	} else {
		r.isErr = false
		r.val = t
	}
	return err
}

func toError(e any) error {
	if err, ok := e.(error); ok {
		return err
	}
	return fmt.Errorf("%v", e)
}

func fromError[E any](e error) E {
	if x, is := e.(E); is {
		return x
	}
	var x E
	return x
}
