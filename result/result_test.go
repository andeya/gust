package result_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/andeya/gust/result"
	"github.com/stretchr/testify/assert"
)

func TestAssert(t *testing.T) {
	r := result.Ok[any]("hello")
	assert.Equal(t, "core.Result[string]", fmt.Sprintf("%T", result.Assert[any, string](r)))
}

func TestXAssert(t *testing.T) {
	r := result.Ok[any]("hello")
	assert.Equal(t, "core.Result[string]", fmt.Sprintf("%T", result.XAssert[string](r)))
}

func TestResult_ContainsErr(t *testing.T) {
	assert.False(t, result.Ok(2).ContainsErr("Some error message"))
	assert.True(t, result.Contains(result.Ok(2), 2))
	assert.False(t, result.Contains(result.Ok(3), 2))
	assert.False(t, result.Contains(result.TryErr[int]("Some error message"), 2))
}

func TestResult_Contains(t *testing.T) {
	assert.True(t, result.Contains(result.Ok(2), 2))
	assert.False(t, result.Contains(result.Ok(3), 2))
	assert.False(t, result.Contains(result.TryErr[int]("Some error message"), 2))
}

func TestResult_MapOrElse_1(t *testing.T) {
	var k = 21
	{
		var x = result.Ok("foo")
		assert.Equal(t, 3, result.MapOrElse(x, func(err error) int {
			return k * 2
		}, func(x string) int { return len(x) }))
	}
	{
		var x = result.TryErr[string]("bar")
		assert.Equal(t, 42, result.MapOrElse(x, func(err error) int {
			return k * 2
		}, func(x string) int { return len(x) }))
	}
}

func TestResult_MapOr_1(t *testing.T) {
	{
		var x = result.Ok("foo")
		assert.Equal(t, 3, result.MapOr(x, 42, func(x string) int { return len(x) }))
	}
	{
		var x = result.TryErr[string]("foo")
		assert.Equal(t, 42, result.MapOr(x, 42, func(x string) int { return len(x) }))
	}
}

func TestResult_Map_3(t *testing.T) {
	var isMyNum = func(s string, search int) result.Result[bool] {
		return result.Map(result.Ret(strconv.Atoi(s)), func(x int) bool { return x == search })
	}
	assert.Equal(t, result.Ok[bool](true), isMyNum("1", 1))
	assert.Equal(t, "Err(strconv.Atoi: parsing \"lol\": invalid syntax)", isMyNum("lol", 1).String())
	assert.Equal(t, "Err(strconv.Atoi: parsing \"NaN\": invalid syntax)", isMyNum("NaN", 1).String())
}

func TestAnd(t *testing.T) {
	// Test with Ok and Ok
	r1 := result.Ok(2)
	r2 := result.Ok(3)
	result1 := result.And(r1, r2)
	assert.True(t, result1.IsOk())
	assert.Equal(t, 3, result1.Unwrap())

	// Test with Ok and Err
	r3 := result.Ok(2)
	r4 := result.TryErr[int]("error")
	result2 := result.And(r3, r4)
	assert.True(t, result2.IsErr())

	// Test with Err and Ok
	r5 := result.TryErr[int]("error")
	r6 := result.Ok(3)
	result3 := result.And(r5, r6)
	assert.True(t, result3.IsErr())
	assert.Equal(t, "error", result3.Err().Error())

	// Test with Err and Err
	r7 := result.TryErr[int]("error1")
	r8 := result.TryErr[int]("error2")
	result4 := result.And(r7, r8)
	assert.True(t, result4.IsErr())
	assert.Equal(t, "error1", result4.Err().Error())
}

func TestAndThen(t *testing.T) {
	// Test with Ok and successful operation
	r1 := result.Ok(2)
	result1 := result.AndThen(r1, func(x int) result.Result[int] {
		return result.Ok(x * 2)
	})
	assert.True(t, result1.IsOk())
	assert.Equal(t, 4, result1.Unwrap())

	// Test with Ok and error operation
	r2 := result.Ok(2)
	result2 := result.AndThen(r2, func(x int) result.Result[int] {
		return result.TryErr[int]("error")
	})
	assert.True(t, result2.IsErr())

	// Test with Err (should not call function)
	r3 := result.TryErr[int]("error")
	result3 := result.AndThen(r3, func(x int) result.Result[int] {
		return result.Ok(x * 2)
	})
	assert.True(t, result3.IsErr())
	assert.Equal(t, "error", result3.Err().Error())
}

func TestFlatten(t *testing.T) {
	// Test with Ok(Ok(value)) - matches documentation example
	{
		r1 := result.Ok(result.Ok(1))
		result1 := result.Flatten(r1)
		assert.Equal(t, result.Ok[int](1), result1)
	}
	// Test with Ok(Err(error)) - matches documentation example
	{
		r2 := result.Ok(result.TryErr[int](errors.New("error")))
		result2 := result.Flatten(r2)
		assert.Equal(t, "error", result2.Err().Error())
	}
	// Test with Err - matches documentation example
	{
		r3 := result.TryErr[result.Result[int]](errors.New("error"))
		result3 := result.Flatten(r3)
		assert.Equal(t, "error", result3.Err().Error())
	}
	// Additional test with Ok(Ok(value)) using different value
	{
		r1 := result.Ok(result.Ok(42))
		result1 := result.Flatten(r1)
		assert.True(t, result1.IsOk())
		assert.Equal(t, 42, result1.Unwrap())
	}
	// Additional test with Ok(Err(error)) using string error
	{
		r2 := result.Ok(result.TryErr[int]("error"))
		result2 := result.Flatten(r2)
		assert.True(t, result2.IsErr())
		assert.Equal(t, "error", result2.Err().Error())
	}
	// Additional test with Err using string error
	{
		r3 := result.TryErr[result.Result[int]]("outer error")
		result3 := result.Flatten(r3)
		assert.True(t, result3.IsErr())
		assert.Equal(t, "outer error", result3.Err().Error())
	}
}

func TestFlatten2(t *testing.T) {
	// Test with Ok result and nil error - matches documentation example
	{
		r1 := result.Ok(1)
		var err1 error = nil
		result1 := result.Flatten2(r1, err1)
		assert.Equal(t, result.Ok[int](1), result1)
	}
	// Test with Ok result and error - matches documentation example
	{
		r2 := result.Ok(1)
		err2 := errors.New("error")
		result2 := result.Flatten2(r2, err2)
		assert.Equal(t, "error", result2.Err().Error())
	}
	// Test with Err result and nil error - matches documentation example
	{
		r3 := result.TryErr[int](errors.New("error"))
		var err3 error = nil
		result3 := result.Flatten2(r3, err3)
		assert.Equal(t, "error", result3.Err().Error())
	}
	// Additional test with Ok result and nil error using different value
	{
		r := result.Ok(42)
		result := result.Flatten2(r, nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, 42, result.Unwrap())
	}
	// Additional test with Ok result and error
	{
		r := result.Ok(42)
		result := result.Flatten2(r, errors.New("test error"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "test error", result.Err().Error())
	}
	// Additional test with Err result and nil error
	{
		r := result.TryErr[int]("original error")
		result := result.Flatten2(r, nil)
		assert.True(t, result.IsErr())
		assert.Equal(t, "original error", result.Err().Error())
	}
	// Additional test with Err result and error (error parameter takes precedence)
	{
		r := result.TryErr[int]("original error")
		result := result.Flatten2(r, errors.New("new error"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "new error", result.Err().Error())
	}
}

func TestAssert_TypeAssertionError(t *testing.T) {
	// Test Assert with type assertion error
	r1 := result.Ok[any](42)
	result1 := result.Assert[any, string](r1)
	assert.True(t, result1.IsErr())
	assert.Contains(t, result1.Err().Error(), "type assert error")
}

func TestXAssert_TypeAssertionError(t *testing.T) {
	// Test XAssert with type assertion error
	r1 := result.Ok[any](42)
	result1 := result.XAssert[string](r1)
	assert.True(t, result1.IsErr())
	assert.Contains(t, result1.Err().Error(), "type assert error")
}

func TestAssert2(t *testing.T) {
	// Test with nil error and successful type assertion
	{
		result := result.Assert2[int, int](42, nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, 42, result.Unwrap())
	}
	// Test with error
	{
		result := result.Assert2[int, int](42, errors.New("test error"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "test error", result.Err().Error())
	}
	// Test with nil error and failed type assertion
	{
		result := result.Assert2[string, int]("hello", nil)
		assert.True(t, result.IsErr())
		assert.Contains(t, result.Err().Error(), "type assert error")
	}
	// Test with successful type assertion from any
	{
		var v any = 42
		result := result.Assert2[any, int](v, nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, 42, result.Unwrap())
	}
}

func TestXAssert2(t *testing.T) {
	// Test with nil error and successful type assertion
	{
		result := result.XAssert2[int](42, nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, 42, result.Unwrap())
	}
	// Test with error
	{
		result := result.XAssert2[int](42, errors.New("test error"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "test error", result.Err().Error())
	}
	// Test with nil error and failed type assertion
	{
		result := result.XAssert2[int]("hello", nil)
		assert.True(t, result.IsErr())
		assert.Contains(t, result.Err().Error(), "type assert error")
	}
	// Test with string type
	{
		result := result.XAssert2[string]("hello", nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, "hello", result.Unwrap())
	}
}

func TestMap2(t *testing.T) {
	// Test with nil error and successful mapping
	{
		result := result.Map2(2, nil, func(x int) int {
			return x * 2
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, 4, result.Unwrap())
	}
	// Test with error
	{
		result := result.Map2(2, errors.New("test error"), func(x int) int {
			return x * 2
		})
		assert.True(t, result.IsErr())
		assert.Equal(t, "test error", result.Err().Error())
	}
	// Test with type conversion
	{
		result := result.Map2(2, nil, func(x int) string {
			return fmt.Sprintf("%d", x)
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, "2", result.Unwrap())
	}
}

func TestMapOr2(t *testing.T) {
	// Test with nil error and successful mapping
	{
		result := result.MapOr2(2, nil, 0, func(x int) int {
			return x * 2
		})
		assert.Equal(t, 4, result)
	}
	// Test with error (should return default)
	{
		result := result.MapOr2(2, errors.New("test error"), 0, func(x int) int {
			return x * 2
		})
		assert.Equal(t, 0, result)
	}
	// Test with different default value
	{
		result := result.MapOr2(2, errors.New("test error"), 100, func(x int) int {
			return x * 2
		})
		assert.Equal(t, 100, result)
	}
	// Test with type conversion
	{
		result := result.MapOr2(2, nil, "default", func(x int) string {
			return fmt.Sprintf("%d", x)
		})
		assert.Equal(t, "2", result)
	}
}

func TestMapOrElse2(t *testing.T) {
	// Test with nil error and successful mapping
	{
		result := result.MapOrElse2(2, nil, func(err error) int {
			return 0
		}, func(x int) int {
			return x * 2
		})
		assert.Equal(t, 4, result)
	}
	// Test with error (should call default function)
	{
		result := result.MapOrElse2(2, errors.New("test error"), func(err error) int {
			return 100
		}, func(x int) int {
			return x * 2
		})
		assert.Equal(t, 100, result)
	}
	// Test with error message in default function
	{
		result := result.MapOrElse2(2, errors.New("test error"), func(err error) string {
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
		result := result.And2(1, nil, 2, nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, 2, result.Unwrap())
	}
	// Test with first error
	{
		result := result.And2(1, errors.New("error1"), 2, nil)
		assert.True(t, result.IsErr())
		assert.Equal(t, "error1", result.Err().Error())
	}
	// Test with second error (first is Ok)
	{
		result := result.And2(1, nil, 2, errors.New("error2"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "error2", result.Err().Error())
	}
	// Test with both errors (should return first error)
	{
		result := result.And2(1, errors.New("error1"), 2, errors.New("error2"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "error1", result.Err().Error())
	}
	// Test with type conversion
	{
		result := result.And2(1, nil, "hello", nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, "hello", result.Unwrap())
	}
}

func TestAndThen2(t *testing.T) {
	// Test with Ok result and successful operation
	{
		r := result.Ok(2)
		result := result.AndThen2(r, func(x int) (int, error) {
			return x * 2, nil
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, 4, result.Unwrap())
	}
	// Test with Ok result and error operation
	{
		r := result.Ok(2)
		result := result.AndThen2(r, func(x int) (int, error) {
			return 0, errors.New("operation error")
		})
		assert.True(t, result.IsErr())
		assert.Equal(t, "operation error", result.Err().Error())
	}
	// Test with Err result (should return original error)
	{
		r := result.TryErr[int]("early error")
		result := result.AndThen2(r, func(x int) (int, error) {
			return x * 2, nil
		})
		assert.True(t, result.IsErr())
		assert.Equal(t, "early error", result.Err().Error())
	}
	// Test with type conversion
	{
		r := result.Ok(2)
		result := result.AndThen2(r, func(x int) (string, error) {
			return fmt.Sprintf("%d", x), nil
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, "2", result.Unwrap())
	}
}

func TestAndThen3(t *testing.T) {
	// Test with nil error and successful operation
	{
		result := result.AndThen3(2, nil, func(x int) (int, error) {
			return x * 2, nil
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, 4, result.Unwrap())
	}
	// Test with error (should return error)
	{
		result := result.AndThen3(2, errors.New("early error"), func(x int) (int, error) {
			return x * 2, nil
		})
		assert.True(t, result.IsErr())
		assert.Equal(t, "early error", result.Err().Error())
	}
	// Test with nil error and error operation
	{
		result := result.AndThen3(2, nil, func(x int) (int, error) {
			return 0, errors.New("operation error")
		})
		assert.True(t, result.IsErr())
		assert.Equal(t, "operation error", result.Err().Error())
	}
	// Test with type conversion
	{
		result := result.AndThen3(2, nil, func(x int) (string, error) {
			return fmt.Sprintf("%d", x), nil
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, "2", result.Unwrap())
	}
}

func TestContains2(t *testing.T) {
	// Test with nil error and matching value
	{
		result := result.Contains2(2, nil, 2)
		assert.True(t, result)
	}
	// Test with nil error and non-matching value
	{
		result := result.Contains2(2, nil, 3)
		assert.False(t, result)
	}
	// Test with error (should return false)
	{
		result := result.Contains2(2, errors.New("test error"), 2)
		assert.False(t, result)
	}
	// Test with string type
	{
		res := result.Contains2("hello", nil, "hello")
		assert.True(t, res)
		res2 := result.Contains2("hello", nil, "world")
		assert.False(t, res2)
	}
}

func TestResult_RetVoid(t *testing.T) {
	// Test RetVoid with nil error
	r1 := result.RetVoid(nil)
	assert.True(t, r1.IsOk())

	// Test RetVoid with error
	err := errors.New("test error")
	r2 := result.RetVoid(err)
	assert.True(t, r2.IsErr())
	assert.Equal(t, "test error", r2.Err().Error())
}

func TestResult_OkVoid(t *testing.T) {
	// Test OkVoid
	r := result.OkVoid()
	assert.True(t, r.IsOk())
	assert.False(t, r.IsErr())
}

func TestResult_TryErrVoid(t *testing.T) {
	// Test TryErrVoid with nil error
	r1 := result.TryErrVoid(nil)
	assert.True(t, r1.IsOk())

	// Test TryErrVoid with error
	err := errors.New("test error")
	r2 := result.TryErrVoid(err)
	assert.True(t, r2.IsErr())
	assert.Equal(t, "test error", r2.Err().Error())
}

func TestResult_FmtErr(t *testing.T) {
	// Test FmtErr with format string
	r1 := result.FmtErr[int]("error: %s", "test")
	assert.True(t, r1.IsErr())
	assert.Contains(t, r1.Err().Error(), "error: test")

	// Test FmtErr with multiple arguments
	r2 := result.FmtErr[string]("error: %d %s", 42, "test")
	assert.True(t, r2.IsErr())
	assert.Contains(t, r2.Err().Error(), "error: 42 test")
}

func TestResult_FmtErrVoid(t *testing.T) {
	// Test FmtErrVoid with format string
	r1 := result.FmtErrVoid("error: %s", "test")
	assert.True(t, r1.IsErr())
	assert.Contains(t, r1.Err().Error(), "error: test")

	// Test FmtErrVoid with multiple arguments
	r2 := result.FmtErrVoid("error: %d %s", 42, "test")
	assert.True(t, r2.IsErr())
	assert.Contains(t, r2.Err().Error(), "error: 42 test")

	// Test FmtErrVoid with no arguments
	r3 := result.FmtErrVoid("simple error")
	assert.True(t, r3.IsErr())
	assert.Equal(t, "simple error", r3.Err().Error())
}

func TestResult_Ret(t *testing.T) {
	// Test Ret with nil error
	r1 := result.Ret(42, nil)
	assert.True(t, r1.IsOk())
	assert.Equal(t, 42, r1.Unwrap())

	// Test Ret with error
	err := errors.New("test error")
	r2 := result.Ret(42, err)
	assert.True(t, r2.IsErr())
	assert.Equal(t, "test error", r2.Err().Error())
}

func TestResult_AssertRet(t *testing.T) {
	// Test AssertRet with valid type
	r1 := result.AssertRet[int](42)
	assert.True(t, r1.IsOk())
	assert.Equal(t, 42, r1.Unwrap())

	// Test AssertRet with invalid type
	r2 := result.AssertRet[int]("string")
	assert.True(t, r2.IsErr())
	assert.Contains(t, r2.Err().Error(), "type assert error")
}
