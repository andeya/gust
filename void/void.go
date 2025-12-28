// Package void provides a Void type for representing the absence of a value.
package void

// Void is a type that represents the absence of a value.
//
// Void is commonly used with Result to create VoidResult (Result[Void]),
// which represents operations that only return success or failure without a value.
// This is equivalent to Rust's Result<(), E>.
//
// # Examples
//
//	// Use with Result
//	var result result.VoidResult = result.RetVoid(err)
//	if result.IsErr() {
//		fmt.Println(result.Err())
//	}
//
//	// Use with OkVoid
//	var success result.VoidResult = result.OkVoid()
//	if success.IsOk() {
//		fmt.Println("Operation succeeded")
//	}
type Void = *struct{}
