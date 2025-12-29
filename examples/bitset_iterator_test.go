package examples_test

import (
	"fmt"

	"github.com/andeya/gust/bitset"
	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/pair"
)

// ExampleBitSet_iterator_compatibility demonstrates how BitSet implements iterator.BitSetLike interface.
func ExampleBitSet_iterator_compatibility() {
	// BitSet implements iterator.BitSetLike interface by providing:
	//   - Size() int: returns the number of bits
	//   - Get(offset int) bool: returns the value of a bit at offset

	bs := bitset.New()
	bs.Set(0, true).Unwrap()
	bs.Set(2, true).Unwrap()
	bs.Set(5, true).Unwrap()
	bs.Set(10, true).Unwrap()

	// 1. Iterate over all bits (offset, value pairs)
	iterator.FromBitSet(bs).ForEach(func(p pair.Pair[int, bool]) {
		if p.B {
			fmt.Printf("Bit %d is set\n", p.A)
		}
	})

	// 2. Get only set bits (just offsets)
	setBits := iterator.FromBitSetOnes(bs).Collect()
	fmt.Println("Set bits:", setBits)

	// 3. Get only unset bits (just offsets)
	unsetBits := iterator.FromBitSetZeros(bs).Take(5).Collect()
	fmt.Println("Unset bits (first 5):", unsetBits)

	// 4. Chain with other iterator methods
	// Filter set bits > 3, then sum their offsets
	sum := iterator.FromBitSetOnes(bs).
		Filter(func(offset int) bool { return offset > 3 }).
		Fold(0, func(acc, offset int) int { return acc + offset })
	fmt.Println("Sum of offsets > 3:", sum)

	// Count set bits
	count := iterator.FromBitSetOnes(bs).Count()
	fmt.Println("Total set bits:", count)

	// Output:
	// Bit 0 is set
	// Bit 2 is set
	// Bit 5 is set
	// Bit 10 is set
	// Set bits: [0 2 5 10]
	// Unset bits (first 5): [1 3 4 6 7]
	// Sum of offsets > 3: 15
	// Total set bits: 4
}

// ExampleBitSet_iterator_advanced demonstrates advanced iterator operations with BitSet.
func ExampleBitSet_iterator_advanced() {
	bs := bitset.New()
	bs.Set(5, true).Unwrap()
	bs.Set(10, true).Unwrap()
	bs.Set(15, true).Unwrap()
	bs.Set(20, true).Unwrap()

	// Find the first set bit greater than 10
	firstLargeBit := iterator.FromBitSetOnes(bs).
		Find(func(offset int) bool { return offset > 10 })
	if firstLargeBit.IsSome() {
		fmt.Printf("First set bit > 10: %d\n", firstLargeBit.Unwrap())
	}

	// Check if all set bits are even
	// Note: 5, 10, 15, 20 - 5 and 15 are odd, so result is false
	allEven := iterator.FromBitSetOnes(bs).
		All(func(offset int) bool { return offset%2 == 0 })
	fmt.Printf("All set bits are even: %v\n", allEven)

	// Get the maximum offset of set bits
	maxOffset := iterator.FromBitSetOnes(bs).
		Fold(0, func(acc, offset int) int {
			if offset > acc {
				return offset
			}
			return acc
		})
	fmt.Printf("Maximum set bit offset: %d\n", maxOffset)

	// Output:
	// First set bit > 10: 15
	// All set bits are even: false
	// Maximum set bit offset: 20
}

// ExampleBitSet_iterator_comparison demonstrates comparing BitSets using iterators.
func ExampleBitSet_iterator_comparison() {
	bs1 := bitset.New()
	bs1.Set(0, true).Unwrap()
	bs1.Set(2, true).Unwrap()
	bs1.Set(4, true).Unwrap()

	bs2 := bitset.New()
	bs2.Set(2, true).Unwrap()
	bs2.Set(4, true).Unwrap()
	bs2.Set(6, true).Unwrap()

	// Find common set bits (intersection)
	commonBits := iterator.FromBitSetOnes(bs1).
		Filter(func(offset int) bool {
			return bs2.Get(offset)
		}).
		Collect()
	fmt.Println("Common set bits:", commonBits)

	// Find bits set in bs1 but not in bs2
	onlyInBs1 := iterator.FromBitSetOnes(bs1).
		Filter(func(offset int) bool {
			return !bs2.Get(offset)
		}).
		Collect()
	fmt.Println("Bits only in bs1:", onlyInBs1)

	// Output:
	// Common set bits: [2 4]
	// Bits only in bs1: [0]
}
