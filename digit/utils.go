package digit

import (
	"math"

	"github.com/andeya/gust"
)

func Abs[T Digit](d T) T {
	var zero T
	if d < zero {
		return -d
	}
	return d
}

func Max[T Integer]() T {
	var t T
	switch any(t).(type) {
	case int:
		var max = math.MaxInt
		return T(max)
	case int8:
		return T(math.MaxInt8)
	case int16:
		var max int16 = math.MaxInt16
		return T(max)
	case int32:
		var max int32 = math.MaxInt32
		return T(max)
	case int64:
		var max int64 = math.MaxInt64
		return T(max)
	case uint:
		var max uint64 = math.MaxUint
		return T(max)
	case uint8:
		var max uint8 = math.MaxUint8
		return T(max)
	case uint16:
		var max uint16 = math.MaxUint16
		return T(max)
	case uint32:
		var max uint32 = math.MaxUint32
		return T(max)
	case uint64:
		var max uint64 = math.MaxUint64
		return T(max)
	default:
		return t
	}
}

func SaturatingAdd[T Integer](a, b T) T {
	if a < Max[T]()-b {
		return a + b
	}
	return Max[T]()
}

func SaturatingSub[T Digit](a, b T) T {
	if a > b {
		return a - b
	}
	return 0
}

func CheckedAdd[T Integer](a, b T) gust.Option[T] {
	if a <= Max[T]()-b {
		return gust.Some(a + b)
	}
	return gust.None[T]()
}

func UintCheckedMul[T Integer](a, b T) gust.Option[T] {
	if a <= Max[T]()/b {
		return gust.Some(a * b)
	}
	return gust.None[T]()
}
