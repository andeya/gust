package iterator_test

import (
	"testing"

	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestCount_1(t *testing.T) {
	assert.Equal(t, uint(3), iter.FromElements(1, 2, 3).Count())
	assert.Equal(t, uint(5), iter.FromElements(1, 2, 3, 4, 5).Count())
}

func TestCount_2(t *testing.T) {
	assert.Equal(t, uint(3), iter.FromRange(1, 3, true).Count())
	assert.Equal(t, uint(5), iter.FromRange(1, 6, false).Count())
	assert.Equal(t, uint(5), iter.FromRange(1, 6).Count())
}

func TestCount_3(t *testing.T) {
	var c = make(chan int, 3)
	c <- 1
	c <- 2
	c <- 3
	assert.Equal(t, uint(3), iter.FromChan(c).Count())
}
