//go:build gc && arm64 && !purego && go1.18

package simd

import (
	"math/big"

	"golang.org/x/sys/cpu"
)

//go:noescape
func _clmul(x, y uint64) (lo, hi uint64)

func init() {
	// arm64 always has NEON
	hasPMULL := cpu.ARM64.HasPMULL

	if hasPMULL {
		CLMUL = func(a, b big.Word) (lo, hi big.Word) {
			l, h := _clmul(uint64(a), uint64(b))
			return big.Word(l), big.Word(h)
		}
	}
}
