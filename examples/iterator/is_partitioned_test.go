package iterator_test

import (
	"testing"

	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestIsPartitioned(t *testing.T) {
	assert.True(t, iter.FromString[rune]("Iterator").IsPartitioned(func(r rune) bool {
		return r >= 'A' && r <= 'Z'
	}))
	assert.False(t, iter.FromString[rune]("IntoIterator").IsPartitioned(func(r rune) bool {
		return r >= 'A' && r <= 'Z'
	}))
}
