package core

import (
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
