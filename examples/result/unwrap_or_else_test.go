package result_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestUnwrapOrElse(t *testing.T) {
	var count = func(x error) int {
		return len(x.Error())
	}
	assert.Equal(t, 2, gust.Ok(2).UnwrapOrElse(count))
	assert.Equal(t, 3, gust.Err[int]("foo").UnwrapOrElse(count))
}
