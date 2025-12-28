package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iterator"
	"github.com/stretchr/testify/assert"
)

func TestPeekable(t *testing.T) {
	xs := []int{1, 2, 3}
	iter := iterator.FromSlice(xs).Peekable()

	// peek() lets us see into the future
	assert.Equal(t, gust.Some(1), iter.Peek())
	assert.Equal(t, gust.Some(1), iter.Next())

	assert.Equal(t, gust.Some(2), iter.Next())

	// we can peek() multiple times, the iterator won't advance
	assert.Equal(t, gust.Some(3), iter.Peek())
	assert.Equal(t, gust.Some(3), iter.Peek())

	assert.Equal(t, gust.Some(3), iter.Next())

	// after the iterator is finished, so is peek()
	assert.Equal(t, gust.None[int](), iter.Peek())
	assert.Equal(t, gust.None[int](), iter.Next())
}

func TestPeekableInheritsIteratorMethods(t *testing.T) {
	// Test that PeekableIterator can use all Iterator methods
	xs := []int{1, 2, 3, 4, 5}
	iter := iterator.FromSlice(xs).Peekable()

	// Test Filter (Iterator method)
	filtered := iter.Filter(func(x int) bool { return x > 2 })
	assert.Equal(t, gust.Some(3), filtered.Next())
	assert.Equal(t, gust.Some(4), filtered.Next())
	assert.Equal(t, gust.Some(5), filtered.Next())
	assert.Equal(t, gust.None[int](), filtered.Next())

	// Test Map (Iterator method)
	xs2 := []int{1, 2, 3}
	iter2 := iterator.FromSlice(xs2).Peekable()
	mapped := iter2.Map(func(x int) int { return x * 2 })
	assert.Equal(t, gust.Some(2), mapped.Next())
	assert.Equal(t, gust.Some(4), mapped.Next())
	assert.Equal(t, gust.Some(6), mapped.Next())

	// Test Take (Iterator method)
	xs3 := []int{1, 2, 3, 4, 5}
	iter3 := iterator.FromSlice(xs3).Peekable()
	taken := iter3.Take(2)
	assert.Equal(t, gust.Some(1), taken.Next())
	assert.Equal(t, gust.Some(2), taken.Next())
	assert.Equal(t, gust.None[int](), taken.Next())

	// Test Collect (Iterator method)
	xs4 := []int{1, 2, 3}
	iter4 := iterator.FromSlice(xs4).Peekable()
	collected := iter4.Collect()
	assert.Equal(t, []int{1, 2, 3}, collected)

	// Test that Peek still works after using Iterator methods
	xs5 := []int{1, 2, 3}
	iter5 := iterator.FromSlice(xs5).Peekable()
	assert.Equal(t, gust.Some(1), iter5.Peek())
	filtered2 := iter5.Filter(func(x int) bool { return x > 1 })
	assert.Equal(t, gust.Some(2), filtered2.Next())
}

// TestCloned_NilPointer tests Cloned with nil pointer (covers special_iterators.go:126-129)
func TestCloned_NilPointer(t *testing.T) {
	ptrs := []*int{nil, intPtr(1), nil, intPtr(2)}
	iter := iterator.Cloned(iterator.FromSlice(ptrs))
	result := iter.Collect()
	// nil pointers should be converted to zero values
	assert.Equal(t, []int{0, 1, 0, 2}, result)
}

func intPtr(i int) *int {
	return &i
}

func TestCloned(t *testing.T) {
	a := []string{"hello", "world"}
	ptrs := make([]*string, len(a))
	for i := range a {
		ptrs[i] = &a[i]
	}
	iter := iterator.Cloned(iterator.FromSlice(ptrs))
	v := iter.Collect()
	assert.Equal(t, []string{"hello", "world"}, v)
}
