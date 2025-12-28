package constraints_test

import (
	"testing"

	"github.com/andeya/gust/constraints"
	"github.com/stretchr/testify/assert"
)

// testPureIntegerHelper is a helper function to test PureInteger constraint
func testPureIntegerHelper[T constraints.PureInteger](val T) T {
	return val
}

// testIntegerHelper is a helper function to test Integer constraint
func testIntegerHelper[T constraints.Integer](val T) T {
	return val
}

// testDigitHelper is a helper function to test Digit constraint
func testDigitHelper[T constraints.Digit](val T) T {
	return val
}

// TestPureInteger tests PureInteger constraint
func TestPureInteger(t *testing.T) {
	// Test that PureInteger accepts pure integer types through generic functions
	assert.Equal(t, int(1), testPureIntegerHelper(int(1)))
	assert.Equal(t, int8(1), testPureIntegerHelper(int8(1)))
	assert.Equal(t, int16(1), testPureIntegerHelper(int16(1)))
	assert.Equal(t, int32(1), testPureIntegerHelper(int32(1)))
	assert.Equal(t, int64(1), testPureIntegerHelper(int64(1)))
	assert.Equal(t, uint(1), testPureIntegerHelper(uint(1)))
	assert.Equal(t, uint8(1), testPureIntegerHelper(uint8(1)))
	assert.Equal(t, uint16(1), testPureIntegerHelper(uint16(1)))
	assert.Equal(t, uint32(1), testPureIntegerHelper(uint32(1)))
	assert.Equal(t, uint64(1), testPureIntegerHelper(uint64(1)))
}

// TestInteger tests Integer constraint
func TestInteger(t *testing.T) {
	// Test that Integer accepts pure integer types
	assert.Equal(t, int(1), testIntegerHelper(int(1)))
	assert.Equal(t, int8(1), testIntegerHelper(int8(1)))
	assert.Equal(t, int16(1), testIntegerHelper(int16(1)))
	assert.Equal(t, int32(1), testIntegerHelper(int32(1)))
	assert.Equal(t, int64(1), testIntegerHelper(int64(1)))
	assert.Equal(t, uint(1), testIntegerHelper(uint(1)))
	assert.Equal(t, uint8(1), testIntegerHelper(uint8(1)))
	assert.Equal(t, uint16(1), testIntegerHelper(uint16(1)))
	assert.Equal(t, uint32(1), testIntegerHelper(uint32(1)))
	assert.Equal(t, uint64(1), testIntegerHelper(uint64(1)))

	// Test that Integer accepts type aliases
	type MyInt int
	assert.Equal(t, MyInt(42), testIntegerHelper(MyInt(42)))

	type MyUint uint
	assert.Equal(t, MyUint(42), testIntegerHelper(MyUint(42)))
}

// TestDigit tests Digit constraint
func TestDigit(t *testing.T) {
	// Test that Digit accepts integer types
	assert.Equal(t, int(1), testDigitHelper(int(1)))
	assert.Equal(t, int8(1), testDigitHelper(int8(1)))
	assert.Equal(t, int16(1), testDigitHelper(int16(1)))
	assert.Equal(t, int32(1), testDigitHelper(int32(1)))
	assert.Equal(t, int64(1), testDigitHelper(int64(1)))
	assert.Equal(t, uint(1), testDigitHelper(uint(1)))
	assert.Equal(t, uint8(1), testDigitHelper(uint8(1)))
	assert.Equal(t, uint16(1), testDigitHelper(uint16(1)))
	assert.Equal(t, uint32(1), testDigitHelper(uint32(1)))
	assert.Equal(t, uint64(1), testDigitHelper(uint64(1)))

	// Test that Digit accepts floating-point types
	assert.Equal(t, float32(1.0), testDigitHelper(float32(1.0)))
	assert.Equal(t, float64(1.0), testDigitHelper(float64(1.0)))

	// Test that Digit accepts type aliases
	type MyInt int
	assert.Equal(t, MyInt(1), testDigitHelper(MyInt(1)))

	type MyFloat float64
	assert.Equal(t, MyFloat(3.14), testDigitHelper(MyFloat(3.14)))
}

// TestConstraintsCompatibility tests that constraints work with generic functions
func TestConstraintsCompatibility(t *testing.T) {
	// Test PureInteger
	assert.Equal(t, int(42), testPureIntegerHelper(int(42)))
	assert.Equal(t, uint(42), testPureIntegerHelper(uint(42)))

	// Test Integer (including type aliases)
	type MyInt int
	assert.Equal(t, MyInt(42), testIntegerHelper(MyInt(42)))

	// Test Digit (including floats)
	assert.Equal(t, float64(3.14), testDigitHelper(float64(3.14)))
	assert.Equal(t, float32(3.14), testDigitHelper(float32(3.14)))

	type MyFloat float64
	assert.Equal(t, MyFloat(3.14), testDigitHelper(MyFloat(3.14)))
}
