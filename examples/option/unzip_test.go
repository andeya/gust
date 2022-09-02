package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/opt"
	"github.com/stretchr/testify/assert"
)

func TestOption_Unzip(t *testing.T) {
	var x = gust.Some[gust.Pair[int, string]](gust.Pair[int, string]{A: 1, B: "hi"})
	var y = gust.None[gust.Pair[int, string]]()
	assert.Equal(t, opt.Unzip(x), gust.Pair[gust.Option[int], gust.Option[string]]{A: gust.Some[int](1), B: gust.Some[string]("hi")})
	assert.Equal(t, opt.Unzip(y), gust.Pair[gust.Option[int], gust.Option[string]]{A: gust.None[int](), B: gust.None[string]()})
}
