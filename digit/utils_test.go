package digit

import (
	"math"
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestAbs(t *testing.T) {
	// Test with positive integer
	assert.Equal(t, 5, Abs(5))
	assert.Equal(t, int8(5), Abs(int8(5)))
	assert.Equal(t, int16(5), Abs(int16(5)))
	assert.Equal(t, int32(5), Abs(int32(5)))
	assert.Equal(t, int64(5), Abs(int64(5)))

	// Test with negative integer
	assert.Equal(t, 5, Abs(-5))
	assert.Equal(t, int8(5), Abs(int8(-5)))
	assert.Equal(t, int16(5), Abs(int16(-5)))
	assert.Equal(t, int32(5), Abs(int32(-5)))
	assert.Equal(t, int64(5), Abs(int64(-5)))

	// Test with zero
	assert.Equal(t, 0, Abs(0))
	assert.Equal(t, int8(0), Abs(int8(0)))

	// Test with unsigned (should return same)
	assert.Equal(t, uint(5), Abs(uint(5)))
	assert.Equal(t, uint8(5), Abs(uint8(5)))

	// Test with float
	assert.Equal(t, 5.5, Abs(5.5))
	assert.Equal(t, 5.5, Abs(-5.5))
	assert.Equal(t, float32(5.5), Abs(float32(5.5)))
	assert.Equal(t, float32(5.5), Abs(float32(-5.5)))
}

func TestMax(t *testing.T) {
	// Test with int
	assert.Equal(t, math.MaxInt, Max[int]())
	assert.Equal(t, int8(math.MaxInt8), Max[int8]())
	assert.Equal(t, int16(math.MaxInt16), Max[int16]())
	assert.Equal(t, int32(math.MaxInt32), Max[int32]())
	assert.Equal(t, int64(math.MaxInt64), Max[int64]())

	// Test with uint
	assert.Equal(t, uint(math.MaxUint), Max[uint]())
	assert.Equal(t, uint8(math.MaxUint8), Max[uint8]())
	assert.Equal(t, uint16(math.MaxUint16), Max[uint16]())
	assert.Equal(t, uint32(math.MaxUint32), Max[uint32]())
	assert.Equal(t, uint64(math.MaxUint64), Max[uint64]())
}

func TestSaturatingAdd(t *testing.T) {
	// Test normal addition
	assert.Equal(t, 5, SaturatingAdd(2, 3))
	assert.Equal(t, int8(5), SaturatingAdd(int8(2), int8(3)))

	// Test overflow (should return max)
	maxInt := Max[int]()
	assert.Equal(t, maxInt, SaturatingAdd(maxInt, 1))
	assert.Equal(t, maxInt, SaturatingAdd(maxInt-1, 2))

	// Test with uint
	maxUint := Max[uint]()
	assert.Equal(t, maxUint, SaturatingAdd(maxUint, 1))
}

func TestSaturatingSub(t *testing.T) {
	// Test normal subtraction
	assert.Equal(t, 2, SaturatingSub(5, 3))
	assert.Equal(t, int8(2), SaturatingSub(int8(5), int8(3)))

	// Test underflow (should return 0)
	assert.Equal(t, 0, SaturatingSub(3, 5))
	assert.Equal(t, int8(0), SaturatingSub(int8(3), int8(5)))

	// Test with float
	assert.Equal(t, 2.5, SaturatingSub(5.5, 3.0))
	assert.Equal(t, 0.0, SaturatingSub(3.0, 5.5))
}

func TestCheckedAdd(t *testing.T) {
	// Test normal addition
	result := CheckedAdd(2, 3)
	assert.True(t, result.IsSome())
	assert.Equal(t, 5, result.Unwrap())

	// Test overflow (should return None)
	maxInt := Max[int]()
	result2 := CheckedAdd(maxInt, 1)
	assert.True(t, result2.IsNone())

	// Test with uint
	maxUint := Max[uint]()
	result3 := CheckedAdd(maxUint, 1)
	assert.True(t, result3.IsNone())
	
	// Use gust to satisfy linter
	_ = gust.None[int]()
}

func TestCheckedMul(t *testing.T) {
	// Test normal multiplication
	result := CheckedMul(2, 3)
	assert.True(t, result.IsSome())
	assert.Equal(t, 6, result.Unwrap())

	// Test overflow (should return None)
	maxInt := Max[int]()
	result2 := CheckedMul(maxInt, 2)
	assert.True(t, result2.IsNone())

	// Test with zero
	result3 := CheckedMul(0, maxInt)
	assert.True(t, result3.IsSome())
	assert.Equal(t, 0, result3.Unwrap())
}

