// Package ret provides helper functions for working with gust.Result types.
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

// Assert2 asserts a value and an error as a gust.Result[U].
//
// # Examples
//
//	var v = 1
//	var err = nil
//	var result = Assert2(v, err)
//	assert.Equal(t, gust.Ok[int](1), result)
func Assert2[T any, U any](v T, err error) gust.Result[U] {
	if err != nil {
		return gust.Err[U](err)
	}
	u, ok := any(v).(U)
	if ok {
		return gust.Ok[U](u)
	}
	return gust.FmtErr[U]("type assert error, got %T, want %T", v, u)
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

// XAssert2 asserts a value and an error as a gust.Result[U].
//
// # Examples
//
//	var v = 1
//	var err = nil
//	var result = XAssert2(v, err)
//	assert.Equal(t, gust.Ok[int](1), result)
func XAssert2[U any](v any, err error) gust.Result[U] {
	if err != nil {
		return gust.Err[U](err)
	}
	u, ok := v.(U)
	if ok {
		return gust.Ok[U](u)
	}
	return gust.FmtErr[U]("type assert error, got %T, want %T", v, u)
}

// Map maps a gust.Result[T] to gust.Result[U] by applying a function to a contained Ok value, leaving an error untouched.
// This function can be used to compose the results of two functions.
func Map[T any, U any](r gust.Result[T], f func(T) U) gust.Result[U] {
	if r.IsOk() {
		return gust.Ok[U](f(r.Unwrap()))
	}
	return gust.Err[U](r.Err())
}

// Map2 maps a value and an error as a gust.Result[U] by applying a function to the value.
//
// # Examples
//
//	var v = 1
//	var err = nil
//	var result = Map2(v, err, func(v int) int { return v * 2 })
//	assert.Equal(t, gust.Ok[int](2), result)
//
//go:inline
func Map2[T any, U any](v T, err error, f func(T) U) gust.Result[U] {
	if err != nil {
		return gust.Err[U](err)
	}
	return gust.Ok[U](f(v))
}

// MapOr returns the provided default (if error), or applies a function to the contained value (if no error),
// Arguments passed to map_or are eagerly evaluated; if you are passing the result of a function call, it is recommended to use MapOrElse, which is lazily evaluated.
func MapOr[T any, U any](r gust.Result[T], defaultOk U, f func(T) U) U {
	if r.IsOk() {
		return f(r.Unwrap())
	}
	return defaultOk
}

// MapOr2 maps a value and an error as a gust.Result[U] by applying a function to the value.
//
// # Examples
//
//	var v = 1
//	var err = nil
//	var result = MapOr2(v, err, func(v int) int { return v * 2 })
//	assert.Equal(t, gust.Ok[int](2), result)
//
//go:inline
func MapOr2[T any, U any](v T, err error, defaultOk U, f func(T) U) U {
	if err != nil {
		return defaultOk
	}
	return f(v)
}

// MapOrElse maps a gust.Result[T] to U by applying fallback function default to a contained error, or function f to a contained Ok value.
// This function can be used to unpack a successful result while handling an error.
func MapOrElse[T any, U any](r gust.Result[T], defaultFn func(error) U, f func(T) U) U {
	if r.IsOk() {
		return f(r.Unwrap())
	}
	return defaultFn(r.Err())
}

// MapOrElse2 maps a value and an error as a gust.Result[U] by applying a function to the value.
//
// # Examples
//
//	var v = 1
//	var err = nil
//	var result = MapOrElse2(v, err, func(err error) int { return 0 }, func(v int) int { return v * 2 })
//	assert.Equal(t, gust.Ok[int](2), result)
//
//go:inline
func MapOrElse2[T any, U any](v T, err error, defaultFn func(error) U, f func(T) U) U {
	if err != nil {
		return defaultFn(err)
	}
	return f(v)
}

// And returns `r1` if `r1` is Err, otherwise returns `r2`.
//
//go:inline
func And[T any, U any](r1 gust.Result[T], r2 gust.Result[U]) gust.Result[U] {
	if r1.IsErr() {
		return gust.Err[U](r1.Err())
	}
	return r2
}

// And2 returns `Ret(v1, err1)` if `r1` is `Err`, otherwise returns `Ret(v2, err2)`.
//
// # Examples
//
//	var v1 = 1
//	var err1 = nil
//	var v2 = 2
//	var err2 = nil
//	var result = And2(v1, err1, v2, err2)
//	assert.Equal(t, gust.Ok[int](2), result)
//
//	var v1 = 1
//	var err1 = errors.New("error1")
//	var v2 = 2
//	var err2 = nil
//	var result = And2(v1, err1, v2, err2)
//	assert.Equal(t, "error1", result.Err().Error())
//
//	var v1 = 1
//	var err1 = nil
//	var v2 = 2
//	var err2 = errors.New("error2")
//	var result = And2(v1, err1, v2, err2)
//	assert.Equal(t, "error2", result.Err().Error())
//
//	var v1 = 1
//	var err1 = errors.New("error1")
//	var v2 = 2
//	var err2 = errors.New("error2")
//	var result = And2(v1, err1, v2, err2)
//	assert.Equal(t, "error1", result.Err().Error())
//
//go:inline
func And2[T any, U any](v1 T, err1 error, v2 U, err2 error) gust.Result[U] {
	if err1 != nil {
		return gust.Err[U](err1)
	}
	return gust.Ret(v2, err2)
}

// AndThen calls op if the result is Ok, otherwise returns the error of self.
// This function can be used for control flow based on gust.Result values.
func AndThen[T any, U any](r gust.Result[T], op func(T) gust.Result[U]) gust.Result[U] {
	if r.IsErr() {
		return gust.Err[U](r.Err())
	}
	return op(r.Unwrap())
}

// AndThen2 calls op if the result is Ok, otherwise returns the error of self.
// This function can be used for control flow based on gust.Result values.
//
//go:inline
func AndThen2[T any, U any](r gust.Result[T], op func(T) (U, error)) gust.Result[U] {
	if r.IsErr() {
		return gust.Err[U](r.Err())
	}
	return gust.Ret[U](op(r.Unwrap()))
}

// AndThen3 calls op if the result is Ok, otherwise returns the error of self.
// This function can be used for control flow based on gust.Result values.
//
//go:inline
func AndThen3[T any, U any](v T, err error, op func(T) (U, error)) gust.Result[U] {
	if err != nil {
		return gust.Err[U](err)
	}
	return gust.Ret[U](op(v))
}

// Contains returns true if the result is an Ok value containing the given value.
//
//go:inline
func Contains[T comparable](r gust.Result[T], x T) bool {
	if r.IsErr() {
		return false
	}
	return r.Unwrap() == x
}

// Contains2 returns true if the result is an Ok value containing the given value.
//
//go:inline
func Contains2[T comparable](v T, err error, x T) bool {
	if err != nil {
		return false
	}
	return v == x
}

// Flatten converts from gust.Result[gust.Result[T]] to gust.Result[T].
//
// # Examples
//
//	var r1 = gust.Ok(gust.Ok(1))
//	var result1 = Flatten(r1)
//	assert.Equal(t, gust.Ok[int](1), result1)
//	var r2 = gust.Ok(gust.Err(errors.New("error")))
//	var result2 = Flatten(r2)
//	assert.Equal(t, "error", result2.Err().Error())
//	var r3 = gust.Err[gust.Result[int]](errors.New("error"))
//	var result3 = Flatten(r3)
//	assert.Equal(t, "error", result3.Err().Error())
//
//go:inline
func Flatten[T any](r gust.Result[gust.Result[T]]) gust.Result[T] {
	if r.IsErr() {
		return gust.Err[T](r.Err())
	}
	return r.Unwrap()
}

// Flatten2 converts from `(gust.Result[T], error)` to gust.Result[T].
//
// # Examples
//
//	var r1 = gust.Ok(1)
//	var err1 = nil
//	var result1 = Flatten2(r1, err1)
//	assert.Equal(t, gust.Ok[int](1), result1)
//	var r2 = gust.Ok(1)
//	var err2 = errors.New("error")
//	var result2 = Flatten2(r2, err2)
//	assert.Equal(t, "error", result2.Err().Error())
//	var r3 = gust.Err(errors.New("error"))
//	var err3 = nil
//	var result3 = Flatten2(r3, err3)
//	assert.Equal(t, "error", result3.Err().Error())
//
//go:inline
func Flatten2[T any](r gust.Result[T], err error) gust.Result[T] {
	if err != nil {
		return gust.Err[T](err)
	}
	return r
}
