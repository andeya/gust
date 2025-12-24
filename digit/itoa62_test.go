// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package digit

import (
	"testing"
)

type itoa64Test struct {
	in   int64
	base int
	out  string
}

var itoa64tests = []itoa64Test{
	{0, 10, "0"},
	{1, 10, "1"},
	{-1, 10, "-1"},
	{12345678, 10, "12345678"},
	{-987654321, 10, "-987654321"},
	{1<<31 - 1, 10, "2147483647"},
	{-1<<31 + 1, 10, "-2147483647"},
	{1 << 31, 10, "2147483648"},
	{-1 << 31, 10, "-2147483648"},
	{1<<31 + 1, 10, "2147483649"},
	{-1<<31 - 1, 10, "-2147483649"},
	{1<<32 - 1, 10, "4294967295"},
	{-1<<32 + 1, 10, "-4294967295"},
	{1 << 32, 10, "4294967296"},
	{-1 << 32, 10, "-4294967296"},
	{1<<32 + 1, 10, "4294967297"},
	{-1<<32 - 1, 10, "-4294967297"},
	{1 << 50, 10, "1125899906842624"},
	{1<<63 - 1, 10, "9223372036854775807"},
	{-1<<63 + 1, 10, "-9223372036854775807"},
	{-1 << 63, 10, "-9223372036854775808"},

	{0, 2, "0"},
	{10, 2, "1010"},
	{-1, 2, "-1"},
	{1 << 15, 2, "1000000000000000"},

	{-8, 8, "-10"},
	{057635436545, 8, "57635436545"},
	{1 << 24, 8, "100000000"},

	{16, 16, "10"},
	{-0x123456789abcdef, 16, "-123456789abcdef"},
	{1<<63 - 1, 16, "7fffffffffffffff"},
	{1<<63 - 1, 2, "111111111111111111111111111111111111111111111111111111111111111"},
	{-1 << 63, 2, "-1000000000000000000000000000000000000000000000000000000000000000"},

	{16, 17, "g"},
	{25, 25, "10"},
	{(((((17*35+24)*35+21)*35+34)*35+12)*35+24)*35 + 32, 35, "holycow"},
	{(((((17*36+24)*36+21)*36+34)*36+12)*36+24)*36 + 32, 36, "holycow"},
}

func TestItoa(t *testing.T) {
	for _, test := range itoa64tests {
		s := FormatInt(test.in, test.base)
		if s != test.out {
			t.Errorf("FormatInt(%v, %v) = %v want %v",
				test.in, test.base, s, test.out)
		}
		x := AppendInt([]byte("abc"), test.in, test.base)
		if string(x) != "abc"+test.out {
			t.Errorf("AppendInt(%q, %v, %v) = %q want %v",
				"abc", test.in, test.base, x, test.out)
		}

		if test.in >= 0 {
			s := FormatUint(uint64(test.in), test.base)
			if s != test.out {
				t.Errorf("FormatUint(%v, %v) = %v want %v",
					test.in, test.base, s, test.out)
			}
			x := AppendUint(nil, uint64(test.in), test.base)
			if string(x) != test.out {
				t.Errorf("AppendUint(%q, %v, %v) = %q want %v",
					"abc", uint64(test.in), test.base, x, test.out)
			}
		}

		if test.base == 10 && int64(int(test.in)) == test.in {
			s := Itoa(int(test.in))
			if s != test.out {
				t.Errorf("Itoa(%v) = %v want %v",
					test.in, s, test.out)
			}
		}
	}
}

type uitoa64Test struct {
	in   uint64
	base int
	out  string
}

var uitoa64tests = []uitoa64Test{
	{1<<63 - 1, 10, "9223372036854775807"},
	{1 << 63, 10, "9223372036854775808"},
	{1<<63 + 1, 10, "9223372036854775809"},
	{1<<64 - 2, 10, "18446744073709551614"},
	{1<<64 - 1, 10, "18446744073709551615"},
	{1<<64 - 1, 2, "1111111111111111111111111111111111111111111111111111111111111111"},
}

func TestUitoa(t *testing.T) {
	for _, test := range uitoa64tests {
		s := FormatUint(test.in, test.base)
		if s != test.out {
			t.Errorf("FormatUint(%v, %v) = %v want %v",
				test.in, test.base, s, test.out)
		}
		x := AppendUint([]byte("abc"), test.in, test.base)
		if string(x) != "abc"+test.out {
			t.Errorf("AppendUint(%q, %v, %v) = %q want %v",
				"abc", test.in, test.base, x, test.out)
		}

	}
}

var varlenUints = []struct {
	in  uint64
	out string
}{
	{1, "1"},
	{12, "12"},
	{123, "123"},
	{1234, "1234"},
	{12345, "12345"},
	{123456, "123456"},
	{1234567, "1234567"},
	{12345678, "12345678"},
	{123456789, "123456789"},
	{1234567890, "1234567890"},
	{12345678901, "12345678901"},
	{123456789012, "123456789012"},
	{1234567890123, "1234567890123"},
	{12345678901234, "12345678901234"},
	{123456789012345, "123456789012345"},
	{1234567890123456, "1234567890123456"},
	{12345678901234567, "12345678901234567"},
	{123456789012345678, "123456789012345678"},
	{1234567890123456789, "1234567890123456789"},
	{12345678901234567890, "12345678901234567890"},
}

func TestFormatUintVarlen(t *testing.T) {
	for _, test := range varlenUints {
		s := FormatUint(test.in, 10)
		if s != test.out {
			t.Errorf("FormatUint(%v, 10) = %v want %v", test.in, s, test.out)
		}
	}
}

func BenchmarkFormatInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range itoa64tests {
			s := FormatInt(test.in, test.base)
			BenchSink += len(s)
		}
	}
}

func BenchmarkAppendInt(b *testing.B) {
	dst := make([]byte, 0, 30)
	for i := 0; i < b.N; i++ {
		for _, test := range itoa64tests {
			dst = AppendInt(dst[:0], test.in, test.base)
			BenchSink += len(dst)
		}
	}
}

func BenchmarkFormatUint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range uitoa64tests {
			s := FormatUint(test.in, test.base)
			BenchSink += len(s)
		}
	}
}

func BenchmarkAppendUint(b *testing.B) {
	dst := make([]byte, 0, 30)
	for i := 0; i < b.N; i++ {
		for _, test := range uitoa64tests {
			dst = AppendUint(dst[:0], test.in, test.base)
			BenchSink += len(dst)
		}
	}
}

func BenchmarkFormatIntSmall(b *testing.B) {
	smallInts := []int64{7, 42}
	for _, smallInt := range smallInts {
		b.Run(Itoa(int(smallInt)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s := FormatInt(smallInt, 10)
				BenchSink += len(s)
			}
		})
	}
}

func BenchmarkAppendIntSmall(b *testing.B) {
	dst := make([]byte, 0, 30)
	const smallInt = 42
	for i := 0; i < b.N; i++ {
		dst = AppendInt(dst[:0], smallInt, 10)
		BenchSink += len(dst)
	}
}

func BenchmarkAppendUintVarlen(b *testing.B) {
	for _, test := range varlenUints {
		b.Run(test.out, func(b *testing.B) {
			dst := make([]byte, 0, 30)
			for j := 0; j < b.N; j++ {
				dst = AppendUint(dst[:0], test.in, 10)
				BenchSink += len(dst)
			}
		})
	}
}

var BenchSink int // make sure compiler cannot optimize away benchmarks

func TestFormatBits_Panic(t *testing.T) {
	// Test panic cases for formatBits
	defer func() {
		if r := recover(); r == nil {
			t.Error("formatBits should panic for base < 2")
		}
	}()
	FormatUint(100, 1) // base < 2 should panic
}

func TestFormatBits_Panic2(t *testing.T) {
	// Test panic cases for formatBits
	defer func() {
		if r := recover(); r == nil {
			t.Error("formatBits should panic for base > 62")
		}
	}()
	FormatUint(100, 63) // base > 62 should panic
}

func TestFormatBits_SingleDigit(t *testing.T) {
	// Test formatBits with single digit results (us < 10 in base 10)
	testCases := []struct {
		in   uint64
		base int
		want string
	}{
		{0, 10, "0"},
		{1, 10, "1"},
		{2, 10, "2"},
		{3, 10, "3"},
		{4, 10, "4"},
		{5, 10, "5"},
		{6, 10, "6"},
		{7, 10, "7"},
		{8, 10, "8"},
		{9, 10, "9"},
	}

	for _, tc := range testCases {
		got := FormatUint(tc.in, tc.base)
		if got != tc.want {
			t.Errorf("FormatUint(%d, %d) = %q, want %q", tc.in, tc.base, got, tc.want)
		}
	}
}

func TestFormatBits_NonPowerOfTwoBase(t *testing.T) {
	// Test formatBits with bases that are not powers of 2 (general case)
	testCases := []struct {
		in   uint64
		base int
		want string
	}{
		{10, 3, "101"},   // base 3
		{10, 5, "20"},    // base 5
		{10, 7, "13"},   // base 7
		{10, 9, "11"},   // base 9
		{10, 11, "a"},   // base 11
		{10, 13, "a"},   // base 13
		{10, 15, "a"},   // base 15
		{10, 17, "a"},   // base 17
		{10, 19, "a"},   // base 19
		{10, 21, "a"},   // base 21
		{10, 23, "a"},   // base 23
		{10, 25, "a"},   // base 25
		{10, 27, "a"},   // base 27
		{10, 29, "a"},   // base 29
		{10, 31, "a"},   // base 31
		{10, 33, "a"},   // base 33
		{10, 35, "a"},   // base 35
		{10, 37, "a"},   // base 37
		{10, 39, "a"},   // base 39
		{10, 41, "a"},   // base 41
		{10, 43, "a"},   // base 43
		{10, 45, "a"},   // base 45
		{10, 47, "a"},   // base 47
		{10, 49, "a"},   // base 49
		{10, 51, "a"},   // base 51
		{10, 53, "a"},   // base 53
		{10, 55, "a"},   // base 55
		{10, 57, "a"},   // base 57
		{10, 59, "a"},   // base 59
		{10, 61, "a"},   // base 61
	}

	for _, tc := range testCases {
		got := FormatUint(tc.in, tc.base)
		if got != tc.want {
			t.Errorf("FormatUint(%d, %d) = %q, want %q", tc.in, tc.base, got, tc.want)
		}
	}
}

func TestFormatBits_AppendMode(t *testing.T) {
	// Test formatBits with append_ = true
	dst := []byte("prefix")
	result := AppendUint(dst, 123, 10)
	expected := "prefix123"
	if string(result) != expected {
		t.Errorf("AppendUint([]byte(\"prefix\"), 123, 10) = %q, want %q", string(result), expected)
	}

	dst2 := []byte("test")
	result2 := AppendInt(dst2, -456, 10)
	expected2 := "test-456"
	if string(result2) != expected2 {
		t.Errorf("AppendInt([]byte(\"test\"), -456, 10) = %q, want %q", string(result2), expected2)
	}
}

func TestFormatBits_LargeNumbers(t *testing.T) {
	// Test formatBits with very large numbers to cover Host32bit branch
	// This will test the u >= 1e9 loop in Host32bit case
	largeNum := uint64(1e9 + 1)
	result := FormatUint(largeNum, 10)
	expected := "1000000001"
	if result != expected {
		t.Errorf("FormatUint(%d, 10) = %q, want %q", largeNum, result, expected)
	}

	// Test with even larger number
	veryLargeNum := uint64(1e18)
	result2 := FormatUint(veryLargeNum, 10)
	if len(result2) == 0 {
		t.Error("FormatUint(1e18, 10) should return non-empty string")
	}
}

func TestSmall_SingleDigit(t *testing.T) {
	// Test small function with single digits (i < 10)
	for i := 0; i < 10; i++ {
		result := FormatUint(uint64(i), 10)
		expected := string(rune('0' + i))
		if result != expected {
			t.Errorf("FormatUint(%d, 10) = %q, want %q", i, result, expected)
		}
	}
}

func TestFormatBits_Base10_EdgeCases(t *testing.T) {
	// Test formatBits with base 10 edge cases
	testCases := []struct {
		in   uint64
		want string
	}{
		{10, "10"},
		{99, "99"},
		{100, "100"},
		{999, "999"},
		{1000, "1000"},
		{9999, "9999"},
		{10000, "10000"},
	}

	for _, tc := range testCases {
		got := FormatUint(tc.in, 10)
		if got != tc.want {
			t.Errorf("FormatUint(%d, 10) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestFormatBits_Base10_SingleDigit(t *testing.T) {
	// Test formatBits with single digit case (us < 10)
	// This covers the branch where us < 10 after the loop
	for i := uint64(0); i < 10; i++ {
		got := FormatUint(i, 10)
		want := string(rune('0' + i))
		if got != want {
			t.Errorf("FormatUint(%d, 10) = %q, want %q", i, got, want)
		}
	}
}

func TestFormatBits_NegativeNumbers(t *testing.T) {
	// Test formatBits with negative numbers (neg = true)
	testCases := []struct {
		in   int64
		base int
		want string
	}{
		{-1, 10, "-1"},
		{-10, 10, "-10"},
		{-100, 10, "-100"},
		{-123, 10, "-123"},
		{-1, 2, "-1"},
		{-10, 16, "-a"},
		{-255, 16, "-ff"},
	}

	for _, tc := range testCases {
		got := FormatInt(tc.in, tc.base)
		if got != tc.want {
			t.Errorf("FormatInt(%d, %d) = %q, want %q", tc.in, tc.base, got, tc.want)
		}
	}
}

func TestFormatBits_PowerOfTwoBases(t *testing.T) {
	// Test formatBits with power of two bases (2, 4, 8, 16, 32)
	testCases := []struct {
		in   uint64
		base int
		want string
	}{
		{15, 2, "1111"},
		{15, 4, "33"},
		{15, 8, "17"},
		{15, 16, "f"},
		{31, 32, "v"},
		{255, 16, "ff"},
		{1023, 32, "vv"},
	}

	for _, tc := range testCases {
		got := FormatUint(tc.in, tc.base)
		if got != tc.want {
			t.Errorf("FormatUint(%d, %d) = %q, want %q", tc.in, tc.base, got, tc.want)
		}
	}
}

func TestFormatBits_Base10_SingleDigitResult(t *testing.T) {
	// Test formatBits with base 10 where result is single digit (us < 10)
	testCases := []struct {
		in   uint64
		want string
	}{
		{0, "0"},
		{1, "1"},
		{2, "2"},
		{3, "3"},
		{4, "4"},
		{5, "5"},
		{6, "6"},
		{7, "7"},
		{8, "8"},
		{9, "9"},
	}

	for _, tc := range testCases {
		got := FormatUint(tc.in, 10)
		if got != tc.want {
			t.Errorf("FormatUint(%d, 10) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestFormatBits_Base10_TwoDigitResult(t *testing.T) {
	// Test formatBits with base 10 where result is two digits (us < 100 && us >= 10)
	testCases := []struct {
		in   uint64
		want string
	}{
		{10, "10"},
		{11, "11"},
		{99, "99"},
	}

	for _, tc := range testCases {
		got := FormatUint(tc.in, 10)
		if got != tc.want {
			t.Errorf("FormatUint(%d, 10) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestFormatBits_AppendInt_Negative(t *testing.T) {
	// Test AppendInt with negative numbers
	dst := []byte("prefix")
	result := AppendInt(dst, -123, 10)
	expected := "prefix-123"
	if string(result) != expected {
		t.Errorf("AppendInt([]byte(\"prefix\"), -123, 10) = %q, want %q", string(result), expected)
	}
}

func TestFormatBits_AppendUint_AllBases(t *testing.T) {
	// Test AppendUint with various bases
	testCases := []struct {
		dst  []byte
		in   uint64
		base int
		want string
	}{
		{[]byte("test"), 255, 16, "testff"},
		{[]byte(""), 15, 2, "1111"},
		{[]byte("x"), 31, 32, "xv"},
	}

	for _, tc := range testCases {
		got := AppendUint(tc.dst, tc.in, tc.base)
		if string(got) != tc.want {
			t.Errorf("AppendUint(%q, %d, %d) = %q, want %q", string(tc.dst), tc.in, tc.base, string(got), tc.want)
		}
	}
}
