package examples_test

import (
	"fmt"

	"github.com/andeya/gust/conv"
)

// Example_conv_bytesString demonstrates zero-copy byte/string conversion.
func Example_conv_bytesString() {
	// Convert bytes to string (zero-copy)
	bytes := []byte{'h', 'e', 'l', 'l', 'o'}
	str := conv.BytesToString[string](bytes)
	fmt.Println("String:", str)

	// Convert string to readonly bytes (zero-copy)
	readonlyBytes := conv.StringToReadonlyBytes("world")
	fmt.Println("Readonly bytes length:", len(readonlyBytes))
	// Note: modifying readonlyBytes will panic

	// Output:
	// String: hello
	// Readonly bytes length: 5
}

// Example_conv_caseConversion demonstrates case conversion utilities.
func Example_conv_caseConversion() {
	// Convert to snake_case
	snake := conv.ToSnakeCase("UserID")
	fmt.Println("Snake case:", snake)

	// Convert to CamelCase
	camel := conv.ToCamelCase("user_id")
	fmt.Println("Camel case:", camel)

	// Convert to PascalCase
	pascal := conv.ToPascalCase("user_id")
	fmt.Println("Pascal case:", pascal)

	// Output:
	// Snake case: user_id
	// Camel case: UserId
	// Pascal case: UserID
}

// Example_conv_jsonQuoting demonstrates JSON string quoting.
func Example_conv_jsonQuoting() {
	// Quote a string for JSON (escape HTML)
	quoted := conv.QuoteJSONString("hello \"world\"", true)
	fmt.Println("Quoted (escape HTML):", string(quoted))

	// Quote a string for JSON (don't escape HTML)
	quoted2 := conv.QuoteJSONString("hello <world>", false)
	fmt.Println("Quoted (no HTML escape):", string(quoted2))

	// Output:
	// Quoted (escape HTML): "hello \"world\""
	// Quoted (no HTML escape): "hello <world>"
}

// Example_conv_safeAssert demonstrates safe type assertions.
func Example_conv_safeAssert() {
	// Safe assertion for slice
	anySlice := []any{1, 2, 3, 4, 5}
	intSlice := conv.SafeAssertSlice[int](anySlice)
	if intSlice.IsOk() {
		fmt.Println("Converted slice:", intSlice.Unwrap())
	}

	// Safe assertion for map
	anyMap := map[string]any{"a": 1, "b": 2, "c": 3}
	intMap := conv.SafeAssertMap[string, int](anyMap)
	if intMap.IsOk() {
		fmt.Println("Converted map:", intMap.Unwrap())
	}

	// Output:
	// Converted slice: [1 2 3 4 5]
	// Converted map: map[a:1 b:2 c:3]
}

// Example_conv_refDeref demonstrates reference and dereference operations.
func Example_conv_refDeref() {
	// Create a reference
	value := 42
	ref := conv.Ref(value)
	fmt.Println("Reference value:", *ref)

	// Dereference (handles nil safely)
	deref := conv.Deref(ref)
	fmt.Println("Dereferenced value:", deref)

	// Dereference nil (returns zero value)
	var nilPtr *int
	zeroValue := conv.Deref(nilPtr)
	fmt.Println("Dereferenced nil:", zeroValue)

	// Output:
	// Reference value: 42
	// Dereferenced value: 42
	// Dereferenced nil: 0
}

// Example_conv_stringManipulation demonstrates string manipulation utilities.
func Example_conv_stringManipulation() {
	// Indent text
	text := "line1\nline2\nline3"
	indented := conv.Indent(text, "  ")
	fmt.Println("Indented text:")
	fmt.Println(indented)

	// Normalize whitespace
	spaced := "  hello   world  \n  test  "
	normalized := conv.NormalizeWhitespace(spaced)
	fmt.Println("Normalized:", normalized)

	// Output:
	// Indented text:
	//   line1
	//   line2
	//   line3
	// Normalized: hello world test
}
