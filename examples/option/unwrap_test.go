package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOption_Unwrap(t *testing.T) {
	{
		var x = gust.Some("air")
		assert.Equal(t, "air", x.Unwrap())
	}
	defer func() {
		assert.Equal(t, gust.ToErrBox("call Option[string].Unwrap() on none"), recover())
	}()
	var x = gust.None[string]()
	x.Unwrap()
}
