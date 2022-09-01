package result_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/ret"
	"github.com/stretchr/testify/assert"
)

func TestResult_ContainsErr(t *testing.T) {
	assert.False(t, gust.Ok(2).ContainsErr("Some error message"))
	assert.True(t, ret.Contains(gust.Ok(2), 2))
	assert.False(t, ret.Contains(gust.Ok(3), 2))
	assert.False(t, ret.Contains(gust.Err[int]("Some error message"), 2))
}
