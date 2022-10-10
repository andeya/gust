package result_test

import (
	"strconv"
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestExpect(t *testing.T) {
	defer func() {
		assert.Equal(t, "failed to parse number: strconv.Atoi: parsing \"4x\": invalid syntax", recover().(error).Error())
	}()
	gust.Ret(strconv.Atoi("4x")).
		Expect("failed to parse number")

}
