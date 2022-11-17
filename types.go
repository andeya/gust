package gust

// For implements iterators, the following methods are available:

type (
	Iterable[T any] interface {
		Next() Option[T]
	}
	SizeIterable[T any] interface {
		Remaining() uint
	}
	DeIterable[T any] interface {
		Iterable[T]
		SizeIterable[T]
		NextBack() Option[T]
	}
	IterableCount interface {
		Count() uint
	}
	IterableSizeHint interface {
		SizeHint() (uint, Option[uint])
	}
)

// Pair is a pair of values.
type Pair[A any, B any] struct {
	A A
	B B
}

// KV is an index-value pair.
type KV[T any] struct {
	Index uint
	Value T
}
