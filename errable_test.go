package gust

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrable(t *testing.T) {
	assert.False(t, ToErrable[any](nil).Ref().HasError())
	assert.False(t, NonErrable[any]().Ref().HasError())
	assert.False(t, (*Errable[any])(nil).HasError())

	assert.False(t, ToErrable[error](nil).Ref().HasError())
	assert.False(t, NonErrable[int]().Ref().HasError())
	assert.False(t, (*Errable[int])(nil).HasError())

	assert.False(t, ToErrable[*int](nil).Ref().HasError())
	assert.False(t, NonErrable[*int]().Ref().HasError())
	assert.False(t, (*Errable[*int])(nil).HasError())
}
