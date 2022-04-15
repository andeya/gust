package gust

import (
	"fmt"
)

// Opt wraps a value.
func Opt[T any](value *T) Option[T] {
	return Option[T]{value: value}
}

// Some wraps a nonnull value.
func Some[T any](value T) Option[T] {
	return Option[T]{value: &value}
}

// None returns a none.
func None[T any]() Option[T] {
	return Option[T]{value: nil}
}

// Option represents an optional value:
// every [`Option`] is either [`Some`](which is nonnull T), or [`None`](which is nil).
type Option[T any] struct {
	value *T
}

// String returns the string representation.
func (o Option[T]) String() string {
	if o.IsNone() {
		return "None"
	}
	return fmt.Sprintf("Some(%v)", *o.value)
}

// IsSome returns `true` if the option has value.
func (o Option[T]) IsSome() bool {
	return !o.IsNone()
}

// IsSomeAnd returns `true` if the option has value and the value inside of it matches a predicate.
func (o Option[T]) IsSomeAnd(f func(T) bool) bool {
	if o.IsSome() {
		return f(*o.value)
	}
	return false
}

// IsNone returns `true` if the option is none.
func (o Option[T]) IsNone() bool {
	return o.value == nil
}

// Expect returns the contained [`Some`] value.
// Panics if the value is null with a custom panic message provided by `msg`.
func (o Option[T]) Expect(msg string) T {
	if o.IsNone() {
		panic(fmt.Errorf("%s", msg))
	}
	return *o.value
}

// Unwrap returns the contained value.
// Panics if the value is null.
func (o Option[T]) Unwrap() T {
	if o.IsSome() {
		return *o.value
	}
	var t T
	panic(fmt.Sprintf("call Option[%T].Unwrap() on nonnull", t))
}

// UnwrapOr returns the contained value or a provided default.
func (o Option[T]) UnwrapOr(defaultSome T) T {
	if o.IsSome() {
		return *o.value
	}
	return defaultSome
}

// UnwrapOrElse returns the contained value or computes it from a closure.
func (o Option[T]) UnwrapOrElse(defaultSome func() T) T {
	if o.IsSome() {
		return *o.value
	}
	return defaultSome()
}

// UnwrapUnchecked returns the contained value.
func (o Option[T]) UnwrapUnchecked() T {
	return *o.value
}

// Map maps an `Option[T]` to `Option[T]` by applying a function to a contained value.
func (o Option[T]) Map(f func(T) T) Option[T] {
	if o.IsSome() {
		return Some[T](f(*o.value))
	}
	return None[T]()
}

// Inspect calls the provided closure with a reference to the contained value (if it has value).
func (o Option[T]) Inspect(f func(T)) Option[T] {
	if o.IsSome() {
		f(*o.value)
	}
	return o
}

// MapOr returns the provided default value (if none),
// or applies a function to the contained value (if any).
func (o Option[T]) MapOr(defaultSome T, f func(T) T) T {
	if o.IsSome() {
		return f(*o.value)
	}
	return defaultSome
}

// MapOrElse computes a default function value (if none), or
// applies a different function to the contained value (if any).
func (o Option[T]) MapOrElse(defaultFn func() T, f func(T) T) T {
	if o.IsSome() {
		return f(*o.value)
	}
	return defaultFn()
}

// OkOr transforms the `Option[T]` into a [`Result[T]`], mapping [`Some(v)`] to
// [`Ok(v)`] and [`None`] to [`Err(err)`].
func (o Option[T]) OkOr(err error) Result[T] {
	if o.IsSome() {
		return Ok(o.Unwrap())
	}
	return Err[T](err)
}

// OkOrElse transforms the `Option[T]` into a [`Result[T]`], mapping [`Some(v)`] to
// [`Ok(v)`] and [`None`] to [`Err(errFn())`].
func (o Option[T]) OkOrElse(errFn func() error) Result[T] {
	if o.IsSome() {
		return Ok(o.Unwrap())
	}
	return Err[T](errFn())
}

// And returns [`None`] if the option is [`None`], otherwise returns `optb`.
func (o Option[T]) And(optb Option[T]) Option[T] {
	if o.IsSome() {
		return optb
	}
	return o
}

// AndThen returns [`None`] if the option is [`None`], otherwise calls `f` with the
func (o Option[T]) AndThen(f func(T) Option[T]) Option[T] {
	if o.IsNone() {
		return o
	}
	return f(*o.value)
}

// Filter returns [`None`] if the option is [`None`], otherwise calls `predicate`
// with the wrapped value and returns.
func (o Option[T]) Filter(predicate func(T) bool) Option[T] {
	if o.IsSome() {
		if predicate(*o.value) {
			return o
		}
	}
	return None[T]()
}

// Or returns the option if it contains a value, otherwise returns `optb`.
func (o Option[T]) Or(optb Option[T]) Option[T] {
	if o.IsNone() {
		return optb
	}
	return o
}

// OrElse returns [`None`] if the option is [`None`], otherwise calls `f` with the returns the result.
func (o Option[T]) OrElse(f func() Option[T]) Option[T] {
	if o.IsNone() {
		return f()
	}
	return o
}

// XorElse [`Some`] if exactly one of `self`, `optb` is [`Some`], otherwise returns [`None`].
func (o Option[T]) XorElse(optb Option[T]) Option[T] {
	if o.IsSome() && optb.IsNone() {
		return o
	}
	if o.IsNone() && optb.IsSome() {
		return optb
	}
	return None[T]()
}

// Insert inserts `value` into the option, then returns a reference to it.
func (o *Option[T]) Insert(some T) T {
	o.value = &some
	return some
}

// GetOrInsert inserts `value` into the option if it is [`None`], then
// returns a reference to the contained value.
func (o *Option[T]) GetOrInsert(some T) T {
	if o.IsNone() {
		o.value = &some
	}
	return *o.value
}

// GetOrInsertWith inserts a value computed from `f` into the option if it is [`None`],
// then returns a mutable reference to the contained value.
func (o *Option[T]) GetOrInsertWith(f func() T) T {
	if o.IsNone() {
		var some = f()
		o.value = &some
	}
	return *o.value
}

// Replace replaces the actual value in the option by the value given in parameter,
// returning the old value if present,
// leaving a [`Some`] in its place without deinitializing either one.
func (o *Option[T]) Replace(some T) *Option[T] {
	o.value = &some
	return o
}
