package field

import (
	"math/big"

	"github.com/RyuaNerin/elliptic2/internal/field/simd"
)

type nat struct {
	words [1 + simd.Words]big.Word
}

func (z *nat) toBigInt(x *big.Int) *big.Int {
	if x == nil {
		x = new(big.Int)
	}
	words := trimLeadingZeros(z.words[:])
	b := make([]big.Word, len(words))
	copy(b, words)
	return x.SetBits(b)
}

/*
- return -1 if z < x

- return 0  if z == x

- return 1  if z > x
*/
func (z *nat) cmp(x *nat) int {
	for idx := range len(z.words) {
		off := len(z.words) - 1 - idx
		zw, xw := z.words[off], x.words[off]
		if zw < xw {
			return -1
		}
		if zw > xw {
			return 1
		}
	}
	return 0
}

func (z *nat) isZero() bool {
	var w big.Word
	for idx := range len(z.words) {
		w |= z.words[idx]
	}
	return w == 0
}

func (z *nat) set(x *nat) *nat {
	copy(z.words[:], x.words[:])
	return z
}

func (z *nat) setZero() {
	clear(z.words[:])
}

func (z *nat) setBits(xWords []big.Word) {
	if len(xWords) > len(z.words) {
		panic("field: input too large")
	}

	// trim leading zeros
	xWords = trimLeadingZeros(xWords)

	n := copy(z.words[:], xWords)
	clear(z.words[n:])
}

func (z *nat) setBigInt(x *big.Int) {
	z.setBits(x.Bits())
}

func (z *nat) setWord(x big.Word) *nat {
	clear(z.words[:])
	z.words[0] = x
	return z
}

func (z *nat) setUint64(x uint64) (zz *nat) {
	// single-word value
	if w := big.Word(x); uint64(w) == x {
		z.setWord(w)
		return z
	}

	// 2-word value
	for idx := range z.words {
		z.words[idx] = 0
	}
	z.words[1] = big.Word(x >> 32)
	z.words[0] = big.Word(x)
	return z
}

func (z *nat) setBytes(b []byte) {
	for idx := range z.words {
		z.words[idx] = 0
	}

	for idx := len(b) - 1; idx >= 0; idx-- {
		byteIdx := len(b) - 1 - idx
		wordIdx := byteIdx / simd.WordByteSize
		byteOffset := (byteIdx % simd.WordByteSize) * 8
		if wordIdx < len(z.words) {
			z.words[wordIdx] |= big.Word(b[idx]) << byteOffset
		}
	}
}

func trimLeadingZeros(w []big.Word) []big.Word {
	for len(w) > 0 && w[len(w)-1] == 0 {
		w = w[:len(w)-1]
	}
	return w
}
