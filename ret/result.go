package ret

import "github.com/andeya/gust"

// Assert asserts gust.Result[T] as gust.Result[U].
func Assert[T any, U any](o gust.Result[T]) gust.Result[U] {
	if o.IsOk() {
		u, ok := any(o.Unwrap()).(U)
		if ok {
			return gust.Ok[U](u)
		}
		return gust.FmtErr[U]("type assert error, got %T, want %T", o.Unwrap(), u)
	}
	return gust.Err[U](o.UnwrapErr())
}

// XAssert asserts gust.Result[any] as gust.Result[U].
func XAssert[U any](o gust.Result[any]) gust.Result[U] {
	if o.IsOk() {
		u, ok := o.Unwrap().(U)
		if ok {
			return gust.Ok[U](u)
		}
		return gust.FmtErr[U]("type assert error, got %T, want %T", o.Unwrap(), u)
	}
	return gust.Err[U](o.UnwrapErr())
}

// Map maps a gust.Result[T] to gust.Result[U] by applying a function to a contained Ok value, leaving an error untouched.
// This function can be used to compose the results of two functions.
func Map[T any, U any](r gust.Result[T], f func(T) U) gust.Result[U] {
	if r.IsOk() {
		return gust.Ok[U](f(r.Unwrap()))
	}
	return gust.Err[U](r.Err())
}

// MapOr returns the provided default (if error), or applies a function to the contained value (if no error),
// Arguments passed to map_or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use MapOrElse, which is lazily evaluated.
func MapOr[T any, U any](r gust.Result[T], defaultOk U, f func(T) U) U {
	if r.IsOk() {
		return f(r.Unwrap())
	}
	return defaultOk
}

// MapOrElse maps a gust.Result[T] to U by applying fallback function default to a contained error, or function f to a contained Ok value.
// This function can be used to unpack a successful result while handling an error.
func MapOrElse[T any, U any](r gust.Result[T], defaultFn func(error) U, f func(T) U) U {
	if r.IsOk() {
		return f(r.Unwrap())
	}
	return defaultFn(r.Err())
}

// And returns r2 if the result is Ok, otherwise returns the error of r.
func And[T any, U any](r gust.Result[T], r2 gust.Result[U]) gust.Result[U] {
	if r.IsErr() {
		return gust.Err[U](r.Err())
	}
	return r2
}

// AndThen calls op if the result is Ok, otherwise returns the error of self.
// This function can be used for control flow based on gust.Result values.
func AndThen[T any, U any](r gust.Result[T], op func(T) gust.Result[U]) gust.Result[U] {
	if r.IsErr() {
		return gust.Err[U](r.Err())
	}
	return op(r.Unwrap())
}

// Contains returns true if the result is an Ok value containing the given value.
func Contains[T comparable](r gust.Result[T], x T) bool {
	if r.IsErr() {
		return false
	}
	return r.Unwrap() == x
}

// Flatten converts from gust.Result[gust.Result[T]] to gust.Result[T].
func Flatten[T any](r gust.Result[gust.Result[T]]) gust.Result[T] {
	return AndThen(r, func(rr gust.Result[T]) gust.Result[T] { return rr })
}
