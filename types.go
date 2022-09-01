package gust

// For implements iterators, the following methods are available:

type (
	Iterable[T any] interface {
		Next() Option[T]
	}
	DeIterable[T any] interface {
		Iterable[T]
		NextBack() Option[T]
	}
	IterableCount interface {
		Count() uint
	}
	IterableSizeHint interface {
		SizeHint() (uint, Option[uint])
	}
	SizeIterable[T any] interface {
		Remaining() uint
	}
	// SizeDeIterable is a double ended iterator that knows the exact size.
	SizeDeIterable[T any] interface {
		DeIterable[T]
		SizeIterable[T]
	}
)

// Pair is a pair of values.
type Pair[A any, B any] struct {
	A A
	B B
}
