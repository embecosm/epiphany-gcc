// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strconv_test

import (
	"math"
	. "strconv"
	"testing"
)

type ftoaTest struct {
	f    float64
	fmt  byte
	prec int
	s    string
}

func fdiv(a, b float64) float64 { return a / b } // keep compiler in the dark

const (
	below1e23 = 99999999999999974834176
	above1e23 = 100000000000000008388608
)

var ftoatests = []ftoaTest{
	{1, 'e', 5, "1.00000e+00"},
	{1, 'f', 5, "1.00000"},
	{1, 'g', 5, "1"},
	{1, 'g', -1, "1"},
	{20, 'g', -1, "20"},
	{1234567.8, 'g', -1, "1.2345678e+06"},
	{200000, 'g', -1, "200000"},
	{2000000, 'g', -1, "2e+06"},

	// g conversion and zero suppression
	{400, 'g', 2, "4e+02"},
	{40, 'g', 2, "40"},
	{4, 'g', 2, "4"},
	{.4, 'g', 2, "0.4"},
	{.04, 'g', 2, "0.04"},
	{.004, 'g', 2, "0.004"},
	{.0004, 'g', 2, "0.0004"},
	{.00004, 'g', 2, "4e-05"},
	{.000004, 'g', 2, "4e-06"},

	{0, 'e', 5, "0.00000e+00"},
	{0, 'f', 5, "0.00000"},
	{0, 'g', 5, "0"},
	{0, 'g', -1, "0"},

	{-1, 'e', 5, "-1.00000e+00"},
	{-1, 'f', 5, "-1.00000"},
	{-1, 'g', 5, "-1"},
	{-1, 'g', -1, "-1"},

	{12, 'e', 5, "1.20000e+01"},
	{12, 'f', 5, "12.00000"},
	{12, 'g', 5, "12"},
	{12, 'g', -1, "12"},

	{123456700, 'e', 5, "1.23457e+08"},
	{123456700, 'f', 5, "123456700.00000"},
	{123456700, 'g', 5, "1.2346e+08"},
	{123456700, 'g', -1, "1.234567e+08"},

	{1.2345e6, 'e', 5, "1.23450e+06"},
	{1.2345e6, 'f', 5, "1234500.00000"},
	{1.2345e6, 'g', 5, "1.2345e+06"},

	{1e23, 'e', 17, "9.99999999999999916e+22"},
	{1e23, 'f', 17, "99999999999999991611392.00000000000000000"},
	{1e23, 'g', 17, "9.9999999999999992e+22"},

	{1e23, 'e', -1, "1e+23"},
	{1e23, 'f', -1, "100000000000000000000000"},
	{1e23, 'g', -1, "1e+23"},

	{below1e23, 'e', 17, "9.99999999999999748e+22"},
	{below1e23, 'f', 17, "99999999999999974834176.00000000000000000"},
	{below1e23, 'g', 17, "9.9999999999999975e+22"},

	{below1e23, 'e', -1, "9.999999999999997e+22"},
	{below1e23, 'f', -1, "99999999999999970000000"},
	{below1e23, 'g', -1, "9.999999999999997e+22"},

	{above1e23, 'e', 17, "1.00000000000000008e+23"},
	{above1e23, 'f', 17, "100000000000000008388608.00000000000000000"},
	{above1e23, 'g', 17, "1.0000000000000001e+23"},

	{above1e23, 'e', -1, "1.0000000000000001e+23"},
	{above1e23, 'f', -1, "100000000000000010000000"},
	{above1e23, 'g', -1, "1.0000000000000001e+23"},

	{fdiv(5e-304, 1e20), 'g', -1, "5e-324"},
	{fdiv(-5e-304, 1e20), 'g', -1, "-5e-324"},

	{32, 'g', -1, "32"},
	{32, 'g', 0, "3e+01"},

	{100, 'x', -1, "%x"},

	{math.NaN(), 'g', -1, "NaN"},
	{-math.NaN(), 'g', -1, "NaN"},
	{math.Inf(0), 'g', -1, "+Inf"},
	{math.Inf(-1), 'g', -1, "-Inf"},
	{-math.Inf(0), 'g', -1, "-Inf"},

	{-1, 'b', -1, "-4503599627370496p-52"},

	// fixed bugs
	{0.9, 'f', 1, "0.9"},
	{0.09, 'f', 1, "0.1"},
	{0.0999, 'f', 1, "0.1"},
	{0.05, 'f', 1, "0.1"},
	{0.05, 'f', 0, "0"},
	{0.5, 'f', 1, "0.5"},
	{0.5, 'f', 0, "0"},
	{1.5, 'f', 0, "2"},

	// http://www.exploringbinary.com/java-hangs-when-converting-2-2250738585072012e-308/
	{2.2250738585072012e-308, 'g', -1, "2.2250738585072014e-308"},
	// http://www.exploringbinary.com/php-hangs-on-numeric-value-2-2250738585072011e-308/
	{2.2250738585072011e-308, 'g', -1, "2.225073858507201e-308"},
}

func TestFtoa(t *testing.T) {
	for i := 0; i < len(ftoatests); i++ {
		test := &ftoatests[i]
		s := FormatFloat(test.f, test.fmt, test.prec, 64)
		if s != test.s {
			t.Error("testN=64", test.f, string(test.fmt), test.prec, "want", test.s, "got", s)
		}
		x := AppendFloat([]byte("abc"), test.f, test.fmt, test.prec, 64)
		if string(x) != "abc"+test.s {
			t.Error("AppendFloat testN=64", test.f, string(test.fmt), test.prec, "want", "abc"+test.s, "got", string(x))
		}
		if float64(float32(test.f)) == test.f && test.fmt != 'b' {
			s := FormatFloat(test.f, test.fmt, test.prec, 32)
			if s != test.s {
				t.Error("testN=32", test.f, string(test.fmt), test.prec, "want", test.s, "got", s)
			}
			x := AppendFloat([]byte("abc"), test.f, test.fmt, test.prec, 32)
			if string(x) != "abc"+test.s {
				t.Error("AppendFloat testN=32", test.f, string(test.fmt), test.prec, "want", "abc"+test.s, "got", string(x))
			}
		}
	}
}

/* This test relies on escape analysis which gccgo does not yet do.

func TestAppendFloatDoesntAllocate(t *testing.T) {
	n := numAllocations(func() {
		var buf [64]byte
		AppendFloat(buf[:0], 1.23, 'g', 5, 64)
	})
	want := 1 // TODO(bradfitz): this might be 0, once escape analysis is better
	if n != want {
		t.Errorf("with local buffer, did %d allocations, want %d", n, want)
	}
	n = numAllocations(func() {
		AppendFloat(globalBuf[:0], 1.23, 'g', 5, 64)
	})
	if n != 0 {
		t.Errorf("with reused buffer, did %d allocations, want 0", n)
	}
}

*/

func BenchmarkFormatFloatDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormatFloat(33909, 'g', -1, 64)
	}
}

func BenchmarkFormatFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormatFloat(339.7784, 'g', -1, 64)
	}
}

func BenchmarkFormatFloatExp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormatFloat(-5.09e75, 'g', -1, 64)
	}
}

func BenchmarkFormatFloatBig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormatFloat(123456789123456789123456789, 'g', -1, 64)
	}
}

func BenchmarkAppendFloatDecimal(b *testing.B) {
	dst := make([]byte, 0, 30)
	for i := 0; i < b.N; i++ {
		AppendFloat(dst, 33909, 'g', -1, 64)
	}
}

func BenchmarkAppendFloat(b *testing.B) {
	dst := make([]byte, 0, 30)
	for i := 0; i < b.N; i++ {
		AppendFloat(dst, 339.7784, 'g', -1, 64)
	}
}

func BenchmarkAppendFloatExp(b *testing.B) {
	dst := make([]byte, 0, 30)
	for i := 0; i < b.N; i++ {
		AppendFloat(dst, -5.09e75, 'g', -1, 64)
	}
}

func BenchmarkAppendFloatBig(b *testing.B) {
	dst := make([]byte, 0, 30)
	for i := 0; i < b.N; i++ {
		AppendFloat(dst, 123456789123456789123456789, 'g', -1, 64)
	}
}
