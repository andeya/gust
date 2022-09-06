package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestEnumerate(t *testing.T) {
	var i = iter.EnumElements[rune]('a', 'b', 'c')
	assert.Equal(t, gust.Some(iter.KV[rune]{
		Index: 0, Value: 'a',
	}), i.Next())
	assert.Equal(t, gust.Some(iter.KV[rune]{
		Index: 1, Value: 'b',
	}), i.Next())
	assert.Equal(t, gust.Some(iter.KV[rune]{
		Index: 2, Value: 'c',
	}), i.Next())
	assert.Equal(t, gust.None[iter.KV[rune]](), i.Next())
}
