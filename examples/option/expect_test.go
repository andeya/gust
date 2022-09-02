package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestOption_Expect(t *testing.T) {
	{
		var x = gust.Some("value")
		assert.Equal(t, "value", x.Expect("fruits are healthy"))
	}
	defer func() {
		assert.Equal(t, "fruits are healthy 1", recover())
		defer func() {
			assert.Equal(t, "fruits are healthy 2", recover())
		}()
		var x gust.Option[string]
		x.Expect("fruits are healthy 2") // panics with `fruits are healthy 2`
	}()
	var x = gust.None[string]()
	x.Expect("fruits are healthy 1") // panics with `fruits are healthy 1`
}
