package ret_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/ret"
	"github.com/stretchr/testify/assert"
)

func TestAssert(t *testing.T) {
	r := gust.Ok[any]("hello")
	assert.Equal(t, "gust.Result[string]", fmt.Sprintf("%T", ret.Assert[any, string](r)))
}

func TestXAssert(t *testing.T) {
	r := gust.Ok[any]("hello")
	assert.Equal(t, "gust.Result[string]", fmt.Sprintf("%T", ret.XAssert[string](r)))
}

func TestResult_ContainsErr(t *testing.T) {
	assert.False(t, gust.Ok(2).ContainsErr("Some error message"))
	assert.True(t, ret.Contains(gust.Ok(2), 2))
	assert.False(t, ret.Contains(gust.Ok(3), 2))
	assert.False(t, ret.Contains(gust.Err[int]("Some error message"), 2))
}

func TestResult_Contains(t *testing.T) {
	assert.True(t, ret.Contains(gust.Ok(2), 2))
	assert.False(t, ret.Contains(gust.Ok(3), 2))
	assert.False(t, ret.Contains(gust.Err[int]("Some error message"), 2))
}

func TestResult_MapOrElse_1(t *testing.T) {
	var k = 21
	{
		var x = gust.Ok("foo")
		assert.Equal(t, 3, ret.MapOrElse(x, func(err error) int {
			return k * 2
		}, func(x string) int { return len(x) }))
	}
	{
		var x = gust.Err[string]("bar")
		assert.Equal(t, 42, ret.MapOrElse(x, func(err error) int {
			return k * 2
		}, func(x string) int { return len(x) }))
	}
}

func TestResult_MapOr_1(t *testing.T) {
	{
		var x = gust.Ok("foo")
		assert.Equal(t, 3, ret.MapOr(x, 42, func(x string) int { return len(x) }))
	}
	{
		var x = gust.Err[string]("foo")
		assert.Equal(t, 42, ret.MapOr(x, 42, func(x string) int { return len(x) }))
	}
}

func TestResult_Map_3(t *testing.T) {
	var isMyNum = func(s string, search int) gust.Result[bool] {
		return ret.Map(gust.Ret(strconv.Atoi(s)), func(x int) bool { return x == search })
	}
	assert.Equal(t, gust.Ok[bool](true), isMyNum("1", 1))
	assert.Equal(t, "Err(strconv.Atoi: parsing \"lol\": invalid syntax)", isMyNum("lol", 1).String())
	assert.Equal(t, "Err(strconv.Atoi: parsing \"NaN\": invalid syntax)", isMyNum("NaN", 1).String())
}

func TestAnd(t *testing.T) {
	// Test with Ok and Ok
	r1 := gust.Ok(2)
	r2 := gust.Ok(3)
	result1 := ret.And(r1, r2)
	assert.True(t, result1.IsOk())
	assert.Equal(t, 3, result1.Unwrap())

	// Test with Ok and Err
	r3 := gust.Ok(2)
	r4 := gust.Err[int]("error")
	result2 := ret.And(r3, r4)
	assert.True(t, result2.IsErr())

	// Test with Err and Ok
	r5 := gust.Err[int]("error")
	r6 := gust.Ok(3)
	result3 := ret.And(r5, r6)
	assert.True(t, result3.IsErr())
	assert.Equal(t, "error", result3.Err().Error())

	// Test with Err and Err
	r7 := gust.Err[int]("error1")
	r8 := gust.Err[int]("error2")
	result4 := ret.And(r7, r8)
	assert.True(t, result4.IsErr())
	assert.Equal(t, "error1", result4.Err().Error())
}

func TestAndThen(t *testing.T) {
	// Test with Ok and successful operation
	r1 := gust.Ok(2)
	result1 := ret.AndThen(r1, func(x int) gust.Result[int] {
		return gust.Ok(x * 2)
	})
	assert.True(t, result1.IsOk())
	assert.Equal(t, 4, result1.Unwrap())

	// Test with Ok and error operation
	r2 := gust.Ok(2)
	result2 := ret.AndThen(r2, func(x int) gust.Result[int] {
		return gust.Err[int]("error")
	})
	assert.True(t, result2.IsErr())

	// Test with Err (should not call function)
	r3 := gust.Err[int]("error")
	result3 := ret.AndThen(r3, func(x int) gust.Result[int] {
		return gust.Ok(x * 2)
	})
	assert.True(t, result3.IsErr())
	assert.Equal(t, "error", result3.Err().Error())
}

func TestFlatten(t *testing.T) {
	// Test with Ok(Ok(value)) - matches documentation example
	{
		r1 := gust.Ok(gust.Ok(1))
		result1 := ret.Flatten(r1)
		assert.Equal(t, gust.Ok[int](1), result1)
	}
	// Test with Ok(Err(error)) - matches documentation example
	{
		r2 := gust.Ok(gust.Err[int](errors.New("error")))
		result2 := ret.Flatten(r2)
		assert.Equal(t, "error", result2.Err().Error())
	}
	// Test with Err - matches documentation example
	{
		r3 := gust.Err[gust.Result[int]](errors.New("error"))
		result3 := ret.Flatten(r3)
		assert.Equal(t, "error", result3.Err().Error())
	}
	// Additional test with Ok(Ok(value)) using different value
	{
		r1 := gust.Ok(gust.Ok(42))
		result1 := ret.Flatten(r1)
		assert.True(t, result1.IsOk())
		assert.Equal(t, 42, result1.Unwrap())
	}
	// Additional test with Ok(Err(error)) using string error
	{
		r2 := gust.Ok(gust.Err[int]("error"))
		result2 := ret.Flatten(r2)
		assert.True(t, result2.IsErr())
		assert.Equal(t, "error", result2.Err().Error())
	}
	// Additional test with Err using string error
	{
		r3 := gust.Err[gust.Result[int]]("outer error")
		result3 := ret.Flatten(r3)
		assert.True(t, result3.IsErr())
		assert.Equal(t, "outer error", result3.Err().Error())
	}
}

func TestFlatten2(t *testing.T) {
	// Test with Ok result and nil error - matches documentation example
	{
		r1 := gust.Ok(1)
		var err1 error = nil
		result1 := ret.Flatten2(r1, err1)
		assert.Equal(t, gust.Ok[int](1), result1)
	}
	// Test with Ok result and error - matches documentation example
	{
		r2 := gust.Ok(1)
		err2 := errors.New("error")
		result2 := ret.Flatten2(r2, err2)
		assert.Equal(t, "error", result2.Err().Error())
	}
	// Test with Err result and nil error - matches documentation example
	{
		r3 := gust.Err[int](errors.New("error"))
		var err3 error = nil
		result3 := ret.Flatten2(r3, err3)
		assert.Equal(t, "error", result3.Err().Error())
	}
	// Additional test with Ok result and nil error using different value
	{
		r := gust.Ok(42)
		result := ret.Flatten2(r, nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, 42, result.Unwrap())
	}
	// Additional test with Ok result and error
	{
		r := gust.Ok(42)
		result := ret.Flatten2(r, errors.New("test error"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "test error", result.Err().Error())
	}
	// Additional test with Err result and nil error
	{
		r := gust.Err[int]("original error")
		result := ret.Flatten2(r, nil)
		assert.True(t, result.IsErr())
		assert.Equal(t, "original error", result.Err().Error())
	}
	// Additional test with Err result and error (error parameter takes precedence)
	{
		r := gust.Err[int]("original error")
		result := ret.Flatten2(r, errors.New("new error"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "new error", result.Err().Error())
	}
}

func TestAssert_TypeAssertionError(t *testing.T) {
	// Test Assert with type assertion error
	r1 := gust.Ok[any](42)
	result1 := ret.Assert[any, string](r1)
	assert.True(t, result1.IsErr())
	assert.Contains(t, result1.Err().Error(), "type assert error")
}

func TestXAssert_TypeAssertionError(t *testing.T) {
	// Test XAssert with type assertion error
	r1 := gust.Ok[any](42)
	result1 := ret.XAssert[string](r1)
	assert.True(t, result1.IsErr())
	assert.Contains(t, result1.Err().Error(), "type assert error")
}

func TestAssert2(t *testing.T) {
	// Test with nil error and successful type assertion
	{
		result := ret.Assert2[int, int](42, nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, 42, result.Unwrap())
	}
	// Test with error
	{
		result := ret.Assert2[int, int](42, errors.New("test error"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "test error", result.Err().Error())
	}
	// Test with nil error and failed type assertion
	{
		result := ret.Assert2[string, int]("hello", nil)
		assert.True(t, result.IsErr())
		assert.Contains(t, result.Err().Error(), "type assert error")
	}
	// Test with successful type assertion from any
	{
		var v any = 42
		result := ret.Assert2[any, int](v, nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, 42, result.Unwrap())
	}
}

func TestXAssert2(t *testing.T) {
	// Test with nil error and successful type assertion
	{
		result := ret.XAssert2[int](42, nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, 42, result.Unwrap())
	}
	// Test with error
	{
		result := ret.XAssert2[int](42, errors.New("test error"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "test error", result.Err().Error())
	}
	// Test with nil error and failed type assertion
	{
		result := ret.XAssert2[int]("hello", nil)
		assert.True(t, result.IsErr())
		assert.Contains(t, result.Err().Error(), "type assert error")
	}
	// Test with string type
	{
		result := ret.XAssert2[string]("hello", nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, "hello", result.Unwrap())
	}
}

func TestMap2(t *testing.T) {
	// Test with nil error and successful mapping
	{
		result := ret.Map2(2, nil, func(x int) int {
			return x * 2
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, 4, result.Unwrap())
	}
	// Test with error
	{
		result := ret.Map2(2, errors.New("test error"), func(x int) int {
			return x * 2
		})
		assert.True(t, result.IsErr())
		assert.Equal(t, "test error", result.Err().Error())
	}
	// Test with type conversion
	{
		result := ret.Map2(2, nil, func(x int) string {
			return fmt.Sprintf("%d", x)
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, "2", result.Unwrap())
	}
}

func TestMapOr2(t *testing.T) {
	// Test with nil error and successful mapping
	{
		result := ret.MapOr2(2, nil, 0, func(x int) int {
			return x * 2
		})
		assert.Equal(t, 4, result)
	}
	// Test with error (should return default)
	{
		result := ret.MapOr2(2, errors.New("test error"), 0, func(x int) int {
			return x * 2
		})
		assert.Equal(t, 0, result)
	}
	// Test with different default value
	{
		result := ret.MapOr2(2, errors.New("test error"), 100, func(x int) int {
			return x * 2
		})
		assert.Equal(t, 100, result)
	}
	// Test with type conversion
	{
		result := ret.MapOr2(2, nil, "default", func(x int) string {
			return fmt.Sprintf("%d", x)
		})
		assert.Equal(t, "2", result)
	}
}

func TestMapOrElse2(t *testing.T) {
	// Test with nil error and successful mapping
	{
		result := ret.MapOrElse2(2, nil, func(err error) int {
			return 0
		}, func(x int) int {
			return x * 2
		})
		assert.Equal(t, 4, result)
	}
	// Test with error (should call default function)
	{
		result := ret.MapOrElse2(2, errors.New("test error"), func(err error) int {
			return 100
		}, func(x int) int {
			return x * 2
		})
		assert.Equal(t, 100, result)
	}
	// Test with error message in default function
	{
		result := ret.MapOrElse2(2, errors.New("test error"), func(err error) string {
			return err.Error()
		}, func(x int) string {
			return fmt.Sprintf("%d", x)
		})
		assert.Equal(t, "test error", result)
	}
}

func TestAnd2(t *testing.T) {
	// Test with nil errors (both Ok)
	{
		result := ret.And2(1, nil, 2, nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, 2, result.Unwrap())
	}
	// Test with first error
	{
		result := ret.And2(1, errors.New("error1"), 2, nil)
		assert.True(t, result.IsErr())
		assert.Equal(t, "error1", result.Err().Error())
	}
	// Test with second error (first is Ok)
	{
		result := ret.And2(1, nil, 2, errors.New("error2"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "error2", result.Err().Error())
	}
	// Test with both errors (should return first error)
	{
		result := ret.And2(1, errors.New("error1"), 2, errors.New("error2"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "error1", result.Err().Error())
	}
	// Test with type conversion
	{
		result := ret.And2(1, nil, "hello", nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, "hello", result.Unwrap())
	}
}

func TestAndThen2(t *testing.T) {
	// Test with Ok result and successful operation
	{
		r := gust.Ok(2)
		result := ret.AndThen2(r, func(x int) (int, error) {
			return x * 2, nil
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, 4, result.Unwrap())
	}
	// Test with Ok result and error operation
	{
		r := gust.Ok(2)
		result := ret.AndThen2(r, func(x int) (int, error) {
			return 0, errors.New("operation error")
		})
		assert.True(t, result.IsErr())
		assert.Equal(t, "operation error", result.Err().Error())
	}
	// Test with Err result (should return original error)
	{
		r := gust.Err[int]("early error")
		result := ret.AndThen2(r, func(x int) (int, error) {
			return x * 2, nil
		})
		assert.True(t, result.IsErr())
		assert.Equal(t, "early error", result.Err().Error())
	}
	// Test with type conversion
	{
		r := gust.Ok(2)
		result := ret.AndThen2(r, func(x int) (string, error) {
			return fmt.Sprintf("%d", x), nil
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, "2", result.Unwrap())
	}
}

func TestAndThen3(t *testing.T) {
	// Test with nil error and successful operation
	{
		result := ret.AndThen3(2, nil, func(x int) (int, error) {
			return x * 2, nil
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, 4, result.Unwrap())
	}
	// Test with error (should return error)
	{
		result := ret.AndThen3(2, errors.New("early error"), func(x int) (int, error) {
			return x * 2, nil
		})
		assert.True(t, result.IsErr())
		assert.Equal(t, "early error", result.Err().Error())
	}
	// Test with nil error and error operation
	{
		result := ret.AndThen3(2, nil, func(x int) (int, error) {
			return 0, errors.New("operation error")
		})
		assert.True(t, result.IsErr())
		assert.Equal(t, "operation error", result.Err().Error())
	}
	// Test with type conversion
	{
		result := ret.AndThen3(2, nil, func(x int) (string, error) {
			return fmt.Sprintf("%d", x), nil
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, "2", result.Unwrap())
	}
}

func TestContains2(t *testing.T) {
	// Test with nil error and matching value
	{
		result := ret.Contains2(2, nil, 2)
		assert.True(t, result)
	}
	// Test with nil error and non-matching value
	{
		result := ret.Contains2(2, nil, 3)
		assert.False(t, result)
	}
	// Test with error (should return false)
	{
		result := ret.Contains2(2, errors.New("test error"), 2)
		assert.False(t, result)
	}
	// Test with string type
	{
		result := ret.Contains2("hello", nil, "hello")
		assert.True(t, result)
		result2 := ret.Contains2("hello", nil, "world")
		assert.False(t, result2)
	}
}
