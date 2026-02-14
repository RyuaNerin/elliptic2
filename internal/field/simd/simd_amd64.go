//go:build gc && amd64 && !purego

package simd

import (
	"math/big"

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
		CLMUL = _clmul
		CLMULWords = _clmulWords
	}

	if hasBMI2 {
		ExpandBits = _expandBitsBMI2
	}
}
