package gust_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/opt"
	"github.com/stretchr/testify/assert"
)

func ExampleOption() {
	type A struct {
		X int
	}
	var a = gust.Some(A{X: 1})
	fmt.Println(a.IsSome(), a.IsNone())

	var b = gust.None[A]()
	fmt.Println(b.IsSome(), b.IsNone())

	var x = b.UnwrapOr(A{X: 2})
	fmt.Println(x)

	var c *A
	fmt.Println(gust.Ptr(c).IsNone())
	c = new(A)
	fmt.Println(gust.Ptr(c).IsNone())

	type B struct {
		Y string
	}
	var d = opt.Map(a, func(t A) B {
		return B{
			Y: strconv.Itoa(t.X),
		}
	})
	fmt.Println(d)

	// Output:
	// true false
	// false true
	// {2}
	// true
	// false
	// Some({1})
}

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
