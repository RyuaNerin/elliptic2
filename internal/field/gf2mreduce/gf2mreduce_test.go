package gf2mreduce

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/RyuaNerin/elliptic2/internal/field/simd"
)

func TestHas(t *testing.T) {
	input := []int{571, 10, 5, 2, 0}
	const want uint32 = 0x002486f5

	h := hash(input)
	if h != want {
		t.Errorf("failed hash function: got 0x%08x, want 0x%08x", h, want)
	}
}

func TestReduce(t *testing.T) {
	for _, poly := range reduces {
		var input [2 * simd.Words]big.Word
		var outputOptimized [simd.Words]big.Word
		var outputGeneric [simd.Words]big.Word

		rnd := rand.New(rand.NewSource(0))
		buf := make([]byte, len(input)*simd.WordByteSize)
		_, _ = rnd.Read(buf)
		for idx := range len(input) {
			for b := range simd.WordByteSize {
				input[idx] |= big.Word(buf[0]) << (8 * b)
				buf = buf[1:]
			}
		}

		fnOptimized := GetReduceFunction(poly)
		fnGeneric := newReduce(poly)

		fnOptimized(outputOptimized[:], input[:])
		fnGeneric(outputGeneric[:], input[:])

		for idx := range len(outputOptimized) {
			if outputOptimized[idx] != outputGeneric[idx] {
				t.Errorf("Reduce failed for poly=%v at word index %d: got %x, want %x",
					poly, idx, outputOptimized[idx], outputGeneric[idx])
				return
			}
		}
	}
}
