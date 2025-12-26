package digit

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParseUint_UnderscoreSkipInBase37 tests underscore character skipping when base > 36
// This should trigger the continue statement at line 58 in parseUint
func TestParseUint_UnderscoreSkipInBase37(t *testing.T) {
	// Test underscore skipping with base 37
	result := ParseUint("1_2_3", 37, 64)
	assert.False(t, result.IsErr())
	expected := uint64(1*37*37 + 2*37 + 3) // 1446
	assert.Equal(t, expected, result.Unwrap())

	// Test multiple underscores
	result2 := ParseUint("1_2_3_4", 37, 64)
	assert.False(t, result2.IsErr())
	expected2 := uint64(1*37*37*37 + 2*37*37 + 3*37 + 4)
	assert.Equal(t, expected2, result2.Unwrap())
}

// TestParseInt_NegativeOverflowExactBoundary tests exact boundary for negative overflow
// This should trigger the check at lines 156-158 in parseInt: if neg && un > cutoff
func TestParseInt_NegativeOverflowExactBoundary(t *testing.T) {
	// Test min int64 - 1 (should overflow)
	result := ParseInt("-9223372036854775809", 10, 64) // -2^63 - 1
	assert.True(t, result.IsErr())
	if nerr, ok := result.UnwrapErr().(*strconv.NumError); ok {
		assert.Equal(t, strconv.ErrRange, nerr.Err)
	}

	// Test boundary value: -2^63 (should succeed, this is the minimum)
	result2 := ParseInt("-9223372036854775808", 10, 64) // -2^63
	assert.False(t, result2.IsErr())
	assert.Equal(t, int64(-9223372036854775808), result2.Unwrap())
}

// TestParseUint_DefaultCase tests default case in parseUint switch (invalid character)
// This covers line 59-60: default case for invalid character
func TestParseUint_DefaultCase(t *testing.T) {
	// Test with invalid character '@' for base > 36
	result := ParseUint("123@", 37, 64)
	assert.True(t, result.IsErr())
	if nerr, ok := result.UnwrapErr().(*strconv.NumError); ok {
		assert.Equal(t, strconv.ErrSyntax, nerr.Err)
		assert.Equal(t, "ParseUint", nerr.Func)
	}
}

// TestParseInt_DefaultCase tests default case in parseInt switch (invalid character)
func TestParseInt_DefaultCase(t *testing.T) {
	// Test with invalid character '@' for base > 36
	result := ParseInt("123@", 37, 64)
	assert.True(t, result.IsErr())
	if nerr, ok := result.UnwrapErr().(*strconv.NumError); ok {
		assert.Equal(t, strconv.ErrSyntax, nerr.Err)
		assert.Equal(t, "ParseInt", nerr.Func)
	}
}
