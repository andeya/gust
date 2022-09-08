package iterator_test

import (
	"strconv"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestFindMap(t *testing.T) {
	var firstNumbe = iter.FromElements("lol", "NaN", "2", "5").
		XFindMap(func(s string) gust.Option[any] {
			return gust.Ret(strconv.Atoi(s)).XOk()
		})
	assert.Equal(t, gust.Some[any](int(2)), firstNumbe)
}
