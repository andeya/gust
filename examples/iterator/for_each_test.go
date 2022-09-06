package iterator_test

import (
	"testing"

	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestForEach_1(t *testing.T) {
	var c = make(chan int, 5)
	iter.FromRange(0, 5).
		Map(func(i int) int { return i*2 + 1 }).
		ForEach(func(i int) {
			c <- i
		})
	var v = iter.FromChan(c).Collect()
	assert.Equal(t, []int{1, 3, 5, 7, 9}, v)
}

func TestForEach_2(t *testing.T) {
	var c = make(chan int)
	go func() {
		iter.FromRange(0, 5).
			Map(func(i int) int { return i*2 + 1 }).
			ForEach(func(i int) {
				c <- i
			})
		close(c)
	}()
	var v = iter.FromChan(c).Collect()
	assert.Equal(t, []int{1, 3, 5, 7, 9}, v)
}
