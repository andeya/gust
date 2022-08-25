package iter

import "github.com/andeya/gust/digit"

func FromVec[T any](slice []T) *Iter[T] {
	return NewVecNext(slice).ToIter()
}

func FromRange[T digit.Integer](start T, end T, rightClosed ...bool) *Iter[T] {
	return NewRangeNext[T](start, end, rightClosed...).ToIter()
}

func FromChan[T any](c <-chan T) *Iter[T] {
	return NewChanNext[T](c).ToIter()
}
