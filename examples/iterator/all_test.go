package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestAll_1(t *testing.T) {
	var a = []int8{1, 2, 3}
	var all = iter.FromVec[int8](a).All(func(v int8) bool {
		return v > 0
	})
	assert.True(t, all)
	all = iter.FromVec[int8](a).All(func(v int8) bool {
		return v > 2
	})
	assert.False(t, all)
}

func TestAll_2(t *testing.T) {
	var a = []int8{1, 2, 3}
	var i = iter.FromVec[int8](a)
	var all = i.All(func(v int8) bool {
		return v != 2
	})
	assert.False(t, all)
	// we can still use `i`, as there are more elements.
	var next = i.Next()
	assert.Equal(t, gust.Some(int8(3)), next)
}
