package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestSkip(t *testing.T) {
	var c = make(chan int, 10)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(1, 2, 3).Inspect(func(v int) {
			c <- v
		}),
		iter.FromRange(1, 4),
		iter.FromChan(c),
	} {
		var iter = i.Skip(2)
		assert.Equal(t, gust.Some(3), iter.Next())
		assert.Equal(t, gust.None[int](), iter.Next())
	}
}
