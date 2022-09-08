package iterator_test

import (
	"testing"

	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestPartition(t *testing.T) {
	var c = make(chan int, 10)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(1, 2, 3).ToInspect(func(v int) {
			c <- v
		}),
		iter.FromRange(1, 4),
		iter.FromChan(c),
	} {
		var even, odd = i.Partition(func(x int) bool { return x%2 == 0 })
		assert.Equal(t, []int{2}, even)
		assert.Equal(t, []int{1, 3}, odd)
	}
}
