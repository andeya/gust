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

// TestOption_Map_None tests Map with None option (covers opt/option.go:60-64)
func TestOption_Map_None(t *testing.T) {
	var x gust.Option[int]
	result := opt.Map(x, func(v int) string {
		return strconv.Itoa(v)
	})
	assert.True(t, result.IsNone())
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

func TestSafeAssert(t *testing.T) {
	// Test with valid type assertion
	opt1 := gust.Some(42)
	result1 := opt.SafeAssert[int, int](opt1)
	assert.True(t, result1.IsOk())
	assert.True(t, result1.Unwrap().IsSome())
	assert.Equal(t, 42, result1.Unwrap().Unwrap())

	// Test with invalid type assertion
	opt2 := gust.Some(42)
	result2 := opt.SafeAssert[int, string](opt2)
	assert.True(t, result2.IsErr())

	// Test with None
	opt3 := gust.None[int]()
	result3 := opt.SafeAssert[int, string](opt3)
	assert.True(t, result3.IsOk())
	assert.True(t, result3.Unwrap().IsNone())
}

func TestXSafeAssert(t *testing.T) {
	// Test with valid type assertion
	opt1 := gust.Some[any](42)
	result1 := opt.XSafeAssert[int](opt1)
	assert.True(t, result1.IsOk())
	assert.True(t, result1.Unwrap().IsSome())
	assert.Equal(t, 42, result1.Unwrap().Unwrap())

	// Test with invalid type assertion
	opt2 := gust.Some[any](42)
	result2 := opt.XSafeAssert[string](opt2)
	assert.True(t, result2.IsErr())

	// Test with None
	opt3 := gust.None[any]()
	result3 := opt.XSafeAssert[int](opt3)
	assert.True(t, result3.IsOk())
	assert.True(t, result3.Unwrap().IsNone())
}

func TestFuzzyAssert(t *testing.T) {
	// Test with valid type assertion
	opt1 := gust.Some(42)
	result1 := opt.FuzzyAssert[int, int](opt1)
	assert.True(t, result1.IsSome())
	assert.Equal(t, 42, result1.Unwrap())

	// Test with invalid type assertion
	opt2 := gust.Some(42)
	result2 := opt.FuzzyAssert[int, string](opt2)
	assert.True(t, result2.IsNone())

	// Test with None
	opt3 := gust.None[int]()
	result3 := opt.FuzzyAssert[int, string](opt3)
	assert.True(t, result3.IsNone())
}

func TestXFuzzyAssert(t *testing.T) {
	// Test with valid type assertion
	opt1 := gust.Some[any](42)
	result1 := opt.XFuzzyAssert[int](opt1)
	assert.True(t, result1.IsSome())
	assert.Equal(t, 42, result1.Unwrap())

	// Test with invalid type assertion
	opt2 := gust.Some[any](42)
	result2 := opt.XFuzzyAssert[string](opt2)
	assert.True(t, result2.IsNone())

	// Test with None
	opt3 := gust.None[any]()
	result3 := opt.XFuzzyAssert[int](opt3)
	assert.True(t, result3.IsNone())
}

func TestAnd(t *testing.T) {
	// Test with Some and Some
	opt1 := gust.Some(2)
	opt2 := gust.Some(3)
	result1 := opt.And(opt1, opt2)
	assert.True(t, result1.IsSome())
	assert.Equal(t, 3, result1.Unwrap())

	// Test with Some and None
	opt3 := gust.Some(2)
	opt4 := gust.None[int]()
	result2 := opt.And(opt3, opt4)
	assert.True(t, result2.IsNone())

	// Test with None and Some
	opt5 := gust.None[int]()
	opt6 := gust.Some(3)
	result3 := opt.And(opt5, opt6)
	assert.True(t, result3.IsNone())

	// Test with None and None
	opt7 := gust.None[int]()
	opt8 := gust.None[int]()
	result4 := opt.And(opt7, opt8)
	assert.True(t, result4.IsNone())
}

func TestSafeAssert_TypeAssertionError(t *testing.T) {
	// Test SafeAssert with type assertion error message
	opt1 := gust.Some(42)
	result1 := opt.SafeAssert[int, string](opt1)
	assert.True(t, result1.IsErr())
	assert.Contains(t, result1.Err().Error(), "type assert error")
}

func TestXSafeAssert_TypeAssertionError(t *testing.T) {
	// Test XSafeAssert with type assertion error message
	opt1 := gust.Some[any](42)
	result1 := opt.XSafeAssert[string](opt1)
	assert.True(t, result1.IsErr())
	assert.Contains(t, result1.Err().Error(), "type assert error")
}
