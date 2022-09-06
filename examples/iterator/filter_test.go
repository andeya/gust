package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	var c = make(chan int, 10)
	c <- 0
	c <- 1
	c <- 2
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(0, 1, 2),
		iter.FromRange(0, 3),
		iter.FromChan(c),
	} {
		var i = i.Filter(func(v int) bool { return v > 0 })
		assert.Equal(t, gust.Some(1), i.Next())
		assert.Equal(t, gust.Some(2), i.Next())
		assert.Equal(t, gust.None[int](), i.Next())
	}
}
