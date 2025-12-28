package iterator_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/option"
	"github.com/andeya/gust/result"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	a := []int{1, 2, 3}
	assert.True(t, iterator.FromSlice(a).All(func(x int) bool { return x > 0 }))
	assert.False(t, iterator.FromSlice(a).All(func(x int) bool { return x > 2 }))
}

func TestAny(t *testing.T) {
	a := []int{1, 2, 3}
	assert.True(t, iterator.FromSlice(a).Any(func(x int) bool { return x > 0 }))
	assert.False(t, iterator.FromSlice(a).Any(func(x int) bool { return x > 5 }))
}

func TestFind(t *testing.T) {
	a := []int{1, 2, 3}
	assert.Equal(t, option.Some(2), iterator.FromSlice(a).Find(func(x int) bool { return x == 2 }))
	assert.Equal(t, option.None[int](), iterator.FromSlice(a).Find(func(x int) bool { return x == 5 }))
}

func TestPosition(t *testing.T) {
	a := []int{1, 2, 3}
	assert.Equal(t, option.Some(uint(1)), iterator.FromSlice(a).Position(func(x int) bool { return x == 2 }))
	assert.Equal(t, option.None[uint](), iterator.FromSlice(a).Position(func(x int) bool { return x == 5 }))
}

// TestAllEmpty tests All with empty iterator
func TestAllEmpty(t *testing.T) {
	iter := iterator.Empty[int]()
	result := iter.All(func(x int) bool { return x > 0 })
	assert.True(t, result) // iterator.Empty iterator returns true
}

// TestAnyEmpty tests Any with empty iterator
func TestAnyEmpty(t *testing.T) {
	iter := iterator.Empty[int]()
	result := iter.Any(func(x int) bool { return x > 0 })
	assert.False(t, result) // iterator.Empty iterator returns false
}

// TestFindEmpty tests Find with empty iterator
func TestFindEmpty(t *testing.T) {
	iter := iterator.Empty[int]()
	result := iter.Find(func(x int) bool { return x > 0 })
	assert.True(t, result.IsNone())
}

// TestPositionEmpty tests Position with empty iterator
func TestPositionEmpty(t *testing.T) {
	iter := iterator.Empty[int]()
	result := iter.Position(func(x int) bool { return x > 0 })
	assert.True(t, result.IsNone())
}

// TestFindMapEmpty tests iter.FindMap with empty iterator
func TestFindMapEmpty(t *testing.T) {
	iter := iterator.Empty[string]()
	result := iterator.FindMap(iter, func(s string) option.Option[int] {
		return option.Some(42)
	})
	assert.True(t, result.IsNone())
}

// TestFindMapBasic tests iter.FindMap with basic usage - finding first non-none result
func TestFindMapBasic(t *testing.T) {
	// Test case from documentation: find first parseable number
	a := []string{"lol", "NaN", "2", "5"}
	firstNumber := iterator.FindMap(iterator.FromSlice(a), func(s string) option.Option[int] {
		if v, err := strconv.Atoi(s); err == nil {
			return option.Some(v)
		}
		return option.None[int]()
	})
	assert.True(t, firstNumber.IsSome())
	assert.Equal(t, 2, firstNumber.Unwrap())
}

// TestFindMapAllNone tests iter.FindMap when all elements return None
func TestFindMapAllNone(t *testing.T) {
	a := []string{"lol", "NaN", "abc", "xyz"}
	result := iterator.FindMap(iterator.FromSlice(a), func(s string) option.Option[int] {
		if v, err := strconv.Atoi(s); err == nil {
			return option.Some(v)
		}
		return option.None[int]()
	})
	assert.True(t, result.IsNone())
}

// TestFindMapFirstElement tests iter.FindMap when first element returns Some
func TestFindMapFirstElement(t *testing.T) {
	a := []string{"1", "NaN", "2", "5"}
	result := iterator.FindMap(iterator.FromSlice(a), func(s string) option.Option[int] {
		if v, err := strconv.Atoi(s); err == nil {
			return option.Some(v)
		}
		return option.None[int]()
	})
	assert.True(t, result.IsSome())
	assert.Equal(t, 1, result.Unwrap())
}

// TestFindMapLastElement tests iter.FindMap when only last element returns Some
func TestFindMapLastElement(t *testing.T) {
	a := []string{"lol", "NaN", "abc", "42"}
	result := iterator.FindMap(iterator.FromSlice(a), func(s string) option.Option[int] {
		if v, err := strconv.Atoi(s); err == nil {
			return option.Some(v)
		}
		return option.None[int]()
	})
	assert.True(t, result.IsSome())
	assert.Equal(t, 42, result.Unwrap())
}

// TestFindMapShortCircuit tests iter.FindMap short-circuits after finding first Some
func TestFindMapShortCircuit(t *testing.T) {
	a := []string{"2", "3", "4", "5"}
	callCount := 0
	result := iterator.FindMap(iterator.FromSlice(a), func(s string) option.Option[int] {
		callCount++
		if v, err := strconv.Atoi(s); err == nil {
			return option.Some(v)
		}
		return option.None[int]()
	})
	assert.True(t, result.IsSome())
	assert.Equal(t, 2, result.Unwrap())
	// Should only call function once (short-circuit after first Some)
	assert.Equal(t, 1, callCount)
}

// TestFindMapTypeConversion tests iter.FindMap with type conversion
func TestFindMapTypeConversion(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	result := iterator.FindMap(iterator.FromSlice(a), func(x int) option.Option[string] {
		if x > 3 {
			return option.Some(strconv.Itoa(x * 2))
		}
		return option.None[string]()
	})
	assert.True(t, result.IsSome())
	assert.Equal(t, "8", result.Unwrap()) // First element > 3 is 4, 4*2 = 8
}

// TestTryFind tests TryFind functionality
func TestTryFind(t *testing.T) {
	a := []string{"1", "2", "lol", "NaN", "5"}
	iter := iterator.FromSlice(a)
	res := iter.TryFind(func(s string) result.Result[bool] {
		if s == "lol" {
			return result.TryErr[bool](errors.New("invalid"))
		}
		if v, err := strconv.Atoi(s); err == nil {
			return result.Ok(v == 2)
		}
		return result.Ok(false)
	})
	assert.True(t, res.IsOk())
}

func TestRfind(t *testing.T) {
	a := []int{1, 2, 3}
	iter := iterator.FromSlice(a)
	deIter := iter.MustToDoubleEnded()
	assert.Equal(t, option.Some(2), deIter.Rfind(func(x int) bool { return x == 2 }))

	b := []int{1, 2, 3}
	iter2 := iterator.FromSlice(b)
	deIter2 := iter2.MustToDoubleEnded()
	assert.Equal(t, option.None[int](), deIter2.Rfind(func(x int) bool { return x == 5 }))

	// Test function version
	iter3 := iterator.FromSlice(a)
	deIter3 := iter3.MustToDoubleEnded()
	assert.Equal(t, option.Some(2), iterator.Rfind(deIter3, func(x int) bool { return x == 2 }))
}

// TestRfindEmptyIterator tests Rfind with empty iterator
func TestRfindEmptyIterator(t *testing.T) {
	iter := iterator.FromSlice([]int{})
	deIter := iter.MustToDoubleEnded()
	result := deIter.Rfind(func(x int) bool { return x == 1 })
	assert.True(t, result.IsNone())
}

// TestIteratorExtConsumerMethods_FindSearch tests find/search methods from TestIteratorExtConsumerMethods
func TestIteratorExtConsumerMethods_FindSearch(t *testing.T) {
	// Test All
	iter3 := iterator.FromSlice([]int{2, 4, 6})
	assert.True(t, iter3.All(func(x int) bool { return x%2 == 0 }))

	// Test Any
	iter4 := iterator.FromSlice([]int{1, 2, 3})
	assert.True(t, iter4.Any(func(x int) bool { return x > 2 }))

	// Test Find
	iter5 := iterator.FromSlice([]int{1, 2, 3})
	assert.Equal(t, option.Some(2), iter5.Find(func(x int) bool { return x > 1 }))

	// Test TryFind
	iter16 := iterator.FromSlice([]string{"1", "2", "lol", "NaN", "5"})
	findResult := iter16.TryFind(func(s string) result.Result[bool] {
		if s == "lol" {
			return result.TryErr[bool](errors.New("invalid"))
		}
		if v, err := strconv.Atoi(s); err == nil {
			return result.Ok[bool](v == 2)
		}
		return result.Ok[bool](false)
	})
	assert.True(t, findResult.IsOk())
	assert.True(t, findResult.Unwrap().IsSome())
	assert.Equal(t, "2", findResult.Unwrap().Unwrap())
}

// TestIterator_WrapperMethods_FindSearch tests find/search wrapper methods from TestIterator_WrapperMethods
func TestIterator_WrapperMethods_FindSearch(t *testing.T) {
	// Test XFindMap (covers iterator_methods.go:416-418)
	iter1 := iterator.FromSlice([]string{"lol", "NaN", "2", "5"})
	result1 := iter1.XFindMap(func(s string) option.Option[any] {
		if v, err := strconv.Atoi(s); err == nil {
			return option.Some[any](v)
		}
		return option.None[any]()
	})
	assert.True(t, result1.IsSome())
	assert.Equal(t, 2, result1.UnwrapUnchecked())

	// Test FindMap (covers iterator_methods.go:434-436)
	iter2 := iterator.FromSlice([]string{"lol", "NaN", "2", "5"})
	result2 := iter2.FindMap(func(s string) option.Option[string] {
		if v, err := strconv.Atoi(s); err == nil {
			return option.Some(strconv.Itoa(v))
		}
		return option.None[string]()
	})
	assert.True(t, result2.IsSome())
	assert.Equal(t, "2", result2.UnwrapUnchecked())
}
