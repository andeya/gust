package result_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestResult_Ok(t *testing.T) {
	{
		var x = gust.Ok[int](2)
		assert.Equal(t, gust.Some[int](2), x.Ok())
	}
	{
		var x = gust.Err[int]("some error message")
		assert.Equal(t, gust.None[int](), x.Ok())
	}
}
