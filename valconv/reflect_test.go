package valconv

import (
	"reflect"
	"testing"
)

func TestRefSliceValue(t *testing.T) {
	v := reflect.ValueOf([]int{1, 2})
	v = RefSliceValue(v, 1)
	ret := v.Interface().([]*int)
	t.Logf("%#v", ret)

	v = reflect.ValueOf([]int{})
	v = RefSliceValue(v, 1)
	ret = v.Interface().([]*int)
	t.Logf("%#v", ret)
}
