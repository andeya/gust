package bitset_test

import (
	"fmt"
	"testing"

	"github.com/andeya/gust/bitset"
	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/pair"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Empty bit set
	bs := bitset.New()
	assert.NotNil(t, bs)
	assert.Equal(t, 0, bs.Size())

	// Bit set with initial bytes
	bs = bitset.New(0xFF, 0x00, 0xAA)
	assert.Equal(t, 24, bs.Size())
	assert.True(t, bs.Get(0))
	assert.False(t, bs.Get(8))
	assert.True(t, bs.Get(16))
}

func TestNewFromString_Hex(t *testing.T) {
	// Valid hex string
	result := bitset.NewFromString("c020", bitset.EncodingHex)
	assert.False(t, result.IsErr())
	bs := result.Unwrap()
	assert.Equal(t, 16, bs.Size())
	assert.True(t, bs.Get(0))
	assert.True(t, bs.Get(1))
	assert.False(t, bs.Get(2))
	assert.True(t, bs.Get(10))

	// Invalid hex string
	result = bitset.NewFromString("invalid", bitset.EncodingHex)
	assert.True(t, result.IsErr())
}

func TestNewFromBase64URL(t *testing.T) {
	// Valid Base64URL string
	result := bitset.NewFromBase64URL("wCA=")
	assert.False(t, result.IsErr())
	bs := result.Unwrap()
	assert.Equal(t, 16, bs.Size())
	assert.True(t, bs.Get(0))
	assert.True(t, bs.Get(1))
	assert.False(t, bs.Get(2))
	assert.True(t, bs.Get(10))

	// Invalid Base64URL string
	result = bitset.NewFromBase64URL("invalid!")
	assert.True(t, result.IsErr())

	// Round-trip test: encode and decode
	original := bitset.New(0xFF, 0x00, 0xAA)
	encoded := original.String() // Base64URL encoding
	decoded := bitset.NewFromBase64URL(encoded)
	assert.False(t, decoded.IsErr())
	assert.Equal(t, original.Bytes(), decoded.Unwrap().Bytes())
}

func TestSet(t *testing.T) {
	bs := bitset.New()

	// Set first bit
	result := bs.Set(0, true)
	assert.False(t, result.IsErr())
	oldValue := result.Unwrap()
	assert.False(t, oldValue) // Was unset
	assert.True(t, bs.Get(0))

	// Set bit with negative index
	result = bs.Set(-1, true)
	assert.False(t, result.IsErr())
	assert.True(t, bs.Get(7)) // Last bit of first byte

	// Set bit beyond current size (auto-grow)
	result = bs.Set(500, true)
	assert.False(t, result.IsErr())
	assert.True(t, bs.Get(500))
	assert.GreaterOrEqual(t, bs.Size(), 504) // At least 63 bytes * 8 bits

	// Invalid negative offset
	result = bs.Set(-1000, true)
	assert.True(t, result.IsErr())
}

func TestGet(t *testing.T) {
	bs := bitset.New(0xFF, 0x00, 0xAA)

	// Get set bit
	assert.True(t, bs.Get(0))
	assert.True(t, bs.Get(7))

	// Get unset bit
	assert.False(t, bs.Get(8))
	assert.False(t, bs.Get(15))

	// Get with negative index
	// Size is 24 bits (3 bytes), so -1 = bit 23, -8 = bit 16
	// 0xAA = 10101010, so bit 16 (MSB of 0xAA) = 1, bit 23 (LSB of 0xAA) = 0
	assert.False(t, bs.Get(-1)) // Last bit (bit 23) = 0
	assert.True(t, bs.Get(-8))  // Bit 16 (first bit of last byte) = 1

	// Out of range
	assert.False(t, bs.Get(1000))
	assert.False(t, bs.Get(-1000))
}

func TestSize(t *testing.T) {
	bs := bitset.New()
	assert.Equal(t, 0, bs.Size())

	bs.Set(7, true).Unwrap()
	assert.Equal(t, 8, bs.Size())

	bs.Set(15, true).Unwrap()
	assert.Equal(t, 16, bs.Size())

	bs.Set(500, true).Unwrap()
	assert.GreaterOrEqual(t, bs.Size(), 504)
}

func TestCount(t *testing.T) {
	bs := bitset.New()
	bs.Set(0, true).Unwrap()
	bs.Set(2, true).Unwrap()
	bs.Set(5, true).Unwrap()
	bs.Set(7, true).Unwrap()

	// Count all bits
	count := bs.Count(0, -1)
	assert.Equal(t, 4, count)

	// Count range
	count = bs.Count(0, 3)
	assert.Equal(t, 2, count) // Bits 0 and 2

	// Count with negative indices (last 8 bits, which are bits 0-7 of first byte)
	// bs has bits set at 0, 2, 5, 7, so counting bits 0-7 should give 4
	count = bs.Count(-8, -1)
	assert.Equal(t, 4, count) // All 8 bits in first byte (bits 0-7)

	// Invalid range
	count = bs.Count(10, 5)
	assert.Equal(t, 0, count)
}

func TestRange(t *testing.T) {
	bs := bitset.New()
	bs.Set(0, true).Unwrap()
	bs.Set(2, true).Unwrap()

	var offsets []int
	bs.Range(func(offset int, value bool) bool {
		if value {
			offsets = append(offsets, offset)
		}
		return true
	})
	assert.Equal(t, []int{0, 2}, offsets)

	// Early termination
	offsets = []int{}
	bs.Range(func(offset int, value bool) bool {
		if offset >= 5 {
			return false // Stop iteration
		}
		if value {
			offsets = append(offsets, offset)
		}
		return true
	})
	assert.Equal(t, []int{0, 2}, offsets)
}

func TestNot(t *testing.T) {
	bs := bitset.NewFromString("c0", bitset.EncodingHex).Unwrap() // 11000000
	not := bs.Not()                                               // 00111111

	assert.False(t, not.Get(0))
	assert.False(t, not.Get(1))
	assert.True(t, not.Get(2))
	assert.True(t, not.Get(7))

	// Original unchanged
	assert.True(t, bs.Get(0))
	assert.True(t, bs.Get(1))
}

func TestAnd(t *testing.T) {
	bs1 := bitset.NewFromString("c0", bitset.EncodingHex).Unwrap() // 11000000
	bs2 := bitset.NewFromString("30", bitset.EncodingHex).Unwrap() // 00110000
	and := bs1.And(bs2)                                            // 00000000

	assert.False(t, and.Get(0))
	assert.False(t, and.Get(1))
	assert.False(t, and.Get(4))
	assert.False(t, and.Get(5))

	// Original unchanged
	// bs1 = 0xc0 = 11000000, so bits 0,1 are set
	// bs2 = 0x30 = 00110000, so bits 2,3 are set
	assert.True(t, bs1.Get(0))
	assert.True(t, bs1.Get(1))
	assert.True(t, bs2.Get(2))
	assert.True(t, bs2.Get(3))

	// Empty arguments
	clone := bs1.And()
	assert.Equal(t, bs1.Size(), clone.Size())
	assert.True(t, clone.Get(0))
}

func TestOr(t *testing.T) {
	bs1 := bitset.NewFromString("c0", bitset.EncodingHex).Unwrap() // 11000000
	bs2 := bitset.NewFromString("30", bitset.EncodingHex).Unwrap() // 00110000
	or := bs1.Or(bs2)                                              // 11110000 (0xF0)

	// 0xF0 = 11110000, so bits 0-3 are set, bits 4-7 are unset
	assert.True(t, or.Get(0))
	assert.True(t, or.Get(1))
	assert.True(t, or.Get(2))
	assert.True(t, or.Get(3))
	assert.False(t, or.Get(4))
	assert.False(t, or.Get(5))
	assert.False(t, or.Get(6))
	assert.False(t, or.Get(7))
}

func TestXor(t *testing.T) {
	bs1 := bitset.NewFromString("c0", bitset.EncodingHex).Unwrap() // 11000000
	bs2 := bitset.NewFromString("30", bitset.EncodingHex).Unwrap() // 00110000
	xor := bs1.Xor(bs2)                                            // 11110000 (F0)

	assert.True(t, xor.Get(0))
	assert.True(t, xor.Get(1))
	assert.True(t, xor.Get(2))
	assert.True(t, xor.Get(3))
	assert.False(t, xor.Get(4))
	assert.False(t, xor.Get(5))
	assert.False(t, xor.Get(6))
	assert.False(t, xor.Get(7))
}

func TestAndNot(t *testing.T) {
	bs1 := bitset.NewFromString("c0", bitset.EncodingHex).Unwrap() // 11000000
	bs2 := bitset.NewFromString("30", bitset.EncodingHex).Unwrap() // 00110000
	andNot := bs1.AndNot(bs2)                                      // 11000000 (bits in bs2 cleared)

	assert.True(t, andNot.Get(0))
	assert.True(t, andNot.Get(1))
	assert.False(t, andNot.Get(4))
	assert.False(t, andNot.Get(5))
}

func TestClear(t *testing.T) {
	bs := bitset.New()
	bs.Set(0, true).Unwrap()
	bs.Set(5, true).Unwrap()
	bs.Set(10, true).Unwrap()

	size := bs.Size()
	bs.Clear()

	assert.Equal(t, size, bs.Size()) // Size unchanged
	assert.False(t, bs.Get(0))
	assert.False(t, bs.Get(5))
	assert.False(t, bs.Get(10))
}

func TestBytes(t *testing.T) {
	bs := bitset.New()
	bs.Set(0, true).Unwrap()
	bs.Set(7, true).Unwrap()

	bytes := bs.Bytes()
	assert.Equal(t, []byte{0x81}, bytes)

	// Modify returned bytes shouldn't affect bit set
	bytes[0] = 0xFF
	assert.False(t, bs.Get(1))
}

func TestBinary(t *testing.T) {
	bs := bitset.NewFromString("c020", bitset.EncodingHex).Unwrap()
	binary := bs.Binary(" ")
	assert.Equal(t, "11000000 00100000", binary)

	// Empty bit set
	bs = bitset.New()
	assert.Equal(t, "", bs.Binary(" "))
}

func TestString(t *testing.T) {
	bs := bitset.New()
	bs.Set(0, true).Unwrap()
	bs.Set(7, true).Unwrap()

	// String() now uses Base64URL encoding by default
	assert.Equal(t, "gQ==", bs.String())
}

func TestEncode(t *testing.T) {
	bs := bitset.New()
	bs.Set(0, true).Unwrap()
	bs.Set(7, true).Unwrap()

	// Test all encoding formats
	// Bit 0 and 7 set = 0x81 = 10000001
	assert.Equal(t, "81", bs.Encode(bitset.EncodingHex))
	assert.Equal(t, "gQ==", bs.Encode(bitset.EncodingBase64))
	assert.Equal(t, "gQ==", bs.Encode(bitset.EncodingBase64URL))
	// Base62: 0x81 = 129 = 2*62 + 5 = "25"
	assert.Equal(t, "25", bs.Encode(bitset.EncodingBase62))

	// Test with larger bit set
	bs2 := bitset.NewFromString("c020", bitset.EncodingHex).Unwrap()
	assert.Equal(t, "c020", bs2.Encode(bitset.EncodingHex))
	assert.Equal(t, "wCA=", bs2.Encode(bitset.EncodingBase64))
	assert.Equal(t, "wCA=", bs2.Encode(bitset.EncodingBase64URL))
	// Base62: [0xc0, 0x20] = 0xc020 = 49184 in base62 = "cNi"
	assert.Equal(t, "cNi", bs2.Encode(bitset.EncodingBase62))
}

func TestNewFromString(t *testing.T) {
	// Test hex format (backward compatibility)
	result := bitset.NewFromString("c020", bitset.EncodingHex)
	assert.False(t, result.IsErr())
	bs := result.Unwrap()
	assert.Equal(t, 16, bs.Size())
	assert.True(t, bs.Get(0))
	assert.True(t, bs.Get(1))

	// Test Base64URL format
	result = bitset.NewFromString("wCA=", bitset.EncodingBase64URL)
	assert.False(t, result.IsErr())
	bs2 := result.Unwrap()
	assert.Equal(t, 16, bs2.Size())
	assert.True(t, bs2.Get(0))
	assert.True(t, bs2.Get(1))

	// Test Base62 format
	result = bitset.NewFromString("cNi", bitset.EncodingBase62)
	assert.False(t, result.IsErr())
	bs3 := result.Unwrap()
	assert.Equal(t, 16, bs3.Size())
	assert.True(t, bs3.Get(0))
	assert.True(t, bs3.Get(1))

	// Verify all formats produce the same bit set
	assert.Equal(t, bs.Bytes(), bs2.Bytes())
	assert.Equal(t, bs.Bytes(), bs3.Bytes())

	// Test invalid format
	result = bitset.NewFromString("invalid", bitset.EncodingBase64URL)
	assert.True(t, result.IsErr())
}

func TestSub(t *testing.T) {
	bs := bitset.New()
	bs.Set(5, true).Unwrap()
	bs.Set(10, true).Unwrap()
	bs.Set(15, true).Unwrap()

	sub := bs.Sub(5, 15)
	assert.True(t, sub.Get(0)) // First bit of sub (was bit 5)
	assert.False(t, sub.Get(1))
	assert.True(t, sub.Get(5))  // Bit 5 of sub (was bit 10)
	assert.True(t, sub.Get(10)) // Bit 10 of sub (was bit 15)

	// Invalid range
	sub = bs.Sub(20, 10)
	assert.Equal(t, 0, sub.Size())
}

func TestClone(t *testing.T) {
	bs1 := bitset.New()
	bs1.Set(0, true).Unwrap()
	bs1.Set(5, true).Unwrap()

	bs2 := bs1.Clone()
	assert.Equal(t, bs1.Size(), bs2.Size())
	assert.True(t, bs2.Get(0))
	assert.True(t, bs2.Get(5))

	// Modify clone shouldn't affect original
	bs2.Set(0, false).Unwrap()
	assert.True(t, bs1.Get(0))
	assert.False(t, bs2.Get(0))
}

func TestIteratorIntegration(t *testing.T) {
	bs := bitset.New()
	bs.Set(0, true).Unwrap()
	bs.Set(2, true).Unwrap()
	bs.Set(4, true).Unwrap()

	// Test FromBitSet
	var pairs []pair.Pair[int, bool]
	iterator.FromBitSet(bs).ForEach(func(p pair.Pair[int, bool]) {
		if p.B {
			pairs = append(pairs, p)
		}
	})
	assert.Len(t, pairs, 3)
	assert.Equal(t, 0, pairs[0].A)
	assert.Equal(t, 2, pairs[1].A)
	assert.Equal(t, 4, pairs[2].A)

	// Test FromBitSetOnes
	setBits := iterator.FromBitSetOnes(bs).Collect()
	assert.Equal(t, []int{0, 2, 4}, setBits)

	// Test FromBitSetZeros
	unsetBits := iterator.FromBitSetZeros(bs).Take(3).Collect()
	assert.Contains(t, unsetBits, 1)
	assert.Contains(t, unsetBits, 3)
	assert.Contains(t, unsetBits, 5)
}

func ExampleBitSet() {
	bs := bitset.NewFromString("c020", bitset.EncodingHex).Unwrap()
	fmt.Println("Origin:", bs.Binary(" "))
	not := bs.Not()
	fmt.Println("Not:", not.Binary(" "))
	fmt.Println("AndNot:", not.AndNot(bitset.New(1, 1)).Binary(" "))
	fmt.Println("And:", not.And(bitset.New(1<<1, 1<<1)).Binary(" "))
	fmt.Println("Or:", not.Or(bitset.New(1<<7, 1<<7)).Binary(" "))
	fmt.Println("Xor:", not.Xor(bitset.New(1<<7, 1<<7)).Binary(" "))

	not.Range(func(k int, v bool) bool {
		fmt.Println(v)
		return true
	})

	// Output:
	// Origin: 11000000 00100000
	// Not: 00111111 11011111
	// AndNot: 00111110 11011110
	// And: 00000010 00000010
	// Or: 10111111 11011111
	// Xor: 10111111 01011111
	// false
	// false
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// true
	// false
	// true
	// true
	// true
	// true
	// true
}
