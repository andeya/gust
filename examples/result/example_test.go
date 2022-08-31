package result

import (
	"fmt"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/ret"
	"github.com/stretchr/testify/assert"
)

// [`Result`] comes with some convenience methods that make working with it more succinct.
func TestExample(t *testing.T) {
	var goodResult1 = gust.Ok(10)
	var badResult1 = gust.Err[int](10)

	// The `IsOk` and `IsErr` methods do what they say.
	assert.True(t, goodResult1.IsOk() && !goodResult1.IsErr())
	assert.True(t, badResult1.IsErr() && !badResult1.IsOk())

	// `map` consumes the `Result` and produces another.
	var goodResult2 = goodResult1.Map(func(i int) int { return i + 1 })
	var badResult2 = badResult1.Map(func(i int) int { return i - 1 })

	// Use `AndThen` to continue the computation.
	var goodResult3 = ret.AndThen(goodResult2, func(i int) gust.Result[bool] { return gust.Ok(i == 11) })

	// Use `OrElse` to handle the error.
	var _ = badResult2.OrElse(func(err error) gust.Result[int] {
		fmt.Println(err)
		return gust.Ok(20)
	})

	// Consume the result and return the contents with `Unwrap`.
	var _ = goodResult3.Unwrap()
}
