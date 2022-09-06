package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOption_Take(t *testing.T) {
	{
		var x = gust.Some(2)
		var y = x.Take()
		assert.True(t, x.IsNone())
		assert.Equal(t, gust.Some(2), y)
	}
	{
		var x gust.Option[int] = gust.None[int]()
		var y = x.Take()
		assert.True(t, x.IsNone())
		assert.True(t, y.IsNone())
	}
}
