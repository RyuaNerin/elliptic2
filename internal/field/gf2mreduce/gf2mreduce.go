package gf2mreduce

import (
	"fmt"
	"math/big"
	"slices"

	"github.com/RyuaNerin/elliptic2/internal/field/simd"
)

type ReduceFunc func(result []big.Word, prod []big.Word)

func GetReduceFunction(polyBits []int) ReduceFunc {
	if f := lookupReduce(polyBits); f != nil {
		return f
	}
	return newGenericReduce(polyBits)
}

func newGenericReduce(poly []int) ReduceFunc {
	if len(poly) < 2 {
		panic(fmt.Sprintf("gf2mreduce: polynomial needs at least 2 terms: %v", poly))
	}
	for i := 1; i < len(poly); i++ {
		if poly[i] >= poly[i-1] {
			panic(fmt.Sprintf("gf2mreduce: polynomial must be strictly descending: %v", poly))
		}
	}
	if poly[len(poly)-1] < 0 {
		panic(fmt.Sprintf("gf2mreduce: negative exponent: %v", poly))
	}

	const W = simd.WordBitSize

	m := poly[0]
	polyBits := slices.Clone(poly[1:])
	n := (m + W - 1) / W
	mWordIdx := m / W
	mBitIdx := uint(m % W)

	return func(dst []big.Word, x []big.Word) {
		_ = dst[n-1]
		_ = x[2*n-1]

		var buf [2 * simd.Words]big.Word
		var work []big.Word
		if 2*n <= len(buf) {
			work = buf[:2*n]
		} else {
			work = make([]big.Word, 2*n)
		}
		copy(work, x[:2*n])

		for idx := 2*n - 1; idx >= mWordIdx; idx-- {
			var s uint
			if idx == mWordIdx {
				s = mBitIdx
			}
			base := idx*W + int(s) - m

			for {
				t := work[idx] >> s
				if t == 0 {
					break
				}
				work[idx] &^= t << s

				for _, p := range polyBits {
					pos := base + p
					wordPos := pos / W
					bitShift := uint(pos % W)

					work[wordPos] ^= t << bitShift
					if bitShift > 0 {
						if hi := t >> (W - bitShift); hi != 0 {
							work[wordPos+1] ^= hi
						}
					}
				}
			}
		}

		copy(dst, work[:n])
		clear(dst[n:])
	}
}
