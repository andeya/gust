// Package gust provides Rust-inspired error handling, optional values, and iteration utilities for Go.
// This file contains numeric type constraints for generic programming.
package gust

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
