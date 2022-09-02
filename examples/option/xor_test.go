package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOption_Xor(t *testing.T) {
	{
		var x = gust.Some[uint32](2)
		var y gust.Option[uint32]
		assert.Equal(t, gust.Some[uint32](2), x.Xor(y))
	}
	{
		var x gust.Option[uint32]
		var y = gust.Some[uint32](100)
		assert.Equal(t, gust.Some[uint32](100), x.Xor(y))
	}
	{
		var x = gust.Some[uint32](2)
		var y = gust.Some[uint32](100)
		assert.Equal(t, gust.None[uint32](), x.Xor(y))
	}
	{
		var x gust.Option[uint32]
		var y = gust.None[uint32]()
		assert.Equal(t, gust.None[uint32](), x.Xor(y))
	}
}
