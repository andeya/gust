package result_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestResult_IsErr(t *testing.T) {
	{
		var x = gust.Ok[int](-3)
		assert.False(t, x.IsErr())
	}
	{
		var x = gust.Err[int]("some error message")
		assert.True(t, x.IsErr())
	}
}
