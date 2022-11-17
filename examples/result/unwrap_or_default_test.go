package result_test

import (
	"testing"
	"time"

	"github.com/andeya/gust"
	"github.com/andeya/gust/valconv"
	"github.com/stretchr/testify/assert"
)

func TestUnwrapOrDefault(t *testing.T) {
	assert.Equal(t, "car", gust.Ok("car").UnwrapOrDefault())
	assert.Equal(t, "", gust.Err[string](nil).UnwrapOrDefault())
	assert.Equal(t, time.Time{}, gust.Err[time.Time](nil).UnwrapOrDefault())
	assert.Equal(t, &time.Time{}, gust.Err[*time.Time](nil).UnwrapOrDefault())
	assert.Equal(t, valconv.Ref(&time.Time{}), gust.Err[**time.Time](nil).UnwrapOrDefault())
}
