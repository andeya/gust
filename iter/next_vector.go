package iter

import (
	"github.com/andeya/gust"
)

var (
	_ NextForIter[any] = (*VecNext[any])(nil)
	_ SizeHintForIter  = (*VecNext[any])(nil)
	_ CountForIter     = (*VecNext[any])(nil)
)

type VecNext[T any] struct {
	slice     []T
	nextIndex int
}

func NewVecNext[T any](slice []T) *VecNext[T] {
	return &VecNext[T]{
		slice:     slice,
		nextIndex: 0,
	}
}

func (v *VecNext[T]) ToIter() *Iter[T] {
	return newIter[T](v)
}

func (v *VecNext[T]) NextForIter() gust.Option[T] {
	if v.nextIndex < len(v.slice) {
		opt := gust.Some(v.slice[v.nextIndex])
		v.nextIndex++
		return opt
	}
	return gust.None[T]()
}

func (v *VecNext[T]) SizeHintForIter() (uint64, gust.Option[uint64]) {
	n := uint64(len(v.slice) - v.nextIndex)
	return n, gust.Some(n)
}

func (v *VecNext[T]) CountForIter() uint64 {
	v.nextIndex = len(v.slice)
	return uint64(len(v.slice) - v.nextIndex)
}
