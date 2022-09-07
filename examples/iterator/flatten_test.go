package iterator_test

import (
	"testing"

	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	var i = iter.FromElements(
		iter.FromElements([]int{1, 2}, []int{3, 4}),
		iter.FromElements([]int{5, 6}, []int{7, 8}),
	)
	var d2 = iter.DeFlatten[iter.DeIterator[[]int], []int](i).Collect()
	assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}}, d2)
}
