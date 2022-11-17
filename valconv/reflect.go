package valconv

import (
	"reflect"
)

// DerefType dereference, get the underlying non-pointer type.
func DerefType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// DerefValue dereference and unpack interface,
// get the underlying non-pointer and non-interface value.
func DerefValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}

// DerefPtrValue returns the underlying non-pointer type value.
func DerefPtrValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

// DerefInterfaceValue returns the value of the underlying type that implements the interface v.
func DerefInterfaceValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}

// DerefImplType returns the underlying type of the value that implements the interface v.
func DerefImplType(v reflect.Value) reflect.Type {
	return DerefType(DerefInterfaceValue(v).Type())
}

// RefType convert T to *T, the ptrDepth is the count of '*'.
func RefType(t reflect.Type, ptrDepth int) reflect.Type {
	switch {
	case ptrDepth > 0:
		for ; ptrDepth > 0; ptrDepth-- {
			t = reflect.PtrTo(t)
		}
	case ptrDepth < 0:
		for ; ptrDepth < 0 && t.Kind() == reflect.Ptr; ptrDepth++ {
			t = t.Elem()
		}
	}
	return t
}

// RefValue convert T to *T, the ptrDepth is the count of '*'.
func RefValue(v reflect.Value, ptrDepth int) reflect.Value {
	switch {
	case ptrDepth > 0:
		for ; ptrDepth > 0; ptrDepth-- {
			vv := reflect.New(v.Type())
			vv.Elem().Set(v)
			v = vv
		}
	case ptrDepth < 0:
		for ; ptrDepth < 0 && v.Kind() == reflect.Ptr; ptrDepth++ {
			v = v.Elem()
		}
	}
	return v
}

// DerefSliceValue convert []*T to []T.
func DerefSliceValue(v reflect.Value) reflect.Value {
	m := v.Len() - 1
	if m < 0 {
		return reflect.New(reflect.SliceOf(DerefType(v.Type().Elem()))).Elem()
	}
	s := make([]reflect.Value, m+1)
	for ; m >= 0; m-- {
		s[m] = DerefValue(v.Index(m))
	}
	v = reflect.New(reflect.SliceOf(s[0].Type())).Elem()
	v = reflect.Append(v, s...)
	return v
}

// RefSliceValue convert []T to []*T, the ptrDepth is the count of '*'.
func RefSliceValue(v reflect.Value, ptrDepth int) reflect.Value {
	if ptrDepth <= 0 {
		return v
	}
	m := v.Len() - 1
	if m < 0 {
		return reflect.New(reflect.SliceOf(RefType(v.Type().Elem(), ptrDepth))).Elem()
	}
	s := make([]reflect.Value, m+1)
	for ; m >= 0; m-- {
		s[m] = RefValue(v.Index(m), ptrDepth)
	}
	v = reflect.New(reflect.SliceOf(s[0].Type())).Elem()
	v = reflect.Append(v, s...)
	return v
}
