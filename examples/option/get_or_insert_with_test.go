package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOption_GetOrInsertWith(t *testing.T) {
	var x = gust.None[int]()
	var y = x.GetOrInsertWith(func() int { return 5 })
	assert.Equal(t, 5, *y)
	*y = 7
	assert.Equal(t, 7, x.Unwrap())
}
