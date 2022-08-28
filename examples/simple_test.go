package examples_test

import (
	"fmt"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/andeya/gust/ret"
)

// ExampleSimple a simple function returning [`Result`] might be defined and used like so:

type Version int8

const (
	Version1 Version = iota + 1
	Version2
)

func ParseVersion(header iter.Iterator[byte]) gust.Result[Version] {
	return ret.AndThen(
		header.Next().
			OkOr("invalid header length"),
		func(b byte) gust.Result[Version] {
			switch b {
			case 1:
				return gust.Ok(Version1)
			case 2:
				return gust.Ok(Version2)
			}
			return gust.Err[Version]("invalid version")
		},
	)
}

func ExampleVersion() {
	ParseVersion(iter.FromElements[byte](1, 2, 3, 4)).
		Inspect(func(v Version) {
			fmt.Printf("working with version: %v\n", v)
		}).
		InspectErr(func(err error) {
			fmt.Printf("error parsing header: %v\n", err)
		})
	// Output:
	// working with version: 1
}
