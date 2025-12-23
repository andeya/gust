package iter

import (
	"github.com/andeya/gust"
)

// NextBack removes and returns an element from the end of the iterator.
//
// Returns None when there are no more elements.
//
// # Examples
//
//	var numbers = []int{1, 2, 3, 4, 5, 6}
//	var iter = FromSlice(numbers)
//	var deIter = AsDoubleEnded(iter)
//
//	assert.Equal(t, gust.Some(1), deIter.Next())
//	assert.Equal(t, gust.Some(6), deIter.NextBack())
//	assert.Equal(t, gust.Some(5), deIter.NextBack())
//
//go:inline
func (de DoubleEndedIterator[T]) NextBack() gust.Option[T] {
	return de.iter.NextBack()
}

// AdvanceBackBy advances the iterator from the back by n elements.
//
// AdvanceBackBy is the reverse version of AdvanceBy. This method will
// eagerly skip n elements starting from the back by calling NextBack up
// to n times until None is encountered.
//
// AdvanceBackBy(n) will return NonErrable if the iterator successfully advances by
// n elements, or a ToErrable(k) with value k if None is encountered, where k
// is remaining number of steps that could not be advanced because the iterator ran out.
// If iter is empty and n is non-zero, then this returns ToErrable(n).
// Otherwise, k is always less than n.
//
// Calling AdvanceBackBy(0) can do meaningful work.
//
// # Examples
//
//	var a = []int{3, 4, 5, 6}
//	var iter = FromSlice(a)
//	var deIter = AsDoubleEnded(iter)
//
//	assert.Equal(t, gust.NonErrable[uint](), deIter.AdvanceBackBy(2))
//	assert.Equal(t, gust.Some(4), deIter.NextBack())
//	assert.Equal(t, gust.NonErrable[uint](), deIter.AdvanceBackBy(0))
//	assert.Equal(t, gust.ToErrable[uint](99), deIter.AdvanceBackBy(100))
func (de DoubleEndedIterator[T]) AdvanceBackBy(n uint) gust.Errable[uint] {
	for i := uint(0); i < n; i++ {
		if de.iter.NextBack().IsNone() {
			return gust.ToErrable[uint](n - i)
		}
	}
	return gust.NonErrable[uint]()
}

// NthBack returns the nth element from the end of the iterator.
//
// This is essentially the reversed version of Nth().
// Although like most indexing operations, the count starts from zero, so
// NthBack(0) returns the first value from the end, NthBack(1) the
// second, and so on.
//
// Note that all elements between the end and the returned element will be
// consumed, including the returned element. This also means that calling
// NthBack(0) multiple times on the same iterator will return different
// elements.
//
// NthBack() will return None if n is greater than or equal to the length of the
// iterator.
//
// # Examples
//
//	var a = []int{1, 2, 3}
//	var iter = FromSlice(a)
//	var deIter = AsDoubleEnded(iter)
//	assert.Equal(t, gust.Some(1), deIter.NthBack(2))
//
//	var b = []int{1, 2, 3}
//	var iter2 = FromSlice(b)
//	var deIter2 = AsDoubleEnded(iter2)
//	assert.Equal(t, gust.Some(2), deIter2.NthBack(1))
//	assert.Equal(t, gust.Some(3), deIter2.NthBack(0))
//
//	var c = []int{1, 2, 3}
//	var iter3 = FromSlice(c)
//	var deIter3 = AsDoubleEnded(iter3)
//	assert.Equal(t, gust.None[int](), deIter3.NthBack(10))
func (de DoubleEndedIterator[T]) NthBack(n uint) gust.Option[T] {
	if de.AdvanceBackBy(n).IsErr() {
		return gust.None[T]()
	}
	return de.iter.NextBack()
}
