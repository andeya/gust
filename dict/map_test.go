package dict

import (
	"testing"

	"github.com/andeya/gust/option"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	var m = map[string]string{"a": "b", "c": "d"}
	assert.Equal(t, option.Some("b"), Get(m, "a"))
	assert.Equal(t, option.None[string](), Get(m, "x"))
	var m2 map[string]string
	assert.Equal(t, option.None[string](), Get(m2, "x"))
}

func TestKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := Keys(m)
	assert.Len(t, keys, 3)
	assert.Contains(t, keys, "a")
	assert.Contains(t, keys, "b")
	assert.Contains(t, keys, "c")

	// Test with nil map
	var nilMap map[string]int
	assert.Nil(t, Keys(nilMap))
}

func TestValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	values := Values(m)
	assert.Len(t, values, 3)
	assert.Contains(t, values, 1)
	assert.Contains(t, values, 2)
	assert.Contains(t, values, 3)

	// Test with nil map
	var nilMap map[string]int
	assert.Nil(t, Values(nilMap))
}

func TestEntries(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	entries := Entries(m)
	assert.Len(t, entries, 2)

	entryMap := make(map[string]int)
	for _, entry := range entries {
		entryMap[entry.Key] = entry.Value
	}
	assert.Equal(t, 1, entryMap["a"])
	assert.Equal(t, 2, entryMap["b"])

	// Test with nil map
	var nilMap map[string]int
	assert.Nil(t, Entries(nilMap))
}

func TestVec(t *testing.T) {
	var m = map[string]string{"a": "b", "c": "d"}
	var s = Vec(m, func(k string, v string) string {
		return k + ":" + v
	})
	assert.Len(t, s, 2)
	assert.Contains(t, s, "a:b")
	assert.Contains(t, s, "c:d")

	// Test with nil map
	var nilMap map[string]string
	assert.Nil(t, Vec(nilMap, func(k string, v string) string {
		return k + ":" + v
	}))
}

func TestCopy(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	copied := Copy(m)
	assert.Equal(t, m, copied)

	// Modify original, copied should not change
	m["d"] = 4
	assert.NotEqual(t, m, copied)
	assert.Len(t, copied, 3)

	// Test with nil map
	var nilMap map[string]int
	assert.Nil(t, Copy(nilMap))
}

func TestEvery(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// All values > 0
	assert.True(t, Every(m, func(k string, v int) bool {
		return v > 0
	}))

	// Not all values > 1
	assert.False(t, Every(m, func(k string, v int) bool {
		return v > 1
	}))

	// Test with empty map (should return true)
	emptyMap := map[string]int{}
	assert.True(t, Every(emptyMap, func(k string, v int) bool {
		return false // Even with false predicate, empty map returns true
	}))
}

func TestSome(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// Some values > 2
	assert.True(t, Some(m, func(k string, v int) bool {
		return v > 2
	}))

	// No values > 5
	assert.False(t, Some(m, func(k string, v int) bool {
		return v > 5
	}))

	// Test with empty map (should return false)
	emptyMap := map[string]int{}
	assert.False(t, Some(emptyMap, func(k string, v int) bool {
		return true // Even with true predicate, empty map returns false
	}))
}

func TestFind(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// Find value == 2
	result := Find(m, func(k string, v int) bool {
		return v == 2
	})
	assert.True(t, result.IsSome())
	entry := result.Unwrap()
	assert.Equal(t, 2, entry.Value)

	// Find non-existent
	result2 := Find(m, func(k string, v int) bool {
		return v == 5
	})
	assert.True(t, result2.IsNone())
}

func TestFilter(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

	// Filter values > 2
	filtered := Filter(m, func(k string, v int) bool {
		return v > 2
	})
	assert.Len(t, filtered, 2)
	assert.Equal(t, 3, filtered["c"])
	assert.Equal(t, 4, filtered["d"])
	assert.NotContains(t, filtered, "a")
	assert.NotContains(t, filtered, "b")
}

func TestFilterMap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	// Filter and map to new key-value types
	filtered := FilterMap(m, func(k string, v int) option.Option[DictEntry[int, string]] {
		if v > 1 {
			return option.Some(DictEntry[int, string]{
				Key:   v,
				Value: k,
			})
		}
		return option.None[DictEntry[int, string]]()
	})
	assert.Len(t, filtered, 2)
	assert.Equal(t, "b", filtered[2])
	assert.Equal(t, "c", filtered[3])
	assert.NotContains(t, filtered, 1)
}

func TestFilterMapKey(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	filtered := FilterMapKey(m, func(k string, v int) option.Option[DictEntry[int, int]] {
		if v > 1 {
			return option.Some(DictEntry[int, int]{
				Key:   v,
				Value: v * 2,
			})
		}
		return option.None[DictEntry[int, int]]()
	})
	assert.Len(t, filtered, 2)
	assert.Equal(t, 4, filtered[2])
	assert.Equal(t, 6, filtered[3])
}

func TestFilterMapValue(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	filtered := FilterMapValue(m, func(k string, v int) option.Option[DictEntry[string, string]] {
		if v > 1 {
			return option.Some(DictEntry[string, string]{
				Key:   k,
				Value: k + ":" + string(rune('0'+v)),
			})
		}
		return option.None[DictEntry[string, string]]()
	})
	assert.Len(t, filtered, 2)
	assert.Contains(t, filtered, "b")
	assert.Contains(t, filtered, "c")
}

func TestMap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	mapped := Map(m, func(k string, v int) DictEntry[int, string] {
		return DictEntry[int, string]{
			Key:   v,
			Value: k,
		}
	})
	assert.Len(t, mapped, 2)
	assert.Equal(t, "a", mapped[1])
	assert.Equal(t, "b", mapped[2])

	// Test with nil map
	var nilMap map[string]int
	assert.Nil(t, Map(nilMap, func(k string, v int) DictEntry[int, string] {
		return DictEntry[int, string]{Key: v, Value: k}
	}))
}

func TestMapCurry(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	curried := MapCurry[string, int, int, string](m, func(k string) int {
		return int(k[0])
	})

	mapped := curried(func(v int) string {
		return string(rune('0' + v))
	})

	assert.Len(t, mapped, 2)
	// Key is int(k[0]), which is 97 for 'a' and 98 for 'b'
	// Value is string('0' + v), which is "1" for v=1 and "2" for v=2
	assert.Equal(t, "1", mapped[int('a')])
	assert.Equal(t, "2", mapped[int('b')])
}

func TestMapKey(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	mapped := MapKey(m, func(k string, v int) int {
		return int(k[0])
	})

	assert.Len(t, mapped, 2)
	assert.Equal(t, 1, mapped[int('a')])
	assert.Equal(t, 2, mapped[int('b')])
}

func TestMapKeyAlone(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	mapped := MapKeyAlone(m, func(k string) int {
		return int(k[0])
	})

	assert.Len(t, mapped, 2)
	assert.Equal(t, 1, mapped[int('a')])
	assert.Equal(t, 2, mapped[int('b')])
}

func TestMapValue(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	mapped := MapValue(m, func(k string, v int) string {
		return k + ":" + string(rune('0'+v))
	})

	assert.Len(t, mapped, 2)
	assert.Equal(t, "a:1", mapped["a"])
	assert.Equal(t, "b:2", mapped["b"])
}

func TestMapValueAlone(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	mapped := MapValueAlone(m, func(v int) string {
		return string(rune('0' + v))
	})

	assert.Len(t, mapped, 2)
	assert.Equal(t, "1", mapped["a"])
	assert.Equal(t, "2", mapped["b"])
}
