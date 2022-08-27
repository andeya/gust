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

func (e Errable[T]) Ref() *Errable[T] {
	return &e
}

func (e *Errable[T]) HasError() bool {
	return e != nil && e.errVal != nil
}

func (e *Errable[T]) ToError() error {
	if !e.HasError() {
		return nil
	}
	if err, ok := any(e.Unwrap()).(error); ok {
		return err
	}
	return fmt.Errorf("%v", e.Unwrap())
}

func (e *Errable[T]) Unwrap() T {
	return *e.errVal
}

func (e *Errable[T]) UnwrapOr(def T) T {
	if e.HasError() {
		return e.Unwrap()
	}
	return def
}
