package ret_test

import (
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
	// Test with Ok(Ok(value))
	r1 := gust.Ok(gust.Ok(42))
	result1 := ret.Flatten(r1)
	assert.True(t, result1.IsOk())
	assert.Equal(t, 42, result1.Unwrap())

	// Test with Ok(Err(error))
	r2 := gust.Ok(gust.Err[int]("error"))
	result2 := ret.Flatten(r2)
	assert.True(t, result2.IsErr())
	assert.Equal(t, "error", result2.Err().Error())

	// Test with Err
	r3 := gust.Err[gust.Result[int]]("outer error")
	result3 := ret.Flatten(r3)
	assert.True(t, result3.IsErr())
	assert.Equal(t, "outer error", result3.Err().Error())
}
