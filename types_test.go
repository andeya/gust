package gust

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPair_Split(t *testing.T) {
	// Test Pair.Split with int types
	pair := Pair[int, string]{A: 42, B: "hello"}
	a, b := pair.Split()
	assert.Equal(t, 42, a)
	assert.Equal(t, "hello", b)

	// Test Pair.Split with different types
	pair2 := Pair[string, int]{"test", 100}
	a2, b2 := pair2.Split()
	assert.Equal(t, "test", a2)
	assert.Equal(t, 100, b2)

	// Test Pair.Split with float types
	pair3 := Pair[float64, bool]{3.14, true}
	a3, b3 := pair3.Split()
	assert.Equal(t, 3.14, a3)
	assert.Equal(t, true, b3)
}

func TestVecEntry_Split(t *testing.T) {
	// Test VecEntry.Split with int
	entry := VecEntry[int]{Index: 5, Elem: 42}
	idx, elem := entry.Split()
	assert.Equal(t, 5, idx)
	assert.Equal(t, 42, elem)

	// Test VecEntry.Split with string
	entry2 := VecEntry[string]{Index: 0, Elem: "first"}
	idx2, elem2 := entry2.Split()
	assert.Equal(t, 0, idx2)
	assert.Equal(t, "first", elem2)

	// Test VecEntry.Split with negative index
	entry3 := VecEntry[bool]{Index: -1, Elem: true}
	idx3, elem3 := entry3.Split()
	assert.Equal(t, -1, idx3)
	assert.Equal(t, true, elem3)
}

func TestDictEntry_Split(t *testing.T) {
	// Test DictEntry.Split with string key and int value
	entry := DictEntry[string, int]{Key: "age", Value: 25}
	key, value := entry.Split()
	assert.Equal(t, "age", key)
	assert.Equal(t, 25, value)

	// Test DictEntry.Split with int key and string value
	entry2 := DictEntry[int, string]{Key: 1, Value: "one"}
	key2, value2 := entry2.Split()
	assert.Equal(t, 1, key2)
	assert.Equal(t, "one", value2)

	// Test DictEntry.Split with complex types
	entry3 := DictEntry[string, []int]{Key: "numbers", Value: []int{1, 2, 3}}
	key3, value3 := entry3.Split()
	assert.Equal(t, "numbers", key3)
	assert.Equal(t, []int{1, 2, 3}, value3)
}

func TestPair_Usage(t *testing.T) {
	// Test Pair usage in a function
	createPair := func(a int, b string) Pair[int, string] {
		return Pair[int, string]{A: a, B: b}
	}

	pair := createPair(10, "test")
	assert.Equal(t, 10, pair.A)
	assert.Equal(t, "test", pair.B)

	// Test accessing fields directly
	pair.A = 20
	pair.B = "updated"
	assert.Equal(t, 20, pair.A)
	assert.Equal(t, "updated", pair.B)
}

func TestVecEntry_Usage(t *testing.T) {
	// Test VecEntry usage
	entries := []VecEntry[string]{
		{Index: 0, Elem: "zero"},
		{Index: 1, Elem: "one"},
		{Index: 2, Elem: "two"},
	}

	for i, entry := range entries {
		idx, elem := entry.Split()
		assert.Equal(t, i, idx)
		assert.Equal(t, entries[i].Elem, elem)
	}
}

func TestDictEntry_Usage(t *testing.T) {
	// Test DictEntry usage
	entries := []DictEntry[string, int]{
		{Key: "a", Value: 1},
		{Key: "b", Value: 2},
		{Key: "c", Value: 3},
	}

	for _, entry := range entries {
		key, value := entry.Split()
		assert.NotEmpty(t, key)
		assert.Greater(t, value, 0)
	}
}
