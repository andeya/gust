package gust_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestCtrlFlow(t *testing.T) {
	// Test Continue
	cf1 := gust.Continue[int, string]("hello")
	assert.True(t, cf1.IsContinue())
	assert.False(t, cf1.IsBreak())
	assert.Equal(t, "hello", cf1.UnwrapContinue())
	assert.Equal(t, gust.Some("hello"), cf1.ContinueValue())
	assert.True(t, cf1.BreakValue().IsNone())

	// Test Break
	cf2 := gust.Break[int, string](42)
	assert.True(t, cf2.IsBreak())
	assert.False(t, cf2.IsContinue())
	assert.Equal(t, 42, cf2.UnwrapBreak())
	assert.Equal(t, gust.Some(42), cf2.BreakValue())
	assert.True(t, cf2.ContinueValue().IsNone())

	// Test MapBreak
	cf3 := gust.Break[int, string](10)
	cf4 := cf3.MapBreak(func(x int) int { return x * 2 })
	assert.True(t, cf4.IsBreak())
	assert.Equal(t, 20, cf4.UnwrapBreak())

	// Test MapContinue
	cf5 := gust.Continue[int, string]("test")
	cf6 := cf5.MapContinue(func(s string) string { return s + "!" })
	assert.True(t, cf6.IsContinue())
	assert.Equal(t, "test!", cf6.UnwrapContinue())

	// Test Map
	cf7 := gust.Break[int, string](5)
	cf8 := cf7.Map(
		func(x int) int { return x * 2 },
		func(s string) string { return s + "!" },
	)
	assert.True(t, cf8.IsBreak())
	assert.Equal(t, 10, cf8.UnwrapBreak())

	// Test XMapBreak
	cf9 := gust.Break[int, string](100)
	cf10 := cf9.XMapBreak(func(x int) any { return x * 3 })
	assert.True(t, cf10.IsBreak())
	assert.Equal(t, 300, cf10.UnwrapBreak())

	// Test XMapContinue
	cf11 := gust.Continue[int, string]("world")
	cf12 := cf11.XMapContinue(func(s string) any { return len(s) })
	assert.True(t, cf12.IsContinue())
	assert.Equal(t, 5, cf12.UnwrapContinue())

	// Test XMap
	cf13 := gust.Break[int, string](7)
	cf14 := cf13.XMap(
		func(x int) any { return x * 2 },
		func(s string) any { return len(s) },
	)
	assert.True(t, cf14.IsBreak())
	assert.Equal(t, 14, cf14.UnwrapBreak())

	// Test Option
	cf15 := gust.Continue[int, string]("option")
	opt := cf15.Option()
	assert.True(t, opt.IsSome())
	assert.Equal(t, "option", opt.Unwrap())

	cf16 := gust.Break[int, string](99)
	opt2 := cf16.Option()
	assert.True(t, opt2.IsNone())

	// Test EnumResult
	cf17 := gust.Continue[int, string]("ok")
	er := cf17.EnumResult()
	assert.True(t, er.IsOk())
	assert.Equal(t, "ok", er.Unwrap())

	cf18 := gust.Break[int, string](123)
	er2 := cf18.EnumResult()
	assert.True(t, er2.IsErr())
	assert.Equal(t, 123, er2.UnwrapErr())

	// Test Result
	cf19 := gust.Continue[int, string]("success")
	res := cf19.Result()
	assert.True(t, res.IsOk())
	assert.Equal(t, "success", res.Unwrap())

	cf20 := gust.Break[int, string](456)
	res2 := cf20.Result()
	assert.True(t, res2.IsErr())

	// Test Errable
	cf21 := gust.Break[int, string](789)
	errable := cf21.Errable()
	assert.True(t, errable.IsErr())
	assert.Equal(t, 789, errable.UnwrapErr())

	cf22 := gust.Continue[int, string]("no error")
	errable2 := cf22.Errable()
	assert.False(t, errable2.IsErr())

	// Test ToX
	cf23 := gust.Break[int, string](111)
	anyCF := cf23.ToX()
	assert.True(t, anyCF.IsBreak())
	assert.Equal(t, 111, anyCF.UnwrapBreak())

	// Test ToXBreak
	cf24 := gust.Break[int, string](222)
	xBreak := cf24.ToXBreak()
	assert.True(t, xBreak.IsBreak())
	assert.Equal(t, 222, xBreak.UnwrapBreak())

	// Test ToXContinue
	cf25 := gust.Continue[int, string]("continue")
	xContinue := cf25.ToXContinue()
	assert.True(t, xContinue.IsContinue())
	assert.Equal(t, "continue", xContinue.UnwrapContinue())

	// Test String
	assert.Contains(t, cf1.String(), "Continue")
	assert.Contains(t, cf2.String(), "Break")
}

func TestSigCtrlFlow(t *testing.T) {
	// Test SigContinue
	scf1 := gust.SigContinue[int](42)
	assert.True(t, scf1.IsContinue())
	assert.Equal(t, 42, scf1.UnwrapContinue())

	// Test SigBreak
	scf2 := gust.SigBreak[string]("error")
	assert.True(t, scf2.IsBreak())
	assert.Equal(t, "error", scf2.UnwrapBreak())
}

func TestAnyCtrlFlow(t *testing.T) {
	// Test AnyContinue
	acf1 := gust.AnyContinue("test")
	assert.True(t, acf1.IsContinue())
	assert.Equal(t, "test", acf1.UnwrapContinue())

	// Test AnyBreak
	acf2 := gust.AnyBreak(123)
	assert.True(t, acf2.IsBreak())
	assert.Equal(t, 123, acf2.UnwrapBreak())
}
