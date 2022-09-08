package option_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/opt"
)

func TestOption(t *testing.T) {
	var divide = func(numerator, denominator float64) gust.Option[float64] {
		if denominator == 0.0 {
			return gust.None[float64]()
		}
		return gust.Some(numerator / denominator)
	}
	// The return value of the function is an option
	divide(2.0, 3.0).
		Inspect(func(x float64) {
			// Pattern match to retrieve the value
			t.Log("Result:", x)
		}).
		InspectNone(func() {
			t.Log("Cannot divide by 0")
		})
}

func ExampleOption() {
	type A struct {
		X int
	}
	var a = gust.Some(A{X: 1})
	fmt.Println(a.IsSome(), a.IsNone())

	var b = gust.None[A]()
	fmt.Println(b.IsSome(), b.IsNone())

	var x = b.UnwrapOr(A{X: 2})
	fmt.Println(x)

	var c *A
	fmt.Println(gust.PtrOpt(c).IsNone())
	c = new(A)
	fmt.Println(gust.PtrOpt(c).IsNone())

	type B struct {
		Y string
	}
	var d = opt.Map(a, func(t A) B {
		return B{
			Y: strconv.Itoa(t.X),
		}
	})
	fmt.Println(d)

	// Output:
	// true false
	// false true
	// {2}
	// true
	// false
	// Some({1})
}
