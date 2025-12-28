package iterator_test

import (
	"iter"
	"testing"

	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/option"
	"github.com/andeya/gust/pair"
	"github.com/stretchr/testify/assert"
)

func TestFromIterable(t *testing.T) {
	// Test with Iterator[T] - should return the same iterator
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	var iterIter1 iterator.Iterable[int] = iter1
	iter2 := iterator.FromIterable(iterIter1)
	assert.Equal(t, option.Some(1), iter2.Next())
	assert.Equal(t, option.Some(2), iter2.Next())
	assert.Equal(t, option.Some(3), iter2.Next())
	assert.Equal(t, option.None[int](), iter2.Next())

	// Test with iterator.Iterable[T] that is not Iterator[T]
	custom := &easyIterable{values: []int{10, 20, 30}, index: 0}
	var iterIter2 iterator.Iterable[int] = custom
	iter3 := iterator.FromIterable(iterIter2)
	assert.Equal(t, option.Some(10), iter3.Next())
	assert.Equal(t, option.Some(20), iter3.Next())
	assert.Equal(t, option.Some(30), iter3.Next())
	assert.Equal(t, option.None[int](), iter3.Next())
}

func TestTryToDoubleEnded(t *testing.T) {
	// Test with double-ended iterator
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	deOpt := iter1.TryToDoubleEnded()
	assert.True(t, deOpt.IsSome())
	deIter := deOpt.Unwrap()
	assert.Equal(t, option.Some(3), deIter.NextBack())

	// Test with non-double-ended iterator (would need a custom iterator)
	// For now, sliceIterator supports double-ended, so this will succeed
}

func TestMustToDoubleEnded(t *testing.T) {
	// Test with double-ended iterator
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	deIter := iter1.MustToDoubleEnded()
	assert.Equal(t, option.Some(3), deIter.NextBack())
	assert.Equal(t, option.Some(2), deIter.NextBack())
	assert.Equal(t, option.Some(1), deIter.NextBack())
	assert.Equal(t, option.None[int](), deIter.NextBack())
}

func TestSeq(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	var results []int
	for v := range iter.Seq() {
		results = append(results, v)
	}
	assert.Equal(t, []int{1, 2, 3}, results)

	// Test early termination
	iter2 := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	var results2 []int
	for v := range iter2.Seq() {
		results2 = append(results2, v)
		if v == 3 {
			break
		}
	}
	assert.Equal(t, []int{1, 2, 3}, results2)
}

func TestPull(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	next, stop := iter.Pull()
	defer stop()

	var results []int
	for {
		v, ok := next()
		if !ok {
			break
		}
		results = append(results, v)
		if v == 3 {
			break // Early termination
		}
	}
	assert.Equal(t, []int{1, 2, 3}, results)
}

func TestSizeHint(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	lower, upper := iter.SizeHint()
	assert.Equal(t, uint(3), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(3), upper.Unwrap())

	iter2 := iterator.FromSlice([]int{})
	lower2, upper2 := iter2.SizeHint()
	assert.Equal(t, uint(0), lower2)
	assert.True(t, upper2.IsSome())
	assert.Equal(t, uint(0), upper2.Unwrap())
}

// TestMustToDoubleEndedPanic tests MustToDoubleEnded panic case
func TestMustToDoubleEndedPanic(t *testing.T) {
	// Create a non-double-ended iterator
	nonDE := &nonDEIterable{values: []int{1, 2, 3}, index: 0}
	var iterable iterator.Iterable[int] = nonDE
	iter := iterator.FromIterable(iterable)

	assert.Panics(t, func() {
		iter.MustToDoubleEnded()
	})
}

// TestTryToDoubleEndedNone tests TryToDoubleEnded returning None
func TestTryToDoubleEndedNone(t *testing.T) {
	// Create a non-double-ended iterator
	nonDE := &nonDEIterable{values: []int{1, 2, 3}, index: 0}
	var iterable iterator.Iterable[int] = nonDE
	iter := iterator.FromIterable(iterable)

	result := iter.TryToDoubleEnded()
	assert.True(t, result.IsNone())
}

type testIterable struct {
	values []int
	index  int
}

func (t *testIterable) Next() option.Option[int] {
	if t.index >= len(t.values) {
		return option.None[int]()
	}
	val := t.values[t.index]
	t.index++
	return option.Some(val)
}

func (t *testIterable) SizeHint() (uint, option.Option[uint]) {
	return iterator.DefaultSizeHint[int]()
}

// TestFromIterableWithIterable tests iterator.FromIterable with iterator.Iterable[T] (not iterator.Iterator[T])
func TestFromIterableWithIterable(t *testing.T) {
	ti := &testIterable{values: []int{10, 20, 30}, index: 0}
	var iterable iterator.Iterable[int] = ti
	var gustIter iterator.Iterable[int] = iterable
	iter := iterator.FromIterable(gustIter)

	assert.Equal(t, option.Some(10), iter.Next())
	assert.Equal(t, option.Some(20), iter.Next())
	assert.Equal(t, option.Some(30), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

type testIterable2 struct {
	values []int
	index  int
}

func (t *testIterable2) Next() option.Option[int] {
	if t.index >= len(t.values) {
		return option.None[int]()
	}
	val := t.values[t.index]
	t.index++
	return option.Some(val)
}

func (t *testIterable2) SizeHint() (uint, option.Option[uint]) {
	return iterator.DefaultSizeHint[int]()
}

// TestFromIterableIterableBranch tests iterator.FromIterable with iterator.Iterable[T] branch
func TestFromIterableIterableBranch(t *testing.T) {
	ti := &testIterable2{values: []int{10, 20, 30}, index: 0}
	var iterable iterator.Iterable[int] = ti
	var gustIter iterator.Iterable[int] = iterable
	iter := iterator.FromIterable(gustIter)

	assert.Equal(t, option.Some(10), iter.Next())
	assert.Equal(t, option.Some(20), iter.Next())
	assert.Equal(t, option.Some(30), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

// TestIterableWrapperSizeHint tests iterableWrapper SizeHint
func TestIterableWrapperSizeHint(t *testing.T) {
	custom := &easyIterable{values: []int{1, 2, 3}, index: 0}
	var gustIter iterator.Iterable[int] = custom
	iter := iterator.FromIterable(gustIter)
	lower, upper := iter.SizeHint()
	// easyIterable implements SizeHint() which returns the actual remaining size
	assert.Equal(t, uint(3), lower)
	assert.True(t, upper.IsSome())
	assert.Equal(t, uint(3), upper.Unwrap())
}

// TestDefaultSizeHint tests iterator.DefaultSizeHint function
func TestDefaultSizeHint(t *testing.T) {
	lower, upper := iterator.DefaultSizeHint[int]()
	assert.Equal(t, uint(0), lower)
	assert.True(t, upper.IsNone())
}

// TestMustToDoubleEnded_Panic tests MustToDoubleEnded with non-double-ended iterator (should panic)
func TestMustToDoubleEnded_Panic(t *testing.T) {
	// Test MustToDoubleEnded with non-double-ended iterator (should panic)
	// Create a non-double-ended iterator using FromIterable
	nonDE := &nonDoubleEndedIterable{values: []int{1, 2, 3}, index: 0}
	var iterable iterator.Iterable[int] = nonDE
	var gustIter iterator.Iterable[int] = iterable
	iter := iterator.FromIterable(gustIter)

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustToDoubleEnded should panic for non-double-ended iterator")
		}
	}()

	_ = iter.MustToDoubleEnded()
}

type nonDoubleEndedIterable struct {
	values []int
	index  int
}

func (n *nonDoubleEndedIterable) Next() option.Option[int] {
	if n.index >= len(n.values) {
		return option.None[int]()
	}
	val := n.values[n.index]
	n.index++
	return option.Some(val)
}

func (n *nonDoubleEndedIterable) SizeHint() (uint, option.Option[uint]) {
	return 0, option.None[uint]()
}

type customIterable struct {
	values []int
	index  int
}

func (c *customIterable) Next() option.Option[int] {
	if c.index >= len(c.values) {
		return option.None[int]()
	}
	val := c.values[c.index]
	c.index++
	return option.Some(val)
}

func (c *customIterable) SizeHint() (uint, option.Option[uint]) {
	return 0, option.None[uint]()
}

// TestFromIterable_IterablePath tests FromIterable with Iterable[T] path (not Iterator[T])
func TestFromIterable_IterablePath(t *testing.T) {
	// Test FromIterable with Iterable[T] path (not Iterator[T])
	custom := &customIterable{values: []int{10, 20, 30}, index: 0}
	var iterable iterator.Iterable[int] = custom
	var gustIter iterator.Iterable[int] = iterable
	iter := iterator.FromIterable(gustIter)
	assert.Equal(t, option.Some(10), iter.Next())
	assert.Equal(t, option.Some(20), iter.Next())
	assert.Equal(t, option.Some(30), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}

func TestSliceIteratorDoubleEnded(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5, 6}
	iter := iterator.FromSlice(numbers)
	deIter := iter.MustToDoubleEnded()

	assert.Equal(t, option.Some(1), deIter.Next())
	assert.Equal(t, option.Some(6), deIter.NextBack())
	assert.Equal(t, option.Some(5), deIter.NextBack())
	assert.Equal(t, option.Some(2), deIter.Next())
	assert.Equal(t, option.Some(3), deIter.Next())
	assert.Equal(t, option.Some(4), deIter.Next())
	assert.Equal(t, option.None[int](), deIter.Next())
	assert.Equal(t, option.None[int](), deIter.NextBack())
}

func TestSliceIteratorRemaining(t *testing.T) {
	var numbers = []int{1, 2, 3, 4, 5, 6}
	var deIter = iterator.FromSlice(numbers).MustToDoubleEnded()
	assert.Equal(t, uint(6), deIter.Remaining())
	deIter.Next()
	assert.Equal(t, uint(5), deIter.Remaining())
	deIter.NextBack()
	assert.Equal(t, uint(4), deIter.Remaining())
	deIter.NextBack()
	assert.Equal(t, uint(3), deIter.Remaining())
	deIter.NextBack()
	assert.Equal(t, uint(2), deIter.Remaining())
	deIter.NextBack()
	assert.Equal(t, uint(1), deIter.Remaining())
	deIter.NextBack()
	assert.Equal(t, uint(0), deIter.Remaining())
	deIter.NextBack()
	assert.Equal(t, uint(0), deIter.Remaining())
}

func TestSliceIteratorNextBack(t *testing.T) {
	numbers := []int{1, 2, 3}
	iter := iterator.FromSlice(numbers)
	deIter := iter.MustToDoubleEnded()

	assert.Equal(t, option.Some(3), deIter.NextBack())
	assert.Equal(t, option.Some(2), deIter.NextBack())
	assert.Equal(t, option.Some(1), deIter.NextBack())
	assert.Equal(t, option.None[int](), deIter.NextBack())
}

func TestSliceIteratorAdvanceBackBy(t *testing.T) {
	a := []int{3, 4, 5, 6}
	iter := iterator.FromSlice(a)
	deIter := iter.MustToDoubleEnded()

	assert.True(t, deIter.AdvanceBackBy(2).IsOk())
	assert.Equal(t, option.Some(4), deIter.NextBack())
	assert.True(t, deIter.AdvanceBackBy(0).IsOk())
	result := deIter.AdvanceBackBy(100)
	assert.True(t, result.IsErr())
	assert.Equal(t, uint(99), result.ErrVal())
}

func TestSliceIteratorNthBack(t *testing.T) {
	a := []int{1, 2, 3}
	iter := iterator.FromSlice(a)
	deIter := iter.MustToDoubleEnded()
	assert.Equal(t, option.Some(1), deIter.NthBack(2))

	b := []int{1, 2, 3}
	iter2 := iterator.FromSlice(b)
	deIter2 := iter2.MustToDoubleEnded()
	assert.Equal(t, option.Some(2), deIter2.NthBack(1))

	// NthBack(0) should return the last element (3)
	c := []int{1, 2, 3}
	iter3 := iterator.FromSlice(c)
	deIter3 := iter3.MustToDoubleEnded()
	assert.Equal(t, option.Some(3), deIter3.NthBack(0))

	d := []int{1, 2, 3}
	iter4 := iterator.FromSlice(d)
	deIter4 := iter4.MustToDoubleEnded()
	assert.Equal(t, option.None[int](), deIter4.NthBack(10))
}

func TestDoubleEndedMixed(t *testing.T) {
	// Test mixing Next() and NextBack() calls
	numbers := []int{1, 2, 3, 4, 5}
	iter := iterator.FromSlice(numbers)
	deIter := iter.MustToDoubleEnded()

	// Start from front
	assert.Equal(t, option.Some(1), deIter.Next())
	assert.Equal(t, option.Some(2), deIter.Next())

	// Switch to back
	assert.Equal(t, option.Some(5), deIter.NextBack())
	assert.Equal(t, option.Some(4), deIter.NextBack())

	// Continue from front
	assert.Equal(t, option.Some(3), deIter.Next())
	assert.Equal(t, option.None[int](), deIter.Next())
	assert.Equal(t, option.None[int](), deIter.NextBack())
}

func TestDoubleEndedEmpty(t *testing.T) {
	iter := iterator.FromSlice([]int{})
	deIter := iter.MustToDoubleEnded()
	assert.Equal(t, option.None[int](), deIter.Next())
	assert.Equal(t, option.None[int](), deIter.NextBack())
}

func TestDoubleEndedSingle(t *testing.T) {
	iter := iterator.FromSlice([]int{42})
	deIter := iter.MustToDoubleEnded()
	assert.Equal(t, option.Some(42), deIter.Next())
	assert.Equal(t, option.None[int](), deIter.Next())
	assert.Equal(t, option.None[int](), deIter.NextBack())
}

func TestRemaining(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	iter := iterator.FromSlice(numbers)
	deIter := iter.MustToDoubleEnded()

	// Remaining is accessed through the underlying iter
	// We can test it indirectly by checking NextBack behavior
	assert.Equal(t, option.Some(5), deIter.NextBack())
	assert.Equal(t, option.Some(1), deIter.Next())
	assert.Equal(t, option.Some(4), deIter.NextBack())
	assert.Equal(t, option.Some(2), deIter.Next())
	assert.Equal(t, option.Some(3), deIter.Next())
	assert.Equal(t, option.None[int](), deIter.Next())
	assert.Equal(t, option.None[int](), deIter.NextBack())
}

func TestDoubleEndedInheritsIteratorMethods(t *testing.T) {
	// Test that DoubleEndedIterator can use all Iterator methods
	numbers := []int{1, 2, 3, 4, 5, 6}
	iter := iterator.FromSlice(numbers)
	deIter := iter.MustToDoubleEnded()

	// Test Filter (Iterator method)
	filtered := deIter.Filter(func(x int) bool { return x > 3 })
	assert.Equal(t, option.Some(4), filtered.Next())
	assert.Equal(t, option.Some(5), filtered.Next())
	assert.Equal(t, option.Some(6), filtered.Next())
	assert.Equal(t, option.None[int](), filtered.Next())

	// Test Skip (Iterator method)
	iter2 := iterator.FromSlice(numbers)
	deIter2 := iter2.MustToDoubleEnded()
	skipped := deIter2.Skip(2)
	assert.Equal(t, option.Some(3), skipped.Next())

	// Test Take (Iterator method)
	iter3 := iterator.FromSlice(numbers)
	deIter3 := iter3.MustToDoubleEnded()
	taken := deIter3.Take(2)
	assert.Equal(t, option.Some(1), taken.Next())
	assert.Equal(t, option.Some(2), taken.Next())
	assert.Equal(t, option.None[int](), taken.Next())

	// Test Collect (Iterator method)
	iter4 := iterator.FromSlice(numbers)
	deIter4 := iter4.MustToDoubleEnded()
	collected := deIter4.Collect()
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, collected)

	// Test Count (Iterator method)
	iter5 := iterator.FromSlice(numbers)
	deIter5 := iter5.MustToDoubleEnded()
	count := deIter5.Count()
	assert.Equal(t, uint(6), count)

	// Test Chain (Iterator method)
	iter6 := iterator.FromSlice([]int{1, 2})
	iter7 := iterator.FromSlice([]int{3, 4})
	deIter6 := iter6.MustToDoubleEnded()
	chained := deIter6.Chain(iter7)
	assert.Equal(t, option.Some(1), chained.Next())
	assert.Equal(t, option.Some(2), chained.Next())
	assert.Equal(t, option.Some(3), chained.Next())
	assert.Equal(t, option.Some(4), chained.Next())
	assert.Equal(t, option.None[int](), chained.Next())

	// Test that we can use Iterator methods and then convert back to DoubleEndedIterator
	iter8 := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	deIter8 := iter8.MustToDoubleEnded()
	// Use Iterator method first
	filtered2 := deIter8.Filter(func(x int) bool { return x%2 == 0 })
	// Filter returns Iterator[T], not DoubleEndedIterator[T]
	// So we can't call NextBack() on it
	assert.Equal(t, option.Some(2), filtered2.Next())
	assert.Equal(t, option.Some(4), filtered2.Next())
	assert.Equal(t, option.None[int](), filtered2.Next())
}

// TestAdvanceBackByZero tests AdvanceBackBy with zero
func TestAdvanceBackByZero(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	deIter := iter.MustToDoubleEnded()
	result := deIter.AdvanceBackBy(0)
	assert.True(t, result.IsOk())
	assert.Equal(t, option.Some(3), deIter.NextBack())
}

// TestAdvanceBackByTooMany tests AdvanceBackBy with too many steps
func TestAdvanceBackByTooMany(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	deIter := iter.MustToDoubleEnded()
	result := deIter.AdvanceBackBy(100)
	assert.True(t, result.IsErr())
	assert.Equal(t, uint(97), result.ErrVal()) // 100 - 3 = 97
}

// TestNthBackEdgeCases tests NthBack edge cases
func TestNthBackEdgeCases(t *testing.T) {
	// Test with empty iterator
	iter := iterator.FromSlice([]int{})
	deIter := iter.MustToDoubleEnded()
	assert.Equal(t, option.None[int](), deIter.NthBack(0))

	// Test with single element
	iter2 := iterator.FromSlice([]int{42})
	deIter2 := iter2.MustToDoubleEnded()
	assert.Equal(t, option.Some(42), deIter2.NthBack(0))
	assert.Equal(t, option.None[int](), deIter2.NthBack(0))
}

// TestSliceIteratorNextBackEmpty tests sliceIterator NextBack with empty slice
func TestSliceIteratorNextBackEmpty(t *testing.T) {
	iter := iterator.FromSlice([]int{})
	deIter := iter.MustToDoubleEnded()
	assert.Equal(t, option.None[int](), deIter.NextBack())
}

// TestSliceIteratorNextBackSingle tests sliceIterator NextBack with single element
func TestSliceIteratorNextBackSingle(t *testing.T) {
	iter := iterator.FromSlice([]int{42})
	deIter := iter.MustToDoubleEnded()
	assert.Equal(t, option.Some(42), deIter.NextBack())
	assert.Equal(t, option.None[int](), deIter.NextBack())
}

// TestAdvanceBackByEmpty tests AdvanceBackBy with empty iterator
func TestAdvanceBackByEmpty(t *testing.T) {
	iter := iterator.FromSlice([]int{})
	deIter := iter.MustToDoubleEnded()
	result := deIter.AdvanceBackBy(5)
	assert.True(t, result.IsErr())
	assert.Equal(t, uint(5), result.ErrVal())
}

// TestNthBackEmpty tests NthBack with empty iterator
func TestNthBackEmpty(t *testing.T) {
	iter := iterator.FromSlice([]int{})
	deIter := iter.MustToDoubleEnded()
	assert.Equal(t, option.None[int](), deIter.NthBack(0))
}

// TestNthBackTooLarge tests NthBack with index too large
func TestNthBackTooLarge(t *testing.T) {
	iter := iterator.FromSlice([]int{1, 2, 3})
	deIter := iter.MustToDoubleEnded()
	assert.Equal(t, option.None[int](), deIter.NthBack(10))
}

// easyIterable is a helper type for testing Iterable interface
type easyIterable struct {
	values []int
	index  int
}

func (c *easyIterable) Next() option.Option[int] {
	if c.index >= len(c.values) {
		return option.None[int]()
	}
	val := c.values[c.index]
	c.index++
	return option.Some(val)
}

func (c *easyIterable) SizeHint() (uint, option.Option[uint]) {
	remaining := uint(len(c.values) - c.index)
	return remaining, option.Some(remaining)
}

// nonDEIterable is a helper type for testing non-double-ended iterators
type nonDEIterable struct {
	values []int
	index  int
}

func (n *nonDEIterable) Next() option.Option[int] {
	if n.index >= len(n.values) {
		return option.None[int]()
	}
	val := n.values[n.index]
	n.index++
	return option.Some(val)
}

func (n *nonDEIterable) SizeHint() (uint, option.Option[uint]) {
	// Use DefaultSizeHint to make it a non-double-ended iterator
	return iterator.DefaultSizeHint[int]()
}

func TestIterator_Seq(t *testing.T) {
	// Test basic conversion
	iter := iterator.FromSlice([]int{1, 2, 3})
	var result []int
	for v := range iter.Seq() {
		result = append(result, v)
	}
	assert.Equal(t, []int{1, 2, 3}, result)

	// Test with filtered iterator
	filtered := iterator.FromSlice([]int{1, 2, 3, 4, 5}).Filter(func(x int) bool { return x%2 == 0 })
	result = nil
	for v := range filtered.Seq() {
		result = append(result, v)
	}
	assert.Equal(t, []int{2, 4}, result)

	// Test with empty iterator
	empty := iterator.Empty[int]()
	result = nil
	for v := range empty.Seq() {
		result = append(result, v)
	}
	assert.Nil(t, result)
	assert.Len(t, result, 0)

	// Test early termination
	iter = iterator.FromSlice([]int{1, 2, 3, 4, 5})
	result = nil
	count := 0
	for v := range iter.Seq() {
		result = append(result, v)
		count++
		if count >= 3 {
			break
		}
	}
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestSeq_RoundTrip(t *testing.T) {
	// Test round trip: gust Iterator -> Seq -> gust Iterator
	// Create two independent iterators: one for seq, one for expected result
	originalForSeq := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	originalForExpected := iterator.FromSlice([]int{1, 2, 3, 4, 5})

	seq := originalForSeq.Seq()
	converted, deferStop := iterator.FromSeq(seq)
	defer deferStop()

	// Get expected result from independent iterator
	var expectedResult []int
	for {
		opt := originalForExpected.Next()
		if opt.IsNone() {
			break
		}
		expectedResult = append(expectedResult, opt.Unwrap())
	}

	// Get actual result from converted iterator
	var convertedResult []int
	for {
		opt := converted.Next()
		if opt.IsNone() {
			break
		}
		convertedResult = append(convertedResult, opt.Unwrap())
	}

	assert.Equal(t, expectedResult, convertedResult)
}

func TestSeq_WithMap(t *testing.T) {
	// Test Seq with mapped iterator
	iter := iterator.FromSlice([]int{1, 2, 3}).Map(func(x int) int { return x * 2 })
	var result []int
	for v := range iter.Seq() {
		result = append(result, v)
	}
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestSeq2(t *testing.T) {
	// Test with Zip iterator
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	iter2 := iterator.FromSlice([]string{"a", "b", "c"})
	zipped := iterator.Zip(iter1, iter2)

	var result []pair.Pair[int, string]
	for k, v := range iterator.Seq2(zipped) {
		result = append(result, pair.Pair[int, string]{A: k, B: v})
	}
	assert.Equal(t, []pair.Pair[int, string]{
		{A: 1, B: "a"},
		{A: 2, B: "b"},
		{A: 3, B: "c"},
	}, result)

	// Test with Enumerate iterator
	enumerated := iterator.Enumerate(iterator.FromSlice([]string{"x", "y", "z"}))
	var enumResult []pair.Pair[uint, string]
	for idx, val := range iterator.Seq2(enumerated) {
		enumResult = append(enumResult, pair.Pair[uint, string]{A: idx, B: val})
	}
	assert.Equal(t, []pair.Pair[uint, string]{
		{A: 0, B: "x"},
		{A: 1, B: "y"},
		{A: 2, B: "z"},
	}, enumResult)

	// Test with empty iterator
	empty := iterator.Zip(iterator.Empty[int](), iterator.Empty[string]())
	var emptyResult []pair.Pair[int, string]
	for k, v := range iterator.Seq2(empty) {
		emptyResult = append(emptyResult, pair.Pair[int, string]{A: k, B: v})
	}
	assert.Nil(t, emptyResult)
	assert.Len(t, emptyResult, 0)

	// Test early termination
	iter1 = iterator.FromSlice([]int{1, 2, 3, 4, 5})
	iter2 = iterator.FromSlice([]string{"a", "b", "c", "d", "e"})
	zipped = iterator.Zip(iter1, iter2)
	var earlyResult []pair.Pair[int, string]
	count := 0
	for k, v := range iterator.Seq2(zipped) {
		earlyResult = append(earlyResult, pair.Pair[int, string]{A: k, B: v})
		count++
		if count >= 2 {
			break
		}
	}
	assert.Equal(t, []pair.Pair[int, string]{
		{A: 1, B: "a"},
		{A: 2, B: "b"},
	}, earlyResult)
}

func TestSeq2_RoundTrip(t *testing.T) {
	// Test round trip: gust Iterator -> Seq2 -> gust Iterator
	// Create two independent iterators: one for seq2, one for expected result
	iter1ForSeq := iterator.FromSlice([]int{1, 2, 3})
	iter2ForSeq := iterator.FromSlice([]string{"a", "b", "c"})
	originalForSeq := iterator.Zip(iter1ForSeq, iter2ForSeq)

	iter1ForExpected := iterator.FromSlice([]int{1, 2, 3})
	iter2ForExpected := iterator.FromSlice([]string{"a", "b", "c"})
	originalForExpected := iterator.Zip(iter1ForExpected, iter2ForExpected)

	seq2 := iterator.Seq2(originalForSeq)
	converted, deferStop := iterator.FromSeq2(seq2)
	defer deferStop()

	// Get expected result from independent iterator
	var expectedResult []pair.Pair[int, string]
	for {
		opt := originalForExpected.Next()
		if opt.IsNone() {
			break
		}
		expectedResult = append(expectedResult, opt.Unwrap())
	}

	// Get actual result from converted iterator
	var convertedResult []pair.Pair[int, string]
	for {
		opt := converted.Next()
		if opt.IsNone() {
			break
		}
		convertedResult = append(convertedResult, opt.Unwrap())
	}

	assert.Equal(t, expectedResult, convertedResult)
}

func TestSeq2_WithGoStandardLibrary(t *testing.T) {
	// Test that Seq2 works with Go's standard library iterator.Pull2
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	iter2 := iterator.FromSlice([]string{"a", "b", "c"})
	zipped := iterator.Zip(iter1, iter2)
	seq2 := iterator.Seq2(zipped)

	// Use iter.Pull2 to pull values manually
	next, stop := iter.Pull2(seq2)
	defer stop()

	var result []pair.Pair[int, string]
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		result = append(result, pair.Pair[int, string]{A: k, B: v})
	}

	assert.Equal(t, []pair.Pair[int, string]{
		{A: 1, B: "a"},
		{A: 2, B: "b"},
		{A: 3, B: "c"},
	}, result)
}

func TestIterator_Pull(t *testing.T) {
	// Test basic Pull functionality
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	next, stop := iter.Pull()
	defer stop()

	var result []int
	for {
		v, ok := next()
		if !ok {
			break
		}
		result = append(result, v)
	}
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)

	// Test early termination
	iter = iterator.FromSlice([]int{1, 2, 3, 4, 5})
	next, stop = iter.Pull()
	defer stop()

	result = nil
	for {
		v, ok := next()
		if !ok {
			break
		}
		result = append(result, v)
		if v == 3 {
			break // Early termination
		}
	}
	assert.Equal(t, []int{1, 2, 3}, result)

	// Test with empty iterator
	empty := iterator.Empty[int]()
	next, stop = empty.Pull()
	defer stop()

	result = nil
	for {
		v, ok := next()
		if !ok {
			break
		}
		result = append(result, v)
	}
	assert.Nil(t, result)
	assert.Len(t, result, 0)

	// Test with filtered iterator
	filtered := iterator.FromSlice([]int{1, 2, 3, 4, 5}).Filter(func(x int) bool { return x%2 == 0 })
	next, stop = filtered.Pull()
	defer stop()

	result = nil
	for {
		v, ok := next()
		if !ok {
			break
		}
		result = append(result, v)
	}
	assert.Equal(t, []int{2, 4}, result)
}

func TestPull2(t *testing.T) {
	// Test basic Pull2 functionality
	iter1 := iterator.FromSlice([]int{1, 2, 3})
	iter2 := iterator.FromSlice([]string{"a", "b", "c"})
	zipped := iterator.Zip(iter1, iter2)

	next, stop := iterator.Pull2(zipped)
	defer stop()

	var result []pair.Pair[int, string]
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		result = append(result, pair.Pair[int, string]{A: k, B: v})
	}

	assert.Equal(t, []pair.Pair[int, string]{
		{A: 1, B: "a"},
		{A: 2, B: "b"},
		{A: 3, B: "c"},
	}, result)

	// Test early termination
	iter1 = iterator.FromSlice([]int{1, 2, 3, 4, 5})
	iter2 = iterator.FromSlice([]string{"a", "b", "c", "d", "e"})
	zipped = iterator.Zip(iter1, iter2)

	next, stop = iterator.Pull2(zipped)
	defer stop()

	result = nil
	count := 0
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		result = append(result, pair.Pair[int, string]{A: k, B: v})
		count++
		if count >= 2 {
			break // Early termination
		}
	}
	assert.Equal(t, []pair.Pair[int, string]{
		{A: 1, B: "a"},
		{A: 2, B: "b"},
	}, result)

	// Test with empty iterator
	empty1 := iterator.Empty[int]()
	empty2 := iterator.Empty[string]()
	zipped = iterator.Zip(empty1, empty2)

	next, stop = iterator.Pull2(zipped)
	defer stop()

	result = nil
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		result = append(result, pair.Pair[int, string]{A: k, B: v})
	}
	assert.Nil(t, result)
	assert.Len(t, result, 0)

	// Test with Enumerate
	enumerated := iterator.Enumerate(iterator.FromSlice([]string{"x", "y", "z"}))
	nextEnum, stopEnum := iterator.Pull2(enumerated)
	defer stopEnum()

	var enumResult []pair.Pair[uint, string]
	for {
		idx, val, ok := nextEnum()
		if !ok {
			break
		}
		enumResult = append(enumResult, pair.Pair[uint, string]{A: idx, B: val})
	}

	assert.Equal(t, []pair.Pair[uint, string]{
		{A: 0, B: "x"},
		{A: 1, B: "y"},
		{A: 2, B: "z"},
	}, enumResult)
}

func TestIterator_Seq2(t *testing.T) {
	// Test basic Seq2 conversion - converts Iterator[T] to iterator.Seq2[uint, T]
	iter := iterator.FromSlice([]int{1, 2, 3})
	var result []pair.Pair[uint, int]
	for k, v := range iter.Seq2() {
		result = append(result, pair.Pair[uint, int]{A: k, B: v})
	}
	assert.Equal(t, []pair.Pair[uint, int]{
		{A: 0, B: 1},
		{A: 1, B: 2},
		{A: 2, B: 3},
	}, result)

	// Test with filtered iterator
	filtered := iterator.FromSlice([]int{10, 20, 30, 40, 50}).Filter(func(x int) bool { return x > 20 })
	var filteredResult []pair.Pair[uint, int]
	for k, v := range filtered.Seq2() {
		filteredResult = append(filteredResult, pair.Pair[uint, int]{A: k, B: v})
	}
	assert.Equal(t, []pair.Pair[uint, int]{
		{A: 0, B: 30},
		{A: 1, B: 40},
		{A: 2, B: 50},
	}, filteredResult)

	// Test with empty iterator
	empty := iterator.Empty[int]()
	var emptyResult []pair.Pair[uint, int]
	for k, v := range empty.Seq2() {
		emptyResult = append(emptyResult, pair.Pair[uint, int]{A: k, B: v})
	}
	assert.Nil(t, emptyResult)
	assert.Len(t, emptyResult, 0)

	// Test early termination
	iter = iterator.FromSlice([]int{1, 2, 3, 4, 5})
	var earlyResult []pair.Pair[uint, int]
	count := 0
	for k, v := range iter.Seq2() {
		earlyResult = append(earlyResult, pair.Pair[uint, int]{A: k, B: v})
		count++
		if count >= 3 {
			break
		}
	}
	assert.Equal(t, []pair.Pair[uint, int]{
		{A: 0, B: 1},
		{A: 1, B: 2},
		{A: 2, B: 3},
	}, earlyResult)

	// Test with string iterator
	strIter := iterator.FromSlice([]string{"hello", "world", "rust"})
	var strResult []pair.Pair[uint, string]
	for k, v := range strIter.Seq2() {
		strResult = append(strResult, pair.Pair[uint, string]{A: k, B: v})
	}
	assert.Equal(t, []pair.Pair[uint, string]{
		{A: 0, B: "hello"},
		{A: 1, B: "world"},
		{A: 2, B: "rust"},
	}, strResult)
}

func TestIterator_Pull2(t *testing.T) {
	// Test basic Pull2 functionality - converts Iterator[T] to pull-style iterator with index-value pairs
	iter := iterator.FromSlice([]int{1, 2, 3, 4, 5})
	next, stop := iter.Pull2()
	defer stop()

	var result []pair.Pair[uint, int]
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		result = append(result, pair.Pair[uint, int]{A: k, B: v})
	}
	assert.Equal(t, []pair.Pair[uint, int]{
		{A: 0, B: 1},
		{A: 1, B: 2},
		{A: 2, B: 3},
		{A: 3, B: 4},
		{A: 4, B: 5},
	}, result)

	// Test early termination
	iter = iterator.FromSlice([]int{1, 2, 3, 4, 5})
	next, stop = iter.Pull2()
	defer stop()

	result = nil
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		result = append(result, pair.Pair[uint, int]{A: k, B: v})
		if v == 3 {
			break // Early termination
		}
	}
	assert.Equal(t, []pair.Pair[uint, int]{
		{A: 0, B: 1},
		{A: 1, B: 2},
		{A: 2, B: 3},
	}, result)

	// Test with empty iterator
	empty := iterator.Empty[int]()
	next, stop = empty.Pull2()
	defer stop()

	result = nil
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		result = append(result, pair.Pair[uint, int]{A: k, B: v})
	}
	assert.Nil(t, result)
	assert.Len(t, result, 0)

	// Test with filtered iterator
	filtered := iterator.FromSlice([]int{10, 20, 30, 40, 50}).Filter(func(x int) bool { return x%20 == 0 })
	next, stop = filtered.Pull2()
	defer stop()

	result = nil
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		result = append(result, pair.Pair[uint, int]{A: k, B: v})
	}
	assert.Equal(t, []pair.Pair[uint, int]{
		{A: 0, B: 20},
		{A: 1, B: 40},
	}, result)

	// Test with string iterator
	strIter := iterator.FromSlice([]string{"a", "b", "c"})
	nextStr, stopStr := strIter.Pull2()
	defer stopStr()

	var strResult []pair.Pair[uint, string]
	for {
		k, v, ok := nextStr()
		if !ok {
			break
		}
		strResult = append(strResult, pair.Pair[uint, string]{A: k, B: v})
	}
	assert.Equal(t, []pair.Pair[uint, string]{
		{A: 0, B: "a"},
		{A: 1, B: "b"},
		{A: 2, B: "c"},
	}, strResult)
}
