package result_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestResult_IsOkAnd(t *testing.T) {
	{
		var x = gust.Ok[int](2)
		assert.True(t, x.IsOkAnd(func(x int) bool { return x > 1 }))
	}
	{
		var x = gust.Ok[int](0)
		assert.False(t, x.IsOkAnd(func(x int) bool { return x > 1 }))
	}
	{
		var x = gust.Err[int]("hey")
		assert.False(t, x.IsOkAnd(func(x int) bool { return x > 1 }))
	}
}
