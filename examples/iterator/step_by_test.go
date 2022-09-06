package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestStepBy(t *testing.T) {
	var c = make(chan int, 6)
	c <- 0
	c <- 1
	c <- 2
	c <- 3
	c <- 4
	c <- 5
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(0, 1, 2, 3, 4, 5),
		iter.FromRange(0, 5, true),
		iter.FromChan(c),
	} {
		var stepBy = i.StepBy(2)
		assert.Equal(t, gust.Some(0), stepBy.Next())
		assert.Equal(t, gust.Some(2), stepBy.Next())
		assert.Equal(t, gust.Some(4), stepBy.Next())
		assert.Equal(t, gust.None[int](), stepBy.Next())
	}
}
