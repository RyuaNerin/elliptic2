//go:build gc && (386 || amd64 || arm || arm64 || ppc64 || ppc64le || s390x) && !purego

package simd

import (
	"io"
	"math/big"
	"runtime"
	"testing"

	"github.com/RyuaNerin/elliptic2/internal"
	"github.com/stretchr/testify/require"
)

const preallocSize = 1024

func randomWords(t testing.TB, dst []big.Word, buf []byte) []byte {
	if cap(buf) < WordByteSize*len(dst) {
		buf = make([]byte, WordByteSize*len(dst))
	}

	_, err := io.ReadFull(internal.Random, buf[:len(dst)*WordByteSize])
	require.NoError(t, err)

	clear(dst)

	idxBuf := 0
	for idxDst := range dst {
		for b := range WordByteSize {
			dst[idxDst] |= big.Word(buf[idxBuf]) << (8 * b)
			idxBuf++
		}
	}

	return buf
}

func FuzzCLMULAsm(f *testing.F) {
	if !isCLMULAsmMode {
		f.Skip("CLMULAsm not available")
	}

	f.Add(uint64(0), uint64(0))
	f.Add(uint64(WordMaxValue), uint64(WordMaxValue))

	f.Fuzz(func(t *testing.T, a uint64, b uint64) {
		wa, wb := big.Word(a), big.Word(b)

		// Generic vs Asm 비교
		wantLo, wantHi := clmulGeneric(wa, wb)
		gotLo, gotHi := CLMUL(wa, wb)

		require.Equal(t, wantLo, gotLo, "CLMUL Lo mismatch")
		require.Equal(t, wantHi, gotHi, "CLMUL Hi mismatch")
	})
}

func FuzzCLMULWordsAsm(f *testing.F) {
	f.Add([]byte{1}, []byte{1})
	f.Add([]byte{0xFF, 0xFF, 0xFF, 0xFF}, []byte{0x01, 0x02})

	toWords := func(b []byte) []big.Word {
		if len(b) == 0 {
			return nil
		}
		n := (len(b) + WordByteSize - 1) / WordByteSize

		lst := make([]big.Word, n)

		for idx := range len(b) {
			wIdx := idx / WordByteSize
			shift := (idx % WordByteSize) * 8
			lst[wIdx] |= big.Word(b[idx]) << shift
		}
		return lst
	}

	f.Fuzz(func(t *testing.T, xb []byte, yb []byte) {
		if len(xb) == 0 || len(yb) == 0 {
			return
		}

		x := toWords(xb)
		y := toWords(yb)

		zLen := len(x) + len(y)

		zWant := make([]big.Word, zLen)
		zGot := make([]big.Word, zLen)

		clmulWordsGeneric(zWant, x, y)
		CLMULWords(zGot, x, y)

		require.Equal(t, zWant, zGot, "CLMULWords mismatch: x=%#v y=%#v", x, y)
	})
}

func FuzzExpandBits64Asm(f *testing.F) {
	if !isExpandBitsAsmMode {
		f.Skip("ExpandBitsAsm not available")
	}

	f.Add(uint64(0))
	f.Add(uint64(WordMaxValue))

	f.Fuzz(func(t *testing.T, x uint64) {
		wa := big.Word(x)

		wantLo, wantHi := expandBitsGeneric(wa)
		gotLo, gotHi := ExpandBits(wa)

		require.Equal(t, wantLo, gotLo, "ExpandBits Lo mismatch")
		require.Equal(t, wantHi, gotHi, "ExpandBits Hi mismatch")
	})
}

func BenchmarkCLMULAsm(b *testing.B) {
	if !isCLMULAsmMode {
		b.Skip("CLMULAsm not available")
	}

	bench := func(fn func(a, b big.Word) (big.Word, big.Word)) func(b *testing.B) {
		return func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			x := make([]big.Word, preallocSize)
			for idx := range x {
				x[idx] = randomWord(b)
			}

			var lo, hi big.Word
			for idx := range b.N {
				lo, hi = fn(x[idx%preallocSize], x[(idx+1)%preallocSize])
			}
			runtime.KeepAlive(lo)
			runtime.KeepAlive(hi)
		}
	}

	b.Run("Generic", bench(clmulGeneric))
	b.Run("Assembly", bench(CLMUL))
}

func BenchmarkCLMULWordsAsm(b *testing.B) {
	if !isCLMULWordsAsmMode {
		b.Skip("CLMULWordsAsm not available")
	}

	bench := func(fn func(z, x, y []big.Word)) func(b *testing.B) {
		return func(b *testing.B) {
			buf := make([]byte, Words*WordByteSize)
			x := make([][]big.Word, preallocSize)
			for idx := range x {
				x[idx] = make([]big.Word, Words)
				randomWords(b, x[idx], buf)
			}

			z := make([]big.Word, 2*Words)

			b.ReportAllocs()
			b.ResetTimer()

			for idx := range b.N {
				fn(z, x[idx%preallocSize], x[(idx+1)%preallocSize])
			}
			runtime.KeepAlive(z)
		}
	}

	b.Run("Generic", bench(clmulWordsGeneric))
	b.Run("Assembly", bench(CLMULWords))
}

func BenchmarkExpandBitsAsm(b *testing.B) {
	if !isExpandBitsAsmMode {
		b.Skip("ExpandBitsAsm not available")
	}

	bench := func(fn func(x big.Word) (big.Word, big.Word)) func(b *testing.B) {
		return func(b *testing.B) {
			x := make([]big.Word, preallocSize)
			for idx := range x {
				x[idx] = randomWord(b)
			}

			b.ReportAllocs()
			b.ResetTimer()

			var lo, hi big.Word
			for idx := range b.N {
				lo, hi = fn(x[idx%preallocSize])
			}
			runtime.KeepAlive(lo)
			runtime.KeepAlive(hi)
		}
	}

	b.Run("Generic", bench(expandBitsGeneric))
	b.Run("Assembly", bench(ExpandBits))
}
