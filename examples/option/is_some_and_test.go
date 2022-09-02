package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOption_IsSomeAnd(t *testing.T) {
	{
		var x = gust.Some[uint32](2)
		assert.True(t, x.IsSomeAnd(func(v uint32) bool { return v > 1 }))
	}
	{
		var x = gust.Some[uint32](0)
		assert.False(t, x.IsSomeAnd(func(v uint32) bool { return v > 1 }))
	}
	{
		var x = gust.None[uint32]()
		assert.False(t, x.IsSomeAnd(func(v uint32) bool { return v > 1 }))
	}
}
