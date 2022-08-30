package result_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/ret"
	"github.com/stretchr/testify/assert"
)

func TestResult_MapOr_1(t *testing.T) {
	{
		var x = gust.Ok("foo")
		assert.Equal(t, 3, ret.MapOr(x, 42, func(x string) int { return len(x) }))
	}
	{
		var x = gust.Err[string]("foo")
		assert.Equal(t, 42, ret.MapOr(x, 42, func(x string) int { return len(x) }))
	}
}

func TestResult_MapOr_2(t *testing.T) {
	{
		var x = gust.Ok("foo")
		assert.Equal(t, "test:foo", x.MapOr("bar", func(x string) string { return "test:" + x }))
	}
	{
		var x = gust.Err[string]("foo")
		assert.Equal(t, "bar", x.MapOr("bar", func(x string) string { return "test:" + x }))
	}
}
