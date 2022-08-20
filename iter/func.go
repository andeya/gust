package iter

import (
	"github.com/andeya/gust"
)

func TryFold[T any, B any](next Nextor[T], init B, f func(B, T) gust.Result[B]) gust.Result[B] {
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

func TryRFold[T any, B any](next Nextor[T], init B, f func(B, T) gust.Result[B]) gust.Result[B] {
	// FIXME
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
