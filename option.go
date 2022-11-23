package gust

import (
	"encoding/json"
	"fmt"
	"unsafe"
)

// BoolOpt wraps a value as an Option.
// NOTE:
//
//	`ok=true` is wrapped as Some,
//	and `ok=false` is wrapped as None.
func BoolOpt[T any](val T, ok bool) Option[T] {
	if !ok {
		return Option[T]{value: nil}
	}
	v := &val
	return Option[T]{value: &v}
}

// AssertOpt wraps a value as an Option.
func AssertOpt[T any](i any) Option[T] {
	val, ok := i.(T)
	if !ok {
		return Option[T]{value: nil}
	}
	v := &val
	return Option[T]{value: &v}
}

// BoolAssertOpt wraps a value as an Option.
func BoolAssertOpt[T any](i any, ok bool) Option[T] {
	if ok {
		if val, ok2 := i.(T); ok2 {
			v := &val
			return Option[T]{value: &v}
		}
	}
	return Option[T]{value: nil}
}

// PtrOpt wraps a pointer value.
// NOTE:
//
//	`non-nil pointer` is wrapped as Some,
//	and `nil pointer` is wrapped as None.
func PtrOpt[U any, T *U](ptr T) Option[T] {
	if ptr == nil {
		return Option[T]{value: nil}
	}
	v := &ptr
	return Option[T]{value: &v}
}

// ElemOpt wraps a value from pointer.
// NOTE:
//
//	`non-nil pointer` is wrapped as Some,
//	and `nil pointer` is wrapped as None.
func ElemOpt[T any](ptr *T) Option[T] {
	if ptr == nil {
		return Option[T]{value: nil}
	}
	return Option[T]{value: &ptr}
}

// ZeroOpt wraps a value as an Option.
// NOTE:
//
//	`non-zero T` is wrapped as Some,
//	and `zero T` is wrapped as None.
func ZeroOpt[T comparable](val T) Option[T] {
	var zero T
	if zero == val {
		return Option[T]{value: nil}
	}
	v := &val
	return Option[T]{value: &v}
}

// Some wraps a non-none value.
// NOTE:
//
//	Option[T].IsSome() returns true.
//	and Option[T].IsNone() returns false.
func Some[T any](value T) Option[T] {
	v := &value
	return Option[T]{value: &v}
}

// None returns a none.
// NOTE:
//
//	Option[T].IsNone() returns true,
//	and Option[T].IsSome() returns false.
func None[T any]() Option[T] {
	return Option[T]{value: nil}
}

// Option can be used to avoid `(T, bool)` and `if *U != nil`,
// represents an optional value:
//
//	every [`Option`] is either [`Some`](which is non-none T), or [`None`](which is none).
type Option[T any] struct {
	value **T
}

// String returns the string representation.
func (o Option[T]) String() string {
	if o.IsNone() {
		return "None"
	}
	return fmt.Sprintf("Some(%v)", o.UnwrapUnchecked())
}

// ToX converts to `Option[any]`.
func (o Option[T]) ToX() Option[any] {
	if o.IsNone() {
		return None[any]()
	}
	return Some[any](o.UnwrapUnchecked())
}

// IsSome returns `true` if the option has value.
func (o Option[T]) IsSome() bool {
	return !o.IsNone()
}

// IsSomeAnd returns `true` if the option has value and the value inside it matches a predicate.
func (o Option[T]) IsSomeAnd(f func(T) bool) bool {
	if o.IsSome() {
		return f(o.UnwrapUnchecked())
	}
	return false
}

// IsNone returns `true` if the option is none.
func (o Option[T]) IsNone() bool {
	return o.value == nil || *o.value == nil
}

// Expect returns the contained [`Some`] value.
// Panics if the value is none with a custom panic message provided by `msg`.
func (o Option[T]) Expect(msg string) T {
	if o.IsNone() {
		panic(ToErrBox(msg))
	}
	return o.UnwrapUnchecked()
}

// Unwrap returns the contained value.
// Panics if the value is none.
func (o Option[T]) Unwrap() T {
	if o.IsSome() {
		return o.UnwrapUnchecked()
	}
	var t T
	panic(ToErrBox(fmt.Sprintf("call Option[%T].Unwrap() on none", t)))
}

// UnwrapOr returns the contained value or a provided fallback value.
func (o Option[T]) UnwrapOr(fallbackValue T) T {
	if o.IsSome() {
		return o.UnwrapUnchecked()
	}
	return fallbackValue
}

// UnwrapOrElse returns the contained value or computes it from a closure.
func (o Option[T]) UnwrapOrElse(defaultSome func() T) T {
	if o.IsSome() {
		return o.UnwrapUnchecked()
	}
	return defaultSome()
}

// UnwrapOrDefault returns the contained value or a non-nil-pointer zero value.
func (o Option[T]) UnwrapOrDefault() T {
	if o.IsSome() {
		return o.UnwrapUnchecked()
	}
	return defaultValue[T]()
}

// Take takes the value out of the option, leaving a [`None`] in its place.
func (o Option[T]) Take() Option[T] {
	if o.IsNone() {
		return None[T]()
	}
	v := *o.value
	*o.value = nil
	return Option[T]{value: &v}
}

// UnwrapUnchecked returns the contained value.
func (o Option[T]) UnwrapUnchecked() T {
	return **o.value
}

// Map maps an `Option[T]` to `Option[T]` by applying a function to a contained value.
func (o Option[T]) Map(f func(T) T) Option[T] {
	if o.IsSome() {
		return Some[T](f(o.UnwrapUnchecked()))
	}
	return None[T]()
}

// XMap maps an `Option[T]` to `Option[any]` by applying a function to a contained value.
func (o Option[T]) XMap(f func(T) any) Option[any] {
	if o.IsSome() {
		return Some[any](f(o.UnwrapUnchecked()))
	}
	return None[any]()
}

// Inspect calls the provided closure with a reference to the contained value (if it has value).
func (o Option[T]) Inspect(f func(T)) Option[T] {
	if o.IsSome() {
		f(o.UnwrapUnchecked())
	}
	return o
}

// InspectNone calls the provided closure (if it is none).
func (o Option[T]) InspectNone(f func()) Option[T] {
	if o.IsNone() {
		f()
	}
	return o
}

// MapOr returns the provided default value (if none),
// or applies a function to the contained value (if any).
func (o Option[T]) MapOr(defaultSome T, f func(T) T) T {
	if o.IsSome() {
		return f(o.UnwrapUnchecked())
	}
	return defaultSome
}

// XMapOr returns the provided default value (if none),
// or applies a function to the contained value (if any).
func (o Option[T]) XMapOr(defaultSome any, f func(T) any) any {
	if o.IsSome() {
		return f(o.UnwrapUnchecked())
	}
	return defaultSome
}

// MapOrElse computes a default function value (if none), or
// applies a different function to the contained value (if any).
func (o Option[T]) MapOrElse(defaultFn func() T, f func(T) T) T {
	if o.IsSome() {
		return f(o.UnwrapUnchecked())
	}
	return defaultFn()
}

// XMapOrElse computes a default function value (if none), or
// applies a different function to the contained value (if any).
func (o Option[T]) XMapOrElse(defaultFn func() any, f func(T) any) any {
	if o.IsSome() {
		return f(o.UnwrapUnchecked())
	}
	return defaultFn()
}

// OkOr transforms the `Option[T]` into a [`Result[T]`], mapping [`Some(v)`] to
// [`Ok(v)`] and [`None`] to [`Err(err)`].
func (o Option[T]) OkOr(err any) Result[T] {
	if o.IsSome() {
		return Ok(o.UnwrapUnchecked())
	}
	return Err[T](err)
}

// XOkOr transforms the `Option[T]` into a [`Result[any]`], mapping [`Some(v)`] to
// [`Ok(v)`] and [`None`] to [`Err(err)`].
func (o Option[T]) XOkOr(err any) Result[any] {
	if o.IsSome() {
		return Ok[any](o.UnwrapUnchecked())
	}
	return Err[any](err)
}

// OkOrElse transforms the `Option[T]` into a [`Result[T]`], mapping [`Some(v)`] to
// [`Ok(v)`] and [`None`] to [`Err(errFn())`].
func (o Option[T]) OkOrElse(errFn func() any) Result[T] {
	if o.IsSome() {
		return Ok(o.UnwrapUnchecked())
	}
	return Err[T](errFn())
}

// XOkOrElse transforms the `Option[T]` into a [`Result[any]`], mapping [`Some(v)`] to
// [`Ok(v)`] and [`None`] to [`Err(errFn())`].
func (o Option[T]) XOkOrElse(errFn func() any) Result[any] {
	if o.IsSome() {
		return Ok[any](o.UnwrapUnchecked())
	}
	return Err[any](errFn())
}

// EnumOkOr transforms the `Option[T]` into a [`EnumResult[T,any]`], mapping [`Some(v)`] to
// [`EnumOk(v)`] and [`None`] to [`EnumErr(err)`].
func (o Option[T]) EnumOkOr(err any) EnumResult[T, any] {
	if o.IsSome() {
		return EnumOk[T, any](o.UnwrapUnchecked())
	}
	return EnumErr[T, any](err)
}

// XEnumOkOr transforms the `Option[T]` into a [`EnumResult[any,any]`], mapping [`Some(v)`] to
// [`EnumOk(v)`] and [`None`] to [`EnumErr(err)`].
func (o Option[T]) XEnumOkOr(err any) EnumResult[any, any] {
	if o.IsSome() {
		return EnumOk[any, any](o.UnwrapUnchecked())
	}
	return EnumErr[any, any](err)
}

// EnumOkOrElse transforms the `Option[T]` into a [`EnumResult[T,any]`], mapping [`Some(v)`] to
// [`EnumOk(v)`] and [`None`] to [`EnumErr(errFn())`].
func (o Option[T]) EnumOkOrElse(errFn func() any) EnumResult[T, any] {
	if o.IsSome() {
		return EnumOk[T, any](o.UnwrapUnchecked())
	}
	return EnumErr[T, any](errFn())
}

// XEnumOkOrElse transforms the `Option[T]` into a [`EnumResult[any,any]`], mapping [`Some(v)`] to
// [`EnumOk(v)`] and [`None`] to [`EnumErr(errFn())`].
func (o Option[T]) XEnumOkOrElse(errFn func() any) EnumResult[any, any] {
	if o.IsSome() {
		return EnumOk[any, any](o.UnwrapUnchecked())
	}
	return EnumErr[any, any](errFn())
}

// AndThen returns [`None`] if the option is [`None`], otherwise calls `f` with the
func (o Option[T]) AndThen(f func(T) Option[T]) Option[T] {
	if o.IsNone() {
		return o
	}
	return f(o.UnwrapUnchecked())
}

// XAndThen returns [`None`] if the option is [`None`], otherwise calls `f` with the
func (o Option[T]) XAndThen(f func(T) Option[any]) Option[any] {
	if o.IsNone() {
		return None[any]()
	}
	return f(o.UnwrapUnchecked())
}

// OrElse returns the option if it contains a value, otherwise calls `f` and returns the result.
func (o Option[T]) OrElse(f func() Option[T]) Option[T] {
	if o.IsNone() {
		return f()
	}
	return o
}

// Filter returns [`None`] if the option is [`None`], otherwise calls `predicate`
// with the wrapped value and returns.
func (o Option[T]) Filter(predicate func(T) bool) Option[T] {
	if o.IsSome() {
		if predicate(o.UnwrapUnchecked()) {
			return o
		}
	}
	return None[T]()
}

// And returns [`None`] if the option is [`None`], otherwise returns `optb`.
func (o Option[T]) And(optb Option[T]) Option[T] {
	if o.IsSome() {
		return optb
	}
	return o
}

// XAnd returns [`None`] if the option is [`None`], otherwise returns `optb`.
func (o Option[T]) XAnd(optb Option[any]) Option[any] {
	if o.IsSome() {
		return optb
	}
	return None[any]()
}

// Or returns the option if it contains a value, otherwise returns `optb`.
func (o Option[T]) Or(optb Option[T]) Option[T] {
	if o.IsNone() {
		return optb
	}
	return o
}

// Xor [`Some`] if exactly one of `self`, `optb` is [`Some`], otherwise returns [`None`].
func (o Option[T]) Xor(optb Option[T]) Option[T] {
	if o.IsSome() && optb.IsNone() {
		return o
	}
	if o.IsNone() && optb.IsSome() {
		return optb
	}
	return None[T]()
}

// Insert inserts `value` into the option, then returns its pointer.
func (o *Option[T]) Insert(some T) *T {
	v := &some
	o.value = &v
	return v
}

// GetOrInsert inserts `value` into the option if it is [`None`], then
// returns the contained value pointer.
func (o *Option[T]) GetOrInsert(some T) *T {
	if o.IsNone() {
		v := &some
		o.value = &v
	}
	return *o.value
}

// GetOrInsertWith inserts a value computed from `f` into the option if it is [`None`],
// then returns the contained value.
func (o *Option[T]) GetOrInsertWith(f func() T) *T {
	if o.IsNone() {
		var v *T
		if f == nil {
			v = defaultValuePtr[T]()
		} else {
			var some = f()
			v = &some
		}
		o.value = &v
	}
	return *o.value
}

// GetOrInsertDefault inserts default value into the option if it is [`None`], then
// returns the contained value pointer.
func (o *Option[T]) GetOrInsertDefault() *T {
	if o.IsNone() {
		v := defaultValuePtr[T]()
		o.value = &v
	}
	return *o.value
}

// AsPtr returns its pointer or nil.
func (o *Option[T]) AsPtr() *T {
	if o.value == nil {
		return nil
	}
	return *o.value
}

// Replace replaces the actual value in the option by the value given in parameter,
// returning the old value if present,
// leaving a [`Some`] in its place without deinitializing either one.
func (o *Option[T]) Replace(some T) (old Option[T]) {
	old.value = o.value
	v := &some
	o.value = &v
	return old
}

const null = "null"

func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.IsNone() {
		return []byte(null), nil
	}
	return json.Marshal(o.value)
}

func (o *Option[T]) UnmarshalJSON(b []byte) error {
	o.value = nil
	if *(*string)(unsafe.Pointer(&b)) == null {
		return nil
	}
	var value = new(T)
	err := json.Unmarshal(b, value)
	if err == nil {
		o.value = &value
	}
	return err
}

var (
	_ Iterable[any]   = Option[any]{}
	_ DeIterable[any] = Option[any]{}
)

func (o Option[T]) Next() Option[T] {
	if o.IsNone() {
		return o
	}
	v := o.Unwrap()
	*o.value = nil
	return Some(v)
}

func (o Option[T]) NextBack() Option[T] {
	return o.Next()
}

func (o Option[T]) Remaining() uint {
	if o.IsNone() {
		return 0
	}
	return 1
}

// CtrlFlow returns the `CtrlFlow[Void, T]`.
func (o Option[T]) CtrlFlow() CtrlFlow[Void, T] {
	if o.IsNone() {
		return Break[Void, T](nil)
	}
	return Continue[Void, T](o.Unwrap())
}

// ToErrable converts from `Option[T]` to `Errable[T]`.
func (o Option[T]) ToErrable() Errable[T] {
	if o.IsSome() {
		return ToErrable[T](o.UnwrapUnchecked())
	}
	return NonErrable[T]()
}
