// Package constraints provides type constraints for generic programming.
//
// This package defines numeric type constraints used throughout the gust library
// for generic type parameters, including PureInteger, Integer, and Digit constraints.
package constraints

// PureInteger represents pure integer types without type aliases.
type PureInteger interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

// Integer represents integer types including type aliases.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Digit represents numeric types including integers and floating-point numbers.
type Digit interface {
	Integer | ~float32 | ~float64
}
