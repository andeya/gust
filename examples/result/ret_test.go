package result_test

import (
	"strconv"
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestResult_Ret(t *testing.T) {
	var w = gust.Ret[int](strconv.Atoi("s"))
	assert.False(t, w.IsOk())
	assert.True(t, w.IsErr())

	var w2 = gust.Ret[any](strconv.Atoi("-1"))
	assert.True(t, w2.IsOk())
	assert.False(t, w2.IsErr())
	assert.Equal(t, -1, w2.Unwrap())
}
