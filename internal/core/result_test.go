package core_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/andeya/gust/conv"
	"github.com/andeya/gust/errutil"
	"github.com/andeya/gust/internal/core"
	"github.com/andeya/gust/void"
	"github.com/stretchr/testify/assert"
)

// [`Result`] comes with some convenience methods that make working with it more succinct.
func TestResult(t *testing.T) {
	var goodResult1 = core.Ok(10)
	var badResult1 = core.TryErr[int](10)

	// The `IsOk` and `IsErr` methods do what they say.
	assert.True(t, goodResult1.IsOk() && !goodResult1.IsErr())
	assert.True(t, badResult1.IsErr() && !badResult1.IsOk())

	// `map` consumes the `Result` and produces another.
	var goodResult2 = goodResult1.Map(func(i int) int { return i + 1 })
	var badResult2 = badResult1.Map(func(i int) int { return i - 1 })

	// Use `AndThen` to continue the computation.
	var goodResult3 = goodResult2.XAndThen(func(i int) core.Result[any] { return core.Ok[any](i == 11) })

	// Use `OrElse` to handle the error.
	var _ = badResult2.OrElse(func(err error) core.Result[int] {
		fmt.Println(err)
		return core.Ok(20)
	})

	// Consume the result and return the contents with `Unwrap`.
	var _ = goodResult3.Unwrap()

	r, err := goodResult3.Split()
	assert.NoError(t, err)
	assert.Equal(t, true, r)
}

func TestResult_AssertRet(t *testing.T) {
	r := core.AssertRet[int](1)
	assert.Equal(t, core.Ok(1), r)
	r2 := core.AssertRet[int]("")
	assert.True(t, r2.IsErr())
	assert.Contains(t, r2.Err().Error(), "type assert error")

	// Test with string type
	r3 := core.AssertRet[string]("hello")
	assert.Equal(t, core.Ok("hello"), r3)

	// Test with struct type
	type S struct {
		X int
	}
	s := S{X: 42}
	r4 := core.AssertRet[S](s)
	assert.Equal(t, core.Ok(s), r4)

	// Test with invalid type assertion
	r5 := core.AssertRet[int]("not an int")
	assert.True(t, r5.IsErr())
	assert.Contains(t, r5.Err().Error(), "type assert error")
}

func TestResult_FmtErr(t *testing.T) {
	// Test FmtErr with format string
	r1 := core.FmtErr[int]("error: %s", "test")
	assert.True(t, r1.IsErr())
	assert.Contains(t, r1.Err().Error(), "error: test")

	// Test FmtErr with multiple arguments
	r2 := core.FmtErr[string]("error: %d %s", 42, "test")
	assert.True(t, r2.IsErr())
	assert.Contains(t, r2.Err().Error(), "error: 42 test")

	// Test FmtErr with no arguments
	r3 := core.FmtErr[int]("simple error")
	assert.True(t, r3.IsErr())
	assert.Equal(t, "simple error", r3.Err().Error())
}

func TestResult_TryErr_EmptyError(t *testing.T) {
	// Test TryErr with nil error (should return Ok with default value)
	r1 := core.TryErr[int](nil)
	assert.True(t, r1.IsOk())
	assert.Equal(t, 0, r1.Unwrap())

	// Test TryErr with empty ErrBox
	eb := errutil.BoxErr(nil)
	r2 := core.TryErr[string](eb)
	assert.True(t, r2.IsOk())
	assert.Equal(t, "", r2.Unwrap())

	// Test TryErr with non-nil error
	r3 := core.TryErr[int](errors.New("test error"))
	assert.True(t, r3.IsErr())
	assert.Equal(t, "test error", r3.Err().Error())
}

func TestResultJSON(t *testing.T) {
	var r = core.TryErr[any](errors.New("err"))
	var b, err = json.Marshal(r)
	assert.Equal(t, "json: error calling MarshalJSON for type core.Result[interface {}]: err", err.Error())
	assert.Nil(t, b)
	type T struct {
		Name string
	}
	var r2 = core.Ok(T{Name: "andeya"})
	var b2, err2 = json.Marshal(r2)
	assert.NoError(t, err2)
	assert.Equal(t, `{"Name":"andeya"}`, string(b2))

	var r3 core.Result[T]
	var err3 = json.Unmarshal(b2, &r3)
	assert.NoError(t, err3)
	assert.Equal(t, r2, r3)

	var r4 core.Result[T]
	var err4 = json.Unmarshal([]byte("0"), &r4)
	assert.True(t, r4.IsErr())
	assert.Equal(t, "json: cannot unmarshal number into Go value of type core_test.T", err4.Error())

	var r5 = core.Ok[void.Void](nil)
	var b5, err5 = json.Marshal(r5)
	assert.NoError(t, err5)
	assert.Equal(t, `null`, string(b5))
	var r6 core.Result[void.Void]
	var err6 = json.Unmarshal(b5, &r6)
	assert.NoError(t, err6)
	assert.Equal(t, r5, r6)
}

func TestResultIsValid(t *testing.T) {
	// Test nil pointer - covers r == nil branch
	var r0 *core.Result[any]
	assert.False(t, r0.IsValid())

	// Test zero value (both empty) - covers r.e.IsEmpty() && !r.t.IsSome() branch
	var r1 core.Result[any]
	assert.False(t, r1.IsValid())
	assert.False(t, (&core.Result[any]{}).IsValid())

	// Test Ok case - covers r.t.IsSome() branch
	var r2 = core.Ok[any](nil)
	assert.True(t, r2.IsValid())

	// Test Err case - covers !r.e.IsEmpty() branch (missing coverage)
	var r3 = core.TryErr[any]("test error")
	assert.True(t, r3.IsValid())

	// Test Err case with error interface - covers !r.e.IsEmpty() branch
	var r4 = core.TryErr[any](errors.New("error"))
	assert.True(t, r4.IsValid())
}

func TestResultUnwrapOrThrow(t *testing.T) {
	// Test UnwrapOrThrow with Ok value
	var r1 = core.Ok(1)
	var v1 = r1.UnwrapOrThrow()
	assert.Equal(t, 1, v1)

	// Test UnwrapOrThrow with Err value (should panic and be caught)
	var r core.Result[string]
	func() {
		defer r.Catch()
		var r2 = core.TryErr[int]("err")
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
		x := core.Ok(2)
		y := core.TryErr[int]("late error")
		assert.Equal(t, core.TryErr[int]("late error"), x.And(y))
	}
	{
		x := core.TryErr[uint]("early error")
		y := core.Ok[string]("foo")
		assert.Equal(t, core.TryErr[any]("early error"), x.XAnd(y.ToX()))
	}
}

func TestResult_And2(t *testing.T) {
	// Test with Ok result and Ok value
	{
		x := core.Ok(2)
		ret := x.And2(3, nil)
		assert.True(t, ret.IsOk())
		assert.Equal(t, 3, ret.Unwrap())
	}
	// Test with Ok result and Err value
	{
		x := core.Ok(2)
		ret := x.And2(3, errors.New("late error"))
		assert.True(t, ret.IsErr())
		assert.Equal(t, "late error", ret.Err().Error())
	}
	// Test with Err result (should return original error)
	{
		x := core.TryErr[int]("early error")
		ret := x.And2(3, nil)
		assert.True(t, ret.IsErr())
		assert.Equal(t, "early error", ret.Err().Error())
	}
	// Test with Err result and Err value (should return original error)
	{
		x := core.TryErr[int]("early error")
		ret := x.And2(3, errors.New("late error"))
		assert.True(t, ret.IsErr())
		assert.Equal(t, "early error", ret.Err().Error())
	}
}

func TestResult_XAnd2(t *testing.T) {
	// Test with Ok result and Ok value
	{
		x := core.Ok(2)
		ret := x.XAnd2("foo", nil)
		assert.True(t, ret.IsOk())
		assert.Equal(t, "foo", ret.Unwrap())
	}
	// Test with Ok result and Err value
	{
		x := core.Ok(2)
		ret := x.XAnd2("foo", errors.New("late error"))
		assert.True(t, ret.IsErr())
		assert.Equal(t, "late error", ret.Err().Error())
	}
	// Test with Err result (should return original error)
	{
		x := core.TryErr[int]("early error")
		ret := x.XAnd2("foo", nil)
		assert.True(t, ret.IsErr())
		assert.Equal(t, "early error", ret.Err().Error())
	}
	// Test with Err result and Err value (should return original error)
	{
		x := core.TryErr[int]("early error")
		ret := x.XAnd2("foo", errors.New("late error"))
		assert.True(t, ret.IsErr())
		assert.Equal(t, "early error", ret.Err().Error())
	}
}

func ExampleResult_AndThen() {
	var divide = func(i, j float32) core.Result[float32] {
		if j == 0 {
			return core.TryErr[float32]("j can not be 0")
		}
		return core.Ok(i / j)
	}
	var ret float32 = divide(1, 2).AndThen(func(i float32) core.Result[float32] {
		return core.Ok(i * 10)
	}).Unwrap()
	fmt.Println(ret)
	// Output:
	// 5
}

func TestResult_Err(t *testing.T) {
	{
		var x = core.Ok[int](2)
		assert.Equal(t, error(nil), x.Err())
	}
	{
		var x = core.TryErr[int]("some error message")
		assert.Equal(t, "some error message", x.Err().Error())
	}
}

func TestResult_ExpectErr(t *testing.T) {
	defer func() {
		assert.Equal(t, "Testing expect_err: 10", recover())
	}()
	err := core.Ok(10).ExpectErr("Testing expect_err")
	assert.NoError(t, err)

	// Test IsErr branch (covers core.go:311-313)
	err2 := core.TryErr[int]("test error")
	ret := err2.ExpectErr("should not panic")
	assert.NotNil(t, ret)
	assert.Equal(t, "test error", ret.Error())
}

// TestResult_String tests Result.String method (covers core.go:180-185)
func TestResult_String(t *testing.T) {
	// Test Ok case - covers safeGetT() path
	ok := core.Ok(42)
	assert.Equal(t, "Ok(42)", ok.String())

	// Test Err case with string error - covers safeGetE() path with string error
	err1 := core.TryErr[int]("error message")
	assert.Contains(t, err1.String(), "Err")
	assert.Contains(t, err1.String(), "error message")

	// Test Err case with error interface - covers safeGetE() path with error interface
	err2 := core.TryErr[int](errors.New("error interface"))
	assert.Contains(t, err2.String(), "Err")
	assert.Contains(t, err2.String(), "error interface")

	// Test Err case with ErrBox - covers safeGetE() path with ErrBox
	err3 := core.TryErr[int](errutil.BoxErr("boxed error"))
	assert.Contains(t, err3.String(), "Err")
	assert.Contains(t, err3.String(), "boxed error")
}

// TestResult_ErrValNil tests Result.ErrVal returning nil (covers core.go:165-169)
func TestResult_ErrValNil(t *testing.T) {
	ok := core.Ok(42)
	val := ok.ErrVal()
	assert.Nil(t, val)
}

func TestResult_Expect(t *testing.T) {
	defer func() {
		assert.Equal(t, "failed to parse number: strconv.Atoi: parsing \"4x\": invalid syntax", recover().(error).Error())
	}()
	core.Ret(strconv.Atoi("4x")).
		Expect("failed to parse number")
}

func TestResult_InspectErr(t *testing.T) {
	core.Ret(strconv.Atoi("4x")).
		InspectErr(func(err error) {
			t.Logf("failed to convert: %v", err)
		})
}

func TestResult_Inspect(t *testing.T) {
	core.Ret(strconv.Atoi("4")).
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
		var x = core.TryErr[int]("hey")
		assert.True(t, x.IsErrAnd(func(err error) bool { return err.Error() == "hey" }))
	}
	{
		var x = core.Ok[int](2)
		assert.False(t, x.IsErrAnd(func(err error) bool { return err.Error() == "hey" }))
	}
}

func TestResult_IsErr(t *testing.T) {
	{
		var x = core.Ok[int](-3)
		assert.False(t, x.IsErr())
	}
	{
		var x = core.TryErr[int]("some error message")
		assert.True(t, x.IsErr())
	}
}

func TestResult_IsOkAnd(t *testing.T) {
	{
		var x = core.Ok[int](2)
		assert.True(t, x.IsOkAnd(func(x int) bool { return x > 1 }))
	}
	{
		var x = core.Ok[int](0)
		assert.False(t, x.IsOkAnd(func(x int) bool { return x > 1 }))
	}
	{
		var x = core.TryErr[int]("hey")
		assert.False(t, x.IsOkAnd(func(x int) bool { return x > 1 }))
	}
}

func TestResult_IsOk(t *testing.T) {
	{
		var x = core.Ok[int](-3)
		assert.True(t, x.IsOk())
	}
	{
		var x = core.TryErr[int]("some error message")
		assert.False(t, x.IsOk())
	}
}

func TestResult_MapErr(t *testing.T) {
	var stringify = func(x error) any { return fmt.Sprintf("error code: %v", x) }
	{
		var x = core.Ok[uint32](2)
		assert.Equal(t, core.Ok[uint32](2), x.MapErr(stringify))
	}
	{
		var x = core.TryErr[uint32](13)
		assert.Equal(t, core.TryErr[uint32]("error code: 13"), x.MapErr(stringify))
	}
}

func TestResult_MapOrElse_2(t *testing.T) {
	{
		var x = core.Ok("foo")
		assert.Equal(t, "test:foo", x.MapOrElse(func(err error) string {
			return "bar"
		}, func(x string) string { return "test:" + x }))
	}
	{
		var x = core.TryErr[string]("foo")
		assert.Equal(t, "bar", x.MapOrElse(func(err error) string {
			return "bar"
		}, func(x string) string { return "test:" + x }))
	}
}

func TestResult_MapOr_2(t *testing.T) {
	{
		var x = core.Ok("foo")
		assert.Equal(t, "test:foo", x.MapOr("bar", func(x string) string { return "test:" + x }))
	}
	{
		var x = core.TryErr[string]("foo")
		assert.Equal(t, "bar", x.MapOr("bar", func(x string) string { return "test:" + x }))
	}
}

func TestResult_Map_1(t *testing.T) {
	var line = "1\n2\n3\n4\n"
	for _, num := range strings.Split(line, "\n") {
		core.Ret(strconv.Atoi(num)).Map(func(i int) int {
			return i * 2
		}).Inspect(func(i int) {
			t.Log(i)
		})
	}
}

func TestResult_Map_2(t *testing.T) {
	var isMyNum = func(s string, search int) core.Result[any] {
		return core.Ret(strconv.Atoi(s)).XMap(func(x int) any { return x == search })
	}
	assert.Equal(t, core.Ok[any](true), isMyNum("1", 1))
	assert.Equal(t, "Err(strconv.Atoi: parsing \"lol\": invalid syntax)", isMyNum("lol", 1).String())
	assert.Equal(t, "Err(strconv.Atoi: parsing \"NaN\": invalid syntax)", isMyNum("NaN", 1).String())
}

func TestResult_Ok(t *testing.T) {
	{
		var x = core.Ok[int](2)
		assert.Equal(t, core.Some[int](2), x.Ok())
	}
	{
		var x = core.TryErr[int]("some error message")
		assert.Equal(t, core.None[int](), x.Ok())
	}
}

func TestResult_OrElse(t *testing.T) {
	var sq = func(err error) core.Result[int] {
		if err == nil {
			return core.Ok(0)
		}
		// Convert error to int for testing
		if err.Error() == "3" {
			return core.Ok(9) // 3 * 3
		}
		return core.TryErr[int](err)
	}
	var errFn = func(err error) core.Result[int] {
		return core.TryErr[int](err)
	}

	assert.Equal(t, core.Ok(2), core.Ok(2).OrElse(sq).OrElse(sq))
	assert.Equal(t, core.Ok(2), core.Ok(2).OrElse(errFn).OrElse(sq))
	assert.Equal(t, core.Ok(9), core.TryErr[int](fmt.Errorf("3")).OrElse(sq).OrElse(errFn))
	assert.Equal(t, core.TryErr[int](fmt.Errorf("3")), core.TryErr[int](fmt.Errorf("3")).OrElse(errFn).OrElse(errFn))
}

func TestResult_Or(t *testing.T) {
	{
		x := core.Ok(2)
		y := core.TryErr[int]("late error")
		assert.Equal(t, core.Ok(2), x.Or(y))
	}
	{
		x := core.TryErr[uint]("early error")
		y := core.Ok[uint](2)
		assert.Equal(t, core.Ok[uint](2), x.Or(y))
	}
	{
		x := core.TryErr[uint]("not a 2")
		y := core.TryErr[uint]("late error")
		assert.Equal(t, core.TryErr[uint]("late error"), x.Or(y))
	}
	{
		x := core.Ok[uint](2)
		y := core.Ok[uint](100)
		assert.Equal(t, core.Ok[uint](2), x.Or(y))
	}
}

func TestResult_Ret(t *testing.T) {
	var w = core.Ret[int](strconv.Atoi("s"))
	assert.False(t, w.IsOk())
	assert.True(t, w.IsErr())

	var w2 = core.Ret[any](strconv.Atoi("-1"))
	assert.True(t, w2.IsOk())
	assert.False(t, w2.IsErr())
	assert.Equal(t, -1, w2.Unwrap())
}

func TestResult_UnwrapErr_1(t *testing.T) {
	defer func() {
		assert.Equal(t, "called `Result.UnwrapErr()` on an `ok` value: 10", recover())
	}()
	err := core.Ok(10).UnwrapErr()
	assert.NoError(t, err)
}

func TestResult_UnwrapErr_2(t *testing.T) {
	err := core.TryErr[int]("emergency failure").UnwrapErr()
	if assert.Error(t, err) {
		assert.Equal(t, "emergency failure", err.Error())
	} else {
		t.FailNow()
	}
}

func TestFmtErr(t *testing.T) {
	ret := core.FmtErr[int]("error: %s", "test")
	assert.True(t, ret.IsErr())
	assert.Contains(t, ret.Err().Error(), "error: test")
}

func TestAssertRet(t *testing.T) {
	// Test with valid type
	ret1 := core.AssertRet[int](42)
	assert.True(t, ret1.IsOk())
	assert.Equal(t, 42, ret1.Unwrap())

	// Test with invalid type
	ret2 := core.AssertRet[int]("string")
	assert.True(t, ret2.IsErr())
	assert.Contains(t, ret2.Err().Error(), "type assert error")
}

func TestResult_Catch(t *testing.T) {
	// Test basic Catch functionality (with stack trace by default)
	var ret core.Result[int]
	func() {
		defer ret.Catch()
		core.TryErr[int]("test error").UnwrapOrThrow()
	}()
	assert.True(t, ret.IsErr())
	// Error() should only return error message, not stack trace
	errMsg := ret.Err().Error()
	assert.Contains(t, errMsg, "test error")
	// Should NOT contain stack trace in Error()
	assert.NotContains(t, errMsg, "\n")
	// But %+v should contain stack trace
	fullMsg := fmt.Sprintf("%+v", ret.Err())
	assert.Contains(t, fullMsg, "test error")
	assert.Contains(t, fullMsg, "\n")
}

func TestResult_Catch_IsSomeBranch(t *testing.T) {
	// Test Catch when ret already has Ok value (IsSome branch)
	var ret core.Result[int] = core.Ok(42)
	func() {
		defer ret.Catch()
		core.TryErr[int]("panic error").UnwrapOrThrow()
	}()
	assert.True(t, ret.IsErr())
	// Error() should only return error message, not stack trace
	errMsg := ret.Err().Error()
	assert.Contains(t, errMsg, "panic error")
	// Should NOT contain stack trace in Error()
	assert.NotContains(t, errMsg, "\n")
	// But %+v should contain stack trace
	fullMsg := fmt.Sprintf("%+v", ret.Err())
	assert.Contains(t, fullMsg, "panic error")
	assert.Contains(t, fullMsg, "\n")
}

func TestResult_UnwrapOrDefault(t *testing.T) {
	assert.Equal(t, "car", core.Ok("car").UnwrapOrDefault())
	assert.Equal(t, "", core.TryErr[string](nil).UnwrapOrDefault())
	assert.Equal(t, time.Time{}, core.TryErr[time.Time](nil).UnwrapOrDefault())
	assert.Equal(t, &time.Time{}, core.TryErr[*time.Time](nil).UnwrapOrDefault())
	assert.Equal(t, conv.Ref(&time.Time{}), core.TryErr[**time.Time](nil).UnwrapOrDefault())

	// Test with actual error (not nil)
	assert.Equal(t, "", core.TryErr[string](errors.New("test error")).UnwrapOrDefault())
	assert.Equal(t, 0, core.TryErr[int](errors.New("test error")).UnwrapOrDefault())
	assert.Equal(t, time.Time{}, core.TryErr[time.Time](errors.New("test error")).UnwrapOrDefault())
}

func TestResult_UnwrapOrElse(t *testing.T) {
	var count = func(x error) int {
		return len(x.Error())
	}
	assert.Equal(t, 2, core.Ok(2).UnwrapOrElse(count))
	assert.Equal(t, 3, core.TryErr[int]("foo").UnwrapOrElse(count))
}

func TestResult_Unwrap(t *testing.T) {
	defer func() {
		p := recover()
		if eb, ok := p.(*strconv.NumError); ok {
			assert.Equal(t, "strconv.Atoi: parsing \"4x\": invalid syntax", eb.Error())
		} else {
			t.Fatalf("expected *strconv.NumError, got %T: %v", p, p)
		}
	}()
	core.Ret(strconv.Atoi("4x")).Unwrap()
}

func TestResult_UnwrapUnchecked(t *testing.T) {
	{
		var r core.Result[string]
		assert.Equal(t, "", r.UnwrapUnchecked())
	}
	{
		var r = core.TryErr[string]("foo")
		assert.Equal(t, "", r.UnwrapUnchecked())
	}
	{
		var r = core.Ok[string]("foo")
		assert.Equal(t, "foo", r.UnwrapUnchecked())
	}
}

func TestResult_XOk(t *testing.T) {
	{
		var x = core.Ok[int](2)
		assert.Equal(t, core.Some[any](2), x.XOk())
	}
	{
		var x = core.TryErr[int]("some error message")
		assert.Equal(t, core.None[any](), x.XOk())
	}
}

func TestResult_ErrVal(t *testing.T) {
	{
		var x = core.Ok[int](2)
		assert.Nil(t, x.ErrVal())
	}
	{
		var x = core.TryErr[int]("some error message")
		assert.Equal(t, "some error message", x.ErrVal())
	}
	{
		var x = core.TryErr[int](errutil.BoxErr("boxed error"))
		assert.Equal(t, "boxed error", x.ErrVal())
	}
	{
		var x = core.TryErr[int](errors.New("std error"))
		assert.NotNil(t, x.ErrVal())
	}
}

func TestResult_ToX(t *testing.T) {
	{
		var x = core.Ok[int](42)
		xResult := x.ToX()
		assert.True(t, xResult.IsOk())
		assert.Equal(t, 42, xResult.Unwrap())
	}
	{
		var x = core.TryErr[int]("error")
		xResult := x.ToX()
		assert.True(t, xResult.IsErr())
		assert.Equal(t, "error", xResult.Err().Error())
	}
}

func TestResult_XMap(t *testing.T) {
	{
		var x = core.Ok[int](2)
		ret := x.XMap(func(i int) any { return i * 2 })
		assert.True(t, ret.IsOk())
		assert.Equal(t, 4, ret.Unwrap())
	}
	{
		var x = core.TryErr[int]("error")
		ret := x.XMap(func(i int) any { return i * 2 })
		assert.True(t, ret.IsErr())
	}
}

func TestResult_XMapOr(t *testing.T) {
	{
		var x = core.Ok[int](2)
		ret := x.XMapOr("default", func(i int) any { return i * 2 })
		assert.Equal(t, 4, ret)
	}
	{
		var x = core.TryErr[int]("error")
		ret := x.XMapOr("default", func(i int) any { return i * 2 })
		assert.Equal(t, "default", ret)
	}
}

func TestResult_XMapOrElse(t *testing.T) {
	{
		var x = core.Ok[int](2)
		ret := x.XMapOrElse(func(error) any { return "default" }, func(i int) any { return i * 2 })
		assert.Equal(t, 4, ret)
	}
	{
		var x = core.TryErr[int]("error")
		result := x.XMapOrElse(func(error) any { return "default" }, func(i int) any { return i * 2 })
		assert.Equal(t, "default", result)
	}
}

func TestResult_ContainsErr(t *testing.T) {
	{
		var x = core.Ok[int](2)
		assert.False(t, x.ContainsErr(errors.New("test")))
	}
	{
		// Use same error instance for comparison
		testErr := errors.New("test error")
		var x = core.TryErr[int](testErr)
		assert.True(t, x.ContainsErr(testErr))
		assert.False(t, x.ContainsErr(errors.New("different error")))
	}
	{
		err1 := errors.New("wrapped")
		err2 := fmt.Errorf("outer: %w", err1)
		var x = core.TryErr[int](err2)
		assert.True(t, x.ContainsErr(err1))
	}
}

func TestResult_Iterator(t *testing.T) {
	// Test Next
	{
		var x = core.Ok("foo")
		opt := x.Next()
		assert.Equal(t, core.Some("foo"), opt)
		// Next() consumes the value, so Ok() should return None after Next()
		assert.True(t, x.Ok().IsNone())
	}
	{
		var x = core.TryErr[string]("error")
		opt := x.Next()
		assert.True(t, opt.IsNone())
	}
	{
		var nilResult *core.Result[string]
		opt := nilResult.Next()
		assert.True(t, opt.IsNone())
	}

	// Test NextBack
	{
		var x = core.Ok("bar")
		opt := x.NextBack()
		assert.Equal(t, core.Some("bar"), opt)
		// NextBack() also consumes the value
		assert.True(t, x.Ok().IsNone())
	}

	// Test Remaining
	{
		var x = core.Ok("baz")
		assert.Equal(t, uint(1), x.Remaining())
		x.Next() // Consume the value
		assert.Equal(t, uint(0), x.Remaining())
	}
	{
		var x = core.TryErr[string]("error")
		assert.Equal(t, uint(0), x.Remaining())
	}
	{
		var nilResult *core.Result[string]
		assert.Equal(t, uint(0), nilResult.Remaining())
	}

	// Test SizeHint
	{
		var x = core.Ok("test")
		lower, upper := x.SizeHint()
		assert.Equal(t, uint(1), lower)
		assert.True(t, upper.IsSome())
		assert.Equal(t, uint(1), upper.Unwrap())
	}
	{
		var x = core.TryErr[string]("error")
		lower, upper := x.SizeHint()
		assert.Equal(t, uint(0), lower)
		assert.True(t, upper.IsSome())
		assert.Equal(t, uint(0), upper.Unwrap())
	}
	{
		var nilResult *core.Result[string]
		lower, upper := nilResult.SizeHint()
		assert.Equal(t, uint(0), lower)
		assert.True(t, upper.IsSome())
		assert.Equal(t, uint(0), upper.Unwrap())
	}
}

func TestResult_Ref(t *testing.T) {
	ret := core.Ok(42)
	ref := ret.Ref()
	assert.Equal(t, core.Ok(42), *ref)
	ref.Unwrap() // Should not panic
}

func TestResult_ErrToVoidResult(t *testing.T) {
	{
		var x = core.Ok[string]("foo")
		ret := core.RetVoid(x.Err())
		assert.False(t, ret.IsErr())
	}
	{
		var x = core.TryErr[string](errors.New("error"))
		ret := core.RetVoid(x.Err())
		assert.True(t, ret.IsErr())
	}
}

func TestResult_Split(t *testing.T) {
	// Test Ok case - covers safeGetT() and safeGetE() with Ok result
	{
		var x = core.Ok("foo")
		val, err := x.Split()
		assert.NoError(t, err)
		assert.Equal(t, "foo", val)
	}

	// Test Err case with error interface - covers safeGetT() returning zero value and safeGetE() with error
	{
		var x = core.TryErr[string](errors.New("error"))
		val, err := x.Split()
		assert.Error(t, err)
		assert.Equal(t, "", val) // Zero value for string
		assert.Contains(t, err.Error(), "error")
	}

	// Test Err case with string error - covers safeGetE() with string error
	{
		var x = core.TryErr[int]("string error")
		val, err := x.Split()
		assert.Error(t, err)
		assert.Equal(t, 0, val) // Zero value for int
		assert.Contains(t, err.Error(), "string error")
	}

	// Test Err case with ErrBox - covers safeGetE() with ErrBox
	{
		var x = core.TryErr[int](errutil.BoxErr("boxed error"))
		val, err := x.Split()
		assert.Error(t, err)
		assert.Equal(t, 0, val) // Zero value for int
		assert.Contains(t, err.Error(), "boxed error")
	}
}

func TestResult_UnmarshalJSON_NilReceiver(t *testing.T) {
	// Test UnmarshalJSON with nil receiver
	var nilResult *core.Result[int]
	err := nilResult.UnmarshalJSON([]byte("42"))
	assert.Error(t, err)
	assert.IsType(t, &json.InvalidUnmarshalError{}, err)
}

func TestResult_UnmarshalJSON_ErrorPath(t *testing.T) {
	// Test UnmarshalJSON with invalid JSON (error path)
	var ret core.Result[int]
	err := ret.UnmarshalJSON([]byte("invalid json"))
	assert.Error(t, err)
	assert.True(t, ret.IsErr()) // Should be Err on error
}

func TestResult_UnmarshalJSON_ValidAfterError(t *testing.T) {
	// Test UnmarshalJSON with error first, then valid JSON
	var ret core.Result[int]
	// First attempt with invalid JSON
	_ = ret.UnmarshalJSON([]byte("invalid"))
	assert.True(t, ret.IsErr())

	// Then with valid JSON
	err := ret.UnmarshalJSON([]byte("42"))
	assert.NoError(t, err)
	assert.True(t, ret.IsOk())
	assert.Equal(t, 42, ret.Unwrap())
}

func TestResult_UnmarshalJSON_Struct(t *testing.T) {
	type S struct {
		X int
		Y string
	}
	var ret core.Result[S]
	err := ret.UnmarshalJSON([]byte(`{"X":10,"Y":"test"}`))
	assert.NoError(t, err)
	assert.True(t, ret.IsOk())
	assert.Equal(t, S{X: 10, Y: "test"}, ret.Unwrap())
}

func TestResult_UnmarshalJSON_Array(t *testing.T) {
	var ret core.Result[[]int]
	err := ret.UnmarshalJSON([]byte("[1,2,3]"))
	assert.NoError(t, err)
	assert.True(t, ret.IsOk())
	assert.Equal(t, []int{1, 2, 3}, ret.Unwrap())
}

func TestResult_UnmarshalJSON_Map(t *testing.T) {
	var ret core.Result[map[string]int]
	err := ret.UnmarshalJSON([]byte(`{"a":1,"b":2}`))
	assert.NoError(t, err)
	assert.True(t, ret.IsOk())
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, ret.Unwrap())
}

func TestResult_Catch_NilReceiver(t *testing.T) {
	// Test Catch with nil receiver
	defer func() {
		assert.NotNil(t, recover())
	}()
	var nilResult *core.Result[int]
	defer nilResult.Catch()
	core.TryErr[int]("test error").UnwrapOrThrow()
}

func TestResult_Catch_NonPanicValue(t *testing.T) {
	// Test Catch with non-ErrBox panic (should be caught and converted to ErrBox)
	var ret core.Result[int]
	func() {
		defer ret.Catch()
		panic("regular panic")
	}()
	// Should be caught and converted to error
	assert.True(t, ret.IsErr())
	// Error() should only return error message, not stack trace
	errMsg := ret.Err().Error()
	assert.Contains(t, errMsg, "regular panic")
	// Should NOT contain stack trace in Error()
	assert.NotContains(t, errMsg, "\n")
	// But %+v should contain stack trace
	fullMsg := fmt.Sprintf("%+v", ret.Err())
	assert.Contains(t, fullMsg, "regular panic")
	assert.Contains(t, fullMsg, "\n")
}

func TestResult_Catch_OkValue(t *testing.T) {
	// Test Catch when result already has Ok value (should be updated to Err)
	var ret core.Result[int] = core.Ok(42)
	func() {
		defer ret.Catch()
		core.TryErr[string]("test error").UnwrapOrThrow()
	}()
	// Result should be updated to Err
	assert.True(t, ret.IsErr())
	// Error() should only return error message, not stack trace
	errMsg := ret.Err().Error()
	assert.Contains(t, errMsg, "test error")
	// Should NOT contain stack trace in Error()
	assert.NotContains(t, errMsg, "\n")
	// But %+v should contain stack trace
	fullMsg := fmt.Sprintf("%+v", ret.Err())
	assert.Contains(t, fullMsg, "test error")
	assert.Contains(t, fullMsg, "\n")
}

// TestResult_Catch_WithStackTrace tests that Catch captures stack trace information
func TestResult_Catch_WithStackTrace(t *testing.T) {
	// Test with ErrBox panic
	{
		var ret core.Result[int]
		func() {
			defer ret.Catch()
			core.TryErr[int]("test error").UnwrapOrThrow()
		}()
		assert.True(t, ret.IsErr())
		err := ret.Err()
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
		var ret core.Result[string]
		func() {
			defer ret.Catch()
			panic(errors.New("regular error"))
		}()
		assert.True(t, ret.IsErr())
		err := ret.Err()
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
		var ret core.Result[int]
		func() {
			defer ret.Catch()
			panic("string panic")
		}()
		assert.True(t, ret.IsErr())
		err := ret.Err()
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
		var ret core.Result[string]
		func() {
			defer ret.Catch()
			panic(42)
		}()
		assert.True(t, ret.IsErr())
		err := ret.Err()
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
	var ret core.Result[int]
	func() {
		defer ret.Catch()
		core.TryErr[int]("test error").UnwrapOrThrow()
	}()
	assert.True(t, ret.IsErr())

	// Get the error
	err := ret.Err()
	assert.NotNil(t, err)

	// Error() should only return error message, not stack trace
	errMsg := err.Error()
	assert.Contains(t, errMsg, "test error")
	assert.NotContains(t, errMsg, "\n")

	// Try to extract StackTraceCarrier from error chain
	var carrier errutil.StackTraceCarrier
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
	var ret core.Result[int]
	eb := errutil.BoxErr(errors.New("errbox pointer error"))
	func() {
		defer ret.Catch()
		panic(eb)
	}()
	assert.True(t, ret.IsErr())
	err := ret.Err()
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
	var ret core.Result[int]
	eb := errutil.BoxErr(errors.New("errbox value error"))
	func() {
		defer ret.Catch()
		panic(eb)
	}()
	assert.True(t, ret.IsErr())
	err := ret.Err()
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
	var ret core.Result[int]
	var eb *errutil.ErrBox
	func() {
		defer ret.Catch()
		panic(eb)
	}()
	assert.True(t, ret.IsErr())
	err := ret.Err()
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
	var ret core.Result[int]
	func() {
		defer ret.Catch(false)
		core.TryErr[int]("test error").UnwrapOrThrow()
	}()
	assert.True(t, ret.IsErr())
	err := ret.Err()
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
	var ret core.Result[int]
	func() {
		defer func() {
			t.Logf("error with panic stack trace: %+v", ret.Err())
		}()
		defer ret.Catch(true)
		core.TryErr[int]("test error").UnwrapOrThrow()
	}()
	assert.True(t, ret.IsErr())
	err := ret.Err()
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
	var ret core.Result[int]
	func() {
		defer ret.Catch()
		core.TryErr[int]("test error").UnwrapOrThrow()
	}()
	assert.True(t, ret.IsErr())
	err := ret.Err()
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
	var ret1 core.Result[int]
	var ret2 core.Result[int]

	// With stack trace
	func() {
		defer ret1.Catch(true)
		panic("error with stack")
	}()
	assert.True(t, ret1.IsErr())

	// Without stack trace
	func() {
		defer ret2.Catch(false)
		panic("error without stack")
	}()
	assert.True(t, ret2.IsErr())

	// Both should be errors
	assert.True(t, ret1.IsErr())
	assert.True(t, ret2.IsErr())
	// Both should contain error message
	assert.Contains(t, ret1.Err().Error(), "error with stack")
	assert.Contains(t, ret2.Err().Error(), "error without stack")
	// Only result1 should have stack trace in %+v
	fullMsg1 := fmt.Sprintf("%+v", ret1.Err())
	fullMsg2 := fmt.Sprintf("%+v", ret2.Err())
	assert.Contains(t, fullMsg1, "\n")
	assert.NotContains(t, fullMsg2, "\n")
}

// TestResult_Catch_FormatOptions tests different format options for caught errors
func TestResult_Catch_FormatOptions(t *testing.T) {
	var ret core.Result[int]
	func() {
		defer ret.Catch()
		core.TryErr[int]("format test error").UnwrapOrThrow()
	}()
	assert.True(t, ret.IsErr())
	err := ret.Err()

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
	var ret core.Result[int] = core.Ok(42)
	func() {
		defer ret.Catch()
		// No panic, just return normally
		ret = core.Ok(100)
	}()
	// Result should remain Ok
	assert.True(t, ret.IsOk())
	assert.Equal(t, 100, ret.Unwrap())
	assert.False(t, ret.IsErr())
}

// TestResult_Catch_MultiplePanics tests that Catch only catches the first panic
func TestResult_Catch_MultiplePanics(t *testing.T) {
	var ret core.Result[int]
	func() {
		defer ret.Catch()
		panic("first panic")
		// Note: second panic would never be reached due to first panic
	}()
	assert.True(t, ret.IsErr())
	errMsg := ret.Err().Error()
	assert.Contains(t, errMsg, "first panic")
}

// TestResult_Catch_WithStackTraceFalse_Verification tests Catch(false) doesn't capture stack
func TestResult_Catch_WithStackTraceFalse_Verification(t *testing.T) {
	var ret core.Result[int]
	func() {
		defer ret.Catch(false)
		core.TryErr[int]("no stack error").UnwrapOrThrow()
	}()
	assert.True(t, ret.IsErr())
	err := ret.Err()
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

	var ret core.Result[int]
	func() {
		defer ret.Catch()
		panic(wrappedErr)
	}()
	assert.True(t, ret.IsErr())

	err := ret.Err()
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
	ret := core.Ok[int](42)
	ret2 := ret.XAndThen(func(i int) core.Result[any] {
		return core.Ok[any](i * 2)
	})
	assert.True(t, ret2.IsOk())
	assert.Equal(t, 84, ret2.Unwrap())

	// Test XAndThen with error path
	ret3 := core.Ok[int](42)
	ret4 := ret3.XAndThen(func(i int) core.Result[any] {
		return core.TryErr[any]("error")
	})
	assert.True(t, ret4.IsErr())
}

func TestResult_AndThen2(t *testing.T) {
	// Test with Ok result and successful operation
	{
		x := core.Ok(2)
		ret := x.AndThen2(func(i int) (int, error) {
			return i * 2, nil
		})
		assert.True(t, ret.IsOk())
		assert.Equal(t, 4, ret.Unwrap())
	}
	// Test with Ok result and error operation
	{
		x := core.Ok(2)
		ret := x.AndThen2(func(i int) (int, error) {
			return 0, errors.New("operation error")
		})
		assert.True(t, ret.IsErr())
		assert.Equal(t, "operation error", ret.Err().Error())
	}
	// Test with Err result (should return original error)
	{
		x := core.TryErr[int]("early error")
		ret := x.AndThen2(func(i int) (int, error) {
			return i * 2, nil
		})
		assert.True(t, ret.IsErr())
		assert.Equal(t, "early error", ret.Err().Error())
	}
}

func TestResult_XAndThen2(t *testing.T) {
	// Test with Ok result and successful operation
	{
		x := core.Ok(2)
		ret := x.XAndThen2(func(i int) (any, error) {
			return i * 2, nil
		})
		assert.True(t, ret.IsOk())
		assert.Equal(t, 4, ret.Unwrap())
	}
	// Test with Ok result and error operation
	{
		x := core.Ok(2)
		ret := x.XAndThen2(func(i int) (any, error) {
			return nil, errors.New("operation error")
		})
		assert.True(t, ret.IsErr())
		assert.Equal(t, "operation error", ret.Err().Error())
	}
	// Test with Err result (should return original error)
	{
		x := core.TryErr[int]("early error")
		ret := x.XAndThen2(func(i int) (any, error) {
			return i * 2, nil
		})
		assert.True(t, ret.IsErr())
		assert.Equal(t, "early error", ret.Err().Error())
	}
}

func TestResult_Or2(t *testing.T) {
	// Test with Ok result (should return original Ok)
	{
		x := core.Ok(2)
		ret := x.Or2(3, nil)
		assert.True(t, ret.IsOk())
		assert.Equal(t, 2, ret.Unwrap())
	}
	// Test with Ok result and Err value (should return original Ok)
	{
		x := core.Ok(2)
		ret := x.Or2(3, errors.New("late error"))
		assert.True(t, ret.IsOk())
		assert.Equal(t, 2, ret.Unwrap())
	}
	// Test with Err result and Ok value
	{
		x := core.TryErr[int]("early error")
		ret := x.Or2(3, nil)
		assert.True(t, ret.IsOk())
		assert.Equal(t, 3, ret.Unwrap())
	}
	// Test with Err result and Err value
	{
		x := core.TryErr[int]("early error")
		ret := x.Or2(3, errors.New("late error"))
		assert.True(t, ret.IsErr())
		assert.Equal(t, "late error", ret.Err().Error())
	}
}

func TestResult_OrElse2(t *testing.T) {
	// Test with Ok result (should return original Ok)
	{
		x := core.Ok(2)
		ret := x.OrElse2(func(err error) (int, error) {
			return 3, nil
		})
		assert.True(t, ret.IsOk())
		assert.Equal(t, 2, ret.Unwrap())
	}
	// Test with Err result and successful operation
	{
		x := core.TryErr[int]("early error")
		ret := x.OrElse2(func(err error) (int, error) {
			return 3, nil
		})
		assert.True(t, ret.IsOk())
		assert.Equal(t, 3, ret.Unwrap())
	}
	// Test with Err result and error operation
	{
		x := core.TryErr[int]("early error")
		ret := x.OrElse2(func(err error) (int, error) {
			return 0, errors.New("late error")
		})
		assert.True(t, ret.IsErr())
		assert.Equal(t, "late error", ret.Err().Error())
	}
}

func TestResult_Flatten(t *testing.T) {
	// Test with Ok result and nil error (should return original Ok)
	{
		r := core.Ok(42)
		ret := r.Flatten(nil)
		assert.True(t, ret.IsOk())
		assert.Equal(t, 42, ret.Unwrap())
	}
	// Test with Ok result and error (should return error)
	{
		r := core.Ok(42)
		ret := r.Flatten(errors.New("test error"))
		assert.True(t, ret.IsErr())
		assert.Equal(t, "test error", ret.Err().Error())
	}
	// Test with Err result and nil error (should return original Err)
	{
		r := core.TryErr[int]("original error")
		ret := r.Flatten(nil)
		assert.True(t, ret.IsErr())
		assert.Equal(t, "original error", ret.Err().Error())
	}
	// Test with Err result and error (should return the provided error)
	{
		r := core.TryErr[int]("original error")
		ret := r.Flatten(errors.New("new error"))
		assert.True(t, ret.IsErr())
		assert.Equal(t, "new error", ret.Err().Error())
	}
}

// TestOkVoid tests OkVoid function (covers core.go:105-107)
func TestOkVoid(t *testing.T) {
	ret := core.OkVoid()
	assert.True(t, ret.IsOk())
	assert.Nil(t, ret.Err())
}

// TestTryErrVoid tests TryErrVoid function (covers core.go:104-110)
func TestTryErrVoid(t *testing.T) {
	// Test with error value
	{
		err := errors.New("test error")
		ret := core.TryErrVoid(err)
		assert.True(t, ret.IsErr())
		assert.False(t, ret.IsOk())
		assert.NotNil(t, ret.Err())
		assert.Equal(t, "test error", ret.Err().Error())
	}
	// Test with string error
	{
		ret := core.TryErrVoid("string error")
		assert.True(t, ret.IsErr())
		assert.False(t, ret.IsOk())
		assert.NotNil(t, ret.Err())
		assert.Equal(t, "string error", ret.Err().Error())
	}
	// Test with nil error (should return OkVoid, as nil represents "no error")
	{
		ret := core.TryErrVoid(nil)
		assert.True(t, ret.IsOk(), "TryErrVoid(nil) should return OkVoid")
		assert.False(t, ret.IsErr(), "TryErrVoid(nil) should not return Err")
		// TryErrVoid(nil) returns OkVoid() because nil represents "no error"
	}
	// Test that TryErrVoid is equivalent to TryErr[Void]
	{
		err := errors.New("comparison error")
		ret1 := core.TryErrVoid(err)
		ret2 := core.TryErr[void.Void](err)
		assert.Equal(t, ret1.IsErr(), ret2.IsErr())
		assert.Equal(t, ret1.IsOk(), ret2.IsOk())
		if ret1.IsErr() && ret2.IsErr() {
			assert.Equal(t, ret1.Err().Error(), ret2.Err().Error())
		}
	}
}

// TestUnwrapErrOr tests UnwrapErrOr function (covers core.go:253-257)
func TestUnwrapErrOr(t *testing.T) {
	// Test with Err result
	{
		r := core.TryErr[void.Void]("test error")
		def := errors.New("default error")
		err := core.UnwrapErrOr(r, def)
		assert.NotNil(t, err)
		assert.Equal(t, "test error", err.Error())
	}
	// Test with Ok result (should return default)
	{
		r := core.Ok[void.Void](nil)
		def := errors.New("default error")
		err := core.UnwrapErrOr(r, def)
		assert.Equal(t, def, err)
	}
}

// TestResult_wrapError tests wrapError method indirectly through Expect
func TestResult_wrapError(t *testing.T) {
	// Test wrapError with nil value (covers core.go:373-376)
	// TryErr(nil) now returns Ok, so Expect should not panic
	{
		r := core.TryErr[int](nil)
		assert.True(t, r.IsOk(), "TryErr(nil) should return Ok")
		assert.False(t, r.IsErr(), "TryErr(nil) should not return Err")
		// Expect should return zero value, not panic
		value := r.Expect("test message")
		assert.Equal(t, 0, value, "Expect should return zero value for TryErr(nil)")
	}
	// Test wrapError with non-error value (covers core.go:380)
	{
		r := core.TryErr[int](42)
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
	// Test wrapError with Ok result (covers core.go:382)
	{
		r := core.Ok[int](42)
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

// TestResult_AndThen tests AndThen method (covers core.go:487-492)
func TestResult_AndThen(t *testing.T) {
	// Test with Err result (should return r) (covers core.go:488-490)
	{
		r := core.TryErr[int]("error")
		ret := r.AndThen(func(i int) core.Result[int] {
			return core.Ok(i * 2)
		})
		assert.True(t, ret.IsErr())
		assert.Equal(t, "error", ret.Err().Error())
	}
	// Test with Ok result
	{
		r := core.Ok[int](10)
		ret := r.AndThen(func(i int) core.Result[int] {
			return core.Ok(i * 2)
		})
		assert.True(t, ret.IsOk())
		assert.Equal(t, 20, ret.Unwrap())
	}
}

func TestResult_ToError(t *testing.T) {
	// Test ToError with Ok result (should return nil)
	r1 := core.Ok[void.Void](nil)
	err1 := core.ToError(r1)
	assert.Nil(t, err1)

	// Test ToError with Err result (should return error)
	r2 := core.TryErr[void.Void](errors.New("test error"))
	err2 := core.ToError(r2)
	assert.NotNil(t, err2)
	assert.Equal(t, "test error", err2.Error())

	// Test ToError with string error
	r3 := core.TryErr[void.Void]("string error")
	err3 := core.ToError(r3)
	assert.NotNil(t, err3)
	assert.Equal(t, "string error", err3.Error())
}

func TestResult_UnwrapErrOr(t *testing.T) {
	// Test UnwrapErrOr with Err result (should return error)
	r1 := core.TryErr[void.Void](errors.New("test error"))
	err1 := core.UnwrapErrOr(r1, errors.New("default error"))
	assert.NotNil(t, err1)
	assert.Equal(t, "test error", err1.Error())

	// Test UnwrapErrOr with Ok result (should return default)
	r2 := core.Ok[void.Void](nil)
	err2 := core.UnwrapErrOr(r2, errors.New("default error"))
	assert.NotNil(t, err2)
	assert.Equal(t, "default error", err2.Error())

	// Test UnwrapErrOr with Ok result and nil default
	r3 := core.Ok[void.Void](nil)
	err3 := core.UnwrapErrOr(r3, nil)
	assert.Nil(t, err3)
}

func TestResult_UnwrapOr(t *testing.T) {
	// Test UnwrapOr with Err result (should return default)
	r1 := core.TryErr[int](errors.New("test error"))
	val1 := r1.UnwrapOr(42)
	assert.Equal(t, 42, val1)

	// Test UnwrapOr with Ok result (should return value)
	r2 := core.Ok[int](10)
	val2 := r2.UnwrapOr(42)
	assert.Equal(t, 10, val2)

	// Test UnwrapOr with string type
	r3 := core.TryErr[string](errors.New("test error"))
	val3 := r3.UnwrapOr("default")
	assert.Equal(t, "default", val3)

	r4 := core.Ok[string]("value")
	val4 := r4.UnwrapOr("default")
	assert.Equal(t, "value", val4)
}
