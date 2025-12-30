package conv

import (
	"reflect"
	"runtime"
	"unicode"
	"unicode/utf8"
)

// DerefType dereference, get the underlying non-pointer type.
//
// # Examples
//
//	```go
//	typ := reflect.TypeOf((*int)(nil))
//	baseType := conv.DerefType(typ) // returns int type
//	```
//
//go:inline
func DerefType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// DerefValue dereference and unpack interface,
// get the underlying non-pointer and non-interface value.
//
// # Examples
//
//	```go
//	x := 42
//	val := reflect.ValueOf(&x)
//	result := conv.DerefValue(val) // returns value of 42
//	```
//
//go:inline
func DerefValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}

// DerefPtrValue returns the underlying non-pointer type value.
//
// # Examples
//
//	```go
//	x := 42
//	val := reflect.ValueOf(&x)
//	result := conv.DerefPtrValue(val) // returns value of 42
//	```
//
//go:inline
func DerefPtrValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

// DerefInterfaceValue returns the value of the underlying type that implements the interface v.
//
// # Examples
//
//	```go
//	var iface interface{} = 42
//	val := reflect.ValueOf(iface)
//	result := conv.DerefInterfaceValue(val) // returns value of 42
//	```
//
//go:inline
func DerefInterfaceValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}

// DerefImplType returns the underlying type of the value that implements the interface v.
//
// # Examples
//
//	```go
//	var iface interface{} = 42
//	val := reflect.ValueOf(iface)
//	typ := conv.DerefImplType(val) // returns int type
//	```
//
//go:inline
func DerefImplType(v reflect.Value) reflect.Type {
	return DerefType(DerefInterfaceValue(v).Type())
}

// RefType convert T to *T, the ptrDepth is the count of '*'.
//
// # Examples
//
//	```go
//	typ := reflect.TypeOf(42)
//	ptrType := conv.RefType(typ, 1) // returns *int type
//	```
//
//go:inline
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
//
// # Examples
//
//	```go
//	val := reflect.ValueOf(42)
//	ptrVal := conv.RefValue(val, 1) // returns *int pointing to 42
//	```
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
//
// # Examples
//
//	```go
//	x, y := 1, 2
//	ptrSlice := []*int{&x, &y}
//	val := reflect.ValueOf(ptrSlice)
//	result := conv.DerefSliceValue(val) // returns []int{1, 2}
//	```
func DerefSliceValue(v reflect.Value) reflect.Value {
	length := v.Len()
	if length == 0 {
		return reflect.New(reflect.SliceOf(DerefType(v.Type().Elem()))).Elem()
	}
	s := make([]reflect.Value, length)
	for i := 0; i < length; i++ {
		s[i] = DerefValue(v.Index(i))
	}
	v = reflect.New(reflect.SliceOf(s[0].Type())).Elem()
	v = reflect.Append(v, s...)
	return v
}

// RefSliceValue convert []T to []*T, the ptrDepth is the count of '*'.
//
// # Examples
//
//	```go
//	slice := []int{1, 2}
//	val := reflect.ValueOf(slice)
//	result := conv.RefSliceValue(val, 1) // returns []*int{&1, &2}
//	```
func RefSliceValue(v reflect.Value, ptrDepth int) reflect.Value {
	if ptrDepth <= 0 {
		return v
	}
	length := v.Len()
	if length == 0 {
		return reflect.New(reflect.SliceOf(RefType(v.Type().Elem(), ptrDepth))).Elem()
	}
	s := make([]reflect.Value, length)
	for i := 0; i < length; i++ {
		s[i] = RefValue(v.Index(i), ptrDepth)
	}
	v = reflect.New(reflect.SliceOf(s[0].Type())).Elem()
	v = reflect.Append(v, s...)
	return v
}

// IsExportedOrBuiltinType checks if the type is exported or a builtin type.
// It dereferences pointer types before checking.
//
// # Examples
//
//	```go
//	typ := reflect.TypeOf(MyStruct{})
//	if conv.IsExportedOrBuiltinType(typ) {
//		// Type is exported or builtin
//	}
//	```
//
//go:inline
func IsExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return IsExportedName(t.Name()) || t.PkgPath() == ""
}

// IsExportedName checks if the name is exported (starts with uppercase letter).
//
// # Examples
//
//	```go
//	if conv.IsExportedName("MyFunction") {
//		// Name is exported
//	}
//	```
//
//go:inline
func IsExportedName(name string) bool {
	if name == "" {
		return false
	}
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

// TypeName gets the type name of the object.
// For functions, it returns the function name from runtime.
// For other types, it returns the type string representation.
//
// # Examples
//
//	```go
//	name := conv.TypeName(42)           // "int"
//	name := conv.TypeName(&MyStruct{})  // "*conv.MyStruct"
//	```
func TypeName(obj any) string {
	v := reflect.ValueOf(obj)
	t := v.Type()
	if t.Kind() == reflect.Func {
		fn := runtime.FuncForPC(v.Pointer())
		if fn != nil {
			return fn.Name()
		}
		// Fallback to type string if runtime.FuncForPC fails
		return t.String()
	}
	return t.String()
}

// IsCompositionMethod determines whether the method is inherited from an anonymous field (composition).
// It checks if the method is auto-generated by Go's composition mechanism.
//
// NOTE: This function relies on the "<autogenerated>" file marker used by Go's compiler
// for composition methods. This may not be reliable in all cases.
//
// # Examples
//
//	```go
//	type Base struct{}
//	func (Base) Method() {}
//
//	type Derived struct {
//		Base // anonymous field
//	}
//
//	method, _ := reflect.TypeOf(Derived{}).MethodByName("Method")
//	if conv.IsCompositionMethod(method) {
//		// Method is inherited from Base
//	}
//	```
func IsCompositionMethod(method reflect.Method) bool {
	// Check if method.Func is valid (not zero value)
	if !method.Func.IsValid() {
		return false
	}
	fn := runtime.FuncForPC(method.Func.Pointer())
	if fn == nil {
		return false
	}
	file, _ := fn.FileLine(fn.Entry())
	if file != "<autogenerated>" {
		return false
	}
	recv := method.Type.In(0)
	var found bool
	var composedMethod reflect.Method
	if recv.Kind() == reflect.Ptr {
		composedMethod, found = recv.Elem().MethodByName(method.Name)
	} else {
		composedMethod, found = reflect.PtrTo(recv).MethodByName(method.Name)
	}
	if !found {
		return true
	}
	// Check if composedMethod.Func is valid
	if !composedMethod.Func.IsValid() {
		return false
	}
	fn = runtime.FuncForPC(composedMethod.Func.Pointer())
	if fn == nil {
		return false
	}
	file, _ = fn.FileLine(fn.Entry())
	return file == "<autogenerated>"
}
