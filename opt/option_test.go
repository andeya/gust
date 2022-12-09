package opt_test

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

func TestOption_Contains(t *testing.T) {
	{
		var x = gust.Some(2)
		assert.Equal(t, true, opt.Contains(x, 2))
	}
	{
		var x = gust.Some(3)
		assert.Equal(t, false, opt.Contains(x, 2))
	}
	{
		var x = gust.None[int]()
		assert.Equal(t, false, opt.Contains(x, 2))
	}
}

func TestOption_MapOrElse_1(t *testing.T) {
	var k = 21
	{
		var x = gust.Some("foo")
		assert.Equal(t, 3, opt.MapOrElse(x, func() int { return 2 * k }, func(v string) int { return len(v) }))
	}
	{
		var x gust.Option[string]
		assert.Equal(t, 42, opt.MapOrElse(x, func() int { return 2 * k }, func(v string) int { return len(v) }))
	}
}

func TestOption_MapOr_1(t *testing.T) {
	{
		var x = gust.Some("foo")
		assert.Equal(t, 3, opt.MapOr(x, 42, func(v string) int { return len(v) }))
	}
	{
		var x gust.Option[string]
		assert.Equal(t, 42, opt.MapOr(x, 42, func(v string) int { return len(v) }))
	}
}

func TestOption_Map_1(t *testing.T) {
	var maybeSomeString = gust.Some("Hello, World!")
	var maybeSomeLen = opt.Map(maybeSomeString, func(s string) int { return len(s) })
	assert.Equal(t, maybeSomeLen, gust.Some(13))
}

func TestOption_Unzip(t *testing.T) {
	var x = gust.Some[gust.Pair[int, string]](gust.Pair[int, string]{A: 1, B: "hi"})
	var y = gust.None[gust.Pair[int, string]]()
	assert.Equal(t, opt.Unzip(x), gust.Pair[gust.Option[int], gust.Option[string]]{A: gust.Some[int](1), B: gust.Some[string]("hi")})
	assert.Equal(t, opt.Unzip(y), gust.Pair[gust.Option[int], gust.Option[string]]{A: gust.None[int](), B: gust.None[string]()})
}

func TestOption_Zip(t *testing.T) {
	var x = gust.Some[byte](1)
	var y = gust.Some("hi")
	var z = gust.None[byte]()
	assert.Equal(t, gust.Some(gust.Pair[byte, string]{1, "hi"}), opt.Zip(x, y))
	assert.Equal(t, gust.None[gust.Pair[byte, byte]](), opt.Zip(x, z))
}

func TestOption_ZipWith(t *testing.T) {
	type Point struct {
		x float64
		y float64
	}
	var newPoint = func(x float64, y float64) Point {
		return Point{x, y}
	}
	var x = gust.Some(17.5)
	var y = gust.Some(42.7)
	assert.Equal(t, opt.ZipWith(x, y, newPoint), gust.Some(Point{x: 17.5, y: 42.7}))
	assert.Equal(t, opt.ZipWith(x, gust.None[float64](), newPoint), gust.None[Point]())
}
