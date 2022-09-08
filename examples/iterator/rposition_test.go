package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestRposition_1(t *testing.T) {
	var a = []int8{1, 2, 3}
	assert.Equal(t, gust.Some(2), iter.FromVec[int8](a).Rposition(func(v int8) bool {
		return v == 3
	}))
	assert.Equal(t, gust.None[int](), iter.FromVec[int8](a).Rposition(func(v int8) bool {
		return v == 5
	}))
}

func TestRposition_2(t *testing.T) {
	var a = iter.FromElements(1, 2, 3)
	assert.Equal(t, gust.Some(1), a.Rposition(func(v int) bool {
		return v == 2
	}))
	assert.Equal(t, gust.Some(1), a.Next())
}
