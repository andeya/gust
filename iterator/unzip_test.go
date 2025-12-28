package iterator_test

import (
	"testing"

	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/pair"
	"github.com/stretchr/testify/assert"
)

func TestUnzip(t *testing.T) {
	a := []pair.Pair[int, string]{
		{A: 1, B: "a"},
		{A: 2, B: "b"},
		{A: 3, B: "c"},
	}
	left, right := iterator.Unzip(iterator.FromSlice(a))
	assert.Equal(t, []int{1, 2, 3}, left)
	assert.Equal(t, []string{"a", "b", "c"}, right)
}
