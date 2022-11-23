package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/valconv"
	"github.com/stretchr/testify/assert"
)

func TestOption_GetOrInsertDefault(t *testing.T) {
	none := gust.None[int]()
	assert.Equal(t, valconv.Ref(0), none.GetOrInsertDefault())
	assert.Equal(t, valconv.Ref(0), none.AsPtr())
	some := gust.Some[int](1)
	assert.Equal(t, valconv.Ref(1), some.GetOrInsertDefault())
	assert.Equal(t, valconv.Ref(1), some.AsPtr())
}
