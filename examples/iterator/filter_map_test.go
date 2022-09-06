package iterator_test

import (
	"strconv"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestFilterMap_1(t *testing.T) {
	var c = make(chan string, 10)
	for idx, i := range []iter.Iterator[string]{
		iter.FromElements("1", "two", "NaN", "four", "5"),
		iter.FromChan(c),
	} {
		var i = iter.FilterMap[string, int](i.Inspect(func(v string) {
			if idx == 0 {
				c <- v
			}
		}), func(v string) gust.Option[int] { return gust.Ret(strconv.Atoi(v)).Ok() })
		assert.Equal(t, gust.Some[int](1), i.Next())
		assert.Equal(t, gust.Some[int](5), i.Next())
		assert.Equal(t, gust.None[int](), i.Next())
	}
}

func TestFilterMap_2(t *testing.T) {
	var c = make(chan string, 10)
	for idx, i := range []iter.Iterator[string]{
		iter.FromElements("1", "two", "NaN", "four", "5"),
		iter.FromChan(c),
	} {
		var i = i.Inspect(func(v string) {
			if idx == 0 {
				c <- v
			}
		}).XFilterMap(func(v string) gust.Option[any] { return gust.Ret(strconv.Atoi(v)).Ok().ToX() })
		assert.Equal(t, gust.Some[any](1), i.Next())
		assert.Equal(t, gust.Some[any](5), i.Next())
		assert.Equal(t, gust.None[any](), i.Next())
	}
}
