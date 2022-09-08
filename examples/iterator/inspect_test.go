package iterator_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestInspect(t *testing.T) {
	var numbers = iter.FromElements[string]("1", "2", "a", "3", "b").
		ToXMap(func(v string) any {
			return gust.Ret(strconv.Atoi(v))
		}).
		ToInspect(func(v any) {
			v.(gust.Result[int]).InspectErr(func(err error) {
				fmt.Println("Parsing error:", err)
			})
		}).
		ToFilterMap(func(v any) gust.Option[any] {
			return v.(gust.Result[int]).Ok().ToX()
		}).
		Collect()
	assert.Equal(t, []interface{}{1, 2, 3}, numbers)
}
