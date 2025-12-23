package iter

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestUnzip(t *testing.T) {
	a := []gust.Pair[int, string]{
		{A: 1, B: "a"},
		{A: 2, B: "b"},
		{A: 3, B: "c"},
	}
	left, right := Unzip(FromSlice(a))
	assert.Equal(t, []int{1, 2, 3}, left)
	assert.Equal(t, []string{"a", "b", "c"}, right)
}

