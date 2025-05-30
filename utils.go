package gust

import (
	"fmt"
	"reflect"
)

func defaultValue[T any]() T {
	var zero T
	_ = initNilPtr(reflect.ValueOf(&zero))
	return zero
}

func defaultValuePtr[T any]() *T {
	var zeroPtr = new(T)
	_ = initNilPtr(reflect.ValueOf(zeroPtr))
	return zeroPtr
}

// initNilPtr initializes nil pointer with zero value.
func initNilPtr(v reflect.Value) (done bool) {
	for {
		kind := v.Kind()
		if kind == reflect.Interface {
			v = v.Elem()
			continue
		}
		if kind != reflect.Ptr {
			return true
		}
		u := v.Elem()
		if u.IsValid() {
			v = u
			continue
		}
		if !v.CanSet() {
			return false
		}
		v2 := reflect.New(v.Type().Elem())
		v.Set(v2)
		v = v.Elem()
	}
}

type panicValue[T any] struct {
	value *T
}

func (p panicValue[T]) ValueOrDefault() T {
	if p.value == nil {
		var t T
		return t
	}
	return *p.value
}

func (p panicValue[T]) String() string {
	if p.value == nil {
		return fmt.Sprintf("%v", p.value)
	}
	return fmt.Sprintf("%v", *p.value)
}

func (p panicValue[T]) GoString() string {
	if p.value == nil {
		return fmt.Sprintf("%#v", p.value)
	}
	return fmt.Sprintf("%#v", *p.value)
}

func (p panicValue[T]) Error() string {
	return fmt.Sprintf("%v", p.value)
}
