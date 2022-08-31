package result_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/andeya/gust"
)

func TestInspectErr(t *testing.T) {
	gust.Ret(strconv.Atoi("4x")).
		InspectErr(func(err error) {
			fmt.Printf("failed to convert: %v", err)
		})
}
