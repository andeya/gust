package result_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestResult_And(t *testing.T) {
	{
		x := gust.Ok(2)
		y := gust.Err[int]("late error")
		assert.Equal(t, gust.Err[int]("late error"), x.And(y))
	}
	{
		x := gust.Err[uint]("early error")
		y := gust.Ok[string]("foo")
		assert.Equal(t, gust.Err[any]("early error"), x.XAnd(y.ToX()))
	}
}
