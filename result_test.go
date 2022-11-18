package gust_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestResultJSON(t *testing.T) {
	var r = gust.Err[any](errors.New("err"))
	var b, err = json.Marshal(r)
	assert.Equal(t, "json: error calling MarshalJSON for type gust.Result[interface {}]: err", err.Error())
	assert.Nil(t, b)
	type T struct {
		Name string
	}
	var r2 = gust.Ok(T{Name: "andeya"})
	var b2, err2 = json.Marshal(r2)
	assert.NoError(t, err2)
	assert.Equal(t, `{"Name":"andeya"}`, string(b2))

	var r3 gust.Result[T]
	var err3 = json.Unmarshal(b2, &r3)
	assert.NoError(t, err3)
	assert.Equal(t, r2, r3)

	var r4 gust.Result[T]
	var err4 = json.Unmarshal([]byte("0"), &r4)
	assert.True(t, r4.IsErr())
	assert.Equal(t, "json: cannot unmarshal number into Go value of type gust_test.T", err4.Error())
}

func TestResultIsValid(t *testing.T) {
	var r0 *gust.Result[any]
	assert.False(t, r0.IsValid())
	var r1 gust.Result[any]
	assert.False(t, r1.IsValid())
	assert.False(t, (&gust.Result[any]{}).IsValid())
	var r2 = gust.Ok[any](nil)
	assert.True(t, r2.IsValid())
}
