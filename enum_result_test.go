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

func TestEnumResultUnwrapOrReturn_1(t *testing.T) {
	var r gust.EnumResult[string, string]
	defer func() {
		assert.Equal(t, gust.EnumErr[string, string]("err"), r)
	}()
	defer gust.CatchEnumResult[string, string](&r)
	var r1 = gust.EnumOk[int, string](1)
	var v1 = r1.UnwrapOrReturn()
	assert.Equal(t, 1, v1)
	var r2 = gust.EnumErr[int, string]("err")
	var v2 = r2.UnwrapOrReturn()
	assert.Equal(t, 0, v2)
}

func TestEnumResultUnwrapOrReturn_2(t *testing.T) {
	defer func() {
		assert.Equal(t, "panic text", recover())
	}()
	var r gust.EnumResult[int, string]
	defer gust.CatchEnumResult[int, string](&r)
	panic("panic text")
}
