package gust

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanicValue_ValueOrDefault(t *testing.T) {
	// Test with nil value
	var pv1 panicValue[int]
	assert.Equal(t, 0, pv1.ValueOrDefault())

	// Test with non-nil value
	val := 42
	pv2 := panicValue[int]{value: &val}
	assert.Equal(t, 42, pv2.ValueOrDefault())

	// Test with string
	var pv3 panicValue[string]
	assert.Equal(t, "", pv3.ValueOrDefault())

	str := "test"
	pv4 := panicValue[string]{value: &str}
	assert.Equal(t, "test", pv4.ValueOrDefault())
}

func TestPanicValue_String(t *testing.T) {
	// Test with nil value
	var pv1 panicValue[int]
	assert.Equal(t, "<nil>", pv1.String())

	// Test with non-nil value
	val := 42
	pv2 := panicValue[int]{value: &val}
	assert.Equal(t, "42", pv2.String())

	// Test with string
	var pv3 panicValue[string]
	assert.Equal(t, "<nil>", pv3.String())

	str := "test"
	pv4 := panicValue[string]{value: &str}
	assert.Equal(t, "test", pv4.String())
}

func TestPanicValue_GoString(t *testing.T) {
	// Test with nil value
	var pv1 panicValue[int]
	assert.Equal(t, "(*int)(nil)", pv1.GoString())

	// Test with non-nil value
	val := 42
	pv2 := panicValue[int]{value: &val}
	assert.Equal(t, "42", pv2.GoString())

	// Test with string
	var pv3 panicValue[string]
	assert.Equal(t, "(*string)(nil)", pv3.GoString())

	str := "test"
	pv4 := panicValue[string]{value: &str}
	assert.Equal(t, "\"test\"", pv4.GoString())
}

func TestPanicValue_Error(t *testing.T) {
	// Test with nil value
	var pv1 panicValue[int]
	assert.Equal(t, "<nil>", pv1.Error())

	// Test with non-nil value
	val := 42
	pv2 := panicValue[int]{value: &val}
	assert.Equal(t, "42", pv2.Error())

	// Test with string
	str := "test error"
	pv3 := panicValue[string]{value: &str}
	assert.Equal(t, "test error", pv3.Error())
}

func TestInitNilPtr(t *testing.T) {
	// Test with non-pointer (should return true immediately)
	var i int
	v := reflect.ValueOf(i)
	done := initNilPtr(v)
	assert.True(t, done)

	// Test with nil pointer
	var nilPtr *int
	v2 := reflect.ValueOf(&nilPtr)
	done2 := initNilPtr(v2)
	assert.True(t, done2)
	assert.NotNil(t, nilPtr)
	assert.Equal(t, 0, *nilPtr) // Verify it's initialized to zero

	// Test with nested nil pointer
	var nilPtr2 **int
	v3 := reflect.ValueOf(&nilPtr2)
	done3 := initNilPtr(v3)
	assert.True(t, done3)
	assert.NotNil(t, nilPtr2)
	assert.NotNil(t, *nilPtr2)
	assert.Equal(t, 0, **nilPtr2) // Verify nested pointer is initialized to zero

	// Test with interface containing nil pointer
	var iface interface{} = (*int)(nil)
	v4 := reflect.ValueOf(&iface)
	done4 := initNilPtr(v4)
	// initNilPtr may not be able to set interface values, so we just check it doesn't panic
	_ = done4

	// Test with non-nil pointer (should continue)
	val := 42
	ptr := &val
	v5 := reflect.ValueOf(&ptr)
	done5 := initNilPtr(v5)
	assert.True(t, done5)

	// Test with unsettable value
	var unsettable *int
	v6 := reflect.ValueOf(unsettable)
	done6 := initNilPtr(v6)
	assert.False(t, done6)
}

func TestDefaultValue(t *testing.T) {
	// Test with int
	var zeroInt int
	assert.Equal(t, zeroInt, defaultValue[int]())

	// Test with string
	var zeroStr string
	assert.Equal(t, zeroStr, defaultValue[string]())

	// Test with pointer type
	result := defaultValue[*int]()
	assert.NotNil(t, result)    // Should initialize nil pointer
	assert.Equal(t, 0, *result) // Verify it's initialized to zero

	// Test with nested pointer
	result2 := defaultValue[**int]()
	assert.NotNil(t, result2)
	assert.NotNil(t, *result2)
	assert.Equal(t, 0, **result2) // Verify nested pointer is initialized to zero

	// Test with struct
	type S struct {
		X int
		Y string
	}
	var zeroS S
	assert.Equal(t, zeroS, defaultValue[S]())
}

func TestDefaultValuePtr(t *testing.T) {
	// Test with int
	ptr := defaultValuePtr[int]()
	assert.NotNil(t, ptr)
	assert.Equal(t, 0, *ptr)

	// Test with string
	ptr2 := defaultValuePtr[string]()
	assert.NotNil(t, ptr2)
	assert.Equal(t, "", *ptr2)

	// Test with pointer type
	ptr3 := defaultValuePtr[*int]()
	assert.NotNil(t, ptr3)
	assert.NotNil(t, *ptr3) // Should initialize nested nil pointer

	// Test with nested pointer
	ptr4 := defaultValuePtr[**int]()
	assert.NotNil(t, ptr4)
	assert.NotNil(t, *ptr4)
	assert.NotNil(t, **ptr4)

	// Test with struct
	type S struct {
		X int
		Y string
	}
	ptr5 := defaultValuePtr[S]()
	assert.NotNil(t, ptr5)
	assert.Equal(t, S{}, *ptr5)
}
