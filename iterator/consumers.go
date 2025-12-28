package iterator

import (
	"github.com/andeya/gust/option"
	"github.com/andeya/gust/result"
	"github.com/andeya/gust/void"
)

//go:inline
func countImpl[T any](iter Iterable[T]) uint {
	var count uint
	it := Iterator[T]{iterable: iter}
	for it.Next().IsSome() {
		count++
	}
	return count
}

//go:inline
func lastImpl[T any](iter Iterable[T]) option.Option[T] {
	var last option.Option[T] = option.None[T]()
	it := Iterator[T]{iterable: iter}
	for {
		item := it.Next()
		if item.IsNone() {
			break
		}
		last = item
	}
	return last
}

//go:inline
func advanceByImpl[T any](iter Iterable[T], n uint) result.VoidResult {
	it := Iterator[T]{iterable: iter}
	for i := uint(0); i < n; i++ {
		if it.Next().IsNone() {
			return result.TryErr[void.Void](n - i)
		}
	}
	return result.Ok[void.Void](nil)
}

//go:inline
func nthImpl[T any](iter Iterable[T], n uint) option.Option[T] {
	it := Iterator[T]{iterable: iter}
	if advanceByImpl(iter, n).IsErr() {
		return option.None[T]()
	}
	return it.Next()
}

//go:inline
func nextChunkImpl[T any](iter Iterable[T], n uint) ChunkResult[[]T] {
	if n == 0 {
		return chunkOk[[]T]([]T{})
	}
	result := make([]T, 0, n)
	it := Iterator[T]{iterable: iter}
	for i := uint(0); i < n; i++ {
		item := it.Next()
		if item.IsNone() {
			// Return error with remaining elements
			return chunkErr[[]T](result)
		}
		result = append(result, item.Unwrap())
	}
	return chunkOk[[]T](result)
}
