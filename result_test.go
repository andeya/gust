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
	var r gust.EnumResult[string, error]
	defer func() {
		assert.Equal(t, gust.EnumErr[string, error](gust.ToErrBox("err")), r)
	}()
	defer gust.CatchEnumResult[string, error](&r)
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

func TestResultUnwrapOrThrow_6(t *testing.T) {
	var r gust.EnumResult[string, error]
	defer func() {
		assert.Equal(t, gust.EnumErr[string, error](gust.ToErrBox("err")), r)
	}()
	defer gust.CatchEnumResult[string, error](&r)
	var r1 = gust.Ok(1)
	var v1 = r1.UnwrapOrThrow()
	assert.Equal(t, 1, v1)
	var r2 = gust.Err[int]("err")
	var v2 = r2.UnwrapOrThrow()
	assert.Equal(t, 0, v2)
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
	var sq = func(x int) gust.EnumResult[int, int] {
		return gust.EnumOk[int, int](x * x)
	}
	var err = func(x int) gust.EnumResult[int, int] {
		return gust.EnumErr[int, int](x)
	}

	assert.Equal(t, gust.EnumOk[int, int](2).OrElse(sq).OrElse(sq), gust.EnumOk[int, int](2))
	assert.Equal(t, gust.EnumOk[int, int](2).OrElse(err).OrElse(sq), gust.EnumOk[int, int](2))
	assert.Equal(t, gust.EnumErr[int, int](3).OrElse(sq).OrElse(err), gust.EnumOk[int, int](9))
	assert.Equal(t, gust.EnumErr[int, int](3).OrElse(err).OrElse(err), gust.EnumErr[int, int](3))
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
