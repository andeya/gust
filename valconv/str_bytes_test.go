package valconv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesToString(t *testing.T) {
	var b = []byte("abc")
	s := BytesToString[string](b)
	assert.Equal(t, string(b), s)
}

func TestStringToReadonlyBytes(t *testing.T) {
	var s = "abc"
	b := StringToReadonlyBytes(s)
	assert.Equal(t, []byte(s), b)
}
