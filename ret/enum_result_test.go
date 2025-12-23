package ret_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/ret"
	"github.com/stretchr/testify/assert"
)

func TestEnumUnsafeAssert(t *testing.T) {
	// Test with Ok - same type (no conversion needed)
	r1 := gust.EnumOk[int, string](42)
	result1 := ret.EnumUnsafeAssert[int, string, int, string](r1)
	assert.True(t, result1.IsOk())
	assert.Equal(t, 42, result1.Unwrap())

	// Test with Err
	r2 := gust.EnumErr[int, string]("error")
	result2 := ret.EnumUnsafeAssert[int, string, int, string](r2)
	assert.True(t, result2.IsErr())
	assert.Equal(t, "error", result2.UnwrapErr())

	// Test with compatible types - use actual compatible conversion
	// int can be asserted to int (same), but not to int64 directly
	// We need to test with types that can actually be asserted
	r3 := gust.EnumOk[any, string](42)
	result3 := ret.EnumUnsafeAssert[any, string, int, string](r3)
	assert.True(t, result3.IsOk())
	assert.Equal(t, 42, result3.Unwrap())
}

func TestEnumXOkUnsafeAssert(t *testing.T) {
	// Test with Ok
	r1 := gust.EnumOk[any, string](42)
	result1 := ret.EnumXOkUnsafeAssert[int, string, int](r1)
	assert.True(t, result1.IsOk())
	assert.Equal(t, 42, result1.Unwrap())

	// Test with Err
	r2 := gust.EnumErr[any, string]("error")
	result2 := ret.EnumXOkUnsafeAssert[int, string, int](r2)
	assert.True(t, result2.IsErr())
	assert.Equal(t, "error", result2.UnwrapErr())
}

func TestEnumXErrUnsafeAssert(t *testing.T) {
	// Test with Ok
	r1 := gust.EnumOk[int, any](42)
	result1 := ret.EnumXErrUnsafeAssert[int, string](r1)
	assert.True(t, result1.IsOk())
	assert.Equal(t, 42, result1.Unwrap())

	// Test with Err
	r2 := gust.EnumErr[int, any]("error")
	result2 := ret.EnumXErrUnsafeAssert[int, string](r2)
	assert.True(t, result2.IsErr())
	assert.Equal(t, "error", result2.UnwrapErr())
}

func TestEnumXUnsafeAssert(t *testing.T) {
	// Test with Ok
	r1 := gust.EnumOk[any, any](42)
	result1 := ret.EnumXUnsafeAssert[int, string](r1)
	assert.True(t, result1.IsOk())
	assert.Equal(t, 42, result1.Unwrap())

	// Test with Err
	r2 := gust.EnumErr[any, any]("error")
	result2 := ret.EnumXUnsafeAssert[int, string](r2)
	assert.True(t, result2.IsErr())
	assert.Equal(t, "error", result2.UnwrapErr())
}

func TestEnumMap(t *testing.T) {
	// Test with Ok
	r1 := gust.EnumOk[int, string](2)
	result1 := ret.EnumMap(r1, func(x int) int64 {
		return int64(x * 2)
	})
	assert.True(t, result1.IsOk())
	assert.Equal(t, int64(4), result1.Unwrap())

	// Test with Err
	r2 := gust.EnumErr[int, string]("error")
	result2 := ret.EnumMap(r2, func(x int) int64 {
		return int64(x * 2)
	})
	assert.True(t, result2.IsErr())
	assert.Equal(t, "error", result2.UnwrapErr())
}

func TestEnumMapOr(t *testing.T) {
	// Test with Ok
	r1 := gust.EnumOk[int, string](2)
	result1 := ret.EnumMapOr(r1, 0, func(x int) int {
		return x * 2
	})
	assert.Equal(t, 4, result1)

	// Test with Err
	r2 := gust.EnumErr[int, string]("error")
	result2 := ret.EnumMapOr(r2, 0, func(x int) int {
		return x * 2
	})
	assert.Equal(t, 0, result2)
}

func TestEnumMapOrElse(t *testing.T) {
	// Test with Ok
	r1 := gust.EnumOk[int, string](2)
	result1 := ret.EnumMapOrElse(r1, func(e string) int {
		return 0
	}, func(x int) int {
		return x * 2
	})
	assert.Equal(t, 4, result1)

	// Test with Err
	r2 := gust.EnumErr[int, string]("error")
	result2 := ret.EnumMapOrElse(r2, func(e string) int {
		return 10
	}, func(x int) int {
		return x * 2
	})
	assert.Equal(t, 10, result2)
}

func TestEnumMapErr(t *testing.T) {
	// Test with Ok
	r1 := gust.EnumOk[int, string](42)
	result1 := ret.EnumMapErr(r1, func(e string) int {
		return len(e)
	})
	assert.True(t, result1.IsOk())
	assert.Equal(t, 42, result1.Unwrap())

	// Test with Err
	r2 := gust.EnumErr[int, string]("error")
	result2 := ret.EnumMapErr(r2, func(e string) int {
		return len(e)
	})
	assert.True(t, result2.IsErr())
	assert.Equal(t, 5, result2.UnwrapErr())
}

func TestEnumAnd(t *testing.T) {
	// Test with Ok and Ok
	r1 := gust.EnumOk[int, string](2)
	r2 := gust.EnumOk[int64, string](3)
	result1 := ret.EnumAnd(r1, r2)
	assert.True(t, result1.IsOk())
	assert.Equal(t, int64(3), result1.Unwrap())

	// Test with Ok and Err
	r3 := gust.EnumOk[int, string](2)
	r4 := gust.EnumErr[int64, string]("error")
	result2 := ret.EnumAnd(r3, r4)
	assert.True(t, result2.IsErr())

	// Test with Err and Ok
	r5 := gust.EnumErr[int, string]("error")
	r6 := gust.EnumOk[int64, string](3)
	result3 := ret.EnumAnd(r5, r6)
	assert.True(t, result3.IsErr())
	assert.Equal(t, "error", result3.UnwrapErr())
}

func TestEnumAndThen(t *testing.T) {
	// Test with Ok and successful operation
	r1 := gust.EnumOk[int, string](2)
	result1 := ret.EnumAndThen(r1, func(x int) gust.EnumResult[int64, string] {
		return gust.EnumOk[int64, string](int64(x * 2))
	})
	assert.True(t, result1.IsOk())
	assert.Equal(t, int64(4), result1.Unwrap())

	// Test with Ok and error operation
	r2 := gust.EnumOk[int, string](2)
	result2 := ret.EnumAndThen(r2, func(x int) gust.EnumResult[int64, string] {
		return gust.EnumErr[int64, string]("error")
	})
	assert.True(t, result2.IsErr())

	// Test with Err
	r3 := gust.EnumErr[int, string]("error")
	result3 := ret.EnumAndThen(r3, func(x int) gust.EnumResult[int64, string] {
		return gust.EnumOk[int64, string](int64(x * 2))
	})
	assert.True(t, result3.IsErr())
	assert.Equal(t, "error", result3.UnwrapErr())
}

func TestEnumOr(t *testing.T) {
	// Test with Ok and Err
	r1 := gust.EnumOk[int, string](42)
	r2 := gust.EnumErr[int, int](100)
	result1 := ret.EnumOr(r1, r2)
	assert.True(t, result1.IsOk())
	assert.Equal(t, 42, result1.Unwrap())

	// Test with Err and Ok
	r3 := gust.EnumErr[int, string]("error")
	r4 := gust.EnumOk[int, int](100)
	result2 := ret.EnumOr(r3, r4)
	assert.True(t, result2.IsOk())
	assert.Equal(t, 100, result2.Unwrap())

	// Test with Err and Err
	r5 := gust.EnumErr[int, string]("error1")
	r6 := gust.EnumErr[int, int](200)
	result3 := ret.EnumOr(r5, r6)
	assert.True(t, result3.IsErr())
}

func TestEnumOrElse(t *testing.T) {
	// Test with Ok
	r1 := gust.EnumOk[int, string](42)
	result1 := ret.EnumOrElse(r1, func(e string) gust.EnumResult[int, int] {
		return gust.EnumOk[int, int](100)
	})
	assert.True(t, result1.IsOk())
	assert.Equal(t, 42, result1.Unwrap())

	// Test with Err
	r2 := gust.EnumErr[int, string]("error")
	result2 := ret.EnumOrElse(r2, func(e string) gust.EnumResult[int, int] {
		return gust.EnumOk[int, int](100)
	})
	assert.True(t, result2.IsOk())
	assert.Equal(t, 100, result2.Unwrap())
}

func TestEnumFlatten(t *testing.T) {
	// Test with Ok(Ok(value))
	r1 := gust.EnumOk[gust.EnumResult[int, string], string](gust.EnumOk[int, string](42))
	result1 := ret.EnumFlatten(r1)
	assert.True(t, result1.IsOk())
	assert.Equal(t, 42, result1.Unwrap())

	// Test with Ok(Err(error))
	r2 := gust.EnumOk[gust.EnumResult[int, string], string](gust.EnumErr[int, string]("error"))
	result2 := ret.EnumFlatten(r2)
	assert.True(t, result2.IsErr())
	assert.Equal(t, "error", result2.UnwrapErr())

	// Test with Err
	r3 := gust.EnumErr[gust.EnumResult[int, string], string]("outer error")
	result3 := ret.EnumFlatten(r3)
	assert.True(t, result3.IsErr())
	assert.Equal(t, "outer error", result3.UnwrapErr())
}
