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
	// &gust.errorWithVal{val:1}
}
