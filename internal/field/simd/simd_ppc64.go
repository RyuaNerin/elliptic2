//go:build gc && (ppc64 || ppc64le) && !purego && go1.18

package simd

import "math/big"

//go:noescape
func _clmul(x, y uint64) (lo, hi uint64)

func init() {
	// power8+ always has VPMSUMD
	CLMUL = func(a, b big.Word) (lo, hi big.Word) {
		l, h := _clmul(uint64(a), uint64(b))
		return big.Word(l), big.Word(h)
	}
}
