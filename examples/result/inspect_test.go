package result_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/andeya/gust"
)

func TestInspect(t *testing.T) {
	gust.Ret(strconv.Atoi("4")).
		Inspect(func(x int) {
			fmt.Println("original: ", x)
		}).
		Map(func(x int) int {
			return x * 3
		}).
		Expect("failed to parse number")
}
