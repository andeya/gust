package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestPosition_1(t *testing.T) {
	var a = []int8{1, 2, 3}
	assert.Equal(t, gust.Some(1), iter.FromVec[int8](a).Position(func(v int8) bool {
		return v == 2
	}))
	assert.Equal(t, gust.None[int](), iter.FromVec[int8](a).Position(func(v int8) bool {
		return v == 5
	}))
}

func TestPosition_2(t *testing.T) {
	// Stopping at the first `true`:
	var a = iter.FromElements(1, 2, 3, 4)
	assert.Equal(t, gust.Some(1), a.Position(func(v int) bool {
		return v >= 2
	}))
	// we can still use `iter`, as there are more elements.
	assert.Equal(t, gust.Some(3), a.Next())
	// The returned index depends on iterator state
	assert.Equal(t, gust.Some(0), a.Position(func(v int) bool {
		return v == 4
	}))
}
