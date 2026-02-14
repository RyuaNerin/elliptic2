//go:build gc && arm && !purego

package simd

import "math/big"

//go:noescape
func _clmul32(a, b uint32) uint64

func init() {
	CLMUL = clmul32Wrapper
}

func clmul32Wrapper(a, b big.Word) (lo, hi big.Word) {
	result := _clmul32(uint32(a), uint32(b))
	lo = big.Word(result)
	hi = big.Word(result >> 32)
	return
}
