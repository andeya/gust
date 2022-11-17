package option_test

import (
	"testing"
	"time"

	"github.com/andeya/gust"
	"github.com/andeya/gust/valconv"
	"github.com/stretchr/testify/assert"
)

func TestOption_UnwrapOrDefault(t *testing.T) {
	assert.Equal(t, "car", gust.Some("car").UnwrapOrDefault())
	assert.Equal(t, "", gust.None[string]().UnwrapOrDefault())
	assert.Equal(t, time.Time{}, gust.None[time.Time]().UnwrapOrDefault())
	assert.Equal(t, &time.Time{}, gust.None[*time.Time]().UnwrapOrDefault())
	assert.Equal(t, valconv.Ref(&time.Time{}), gust.None[**time.Time]().UnwrapOrDefault())
}
