package gust

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrable(t *testing.T) {
	assert.False(t, ToErrable[any](nil).HasError())
	assert.False(t, NonErrable[any]().HasError())

	assert.False(t, ToErrable[error](nil).HasError())
	assert.False(t, NonErrable[int]().HasError())

	assert.False(t, ToErrable[*int](nil).HasError())
	assert.False(t, NonErrable[*int]().HasError())
}
