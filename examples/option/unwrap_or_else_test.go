package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOption_UnwrapOrElse(t *testing.T) {
	var k = 10
	assert.Equal(t, 4, gust.Some(4).UnwrapOrElse(func() int { return 2 * k }))
	assert.Equal(t, 20, gust.None[int]().UnwrapOrElse(func() int { return 2 * k }))
}
