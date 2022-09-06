package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestPeekable(t *testing.T) {
	var i = iter.FromElements(1, 2, 3).Peekable()
	assert.Equal(t, gust.Some(1), i.Peek())
	assert.Equal(t, gust.Some(1), i.Peek())
	assert.Equal(t, gust.Some(1), i.Next())
	peeked := i.Peek()
	if peeked.IsSome() {
		p := peeked.GetOrInsertWith(nil)
		assert.Equal(t, 2, *p)
		*p = 1000
	}
	// The value reappears as the iterator continues
	assert.Equal(t, []int{1000, 3}, i.Collect())
}
