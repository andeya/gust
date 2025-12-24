package gust_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestEnumResultJSON(t *testing.T) {
	var r = gust.EnumErr[any, error](errors.New("err"))
	var b, err = json.Marshal(r)
	assert.Equal(t, "json: error calling MarshalJSON for type gust.EnumResult[interface {},error]: err", err.Error())
	assert.Nil(t, b)
	type T struct {
		Name string
	}
	var r2 = gust.EnumOk[T, error](T{Name: "andeya"})
	var b2, err2 = json.Marshal(r2)
	assert.NoError(t, err2)
	assert.Equal(t, `{"Name":"andeya"}`, string(b2))

	var r3 gust.EnumResult[T, error]
	var err3 = json.Unmarshal(b2, &r3)
	assert.NoError(t, err3)
	assert.Equal(t, r2, r3)

	var r4 gust.EnumResult[T, error]
	var err4 = json.Unmarshal([]byte("0"), &r4)
	assert.True(t, r4.IsErr())
	assert.Equal(t, "json: cannot unmarshal number into Go value of type gust_test.T", err4.Error())
}

func TestEnumResultIsValid(t *testing.T) {
	var r0 *gust.EnumResult[any, any]
	assert.False(t, r0.IsValid())
	var r1 gust.EnumResult[any, any]
	assert.False(t, r1.IsValid())
	assert.False(t, (&gust.EnumResult[any, any]{}).IsValid())
	var r2 = gust.EnumOk[any, any](nil)
	assert.True(t, r2.IsValid())
}

func TestEnumResultUnwrapOrThrow_1(t *testing.T) {
	var r gust.EnumResult[string, string]
	defer func() {
		assert.Equal(t, gust.EnumErr[string, string]("err"), r)
	}()
	defer gust.CatchEnumResult[string, string](&r)
	var r1 = gust.EnumOk[int, string](1)
	var v1 = r1.UnwrapOrThrow()
	assert.Equal(t, 1, v1)
	var r2 = gust.EnumErr[int, string]("err")
	var v2 = r2.UnwrapOrThrow()
	assert.Equal(t, 0, v2)
}

func TestEnumResultUnwrapOrThrow_2(t *testing.T) {
	defer func() {
		assert.Equal(t, "panic text", recover())
	}()
	var r gust.EnumResult[int, string]
	defer gust.CatchEnumResult[int, string](&r)
	panic("panic text")
}

func TestEnumResultUnwrapOrThrow_3(t *testing.T) {
	var r gust.EnumResult[string, string]
	defer func() {
		assert.Equal(t, gust.EnumErr[string, string]("err"), r)
	}()
	defer r.Catch()
	var r1 = gust.EnumOk[int, string](1)
	var v1 = r1.UnwrapOrThrow()
	assert.Equal(t, 1, v1)
	var r2 = gust.EnumErr[int, string]("err")
	var v2 = r2.UnwrapOrThrow()
	assert.Equal(t, 0, v2)
}

func TestEnumResultUnwrapOrThrow_4(t *testing.T) {
	defer func() {
		assert.Equal(t, "panic text", recover())
	}()
	var r gust.EnumResult[int, string]
	defer r.Catch()
	panic("panic text")
}

func TestEnumResult_Err(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](2)
		assert.Equal(t, gust.None[string](), x.Err())
	}
	{
		var x = gust.EnumErr[int, string]("some error message")
		assert.Equal(t, gust.Some[string]("some error message"), x.Err())
	}
}

func TestEnumResult_IsErrAnd(t *testing.T) {
	{
		var x = gust.EnumErr[int, int8](-1)
		assert.True(t, x.IsErrAnd(func(x int8) bool { return x == -1 }))
	}
}

func TestEnumResult_XOk(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](2)
		assert.Equal(t, gust.Some[any](2), x.XOk())
	}
	{
		var x = gust.EnumErr[int, string]("error")
		assert.Equal(t, gust.None[any](), x.XOk())
	}
}

func TestEnumResult_XErr(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](2)
		assert.Equal(t, gust.None[any](), x.XErr())
	}
	{
		var x = gust.EnumErr[int, string]("error")
		assert.Equal(t, gust.Some[any]("error"), x.XErr())
	}
}

func TestEnumResult_ToXOk(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](42)
		xResult := x.ToXOk()
		assert.True(t, xResult.IsOk())
		assert.Equal(t, 42, xResult.Unwrap())
	}
	{
		var x = gust.EnumErr[int, string]("error")
		xResult := x.ToXOk()
		assert.True(t, xResult.IsErr())
		assert.Equal(t, "error", xResult.UnwrapErr())
	}
}

func TestEnumResult_ToXErr(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](42)
		xResult := x.ToXErr()
		assert.True(t, xResult.IsOk())
		assert.Equal(t, 42, xResult.Unwrap())
	}
	{
		var x = gust.EnumErr[int, string]("error")
		xResult := x.ToXErr()
		assert.True(t, xResult.IsErr())
		assert.Equal(t, "error", xResult.UnwrapErr())
	}
}

func TestEnumResult_ToX(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](42)
		xResult := x.ToX()
		assert.True(t, xResult.IsOk())
		assert.Equal(t, 42, xResult.Unwrap())
	}
	{
		var x = gust.EnumErr[int, string]("error")
		xResult := x.ToX()
		assert.True(t, xResult.IsErr())
		assert.Equal(t, "error", xResult.UnwrapErr())
	}
}

func TestEnumResult_XMap(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](2)
		result := x.XMap(func(i int) any { return i * 2 })
		assert.True(t, result.IsOk())
		assert.Equal(t, 4, result.Unwrap())
	}
	{
		var x = gust.EnumErr[int, string]("error")
		result := x.XMap(func(i int) any { return i * 2 })
		assert.True(t, result.IsErr())
	}
}

func TestEnumResult_XMapOr(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](2)
		result := x.XMapOr("default", func(i int) any { return i * 2 })
		assert.Equal(t, 4, result)
	}
	{
		var x = gust.EnumErr[int, string]("error")
		result := x.XMapOr("default", func(i int) any { return i * 2 })
		assert.Equal(t, "default", result)
	}
}

func TestEnumResult_XMapOrElse(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](2)
		result := x.XMapOrElse(func(string) any { return "default" }, func(i int) any { return i * 2 })
		assert.Equal(t, 4, result)
	}
	{
		var x = gust.EnumErr[int, string]("error")
		result := x.XMapOrElse(func(string) any { return "default" }, func(i int) any { return i * 2 })
		assert.Equal(t, "default", result)
	}
}

func TestEnumResult_XMapErr(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](2)
		result := x.XMapErr(func(s string) any { return "mapped: " + s })
		assert.True(t, result.IsOk())
		assert.Equal(t, 2, result.Unwrap())
	}
	{
		var x = gust.EnumErr[int, string]("error")
		result := x.XMapErr(func(s string) any { return "mapped: " + s })
		assert.True(t, result.IsErr())
		assert.Equal(t, "mapped: error", result.UnwrapErr())
	}
}

func TestEnumResult_XAnd(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](2)
		var y = gust.EnumOk[any, string]("foo")
		result := x.XAnd(y)
		assert.True(t, result.IsOk())
		assert.Equal(t, "foo", result.Unwrap())
	}
	{
		var x = gust.EnumErr[int, string]("error")
		var y = gust.EnumOk[any, string]("foo")
		result := x.XAnd(y)
		assert.True(t, result.IsErr())
		assert.Equal(t, "error", result.UnwrapErr())
	}
}

func TestEnumResult_XAndThen(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](2)
		result := x.XAndThen(func(i int) gust.EnumResult[any, string] {
			return gust.EnumOk[any, string](i * 2)
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, 4, result.Unwrap())
	}
	{
		var x = gust.EnumErr[int, string]("error")
		result := x.XAndThen(func(i int) gust.EnumResult[any, string] {
			return gust.EnumOk[any, string](i * 2)
		})
		assert.True(t, result.IsErr())
		assert.Equal(t, "error", result.UnwrapErr())
	}
}

func TestEnumResult_XOr(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](2)
		var y = gust.EnumOk[int, any](42)
		result := x.XOr(y)
		assert.True(t, result.IsOk())
		assert.Equal(t, 2, result.Unwrap())
	}
	{
		var x = gust.EnumErr[int, string]("error")
		var y = gust.EnumOk[int, any](42)
		result := x.XOr(y)
		assert.True(t, result.IsOk())
		assert.Equal(t, 42, result.Unwrap())
	}
}

func TestEnumResult_XOrElse(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](2)
		result := x.XOrElse(func(s string) gust.EnumResult[int, any] {
			return gust.EnumOk[int, any](42)
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, 2, result.Unwrap())
	}
	{
		var x = gust.EnumErr[int, string]("error")
		result := x.XOrElse(func(s string) gust.EnumResult[int, any] {
			return gust.EnumOk[int, any](42)
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, 42, result.Unwrap())
	}
}

func TestEnumResult_Iterator(t *testing.T) {
	// Test Next
	{
		var x = gust.EnumOk[string, string]("foo")
		opt := x.Next()
		assert.Equal(t, gust.Some("foo"), opt)
	}
	{
		var x = gust.EnumErr[string, string]("error")
		opt := x.Next()
		assert.True(t, opt.IsNone())
	}
	{
		var nilResult *gust.EnumResult[string, string]
		opt := nilResult.Next()
		assert.True(t, opt.IsNone())
	}

	// Test NextBack
	{
		var x = gust.EnumOk[string, string]("bar")
		opt := x.NextBack()
		assert.Equal(t, gust.Some("bar"), opt)
	}

	// Test Remaining
	{
		var x = gust.EnumOk[string, string]("baz")
		assert.Equal(t, uint(1), x.Remaining())
	}
	{
		var x = gust.EnumErr[string, string]("error")
		assert.Equal(t, uint(0), x.Remaining())
	}
	{
		var nilResult *gust.EnumResult[string, string]
		assert.Equal(t, uint(0), nilResult.Remaining())
	}
}

func TestEnumResult_CtrlFlow(t *testing.T) {
	{
		var x = gust.EnumOk[string, string]("foo")
		cf := x.CtrlFlow()
		assert.True(t, cf.IsContinue())
		assert.Equal(t, "foo", cf.UnwrapContinue())
	}
	{
		var x = gust.EnumErr[string, string]("error")
		cf := x.CtrlFlow()
		assert.True(t, cf.IsBreak())
		assert.Equal(t, "error", cf.UnwrapBreak())
	}
}

func TestEnumResult_Ref(t *testing.T) {
	result := gust.EnumOk[int, string](42)
	ref := result.Ref()
	assert.Equal(t, gust.EnumOk[int, string](42), *ref)
	ref.Unwrap() // Should not panic
}

func TestEnumResult_Split(t *testing.T) {
	{
		var x = gust.EnumOk[string, string]("foo")
		val, err := x.Split()
		assert.Equal(t, "foo", val)
		assert.Equal(t, "", err)
	}
	{
		var x = gust.EnumErr[string, string]("error")
		val, err := x.Split()
		assert.Equal(t, "", val)
		assert.Equal(t, "error", err)
	}
}

func TestEnumResult_String(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](42)
		assert.Contains(t, x.String(), "Ok")
		assert.Contains(t, x.String(), "42")
	}
	{
		var x = gust.EnumErr[int, string]("error")
		assert.Contains(t, x.String(), "Err")
		assert.Contains(t, x.String(), "error")
	}
}

func TestEnumResult_IsOkAnd(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](2)
		assert.True(t, x.IsOkAnd(func(i int) bool { return i > 1 }))
	}
	{
		var x = gust.EnumOk[int, string](0)
		assert.False(t, x.IsOkAnd(func(i int) bool { return i > 1 }))
	}
	{
		var x = gust.EnumErr[int, string]("error")
		assert.False(t, x.IsOkAnd(func(i int) bool { return i > 1 }))
	}
}

func TestEnumResult_MapErr(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](2)
		result := x.MapErr(func(s string) string { return "mapped: " + s })
		assert.True(t, result.IsOk())
		assert.Equal(t, 2, result.Unwrap())
	}
	{
		var x = gust.EnumErr[int, string]("error")
		result := x.MapErr(func(s string) string { return "mapped: " + s })
		assert.True(t, result.IsErr())
		assert.Equal(t, "mapped: error", result.UnwrapErr())
	}
}

func TestEnumResult_Inspect(t *testing.T) {
	called := false
	{
		var x = gust.EnumOk[int, string](42)
		result := x.Inspect(func(i int) {
			called = true
			assert.Equal(t, 42, i)
		})
		assert.True(t, called)
		assert.True(t, result.IsOk())
	}
	called = false
	{
		var x = gust.EnumErr[int, string]("error")
		result := x.Inspect(func(i int) {
			called = true
		})
		assert.False(t, called)
		assert.True(t, result.IsErr())
	}
}

func TestEnumResult_InspectErr(t *testing.T) {
	called := false
	{
		var x = gust.EnumErr[int, string]("error")
		result := x.InspectErr(func(s string) {
			called = true
			assert.Equal(t, "error", s)
		})
		assert.True(t, called)
		assert.True(t, result.IsErr())
	}
	called = false
	{
		var x = gust.EnumOk[int, string](42)
		result := x.InspectErr(func(s string) {
			called = true
		})
		assert.False(t, called)
		assert.True(t, result.IsOk())
	}
}

func TestEnumResult_UnwrapOr(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](2)
		assert.Equal(t, 2, x.UnwrapOr(0))
	}
	{
		var x = gust.EnumErr[int, string]("error")
		assert.Equal(t, 0, x.UnwrapOr(0))
	}
}

func TestEnumResult_UnwrapOrElse(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](2)
		assert.Equal(t, 2, x.UnwrapOrElse(func(s string) int { return 0 }))
	}
	{
		var x = gust.EnumErr[int, string]("error")
		assert.Equal(t, 0, x.UnwrapOrElse(func(s string) int { return 0 }))
	}
}

func TestEnumResult_Expect(t *testing.T) {
	defer func() {
		assert.Contains(t, recover().(error).Error(), "Testing expect")
	}()
	gust.EnumErr[int, string]("error").Expect("Testing expect")
}

func TestEnumResult_Unwrap(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()
	gust.EnumErr[int, string]("error").Unwrap()
}

func TestEnumResult_ExpectErr(t *testing.T) {
	defer func() {
		assert.Contains(t, recover().(error).Error(), "Testing expect_err")
	}()
	gust.EnumOk[int, string](10).ExpectErr("Testing expect_err")
}

func TestEnumResult_UnwrapErr(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()
	gust.EnumOk[int, string](10).UnwrapErr()
}

func TestEnumResult_AsPtr(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](42)
		ptr := x.AsPtr()
		assert.NotNil(t, ptr)
		assert.Equal(t, 42, *ptr)
	}
	{
		var x = gust.EnumErr[int, string]("error")
		ptr := x.AsPtr()
		assert.Nil(t, ptr)
	}
}

func TestEnumResult_Result(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](42)
		result := x.Result()
		assert.True(t, result.IsOk())
		assert.Equal(t, 42, result.Unwrap())
	}
	{
		var x = gust.EnumErr[int, string]("error")
		result := x.Result()
		assert.True(t, result.IsErr())
	}
}

func TestEnumResult_Errable(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](42)
		errable := x.Errable()
		assert.False(t, errable.IsErr())
	}
	{
		var x = gust.EnumErr[int, string]("error")
		errable := x.Errable()
		assert.True(t, errable.IsErr())
		assert.Equal(t, "error", errable.UnwrapErr())
	}
}

func TestEnumResult_UnmarshalJSON_NilReceiver(t *testing.T) {
	// Test UnmarshalJSON with nil receiver
	var nilResult *gust.EnumResult[int, string]
	err := nilResult.UnmarshalJSON([]byte("42"))
	assert.Error(t, err)
	assert.IsType(t, &json.InvalidUnmarshalError{}, err)
}

func TestEnumResult_UnmarshalJSON_ErrorPath(t *testing.T) {
	// Test UnmarshalJSON with invalid JSON (error path)
	var result gust.EnumResult[int, string]
	err := result.UnmarshalJSON([]byte("invalid json"))
	assert.Error(t, err)
	assert.True(t, result.IsErr()) // Should be Err on error
}

func TestEnumResult_UnmarshalJSON_ValidAfterError(t *testing.T) {
	// Test UnmarshalJSON with error first, then valid JSON
	var result gust.EnumResult[int, string]
	// First attempt with invalid JSON
	_ = result.UnmarshalJSON([]byte("invalid"))
	assert.True(t, result.IsErr())
	
	// Then with valid JSON
	err := result.UnmarshalJSON([]byte("42"))
	assert.NoError(t, err)
	assert.True(t, result.IsOk())
	assert.Equal(t, 42, result.Unwrap())
}

func TestEnumResult_UnmarshalJSON_Struct(t *testing.T) {
	type S struct {
		X int
		Y string
	}
	var result gust.EnumResult[S, string]
	err := result.UnmarshalJSON([]byte(`{"X":10,"Y":"test"}`))
	assert.NoError(t, err)
	assert.True(t, result.IsOk())
	assert.Equal(t, S{X: 10, Y: "test"}, result.Unwrap())
}

func TestEnumResult_UnmarshalJSON_Array(t *testing.T) {
	var result gust.EnumResult[[]int, string]
	err := result.UnmarshalJSON([]byte("[1,2,3]"))
	assert.NoError(t, err)
	assert.True(t, result.IsOk())
	assert.Equal(t, []int{1, 2, 3}, result.Unwrap())
}

func TestEnumResult_UnmarshalJSON_Map(t *testing.T) {
	var result gust.EnumResult[map[string]int, string]
	err := result.UnmarshalJSON([]byte(`{"a":1,"b":2}`))
	assert.NoError(t, err)
	assert.True(t, result.IsOk())
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, result.Unwrap())
}

func TestEnumResult_Catch_NilReceiver(t *testing.T) {
	// Test Catch with nil receiver
	defer func() {
		assert.NotNil(t, recover())
	}()
	var nilResult *gust.EnumResult[int, string]
	defer nilResult.Catch()
	gust.EnumErr[int, string]("test error").UnwrapOrThrow()
}

func TestCatchEnumResult_NilReceiver(t *testing.T) {
	// Test CatchEnumResult with nil receiver
	defer func() {
		assert.NotNil(t, recover())
	}()
	defer gust.CatchEnumResult[int, string](nil)
	gust.EnumErr[int, string]("test error").UnwrapOrThrow()
}

func TestEnumResult_Catch_NonPanicValue(t *testing.T) {
	// Test Catch with non-panicValue panic
	defer func() {
		assert.Equal(t, "regular panic", recover())
	}()
	var result gust.EnumResult[int, string]
	defer result.Catch()
	panic("regular panic")
}

func TestCatchEnumResult_NonPanicValue(t *testing.T) {
	// Test CatchEnumResult with non-panicValue panic
	defer func() {
		assert.Equal(t, "regular panic", recover())
	}()
	var result gust.EnumResult[int, string]
	defer gust.CatchEnumResult(&result)
	panic("regular panic")
}

func TestEnumResult_Catch_OkValue(t *testing.T) {
	// Test Catch when result already has Ok value
	var result gust.EnumResult[int, string] = gust.EnumOk[int, string](42)
	defer result.Catch()
	gust.EnumErr[int, string]("test error").UnwrapOrThrow()
	// Result should be updated to Err
	assert.True(t, result.IsErr())
}

func TestCatchEnumResult_OkValue(t *testing.T) {
	// Test CatchEnumResult when result already has Ok value
	var result gust.EnumResult[int, string] = gust.EnumOk[int, string](42)
	defer gust.CatchEnumResult(&result)
	gust.EnumErr[int, string]("test error").UnwrapOrThrow()
	// Result should be updated to Err
	assert.True(t, result.IsErr())
}

func TestEnumResult_XAndThen_ErrorPath(t *testing.T) {
	// Test XAndThen with error path
	result := gust.EnumOk[int, string](42)
	result2 := result.XAndThen(func(i int) gust.EnumResult[any, string] {
		return gust.EnumErr[any, string]("error")
	})
	assert.True(t, result2.IsErr())
}

func TestEnumResult_XAndThen_OkPath(t *testing.T) {
	// Test XAndThen with Ok path
	result := gust.EnumOk[int, string](42)
	result2 := result.XAndThen(func(i int) gust.EnumResult[any, string] {
		return gust.EnumOk[any, string](i * 2)
	})
	assert.True(t, result2.IsOk())
	assert.Equal(t, 84, result2.Unwrap())
}

func TestEnumResult_XOr_ErrorPath(t *testing.T) {
	// Test XOr with error path
	result := gust.EnumErr[int, string]("error")
	result2 := gust.EnumOk[int, any](42)
	result3 := result.XOr(result2)
	assert.True(t, result3.IsOk())
	assert.Equal(t, 42, result3.Unwrap())
}

func TestEnumResult_XOr_OkPath(t *testing.T) {
	// Test XOr with Ok path
	result := gust.EnumOk[int, string](10)
	result2 := gust.EnumOk[int, any](42)
	result3 := result.XOr(result2)
	assert.True(t, result3.IsOk())
	assert.Equal(t, 10, result3.Unwrap())
}

func TestEnumResult_XOrElse_ErrorPath(t *testing.T) {
	// Test XOrElse with error path
	result := gust.EnumErr[int, string]("error")
	result2 := result.XOrElse(func(s string) gust.EnumResult[int, any] {
		return gust.EnumOk[int, any](42)
	})
	assert.True(t, result2.IsOk())
	assert.Equal(t, 42, result2.Unwrap())
}

func TestEnumResult_XOrElse_OkPath(t *testing.T) {
	// Test XOrElse with Ok path
	result := gust.EnumOk[int, string](10)
	result2 := result.XOrElse(func(s string) gust.EnumResult[int, any] {
		return gust.EnumOk[int, any](42)
	})
	assert.True(t, result2.IsOk())
	assert.Equal(t, 10, result2.Unwrap())
}

func TestEnumResult_wrapError_WithError(t *testing.T) {
	// Test wrapError with error type
	result := gust.EnumErr[int, error](errors.New("test error"))
	defer func() {
		err := recover()
		assert.NotNil(t, err)
		assert.Contains(t, err.(error).Error(), "Testing")
		assert.Contains(t, err.(error).Error(), "test error")
	}()
	result.Expect("Testing")
}

func TestEnumResult_wrapError_WithNonError(t *testing.T) {
	// Test wrapError with non-error type
	result := gust.EnumErr[int, string]("test error")
	defer func() {
		err := recover()
		assert.NotNil(t, err)
		assert.Contains(t, err.(error).Error(), "Testing")
		assert.Contains(t, err.(error).Error(), "test error")
	}()
	result.Expect("Testing")
}

func TestFromError_WithMatchingType(t *testing.T) {
	// Test fromError with matching type
	// UnmarshalJSON will use fromError internally
	var result gust.EnumResult[int, error]
	err := result.UnmarshalJSON([]byte("invalid"))
	assert.Error(t, err)
	assert.True(t, result.IsErr())
	// The error should be converted using fromError
}

func TestFromError_WithNonMatchingType(t *testing.T) {
	// Test fromError with non-matching type
	var result gust.EnumResult[int, string]
	err := result.UnmarshalJSON([]byte("invalid"))
	assert.Error(t, err)
	assert.True(t, result.IsErr())
	// The error should be converted to zero value of string
	assert.Equal(t, "", result.UnwrapErr())
}
