package digit

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRandom(t *testing.T) {
	// Test case-insensitive (base62)
	r1 := NewRandom(false)
	assert.NotNil(t, r1)
	assert.False(t, r1.caseSensitive)
	assert.Equal(t, caseInsensitiveBase, r1.base)
	assert.Equal(t, len(Digits), len(r1.substitute))

	// Test case-sensitive (base36)
	r2 := NewRandom(true)
	assert.NotNil(t, r2)
	assert.True(t, r2.caseSensitive)
	assert.Equal(t, caseSensitiveBase, r2.base)
	assert.Equal(t, caseSensitiveBase, len(r2.substitute))
}

func TestRandomString(t *testing.T) {
	// Test case-insensitive random string
	r1 := NewRandom(false)

	// Test different lengths
	for _, length := range []int{0, 1, 5, 10, 20, 100} {
		result := r1.RandomString(length)
		assert.False(t, result.IsErr(), "RandomString(%d) should succeed", length)
		str := result.Unwrap()
		assert.Equal(t, length, len(str))

		// Verify all characters are from the expected set (base62)
		for _, char := range str {
			assert.True(t, strings.ContainsRune(Digits, char), "character %c not in base62 set", char)
		}
	}

	// Test case-sensitive random string
	r2 := NewRandom(true)
	result2 := r2.RandomString(50)
	assert.False(t, result2.IsErr())
	str2 := result2.Unwrap()
	assert.Equal(t, 50, len(str2))

	// Verify all characters are from base36 set (lowercase only)
	base36Set := Digits[:caseSensitiveBase]
	for _, char := range str2 {
		assert.True(t, strings.ContainsRune(base36Set, char), "character %c not in base36 set", char)
	}

	// Test uniqueness (very unlikely to get same string twice)
	result3 := r1.RandomString(20)
	result4 := r1.RandomString(20)
	assert.False(t, result3.IsErr())
	assert.False(t, result4.IsErr())
	str3 := result3.Unwrap()
	str4 := result4.Unwrap()
	// While theoretically possible, it's extremely unlikely
	assert.NotEqual(t, str3, str4, "random strings should be different")

	// Test negative length
	result5 := r1.RandomString(-1)
	assert.True(t, result5.IsErr())
}

func TestRandomStringWithTime(t *testing.T) {
	r1 := NewRandom(false)
	r2 := NewRandom(true)

	// Test valid cases
	testCases := []struct {
		name    string
		random  *Random
		length  int
		unixTs  int64
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid base62 with timestamp",
			random:  r1,
			length:  10,
			unixTs:  1609459200, // 2021-01-01 00:00:00 UTC
			wantErr: false,
		},
		{
			name:    "valid base36 with timestamp",
			random:  r2,
			length:  10,
			unixTs:  1609459200,
			wantErr: false,
		},
		{
			name:    "minimum valid length",
			random:  r1,
			length:  7, // timestampSecondLength + 1
			unixTs:  0,
			wantErr: false,
		},
		{
			name:    "maximum timestamp",
			random:  r1,
			length:  10,
			unixTs:  timestampSecondMax,
			wantErr: false,
		},
		{
			name:    "zero timestamp",
			random:  r1,
			length:  10,
			unixTs:  0,
			wantErr: false,
		},
		{
			name:    "length too short",
			random:  r1,
			length:  6, // equal to timestampSecondLength
			unixTs:  1609459200,
			wantErr: true,
			errMsg:  "length must be greater than 6",
		},
		{
			name:    "length less than timestamp length",
			random:  r1,
			length:  5,
			unixTs:  1609459200,
			wantErr: true,
			errMsg:  "length must be greater than 6",
		},
		{
			name:    "negative timestamp",
			random:  r1,
			length:  10,
			unixTs:  -1,
			wantErr: true,
			errMsg:  "unixTs is out of range",
		},
		{
			name:    "timestamp too large",
			random:  r1,
			length:  10,
			unixTs:  timestampSecondMax + 1,
			wantErr: true,
			errMsg:  "unixTs is out of range",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.random.RandomStringWithTime(tc.length, tc.unixTs)

			if tc.wantErr {
				assert.True(t, result.IsErr(), "expected error but got success")
				if tc.errMsg != "" {
					assert.Contains(t, result.Err().Error(), tc.errMsg)
				}
			} else {
				assert.False(t, result.IsErr(), "unexpected error: %v", result.Err())
				str := result.Unwrap()
				assert.Equal(t, tc.length, len(str))

				// Verify timestamp suffix can be parsed back
				parsedResult := tc.random.ParseTime(str)
				assert.False(t, parsedResult.IsErr(), "failed to parse timestamp: %v", parsedResult.Err())
				assert.Equal(t, tc.unixTs, parsedResult.Unwrap())
			}
		})
	}
}

func TestParseTime(t *testing.T) {
	r1 := NewRandom(false)
	r2 := NewRandom(true)

	// Test valid parsing
	t.Run("valid base62 string", func(t *testing.T) {
		// Generate a string with known timestamp
		unixTs := int64(1609459200) // 2021-01-01 00:00:00 UTC
		result := r1.RandomStringWithTime(10, unixTs)
		assert.False(t, result.IsErr())

		str := result.Unwrap()
		parsedResult := r1.ParseTime(str)
		assert.False(t, parsedResult.IsErr())
		assert.Equal(t, unixTs, parsedResult.Unwrap())
	})

	t.Run("valid base36 string", func(t *testing.T) {
		unixTs := int64(1609459200)
		result := r2.RandomStringWithTime(10, unixTs)
		assert.False(t, result.IsErr())

		str := result.Unwrap()
		parsedResult := r2.ParseTime(str)
		assert.False(t, parsedResult.IsErr())
		assert.Equal(t, unixTs, parsedResult.Unwrap())
	})

	t.Run("zero timestamp", func(t *testing.T) {
		result := r1.RandomStringWithTime(10, 0)
		assert.False(t, result.IsErr())

		str := result.Unwrap()
		parsedResult := r1.ParseTime(str)
		assert.False(t, parsedResult.IsErr())
		assert.Equal(t, int64(0), parsedResult.Unwrap())
	})

	t.Run("maximum timestamp", func(t *testing.T) {
		result := r1.RandomStringWithTime(10, timestampSecondMax)
		assert.False(t, result.IsErr())

		str := result.Unwrap()
		parsedResult := r1.ParseTime(str)
		assert.False(t, parsedResult.IsErr())
		assert.Equal(t, int64(timestampSecondMax), parsedResult.Unwrap())
	})

	// Test error cases
	t.Run("string too short", func(t *testing.T) {
		result := r1.ParseTime("short")
		assert.True(t, result.IsErr())
		assert.Contains(t, result.Err().Error(), "stringWithTime length must be greater than 6")
	})

	t.Run("string equal to timestamp length", func(t *testing.T) {
		result := r1.ParseTime("123456") // exactly 6 characters
		assert.True(t, result.IsErr())
		assert.Contains(t, result.Err().Error(), "stringWithTime length must be greater than 6")
	})

	t.Run("empty string", func(t *testing.T) {
		result := r1.ParseTime("")
		assert.True(t, result.IsErr())
		assert.Contains(t, result.Err().Error(), "stringWithTime length must be greater than 6")
	})

	// Test round-trip with various timestamps
	t.Run("round-trip various timestamps", func(t *testing.T) {
		timestamps := []int64{
			0,
			1,
			100,
			1000,
			1000000,
			1609459200, // 2021-01-01
			2147483647, // Max int32
			timestampSecondMax,
		}

		for _, ts := range timestamps {
			result := r1.RandomStringWithTime(20, ts)
			assert.False(t, result.IsErr(), "failed to generate string for timestamp %d", ts)

			str := result.Unwrap()
			parsedResult := r1.ParseTime(str)
			assert.False(t, parsedResult.IsErr(), "failed to parse timestamp from string %s", str)
			assert.Equal(t, ts, parsedResult.Unwrap(), "timestamp mismatch for %d", ts)
		}
	})

	// Test with current timestamp
	t.Run("current timestamp", func(t *testing.T) {
		now := time.Now().Unix()
		if now <= timestampSecondMax {
			result := r1.RandomStringWithTime(20, now)
			assert.False(t, result.IsErr())

			str := result.Unwrap()
			parsedResult := r1.ParseTime(str)
			assert.False(t, parsedResult.IsErr())
			assert.Equal(t, now, parsedResult.Unwrap())
		}
	})
}

func TestRandomStringWithCurrentTime(t *testing.T) {
	r1 := NewRandom(false)

	// Test valid length
	result := r1.RandomStringWithCurrentTime(10)
	assert.False(t, result.IsErr())
	str := result.Unwrap()
	assert.Equal(t, 10, len(str))

	// Verify timestamp can be parsed
	parsedResult := r1.ParseTime(str)
	assert.False(t, parsedResult.IsErr())
	parsedTs := parsedResult.Unwrap()

	// Verify timestamp is recent (within last minute)
	now := time.Now().Unix()
	assert.True(t, parsedTs >= now-60 && parsedTs <= now+1, "timestamp should be recent")

	// Test invalid length
	result2 := r1.RandomStringWithCurrentTime(5)
	assert.True(t, result2.IsErr())
}

func TestRandomBytes(t *testing.T) {
	// Test different lengths
	lengths := []int{0, 1, 5, 10, 20, 100, 1000}

	for _, length := range lengths {
		t.Run(string(rune(length)), func(t *testing.T) {
			result := RandomBytes(length)
			assert.False(t, result.IsErr(), "RandomBytes(%d) should succeed", length)
			bytes := result.Unwrap()
			assert.Equal(t, length, len(bytes))
		})
	}

	// Test uniqueness
	result1 := RandomBytes(100)
	result2 := RandomBytes(100)
	assert.False(t, result1.IsErr())
	assert.False(t, result2.IsErr())
	bytes1 := result1.Unwrap()
	bytes2 := result2.Unwrap()
	assert.NotEqual(t, bytes1, bytes2, "random bytes should be different")

	// Test that bytes are actually random (check distribution)
	// Generate many bytes and check that we get variety
	allBytes := make(map[byte]bool)
	for i := 0; i < 1000; i++ {
		result := RandomBytes(1)
		assert.False(t, result.IsErr())
		b := result.Unwrap()
		allBytes[b[0]] = true
	}
	// With 1000 random bytes, we should have many different values
	// (at least 200 different byte values is reasonable)
	assert.Greater(t, len(allBytes), 200, "random bytes should have good distribution")

	// Test negative length
	result3 := RandomBytes(-1)
	assert.True(t, result3.IsErr())
}

func TestRandomStringWithTime_EdgeCases(t *testing.T) {
	r1 := NewRandom(false)

	// Test boundary length
	t.Run("boundary length 7", func(t *testing.T) {
		result := r1.RandomStringWithTime(7, 1000)
		assert.False(t, result.IsErr())
		str := result.Unwrap()
		assert.Equal(t, 7, len(str))

		// Should have 1 random char + 6 timestamp chars
		parsedResult := r1.ParseTime(str)
		assert.False(t, parsedResult.IsErr())
		assert.Equal(t, int64(1000), parsedResult.Unwrap())
	})

	// Test very long string
	t.Run("very long string", func(t *testing.T) {
		result := r1.RandomStringWithTime(1000, 1609459200)
		assert.False(t, result.IsErr())
		str := result.Unwrap()
		assert.Equal(t, 1000, len(str))

		parsedResult := r1.ParseTime(str)
		assert.False(t, parsedResult.IsErr())
		assert.Equal(t, int64(1609459200), parsedResult.Unwrap())
	})
}

func TestRandomString_DifferentBases(t *testing.T) {
	r1 := NewRandom(false) // base62
	r2 := NewRandom(true)  // base36

	// Generate strings and verify character sets
	_ = r1.RandomString(100) // base62 can have uppercase letters (A-Z), but we can't guarantee it
	result2 := r2.RandomString(100)
	assert.False(t, result2.IsErr())
	str2 := result2.Unwrap()

	// base36 should NOT have uppercase letters
	for _, char := range str2 {
		assert.False(t, char >= 'A' && char <= 'Z', "base36 should not contain uppercase: %c", char)
	}
}

func TestRandomStringWithTime_FormatIntIntegration(t *testing.T) {
	r1 := NewRandom(false)
	r2 := NewRandom(true)

	// Test that FormatInt is used correctly
	testTimestamps := []int64{0, 1, 36, 62, 100, 1000, 10000}

	for _, ts := range testTimestamps {
		// base62
		result1 := r1.RandomStringWithTime(10, ts)
		assert.False(t, result1.IsErr())
		str1 := result1.Unwrap()

		// Extract timestamp suffix
		tsSuffix := str1[len(str1)-6:]

		// Parse it back
		parsed1 := ParseInt(tsSuffix, caseInsensitiveBase, 64)
		assert.False(t, parsed1.IsErr())
		assert.Equal(t, ts, parsed1.Unwrap())

		// base36
		result2 := r2.RandomStringWithTime(10, ts)
		assert.False(t, result2.IsErr())
		str2 := result2.Unwrap()

		tsSuffix2 := str2[len(str2)-6:]
		parsed2 := ParseInt(tsSuffix2, caseSensitiveBase, 64)
		assert.False(t, parsed2.IsErr())
		assert.Equal(t, ts, parsed2.Unwrap())
	}
}

func TestRandom_SubstituteInitialization(t *testing.T) {
	// Verify that substitute slices are correctly initialized
	assert.Equal(t, len(Digits), len(caseInsensitiveSubstitute))
	assert.Equal(t, caseSensitiveBase, len(caseSensitiveSubstitute))

	// Verify caseInsensitiveSubstitute contains all Digits
	for i, char := range Digits {
		assert.Equal(t, byte(char), caseInsensitiveSubstitute[i], "mismatch at index %d: expected %c (%d), got %c (%d)", i, char, char, caseInsensitiveSubstitute[i], caseInsensitiveSubstitute[i])
	}

	// Verify caseSensitiveSubstitute contains only first 36 characters
	for i, char := range Digits[:caseSensitiveBase] {
		assert.Equal(t, byte(char), caseSensitiveSubstitute[i], "mismatch at index %d: expected %c (%d), got %c (%d)", i, char, char, caseSensitiveSubstitute[i], caseSensitiveSubstitute[i])
	}
}
