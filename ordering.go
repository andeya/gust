// Package gust provides Rust-inspired error handling, optional values, and iteration utilities for Go.
// This file contains ordering types and comparison utilities.
package gust

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
