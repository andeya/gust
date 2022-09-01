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
