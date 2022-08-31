package result_test

import (
	"fmt"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/ret"
	"github.com/stretchr/testify/assert"
)

func TestAssert(t *testing.T) {
	r := gust.Ok[any]("hello")
	assert.Equal(t, "gust.Result[string]", fmt.Sprintf("%T", ret.Assert[any, string](r)))
}

func TestXAssert(t *testing.T) {
	r := gust.Ok[any]("hello")
	assert.Equal(t, "gust.Result[string]", fmt.Sprintf("%T", ret.XAssert[string](r)))
}
