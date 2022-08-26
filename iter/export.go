package iter

import (
	"github.com/andeya/gust"
	"github.com/andeya/gust/digit"
)

func FromVec[T any](slice []T) *Iter[T] {
	return NewVecNext(slice).ToIter()
}

func FromRange[T digit.Integer](start T, end T, rightClosed ...bool) *Iter[T] {
	return NewRangeNext[T](start, end, rightClosed...).ToIter()
}

func FromChan[T any](c <-chan T) *Iter[T] {
	return NewChanNext[T](c).ToIter()
}

func TryFold[T any, B any](next iNext[T], init B, f func(B, T) gust.Result[B]) gust.Result[B] {
	var accum = gust.Ok(init)
	for {
		x := next.Next()
		if x.IsNone() {
			break
		}
		accum = f(accum.Unwrap(), x.Unwrap())
		if accum.IsErr() {
			return accum
		}
	}
	return accum
}

func Fold[T any, B any](next iNext[T], init B, f func(B, T) B) B {
	var accum = init
	for {
		x := next.Next()
		if x.IsNone() {
			break
		}
		accum = f(accum, x.Unwrap())
	}
	return accum
}

func Map[T any, B any](iter Iterator[T], f func(T) B) *MapIterator[T, B] {
	return newMapIterator(iter, f)
}

func FindMap[T any, B any](iter Iterator[T], f func(T) gust.Option[B]) gust.Option[B] {
	for {
		x := iter.Next()
		if x.IsNone() {
			break
		}
		y := f(x.Unwrap())
		if y.IsSome() {
			return y
		}
	}
	return gust.None[B]()
}

// FromIterator conversion from an [`Iterator`].
//
// By implementing `FromIterator` for a type, you define how it will be
// created from an iterator. This is common for types which describe a
// collection of some kind.
type FromIterator[T any, R any] interface {
	FromIter(Iterator[T]) R
}

// Collect collects all the items in the iterator into a slice.
func Collect[T any, R any](iter Iterator[T], fromIter FromIterator[T, R]) R {
	return fromIter.FromIter(iter)
}
