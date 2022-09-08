package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestRev(t *testing.T) {
	for _, i := range []iter.DeIterator[int]{
		iter.FromElements(1, 2, 3),
		iter.FromRange(1, 3, true),
	} {
		var rev = i.ToRev()
		assert.Equal(t, gust.Some(3), rev.Next())
		assert.Equal(t, gust.Some(2), rev.Next())
		assert.Equal(t, gust.Some(1), rev.Next())
		assert.Equal(t, gust.None[int](), rev.Next())
	}
}
