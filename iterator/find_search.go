package iterator

import (
	"github.com/andeya/gust/option"
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
