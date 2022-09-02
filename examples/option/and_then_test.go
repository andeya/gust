package option_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/opt"
	"github.com/stretchr/testify/assert"
)

func TestOption_AndThen_1(t *testing.T) {
	var sqThenToString = func(x uint32) gust.Option[string] {
		if x <= math.MaxUint32/x {
			return gust.Some(strconv.FormatUint(uint64(x*x), 10))
		}
		return gust.None[string]()
	}
	assert.Equal(t, gust.Some("4"), opt.AndThen(gust.Some[uint32](2), sqThenToString))
	assert.Equal(t, gust.None[string](), opt.AndThen(gust.Some[uint32](1000000), sqThenToString))
	assert.Equal(t, gust.None[string](), opt.AndThen(gust.None[uint32](), sqThenToString))
}

func TestOption_AndThen_2(t *testing.T) {
	var sqThenToString = func(x uint32) gust.Option[any] {
		if x <= math.MaxUint32/x {
			return gust.Some[any](strconv.FormatUint(uint64(x*x), 10))
		}
		return gust.None[any]()
	}
	assert.Equal(t, gust.Some("4").ToX(), gust.Some[uint32](2).XAndThen(sqThenToString))
	assert.Equal(t, gust.None[string]().ToX(), gust.Some[uint32](1000000).XAndThen(sqThenToString))
	assert.Equal(t, gust.None[string]().ToX(), gust.None[uint32]().XAndThen(sqThenToString))
}
