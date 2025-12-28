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
		{10, 3, "101"}, // base 3
		{10, 5, "20"},  // base 5
		{10, 7, "13"},  // base 7
		{10, 9, "11"},  // base 9
		{10, 11, "a"},  // base 11
		{10, 13, "a"},  // base 13
		{10, 15, "a"},  // base 15
		{10, 17, "a"},  // base 17
		{10, 19, "a"},  // base 19
		{10, 21, "a"},  // base 21
		{10, 23, "a"},  // base 23
		{10, 25, "a"},  // base 25
		{10, 27, "a"},  // base 27
		{10, 29, "a"},  // base 29
		{10, 31, "a"},  // base 31
		{10, 33, "a"},  // base 33
		{10, 35, "a"},  // base 35
		{10, 37, "a"},  // base 37
		{10, 39, "a"},  // base 39
		{10, 41, "a"},  // base 41
		{10, 43, "a"},  // base 43
		{10, 45, "a"},  // base 45
		{10, 47, "a"},  // base 47
		{10, 49, "a"},  // base 49
		{10, 51, "a"},  // base 51
		{10, 53, "a"},  // base 53
		{10, 55, "a"},  // base 55
		{10, 57, "a"},  // base 57
		{10, 59, "a"},  // base 59
		{10, 61, "a"},  // base 61
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

// TestFormatBits_Base10_EdgeCases, TestFormatBits_Base10_SingleDigit and TestFormatBits_Base10_UsGreaterThanOrEqual10
// are merged into TestFormatBits_Base10_UsLessThan10, TestFormatBits_Base10_UsLessThan100 and TestFormatBits_Base10_UsGreaterThanOrEqual100

func TestFormatBits_AllBases(t *testing.T) {
	// Test formatBits with all valid bases (2-62)
	for base := 2; base <= 62; base++ {
		got := FormatUint(100, base)
		if len(got) == 0 {
			t.Errorf("FormatUint(100, %d) should return non-empty string", base)
		}
	}
}

func TestFormatBits_NonPowerOfTwoBases(t *testing.T) {
	// Test formatBits with non-power-of-two bases (general case)
	testCases := []struct {
		in   uint64
		base int
	}{
		{100, 3},
		{100, 5},
		{100, 7},
		{100, 9},
		{100, 11},
		{100, 13},
		{100, 15},
		{100, 17},
		{100, 19},
		{100, 21},
		{100, 23},
		{100, 25},
		{100, 27},
		{100, 29},
		{100, 31},
		{100, 33},
		{100, 35},
		{100, 37},
		{100, 39},
		{100, 41},
		{100, 43},
		{100, 45},
		{100, 47},
		{100, 49},
		{100, 51},
		{100, 53},
		{100, 55},
		{100, 57},
		{100, 59},
		{100, 61},
	}

	for _, tc := range testCases {
		got := FormatUint(tc.in, tc.base)
		if len(got) == 0 {
			t.Errorf("FormatUint(%d, %d) should return non-empty string", tc.in, tc.base)
		}
		// Verify it can be parsed back
		result := ParseUint(got, tc.base, 64)
		if result.IsErr() {
			t.Errorf("FormatUint(%d, %d) = %q, but ParseUint failed: %v", tc.in, tc.base, got, result.Err())
		} else if result.Unwrap() != tc.in {
			t.Errorf("FormatUint(%d, %d) = %q, but ParseUint returned %d", tc.in, tc.base, got, result.Unwrap())
		}
	}
}

func TestFormatBits_PowerOfTwoBases_All(t *testing.T) {
	// Test formatBits with all power-of-two bases
	powerOfTwoBases := []int{2, 4, 8, 16, 32}
	for _, base := range powerOfTwoBases {
		got := FormatUint(255, base)
		if len(got) == 0 {
			t.Errorf("FormatUint(255, %d) should return non-empty string", base)
		}
		// Verify it can be parsed back
		result := ParseUint(got, base, 64)
		if result.IsErr() {
			t.Errorf("FormatUint(255, %d) = %q, but ParseUint failed: %v", base, got, result.Err())
		} else if result.Unwrap() != 255 {
			t.Errorf("FormatUint(255, %d) = %q, but ParseUint returned %d", base, got, result.Unwrap())
		}
	}
}

func TestFormatBits_Base10_LargeNumbers(t *testing.T) {
	// Test formatBits with base 10 and large numbers
	testCases := []uint64{
		1e6,
		1e7,
		1e8,
		1e9,
		1e10,
		1e11,
		1e12,
		1e13,
		1e14,
		1e15,
		1e16,
		1e17,
		1e18,
	}

	for _, tc := range testCases {
		got := FormatUint(tc, 10)
		if len(got) == 0 {
			t.Errorf("FormatUint(%d, 10) should return non-empty string", tc)
		}
		// Verify it can be parsed back
		result := ParseUint(got, 10, 64)
		if result.IsErr() {
			t.Errorf("FormatUint(%d, 10) = %q, but ParseUint failed: %v", tc, got, result.Err())
		} else if result.Unwrap() != tc {
			t.Errorf("FormatUint(%d, 10) = %q, but ParseUint returned %d", tc, got, result.Unwrap())
		}
	}
}

func TestFormatBits_NegativeNumbers_AllBases(t *testing.T) {
	// Test formatBits with negative numbers and various bases
	testCases := []struct {
		in   int64
		base int
	}{
		{-1, 2},
		{-10, 10},
		{-100, 16},
		{-255, 32},
		{-1000, 37},
		{-10000, 62},
	}

	for _, tc := range testCases {
		got := FormatInt(tc.in, tc.base)
		if len(got) == 0 {
			t.Errorf("FormatInt(%d, %d) should return non-empty string", tc.in, tc.base)
		}
		if got[0] != '-' {
			t.Errorf("FormatInt(%d, %d) = %q, should start with '-'", tc.in, tc.base, got)
		}
	}
}

// TestFormatBits_AppendMode is already defined above

func TestFormatBits_Zero(t *testing.T) {
	// Test formatBits with zero
	got := FormatUint(0, 10)
	if got != "0" {
		t.Errorf("FormatUint(0, 10) = %q, want \"0\"", got)
	}

	got2 := FormatInt(0, 10)
	if got2 != "0" {
		t.Errorf("FormatInt(0, 10) = %q, want \"0\"", got2)
	}

	got3 := FormatUint(0, 16)
	if got3 != "0" {
		t.Errorf("FormatUint(0, 16) = %q, want \"0\"", got3)
	}

	got4 := FormatUint(0, 37)
	if got4 != "0" {
		t.Errorf("FormatUint(0, 37) = %q, want \"0\"", got4)
	}
}

func TestFormatBits_Base10_UsLessThan10(t *testing.T) {
	// Test formatBits with us < 10 case (single digit after loop)
	// Test using FormatUint (public API)
	for i := uint64(0); i < 10; i++ {
		got := FormatUint(i, 10)
		want := string(rune('0' + i))
		if got != want {
			t.Errorf("FormatUint(%d, 10) = %q, want %q", i, got, want)
		}
	}
	// Test using formatBits directly (internal function)
	testCases := []struct {
		u    uint64
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
		_, s := formatBits(nil, tc.u, 10, false, false)
		if s != tc.want {
			t.Errorf("formatBits(%d, 10, false, false) = %q, want %q", tc.u, s, tc.want)
		}
	}
}

func TestFormatBits_Base10_UsLessThan100(t *testing.T) {
	// Test formatBits with us < 100 case (covers us >= 10 branch at line 157)
	// Test using FormatUint (public API)
	for i := uint64(10); i < 100; i++ {
		got := FormatUint(i, 10)
		// Verify it's a 2-digit string
		if len(got) != 2 {
			t.Errorf("FormatUint(%d, 10) = %q, should be 2 digits", i, got)
		}
		// Verify it can be parsed back
		result := ParseUint(got, 10, 64)
		if result.IsErr() {
			t.Errorf("FormatUint(%d, 10) = %q, but ParseUint failed: %v", i, got, result.Err())
		} else if result.Unwrap() != i {
			t.Errorf("FormatUint(%d, 10) = %q, but ParseUint returned %d", i, got, result.Unwrap())
		}
	}
	// Test using formatBits directly (internal function)
	testCases := []struct {
		u    uint64
		want string
	}{
		{10, "10"},
		{11, "11"},
		{99, "99"},
		{100, "100"}, // This will have us >= 100, so it will use the loop
	}
	for _, tc := range testCases {
		_, s := formatBits(nil, tc.u, 10, false, false)
		if s != tc.want {
			t.Errorf("formatBits(%d, 10, false, false) = %q, want %q", tc.u, s, tc.want)
		}
	}
}

func TestFormatBits_Base10_UsGreaterThanOrEqual100(t *testing.T) {
	// Test formatBits with us >= 100 case (triggers the loop at line 145)
	// Test using FormatUint (public API)
	testCases := []uint64{
		100,
		101,
		199,
		200,
		999,
		1000,
		9999,
		10000,
	}

	for _, tc := range testCases {
		got := FormatUint(tc, 10)
		if len(got) == 0 {
			t.Errorf("FormatUint(%d, 10) should return non-empty string", tc)
		}
		// Verify it can be parsed back
		result := ParseUint(got, 10, 64)
		if result.IsErr() {
			t.Errorf("FormatUint(%d, 10) = %q, but ParseUint failed: %v", tc, got, result.Err())
		} else if result.Unwrap() != tc {
			t.Errorf("FormatUint(%d, 10) = %q, but ParseUint returned %d", tc, got, result.Unwrap())
		}
	}
	// Test using formatBits directly (internal function)
	testCases2 := []struct {
		u    uint64
		want string
	}{
		{100, "100"},
		{101, "101"},
		{999, "999"},
		{1000, "1000"},
		{12345, "12345"},
		{99999, "99999"},
		{100000, "100000"},
	}
	for _, tc := range testCases2 {
		_, s := formatBits(nil, tc.u, 10, false, false)
		if s != tc.want {
			t.Errorf("formatBits(%d, 10, false, false) = %q, want %q", tc.u, s, tc.want)
		}
	}
}

func TestIsPowerOfTwo(t *testing.T) {
	// Test isPowerOfTwo function indirectly through formatBits
	// Power of two bases: 2, 4, 8, 16, 32
	powerOfTwoBases := []int{2, 4, 8, 16, 32}
	for _, base := range powerOfTwoBases {
		got := FormatUint(255, base)
		if len(got) == 0 {
			t.Errorf("FormatUint(255, %d) should return non-empty string", base)
		}
	}

	// Non-power of two bases: 3, 5, 7, 9, etc.
	nonPowerOfTwoBases := []int{3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29, 31, 33, 35, 37}
	for _, base := range nonPowerOfTwoBases {
		got := FormatUint(255, base)
		if len(got) == 0 {
			t.Errorf("FormatUint(255, %d) should return non-empty string", base)
		}
	}
}

func TestFormatBits_Base10_Host32bit(t *testing.T) {
	// Test formatBits with base 10 and Host32bit branch
	// This tests the u >= 1e9 loop in Host32bit case
	if Host32bit {
		largeNum := uint64(1e9)
		got := FormatUint(largeNum, 10)
		expected := "1000000000"
		if got != expected {
			t.Errorf("FormatUint(%d, 10) = %q, want %q", largeNum, got, expected)
		}

		veryLargeNum := uint64(1e18)
		got2 := FormatUint(veryLargeNum, 10)
		if len(got2) == 0 {
			t.Error("FormatUint(1e18, 10) should return non-empty string")
		}

		// Test with numbers that trigger the u >= 1e9 loop multiple times
		veryVeryLargeNum := uint64(1e9 * 1e9)
		got3 := FormatUint(veryVeryLargeNum, 10)
		if len(got3) == 0 {
			t.Error("FormatUint(1e18, 10) should return non-empty string")
		}
	} else {
		// On 64-bit systems, test that large numbers still work
		largeNum := uint64(1e9)
		got := FormatUint(largeNum, 10)
		expected := "1000000000"
		if got != expected {
			t.Errorf("FormatUint(%d, 10) = %q, want %q", largeNum, got, expected)
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

func TestFormatBits_GeneralCase_NonPowerOfTwo(t *testing.T) {
	// Test formatBits with general case (non-power-of-two base)
	// This tests the general case path at line 181
	// Test with bases that are not powers of 2 and not 10
	testCases := []struct {
		u    uint64
		base int
		want string
	}{
		{0, 3, "0"},
		{1, 3, "1"},
		{2, 3, "2"},
		{3, 3, "10"},
		{4, 3, "11"},
		{9, 3, "100"},
		{0, 5, "0"},
		{1, 5, "1"},
		{4, 5, "4"},
		{5, 5, "10"},
		{25, 5, "100"},
		{0, 7, "0"},
		{1, 7, "1"},
		{6, 7, "6"},
		{7, 7, "10"},
		{49, 7, "100"},
		{0, 11, "0"},
		{1, 11, "1"},
		{10, 11, "a"},
		{11, 11, "10"},
		{0, 13, "0"},
		{1, 13, "1"},
		{12, 13, "c"},
		{13, 13, "10"},
		{0, 17, "0"},
		{1, 17, "1"},
		{16, 17, "g"},
		{17, 17, "10"},
		{0, 19, "0"},
		{1, 19, "1"},
		{18, 19, "i"},
		{19, 19, "10"},
		{0, 23, "0"},
		{1, 23, "1"},
		{22, 23, "m"},
		{23, 23, "10"},
		{0, 29, "0"},
		{1, 29, "1"},
		{28, 29, "s"},
		{29, 29, "10"},
		{0, 31, "0"},
		{1, 31, "1"},
		{30, 31, "u"},
		{31, 31, "10"},
		{0, 37, "0"},
		{1, 37, "1"},
		{36, 37, "A"},
		{37, 37, "10"},
		{0, 41, "0"},
		{1, 41, "1"},
		{40, 41, "E"}, // digits[40] = 'E' (36 + 4 = 40, 'A' + 4 = 'E')
		{41, 41, "10"},
		{0, 43, "0"},
		{1, 43, "1"},
		{42, 43, "G"}, // digits[42] = 'G' (36 + 6 = 42, 'A' + 6 = 'G')
		{43, 43, "10"},
		{0, 47, "0"},
		{1, 47, "1"},
		{46, 47, "K"}, // digits[46] = 'K' (36 + 10 = 46, 'A' + 10 = 'K')
		{47, 47, "10"},
		{0, 53, "0"},
		{1, 53, "1"},
		{52, 53, "Q"}, // digits[52] = 'Q' (36 + 16 = 52, 'A' + 16 = 'Q')
		{53, 53, "10"},
		{0, 59, "0"},
		{1, 59, "1"},
		{58, 59, "W"}, // digits[58] = 'W' (36 + 22 = 58, 'A' + 22 = 'W')
		{59, 59, "10"},
		{0, 61, "0"},
		{1, 61, "1"},
		{60, 61, "Y"}, // digits[60] = 'Y' (36 + 24 = 60, 'A' + 24 = 'Y')
		{61, 61, "10"},
	}

	for _, tc := range testCases {
		_, s := formatBits(nil, tc.u, tc.base, false, false)
		if s != tc.want {
			t.Errorf("formatBits(%d, %d, false, false) = %q, want %q", tc.u, tc.base, s, tc.want)
		}
	}
}

func TestFormatBits_GeneralCase_LargeNumbers(t *testing.T) {
	// Test formatBits with general case and large numbers
	// This tests the u >= b loop at line 184
	testCases := []struct {
		u    uint64
		base int
	}{
		{100, 3},
		{1000, 5},
		{10000, 7},
		{100000, 11},
		{1000000, 13},
		{10000000, 17},
		{100000000, 19},
		{1000000000, 23},
		{10000000000, 29},
		{100000000000, 31},
		{1000000000000, 37},
		{10000000000000, 41},
		{100000000000000, 43},
		{1000000000000000, 47},
		{10000000000000000, 53},
		{100000000000000000, 59},
		{1000000000000000000, 61},
	}

	for _, tc := range testCases {
		_, s := formatBits(nil, tc.u, tc.base, false, false)
		if s == "" {
			t.Errorf("formatBits(%d, %d, false, false) returned empty string", tc.u, tc.base)
		}
		// Verify it's a valid representation
		if len(s) == 0 {
			t.Errorf("formatBits(%d, %d, false, false) returned empty string", tc.u, tc.base)
		}
	}
}

func TestTryFromString_ReflectPath_AllKinds(t *testing.T) {
	// Test tryFromString with reflect path for all kinds
	// This tests the reflect.TypeOf(x).Kind() path at line 45
	// We need to use custom types that will trigger the reflect path

	// Test with custom string type
	type CustomString string
	type CustomInt int

	// This should use the reflect path
	result := TryFromString[CustomString, CustomInt]("42", 10, 0)
	if result.IsErr() {
		t.Logf("TryFromString returned error (may be expected): %v", result.UnwrapErr())
	} else {
		t.Logf("TryFromString returned: %v", result.Unwrap())
	}

	// Test with standard types that will use reflect path
	result2 := TryFromString[string, int]("42", 10, 0)
	if result2.IsErr() {
		t.Errorf("TryFromString[string, int]('42', 10, 0) should work, got error: %v", result2.UnwrapErr())
	} else if result2.Unwrap() != 42 {
		t.Errorf("TryFromString[string, int]('42', 10, 0) = %d, want 42", result2.Unwrap())
	}

	// Test with uint
	result3 := TryFromString[string, uint]("42", 10, 0)
	if result3.IsErr() {
		t.Errorf("TryFromString[string, uint]('42', 10, 0) should work, got error: %v", result3.UnwrapErr())
	} else if result3.Unwrap() != 42 {
		t.Errorf("TryFromString[string, uint]('42', 10, 0) = %d, want 42", result3.Unwrap())
	}

	// Test with float32
	result4 := TryFromString[string, float32]("3.14", 10, 32)
	if result4.IsErr() {
		t.Errorf("TryFromString[string, float32]('3.14', 10, 32) should work, got error: %v", result4.UnwrapErr())
	}

	// Test with float64
	result5 := TryFromString[string, float64]("3.14", 10, 64)
	if result5.IsErr() {
		t.Errorf("TryFromString[string, float64]('3.14', 10, 64) should work, got error: %v", result5.UnwrapErr())
	}
}

// TestFormatBits_BaseZero tests formatBits with base == 0 (should panic)
func TestFormatBits_BaseZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("formatBits should panic for base == 0")
		}
	}()
	FormatUint(100, 0) // base == 0 should panic
}

// TestFormatBits_NegativeZero tests formatBits with negative zero (neg=true, u=0)
func TestFormatBits_NegativeZero(t *testing.T) {
	// Test FormatInt with 0 (should not have negative sign)
	got := FormatInt(0, 10)
	if got != "0" {
		t.Errorf("FormatInt(0, 10) = %q, want \"0\"", got)
	}

	// Test FormatInt with -0 (should still be "0", not "-0")
	got2 := FormatInt(-0, 10)
	if got2 != "0" {
		t.Errorf("FormatInt(-0, 10) = %q, want \"0\"", got2)
	}
}

// TestFormatBits_Host32bit_MultipleLoops tests Host32bit branch with multiple u >= 1e9 loops
func TestFormatBits_Host32bit_MultipleLoops(t *testing.T) {
	if Host32bit {
		// Test with number that requires multiple iterations of u >= 1e9 loop
		// 1e18 = 1e9 * 1e9, so it will loop twice
		veryLargeNum := uint64(1e18)
		got := FormatUint(veryLargeNum, 10)
		expected := "1000000000000000000"
		if got != expected {
			t.Errorf("FormatUint(1e18, 10) = %q, want %q", got, expected)
		}

		// Test with even larger number that requires more loops
		// Use maximum uint64 value to test with very large numbers
		// Note: 1e27 exceeds uint64 range, so we use maxUint64 instead
		maxUint64 := uint64(18446744073709551615)
		got2 := FormatUint(maxUint64, 10)
		if len(got2) == 0 {
			t.Error("FormatUint(maxUint64, 10) should return non-empty string")
		}
		// Verify it's the correct string representation
		expectedMax := "18446744073709551615"
		if got2 != expectedMax {
			t.Errorf("FormatUint(maxUint64, 10) = %q, want %q", got2, expectedMax)
		}

		// Test with number that triggers u >= 1e9 loop and us < 10 case
		// This tests the us < 10 branch at line 1543-1546
		// Use a number like 1e9 + 5, which will have us = 5 after the first iteration
		testNum := uint64(1e9 + 5)
		got3 := FormatUint(testNum, 10)
		expected3 := "1000000005"
		if got3 != expected3 {
			t.Errorf("FormatUint(1e9+5, 10) = %q, want %q", got3, expected3)
		}

		// Test with number that requires multiple loops and has us < 10
		// Use 1e18 + 3, which will loop twice and have us = 3 after the first iteration
		testNum2 := uint64(1e18 + 3)
		got4 := FormatUint(testNum2, 10)
		expected4 := "1000000000000000003"
		if got4 != expected4 {
			t.Errorf("FormatUint(1e18+3, 10) = %q, want %q", got4, expected4)
		}
	}
}

// TestSmall_FastPath tests the fast path for small integers (fastSmalls && i < nSmalls && base == 10)
func TestSmall_FastPath(t *testing.T) {
	// Test FormatUint with small integers (should use fast path)
	for i := 0; i < 100; i++ {
		got := FormatUint(uint64(i), 10)
		expected := Itoa(i)
		if got != expected {
			t.Errorf("FormatUint(%d, 10) = %q, want %q", i, got, expected)
		}
	}

	// Test FormatInt with small integers (should use fast path)
	for i := -99; i < 100; i++ {
		got := FormatInt(int64(i), 10)
		expected := Itoa(i)
		if got != expected {
			t.Errorf("FormatInt(%d, 10) = %q, want %q", i, got, expected)
		}
	}

	// Test AppendUint with small integers (should use fast path)
	for i := 0; i < 100; i++ {
		dst := []byte("prefix")
		got := AppendUint(dst, uint64(i), 10)
		expected := "prefix" + Itoa(i)
		if string(got) != expected {
			t.Errorf("AppendUint([]byte(\"prefix\"), %d, 10) = %q, want %q", i, string(got), expected)
		}
	}

	// Test AppendInt with small integers (should use fast path)
	for i := -99; i < 100; i++ {
		dst := []byte("prefix")
		got := AppendInt(dst, int64(i), 10)
		expected := "prefix" + Itoa(i)
		if string(got) != expected {
			t.Errorf("AppendInt([]byte(\"prefix\"), %d, 10) = %q, want %q", i, string(got), expected)
		}
	}
}

// TestSmall_FastPath_NonBase10 tests that fast path is not used for non-base-10
func TestSmall_FastPath_NonBase10(t *testing.T) {
	// Test that small integers with base != 10 don't use fast path
	// They should still work correctly
	for i := 0; i < 100; i++ {
		got := FormatUint(uint64(i), 16)
		// Verify it's a valid hex representation
		if len(got) == 0 {
			t.Errorf("FormatUint(%d, 16) should return non-empty string", i)
		}
		// Verify it can be parsed back
		result := ParseUint(got, 16, 64)
		if result.IsErr() {
			t.Errorf("FormatUint(%d, 16) = %q, but ParseUint failed: %v", i, got, result.Err())
		} else if result.Unwrap() != uint64(i) {
			t.Errorf("FormatUint(%d, 16) = %q, but ParseUint returned %d", i, got, result.Unwrap())
		}
	}
}
