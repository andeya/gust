package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestSkipWhile(t *testing.T) {
	var c = make(chan int, 10)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(-1, 0, 1).ToInspect(func(v int) {
			c <- v
		}),
		iter.FromRange(-1, 2),
		iter.FromChan(c),
	} {
		var iter = i.ToSkipWhile(func(v int) bool {
			return v < 0
		})
		assert.Equal(t, gust.Some(0), iter.Next())
		assert.Equal(t, gust.Some(1), iter.Next())
		assert.Equal(t, gust.None[int](), iter.Next())
	}

}
