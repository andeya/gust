package gust_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOrdering(t *testing.T) {
	// Test Less
	less := gust.Less()
	assert.True(t, less.IsLess())
	assert.False(t, less.IsEqual())
	assert.False(t, less.IsGreater())
	assert.True(t, less.Is(gust.Less()))

	// Test Equal
	equal := gust.Equal()
	assert.False(t, equal.IsLess())
	assert.True(t, equal.IsEqual())
	assert.False(t, equal.IsGreater())
	assert.True(t, equal.Is(gust.Equal()))

	// Test Greater
	greater := gust.Greater()
	assert.False(t, greater.IsLess())
	assert.False(t, greater.IsEqual())
	assert.True(t, greater.IsGreater())
	assert.True(t, greater.Is(gust.Greater()))

	// Test Compare with integers
	assert.True(t, gust.Compare(1, 2).IsLess())
	assert.True(t, gust.Compare(2, 2).IsEqual())
	assert.True(t, gust.Compare(3, 2).IsGreater())

	// Test Compare with strings
	assert.True(t, gust.Compare("a", "b").IsLess())
	assert.True(t, gust.Compare("b", "b").IsEqual())
	assert.True(t, gust.Compare("c", "b").IsGreater())

	// Test Compare with floats
	assert.True(t, gust.Compare(1.5, 2.5).IsLess())
	assert.True(t, gust.Compare(2.5, 2.5).IsEqual())
	assert.True(t, gust.Compare(3.5, 2.5).IsGreater())

	// Test Compare with uint
	assert.True(t, gust.Compare(uint(10), uint(20)).IsLess())
	assert.True(t, gust.Compare(uint(20), uint(20)).IsEqual())
	assert.True(t, gust.Compare(uint(30), uint(20)).IsGreater())
}
