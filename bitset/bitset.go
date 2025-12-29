// Package bitset provides a thread-safe, efficient bit set implementation with comprehensive bit manipulation operations.
//
// The package offers a BitSet type that stores bits efficiently using byte slices,
// supporting operations like setting/getting bits, bitwise operations (AND, OR, XOR, ANDNOT, NOT),
// counting set bits, range operations, and conversion to various string formats.
//
// # Features
//
//   - Thread-safe: All operations are protected by read-write mutexes
//   - Efficient storage: Uses byte slices for compact bit representation
//   - Negative indexing: Supports negative offsets (e.g., -1 for last bit)
//   - Auto-growing: Automatically expands when setting bits beyond current size
//   - Iterator integration: Implements iterator.BitSetLike interface for use with gust iterators
//
// # Basic Usage
//
//	// Create a new bit set
//	bs := bitset.New()
//	bs.Set(0, true).Unwrap()  // Set first bit
//	bs.Set(5, true).Unwrap()   // Set 6th bit
//
//	// Get bit values
//	if bs.Get(0) {
//		fmt.Println("First bit is set")
//	}
//
//	// Count set bits
//	count := bs.Count(0, -1)  // Count all bits
//
//	// Convert to hex string
//	fmt.Println(bs.String())  // Output: "21" (binary: 00100001)
//
// # Iterator Integration
//
// The BitSet implements iterator.BitSetLike interface, allowing seamless integration
// with gust iterators:
//
//	import "github.com/andeya/gust/iterator"
//
//	bs := bitset.New()
//	bs.Set(0, true).Unwrap()
//	bs.Set(2, true).Unwrap()
//
//	// Iterate over all bits
//	iterator.FromBitSet(bs).ForEach(func(p pair.Pair[int, bool]) {
//		fmt.Printf("Bit %d: %v\n", p.A, p.B)
//	})
//
//	// Get only set bits
//	setBits := iterator.FromBitSetOnes(bs).Collect()
//	fmt.Println(setBits)  // Output: [0 2]
//
// # Bitwise Operations
//
//	bs1 := bitset.NewFromString("c0", bitset.EncodingHex).Unwrap()  // 11000000
//	bs2 := bitset.NewFromString("30", bitset.EncodingHex).Unwrap()  // 00110000
//
//	and := bs1.And(bs2)   // 00000000
//	or := bs1.Or(bs2)     // 11110000
//	xor := bs1.Xor(bs2)   // 11110000
//	not := bs1.Not()      // 00111111
//
// # Examples
//
// See the examples package for more detailed usage examples.
package bitset

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/bits"
	"sync"

	"github.com/andeya/gust/conv"
	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/result"
)

// BitSet represents a thread-safe bit set that stores bits efficiently using byte slices.
// Bits are stored with the most significant bit (MSB) first within each byte.
//
// The BitSet implements iterator.BitSetLike interface, making it compatible with
// gust iterator functions like FromBitSet, FromBitSetOnes, and FromBitSetZeros.
type BitSet struct {
	mu  sync.RWMutex
	set []byte
}

// Compile-time check that BitSet implements iterator.BitSetLike interface.
var _ iterator.BitSetLike = (*BitSet)(nil)

// EncodingFormat represents the encoding format for string representation of BitSet.
type EncodingFormat int

const (
	// EncodingHex represents hexadecimal encoding (base16).
	// Compression ratio: 2.00 (2 characters per byte).
	// Best for: Human readability, debugging, backward compatibility.
	EncodingHex EncodingFormat = iota

	// EncodingBase64 represents standard Base64 encoding.
	// Compression ratio: ~1.33 for 3-byte aligned data (best theoretical ratio).
	// Best for: General purpose encoding, standard compatibility.
	// Note: Uses '+' and '/' characters, may need URL encoding.
	EncodingBase64

	// EncodingBase64URL represents URL-safe Base64 encoding (RFC 4648).
	// Compression ratio: ~1.33 for 3-byte aligned data (best theoretical ratio).
	// Best for: URL embedding, API responses, modern applications (recommended default).
	// Note: Uses '-' and '_' instead of '+' and '/', no padding '='.
	EncodingBase64URL

	// EncodingBase62 represents base62 encoding (0-9, a-z, A-Z).
	// Compression ratio: ~1.38-1.50 for large byte arrays.
	// Best for: Maximum compression without special characters, URL-friendly.
	//
	// Implementation uses math/big.Int.Text(62) instead of digit.FormatUint(62):
	//
	//	Feature              big.Int.Text(62)        digit.FormatUint(62)
	//	Value Range          Unlimited               Limited to uint64
	//	Byte Array           Direct (SetBytes)       Must convert to uint64 (may overflow)
	//	Performance          General algorithm       Fast path for small integers
	//	Interface            Independent             Compatible with strconv
	//	Append               Not supported           Supported (zero-copy)
	//
	// IMPORTANT: Size Limitations
	//   - digit.FormatUint(62) only supports uint64 range (0 to 18,446,744,073,709,551,615)
	//   - For byte arrays larger than 8 bytes, converting to uint64 will cause overflow
	//   - big.Int.Text(62) has no size limit and can handle arbitrarily large byte arrays
	//
	// This implementation chooses big.Int because:
	//   - BitSet byte arrays can be arbitrarily large (not limited to uint64)
	//   - big.Int.SetBytes() directly handles byte slices without conversion
	//   - Avoids potential overflow when converting large byte arrays to uint64
	EncodingBase62
)

// New creates a new BitSet with optional initial bytes.
// If no bytes are provided, creates an empty bit set.
//
// # Examples
//
//	// Empty bit set
//	bs := bitset.New()
//
//	// Bit set with initial bytes
//	bs := bitset.New(0xFF, 0x00, 0xAA)
func New(initialBytes ...byte) *BitSet {
	if len(initialBytes) == 0 {
		return &BitSet{set: make([]byte, 0)}
	}
	set := make([]byte, len(initialBytes))
	copy(set, initialBytes)
	return &BitSet{set: set}
}

// encodeBytes encodes a byte slice using the specified encoding format.
func encodeBytes(data []byte, format EncodingFormat) string {
	if len(data) == 0 {
		return ""
	}

	switch format {
	case EncodingHex:
		return hex.EncodeToString(data)
	case EncodingBase64:
		return base64.StdEncoding.EncodeToString(data)
	case EncodingBase64URL:
		return base64.URLEncoding.EncodeToString(data)
	case EncodingBase62:
		// Convert byte slice to big.Int, then to base62 string.
		// Uses big.Int instead of digit.FormatUint because:
		//   - Supports arbitrarily large byte arrays (not limited to uint64)
		//   - Direct byte slice handling via SetBytes() without conversion
		//   - Avoids overflow risk when converting large arrays to uint64
		// NOTE: digit.FormatUint(62) only supports uint64 range, will overflow for data larger than 8 bytes
		bigInt := new(big.Int).SetBytes(data)
		return bigInt.Text(62)
	default:
		return hex.EncodeToString(data) // fallback to hex
	}
}

// decodeString decodes a string using the specified encoding format.
func decodeString(s string, format EncodingFormat) ([]byte, error) {
	if s == "" {
		return []byte{}, nil
	}

	switch format {
	case EncodingHex:
		return hex.DecodeString(s)
	case EncodingBase64:
		return base64.StdEncoding.DecodeString(s)
	case EncodingBase64URL:
		return base64.URLEncoding.DecodeString(s)
	case EncodingBase62:
		// Parse base62 string to big.Int, then to byte slice.
		// NOTE: Using big.Int can handle arbitrarily large data, not limited by uint64
		bigInt := new(big.Int)
		bigInt, ok := bigInt.SetString(s, 62)
		if !ok {
			return nil, fmt.Errorf("invalid base62 string: %q", s)
		}
		return bigInt.Bytes(), nil
	default:
		return hex.DecodeString(s) // fallback to hex
	}
}

// NewFromString creates a BitSet from an encoded string using the specified format.
// Returns an error if the string is invalid for the given format.
//
// # Examples
//
//	// From hex string
//	bs := bitset.NewFromString("c020", bitset.EncodingHex).Unwrap()
//
//	// From base64url string (recommended for modern applications)
//	bs := bitset.NewFromString("wCA", bitset.EncodingBase64URL).Unwrap()
//
//	// From base62 string
//	bs := bitset.NewFromString("lTn7eSlSalC", bitset.EncodingBase62).Unwrap()
func NewFromString(encodedStr string, format EncodingFormat) result.Result[*BitSet] {
	decoded, err := decodeString(encodedStr, format)
	if err != nil {
		return result.TryErr[*BitSet](fmt.Errorf("invalid %v string: %w", format, err))
	}
	return result.Ok(New(decoded...))
}

// NewFromBase64URL creates a BitSet from a Base64URL encoded string.
// This is a convenience function for NewFromString(str, EncodingBase64URL).
// Base64URL is the default encoding format used by String() method.
//
// # Examples
//
//	bs := bitset.NewFromBase64URL("wCA=").Unwrap()
//	// Creates a bit set with bytes: [0xc0, 0x20]
//
//	// Round-trip: encode and decode
//	original := bitset.New(0xFF, 0x00)
//	encoded := original.String()  // Base64URL encoding
//	decoded := bitset.NewFromBase64URL(encoded).Unwrap()
//	// decoded equals original
func NewFromBase64URL(base64URLStr string) result.Result[*BitSet] {
	return NewFromString(base64URLStr, EncodingBase64URL)
}

// Set sets the bit at the specified offset to the given value.
// Returns the previous value of the bit and any error encountered.
//
// Offset semantics:
//   - Positive offsets: 0 = first bit, 1 = second bit, etc.
//   - Negative offsets: -1 = last bit, -2 = second-to-last bit, etc.
//   - Out of range: If offset is beyond current size, the bit set automatically grows
//
// # Examples
//
//	bs := bitset.New()
//	oldValue, err := bs.Set(0, true).Unwrap()
//	// oldValue = false (bit was unset), err = nil
//
//	// Set last bit using negative index
//	bs.Set(-1, true).Unwrap()
func (b *BitSet) Set(offset int, value bool) result.Result[bool] {
	b.mu.Lock()
	defer b.mu.Unlock()

	size := b.sizeUnsafe()
	normalizedOffset := b.normalizeOffset(offset, size)
	if normalizedOffset < 0 {
		return result.TryErr[bool](fmt.Errorf("bit offset %d is out of range", offset))
	}

	byteIdx := normalizedOffset / 8
	bitIdx := normalizedOffset % 8

	// Auto-grow if necessary
	if byteIdx >= len(b.set) {
		newSet := make([]byte, byteIdx+1)
		copy(newSet, b.set)
		b.set = newSet
	}

	// Get old value before modifying
	oldValue := b.getBitUnsafe(byteIdx, bitIdx)

	// Set or clear the bit
	mask := byte(1 << (7 - bitIdx))
	if value {
		b.set[byteIdx] |= mask
	} else {
		b.set[byteIdx] &^= mask
	}

	return result.Ok(oldValue)
}

// Get returns the value of the bit at the specified offset.
// Returns false if the offset is out of range.
//
// Offset semantics:
//   - Positive offsets: 0 = first bit, 1 = second bit, etc.
//   - Negative offsets: -1 = last bit, -2 = second-to-last bit, etc.
//
// # Examples
//
//	bs := bitset.New()
//	bs.Set(5, true).Unwrap()
//	if bs.Get(5) {
//		fmt.Println("Bit 5 is set")
//	}
//
//	// Get last bit using negative index
//	lastBit := bs.Get(-1)
func (b *BitSet) Get(offset int) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	size := b.sizeUnsafe()
	normalizedOffset := b.normalizeOffset(offset, size)
	if normalizedOffset < 0 || normalizedOffset >= size {
		return false
	}

	byteIdx := normalizedOffset / 8
	bitIdx := normalizedOffset % 8
	return b.getBitUnsafe(byteIdx, bitIdx)
}

// Size returns the total number of bits in the bit set.
// The size is always a multiple of 8 (aligned to byte boundaries).
//
// # Examples
//
//	bs := bitset.New()
//	bs.Set(15, true).Unwrap()
//	fmt.Println(bs.Size())  // Output: 16 (2 bytes * 8 bits)
func (b *BitSet) Size() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.sizeUnsafe()
}

// Count returns the number of bits set to 1 within the specified range [start, end].
// Both start and end are inclusive and support negative indexing.
//
// # Examples
//
//	bs := bitset.New()
//	bs.Set(0, true).Unwrap()
//	bs.Set(2, true).Unwrap()
//	bs.Set(5, true).Unwrap()
//
//	count := bs.Count(0, -1)  // Count all bits: 3
//	count := bs.Count(0, 3)   // Count bits 0-3: 2
func (b *BitSet) Count(start, end int) int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	size := b.sizeUnsafe()
	if size == 0 {
		return 0
	}

	startOffset := b.normalizeOffset(start, size)
	endOffset := b.normalizeOffset(end, size)

	if startOffset < 0 {
		startOffset = 0
	}
	if endOffset >= size {
		endOffset = size - 1
	}
	if startOffset > endOffset {
		return 0
	}

	startByteIdx := startOffset / 8
	startBitIdx := startOffset % 8
	endByteIdx := endOffset / 8
	endBitIdx := endOffset % 8

	var count int

	// Count bits in the first partial byte
	if startByteIdx == endByteIdx {
		// Range is within a single byte
		mask := byte(0xFF >> startBitIdx)
		mask &= byte(0xFF << (7 - endBitIdx))
		count += bits.OnesCount8(b.set[startByteIdx] & mask)
	} else {
		// Count bits in first partial byte
		count += bits.OnesCount8(b.set[startByteIdx] << startBitIdx)

		// Count bits in full bytes
		for i := startByteIdx + 1; i < endByteIdx; i++ {
			count += bits.OnesCount8(b.set[i])
		}

		// Count bits in last partial byte
		mask := byte(0xFF << (7 - endBitIdx))
		count += bits.OnesCount8(b.set[endByteIdx] & mask)
	}

	return count
}

// Range calls the function f sequentially for each bit in the bit set.
// If f returns false, Range stops the iteration.
//
// The function receives the bit offset and its value (true if set, false if unset).
//
// # Examples
//
//	bs := bitset.New()
//	bs.Set(0, true).Unwrap()
//	bs.Set(2, true).Unwrap()
//
//	bs.Range(func(offset int, value bool) bool {
//		if value {
//			fmt.Printf("Bit %d is set\n", offset)
//		}
//		return true  // Continue iteration
//	})
func (b *BitSet) Range(f func(offset int, value bool) bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	size := b.sizeUnsafe()
	for offset := 0; offset < size; offset++ {
		byteIdx := offset / 8
		bitIdx := offset % 8
		value := b.getBitUnsafe(byteIdx, bitIdx)
		if !f(offset, value) {
			return
		}
	}
}

// Not returns a new BitSet that is the bitwise NOT of this bit set.
// The original bit set is not modified.
//
// # Examples
//
//	bs := bitset.NewFromHex("c0").Unwrap()  // 11000000
//	not := bs.Not()                          // 00111111
func (b *BitSet) Not() *BitSet {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := &BitSet{
		set: make([]byte, len(b.set)),
	}
	for i, b := range b.set {
		result.set[i] = ^b
	}
	return result
}

// And returns a new BitSet that is the bitwise AND of this bit set and the provided bit sets.
// If no bit sets are provided, returns a copy of this bit set.
// The original bit set is not modified.
//
// # Examples
//
//	bs1 := bitset.NewFromString("c0", bitset.EncodingHex).Unwrap()  // 11000000
//	bs2 := bitset.NewFromString("30", bitset.EncodingHex).Unwrap()  // 00110000
//	and := bs1.And(bs2)                      // 00000000
func (b *BitSet) And(others ...*BitSet) *BitSet {
	if len(others) == 0 {
		return b.Clone()
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	// Determine result size
	maxLen := len(b.set)
	for _, other := range others {
		other.mu.RLock()
		if len(other.set) > maxLen {
			maxLen = len(other.set)
		}
		other.mu.RUnlock()
	}

	result := &BitSet{
		set: make([]byte, maxLen),
	}

	// Copy this bit set
	copy(result.set, b.set)

	// Apply AND operation
	// For AND, we only need to process up to the minimum length
	// Bits beyond the shorter bit set are implicitly 0, so result bits become 0
	minLen := len(b.set)
	for _, other := range others {
		other.mu.RLock()
		otherLen := len(other.set)
		if otherLen < minLen {
			minLen = otherLen
			// Clear bits beyond the shorter bit set
			for i := minLen; i < len(result.set); i++ {
				result.set[i] = 0
			}
		}
		for i := 0; i < minLen; i++ {
			result.set[i] &= other.set[i]
		}
		other.mu.RUnlock()
	}

	return result
}

// Or returns a new BitSet that is the bitwise OR of this bit set and the provided bit sets.
// If no bit sets are provided, returns a copy of this bit set.
// The original bit set is not modified.
//
// # Examples
//
//	bs1 := bitset.NewFromString("c0", bitset.EncodingHex).Unwrap()  // 11000000
//	bs2 := bitset.NewFromString("30", bitset.EncodingHex).Unwrap()  // 00110000
//	or := bs1.Or(bs2)                        // 11110000
func (b *BitSet) Or(others ...*BitSet) *BitSet {
	return b.bitwiseOp("|", others)
}

// Xor returns a new BitSet that is the bitwise XOR of this bit set and the provided bit sets.
// If no bit sets are provided, returns a copy of this bit set.
// The original bit set is not modified.
//
// # Examples
//
//	bs1 := bitset.NewFromString("c0", bitset.EncodingHex).Unwrap()  // 11000000
//	bs2 := bitset.NewFromString("30", bitset.EncodingHex).Unwrap()  // 00110000
//	xor := bs1.Xor(bs2)                      // 11110000
func (b *BitSet) Xor(others ...*BitSet) *BitSet {
	return b.bitwiseOp("^", others)
}

// AndNot returns a new BitSet that is the bitwise AND NOT (bit clear) of this bit set and the provided bit sets.
// If no bit sets are provided, returns a copy of this bit set.
// The original bit set is not modified.
//
// # Examples
//
//	bs1 := bitset.NewFromString("c0", bitset.EncodingHex).Unwrap()  // 11000000
//	bs2 := bitset.NewFromString("30", bitset.EncodingHex).Unwrap()  // 00110000
//	andNot := bs1.AndNot(bs2)                // 11000000 (bits in bs2 are cleared)
func (b *BitSet) AndNot(others ...*BitSet) *BitSet {
	return b.bitwiseOp("&^", others)
}

// Clear sets all bits in the bit set to 0.
// The size of the bit set remains unchanged.
//
// # Examples
//
//	bs := bitset.New()
//	bs.Set(0, true).Unwrap()
//	bs.Set(5, true).Unwrap()
//	bs.Clear()
//	// All bits are now 0
func (b *BitSet) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i := range b.set {
		b.set[i] = 0
	}
}

// Bytes returns a copy of the underlying byte slice.
// Modifying the returned slice does not affect the bit set.
//
// # Examples
//
//	bs := bitset.New()
//	bs.Set(0, true).Unwrap()
//	bytes := bs.Bytes()
//	fmt.Printf("%x\n", bytes)  // Output: "01"
func (b *BitSet) Bytes() []byte {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := make([]byte, len(b.set))
	copy(result, b.set)
	return result
}

// Binary returns a binary string representation of the bit set.
// The sep parameter specifies the separator between bytes.
//
// # Examples
//
//	bs := bitset.NewFromString("c020", bitset.EncodingHex).Unwrap()
//	fmt.Println(bs.Binary(" "))  // Output: "11000000 00100000"
func (b *BitSet) Binary(sep string) string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if len(b.set) == 0 {
		return ""
	}

	var buf bytes.Buffer
	sepBytes := conv.StringToReadonlyBytes(sep)
	for i, b := range b.set {
		if i > 0 {
			buf.Write(sepBytes)
		}
		fmt.Fprintf(&buf, "%08b", b)
	}
	return conv.BytesToString[string](buf.Bytes())
}

// String returns a string representation of the bit set using Base64URL encoding (default).
// Implements fmt.Stringer interface.
//
// Base64URL is chosen as the default because it provides:
//   - Best compression ratio (~1.33 for aligned data)
//   - URL-safe characters (no need for URL encoding)
//   - Standard library support (no external dependencies)
//
// For other encoding formats, use Encode() method.
//
// # Examples
//
//	bs := bitset.New()
//	bs.Set(0, true).Unwrap()
//	fmt.Println(bs.String())  // Output: "AQ==" (Base64URL)
func (b *BitSet) String() string {
	return b.Encode(EncodingBase64URL)
}

// Encode returns a string representation of the bit set using the specified encoding format.
//
// # Encoding Formats
//
//   - EncodingHex: Hexadecimal (base16), compression ratio 2.00
//   - EncodingBase64: Standard Base64, compression ratio ~1.33
//   - EncodingBase64URL: URL-safe Base64 (recommended), compression ratio ~1.33
//   - EncodingBase62: Base62 encoding, compression ratio ~1.38-1.50
//
// # Examples
//
//	bs := bitset.NewFromString("c020", bitset.EncodingHex).Unwrap()
//	fmt.Println(bs.Encode(bitset.EncodingHex))        // "c020"
//	fmt.Println(bs.Encode(bitset.EncodingBase64URL))  // "wCA"
//	fmt.Println(bs.Encode(bitset.EncodingBase62))     // "gYU"
func (b *BitSet) Encode(format EncodingFormat) string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return encodeBytes(b.set, format)
}

// Sub returns a new BitSet containing bits from the specified range [start, end].
// Both start and end are inclusive and support negative indexing.
// The original bit set is not modified.
//
// The extracted bits are packed into the result bit set, starting from offset 0.
//
// # Examples
//
//	bs := bitset.New()
//	bs.Set(5, true).Unwrap()
//	bs.Set(10, true).Unwrap()
//	bs.Set(15, true).Unwrap()
//
//	sub := bs.Sub(5, 15)  // Extract bits 5-15
func (b *BitSet) Sub(start, end int) *BitSet {
	b.mu.RLock()
	defer b.mu.RUnlock()

	size := b.sizeUnsafe()
	if size == 0 {
		return New()
	}

	startOffset := b.normalizeOffset(start, size)
	endOffset := b.normalizeOffset(end, size)

	if startOffset < 0 {
		startOffset = 0
	}
	if endOffset >= size {
		endOffset = size - 1
	}
	if startOffset > endOffset {
		return New()
	}

	// Create result bit set with enough capacity
	bitCount := endOffset - startOffset + 1
	resultByteCount := (bitCount + 7) / 8
	result := &BitSet{
		set: make([]byte, resultByteCount),
	}

	// Extract bits and pack them into result
	for i := 0; i < bitCount; i++ {
		srcOffset := startOffset + i
		srcByteIdx := srcOffset / 8
		srcBitIdx := srcOffset % 8
		srcValue := b.getBitUnsafe(srcByteIdx, srcBitIdx)

		if srcValue {
			dstByteIdx := i / 8
			dstBitIdx := i % 8
			mask := byte(1 << (7 - dstBitIdx))
			result.set[dstByteIdx] |= mask
		}
	}

	return result
}

// Clone returns a deep copy of the bit set.
func (b *BitSet) Clone() *BitSet {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := &BitSet{
		set: make([]byte, len(b.set)),
	}
	copy(result.set, b.set)
	return result
}

// Private helper methods

func (b *BitSet) sizeUnsafe() int {
	return len(b.set) * 8
}

func (b *BitSet) normalizeOffset(offset int, size int) int {
	if offset < 0 {
		return offset + size
	}
	return offset
}

func (b *BitSet) getBitUnsafe(byteIdx, bitIdx int) bool {
	mask := byte(1 << (7 - bitIdx))
	return (b.set[byteIdx] & mask) != 0
}

func (b *BitSet) bitwiseOp(op string, others []*BitSet) *BitSet {
	if len(others) == 0 {
		return b.Clone()
	}

	// Collect all byte slices first (with locks held)
	b.mu.RLock()
	bBytes := make([]byte, len(b.set))
	copy(bBytes, b.set)
	bLen := len(b.set)
	b.mu.RUnlock()

	// Determine result size and collect other byte slices
	maxLen := bLen
	otherBytesList := make([][]byte, len(others))
	for i, other := range others {
		other.mu.RLock()
		otherBytes := make([]byte, len(other.set))
		copy(otherBytes, other.set)
		otherBytesList[i] = otherBytes
		if len(otherBytes) > maxLen {
			maxLen = len(otherBytes)
		}
		other.mu.RUnlock()
	}

	result := &BitSet{
		set: make([]byte, maxLen),
	}
	copy(result.set, bBytes)

	// Apply operation
	for _, otherBytes := range otherBytesList {
		processLen := len(otherBytes)
		if processLen > len(result.set) {
			processLen = len(result.set)
		}
		for i := 0; i < processLen; i++ {
			switch op {
			case "|":
				result.set[i] |= otherBytes[i]
			case "^":
				result.set[i] ^= otherBytes[i]
			case "&^":
				result.set[i] &^= otherBytes[i]
			}
		}
	}

	return result
}
