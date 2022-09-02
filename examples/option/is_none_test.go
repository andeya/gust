package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOption_IsNone(t *testing.T) {
	{
		var x = gust.Some[uint32](2)
		assert.False(t, x.IsNone())
	}
	{
		var x = gust.None[uint32]()
		assert.True(t, x.IsNone())
	}
}
