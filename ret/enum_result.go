package ret

import "github.com/andeya/gust"

// EnumMap maps a gust.EnumResult[T,E] to gust.EnumResult[U,E] by applying a function to a contained Ok value, leaving an error untouched.
// This function can be used to compose the results of two functions.
func EnumMap[T any, U any, E any](r gust.EnumResult[T, E], f func(T) U) gust.EnumResult[U, E] {
	if r.IsOk() {
		return gust.EnumOk[U, E](f(r.Unwrap()))
	}
	return gust.EnumErr[U, E](r.UnwrapErr())
}

// EnumMapErr maps a EnumResult[T,E] to EnumResult[T,F] by applying a function to a contained E, leaving an T value untouched.
// This function can be used to pass through a successful result while handling an error.
func EnumMapErr[T any, E any, F any](r gust.EnumResult[T, E], op func(E) F) gust.EnumResult[T, F] {
	if r.IsErr() {
		return gust.EnumErr[T, F](op(r.UnwrapErr()))
	}
	return gust.EnumOk[T, F](r.Unwrap())
}

// EnumMapOr returns the provided default (if error), or applies a function to the contained value (if no error),
// Arguments passed to map_or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use MapOrElse, which is lazily evaluated.
func EnumMapOr[T any, U any, E any](r gust.EnumResult[T, E], defaultOk U, f func(T) U) U {
	if r.IsOk() {
		return f(r.Unwrap())
	}
	return defaultOk
}

// EnumMapOrElse maps a gust.EnumResult[T,E] to U by applying fallback function default to a contained error, or function f to a contained Ok value.
// This function can be used to unpack a successful result while handling an error.
func EnumMapOrElse[T any, U any, E any](r gust.EnumResult[T, E], defaultFn func(E) U, f func(T) U) U {
	if r.IsOk() {
		return f(r.Unwrap())
	}
	return defaultFn(r.UnwrapErr())
}

// EnumAnd returns r2 if the result is Ok, otherwise returns the error of r.
func EnumAnd[T any, U any, E any](r gust.EnumResult[T, E], r2 gust.EnumResult[U, E]) gust.EnumResult[U, E] {
	if r.IsErr() {
		return gust.EnumErr[U, E](r.UnwrapErr())
	}
	return r2
}

// EnumAndThen calls op if the result is Ok, otherwise returns the error of self.
// This function can be used for control flow based on gust.EnumResult values.
func EnumAndThen[T any, U any, E any](r gust.EnumResult[T, E], op func(T) gust.EnumResult[U, E]) gust.EnumResult[U, E] {
	if r.IsErr() {
		return gust.EnumErr[U, E](r.UnwrapErr())
	}
	return op(r.Unwrap())
}

// EnumFlatten converts from gust.EnumResult[gust.EnumResult[T,E]] to gust.EnumResult[T,E].
func EnumFlatten[T any, E any](r gust.EnumResult[gust.EnumResult[T, E], E]) gust.EnumResult[T, E] {
	return EnumAndThen(r, func(rr gust.EnumResult[T, E]) gust.EnumResult[T, E] { return rr })
}
