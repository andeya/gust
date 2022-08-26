package iter

import (
	"math"

	"github.com/andeya/gust"
)

func saturatingAdd(a, b uint64) uint64 {
	if a < math.MaxUint64-b {
		return a + b
	}
	return math.MaxUint64
}

func checkedAdd(a, b uint64) gust.Option[uint64] {
	if a <= math.MaxUint64-b {
		return gust.Some(a + b)
	}
	return gust.None[uint64]()
}

func uintCheckedMul(a, b uint) gust.Option[uint] {
	if a <= math.MaxUint64/b {
		return gust.Some(a * b)
	}
	return gust.None[uint]()
}

func uint64CheckedMul(a, b uint64) gust.Option[uint64] {
	if a <= math.MaxUint64/b {
		return gust.Some(a * b)
	}
	return gust.None[uint64]()
}
