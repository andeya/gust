package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestLast_1(t *testing.T) {
	var a = []int{1, 2, 3}
	var i = iter.FromVec(a)
	assert.Equal(t, gust.Some(3), i.Last())
	assert.Equal(t, gust.None[int](), i.Last())
	assert.Equal(t, gust.None[int](), i.Next())
}

func TestLast_2(t *testing.T) {
	var i = iter.FromRange(1, 3, true)
	assert.Equal(t, gust.Some(3), i.Last())
	assert.Equal(t, gust.None[int](), i.Last())
	assert.Equal(t, gust.None[int](), i.Next())
}

func TestLast_3(t *testing.T) {
	var c = make(chan int, 3)
	c <- 1
	c <- 2
	c <- 3
	close(c)
	var i = iter.FromChan(c)
	assert.Equal(t, gust.Some(3), i.Last())
	assert.Equal(t, gust.None[int](), i.Last())
	assert.Equal(t, gust.None[int](), i.Next())
}
