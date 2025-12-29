package examples_test

import (
	"fmt"

	"github.com/andeya/gust/bitset"
	"github.com/andeya/gust/iterator"
)

// ExampleBitSet_basic demonstrates basic bit set operations.
func ExampleBitSet_basic() {
	// Create a new bit set
	bs := bitset.New()
	bs.Set(0, true).Unwrap()  // Set first bit
	bs.Set(5, true).Unwrap()  // Set 6th bit
	bs.Set(10, true).Unwrap() // Set 11th bit

	// Check bit values
	fmt.Println("Bit 0:", bs.Get(0))
	fmt.Println("Bit 5:", bs.Get(5))
	fmt.Println("Bit 10:", bs.Get(10))

	// Count set bits
	count := bs.Count(0, -1)
	fmt.Println("Total set bits:", count)

	// Output:
	// Bit 0: true
	// Bit 5: true
	// Bit 10: true
	// Total set bits: 3
}

// ExampleBitSet_fromHex demonstrates creating a bit set from hex string.
func ExampleBitSet_fromHex() {
	// Create from hex string
	bs := bitset.NewFromString("c020", bitset.EncodingHex).Unwrap()
	fmt.Println("Hex:", bs.Encode(bitset.EncodingHex))
	fmt.Println("Base64URL (default):", bs.String())
	fmt.Println("Binary:", bs.Binary(" "))

	// Output:
	// Hex: c020
	// Base64URL (default): wCA=
	// Binary: 11000000 00100000
}

// ExampleBitSet_bitwiseOperations demonstrates bitwise operations.
func ExampleBitSet_bitwiseOperations() {
	bs1 := bitset.NewFromString("c0", bitset.EncodingHex).Unwrap() // 11000000
	bs2 := bitset.NewFromString("30", bitset.EncodingHex).Unwrap() // 00110000

	and := bs1.And(bs2)
	fmt.Println("AND:", and.Binary(" "))

	or := bs1.Or(bs2)
	fmt.Println("OR:", or.Binary(" "))

	xor := bs1.Xor(bs2)
	fmt.Println("XOR:", xor.Binary(" "))

	not := bs1.Not()
	fmt.Println("NOT:", not.Binary(" "))

	// Output:
	// AND: 00000000
	// OR: 11110000
	// XOR: 11110000
	// NOT: 00111111
}

// ExampleBitSet_iterator demonstrates using bit set with gust iterators.
func ExampleBitSet_iterator() {
	bs := bitset.New()
	bs.Set(0, true).Unwrap()
	bs.Set(2, true).Unwrap()
	bs.Set(4, true).Unwrap()

	// Get all set bits using iterator
	setBits := iterator.FromBitSetOnes(bs).Collect()
	fmt.Println("Set bits:", setBits)

	// Count set bits
	count := iterator.FromBitSetOnes(bs).Count()
	fmt.Println("Count:", count)

	// Filter and process
	sum := iterator.FromBitSetOnes(bs).
		Filter(func(offset int) bool { return offset > 1 }).
		Fold(0, func(acc, offset int) int { return acc + offset })
	fmt.Println("Sum of offsets > 1:", sum)

	// Output:
	// Set bits: [0 2 4]
	// Count: 3
	// Sum of offsets > 1: 6
}

// ExampleBitSet_range demonstrates iterating over bits using Range.
func ExampleBitSet_range() {
	bs := bitset.New()
	bs.Set(0, true).Unwrap()
	bs.Set(3, true).Unwrap()
	bs.Set(7, true).Unwrap()

	// Iterate and print set bits
	fmt.Println("Set bits:")
	bs.Range(func(offset int, value bool) bool {
		if value {
			fmt.Printf("  Bit %d is set\n", offset)
		}
		return true // Continue iteration
	})

	// Output:
	// Set bits:
	//   Bit 0 is set
	//   Bit 3 is set
	//   Bit 7 is set
}

// ExampleBitSet_sub demonstrates extracting a subset of bits.
func ExampleBitSet_sub() {
	bs := bitset.New()
	bs.Set(5, true).Unwrap()
	bs.Set(10, true).Unwrap()
	bs.Set(15, true).Unwrap()

	// Extract bits 5-15 (11 bits, aligned to 16 bits = 2 bytes)
	sub := bs.Sub(5, 15)
	fmt.Println("Original size:", bs.Size())
	fmt.Println("Subset size:", sub.Size(), "(aligned to byte boundary)")
	fmt.Println("Subset set bits:", iterator.FromBitSetOnes(sub).Collect())

	// Output:
	// Original size: 16
	// Subset size: 16 (aligned to byte boundary)
	// Subset set bits: [0 5 10]
}

// ExampleBitSet_negativeIndex demonstrates using negative indices.
func ExampleBitSet_negativeIndex() {
	bs := bitset.New()
	bs.Set(0, true).Unwrap()
	bs.Set(7, true).Unwrap()
	bs.Set(8, true).Unwrap() // Set bit 8 for demonstration
	bs.Set(15, true).Unwrap()

	// Use negative indices
	fmt.Println("Last bit (-1):", bs.Get(-1))
	fmt.Println("Second to last (-2):", bs.Get(-2))
	fmt.Println("Bit 8 (-8):", bs.Get(-8))

	// Set using negative index
	bs.Set(-1, false).Unwrap()
	fmt.Println("After clearing last bit:", bs.Get(15))

	// Output:
	// Last bit (-1): true
	// Second to last (-2): false
	// Bit 8 (-8): true
	// After clearing last bit: false
}

// ExampleBitSet_encoding demonstrates encoding and decoding bit sets.
func ExampleBitSet_encoding() {
	bs := bitset.New()
	bs.Set(0, true).Unwrap()
	bs.Set(5, true).Unwrap()

	// Default encoding is Base64URL
	encoded := bs.String()
	fmt.Println("Encoded (Base64URL):", encoded)

	// Decode from Base64URL (round-trip)
	decoded := bitset.NewFromBase64URL(encoded).Unwrap()
	fmt.Println("Round-trip successful:", decoded.Bytes()[0] == bs.Bytes()[0])

	// Encode in different formats
	fmt.Println("Hex:", bs.Encode(bitset.EncodingHex))
	fmt.Println("Base64:", bs.Encode(bitset.EncodingBase64))
	fmt.Println("Base62:", bs.Encode(bitset.EncodingBase62))

	// Output:
	// Encoded (Base64URL): hA==
	// Round-trip successful: true
	// Hex: 84
	// Base64: hA==
	// Base62: 28
}
