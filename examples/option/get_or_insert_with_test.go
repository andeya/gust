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
	assert.Equal(t, 5, x.Unwrap())
	*y = 7
	assert.Equal(t, 7, x.Unwrap())

	var x2 = gust.None[int]()
	var y2 = x2.GetOrInsertWith(nil)
	assert.Equal(t, 0, *y2)
	assert.Equal(t, 0, x2.Unwrap())
	*y2 = 7
	assert.Equal(t, 7, x2.Unwrap())
}
