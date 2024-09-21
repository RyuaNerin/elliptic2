//go:build gc && 386 && !purego && go1.18

package simd

import (
	"math/big"

	"golang.org/x/sys/cpu"
)

//go:noescape
func _clmul32(a, b uint32) (lo, hi uint32)

func init() {
	hasPCLMULQDQ := cpu.X86.HasPCLMULQDQ

	if hasPCLMULQDQ {
		CLMUL = func(a, b big.Word) (lo, hi big.Word) {
			l, h := _clmul32(uint32(a), uint32(b))
			return big.Word(l), big.Word(h)
		}
	}
}
