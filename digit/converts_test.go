package digit

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToBool(t *testing.T) {
	assert.True(t, ToBool(1.1))
	assert.False(t, ToBool(0.0))
	assert.True(t, ToBool(1))
	assert.False(t, ToBool(0))
	assert.True(t, ToBool(-1))
}

func TestToBools(t *testing.T) {
	result := ToBools([]int{0, 1, 2, -1})
	assert.Equal(t, []bool{false, true, true, true}, result)

	result2 := ToBools([]float64{0.0, 1.5, -2.3})
	assert.Equal(t, []bool{false, true, true}, result2)
}

func TestFromBool(t *testing.T) {
	assert.Equal(t, int(1), FromBool[bool, int](true))
	assert.Equal(t, int(0), FromBool[bool, int](false))
	assert.Equal(t, float64(1), FromBool[bool, float64](true))
	assert.Equal(t, float64(0), FromBool[bool, float64](false))
}

func TestFromBools(t *testing.T) {
	result := FromBools[bool, int]([]bool{true, false, true})
	assert.Equal(t, []int{1, 0, 1}, result)

	result2 := FromBools[bool, float64]([]bool{false, true})
	assert.Equal(t, []float64{0, 1}, result2)
}

func TestTryFromString(t *testing.T) {
	// Test with int
	result := TryFromString[string, int]("42", 10, 0)
	assert.True(t, result.IsOk())
	assert.Equal(t, 42, result.Unwrap())

	// Test with int8
	result_int8 := TryFromString[string, int8]("42", 10, 0)
	assert.True(t, result_int8.IsOk())
	assert.Equal(t, int8(42), result_int8.Unwrap())

	// Test with int16
	result_int16 := TryFromString[string, int16]("42", 10, 0)
	assert.True(t, result_int16.IsOk())
	assert.Equal(t, int16(42), result_int16.Unwrap())

	// Test with int32
	result_int32 := TryFromString[string, int32]("42", 10, 0)
	assert.True(t, result_int32.IsOk())
	assert.Equal(t, int32(42), result_int32.Unwrap())

	// Test with int64
	result_int64 := TryFromString[string, int64]("42", 10, 0)
	assert.True(t, result_int64.IsOk())
	assert.Equal(t, int64(42), result_int64.Unwrap())

	// Test with uint
	result2 := TryFromString[string, uint]("42", 10, 0)
	assert.True(t, result2.IsOk())
	assert.Equal(t, uint(42), result2.Unwrap())

	// Test with uint8
	result_uint8 := TryFromString[string, uint8]("42", 10, 0)
	assert.True(t, result_uint8.IsOk())
	assert.Equal(t, uint8(42), result_uint8.Unwrap())

	// Test with uint16
	result_uint16 := TryFromString[string, uint16]("42", 10, 0)
	assert.True(t, result_uint16.IsOk())
	assert.Equal(t, uint16(42), result_uint16.Unwrap())

	// Test with uint32
	result_uint32 := TryFromString[string, uint32]("42", 10, 0)
	assert.True(t, result_uint32.IsOk())
	assert.Equal(t, uint32(42), result_uint32.Unwrap())

	// Test with uint64
	result_uint64 := TryFromString[string, uint64]("42", 10, 0)
	assert.True(t, result_uint64.IsOk())
	assert.Equal(t, uint64(42), result_uint64.Unwrap())

	// Test with float32
	result_float32 := TryFromString[string, float32]("3.14", 10, 32)
	assert.True(t, result_float32.IsOk())
	assert.InDelta(t, float32(3.14), result_float32.Unwrap(), 0.001)

	// Test with float64
	result3 := TryFromString[string, float64]("3.14", 10, 64)
	assert.True(t, result3.IsOk())
	assert.InDelta(t, 3.14, result3.Unwrap(), 0.001)

	// Test with invalid string
	result4 := TryFromString[string, int]("invalid", 10, 0)
	assert.True(t, result4.IsErr())

	// Test with base 16
	result5 := TryFromString[string, int]("FF", 16, 0)
	assert.True(t, result5.IsOk())
	assert.Equal(t, 255, result5.Unwrap())

	// Test with base 2
	result6 := TryFromString[string, int]("1010", 2, 0)
	assert.True(t, result6.IsOk())
	assert.Equal(t, 10, result6.Unwrap())

	// Test with base 8
	result7 := TryFromString[string, int]("10", 8, 0)
	assert.True(t, result7.IsOk())
	assert.Equal(t, 8, result7.Unwrap())
}

func TestTryFromStrings(t *testing.T) {
	// Test with valid strings
	result := TryFromStrings[string, int]([]string{"1", "2", "3"}, 10, 0)
	assert.True(t, result.IsOk())
	assert.Equal(t, []int{1, 2, 3}, result.Unwrap())

	// Test with invalid string
	result2 := TryFromStrings[string, int]([]string{"1", "invalid", "3"}, 10, 0)
	assert.True(t, result2.IsErr())

	// Test with empty slice
	result3 := TryFromStrings[string, int]([]string{}, 10, 0)
	assert.True(t, result3.IsOk())
	assert.Equal(t, []int{}, result3.Unwrap())
}

func TestAs(t *testing.T) {
	// Test int to int
	result1 := As[int, int](42)
	assert.True(t, result1.IsOk())
	assert.Equal(t, 42, result1.Unwrap())

	// Test int to int8
	result_int8 := As[int, int8](42)
	assert.True(t, result_int8.IsOk())
	assert.Equal(t, int8(42), result_int8.Unwrap())

	// Test int to int16
	result_int16 := As[int, int16](42)
	assert.True(t, result_int16.IsOk())
	assert.Equal(t, int16(42), result_int16.Unwrap())

	// Test int to int32
	result_int32 := As[int, int32](42)
	assert.True(t, result_int32.IsOk())
	assert.Equal(t, int32(42), result_int32.Unwrap())

	// Test int to int64
	result2 := As[int, int64](42)
	assert.True(t, result2.IsOk())
	assert.Equal(t, int64(42), result2.Unwrap())

	// Test int to uint
	result3 := As[int, uint](42)
	assert.True(t, result3.IsOk())
	assert.Equal(t, uint(42), result3.Unwrap())

	// Test int to uint8
	result_uint8 := As[int, uint8](42)
	assert.True(t, result_uint8.IsOk())
	assert.Equal(t, uint8(42), result_uint8.Unwrap())

	// Test int to uint16
	result_uint16 := As[int, uint16](42)
	assert.True(t, result_uint16.IsOk())
	assert.Equal(t, uint16(42), result_uint16.Unwrap())

	// Test int to uint32
	result_uint32 := As[int, uint32](42)
	assert.True(t, result_uint32.IsOk())
	assert.Equal(t, uint32(42), result_uint32.Unwrap())

	// Test int to uint64
	result_uint64 := As[int, uint64](42)
	assert.True(t, result_uint64.IsOk())
	assert.Equal(t, uint64(42), result_uint64.Unwrap())

	// Test negative int to uint (should fail)
	result4 := As[int, uint](-1)
	assert.True(t, result4.IsErr())

	// Test negative int to uint8 (should fail)
	result_neg_uint8 := As[int, uint8](-1)
	assert.True(t, result_neg_uint8.IsErr())

	// Test int to float32
	result_float32 := As[int, float32](42)
	assert.True(t, result_float32.IsOk())
	assert.Equal(t, float32(42), result_float32.Unwrap())

	// Test int to float64
	result5 := As[int, float64](42)
	assert.True(t, result5.IsOk())
	assert.Equal(t, float64(42), result5.Unwrap())

	// Test float64 to int
	result6 := As[float64, int](42.5)
	assert.True(t, result6.IsOk())
	assert.Equal(t, 42, result6.Unwrap())

	// Test float64 to int8
	result_float_int8 := As[float64, int8](42.5)
	assert.True(t, result_float_int8.IsOk())
	assert.Equal(t, int8(42), result_float_int8.Unwrap())

	// Test overflow - float64 to int8
	result7 := As[float64, int8](1000.0)
	assert.True(t, result7.IsErr())

	// Test overflow - float64 to int16
	result_overflow_int16 := As[float64, int16](100000.0)
	assert.True(t, result_overflow_int16.IsErr())

	// Test overflow - float64 to float32 (too large)
	result_overflow_float32 := As[float64, float32](math.MaxFloat64)
	assert.True(t, result_overflow_float32.IsErr())

	// Test uint to int
	result8 := As[uint, int](42)
	assert.True(t, result8.IsOk())
	assert.Equal(t, 42, result8.Unwrap())

	// Test uint8 to int
	result9 := As[uint8, int](255)
	assert.True(t, result9.IsOk())
	assert.Equal(t, 255, result9.Unwrap())

	// Test uint16 to int
	result_uint16_int := As[uint16, int](1000)
	assert.True(t, result_uint16_int.IsOk())
	assert.Equal(t, 1000, result_uint16_int.Unwrap())

	// Test uint32 to int
	result_uint32_int := As[uint32, int](1000)
	assert.True(t, result_uint32_int.IsOk())
	assert.Equal(t, 1000, result_uint32_int.Unwrap())

	// Test uint64 to int
	result_uint64_int := As[uint64, int](1000)
	assert.True(t, result_uint64_int.IsOk())
	assert.Equal(t, 1000, result_uint64_int.Unwrap())

	// Test float32 to float64
	result_float32_float64 := As[float32, float64](3.14)
	assert.True(t, result_float32_float64.IsOk())
	assert.InDelta(t, 3.14, result_float32_float64.Unwrap(), 0.001)
}

func TestSliceAs(t *testing.T) {
	// Test int slice to int64 slice
	result, err := SliceAs[int, int64]([]int{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, []int64{1, 2, 3}, result)

	// Test int slice to float64 slice
	result2, err2 := SliceAs[int, float64]([]int{1, 2, 3})
	assert.NoError(t, err2)
	assert.Equal(t, []float64{1.0, 2.0, 3.0}, result2)

	// Test with overflow (should fail)
	result3, err3 := SliceAs[float64, int8]([]float64{1000.0})
	assert.Error(t, err3)
	assert.Nil(t, result3)

	// Test empty slice
	result4, err4 := SliceAs[int, int64]([]int{})
	assert.NoError(t, err4)
	assert.Equal(t, []int64{}, result4)
}
