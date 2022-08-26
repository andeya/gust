package iter

import (
	"github.com/andeya/gust"
)

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
