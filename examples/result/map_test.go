package result_test

import (
	"strconv"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/ret"
	"github.com/stretchr/testify/assert"
)

func TestResult_Map(t *testing.T) {
	var isMyNum = func(s string, search int) gust.Result[bool] {
		return ret.Map(gust.Ret(strconv.Atoi(s)), func(x int) bool { return x == search })
	}
	assert.Equal(t, gust.Ok[bool](true), isMyNum("1", 1))
	assert.Equal(t, "Err(strconv.Atoi: parsing \"lol\": invalid syntax)", isMyNum("lol", 1).String())
	assert.Equal(t, "Err(strconv.Atoi: parsing \"NaN\": invalid syntax)", isMyNum("NaN", 1).String())
}
