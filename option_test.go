package gust_test

import (
	"fmt"
	"strconv"

	"github.com/henrylee2cn/gust"
	"github.com/henrylee2cn/gust/opt"
)

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

	type B struct {
		Y string
	}
	var c = opt.Map(a, func(t A) B {
		return B{
			Y: strconv.Itoa(t.X),
		}
	})
	fmt.Println(c)

	// Output:
	// true false
	// false true
	// {2}
	// Some({1})
}
