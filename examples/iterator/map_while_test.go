package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestMapWhile_1(t *testing.T) {
	var c = make(chan int, 10)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(-1, 4, 0, 1).Inspect(func(v int) {
			c <- v
		}),
		iter.FromChan(c),
	} {
		i := i.MapWhile(func(x int) gust.Option[int] { return checkedDivide(16, x) })
		assert.Equal(t, gust.Some(-16), i.Next())
		assert.Equal(t, gust.Some(4), i.Next())
		assert.Equal(t, gust.None[int](), i.Next())
		assert.Equal(t, gust.Some(16), i.Next())
	}
}

func checkedDivide(x, y int) gust.Option[int] {
	if y == 0 {
		return gust.None[int]()
	}
	return gust.Some(x / y)
}

func TestMapWhile_2(t *testing.T) {
	var c = make(chan int, 10)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(1, 2, -3, 4).Inspect(func(v int) {
			c <- v
		}),
		iter.FromChan(c),
	} {
		a := i.XMapWhile(func(x int) gust.Option[any] {
			if x < 0 {
				return gust.None[any]()
			}
			return gust.Some[any](uint(x))
		}).Collect()
		assert.Equal(t, []any{uint(1), uint(2)}, a)
	}
}
