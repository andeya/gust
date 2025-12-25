package gust

import "fmt"

// AnyCtrlFlow is a placeholder for wildcard control flow statements.
type AnyCtrlFlow = SigCtrlFlow[any]

// AnyContinue returns a AnyCtrlFlow that tells the operation to continue.
func AnyContinue(c any) AnyCtrlFlow {
	return SigContinue[any](c)
}

// AnyBreak returns a AnyCtrlFlow that tells the operation to break.
func AnyBreak(b any) AnyCtrlFlow {
	return SigBreak[any](b)
}

// SigCtrlFlow is a placeholder for single type control flow statements.
type SigCtrlFlow[T any] struct {
	CtrlFlow[T, T]
}

// SigContinue returns a `SigCtrlFlow[T]` that tells the operation to continue.
func SigContinue[T any](c T) SigCtrlFlow[T] {
	return SigCtrlFlow[T]{
		CtrlFlow[T, T]{_continue: Some(c)},
	}
}

// SigBreak returns a `SigCtrlFlow[T]` that tells the operation to break.
func SigBreak[T any](b T) SigCtrlFlow[T] {
	return SigCtrlFlow[T]{
		CtrlFlow[T, T]{_break: Some(b)},
	}
}

// CtrlFlow is used to tell an operation whether it should exit early or go on as usual.
//
// This is used when exposing things (like graph traversals or visitors) where
// you want the user to be able to choose whether to exit early.
type CtrlFlow[B any, C any] struct {
	_break    Option[B]
	_continue Option[C]
}

// Continue returns a CtrlFlow that tells the operation to continue.
func Continue[B any, C any](c C) CtrlFlow[B, C] {
	return CtrlFlow[B, C]{_continue: Some(c)}
}

// Break returns a CtrlFlow that tells the operation to break.
func Break[B any, C any](b B) CtrlFlow[B, C] {
	return CtrlFlow[B, C]{_break: Some(b)}
}

// String returns the string representation.
func (c CtrlFlow[B, C]) String() string {
	if c.IsBreak() {
		return fmt.Sprintf("Break(%v)", c.UnwrapBreak())
	}
	return fmt.Sprintf("Continue(%v)", c.UnwrapContinue())
}

// IsBreak returns `true` if this is a `Break` variant.
func (c CtrlFlow[B, C]) IsBreak() bool {
	return c._break.IsSome()
}

// IsContinue returns `true` if this is a `Continue` variant.
func (c CtrlFlow[B, C]) IsContinue() bool {
	return c._continue.IsSome()
}

// BreakValue returns the inner value of a `Break` variant.
func (c CtrlFlow[B, C]) BreakValue() Option[B] {
	return c._break
}

// ContinueValue returns the inner value of a `Continue` variant.
func (c CtrlFlow[B, C]) ContinueValue() Option[C] {
	return c._continue
}

// MapBreak maps `Break` value by applying a function
// to the break value in case it exists.
func (c CtrlFlow[B, C]) MapBreak(f func(B) B) CtrlFlow[B, C] {
	if c.IsBreak() {
		return CtrlFlow[B, C]{
			_break:    c._break.Map(f),
			_continue: c._continue,
		}
	}
	return c
}

// XMapBreak maps `Break` value by applying a function
// to the break value in case it exists.
func (c CtrlFlow[B, C]) XMapBreak(f func(B) any) CtrlFlow[any, C] {
	if c.IsBreak() {
		return CtrlFlow[any, C]{
			_break:    c._break.XMap(f),
			_continue: c._continue,
		}
	}
	return CtrlFlow[any, C]{
		_break:    c._break.ToX(),
		_continue: c._continue,
	}
}

// MapContinue maps `Continue` value by applying a function
// to the continue value in case it exists.
func (c CtrlFlow[B, C]) MapContinue(f func(C) C) CtrlFlow[B, C] {
	if c.IsContinue() {
		return CtrlFlow[B, C]{
			_break:    c._break,
			_continue: c._continue.Map(f),
		}
	}
	return c
}

// XMapContinue maps `Continue` value by applying a function
// to the continue value in case it exists.
func (c CtrlFlow[B, C]) XMapContinue(f func(C) any) CtrlFlow[B, any] {
	if c.IsContinue() {
		return CtrlFlow[B, any]{
			_break:    c._break,
			_continue: c._continue.XMap(f),
		}
	}
	return CtrlFlow[B, any]{
		_break:    c._break,
		_continue: c._continue.ToX(),
	}
}

// Map maps both `Break` and `Continue` values by applying a function
// to the respective value in case it exists.
func (c CtrlFlow[B, C]) Map(f func(B) B, g func(C) C) CtrlFlow[B, C] {
	return c.MapBreak(f).MapContinue(g)
}

// XMap maps both `Break` and `Continue` values by applying functions
// to the break and continue values in case they exist.
func (c CtrlFlow[B, C]) XMap(f func(B) any, g func(C) any) CtrlFlow[any, any] {
	return c.XMapBreak(f).XMapContinue(g)
}

// UnwrapBreak returns the inner value of a `Break` variant.
//
// Panics if the variant is not a `Break`.
func (c CtrlFlow[B, C]) UnwrapBreak() B {
	return c._break.Unwrap()
}

// UnwrapContinue returns the inner value of a `Continue` variant.
//
// Panics if the variant is not a `Continue`.
func (c CtrlFlow[B, C]) UnwrapContinue() C {
	return c._continue.Unwrap()
}

// Option converts the `CtrlFlow` to an `Option`.
func (c CtrlFlow[B, C]) Option() Option[C] {
	if c.IsBreak() {
		return None[C]()
	}
	return c._continue
}

// EnumResult converts the `CtrlFlow[B,C]` to an `EnumResult[C,B]`.
func (c CtrlFlow[B, C]) EnumResult() EnumResult[C, B] {
	if c.IsBreak() {
		return EnumErr[C, B](c._break.UnwrapUnchecked())
	}
	return EnumOk[C, B](c._continue.UnwrapUnchecked())
}

// Result converts the `CtrlFlow[B,C]` to an `Result[C]`.
func (c CtrlFlow[B, C]) Result() Result[C] {
	if c.IsBreak() {
		return Err[C](c._break.UnwrapUnchecked())
	}
	return Ok[C](c._continue.UnwrapUnchecked())
}

// Errable converts the `CtrlFlow[B,C]` to an `Errable[B]`.
func (c CtrlFlow[B, C]) Errable() Errable[B] {
	if c.IsBreak() {
		return ToErrable[B](c._break.UnwrapUnchecked())
	}
	return NonErrable[B]()
}

// ToX converts the `CtrlFlow[B, C]` to `AnyCtrlFlow`.
func (c CtrlFlow[B, C]) ToX() AnyCtrlFlow {
	return SigCtrlFlow[any]{
		CtrlFlow: CtrlFlow[any, any]{
			_break:    c._break.ToX(),
			_continue: c._continue.ToX(),
		},
	}
}

// ToXBreak converts the break type of `CtrlFlow[B, C]` to `any`.
func (c CtrlFlow[B, C]) ToXBreak() CtrlFlow[any, C] {
	return CtrlFlow[any, C]{
		_break:    c._break.ToX(),
		_continue: c._continue,
	}
}

// ToXContinue converts the continue type of `CtrlFlow[B, C]` to `any`.
func (c CtrlFlow[B, C]) ToXContinue() CtrlFlow[B, any] {
	return CtrlFlow[B, any]{
		_break:    c._break,
		_continue: c._continue.ToX(),
	}
}
