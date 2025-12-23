package iter

import (
	"github.com/andeya/gust"
)

//go:inline
func countImpl[T any](iter Iterable[T]) uint {
	var count uint
	for iter.Next().IsSome() {
		count++
	}
	return count
}

//go:inline
func lastImpl[T any](iter Iterable[T]) gust.Option[T] {
	var last gust.Option[T] = gust.None[T]()
	for {
		item := iter.Next()
		if item.IsNone() {
			break
		}
		last = item
	}
	return last
}

//go:inline
func advanceByImpl[T any](iter Iterable[T], n uint) gust.Errable[uint] {
	for i := uint(0); i < n; i++ {
		if iter.Next().IsNone() {
			return gust.ToErrable[uint](n - i)
		}
	}
	return gust.NonErrable[uint]()
}

//go:inline
func nthImpl[T any](iter Iterable[T], n uint) gust.Option[T] {
	if advanceByImpl(iter, n).IsErr() {
		return gust.None[T]()
	}
	return iter.Next()
}

//go:inline
func nextChunkImpl[T any](iter Iterable[T], n uint) gust.EnumResult[[]T, []T] {
	if n == 0 {
		return gust.EnumOk[[]T, []T]([]T{})
	}
	result := make([]T, 0, n)
	for i := uint(0); i < n; i++ {
		item := iter.Next()
		if item.IsNone() {
			// Return error with remaining elements
			return gust.EnumErr[[]T, []T](result)
		}
		result = append(result, item.Unwrap())
	}
	return gust.EnumOk[[]T, []T](result)
}
