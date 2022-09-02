package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOption_OkOr(t *testing.T) {
	{
		var x = gust.Some("foo")
		assert.Equal(t, gust.Ok("foo"), x.OkOr(0))
	}
	{
		var x gust.Option[string]
		assert.Equal(t, gust.Err[string](0), x.OkOr(0))
	}
}
