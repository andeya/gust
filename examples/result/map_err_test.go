package result_test

import (
	"fmt"
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestResult_MapErr(t *testing.T) {
	var stringify = func(x error) any { return fmt.Sprintf("error code: %v", x) }
	{
		var x = gust.Ok[uint32](2)
		assert.Equal(t, gust.Ok[uint32](2), x.MapErr(stringify))
	}
	{
		var x = gust.Err[uint32](13)
		assert.Equal(t, gust.Err[uint32]("error code: 13"), x.MapErr(stringify))
	}
}
