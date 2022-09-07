package gust

import (
	"fmt"
	"reflect"
)

// Errable is the type that indicates whether there is an error.
type Errable[E any] struct {
	errVal *E
}

// NonErrable returns no error object.
func NonErrable[E any]() Errable[E] {
	return Errable[E]{}
}

// ToErrable converts an error value (E) to `Errable[T]`.
func ToErrable[E any](errVal E) Errable[E] {
	if any(errVal) == nil {
		return Errable[E]{}
	} else {
		v := reflect.ValueOf(errVal)
		if v.Kind() == reflect.Ptr && v.IsNil() {
			return Errable[E]{}
		}
	}
	return Errable[E]{errVal: &errVal}
}

func (e Errable[E]) IsErr() bool {
	return e.errVal != nil
}

func (e Errable[E]) IsOk() bool {
	return e.errVal == nil
}

func (e Errable[E]) ToError() error {
	if !e.IsErr() {
		return nil
	}
	return newAnyError(e.UnwrapErr())
}

func (e Errable[E]) UnwrapErr() E {
	return *e.errVal
}

func (e Errable[E]) UnwrapErrOr(def E) E {
	if e.IsErr() {
		return e.UnwrapErr()
	}
	return def
}

func (e Errable[E]) EnumResult() EnumResult[Void, E] {
	if e.IsErr() {
		return EnumErr[Void, E](e.UnwrapErr())
	}
	return EnumOk[Void, E](nil)
}

func (e Errable[E]) Result() Result[Void] {
	if e.IsErr() {
		return Err[Void](e.UnwrapErr())
	}
	return Ok[Void](nil)
}

func (e Errable[E]) Option() Option[E] {
	if e.IsErr() {
		return Some[E](e.UnwrapErr())
	}
	return None[E]()
}

// CtrlFlow returns the `CtrlFlow[E, Void]`.
func (e Errable[E]) CtrlFlow() CtrlFlow[E, Void] {
	if e.IsErr() {
		return Break[E, Void](e.UnwrapErr())
	}
	return Continue[E, Void](nil)
}

type errorWithVal struct {
	val any
}

func newAnyError(val any) error {
	if err, _ := val.(error); err != nil {
		return err
	}
	return &errorWithVal{val: val}
}

func (a *errorWithVal) Error() string {
	return fmt.Sprintf("%v", a.val)
}
