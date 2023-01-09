package dict

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	var m = map[string]string{"a": "b", "c": "d"}
	assert.Equal(t, gust.Some("b"), Get(m, "a"))
	assert.Equal(t, gust.None[string](), Get(m, "x"))
	var m2 map[string]string
	assert.Equal(t, gust.None[string](), Get(m2, "x"))
}
