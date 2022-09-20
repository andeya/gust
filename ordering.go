package gust

type Ord interface {
	Digit | ~string | ~uintptr
}

type Ordering struct {
	cmp int8
}

func Less() Ordering {
	return Ordering{cmp: -1}
}

func Equal() Ordering {
	return Ordering{cmp: 0}
}

func Greater() Ordering {
	return Ordering{cmp: 1}
}

func Compare[T Ord](a, b T) Ordering {
	if a < b {
		return Less()
	}
	if a == b {
		return Equal()
	}
	return Greater()
}

func (o Ordering) Is(ord Ordering) bool {
	return o.cmp == ord.cmp
}

func (o Ordering) IsLess() bool {
	return o.cmp == -1
}

func (o Ordering) IsEqual() bool {
	return o.cmp == 0
}

func (o Ordering) IsGreater() bool {
	return o.cmp == 1
}
