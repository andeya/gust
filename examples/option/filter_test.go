package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOption_Filter(t *testing.T) {
	var isEven = func(n int32) bool {
		return n%2 == 0
	}
	assert.Equal(t, gust.None[int32](), gust.None[int32]().Filter(isEven))
	assert.Equal(t, gust.None[int32](), gust.Some[int32](3).Filter(isEven))
	assert.Equal(t, gust.Some[int32](4), gust.Some[int32](4).Filter(isEven))
}
