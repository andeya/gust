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
	segments1 := SliceSegment([]string{"1", "2", "3", "4", "5", "6", "7"}, 2)
	segments2 := SliceSegment([]string{"1", "2", "3", "4", "5", "6", "7"}, 2, true)
	slice := [][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}, {"7"}}
	assert.Equal(t, slice, segments1)
	assert.Equal(t, slice, segments2)
	assert.Equal(t, [][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}}, SliceSegment([]string{"1", "2", "3", "4", "5", "6"}, 2))
}
