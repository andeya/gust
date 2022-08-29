package gust

// For implements iterators, the following methods are available:

type (
	DataForIter[T any] interface {
		NextForIter() Option[T]
	}
	DataForDoubleEndedIter[T any] interface {
		DataForIter[T]
		NextBackForIter() Option[T]
		RemainingLenForIter() uint
	}
	CountForIter interface {
		CountForIter() uint
	}
	SizeHintForIter interface {
		SizeHintForIter() (uint, Option[uint])
	}
)
