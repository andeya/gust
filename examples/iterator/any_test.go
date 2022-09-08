package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestAny_1(t *testing.T) {
	var a = []int8{1, 2, 3}
	var hasAny = iter.FromVec[int8](a).Any(func(v int8) bool {
		return v > 0
	})
	assert.True(t, hasAny)
	hasAny = iter.FromVec[int8](a).Any(func(v int8) bool {
		return v > 5
	})
	assert.False(t, hasAny)
}

func TestAny_2(t *testing.T) {
	var a = []int8{1, 2, 3}
	var i = iter.FromVec[int8](a)
	var hasAny = i.Any(func(v int8) bool {
		return v != 2
	})
	assert.True(t, hasAny)
	// we can still use `i`, as there are more elements.
	var next = i.Next()
	assert.Equal(t, gust.Some(int8(2)), next)
}
