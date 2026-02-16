//go:build gc && (amd64 || 386 || arm64 || arm || riscv64 || ppc64 || ppc64le || s390x) && !purego

package simd

import (
	"io"
	"math/big"
	"testing"

	"github.com/RyuaNerin/elliptic2/internal"
	"github.com/stretchr/testify/require"
)

func initWord() (w big.Word) {
	for b := range WordByteSize {
		w |= big.Word(b) << (8 * b)
	}
	return w
}

func randomWord(t testing.TB) (w big.Word) {
	var buf [WordBitSize / 8]byte
	_, err := io.ReadFull(internal.Random, buf[:])
	require.NoError(t, err)

	for b := range WordByteSize {
		w |= big.Word(buf[b]) << (8 * b)
	}
	return w
}

func randomWords(t testing.TB, dst []big.Word, words int, buf []byte) []byte {
	if cap(buf) < 8*len(dst) {
		buf = make([]byte, 8*len(dst))
	}

	_, err := io.ReadFull(internal.Random, buf[:words*WordBitSize/8])
	require.NoError(t, err)

	idx := 0
	for j := range words {
		dst[j] = 0
		for b := range WordByteSize {
			dst[j] |= big.Word(buf[idx]) << (8 * b)
			idx++
		}
	}

	return buf
}

func TestCLMULWords(t *testing.T) {
	buf := make([]byte, 128)

	oGeneric := make([]big.Word, Words*2)
	oAsm := make([]big.Word, Words*2)
	xWords := make([]big.Word, Words)
	yWords := make([]big.Word, Words)

	for range 1_000 {
		_, err := io.ReadFull(internal.Random, buf[:2])
		require.NoError(t, err)

		xWordsLen := int(buf[0]) % cap(xWords)
		yWordsLen := int(buf[1]) % cap(yWords)
		oLen := xWordsLen + yWordsLen

		buf = randomWords(t, xWords, xWordsLen, buf)
		buf = randomWords(t, yWords, yWordsLen, buf)

		for idx := range oGeneric {
			oGeneric[idx] = 0
		}
		for idx := range oAsm {
			oAsm[idx] = 0
		}

		clmulWordsGeneric(oGeneric[:oLen], xWords[:xWordsLen], yWords[:yWordsLen])
		CLMULWords(oAsm[:oLen], xWords[:xWordsLen], yWords[:yWordsLen])

		for idx := range oLen {
			if oGeneric[idx] != oAsm[idx] {
				require.Equal(t, oGeneric[idx], oAsm[idx], "CLMULWords mismatch at idx=%d", idx)
			}
		}
	}
}

func TestCLMUL(t *testing.T) {
	for range 10_000 {
		lo := randomWord(t)
		hi := randomWord(t)

		wantLo, wantHi := clmulGeneric(lo, hi)
		gotLo, gotHi := CLMUL(lo, hi)

		require.Equal(t, wantLo, gotLo, "CLMUL Lo mismatch")
		require.Equal(t, wantHi, gotHi, "CLMUL Hi mismatch")
	}
}

func TestExpandBits64(t *testing.T) {
	for range 10_000 {
		x := randomWord(t)

		wantLo, wantHi := expandBitsGeneric(x)
		gotLo, gotHi := ExpandBits(x)

		require.Equal(t, wantLo, gotLo, "ExpandBits Lo mismatch")
		require.Equal(t, wantHi, gotHi, "ExpandBits Hi mismatch")
	}
}

func BenchmarkCLMUL(b *testing.B) {
	bench := func(fn func(a, b big.Word) (big.Word, big.Word)) func(b *testing.B) {
		return func(b *testing.B) {
			x := initWord()
			y := initWord()

			b.SetBytes(WordBitSize / 8)
			b.ReportAllocs()
			b.ResetTimer()

			for idx := range b.N {
				lo, hi := fn(x, y)

				x = x ^ hi ^ big.Word(idx)
				y = y ^ lo ^ big.Word(idx)
			}
		}
	}

	b.Run("Generic", bench(clmulGeneric))
	b.Run("Assembly", bench(CLMUL))
}

func BenchmarkExpandBits(b *testing.B) {
	bench := func(fn func(x big.Word) (big.Word, big.Word)) func(b *testing.B) {
		return func(b *testing.B) {
			x := initWord()

			b.SetBytes(8)
			b.ReportAllocs()
			b.ResetTimer()

			for idx := range b.N {
				lo, hi := fn(x)

				x = x ^ hi ^ lo ^ big.Word(idx)
			}
		}
	}

	b.Run("Generic", bench(expandBitsGeneric))
	b.Run("Assembly", bench(ExpandBits))
}
