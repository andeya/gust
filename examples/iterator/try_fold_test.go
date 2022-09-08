package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/digit"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestTryFold_1(t *testing.T) {
	var c = make(chan int8, 10)
	for _, i := range []iter.Iterator[int8]{
		iter.FromElements[int8](1, 2, 3).ToInspect(func(v int8) {
			c <- v
		}),
		iter.FromRange[int8](1, 4),
		iter.FromChan(c),
	} {
		// the checked sum of all the elements of the array
		var sum = i.TryFold(int8(0), func(acc any, v int8) gust.AnyCtrlFlow {
			return digit.CheckedAdd[int8](acc.(int8), v).CtrlFlow().ToX()
		})
		assert.Equal(t, gust.AnyContinue(int8(6)), sum)
	}
}

func TestTryFold_2(t *testing.T) {
	var i = iter.FromElements[int8](10, 20, 30, 100, 40, 50)
	// This sum overflows when adding the 100 element
	var sum = i.TryFold(int8(0), func(acc any, v int8) gust.AnyCtrlFlow {
		return digit.CheckedAdd[int8](acc.(int8), v).CtrlFlow().ToX()
	})
	assert.Equal(t, gust.AnyBreak(gust.Void(nil)), sum)
	// Because it short-circuited, the remaining elements are still
	// available through the iterator.
	assert.Equal(t, uint(2), i.Remaining())
	assert.Equal(t, gust.Some[int8](40), i.Next())
}

func TestTryFold_3(t *testing.T) {
	var triangular8 = iter.FromRange[int8](1, 30).TryFold(int8(0), func(acc any, v int8) gust.AnyCtrlFlow {
		return digit.CheckedAdd[int8](acc.(int8), v).XMapOrElse(
			func() any {
				return gust.AnyBreak(acc)
			}, func(sum int8) any {
				return gust.AnyContinue(sum)
			}).(gust.AnyCtrlFlow)
	})
	assert.Equal(t, gust.AnyBreak(int8(120)), triangular8)

	var triangular64 = iter.FromRange[uint64](1, 30).TryFold(uint64(0), func(acc any, v uint64) gust.AnyCtrlFlow {
		return digit.CheckedAdd(acc.(uint64), v).XMapOrElse(
			func() any {
				return gust.AnyBreak(acc)
			}, func(sum uint64) any {
				return gust.AnyContinue(sum)
			}).(gust.AnyCtrlFlow)
	})
	assert.Equal(t, gust.AnyContinue(uint64(435)), triangular64)
}
