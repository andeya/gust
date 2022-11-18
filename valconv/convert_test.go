package valconv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesToString(t *testing.T) {
	var b = []byte("abc")
	s := BytesToString[string](b)
	assert.Equal(t, string(b), s)
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
	var s = []string{"a", "b", "c"}
	a := ToAnySlice(s)
	assert.Equal(t, []interface{}{"a", "b", "c"}, a)
}

func TestToAnyMap(t *testing.T) {
	var s = map[string]int{"a": 1, "b": 2, "c": 3}
	a := ToAnyMap(s)
	assert.Equal(t, map[string]interface{}{"a": 1, "b": 2, "c": 3}, a)
}

func TestSafeAssert(t *testing.T) {
	var s any = "abc"
	assert.Equal(t, s, SafeAssert[string](s))
}

func TestUnsafeAssertSlice(t *testing.T) {
	var a = []interface{}{"a", "b", "c"}
	s := UnsafeAssertSlice[string](a)
	assert.Equal(t, []string{"a", "b", "c"}, s)
}

func TestUnsafeAssertMap(t *testing.T) {
	var a = map[string]interface{}{"a": 1, "b": 2, "c": 3}
	s := UnsafeAssertMap[string, int](a)
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, s)
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
