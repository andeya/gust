package gust

import (
	"fmt"
)

// Ptr wraps a value pointer.
func Ptr[T any](value *T) OptNil[T] {
	return OptNil[T]{value: value}
}

// Nil returns a nil.
func Nil[T any]() OptNil[T] {
	return OptNil[T]{value: nil}
}

// OptPtr converts `Option[T]` to `OptNil[T]`.
func OptPtr[T any](o Option[T]) OptNil[T] {
	return Ptr[T](o.value)
}

// OptNil represents an optional value:
// every [`OptNil`] is either [`NonNil`](which is nonnil *T), or [`Nil`](which is nil).
type OptNil[T any] struct {
	value *T
}

// String returns the string representation.
func (o OptNil[T]) String() string {
	if o.IsNil() {
		return "Nil"
	}
	return fmt.Sprintf("NonNil(%v)", o.value)
}

// ToOption converts to Option[T].
func (o OptNil[T]) ToOption() Option[T] {
	return Opt[T](o.value)
}

// NotNil returns `true` if the value is not nil.
func (o OptNil[T]) NotNil() bool {
	return !o.IsNil()
}

// NotNilAnd returns `true` if the option has value and the value inside of it matches a predicate.
func (o OptNil[T]) NotNilAnd(f func(*T) bool) bool {
	if o.NotNil() {
		return f(o.value)
	}
	return false
}

// IsNil returns `true` if the value is nil.
func (o OptNil[T]) IsNil() bool {
	return o.value == nil
}

// Expect returns the contained [`NonNil`] value.
// Panics if the value is nil with a custom panic message provided by `msg`.
func (o OptNil[T]) Expect(msg string) *T {
	if o.IsNil() {
		panic(fmt.Errorf("%s", msg))
	}
	return o.value
}

// Unwrap returns the contained value.
// Panics if the value is nil.
func (o OptNil[T]) Unwrap() *T {
	if o.NotNil() {
		return o.value
	}
	var t T
	panic(fmt.Sprintf("call OptNil[%T].Unwrap() on nonnil", t))
}

// UnwrapOr returns the contained value or a provided default.
func (o OptNil[T]) UnwrapOr(defaultPtr *T) *T {
	if o.NotNil() {
		return o.value
	}
	return defaultPtr
}

// UnwrapOrElse returns the contained value or computes it from a closure.
func (o OptNil[T]) UnwrapOrElse(defaultPtr func() *T) *T {
	if o.NotNil() {
		return o.value
	}
	return defaultPtr()
}

// UnwrapUnchecked returns the contained value.
func (o OptNil[T]) UnwrapUnchecked() *T {
	return o.value
}

// Map maps an `OptNil[T]` to `OptNil[T]` by applying a function to a contained value.
func (o OptNil[T]) Map(f func(*T) *T) OptNil[T] {
	if o.NotNil() {
		return Ptr[T](f(o.value))
	}
	return Nil[T]()
}

// Inspect calls the provided closure with a reference to the contained value (if it has value).
func (o OptNil[T]) Inspect(f func(*T)) OptNil[T] {
	if o.NotNil() {
		f(o.value)
	}
	return o
}

// MapOr returns the provided default value (if none),
// or applies a function to the contained value (if any).
func (o OptNil[T]) MapOr(defaultPtr *T, f func(*T) *T) *T {
	if o.NotNil() {
		return f(o.value)
	}
	return defaultPtr
}

// MapOrElse computes a default function value (if none), or
// applies a different function to the contained value (if any).
func (o OptNil[T]) MapOrElse(defaultFn func() *T, f func(*T) *T) *T {
	if o.NotNil() {
		return f(o.value)
	}
	return defaultFn()
}

// OkOr transforms the `Option[T]` into a [`Result[T]`], mapping [`Some(v)`] to
// [`Ok(v)`] and [`None`] to [`Err(err)`].
func (o OptNil[T]) OkOr(err error) Result[*T] {
	if o.NotNil() {
		return Ok(o.Unwrap())
	}
	return Err[*T](err)
}

// OkOrElse transforms the `Option[T]` into a [`Result[T]`], mapping [`Some(v)`] to
// [`Ok(v)`] and [`None`] to [`Err(errFn())`].
func (o OptNil[T]) OkOrElse(errFn func() error) Result[*T] {
	if o.NotNil() {
		return Ok(o.Unwrap())
	}
	return Err[*T](errFn())
}

// And returns [`Nil`] if the option is [`Nil`], otherwise returns `optb`.
func (o OptNil[T]) And(optb OptNil[T]) OptNil[T] {
	if o.NotNil() {
		return optb
	}
	return o
}

// AndThen returns [`Nil`] if the option is [`Nil`], otherwise calls `f` with the
func (o OptNil[T]) AndThen(f func(*T) OptNil[T]) OptNil[T] {
	if o.IsNil() {
		return o
	}
	return f(o.value)
}

// Filter returns [`Nil`] if the option is [`Nil`], otherwise calls `predicate`
// with the wrapped value and returns.
func (o OptNil[T]) Filter(predicate func(*T) bool) OptNil[T] {
	if o.NotNil() {
		if predicate(o.value) {
			return o
		}
	}
	return Nil[T]()
}

// Or returns the option if it contains a value, otherwise returns `optb`.
func (o OptNil[T]) Or(optb OptNil[T]) OptNil[T] {
	if o.IsNil() {
		return optb
	}
	return o
}

// OrElse returns [`Nil`] if the option is [`Nil`], otherwise calls `f` with the returns the result.
func (o OptNil[T]) OrElse(f func() OptNil[T]) OptNil[T] {
	if o.IsNil() {
		return f()
	}
	return o
}

// XorElse [`NonNil`] if exactly one of `self`, `optb` is [`NonNil`], otherwise returns [`Nil`].
func (o OptNil[T]) XorElse(optb OptNil[T]) OptNil[T] {
	if o.NotNil() && optb.IsNil() {
		return o
	}
	if o.IsNil() && optb.NotNil() {
		return optb
	}
	return Nil[T]()
}

// Insert inserts `value` into the option, then returns a reference to it.
func (o *OptNil[T]) Insert(some *T) *T {
	o.value = some
	return o.value
}

// GetOrInsert inserts `value` into the option if it is [`Nil`], then
// returns a reference to the contained value.
func (o *OptNil[T]) GetOrInsert(some *T) *T {
	if o.IsNil() {
		o.value = some
	}
	return o.value
}

// GetOrInsertWith inserts a value computed from `f` into the option if it is [`Nil`],
// then returns a mutable reference to the contained value.
func (o *OptNil[T]) GetOrInsertWith(f func() *T) *T {
	if o.IsNil() {
		o.value = f()
	}
	return o.value
}

// Replace replaces the actual value in the option by the value given in parameter,
// returning the old value if present,
// leaving a [`NonNil`] in its place without deinitializing either one.
func (o *OptNil[T]) Replace(some *T) *OptNil[T] {
	o.value = some
	return o
}
