package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestScan(t *testing.T) {
	var c = make(chan int, 10)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(1, 2, 3).Inspect(func(v int) {
			c <- v
		}),
		iter.FromRange(1, 4),
		iter.FromChan(c),
	} {
		j := i.Scan(1, func(state *any, x int) gust.Option[any] {
			// each iteration, we'll multiply the state by the element
			*state = (*state).(int) * x
			// then, we'll yield the negation of the state
			return gust.Some[any](-(*state).(int))
		})
		assert.Equal(t, gust.Some[any](-1), j.Next())
		assert.Equal(t, gust.Some[any](-2), j.Next())
		assert.Equal(t, gust.Some[any](-6), j.Next())
		assert.Equal(t, gust.None[any](), j.Next())
	}
}
