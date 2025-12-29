// Package random provides secure random string generation with optional timestamp encoding.
//
// The package offers two main features:
//  1. Random string generation using cryptographically secure random number generator
//  2. Timestamped random strings that embed UNIX timestamps for tracking or expiration
//
// # Basic Usage
//
// Generate a simple random string:
//
//	gen := random.NewGenerator(false) // case-insensitive (base62: 0-9, a-z, A-Z)
//	str := gen.RandomString(16).Unwrap() // e.g., "a3B9kL2mN8pQ4rS"
//
// Generate a random string with current timestamp embedded:
//
//	gen := random.NewGenerator(false)
//	str := gen.StringWithNow(20).Unwrap() // e.g., "x7K9mN2pQ4rS5tUvW8yZa3Bc"
//	// The last 6 characters encode the current timestamp
//
// Parse timestamp from a string:
//
//	gen := random.NewGenerator(false)
//	timestamp := gen.ParseTimestamp(str).Unwrap() // Returns absolute UNIX timestamp
//
// # Character Sets
//
// The Generator supports two character sets based on caseSensitive parameter:
//   - caseSensitive=false (default): Uses base62 encoding (0-9, a-z, A-Z) - 62 characters
//     Suitable for most use cases, provides longer timestamp range (until 3825)
//   - caseSensitive=true: Uses base36 encoding (0-9, a-z) - 36 characters
//     Suitable for case-insensitive systems, shorter timestamp range (until 2094)
//
// # Timestamp Encoding
//
// When generating strings with timestamps, the last 6 characters encode the timestamp offset
// from the epoch (2026-01-01 00:00:00 UTC). The timestamp is encoded in the configured base:
//   - Base62: Can encode timestamps from 2026-01-01 to 3825-12-06
//   - Base36: Can encode timestamps from 2026-01-01 to 2094-12-24
//
// # Security
//
// The package uses crypto/rand for cryptographically secure random number generation,
// ensuring that generated strings are suitable for security-sensitive applications
// such as session tokens, API keys, or password reset tokens.
//
// # Examples
//
// See the examples package for more detailed usage examples.
package random

import (
	"crypto/rand"
	"strings"
	"time"

	"github.com/andeya/gust/digit"
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
	caseInsensitiveSubstitute = []byte(digit.Digits)
	caseSensitiveSubstitute   = []byte(digit.Digits[:caseSensitiveBase])
)

// Generator provides random string generation with optional timestamp encoding
type Generator struct {
	// caseSensitive determines whether to use case-sensitive (base36) or case-insensitive (base62) encoding
	caseSensitive bool
	// base is the numeric base used for encoding (36 or 62)
	base int
	// substitute is the character set used for encoding
	substitute []byte
}

// NewGenerator creates a new Generator instance for generating random strings with optional timestamp encoding
// caseSensitive: determines the character set used for encoding:
//   - true:  uses lowercase-only characters (0-9, a-z), suitable for case-insensitive systems
//   - false: uses mixed case characters (0-9, a-z, A-Z), provides more encoding capacity
//
// Timestamp encoding range (when using StringWithTimestamp or StringWithNow):
//   - caseSensitive=true:  can encode timestamps from 2026-01-01 to 2094-12-24 (approximately 69 years)
//   - caseSensitive=false: can encode timestamps from 2026-01-01 to 3825-12-06 (approximately 1800 years)
func NewGenerator(caseSensitive bool) *Generator {
	base := caseInsensitiveBase
	substitute := caseInsensitiveSubstitute
	if caseSensitive {
		base = caseSensitiveBase
		substitute = caseSensitiveSubstitute
	}
	return &Generator{
		caseSensitive: caseSensitive,
		base:          base,
		substitute:    substitute,
	}
}

// RandomString generates a random string of the specified length
func (g *Generator) RandomString(length int) result.Result[string] {
	if length < 0 {
		return result.TryErr[string]("length must be non-negative")
	}
	if length == 0 {
		return result.Ok("")
	}

	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		charResult := g.randomChar()
		if charResult.IsErr() {
			return result.TryErr[string](charResult.Err())
		}
		buf[i] = charResult.Unwrap()
	}
	return result.Ok(string(buf))
}

// randomChar generates a uniformly distributed random character from the substitute set
// using rejection sampling to ensure uniform distribution
func (g *Generator) randomChar() result.Result[byte] {
	base := len(g.substitute)
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
			return result.Ok(g.substitute[int(b)%base])
		}
		// Reject and try again to avoid modulo bias
	}
}

// maxTimestampForBase returns the maximum timestamp that can be encoded with the current base
func (g *Generator) maxTimestampForBase() int64 {
	if g.base == caseSensitiveBase {
		return timestampSecondMaxBase36
	}
	return timestampSecondMaxBase62
}

// StringWithTimestamp returns a random string with UNIX timestamp (in seconds) appended as suffix
// length: total length of the returned string, must be greater than timestampSecondLength (6)
// timestamp: absolute UNIX timestamp in seconds, will be converted to offset from timestampEpoch (2026-01-01)
// Supported timestamp range depends on the caseSensitive setting used when creating the Generator instance:
//   - caseSensitive=true:  timestamps from 2026-01-01 to 2094-12-24 (approximately 69 years)
//   - caseSensitive=false: timestamps from 2026-01-01 to 3825-12-06 (approximately 1800 years)
func (g *Generator) StringWithTimestamp(length int, timestamp int64) result.Result[string] {
	if length <= timestampSecondLength {
		return result.FmtErr[string]("length must be greater than %d", timestampSecondLength)
	}
	// Convert absolute timestamp to offset from epoch
	if timestamp < timestampEpoch {
		return result.FmtErr[string]("timestamp must be >= timestampEpoch (%d)", timestampEpoch)
	}
	offset := timestamp - timestampEpoch
	maxTs := g.maxTimestampForBase()
	if offset < 0 || offset > maxTs {
		return result.FmtErr[string]("timestamp offset is out of range [0,%d] for base %d (absolute timestamp range: [%d,%d])", maxTs, g.base, timestampEpoch, timestampEpoch+maxTs)
	}
	return g.RandomString(length - timestampSecondLength).AndThen(func(randomPart string) result.Result[string] {
		return g.formatTimestamp(offset).Map(func(timestampPart string) string {
			return randomPart + timestampPart
		})
	})
}

// StringWithNow returns a random string with current UNIX timestamp (in seconds) appended as suffix
// length: total length of the returned string, must be greater than timestampSecondLength (6)
// The current time will be converted to offset from timestampEpoch (2026-01-01)
func (g *Generator) StringWithNow(length int) result.Result[string] {
	now := time.Now().Unix()
	return g.StringWithTimestamp(length, now)
}

// formatTimestamp formats a timestamp offset (from timestampEpoch) to a fixed-length string using the specified base
// The result is always timestampSecondLength characters long, padded with '0' on the left if necessary
func (g *Generator) formatTimestamp(offset int64) result.Result[string] {
	formatted := digit.FormatInt(offset, g.base)
	if len(formatted) > timestampSecondLength {
		// Timestamp exceeds maximum encodable length - return error instead of truncating
		return result.FmtErr[string]("timestamp offset %d exceeds maximum encodable length %d for base %d", offset, timestampSecondLength, g.base)
	}
	if len(formatted) < timestampSecondLength {
		// Pad with '0' on the left
		padding := strings.Repeat("0", timestampSecondLength-len(formatted))
		return result.Ok(padding + formatted)
	}
	return result.Ok(formatted)
}

// ParseTimestamp parses absolute UNIX timestamp (in seconds) from the suffix of s
// The timestamp offset is expected to be encoded in the last 6 characters of the string
// Returns the absolute timestamp by adding the offset to timestampEpoch
func (g *Generator) ParseTimestamp(s string) result.Result[int64] {
	length := len(s)
	if length <= timestampSecondLength {
		return result.FmtErr[int64]("string length must be greater than %d", timestampSecondLength)
	}
	offsetResult := digit.ParseInt(s[length-timestampSecondLength:], g.base, 64)
	return offsetResult.Map(func(offset int64) int64 {
		return timestampEpoch + offset
	})
}

// RandomBytes returns securely generated random bytes
// count: number of bytes to generate
// Returns an error if the system's secure random number generator fails to function correctly
func RandomBytes(count int) result.Result[[]byte] {
	if count < 0 {
		return result.TryErr[[]byte]("count must be non-negative")
	}
	if count == 0 {
		return result.Ok([]byte{})
	}

	b := make([]byte, count)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return result.TryErr[[]byte](err)
	}
	return result.Ok(b)
}
