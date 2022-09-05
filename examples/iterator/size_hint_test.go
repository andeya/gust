package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestSizeHint_1(t *testing.T) {
	var i = iter.FromElements(1, 2, 3)
	var lo, hi = i.SizeHint()
	assert.Equal(t, uint(3), lo)
	assert.Equal(t, gust.Some[uint](3), hi)
	_ = i.Next()
	lo, hi = i.SizeHint()
	assert.Equal(t, uint(2), lo)
	assert.Equal(t, gust.Some[uint](2), hi)
}

func TestSizeHint_2(t *testing.T) {
	var i = iter.FromRange(1, 3, true)
	var lo, hi = i.SizeHint()
	assert.Equal(t, uint(3), lo)
	assert.Equal(t, gust.Some[uint](3), hi)
	_ = i.Next()
	lo, hi = i.SizeHint()
	assert.Equal(t, uint(2), lo)
	assert.Equal(t, gust.Some[uint](2), hi)
}

func TestSizeHint_3(t *testing.T) {
	var c = make(chan int, 3)
	c <- 1
	c <- 2
	c <- 3
	var i = iter.FromChan(c)
	var lo, hi = i.SizeHint()
	assert.Equal(t, uint(3), lo)
	assert.Equal(t, gust.Some[uint](3), hi)
	_ = i.Next()
	lo, hi = i.SizeHint()
	assert.Equal(t, uint(2), lo)
	assert.Equal(t, gust.Some[uint](3), hi)
}
