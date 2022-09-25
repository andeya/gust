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
	switch t := any(errVal).(type) {
	case error:
		if t == nil {
			return Errable[E]{}
		}
	case nil:
		return Errable[E]{}
	default:
		v := reflect.ValueOf(errVal)
		if v.Kind() == reflect.Ptr && v.IsNil() {
			return Errable[E]{}
		}
	}
	return Errable[E]{errVal: &errVal}
}

// TryPanic panics if the errVal is not nil.
func TryPanic[E any](errVal E) {
	ToErrable(errVal).TryPanic()
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

// TryPanic panics if the errVal is not nil.
func (e Errable[E]) TryPanic() {
	if e.IsErr() {
		panic(e.UnwrapErr())
	}
}

type errorWithVal struct {
	val any
}

func newAnyError(val any) error {
	if err, ok := val.(error); ok {
		return err
	}
	return &errorWithVal{val: val}
}

func (a *errorWithVal) Error() string {
	return fmt.Sprintf("%v", a.val)
}
