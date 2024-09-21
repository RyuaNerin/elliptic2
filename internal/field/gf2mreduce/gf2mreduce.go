package gf2mreduce

import (
	"math/big"

	"github.com/RyuaNerin/elliptic2/internal/field/simd"
)

type ReduceFunc func(result []big.Word, prod []big.Word)

var (
	funcMap = map[uint32]ReduceFunc{}
	reduces [][]int
)

func GetReduceFunction(polyBits []int) ReduceFunc {
	f, ok := funcMap[hash(polyBits)]
	if !ok {
		return newReduce(polyBits)
	}
	return f
}

func hash(polyBits []int) uint32 {
	var fnvPrime uint32 = 0x01000193
	var fnvOffset uint32 = 0x811C9DC5

	h := fnvOffset
	for _, b := range polyBits {
		if b < 256 {
			h ^= uint32(b)
			h *= fnvPrime
		} else if b < 65536 {
			hi := (b >> 8) & 0xFF
			lo := b & 0xFF

			h ^= uint32(hi)
			h *= fnvPrime
			h ^= uint32(lo)
			h *= fnvPrime
		}
	}

	return h
}

func newReduce(poly []int) func(dst []big.Word, x []big.Word) {
	n := (poly[0] + simd.WordBitSize - 1) / simd.WordBitSize
	polyBits := poly[1:]

	return func(dst []big.Word, x []big.Word) {
		m := poly[0]

		var work [2 * simd.Words]big.Word
		copy(work[:2*n], x[:])

		mWordIdx := m / simd.WordBitSize
		mBitIdx := uint(m % simd.WordBitSize)

		for idx := 2*n - 1; idx > mWordIdx; idx-- {
			w := work[idx]
			if w == 0 {
				continue
			}
			work[idx] = 0

			baseShift := idx*simd.WordBitSize - m

			for _, polyBit := range polyBits {
				totalShift := baseShift + polyBit
				wordPos := totalShift / simd.WordBitSize
				bitShift := uint(totalShift % simd.WordBitSize)

				work[wordPos] ^= w << bitShift
				if bitShift > 0 && wordPos+1 < 2*n {
					work[wordPos+1] ^= w >> (simd.WordBitSize - bitShift)
				}
			}
		}

		if mBitIdx > 0 {
			mask := ^big.Word(0) << mBitIdx
			w := work[mWordIdx] & mask

			if w != 0 {
				work[mWordIdx] ^= w
				shiftedW := w >> mBitIdx
				for _, polyBit := range polyBits {
					wordPos := polyBit / simd.WordBitSize
					bitShift := uint(polyBit % simd.WordBitSize)

					work[wordPos] ^= shiftedW << bitShift
					if bitShift > 0 && wordPos+1 <= mWordIdx {
						work[wordPos+1] ^= shiftedW >> (simd.WordBitSize - bitShift)
					}
				}
			}
		}

		copy(dst[:n], work[:n])
	}
}
