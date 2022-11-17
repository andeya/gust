package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestEnumerate(t *testing.T) {
	var i = iter.EnumElements[rune]('a', 'b', 'c')
	assert.Equal(t, gust.Some(gust.VecEntry[rune]{
		Index: 0, Elem: 'a',
	}), i.Next())
	assert.Equal(t, gust.Some(gust.VecEntry[rune]{
		Index: 1, Elem: 'b',
	}), i.Next())
	assert.Equal(t, gust.Some(gust.VecEntry[rune]{
		Index: 2, Elem: 'c',
	}), i.Next())
	assert.Equal(t, gust.None[gust.VecEntry[rune]](), i.Next())
}
