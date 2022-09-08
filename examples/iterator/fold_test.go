package iterator_test

import (
	"fmt"
	"testing"

	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestFold_1(t *testing.T) {
	var c = make(chan int8, 10)
	for _, i := range []iter.Iterator[int8]{
		iter.FromElements[int8](1, 2, 3).ToInspect(func(v int8) {
			c <- v
		}),
		iter.FromRange[int8](1, 3, true),
		iter.FromChan(c),
	} {
		// the sum of all of the elements of the array
		var sum = i.Fold(int8(0), func(acc any, v int8) any {
			return acc.(int8) + v
		}).(int8)
		assert.Equal(t, int8(6), sum)
	}
}

func TestFold_2(t *testing.T) {
	var c = make(chan int8, 10)
	for _, i := range []iter.Iterator[int8]{
		iter.FromElements[int8](1, 2, 3, 4, 5).ToInspect(func(v int8) {
			c <- v
		}),
		iter.FromRange[int8](1, 5, true),
		iter.FromChan(c),
	} {
		var zero = "0"
		var result = i.Fold(zero, func(acc any, v int8) any {
			return fmt.Sprintf("(%s + %d)", acc.(string), v)
		}).(string)
		assert.Equal(t, "(((((0 + 1) + 2) + 3) + 4) + 5)", result)
	}
}

func TestFold_3(t *testing.T) {
	var numbers = []int{1, 2, 3, 4, 5}

	// for loop:
	var result = 0
	for _, x := range numbers {
		result = result + x
	}

	// fold:
	var result2 = iter.FromVec(numbers).Fold(0, func(acc any, x int) any {
		return acc.(int) + x
	}).(int)

	// they're the same
	assert.Equal(t, result, result2)
}
