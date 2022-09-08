package iterator_test

import (
	"testing"

	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestFlatMap(t *testing.T) {
	var c = make(chan string, 10)
	for _, i := range []iter.Iterator[string]{
		iter.FromElements("alpha", "beta", "gamma").ToInspect(func(v string) {
			c <- v
		}),
		iter.FromChan(c),
	} {
		var merged = iter.ToFlatMap(i, func(t string) iter.Iterator[rune] {
			return iter.FromString[rune](t)
		}).Collect()
		assert.Equal(t, "alphabetagamma", string(merged))
	}
}
