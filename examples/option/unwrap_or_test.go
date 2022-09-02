package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOption_UnwrapOr(t *testing.T) {
	assert.Equal(t, "car", gust.Some("car").UnwrapOr("bike"))
	assert.Equal(t, "bike", gust.None[string]().UnwrapOr("bike"))
}
