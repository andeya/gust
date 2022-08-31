package result_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/ret"
	"github.com/stretchr/testify/assert"
)

func TestResult_Map_1(t *testing.T) {
	var line = "1\n2\n3\n4\n"
	for _, num := range strings.Split(line, "\n") {
		gust.Ret(strconv.Atoi(num)).Map(func(i int) any {
			return i * 2
		}).Inspect(func(i any) {
			t.Log(i)
		})
	}
}

func TestResult_Map_2(t *testing.T) {
	var isMyNum = func(s string, search int) gust.Result[bool] {
		return ret.Map(gust.Ret(strconv.Atoi(s)), func(x int) bool { return x == search })
	}
	assert.Equal(t, gust.Ok[bool](true), isMyNum("1", 1))
	assert.Equal(t, "Err(strconv.Atoi: parsing \"lol\": invalid syntax)", isMyNum("lol", 1).String())
	assert.Equal(t, "Err(strconv.Atoi: parsing \"NaN\": invalid syntax)", isMyNum("NaN", 1).String())
}
