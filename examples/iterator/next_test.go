package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestNext_1(t *testing.T) {
	var a = []int{1, 2, 3}
	var i = iter.FromVec(a)
	// A call to Next() returns the next value...
	assert.Equal(t, gust.Some(1), i.Next())
	assert.Equal(t, gust.Some(2), i.Next())
	assert.Equal(t, gust.Some(3), i.Next())
	// ... and then None once it's over.
	assert.Equal(t, gust.None[int](), i.Next())
	// More calls may or may not return `None`. Here, they always will.
	assert.Equal(t, gust.None[int](), i.Next())
	assert.Equal(t, gust.None[int](), i.Next())
}

func TestNext_2(t *testing.T) {
	var i = iter.FromRange(1, 3, true)
	// A call to Next() returns the next value...
	assert.Equal(t, gust.Some(1), i.Next())
	assert.Equal(t, gust.Some(2), i.Next())
	assert.Equal(t, gust.Some(3), i.Next())
	// ... and then None once it's over.
	assert.Equal(t, gust.None[int](), i.Next())
	// More calls may or may not return `None`. Here, they always will.
	assert.Equal(t, gust.None[int](), i.Next())
	assert.Equal(t, gust.None[int](), i.Next())
}

func TestNext_3(t *testing.T) {
	var c = make(chan int, 3)
	c <- 1
	c <- 2
	c <- 3
	close(c)
	var i = iter.FromChan(c)
	// A call to Next() returns the next value...
	assert.Equal(t, gust.Some(1), i.Next())
	assert.Equal(t, gust.Some(2), i.Next())
	assert.Equal(t, gust.Some(3), i.Next())
	// ... and then None once it's over.
	assert.Equal(t, gust.None[int](), i.Next())
	// More calls may or may not return `None`. Here, they always will.
	assert.Equal(t, gust.None[int](), i.Next())
	assert.Equal(t, gust.None[int](), i.Next())
}
