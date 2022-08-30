package result_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestResult_IsErrAnd(t *testing.T) {
	{
		var x = gust.Err[int]("hey")
		assert.True(t, x.IsErrAnd(func(err error) bool { return err.Error() == "hey" }))
	}
	{
		var x = gust.Ok[int](2)
		assert.False(t, x.IsErrAnd(func(err error) bool { return err.Error() == "hey" }))
	}
}

func TestEnumResult_IsErrAnd(t *testing.T) {
	{
		var x = gust.EnumErr[int, int8](-1)
		assert.True(t, x.IsErrAnd(func(x int8) bool { return x == -1 }))
	}
}
