package optnil

import "github.com/henrylee2cn/gust"

// Map maps an `gust.OptNil[T]` to `gust.OptNil[U]` by applying a function to a contained value.
func Map[T any, U any](o gust.OptNil[T], f func(*T) *U) gust.OptNil[U] {
	if o.NotNil() {
		return gust.Ptr[U](f(o.Unwrap()))
	}
	return gust.Nil[U]()
}

// MapOr returns the provided default value (if none),
// or applies a function to the contained value (if any).
func MapOr[T any, U any](o gust.OptNil[T], defaultPtr *U, f func(*T) *U) *U {
	if o.NotNil() {
		return f(o.Unwrap())
	}
	return defaultPtr
}

// MapOrElse computes a default function value (if none), or
// applies a different function to the contained value (if any).
func MapOrElse[T any, U any](o gust.OptNil[T], defaultFn func() *U, f func(*T) *U) *U {
	if o.NotNil() {
		return f(o.Unwrap())
	}
	return defaultFn()
}

// And returns [`Nil`] if the option is [`Nil`], otherwise returns `optb`.
func And[T any, U any](o gust.OptNil[T], optb gust.OptNil[U]) gust.OptNil[U] {
	if o.NotNil() {
		return optb
	}
	return gust.Nil[U]()
}

// AndThen returns [`Nil`] if the option is [`Nil`], otherwise calls `f` with the
func AndThen[T any, U any](o gust.OptNil[T], f func(*T) gust.OptNil[U]) gust.OptNil[U] {
	if o.IsNil() {
		return gust.Nil[U]()
	}
	return f(o.Unwrap())
}

// Contains returns `true` if the option is a [`NonNil`] value containing the given value.
func Contains[T comparable](o gust.OptNil[T], x *T) bool {
	if o.IsNil() {
		return false
	}
	return o.Unwrap() == x
}

// ZipWith zips `value` and another `gust.OptNil` with function `f`.
//
// If `value` is `Ptr(s)` and `other` is `Ptr(o)`, this method returns `Ptr(f(s, o))`.
// Otherwise, `Nil` is returned.
func ZipWith[T any, U any, R any](some gust.OptNil[T], other gust.OptNil[U], f func(*T, *U) *R) gust.OptNil[R] {
	if some.NotNil() && other.NotNil() {
		return gust.Ptr(f(some.Unwrap(), other.Unwrap()))
	}
	return gust.Nil[R]()
}
