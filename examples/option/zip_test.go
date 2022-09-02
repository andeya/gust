package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/opt"
	"github.com/stretchr/testify/assert"
)

func TestOption_Zip(t *testing.T) {
	var x = gust.Some[byte](1)
	var y = gust.Some("hi")
	var z = gust.None[byte]()
	assert.Equal(t, gust.Some(gust.Pair[byte, string]{1, "hi"}), opt.Zip(x, y))
	assert.Equal(t, gust.None[gust.Pair[byte, byte]](), opt.Zip(x, z))
}
