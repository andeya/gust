package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/opt"
	"github.com/stretchr/testify/assert"
)

func TestOption_MapOrElse_1(t *testing.T) {
	var k = 21
	{
		var x = gust.Some("foo")
		assert.Equal(t, 3, opt.MapOrElse(x, func() int { return 2 * k }, func(v string) int { return len(v) }))
	}
	{
		var x gust.Option[string]
		assert.Equal(t, 42, opt.MapOrElse(x, func() int { return 2 * k }, func(v string) int { return len(v) }))
	}
}

func TestOption_MapOrElse_2(t *testing.T) {
	var k = 21
	{
		var x = gust.Some("foo")
		assert.Equal(t, 3, x.XMapOrElse(func() any { return 2 * k }, func(v string) any { return len(v) }))
	}
	{
		var x gust.Option[string]
		assert.Equal(t, 42, x.XMapOrElse(func() any { return 2 * k }, func(v string) any { return len(v) }))
	}
}
