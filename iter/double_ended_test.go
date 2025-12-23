package iter

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestSliceIteratorDoubleEnded(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5, 6}
	iter := FromSlice(numbers)
	deIter := iter.MustToDoubleEnded()

	assert.Equal(t, gust.Some(1), deIter.Next())
	assert.Equal(t, gust.Some(6), deIter.NextBack())
	assert.Equal(t, gust.Some(5), deIter.NextBack())
	assert.Equal(t, gust.Some(2), deIter.Next())
	assert.Equal(t, gust.Some(3), deIter.Next())
	assert.Equal(t, gust.Some(4), deIter.Next())
	assert.Equal(t, gust.None[int](), deIter.Next())
	assert.Equal(t, gust.None[int](), deIter.NextBack())
}

func TestSliceIteratorNextBack(t *testing.T) {
	numbers := []int{1, 2, 3}
	iter := FromSlice(numbers)
	deIter := iter.MustToDoubleEnded()

	assert.Equal(t, gust.Some(3), deIter.NextBack())
	assert.Equal(t, gust.Some(2), deIter.NextBack())
	assert.Equal(t, gust.Some(1), deIter.NextBack())
	assert.Equal(t, gust.None[int](), deIter.NextBack())
}

func TestSliceIteratorAdvanceBackBy(t *testing.T) {
	a := []int{3, 4, 5, 6}
	iter := FromSlice(a)
	deIter := iter.MustToDoubleEnded()

	assert.Equal(t, gust.NonErrable[uint](), deIter.AdvanceBackBy(2))
	assert.Equal(t, gust.Some(4), deIter.NextBack())
	assert.Equal(t, gust.NonErrable[uint](), deIter.AdvanceBackBy(0))
	assert.Equal(t, gust.ToErrable[uint](99), deIter.AdvanceBackBy(100))
}

func TestSliceIteratorNthBack(t *testing.T) {
	a := []int{1, 2, 3}
	iter := FromSlice(a)
	deIter := iter.MustToDoubleEnded()
	assert.Equal(t, gust.Some(1), deIter.NthBack(2))

	b := []int{1, 2, 3}
	iter2 := FromSlice(b)
	deIter2 := iter2.MustToDoubleEnded()
	assert.Equal(t, gust.Some(2), deIter2.NthBack(1))

	// NthBack(0) should return the last element (3)
	c := []int{1, 2, 3}
	iter3 := FromSlice(c)
	deIter3 := iter3.MustToDoubleEnded()
	assert.Equal(t, gust.Some(3), deIter3.NthBack(0))

	d := []int{1, 2, 3}
	iter4 := FromSlice(d)
	deIter4 := iter4.MustToDoubleEnded()
	assert.Equal(t, gust.None[int](), deIter4.NthBack(10))
}

func TestRfold(t *testing.T) {
	a := []int{1, 2, 3}
	iter := FromSlice(a)
	deIter := iter.MustToDoubleEnded()

	// Test XRfold (any version)
	sum := deIter.XRfold(0, func(acc any, x int) any { return acc.(int) + x })
	assert.Equal(t, 6, sum.(int))

	// Test Rfold (T version)
	iter2 := FromSlice(a)
	deIter2 := iter2.MustToDoubleEnded()
	sum2 := deIter2.Rfold(0, func(acc int, x int) int { return acc + x })
	assert.Equal(t, 6, sum2)

	// Test generic function version
	iter3 := FromSlice(a)
	deIter3 := iter3.MustToDoubleEnded()
	sum3 := Rfold(deIter3, 0, func(acc int, x int) int { return acc + x })
	assert.Equal(t, 6, sum3)
}

func TestTryRfold(t *testing.T) {
	a := []int{1, 2, 3}
	iter := FromSlice(a)
	deIter := iter.MustToDoubleEnded()

	// Test XTryRfold (any version)
	sum := deIter.XTryRfold(0, func(acc any, x int) gust.Result[any] {
		return gust.Ok(any(acc.(int) + x))
	})
	assert.True(t, sum.IsOk())
	assert.Equal(t, 6, sum.Unwrap().(int))

	// Test TryRfold (T version)
	iter2 := FromSlice(a)
	deIter2 := iter2.MustToDoubleEnded()
	sum2 := deIter2.TryRfold(0, func(acc int, x int) gust.Result[int] {
		return gust.Ok(acc + x)
	})
	assert.True(t, sum2.IsOk())
	assert.Equal(t, 6, sum2.Unwrap())

	// Test generic function version
	iter3 := FromSlice(a)
	deIter3 := iter3.MustToDoubleEnded()
	sum3 := TryRfold(deIter3, 0, func(acc int, x int) gust.Result[int] {
		return gust.Ok(acc + x)
	})
	assert.True(t, sum3.IsOk())
	assert.Equal(t, 6, sum3.Unwrap())
}

func TestRfind(t *testing.T) {
	a := []int{1, 2, 3}
	iter := FromSlice(a)
	deIter := iter.MustToDoubleEnded()
	assert.Equal(t, gust.Some(2), deIter.Rfind(func(x int) bool { return x == 2 }))

	b := []int{1, 2, 3}
	iter2 := FromSlice(b)
	deIter2 := iter2.MustToDoubleEnded()
	assert.Equal(t, gust.None[int](), deIter2.Rfind(func(x int) bool { return x == 5 }))

	// Test function version
	iter3 := FromSlice(a)
	deIter3 := iter3.MustToDoubleEnded()
	assert.Equal(t, gust.Some(2), Rfind(deIter3, func(x int) bool { return x == 2 }))
}

func TestDoubleEndedMixed(t *testing.T) {
	// Test mixing Next() and NextBack() calls
	numbers := []int{1, 2, 3, 4, 5}
	iter := FromSlice(numbers)
	deIter := iter.MustToDoubleEnded()

	// Start from front
	assert.Equal(t, gust.Some(1), deIter.Next())
	assert.Equal(t, gust.Some(2), deIter.Next())

	// Switch to back
	assert.Equal(t, gust.Some(5), deIter.NextBack())
	assert.Equal(t, gust.Some(4), deIter.NextBack())

	// Continue from front
	assert.Equal(t, gust.Some(3), deIter.Next())
	assert.Equal(t, gust.None[int](), deIter.Next())
	assert.Equal(t, gust.None[int](), deIter.NextBack())
}

func TestDoubleEndedEmpty(t *testing.T) {
	iter := FromSlice([]int{})
	deIter := iter.MustToDoubleEnded()
	assert.Equal(t, gust.None[int](), deIter.Next())
	assert.Equal(t, gust.None[int](), deIter.NextBack())
}

func TestDoubleEndedSingle(t *testing.T) {
	iter := FromSlice([]int{42})
	deIter := iter.MustToDoubleEnded()
	assert.Equal(t, gust.Some(42), deIter.Next())
	assert.Equal(t, gust.None[int](), deIter.Next())
	assert.Equal(t, gust.None[int](), deIter.NextBack())
}

func TestRemaining(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	iter := FromSlice(numbers)
	deIter := iter.MustToDoubleEnded()

	// Remaining is accessed through the underlying iter
	// We can test it indirectly by checking NextBack behavior
	assert.Equal(t, gust.Some(5), deIter.NextBack())
	assert.Equal(t, gust.Some(1), deIter.Next())
	assert.Equal(t, gust.Some(4), deIter.NextBack())
	assert.Equal(t, gust.Some(2), deIter.Next())
	assert.Equal(t, gust.Some(3), deIter.Next())
	assert.Equal(t, gust.None[int](), deIter.Next())
	assert.Equal(t, gust.None[int](), deIter.NextBack())
}

func TestDoubleEndedInheritsIteratorMethods(t *testing.T) {
	// Test that DoubleEndedIterator can use all Iterator methods
	numbers := []int{1, 2, 3, 4, 5, 6}
	iter := FromSlice(numbers)
	deIter := iter.MustToDoubleEnded()

	// Test Filter (Iterator method)
	filtered := deIter.Filter(func(x int) bool { return x > 3 })
	assert.Equal(t, gust.Some(4), filtered.Next())
	assert.Equal(t, gust.Some(5), filtered.Next())
	assert.Equal(t, gust.Some(6), filtered.Next())
	assert.Equal(t, gust.None[int](), filtered.Next())

	// Test Skip (Iterator method)
	iter2 := FromSlice(numbers)
	deIter2 := iter2.MustToDoubleEnded()
	skipped := deIter2.Skip(2)
	assert.Equal(t, gust.Some(3), skipped.Next())

	// Test Take (Iterator method)
	iter3 := FromSlice(numbers)
	deIter3 := iter3.MustToDoubleEnded()
	taken := deIter3.Take(2)
	assert.Equal(t, gust.Some(1), taken.Next())
	assert.Equal(t, gust.Some(2), taken.Next())
	assert.Equal(t, gust.None[int](), taken.Next())

	// Test Collect (Iterator method)
	iter4 := FromSlice(numbers)
	deIter4 := iter4.MustToDoubleEnded()
	collected := deIter4.Collect()
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, collected)

	// Test Count (Iterator method)
	iter5 := FromSlice(numbers)
	deIter5 := iter5.MustToDoubleEnded()
	count := deIter5.Count()
	assert.Equal(t, uint(6), count)

	// Test Chain (Iterator method)
	iter6 := FromSlice([]int{1, 2})
	iter7 := FromSlice([]int{3, 4})
	deIter6 := iter6.MustToDoubleEnded()
	chained := deIter6.Chain(iter7)
	assert.Equal(t, gust.Some(1), chained.Next())
	assert.Equal(t, gust.Some(2), chained.Next())
	assert.Equal(t, gust.Some(3), chained.Next())
	assert.Equal(t, gust.Some(4), chained.Next())
	assert.Equal(t, gust.None[int](), chained.Next())

	// Test that we can use Iterator methods and then convert back to DoubleEndedIterator
	iter8 := FromSlice([]int{1, 2, 3, 4, 5})
	deIter8 := iter8.MustToDoubleEnded()
	// Use Iterator method first
	filtered2 := deIter8.Filter(func(x int) bool { return x%2 == 0 })
	// Filter returns Iterator[T], not DoubleEndedIterator[T]
	// So we can't call NextBack() on it
	assert.Equal(t, gust.Some(2), filtered2.Next())
	assert.Equal(t, gust.Some(4), filtered2.Next())
	assert.Equal(t, gust.None[int](), filtered2.Next())
}
