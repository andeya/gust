package option_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/opt"
	"github.com/stretchr/testify/assert"
)

func TestOption_ZipWith(t *testing.T) {
	type Point struct {
		x float64
		y float64
	}
	var newPoint = func(x float64, y float64) Point {
		return Point{x, y}
	}
	var x = gust.Some(17.5)
	var y = gust.Some(42.7)
	assert.Equal(t, opt.ZipWith(x, y, newPoint), gust.Some(Point{x: 17.5, y: 42.7}))
	assert.Equal(t, opt.ZipWith(x, gust.None[float64](), newPoint), gust.None[Point]())
}
