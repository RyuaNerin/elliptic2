// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !gc || purego || !(amd64 || 386 || arm || arm64)

package simd

import (
	"math/big"
	"math/bits"
)

func addVV(z, x, y []big.Word) (c big.Word) {
	var cc uint
	for i := range z {
		zi, carry := bits.Add(uint(x[i]), uint(y[i]), cc)
		z[i] = big.Word(zi)
		cc = carry
	}
	return big.Word(cc)
}

func subVV(z, x, y []big.Word) (c big.Word) {
	var cc uint
	for i := range z {
		zi, borrow := bits.Sub(uint(x[i]), uint(y[i]), cc)
		z[i] = big.Word(zi)
		cc = borrow
	}
	return big.Word(cc)
}

func lshVU(z, x []big.Word, s uint) (c big.Word) {
	if len(z) == 0 {
		return 0
	}
	if len(x) != len(z) {
		panic("simd: lshVU len")
	}

	ŝ := WordBitSize - s

	c = x[len(z)-1] >> ŝ
	for i := len(z) - 1; i > 0; i-- {
		z[i] = x[i]<<s | x[i-1]>>ŝ
	}
	z[0] = x[0] << s
	return c
}

func rshVU(z, x []big.Word, s uint) (c big.Word) {
	if len(z) == 0 {
		return 0
	}
	if len(x) != len(z) {
		panic("simd: rshVU len")
	}

	ŝ := WordBitSize - s

	c = x[0] << ŝ
	for i := range z {
		z[i] = x[i]>>s | x[i+1]<<ŝ
	}
	z[len(z)-1] = x[len(z)-1] >> s
	return c
}

func mulAddVWW(z, x []big.Word, m, a big.Word) (c big.Word) {
	c = a
	for i := range z {
		c, z[i] = mulAddWWW(x[i], m, c)
	}
	return c
}

func addMulVVWW(z, x, y []big.Word, m, a big.Word) (c big.Word) {
	c = a
	for i := range z {
		hi, lo := mulAddWWW(y[i], m, x[i])
		lo2, carry := bits.Add(uint(lo), uint(c), 0)
		c = hi + big.Word(carry)
		z[i] = big.Word(lo2)
	}
	return c
}
