package valconv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZero(t *testing.T) {
	// Test with int
	assert.Equal(t, 0, Zero[int]())

	// Test with string
	assert.Equal(t, "", Zero[string]())

	// Test with pointer
	var nilPtr *int
	assert.Equal(t, nilPtr, Zero[*int]())

	// Test with struct
	type S struct {
		X int
		Y string
	}
	assert.Equal(t, S{}, Zero[S]())
}

func TestRef(t *testing.T) {
	// Test with int
	x := 42
	ptr := Ref(x)
	assert.NotNil(t, ptr)
	assert.Equal(t, 42, *ptr)

	// Test with string
	s := "test"
	ptrS := Ref(s)
	assert.NotNil(t, ptrS)
	assert.Equal(t, "test", *ptrS)

	// Test with struct
	type S struct {
		X int
	}
	val := S{X: 10}
	ptrVal := Ref(val)
	assert.NotNil(t, ptrVal)
	assert.Equal(t, 10, ptrVal.X)
}

func TestDeref(t *testing.T) {
	// Test with non-nil pointer
	x := 42
	ptr := &x
	assert.Equal(t, 42, Deref(ptr))

	// Test with nil pointer
	var nilPtr *int
	assert.Equal(t, 0, Deref(nilPtr))

	// Test with string pointer
	s := "test"
	ptrS := &s
	assert.Equal(t, "test", Deref(ptrS))

	// Test with nil string pointer
	var nilStrPtr *string
	assert.Equal(t, "", Deref(nilStrPtr))
}

func TestRefSlice(t *testing.T) {
	// Test with non-nil slice
	slice := []int{1, 2, 3}
	refSlice := RefSlice(slice)
	assert.Len(t, refSlice, 3)
	assert.Equal(t, 1, *refSlice[0])
	assert.Equal(t, 2, *refSlice[1])
	assert.Equal(t, 3, *refSlice[2])

	// Test with nil slice
	var nilSlice []int
	refNilSlice := RefSlice(nilSlice)
	assert.Nil(t, refNilSlice)

	// Test with empty slice
	emptySlice := []int{}
	refEmptySlice := RefSlice(emptySlice)
	assert.Len(t, refEmptySlice, 0)
	assert.NotNil(t, refEmptySlice)

	// Test with string slice
	strSlice := []string{"a", "b"}
	refStrSlice := RefSlice(strSlice)
	assert.Len(t, refStrSlice, 2)
	assert.Equal(t, "a", *refStrSlice[0])
	assert.Equal(t, "b", *refStrSlice[1])
}

func TestDerefSlice(t *testing.T) {
	// Test with non-nil pointers
	x, y, z := 1, 2, 3
	ptrSlice := []*int{&x, &y, &z}
	derefSlice := DerefSlice(ptrSlice)
	assert.Len(t, derefSlice, 3)
	assert.Equal(t, 1, derefSlice[0])
	assert.Equal(t, 2, derefSlice[1])
	assert.Equal(t, 3, derefSlice[2])

	// Test with nil pointer in slice
	nilPtrSlice := []*int{&x, nil, &z}
	derefNilSlice := DerefSlice(nilPtrSlice)
	assert.Len(t, derefNilSlice, 3)
	assert.Equal(t, 1, derefNilSlice[0])
	assert.Equal(t, 0, derefNilSlice[1]) // zero value
	assert.Equal(t, 3, derefNilSlice[2])

	// Test with nil slice
	var nilSlice []*int
	derefNilSlice2 := DerefSlice(nilSlice)
	assert.Nil(t, derefNilSlice2)

	// Test with empty slice
	emptySlice := []*int{}
	derefEmptySlice := DerefSlice(emptySlice)
	assert.Len(t, derefEmptySlice, 0)
	assert.NotNil(t, derefEmptySlice)

	// Test with string pointers
	s1, s2 := "a", "b"
	strPtrSlice := []*string{&s1, &s2}
	derefStrSlice := DerefSlice(strPtrSlice)
	assert.Len(t, derefStrSlice, 2)
	assert.Equal(t, "a", derefStrSlice[0])
	assert.Equal(t, "b", derefStrSlice[1])
}

func TestRefMap(t *testing.T) {
	// Test with non-nil map
	m := map[string]int{"a": 1, "b": 2}
	refMap := RefMap(m)
	assert.Len(t, refMap, 2)
	assert.Equal(t, 1, *refMap["a"])
	assert.Equal(t, 2, *refMap["b"])

	// Test with nil map
	var nilMap map[string]int
	refNilMap := RefMap(nilMap)
	assert.Nil(t, refNilMap)

	// Test with empty map
	emptyMap := map[string]int{}
	refEmptyMap := RefMap(emptyMap)
	assert.Len(t, refEmptyMap, 0)
	assert.NotNil(t, refEmptyMap)
}

func TestDerefMap(t *testing.T) {
	// Test with non-nil pointers
	x, y := 1, 2
	ptrMap := map[string]*int{"a": &x, "b": &y}
	derefMap := DerefMap(ptrMap)
	assert.Len(t, derefMap, 2)
	assert.Equal(t, 1, derefMap["a"])
	assert.Equal(t, 2, derefMap["b"])

	// Test with nil pointer in map
	nilPtrMap := map[string]*int{"a": &x, "b": nil, "c": &y}
	derefNilMap := DerefMap(nilPtrMap)
	assert.Len(t, derefNilMap, 3)
	assert.Equal(t, 1, derefNilMap["a"])
	assert.Equal(t, 0, derefNilMap["b"]) // zero value
	assert.Equal(t, 2, derefNilMap["c"])

	// Test with nil map
	var nilMap map[string]*int
	derefNilMap2 := DerefMap(nilMap)
	assert.Nil(t, derefNilMap2)

	// Test with empty map
	emptyMap := map[string]*int{}
	derefEmptyMap := DerefMap(emptyMap)
	assert.Len(t, derefEmptyMap, 0)
	assert.NotNil(t, derefEmptyMap)
}
