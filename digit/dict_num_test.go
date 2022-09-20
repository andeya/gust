package digit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatByDict(t *testing.T) {
	dict := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := uint64(0); i < 100; i++ {
		numStr := FormatByDict(dict, i)
		t.Logf("i=%d, s=%s", i, numStr)
		r := ParseByDict[uint64](dict, numStr)
		assert.False(t, r.IsErr())
		assert.Equal(t, i, r.Unwrap())
	}
}

func TestParseByDict(t *testing.T) {
	dict := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numStr := "DDEZQ"
	r := ParseByDict[uint64](dict, numStr)
	assert.False(t, r.IsErr())
	t.Logf("DDEZQ=%d", r.Unwrap()) // DDEZQ=1427026
	numStr2 := FormatByDict(dict, r.Unwrap())
	assert.Equal(t, numStr2, numStr)
}
