package gf2mreduce

import (
	"fmt"
	"io"
	"math/big"
	"testing"

	"github.com/RyuaNerin/elliptic2/internal/curvetesting/testrand"
	"github.com/RyuaNerin/elliptic2/internal/field/simd"

	"github.com/stretchr/testify/require"
)

func BenchmarkReduce(b *testing.B) {
	bench := func(fn func(output, input []big.Word)) func(b *testing.B) {
		return func(b *testing.B) {
			var input [2 * simd.Words]big.Word
			var output [simd.Words]big.Word
			buf := make([]byte, len(input)*simd.WordByteSize)

			fill(b, testrand.Random, buf, input[:])

			b.ResetTimer()
			for range b.N {
				fn(output[:], input[:])
			}
		}
	}

	for _, poly := range reduces {
		b.Run(fmt.Sprintf("Optimized/poly=%v", poly), bench(GetReduceFunction(poly)))
		b.Run(fmt.Sprintf("Generic/poly=%v", poly), bench(newGenericReduce(poly)))
	}
}

func fill(t testing.TB, r io.Reader, buf []byte, input []big.Word) {
	_, err := r.Read(buf)
	require.NoError(t, err)

	clear(input)

	bufIdx := 0
	for idx := range input {
		for idxByte := range simd.WordByteSize {
			input[idx] |= big.Word(buf[bufIdx]) << (8 * idxByte)
			bufIdx++
		}
	}
}
