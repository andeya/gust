package digit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToBool(t *testing.T) {
	assert.True(t, ToBool("a"))
	assert.False(t, ToBool(""))
	assert.True(t, ToBool(1.1))
	assert.False(t, ToBool(0.0))
}
