package iter

import (
	"github.com/andeya/gust"
)

//go:inline
func allImpl[T any](iter Iterable[T], f func(T) bool) bool {
	for {
		item := iter.Next()
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
	for {
		item := iter.Next()
		if item.IsNone() {
			return false
		}
		if f(item.Unwrap()) {
			return true
		}
	}
}

//go:inline
func findImpl[T any](iter Iterable[T], predicate func(T) bool) gust.Option[T] {
	for {
		item := iter.Next()
		if item.IsNone() {
			return gust.None[T]()
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
//	var firstNumber = FindMap(FromSlice(a), func(s string) gust.Option[int] {
//		if v, err := strconv.Atoi(s); err == nil {
//			return gust.Some(v)
//		}
//		return gust.None[int]()
//	})
//	assert.Equal(t, gust.Some(2), firstNumber)
func FindMap[T any, U any](iter Iterator[T], f func(T) gust.Option[U]) gust.Option[U] {
	return findMapImpl(iter.iterable, f)
}

//go:inline
func findMapImpl[T any, U any](iter Iterable[T], f func(T) gust.Option[U]) gust.Option[U] {
	for {
		item := iter.Next()
		if item.IsNone() {
			return gust.None[U]()
		}
		if result := f(item.Unwrap()); result.IsSome() {
			return result
		}
	}
}

//go:inline
func positionImpl[T any](iter Iterable[T], predicate func(T) bool) gust.Option[uint] {
	var index uint
	for {
		item := iter.Next()
		if item.IsNone() {
			return gust.None[uint]()
		}
		if predicate(item.Unwrap()) {
			return gust.Some(index)
		}
		index++
	}
}
