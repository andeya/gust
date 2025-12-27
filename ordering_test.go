package gust_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOrdering(t *testing.T) {
	// Test Less
	less := gust.Less()
	assert.True(t, less.IsLess())
	assert.False(t, less.IsEqual())
	assert.False(t, less.IsGreater())
	assert.True(t, less.Is(gust.Less()))

	// Test Equal
	equal := gust.Equal()
	assert.False(t, equal.IsLess())
	assert.True(t, equal.IsEqual())
	assert.False(t, equal.IsGreater())
	assert.True(t, equal.Is(gust.Equal()))

	// Test Greater
	greater := gust.Greater()
	assert.False(t, greater.IsLess())
	assert.False(t, greater.IsEqual())
	assert.True(t, greater.IsGreater())
	assert.True(t, greater.Is(gust.Greater()))

	// Test Compare with integers
	assert.True(t, gust.Compare(1, 2).IsLess())
	assert.True(t, gust.Compare(2, 2).IsEqual())
	assert.True(t, gust.Compare(3, 2).IsGreater())

	// Test Compare with strings
	assert.True(t, gust.Compare("a", "b").IsLess())
	assert.True(t, gust.Compare("b", "b").IsEqual())
	assert.True(t, gust.Compare("c", "b").IsGreater())

	// Test Compare with floats
	assert.True(t, gust.Compare(1.5, 2.5).IsLess())
	assert.True(t, gust.Compare(2.5, 2.5).IsEqual())
	assert.True(t, gust.Compare(3.5, 2.5).IsGreater())

	// Test Compare with uint
	assert.True(t, gust.Compare(uint(10), uint(20)).IsLess())
	assert.True(t, gust.Compare(uint(20), uint(20)).IsEqual())
	assert.True(t, gust.Compare(uint(30), uint(20)).IsGreater())

	// Test Compare with uintptr
	assert.True(t, gust.Compare(uintptr(10), uintptr(20)).IsLess())
	assert.True(t, gust.Compare(uintptr(20), uintptr(20)).IsEqual())
	assert.True(t, gust.Compare(uintptr(30), uintptr(20)).IsGreater())

	// Test Compare with int8
	assert.True(t, gust.Compare(int8(1), int8(2)).IsLess())
	assert.True(t, gust.Compare(int8(2), int8(2)).IsEqual())
	assert.True(t, gust.Compare(int8(3), int8(2)).IsGreater())

	// Test Compare with int16
	assert.True(t, gust.Compare(int16(1), int16(2)).IsLess())
	assert.True(t, gust.Compare(int16(2), int16(2)).IsEqual())
	assert.True(t, gust.Compare(int16(3), int16(2)).IsGreater())

	// Test Compare with int32
	assert.True(t, gust.Compare(int32(1), int32(2)).IsLess())
	assert.True(t, gust.Compare(int32(2), int32(2)).IsEqual())
	assert.True(t, gust.Compare(int32(3), int32(2)).IsGreater())

	// Test Compare with int64
	assert.True(t, gust.Compare(int64(1), int64(2)).IsLess())
	assert.True(t, gust.Compare(int64(2), int64(2)).IsEqual())
	assert.True(t, gust.Compare(int64(3), int64(2)).IsGreater())

	// Test Compare with uint8
	assert.True(t, gust.Compare(uint8(10), uint8(20)).IsLess())
	assert.True(t, gust.Compare(uint8(20), uint8(20)).IsEqual())
	assert.True(t, gust.Compare(uint8(30), uint8(20)).IsGreater())

	// Test Compare with uint16
	assert.True(t, gust.Compare(uint16(10), uint16(20)).IsLess())
	assert.True(t, gust.Compare(uint16(20), uint16(20)).IsEqual())
	assert.True(t, gust.Compare(uint16(30), uint16(20)).IsGreater())

	// Test Compare with uint32
	assert.True(t, gust.Compare(uint32(10), uint32(20)).IsLess())
	assert.True(t, gust.Compare(uint32(20), uint32(20)).IsEqual())
	assert.True(t, gust.Compare(uint32(30), uint32(20)).IsGreater())

	// Test Compare with uint64
	assert.True(t, gust.Compare(uint64(10), uint64(20)).IsLess())
	assert.True(t, gust.Compare(uint64(20), uint64(20)).IsEqual())
	assert.True(t, gust.Compare(uint64(30), uint64(20)).IsGreater())

	// Test Compare with float32
	assert.True(t, gust.Compare(float32(1.5), float32(2.5)).IsLess())
	assert.True(t, gust.Compare(float32(2.5), float32(2.5)).IsEqual())
	assert.True(t, gust.Compare(float32(3.5), float32(2.5)).IsGreater())
}
