//go:build gc && riscv64 && !purego && go1.18

package simd

import (
	"math/big"

	"golang.org/x/sys/cpu"
)

//go:noescape
func _clmul(x, y uint64) (lo, hi uint64)

func init() {
	hasV := cpu.RISCV64.HasV

	if hasV {
		CLMUL = func(a, b big.Word) (lo, hi big.Word) {
			l, h := _clmul(uint64(a), uint64(b))
			return big.Word(l), big.Word(h)
		}
	}
}
