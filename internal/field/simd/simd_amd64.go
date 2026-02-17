//go:build gc && amd64 && !purego

package simd

import (
	"math/big"

	"github.com/RyuaNerin/elliptic2/internal"
	"golang.org/x/sys/cpu"
)

//go:noescape
func _clmul(a, b big.Word) (lo, hi big.Word)

//go:noescape
func _clmulWords(z, x, y []big.Word)

//go:noescape
func _expandBitsBMI2(x big.Word) (lo, hi big.Word)

func init() {
	hasPCLMULQDQ := cpu.X86.HasPCLMULQDQ
	hasBMI2 := cpu.X86.HasBMI2

	if hasPCLMULQDQ {
		isCLMULAsmMode = true
		isCLMULWordsAsmMode = true

		CLMUL = _clmul
		CLMULWords = func(z, x, y []big.Word) {
			if len(z) < len(x)+len(y)-1 {
				panic("simd: output slice too small")
			}
			if internal.Overlaps(z, x) || internal.Overlaps(z, y) {
				panic("simd: output slice overlaps input")
			}
			clear(z)
			_clmulWords(z, x, y)
		}
	}

	if hasBMI2 {
		isExpandBitsAsmMode = true
		ExpandBits = _expandBitsBMI2
	}
}
