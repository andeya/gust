package iter

import "github.com/andeya/gust/digit"

func IterAny[T any](next Nextor[T]) *AnyIter[T] {
	return newAnyIter[T](1, nil, next)
}

func IterAnyFromVec[T any](slice []T) *AnyIter[T] {
	return NewVecNext(slice).ToIter()
}

func IterAnyFromRange[T digit.Integer](start T, end T, rightClosed ...bool) *AnyIter[T] {
	return NewRangeNext[T](start, end, rightClosed...).ToIter()
}

func IterAnyFromChan[T any](c <-chan T) *AnyIter[T] {
	return NewChanNext[T](c).ToIter()
}
