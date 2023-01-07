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
