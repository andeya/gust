package result_test

import (
	"strconv"
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestUnwrap(t *testing.T) {
	defer func() {
		assert.Equal(t, "strconv.Atoi: parsing \"4x\": invalid syntax", recover().(error).Error())
	}()
	gust.Ret(strconv.Atoi("4x")).Unwrap()
}
