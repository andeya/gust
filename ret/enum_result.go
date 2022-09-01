package ret

import "github.com/andeya/gust"

// EnumAssert asserts gust.EnumResult[T,E] as gust.EnumResult[U,F].
func EnumAssert[T any, E any, U any, F any](o gust.EnumResult[T, E]) gust.EnumResult[U, F] {
	if o.IsOk() {
		return gust.EnumOk[U, F](any(o.Unwrap()).(U))
	}
	return gust.EnumErr[U, F](any(o.UnwrapErr()).(F))
}

// EnumXOkAssert asserts gust.EnumResult[any, E] as gust.EnumResult[U, E].
func EnumXOkAssert[T any, E any, U any](o gust.EnumResult[any, E]) gust.EnumResult[U, E] {
	if o.IsOk() {
		return gust.EnumOk[U, E](o.Unwrap().(U))
	}
	return gust.EnumErr[U, E](o.UnwrapErr())
}

// EnumXErrAssert asserts gust.EnumResult[T, any] as gust.EnumResult[T, E].
func EnumXErrAssert[T any, E any](o gust.EnumResult[T, any]) gust.EnumResult[T, E] {
	if o.IsOk() {
		return gust.EnumOk[T, E](o.Unwrap())
	}
	return gust.EnumErr[T, E](o.UnwrapErr().(E))
}

// EnumXAssert asserts gust.EnumResult[any, any] as gust.EnumResult[T, E].
func EnumXAssert[T any, E any](o gust.EnumResult[any, any]) gust.EnumResult[T, E] {
	if o.IsOk() {
		return gust.EnumOk[T, E](o.Unwrap().(T))
	}
	return gust.EnumErr[T, E](o.UnwrapErr().(E))
}

// EnumMap maps a gust.EnumResult[T,E] to gust.EnumResult[U,E] by applying a function to a contained Ok value, leaving an error untouched.
// This function can be used to compose the results of two functions.
func EnumMap[T any, U any, E any](r gust.EnumResult[T, E], f func(T) U) gust.EnumResult[U, E] {
	if r.IsOk() {
		return gust.EnumOk[U, E](f(r.Unwrap()))
	}
	return gust.EnumErr[U, E](r.UnwrapErr())
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

// EnumMapErr maps a EnumResult[T,E] to EnumResult[T,F] by applying a function to a contained E, leaving an T value untouched.
// This function can be used to pass through a successful result while handling an error.
func EnumMapErr[T any, E any, F any](r gust.EnumResult[T, E], op func(E) F) gust.EnumResult[T, F] {
	if r.IsErr() {
		return gust.EnumErr[T, F](op(r.UnwrapErr()))
	}
	return gust.EnumOk[T, F](r.Unwrap())
}

// EnumAnd returns `b` if the `a` is Ok, otherwise returns the error of `a`.
func EnumAnd[T any, U any, E any](a gust.EnumResult[T, E], b gust.EnumResult[U, E]) gust.EnumResult[U, E] {
	if a.IsErr() {
		return gust.EnumErr[U, E](a.UnwrapErr())
	}
	return b
}

// EnumAndThen calls op if the result is Ok, otherwise returns the error of self.
// This function can be used for control flow based on gust.EnumResult values.
func EnumAndThen[T any, U any, E any](r gust.EnumResult[T, E], op func(T) gust.EnumResult[U, E]) gust.EnumResult[U, E] {
	if r.IsErr() {
		return gust.EnumErr[U, E](r.UnwrapErr())
	}
	return op(r.Unwrap())
}

// EnumOr returns `b` if `a` is E, otherwise returns the T value of `a`.
// Arguments passed to or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use EnumOrElse, which is lazily evaluated.
func EnumOr[T any, E any, F any](a gust.EnumResult[T, E], b gust.EnumResult[T, F]) gust.EnumResult[T, F] {
	if a.IsErr() {
		return b
	}
	return gust.EnumOk[T, F](a.Unwrap())
}

// EnumOrElse calls op if the result is E, otherwise returns the T value of result.
// This function can be used for control flow based on result values.
func EnumOrElse[T any, E any, F any](result gust.EnumResult[T, E], op func(E) gust.EnumResult[T, F]) gust.EnumResult[T, F] {
	if result.IsErr() {
		return op(result.UnwrapErr())
	}
	return gust.EnumOk[T, F](result.Unwrap())
}

// EnumFlatten converts from gust.EnumResult[gust.EnumResult[T,E]] to gust.EnumResult[T,E].
func EnumFlatten[T any, E any](r gust.EnumResult[gust.EnumResult[T, E], E]) gust.EnumResult[T, E] {
	return EnumAndThen(r, func(rr gust.EnumResult[T, E]) gust.EnumResult[T, E] { return rr })
}
