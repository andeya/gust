package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOption_Replace(t *testing.T) {
	{
		var x = gust.Some(2)
		var old = x.Replace(5)
		assert.Equal(t, gust.Some(5), x)
		assert.Equal(t, gust.Some(2), old)
	}
	{
		var x = gust.None[int]()
		var old = x.Replace(3)
		assert.Equal(t, gust.Some(3), x)
		assert.Equal(t, gust.None[int](), old)
	}
}
