package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestAdvanceBy(t *testing.T) {
	var c = make(chan int, 4)
	c <- 1
	c <- 2
	c <- 3
	c <- 4
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(1, 2, 3, 4),
		iter.FromRange(1, 4, true),
		iter.FromChan(c),
	} {
		assert.Equal(t, gust.NonErrable[uint](), i.AdvanceBy(2))
		assert.Equal(t, gust.Some(3), i.Next())
		assert.Equal(t, gust.NonErrable[uint](), i.AdvanceBy(0))
		assert.Equal(t, gust.ToErrable[uint](1), i.AdvanceBy(100)) // only `4` was skipped
	}
}
