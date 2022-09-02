package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/opt"
	"github.com/stretchr/testify/assert"
)

func TestOption_Contains(t *testing.T) {
	{
		var x = gust.Some(2)
		assert.Equal(t, true, opt.Contains(x, 2))
	}
	{
		var x = gust.Some(3)
		assert.Equal(t, false, opt.Contains(x, 2))
	}
	{
		var x = gust.None[int]()
		assert.Equal(t, false, opt.Contains(x, 2))
	}
}
