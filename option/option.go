// Package option provides helper functions for working with Option types.
package option

import (
	"github.com/andeya/gust/internal/core"
	"github.com/andeya/gust/pair"
	"github.com/andeya/gust/result"
)

// Option is an alias for core.Option[T].
// This allows using option.Option[T] instead of core.Option[T].
//
// Option represents an optional value: every Option is either Some (which is non-none T),
// or None (which is none). Option types are very common in Rust code, as they have
// a number of uses:
//
//   - Initial values
//   - Return values for functions that are not defined over their entire input range (partial functions)
//   - Return value for otherwise reporting simple errors, where None is returned on error
//   - Optional struct fields
//   - Optional function arguments
//   - Nullable pointers
//   - Swapping things out of difficult situations
type Option[T any] = core.Option[T]

// Some wraps a non-none value.
// NOTE:
//
//	Option[T].IsSome() returns true.
//	and Option[T].IsNone() returns false.
//
//go:inline
func Some[T any](value T) Option[T] {
	return core.Some(value)
}

// None returns a none.
// NOTE:
//
//	Option[T].IsNone() returns true,
//	and Option[T].IsSome() returns false.
//
//go:inline
func None[T any]() Option[T] {
	return core.None[T]()
}

// BoolOpt wraps a value as an Option.
// NOTE:
//
//	`ok=true` is wrapped as Some,
//	and `ok=false` is wrapped as None.
//
//go:inline
func BoolOpt[T any](v T, ok bool) Option[T] {
	return core.BoolOpt(v, ok)
}

// AssertOpt returns the Option[T] of asserting `i` to type `T`
//
//go:inline
func AssertOpt[T any](v any) Option[T] {
	return core.AssertOpt[T](v)
}

// BoolAssertOpt wraps a value as an Option.
// NOTE:
//
//	`ok=true` is wrapped as Some,
//	and `ok=false` is wrapped as None.
//
//go:inline
func BoolAssertOpt[T any](i any, ok bool) Option[T] {
	return core.BoolAssertOpt[T](i, ok)
}

// PtrOpt wraps a pointer value.
// NOTE:
//
//	`non-nil pointer` is wrapped as Some,
//	and `nil pointer` is wrapped as None.
//
//go:inline
func PtrOpt[U any, T *U](ptr T) Option[T] {
	return core.PtrOpt[U, T](ptr)
}

// ElemOpt wraps a value from pointer.
// NOTE:
//
//	`non-nil pointer` is wrapped as Some,
//	and `nil pointer` is wrapped as None.
//
//go:inline
func ElemOpt[T any](ptr *T) Option[T] {
	return core.ElemOpt(ptr)
}

// ZeroOpt wraps a value as an Option.
// NOTE:
//
//	`non-zero T` is wrapped as Some,
//	and `zero T` is wrapped as None.
//
//go:inline
func ZeroOpt[T comparable](v T) Option[T] {
	return core.ZeroOpt(v)
}

// RetOpt wraps a value as an `Option[T]`.
// NOTE:
//
//	`err != nil` is wrapped as None,
//	and `err == nil` is wrapped as Some.
//
//go:inline
func RetOpt[T any](v T, err error) Option[T] {
	return core.RetOpt(v, err)
}

// RetAnyOpt wraps a value as an `Option[any]`.
// NOTE:
//
//	`err != nil` or `value`==nil is wrapped as None,
//	and `err == nil` and `value != nil` is wrapped as Some.
//
//go:inline
func RetAnyOpt[T any](v any, err error) Option[any] {
	return core.RetAnyOpt[T](v, err)
}

// SafeAssert asserts Option[T] as result.Result[Option[U]].
// NOTE:
//
//	If the assertion fails, return error.
//
//go:inline
func SafeAssert[T any, U any](o Option[T]) result.Result[Option[U]] {
	if o.IsSome() {
		u, ok := any(o.UnwrapUnchecked()).(U)
		if ok {
			return result.Ok(Some[U](u))
		}
		return result.FmtErr[Option[U]]("type assert error, got %T, want %T", o.UnwrapUnchecked(), u)
	}
	return result.Ok(None[U]())
}

// XSafeAssert asserts Option[any] as result.Result[Option[U]].
// NOTE:
//
//	If the assertion fails, return error.
//
//go:inline
func XSafeAssert[U any](o Option[any]) result.Result[Option[U]] {
	if o.IsSome() {
		u, ok := o.UnwrapUnchecked().(U)
		if ok {
			return result.Ok(Some[U](u))
		}
		return result.FmtErr[Option[U]]("type assert error, got %T, want %T", o.UnwrapUnchecked(), u)
	}
	return result.Ok(None[U]())
}

// FuzzyAssert asserts Option[T] as Option[U].
// NOTE:
//
//	If the assertion fails, return none.
//
//go:inline
func FuzzyAssert[T any, U any](o Option[T]) Option[U] {
	if o.IsSome() {
		u, ok := any(o.UnwrapUnchecked()).(U)
		uVal := any(u)
		return BoolAssertOpt[U](uVal, ok)
	}
	return None[U]()
}

// XFuzzyAssert asserts Option[any] as Option[U].
// NOTE:
//
//	If the assertion fails, return none.
//
//go:inline
func XFuzzyAssert[U any](o Option[any]) Option[U] {
	if o.IsSome() {
		u, ok := o.UnwrapUnchecked().(U)
		uVal := any(u)
		return BoolAssertOpt[U](uVal, ok)
	}
	return None[U]()
}

// Map maps an `Option[T]` to `Option[U]` by applying a function to a contained value.
//
//go:inline
func Map[T any, U any](o Option[T], f func(T) U) Option[U] {
	if o.IsSome() {
		return Some[U](f(o.UnwrapUnchecked()))
	}
	return None[U]()
}

// MapOr returns the provided default value (if none),
// or applies a function to the contained value (if any).
//
//go:inline
func MapOr[T any, U any](o Option[T], defaultSome U, f func(T) U) U {
	if o.IsSome() {
		return f(o.UnwrapUnchecked())
	}
	return defaultSome
}

// MapOrElse computes a default function value (if none), or
// applies a different function to the contained value (if any).
//
//go:inline
func MapOrElse[T any, U any](o Option[T], defaultFn func() U, f func(T) U) U {
	if o.IsSome() {
		return f(o.UnwrapUnchecked())
	}
	return defaultFn()
}

// And returns [`None`] if the option is [`None`], otherwise returns `optb`.
//
//go:inline
func And[T any, U any](o Option[T], optb Option[U]) Option[U] {
	if o.IsSome() {
		return optb
	}
	return None[U]()
}

// AndThen returns [`None`] if the option is [`None`], otherwise calls `f` with the wrapped value.
//
//go:inline
func AndThen[T any, U any](o Option[T], f func(T) Option[U]) Option[U] {
	if o.IsNone() {
		return None[U]()
	}
	return f(o.UnwrapUnchecked())
}

// Contains returns `true` if the option is a [`Some`] value containing the given value.
//
//go:inline
func Contains[T comparable](o Option[T], x T) bool {
	if o.IsNone() {
		return false
	}
	return o.UnwrapUnchecked() == x
}

// Zip zips `a` with b `Option`.
//
// If `a` is `Some(s)` and `b` is `Some(o)`, this method returns `Some(Pair{A:s, B:o})`.
// Otherwise, `None` is returned.
//
//go:inline
func Zip[A any, B any](a Option[A], b Option[B]) Option[pair.Pair[A, B]] {
	if a.IsSome() && b.IsSome() {
		return Some[pair.Pair[A, B]](pair.Pair[A, B]{A: a.UnwrapUnchecked(), B: b.UnwrapUnchecked()})
	}
	return None[pair.Pair[A, B]]()
}

// ZipWith zips `value` and another `Option` with function `f`.
//
// If `value` is `Some(s)` and `other` is `Some(o)`, this method returns `Some(f(s, o))`.
// Otherwise, `None` is returned.
//
//go:inline
func ZipWith[T any, U any, R any](some Option[T], other Option[U], f func(T, U) R) Option[R] {
	if some.IsSome() && other.IsSome() {
		return Some(f(some.UnwrapUnchecked(), other.UnwrapUnchecked()))
	}
	return None[R]()
}

// Unzip unzips an option containing a `Pair` of two values.
//
//go:inline
func Unzip[T any, U any](p Option[pair.Pair[T, U]]) pair.Pair[Option[T], Option[U]] {
	if p.IsSome() {
		v := p.UnwrapUnchecked()
		return pair.Pair[Option[T], Option[U]]{A: Some[T](v.A), B: Some[U](v.B)}
	}
	return pair.Pair[Option[T], Option[U]]{A: None[T](), B: None[U]()}
}
