package examples_test

import (
	"fmt"
	"strconv"

	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/option"
	"github.com/andeya/gust/pair"
	"github.com/andeya/gust/vec"
)

// ExampleIterator demonstrates Rust-like iterator chains in Go.
func ExampleIterator() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	sum := iterator.FromSlice(numbers).
		Filter(func(x int) bool { return x%2 == 0 }).
		Map(func(x int) int { return x * x }).
		Take(3).
		Fold(0, func(acc int, x int) int {
			return acc + x
		})

	fmt.Println("Sum of first 3 even squares:", sum)
	// Output: Sum of first 3 even squares: 56
}

// ExampleDoubleEndedIterator demonstrates bidirectional iteration.
func ExampleDoubleEndedIterator() {
	numbers := []int{1, 2, 3, 4, 5}
	deIter := iterator.FromSlice(numbers).MustToDoubleEnded()

	// Iterate from front
	fmt.Print("Front: ")
	for i := 0; i < 2; i++ {
		if val := deIter.Next(); val.IsSome() {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(val.Unwrap())
		}
	}

	// Iterate from back
	fmt.Print("\nBack: ")
	for i := 0; i < 2; i++ {
		if val := deIter.NextBack(); val.IsSome() {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(val.Unwrap())
		}
	}
	fmt.Println()
	// Output:
	// Front: 1 2
	// Back: 5 4
}

// ExampleEnumerate demonstrates advanced iterator operations.
func ExampleEnumerate() {
	data := []string{"hello", "world", "rust", "go", "iterator"}

	// Chain operations elegantly - use function API for Enumerate, then continue chaining
	enumerated := iterator.Enumerate(
		iterator.FromSlice(data).
			Filter(func(s string) bool { return len(s) > 2 }).
			XMap(func(s string) any { return len(s) }),
	)
	// Enumerate returns Iterator[pair.Pair[uint, any]], so we need to use Map with proper type
	result := iterator.Map(enumerated, func(p pair.Pair[uint, any]) string {
		return fmt.Sprintf("%d: %d", p.A, p.B)
	}).
		Collect()

	fmt.Println(result)
	// Output: [0: 5 1: 5 2: 4 3: 8]
}

// ExampleFlatMap demonstrates flattening nested structures.
func ExampleFlatMap() {
	words := []string{"hello", "world"}

	chars := iterator.FromSlice(words).
		XFlatMap(func(s string) iterator.Iterator[any] {
			return iterator.FromSlice([]rune(s)).XMap(func(r rune) any { return r })
		}).
		Collect()

	// Convert []any to []rune using vec.MapAlone
	runeSlice := vec.MapAlone(chars, func(v any) rune {
		return v.(rune)
	})

	fmt.Println("Characters:", string(runeSlice))
	// Output: Characters: helloworld
}

// ExampleFindMap demonstrates finding and mapping in one operation.
func ExampleFindMap() {
	numbers := []string{"lol", "NaN", "2", "5"}

	result := iterator.FromSlice(numbers).
		XFilterMap(func(s string) option.Option[any] {
			return option.RetAnyOpt[int](strconv.Atoi(s))
		}).
		Take(1).
		Collect()

	if len(result) > 0 {
		fmt.Println("First number:", result[0].(int))
	} else {
		fmt.Println("First number: 0")
	}
	// Output: First number: 2
}

// Example_iteratorPartition demonstrates splitting iterators.
func Example_iteratorPartition() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	evens, odds := iterator.FromSlice(numbers).
		Partition(func(x int) bool {
			return x%2 == 0
		})

	fmt.Println("Evens:", evens)
	fmt.Println("Odds:", odds)
	// Output:
	// Evens: [2 4 6 8 10]
	// Odds: [1 3 5 7 9]
}

// ExampleIterator_Chain demonstrates chaining multiple iterator operations.
func ExampleIterator_Chain() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Chain multiple operations: filter, map, take, fold
	result := iterator.FromSlice(numbers).
		Filter(func(x int) bool { return x%2 == 0 }).
		Map(func(x int) int { return x * x }).
		Take(3).
		Fold(0, func(acc int, x int) int {
			return acc + x
		})

	fmt.Println("Result:", result)
	// Output: Result: 56
}

// ExampleIterator_Seq demonstrates converting gust Iterator to Go's standard iterator.Seq.
func ExampleIterator_Seq() {
	numbers := []int{1, 2, 3, 4, 5}
	gustIter := iterator.FromSlice(numbers).Filter(func(x int) bool { return x%2 == 0 })

	// Use gust Iterator in Go's standard for-range loop
	fmt.Print("Even numbers: ")
	for v := range gustIter.Seq() {
		fmt.Print(v, " ")
	}
	fmt.Println()
	// Output: Even numbers: 2 4
}

// ExampleFromSeq demonstrates converting Go's standard iterator.Seq to gust Iterator.
func ExampleFromSeq() {
	// Create a Go standard iterator sequence
	goSeq := func(yield func(int) bool) {
		for i := 0; i < 5; i++ {
			if !yield(i) {
				return
			}
		}
	}

	// Convert to gust Iterator and use gust methods
	gustIter, deferStop := iterator.FromSeq(goSeq)
	defer deferStop()
	result := gustIter.
		Filter(func(x int) bool { return x > 1 }).
		Map(func(x int) int { return x * x }).
		Collect()

	fmt.Println("Squares of numbers > 1:", result)
	// Output: Squares of numbers > 1: [4 9 16]
}

// Example_iteratorConstructors demonstrates creating iterators from various sources.
func Example_iteratorConstructors() {
	// From slice
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	fmt.Println("FromSlice:", iter1.Collect())

	// From individual elements
	iter2 := iterator.FromElements(1, 2, 3)
	fmt.Println("FromElements:", iter2.Collect())

	// From range [start, end)
	iter3 := iterator.FromRange(0, 5) // 0, 1, 2, 3, 4
	fmt.Println("FromRange:", iter3.Collect())

	// From function
	count := 0
	iter4 := iterator.FromFunc(func() option.Option[int] {
		if count < 3 {
			count++
			return option.Some(count)
		}
		return option.None[int]()
	})
	fmt.Println("FromFunc:", iter4.Collect())

	// Empty iterator
	iter5 := iterator.Empty[int]()
	fmt.Println("Empty:", iter5.Collect())

	// Single value
	iter6 := iterator.Once(42)
	fmt.Println("Once:", iter6.Collect())

	// Infinite repeat (take first 3)
	iter7 := iterator.Repeat("hello").Take(3)
	fmt.Println("Repeat:", iter7.Collect())
	// Output:
	// FromSlice: [1 2 3]
	// FromElements: [1 2 3]
	// FromRange: [0 1 2 3 4]
	// FromFunc: [1 2 3]
	// Empty: []
	// Once: [42]
	// Repeat: [hello hello hello]
}

// Example_bitSetIteration demonstrates iterating over bits in bit sets or byte slices.
func Example_bitSetIteration() {
	// Iterate over bits in a byte slice
	bytes := []byte{0b10101010, 0b11001100}

	// Get all set bit offsets
	setBits := iterator.FromBitSetBytesOnes(bytes).
		Filter(func(offset int) bool { return offset > 5 }).
		Collect()
	fmt.Println("Set bits (offset > 5):", setBits)

	// Count set bits
	count := iterator.FromBitSetBytesOnes(bytes).Count()
	fmt.Println("Total set bits:", count)

	// Sum of offsets of set bits
	sum := iterator.FromBitSetBytesOnes(bytes).
		Fold(0, func(acc, offset int) int { return acc + offset })
	fmt.Println("Sum of offsets:", sum)
	// Output:
	// Set bits (offset > 5): [6 8 9 12 13]
	// Total set bits: 8
	// Sum of offsets: 54
}
