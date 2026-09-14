// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package simd

import (
	"math/big"
	"math/bits"
)

func AddVV(z, x, y []big.Word) (c big.Word)                     { return addVV(z, x, y) }
func SubVV(z, x, y []big.Word) (c big.Word)                     { return subVV(z, x, y) }
func MulAddVWW(z, x []big.Word, m, a big.Word) (c big.Word)     { return mulAddVWW(z, x, m, a) }
func AddMulVVW(z, x []big.Word, m big.Word) (c big.Word)        { return addMulVVWW(z, z, x, m, 0) }
func AddMulVVWW(z, x, y []big.Word, m, a big.Word) (c big.Word) { return addMulVVWW(z, x, y, m, a) }

func ShlVU(z, x []big.Word, s uint) (c big.Word) {
	if s == 0 {
		copy(z, x)
		return 0
	}
	return lshVU(z, x, s)
}

func ShrVU(z, x []big.Word, s uint) (c big.Word) {
	if s == 0 {
		copy(z, x)
		return 0
	}
	return rshVU(z, x, s)
}

func AddVW(z, x []big.Word, y big.Word) (c big.Word) {
	if len(x) != len(z) {
		panic("simd: AddVW len")
	}
	if len(z) == 0 {
		return y
	}

	zi, cc := bits.Add(uint(x[0]), uint(y), 0)
	z[0] = big.Word(zi)
	if cc == 0 {
		if &z[0] != &x[0] {
			copy(z[1:], x[1:])
		}
		return 0
	}

	for i := 1; i < len(z); i++ {
		xi := x[i]
		if xi != WordMaxValue {
			z[i] = xi + 1
			if &z[0] != &x[0] {
				copy(z[i+1:], x[i+1:])
			}
			return 0
		}
		z[i] = 0
	}
	return 1
}

func SubVW(z, x []big.Word, y big.Word) (c big.Word) {
	if len(x) != len(z) {
		panic("simd: SubVW len")
	}
	if len(z) == 0 {
		return y
	}

	zi, cc := bits.Sub(uint(x[0]), uint(y), 0)
	z[0] = big.Word(zi)
	if cc == 0 {
		if &z[0] != &x[0] {
			copy(z[1:], x[1:])
		}
		return 0
	}

	for i := 1; i < len(z); i++ {
		xi := x[i]
		if xi != 0 {
			z[i] = xi - 1
			if &z[0] != &x[0] {
				copy(z[i+1:], x[i+1:])
			}
			return 0
		}
		z[i] = WordMaxValue
	}
	return 1
}

func MulWW(x, y big.Word) (hi, lo big.Word) {
	h, l := bits.Mul(uint(x), uint(y))
	return big.Word(h), big.Word(l)
}

func DivWW(hi, lo, y big.Word) (q, r big.Word) {
	qq, rr := bits.Div(uint(hi), uint(lo), uint(y))
	return big.Word(qq), big.Word(rr)
}

func Nlz(x big.Word) uint { return uint(bits.LeadingZeros(uint(x))) }

func Reciprocal(d big.Word) big.Word {
	rec, _ := bits.Div(^uint(d), ^uint(0), uint(d))
	return big.Word(rec)
}

func DivWWRecip(u1, u0, d, v big.Word) (q, r big.Word) {
	dd := uint(d)

	q1, q0 := bits.Mul(uint(v), uint(u1))
	q0, c := bits.Add(q0, uint(u0), 0)
	q1, _ = bits.Add(q1, uint(u1), c)
	q1++

	rr := uint(u0) - q1*dd
	if rr > q0 {
		q1--
		rr += dd
	}
	if rr >= dd {
		q1++
		rr -= dd
	}

	return big.Word(q1), big.Word(rr)
}

func mulAddWWW(x, y, c big.Word) (hi, lo big.Word) {
	h, l := bits.Mul(uint(x), uint(y))
	ll, cc := bits.Add(l, uint(c), 0)
	return big.Word(h + cc), big.Word(ll)
}
