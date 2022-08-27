package gust_test

import (
	"fmt"
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestErrable(t *testing.T) {
	assert.False(t, gust.ToErrable[any](nil).HasError())
	assert.False(t, gust.NonErrable[any]().HasError())

	assert.False(t, gust.ToErrable[error](nil).HasError())
	assert.False(t, gust.NonErrable[int]().HasError())

	assert.False(t, gust.ToErrable[*int](nil).HasError())
	assert.False(t, gust.NonErrable[*int]().HasError())
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
	fmt.Println(r.HasError())
	fmt.Println(r.Unwrap())
	fmt.Printf("%#v", r.ToError())
	// Output:
	// true
	// 1
	// &errors.errorString{s:"1"}
}
