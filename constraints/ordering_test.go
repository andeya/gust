package constraints_test

import (
	"testing"

	"github.com/andeya/gust/constraints"
	"github.com/stretchr/testify/assert"
)

func TestOrdering(t *testing.T) {
	// Test Less
	less := constraints.Less()
	assert.True(t, less.IsLess())
	assert.False(t, less.IsEqual())
	assert.False(t, less.IsGreater())
	assert.True(t, less.Is(constraints.Less()))

	// Test Equal
	equal := constraints.Equal()
	assert.False(t, equal.IsLess())
	assert.True(t, equal.IsEqual())
	assert.False(t, equal.IsGreater())
	assert.True(t, equal.Is(constraints.Equal()))

	// Test Greater
	greater := constraints.Greater()
	assert.False(t, greater.IsLess())
	assert.False(t, greater.IsEqual())
	assert.True(t, greater.IsGreater())
	assert.True(t, greater.Is(constraints.Greater()))

	// Test Compare with integers
	assert.True(t, constraints.Compare(1, 2).IsLess())
	assert.True(t, constraints.Compare(2, 2).IsEqual())
	assert.True(t, constraints.Compare(3, 2).IsGreater())

	// Test Compare with strings
	assert.True(t, constraints.Compare("a", "b").IsLess())
	assert.True(t, constraints.Compare("b", "b").IsEqual())
	assert.True(t, constraints.Compare("c", "b").IsGreater())

	// Test Compare with floats
	assert.True(t, constraints.Compare(1.5, 2.5).IsLess())
	assert.True(t, constraints.Compare(2.5, 2.5).IsEqual())
	assert.True(t, constraints.Compare(3.5, 2.5).IsGreater())

	// Test Compare with uint
	assert.True(t, constraints.Compare(uint(10), uint(20)).IsLess())
	assert.True(t, constraints.Compare(uint(20), uint(20)).IsEqual())
	assert.True(t, constraints.Compare(uint(30), uint(20)).IsGreater())

	// Test Compare with uintptr
	assert.True(t, constraints.Compare(uintptr(10), uintptr(20)).IsLess())
	assert.True(t, constraints.Compare(uintptr(20), uintptr(20)).IsEqual())
	assert.True(t, constraints.Compare(uintptr(30), uintptr(20)).IsGreater())

	// Test Compare with int8
	assert.True(t, constraints.Compare(int8(1), int8(2)).IsLess())
	assert.True(t, constraints.Compare(int8(2), int8(2)).IsEqual())
	assert.True(t, constraints.Compare(int8(3), int8(2)).IsGreater())

	// Test Compare with int16
	assert.True(t, constraints.Compare(int16(1), int16(2)).IsLess())
	assert.True(t, constraints.Compare(int16(2), int16(2)).IsEqual())
	assert.True(t, constraints.Compare(int16(3), int16(2)).IsGreater())

	// Test Compare with int32
	assert.True(t, constraints.Compare(int32(1), int32(2)).IsLess())
	assert.True(t, constraints.Compare(int32(2), int32(2)).IsEqual())
	assert.True(t, constraints.Compare(int32(3), int32(2)).IsGreater())

	// Test Compare with int64
	assert.True(t, constraints.Compare(int64(1), int64(2)).IsLess())
	assert.True(t, constraints.Compare(int64(2), int64(2)).IsEqual())
	assert.True(t, constraints.Compare(int64(3), int64(2)).IsGreater())

	// Test Compare with uint8
	assert.True(t, constraints.Compare(uint8(10), uint8(20)).IsLess())
	assert.True(t, constraints.Compare(uint8(20), uint8(20)).IsEqual())
	assert.True(t, constraints.Compare(uint8(30), uint8(20)).IsGreater())

	// Test Compare with uint16
	assert.True(t, constraints.Compare(uint16(10), uint16(20)).IsLess())
	assert.True(t, constraints.Compare(uint16(20), uint16(20)).IsEqual())
	assert.True(t, constraints.Compare(uint16(30), uint16(20)).IsGreater())

	// Test Compare with uint32
	assert.True(t, constraints.Compare(uint32(10), uint32(20)).IsLess())
	assert.True(t, constraints.Compare(uint32(20), uint32(20)).IsEqual())
	assert.True(t, constraints.Compare(uint32(30), uint32(20)).IsGreater())

	// Test Compare with uint64
	assert.True(t, constraints.Compare(uint64(10), uint64(20)).IsLess())
	assert.True(t, constraints.Compare(uint64(20), uint64(20)).IsEqual())
	assert.True(t, constraints.Compare(uint64(30), uint64(20)).IsGreater())

	// Test Compare with float32
	assert.True(t, constraints.Compare(float32(1.5), float32(2.5)).IsLess())
	assert.True(t, constraints.Compare(float32(2.5), float32(2.5)).IsEqual())
	assert.True(t, constraints.Compare(float32(3.5), float32(2.5)).IsGreater())
}
