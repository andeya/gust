package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestNth(t *testing.T) {
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
		assert.Equal(t, gust.Some(2), i.Nth(1))
		// Calling `Nth()` multiple times doesn't rewind the iterator:
		assert.Equal(t, gust.None[int](), i.Nth(1))
		// Returning `None` if there are less than `n + 1` elements:
		assert.Equal(t, gust.None[int](), i.Nth(10))
	}
}
