package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOption_OrElse(t *testing.T) {
	var nobody = func() gust.Option[string] { return gust.None[string]() }
	var vikings = func() gust.Option[string] { return gust.Some("vikings") }
	assert.Equal(t, gust.Some("barbarians"), gust.Some("barbarians").OrElse(vikings))
	assert.Equal(t, gust.Some("vikings"), gust.None[string]().OrElse(vikings))
	assert.Equal(t, gust.None[string](), gust.None[string]().OrElse(nobody))
}
