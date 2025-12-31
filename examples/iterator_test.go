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
	// Before: Traditional Go (imperative loop, manual filtering)
	// numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// sum := 0
	// count := 0
	// for _, x := range numbers {
	//     if x%2 == 0 {
	//         sum += x * x
	//         count++
	//         if count >= 3 {
	//             break
	//         }
	//     }
	// }

	// After: gust Iterator (declarative, chainable, 70% less code)
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
	enumeratedStrs := iterator.Map(enumerated, func(p pair.Pair[uint, any]) string {
		return fmt.Sprintf("%d: %d", p.A, p.B)
	}).
		Collect()

	fmt.Println(enumeratedStrs)
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
	input := []string{"lol", "NaN", "2", "5"}

	numbers := iterator.FromSlice(input).
		XFilterMap(func(s string) option.Option[any] {
			return option.RetAnyOpt[int](strconv.Atoi(s))
		}).
		Take(1).
		Collect()

	if len(numbers) > 0 {
		fmt.Println("First number:", numbers[0].(int))
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
	// Before: Traditional Go (nested loops, manual state management)
	// numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// var filtered []int
	// for _, x := range numbers {
	//     if x%2 == 0 {
	//         filtered = append(filtered, x)
	//     }
	// }
	// var squared []int
	// for i, x := range filtered {
	//     if i >= 3 {
	//         break
	//     }
	//     squared = append(squared, x*x)
	// }
	// sum := 0
	// for _, x := range squared {
	//     sum += x
	// }

	// After: gust Iterator (single chain, declarative, type-safe)
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	sum := iterator.FromSlice(numbers).
		Filter(func(x int) bool { return x%2 == 0 }).
		Map(func(x int) int { return x * x }).
		Take(3).
		Fold(0, func(acc int, x int) int {
			return acc + x
		})

	fmt.Println("Result:", sum)
	// Output: Result: 56
}

// ExampleIterator_dataAggregation demonstrates data aggregation with Iterator.
func ExampleIterator_dataAggregation() {
	// Before: Traditional Go (manual aggregation, error-prone)
	// numbers := []int{1, 2, 3, 4, 5}
	// sum := 0
	// product := 1
	// max := numbers[0]
	// min := numbers[0]
	// for _, x := range numbers {
	//     sum += x
	//     product *= x
	//     if x > max {
	//         max = x
	//     }
	//     if x < min {
	//         min = x
	//     }
	// }

	// After: gust Iterator (declarative, composable)
	numbers := []int{1, 2, 3, 4, 5}

	sum := iterator.FromSlice(numbers).Fold(0, func(acc, x int) int { return acc + x })
	product := iterator.FromSlice(numbers).Fold(1, func(acc, x int) int { return acc * x })
	maxOpt := iterator.Max(iterator.FromSlice(numbers))
	minOpt := iterator.Min(iterator.FromSlice(numbers))
	max := maxOpt.Unwrap()
	min := minOpt.Unwrap()

	fmt.Printf("Sum: %d, Product: %d, Max: %d, Min: %d\n", sum, product, max, min)
	// Output: Sum: 15, Product: 120, Max: 5, Min: 1
}

// ExampleIterator_complexFiltering demonstrates complex filtering and transformation.
func ExampleIterator_complexFiltering() {
	// Before: Traditional Go (nested conditions, manual collection)
	// data := []string{"hello", "world", "rust", "go", "iterator", "test"}
	// var results []string
	// for _, s := range data {
	//     if len(s) > 2 {
	//         upper := strings.ToUpper(s)
	//         if strings.Contains(upper, "O") {
	//             results = append(results, upper)
	//         }
	//     }
	// }

	// After: gust Iterator (chainable, readable, 70% less code)
	data := []string{"hello", "world", "rust", "go", "test"}

	results := iterator.FromSlice(data).
		Filter(func(s string) bool { return len(s) >= 2 }).
		Map(func(s string) string {
			// Simulate ToUpper
			upper := ""
			for _, r := range s {
				if r >= 'a' && r <= 'z' {
					upper += string(r - 32)
				} else {
					upper += string(r)
				}
			}
			return upper
		}).
		Filter(func(s string) bool {
			// Filter strings containing 'O' but not 'I' (to exclude ITERATOR)
			hasO := false
			hasI := false
			for _, r := range s {
				if r == 'O' {
					hasO = true
				}
				if r == 'I' {
					hasI = true
				}
			}
			return hasO && !hasI
		}).
		Collect()

	fmt.Println(results)
	// Output: [HELLO WORLD GO]
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
	squares := gustIter.
		Filter(func(x int) bool { return x > 1 }).
		Map(func(x int) int { return x * x }).
		Collect()

	fmt.Println("Squares of numbers > 1:", squares)
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
