//go:build gc && (amd64 || 386 || arm64 || arm || riscv64 || riscv64 || ppc64 || ppc64le || s390x) && !purego

package simd

import (
	"bufio"
	"math/big"
	"math/rand"
	"testing"
)

var Random = bufio.NewReaderSize(rand.New(rand.NewSource(0)), 1<<15)

func initWord() (w big.Word) {
	for b := 0; b < WordByteSize; b++ {
		w |= big.Word(b) << (8 * b)
	}
	return w
}

func randomWord() (w big.Word) {
	var buf [WordBitSize / 8]byte
	Random.Read(buf[:])
	for b := 0; b < WordByteSize; b++ {
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

	for i := 0; i < 1_000; i++ {
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

		for idx := 0; idx < oLen; idx++ {
			if oGeneric[idx] != oAsm[idx] {
				t.Fatalf("mismatch at idx=%d:\ngot:  %016x\nwant: %016x",
					idx,
					oAsm[idx],
					oGeneric[idx],
				)
			}
		}
	}
}

func TestCLMUL(t *testing.T) {
	for i := 0; i < 10_000; i++ {
		lo := randomWord()
		hi := randomWord()

		wantLo, wantHi := clmulGeneric(lo, hi)
		gotLo, gotHi := CLMUL(lo, hi)

		if wantLo != gotLo || wantHi != gotHi {
			t.Fatalf("mismatch:\ngot:  %016x %016x\nwant: %016x %016x",
				gotLo, gotHi,
				wantLo, wantHi,
			)
			return
		}
	}
}

func TestExpandBits64(t *testing.T) {
	for i := 0; i < 10_000; i++ {
		x := randomWord()

		wantLo, wantHi := expandBitsGeneric(x)
		gotLo, gotHi := ExpandBits(x)

		if wantLo != gotLo || wantHi != gotHi {
			t.Fatalf("mismatch:\ngot: %016x %016x\nwant:  %016x %016x",
				gotLo, gotHi,
				wantLo, wantHi,
			)
			return
		}
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
