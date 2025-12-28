package iterator_test

import (
	"testing"

	"github.com/andeya/gust/iterator"
	"github.com/andeya/gust/option"
	"github.com/stretchr/testify/assert"
)

func TestSum(t *testing.T) {
	a := []int{1, 2, 3}
	sum := iterator.Sum(iterator.FromSlice(a))
	assert.Equal(t, 6, sum)

	b := []float64{}
	sumFloat := iterator.Sum(iterator.FromSlice(b))
	assert.Equal(t, 0.0, sumFloat)
}

func TestProduct(t *testing.T) {
	factorial := func(n int) int {
		return iterator.Product(iterator.FromRange(1, n+1))
	}
	assert.Equal(t, 1, factorial(0))
	assert.Equal(t, 1, factorial(1))
	assert.Equal(t, 120, factorial(5))
}

func TestProduct_AllTypes(t *testing.T) {
	// Test Product with all numeric types to cover getProductIdentity
	// int
	assert.Equal(t, 1, iterator.Product(iterator.FromSlice([]int{})))
	assert.Equal(t, 6, iterator.Product(iterator.FromSlice([]int{1, 2, 3})))

	// int8
	assert.Equal(t, int8(1), iterator.Product(iterator.FromSlice([]int8{})))
	assert.Equal(t, int8(6), iterator.Product(iterator.FromSlice([]int8{1, 2, 3})))

	// int16
	assert.Equal(t, int16(1), iterator.Product(iterator.FromSlice([]int16{})))
	assert.Equal(t, int16(6), iterator.Product(iterator.FromSlice([]int16{1, 2, 3})))

	// int32
	assert.Equal(t, int32(1), iterator.Product(iterator.FromSlice([]int32{})))
	assert.Equal(t, int32(6), iterator.Product(iterator.FromSlice([]int32{1, 2, 3})))

	// int64
	assert.Equal(t, int64(1), iterator.Product(iterator.FromSlice([]int64{})))
	assert.Equal(t, int64(6), iterator.Product(iterator.FromSlice([]int64{1, 2, 3})))

	// uint
	assert.Equal(t, uint(1), iterator.Product(iterator.FromSlice([]uint{})))
	assert.Equal(t, uint(6), iterator.Product(iterator.FromSlice([]uint{1, 2, 3})))

	// uint8
	assert.Equal(t, uint8(1), iterator.Product(iterator.FromSlice([]uint8{})))
	assert.Equal(t, uint8(6), iterator.Product(iterator.FromSlice([]uint8{1, 2, 3})))

	// uint16
	assert.Equal(t, uint16(1), iterator.Product(iterator.FromSlice([]uint16{})))
	assert.Equal(t, uint16(6), iterator.Product(iterator.FromSlice([]uint16{1, 2, 3})))

	// uint32
	assert.Equal(t, uint32(1), iterator.Product(iterator.FromSlice([]uint32{})))
	assert.Equal(t, uint32(6), iterator.Product(iterator.FromSlice([]uint32{1, 2, 3})))

	// uint64
	assert.Equal(t, uint64(1), iterator.Product(iterator.FromSlice([]uint64{})))
	assert.Equal(t, uint64(6), iterator.Product(iterator.FromSlice([]uint64{1, 2, 3})))

	// float32
	assert.Equal(t, float32(1.0), iterator.Product(iterator.FromSlice([]float32{})))
	assert.Equal(t, float32(6.0), iterator.Product(iterator.FromSlice([]float32{1.0, 2.0, 3.0})))

	// float64
	assert.Equal(t, 1.0, iterator.Product(iterator.FromSlice([]float64{})))
	assert.Equal(t, 6.0, iterator.Product(iterator.FromSlice([]float64{1.0, 2.0, 3.0})))
}

func TestSum_AllTypes(t *testing.T) {
	// Test Sum with all numeric types
	// int
	assert.Equal(t, 0, iterator.Sum(iterator.FromSlice([]int{})))
	assert.Equal(t, 6, iterator.Sum(iterator.FromSlice([]int{1, 2, 3})))

	// float64
	assert.Equal(t, 0.0, iterator.Sum(iterator.FromSlice([]float64{})))
	assert.Equal(t, 6.0, iterator.Sum(iterator.FromSlice([]float64{1.0, 2.0, 3.0})))
}

type nonDoubleEndedIterable struct {
	values []int
	index  int
}

func (n *nonDoubleEndedIterable) Next() option.Option[int] {
	if n.index >= len(n.values) {
		return option.None[int]()
	}
	val := n.values[n.index]
	n.index++
	return option.Some(val)
}

func (n *nonDoubleEndedIterable) SizeHint() (uint, option.Option[uint]) {
	return 0, option.None[uint]()
}

func TestMustToDoubleEnded_Panic(t *testing.T) {
	// Test MustToDoubleEnded with non-double-ended iterator (should panic)
	// Create a non-double-ended iterator using FromIterable
	nonDE := &nonDoubleEndedIterable{values: []int{1, 2, 3}, index: 0}
	var iterable iterator.Iterable[int] = nonDE
	var gustIter iterator.Iterable[int] = iterable
	iter := iterator.FromIterable(gustIter)

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustToDoubleEnded should panic for non-double-ended iterator")
		}
	}()

	_ = iter.MustToDoubleEnded()
}

type customIterable struct {
	values []int
	index  int
}

func (c *customIterable) Next() option.Option[int] {
	if c.index >= len(c.values) {
		return option.None[int]()
	}
	val := c.values[c.index]
	c.index++
	return option.Some(val)
}

func (c *customIterable) SizeHint() (uint, option.Option[uint]) {
	return 0, option.None[uint]()
}

func TestFromIterable_IterablePath(t *testing.T) {
	// Test FromIterable with Iterable[T] path (not Iterator[T])
	custom := &customIterable{values: []int{10, 20, 30}, index: 0}
	var iterable iterator.Iterable[int] = custom
	var gustIter iterator.Iterable[int] = iterable
	iter := iterator.FromIterable(gustIter)
	assert.Equal(t, option.Some(10), iter.Next())
	assert.Equal(t, option.Some(20), iter.Next())
	assert.Equal(t, option.Some(30), iter.Next())
	assert.Equal(t, option.None[int](), iter.Next())
}
