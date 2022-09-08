package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/digit"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func findMax[T digit.Integer](i iter.Iterator[T]) gust.Option[T] {
	return i.Reduce(func(acc T, v T) T {
		if acc >= v {
			return acc
		}
		return v
	})
}

func TestReduce(t *testing.T) {
	var a = []int{10, 20, 5, -23, 0}
	var b []uint
	assert.Equal(t, gust.Some(20), findMax[int](iter.FromVec(a)))
	assert.Equal(t, gust.None[uint](), findMax[uint](iter.FromVec(b)))
}
