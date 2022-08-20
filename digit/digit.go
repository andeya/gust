package digit

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Digit interface {
	Integer | ~float32 | ~float64
}

func Abs[T Digit](d T) T {
	var zero T
	if d < zero {
		return -d
	}
	return d
}
