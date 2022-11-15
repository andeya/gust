package iter

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestAny(t *testing.T) {
	var iter = FromVec([]int{1, 2, 3})
	if !iter.Any(func(x int) bool {
		return x > 1
	}) {
		t.Error("Any failed")
	}
}

func TestNextChunk(t *testing.T) {
	var iter = FromVec([]int{1, 2, 3})
	assert.Equal(t, []int{1, 2}, iter.NextChunk(2).Unwrap())
	assert.Equal(t, []int{3}, iter.NextChunk(2).UnwrapErr())
	assert.Equal(t, []int{}, iter.NextChunk(2).UnwrapErr())
}

func TestZip(t *testing.T) {
	var a = FromVec([]string{"x", "y", "z"})
	var b = FromVec([]int{1, 2})
	var iter = ToZip[string, int](a, b)
	var pairs = Fold[gust.Pair[string, int]](iter, nil, func(acc []gust.Pair[string, int], t gust.Pair[string, int]) []gust.Pair[string, int] {
		return append(acc, t)
	})
	assert.Equal(t, []gust.Pair[string, int]{{"x", 1}, {"y", 2}}, pairs)
}

func TestToUnique(t *testing.T) {
	var data = FromElements(10, 20, 30, 20, 40, 10, 50)
	assert.Equal(t, []int{10, 20, 30, 40, 50}, ToUnique[int](data).Collect())
}

func TestToDeUnique(t *testing.T) {
	var data = FromElements(10, 20, 30, 20, 40, 10, 50)
	assert.Equal(t, []int{10, 20, 30, 40, 50}, ToDeUnique[int](data).Collect())
	var data2 = FromElements(10, 20, 30, 20, 40, 10, 50)
	assert.Equal(t, []int{50, 10, 40, 20, 30}, ToDeUnique[int](data2).ToRev().Collect())
	var data3 = FromElements(10, 20, 30, 20, 40, 10, 50)
	assert.Equal(t, []int{50, 10, 40, 20, 30}, ToDeUnique(data3.ToRev()).Collect())
}

func TestToUniqueBy(t *testing.T) {
	var data = FromElements("a", "bb", "aa", "c", "ccc")
	assert.Equal(t, []string{"a", "bb", "ccc"}, ToUniqueBy[string, int](data, func(s string) int { return len(s) }).Collect())
}

func TestToDeUniqueBy(t *testing.T) {
	var f = func(s string) int { return len(s) }
	var data = FromElements("a", "bb", "aa", "c", "ccc")
	assert.Equal(t, []string{"a", "bb", "ccc"}, ToDeUniqueBy[string, int](data, f).Collect())
	var data2 = FromElements("a", "bb", "aa", "c", "ccc")
	assert.Equal(t, []string{"ccc", "c", "aa"}, ToDeUniqueBy[string, int](data2, f).ToRev().Collect())
	var data3 = FromElements("a", "bb", "aa", "c", "ccc")
	assert.Equal(t, []string{"ccc", "c", "aa"}, ToDeUniqueBy[string, int](data3.ToRev(), f).Collect())
}
