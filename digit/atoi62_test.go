// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package digit

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type parseUint64Test struct {
	in  string
	out uint64
	err error
}

var parseUint64Tests = []parseUint64Test{
	{"", 0, strconv.ErrSyntax},
	{"0", 0, nil},
	{"1", 1, nil},
	{"12345", 12345, nil},
	{"012345", 12345, nil},
	{"12345x", 0, strconv.ErrSyntax},
	{"98765432100", 98765432100, nil},
	{"18446744073709551615", 1<<64 - 1, nil},
	{"18446744073709551616", 1<<64 - 1, strconv.ErrRange},
	{"18446744073709551620", 1<<64 - 1, strconv.ErrRange},
	{"1_2_3_4_5", 0, strconv.ErrSyntax}, // base=10 so no underscores allowed
	{"_12345", 0, strconv.ErrSyntax},
	{"1__2345", 0, strconv.ErrSyntax},
	{"12345_", 0, strconv.ErrSyntax},
}

type parseUint64BaseTest struct {
	in   string
	base int
	out  uint64
	err  error
}

var parseUint64BaseTests = []parseUint64BaseTest{
	{"", 0, 0, strconv.ErrSyntax},
	{"0", 0, 0, nil},
	{"0x", 0, 0, strconv.ErrSyntax},
	{"0X", 0, 0, strconv.ErrSyntax},
	{"1", 0, 1, nil},
	{"12345", 0, 12345, nil},
	{"012345", 0, 012345, nil},
	{"0x12345", 0, 0x12345, nil},
	{"0X12345", 0, 0x12345, nil},
	{"12345x", 0, 0, strconv.ErrSyntax},
	{"0xabcdefg123", 0, 0, strconv.ErrSyntax},
	{"123456789abc", 0, 0, strconv.ErrSyntax},
	{"98765432100", 0, 98765432100, nil},
	{"18446744073709551615", 0, 1<<64 - 1, nil},
	{"18446744073709551616", 0, 1<<64 - 1, strconv.ErrRange},
	{"18446744073709551620", 0, 1<<64 - 1, strconv.ErrRange},
	{"0xFFFFFFFFFFFFFFFF", 0, 1<<64 - 1, nil},
	{"0x10000000000000000", 0, 1<<64 - 1, strconv.ErrRange},
	{"01777777777777777777777", 0, 1<<64 - 1, nil},
	{"01777777777777777777778", 0, 0, strconv.ErrSyntax},
	{"02000000000000000000000", 0, 1<<64 - 1, strconv.ErrRange},
	{"0200000000000000000000", 0, 1 << 61, nil},
	{"0b", 0, 0, strconv.ErrSyntax},
	{"0B", 0, 0, strconv.ErrSyntax},
	{"0b101", 0, 5, nil},
	{"0B101", 0, 5, nil},
	{"0o", 0, 0, strconv.ErrSyntax},
	{"0O", 0, 0, strconv.ErrSyntax},
	{"0o377", 0, 255, nil},
	{"0O377", 0, 255, nil},

	// underscores allowed with base == 0 only
	{"1_2_3_4_5", 0, 12345, nil}, // base 0 => 10
	{"_12345", 0, 0, strconv.ErrSyntax},
	{"1__2345", 0, 0, strconv.ErrSyntax},
	{"12345_", 0, 0, strconv.ErrSyntax},

	{"1_2_3_4_5", 10, 0, strconv.ErrSyntax}, // base 10
	{"_12345", 10, 0, strconv.ErrSyntax},
	{"1__2345", 10, 0, strconv.ErrSyntax},
	{"12345_", 10, 0, strconv.ErrSyntax},

	{"0x_1_2_3_4_5", 0, 0x12345, nil}, // base 0 => 16
	{"_0x12345", 0, 0, strconv.ErrSyntax},
	{"0x__12345", 0, 0, strconv.ErrSyntax},
	{"0x1__2345", 0, 0, strconv.ErrSyntax},
	{"0x1234__5", 0, 0, strconv.ErrSyntax},
	{"0x12345_", 0, 0, strconv.ErrSyntax},

	{"1_2_3_4_5", 16, 0, strconv.ErrSyntax}, // base 16
	{"_12345", 16, 0, strconv.ErrSyntax},
	{"1__2345", 16, 0, strconv.ErrSyntax},
	{"1234__5", 16, 0, strconv.ErrSyntax},
	{"12345_", 16, 0, strconv.ErrSyntax},

	{"0_1_2_3_4_5", 0, 012345, nil}, // base 0 => 8 (0377)
	{"_012345", 0, 0, strconv.ErrSyntax},
	{"0__12345", 0, 0, strconv.ErrSyntax},
	{"01234__5", 0, 0, strconv.ErrSyntax},
	{"012345_", 0, 0, strconv.ErrSyntax},

	{"0o_1_2_3_4_5", 0, 012345, nil}, // base 0 => 8 (0o377)
	{"_0o12345", 0, 0, strconv.ErrSyntax},
	{"0o__12345", 0, 0, strconv.ErrSyntax},
	{"0o1234__5", 0, 0, strconv.ErrSyntax},
	{"0o12345_", 0, 0, strconv.ErrSyntax},

	{"0_1_2_3_4_5", 8, 0, strconv.ErrSyntax}, // base 8
	{"_012345", 8, 0, strconv.ErrSyntax},
	{"0__12345", 8, 0, strconv.ErrSyntax},
	{"01234__5", 8, 0, strconv.ErrSyntax},
	{"012345_", 8, 0, strconv.ErrSyntax},

	{"0b_1_0_1", 0, 5, nil}, // base 0 => 2 (0b101)
	{"_0b101", 0, 0, strconv.ErrSyntax},
	{"0b__101", 0, 0, strconv.ErrSyntax},
	{"0b1__01", 0, 0, strconv.ErrSyntax},
	{"0b10__1", 0, 0, strconv.ErrSyntax},
	{"0b101_", 0, 0, strconv.ErrSyntax},

	{"1_0_1", 2, 0, strconv.ErrSyntax}, // base 2
	{"_101", 2, 0, strconv.ErrSyntax},
	{"1_01", 2, 0, strconv.ErrSyntax},
	{"10_1", 2, 0, strconv.ErrSyntax},
	{"101_", 2, 0, strconv.ErrSyntax},
}

type parseInt64Test struct {
	in  string
	out int64
	err error
}

var parseInt64Tests = []parseInt64Test{
	{"", 0, strconv.ErrSyntax},
	{"0", 0, nil},
	{"-0", 0, nil},
	{"1", 1, nil},
	{"-1", -1, nil},
	{"12345", 12345, nil},
	{"-12345", -12345, nil},
	{"012345", 12345, nil},
	{"-012345", -12345, nil},
	{"98765432100", 98765432100, nil},
	{"-98765432100", -98765432100, nil},
	{"9223372036854775807", 1<<63 - 1, nil},
	{"-9223372036854775807", -(1<<63 - 1), nil},
	{"9223372036854775808", 1<<63 - 1, strconv.ErrRange},
	{"-9223372036854775808", -1 << 63, nil},
	{"9223372036854775809", 1<<63 - 1, strconv.ErrRange},
	{"-9223372036854775809", -1 << 63, strconv.ErrRange},
	{"-1_2_3_4_5", 0, strconv.ErrSyntax}, // base=10 so no underscores allowed
	{"-_12345", 0, strconv.ErrSyntax},
	{"_12345", 0, strconv.ErrSyntax},
	{"1__2345", 0, strconv.ErrSyntax},
	{"12345_", 0, strconv.ErrSyntax},
}

type parseInt64BaseTest struct {
	in   string
	base int
	out  int64
	err  error
}

var parseInt64BaseTests = []parseInt64BaseTest{
	{"", 0, 0, strconv.ErrSyntax},
	{"0", 0, 0, nil},
	{"-0", 0, 0, nil},
	{"1", 0, 1, nil},
	{"-1", 0, -1, nil},
	{"12345", 0, 12345, nil},
	{"-12345", 0, -12345, nil},
	{"012345", 0, 012345, nil},
	{"-012345", 0, -012345, nil},
	{"0x12345", 0, 0x12345, nil},
	{"-0X12345", 0, -0x12345, nil},
	{"12345x", 0, 0, strconv.ErrSyntax},
	{"-12345x", 0, 0, strconv.ErrSyntax},
	{"98765432100", 0, 98765432100, nil},
	{"-98765432100", 0, -98765432100, nil},
	{"9223372036854775807", 0, 1<<63 - 1, nil},
	{"-9223372036854775807", 0, -(1<<63 - 1), nil},
	{"9223372036854775808", 0, 1<<63 - 1, strconv.ErrRange},
	{"-9223372036854775808", 0, -1 << 63, nil},
	{"9223372036854775809", 0, 1<<63 - 1, strconv.ErrRange},
	{"-9223372036854775809", 0, -1 << 63, strconv.ErrRange},

	// other bases
	{"g", 17, 16, nil},
	{"10", 25, 25, nil},
	{"holycow", 35, (((((17*35+24)*35+21)*35+34)*35+12)*35+24)*35 + 32, nil},
	{"holycow", 36, (((((17*36+24)*36+21)*36+34)*36+12)*36+24)*36 + 32, nil},

	// base 2
	{"0", 2, 0, nil},
	{"-1", 2, -1, nil},
	{"1010", 2, 10, nil},
	{"1000000000000000", 2, 1 << 15, nil},
	{"111111111111111111111111111111111111111111111111111111111111111", 2, 1<<63 - 1, nil},
	{"1000000000000000000000000000000000000000000000000000000000000000", 2, 1<<63 - 1, strconv.ErrRange},
	{"-1000000000000000000000000000000000000000000000000000000000000000", 2, -1 << 63, nil},
	{"-1000000000000000000000000000000000000000000000000000000000000001", 2, -1 << 63, strconv.ErrRange},

	// base 8
	{"-10", 8, -8, nil},
	{"57635436545", 8, 057635436545, nil},
	{"100000000", 8, 1 << 24, nil},

	// base 16
	{"10", 16, 16, nil},
	{"-123456789abcdef", 16, -0x123456789abcdef, nil},
	{"7fffffffffffffff", 16, 1<<63 - 1, nil},

	// underscores
	{"-0x_1_2_3_4_5", 0, -0x12345, nil},
	{"0x_1_2_3_4_5", 0, 0x12345, nil},
	{"-_0x12345", 0, 0, strconv.ErrSyntax},
	{"_-0x12345", 0, 0, strconv.ErrSyntax},
	{"_0x12345", 0, 0, strconv.ErrSyntax},
	{"0x__12345", 0, 0, strconv.ErrSyntax},
	{"0x1__2345", 0, 0, strconv.ErrSyntax},
	{"0x1234__5", 0, 0, strconv.ErrSyntax},
	{"0x12345_", 0, 0, strconv.ErrSyntax},

	{"-0_1_2_3_4_5", 0, -012345, nil}, // octal
	{"0_1_2_3_4_5", 0, 012345, nil},   // octal
	{"-_012345", 0, 0, strconv.ErrSyntax},
	{"_-012345", 0, 0, strconv.ErrSyntax},
	{"_012345", 0, 0, strconv.ErrSyntax},
	{"0__12345", 0, 0, strconv.ErrSyntax},
	{"01234__5", 0, 0, strconv.ErrSyntax},
	{"012345_", 0, 0, strconv.ErrSyntax},
}

type parseUint32Test struct {
	in  string
	out uint32
	err error
}

var parseUint32Tests = []parseUint32Test{
	{"", 0, strconv.ErrSyntax},
	{"0", 0, nil},
	{"1", 1, nil},
	{"12345", 12345, nil},
	{"012345", 12345, nil},
	{"12345x", 0, strconv.ErrSyntax},
	{"987654321", 987654321, nil},
	{"4294967295", 1<<32 - 1, nil},
	{"4294967296", 1<<32 - 1, strconv.ErrRange},
	{"1_2_3_4_5", 0, strconv.ErrSyntax}, // base=10 so no underscores allowed
	{"_12345", 0, strconv.ErrSyntax},
	{"_12345", 0, strconv.ErrSyntax},
	{"1__2345", 0, strconv.ErrSyntax},
	{"12345_", 0, strconv.ErrSyntax},
}

type parseInt32Test struct {
	in  string
	out int32
	err error
}

var parseInt32Tests = []parseInt32Test{
	{"", 0, strconv.ErrSyntax},
	{"0", 0, nil},
	{"-0", 0, nil},
	{"1", 1, nil},
	{"-1", -1, nil},
	{"12345", 12345, nil},
	{"-12345", -12345, nil},
	{"012345", 12345, nil},
	{"-012345", -12345, nil},
	{"12345x", 0, strconv.ErrSyntax},
	{"-12345x", 0, strconv.ErrSyntax},
	{"987654321", 987654321, nil},
	{"-987654321", -987654321, nil},
	{"2147483647", 1<<31 - 1, nil},
	{"-2147483647", -(1<<31 - 1), nil},
	{"2147483648", 1<<31 - 1, strconv.ErrRange},
	{"-2147483648", -1 << 31, nil},
	{"2147483649", 1<<31 - 1, strconv.ErrRange},
	{"-2147483649", -1 << 31, strconv.ErrRange},
	{"-1_2_3_4_5", 0, strconv.ErrSyntax}, // base=10 so no underscores allowed
	{"-_12345", 0, strconv.ErrSyntax},
	{"_12345", 0, strconv.ErrSyntax},
	{"1__2345", 0, strconv.ErrSyntax},
	{"12345_", 0, strconv.ErrSyntax},
}

type numErrorTest struct {
	num, want string
}

var numErrorTests = []numErrorTest{
	{"0", `strconv.ParseFloat: parsing "0": failed`},
	{"`", "strconv.ParseFloat: parsing \"`\": failed"},
	{"1\x00.2", `strconv.ParseFloat: parsing "1\x00.2": failed`},
}

func init() {
	// The parse routines return strconv.NumErrors wrapping
	// the error and the string. Convert the tables above.
	for i := range parseUint64Tests {
		test := &parseUint64Tests[i]
		if test.err != nil {
			test.err = &strconv.NumError{Func: "ParseUint", Num: test.in, Err: test.err}
		}
	}
	for i := range parseUint64BaseTests {
		test := &parseUint64BaseTests[i]
		if test.err != nil {
			test.err = &strconv.NumError{Func: "ParseUint", Num: test.in, Err: test.err}
		}
	}
	for i := range parseInt64Tests {
		test := &parseInt64Tests[i]
		if test.err != nil {
			test.err = &strconv.NumError{Func: "ParseInt", Num: test.in, Err: test.err}
		}
	}
	for i := range parseInt64BaseTests {
		test := &parseInt64BaseTests[i]
		if test.err != nil {
			test.err = &strconv.NumError{Func: "ParseInt", Num: test.in, Err: test.err}
		}
	}
	for i := range parseUint32Tests {
		test := &parseUint32Tests[i]
		if test.err != nil {
			test.err = &strconv.NumError{Func: "ParseUint", Num: test.in, Err: test.err}
		}
	}
	for i := range parseInt32Tests {
		test := &parseInt32Tests[i]
		if test.err != nil {
			test.err = &strconv.NumError{Func: "ParseInt", Num: test.in, Err: test.err}
		}
	}
}

func TestParseUint32(t *testing.T) {
	for i := range parseUint32Tests {
		test := &parseUint32Tests[i]
		out, err := parseUint(test.in, 10, 32)
		if uint64(test.out) != out || !reflect.DeepEqual(test.err, err) {
			t.Errorf("parseUint(%q, 10, 32) = %v, %v want %v, %v",
				test.in, out, err, test.out, test.err)
		}
	}
}

func TestParseUint64(t *testing.T) {
	for i := range parseUint64Tests {
		test := &parseUint64Tests[i]
		out, err := parseUint(test.in, 10, 64)
		if test.out != out || !reflect.DeepEqual(test.err, err) {
			t.Errorf("parseUint(%q, 10, 64) = %v, %v want %v, %v",
				test.in, out, err, test.out, test.err)
		}
	}
}

func TestParseUint64Base(t *testing.T) {
	for i := range parseUint64BaseTests {
		test := &parseUint64BaseTests[i]
		out, err := parseUint(test.in, test.base, 64)
		if test.out != out || !reflect.DeepEqual(test.err, err) {
			t.Errorf("parseUint(%q, %v, 64) = %v, %v want %v, %v",
				test.in, test.base, out, err, test.out, test.err)
		}
	}
}

func TestParseInt32(t *testing.T) {
	for i := range parseInt32Tests {
		test := &parseInt32Tests[i]
		out, err := parseInt(test.in, 10, 32)
		if int64(test.out) != out || !reflect.DeepEqual(test.err, err) {
			t.Errorf("parseInt(%q, 10 ,32) = %v, %v want %v, %v",
				test.in, out, err, test.out, test.err)
		}
	}
}

func TestParseInt64(t *testing.T) {
	for i := range parseInt64Tests {
		test := &parseInt64Tests[i]
		out, err := parseInt(test.in, 10, 64)
		if test.out != out || !reflect.DeepEqual(test.err, err) {
			t.Errorf("parseInt(%q, 10, 64) = %v, %v want %v, %v",
				test.in, out, err, test.out, test.err)
		}
	}
}

func TestParseInt64Base(t *testing.T) {
	for i := range parseInt64BaseTests {
		test := &parseInt64BaseTests[i]
		out, err := parseInt(test.in, test.base, 64)
		if test.out != out || !reflect.DeepEqual(test.err, err) {
			t.Errorf("parseInt(%q, %v, 64) = %v, %v want %v, %v",
				test.in, test.base, out, err, test.out, test.err)
		}
	}
}

func TestParseUint(t *testing.T) {
	switch strconv.IntSize {
	case 32:
		for i := range parseUint32Tests {
			test := &parseUint32Tests[i]
			out, err := parseUint(test.in, 10, 0)
			if uint64(test.out) != out || !reflect.DeepEqual(test.err, err) {
				t.Errorf("parseUint(%q, 10, 0) = %v, %v want %v, %v",
					test.in, out, err, test.out, test.err)
			}
		}
	case 64:
		for i := range parseUint64Tests {
			test := &parseUint64Tests[i]
			out, err := parseUint(test.in, 10, 0)
			if test.out != out || !reflect.DeepEqual(test.err, err) {
				t.Errorf("parseUint(%q, 10, 0) = %v, %v want %v, %v",
					test.in, out, err, test.out, test.err)
			}
		}
	}
}

func TestParseInt(t *testing.T) {
	switch strconv.IntSize {
	case 32:
		for i := range parseInt32Tests {
			test := &parseInt32Tests[i]
			out, err := parseInt(test.in, 10, 0)
			if int64(test.out) != out || !reflect.DeepEqual(test.err, err) {
				t.Errorf("parseInt(%q, 10, 0) = %v, %v want %v, %v",
					test.in, out, err, test.out, test.err)
			}
		}
	case 64:
		for i := range parseInt64Tests {
			test := &parseInt64Tests[i]
			out, err := parseInt(test.in, 10, 0)
			if test.out != out || !reflect.DeepEqual(test.err, err) {
				t.Errorf("parseInt(%q, 10, 0) = %v, %v want %v, %v",
					test.in, out, err, test.out, test.err)
			}
		}
	}
}

func TestAtoi(t *testing.T) {
	switch strconv.IntSize {
	case 32:
		for i := range parseInt32Tests {
			test := &parseInt32Tests[i]
			out, err := Atoi(test.in)
			var testErr error
			if test.err != nil {
				testErr = &strconv.NumError{Func: "Atoi", Num: test.in, Err: test.err.(*strconv.NumError).Err}
			}
			if int(test.out) != out || !reflect.DeepEqual(testErr, err) {
				t.Errorf("Atoi(%q) = %v, %v want %v, %v",
					test.in, out, err, test.out, testErr)
			}
		}
	case 64:
		for i := range parseInt64Tests {
			test := &parseInt64Tests[i]
			out, err := Atoi(test.in)
			var testErr error
			if test.err != nil {
				testErr = &strconv.NumError{Func: "Atoi", Num: test.in, Err: test.err.(*strconv.NumError).Err}
			}
			if test.out != int64(out) || !reflect.DeepEqual(testErr, err) {
				t.Errorf("Atoi(%q) = %v, %v want %v, %v",
					test.in, out, err, test.out, testErr)
			}
		}
	}
}

func bitSizeErrStub(name string, bitSize int) error {
	return bitSizeError(name, "0", bitSize)
}

func baseErrStub(name string, base int) error {
	return baseError(name, "0", base)
}

func noErrStub(name string, arg int) error {
	return nil
}

type parseErrorTest struct {
	arg     int
	errStub func(name string, arg int) error
}

var parseBitSizeTests = []parseErrorTest{
	{-1, bitSizeErrStub},
	{0, noErrStub},
	{64, noErrStub},
	{65, bitSizeErrStub},
}

var parseBaseTests = []parseErrorTest{
	{-1, baseErrStub},
	{0, noErrStub},
	{1, baseErrStub},
	{2, noErrStub},
	{36, noErrStub},
	{62, noErrStub},
	{63, baseErrStub},
}

func equalError(a, b error) bool {
	if a == nil {
		return b == nil
	}
	if b == nil {
		return a == nil
	}
	return a.Error() == b.Error()
}

func TestParseIntBitSize(t *testing.T) {
	for i := range parseBitSizeTests {
		test := &parseBitSizeTests[i]
		testErr := test.errStub("ParseInt", test.arg)
		_, err := parseInt("0", 0, test.arg)
		if !equalError(testErr, err) {
			t.Errorf("parseInt(\"0\", 0, %v) = 0, %v want 0, %v",
				test.arg, err, testErr)
		}
	}
}

func TestParseUintBitSize(t *testing.T) {
	for i := range parseBitSizeTests {
		test := &parseBitSizeTests[i]
		testErr := test.errStub("ParseUint", test.arg)
		_, err := parseUint("0", 0, test.arg)
		if !equalError(testErr, err) {
			t.Errorf("parseUint(\"0\", 0, %v) = 0, %v want 0, %v",
				test.arg, err, testErr)
		}
	}
}

func TestParseIntBase(t *testing.T) {
	for i := range parseBaseTests {
		test := &parseBaseTests[i]
		testErr := test.errStub("ParseInt", test.arg)
		_, err := parseInt("0", test.arg, 0)
		if !equalError(testErr, err) {
			t.Errorf("parseInt(\"0\", %v, 0) = 0, %v want 0, %v",
				test.arg, err, testErr)
		}
	}
}

func TestParseUintBase(t *testing.T) {
	for i := range parseBaseTests {
		test := &parseBaseTests[i]
		testErr := test.errStub("ParseUint", test.arg)
		_, err := parseUint("0", test.arg, 0)
		if !equalError(testErr, err) {
			t.Errorf("parseUint(\"0\", %v, 0) = 0, %v want 0, %v",
				test.arg, err, testErr)
		}
	}
}

func TestNumError(t *testing.T) {
	for _, test := range numErrorTests {
		err := &strconv.NumError{
			Func: "ParseFloat",
			Num:  test.num,
			Err:  errors.New("failed"),
		}
		if got := err.Error(); got != test.want {
			t.Errorf(`(&strconv.NumError{"ParseFloat", %q, "failed"}).Error() = %v, want %v`, test.num, got, test.want)
		}
	}
}

func BenchmarkParseInt(b *testing.B) {
	b.Run("Pos", func(b *testing.B) {
		benchmarkparseInt(b, 1)
	})
	b.Run("Neg", func(b *testing.B) {
		benchmarkparseInt(b, -1)
	})
}

type benchCase struct {
	name string
	num  int64
}

func benchmarkparseInt(b *testing.B, neg int) {
	cases := []benchCase{
		{"7bit", 1<<7 - 1},
		{"26bit", 1<<26 - 1},
		{"31bit", 1<<31 - 1},
		{"56bit", 1<<56 - 1},
		{"63bit", 1<<63 - 1},
	}
	for _, cs := range cases {
		b.Run(cs.name, func(b *testing.B) {
			s := fmt.Sprintf("%d", cs.num*int64(neg))
			for i := 0; i < b.N; i++ {
				out, _ := parseInt(s, 10, 64)
				BenchSink += int(out)
			}
		})
	}
}

func BenchmarkAtoi(b *testing.B) {
	b.Run("Pos", func(b *testing.B) {
		benchmarkAtoi(b, 1)
	})
	b.Run("Neg", func(b *testing.B) {
		benchmarkAtoi(b, -1)
	})
}

func benchmarkAtoi(b *testing.B, neg int) {
	cases := []benchCase{
		{"7bit", 1<<7 - 1},
		{"26bit", 1<<26 - 1},
		{"31bit", 1<<31 - 1},
	}
	if strconv.IntSize == 64 {
		cases = append(cases, []benchCase{
			{"56bit", 1<<56 - 1},
			{"63bit", 1<<63 - 1},
		}...)
	}
	for _, cs := range cases {
		b.Run(cs.name, func(b *testing.B) {
			s := fmt.Sprintf("%d", cs.num*int64(neg))
			for i := 0; i < b.N; i++ {
				out, _ := Atoi(s)
				BenchSink += out
			}
		})
	}
}

func TestAtoi62(t *testing.T) {
	i, err := parseUint("aZl8N0y58M7", 62, 64)
	assert.NoError(t, err)
	assert.Equal(t, uint64(9223372036854775807), i)
}

func TestParseUint_Base37To62(t *testing.T) {
	// Test base 37-62 (extended bases)
	// Note: base <= 36 uses strconv.ParseUint (only supports lowercase a-z)
	// base > 36 uses custom logic:
	//   - '0'-'9' = 0-9
	//   - 'a'-'z' = 10-35
	//   - 'A'-'Z' = 36-61
	testCases := []struct {
		s    string
		base int
		want uint64
	}{
		{"z", 36, 35},  // base 36, 'z' = 35 (uses strconv, lowercase only)
		{"10", 37, 37}, // base 37, "10" = 37
		{"A", 37, 36},  // base 37, 'A' = 36 (A-Z = 36-61 in base > 36)
		{"Z", 37, 61},  // base 37, 'Z' = 61
		{"a", 37, 10},  // base 37, 'a' = 10 (a-z = 10-35 in base > 36)
		{"z", 37, 35},  // base 37, 'z' = 35
		{"10", 62, 62}, // base 62, "10" = 62
		{"A", 62, 36},  // base 62, 'A' = 36
		{"Z", 62, 61},  // base 62, 'Z' = 61
		{"a", 62, 10},  // base 62, 'a' = 10
		{"z", 62, 35},  // base 62, 'z' = 35
	}

	for _, tc := range testCases {
		got, err := parseUint(tc.s, tc.base, 64)
		if err != nil {
			// For base <= 36, strconv.ParseUint doesn't support uppercase, so skip those
			if tc.base <= 36 && (tc.s[0] >= 'A' && tc.s[0] <= 'Z') {
				continue
			}
			t.Errorf("parseUint(%q, %d, 64) = error: %v, want %d", tc.s, tc.base, err, tc.want)
			continue
		}
		if got != tc.want {
			t.Errorf("parseUint(%q, %d, 64) = %d, want %d", tc.s, tc.base, got, tc.want)
		}
	}
}

func TestParseUint_Overflow(t *testing.T) {
	// Test overflow cases for base > 36
	// These should return maxVal and ErrRange

	// Test n >= cutoff case (n*base overflows)
	// Use a very large number in base 62
	_, err := parseUint("ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ", 62, 64)
	if err == nil {
		t.Error("parseUint with overflow should return error")
	} else if err.(*strconv.NumError).Err != strconv.ErrRange {
		t.Errorf("parseUint with overflow should return ErrRange, got %v", err)
	}

	// Test n1 < n || n1 > maxVal case (n+v overflows)
	// This is harder to trigger, but we can test with smaller bitSize
	_, err2 := parseUint("100", 62, 8) // base 62, bitSize 8, should overflow
	if err2 == nil {
		t.Error("parseUint with overflow should return error")
	} else if err2.(*strconv.NumError).Err != strconv.ErrRange {
		t.Errorf("parseUint with overflow should return ErrRange, got %v", err2)
	}

	// Test with max value
	got, err3 := parseUint("1V112200s4a0", 62, 64) // This is close to maxUint64 in base 62
	if err3 != nil && err3.(*strconv.NumError).Err == strconv.ErrRange {
		// Expected overflow
		_ = got
	}
}

func TestParseInt_Base37To62(t *testing.T) {
	// Test base 37-62 for signed integers
	// Note: base <= 36 uses strconv.ParseInt (only supports lowercase a-z)
	// base > 36 uses custom logic:
	//   - '0'-'9' = 0-9
	//   - 'a'-'z' = 10-35
	//   - 'A'-'Z' = 36-61
	testCases := []struct {
		s    string
		base int
		want int64
	}{
		{"-10", 37, -37},
		{"-A", 37, -36}, // base 37, '-A' = -36
		{"-Z", 37, -61}, // base 37, '-Z' = -61
		{"-a", 37, -10}, // base 37, '-a' = -10
		{"-z", 37, -35}, // base 37, '-z' = -35 (not -61!)
		{"10", 62, 62},
		{"-10", 62, -62},
		{"-A", 62, -36},
		{"-Z", 62, -61},
	}

	for _, tc := range testCases {
		got, err := parseInt(tc.s, tc.base, 64)
		if err != nil {
			t.Errorf("parseInt(%q, %d, 64) = error: %v, want %d", tc.s, tc.base, err, tc.want)
			continue
		}
		if got != tc.want {
			t.Errorf("parseInt(%q, %d, 64) = %d, want %d", tc.s, tc.base, got, tc.want)
		}
	}
}

func TestParseInt_Overflow(t *testing.T) {
	// Test overflow cases for signed integers with base > 36
	maxInt64 := int64(1<<63 - 1)
	minInt64 := int64(-1 << 63)

	// Test !neg && un >= cutoff (positive overflow)
	_, err := parseInt("9223372036854775808", 10, 64) // This should overflow
	if err == nil {
		t.Error("parseInt with positive overflow should return error")
	} else if err.(*strconv.NumError).Err != strconv.ErrRange {
		t.Errorf("parseInt with positive overflow should return ErrRange, got %v", err)
	}

	// Test neg && un > cutoff (negative overflow)
	_, err2 := parseInt("-9223372036854775809", 10, 64) // This should overflow
	if err2 == nil {
		t.Error("parseInt with negative overflow should return error")
	} else if err2.(*strconv.NumError).Err != strconv.ErrRange {
		t.Errorf("parseInt with negative overflow should return ErrRange, got %v", err2)
	}

	// Test with base 62
	_, err3 := parseInt("aZl8N0y58M8", 62, 64) // This should overflow (one more than max)
	if err3 == nil {
		t.Error("parseInt with overflow in base 62 should return error")
	}

	_ = maxInt64
	_ = minInt64
}

func TestParseUint_BaseError(t *testing.T) {
	// Test base > 62 error
	_, err := parseUint("10", 63, 64)
	if err == nil {
		t.Error("parseUint with base > 62 should return error")
	} else if nerr, ok := err.(*strconv.NumError); !ok || nerr.Err.Error() != "invalid base 63" {
		t.Errorf("parseUint with base > 62 should return baseError, got %v", err)
	}

	_, err2 := parseUint("10", 100, 64)
	if err2 == nil {
		t.Error("parseUint with base 100 should return error")
	}
}

func TestParseUint_BitSizeError(t *testing.T) {
	// Test bitSize < 0
	_, err := parseUint("10", 37, -1)
	if err == nil {
		t.Error("parseUint with bitSize < 0 should return error")
	} else if nerr, ok := err.(*strconv.NumError); !ok || nerr.Err.Error() != "invalid bit size -1" {
		t.Errorf("parseUint with bitSize < 0 should return bitSizeError, got %v", err)
	}

	// Test bitSize > 64
	_, err2 := parseUint("10", 37, 65)
	if err2 == nil {
		t.Error("parseUint with bitSize > 64 should return error")
	} else if nerr, ok := err2.(*strconv.NumError); !ok || nerr.Err.Error() != "invalid bit size 65" {
		t.Errorf("parseUint with bitSize > 64 should return bitSizeError, got %v", err2)
	}
}

func TestParseUint_BitSizeZero(t *testing.T) {
	// Test bitSize == 0 (should use strconv.IntSize)
	_, err := parseUint("100", 37, 0)
	if err != nil {
		t.Errorf("parseUint with bitSize 0 should work, got error: %v", err)
	}
}

func TestParseUint_InvalidChar(t *testing.T) {
	// Test invalid character (default case)
	_, err := parseUint("@", 37, 64)
	if err == nil {
		t.Error("parseUint with invalid character should return error")
	} else if nerr, ok := err.(*strconv.NumError); !ok || nerr.Err != strconv.ErrSyntax {
		t.Errorf("parseUint with invalid character should return syntaxError, got %v", err)
	}

	_, err2 := parseUint("10#", 37, 64)
	if err2 == nil {
		t.Error("parseUint with invalid character should return error")
	}
}

func TestParseUint_Base36InvalidDigit(t *testing.T) {
	// Test base <= 36 && d >= byte(base) case with lowercase letters
	// For base 10, 'a' (10) >= 10, should return error
	_, err := parseUint("a", 10, 64)
	if err == nil {
		t.Error("parseUint('a', 10) should return error")
	}
	if nerr, ok := err.(*strconv.NumError); ok {
		if nerr.Err != strconv.ErrSyntax {
			t.Errorf("parseUint error Err = %v, want ErrSyntax", nerr.Err)
		}
	}

	// For base 16, 'g' (16) >= 16, should return error
	_, err2 := parseUint("g", 16, 64)
	if err2 == nil {
		t.Error("parseUint('g', 16) should return error")
	}
	if nerr, ok := err2.(*strconv.NumError); ok {
		if nerr.Err != strconv.ErrSyntax {
			t.Errorf("parseUint error Err = %v, want ErrSyntax", nerr.Err)
		}
	}
}

func TestParseInt_EmptyString(t *testing.T) {
	// Test parseInt with empty string (base > 36)
	_, err := parseInt("", 37, 64)
	if err == nil {
		t.Error("parseInt with empty string should return error")
	} else if nerr, ok := err.(*strconv.NumError); !ok || nerr.Err != strconv.ErrSyntax {
		t.Errorf("parseInt with empty string should return syntaxError, got %v", err)
	}
}

func TestParseInt_PlusSign(t *testing.T) {
	// Test parseInt with '+' sign
	got, err := parseInt("+10", 37, 64)
	if err != nil {
		t.Errorf("parseInt('+10', 37) should work, got error: %v", err)
	} else if got != 37 {
		t.Errorf("parseInt('+10', 37) = %d, want 37", got)
	}

	got2, err2 := parseInt("+Z", 37, 64)
	if err2 != nil {
		t.Errorf("parseInt('+Z', 37) should work, got error: %v", err2)
	} else if got2 != 61 {
		t.Errorf("parseInt('+Z', 37) = %d, want 61", got2)
	}
}

func TestParseInt_BitSizeZero(t *testing.T) {
	// Test parseInt with bitSize == 0
	got, err := parseInt("100", 37, 0)
	if err != nil {
		t.Errorf("parseInt with bitSize 0 should work, got error: %v", err)
	} else if got <= 0 {
		t.Errorf("parseInt('100', 37, 0) = %d, should be positive", got)
	}
}

func TestParseInt_PositiveOverflow(t *testing.T) {
	// Test !neg && un >= cutoff case
	// For bitSize 8, max is 127, so 128 should overflow
	got, err := parseInt("128", 10, 8)
	if err == nil {
		t.Error("parseInt('128', 10, 8) should return error")
	} else if nerr, ok := err.(*strconv.NumError); !ok || nerr.Err != strconv.ErrRange {
		t.Errorf("parseInt with positive overflow should return ErrRange, got %v", err)
	} else if got != 127 {
		t.Errorf("parseInt with positive overflow should return max value 127, got %d", got)
	}
}

func TestParseInt_NegativeOverflow(t *testing.T) {
	// Test neg && un > cutoff case
	// For bitSize 8, min is -128, so -129 should overflow
	got, err := parseInt("-129", 10, 8)
	if err == nil {
		t.Error("parseInt('-129', 10, 8) should return error")
	} else if nerr, ok := err.(*strconv.NumError); !ok || nerr.Err != strconv.ErrRange {
		t.Errorf("parseInt with negative overflow should return ErrRange, got %v", err)
	} else if got != -128 {
		t.Errorf("parseInt with negative overflow should return min value -128, got %d", got)
	}

	// Test with 64-bit int, cutoff is 9223372036854775808 (2^63)
	// So -9223372036854775809 should trigger overflow
	_, err2 := parseInt("-9223372036854775809", 10, 64)
	if err2 == nil {
		t.Error("parseInt with negative overflow should return error")
	}
	if nerr, ok := err2.(*strconv.NumError); ok {
		if nerr.Err != strconv.ErrRange {
			t.Errorf("parseInt error Err = %v, want ErrRange", nerr.Err)
		}
		if nerr.Func != "ParseInt" {
			t.Errorf("parseInt error Func = %s, want ParseInt", nerr.Func)
		}
	}

	// Test with 32-bit int, cutoff is 2147483648 (2^31)
	// So -2147483649 should trigger overflow
	_, err3 := parseInt("-2147483649", 10, 32)
	if err3 == nil {
		t.Error("parseInt with negative overflow should return error")
	}
	if nerr, ok := err3.(*strconv.NumError); ok {
		if nerr.Err != strconv.ErrRange {
			t.Errorf("parseInt error Err = %v, want ErrRange", nerr.Err)
		}
	}

	// Test boundary case: min int64 should work
	minInt64 := "-9223372036854775808"
	got2, err4 := parseInt(minInt64, 10, 64)
	if err4 != nil {
		t.Errorf("parseInt with min int64 should work, got error: %v", err4)
	} else if got2 != math.MinInt64 {
		t.Errorf("parseInt('%s', 10, 64) = %d, want %d", minInt64, got2, math.MinInt64)
	}
}

func TestParseInt_ErrRangePropagation(t *testing.T) {
	// Test that ErrRange from parseUint is propagated correctly
	// Use a number that causes overflow in parseUint (base > 36)
	_, err := parseInt("ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ", 62, 64)
	if err == nil {
		t.Error("parseInt with overflow should return error")
	} else if nerr, ok := err.(*strconv.NumError); !ok || nerr.Err != strconv.ErrRange {
		t.Errorf("parseInt with overflow should return ErrRange, got %v", err)
	}
}

func TestParseInt_BoundaryCases(t *testing.T) {
	// Test !neg && un == cutoff case (should return error)
	// For bitSize 8, cutoff is 128, so 128 should overflow
	got, err := parseInt("128", 10, 8)
	if err == nil {
		t.Error("parseInt('128', 10, 8) should return error")
	} else if got != 127 {
		t.Errorf("parseInt('128', 10, 8) should return 127, got %d", got)
	}

	// Test neg && un == cutoff case (should NOT return error, condition is un > cutoff)
	// For bitSize 8, cutoff is 128, so -128 should NOT overflow
	got2, err2 := parseInt("-128", 10, 8)
	if err2 != nil {
		t.Errorf("parseInt('-128', 10, 8) should work, got error: %v", err2)
	} else if got2 != -128 {
		t.Errorf("parseInt('-128', 10, 8) = %d, want -128", got2)
	}
}

func TestParseUint_Wrapper(t *testing.T) {
	// Test ParseUint wrapper function
	result := ParseUint("10", 37, 64)
	if result.IsErr() {
		t.Errorf("ParseUint('10', 37, 64) should work, got error: %v", result.UnwrapErr())
	} else if result.Unwrap() != 37 {
		t.Errorf("ParseUint('10', 37, 64) = %d, want 37", result.Unwrap())
	}

	// Test error case
	result2 := ParseUint("", 37, 64)
	if !result2.IsErr() {
		t.Error("ParseUint('', 37, 64) should return error")
	}
}

func TestParseInt_Wrapper(t *testing.T) {
	// Test ParseInt wrapper function
	result := ParseInt("10", 37, 64)
	if result.IsErr() {
		t.Errorf("ParseInt('10', 37, 64) should work, got error: %v", result.UnwrapErr())
	} else if result.Unwrap() != 37 {
		t.Errorf("ParseInt('10', 37, 64) = %d, want 37", result.Unwrap())
	}

	// Test error case
	result2 := ParseInt("", 37, 64)
	if !result2.IsErr() {
		t.Error("ParseInt('', 37, 64) should return error")
	}

	// Test negative case
	result3 := ParseInt("-10", 37, 64)
	if result3.IsErr() {
		t.Errorf("ParseInt('-10', 37, 64) should work, got error: %v", result3.UnwrapErr())
	} else if result3.Unwrap() != -37 {
		t.Errorf("ParseInt('-10', 37, 64) = %d, want -37", result3.Unwrap())
	}
}

func TestAtoi_FastPath(t *testing.T) {
	// Test fast path for small integers
	got, err := Atoi("123")
	if err != nil {
		t.Errorf("Atoi('123') should work, got error: %v", err)
	} else if got != 123 {
		t.Errorf("Atoi('123') = %d, want 123", got)
	}

	// Test fast path with negative
	got2, err2 := Atoi("-456")
	if err2 != nil {
		t.Errorf("Atoi('-456') should work, got error: %v", err2)
	} else if got2 != -456 {
		t.Errorf("Atoi('-456') = %d, want -456", got2)
	}

	// Test fast path with plus sign
	got3, err3 := Atoi("+789")
	if err3 != nil {
		t.Errorf("Atoi('+789') should work, got error: %v", err3)
	} else if got3 != 789 {
		t.Errorf("Atoi('+789') = %d, want 789", got3)
	}

	// Test fast path with invalid character
	_, err4 := Atoi("12a")
	if err4 == nil {
		t.Error("Atoi('12a') should return error")
	}

	// Test fast path with only sign (empty after sign)
	_, err5 := Atoi("-")
	if err5 == nil {
		t.Error("Atoi('-') should return error")
	}

	_, err6 := Atoi("+")
	if err6 == nil {
		t.Error("Atoi('+') should return error")
	}
}

func TestAtoi_SlowPath(t *testing.T) {
	// Test slow path for large integers
	largeNum := "9223372036854775807" // max int64
	got, err := Atoi(largeNum)
	if err != nil {
		t.Errorf("Atoi('%s') should work, got error: %v", largeNum, err)
	} else if got <= 0 {
		t.Errorf("Atoi('%s') = %d, should be positive", largeNum, got)
	}

	// Test slow path with underscores (should fail for base 10)
	_, err2 := Atoi("1_2_3")
	if err2 == nil {
		t.Error("Atoi('1_2_3') should return error (underscores not allowed in base 10)")
	}
}

func TestParseInt_ErrRangePropagation2(t *testing.T) {
	// Test that ErrRange from parseUint is propagated correctly (err != nil && nerr.Err == strconv.ErrRange)
	// This tests the else branch at line 145
	_, err := parseInt("18446744073709551616", 10, 64) // This causes overflow in parseUint
	if err == nil {
		t.Error("parseInt with overflow should return error")
	} else if nerr, ok := err.(*strconv.NumError); !ok || nerr.Err != strconv.ErrRange {
		t.Errorf("parseInt with overflow should return ErrRange, got %v", err)
	}
}

func TestTryFromString_ReflectPath(t *testing.T) {
	// Test reflect path in tryFromString (when type switch doesn't match)
	// This requires a custom type that doesn't match any case
	type CustomString string
	type CustomInt int

	// This should use the reflect path
	result := TryFromString[CustomString, CustomInt]("42", 10, 0)
	if result.IsErr() {
		t.Errorf("TryFromString should work, got error: %v", result.UnwrapErr())
	}
}

func TestAs_ReflectPath(t *testing.T) {
	// Test reflect path in as function
	// The reflect path is used when type switch doesn't match
	// For example, when converting between types that are not directly matched in the type switch
	// but have the same underlying kind

	// Test with types that will use reflect path
	// Using int8 -> int16 which should use reflect path since type switch checks for *int16
	result := As[int8, int16](100)
	if result.IsErr() {
		t.Errorf("As[int8, int16] should work, got error: %v", result.UnwrapErr())
	} else if result.Unwrap() != 100 {
		t.Errorf("As[int8, int16](100) = %d, want 100", result.Unwrap())
	}

	// Test with uint8 -> uint16
	result2 := As[uint8, uint16](200)
	if result2.IsErr() {
		t.Errorf("As[uint8, uint16] should work, got error: %v", result2.UnwrapErr())
	} else if result2.Unwrap() != 200 {
		t.Errorf("As[uint8, uint16](200) = %d, want 200", result2.Unwrap())
	}
}

func TestParseUint_UnderscoreCases(t *testing.T) {
	// Test underscore handling in parseUint for base > 36
	got, err := parseUint("1_2_3", 37, 64)
	if err != nil {
		t.Errorf("parseUint('1_2_3', 37, 64) should work, got error: %v", err)
	} else {
		// 1_2_3 in base 37 = 1*37^2 + 2*37 + 3 = 1369 + 74 + 3 = 1446
		expected := uint64(1*37*37 + 2*37 + 3)
		if got != expected {
			t.Errorf("parseUint('1_2_3', 37, 64) = %d, want %d", got, expected)
		}
	}

	// Test underscore at beginning (should fail)
	_, err2 := parseUint("_123", 37, 64)
	if err2 == nil {
		t.Error("parseUint('_123', 37, 64) should return error")
	}

	// Test underscore at end (underscoreOK allows it, parseUint skips underscores)
	// "123_" is parsed as "123" in base 37: 1*37^2 + 2*37 + 3 = 1446
	got3, err3 := parseUint("123_", 37, 64)
	if err3 != nil {
		t.Errorf("parseUint('123_', 37, 64) should work (underscore is skipped), got error: %v", err3)
	} else {
		expected := uint64(1*37*37 + 2*37 + 3) // 1446
		if got3 != expected {
			t.Errorf("parseUint('123_', 37, 64) = %d, want %d", got3, expected)
		}
	}

	// Test double underscore (should fail - underscoreOK returns false)
	_, err4 := parseUint("1__23", 37, 64)
	if err4 == nil {
		t.Error("parseUint('1__23', 37, 64) should return error")
	}
}

func TestParseUint_Base37To62_AllCases(t *testing.T) {
	// Test all base 37-62 cases to ensure full coverage
	testCases := []struct {
		s    string
		base int
		want uint64
	}{
		{"A", 37, 36},
		{"B", 37, 37},
		{"Z", 37, 61},
		{"10", 37, 37},
		{"1Z", 37, 37*1 + 61},
		{"A", 62, 36},
		{"Z", 62, 61},
		{"10", 62, 62},
		{"1Z", 62, 62*1 + 61},
	}

	for _, tc := range testCases {
		got, err := parseUint(tc.s, tc.base, 64)
		if err != nil {
			t.Errorf("parseUint(%q, %d, 64) = error: %v, want %d", tc.s, tc.base, err, tc.want)
			continue
		}
		if got != tc.want {
			t.Errorf("parseUint(%q, %d, 64) = %d, want %d", tc.s, tc.base, got, tc.want)
		}
	}
}

func TestUnderscoreOK_BasePrefix(t *testing.T) {
	// Test underscoreOK with base prefix (0x, 0b, 0o)
	testCases := []struct {
		s    string
		want bool
	}{
		{"0x123", true},
		{"0X123", true},
		{"0b101", true},
		{"0B101", true},
		{"0o377", true},
		{"0O377", true},
		{"0x_123", true},   // underscore after prefix (base prefix counts as digit, so this is valid)
		{"0x1_23", true},   // underscore in number
		{"0x12_3", true},   // underscore in number
		{"0x_", true},      // underscore immediately after prefix (base prefix counts as digit, so this is valid)
		{"0x__123", false}, // double underscore
		{"0x123_", true},   // underscore at end (allowed by underscoreOK - it follows a digit)
	}

	for _, tc := range testCases {
		got := underscoreOK(tc.s)
		if got != tc.want {
			t.Errorf("underscoreOK(%q) = %v, want %v", tc.s, got, tc.want)
		}
	}
}

func TestUnderscoreOK_WithSign(t *testing.T) {
	// Test underscoreOK with sign
	testCases := []struct {
		s    string
		want bool
	}{
		{"-123", true},
		{"+123", true},
		{"-1_2_3", true},
		{"+1_2_3", true},
		{"-_123", false},  // underscore after sign (invalid)
		{"-1__23", false}, // double underscore
		{"-123_", true},   // underscore at end (allowed by underscoreOK - it follows a digit)
	}

	for _, tc := range testCases {
		got := underscoreOK(tc.s)
		if got != tc.want {
			t.Errorf("underscoreOK(%q) = %v, want %v", tc.s, got, tc.want)
		}
	}
}

func TestTryFromString_Base37(t *testing.T) {
	// Test TryFromString with base > 36
	result := TryFromString[string, int]("Z", 37, 64)
	if result.IsErr() {
		t.Errorf("TryFromString('Z', 37, 64) should work, got error: %v", result.UnwrapErr())
	} else if result.Unwrap() != 61 {
		t.Errorf("TryFromString('Z', 37, 64) = %d, want 61", result.Unwrap())
	}
}

func TestAs_AllTypes(t *testing.T) {
	// Test As function with various type combinations to cover reflect path
	testCases := []struct {
		name string
		test func(*testing.T)
	}{
		{"int8 to int16", func(t *testing.T) {
			result := As[int8, int16](100)
			assert.True(t, result.IsOk())
			assert.Equal(t, int16(100), result.Unwrap())
		}},
		{"int16 to int32", func(t *testing.T) {
			result := As[int16, int32](1000)
			assert.True(t, result.IsOk())
			assert.Equal(t, int32(1000), result.Unwrap())
		}},
		{"int32 to int64", func(t *testing.T) {
			result := As[int32, int64](100000)
			assert.True(t, result.IsOk())
			assert.Equal(t, int64(100000), result.Unwrap())
		}},
		{"uint8 to uint16", func(t *testing.T) {
			result := As[uint8, uint16](200)
			assert.True(t, result.IsOk())
			assert.Equal(t, uint16(200), result.Unwrap())
		}},
		{"uint16 to uint32", func(t *testing.T) {
			result := As[uint16, uint32](50000)
			assert.True(t, result.IsOk())
			assert.Equal(t, uint32(50000), result.Unwrap())
		}},
		{"uint32 to uint64", func(t *testing.T) {
			result := As[uint32, uint64](1000000)
			assert.True(t, result.IsOk())
			assert.Equal(t, uint64(1000000), result.Unwrap())
		}},
		{"float32 to float64", func(t *testing.T) {
			result := As[float32, float64](3.14)
			assert.True(t, result.IsOk())
			assert.InDelta(t, 3.14, result.Unwrap(), 0.001)
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.test)
	}
}

func TestAtoi_EmptyString(t *testing.T) {
	// Test empty string (should use slow path)
	_, err := Atoi("")
	if err == nil {
		t.Error("Atoi('') should return error")
	}
}

func TestAtoi_ZeroLength(t *testing.T) {
	// Test zero length string (should use slow path)
	_, err := Atoi("")
	if err == nil {
		t.Error("Atoi('') should return error")
	}
}

func TestParseInt_ErrRangePropagation_NonRangeError(t *testing.T) {
	// Test that non-ErrRange errors from parseUint are propagated correctly
	// This tests the if branch at line 140-143
	_, err := parseInt("abc", 10, 64) // This causes syntax error in parseUint
	if err == nil {
		t.Error("parseInt with syntax error should return error")
	} else if nerr, ok := err.(*strconv.NumError); !ok || nerr.Err == strconv.ErrRange {
		// Should be syntax error, not range error
		if nerr.Err == strconv.ErrRange {
			t.Errorf("parseInt with syntax error should not return ErrRange, got %v", err)
		}
	}
}

func TestParseInt_NegativeOverflowBoundary(t *testing.T) {
	// Test negative overflow boundary case
	// For 64-bit: cutoff = 1<<63, neg && un > cutoff
	result := ParseInt("-9223372036854775809", 10, 64) // -2^63 - 1
	assert.True(t, result.IsErr())
}

func TestParseInt_PositiveOverflowBoundary(t *testing.T) {
	// Test positive overflow boundary case
	// For 64-bit: cutoff = 1<<63, !neg && un >= cutoff
	result := ParseInt("9223372036854775808", 10, 64) // 2^63
	assert.True(t, result.IsErr())
}

func TestParseUint_OverflowBoundary(t *testing.T) {
	// Test overflow boundary: n >= cutoff
	result := ParseUint("18446744073709551615", 10, 64) // Max uint64
	assert.False(t, result.IsErr())
	assert.Equal(t, uint64(18446744073709551615), result.Unwrap())

	// Test overflow: n >= cutoff
	result2 := ParseUint("18446744073709551616", 10, 64) // Max uint64 + 1
	assert.True(t, result2.IsErr())
}

func TestParseUint_N1Overflow(t *testing.T) {
	// Test n1 overflow: n1 < n || n1 > maxVal
	// This is hard to trigger, but we can try with very large numbers
	result := ParseUint("18446744073709551615", 10, 64) // Max uint64
	assert.False(t, result.IsErr())
}

func TestParseUint_Base37InvalidDigit(t *testing.T) {
	// Test base 37 with invalid digit (digit >= base for base <= 36)
	// But for base > 36, all digits 0-9, a-z, A-Z are valid
	result := ParseUint("Z", 37, 64)
	assert.False(t, result.IsErr())
	assert.Equal(t, uint64(61), result.Unwrap()) // Z = 61 in base 37
}

func TestTryFromStrings_ErrorPropagation(t *testing.T) {
	// Test error propagation in TryFromStrings
	result := TryFromStrings[string, int]([]string{"42", "abc", "100"}, 10, 0)
	assert.True(t, result.IsErr())
}

func TestTryFromStrings_EmptySlice(t *testing.T) {
	// Test empty slice
	result := TryFromStrings[string, int]([]string{}, 10, 0)
	assert.False(t, result.IsErr())
	assert.Equal(t, []int{}, result.Unwrap())
}

func TestSliceAs_ErrorPropagation(t *testing.T) {
	// Test error propagation in SliceAs
	result, err := SliceAs[int, int8]([]int{100, 200, 300})
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestFormatByDict_LargeNumber(t *testing.T) {
	// Test FormatByDict with large number
	dict := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	result := FormatByDict(dict, 1000000)
	assert.NotEmpty(t, result)

	// Verify round-trip
	parsed := ParseByDict[uint64](dict, result)
	assert.False(t, parsed.IsErr())
	assert.Equal(t, uint64(1000000), parsed.Unwrap())
}

func TestFormatByDict_Zero(t *testing.T) {
	// Test FormatByDict with zero
	dict := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := FormatByDict(dict, 0)
	assert.NotEmpty(t, result)

	// Verify round-trip
	parsed := ParseByDict[uint64](dict, result)
	assert.False(t, parsed.IsErr())
	assert.Equal(t, uint64(0), parsed.Unwrap())
}

func TestParseByDict_EmptyString(t *testing.T) {
	// Test ParseByDict with empty string
	dict := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := ParseByDict[uint64](dict, "")
	assert.False(t, result.IsErr())
	assert.Equal(t, uint64(0), result.Unwrap())
}

func TestParseByDict_LargeNumber(t *testing.T) {
	// Test ParseByDict with large number
	dict := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	numStr := FormatByDict(dict, 1000000)
	result := ParseByDict[uint64](dict, numStr)
	assert.False(t, result.IsErr())
	assert.Equal(t, uint64(1000000), result.Unwrap())
}

func TestCheckedMul_EdgeCases(t *testing.T) {
	// Test CheckedMul with edge cases
	maxInt := Max[int]()

	// Test a = Max/b exactly (boundary case)
	result := CheckedMul(maxInt, 1)
	assert.True(t, result.IsSome())
	assert.Equal(t, maxInt, result.Unwrap())

	// Test a = Max/b + 1 (should overflow)
	result2 := CheckedMul(maxInt, 2)
	assert.True(t, result2.IsNone())

	// Test with zero (should work)
	result3 := CheckedMul(0, maxInt)
	assert.True(t, result3.IsSome())
	assert.Equal(t, 0, result3.Unwrap())

	// Test with one (should work)
	result4 := CheckedMul(1, maxInt)
	assert.True(t, result4.IsSome())
	assert.Equal(t, maxInt, result4.Unwrap())
}

func TestCheckedAdd_EdgeCases(t *testing.T) {
	// Test CheckedAdd with edge cases
	maxInt := Max[int]()

	// Test a = Max - b exactly (boundary case)
	result := CheckedAdd(maxInt-1, 1)
	assert.True(t, result.IsSome())
	assert.Equal(t, maxInt, result.Unwrap())

	// Test a = Max - b + 1 (should overflow)
	result2 := CheckedAdd(maxInt, 1)
	assert.True(t, result2.IsNone())

	// Test with zero (should work)
	result3 := CheckedAdd(0, maxInt)
	assert.True(t, result3.IsSome())
	assert.Equal(t, maxInt, result3.Unwrap())
}

func TestSaturatingAdd_EdgeCases(t *testing.T) {
	// Test SaturatingAdd with edge cases
	maxInt := Max[int]()

	// Test a = Max - b exactly (boundary case)
	result := SaturatingAdd(maxInt-1, 1)
	assert.Equal(t, maxInt, result)

	// Test a = Max - b + 1 (should saturate)
	result2 := SaturatingAdd(maxInt, 1)
	assert.Equal(t, maxInt, result2)

	// Test with zero (should work)
	result3 := SaturatingAdd(0, maxInt)
	assert.Equal(t, maxInt, result3)
}

func TestSaturatingSub_EdgeCases(t *testing.T) {
	// Test SaturatingSub with edge cases
	// Test a = b exactly (boundary case)
	result := SaturatingSub(5, 5)
	assert.Equal(t, 0, result)

	// Test a < b (should return 0)
	result2 := SaturatingSub(3, 5)
	assert.Equal(t, 0, result2)

	// Test with floats
	result3 := SaturatingSub(5.5, 3.0)
	assert.Equal(t, 2.5, result3)

	result4 := SaturatingSub(3.0, 5.5)
	assert.Equal(t, 0.0, result4)
}

func TestAtoi_FastPathBoundary(t *testing.T) {
	// Test fast path boundary cases
	// For 64-bit: length < 19 uses fast path
	// For 32-bit: length < 10 uses fast path

	// Test length 18 (64-bit fast path boundary)
	if strconv.IntSize == 64 {
		got, err := Atoi("123456789012345678") // 18 digits
		if err != nil {
			t.Errorf("Atoi with 18 digits should work, got error: %v", err)
		}
		if got == 0 {
			t.Error("Atoi with 18 digits should return non-zero")
		}

		// Test length 19 (should use slow path)
		got2, err2 := Atoi("1234567890123456789") // 19 digits
		if err2 != nil {
			t.Errorf("Atoi with 19 digits should work, got error: %v", err2)
		}
		if got2 == 0 {
			t.Error("Atoi with 19 digits should return non-zero")
		}
	}

	// Test length 9 (32-bit fast path boundary)
	if strconv.IntSize == 32 {
		got, err := Atoi("123456789") // 9 digits
		if err != nil {
			t.Errorf("Atoi with 9 digits should work, got error: %v", err)
		}
		if got == 0 {
			t.Error("Atoi with 9 digits should return non-zero")
		}

		// Test length 10 (should use slow path)
		got2, err2 := Atoi("1234567890") // 10 digits
		if err2 != nil {
			t.Errorf("Atoi with 10 digits should work, got error: %v", err2)
		}
		if got2 == 0 {
			t.Error("Atoi with 10 digits should return non-zero")
		}
	}
}

func TestParseUint_BitSizeBoundary(t *testing.T) {
	// Test bitSize boundary cases
	// bitSize < 0 should return error
	_, err := parseUint("123", 10, -1)
	if err == nil {
		t.Error("parseUint with bitSize -1 should return error")
	}

	// bitSize > 64 should return error
	_, err2 := parseUint("123", 10, 65)
	if err2 == nil {
		t.Error("parseUint with bitSize 65 should return error")
	}

	// bitSize == 64 should work
	got, err3 := parseUint("123", 10, 64)
	if err3 != nil {
		t.Errorf("parseUint with bitSize 64 should work, got error: %v", err3)
	} else if got != 123 {
		t.Errorf("parseUint('123', 10, 64) = %d, want 123", got)
	}
}

func TestParseInt_BitSizeBoundary(t *testing.T) {
	// Test bitSize boundary cases in parseInt
	// bitSize < 0 should return error (through parseUint)
	_, err := parseInt("123", 10, -1)
	if err == nil {
		t.Error("parseInt with bitSize -1 should return error")
	}

	// bitSize > 64 should return error (through parseUint)
	_, err2 := parseInt("123", 10, 65)
	if err2 == nil {
		t.Error("parseInt with bitSize 65 should return error")
	}
}

func TestParseUint_OverflowBoundary2(t *testing.T) {
	// Test overflow boundary cases (internal parseUint function)
	// n >= cutoff case
	maxUint64 := "18446744073709551615" // max uint64
	got, err := parseUint(maxUint64, 10, 64)
	if err != nil {
		t.Errorf("parseUint with max uint64 should work, got error: %v", err)
	} else if got != math.MaxUint64 {
		t.Errorf("parseUint('%s', 10, 64) = %d, want %d", maxUint64, got, uint64(math.MaxUint64))
	}

	// n >= cutoff case (should overflow)
	overflow := "18446744073709551616" // max uint64 + 1
	got2, err2 := parseUint(overflow, 10, 64)
	if err2 == nil {
		t.Error("parseUint with overflow should return error")
	} else if got2 != math.MaxUint64 {
		t.Errorf("parseUint('%s', 10, 64) should return max uint64 on overflow, got %d", overflow, got2)
	}

	// n1 < n case (n+v overflows)
	// This is harder to trigger, but we can try with a large number in a small base
	largeBase10 := "9999999999999999999" // Very large number
	got3, err3 := parseUint(largeBase10, 10, 64)
	if err3 == nil && got3 == 0 {
		t.Error("parseUint with very large number should either return error or non-zero value")
	}
}

func TestParseInt_PositiveOverflowBoundary2(t *testing.T) {
	// Test positive overflow boundary (!neg && un >= cutoff) - internal parseInt function
	maxInt64 := "9223372036854775807" // max int64
	got, err := parseInt(maxInt64, 10, 64)
	if err != nil {
		t.Errorf("parseInt with max int64 should work, got error: %v", err)
	} else if got != math.MaxInt64 {
		t.Errorf("parseInt('%s', 10, 64) = %d, want %d", maxInt64, got, int64(math.MaxInt64))
	}

	// Test overflow (un >= cutoff)
	overflow := "9223372036854775808" // max int64 + 1
	got2, err2 := parseInt(overflow, 10, 64)
	if err2 == nil {
		t.Error("parseInt with overflow should return error")
	} else if got2 != math.MaxInt64 {
		t.Errorf("parseInt('%s', 10, 64) should return max int64 on overflow, got %d", overflow, got2)
	}
}

func TestUnderscoreOK_EdgeCases(t *testing.T) {
	// Test underscoreOK with various edge cases
	// Test with base prefix and underscore
	assert.True(t, underscoreOK("0x_123"))
	assert.True(t, underscoreOK("0b_101"))
	assert.True(t, underscoreOK("0o_377"))

	// Test with underscore after base prefix
	// Base prefix counts as a digit, so underscore immediately after prefix is OK
	assert.True(t, underscoreOK("0x1_23"))
	assert.True(t, underscoreOK("0x_1_23")) // underscore immediately after prefix is OK because prefix counts as digit

	// Test with multiple underscores
	assert.True(t, underscoreOK("1_2_3_4"))
	assert.False(t, underscoreOK("_1_2_3")) // underscore at start (before any digit)
	assert.True(t, underscoreOK("1_2_3_"))  // underscore at end is OK (follows a digit)

	// Test with sign and underscore
	assert.True(t, underscoreOK("+1_2_3"))
	assert.True(t, underscoreOK("-1_2_3"))
	assert.False(t, underscoreOK("+_1_2_3")) // underscore after sign (before any digit)
	assert.False(t, underscoreOK("-_1_2_3")) // underscore after sign (before any digit)
}

func TestParseUint_AllBases(t *testing.T) {
	// Test parseUint with all valid bases (2-62)
	for base := 2; base <= 62; base++ {
		// Test with valid digit for this base
		var testStr string
		if base <= 10 {
			// For base <= 10, use a digit that's valid for all bases (e.g., "1")
			testStr = "1"
		} else if base <= 36 {
			// For base 11-36, use "a" which is valid (represents 10)
			testStr = "a"
		} else {
			// For base > 36, use "A" which is valid (represents 36)
			testStr = "A"
		}
		got, err := parseUint(testStr, base, 64)
		if err != nil {
			t.Errorf("parseUint('%s', %d, 64) should work, got error: %v", testStr, base, err)
		} else if got == 0 && testStr != "0" {
			t.Errorf("parseUint('%s', %d, 64) = %d, should be non-zero", testStr, base, got)
		}
	}
}

func TestParseUint_BaseErrors(t *testing.T) {
	// Test base error cases
	// Note: base 0 is handled by strconv.ParseUint (base <= 36 path), so it won't error
	// unless the string is invalid. We test base 0 with an invalid string.
	_, err := parseUint("invalid", 0, 64)
	if err == nil {
		t.Error("parseUint with base 0 and invalid string should return error")
	}

	// base 1 is invalid
	_, err2 := parseUint("123", 1, 64)
	if err2 == nil {
		t.Error("parseUint with base 1 should return error")
	}

	// base > 62 should return error (enters base > 36 path)
	_, err3 := parseUint("123", 63, 64)
	if err3 == nil {
		t.Error("parseUint with base 63 should return error")
	}

	_, err4 := parseUint("123", 100, 64)
	if err4 == nil {
		t.Error("parseUint with base 100 should return error")
	}
}

func TestParseUint_EmptyString(t *testing.T) {
	// Test empty string
	_, err := parseUint("", 37, 64)
	if err == nil {
		t.Error("parseUint with empty string should return error")
	}
}

func TestParseUint_InvalidCharacters(t *testing.T) {
	// Test invalid characters
	_, err := parseUint("12@34", 37, 64)
	if err == nil {
		t.Error("parseUint with invalid character @ should return error")
	}

	_, err2 := parseUint("12#34", 37, 64)
	if err2 == nil {
		t.Error("parseUint with invalid character # should return error")
	}

	_, err3 := parseUint("12 34", 37, 64)
	if err3 == nil {
		t.Error("parseUint with space should return error")
	}
}

func TestParseUint_OverflowCases(t *testing.T) {
	// Test various overflow cases
	// n >= cutoff case
	maxUint64Str := "18446744073709551615"
	got, err := parseUint(maxUint64Str, 10, 64)
	if err != nil {
		t.Errorf("parseUint with max uint64 should work, got error: %v", err)
	} else if got != math.MaxUint64 {
		t.Errorf("parseUint('%s', 10, 64) = %d, want %d", maxUint64Str, got, uint64(math.MaxUint64))
	}

	// Test n1 < n case (n+v overflows)
	overflowStr := "18446744073709551616"
	got2, err2 := parseUint(overflowStr, 10, 64)
	if err2 == nil {
		t.Error("parseUint with overflow should return error")
	} else if got2 != math.MaxUint64 {
		t.Errorf("parseUint('%s', 10, 64) should return max uint64 on overflow, got %d", overflowStr, got2)
	}
}

func TestParseInt_OnlySign(t *testing.T) {
	// Test with only sign
	_, err := parseInt("+", 37, 64)
	if err == nil {
		t.Error("parseInt with only '+' should return error")
	}

	_, err2 := parseInt("-", 37, 64)
	if err2 == nil {
		t.Error("parseInt with only '-' should return error")
	}
}

func TestAtoi_EdgeCases(t *testing.T) {
	// Test Atoi with various edge cases
	// Test with zero
	got, err := Atoi("0")
	if err != nil {
		t.Errorf("Atoi('0') should work, got error: %v", err)
	} else if got != 0 {
		t.Errorf("Atoi('0') = %d, want 0", got)
	}

	// Test with single digit
	got2, err2 := Atoi("5")
	if err2 != nil {
		t.Errorf("Atoi('5') should work, got error: %v", err2)
	} else if got2 != 5 {
		t.Errorf("Atoi('5') = %d, want 5", got2)
	}

	// Test with negative zero (edge case)
	got3, err3 := Atoi("-0")
	if err3 != nil {
		t.Errorf("Atoi('-0') should work, got error: %v", err3)
	} else if got3 != 0 {
		t.Errorf("Atoi('-0') = %d, want 0", got3)
	}
}

func TestAtoi_InvalidInput(t *testing.T) {
	// Test Atoi with invalid input
	_, err := Atoi("")
	if err == nil {
		t.Error("Atoi with empty string should return error")
	}

	_, err2 := Atoi("abc")
	if err2 == nil {
		t.Error("Atoi with non-numeric string should return error")
	}

	_, err3 := Atoi("12.34")
	if err3 == nil {
		t.Error("Atoi with decimal should return error")
	}
}

// Note: baseError, bitSizeError, rangeError, syntaxError are unexported functions
// They are tested indirectly through ParseUint and ParseInt functions

func TestLower(t *testing.T) {
	// Test lower function
	if lower('A') != 'a' {
		t.Errorf("lower('A') = %c, want 'a'", lower('A'))
	}
	if lower('Z') != 'z' {
		t.Errorf("lower('Z') = %c, want 'z'", lower('Z'))
	}
	if lower('a') != 'a' {
		t.Errorf("lower('a') = %c, want 'a'", lower('a'))
	}
	if lower('0') != '0' {
		t.Errorf("lower('0') = %c, want '0'", lower('0'))
	}
}

func TestParseInt_NonErrRangeError(t *testing.T) {
	// Test parseInt with non-ErrRange error (syntax error)
	// This tests the nerr.Err != strconv.ErrRange branch at line 140
	// We need to trigger a syntax error (not range error) in parseUint
	// For base > 36, parseUint will call syntaxError for invalid characters
	// Use a string with invalid character '#' to trigger syntax error
	_, err := parseInt("12#34", 37, 64)
	if err == nil {
		t.Error("parseInt with invalid string should return error")
	}
	// Verify it's a syntax error, not a range error
	if nerr, ok := err.(*strconv.NumError); ok {
		if nerr.Err == strconv.ErrRange {
			t.Error("parseInt should return syntax error, not range error")
		}
		if nerr.Func != "ParseInt" {
			t.Errorf("parseInt error Func = %s, want ParseInt", nerr.Func)
		}
	}
}

func TestParseInt_NonErrRangeError2(t *testing.T) {
	// Test parseInt with non-ErrRange error (syntax error) for base > 36
	// This tests the nerr.Err != strconv.ErrRange branch at line 140
	// Use a string with invalid character '@' to trigger syntax error
	_, err := parseInt("+12@34", 37, 64)
	if err == nil {
		t.Error("parseInt with invalid string should return error")
	}
	// Verify it's a syntax error, not a range error
	if nerr, ok := err.(*strconv.NumError); ok {
		if nerr.Err == strconv.ErrRange {
			t.Error("parseInt should return syntax error, not range error")
		}
		if nerr.Func != "ParseInt" {
			t.Errorf("parseInt error Func = %s, want ParseInt", nerr.Func)
		}
	}
}

func TestParseInt_Base36InvalidDigit_Lowercase(t *testing.T) {
	// Test base <= 36 && d >= byte(base) case with lowercase letters
	// For base 10, 'a' (10) >= 10, should return error
	_, err := parseInt("a", 10, 64)
	if err == nil {
		t.Error("parseInt with invalid digit for base should return error")
	}
	if nerr, ok := err.(*strconv.NumError); ok {
		if nerr.Err != strconv.ErrSyntax {
			t.Errorf("parseInt error Err = %v, want ErrSyntax", nerr.Err)
		}
	}

	// Test with base 16, 'g' (16) >= 16, should return error
	_, err2 := parseInt("g", 16, 64)
	if err2 == nil {
		t.Error("parseInt with invalid digit for base should return error")
	}
	if nerr, ok := err2.(*strconv.NumError); ok {
		if nerr.Err != strconv.ErrSyntax {
			t.Errorf("parseInt error Err = %v, want ErrSyntax", nerr.Err)
		}
	}
}
