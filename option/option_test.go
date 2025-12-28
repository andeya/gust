package option_test

import (
	"errors"
	"math"
	"strconv"
	"testing"

	"github.com/andeya/gust/option"
	"github.com/andeya/gust/pair"
	"github.com/stretchr/testify/assert"
)

func TestOption_AndThen_1(t *testing.T) {
	var sqThenToString = func(x uint32) option.Option[string] {
		if x <= math.MaxUint32/x {
			return option.Some(strconv.FormatUint(uint64(x*x), 10))
		}
		return option.None[string]()
	}
	assert.Equal(t, option.Some("4"), option.AndThen(option.Some[uint32](2), sqThenToString))
	assert.Equal(t, option.None[string](), option.AndThen(option.Some[uint32](1000000), sqThenToString))
	assert.Equal(t, option.None[string](), option.AndThen(option.None[uint32](), sqThenToString))
}

func TestOption_Contains(t *testing.T) {
	{
		var x = option.Some(2)
		assert.Equal(t, true, option.Contains(x, 2))
	}
	{
		var x = option.Some(3)
		assert.Equal(t, false, option.Contains(x, 2))
	}
	{
		var x = option.None[int]()
		assert.Equal(t, false, option.Contains(x, 2))
	}
}

func TestOption_MapOrElse_1(t *testing.T) {
	var k = 21
	{
		var x = option.Some("foo")
		assert.Equal(t, 3, option.MapOrElse(x, func() int { return 2 * k }, func(v string) int { return len(v) }))
	}
	{
		var x option.Option[string]
		assert.Equal(t, 42, option.MapOrElse(x, func() int { return 2 * k }, func(v string) int { return len(v) }))
	}
}

func TestOption_MapOr_1(t *testing.T) {
	{
		var x = option.Some("foo")
		assert.Equal(t, 3, option.MapOr(x, 42, func(v string) int { return len(v) }))
	}
	{
		var x option.Option[string]
		assert.Equal(t, 42, option.MapOr(x, 42, func(v string) int { return len(v) }))
	}
}

// TestOption_Map_None tests Map with None option (covers opt/option.go:60-64)
func TestOption_Map_None(t *testing.T) {
	var x option.Option[int]
	result := option.Map(x, func(v int) string {
		return strconv.Itoa(v)
	})
	assert.True(t, result.IsNone())
}

func TestOption_Map_1(t *testing.T) {
	var maybeSomeString = option.Some("Hello, World!")
	var maybeSomeLen = option.Map(maybeSomeString, func(s string) int { return len(s) })
	assert.Equal(t, maybeSomeLen, option.Some(13))
}

func TestOption_Unzip(t *testing.T) {
	var x = option.Some[pair.Pair[int, string]](pair.Pair[int, string]{A: 1, B: "hi"})
	var y = option.None[pair.Pair[int, string]]()
	assert.Equal(t, option.Unzip(x), pair.Pair[option.Option[int], option.Option[string]]{A: option.Some[int](1), B: option.Some[string]("hi")})
	assert.Equal(t, option.Unzip(y), pair.Pair[option.Option[int], option.Option[string]]{A: option.None[int](), B: option.None[string]()})
}

func TestOption_Zip(t *testing.T) {
	var x = option.Some[byte](1)
	var y = option.Some("hi")
	var z = option.None[byte]()
	assert.Equal(t, option.Some(pair.Pair[byte, string]{1, "hi"}), option.Zip(x, y))
	assert.Equal(t, option.None[pair.Pair[byte, byte]](), option.Zip(x, z))
}

func TestOption_ZipWith(t *testing.T) {
	type Point struct {
		x float64
		y float64
	}
	var newPoint = func(x float64, y float64) Point {
		return Point{x, y}
	}
	var x = option.Some(17.5)
	var y = option.Some(42.7)
	assert.Equal(t, option.ZipWith(x, y, newPoint), option.Some(Point{x: 17.5, y: 42.7}))
	assert.Equal(t, option.ZipWith(x, option.None[float64](), newPoint), option.None[Point]())
}

func TestSafeAssert(t *testing.T) {
	// Test with valid type assertion
	opt1 := option.Some(42)
	result1 := option.SafeAssert[int, int](opt1)
	assert.True(t, result1.IsOk())
	assert.True(t, result1.Unwrap().IsSome())
	assert.Equal(t, 42, result1.Unwrap().Unwrap())

	// Test with invalid type assertion
	opt2 := option.Some(42)
	result2 := option.SafeAssert[int, string](opt2)
	assert.True(t, result2.IsErr())

	// Test with None
	opt3 := option.None[int]()
	result3 := option.SafeAssert[int, string](opt3)
	assert.True(t, result3.IsOk())
	assert.True(t, result3.Unwrap().IsNone())
}

func TestXSafeAssert(t *testing.T) {
	// Test with valid type assertion
	opt1 := option.Some[any](42)
	result1 := option.XSafeAssert[int](opt1)
	assert.True(t, result1.IsOk())
	assert.True(t, result1.Unwrap().IsSome())
	assert.Equal(t, 42, result1.Unwrap().Unwrap())

	// Test with invalid type assertion
	opt2 := option.Some[any](42)
	result2 := option.XSafeAssert[string](opt2)
	assert.True(t, result2.IsErr())

	// Test with None
	opt3 := option.None[any]()
	result3 := option.XSafeAssert[int](opt3)
	assert.True(t, result3.IsOk())
	assert.True(t, result3.Unwrap().IsNone())
}

func TestFuzzyAssert(t *testing.T) {
	// Test with valid type assertion
	opt1 := option.Some(42)
	result1 := option.FuzzyAssert[int, int](opt1)
	assert.True(t, result1.IsSome())
	assert.Equal(t, 42, result1.Unwrap())

	// Test with invalid type assertion
	opt2 := option.Some(42)
	result2 := option.FuzzyAssert[int, string](opt2)
	assert.True(t, result2.IsNone())

	// Test with None
	opt3 := option.None[int]()
	result3 := option.FuzzyAssert[int, string](opt3)
	assert.True(t, result3.IsNone())
}

func TestXFuzzyAssert(t *testing.T) {
	// Test with valid type assertion
	opt1 := option.Some[any](42)
	result1 := option.XFuzzyAssert[int](opt1)
	assert.True(t, result1.IsSome())
	assert.Equal(t, 42, result1.Unwrap())

	// Test with invalid type assertion
	opt2 := option.Some[any](42)
	result2 := option.XFuzzyAssert[string](opt2)
	assert.True(t, result2.IsNone())

	// Test with None
	opt3 := option.None[any]()
	result3 := option.XFuzzyAssert[int](opt3)
	assert.True(t, result3.IsNone())
}

func TestAnd(t *testing.T) {
	// Test with Some and Some
	opt1 := option.Some(2)
	opt2 := option.Some(3)
	result1 := option.And(opt1, opt2)
	assert.True(t, result1.IsSome())
	assert.Equal(t, 3, result1.Unwrap())

	// Test with Some and None
	opt3 := option.Some(2)
	opt4 := option.None[int]()
	result2 := option.And(opt3, opt4)
	assert.True(t, result2.IsNone())

	// Test with None and Some
	opt5 := option.None[int]()
	opt6 := option.Some(3)
	result3 := option.And(opt5, opt6)
	assert.True(t, result3.IsNone())

	// Test with None and None
	opt7 := option.None[int]()
	opt8 := option.None[int]()
	result4 := option.And(opt7, opt8)
	assert.True(t, result4.IsNone())
}

func TestSafeAssert_TypeAssertionError(t *testing.T) {
	// Test SafeAssert with type assertion error message
	opt1 := option.Some(42)
	result1 := option.SafeAssert[int, string](opt1)
	assert.True(t, result1.IsErr())
	assert.Contains(t, result1.Err().Error(), "type assert error")
}

func TestXSafeAssert_TypeAssertionError(t *testing.T) {
	// Test XSafeAssert with type assertion error message
	opt1 := option.Some[any](42)
	result1 := option.XSafeAssert[string](opt1)
	assert.True(t, result1.IsErr())
	assert.Contains(t, result1.Err().Error(), "type assert error")
}

func TestOption_ZeroOpt(t *testing.T) {
	// Test ZeroOpt with non-zero value
	opt1 := option.ZeroOpt(42)
	assert.True(t, opt1.IsSome())
	assert.Equal(t, 42, opt1.Unwrap())

	// Test ZeroOpt with zero value
	opt2 := option.ZeroOpt(0)
	assert.True(t, opt2.IsNone())

	// Test ZeroOpt with zero string
	opt3 := option.ZeroOpt("")
	assert.True(t, opt3.IsNone())

	// Test ZeroOpt with non-zero string
	opt4 := option.ZeroOpt("test")
	assert.True(t, opt4.IsSome())
	assert.Equal(t, "test", opt4.Unwrap())
}

func TestOption_RetOpt(t *testing.T) {
	// Test RetOpt with err == nil
	opt1 := option.RetOpt(42, nil)
	assert.True(t, opt1.IsSome())
	assert.Equal(t, 42, opt1.Unwrap())

	// Test RetOpt with err != nil
	err := errors.New("test error")
	opt2 := option.RetOpt(42, err)
	assert.True(t, opt2.IsNone())
}

func TestOption_RetAnyOpt(t *testing.T) {
	// Test RetAnyOpt with err == nil and value != nil
	opt1 := option.RetAnyOpt[int](42, nil)
	assert.True(t, opt1.IsSome())
	assert.Equal(t, 42, opt1.Unwrap())

	// Test RetAnyOpt with err != nil
	err := errors.New("test error")
	opt2 := option.RetAnyOpt[int](42, err)
	assert.True(t, opt2.IsNone())

	// Test RetAnyOpt with nil value
	opt3 := option.RetAnyOpt[*int](nil, nil)
	assert.True(t, opt3.IsNone())
}

func TestOption_PtrOpt(t *testing.T) {
	// Test PtrOpt with nil pointer
	var nilPtr *int
	opt1 := option.PtrOpt(nilPtr)
	assert.True(t, opt1.IsNone())

	// Test PtrOpt with non-nil pointer
	val := 42
	opt2 := option.PtrOpt(&val)
	assert.True(t, opt2.IsSome())
	assert.Equal(t, &val, opt2.Unwrap())
}

func TestOption_ElemOpt(t *testing.T) {
	// Test ElemOpt with nil pointer
	var nilPtr *int
	opt1 := option.ElemOpt(nilPtr)
	assert.True(t, opt1.IsNone())

	// Test ElemOpt with non-nil pointer
	val := 42
	opt2 := option.ElemOpt(&val)
	assert.True(t, opt2.IsSome())
	assert.Equal(t, 42, opt2.Unwrap())
}

func TestOption_BoolOpt(t *testing.T) {
	// Test BoolOpt with ok=true
	opt1 := option.BoolOpt(42, true)
	assert.True(t, opt1.IsSome())
	assert.Equal(t, 42, opt1.Unwrap())

	// Test BoolOpt with ok=false
	opt2 := option.BoolOpt(42, false)
	assert.True(t, opt2.IsNone())
}

func TestOption_AssertOpt(t *testing.T) {
	// Test AssertOpt with valid type
	opt1 := option.AssertOpt[int](42)
	assert.True(t, opt1.IsSome())
	assert.Equal(t, 42, opt1.Unwrap())

	// Test AssertOpt with invalid type
	opt2 := option.AssertOpt[int]("string")
	assert.True(t, opt2.IsNone())
}

func TestOption_BoolAssertOpt(t *testing.T) {
	// Test BoolAssertOpt with ok=true and valid type
	opt1 := option.BoolAssertOpt[int](42, true)
	assert.True(t, opt1.IsSome())
	assert.Equal(t, 42, opt1.Unwrap())

	// Test BoolAssertOpt with ok=false
	opt2 := option.BoolAssertOpt[int](42, false)
	assert.True(t, opt2.IsNone())

	// Test BoolAssertOpt with ok=true but invalid type
	opt3 := option.BoolAssertOpt[int]("string", true)
	assert.True(t, opt3.IsNone())
}
