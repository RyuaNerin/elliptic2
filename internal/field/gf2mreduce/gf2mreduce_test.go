package gf2mreduce

import (
	"fmt"
	"io"
	"math/big"
	"testing"

	"github.com/RyuaNerin/elliptic2/internal"
	"github.com/RyuaNerin/elliptic2/internal/field/simd"
	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	input := []int{571, 10, 5, 2, 0}
	const want uint32 = 0x002486f5

	h := hash(input)
	require.Equal(t, want, h, "Hash mismatch for input=%v", input)
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

func TestReduce(t *testing.T) {
	for _, poly := range reduces {
		t.Run(fmt.Sprintf("poly=%v", poly), func(t *testing.T) {
			fnOptimized := GetReduceFunction(poly)
			fnGeneric := newGenericReduce(poly)

			test := func(input []big.Word) {
				var outputGeneric [simd.Words]big.Word
				var outputOptimized [simd.Words]big.Word

				fnGeneric(outputGeneric[:], input)
				fnOptimized(outputOptimized[:], input)

				require.Equal(t, outputGeneric[:], outputOptimized[:], "Reduce failed for poly=%v", poly)
			}

			testReduce(test)
			testReduceRandom(t, test)
		})
	}
}

func testReduce(test func(input []big.Word)) {
	var input [2 * simd.Words]big.Word

	// Zero
	clear(input[:])
	test(input[:])

	// Max
	for idx := range input {
		input[idx] = ^big.Word(0)
	}
	test(input[:])

	// Single bit
	for idx := range input {
		input[idx] = 1 << (idx % (8 * simd.WordByteSize))
	}
	test(input[:])

	// 0 | max
	for idx := range len(input) / 2 {
		input[idx] = 0
	}
	for idx := len(input) / 2; idx < len(input); idx++ {
		input[idx] = ^big.Word(0)
	}
	test(input[:])

	// max | 0
	for idx := range len(input) / 2 {
		input[idx] = ^big.Word(0)
	}
	for idx := len(input) / 2; idx < len(input); idx++ {
		input[idx] = 0
	}
	test(input[:])
}

func testReduceRandom(t *testing.T, test func(input []big.Word)) {
	var input [2 * simd.Words]big.Word
	buf := make([]byte, len(input)*simd.WordByteSize)

	for range 1_000 {
		fill(t, internal.Random, buf, input[:])
		test(input[:])
	}
}

func BenchmarkReduce(b *testing.B) {
	bench := func(fn func(output, input []big.Word)) func(b *testing.B) {
		return func(b *testing.B) {
			var input [2 * simd.Words]big.Word
			var output [simd.Words]big.Word
			buf := make([]byte, len(input)*simd.WordByteSize)

			fill(b, internal.Random, buf, input[:])

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
