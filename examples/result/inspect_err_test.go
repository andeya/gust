package result_test

import (
	"strconv"
	"testing"

	"github.com/andeya/gust"
)

func TestInspectErr(t *testing.T) {
	gust.Ret(strconv.Atoi("4x")).
		InspectErr(func(err error) {
			t.Logf("failed to convert: %v", err)
		})
}
