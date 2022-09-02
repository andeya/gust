package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOption_IsSome(t *testing.T) {
	{
		var x = gust.Some[uint32](2)
		assert.True(t, x.IsSome())
	}
	{
		var x = gust.None[uint32]()
		assert.False(t, x.IsSome())
	}
}
