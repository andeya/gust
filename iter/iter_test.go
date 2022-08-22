package iter

import "testing"

func TestAnyIter(t *testing.T) {
	var iter = IterAnyFromVec([]int{1, 2, 3})
	if !iter.Any(func(x int) bool {
		return x > 1
	}) {
		t.Error("Any failed")
	}
}
