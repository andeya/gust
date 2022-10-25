package gust_test

import (
	"encoding/json"
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOptionJSON(t *testing.T) {
	var r = gust.None[any]()
	var b, err = json.Marshal(r)
	a, _ := json.Marshal(nil)
	assert.Equal(t, a, b)
	assert.NoError(t, err)
	type T struct {
		Name string
	}
	var r2 = gust.Some(T{Name: "andeya"})
	var b2, err2 = json.Marshal(r2)
	assert.NoError(t, err2)
	assert.Equal(t, `{"Name":"andeya"}`, string(b2))

	var r3 gust.Option[T]
	var err3 = json.Unmarshal(b2, &r3)
	assert.NoError(t, err3)
	assert.Equal(t, r2, r3)

	var r4 gust.Option[T]
	var err4 = json.Unmarshal([]byte("0"), &r4)
	assert.True(t, r4.IsNone())
	assert.Equal(t, "json: cannot unmarshal number into Go value of type gust_test.T", err4.Error())
}

func TestOptionJSON2(t *testing.T) {
	type T struct {
		Name gust.Option[string]
	}
	var r = T{Name: gust.Some("andeya")}
	var b, err = json.Marshal(r)
	assert.NoError(t, err)
	assert.Equal(t, `{"Name":"andeya"}`, string(b))
	var r2 T
	err2 := json.Unmarshal(b, &r2)
	assert.NoError(t, err2)
	assert.Equal(t, r, r2)

	var r3 = T{Name: gust.Some("")}
	var b3, err3 = json.Marshal(r3)
	assert.NoError(t, err3)
	assert.Equal(t, `{"Name":""}`, string(b3))
	var r4 T
	err4 := json.Unmarshal(b3, &r4)
	assert.NoError(t, err4)
	assert.Equal(t, r3, r4)

	var r5 = T{Name: gust.None[string]()}
	var b5, err5 = json.Marshal(r5)
	assert.NoError(t, err5)
	assert.Equal(t, `{"Name":null}`, string(b5))
	var r6 T
	err6 := json.Unmarshal(b5, &r6)
	assert.NoError(t, err6)
	assert.Equal(t, r5, r6)
}
