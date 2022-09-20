package digit

import (
	"errors"
	"math"
	"reflect"
	"strconv"

	"github.com/andeya/gust"
)

// TryFromString converts ~string to digit.
// If base == 0, the base is implied by the string's prefix:
// base 2 for "0b", base 8 for "0" or "0o", base 16 for "0x",
// and base 10 otherwise. Also, for base == 0 only, underscore
// characters are permitted per the Go integer literal syntax.
// If base is below 0, is 1, or is above 62, an error is returned.
func TryFromString[T ~string, D gust.Digit](v T, base int, bitSize int) gust.Result[D] {
	return gust.Ret[D](tryFromString[T, D](v, base, bitSize))
}

func tryFromString[T ~string, D gust.Digit](v T, base int, bitSize int) (D, error) {
	var d *D
	var x interface{} = d
	switch x.(type) {
	case *int, *int8, *int16, *int32, *int64:
		r, err := parseInt(string(v), base, bitSize)
		if err != nil {
			return 0, err
		}
		return as[int64, D](r)
	case *uint, *uint8, *uint16, *uint32, *uint64:
		r, err := parseUint(string(v), base, bitSize)
		if err != nil {
			return 0, err
		}
		return as[uint64, D](r)
	case *float32, *float64:
		r, err := strconv.ParseFloat(string(v), bitSize)
		if err != nil {
			return 0, err
		}
		return as[float64, D](r)
	}
	switch reflect.TypeOf(x).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r, err := parseInt(string(v), base, bitSize)
		if err != nil {
			return 0, err
		}
		return as[int64, D](r)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r, err := parseUint(string(v), base, bitSize)
		if err != nil {
			return 0, err
		}
		return as[uint64, D](r)
	case reflect.Float32, reflect.Float64:
		r, err := strconv.ParseFloat(string(v), bitSize)
		if err != nil {
			return 0, err
		}
		return as[float64, D](r)
	}
	return 0, nil
}

func TryFromStrings[T ~string, D gust.Digit](a []T, base int, bitSize int) gust.Result[[]D] {
	return gust.Ret(tryFromStrings[T, D](a, base, bitSize))
}

func tryFromStrings[T ~string, D gust.Digit](a []T, base int, bitSize int) (b []D, err error) {
	b = make([]D, len(a))
	for i, t := range a {
		b[i], err = tryFromString[T, D](t, base, bitSize)
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

// ToBool converts D to bool.
func ToBool[T gust.Digit](v T) bool {
	zero := new(T)
	return v != *zero
}

// ToBools converts []D to []bool.
func ToBools[T gust.Digit](a []T) []bool {
	b := make([]bool, len(a))
	for i, t := range a {
		b[i] = ToBool(t)
	}
	return b
}

// FromBool converts bool to digit.
func FromBool[T ~bool, D gust.Digit](v T) D {
	if v == true {
		return D(1)
	}
	return D(0)
}

func FromBools[T ~bool, D gust.Digit](a []T) (b []D) {
	b = make([]D, len(a))
	for i, t := range a {
		b[i] = FromBool[T, D](t)
	}
	return b
}

func As[T gust.Digit, D gust.Digit](v T) gust.Result[D] {
	return gust.Ret(as[T, D](v))
}

func as[T gust.Digit, D gust.Digit](v T) (D, error) {
	var d *D
	var x interface{} = d
	switch x.(type) {
	case *int:
		r, err := digitToInt(v)
		return D(r), err
	case *int8:
		r, err := digitToInt8[T](v)
		return D(r), err
	case *int16:
		r, err := digitToInt16[T](v)
		return D(r), err
	case *int32:
		r, err := digitToInt32[T](v)
		return D(r), err
	case *int64:
		r, err := digitToInt64[T](v)
		return D(r), err
	case *uint:
		r, err := digitToUint[T](v)
		return D(r), err
	case *uint8:
		r, err := digitToUint8[T](v)
		return D(r), err
	case *uint16:
		r, err := digitToUint16[T](v)
		return D(r), err
	case *uint32:
		r, err := digitToUint32[T](v)
		return D(r), err
	case *uint64:
		r, err := digitToUint64[T](v)
		return D(r), err
	case *float32:
		r, err := digitToFloat32[T](v)
		return D(r), err
	case *float64:
		r := digitToFloat64[T](v)
		return D(r), nil
	}
	switch reflect.TypeOf(x).Kind() {
	case reflect.Int:
		r, err := digitToInt(v)
		return D(r), err
	case reflect.Int8:
		r, err := digitToInt8[T](v)
		return D(r), err
	case reflect.Int16:
		r, err := digitToInt16[T](v)
		return D(r), err
	case reflect.Int32:
		r, err := digitToInt32[T](v)
		return D(r), err
	case reflect.Int64:
		r, err := digitToInt64[T](v)
		return D(r), err
	case reflect.Uint:
		r, err := digitToUint[T](v)
		return D(r), err
	case reflect.Uint8:
		r, err := digitToUint8[T](v)
		return D(r), err
	case reflect.Uint16:
		r, err := digitToUint16[T](v)
		return D(r), err
	case reflect.Uint32:
		r, err := digitToUint32[T](v)
		return D(r), err
	case reflect.Uint64:
		r, err := digitToUint64[T](v)
		return D(r), err
	case reflect.Float32:
		r, err := digitToFloat32[T](v)
		return D(r), err
	case reflect.Float64:
		r := digitToFloat64[T](v)
		return D(r), nil
	}
	return 0, nil
}

// SliceAs creates a copy of the digit slice.
func SliceAs[T gust.Digit, D gust.Digit](a []T) (b []D, err error) {
	b = make([]D, len(a))
	for i, t := range a {
		b[i], err = as[T, D](t)
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

// digitToFloat64 converts digit to float64.
func digitToFloat64[T gust.Digit](v T) float64 {
	return float64(v)
}

// digitToFloat32 converts digit to float32.
func digitToFloat32[T gust.Digit](v T) (float32, error) {
	f := float64(v)
	if f > math.MaxFloat32 || f < -math.MaxFloat32 {
		return 0, errOverflowValue
	}
	return float32(v), nil
}

// digitToInt converts digit to int.
func digitToInt[T gust.Digit](v T) (int, error) {
	f := float64(v)
	if f > math.MaxInt || f < math.MinInt {
		return 0, errOverflowValue
	}
	return int(v), nil
}

// digitToInt8 converts digit to int8.
func digitToInt8[T gust.Digit](v T) (int8, error) {
	if v > 0 {
		if v > math.MaxInt8 {
			return 0, errOverflowValue
		}
	} else {
		if float64(v) < math.MinInt8 {
			return 0, errOverflowValue
		}
	}
	return int8(v), nil
}

// digitToInt16 converts digit to int16.
func digitToInt16[T gust.Digit](v T) (int16, error) {
	f := float64(v)
	if f > math.MaxInt16 || f < math.MinInt16 {
		return 0, errOverflowValue
	}
	return int16(v), nil
}

// digitToInt32 converts digit to int32.
func digitToInt32[T gust.Digit](v T) (int32, error) {
	f := float64(v)
	if f > math.MaxInt32 || f < math.MinInt32 {
		return 0, errOverflowValue
	}
	return int32(v), nil
}

// digitToInt64 converts digit to int64.
func digitToInt64[T gust.Digit](v T) (int64, error) {
	f := float64(v)
	if f > math.MaxInt64 || f < math.MinInt64 {
		return 0, errOverflowValue
	}
	return int64(v), nil
}

// digitToUint converts digit to uint.
func digitToUint[T gust.Digit](v T) (uint, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if float64(v) > math.MaxUint {
		return 0, errOverflowValue
	}
	return uint(v), nil
}

// digitToUint8 converts digit to uint8.
func digitToUint8[T gust.Digit](v T) (uint8, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if float64(v) > math.MaxUint8 {
		return 0, errOverflowValue
	}
	return uint8(v), nil
}

// digitToUint16 converts digit to uint16.
func digitToUint16[T gust.Digit](v T) (uint16, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if float64(v) > math.MaxUint16 {
		return 0, errOverflowValue
	}
	return uint16(v), nil
}

// digitToUint32 converts digit to uint32.
func digitToUint32[T gust.Digit](v T) (uint32, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if float64(v) > math.MaxUint32 {
		return 0, errOverflowValue
	}
	return uint32(v), nil
}

// digitToUint64 converts digit to uint64.
func digitToUint64[T gust.Digit](v T) (uint64, error) {
	if v < 0 {
		return 0, errNegativeValue
	}
	if float64(v) > math.MaxUint64 {
		return 0, errOverflowValue
	}
	return uint64(v), nil
}

var (
	errNegativeValue = errors.New("contains negative value")
	errOverflowValue = errors.New("contains overflow value")
)
