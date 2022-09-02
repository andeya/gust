package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOption_And_1(t *testing.T) {
	{
		var x = gust.Some[uint32](2)
		var y gust.Option[uint32]
		assert.Equal(t, gust.None[uint32](), x.And(y))
	}
	{
		var x gust.Option[uint32]
		var y = gust.Some[uint32](3)
		assert.Equal(t, gust.None[uint32](), x.And(y))
	}
	{
		var x = gust.Some[uint32](2)
		var y = gust.Some[uint32](3)
		assert.Equal(t, gust.Some[uint32](3), x.And(y))
	}
	{
		var x gust.Option[uint32]
		var y gust.Option[uint32]
		assert.Equal(t, gust.None[uint32](), x.And(y))
	}
}

func TestOption_And_2(t *testing.T) {
	{
		var x = gust.Some[uint32](2)
		var y gust.Option[string]
		assert.Equal(t, gust.None[string]().ToX(), x.XAnd(y.ToX()))
	}
	{
		var x gust.Option[uint32]
		var y = gust.Some[string]("foo")
		assert.Equal(t, gust.None[string]().ToX(), x.XAnd(y.ToX()))
	}
	{
		var x = gust.Some[uint32](2)
		var y = gust.Some[string]("foo")
		assert.Equal(t, gust.Some[string]("foo").ToX(), x.XAnd(y.ToX()))
	}
	{
		var x gust.Option[uint32]
		var y gust.Option[string]
		assert.Equal(t, gust.None[string]().ToX(), x.XAnd(y.ToX()))
	}
}
