package iterator

import (
	"github.com/andeya/gust/option"
	"github.com/andeya/gust/result"
)

//go:inline
func allImpl[T any](iter Iterable[T], f func(T) bool) bool {
	it := Iterator[T]{iterable: iter}
	for {
		item := it.Next()
		if item.IsNone() {
			return true
		}
		if !f(item.Unwrap()) {
			return false
		}
	}
}

//go:inline
func anyImpl[T any](iter Iterable[T], f func(T) bool) bool {
	it := Iterator[T]{iterable: iter}
	for {
		item := it.Next()
		if item.IsNone() {
			return false
		}
		if f(item.Unwrap()) {
			return true
		}
	}
}

//go:inline
func findImpl[T any](iter Iterable[T], predicate func(T) bool) option.Option[T] {
	it := Iterator[T]{iterable: iter}
	for {
		item := it.Next()
		if item.IsNone() {
			return option.None[T]()
		}
		if predicate(item.Unwrap()) {
			return item
		}
	}
}

// FindMap applies function to the elements of iterator and returns
// the first non-none result.
//
// FindMap(f) is equivalent to FilterMap(iter, f).Next().
//
// # Examples
//
//	var a = []string{"lol", "NaN", "2", "5"}
//	var firstNumber = FindMap(FromSlice(a), func(s string) option.Option[int] {
//		if v, err := strconv.Atoi(s); err == nil {
//			return option.Some(v)
//		}
//		return option.None[int]()
//	})
//	assert.Equal(t, option.Some(2), firstNumber)
//
//go:inline
func FindMap[T any, U any](iter Iterator[T], f func(T) option.Option[U]) option.Option[U] {
	return findMapImpl(iter.iterable, f)
}

//go:inline
func findMapImpl[T any, U any](iter Iterable[T], f func(T) option.Option[U]) option.Option[U] {
	it := Iterator[T]{iterable: iter}
	for {
		item := it.Next()
		if item.IsNone() {
			return option.None[U]()
		}
		if result := f(item.Unwrap()); result.IsSome() {
			return result
		}
	}
}

//go:inline
func positionImpl[T any](iter Iterable[T], predicate func(T) bool) option.Option[uint] {
	var index uint
	it := Iterator[T]{iterable: iter}
	for {
		item := it.Next()
		if item.IsNone() {
			return option.None[uint]()
		}
		if predicate(item.Unwrap()) {
			return option.Some(index)
		}
		index++
	}
}

//go:inline
func tryFindImpl[T any](iter Iterator[T], f func(T) result.Result[bool]) result.Result[option.Option[T]] {
	for {
		item := iter.Next()
		if item.IsNone() {
			return result.Ok(option.None[T]())
		}
		res := f(item.Unwrap())
		if res.IsErr() {
			return result.TryErr[option.Option[T]](res.UnwrapErr())
		}
		if res.Unwrap() {
			return result.Ok(item)
		}
	}
}

// All tests if all elements satisfy a predicate.
//
// # Examples
//
//	var iter = FromSlice([]int{2, 4, 6})
//	assert.True(t, iterator.All(func(x int) bool { return x%2 == 0 }))
//
//go:inline
func (it Iterator[T]) All(predicate func(T) bool) bool {
	return allImpl(it.iterable, predicate)
}

// Any tests if any element satisfies a predicate.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	assert.True(t, iterator.Any(func(x int) bool { return x > 2 }))
//
//go:inline
func (it Iterator[T]) Any(predicate func(T) bool) bool {
	return anyImpl(it.iterable, predicate)
}

// Find searches for an element that satisfies a predicate.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	assert.Equal(t, option.Some(2), iterator.Find(func(x int) bool { return x > 1 }))
//
//go:inline
func (it Iterator[T]) Find(predicate func(T) bool) option.Option[T] {
	return findImpl(it.iterable, predicate)
}

// XFindMap searches for an element and maps it (any version).
//
// # Examples
//
//	var iter = FromSlice([]string{"lol", "NaN", "2", "5"})
//	var firstNumber = iterator.XFindMap(func(s string) option.Option[any] {
//		if v, err := strconv.Atoi(s); err == nil {
//			return option.Some(any(v))
//		}
//		return option.None[any]()
//	})
//	assert.True(t, firstNumber.IsSome())
//	assert.Equal(t, 2, firstNumber.Unwrap().(int))
//
//go:inline
func (it Iterator[T]) XFindMap(f func(T) option.Option[any]) option.Option[any] {
	return FindMap(it, f)
}

// FindMap searches for an element and maps it.
//
// # Examples
//
//	var iter = FromSlice([]string{"lol", "NaN", "2", "5"})
//	var firstNumber = iterator.FindMap(func(s string) option.Option[int] {
//		if v, err := strconv.Atoi(s); err == nil {
//			return option.Some(v)
//		}
//		return option.None[int]()
//	})
//	assert.Equal(t, option.Some(2), firstNumber)
//
//go:inline
func (it Iterator[T]) FindMap(f func(T) option.Option[T]) option.Option[T] {
	return FindMap(it, f)
}

// Position searches for an element in an iterator, returning its index.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	assert.Equal(t, option.Some(uint(1)), iterator.Position(func(x int) bool { return x == 2 }))
//
//go:inline
func (it Iterator[T]) Position(predicate func(T) bool) option.Option[uint] {
	return positionImpl(it.iterable, predicate)
}

// TryFind applies function to the elements of iterator and returns
// the first true result or the first error.
//
// This is the fallible version of Find().
//
// # Examples
//
//	var a = []string{"1", "2", "lol", "NaN", "5"}
//	var res = iterator.TryFind(func(s string) result.Result[bool] {
//		if s == "lol" {
//			return result.TryErr[bool](errors.New("invalid"))
//		}
//		if v, err := strconv.Atoi(s); err == nil {
//			return result.Ok(v == 2)
//		}
//		return result.Ok(false)
//	})
//	assert.True(t, result.IsOk())
//
//go:inline
func (it Iterator[T]) TryFind(f func(T) result.Result[bool]) result.Result[option.Option[T]] {
	return tryFindImpl(it, f)
}

// Rfind searches for an element of an iterator from the back that satisfies a predicate.
//
// Rfind() takes a closure that returns true or false. It applies
// this closure to each element of the iterator, starting at the end, and if any
// of them return true, then Rfind() returns Some(element). If they all return
// false, it returns None.
//
// Rfind() is short-circuiting; in other words, it will stop processing
// as soon as the closure returns true.
//
// # Examples
//
//	var a = []int{1, 2, 3}
//	var deIter = FromSlice(a).MustToDoubleEnded()
//	assert.Equal(t, option.Some(2), deIter.Rfind(func(x int) bool { return x == 2 }))
//
//	var b = []int{1, 2, 3}
//	var deIter2 = FromSlice(b).MustToDoubleEnded()
//	assert.Equal(t, option.None[int](), deIter2.Rfind(func(x int) bool { return x == 5 }))
func (de DoubleEndedIterator[T]) Rfind(predicate func(T) bool) option.Option[T] {
	return rfindImpl(de.iterable, predicate)
}

// Rfind searches for an element of an iterator from the back that satisfies a predicate (function version).
//
// # Examples
//
//	var a = []int{1, 2, 3}
//	var deIter = FromSlice(a).MustToDoubleEnded()
//	assert.Equal(t, option.Some(2), Rfind(deIter, func(x int) bool { return x == 2 }))
func Rfind[T any](de DoubleEndedIterator[T], predicate func(T) bool) option.Option[T] {
	return de.Rfind(predicate)
}

// rfindImpl is the internal implementation of Rfind.
//
//go:inline
func rfindImpl[T any](iter DoubleEndedIterable[T], predicate func(T) bool) option.Option[T] {
	de := DoubleEndedIterator[T]{iterable: iter}
	for {
		item := de.NextBack()
		if item.IsNone() {
			return option.None[T]()
		}
		if predicate(item.Unwrap()) {
			return item
		}
	}
}
