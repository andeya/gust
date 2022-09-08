package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestFind_1(t *testing.T) {
	var a = []int8{1, 2, 3}
	var i = iter.FromVec[int8](a)
	assert.Equal(t, gust.Some[int8](2), i.Find(func(v int8) bool {
		return v == 2
	}))
	assert.Equal(t, gust.None[int8](), i.Find(func(v int8) bool {
		return v == 5
	}))
}

func TestFind_2(t *testing.T) {
	var a = []int8{1, 2, 3}
	var i = iter.FromVec[int8](a)
	// / assert_eq!(iter.find(|&&x| x == 2), Some(&2));
	assert.Equal(t, gust.Some[int8](2), i.Find(func(v int8) bool {
		return v == 2
	}))
	// we can still use `iter`, as there are more elements.
	assert.Equal(t, gust.Some[int8](3), i.Next())
}
