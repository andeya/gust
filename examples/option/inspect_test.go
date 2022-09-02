package option_test

import (
	"fmt"

	"github.com/andeya/gust"
)

func ExampleOption_Inspect() {
	// prints "got: 3"
	_ = gust.Some(3).Inspect(func(x int) {
		fmt.Println("got:", x)
	})

	// prints nothing
	_ = gust.None[int]().Inspect(func(x int) {
		fmt.Println("got:", x)
	})

	// Output:
	// got: 3
}
