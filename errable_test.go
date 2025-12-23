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
