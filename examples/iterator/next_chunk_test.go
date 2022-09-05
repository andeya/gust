package iterator_test

import (
	"testing"

	"github.com/andeya/gust/iter"
	"github.com/stretchr/testify/assert"
)

func TestNextChunk_1(t *testing.T) {
	var i = iter.FromString[rune]("中国-CN")
	assert.Equal(t, []rune{'中', '国'}, i.NextChunk(2).Unwrap())
	assert.Equal(t, []rune{'-', 'C'}, i.NextChunk(2).Unwrap())
	assert.Equal(t, []rune{'N'}, i.NextChunk(4).UnwrapErr())
	assert.Equal(t, []rune{}, i.NextChunk(1).UnwrapErr())
	assert.Equal(t, []rune{}, i.NextChunk(0).Unwrap())
}

func TestNextChunk_2(t *testing.T) {
	var i = iter.FromString[byte]("中国-CN")
	assert.Equal(t, []byte{0xe4, 0xb8, 0xad}, i.NextChunk(3).Unwrap()) // '中'
	assert.Equal(t, []byte{0xe5, 0x9b, 0xbd}, i.NextChunk(3).Unwrap()) // '国'
	assert.Equal(t, []byte{'-'}, i.NextChunk(1).Unwrap())
	assert.Equal(t, []byte{'C'}, i.NextChunk(1).Unwrap())
	assert.Equal(t, []byte{'N'}, i.NextChunk(3).UnwrapErr())
	assert.Equal(t, []byte{}, i.NextChunk(1).UnwrapErr())
	assert.Equal(t, []byte{}, i.NextChunk(0).Unwrap())
}
