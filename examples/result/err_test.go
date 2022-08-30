package result_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestResult_Err(t *testing.T) {
	{
		var x = gust.Ok[int](2)
		assert.Equal(t, error(nil), x.Err())
	}
	{
		var x = gust.Err[int]("some error message")
		assert.Equal(t, "some error message", x.Err().Error())
	}
}

func TestEnumResult_Err(t *testing.T) {
	{
		var x = gust.EnumOk[int, string](2)
		assert.Equal(t, gust.None[string](), x.Err())
	}
	{
		var x = gust.EnumErr[int, string]("some error message")
		assert.Equal(t, gust.Some[string]("some error message"), x.Err())
	}
}
