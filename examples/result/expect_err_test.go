package result_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestExpectErr(t *testing.T) {
	defer func() {
		assert.Equal(t, gust.ToErrBox("Testing expect_err: 10"), recover())
	}()
	err := gust.Ok(10).ExpectErr("Testing expect_err")
	assert.NoError(t, err)
}
