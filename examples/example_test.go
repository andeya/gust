package examples_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/iter"
	"github.com/andeya/gust/ret"
	"github.com/stretchr/testify/assert"
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

// You might want to use an iterator chain to do multiple instances of an
// operation that can fail, but would like to ignore failures while
// continuing to process the successful results. In this example, we take
// advantage of the iterable nature of [`gust.Result`] to select only the
// [`gust.Ok`] values using [`iter.Flatten`].
func TestResultFlatten(t *testing.T) {
	var results []gust.Result[uint64]
	var errs []error
	var nums = iter.
		Flatten[uint64, gust.Result[uint64]](
		iter.Map[string, gust.Result[uint64]](
			iter.FromElements("17", "not a number", "99", "-27", "768"),
			func(s string) gust.Result[uint64] { return gust.Ret(strconv.ParseUint(s, 10, 64)) },
		).
			Inspect(func(x gust.Result[uint64]) {
				// Save clones of the raw `Result` values to inspect
				results = append(results,
					x.InspectErr(func(err error) {
						// Challenge: explain how this captures only the `Err` values
						errs = append(errs, err)
					}))
			}),
	).Collect()
	assert.Equal(t, 2, len(errs))
	assert.Equal(t, []uint64{17, 99, 768}, nums)
	fmt.Printf("results %v\n", results)
	fmt.Printf("errs %v\n", errs)
	fmt.Printf("nums %v\n", nums)
}

func TestResultOption(t *testing.T) {
	var okStr = "20"
	var errStr = "not a number"
	var emptyStr string
	// Ret encapsulates value and err into result
	okResult := gust.Ret(strconv.ParseUint(okStr, 10, 64))
	errResult := gust.Ret(strconv.ParseUint(errStr, 10, 64))
	emptyResult := gust.Ret(strconv.ParseUint(emptyStr, 10, 64))
	assert.Equal(t, uint64(20), okResult.Unwrap())
	assert.Equal(t, true, errResult.IsErr())
	assert.Equal(t, true, emptyResult.IsErr())
	// If an error occurs Ok will set it to empty
	// If you encounter the return value is empty but err is also empty, You can use ok to convert result to option
	okOption := okResult.Ok()
	errOption := errResult.Ok()
	emptyOption := emptyResult.Ok()
	assert.Equal(t, uint64(20), okOption.Unwrap())
	assert.Equal(t, true, errOption.IsNone())
	assert.Equal(t, true, emptyOption.IsNone())
}
