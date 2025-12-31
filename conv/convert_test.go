package conv

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesToString(t *testing.T) {
	var b = []byte("abc")
	s := BytesToString[string](b)
	assert.Equal(t, string(b), s)
}

func TestSystemEndian(t *testing.T) {
	// SystemEndian should be initialized
	assert.NotNil(t, SystemEndian())

	// Test it works with binary functions
	var data [4]byte
	SystemEndian().PutUint32(data[:], 0x01020304)

	value := SystemEndian().Uint32(data[:])
	assert.Equal(t, uint32(0x01020304), value)

	// Verify it's either LittleEndian or BigEndian
	isLittleEndian := SystemEndian() == binary.LittleEndian
	isBigEndian := SystemEndian() == binary.BigEndian
	assert.True(t, isLittleEndian || isBigEndian, "SystemEndian should be either LittleEndian or BigEndian")
}

func TestStringToReadonlyBytes(t *testing.T) {
	var s = "abc"
	b := StringToReadonlyBytes(s)
	assert.Equal(t, []byte(s), b)
}

func TestUnsafeConvert(t *testing.T) {
	var s = "abc"
	b := UnsafeConvert[string, []byte](s)
	assert.Equal(t, []byte(s), b)
}

func TestToAnySlice(t *testing.T) {
	// Test with non-nil slice
	var s = []string{"a", "b", "c"}
	a := ToAnySlice(s)
	assert.Equal(t, []interface{}{"a", "b", "c"}, a)

	// Test with nil slice
	var nilSlice []string
	a2 := ToAnySlice(nilSlice)
	assert.Nil(t, a2)

	// Test with empty slice
	emptySlice := []string{}
	a3 := ToAnySlice(emptySlice)
	assert.Len(t, a3, 0)
	assert.NotNil(t, a3)
}

func TestToAnyMap(t *testing.T) {
	// Test with non-nil map
	var s = map[string]int{"a": 1, "b": 2, "c": 3}
	a := ToAnyMap(s)
	assert.Equal(t, map[string]interface{}{"a": 1, "b": 2, "c": 3}, a)

	// Test with nil map
	var nilMap map[string]int
	a2 := ToAnyMap(nilMap)
	assert.Nil(t, a2)

	// Test with empty map
	emptyMap := map[string]int{}
	a3 := ToAnyMap(emptyMap)
	assert.Len(t, a3, 0)
	assert.NotNil(t, a3)
}

func TestSafeAssert(t *testing.T) {
	var s any = "abc"
	assert.Equal(t, s, SafeAssert[string](s))
	assert.Equal(t, 0, SafeAssert[int](s))
}

func TestSafeAssertSlice(t *testing.T) {
	// Test with valid types
	var a = []interface{}{"a", "b", "c"}
	s := SafeAssertSlice[string](a).Unwrap()
	assert.Equal(t, []string{"a", "b", "c"}, s)

	// Test with invalid types
	result := SafeAssertSlice[int](a)
	assert.True(t, result.IsErr())

	// Test with nil slice
	var nilSlice []interface{}
	result2 := SafeAssertSlice[string](nilSlice)
	assert.True(t, result2.IsOk())
	assert.Nil(t, result2.Unwrap())

	// Test with empty slice
	emptySlice := []interface{}{}
	result3 := SafeAssertSlice[string](emptySlice)
	assert.True(t, result3.IsOk())
	assert.Equal(t, []string{}, result3.Unwrap())
}

func TestSafeAssertMap(t *testing.T) {
	// Test with valid types
	var a = map[string]interface{}{"a": 1, "b": 2, "c": 3}
	s := SafeAssertMap[string, int](a).Unwrap()
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, s)

	// Test with invalid types
	result := SafeAssertMap[string, string](a)
	assert.True(t, result.IsErr())

	// Test with nil map
	var nilMap map[string]interface{}
	result2 := SafeAssertMap[string, int](nilMap)
	assert.True(t, result2.IsOk())
	assert.Nil(t, result2.Unwrap())

	// Test with empty map
	emptyMap := map[string]interface{}{}
	result3 := SafeAssertMap[string, int](emptyMap)
	assert.True(t, result3.IsOk())
	assert.Equal(t, map[string]int{}, result3.Unwrap())
}

func TestUnsafeAssertSlice(t *testing.T) {
	var a = []interface{}{"a", "b", "c"}
	s := UnsafeAssertSlice[string](a)
	assert.Equal(t, []string{"a", "b", "c"}, s)
	assert.Panics(t, func() {
		UnsafeAssertSlice[int](a)
	})
}

func TestUnsafeAssertMap(t *testing.T) {
	var a = map[string]interface{}{"a": 1, "b": 2, "c": 3}
	s := UnsafeAssertMap[string, int](a)
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, s)
	assert.Panics(t, func() {
		UnsafeAssertMap[string, string](a)
	})
}

func BenchmarkBytesToString(b *testing.B) {
	var bs = []byte("abc")
	for i := 0; i < b.N; i++ {
		BytesToString[string](bs)
	}
}

func BenchmarkStringToReadonlyBytes(b *testing.B) {
	var s = "abc"
	for i := 0; i < b.N; i++ {
		StringToReadonlyBytes(s)
	}
}

func BenchmarkUnsafeConvert(b *testing.B) {
	var s = "abc"
	for i := 0; i < b.N; i++ {
		UnsafeConvert[string, []byte](s)
	}
}

func BenchmarkToAnySlice(b *testing.B) {
	var s = []string{"a", "b", "c"}
	for i := 0; i < b.N; i++ {
		ToAnySlice(s)
	}
}

func BenchmarkToAnyMap(b *testing.B) {
	var s = map[string]int{"a": 1, "b": 2, "c": 3}
	for i := 0; i < b.N; i++ {
		ToAnyMap(s)
	}
}
