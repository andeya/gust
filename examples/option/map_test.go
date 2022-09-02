package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/opt"
	"github.com/stretchr/testify/assert"
)

func TestOption_Map_1(t *testing.T) {
	var maybeSomeString = gust.Some("Hello, World!")
	var maybeSomeLen = opt.Map(maybeSomeString, func(s string) int { return len(s) })
	assert.Equal(t, maybeSomeLen, gust.Some(13))
}

func TestOption_Map_2(t *testing.T) {
	var maybeSomeString = gust.Some("Hello, World!")
	var maybeSomeLen = maybeSomeString.XMap(func(s string) any { return len(s) })
	assert.Equal(t, maybeSomeLen, gust.Some(13).ToX())
}
