// Package gust provides Rust-inspired error handling, optional values, and iteration utilities for Go.
// This file contains the Void type for representing the absence of a value.
package gust

// Void is a type that represents the absence of a value.
type Void = *struct{}
