package result_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestUnwrapErr_1(t *testing.T) {
	defer func() {
		assert.Equal(t, "called `Result.UnwrapErr()` on an `ok` value: 10", recover())
	}()
	err := gust.Ok(10).UnwrapErr()
	assert.NoError(t, err)
}

func TestUnwrapErr_2(t *testing.T) {
	err := gust.Err[int]("emergency failure").UnwrapErr()
	if assert.Error(t, err) {
		assert.Equal(t, "emergency failure", err.Error())
	} else {
		t.FailNow()
	}
}
