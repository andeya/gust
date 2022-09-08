package iterator_test

import (
	"strconv"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/andeya/gust/ret"
	"github.com/stretchr/testify/assert"
)

func TestTryFind(t *testing.T) {
	var a = []string{"1", "2", "lol", "NaN", "5"}
	var isMyNum = func(s string, search int) gust.Result[bool] {
		return ret.Map[int, bool](gust.Ret(strconv.Atoi(s)), func(x int) bool {
			return x == search
		})
	}
	var result = iter.FromVec[string](a).TryFind(func(s string) gust.Result[bool] {
		return isMyNum(s, 2)
	})
	assert.Equal(t, gust.Ok(gust.Some[string]("2")), result)
	result = iter.FromVec[string](a).TryFind(func(s string) gust.Result[bool] {
		return isMyNum(s, 5)
	})
	assert.True(t, result.IsErr())
}
