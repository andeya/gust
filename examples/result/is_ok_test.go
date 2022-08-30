package result_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestResult_IsOk(t *testing.T) {
	{
		var x = gust.Ok[int](-3)
		assert.True(t, x.IsOk())
	}
	{
		var x = gust.Err[int]("some error message")
		assert.False(t, x.IsOk())
	}
}
