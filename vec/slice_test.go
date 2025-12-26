package vec

import (
	"strconv"
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	slice := []string{"Dodo", "Tiger", "Penguin", "Dodo"}
	val := Get(slice, 1)
	assert.Equal(t, gust.Some("Tiger"), val)
	assert.Equal(t, gust.None[string](), Get(slice, 10))
	assert.Equal(t, gust.None[string](), Get(slice, -10))
}

func TestDict(t *testing.T) {
	slice := []string{"Dodo", "Tiger", "Penguin", "Dodo"}
	dict := Dict(slice, func(m map[string]int, k int, v string) {
		if _, ok := m[v]; !ok {
			m[v] = k
		}
	})
	assert.Equal(t, map[string]int{"Dodo": 0, "Tiger": 1, "Penguin": 2}, dict)
}

func TestConcat(t *testing.T) {
	a := []string{"a", "0"}
	b := []string{"b", "1"}
	c := []string{"c", "2"}
	r := Concat(a, b, c)
	assert.Equal(t, []string{"a", "0", "b", "1", "c", "2"}, r)
}

func TestIntersect(t *testing.T) {
	x := []string{"a", "b", "a", "b", "b", "a", "a"}
	y := []string{"a", "b", "c", "a", "b", "c", "b", "c", "c"}
	z := []string{"a", "b", "c", "d", "a", "b", "c", "d", "b", "c", "d", "c", "d", "d"}
	r := Intersect(x, y, z)
	assert.Equal(t, map[string]int{"a": 2, "b": 3}, r)
}

func TestCopyWithin(t *testing.T) {
	slice := []string{"a", "b", "c", "d", "e"}
	CopyWithin(slice, 0, 3, 4)
	assert.Equal(t, []string{"d", "b", "c", "d", "e"}, slice)
	CopyWithin(slice, 1, -2)
	assert.Equal(t, []string{"d", "d", "e", "d", "e"}, slice)
}

func TestEvery(t *testing.T) {
	slice := []string{"1", "30", "39", "29", "10", "13"}
	isBelowThreshold := Every(slice, func(k int, v string) bool {
		i, _ := strconv.Atoi(v)
		return i < 40
	})
	assert.Equal(t, true, isBelowThreshold)

	// Test Every with return false branch (covers slice.go:78-80)
	slice2 := []string{"1", "30", "50", "29", "10", "13"}
	isBelowThreshold2 := Every(slice2, func(k int, v string) bool {
		i, _ := strconv.Atoi(v)
		return i < 40
	})
	assert.Equal(t, false, isBelowThreshold2) // "50" >= 40, so should return false
}

func TestFill(t *testing.T) {
	slice := []string{"a", "b", "c", "d"}
	Fill(slice, "?", 2, 4)
	assert.Equal(t, []string{"a", "b", "?", "?"}, slice)
	Fill(slice, "e", -1)
	assert.Equal(t, []string{"a", "b", "?", "e"}, slice)
}

func TestFilter(t *testing.T) {
	slice := []string{"spray", "limit", "elite", "exuberant", "destruction", "present"}
	result := Filter(slice, func(k int, v string) bool {
		return len(v) > 6
	})
	assert.Equal(t, []string{"exuberant", "destruction", "present"}, result)
}

func TestFilterMap(t *testing.T) {
	slice := []int8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result := FilterMap[int8, uint](slice, func(k int, v int8) gust.Option[uint] {
		if v > 6 {
			return gust.Some(uint(v))
		}
		return gust.None[uint]()
	})
	assert.Equal(t, []uint{7, 8, 9, 10}, result)
}

func TestFind(t *testing.T) {
	slice := []string{"spray", "limit", "elite", "exuberant", "destruction", "present"}
	entry := Find(slice, func(k int, v string) bool {
		return len(v) > 6
	})
	assert.Equal(t, gust.Some(gust.VecEntry[string]{Index: 3, Elem: "exuberant"}), entry)
}

func TestIncludes(t *testing.T) {
	slice := []string{"spray", "limit", "elite", "exuberant", "destruction", "present"}
	had := Includes(slice, "limit")
	assert.True(t, had)
	had = Includes(slice, "limit", 1)
	assert.True(t, had)
	had = Includes(slice, "limit", 2)
	assert.False(t, had)
}

func TestIndexOf(t *testing.T) {
	slice := []string{"spray", "limit", "elite", "exuberant", "destruction", "present"}
	idx := IndexOf(slice, "limit")
	assert.Equal(t, 1, idx)
	idx = IndexOf(slice, "limit", 1)
	assert.Equal(t, 1, idx)
	idx = IndexOf(slice, "limit", 10)
	assert.Equal(t, -1, idx)
}

func TestLastIndexOf(t *testing.T) {
	slice := []string{"Dodo", "Tiger", "Penguin", "Dodo"}
	idx := LastIndexOf(slice, "Dodo")
	assert.Equal(t, 3, idx)
	idx = LastIndexOf(slice, "Dodo", 1)
	assert.Equal(t, 3, idx)
	idx = LastIndexOf(slice, "Dodo", 10)
	assert.Equal(t, -1, idx)
	idx = LastIndexOf(slice, "?")
	assert.Equal(t, -1, idx)
}

func TestMap(t *testing.T) {
	slice := []string{"Dodo", "Tiger", "Penguin", "Dodo"}
	ret := Map(slice, func(k int, v string) string {
		return strconv.Itoa(k+1) + ":" + v
	})
	assert.Equal(t, []string{"1:Dodo", "2:Tiger", "3:Penguin", "4:Dodo"}, ret)
}

func TestPop(t *testing.T) {
	slice := []string{"kale", "tomato"}
	last := Pop(&slice)
	assert.Equal(t, gust.Some("tomato"), last)
	last = Pop(&slice)
	assert.Equal(t, gust.Some("kale"), last)
	last = Pop(&slice)
	assert.Equal(t, gust.None[string](), last)
}

func TestPushDistinct(t *testing.T) {
	slice := []string{"1", "2", "3", "4"}
	slice = PushDistinct(slice, "1", "5", "6", "1", "5", "6")
	assert.Equal(t, []string{"1", "2", "3", "4", "5", "6"}, slice)
}

func TestReduce(t *testing.T) {
	slice := []string{"1", "2", "3", "4"}
	reducer := Reduce(slice, func(k int, v string, accumulator string) string {
		return accumulator + "+" + v
	})
	assert.Equal(t, "1+2+3+4", reducer)
	reducer = Reduce(slice, func(k int, v string, accumulator string) string {
		return accumulator + "+" + v
	}, "100")
	assert.Equal(t, "100+1+2+3+4", reducer)
}

func TestReduceRight(t *testing.T) {
	slice := []string{"1", "2", "3", "4"}
	reducer := ReduceRight(slice, func(k int, v string, accumulator string) string {
		return accumulator + "+" + v
	})
	assert.Equal(t, "4+3+2+1", reducer)
	reducer = ReduceRight(slice, func(k int, v string, accumulator string) string {
		return accumulator + "+" + v
	}, "100")
	assert.Equal(t, "100+4+3+2+1", reducer)
}

func TestReverse(t *testing.T) {
	slice := []string{"1", "2", "3", "4"}
	Reverse(slice)
	assert.Equal(t, []string{"4", "3", "2", "1"}, slice)
}

func TestShift(t *testing.T) {
	slice := []string{"kale", "tomato"}
	first := Shift(&slice)
	assert.Equal(t, gust.Some("kale"), first)
	first = Pop(&slice)
	assert.Equal(t, gust.Some("tomato"), first)
	first = Pop(&slice)
	assert.Equal(t, gust.None[string](), first)
}

func TestSlice(t *testing.T) {
	slice := []string{"a", "b", "c", "d", "e"}
	sub := Slice(slice, 3)
	assert.Equal(t, []string{"d", "e"}, sub)
	sub = Slice(slice, 3, 4)
	assert.Equal(t, []string{"d"}, sub)
	sub = Slice(slice, 1, -2)
	assert.Equal(t, []string{"b", "c"}, sub)
	sub[0] = "x"
	assert.Equal(t, []string{"x", "c"}, sub)
	assert.Equal(t, []string{"a", "b", "c", "d", "e"}, slice)
}

func TestSome(t *testing.T) {
	slice := []string{"1", "30", "39", "29", "10", "13"}
	even := Some(slice, func(k int, v string) bool {
		i, _ := strconv.Atoi(v)
		return i%2 == 0
	})
	assert.Equal(t, true, even)
}

func TestSplice(t *testing.T) {
	slice := []string{"0", "1", "2", "3", "4"}
	Splice(&slice, 0, 0, "a", "b")
	assert.Equal(t, []string{"a", "b", "0", "1", "2", "3", "4"}, slice)

	slice = []string{"0", "1", "2", "3", "4"}
	Splice(&slice, 10, 0, "a", "b")
	assert.Equal(t, []string{"0", "1", "2", "3", "4", "a", "b"}, slice)

	slice = []string{"0", "1", "2", "3", "4"}
	Splice(&slice, 1, 0, "a", "b")
	assert.Equal(t, []string{"0", "a", "b", "1", "2", "3", "4"}, slice)

	slice = []string{"0", "1", "2", "3", "4"}
	Splice(&slice, 1, 2, "a", "b")
	assert.Equal(t, []string{"0", "a", "b", "3", "4"}, slice)

	slice = []string{"0", "1", "2", "3", "4"}
	Splice(&slice, 1, 10, "a", "b")
	assert.Equal(t, []string{"0", "a", "b"}, slice)
}

func TestUnshift(t *testing.T) {
	slice := []string{"0", "1", "2", "3", "4"}
	n := Unshift(&slice, "a", "b")
	assert.Equal(t, len(slice), n)
	assert.Equal(t, []string{"a", "b", "0", "1", "2", "3", "4"}, slice)
}

func TestUnshiftDistinct(t *testing.T) {
	slice := []string{"1", "2", "3", "4"}
	n := UnshiftDistinct(&slice, "-1", "0", "-1", "0", "1", "1")
	assert.Equal(t, len(slice), n)
	assert.Equal(t, []string{"-1", "0", "1", "2", "3", "4"}, slice)
}

func TestDistinct(t *testing.T) {
	slice := []string{"-1", "0", "-1", "0", "1", "1"}
	distinctCount := Distinct(&slice, true)
	assert.Equal(t, len(slice), len(distinctCount))
	assert.Equal(t, []string{"-1", "0", "1"}, slice)
	assert.Equal(t, map[string]int{"-1": 2, "0": 2, "1": 2}, distinctCount)
	var null []string
	m := Distinct(&null, true)
	assert.Equal(t, map[string]int{}, m)
	assert.Equal(t, []string(nil), null)
	assert.Equal(t, map[string]int{}, Distinct(&null, false))
}

func TestDistinctBy(t *testing.T) {
	data := []string{"-1", "0", "-12", "0", "1", "1234"}
	slice := Copy(data)
	mapping := func(k int, v string) int { return len(v) }
	DistinctBy(&slice, mapping)
	assert.Equal(t, []string{"-1", "0", "-12", "1234"}, slice)

	slice = Copy(data)
	DistinctBy(&slice, mapping, func(a, b string) string {
		if a < b {
			return b
		}
		return a
	})
	assert.Equal(t, []string{"-1", "1", "-12", "1234"}, slice)
}

func TestDistinctMap(t *testing.T) {
	slice := []string{"-1", "0", "-1", "0", "1", "1"}
	newSlice := DistinctMap(slice, func(k int, v string) string { return v + "0" })
	assert.Equal(t, []string{"-10", "00", "10"}, newSlice)
	newSlice2 := DistinctMap(slice, func(k int, v string) int { return len(v) })
	assert.Equal(t, []int{2, 1}, newSlice2)
}

func TestRemoveFirst(t *testing.T) {
	slice := []string{"-1", "0", "-1", "0", "1", "1"}
	n := RemoveFirst(&slice, "0")
	assert.Equal(t, len(slice), n)
	assert.Equal(t, []string{"-1", "-1", "0", "1", "1"}, slice)
}

func TestRemoveEvery(t *testing.T) {
	slice := []string{"-1", "0", "-1", "0", "1", "1"}
	n := RemoveEvery(&slice, "0")
	assert.Equal(t, len(slice), n)
	assert.Equal(t, []string{"-1", "-1", "1", "1"}, slice)
}

func TestStringSet(t *testing.T) {
	set1 := []string{"1", "1", "1", "2", "3", "6", "8"}
	set2 := []string{"2", "3", "5", "2", "2", "0"}
	set3 := []string{"2", "6", "6", "6", "7"}
	un := SetsUnion(set1, set2, set3)
	assert.Equal(t, []string{"1", "2", "3", "6", "8", "5", "0", "7"}, un)
	in := SetsIntersect(set1, set2, set3)
	assert.Equal(t, []string{"2"}, in)
	di := SetsDifference(set1, set2, set3)
	assert.Equal(t, []string{"1", "8"}, di)
}

func TestFlatten(t *testing.T) {
	slice := [][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}}
	flatten := Flatten(slice)
	assert.Equal(t, []string{"1", "2", "3", "4", "5", "6"}, flatten)
}

func TestFlatMap(t *testing.T) {
	slice := [][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}}
	flatten := FlatMap(slice, func(s string) string {
		return "-" + s
	})
	assert.Equal(t, []string{"-1", "-2", "-3", "-4", "-5", "-6"}, flatten)
}

func TestSliceSegment(t *testing.T) {
	a := []string{"1", "2", "3", "4", "5", "6"}
	b := append(a, "7")
	segments1 := SliceSegment(b, 2)
	segments2 := SliceSegment(b, 2, true)
	segments3 := SliceSegment(b, len(b)+3)
	segments4 := SliceSegment(b, len(b))
	segments5 := SliceSegment(b, -1)
	segments6 := SliceSegment(b, 0)
	x := [][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}}
	y := append(x, []string{"7"})
	assert.Equal(t, y, segments1)
	assert.Equal(t, y, segments2)
	assert.Equal(t, [][]string{b}, segments3)
	assert.Equal(t, [][]string{b}, segments4)
	assert.Equal(t, [][]string{b}, segments5)
	assert.Equal(t, [][]string(nil), segments6)
	assert.Equal(t, x, SliceSegment(a, 2))
}

func TestCopyWithin_EdgeCases(t *testing.T) {
	// Test CopyWithin with target == len(s) (should return early)
	slice := []string{"a", "b", "c"}
	CopyWithin(slice, 3, 0)
	assert.Equal(t, []string{"a", "b", "c"}, slice) // Should not change

	// Test CopyWithin with negative indices
	slice2 := []string{"a", "b", "c", "d", "e"}
	CopyWithin(slice2, -2, -3, -1)
	// target = fixIndex(5, -2, true) = 3
	// start = fixIndex(5, -3, true) = 2
	// end = fixIndex(5, -1, true) = 4
	// sub = Slice(slice2, 2, 4) = slice2[2:4] = {"c", "d"}
	// Copy to target=3: slice2[3] = "c", slice2[4] = "d"
	assert.Equal(t, []string{"a", "b", "c", "c", "d"}, slice2)
}

func TestFill_EdgeCases(t *testing.T) {
	// Test Fill with invalid range (should return early)
	slice := []string{"a", "b", "c"}
	Fill(slice, "x", 5, 10)                         // Invalid range
	assert.Equal(t, []string{"a", "b", "c"}, slice) // Should not change

	// Test Fill with negative indices
	slice2 := []string{"a", "b", "c", "d"}
	Fill(slice2, "x", -2, -1)
	assert.Equal(t, []string{"a", "b", "x", "d"}, slice2)
}

func TestSplice_ReplaceMode(t *testing.T) {
	// Test Splice with replace mode (deleteCount > 0)
	slice := []string{"a", "b", "c", "d"}
	Splice(&slice, 1, 2, "x", "y")
	assert.Equal(t, []string{"a", "x", "y", "d"}, slice)
}

func TestSplice_InsertMode(t *testing.T) {
	// Test Splice with insert mode (deleteCount == 0)
	slice := []string{"a", "b", "c"}
	Splice(&slice, 1, 0, "x", "y")
	assert.Equal(t, []string{"a", "x", "y", "b", "c"}, slice)
}

func TestSplice_DeleteOnly(t *testing.T) {
	// Test Splice with delete only (no items)
	slice := []string{"a", "b", "c", "d"}
	Splice(&slice, 1, 2)
	assert.Equal(t, []string{"a", "d"}, slice)
}

func TestSplice_NegativeDeleteCount(t *testing.T) {
	// Test Splice with negative deleteCount (should be set to 0)
	slice := []string{"a", "b", "c"}
	Splice(&slice, 1, -1, "x")
	assert.Equal(t, []string{"a", "x", "b", "c"}, slice)
}

func TestForEachSegment_Error(t *testing.T) {
	// Test ForEachSegment with callback error
	slice := []string{"a", "b", "c", "d"}
	err := ForEachSegment(slice, 2, func(slice []string) error {
		return assert.AnError
	})
	assert.Error(t, err)
}

func TestForEachSegment_Clone(t *testing.T) {
	// Test ForEachSegment with clone = true
	slice := []string{"a", "b", "c", "d"}
	var segments [][]string
	err := ForEachSegment(slice, 2, func(s []string) error {
		segments = append(segments, s)
		return nil
	}, true)
	assert.NoError(t, err)
	assert.Len(t, segments, 2)
}

func TestForEachSegment_NegativeLength(t *testing.T) {
	// Test ForEachSegment with negative maxSegmentLength
	slice := []string{"a", "b", "c"}
	var segments [][]string
	err := ForEachSegment(slice, -1, func(s []string) error {
		segments = append(segments, s)
		return nil
	})
	assert.NoError(t, err)
	assert.Len(t, segments, 1)
	assert.Len(t, segments[0], 3)
}

func TestSliceSegment_EdgeCases(t *testing.T) {
	// Test SliceSegment with maxSegmentLength = 0
	slice := []string{"a", "b", "c"}
	result := SliceSegment(slice, 0)
	assert.Nil(t, result)

	// Test SliceSegment with negative maxSegmentLength
	result2 := SliceSegment(slice, -1)
	assert.Len(t, result2, 1)
	assert.Len(t, result2[0], 3)
}

func TestDistinct_ChangeSlice(t *testing.T) {
	// Test Distinct with changeSlice = true
	slice := []string{"a", "b", "a", "c", "b"}
	count := Distinct(&slice, true)
	assert.Equal(t, map[string]int{"a": 2, "b": 2, "c": 1}, count)
	assert.Len(t, slice, 3)
}

func TestDistinctBy_WithWhoToKeep(t *testing.T) {
	// Test DistinctBy with whoToKeep function
	slice := []int{1, 2, 3, 2, 4}
	DistinctBy(&slice, func(k int, v int) int { return v }, func(a, b int) int {
		if a > b {
			return a
		}
		return b
	})
	assert.Len(t, slice, 4)
}

func TestIntersect_EmptySlice(t *testing.T) {
	// Test Intersect with empty slice
	result := Intersect([]string{}, []string{"a"})
	assert.Nil(t, result)

	result2 := Intersect([]string{"a"}, []string{})
	assert.Nil(t, result2)
}

func TestIntersect_NoArgs(t *testing.T) {
	// Test Intersect with no arguments
	result := Intersect[string]()
	assert.Nil(t, result)
}

func TestSetsDifference_MultipleOthers(t *testing.T) {
	// Test SetsDifference with multiple others
	set1 := []string{"a", "b", "c", "d"}
	set2 := []string{"b", "c"}
	set3 := []string{"c", "e"}
	result := SetsDifference(set1, set2, set3)
	assert.Equal(t, []string{"a", "d"}, result)
}

func TestUnshiftDistinct_EmptyElement(t *testing.T) {
	// Test UnshiftDistinct with empty element slice
	slice := []string{"a", "b"}
	length := UnshiftDistinct(&slice)
	assert.Equal(t, 2, length)
	assert.Equal(t, []string{"a", "b"}, slice)
}

func TestRemoveFirst_DuplicateElements(t *testing.T) {
	// Test RemoveFirst with duplicate elements in elements slice
	slice := []string{"a", "b", "c", "d"}
	length := RemoveFirst(&slice, "b", "b", "c")
	assert.Equal(t, 2, length)
	assert.Equal(t, []string{"a", "d"}, slice)
}

func TestRemoveEvery_DuplicateElements(t *testing.T) {
	// Test RemoveEvery with duplicate elements in elements slice
	slice := []string{"a", "b", "b", "c", "b"}
	length := RemoveEvery(&slice, "b", "b")
	assert.Equal(t, 2, length)
	assert.Equal(t, []string{"a", "c"}, slice)
}

func TestIndexOf_WithFromIndex(t *testing.T) {
	// Test IndexOf with fromIndex parameter
	slice := []string{"a", "b", "c", "b", "d"}
	idx := IndexOf(slice, "b", 2)
	assert.Equal(t, 3, idx) // Should find "b" at index 3, not index 1

	idx2 := IndexOf(slice, "b", 4)
	assert.Equal(t, -1, idx2) // Should not find "b" starting from index 4

	idx3 := IndexOf(slice, "a", 0)
	assert.Equal(t, 0, idx3)
}

func TestLastIndexOf_WithFromIndex(t *testing.T) {
	// Test LastIndexOf with fromIndex parameter
	// LastIndexOf searches backwards from the end to fromIndex (inclusive)
	slice := []string{"a", "b", "c", "b", "d"}

	// LastIndexOf(slice, "b", 2) searches from index 4 backwards to index 2 (inclusive)
	// Checks: index 4 ("d"), index 3 ("b" - found!), returns 3
	idx := LastIndexOf(slice, "b", 2)
	assert.Equal(t, 3, idx)

	// LastIndexOf(slice, "b", 4) searches from index 4 backwards to index 4 (inclusive)
	// Checks: index 4 ("d"), doesn't find "b", returns -1
	idx2 := LastIndexOf(slice, "b", 4)
	assert.Equal(t, -1, idx2)

	// LastIndexOf(slice, "d", 4) searches from index 4 backwards to index 4 (inclusive)
	// Checks: index 4 ("d" - found!), returns 4
	idx3 := LastIndexOf(slice, "d", 4)
	assert.Equal(t, 4, idx3)

	// LastIndexOf(slice, "b", 3) searches from index 4 backwards to index 3 (inclusive)
	// Checks: index 4 ("d"), index 3 ("b" - found!), returns 3
	idx4 := LastIndexOf(slice, "b", 3)
	assert.Equal(t, 3, idx4)

	// LastIndexOf(slice, "b", 1) searches from index 4 backwards to index 1 (inclusive)
	// Checks: index 4 ("d"), index 3 ("b" - found!), returns 3 (stops at first match)
	idx5 := LastIndexOf(slice, "b", 1)
	assert.Equal(t, 3, idx5)

	// LastIndexOf(slice, "b", 0) searches from index 4 backwards to index 0 (inclusive)
	// Checks: index 4 ("d"), index 3 ("b" - found!), returns 3
	idx6 := LastIndexOf(slice, "b", 0)
	assert.Equal(t, 3, idx6)
}

func TestIncludes_WithFromIndex(t *testing.T) {
	// Test Includes with fromIndex parameter
	slice := []string{"a", "b", "c", "b", "d"}
	assert.True(t, Includes(slice, "b", 2))
	assert.False(t, Includes(slice, "b", 4))
	assert.True(t, Includes(slice, "a", 0))
}

func TestMapAlone(t *testing.T) {
	// Test MapAlone function
	slice := []int{1, 2, 3}
	result := MapAlone(slice, func(v int) int { return v * 2 })
	assert.Equal(t, []int{2, 4, 6}, result)

	// Test with nil slice
	var nilSlice []int
	result2 := MapAlone(nilSlice, func(v int) int { return v * 2 })
	assert.Nil(t, result2)
}

func TestPushDistinctBy(t *testing.T) {
	// Test PushDistinctBy function
	slice := []int{1, 2, 3}
	result := PushDistinctBy(slice, func(a, b int) bool { return a == b }, 2, 4, 3)
	assert.Equal(t, []int{1, 2, 3, 4}, result) // 2 and 3 already exist, only 4 is added
}

func TestSetsIntersect(t *testing.T) {
	// Test SetsIntersect function
	set1 := []string{"a", "b", "c"}
	set2 := []string{"b", "c", "d"}
	set3 := []string{"c", "d", "e"}
	result := SetsIntersect(set1, set2, set3)
	assert.Equal(t, []string{"c"}, result) // Only "c" is in all three sets

	// Test with no intersection
	set4 := []string{"x", "y", "z"}
	result2 := SetsIntersect(set1, set4)
	assert.Equal(t, []string{}, result2)
}

func TestSetsDifference(t *testing.T) {
	// Test SetsDifference function
	set1 := []string{"a", "b", "c", "d"}
	set2 := []string{"b", "c"}
	set3 := []string{"d"}
	result := SetsDifference(set1, set2, set3)
	assert.Equal(t, []string{"a"}, result) // Only "a" is in set1 but not in set2 or set3

	// Test with empty difference
	result2 := SetsDifference(set1, set1)
	assert.Equal(t, []string{}, result2)
}

func TestFlatten_EdgeCases(t *testing.T) {
	// Test Flatten function with edge cases
	nested := [][]int{{1, 2}, {3, 4}, {5}}
	result := Flatten(nested)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)

	// Test with nil
	var nilSlice [][]int
	result2 := Flatten(nilSlice)
	assert.Nil(t, result2)

	// Test with empty slices
	empty := [][]int{{}, {1}, {}}
	result3 := Flatten(empty)
	assert.Equal(t, []int{1}, result3)
}

func TestFlatMap_EdgeCases(t *testing.T) {
	// Test FlatMap function with edge cases
	nested := [][]int{{1, 2}, {3, 4}}
	result := FlatMap(nested, func(x int) int { return x * 2 })
	assert.Equal(t, []int{2, 4, 6, 8}, result)

	// Test with nil
	var nilSlice [][]int
	result2 := FlatMap(nilSlice, func(x int) int { return x * 2 })
	assert.Nil(t, result2)

	// Test with empty slices
	empty := [][]int{{}, {1}, {}}
	result3 := FlatMap(empty, func(x int) string { return "a" })
	assert.Equal(t, []string{"a"}, result3)
}
