package opt

import "github.com/andeya/gust"

// Assert asserts gust.Option[T] as gust.Option[U].
func Assert[T any, U any](o gust.Option[T]) gust.Option[U] {
	if o.IsSome() {
		return gust.Some[U](any(o.Unwrap()).(U))
	}
	return gust.None[U]()
}

// XAssert asserts gust.Option[any] as gust.Option[U].
func XAssert[U any](o gust.Option[any]) gust.Option[U] {
	if o.IsSome() {
		return gust.Some[U](o.Unwrap().(U))
	}
	return gust.None[U]()
}

// Map maps an `gust.Option[T]` to `gust.Option[U]` by applying a function to a contained value.
func Map[T any, U any](o gust.Option[T], f func(T) U) gust.Option[U] {
	if o.IsSome() {
		return gust.Some[U](f(o.Unwrap()))
	}
	return gust.None[U]()
}

// MapOr returns the provided default value (if none),
// or applies a function to the contained value (if any).
func MapOr[T any, U any](o gust.Option[T], defaultSome U, f func(T) U) U {
	if o.IsSome() {
		return f(o.Unwrap())
	}
	return defaultSome
}

// MapOrElse computes a default function value (if none), or
// applies a different function to the contained value (if any).
func MapOrElse[T any, U any](o gust.Option[T], defaultFn func() U, f func(T) U) U {
	if o.IsSome() {
		return f(o.Unwrap())
	}
	return defaultFn()
}

// And returns [`None`] if the option is [`None`], otherwise returns `optb`.
func And[T any, U any](o gust.Option[T], optb gust.Option[U]) gust.Option[U] {
	if o.IsSome() {
		return optb
	}
	return gust.None[U]()
}

// AndThen returns [`None`] if the option is [`None`], otherwise calls `f` with the
func AndThen[T any, U any](o gust.Option[T], f func(T) gust.Option[U]) gust.Option[U] {
	if o.IsNone() {
		return gust.None[U]()
	}
	return f(o.Unwrap())
}

// Contains returns `true` if the option is a [`Some`] value containing the given value.
func Contains[T comparable](o gust.Option[T], x T) bool {
	if o.IsNone() {
		return false
	}
	return o.Unwrap() == x
}

// Zip zips `a` with b `Option`.
//
// If `a` is `gust.Some(s)` and `b` is `gust.Some(o)`, this method returns `gust.Some(gust.Pair{A:s, B:o})`.
// Otherwise, `None` is returned.
func Zip[A any, B any](a gust.Option[A], b gust.Option[B]) gust.Option[gust.Pair[A, B]] {
	if a.IsSome() && b.IsSome() {
		return gust.Some[gust.Pair[A, B]](gust.Pair[A, B]{A: a.Unwrap(), B: b.Unwrap()})
	}
	return gust.None[gust.Pair[A, B]]()
}

// ZipWith zips `value` and another `gust.Option` with function `f`.
//
// If `value` is `Some(s)` and `other` is `Some(o)`, this method returns `Some(f(s, o))`.
// Otherwise, `None` is returned.
func ZipWith[T any, U any, R any](some gust.Option[T], other gust.Option[U], f func(T, U) R) gust.Option[R] {
	if some.IsSome() && other.IsSome() {
		return gust.Some(f(some.Unwrap(), other.Unwrap()))
	}
	return gust.None[R]()
}

// Unzip unzips an option containing a `Pair` of two values.
func Unzip[T any, U any](p gust.Option[gust.Pair[T, U]]) gust.Pair[gust.Option[T], gust.Option[U]] {
	if p.IsSome() {
		v := p.Unwrap()
		return gust.Pair[gust.Option[T], gust.Option[U]]{A: gust.Some[T](v.A), B: gust.Some[U](v.B)}
	}
	return gust.Pair[gust.Option[T], gust.Option[U]]{A: gust.None[T](), B: gust.None[U]()}
}
