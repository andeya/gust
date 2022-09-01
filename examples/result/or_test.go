package result_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestResult_Or(t *testing.T) {
	{
		x := gust.Ok(2)
		y := gust.Err[int]("late error")
		assert.Equal(t, gust.Ok(2), x.Or(y))
	}
	{
		x := gust.Err[uint]("early error")
		y := gust.Ok[uint](2)
		assert.Equal(t, gust.Ok[uint](2), x.Or(y))
	}
	{
		x := gust.Err[uint]("not a 2")
		y := gust.Err[uint]("late error")
		assert.Equal(t, gust.Err[uint]("late error"), x.Or(y))
	}
	{
		x := gust.Ok[uint](2)
		y := gust.Ok[uint](100)
		assert.Equal(t, gust.Ok[uint](2), x.Or(y))
	}
}
