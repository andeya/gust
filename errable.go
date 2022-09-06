package gust

import (
	"fmt"
	"reflect"
)

type Errable[T any] struct {
	errVal *T
}

func NonErrable[T any]() Errable[T] {
	return Errable[T]{}
}

func ToErrable[T any](errVal T) Errable[T] {
	if any(errVal) == nil {
		return Errable[T]{}
	} else {
		v := reflect.ValueOf(errVal)
		if v.Kind() == reflect.Ptr && v.IsNil() {
			return Errable[T]{}
		}
	}
	return Errable[T]{errVal: &errVal}
}

func (e Errable[T]) IsErr() bool {
	return e.errVal != nil
}

func (e Errable[T]) ToError() error {
	if !e.IsErr() {
		return nil
	}
	return newAnyError(e.UnwrapErr())
}

func (e Errable[T]) UnwrapErr() T {
	return *e.errVal
}

func (e Errable[T]) UnwrapErrOr(def T) T {
	if e.IsErr() {
		return e.UnwrapErr()
	}
	return def
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
