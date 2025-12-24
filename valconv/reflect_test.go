package valconv

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDerefType(t *testing.T) {
	// Test with non-pointer type
	typ := reflect.TypeOf(42)
	assert.Equal(t, typ, DerefType(typ))

	// Test with pointer type
	ptrTyp := reflect.TypeOf((*int)(nil))
	assert.Equal(t, typ, DerefType(ptrTyp))

	// Test with double pointer
	ptrPtrTyp := reflect.TypeOf((**int)(nil))
	assert.Equal(t, typ, DerefType(ptrPtrTyp))
}

func TestDerefValue(t *testing.T) {
	// Test with non-pointer value
	val := reflect.ValueOf(42)
	result := DerefValue(val)
	assert.Equal(t, 42, result.Interface())

	// Test with pointer value
	x := 42
	ptrVal := reflect.ValueOf(&x)
	result2 := DerefValue(ptrVal)
	assert.Equal(t, 42, result2.Interface())

	// Test with interface value
	var iface interface{} = 42
	ifaceVal := reflect.ValueOf(iface)
	result3 := DerefValue(ifaceVal)
	assert.Equal(t, 42, result3.Interface())

	// Test with pointer to interface
	var ifacePtr interface{} = &x
	ifacePtrVal := reflect.ValueOf(ifacePtr)
	result4 := DerefValue(ifacePtrVal)
	assert.Equal(t, 42, result4.Interface())
}

func TestDerefPtrValue(t *testing.T) {
	// Test with non-pointer value
	val := reflect.ValueOf(42)
	result := DerefPtrValue(val)
	assert.Equal(t, 42, result.Interface())

	// Test with pointer value
	x := 42
	ptrVal := reflect.ValueOf(&x)
	result2 := DerefPtrValue(ptrVal)
	assert.Equal(t, 42, result2.Interface())

	// Test with double pointer
	xx := 42
	ptr := &xx
	ptrPtr := &ptr
	ptrPtrVal := reflect.ValueOf(ptrPtr)
	result3 := DerefPtrValue(ptrPtrVal)
	assert.Equal(t, 42, result3.Interface())
}

func TestDerefInterfaceValue(t *testing.T) {
	// Test with non-interface value
	val := reflect.ValueOf(42)
	result := DerefInterfaceValue(val)
	assert.Equal(t, 42, result.Interface())

	// Test with interface value
	var iface interface{} = 42
	ifaceVal := reflect.ValueOf(iface)
	result2 := DerefInterfaceValue(ifaceVal)
	assert.Equal(t, 42, result2.Interface())
}

func TestDerefImplType(t *testing.T) {
	// Test with interface value
	var iface interface{} = 42
	ifaceVal := reflect.ValueOf(iface)
	result := DerefImplType(ifaceVal)
	assert.Equal(t, reflect.TypeOf(42), result)
}

func TestRefType(t *testing.T) {
	typ := reflect.TypeOf(42)

	// Test with positive ptrDepth
	result1 := RefType(typ, 1)
	assert.Equal(t, reflect.PtrTo(typ), result1)

	result2 := RefType(typ, 2)
	assert.Equal(t, reflect.PtrTo(reflect.PtrTo(typ)), result2)

	// Test with zero ptrDepth
	result3 := RefType(typ, 0)
	assert.Equal(t, typ, result3)

	// Test with negative ptrDepth
	ptrTyp := reflect.PtrTo(typ)
	result4 := RefType(ptrTyp, -1)
	assert.Equal(t, typ, result4)

	// Test with negative ptrDepth but not enough pointers
	result5 := RefType(typ, -1)
	assert.Equal(t, typ, result5)
}

func TestRefValue(t *testing.T) {
	val := reflect.ValueOf(42)

	// Test with positive ptrDepth
	result1 := RefValue(val, 1)
	assert.Equal(t, reflect.Ptr, result1.Kind())
	assert.Equal(t, 42, result1.Elem().Interface())

	result2 := RefValue(val, 2)
	assert.Equal(t, reflect.Ptr, result2.Kind())
	assert.Equal(t, 42, result2.Elem().Elem().Interface())

	// Test with zero ptrDepth
	result3 := RefValue(val, 0)
	assert.Equal(t, 42, result3.Interface())

	// Test with negative ptrDepth
	x := 42
	ptrVal := reflect.ValueOf(&x)
	result4 := RefValue(ptrVal, -1)
	assert.Equal(t, 42, result4.Interface())

	// Test with negative ptrDepth but not enough pointers
	result5 := RefValue(val, -1)
	assert.Equal(t, 42, result5.Interface())
}

func TestDerefSliceValue(t *testing.T) {
	// Test with slice of pointers
	x, y := 1, 2
	ptrSlice := []*int{&x, &y}
	val := reflect.ValueOf(ptrSlice)
	result := DerefSliceValue(val)
	resultSlice := result.Interface().([]int)
	assert.Len(t, resultSlice, 2)
	assert.Equal(t, 1, resultSlice[0])
	assert.Equal(t, 2, resultSlice[1])

	// Test with empty slice
	emptySlice := []*int{}
	emptyVal := reflect.ValueOf(emptySlice)
	result2 := DerefSliceValue(emptyVal)
	assert.Equal(t, reflect.Slice, result2.Kind())
	assert.Equal(t, 0, result2.Len())
}

func TestRefSliceValue(t *testing.T) {
	// Test with positive ptrDepth
	v := reflect.ValueOf([]int{1, 2})
	v = RefSliceValue(v, 1)
	ret := v.Interface().([]*int)
	assert.Len(t, ret, 2)
	assert.Equal(t, 1, *ret[0])
	assert.Equal(t, 2, *ret[1])

	// Test with empty slice
	v = reflect.ValueOf([]int{})
	v = RefSliceValue(v, 1)
	ret = v.Interface().([]*int)
	assert.Len(t, ret, 0)

	// Test with zero ptrDepth
	v2 := reflect.ValueOf([]int{1, 2})
	v2 = RefSliceValue(v2, 0)
	ret2 := v2.Interface().([]int)
	assert.Len(t, ret2, 2)
	assert.Equal(t, 1, ret2[0])
	assert.Equal(t, 2, ret2[1])

	// Test with negative ptrDepth
	v3 := reflect.ValueOf([]int{1, 2})
	v3 = RefSliceValue(v3, -1)
	ret3 := v3.Interface().([]int)
	assert.Len(t, ret3, 2)
	assert.Equal(t, 1, ret3[0])
	assert.Equal(t, 2, ret3[1])
}

func TestRefType_NegativePtrDepth_NotEnoughPointers(t *testing.T) {
	// Test RefType with negative ptrDepth but not enough pointers
	typ := reflect.TypeOf(42)
	result := RefType(typ, -2) // More negative than available pointers
	assert.Equal(t, typ, result)
}

func TestRefValue_NegativePtrDepth_NotEnoughPointers(t *testing.T) {
	// Test RefValue with negative ptrDepth but not enough pointers
	val := reflect.ValueOf(42)
	result := RefValue(val, -2) // More negative than available pointers
	assert.Equal(t, 42, result.Interface())
}

func TestDerefSliceValue_EmptySlice(t *testing.T) {
	// Test DerefSliceValue with empty slice (m < 0)
	emptySlice := []*int{}
	val := reflect.ValueOf(emptySlice)
	result := DerefSliceValue(val)
	assert.Equal(t, reflect.Slice, result.Kind())
	assert.Equal(t, 0, result.Len())
}

func TestRefSliceValue_EmptySlice(t *testing.T) {
	// Test RefSliceValue with empty slice (m < 0)
	emptySlice := []int{}
	val := reflect.ValueOf(emptySlice)
	result := RefSliceValue(val, 1)
	assert.Equal(t, reflect.Slice, result.Kind())
	assert.Equal(t, 0, result.Len())
}

func TestDerefValue_MultiplePointers(t *testing.T) {
	// Test DerefValue with multiple levels of pointers
	x := 42
	ptr := &x
	ptrPtr := &ptr
	val := reflect.ValueOf(ptrPtr)
	result := DerefValue(val)
	assert.Equal(t, 42, result.Interface())
}

func TestDerefValue_InterfaceWithPointer(t *testing.T) {
	// Test DerefValue with interface containing pointer
	x := 42
	var iface interface{} = &x
	val := reflect.ValueOf(iface)
	result := DerefValue(val)
	assert.Equal(t, 42, result.Interface())
}

func TestDerefPtrValue_MultiplePointers(t *testing.T) {
	// Test DerefPtrValue with multiple levels of pointers
	x := 42
	ptr := &x
	ptrPtr := &ptr
	val := reflect.ValueOf(ptrPtr)
	result := DerefPtrValue(val)
	assert.Equal(t, 42, result.Interface())
}

func TestRefType_MultipleLevels(t *testing.T) {
	// Test RefType with multiple positive ptrDepth
	typ := reflect.TypeOf(42)
	result := RefType(typ, 3)
	expected := reflect.PtrTo(reflect.PtrTo(reflect.PtrTo(typ)))
	assert.Equal(t, expected, result)
}

func TestRefValue_MultipleLevels(t *testing.T) {
	// Test RefValue with multiple positive ptrDepth
	val := reflect.ValueOf(42)
	result := RefValue(val, 3)
	assert.Equal(t, reflect.Ptr, result.Kind())
	assert.Equal(t, 42, result.Elem().Elem().Elem().Interface())
}

func TestRefSliceValue_MultipleLevels(t *testing.T) {
	// Test RefSliceValue with multiple positive ptrDepth
	slice := []int{1, 2}
	val := reflect.ValueOf(slice)
	result := RefSliceValue(val, 2)
	ret := result.Interface().([]**int)
	assert.Len(t, ret, 2)
	assert.Equal(t, 1, **ret[0])
	assert.Equal(t, 2, **ret[1])
}

