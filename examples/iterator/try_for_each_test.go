package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestTryForEach(t *testing.T) {
	var c = make(chan int, 1000)
	for _, i := range []iter.Iterator[int]{
		iter.FromRange[int](2, 100).ToInspect(func(v int) {
			c <- v
		}),
		iter.FromChan(c),
	} {
		var r = i.TryForEach(func(v int) gust.AnyCtrlFlow {
			if 323%v == 0 {
				return gust.AnyBreak(v)
			}
			return gust.AnyContinue(nil)
		})
		assert.Equal(t, gust.AnyBreak(17), r)
	}
}
