package gust_test

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/andeya/gust"
	"github.com/andeya/gust/opt"
	"github.com/andeya/gust/valconv"
	"github.com/stretchr/testify/assert"
)

func ExampleOption() {
	type A struct {
		X int
	}
	var a = gust.Some(A{X: 1})
	fmt.Println(a.IsSome(), a.IsNone())

	var b = gust.None[A]()
	fmt.Println(b.IsSome(), b.IsNone())

	var x = b.UnwrapOr(A{X: 2})
	fmt.Println(x)

	var c *A
	fmt.Println(gust.PtrOpt(c).IsNone())
	fmt.Println(gust.ElemOpt(c).IsNone())
	c = new(A)
	fmt.Println(gust.PtrOpt(c).IsNone())
	fmt.Println(gust.ElemOpt(c).IsNone())

	type B struct {
		Y string
	}
	var d = opt.Map(a, func(t A) B {
		return B{
			Y: strconv.Itoa(t.X),
		}
	})
	fmt.Println(d)

	// Output:
	// true false
	// false true
	// {2}
	// true
	// true
	// false
	// false
	// Some({1})
}

func ExampleOption_Inspect() {
	// prints "got: 3"
	_ = gust.Some(3).Inspect(func(x int) {
		fmt.Println("got:", x)
	})

	// prints nothing
	_ = gust.None[int]().Inspect(func(x int) {
		fmt.Println("got:", x)
	})

	// Output:
	// got: 3
}

func TestOption(t *testing.T) {
	var divide = func(numerator, denominator float64) gust.Option[float64] {
		if denominator == 0.0 {
			return gust.None[float64]()
		}
		return gust.Some(numerator / denominator)
	}
	// The return value of the function is an option
	divide(2.0, 3.0).
		Inspect(func(x float64) {
			// Pattern match to retrieve the value
			t.Log("Result:", x)
		}).
		InspectNone(func() {
			t.Log("Cannot divide by 0")
		})
}

func TestOption_AssertOpt(t *testing.T) {
	opt := gust.AssertOpt[int](1)
	assert.Equal(t, gust.Some(1), opt)
	opt2 := gust.AssertOpt[int]("")
	assert.Equal(t, gust.None[int](), opt2)
}

func TestBoolOpt(t *testing.T) {
	// Test with ok=true
	opt1 := gust.BoolOpt(42, true)
	assert.True(t, opt1.IsSome())
	assert.Equal(t, 42, opt1.Unwrap())

	// Test with ok=false
	opt2 := gust.BoolOpt(42, false)
	assert.True(t, opt2.IsNone())
}

func TestBoolAssertOpt(t *testing.T) {
	// Test with ok=true and valid type
	opt1 := gust.BoolAssertOpt[int](42, true)
	assert.True(t, opt1.IsSome())
	assert.Equal(t, 42, opt1.Unwrap())

	// Test with ok=false
	opt2 := gust.BoolAssertOpt[int](42, false)
	assert.True(t, opt2.IsNone())

	// Test with ok=true but invalid type
	opt3 := gust.BoolAssertOpt[int]("string", true)
	assert.True(t, opt3.IsNone())
}

func TestZeroOpt(t *testing.T) {
	// Test with non-zero value
	opt1 := gust.ZeroOpt(42)
	assert.True(t, opt1.IsSome())
	assert.Equal(t, 42, opt1.Unwrap())

	// Test with zero value
	opt2 := gust.ZeroOpt(0)
	assert.True(t, opt2.IsNone())

	// Test with zero string
	opt3 := gust.ZeroOpt("")
	assert.True(t, opt3.IsNone())

	// Test with non-zero string
	opt4 := gust.ZeroOpt("test")
	assert.True(t, opt4.IsSome())
	assert.Equal(t, "test", opt4.Unwrap())
}

func TestOptionJSON(t *testing.T) {
	var r = gust.None[any]()
	var b, err = json.Marshal(r)
	a, _ := json.Marshal(nil)
	assert.Equal(t, a, b)
	assert.NoError(t, err)
	type T struct {
		Name string
	}
	var r2 = gust.Some(T{Name: "andeya"})
	var b2, err2 = json.Marshal(r2)
	assert.NoError(t, err2)
	assert.Equal(t, `{"Name":"andeya"}`, string(b2))

	var r3 gust.Option[T]
	var err3 = json.Unmarshal(b2, &r3)
	assert.NoError(t, err3)
	assert.Equal(t, r2, r3)

	var r4 gust.Option[T]
	var err4 = json.Unmarshal([]byte("0"), &r4)
	assert.True(t, r4.IsNone())
	assert.Equal(t, "json: cannot unmarshal number into Go value of type gust_test.T", err4.Error())
}

func TestOptionJSON2(t *testing.T) {
	type T struct {
		Name gust.Option[string]
	}
	var r = T{Name: gust.Some("andeya")}
	var b, err = json.Marshal(r)
	assert.NoError(t, err)
	assert.Equal(t, `{"Name":"andeya"}`, string(b))
	var r2 T
	err2 := json.Unmarshal(b, &r2)
	assert.NoError(t, err2)
	assert.Equal(t, r, r2)

	var r3 = T{Name: gust.Some("")}
	var b3, err3 = json.Marshal(r3)
	assert.NoError(t, err3)
	assert.Equal(t, `{"Name":""}`, string(b3))
	var r4 T
	err4 := json.Unmarshal(b3, &r4)
	assert.NoError(t, err4)
	assert.Equal(t, r3, r4)

	var r5 = T{Name: gust.None[string]()}
	var b5, err5 = json.Marshal(r5)
	assert.NoError(t, err5)
	assert.Equal(t, `{"Name":null}`, string(b5))
	var r6 T
	err6 := json.Unmarshal(b5, &r6)
	assert.NoError(t, err6)
	assert.Equal(t, r5, r6)
}

func TestOption_And_1(t *testing.T) {
	{
		var x = gust.Some[uint32](2)
		var y gust.Option[uint32]
		assert.Equal(t, gust.None[uint32](), x.And(y))
	}
	{
		var x gust.Option[uint32]
		var y = gust.Some[uint32](3)
		assert.Equal(t, gust.None[uint32](), x.And(y))
	}
	{
		var x = gust.Some[uint32](2)
		var y = gust.Some[uint32](3)
		assert.Equal(t, gust.Some[uint32](3), x.And(y))
	}
	{
		var x gust.Option[uint32]
		var y gust.Option[uint32]
		assert.Equal(t, gust.None[uint32](), x.And(y))
	}
}

func TestOption_And_2(t *testing.T) {
	{
		var x = gust.Some[uint32](2)
		var y gust.Option[string]
		assert.Equal(t, gust.None[string]().ToX(), x.XAnd(y.ToX()))
	}
	{
		var x gust.Option[uint32]
		var y = gust.Some[string]("foo")
		assert.Equal(t, gust.None[string]().ToX(), x.XAnd(y.ToX()))
	}
	{
		var x = gust.Some[uint32](2)
		var y = gust.Some[string]("foo")
		assert.Equal(t, gust.Some[string]("foo").ToX(), x.XAnd(y.ToX()))
	}
	{
		var x gust.Option[uint32]
		var y gust.Option[string]
		assert.Equal(t, gust.None[string]().ToX(), x.XAnd(y.ToX()))
	}
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

func TestOption_Expect(t *testing.T) {
	{
		var x = gust.Some("value")
		assert.Equal(t, "value", x.Expect("fruits are healthy"))
	}
	defer func() {
		assert.Equal(t, gust.ToErrBox("fruits are healthy 1"), recover())
		defer func() {
			assert.Equal(t, gust.ToErrBox("fruits are healthy 2"), recover())
		}()
		var x gust.Option[string]
		x.Expect("fruits are healthy 2") // panics with `fruits are healthy 2`
	}()
	var x = gust.None[string]()
	x.Expect("fruits are healthy 1") // panics with `fruits are healthy 1`
}

func TestOption_Filter(t *testing.T) {
	var isEven = func(n int32) bool {
		return n%2 == 0
	}
	assert.Equal(t, gust.None[int32](), gust.None[int32]().Filter(isEven))
	assert.Equal(t, gust.None[int32](), gust.Some[int32](3).Filter(isEven))
	assert.Equal(t, gust.Some[int32](4), gust.Some[int32](4).Filter(isEven))
}

func TestOption_GetOrInsertDefault(t *testing.T) {
	none := gust.None[int]()
	assert.Equal(t, valconv.Ref(0), none.GetOrInsertDefault())
	assert.Equal(t, valconv.Ref(0), none.AsPtr())
	some := gust.Some[int](1)
	assert.Equal(t, valconv.Ref(1), some.GetOrInsertDefault())
	assert.Equal(t, valconv.Ref(1), some.AsPtr())
}

func TestOption_GetOrInsert(t *testing.T) {
	var x = gust.None[int]()
	var y = x.GetOrInsert(5)
	assert.Equal(t, 5, *y)
	*y = 7
	assert.Equal(t, 7, x.Unwrap())
}

func TestOption_GetOrInsertWith(t *testing.T) {
	var x = gust.None[int]()
	var y = x.GetOrInsertWith(func() int { return 5 })
	assert.Equal(t, 5, *y)
	assert.Equal(t, 5, x.Unwrap())
	*y = 7
	assert.Equal(t, 7, x.Unwrap())

	var x2 = gust.None[int]()
	var y2 = x2.GetOrInsertWith(nil)
	assert.Equal(t, 0, *y2)
	assert.Equal(t, 0, x2.Unwrap())
	*y2 = 7
	assert.Equal(t, 7, x2.Unwrap())
}

func TestOption_Insert(t *testing.T) {
	var opt = gust.None[int]()
	var val = opt.Insert(1)
	assert.Equal(t, 1, *val)
	assert.Equal(t, 1, opt.Unwrap())
	val = opt.Insert(2)
	assert.Equal(t, 2, *val)
	*val = 3
	assert.Equal(t, 3, opt.Unwrap())
}

func TestOption_IsNone(t *testing.T) {
	{
		var x = gust.Some[uint32](2)
		assert.False(t, x.IsNone())
	}
	{
		var x = gust.None[uint32]()
		assert.True(t, x.IsNone())
	}
}

func TestOption_IsSomeAnd(t *testing.T) {
	{
		var x = gust.Some[uint32](2)
		assert.True(t, x.IsSomeAnd(func(v uint32) bool { return v > 1 }))
	}
	{
		var x = gust.Some[uint32](0)
		assert.False(t, x.IsSomeAnd(func(v uint32) bool { return v > 1 }))
	}
	{
		var x = gust.None[uint32]()
		assert.False(t, x.IsSomeAnd(func(v uint32) bool { return v > 1 }))
	}
}

func TestOption_IsSome(t *testing.T) {
	{
		var x = gust.Some[uint32](2)
		assert.True(t, x.IsSome())
	}
	{
		var x = gust.None[uint32]()
		assert.False(t, x.IsSome())
	}
}

func TestOption_MapOrElse_2(t *testing.T) {
	var k = 21
	{
		var x = gust.Some("foo")
		assert.Equal(t, 3, x.XMapOrElse(func() any { return 2 * k }, func(v string) any { return len(v) }))
	}
	{
		var x gust.Option[string]
		assert.Equal(t, 42, x.XMapOrElse(func() any { return 2 * k }, func(v string) any { return len(v) }))
	}
}

func TestOption_MapOr_2(t *testing.T) {
	{
		var x = gust.Some("foo")
		assert.Equal(t, 3, x.XMapOr(42, func(v string) any { return len(v) }))
	}
	{
		var x gust.Option[string]
		assert.Equal(t, 42, x.XMapOr(42, func(v string) any { return len(v) }))
	}
}

func TestOption_Map_2(t *testing.T) {
	var maybeSomeString = gust.Some("Hello, World!")
	var maybeSomeLen = maybeSomeString.XMap(func(s string) any { return len(s) })
	assert.Equal(t, maybeSomeLen, gust.Some(13).ToX())
}

func TestOption_OkOrElse(t *testing.T) {
	{
		var x = gust.Some("foo")
		assert.Equal(t, gust.Ok("foo"), x.OkOrElse(func() any { return 0 }))
	}
	{
		var x gust.Option[string]
		assert.Equal(t, gust.Err[string](0), x.OkOrElse(func() any { return 0 }))
	}
}

func TestOption_OkOr(t *testing.T) {
	{
		var x = gust.Some("foo")
		assert.Equal(t, gust.Ok("foo"), x.OkOr(0))
	}
	{
		var x gust.Option[string]
		assert.Equal(t, gust.Err[string](0), x.OkOr(0))
	}
}

func TestOption_OrElse(t *testing.T) {
	var nobody = func() gust.Option[string] { return gust.None[string]() }
	var vikings = func() gust.Option[string] { return gust.Some("vikings") }
	assert.Equal(t, gust.Some("barbarians"), gust.Some("barbarians").OrElse(vikings))
	assert.Equal(t, gust.Some("vikings"), gust.None[string]().OrElse(vikings))
	assert.Equal(t, gust.None[string](), gust.None[string]().OrElse(nobody))
}

func TestOption_Or(t *testing.T) {
	{
		var x = gust.Some[uint32](2)
		var y gust.Option[uint32]
		assert.Equal(t, gust.Some[uint32](2), x.Or(y))
	}
	{
		var x gust.Option[uint32]
		var y = gust.Some[uint32](100)
		assert.Equal(t, gust.Some[uint32](100), x.Or(y))
	}
	{
		var x = gust.Some[uint32](2)
		var y = gust.Some[uint32](100)
		assert.Equal(t, gust.Some[uint32](2), x.Or(y))
	}
	{
		var x gust.Option[uint32]
		var y = gust.None[uint32]()
		assert.Equal(t, gust.None[uint32](), x.Or(y))
	}
}

func TestOption_Replace(t *testing.T) {
	{
		var x = gust.Some(2)
		var old = x.Replace(5)
		assert.Equal(t, gust.Some(5), x)
		assert.Equal(t, gust.Some(2), old)
	}
	{
		var x = gust.None[int]()
		var old = x.Replace(3)
		assert.Equal(t, gust.Some(3), x)
		assert.Equal(t, gust.None[int](), old)
	}
}

func TestOption_Take(t *testing.T) {
	{
		var x = gust.Some(2)
		var y = x.Take()
		assert.True(t, x.IsNone())
		assert.Equal(t, gust.Some(2), y)
		a, ok := x.Split()
		assert.False(t, ok)
		assert.Equal(t, 0, a)
		b, ok2 := y.Split()
		assert.True(t, ok2)
		assert.Equal(t, 2, b)
	}
	{
		var x gust.Option[int] = gust.None[int]()
		var y = x.Take()
		assert.True(t, x.IsNone())
		assert.True(t, y.IsNone())
	}
}

func TestOption_UnwrapOrDefault(t *testing.T) {
	assert.Equal(t, "car", gust.Some("car").UnwrapOrDefault())
	assert.Equal(t, "", gust.None[string]().UnwrapOrDefault())
	assert.Equal(t, time.Time{}, gust.None[time.Time]().UnwrapOrDefault())
	assert.Equal(t, &time.Time{}, gust.None[*time.Time]().UnwrapOrDefault())
	assert.Equal(t, valconv.Ref(&time.Time{}), gust.None[**time.Time]().UnwrapOrDefault())
}

func TestOption_UnwrapOrElse(t *testing.T) {
	var k = 10
	assert.Equal(t, 4, gust.Some(4).UnwrapOrElse(func() int { return 2 * k }))
	assert.Equal(t, 20, gust.None[int]().UnwrapOrElse(func() int { return 2 * k }))
}

func TestOption_UnwrapOr(t *testing.T) {
	assert.Equal(t, "car", gust.Some("car").UnwrapOr("bike"))
	assert.Equal(t, "bike", gust.None[string]().UnwrapOr("bike"))
}

func TestOption_Unwrap(t *testing.T) {
	{
		var x = gust.Some("air")
		assert.Equal(t, "air", x.Unwrap())
	}
	defer func() {
		assert.Equal(t, gust.ToErrBox("call Option[string].Unwrap() on none"), recover())
	}()
	var x = gust.None[string]()
	x.Unwrap()
}

func TestOption_Xor(t *testing.T) {
	{
		var x = gust.Some[uint32](2)
		var y gust.Option[uint32]
		assert.Equal(t, gust.Some[uint32](2), x.Xor(y))
	}
	{
		var x gust.Option[uint32]
		var y = gust.Some[uint32](100)
		assert.Equal(t, gust.Some[uint32](100), x.Xor(y))
	}
	{
		var x = gust.Some[uint32](2)
		var y = gust.Some[uint32](100)
		assert.Equal(t, gust.None[uint32](), x.Xor(y))
	}
	{
		var x gust.Option[uint32]
		var y = gust.None[uint32]()
		assert.Equal(t, gust.None[uint32](), x.Xor(y))
	}
}

func TestOption_UnwrapUnchecked(t *testing.T) {
	{
		var x = gust.Some("foo")
		assert.Equal(t, "foo", x.UnwrapUnchecked())
	}
	{
		var x gust.Option[string]
		assert.Equal(t, "", x.UnwrapUnchecked())
	}
	{
		var r gust.Result[string]
		assert.Equal(t, "", r.Ok().UnwrapUnchecked())
	}
	{
		var r = gust.Err[string]("foo")
		assert.Equal(t, "", r.Ok().UnwrapUnchecked())
	}
	{
		var r = gust.Ok[string]("foo")
		assert.Equal(t, "foo", r.Ok().UnwrapUnchecked())
	}
}

func TestOption_UnwrapOrThrow(t *testing.T) {
	// Test with Some value
	var opt1 = gust.Some("value")
	var result gust.Result[string]
	defer gust.CatchResult(&result)
	val := opt1.UnwrapOrThrow("error message")
	assert.Equal(t, "value", val)
	assert.True(t, result.IsOk())

	// Test with None (should panic)
	var opt2 = gust.None[string]()
	var result2 gust.Result[int]
	defer func() {
		assert.True(t, result2.IsErr())
		assert.Equal(t, "error message", result2.Err().Error())
	}()
	defer gust.CatchResult(&result2)
	_ = opt2.UnwrapOrThrow("error message")
}
