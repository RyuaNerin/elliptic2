//go:build gc && (amd64 || 386 || arm64 || arm || riscv64 || riscv64 || ppc64 || ppc64le || s390x) && !purego

package simd

import (
	"bufio"
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

var Random = bufio.NewReaderSize(rand.New(rand.NewSource(0)), 1<<15)

func initWord() (w big.Word) {
	for b := 0; b < WordByteSize; b++ {
		w |= big.Word(b) << (8 * b)
	}
	return w
}

func randomWord(t testing.TB) (w big.Word) {
	var buf [WordBitSize / 8]byte
	_, err := Random.Read(buf[:])
	require.NoError(t, err)

	for b := range WordByteSize {
		w |= big.Word(buf[b]) << (8 * b)
	}
	return w
}

func randomWords(dst []big.Word, words int, buf []byte) []byte {
	if cap(buf) < 8*len(dst) {
		buf = make([]byte, 8*len(dst))
	}

	Random.Read(buf[:words*WordBitSize/8])
	idx := 0
	for j := 0; j < words; j++ {
		dst[j] = 0
		for b := 0; b < WordByteSize; b++ {
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
		Random.Read(buf[:2])
		xWordsLen := int(buf[0]) % cap(xWords)
		yWordsLen := int(buf[1]) % cap(yWords)
		oLen := xWordsLen + yWordsLen

		buf = randomWords(xWords, xWordsLen, buf)
		buf = randomWords(yWords, yWordsLen, buf)

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
				require.Equal(t, oAsm[idx], oGeneric[idx], "CLMULWords mismatch at idx=%d", idx)
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

		require.True(t, wantLo == gotLo, "CLMUL Lo mismatch:\ngot:  %016x\nwant: %016x", gotLo, wantLo)
		require.True(t, wantHi == gotHi, "CLMUL Hi mismatch:\ngot:  %016x\nwant: %016x", gotHi, wantHi)
	}
}

func TestExpandBits64(t *testing.T) {
	for range 10_000 {
		x := randomWord(t)

		wantLo, wantHi := expandBitsGeneric(x)
		gotLo, gotHi := ExpandBits(x)

		require.True(t, wantLo == gotLo, "ExpandBits Lo mismatch:\ngot:  %016x\nwant: %016x", gotLo, wantLo)
		require.True(t, wantHi == gotHi, "ExpandBits Hi mismatch:\ngot:  %016x\nwant: %016x", gotHi, wantHi)
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

			for idx := 0; idx < b.N; idx++ {
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

			for idx := 0; idx < b.N; idx++ {
				lo, hi := fn(x)

				x = x ^ hi ^ lo ^ big.Word(idx)
			}
		}
	}

	b.Run("Generic", bench(expandBitsGeneric))
	b.Run("Assembly", bench(ExpandBits))
}
