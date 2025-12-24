package iter

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestSum(t *testing.T) {
	a := []int{1, 2, 3}
	sum := Sum(FromSlice(a))
	assert.Equal(t, 6, sum)

	b := []float64{}
	sumFloat := Sum(FromSlice(b))
	assert.Equal(t, 0.0, sumFloat)
}

func TestProduct(t *testing.T) {
	factorial := func(n int) int {
		return Product(FromRange(1, n+1))
	}
	assert.Equal(t, 1, factorial(0))
	assert.Equal(t, 1, factorial(1))
	assert.Equal(t, 120, factorial(5))
}

func TestProduct_AllTypes(t *testing.T) {
	// Test Product with all numeric types to cover getProductIdentity
	// int
	assert.Equal(t, 1, Product(FromSlice([]int{})))
	assert.Equal(t, 6, Product(FromSlice([]int{1, 2, 3})))
	
	// int8
	assert.Equal(t, int8(1), Product(FromSlice([]int8{})))
	assert.Equal(t, int8(6), Product(FromSlice([]int8{1, 2, 3})))
	
	// int16
	assert.Equal(t, int16(1), Product(FromSlice([]int16{})))
	assert.Equal(t, int16(6), Product(FromSlice([]int16{1, 2, 3})))
	
	// int32
	assert.Equal(t, int32(1), Product(FromSlice([]int32{})))
	assert.Equal(t, int32(6), Product(FromSlice([]int32{1, 2, 3})))
	
	// int64
	assert.Equal(t, int64(1), Product(FromSlice([]int64{})))
	assert.Equal(t, int64(6), Product(FromSlice([]int64{1, 2, 3})))
	
	// uint
	assert.Equal(t, uint(1), Product(FromSlice([]uint{})))
	assert.Equal(t, uint(6), Product(FromSlice([]uint{1, 2, 3})))
	
	// uint8
	assert.Equal(t, uint8(1), Product(FromSlice([]uint8{})))
	assert.Equal(t, uint8(6), Product(FromSlice([]uint8{1, 2, 3})))
	
	// uint16
	assert.Equal(t, uint16(1), Product(FromSlice([]uint16{})))
	assert.Equal(t, uint16(6), Product(FromSlice([]uint16{1, 2, 3})))
	
	// uint32
	assert.Equal(t, uint32(1), Product(FromSlice([]uint32{})))
	assert.Equal(t, uint32(6), Product(FromSlice([]uint32{1, 2, 3})))
	
	// uint64
	assert.Equal(t, uint64(1), Product(FromSlice([]uint64{})))
	assert.Equal(t, uint64(6), Product(FromSlice([]uint64{1, 2, 3})))
	
	// float32
	assert.Equal(t, float32(1.0), Product(FromSlice([]float32{})))
	assert.Equal(t, float32(6.0), Product(FromSlice([]float32{1.0, 2.0, 3.0})))
	
	// float64
	assert.Equal(t, 1.0, Product(FromSlice([]float64{})))
	assert.Equal(t, 6.0, Product(FromSlice([]float64{1.0, 2.0, 3.0})))
}

func TestSum_AllTypes(t *testing.T) {
	// Test Sum with all numeric types
	// int
	assert.Equal(t, 0, Sum(FromSlice([]int{})))
	assert.Equal(t, 6, Sum(FromSlice([]int{1, 2, 3})))
	
	// float64
	assert.Equal(t, 0.0, Sum(FromSlice([]float64{})))
	assert.Equal(t, 6.0, Sum(FromSlice([]float64{1.0, 2.0, 3.0})))
}

type nonDoubleEndedIter struct {
	values []int
	index  int
}

func (n *nonDoubleEndedIter) Next() gust.Option[int] {
	if n.index >= len(n.values) {
		return gust.None[int]()
	}
	val := n.values[n.index]
	n.index++
	return gust.Some(val)
}

func (n *nonDoubleEndedIter) SizeHint() (uint, gust.Option[uint]) {
	return 0, gust.None[uint]()
}

func TestMustToDoubleEnded_Panic(t *testing.T) {
	// Test MustToDoubleEnded with non-double-ended iterator (should panic)
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustToDoubleEnded should panic for non-double-ended iterator")
		}
	}()
	
	iter := Iterator[int]{
		iter: &nonDoubleEndedIter{values: []int{1, 2, 3}, index: 0},
	}
	
	_ = iter.MustToDoubleEnded()
}

type customIterable struct {
	values []int
	index  int
}

func (c *customIterable) Next() gust.Option[int] {
	if c.index >= len(c.values) {
		return gust.None[int]()
	}
	val := c.values[c.index]
	c.index++
	return gust.Some(val)
}

func (c *customIterable) SizeHint() (uint, gust.Option[uint]) {
	return 0, gust.None[uint]()
}

func TestFromIterable_IterablePath(t *testing.T) {
	// Test FromIterable with Iterable[T] path (not Iterator[T])
	custom := &customIterable{values: []int{10, 20, 30}, index: 0}
	var iterable Iterable[int] = custom
	var gustIter gust.Iterable[int] = iterable
	iter := FromIterable(gustIter)
	assert.Equal(t, gust.Some(10), iter.Next())
	assert.Equal(t, gust.Some(20), iter.Next())
	assert.Equal(t, gust.Some(30), iter.Next())
	assert.Equal(t, gust.None[int](), iter.Next())
}

