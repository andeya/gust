package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/opt"
	"github.com/stretchr/testify/assert"
)

func TestOption_MapOr_1(t *testing.T) {
	{
		var x = gust.Some("foo")
		assert.Equal(t, 3, opt.MapOr(x, 42, func(v string) int { return len(v) }))
	}
	{
		var x gust.Option[string]
		assert.Equal(t, 42, opt.MapOr(x, 42, func(v string) int { return len(v) }))
	}
}

func TestOption_MapOr_2(t *testing.T) {
	{
		var x = gust.Some("foo")
		assert.Equal(t, 3, x.XMapOr(42, func(v string) any { return len(v) }))
	}
	{
		var x gust.Option[string]
		assert.Equal(t, 42, x.XMapOr(42, func(v string) any { return len(v) }))
	}
}
