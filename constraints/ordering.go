// Package constraints provides type constraints for generic programming.
//
// This package defines type constraints used throughout the gust library for
// generic type parameters, including ordering, numeric, and comparison constraints.
//
// # Examples
//
//	// Use Ordering for comparisons
//	ord := constraints.Compare(1, 2)
//	if ord.IsLess() {
//		fmt.Println("1 is less than 2")
//	}
//
//	// Use constraints in generic functions
//	func max[T constraints.Ord](a, b T) T {
//		if constraints.Compare(a, b).IsGreater() {
//			return a
//		}
//		return b
//	}
//
//	// Use numeric constraints
//	func sum[T constraints.Digit](values []T) T {
//		var total T
//		for _, v := range values {
//			total += v
//		}
//		return total
//	}
package constraints

// Ord represents types that can be ordered (compared).
type Ord interface {
	Digit | ~string | ~uintptr
}

// Ordering represents the result of a comparison between two values.
type Ordering struct {
	cmp int8
}

// Less returns an Ordering representing "less than".
func Less() Ordering {
	return Ordering{cmp: -1}
}

// Equal returns an Ordering representing "equal".
func Equal() Ordering {
	return Ordering{cmp: 0}
}

// Greater returns an Ordering representing "greater than".
func Greater() Ordering {
	return Ordering{cmp: 1}
}

// Compare compares two values and returns their Ordering.
func Compare[T Ord](a, b T) Ordering {
	if a < b {
		return Less()
	}
	if a == b {
		return Equal()
	}
	return Greater()
}

// Is checks if this Ordering matches the given Ordering.
func (o Ordering) Is(ord Ordering) bool {
	return o.cmp == ord.cmp
}

// IsLess returns true if this Ordering represents "less than".
func (o Ordering) IsLess() bool {
	return o.cmp == -1
}

// IsEqual returns true if this Ordering represents "equal".
func (o Ordering) IsEqual() bool {
	return o.cmp == 0
}

// IsGreater returns true if this Ordering represents "greater than".
func (o Ordering) IsGreater() bool {
	return o.cmp == 1
}
