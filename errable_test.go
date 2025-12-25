package gust_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestErrable(t *testing.T) {
	assert.False(t, gust.ToErrable[any](nil).IsErr())
	assert.False(t, gust.NonErrable[any]().IsErr())

	assert.False(t, gust.ToErrable[error](nil).IsErr())
	assert.False(t, gust.NonErrable[int]().IsErr())

	assert.False(t, gust.ToErrable[*int](nil).IsErr())
	assert.False(t, gust.NonErrable[*int]().IsErr())

	assert.True(t, gust.ToErrable[any](1).IsErr())
	assert.True(t, gust.ToErrable[error](fmt.Errorf("")).IsErr())
	assert.PanicsWithError(t, "test TryPanic", gust.ToErrable[error](errors.New("test TryPanic")).TryPanic)
}

func ExampleErrable() {
	var hasErr = true
	var f = func() gust.Errable[int] {
		if hasErr {
			return gust.ToErrable(1)
		}
		return gust.NonErrable[int]()
	}
	var r = f()
	fmt.Println(r.IsErr())
	fmt.Println(r.UnwrapErr())
	fmt.Printf("%#v", r.ToError())
	// Output:
	// true
	// 1
	// &gust.ErrBox{val:1}
}

func TestErrableTryThrow_1(t *testing.T) {
	var e gust.Errable[int]
	defer func() {
		assert.Equal(t, gust.ToErrable[int](1), e)
	}()
	defer gust.CatchErrable[int](&e)
	gust.ToErrable(1).TryThrow()
}

func TestErrableTryThrow_2(t *testing.T) {
	defer func() {
		assert.Equal(t, "panic text", recover())
	}()
	var e gust.Errable[string]
	defer gust.CatchErrable[string](&e)
	panic("panic text")
}

func TestErrableTryThrow_3(t *testing.T) {
	var r gust.Result[string]
	defer func() {
		assert.Equal(t, gust.Err[string]("err"), r)
	}()
	defer gust.CatchResult[string](&r)
	assert.Equal(t, gust.Void(nil), gust.ToErrable("err").Result().UnwrapOrThrow())
}

func TestErrableTryThrow_4(t *testing.T) {
	var r gust.EnumResult[int, string]
	defer func() {
		assert.Equal(t, gust.EnumErr[int, string]("err"), r)
	}()
	defer gust.CatchEnumResult[int, string](&r)
	assert.Equal(t, gust.Void(nil), gust.ToErrable("err").EnumResult().UnwrapOrThrow())
}

func TestErrableTryThrow_5(t *testing.T) {
	var r gust.EnumResult[int, string]
	defer func() {
		assert.Equal(t, gust.EnumErr[int, string]("err"), r)
	}()
	defer gust.CatchEnumResult[int, string](&r)
	gust.TryThrow("err")
}

func TestErrableTryThrow_6(t *testing.T) {
	var e gust.Errable[int]
	defer func() {
		assert.Equal(t, gust.ToErrable[int](1), e)
	}()
	defer e.Catch()
	gust.ToErrable(1).TryThrow()
}

func TestErrableTryThrow_7(t *testing.T) {
	defer func() {
		assert.Equal(t, "panic text", recover())
	}()
	var e gust.Errable[string]
	defer e.Catch()
	panic("panic text")
}

func TestErrableTryThrow_8(t *testing.T) {
	var r gust.Result[string]
	defer func() {
		assert.Equal(t, gust.Err[string]("err"), r)
	}()
	defer r.Catch()
	assert.Equal(t, gust.Void(nil), gust.ToErrable("err").Result().UnwrapOrThrow())
}

func TestErrableTryThrow_9(t *testing.T) {
	var r gust.EnumResult[int, string]
	defer func() {
		assert.Equal(t, gust.EnumErr[int, string]("err"), r)
	}()
	defer r.Catch()
	assert.Equal(t, gust.Void(nil), gust.ToErrable("err").EnumResult().UnwrapOrThrow())
}

func TestErrableTryThrow_10(t *testing.T) {
	var r gust.EnumResult[int, string]
	defer func() {
		assert.Equal(t, gust.EnumErr[int, string]("err"), r)
	}()
	defer r.Catch()
	gust.TryThrow("err")
}

func TestFmtErrable(t *testing.T) {
	errable := gust.FmtErrable("error: %s", "test")
	assert.True(t, errable.IsErr())
	assert.Contains(t, errable.ToError().Error(), "error: test")
}

func TestToErrable_NilPointers(t *testing.T) {
	// Test various nil pointer types
	var nilInt *int
	assert.False(t, gust.ToErrable[*int](nilInt).IsErr())

	var nilInt64 *int64
	assert.False(t, gust.ToErrable[*int64](nilInt64).IsErr())

	var nilInt32 *int32
	assert.False(t, gust.ToErrable[*int32](nilInt32).IsErr())

	var nilInt16 *int16
	assert.False(t, gust.ToErrable[*int16](nilInt16).IsErr())

	var nilInt8 *int8
	assert.False(t, gust.ToErrable[*int8](nilInt8).IsErr())

	var nilUint *uint
	assert.False(t, gust.ToErrable[*uint](nilUint).IsErr())

	var nilUint64 *uint64
	assert.False(t, gust.ToErrable[*uint64](nilUint64).IsErr())

	var nilUint32 *uint32
	assert.False(t, gust.ToErrable[*uint32](nilUint32).IsErr())

	var nilUint16 *uint16
	assert.False(t, gust.ToErrable[*uint16](nilUint16).IsErr())

	var nilUint8 *uint8
	assert.False(t, gust.ToErrable[*uint8](nilUint8).IsErr())

	var nilFloat32 *float32
	assert.False(t, gust.ToErrable[*float32](nilFloat32).IsErr())

	var nilFloat64 *float64
	assert.False(t, gust.ToErrable[*float64](nilFloat64).IsErr())

	var nilComplex64 *complex64
	assert.False(t, gust.ToErrable[*complex64](nilComplex64).IsErr())

	var nilComplex128 *complex128
	assert.False(t, gust.ToErrable[*complex128](nilComplex128).IsErr())

	var nilString *string
	assert.False(t, gust.ToErrable[*string](nilString).IsErr())

	var nilBool *bool
	assert.False(t, gust.ToErrable[*bool](nilBool).IsErr())

	// Test non-nil pointers
	intVal := 42
	assert.True(t, gust.ToErrable[*int](&intVal).IsErr())
	assert.Equal(t, &intVal, gust.ToErrable[*int](&intVal).UnwrapErr())
}

func TestToErrable_NonPointerTypes(t *testing.T) {
	// Test various non-pointer types
	assert.True(t, gust.ToErrable[int](42).IsErr())
	assert.True(t, gust.ToErrable[int64](42).IsErr())
	assert.True(t, gust.ToErrable[int32](42).IsErr())
	assert.True(t, gust.ToErrable[int16](42).IsErr())
	assert.True(t, gust.ToErrable[int8](42).IsErr())
	assert.True(t, gust.ToErrable[uint](42).IsErr())
	assert.True(t, gust.ToErrable[uint64](42).IsErr())
	assert.True(t, gust.ToErrable[uint32](42).IsErr())
	assert.True(t, gust.ToErrable[uint16](42).IsErr())
	assert.True(t, gust.ToErrable[uint8](42).IsErr())
	assert.True(t, gust.ToErrable[float32](42.0).IsErr())
	assert.True(t, gust.ToErrable[float64](42.0).IsErr())
	assert.True(t, gust.ToErrable[complex64](42+0i).IsErr())
	assert.True(t, gust.ToErrable[complex128](42+0i).IsErr())
	assert.True(t, gust.ToErrable[string]("test").IsErr())
	assert.True(t, gust.ToErrable[bool](true).IsErr())
}

func TestToErrable_CustomPointerType(t *testing.T) {
	type CustomStruct struct {
		Value int
	}
	var nilCustom *CustomStruct
	errable := gust.ToErrable[*CustomStruct](nilCustom)
	assert.False(t, errable.IsErr())

	custom := &CustomStruct{Value: 42}
	errable2 := gust.ToErrable[*CustomStruct](custom)
	assert.True(t, errable2.IsErr())
	assert.Equal(t, custom, errable2.UnwrapErr())
}

func TestErrable_UnwrapErrOr(t *testing.T) {
	{
		var e = gust.ToErrable[int](42)
		assert.Equal(t, 42, e.UnwrapErrOr(0))
	}
	{
		var e = gust.NonErrable[int]()
		assert.Equal(t, 0, e.UnwrapErrOr(0))
	}
}

func TestErrable_EnumResult(t *testing.T) {
	{
		var e = gust.ToErrable[string]("error")
		er := e.EnumResult()
		assert.True(t, er.IsErr())
		assert.Equal(t, "error", er.UnwrapErr())
	}
	{
		var e = gust.NonErrable[string]()
		er := e.EnumResult()
		assert.True(t, er.IsOk())
		assert.Equal(t, gust.Void(nil), er.Unwrap())
	}
}

func TestErrable_Result(t *testing.T) {
	{
		var e = gust.ToErrable[string]("error")
		result := e.Result()
		assert.True(t, result.IsErr())
	}
	{
		var e = gust.NonErrable[string]()
		result := e.Result()
		assert.True(t, result.IsOk())
		assert.Equal(t, gust.Void(nil), result.Unwrap())
	}
}

func TestErrable_Option(t *testing.T) {
	{
		var e = gust.ToErrable[string]("error")
		opt := e.Option()
		assert.True(t, opt.IsSome())
		assert.Equal(t, "error", opt.Unwrap())
	}
	{
		var e = gust.NonErrable[string]()
		opt := e.Option()
		assert.True(t, opt.IsNone())
	}
}

func TestErrable_CtrlFlow(t *testing.T) {
	{
		var e = gust.ToErrable[string]("error")
		cf := e.CtrlFlow()
		assert.True(t, cf.IsBreak())
		assert.Equal(t, "error", cf.UnwrapBreak())
	}
	{
		var e = gust.NonErrable[string]()
		cf := e.CtrlFlow()
		assert.True(t, cf.IsContinue())
		assert.Equal(t, gust.Void(nil), cf.UnwrapContinue())
	}
}

func TestErrable_InspectErr(t *testing.T) {
	called := false
	{
		var e = gust.ToErrable[string]("error")
		result := e.InspectErr(func(s string) {
			called = true
			assert.Equal(t, "error", s)
		})
		assert.True(t, called)
		assert.True(t, result.IsErr())
	}
	called = false
	{
		var e = gust.NonErrable[string]()
		result := e.InspectErr(func(s string) {
			called = true
		})
		assert.False(t, called)
		assert.False(t, result.IsErr())
	}
}

func TestErrable_Inspect(t *testing.T) {
	called := false
	{
		var e = gust.NonErrable[string]()
		result := e.Inspect(func() {
			called = true
		})
		assert.True(t, called)
		assert.False(t, result.IsErr())
	}
	called = false
	{
		var e = gust.ToErrable[string]("error")
		result := e.Inspect(func() {
			called = true
		})
		assert.False(t, called)
		assert.True(t, result.IsErr())
	}
}

func TestErrable_UnwrapErr(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()
	var e = gust.NonErrable[string]()
	_ = e.UnwrapErr() // Should panic
}

func TestErrable_ToError(t *testing.T) {
	{
		var e = gust.NonErrable[string]()
		assert.Nil(t, e.ToError())
	}
	{
		var e = gust.ToErrable[string]("error")
		assert.NotNil(t, e.ToError())
		assert.Equal(t, "error", e.ToError().Error())
	}
	{
		var e = gust.ToErrable[error](errors.New("std error"))
		assert.NotNil(t, e.ToError())
		assert.Equal(t, "std error", e.ToError().Error())
	}
}

func TestErrable_Catch_NilReceiver(t *testing.T) {
	// Test Catch with nil receiver
	defer func() {
		assert.NotNil(t, recover())
	}()
	var nilErrable *gust.Errable[string]
	defer nilErrable.Catch()
	gust.ToErrable("test error").TryThrow()
}

func TestCatchErrable_NilReceiver(t *testing.T) {
	// Test CatchErrable with nil receiver
	defer func() {
		assert.NotNil(t, recover())
	}()
	defer gust.CatchErrable[string](nil)
	gust.ToErrable("test error").TryThrow()
}

func TestErrable_Catch_NonPanicValue(t *testing.T) {
	// Test Catch with non-panicValue panic
	defer func() {
		assert.Equal(t, "regular panic", recover())
	}()
	var errable gust.Errable[string]
	defer errable.Catch()
	panic("regular panic")
}

func TestCatchErrable_NonPanicValue(t *testing.T) {
	// Test CatchErrable with non-panicValue panic
	defer func() {
		assert.Equal(t, "regular panic", recover())
	}()
	var errable gust.Errable[string]
	defer gust.CatchErrable(&errable)
	panic("regular panic")
}

func TestToErrable_ReflectPath(t *testing.T) {
	// Test ToErrable with reflect path (custom pointer type)
	type CustomStruct struct {
		Value int
	}
	var nilCustom *CustomStruct
	errable := gust.ToErrable[*CustomStruct](nilCustom)
	assert.False(t, errable.IsErr())

	custom := &CustomStruct{Value: 42}
	errable2 := gust.ToErrable[*CustomStruct](custom)
	assert.True(t, errable2.IsErr())
	assert.Equal(t, custom, errable2.UnwrapErr())
}

func TestToErrable_InterfaceNil(t *testing.T) {
	// Test ToErrable with interface containing nil
	var iface interface{} = (*int)(nil)
	errable := gust.ToErrable[interface{}](iface)
	// Should use reflect path and detect nil
	assert.False(t, errable.IsErr())
}

func TestToErrable_InterfaceNonNil(t *testing.T) {
	// Test ToErrable with interface containing non-nil
	val := 42
	var iface interface{} = &val
	errable := gust.ToErrable[interface{}](iface)
	assert.True(t, errable.IsErr())
}
