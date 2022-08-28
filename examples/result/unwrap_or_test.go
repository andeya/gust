package result_test

import (
	"fmt"
	"strconv"

	"github.com/andeya/gust"
)

func ExampleResult_UnwrapOr() {
	const def int = 10

	// before
	i, err := strconv.Atoi("1")
	if err != nil {
		i = def
	}
	fmt.Println(i * 2)

	// now
	fmt.Println(gust.Ret(strconv.Atoi("1")).UnwrapOr(def) * 2)

	// Output:
	// 2
	// 2
}
