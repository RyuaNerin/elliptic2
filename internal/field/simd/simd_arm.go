//go:build gc && arm && !purego && go1.18

package simd

import "math/big"

//go:noescape
func _clmul32(a, b uint32) uint64

func init() {
	// ARM 32-bit: Word is 32-bit
	// Use 32x32->64 carry-less multiplication
	CLMUL = clmul32Wrapper
}

// clmul32Wrapper wraps the 32-bit clmul for the CLMUL interface
func clmul32Wrapper(a, b big.Word) (lo, hi big.Word) {
	result := _clmul32(uint32(a), uint32(b))
	lo = big.Word(result)
	hi = big.Word(result >> 32)
	return
}

// clmul64Via32 performs 64-bit carry-less multiplication using 32-bit operations
// This is useful if someone needs 64-bit clmul on 32-bit ARM
func clmul64Via32(a, b uint64) (lo, hi uint64) {
	a0, a1 := uint32(a), uint32(a>>32)
	b0, b1 := uint32(b), uint32(b>>32)

	z0 := _clmul32(a0, b0)
	z2 := _clmul32(a1, b1)
	z1a := _clmul32(a0, b1)
	z1b := _clmul32(a1, b0)
	z1 := z1a ^ z1b

	lo = z0 ^ (z1 << 32)
	hi = z2 ^ (z1 >> 32)
	return
}
