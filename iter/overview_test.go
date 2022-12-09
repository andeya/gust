package iter_test

import (
	"fmt"

	"github.com/andeya/gust/iter"
)

func ExampleIterator() {
	var v = []int{1, 2, 3, 4, 5}
	iter.FromVec(v).ForEach(func(x int) {
		fmt.Printf("%d\n", x)
	})
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}
