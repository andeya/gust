package gust_test

import (
	"fmt"
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestErrable(t *testing.T) {
	assert.False(t, gust.ToErrable[any](nil).AsError())
	assert.False(t, gust.NonErrable[any]().AsError())

	assert.False(t, gust.ToErrable[error](nil).AsError())
	assert.False(t, gust.NonErrable[int]().AsError())

	assert.False(t, gust.ToErrable[*int](nil).AsError())
	assert.False(t, gust.NonErrable[*int]().AsError())
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
	fmt.Println(r.AsError())
	fmt.Println(r.Unwrap())
	fmt.Printf("%#v", r.ToError())
	// Output:
	// true
	// 1
	// &gust.errorWithVal{val:1}
}
