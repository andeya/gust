package pair

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPair_Split(t *testing.T) {
	// Test Pair.Split with int types
	pair := Pair[int, string]{A: 42, B: "hello"}
	a, b := pair.Split()
	assert.Equal(t, 42, a)
	assert.Equal(t, "hello", b)

	// Test Pair.Split with different types
	pair2 := Pair[string, int]{A: "test", B: 100}
	a2, b2 := pair2.Split()
	assert.Equal(t, "test", a2)
	assert.Equal(t, 100, b2)

	// Test Pair.Split with float types
	pair3 := Pair[float64, bool]{A: 3.14, B: true}
	a3, b3 := pair3.Split()
	assert.Equal(t, 3.14, a3)
	assert.Equal(t, true, b3)
}
