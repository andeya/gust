package iter

import (
	"testing"

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

