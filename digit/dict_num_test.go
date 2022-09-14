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
		i2, err := ParseByDict[uint64](dict, numStr)
		assert.NoError(t, err)
		assert.Equal(t, i, i2)
	}
}

func TestParseByDict(t *testing.T) {
	dict := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numStr := "DDEZQ"
	num, err := ParseByDict[uint64](dict, numStr)
	assert.NoError(t, err)
	t.Logf("DDEZQ=%d", num) // DDEZQ=1427026
	numStr2 := FormatByDict(dict, num)
	assert.Equal(t, numStr2, numStr)
}
