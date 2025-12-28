// Package pair provides a generic Pair type for representing pairs of values.
package pair

// Pair is a pair of values.
//
// Pair is commonly used with iterators for operations like Zip and Enumerate,
// as well as with Option for Zip and Unzip operations.
//
// # Examples
//
//	// Create a pair
//	p := pair.Pair[int, string]{A: 42, B: "hello"}
//	a, b := p.Split()
//	fmt.Println(a, b) // Output: 42 hello
//
//	// Use with iterators
//	iter1 := iterator.FromSlice([]int{1, 2, 3})
//	iter2 := iterator.FromSlice([]string{"a", "b", "c"})
//	zipped := iterator.Zip(iter1, iter2)
//	for opt := zipped.Next(); opt.IsSome(); opt = zipped.Next() {
//		p := opt.Unwrap()
//		fmt.Println(p.A, p.B) // Output: 1 a, 2 b, 3 c
//	}
type Pair[A any, B any] struct {
	A A
	B B
}

// Split splits the pair into its two components.
//
// # Examples
//
//	p := pair.Pair[int, string]{A: 42, B: "hello"}
//	a, b := p.Split()
//	fmt.Println(a, b) // Output: 42 hello
//
//go:inline
func (p Pair[A, B]) Split() (A, B) {
	return p.A, p.B
}
