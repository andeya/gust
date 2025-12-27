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

func TestRetOpt(t *testing.T) {
	// Test with err == nil, should return Some(value)
	opt1 := gust.RetOpt(42, nil)
	assert.True(t, opt1.IsSome())
	assert.Equal(t, 42, opt1.Unwrap())

	// Test with err != nil, should return None
	err := fmt.Errorf("test error")
	opt2 := gust.RetOpt(42, err)
	assert.True(t, opt2.IsNone())

	// Test with string value and err == nil
	opt3 := gust.RetOpt("hello", nil)
	assert.True(t, opt3.IsSome())
	assert.Equal(t, "hello", opt3.Unwrap())

	// Test with string value and err != nil
	opt4 := gust.RetOpt("hello", err)
	assert.True(t, opt4.IsNone())

	// Test with zero value and err == nil
	opt5 := gust.RetOpt(0, nil)
	assert.True(t, opt5.IsSome())
	assert.Equal(t, 0, opt5.Unwrap())

	// Test with zero value and err != nil
	opt6 := gust.RetOpt(0, err)
	assert.True(t, opt6.IsNone())
}

func TestRetAnyOpt(t *testing.T) {
	err := fmt.Errorf("test error")

	// Test with err == nil and value != nil, should return Some(value)
	opt1 := gust.RetAnyOpt[int](42, nil)
	assert.True(t, opt1.IsSome())
	assert.Equal(t, 42, opt1.Unwrap())

	// Test with err != nil, should return None (even if value != nil)
	opt2 := gust.RetAnyOpt[int](42, err)
	assert.True(t, opt2.IsNone())

	// Test with string value and err == nil
	opt3 := gust.RetAnyOpt[string]("hello", nil)
	assert.True(t, opt3.IsSome())
	assert.Equal(t, "hello", opt3.Unwrap())

	// Test with string value and err != nil
	opt4 := gust.RetAnyOpt[string]("hello", err)
	assert.True(t, opt4.IsNone())

	// Test with different types - float64 and err == nil
	opt5 := gust.RetAnyOpt[float64](3.14, nil)
	assert.True(t, opt5.IsSome())
	assert.Equal(t, 3.14, opt5.Unwrap())

	// Test with different types - float64 and err != nil
	opt6 := gust.RetAnyOpt[float64](3.14, err)
	assert.True(t, opt6.IsNone())

	// Test with nil value and err == nil, should return None (value == nil)
	opt7 := gust.RetAnyOpt[*int](nil, nil)
	assert.True(t, opt7.IsNone())

	// Test with nil value and err != nil, should return None (both conditions)
	opt8 := gust.RetAnyOpt[*int](nil, err)
	assert.True(t, opt8.IsNone())

	// Test with nil string and err == nil
	opt9 := gust.RetAnyOpt[*string](nil, nil)
	assert.True(t, opt9.IsNone())

	// Test with nil string and err != nil
	opt10 := gust.RetAnyOpt[*string](nil, err)
	assert.True(t, opt10.IsNone())

	// Test with struct value and err == nil
	type TestStruct struct {
		X int
		Y string
	}
	testStruct := TestStruct{X: 1, Y: "test"}
	opt11 := gust.RetAnyOpt[TestStruct](testStruct, nil)
	assert.True(t, opt11.IsSome())
	assert.Equal(t, testStruct, opt11.Unwrap())

	// Test with struct value and err != nil
	opt12 := gust.RetAnyOpt[TestStruct](testStruct, err)
	assert.True(t, opt12.IsNone())

	// Test with nil struct pointer and err == nil
	opt13 := gust.RetAnyOpt[*TestStruct](nil, nil)
	assert.True(t, opt13.IsNone())

	// Test with zero value (non-nil) and err == nil
	opt14 := gust.RetAnyOpt[int](0, nil)
	assert.True(t, opt14.IsSome())
	assert.Equal(t, 0, opt14.Unwrap())

	// Test with empty string (non-nil) and err == nil
	opt15 := gust.RetAnyOpt[string]("", nil)
	assert.True(t, opt15.IsSome())
	assert.Equal(t, "", opt15.Unwrap())
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
		assert.Equal(t, gust.BoxErr("fruits are healthy 1"), recover())
		defer func() {
			assert.Equal(t, gust.BoxErr("fruits are healthy 2"), recover())
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
		assert.Equal(t, gust.BoxErr("call Option[string].Unwrap() on none"), recover())
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
	defer result.Catch()
	val := opt1.UnwrapOrThrow("error message")
	assert.Equal(t, "value", val)
	// When UnwrapOrThrow succeeds (no panic), Catch() doesn't modify result,
	// so result remains in zero state (IsOk() returns false).
	// This is expected behavior - Catch() only handles panics, not success cases.
	assert.False(t, result.IsOk())

	// Test with None (should panic)
	var opt2 = gust.None[string]()
	var result2 gust.Result[int]
	defer func() {
		assert.True(t, result2.IsErr())
		assert.Equal(t, "error message", result2.Err().Error())
	}()
	defer result2.Catch()
	_ = opt2.UnwrapOrThrow("error message")
}

func TestOption_XOkOr(t *testing.T) {
	{
		var x = gust.Some("foo")
		assert.Equal(t, gust.Ok[any]("foo"), x.XOkOr(0))
	}
	{
		var x gust.Option[string]
		assert.Equal(t, gust.Err[any](0), x.XOkOr(0))
	}
}

func TestOption_XOkOrElse(t *testing.T) {
	{
		var x = gust.Some("foo")
		assert.Equal(t, gust.Ok[any]("foo"), x.XOkOrElse(func() any { return 0 }))
	}
	{
		var x gust.Option[string]
		assert.Equal(t, gust.Err[any](0), x.XOkOrElse(func() any { return 0 }))
	}
}

func TestOption_ToResult(t *testing.T) {
	{
		var x = gust.Some("foo")
		result := x.ToResult()
		assert.True(t, result.IsErr())
		assert.Equal(t, "foo", result.ErrVal())
	}
	{
		var x gust.Option[string]
		result := x.ToResult()
		assert.False(t, result.IsErr())
	}
}

// TestOption_InspectNone tests InspectNone method (covers option.go:285-289)
func TestOption_InspectNone(t *testing.T) {
	// Test with None
	{
		var opt gust.Option[string]
		called := false
		result := opt.InspectNone(func() {
			called = true
		})
		assert.True(t, called)
		assert.True(t, result.IsNone())
	}
	// Test with Some (should not call the function)
	{
		opt := gust.Some("value")
		called := false
		result := opt.InspectNone(func() {
			called = true
		})
		assert.False(t, called)
		assert.True(t, result.IsSome())
		assert.Equal(t, "value", result.UnwrapUnchecked())
	}
}

// TestOption_String tests Option.String method (covers option.go:148-152)
func TestOption_String(t *testing.T) {
	none := gust.None[int]()
	assert.Equal(t, "None", none.String())

	some := gust.Some(42)
	assert.Contains(t, some.String(), "Some")
}

func TestOption_Iterator(t *testing.T) {
	// Test Next
	{
		var x = gust.Some("foo")
		opt := x.Next()
		assert.Equal(t, gust.Some("foo"), opt)
		assert.True(t, x.IsNone()) // Should consume the value
	}
	{
		var x gust.Option[string]
		opt := x.Next()
		assert.True(t, opt.IsNone())
	}
	{
		var nilOpt *gust.Option[string]
		opt := nilOpt.Next()
		assert.True(t, opt.IsNone())
	}

	// Test NextBack
	{
		var x = gust.Some("bar")
		opt := x.NextBack()
		assert.Equal(t, gust.Some("bar"), opt)
		assert.True(t, x.IsNone())
	}

	// Test Remaining
	{
		var x = gust.Some("baz")
		assert.Equal(t, uint(1), x.Remaining())
		x.Next()
		assert.Equal(t, uint(0), x.Remaining())
	}
	{
		var x gust.Option[string]
		assert.Equal(t, uint(0), x.Remaining())
	}
	{
		var nilOpt *gust.Option[string]
		assert.Equal(t, uint(0), nilOpt.Remaining())
	}
}

func TestOption_InsertNil(t *testing.T) {
	var nilOpt *gust.Option[int]
	val := nilOpt.Insert(42)
	assert.Nil(t, val)
}

func TestOption_GetOrInsertNil(t *testing.T) {
	var nilOpt *gust.Option[int]
	val := nilOpt.GetOrInsert(42)
	assert.Nil(t, val)
}

func TestOption_GetOrInsertWithNil(t *testing.T) {
	var nilOpt *gust.Option[int]
	val := nilOpt.GetOrInsertWith(func() int { return 42 })
	assert.Nil(t, val)
}

func TestOption_GetOrInsertDefaultNil(t *testing.T) {
	var nilOpt *gust.Option[int]
	val := nilOpt.GetOrInsertDefault()
	assert.Nil(t, val)
}

func TestOption_Ref(t *testing.T) {
	opt := gust.Some(42)
	ref := opt.Ref()
	assert.Equal(t, gust.Some(42), *ref)
	ref.Unwrap() // Should not panic
}

func TestOption_Split(t *testing.T) {
	{
		var x = gust.Some("foo")
		val, ok := x.Split()
		assert.True(t, ok)
		assert.Equal(t, "foo", val)
	}
	{
		var x gust.Option[string]
		val, ok := x.Split()
		assert.False(t, ok)
		assert.Equal(t, "", val)
	}
}

func TestOption_ToX(t *testing.T) {
	{
		var x = gust.Some("foo")
		xOpt := x.ToX()
		assert.True(t, xOpt.IsSome())
		assert.Equal(t, "foo", xOpt.Unwrap())
	}
	{
		var x gust.Option[string]
		xOpt := x.ToX()
		assert.True(t, xOpt.IsNone())
	}
}

func TestOption_UnmarshalJSON_NilReceiver(t *testing.T) {
	// Test UnmarshalJSON with nil receiver
	var nilOpt *gust.Option[int]
	err := nilOpt.UnmarshalJSON([]byte("42"))
	assert.Error(t, err)
	assert.IsType(t, &json.InvalidUnmarshalError{}, err)
}

func TestOption_UnmarshalJSON_ErrorPath(t *testing.T) {
	// Test UnmarshalJSON with invalid JSON (error path)
	var opt gust.Option[int]
	err := opt.UnmarshalJSON([]byte("invalid json"))
	assert.Error(t, err)
	assert.True(t, opt.IsNone()) // Should remain None on error
}

func TestOption_UnmarshalJSON_ValidAfterError(t *testing.T) {
	// Test UnmarshalJSON with error first, then valid JSON
	var opt gust.Option[int]
	// First attempt with invalid JSON
	_ = opt.UnmarshalJSON([]byte("invalid"))
	assert.True(t, opt.IsNone())

	// Then with valid JSON
	err := opt.UnmarshalJSON([]byte("42"))
	assert.NoError(t, err)
	assert.True(t, opt.IsSome())
	assert.Equal(t, 42, opt.Unwrap())
}

func TestOption_UnmarshalJSON_NullString(t *testing.T) {
	// Test UnmarshalJSON with "null" string (unsafe.Pointer path)
	var opt gust.Option[int]
	err := opt.UnmarshalJSON([]byte("null"))
	assert.NoError(t, err)
	assert.True(t, opt.IsNone())
}

func TestOption_UnmarshalJSON_Struct(t *testing.T) {
	type S struct {
		X int
		Y string
	}
	var opt gust.Option[S]
	err := opt.UnmarshalJSON([]byte(`{"X":10,"Y":"test"}`))
	assert.NoError(t, err)
	assert.True(t, opt.IsSome())
	assert.Equal(t, S{X: 10, Y: "test"}, opt.Unwrap())
}

func TestOption_UnmarshalJSON_Array(t *testing.T) {
	var opt gust.Option[[]int]
	err := opt.UnmarshalJSON([]byte("[1,2,3]"))
	assert.NoError(t, err)
	assert.True(t, opt.IsSome())
	assert.Equal(t, []int{1, 2, 3}, opt.Unwrap())
}

func TestOption_UnmarshalJSON_Map(t *testing.T) {
	var opt gust.Option[map[string]int]
	err := opt.UnmarshalJSON([]byte(`{"a":1,"b":2}`))
	assert.NoError(t, err)
	assert.True(t, opt.IsSome())
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, opt.Unwrap())
}

func TestPtrOpt_EdgeCases(t *testing.T) {
	// Test PtrOpt with nil pointer
	var nilPtr *int
	opt1 := gust.PtrOpt(nilPtr)
	assert.True(t, opt1.IsNone())

	// Test PtrOpt with non-nil pointer
	val := 42
	opt2 := gust.PtrOpt(&val)
	assert.True(t, opt2.IsSome())
	assert.Equal(t, &val, opt2.Unwrap())

	// Test PtrOpt with nested pointer
	var nilPtr2 **int
	opt3 := gust.PtrOpt(nilPtr2)
	assert.True(t, opt3.IsNone())

	val2 := 100
	ptr2 := &val2
	opt4 := gust.PtrOpt(&ptr2)
	assert.True(t, opt4.IsSome())
	assert.Equal(t, &ptr2, opt4.Unwrap())
}

func TestElemOpt_EdgeCases(t *testing.T) {
	// Test ElemOpt with nil pointer
	var nilPtr *int
	opt1 := gust.ElemOpt(nilPtr)
	assert.True(t, opt1.IsNone())

	// Test ElemOpt with non-nil pointer
	val := 42
	opt2 := gust.ElemOpt(&val)
	assert.True(t, opt2.IsSome())
	assert.Equal(t, 42, opt2.Unwrap())

	// Test ElemOpt with struct pointer
	type S struct {
		X int
		Y string
	}
	var nilStruct *S
	opt3 := gust.ElemOpt(nilStruct)
	assert.True(t, opt3.IsNone())

	s := &S{X: 10, Y: "test"}
	opt4 := gust.ElemOpt(s)
	assert.True(t, opt4.IsSome())
	assert.Equal(t, S{X: 10, Y: "test"}, opt4.Unwrap())
}

func TestBoolAssertOpt_EdgeCases(t *testing.T) {
	// Test BoolAssertOpt with ok=false
	opt1 := gust.BoolAssertOpt[int](42, false)
	assert.True(t, opt1.IsNone())

	// Test BoolAssertOpt with ok=true and valid type
	opt2 := gust.BoolAssertOpt[int](42, true)
	assert.True(t, opt2.IsSome())
	assert.Equal(t, 42, opt2.Unwrap())

	// Test BoolAssertOpt with ok=true but invalid type
	opt3 := gust.BoolAssertOpt[int]("string", true)
	assert.True(t, opt3.IsNone())

	// Test BoolAssertOpt with ok=true and valid string type
	opt4 := gust.BoolAssertOpt[string]("test", true)
	assert.True(t, opt4.IsSome())
	assert.Equal(t, "test", opt4.Unwrap())
}

func TestOption_MarshalJSON_EdgeCases(t *testing.T) {
	// Test MarshalJSON with None
	opt1 := gust.None[int]()
	b1, err1 := opt1.MarshalJSON()
	assert.NoError(t, err1)
	assert.Equal(t, []byte("null"), b1)

	// Test MarshalJSON with Some containing zero value
	opt2 := gust.Some(0)
	b2, err2 := opt2.MarshalJSON()
	assert.NoError(t, err2)
	assert.Equal(t, []byte("0"), b2)

	// Test MarshalJSON with Some containing struct
	type S struct {
		X int
		Y string
	}
	opt3 := gust.Some(S{X: 10, Y: "test"})
	b3, err3 := opt3.MarshalJSON()
	assert.NoError(t, err3)
	assert.Contains(t, string(b3), "X")
	assert.Contains(t, string(b3), "Y")
}

func TestOption_AndThen_EdgeCases(t *testing.T) {
	// Test AndThen with None
	opt1 := gust.None[int]()
	opt2 := opt1.AndThen(func(x int) gust.Option[int] {
		return gust.Some(x * 2)
	})
	assert.True(t, opt2.IsNone())

	// Test AndThen with Some that returns None
	opt3 := gust.Some(42)
	opt4 := opt3.AndThen(func(x int) gust.Option[int] {
		return gust.None[int]()
	})
	assert.True(t, opt4.IsNone())

	// Test AndThen with Some that returns Some
	opt5 := gust.Some(21)
	opt6 := opt5.AndThen(func(x int) gust.Option[int] {
		return gust.Some(x * 2)
	})
	assert.True(t, opt6.IsSome())
	assert.Equal(t, 42, opt6.Unwrap())
}

func TestOption_Or_EdgeCases(t *testing.T) {
	// Test Or with None and None
	opt1 := gust.None[int]()
	opt2 := gust.None[int]()
	opt3 := opt1.Or(opt2)
	assert.True(t, opt3.IsNone())

	// Test Or with None and Some
	opt4 := gust.None[int]()
	opt5 := gust.Some(42)
	opt6 := opt4.Or(opt5)
	assert.True(t, opt6.IsSome())
	assert.Equal(t, 42, opt6.Unwrap())

	// Test Or with Some and None
	opt7 := gust.Some(10)
	opt8 := gust.None[int]()
	opt9 := opt7.Or(opt8)
	assert.True(t, opt9.IsSome())
	assert.Equal(t, 10, opt9.Unwrap())

	// Test Or with Some and Some
	opt10 := gust.Some(10)
	opt11 := gust.Some(20)
	opt12 := opt10.Or(opt11)
	assert.True(t, opt12.IsSome())
	assert.Equal(t, 10, opt12.Unwrap())
}

func TestOption_MapOrElse_EdgeCases(t *testing.T) {
	// Test MapOrElse with None
	opt1 := gust.None[int]()
	result1 := opt1.MapOrElse(func() int {
		return 100
	}, func(x int) int {
		return x * 2
	})
	assert.Equal(t, 100, result1)

	// Test MapOrElse with Some
	opt2 := gust.Some(21)
	result2 := opt2.MapOrElse(func() int {
		return 100
	}, func(x int) int {
		return x * 2
	})
	assert.Equal(t, 42, result2)
}

func TestOption_MapOr_EdgeCases(t *testing.T) {
	// Test MapOr with None
	opt1 := gust.None[int]()
	result1 := opt1.MapOr(100, func(x int) int {
		return x * 2
	})
	assert.Equal(t, 100, result1)

	// Test MapOr with Some
	opt2 := gust.Some(21)
	result2 := opt2.MapOr(100, func(x int) int {
		return x * 2
	})
	assert.Equal(t, 42, result2)
}

func TestOption_Map_EdgeCases(t *testing.T) {
	// Test Map with None
	opt1 := gust.None[int]()
	opt2 := opt1.Map(func(x int) int {
		return x * 2
	})
	assert.True(t, opt2.IsNone())

	// Test Map with Some
	opt3 := gust.Some(21)
	opt4 := opt3.Map(func(x int) int {
		return x * 2
	})
	assert.True(t, opt4.IsSome())
	assert.Equal(t, 42, opt4.Unwrap())
}

func TestOption_XMap_EdgeCases(t *testing.T) {
	// Test XMap with None
	opt1 := gust.None[int]()
	opt2 := opt1.XMap(func(x int) any {
		return x * 2
	})
	assert.True(t, opt2.IsNone())

	// Test XMap with Some
	opt3 := gust.Some(21)
	opt4 := opt3.XMap(func(x int) any {
		return x * 2
	})
	assert.True(t, opt4.IsSome())
	assert.Equal(t, 42, opt4.Unwrap())
}

func TestOption_XMapOr_EdgeCases(t *testing.T) {
	// Test XMapOr with None
	opt1 := gust.None[int]()
	result1 := opt1.XMapOr(100, func(x int) any {
		return x * 2
	})
	assert.Equal(t, 100, result1)

	// Test XMapOr with Some
	opt2 := gust.Some(21)
	result2 := opt2.XMapOr(100, func(x int) any {
		return x * 2
	})
	assert.Equal(t, 42, result2)
}

func TestOption_XMapOrElse_EdgeCases(t *testing.T) {
	// Test XMapOrElse with None
	opt1 := gust.None[int]()
	result1 := opt1.XMapOrElse(func() any {
		return 100
	}, func(x int) any {
		return x * 2
	})
	assert.Equal(t, 100, result1)

	// Test XMapOrElse with Some
	opt2 := gust.Some(21)
	result2 := opt2.XMapOrElse(func() any {
		return 100
	}, func(x int) any {
		return x * 2
	})
	assert.Equal(t, 42, result2)
}
