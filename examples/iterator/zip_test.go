package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestZip(t *testing.T) {
	var c = make(chan int, 3)
	c <- 4
	c <- 5
	c <- 6
	close(c)
	for _, x := range [][2]iter.Iterator[int]{
		{iter.FromElements(1, 2, 3), iter.FromElements(4, 5, 6)},
		{iter.FromRange(1, 4), iter.FromChan(c)},
	} {
		var i = iter.Zip(x[0], x[1])
		assert.Equal(t, gust.Some(gust.Pair[int, int]{A: 1, B: 4}), i.Next())
		assert.Equal(t, gust.Some(gust.Pair[int, int]{A: 2, B: 5}), i.Next())
		assert.Equal(t, gust.Some(gust.Pair[int, int]{A: 3, B: 6}), i.Next())
		assert.Equal(t, gust.None[gust.Pair[int, int]](), i.Next())
	}
}
