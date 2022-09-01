package result_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestResult_OrElse(t *testing.T) {
	var sq = func(x int) gust.EnumResult[int, int] {
		return gust.EnumOk[int, int](x * x)
	}
	var err = func(x int) gust.EnumResult[int, int] {
		return gust.EnumErr[int, int](x)
	}

	assert.Equal(t, gust.EnumOk[int, int](2).OrElse(sq).OrElse(sq), gust.EnumOk[int, int](2))
	assert.Equal(t, gust.EnumOk[int, int](2).OrElse(err).OrElse(sq), gust.EnumOk[int, int](2))
	assert.Equal(t, gust.EnumErr[int, int](3).OrElse(sq).OrElse(err), gust.EnumOk[int, int](9))
	assert.Equal(t, gust.EnumErr[int, int](3).OrElse(err).OrElse(err), gust.EnumErr[int, int](3))
}
