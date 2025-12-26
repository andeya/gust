package gust_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/andeya/gust"
	"github.com/andeya/gust/ret"
	"github.com/andeya/gust/valconv"
	"github.com/stretchr/testify/assert"
)

func ExampleResult_UnwrapOr() {
	const def int = 10

	// before
	i, err := strconv.Atoi("1")
	if err != nil {
		i = def
	}
	fmt.Println(i * 2)

	// now
	fmt.Println(gust.Ret(strconv.Atoi("1")).UnwrapOr(def) * 2)

	// Output:
	// 2
	// 2
}

// [`Result`] comes with some convenience methods that make working with it more succinct.
func TestResult(t *testing.T) {
	var goodResult1 = gust.Ok(10)
	var badResult1 = gust.Err[int](10)

	// The `IsOk` and `IsErr` methods do what they say.
	assert.True(t, goodResult1.IsOk() && !goodResult1.IsErr())
	assert.True(t, badResult1.IsErr() && !badResult1.IsOk())

	// `map` consumes the `Result` and produces another.
	var goodResult2 = goodResult1.Map(func(i int) int { return i + 1 })
	var badResult2 = badResult1.Map(func(i int) int { return i - 1 })

	// Use `AndThen` to continue the computation.
	var goodResult3 = ret.AndThen(goodResult2, func(i int) gust.Result[bool] { return gust.Ok(i == 11) })

	// Use `OrElse` to handle the error.
	var _ = badResult2.OrElse(func(err error) gust.Result[int] {
		fmt.Println(err)
		return gust.Ok(20)
	})

	// Consume the result and return the contents with `Unwrap`.
	var _ = goodResult3.Unwrap()

	r, err := goodResult3.Split()
	assert.NoError(t, err)
	assert.Equal(t, true, r)
}

func TestResult_AssertRet(t *testing.T) {
	r := gust.AssertRet[int](1)
	assert.Equal(t, gust.Ok(1), r)
	r2 := gust.AssertRet[int]("")
	assert.Equal(t, gust.FmtErr[int]("type assert error, got string, want int"), r2)
}

func TestResultJSON(t *testing.T) {
	var r = gust.Err[any](errors.New("err"))
	var b, err = json.Marshal(r)
	assert.Equal(t, "json: error calling MarshalJSON for type gust.Result[interface {}]: err", err.Error())
	assert.Nil(t, b)
	type T struct {
		Name string
	}
	var r2 = gust.Ok(T{Name: "andeya"})
	var b2, err2 = json.Marshal(r2)
	assert.NoError(t, err2)
	assert.Equal(t, `{"Name":"andeya"}`, string(b2))

	var r3 gust.Result[T]
	var err3 = json.Unmarshal(b2, &r3)
	assert.NoError(t, err3)
	assert.Equal(t, r2, r3)

	var r4 gust.Result[T]
	var err4 = json.Unmarshal([]byte("0"), &r4)
	assert.True(t, r4.IsErr())
	assert.Equal(t, "json: cannot unmarshal number into Go value of type gust_test.T", err4.Error())

	var r5 = gust.Ok[gust.Void](nil)
	var b5, err5 = json.Marshal(r5)
	assert.NoError(t, err5)
	assert.Equal(t, `null`, string(b5))
	var r6 gust.Result[gust.Void]
	var err6 = json.Unmarshal(b5, &r6)
	assert.NoError(t, err6)
	assert.Equal(t, r5, r6)
}

func TestResultIsValid(t *testing.T) {
	var r0 *gust.Result[any]
	assert.False(t, r0.IsValid())
	var r1 gust.Result[any]
	assert.False(t, r1.IsValid())
	assert.False(t, (&gust.Result[any]{}).IsValid())
	var r2 = gust.Ok[any](nil)
	assert.True(t, r2.IsValid())
}

func TestResultUnwrapOrThrow_1(t *testing.T) {
	var r gust.Result[string]
	defer func() {
		assert.Equal(t, gust.Err[string]("err"), r)
	}()
	defer gust.CatchResult[string](&r)
	var r1 = gust.Ok(1)
	var v1 = r1.UnwrapOrThrow()
	assert.Equal(t, 1, v1)
	var r2 = gust.Err[int]("err")
	var v2 = r2.UnwrapOrThrow()
	assert.Equal(t, 0, v2)
}

func TestResultUnwrapOrThrow_2(t *testing.T) {
	defer func() {
		assert.Equal(t, "panic text", recover())
	}()
	var r gust.Result[string]
	defer gust.CatchResult[string](&r)
	panic("panic text")
}

func TestResultUnwrapOrThrow_3(t *testing.T) {
	var r gust.Result[string]
	defer func() {
		assert.Equal(t, gust.Err[string]("err"), r)
	}()
	defer r.Catch()
	var r1 = gust.Ok(1)
	var v1 = r1.UnwrapOrThrow()
	assert.Equal(t, 1, v1)
	var r2 = gust.Err[int]("err")
	var v2 = r2.UnwrapOrThrow()
	assert.Equal(t, 0, v2)
}

func TestResultUnwrapOrThrow_4(t *testing.T) {
	var r gust.Result[string]
	defer func() {
		assert.Equal(t, gust.Err[string]("err"), r)
	}()
	defer r.Catch()
	var r1 = gust.Ok(1)
	var v1 = r1.UnwrapOrThrow()
	assert.Equal(t, 1, v1)
	var r2 = gust.Err[int]("err")
	var v2 = r2.UnwrapOrThrow()
	assert.Equal(t, 0, v2)
}

func TestResultUnwrapOrThrow_5(t *testing.T) {
	defer func() {
		assert.Equal(t, "panic text", recover())
	}()
	var r gust.Result[string]
	defer r.Catch()
	panic("panic text")
}

func TestResult_And(t *testing.T) {
	{
		x := gust.Ok(2)
		y := gust.Err[int]("late error")
		assert.Equal(t, gust.Err[int]("late error"), x.And(y))
	}
	{
		x := gust.Err[uint]("early error")
		y := gust.Ok[string]("foo")
		assert.Equal(t, gust.Err[any]("early error"), x.XAnd(y.ToX()))
	}
}

func TestResult_And2(t *testing.T) {
	// Test with Ok result and Ok value
	{
		x := gust.Ok(2)
		result := x.And2(3, nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, 3, result.Unwrap())
	}
	// Test with Ok result and Err value
	{
		x := gust.Ok(2)
		result := x.And2(3, errors.New("late error"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "late error", result.Err().Error())
	}
	// Test with Err result (should return original error)
	{
		x := gust.Err[int]("early error")
		result := x.And2(3, nil)
		assert.True(t, result.IsErr())
		assert.Equal(t, "early error", result.Err().Error())
	}
	// Test with Err result and Err value (should return original error)
	{
		x := gust.Err[int]("early error")
		result := x.And2(3, errors.New("late error"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "early error", result.Err().Error())
	}
}

func TestResult_XAnd2(t *testing.T) {
	// Test with Ok result and Ok value
	{
		x := gust.Ok(2)
		result := x.XAnd2("foo", nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, "foo", result.Unwrap())
	}
	// Test with Ok result and Err value
	{
		x := gust.Ok(2)
		result := x.XAnd2("foo", errors.New("late error"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "late error", result.Err().Error())
	}
	// Test with Err result (should return original error)
	{
		x := gust.Err[int]("early error")
		result := x.XAnd2("foo", nil)
		assert.True(t, result.IsErr())
		assert.Equal(t, "early error", result.Err().Error())
	}
	// Test with Err result and Err value (should return original error)
	{
		x := gust.Err[int]("early error")
		result := x.XAnd2("foo", errors.New("late error"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "early error", result.Err().Error())
	}
}

func ExampleResult_AndThen() {
	var divide = func(i, j float32) gust.Result[float32] {
		if j == 0 {
			return gust.Err[float32]("j can not be 0")
		}
		return gust.Ok(i / j)
	}
	var ret float32 = divide(1, 2).AndThen(func(i float32) gust.Result[float32] {
		return gust.Ok(i * 10)
	}).Unwrap()
	fmt.Println(ret)
	// Output:
	// 5
}

func TestResult_Err(t *testing.T) {
	{
		var x = gust.Ok[int](2)
		assert.Equal(t, error(nil), x.Err())
	}
	{
		var x = gust.Err[int]("some error message")
		assert.Equal(t, "some error message", x.Err().Error())
	}
}

func TestResult_ExpectErr(t *testing.T) {
	defer func() {
		assert.Equal(t, gust.ToErrBox("Testing expect_err: 10"), recover())
	}()
	err := gust.Ok(10).ExpectErr("Testing expect_err")
	assert.NoError(t, err)

	// Test IsErr branch (covers result.go:311-313)
	err2 := gust.Err[int]("test error")
	result := err2.ExpectErr("should not panic")
	assert.NotNil(t, result)
	assert.Equal(t, "test error", result.Error())
}

// TestResult_String tests Result.String method (covers result.go:102-107)
func TestResult_String(t *testing.T) {
	ok := gust.Ok(42)
	assert.Equal(t, "Ok(42)", ok.String())

	err := gust.Err[int]("error message")
	assert.Contains(t, err.String(), "Err")
}

// TestResult_ErrValNil tests Result.ErrVal returning nil (covers result.go:165-169)
func TestResult_ErrValNil(t *testing.T) {
	ok := gust.Ok(42)
	val := ok.ErrVal()
	assert.Nil(t, val)
}

func TestResult_Expect(t *testing.T) {
	defer func() {
		assert.Equal(t, "failed to parse number: strconv.Atoi: parsing \"4x\": invalid syntax", recover().(error).Error())
	}()
	gust.Ret(strconv.Atoi("4x")).
		Expect("failed to parse number")
}

func TestResult_InspectErr(t *testing.T) {
	gust.Ret(strconv.Atoi("4x")).
		InspectErr(func(err error) {
			t.Logf("failed to convert: %v", err)
		})
}

func TestResult_Inspect(t *testing.T) {
	gust.Ret(strconv.Atoi("4")).
		Inspect(func(x int) {
			fmt.Println("original: ", x)
		}).
		Map(func(x int) int {
			return x * 3
		}).
		Expect("failed to parse number")
}

func TestResult_IsErrAnd(t *testing.T) {
	{
		var x = gust.Err[int]("hey")
		assert.True(t, x.IsErrAnd(func(err error) bool { return err.Error() == "hey" }))
	}
	{
		var x = gust.Ok[int](2)
		assert.False(t, x.IsErrAnd(func(err error) bool { return err.Error() == "hey" }))
	}
}

func TestResult_IsErr(t *testing.T) {
	{
		var x = gust.Ok[int](-3)
		assert.False(t, x.IsErr())
	}
	{
		var x = gust.Err[int]("some error message")
		assert.True(t, x.IsErr())
	}
}

func TestResult_IsOkAnd(t *testing.T) {
	{
		var x = gust.Ok[int](2)
		assert.True(t, x.IsOkAnd(func(x int) bool { return x > 1 }))
	}
	{
		var x = gust.Ok[int](0)
		assert.False(t, x.IsOkAnd(func(x int) bool { return x > 1 }))
	}
	{
		var x = gust.Err[int]("hey")
		assert.False(t, x.IsOkAnd(func(x int) bool { return x > 1 }))
	}
}

func TestResult_IsOk(t *testing.T) {
	{
		var x = gust.Ok[int](-3)
		assert.True(t, x.IsOk())
	}
	{
		var x = gust.Err[int]("some error message")
		assert.False(t, x.IsOk())
	}
}

func TestResult_MapErr(t *testing.T) {
	var stringify = func(x error) any { return fmt.Sprintf("error code: %v", x) }
	{
		var x = gust.Ok[uint32](2)
		assert.Equal(t, gust.Ok[uint32](2), x.MapErr(stringify))
	}
	{
		var x = gust.Err[uint32](13)
		assert.Equal(t, gust.Err[uint32]("error code: 13"), x.MapErr(stringify))
	}
}

func TestResult_MapOrElse_2(t *testing.T) {
	{
		var x = gust.Ok("foo")
		assert.Equal(t, "test:foo", x.MapOrElse(func(err error) string {
			return "bar"
		}, func(x string) string { return "test:" + x }))
	}
	{
		var x = gust.Err[string]("foo")
		assert.Equal(t, "bar", x.MapOrElse(func(err error) string {
			return "bar"
		}, func(x string) string { return "test:" + x }))
	}
}

func TestResult_MapOr_2(t *testing.T) {
	{
		var x = gust.Ok("foo")
		assert.Equal(t, "test:foo", x.MapOr("bar", func(x string) string { return "test:" + x }))
	}
	{
		var x = gust.Err[string]("foo")
		assert.Equal(t, "bar", x.MapOr("bar", func(x string) string { return "test:" + x }))
	}
}

func TestResult_Map_1(t *testing.T) {
	var line = "1\n2\n3\n4\n"
	for _, num := range strings.Split(line, "\n") {
		gust.Ret(strconv.Atoi(num)).Map(func(i int) int {
			return i * 2
		}).Inspect(func(i int) {
			t.Log(i)
		})
	}
}

func TestResult_Map_2(t *testing.T) {
	var isMyNum = func(s string, search int) gust.Result[any] {
		return gust.Ret(strconv.Atoi(s)).XMap(func(x int) any { return x == search })
	}
	assert.Equal(t, gust.Ok[any](true), isMyNum("1", 1))
	assert.Equal(t, gust.Ok[bool](true), ret.XAssert[bool](isMyNum("1", 1)))
	assert.Equal(t, "Err(strconv.Atoi: parsing \"lol\": invalid syntax)", isMyNum("lol", 1).String())
	assert.Equal(t, "Err(strconv.Atoi: parsing \"NaN\": invalid syntax)", isMyNum("NaN", 1).String())
}

func TestResult_Ok(t *testing.T) {
	{
		var x = gust.Ok[int](2)
		assert.Equal(t, gust.Some[int](2), x.Ok())
	}
	{
		var x = gust.Err[int]("some error message")
		assert.Equal(t, gust.None[int](), x.Ok())
	}
}

func TestResult_OrElse(t *testing.T) {
	var sq = func(err error) gust.Result[int] {
		if err == nil {
			return gust.Ok(0)
		}
		// Convert error to int for testing
		if err.Error() == "3" {
			return gust.Ok(9) // 3 * 3
		}
		return gust.Err[int](err)
	}
	var errFn = func(err error) gust.Result[int] {
		return gust.Err[int](err)
	}

	assert.Equal(t, gust.Ok(2), gust.Ok(2).OrElse(sq).OrElse(sq))
	assert.Equal(t, gust.Ok(2), gust.Ok(2).OrElse(errFn).OrElse(sq))
	assert.Equal(t, gust.Ok(9), gust.Err[int](fmt.Errorf("3")).OrElse(sq).OrElse(errFn))
	assert.Equal(t, gust.Err[int](fmt.Errorf("3")), gust.Err[int](fmt.Errorf("3")).OrElse(errFn).OrElse(errFn))
}

func TestResult_Or(t *testing.T) {
	{
		x := gust.Ok(2)
		y := gust.Err[int]("late error")
		assert.Equal(t, gust.Ok(2), x.Or(y))
	}
	{
		x := gust.Err[uint]("early error")
		y := gust.Ok[uint](2)
		assert.Equal(t, gust.Ok[uint](2), x.Or(y))
	}
	{
		x := gust.Err[uint]("not a 2")
		y := gust.Err[uint]("late error")
		assert.Equal(t, gust.Err[uint]("late error"), x.Or(y))
	}
	{
		x := gust.Ok[uint](2)
		y := gust.Ok[uint](100)
		assert.Equal(t, gust.Ok[uint](2), x.Or(y))
	}
}

func TestResult_Ret(t *testing.T) {
	var w = gust.Ret[int](strconv.Atoi("s"))
	assert.False(t, w.IsOk())
	assert.True(t, w.IsErr())

	var w2 = gust.Ret[any](strconv.Atoi("-1"))
	assert.True(t, w2.IsOk())
	assert.False(t, w2.IsErr())
	assert.Equal(t, -1, w2.Unwrap())
}

func TestResult_UnwrapErr_1(t *testing.T) {
	defer func() {
		assert.Equal(t, gust.ToErrBox("called `Result.UnwrapErr()` on an `ok` value: 10"), recover())
	}()
	err := gust.Ok(10).UnwrapErr()
	assert.NoError(t, err)
}

func TestResult_UnwrapErr_2(t *testing.T) {
	err := gust.Err[int]("emergency failure").UnwrapErr()
	if assert.Error(t, err) {
		assert.Equal(t, "emergency failure", err.Error())
	} else {
		t.FailNow()
	}
}

func TestFmtErr(t *testing.T) {
	result := gust.FmtErr[int]("error: %s", "test")
	assert.True(t, result.IsErr())
	assert.Contains(t, result.Err().Error(), "error: test")
}

func TestAssertRet(t *testing.T) {
	// Test with valid type
	result1 := gust.AssertRet[int](42)
	assert.True(t, result1.IsOk())
	assert.Equal(t, 42, result1.Unwrap())

	// Test with invalid type
	result2 := gust.AssertRet[int]("string")
	assert.True(t, result2.IsErr())
	assert.Contains(t, result2.Err().Error(), "type assert error")
}

func TestCatchResult(t *testing.T) {
	var result gust.Result[int]
	defer gust.CatchResult(&result)
	gust.Err[int]("test error").UnwrapOrThrow()
	assert.True(t, result.IsErr())
	assert.Equal(t, "test error", result.Err().Error())
}

// TestCatchResult_IsSomeBranch tests CatchResult with IsSome branch (covers result.go:609-611)
func TestCatchResult_IsSomeBranch(t *testing.T) {
	// Test result.t.IsSome() branch
	var result gust.Result[int]
	result = gust.Ok(42)
	func() {
		defer gust.CatchResult(&result)
		gust.Err[int]("panic error").UnwrapOrThrow()
	}()
	assert.True(t, result.IsErr())
	assert.Equal(t, "panic error", result.Err().Error())
}

func TestResult_UnwrapOrDefault(t *testing.T) {
	assert.Equal(t, "car", gust.Ok("car").UnwrapOrDefault())
	assert.Equal(t, "", gust.Err[string](nil).UnwrapOrDefault())
	assert.Equal(t, time.Time{}, gust.Err[time.Time](nil).UnwrapOrDefault())
	assert.Equal(t, &time.Time{}, gust.Err[*time.Time](nil).UnwrapOrDefault())
	assert.Equal(t, valconv.Ref(&time.Time{}), gust.Err[**time.Time](nil).UnwrapOrDefault())
}

func TestResult_UnwrapOrElse(t *testing.T) {
	var count = func(x error) int {
		return len(x.Error())
	}
	assert.Equal(t, 2, gust.Ok(2).UnwrapOrElse(count))
	assert.Equal(t, 3, gust.Err[int]("foo").UnwrapOrElse(count))
}

func TestResult_Unwrap(t *testing.T) {
	defer func() {
		assert.Equal(t, "strconv.Atoi: parsing \"4x\": invalid syntax", recover().(error).Error())
	}()
	gust.Ret(strconv.Atoi("4x")).Unwrap()
}

func TestResult_UnwrapUnchecked(t *testing.T) {
	{
		var r gust.Result[string]
		assert.Equal(t, "", r.UnwrapUnchecked())
	}
	{
		var r = gust.Err[string]("foo")
		assert.Equal(t, "", r.UnwrapUnchecked())
	}
	{
		var r = gust.Ok[string]("foo")
		assert.Equal(t, "foo", r.UnwrapUnchecked())
	}
}

func TestResult_XOk(t *testing.T) {
	{
		var x = gust.Ok[int](2)
		assert.Equal(t, gust.Some[any](2), x.XOk())
	}
	{
		var x = gust.Err[int]("some error message")
		assert.Equal(t, gust.None[any](), x.XOk())
	}
}

func TestResult_ErrVal(t *testing.T) {
	{
		var x = gust.Ok[int](2)
		assert.Nil(t, x.ErrVal())
	}
	{
		var x = gust.Err[int]("some error message")
		assert.Equal(t, "some error message", x.ErrVal())
	}
	{
		var x = gust.Err[int](gust.ToErrBox("boxed error"))
		assert.Equal(t, "boxed error", x.ErrVal())
	}
	{
		var x = gust.Err[int](errors.New("std error"))
		assert.NotNil(t, x.ErrVal())
	}
}

func TestResult_ToX(t *testing.T) {
	{
		var x = gust.Ok[int](42)
		xResult := x.ToX()
		assert.True(t, xResult.IsOk())
		assert.Equal(t, 42, xResult.Unwrap())
	}
	{
		var x = gust.Err[int]("error")
		xResult := x.ToX()
		assert.True(t, xResult.IsErr())
		assert.Equal(t, "error", xResult.Err().Error())
	}
}

func TestResult_XMap(t *testing.T) {
	{
		var x = gust.Ok[int](2)
		result := x.XMap(func(i int) any { return i * 2 })
		assert.True(t, result.IsOk())
		assert.Equal(t, 4, result.Unwrap())
	}
	{
		var x = gust.Err[int]("error")
		result := x.XMap(func(i int) any { return i * 2 })
		assert.True(t, result.IsErr())
	}
}

func TestResult_XMapOr(t *testing.T) {
	{
		var x = gust.Ok[int](2)
		result := x.XMapOr("default", func(i int) any { return i * 2 })
		assert.Equal(t, 4, result)
	}
	{
		var x = gust.Err[int]("error")
		result := x.XMapOr("default", func(i int) any { return i * 2 })
		assert.Equal(t, "default", result)
	}
}

func TestResult_XMapOrElse(t *testing.T) {
	{
		var x = gust.Ok[int](2)
		result := x.XMapOrElse(func(error) any { return "default" }, func(i int) any { return i * 2 })
		assert.Equal(t, 4, result)
	}
	{
		var x = gust.Err[int]("error")
		result := x.XMapOrElse(func(error) any { return "default" }, func(i int) any { return i * 2 })
		assert.Equal(t, "default", result)
	}
}

func TestResult_ContainsErr(t *testing.T) {
	{
		var x = gust.Ok[int](2)
		assert.False(t, x.ContainsErr(errors.New("test")))
	}
	{
		// Use same error instance for comparison
		testErr := errors.New("test error")
		var x = gust.Err[int](testErr)
		assert.True(t, x.ContainsErr(testErr))
		assert.False(t, x.ContainsErr(errors.New("different error")))
	}
	{
		err1 := errors.New("wrapped")
		err2 := fmt.Errorf("outer: %w", err1)
		var x = gust.Err[int](err2)
		assert.True(t, x.ContainsErr(err1))
	}
}

func TestResult_Iterator(t *testing.T) {
	// Test Next
	{
		var x = gust.Ok("foo")
		opt := x.Next()
		assert.Equal(t, gust.Some("foo"), opt)
		// Next() consumes the value, so Ok() should return None after Next()
		assert.True(t, x.Ok().IsNone())
	}
	{
		var x = gust.Err[string]("error")
		opt := x.Next()
		assert.True(t, opt.IsNone())
	}
	{
		var nilResult *gust.Result[string]
		opt := nilResult.Next()
		assert.True(t, opt.IsNone())
	}

	// Test NextBack
	{
		var x = gust.Ok("bar")
		opt := x.NextBack()
		assert.Equal(t, gust.Some("bar"), opt)
		// NextBack() also consumes the value
		assert.True(t, x.Ok().IsNone())
	}

	// Test Remaining
	{
		var x = gust.Ok("baz")
		assert.Equal(t, uint(1), x.Remaining())
		x.Next() // Consume the value
		assert.Equal(t, uint(0), x.Remaining())
	}
	{
		var x = gust.Err[string]("error")
		assert.Equal(t, uint(0), x.Remaining())
	}
	{
		var nilResult *gust.Result[string]
		assert.Equal(t, uint(0), nilResult.Remaining())
	}
}

func TestResult_Ref(t *testing.T) {
	result := gust.Ok(42)
	ref := result.Ref()
	assert.Equal(t, gust.Ok(42), *ref)
	ref.Unwrap() // Should not panic
}

func TestResult_Errable(t *testing.T) {
	{
		var x = gust.Ok[string]("foo")
		errable := x.Errable()
		assert.False(t, errable.IsErr())
	}
	{
		var x = gust.Err[string](errors.New("error"))
		errable := x.Errable()
		assert.True(t, errable.IsErr())
	}
}

func TestResult_Split(t *testing.T) {
	{
		var x = gust.Ok("foo")
		val, err := x.Split()
		assert.NoError(t, err)
		assert.Equal(t, "foo", val)
	}
	{
		var x = gust.Err[string](errors.New("error"))
		val, err := x.Split()
		assert.Error(t, err)
		assert.Equal(t, "", val)
	}
}

func TestResult_UnmarshalJSON_NilReceiver(t *testing.T) {
	// Test UnmarshalJSON with nil receiver
	var nilResult *gust.Result[int]
	err := nilResult.UnmarshalJSON([]byte("42"))
	assert.Error(t, err)
	assert.IsType(t, &json.InvalidUnmarshalError{}, err)
}

func TestResult_UnmarshalJSON_ErrorPath(t *testing.T) {
	// Test UnmarshalJSON with invalid JSON (error path)
	var result gust.Result[int]
	err := result.UnmarshalJSON([]byte("invalid json"))
	assert.Error(t, err)
	assert.True(t, result.IsErr()) // Should be Err on error
}

func TestResult_UnmarshalJSON_ValidAfterError(t *testing.T) {
	// Test UnmarshalJSON with error first, then valid JSON
	var result gust.Result[int]
	// First attempt with invalid JSON
	_ = result.UnmarshalJSON([]byte("invalid"))
	assert.True(t, result.IsErr())

	// Then with valid JSON
	err := result.UnmarshalJSON([]byte("42"))
	assert.NoError(t, err)
	assert.True(t, result.IsOk())
	assert.Equal(t, 42, result.Unwrap())
}

func TestResult_UnmarshalJSON_Struct(t *testing.T) {
	type S struct {
		X int
		Y string
	}
	var result gust.Result[S]
	err := result.UnmarshalJSON([]byte(`{"X":10,"Y":"test"}`))
	assert.NoError(t, err)
	assert.True(t, result.IsOk())
	assert.Equal(t, S{X: 10, Y: "test"}, result.Unwrap())
}

func TestResult_UnmarshalJSON_Array(t *testing.T) {
	var result gust.Result[[]int]
	err := result.UnmarshalJSON([]byte("[1,2,3]"))
	assert.NoError(t, err)
	assert.True(t, result.IsOk())
	assert.Equal(t, []int{1, 2, 3}, result.Unwrap())
}

func TestResult_UnmarshalJSON_Map(t *testing.T) {
	var result gust.Result[map[string]int]
	err := result.UnmarshalJSON([]byte(`{"a":1,"b":2}`))
	assert.NoError(t, err)
	assert.True(t, result.IsOk())
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, result.Unwrap())
}

func TestResult_Catch_NilReceiver(t *testing.T) {
	// Test Catch with nil receiver
	defer func() {
		assert.NotNil(t, recover())
	}()
	var nilResult *gust.Result[int]
	defer nilResult.Catch()
	gust.Err[int]("test error").UnwrapOrThrow()
}

func TestCatchResult_NilReceiver(t *testing.T) {
	// Test CatchResult with nil receiver
	defer func() {
		assert.NotNil(t, recover())
	}()
	defer gust.CatchResult[int](nil)
	gust.Err[int]("test error").UnwrapOrThrow()
}

func TestResult_Catch_NonPanicValue(t *testing.T) {
	// Test Catch with non-panicValue panic
	defer func() {
		assert.Equal(t, "regular panic", recover())
	}()
	var result gust.Result[int]
	defer result.Catch()
	panic("regular panic")
}

func TestCatchResult_NonPanicValue(t *testing.T) {
	// Test CatchResult with non-panicValue panic
	defer func() {
		assert.Equal(t, "regular panic", recover())
	}()
	var result gust.Result[int]
	defer gust.CatchResult(&result)
	panic("regular panic")
}

func TestResult_Catch_OkValue(t *testing.T) {
	// Test Catch when result already has Ok value
	var result gust.Result[int] = gust.Ok(42)
	defer result.Catch()
	gust.Err[string]("test error").UnwrapOrThrow()
	// Result should be updated to Err
	assert.True(t, result.IsErr())

	// Test CatchResult when result already has Ok value
	var result2 gust.Result[int] = gust.Ok(42)
	defer gust.CatchResult(&result2)
	gust.Err[string]("test error").UnwrapOrThrow()
	// Result should be updated to Err
	assert.True(t, result2.IsErr())
}

func TestResult_XAndThen(t *testing.T) {
	// Test XAndThen with Ok path
	result := gust.Ok[int](42)
	result2 := result.XAndThen(func(i int) gust.Result[any] {
		return gust.Ok[any](i * 2)
	})
	assert.True(t, result2.IsOk())
	assert.Equal(t, 84, result2.Unwrap())

	// Test XAndThen with error path
	result3 := gust.Ok[int](42)
	result4 := result3.XAndThen(func(i int) gust.Result[any] {
		return gust.Err[any]("error")
	})
	assert.True(t, result4.IsErr())
}

func TestResult_AndThen2(t *testing.T) {
	// Test with Ok result and successful operation
	{
		x := gust.Ok(2)
		result := x.AndThen2(func(i int) (int, error) {
			return i * 2, nil
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, 4, result.Unwrap())
	}
	// Test with Ok result and error operation
	{
		x := gust.Ok(2)
		result := x.AndThen2(func(i int) (int, error) {
			return 0, errors.New("operation error")
		})
		assert.True(t, result.IsErr())
		assert.Equal(t, "operation error", result.Err().Error())
	}
	// Test with Err result (should return original error)
	{
		x := gust.Err[int]("early error")
		result := x.AndThen2(func(i int) (int, error) {
			return i * 2, nil
		})
		assert.True(t, result.IsErr())
		assert.Equal(t, "early error", result.Err().Error())
	}
}

func TestResult_XAndThen2(t *testing.T) {
	// Test with Ok result and successful operation
	{
		x := gust.Ok(2)
		result := x.XAndThen2(func(i int) (any, error) {
			return i * 2, nil
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, 4, result.Unwrap())
	}
	// Test with Ok result and error operation
	{
		x := gust.Ok(2)
		result := x.XAndThen2(func(i int) (any, error) {
			return nil, errors.New("operation error")
		})
		assert.True(t, result.IsErr())
		assert.Equal(t, "operation error", result.Err().Error())
	}
	// Test with Err result (should return original error)
	{
		x := gust.Err[int]("early error")
		result := x.XAndThen2(func(i int) (any, error) {
			return i * 2, nil
		})
		assert.True(t, result.IsErr())
		assert.Equal(t, "early error", result.Err().Error())
	}
}

func TestResult_Or2(t *testing.T) {
	// Test with Ok result (should return original Ok)
	{
		x := gust.Ok(2)
		result := x.Or2(3, nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, 2, result.Unwrap())
	}
	// Test with Ok result and Err value (should return original Ok)
	{
		x := gust.Ok(2)
		result := x.Or2(3, errors.New("late error"))
		assert.True(t, result.IsOk())
		assert.Equal(t, 2, result.Unwrap())
	}
	// Test with Err result and Ok value
	{
		x := gust.Err[int]("early error")
		result := x.Or2(3, nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, 3, result.Unwrap())
	}
	// Test with Err result and Err value
	{
		x := gust.Err[int]("early error")
		result := x.Or2(3, errors.New("late error"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "late error", result.Err().Error())
	}
}

func TestResult_OrElse2(t *testing.T) {
	// Test with Ok result (should return original Ok)
	{
		x := gust.Ok(2)
		result := x.OrElse2(func(err error) (int, error) {
			return 3, nil
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, 2, result.Unwrap())
	}
	// Test with Err result and successful operation
	{
		x := gust.Err[int]("early error")
		result := x.OrElse2(func(err error) (int, error) {
			return 3, nil
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, 3, result.Unwrap())
	}
	// Test with Err result and error operation
	{
		x := gust.Err[int]("early error")
		result := x.OrElse2(func(err error) (int, error) {
			return 0, errors.New("late error")
		})
		assert.True(t, result.IsErr())
		assert.Equal(t, "late error", result.Err().Error())
	}
}

func TestResult_Flatten(t *testing.T) {
	// Test with Ok result and nil error (should return original Ok)
	{
		r := gust.Ok(42)
		result := r.Flatten(nil)
		assert.True(t, result.IsOk())
		assert.Equal(t, 42, result.Unwrap())
	}
	// Test with Ok result and error (should return error)
	{
		r := gust.Ok(42)
		result := r.Flatten(errors.New("test error"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "test error", result.Err().Error())
	}
	// Test with Err result and nil error (should return original Err)
	{
		r := gust.Err[int]("original error")
		result := r.Flatten(nil)
		assert.True(t, result.IsErr())
		assert.Equal(t, "original error", result.Err().Error())
	}
	// Test with Err result and error (should return the provided error)
	{
		r := gust.Err[int]("original error")
		result := r.Flatten(errors.New("new error"))
		assert.True(t, result.IsErr())
		assert.Equal(t, "new error", result.Err().Error())
	}
}
