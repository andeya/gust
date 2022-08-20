package iter

import "github.com/andeya/gust"

type VecNext[T any] struct {
	slice     []T
	nextIndex int
}

var (
	_ Nextor[any] = new(VecNext[any])
	_ SizeHint    = new(VecNext[any])
	_ counter     = new(VecNext[any])
)

func NewVecNext[T any](slice []T) *VecNext[T] {
	return &VecNext[T]{
		slice:     slice,
		nextIndex: 0,
	}
}

func (v *VecNext[T]) ToIter() *AnyIter[T] {
	return IterAny[T](v)
}

func (v *VecNext[T]) Next() gust.Option[T] {
	if v.nextIndex < len(v.slice) {
		v.nextIndex++
		return gust.Some(v.slice[v.nextIndex])
	}
	return gust.None[T]()
}

func (v *VecNext[T]) SizeHint() (uint64, gust.Option[uint64]) {
	n := uint64(len(v.slice) - v.nextIndex)
	return n, gust.Some(n)
}

func (v *VecNext[T]) count() uint64 {
	v.nextIndex = len(v.slice)
	return uint64(len(v.slice) - v.nextIndex)
}
