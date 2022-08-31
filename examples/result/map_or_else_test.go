package result_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/ret"
	"github.com/stretchr/testify/assert"
)

func TestResult_MapOrElse_1(t *testing.T) {
	var k = 21
	{
		var x = gust.Ok("foo")
		assert.Equal(t, 3, ret.MapOrElse(x, func(err error) int {
			return k * 2
		}, func(x string) int { return len(x) }))
	}
	{
		var x = gust.Err[string]("bar")
		assert.Equal(t, 42, ret.MapOrElse(x, func(err error) int {
			return k * 2
		}, func(x string) int { return len(x) }))
	}
}

func TestResult_MapOrElse_2(t *testing.T) {
	{
		var x = gust.Ok("foo")
		assert.Equal(t, "test:foo", x.MapOrElse(func(err error) any {
			return "bar"
		}, func(x string) any { return "test:" + x }).(string))
	}
	{
		var x = gust.Err[string]("foo")
		assert.Equal(t, "bar", x.MapOrElse(func(err error) any {
			return "bar"
		}, func(x string) any { return "test:" + x }).(string))
	}
}
