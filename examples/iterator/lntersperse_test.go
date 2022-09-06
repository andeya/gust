package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestIntersperse_1(t *testing.T) {
	var c = make(chan int, 4)
	c <- 0
	c <- 1
	c <- 2
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(0, 1, 2),
		iter.FromRange(0, 3),
		iter.FromChan(c),
	} {
		i := i.Intersperse(100)
		assert.Equal(t, gust.Some(0), i.Next())     // The first element from `a`.
		assert.Equal(t, gust.Some(100), i.Next())   // The separator.
		assert.Equal(t, gust.Some(1), i.Next())     // The next element from `a`.
		assert.Equal(t, gust.Some(100), i.Next())   // The separator.
		assert.Equal(t, gust.Some(2), i.Next())     // The last element from `a`.
		assert.Equal(t, gust.None[int](), i.Next()) // The iterator is finished.

	}
}

func TestIntersperse_2(t *testing.T) {
	var c = make(chan int, 4)
	c <- 0
	c <- 1
	c <- 2
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(0, 1, 2),
		iter.FromRange(0, 3),
		iter.FromChan(c),
	} {
		i := i.Peekable().Intersperse(100)
		assert.Equal(t, gust.Some(0), i.Next())     // The first element from `a`.
		assert.Equal(t, gust.Some(100), i.Next())   // The separator.
		assert.Equal(t, gust.Some(1), i.Next())     // The next element from `a`.
		assert.Equal(t, gust.Some(100), i.Next())   // The separator.
		assert.Equal(t, gust.Some(2), i.Next())     // The last element from `a`.
		assert.Equal(t, gust.None[int](), i.Next()) // The iterator is finished.
	}
}

func TestIntersperse_3(t *testing.T) {
	var c = make(chan int, 4)
	c <- 0
	c <- 1
	c <- 2
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(0, 1, 2),
		iter.FromRange(0, 3),
		iter.FromChan(c),
	} {
		i := i.IntersperseWith(func() int { return 100 })
		assert.Equal(t, gust.Some(0), i.Next())     // The first element from `a`.
		assert.Equal(t, gust.Some(100), i.Next())   // The separator.
		assert.Equal(t, gust.Some(1), i.Next())     // The next element from `a`.
		assert.Equal(t, gust.Some(100), i.Next())   // The separator.
		assert.Equal(t, gust.Some(2), i.Next())     // The last element from `a`.
		assert.Equal(t, gust.None[int](), i.Next()) // The iterator is finished.
	}
}
