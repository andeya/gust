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

func TestResultUnwrapOrThrow(t *testing.T) {
	// Test UnwrapOrThrow with Ok value
	var r1 = gust.Ok(1)
	var v1 = r1.UnwrapOrThrow()
	assert.Equal(t, 1, v1)

	// Test UnwrapOrThrow with Err value (should panic and be caught)
	var r gust.Result[string]
	func() {
		defer r.Catch()
		var r2 = gust.Err[int]("err")
		_ = r2.UnwrapOrThrow() // This will panic
	}()
	assert.True(t, r.IsErr())
	// Error() should only return error message, not stack trace
	errMsg := r.Err().Error()
	assert.Contains(t, errMsg, "err")
	// Should NOT contain stack trace in Error()
	assert.NotContains(t, errMsg, "\n")
	// But %+v should contain stack trace
	fullMsg := fmt.Sprintf("%+v", r.Err())
	assert.Contains(t, fullMsg, "err")
	assert.Contains(t, fullMsg, "\n")
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
		assert.Equal(t, gust.BoxErr("Testing expect_err: 10"), recover())
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
		assert.Equal(t, gust.BoxErr("called `Result.UnwrapErr()` on an `ok` value: 10"), recover())
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

func TestResult_Catch(t *testing.T) {
	// Test basic Catch functionality (with stack trace by default)
	var result gust.Result[int]
	func() {
		defer result.Catch()
		gust.Err[int]("test error").UnwrapOrThrow()
	}()
	assert.True(t, result.IsErr())
	// Error() should only return error message, not stack trace
	errMsg := result.Err().Error()
	assert.Contains(t, errMsg, "test error")
	// Should NOT contain stack trace in Error()
	assert.NotContains(t, errMsg, "\n")
	// But %+v should contain stack trace
	fullMsg := fmt.Sprintf("%+v", result.Err())
	assert.Contains(t, fullMsg, "test error")
	assert.Contains(t, fullMsg, "\n")
}

func TestResult_Catch_IsSomeBranch(t *testing.T) {
	// Test Catch when result already has Ok value (IsSome branch)
	var result gust.Result[int] = gust.Ok(42)
	func() {
		defer result.Catch()
		gust.Err[int]("panic error").UnwrapOrThrow()
	}()
	assert.True(t, result.IsErr())
	// Error() should only return error message, not stack trace
	errMsg := result.Err().Error()
	assert.Contains(t, errMsg, "panic error")
	// Should NOT contain stack trace in Error()
	assert.NotContains(t, errMsg, "\n")
	// But %+v should contain stack trace
	fullMsg := fmt.Sprintf("%+v", result.Err())
	assert.Contains(t, fullMsg, "panic error")
	assert.Contains(t, fullMsg, "\n")
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
		p := recover()
		if eb, ok := p.(*gust.ErrBox); ok {
			// Unwrap() panics *ErrBox, convert to error for comparison
			assert.Equal(t, "strconv.Atoi: parsing \"4x\": invalid syntax", eb.String())
		} else {
			t.Fatalf("expected *gust.ErrBox, got %T: %v", p, p)
		}
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
		var x = gust.Err[int](gust.BoxErr("boxed error"))
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

func TestResult_ErrToVoidResult(t *testing.T) {
	{
		var x = gust.Ok[string]("foo")
		result := gust.RetVoid(x.Err())
		assert.False(t, result.IsErr())
	}
	{
		var x = gust.Err[string](errors.New("error"))
		result := gust.RetVoid(x.Err())
		assert.True(t, result.IsErr())
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

func TestResult_Catch_NonPanicValue(t *testing.T) {
	// Test Catch with non-ErrBox panic (should be caught and converted to ErrBox)
	var result gust.Result[int]
	func() {
		defer result.Catch()
		panic("regular panic")
	}()
	// Should be caught and converted to error
	assert.True(t, result.IsErr())
	// Error() should only return error message, not stack trace
	errMsg := result.Err().Error()
	assert.Contains(t, errMsg, "regular panic")
	// Should NOT contain stack trace in Error()
	assert.NotContains(t, errMsg, "\n")
	// But %+v should contain stack trace
	fullMsg := fmt.Sprintf("%+v", result.Err())
	assert.Contains(t, fullMsg, "regular panic")
	assert.Contains(t, fullMsg, "\n")
}

func TestResult_Catch_OkValue(t *testing.T) {
	// Test Catch when result already has Ok value (should be updated to Err)
	var result gust.Result[int] = gust.Ok(42)
	func() {
		defer result.Catch()
		gust.Err[string]("test error").UnwrapOrThrow()
	}()
	// Result should be updated to Err
	assert.True(t, result.IsErr())
	// Error() should only return error message, not stack trace
	errMsg := result.Err().Error()
	assert.Contains(t, errMsg, "test error")
	// Should NOT contain stack trace in Error()
	assert.NotContains(t, errMsg, "\n")
	// But %+v should contain stack trace
	fullMsg := fmt.Sprintf("%+v", result.Err())
	assert.Contains(t, fullMsg, "test error")
	assert.Contains(t, fullMsg, "\n")
}

// TestResult_Catch_WithStackTrace tests that Catch captures stack trace information
func TestResult_Catch_WithStackTrace(t *testing.T) {
	// Test with ErrBox panic
	{
		var result gust.Result[int]
		func() {
			defer result.Catch()
			gust.Err[int]("test error").UnwrapOrThrow()
		}()
		assert.True(t, result.IsErr())
		err := result.Err()
		assert.NotNil(t, err)
		// Error() should only return error message, not stack trace
		errMsg := err.Error()
		assert.Contains(t, errMsg, "test error")
		assert.NotContains(t, errMsg, "\n")
		// But %+v should contain stack trace
		fullMsg := fmt.Sprintf("%+v", err)
		assert.Contains(t, fullMsg, "test error")
		assert.Contains(t, fullMsg, "\n")
	}

	// Test with regular error panic
	{
		var result gust.Result[string]
		func() {
			defer result.Catch()
			panic(errors.New("regular error"))
		}()
		assert.True(t, result.IsErr())
		err := result.Err()
		assert.NotNil(t, err)
		errMsg := err.Error()
		assert.Contains(t, errMsg, "regular error")
		assert.NotContains(t, errMsg, "\n")
		// But %+v should contain stack trace
		fullMsg := fmt.Sprintf("%+v", err)
		assert.Contains(t, fullMsg, "regular error")
		assert.Contains(t, fullMsg, "\n")
	}

	// Test with string panic
	{
		var result gust.Result[int]
		func() {
			defer result.Catch()
			panic("string panic")
		}()
		assert.True(t, result.IsErr())
		err := result.Err()
		assert.NotNil(t, err)
		errMsg := err.Error()
		assert.Contains(t, errMsg, "string panic")
		assert.NotContains(t, errMsg, "\n")
		// But %+v should contain stack trace
		fullMsg := fmt.Sprintf("%+v", err)
		assert.Contains(t, fullMsg, "string panic")
		assert.Contains(t, fullMsg, "\n")
	}

	// Test with int panic
	{
		var result gust.Result[string]
		func() {
			defer result.Catch()
			panic(42)
		}()
		assert.True(t, result.IsErr())
		err := result.Err()
		assert.NotNil(t, err)
		errMsg := err.Error()
		assert.Contains(t, errMsg, "42")
		assert.NotContains(t, errMsg, "\n")
		// But %+v should contain stack trace
		fullMsg := fmt.Sprintf("%+v", err)
		assert.Contains(t, fullMsg, "42")
		assert.Contains(t, fullMsg, "\n")
	}
}

// TestResult_Catch_StackTraceAccess tests accessing stack trace from caught panic
func TestResult_Catch_StackTraceAccess(t *testing.T) {
	var result gust.Result[int]
	func() {
		defer result.Catch()
		gust.Err[int]("test error").UnwrapOrThrow()
	}()
	assert.True(t, result.IsErr())

	// Get the error
	err := result.Err()
	assert.NotNil(t, err)

	// Error() should only return error message, not stack trace
	errMsg := err.Error()
	assert.Contains(t, errMsg, "test error")
	assert.NotContains(t, errMsg, "\n")

	// Try to extract StackTraceCarrier from error chain
	var carrier gust.StackTraceCarrier
	if errors.As(err, &carrier) {
		stack := carrier.StackTrace()
		assert.True(t, len(stack) > 0, "Stack trace should not be empty")
		// Verify stack trace contains frames
		assert.Greater(t, len(stack), 0)
	}

	// Verify %+v contains stack trace
	fullMsg := fmt.Sprintf("%+v", err)
	assert.Contains(t, fullMsg, "test error")
	assert.Contains(t, fullMsg, "\n")
}

// TestResult_Catch_ErrBoxPointer tests Catch with *ErrBox panic
func TestResult_Catch_ErrBoxPointer(t *testing.T) {
	var result gust.Result[int]
	eb := gust.BoxErr(errors.New("errbox pointer error"))
	func() {
		defer result.Catch()
		panic(eb)
	}()
	assert.True(t, result.IsErr())
	err := result.Err()
	assert.NotNil(t, err)
	// Error() should only return error message, not stack trace
	errMsg := err.Error()
	assert.Contains(t, errMsg, "errbox pointer error")
	assert.NotContains(t, errMsg, "\n")
	// But %+v should contain stack trace
	fullMsg := fmt.Sprintf("%+v", err)
	assert.Contains(t, fullMsg, "errbox pointer error")
	assert.Contains(t, fullMsg, "\n")
}

// TestResult_Catch_ErrBoxValue tests Catch with ErrBox value panic
func TestResult_Catch_ErrBoxValue(t *testing.T) {
	var result gust.Result[int]
	eb := gust.BoxErr(errors.New("errbox value error"))
	func() {
		defer result.Catch()
		panic(*eb)
	}()
	assert.True(t, result.IsErr())
	err := result.Err()
	assert.NotNil(t, err)
	// Error() should only return error message, not stack trace
	errMsg := err.Error()
	assert.Contains(t, errMsg, "errbox value error")
	assert.NotContains(t, errMsg, "\n")
	// But %+v should contain stack trace
	fullMsg := fmt.Sprintf("%+v", err)
	assert.Contains(t, fullMsg, "errbox value error")
	assert.Contains(t, fullMsg, "\n")
}

// TestResult_Catch_NilErrBoxPointer tests Catch with nil *ErrBox panic
func TestResult_Catch_NilErrBoxPointer(t *testing.T) {
	var result gust.Result[int]
	var eb *gust.ErrBox
	func() {
		defer result.Catch()
		panic(eb)
	}()
	assert.True(t, result.IsErr())
	err := result.Err()
	assert.NotNil(t, err)
	errMsg := err.Error()
	// With nil ErrBox, error message might be "<nil>" or contain stack trace
	// Both cases are acceptable
	if errMsg == "<nil>" {
		// Without stack trace (if Catch(false) was used or stack is empty)
		assert.Equal(t, "<nil>", errMsg)
	} else {
		// With stack trace (default behavior)
		assert.Contains(t, errMsg, "\n")
	}
}

// TestResult_Catch_WithStackTraceFalse tests Catch without stack trace capture
func TestResult_Catch_WithStackTraceFalse(t *testing.T) {
	// Test with withStackTrace=false
	var result gust.Result[int]
	func() {
		defer result.Catch(false)
		gust.Err[int]("test error").UnwrapOrThrow()
	}()
	assert.True(t, result.IsErr())
	err := result.Err()
	assert.NotNil(t, err)
	errMsg := err.Error()
	// Should contain error message but not stack trace (no newlines from stack)
	assert.Contains(t, errMsg, "test error")
	// Error message should not contain stack trace formatting
	// (Note: panicError.Error() will still format stack, but it will be empty)
}

// TestResult_Catch_WithStackTraceTrue tests Catch with stack trace capture (default)
func TestResult_Catch_WithStackTraceTrue(t *testing.T) {
	// Test with withStackTrace=true (explicit)
	var result gust.Result[int]
	func() {
		defer func() {
			t.Logf("error with panic stack trace: %+v", result.Err())
		}()
		defer result.Catch(true)
		gust.Err[int]("test error").UnwrapOrThrow()
	}()
	assert.True(t, result.IsErr())
	err := result.Err()
	assert.NotNil(t, err)
	// Error() should only return error message, not stack trace
	errMsg := err.Error()
	assert.Contains(t, errMsg, "test error")
	assert.NotContains(t, errMsg, "\n")
	// But %+v should contain stack trace
	fullMsg := fmt.Sprintf("%+v", err)
	assert.Contains(t, fullMsg, "test error")
	assert.Contains(t, fullMsg, "\n")
}

// TestResult_Catch_WithStackTraceDefault tests Catch with default behavior (stack trace enabled)
func TestResult_Catch_WithStackTraceDefault(t *testing.T) {
	// Test with default (no parameter, should default to true)
	var result gust.Result[int]
	func() {
		defer result.Catch()
		gust.Err[int]("test error").UnwrapOrThrow()
	}()
	assert.True(t, result.IsErr())
	err := result.Err()
	assert.NotNil(t, err)
	// Error() should only return error message, not stack trace
	errMsg := err.Error()
	assert.Contains(t, errMsg, "test error")
	assert.NotContains(t, errMsg, "\n")
	// But %+v should contain stack trace
	fullMsg := fmt.Sprintf("%+v", err)
	assert.Contains(t, fullMsg, "test error")
	assert.Contains(t, fullMsg, "\n")
}

// TestResult_Catch_WithStackTracePerformance tests that disabling stack trace improves performance
func TestResult_Catch_WithStackTracePerformance(t *testing.T) {
	// Test that both modes work correctly
	var result1 gust.Result[int]
	var result2 gust.Result[int]

	// With stack trace
	func() {
		defer result1.Catch(true)
		panic("error with stack")
	}()
	assert.True(t, result1.IsErr())

	// Without stack trace
	func() {
		defer result2.Catch(false)
		panic("error without stack")
	}()
	assert.True(t, result2.IsErr())

	// Both should be errors
	assert.True(t, result1.IsErr())
	assert.True(t, result2.IsErr())
	// Both should contain error message
	assert.Contains(t, result1.Err().Error(), "error with stack")
	assert.Contains(t, result2.Err().Error(), "error without stack")
	// Only result1 should have stack trace in %+v
	fullMsg1 := fmt.Sprintf("%+v", result1.Err())
	fullMsg2 := fmt.Sprintf("%+v", result2.Err())
	assert.Contains(t, fullMsg1, "\n")
	assert.NotContains(t, fullMsg2, "\n")
}

// TestResult_Catch_FormatOptions tests different format options for caught errors
func TestResult_Catch_FormatOptions(t *testing.T) {
	var result gust.Result[int]
	func() {
		defer result.Catch()
		gust.Err[int]("format test error").UnwrapOrThrow()
	}()
	assert.True(t, result.IsErr())
	err := result.Err()

	// Test %v (error message only)
	errMsgV := fmt.Sprintf("%v", err)
	t.Logf("%%v format: %v", errMsgV)
	assert.Equal(t, "format test error", errMsgV)
	assert.NotContains(t, errMsgV, "\n")

	// Test %+v (error message with detailed stack trace - shows function name, full path, and line number)
	// Format: each frame on a separate line with "\n函数名\n\t完整路径:行号"
	errMsgPlusV := fmt.Sprintf("%+v", err)
	t.Logf("=== %%+v format (detailed - each frame on separate line WITH line numbers) ===")
	t.Logf("%+v", err)
	t.Logf("=== End of %%+v format ===")
	assert.Contains(t, errMsgPlusV, "format test error")
	assert.Contains(t, errMsgPlusV, "\n")
	// %+v should contain line numbers in format "file:line" or "path:line"
	assert.Contains(t, errMsgPlusV, ":") // Should contain line numbers

	// Test %s (error message only)
	errMsgS := fmt.Sprintf("%s", err)
	t.Logf("=== %%s format (error message only) ===")
	t.Logf("%s", errMsgS)
	t.Logf("=== End of %%s format ===")
	assert.Equal(t, "format test error", errMsgS)
	assert.NotContains(t, errMsgS, "\n")

	// Test %+s (error message with stack trace in %s format - shows function name + full path, NO line numbers)
	// Format: "[函数名\n\t完整路径 函数名\n\t完整路径 ...]" (wrapped in brackets, space-separated, NO line numbers)
	errMsgPlusS := fmt.Sprintf("%+s", err)
	t.Logf("=== %%+s format (function name + full path, wrapped in brackets, NO line numbers) ===")
	t.Logf("%+s", err)
	t.Logf("=== End of %%+s format ===")
	assert.Contains(t, errMsgPlusS, "format test error")
	assert.Contains(t, errMsgPlusS, "\n")
	// %+s should be wrapped in brackets and space-separated
	assert.True(t, strings.HasPrefix(errMsgPlusS[strings.Index(errMsgPlusS, "\n")+1:], "["), "%%+s should start with '[' after error message")
	// Note: %+s may still contain ":" in file paths, but should NOT have ":行号" format
	// The key difference is that %+v shows each frame on a separate line with line numbers,
	// while %+s shows frames in brackets separated by spaces without explicit line numbers
}

// TestResult_Catch_NoPanic tests Catch when no panic occurs
func TestResult_Catch_NoPanic(t *testing.T) {
	var result gust.Result[int] = gust.Ok(42)
	func() {
		defer result.Catch()
		// No panic, just return normally
		result = gust.Ok(100)
	}()
	// Result should remain Ok
	assert.True(t, result.IsOk())
	assert.Equal(t, 100, result.Unwrap())
	assert.False(t, result.IsErr())
}

// TestResult_Catch_MultiplePanics tests that Catch only catches the first panic
func TestResult_Catch_MultiplePanics(t *testing.T) {
	var result gust.Result[int]
	func() {
		defer result.Catch()
		panic("first panic")
		// Note: second panic would never be reached due to first panic
	}()
	assert.True(t, result.IsErr())
	errMsg := result.Err().Error()
	assert.Contains(t, errMsg, "first panic")
}

// TestResult_Catch_WithStackTraceFalse_Verification tests Catch(false) doesn't capture stack
func TestResult_Catch_WithStackTraceFalse_Verification(t *testing.T) {
	var result gust.Result[int]
	func() {
		defer result.Catch(false)
		gust.Err[int]("no stack error").UnwrapOrThrow()
	}()
	assert.True(t, result.IsErr())
	err := result.Err()
	assert.NotNil(t, err)

	// Error() should only return error message
	errMsg := err.Error()
	assert.Contains(t, errMsg, "no stack error")
	assert.NotContains(t, errMsg, "\n")

	// %+v should also not contain stack trace (because Catch(false) was used)
	fullMsg := fmt.Sprintf("%+v", err)
	assert.Contains(t, fullMsg, "no stack error")
	// Without stack trace capture, %+v should be same as %v
	assert.Equal(t, errMsg, fullMsg)
}

// TestResult_Catch_ErrorChain tests error unwrapping with caught panics
func TestResult_Catch_ErrorChain(t *testing.T) {
	baseErr := errors.New("base error")
	wrappedErr := fmt.Errorf("wrapped: %w", baseErr)

	var result gust.Result[int]
	func() {
		defer result.Catch()
		panic(wrappedErr)
	}()
	assert.True(t, result.IsErr())

	err := result.Err()
	// Verify error chain is preserved
	assert.True(t, errors.Is(err, baseErr))
	assert.True(t, errors.Is(err, wrappedErr))

	// Error() should only return the wrapped error message
	errMsg := err.Error()
	assert.Contains(t, errMsg, "wrapped")
	assert.NotContains(t, errMsg, "\n")
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

// TestOkVoid tests OkVoid function (covers result.go:105-107)
func TestOkVoid(t *testing.T) {
	result := gust.OkVoid()
	assert.True(t, result.IsOk())
	assert.Nil(t, result.Err())
}

// TestErrVoid tests ErrVoid function (covers result.go:102-107)
func TestErrVoid(t *testing.T) {
	// Test with error value
	{
		err := errors.New("test error")
		result := gust.ErrVoid(err)
		assert.True(t, result.IsErr())
		assert.False(t, result.IsOk())
		assert.NotNil(t, result.Err())
		assert.Equal(t, "test error", result.Err().Error())
	}
	// Test with string error
	{
		result := gust.ErrVoid("string error")
		assert.True(t, result.IsErr())
		assert.False(t, result.IsOk())
		assert.NotNil(t, result.Err())
		assert.Equal(t, "string error", result.Err().Error())
	}
	// Test with nil error (should still be error state, following declarative programming principles)
	{
		result := gust.ErrVoid(nil)
		assert.True(t, result.IsErr())
		assert.False(t, result.IsOk())
		// Even if err is nil, ErrVoid(nil) still represents an error state
		// This follows declarative programming principles consistent with Err[T](nil)
	}
	// Test that ErrVoid is equivalent to Err[Void]
	{
		err := errors.New("comparison error")
		result1 := gust.ErrVoid(err)
		result2 := gust.Err[gust.Void](err)
		assert.Equal(t, result1.IsErr(), result2.IsErr())
		assert.Equal(t, result1.IsOk(), result2.IsOk())
		if result1.IsErr() && result2.IsErr() {
			assert.Equal(t, result1.Err().Error(), result2.Err().Error())
		}
	}
}

// TestUnwrapErrOr tests UnwrapErrOr function (covers result.go:253-257)
func TestUnwrapErrOr(t *testing.T) {
	// Test with Err result
	{
		r := gust.Err[gust.Void]("test error")
		def := errors.New("default error")
		err := gust.UnwrapErrOr(r, def)
		assert.NotNil(t, err)
		assert.Equal(t, "test error", err.Error())
	}
	// Test with Ok result (should return default)
	{
		r := gust.Ok[gust.Void](nil)
		def := errors.New("default error")
		err := gust.UnwrapErrOr(r, def)
		assert.Equal(t, def, err)
	}
}

// TestResult_wrapError tests wrapError method indirectly through Expect
func TestResult_wrapError(t *testing.T) {
	// Test wrapError with nil value (covers result.go:373-376)
	{
		r := gust.Err[int](nil)
		defer func() {
			p := recover()
			if p != nil {
				if err, ok := p.(error); ok {
					assert.Contains(t, err.Error(), "gust.Err(nil)")
				} else {
					t.Errorf("Expected error panic, got %T: %v", p, p)
				}
			}
		}()
		r.Expect("test message")
		t.Fatal("Expected panic")
	}
	// Test wrapError with non-error value (covers result.go:380)
	{
		r := gust.Err[int](42)
		defer func() {
			p := recover()
			if p != nil {
				if err, ok := p.(error); ok {
					assert.Contains(t, err.Error(), "42")
				} else {
					t.Errorf("Expected error panic, got %T: %v", p, p)
				}
			}
		}()
		r.Expect("test message")
		t.Fatal("Expected panic")
	}
	// Test wrapError with Ok result (covers result.go:382)
	{
		r := gust.Ok[int](42)
		defer func() {
			p := recover()
			if p != nil {
				if err, ok := p.(error); ok {
					assert.Contains(t, err.Error(), "test message")
				} else {
					t.Errorf("Expected error panic, got %T: %v", p, p)
				}
			}
		}()
		r.Expect("test message")
		t.Fatal("Expected panic")
	}
}

// TestResult_AndThen tests AndThen method (covers result.go:487-492)
func TestResult_AndThen(t *testing.T) {
	// Test with Err result (should return r) (covers result.go:488-490)
	{
		r := gust.Err[int]("error")
		result := r.AndThen(func(i int) gust.Result[int] {
			return gust.Ok(i * 2)
		})
		assert.True(t, result.IsErr())
		assert.Equal(t, "error", result.Err().Error())
	}
	// Test with Ok result
	{
		r := gust.Ok[int](10)
		result := r.AndThen(func(i int) gust.Result[int] {
			return gust.Ok(i * 2)
		})
		assert.True(t, result.IsOk())
		assert.Equal(t, 20, result.Unwrap())
	}
}
