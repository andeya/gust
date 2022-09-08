package iterator_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

var _ gust.Iterable[int] = (*Alternate)(nil)

type Alternate struct {
	state int
}

func (a *Alternate) Next() gust.Option[int] {
	var val = a.state
	a.state = a.state + 1
	// if it's even, Some(i32), else None
	if val%2 == 0 {
		return gust.Some(val)
	}
	return gust.None[int]()
}

func TestFuse(t *testing.T) {
	var a = &Alternate{state: 0}
	var i = iter.FromIterable[int](a)
	// we can see our iterator going back and forth
	assert.Equal(t, gust.Some(0), i.Next())
	assert.Equal(t, gust.None[int](), i.Next())
	assert.Equal(t, gust.Some(2), i.Next())
	assert.Equal(t, gust.None[int](), i.Next())
	// however, once we fuse it...
	var j = i.ToFuse()
	assert.Equal(t, gust.Some(4), j.Next())
	assert.Equal(t, gust.None[int](), j.Next())
	// it will always return `None` after the first time.
	assert.Equal(t, gust.None[int](), j.Next())
	assert.Equal(t, gust.None[int](), j.Next())
	assert.Equal(t, gust.None[int](), j.Next())

}
