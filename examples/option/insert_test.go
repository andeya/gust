package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOption_Insert(t *testing.T) {
	var opt = gust.None[int]()
	var val = opt.Insert(1)
	assert.Equal(t, 1, *val)
	assert.Equal(t, 1, opt.Unwrap())
	val = opt.Insert(2)
	assert.Equal(t, 2, *val)
	*val = 3
	assert.Equal(t, 3, opt.Unwrap())
}
