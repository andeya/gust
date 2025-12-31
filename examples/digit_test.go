package examples_test

import (
	"fmt"

	"github.com/andeya/gust/digit"
)

// Example_digit_baseConversion demonstrates base 2-62 conversion.
func Example_digit_baseConversion() {
	// Format integer in different bases
	value := uint64(255)
	fmt.Println("Base 2 (binary):", digit.FormatUint(value, 2))
	fmt.Println("Base 8 (octal):", digit.FormatUint(value, 8))
	fmt.Println("Base 16 (hex):", digit.FormatUint(value, 16))
	fmt.Println("Base 36:", digit.FormatUint(value, 36))
	fmt.Println("Base 62:", digit.FormatUint(value, 62))

	// Parse from different bases
	parsed := digit.ParseUint("ff", 16, 64)
	if parsed.IsOk() {
		fmt.Println("Parsed from hex:", parsed.Unwrap())
	}

	// Output:
	// Base 2 (binary): 11111111
	// Base 8 (octal): 377
	// Base 16 (hex): ff
	// Base 36: 73
	// Base 62: 47
	// Parsed from hex: 255
}

// Example_digit_formatParse demonstrates formatting and parsing with Result types.
func Example_digit_formatParse() {
	// Format with error handling
	formatted := digit.FormatInt(-42, 10)
	fmt.Println("Formatted:", formatted)

	// Parse with error handling
	parsed := digit.ParseInt("-42", 10, 64)
	if parsed.IsOk() {
		fmt.Println("Parsed:", parsed.Unwrap())
	}

	// Parse with automatic base detection (0 means auto-detect)
	autoParsed := digit.ParseInt("0xff", 0, 64)
	if autoParsed.IsOk() {
		fmt.Println("Auto-parsed hex:", autoParsed.Unwrap())
	}

	// Output:
	// Formatted: -42
	// Parsed: -42
	// Auto-parsed hex: 255
}

// Example_digit_checkedOperations demonstrates checked arithmetic operations.
func Example_digit_checkedOperations() {
	// Checked addition (prevents overflow)
	result := digit.CheckedAdd(int32(100), int32(50))
	if result.IsSome() {
		fmt.Println("Checked add result:", result.Unwrap())
	}

	// Checked multiplication (prevents overflow)
	result2 := digit.CheckedMul(int8(10), int8(12))
	if result2.IsSome() {
		fmt.Println("Checked mul result:", result2.Unwrap())
	} else {
		fmt.Println("Checked mul overflowed")
	}

	// Checked multiplication that overflows
	result3 := digit.CheckedMul(int8(100), int8(2))
	if result3.IsSome() {
		fmt.Println("Checked mul result:", result3.Unwrap())
	} else {
		fmt.Println("Checked mul overflowed")
	}

	// Saturating operations (clamp to min/max on overflow)
	saturating := digit.SaturatingAdd(int32(2147483640), int32(10))
	fmt.Println("Saturating add:", saturating)

	// Output:
	// Checked add result: 150
	// Checked mul result: 120
	// Checked mul overflowed
	// Saturating add: 2147483647
}

// Example_digit_typeConversion demonstrates type-safe conversions.
func Example_digit_typeConversion() {
	// Convert between numeric types with error handling
	intValue := int32(100)
	uintResult := digit.As[int32, uint32](intValue)
	if uintResult.IsOk() {
		fmt.Println("Converted to uint32:", uintResult.Unwrap())
	}

	// Convert negative value (will fail for unsigned types)
	negativeResult := digit.As[int32, uint32](-10)
	if negativeResult.IsErr() {
		fmt.Println("Conversion failed (negative to unsigned)")
	}

	// Convert slice of types
	intSlice := []int32{1, 2, 3, 4, 5}
	uintSlice, err := digit.SliceAs[int32, uint32](intSlice)
	if err == nil {
		fmt.Println("Converted slice:", uintSlice)
	}

	// Output:
	// Converted to uint32: 100
	// Conversion failed (negative to unsigned)
	// Converted slice: [1 2 3 4 5]
}

// Example_digit_boolConversion demonstrates boolean conversions.
func Example_digit_boolConversion() {
	// Convert digit to bool
	zero := digit.ToBool(0)
	nonZero := digit.ToBool(42)
	fmt.Println("Zero to bool:", zero)
	fmt.Println("Non-zero to bool:", nonZero)

	// Convert bool to digit
	trueVal := digit.FromBool[bool, int](true)
	falseVal := digit.FromBool[bool, int](false)
	fmt.Println("True to digit:", trueVal)
	fmt.Println("False to digit:", falseVal)

	// Convert slices
	digits := []int{0, 1, 2, 0, 3}
	bools := digit.ToBools(digits)
	fmt.Println("Digits to bools:", bools)

	// Output:
	// Zero to bool: false
	// Non-zero to bool: true
	// True to digit: 1
	// False to digit: 0
	// Digits to bools: [false true true false true]
}

// Example_digit_tryFromString demonstrates parsing strings to digits with Result types.
func Example_digit_tryFromString() {
	// Parse string to int with automatic base detection
	intResult := digit.TryFromString[string, int]("42", 0, 64)
	if intResult.IsOk() {
		fmt.Println("Parsed int:", intResult.Unwrap())
	}

	// Parse string to float
	floatResult := digit.TryFromString[string, float64]("3.14", 0, 64)
	if floatResult.IsOk() {
		fmt.Println("Parsed float:", floatResult.Unwrap())
	}

	// Parse multiple strings
	strings := []string{"10", "20", "30"}
	intsResult := digit.TryFromStrings[string, int](strings, 10, 64)
	if intsResult.IsOk() {
		fmt.Println("Parsed ints:", intsResult.Unwrap())
	}

	// Output:
	// Parsed int: 42
	// Parsed float: 3.14
	// Parsed ints: [10 20 30]
}
