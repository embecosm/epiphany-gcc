// Copyright 2009-2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

// Floor returns the greatest integer value less than or equal to x.
//
// Special cases are:
//	Floor(±0) = ±0
//	Floor(±Inf) = ±Inf
//	Floor(NaN) = NaN
func libc_floor(float64) float64 __asm__("floor")
func Floor(x float64) float64 {
	return libc_floor(x)
}

func floor(x float64) float64 {
	// TODO(rsc): Remove manual inlining of IsNaN, IsInf
	// when compiler does it for us
	if x == 0 || x != x || x > MaxFloat64 || x < -MaxFloat64 { // x == 0 || IsNaN(x) || IsInf(x, 0)
		return x
	}
	if x < 0 {
		d, fract := Modf(-x)
		if fract != 0.0 {
			d = d + 1
		}
		return -d
	}
	d, _ := Modf(x)
	return d
}

// Ceil returns the least integer value greater than or equal to x.
//
// Special cases are:
//	Ceil(±0) = ±0
//	Ceil(±Inf) = ±Inf
//	Ceil(NaN) = NaN
func libc_ceil(float64) float64 __asm__("ceil")
func Ceil(x float64) float64 {
	return libc_ceil(x)
}

func ceil(x float64) float64 {
	return -Floor(-x)
}

// Trunc returns the integer value of x.
//
// Special cases are:
//	Trunc(±0) = ±0
//	Trunc(±Inf) = ±Inf
//	Trunc(NaN) = NaN
func libc_trunc(float64) float64 __asm__("trunc")
func Trunc(x float64) float64 {
	return libc_trunc(x)
}

func trunc(x float64) float64 {
	// TODO(rsc): Remove manual inlining of IsNaN, IsInf
	// when compiler does it for us
	if x == 0 || x != x || x > MaxFloat64 || x < -MaxFloat64 { // x == 0 || IsNaN(x) || IsInf(x, 0)
		return x
	}
	d, _ := Modf(x)
	return d
}
