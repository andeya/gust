package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/digit"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestTryReduce_1(t *testing.T) {
	// Safely calculate the sum of a series of numbers:
	var numbers = []uint{10, 20, 5, 23, 0}
	var sum = iter.FromVec(numbers).TryReduce(func(x, y uint) gust.Result[uint] {
		return digit.CheckedAdd(x, y).OkOr("overflow")
	})
	assert.Equal(t, gust.Ok(gust.Some[uint](58)), sum)
}

func TestTryReduce_2(t *testing.T) {
	// Determine when a reduction short circuited:
	var numbers = []uint{1, 2, 3, ^uint(0), 4, 5}
	var sum = iter.FromVec(numbers).TryReduce(func(x, y uint) gust.Result[uint] {
		return digit.CheckedAdd(x, y).OkOr("overflow")
	})
	assert.Equal(t, gust.Err[gust.Option[uint]]("overflow"), sum)
}

func TestTryReduce_3(t *testing.T) {
	// Determine when a reduction was not performed because there are no elements:
	var numbers = []uint{}
	var sum = iter.FromVec(numbers).TryReduce(func(x, y uint) gust.Result[uint] {
		return digit.CheckedAdd(x, y).OkOr("overflow")
	})
	assert.Equal(t, gust.Ok(gust.None[uint]()), sum)
}
