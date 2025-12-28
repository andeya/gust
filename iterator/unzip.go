package iterator

import (
	"github.com/andeya/gust/pair"
)

// Unzip converts an iterator of pairs into a pair of containers.
//
// Unzip() consumes an entire iterator of pairs, producing two
// collections: one from the left elements of the pairs, and one
// from the right elements.
//
// This function is, in some sense, the opposite of Zip.
//
// # Examples
//
//	var a = []pair.Pair[int, string]{
//		{A: 1, B: "a"},
//		{A: 2, B: "b"},
//		{A: 3, B: "c"},
//	}
//	var (left, right) = Unzip(FromSlice(a))
//	assert.Equal(t, []int{1, 2, 3}, left)
//	assert.Equal(t, []string{"a", "b", "c"}, right)
func Unzip[T any, U any](iter Iterator[pair.Pair[T, U]]) ([]T, []U) {
	var left []T
	var right []U
	for {
		item := iter.Next()
		if item.IsNone() {
			break
		}
		pair := item.Unwrap()
		left = append(left, pair.A)
		right = append(right, pair.B)
	}
	return left, right
}
