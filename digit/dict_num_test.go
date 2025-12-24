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

func TestFormatByDict_EmptyDict(t *testing.T) {
	// Test with empty dict (base == 0)
	emptyDict := []byte{}
	result := FormatByDict(emptyDict, 42)
	assert.Equal(t, "", result)
}

func TestParseByDict_EmptyDict(t *testing.T) {
	// Test with empty dict
	emptyDict := []byte{}
	r := ParseByDict[uint64](emptyDict, "ABC")
	assert.True(t, r.IsErr())
	assert.Contains(t, r.Err().Error(), "dict is empty")
}

func TestParseByDict_InvalidChar(t *testing.T) {
	// Test with character not in dict
	dict := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	r := ParseByDict[uint64](dict, "ABC1") // '1' is not in dict
	assert.True(t, r.IsErr())
	assert.Contains(t, r.Err().Error(), "found a char not included in the dict")
}
