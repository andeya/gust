package digit

import (
	"crypto/rand"
	"strings"
	"time"

	"github.com/andeya/gust/result"
)

const (
	// caseSensitiveBase is the base for case-sensitive encoding (36 characters: 0-9, a-z)
	caseSensitiveBase = 36
	// caseInsensitiveBase is the base for case-insensitive encoding (62 characters: 0-9, a-z, A-Z)
	caseInsensitiveBase = 62
	// timestampEpoch is the epoch time (2026-01-01 00:00:00 UTC) used as the starting point for timestamp encoding
	// Unix timestamp: 1767225600
	timestampEpoch = 1767225600
	// timestampSecondLength is the length of timestamp suffix in seconds
	// MAX: base62=ZZZZZZ (56800235583, 3825-12-06), base36=zzzzzz (2176782335, 2094-12-24)
	timestampSecondLength = 6
	// timestampSecondMaxBase36 is the maximum offset in seconds from timestampEpoch that can be encoded with base36 (6 digits)
	// Value: 36^6 - 1 = 2176782335
	// Max date: 2026-01-01 + 2176782335 seconds = 2094-12-24 05:45:35 UTC
	timestampSecondMaxBase36 = 2176782335
	// timestampSecondMaxBase62 is the maximum offset in seconds from timestampEpoch that can be encoded with base62 (6 digits)
	// Value: 62^6 - 1 = 56800235583
	// Max date: 2026-01-01 + 56800235583 seconds = 3825-12-06 03:13:03 UTC
	timestampSecondMaxBase62 = 56800235583
)

var (
	caseInsensitiveSubstitute = []byte(Digits)
	caseSensitiveSubstitute   = []byte(Digits[:caseSensitiveBase])
)

// Random provides random string generation with optional timestamp encoding
type Random struct {
	// caseSensitive determines whether to use case-sensitive (base36) or case-insensitive (base62) encoding
	caseSensitive bool
	// base is the numeric base used for encoding (36 or 62)
	base int
	// substitute is the character set used for encoding
	substitute []byte
}

// NewRandom creates a new Random instance for generating random strings with optional timestamp encoding
// caseSensitive: determines the character set used for encoding:
//   - true:  uses lowercase-only characters (0-9, a-z), suitable for case-insensitive systems
//   - false: uses mixed case characters (0-9, a-z, A-Z), provides more encoding capacity
//
// Timestamp encoding range (when using RandomStringWithTime or RandomStringWithCurrentTime):
//   - caseSensitive=true:  can encode timestamps from 2026-01-01 to 2094-12-24 (approximately 69 years)
//   - caseSensitive=false: can encode timestamps from 2026-01-01 to 3825-12-06 (approximately 1800 years)
func NewRandom(caseSensitive bool) *Random {
	base := caseInsensitiveBase
	substitute := caseInsensitiveSubstitute
	if caseSensitive {
		base = caseSensitiveBase
		substitute = caseSensitiveSubstitute
	}
	return &Random{
		caseSensitive: caseSensitive,
		base:          base,
		substitute:    substitute,
	}
}

// RandomString generates a random string of the specified length
func (r *Random) RandomString(length int) result.Result[string] {
	if length < 0 {
		return result.TryErr[string]("length must be non-negative")
	}
	if length == 0 {
		return result.Ok("")
	}

	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		charResult := r.randomChar()
		if charResult.IsErr() {
			return result.TryErr[string](charResult.Err())
		}
		buf[i] = charResult.Unwrap()
	}
	return result.Ok(string(buf))
}

// randomChar generates a uniformly distributed random character from the substitute set
// using rejection sampling to ensure uniform distribution
func (r *Random) randomChar() result.Result[byte] {
	base := len(r.substitute)
	// Calculate the maximum value that ensures uniform distribution
	// We reject values >= maxValid to avoid modulo bias
	maxValid := (256 / base) * base

	for {
		bytesResult := RandomBytes(1)
		if bytesResult.IsErr() {
			return result.TryErr[byte](bytesResult.Err())
		}
		b := bytesResult.Unwrap()[0]
		if int(b) < maxValid {
			return result.Ok(r.substitute[int(b)%base])
		}
		// Reject and try again to avoid modulo bias
	}
}

// maxTimestampForBase returns the maximum timestamp that can be encoded with the current base
func (r *Random) maxTimestampForBase() int64 {
	if r.base == caseSensitiveBase {
		return timestampSecondMaxBase36
	}
	return timestampSecondMaxBase62
}

// RandomStringWithTime returns a random string with UNIX timestamp (in seconds) appended as suffix
// length: total length of the returned string, must be greater than timestampSecondLength (6)
// unixTs: absolute UNIX timestamp in seconds, will be converted to offset from timestampEpoch (2026-01-01)
// Supported timestamp range depends on the caseSensitive setting used when creating the Random instance:
//   - caseSensitive=true:  timestamps from 2026-01-01 to 2094-12-24 (approximately 69 years)
//   - caseSensitive=false: timestamps from 2026-01-01 to 3825-12-06 (approximately 1800 years)
func (r *Random) RandomStringWithTime(length int, unixTs int64) result.Result[string] {
	if length <= timestampSecondLength {
		return result.FmtErr[string]("length must be greater than %d", timestampSecondLength)
	}
	// Convert absolute timestamp to offset from epoch
	if unixTs < timestampEpoch {
		return result.FmtErr[string]("unixTs must be >= timestampEpoch (%d)", timestampEpoch)
	}
	offset := unixTs - timestampEpoch
	maxTs := r.maxTimestampForBase()
	if offset < 0 || offset > maxTs {
		return result.FmtErr[string]("timestamp offset is out of range [0,%d] for base %d (absolute timestamp range: [%d,%d])", maxTs, r.base, timestampEpoch, timestampEpoch+maxTs)
	}
	return r.RandomString(length - timestampSecondLength).AndThen(func(randomPart string) result.Result[string] {
		return r.formatTimestamp(offset).Map(func(timestampPart string) string {
			return randomPart + timestampPart
		})
	})
}

// RandomStringWithCurrentTime returns a random string with current UNIX timestamp (in seconds) appended as suffix
// length: total length of the returned string, must be greater than timestampSecondLength (6)
// The current time will be converted to offset from timestampEpoch (2026-01-01)
func (r *Random) RandomStringWithCurrentTime(length int) result.Result[string] {
	now := time.Now().Unix()
	return r.RandomStringWithTime(length, now)
}

// formatTimestamp formats a timestamp offset (from timestampEpoch) to a fixed-length string using the specified base
// The result is always timestampSecondLength characters long, padded with '0' on the left if necessary
func (r *Random) formatTimestamp(offset int64) result.Result[string] {
	formatted := FormatInt(offset, r.base)
	if len(formatted) > timestampSecondLength {
		// Timestamp exceeds maximum encodable length - return error instead of truncating
		return result.FmtErr[string]("timestamp offset %d exceeds maximum encodable length %d for base %d", offset, timestampSecondLength, r.base)
	}
	if len(formatted) < timestampSecondLength {
		// Pad with '0' on the left
		padding := strings.Repeat("0", timestampSecondLength-len(formatted))
		return result.Ok(padding + formatted)
	}
	return result.Ok(formatted)
}

// ParseTime parses absolute UNIX timestamp (in seconds) from the suffix of stringWithTime
// The timestamp offset is expected to be encoded in the last 6 characters of the string
// Returns the absolute timestamp by adding the offset to timestampEpoch
func (r *Random) ParseTime(stringWithTime string) result.Result[int64] {
	length := len(stringWithTime)
	if length <= timestampSecondLength {
		return result.FmtErr[int64]("stringWithTime length must be greater than %d", timestampSecondLength)
	}
	offsetResult := ParseInt(stringWithTime[length-timestampSecondLength:], r.base, 64)
	return offsetResult.Map(func(offset int64) int64 {
		return timestampEpoch + offset
	})
}

// RandomBytes returns securely generated random bytes
// Returns an error if the system's secure random number generator fails to function correctly
func RandomBytes(n int) result.Result[[]byte] {
	if n < 0 {
		return result.TryErr[[]byte]("n must be non-negative")
	}
	if n == 0 {
		return result.Ok([]byte{})
	}

	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return result.TryErr[[]byte](err)
	}
	return result.Ok(b)
}
