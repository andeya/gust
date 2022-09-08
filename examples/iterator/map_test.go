package iterator_test

import (
	"fmt"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestMap_1(t *testing.T) {
	var c = make(chan int, 4)
	c <- 1
	c <- 2
	c <- 3
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(1, 2, 3),
		iter.FromRange(1, 4),
		iter.FromChan(c),
	} {
		i := i.ToMap(func(v int) int { return v * 2 })
		assert.Equal(t, gust.Some(2), i.Next())
		assert.Equal(t, gust.Some(4), i.Next())
		assert.Equal(t, gust.Some(6), i.Next())
		assert.Equal(t, gust.None[int](), i.Next())
	}
}

func TestMap_2(t *testing.T) {
	var c = make(chan int, 4)
	c <- 1
	c <- 2
	c <- 3
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(1, 2, 3),
		iter.FromRange(1, 4),
		iter.FromChan(c),
	} {
		i := i.ToXMap(func(v int) any { return fmt.Sprintf("%d", v*2) })
		assert.Equal(t, gust.Some[any]("2"), i.Next())
		assert.Equal(t, gust.Some[any]("4"), i.Next())
		assert.Equal(t, gust.Some[any]("6"), i.Next())
		assert.Equal(t, gust.None[any](), i.Next())
	}
}

func TestMap_3(t *testing.T) {
	var c = make(chan int, 4)
	c <- 1
	c <- 2
	c <- 3
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(1, 2, 3),
		iter.FromRange(1, 4),
		iter.FromChan(c),
	} {
		i := iter.ToMap(i, func(v int) string { return fmt.Sprintf("%d", v*2) })
		assert.Equal(t, gust.Some[string]("2"), i.Next())
		assert.Equal(t, gust.Some[string]("4"), i.Next())
		assert.Equal(t, gust.Some[string]("6"), i.Next())
		assert.Equal(t, gust.None[string](), i.Next())
	}
}
