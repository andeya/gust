package result_test

import (
	"fmt"

	"github.com/andeya/gust"
)

func ExampleResult_AndThen() {
	var divide = func(i, j float32) gust.Result[float32] {
		if j == 0 {
			return gust.Err[float32]("j can not be 0")
		}
		return gust.Ok(i / j)
	}
	var ret = divide(1, 2).AndThen(func(i float32) gust.Result[any] {
		return gust.Ok[any](i * 10)
	}).Unwrap().(float32)
	fmt.Println(ret)
	// Output:
	// 5
}
