// Package digit is a package that provides generic-type digit operations.
package digit

import (
	"bytes"
	"errors"
	"fmt"
	"math"

	"github.com/andeya/gust"
)

// FormatByDict convert num into corresponding string according to dict.
func FormatByDict(dict []byte, num uint64) string {
	var base = uint64(len(dict))
	if base == 0 {
		return ""
	}
	var str []byte
	for {
		tmp := make([]byte, len(str)+1)
		tmp[0] = dict[num%base]
		copy(tmp[1:], str)
		str = tmp
		num = num / base
		if num == 0 {
			break
		}
	}
	return string(str)
}

// ParseByDict convert numStr into corresponding digit according to dict.
func ParseByDict[D gust.Digit](dict []byte, numStr string) gust.Result[D] {
	if len(dict) == 0 {
		return gust.Err[D](errors.New("dict is empty"))
	}
	base := float64(len(dict))
	len := len(numStr)
	var number float64
	for i := 0; i < len; i++ {
		char := numStr[i : i+1]
		pos := bytes.IndexAny(dict, char)
		if pos == -1 {
			return gust.Err[D](fmt.Errorf("found a char not included in the dict: %q", char))
		}
		number = math.Pow(base, float64(len-i-1))*float64(pos) + number
	}
	return gust.Ok(D(number))
}
