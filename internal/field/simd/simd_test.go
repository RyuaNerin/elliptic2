package simd

import (
	"io"
	"math/big"
	"math/bits"
	"runtime"
	"testing"

	"github.com/RyuaNerin/elliptic2/internal"
	"github.com/stretchr/testify/require"
)

const preallocSize = 1024

func randomWord(t testing.TB) (w big.Word) {
	var buf [WordByteSize]byte
	_, err := io.ReadFull(internal.Random, buf[:])
	require.NoError(t, err)

	for i := range WordByteSize {
		w |= big.Word(buf[i]) << (8 * i)
	}
	return w
}

func randomWordSlice(t testing.TB, n int) []big.Word {
	buf := make([]byte, WordByteSize*n)
	_, err := io.ReadFull(internal.Random, buf)
	require.NoError(t, err)

	ws := make([]big.Word, n)
	for i := range n {
		for b := range WordByteSize {
			ws[i] |= big.Word(buf[i*WordByteSize+b]) << (8 * b)
		}
	}
	return ws
}

func TestCLMUL(t *testing.T) {
	t.Run("known_values", func(t *testing.T) {
		cases := []struct {
			a, b, lo, hi big.Word
		}{
			{0x0000, 0x0000, 0x0000, 0x0000},
			{0x1234, 0x0001, 0x1234, 0x0000},
			{0x0001, 0x1234, 0x1234, 0x0000},
			{0b11, 0b11, 0b101, 0},
			{0b1111, 0b1111, 0b1010101, 0},
		}
		for _, tc := range cases {
			lo, hi := clmulGeneric(tc.a, tc.b)
			require.Equal(t, tc.lo, lo, "lo: a=%#x b=%#x", tc.a, tc.b)
			require.Equal(t, tc.hi, hi, "hi: a=%#x b=%#x", tc.a, tc.b)
		}
	})

	t.Run("identity_and_zero", func(t *testing.T) {
		for range 10_000 {
			a := randomWord(t)

			lo, hi := clmulGeneric(a, 0)
			require.Zero(t, lo)
			require.Zero(t, hi)

			lo, hi = clmulGeneric(0, a)
			require.Zero(t, lo)
			require.Zero(t, hi)

			lo, hi = clmulGeneric(a, 1)
			require.Equal(t, a, lo)
			require.Zero(t, hi)

			lo, hi = clmulGeneric(1, a)
			require.Equal(t, a, lo)
			require.Zero(t, hi)
		}
	})

	t.Run("commutativity", func(t *testing.T) {
		for range 10_000 {
			a := randomWord(t)
			b := randomWord(t)

			lo1, hi1 := clmulGeneric(a, b)
			lo2, hi2 := clmulGeneric(b, a)
			require.Equal(t, lo1, lo2, "a=%#x b=%#x", a, b)
			require.Equal(t, hi1, hi2, "a=%#x b=%#x", a, b)
		}
	})

	t.Run("distributivity", func(t *testing.T) {
		for range 10_000 {
			a := randomWord(t)
			b := randomWord(t)
			c := randomWord(t)

			// a * (b ^ c) == (a*b) ^ (a*c)
			lo1, hi1 := clmulGeneric(a, b^c)
			lob, hib := clmulGeneric(a, b)
			loc, hic := clmulGeneric(a, c)
			require.Equal(t, lob^loc, lo1)
			require.Equal(t, hib^hic, hi1)

			// (a ^ b) * c == (a*c) ^ (b*c)
			lo2, hi2 := clmulGeneric(a^b, c)
			loa, hia := clmulGeneric(a, c)
			lob2, hib2 := clmulGeneric(b, c)
			require.Equal(t, loa^lob2, lo2)
			require.Equal(t, hia^hib2, hi2)
		}
	})

	t.Run("single_bit_shift", func(t *testing.T) {
		for range 5_000 {
			a := randomWord(t)
			k := int(randomWord(t) % big.Word(WordBitSize))
			b := big.Word(1) << k

			lo, hi := clmulGeneric(a, b)

			loWant := a << k
			hiWant := big.Word(0)
			if k != 0 {
				hiWant = a >> (WordBitSize - k)
			}
			require.Equal(t, loWant, lo, "a=%#x k=%d", a, k)
			require.Equal(t, hiWant, hi, "a=%#x k=%d", a, k)
		}
	})

	t.Run("msb_carry", func(t *testing.T) {
		a := big.Word(1) << (WordBitSize - 1)
		b := big.Word(2)
		lo, hi := clmulGeneric(a, b)
		require.Zero(t, lo)
		require.Equal(t, big.Word(1), hi)
	})

	t.Run("squaring_no_odd_bits", func(t *testing.T) {
		var oddMask uint64
		if WordBitSize == 64 {
			oddMask = 0xAAAAAAAAAAAAAAAA
		} else {
			oddMask = 0xAAAAAAAA
		}

		for range 10_000 {
			a := randomWord(t)
			lo, hi := clmulGeneric(a, a)
			require.Zero(t, uint64(lo)&oddMask, "a=%#x", a)
			require.Zero(t, uint64(hi)&oddMask, "a=%#x", a)
		}
	})

	t.Run("squaring_matches_expand", func(t *testing.T) {
		for range 10_000 {
			a := randomWord(t)
			sqLo, sqHi := clmulGeneric(a, a)
			exLo, exHi := expandBitsGeneric(a)
			require.Equal(t, exLo, sqLo, "a=%#x", a)
			require.Equal(t, exHi, sqHi, "a=%#x", a)
		}
	})
}

func TestCLMULWordsGeneric(t *testing.T) {
	t.Run("known_values", func(t *testing.T) {
		cases := []struct {
			x, y, z []big.Word
		}{
			{[]big.Word{1, 1}, []big.Word{1, 1}, []big.Word{1, 0, 1, 0}},
			{[]big.Word{3}, []big.Word{2}, []big.Word{6, 0}},
		}
		for _, tc := range cases {
			z := make([]big.Word, len(tc.z))
			clmulWordsGeneric(z, tc.x, tc.y)
			require.Equal(t, tc.z, z)
		}
	})

	t.Run("matches_single_word_clmul", func(t *testing.T) {
		for range 5_000 {
			a := randomWord(t)
			b := randomWord(t)

			z := make([]big.Word, 2)
			clmulWordsGeneric(z, []big.Word{a}, []big.Word{b})

			lo, hi := clmulGeneric(a, b)
			require.Equal(t, lo, z[0])
			require.Equal(t, hi, z[1])
		}
	})

	t.Run("distributivity", func(t *testing.T) {
		for range 3_000 {
			xLen := 1 + int(randomWord(t)%4)
			yLen := 1 + int(randomWord(t)%4)
			zLen := xLen + yLen

			x := randomWordSlice(t, xLen)
			y1 := randomWordSlice(t, yLen)
			y2 := randomWordSlice(t, yLen)

			y12 := make([]big.Word, yLen)
			for i := range yLen {
				y12[i] = y1[i] ^ y2[i]
			}

			z1 := make([]big.Word, zLen)
			z2 := make([]big.Word, zLen)
			z12 := make([]big.Word, zLen)

			clmulWordsGeneric(z1, x, y1)
			clmulWordsGeneric(z2, x, y2)
			clmulWordsGeneric(z12, x, y12)

			for i := range zLen {
				require.Equal(t, z1[i]^z2[i], z12[i], "idx=%d", i)
			}
		}
	})

	t.Run("commutativity", func(t *testing.T) {
		for range 3_000 {
			xLen := 1 + int(randomWord(t)%4)
			yLen := 1 + int(randomWord(t)%4)

			x := randomWordSlice(t, xLen)
			y := randomWordSlice(t, yLen)

			zxy := make([]big.Word, xLen+yLen)
			zyx := make([]big.Word, xLen+yLen)

			clmulWordsGeneric(zxy, x, y)
			clmulWordsGeneric(zyx, y, x)

			require.Equal(t, zxy, zyx)
		}
	})
}

func TestExpandBitsGeneric(t *testing.T) {
	t.Run("known_values", func(t *testing.T) {
		lo, hi := expandBitsGeneric(0)
		require.Zero(t, lo)
		require.Zero(t, hi)

		lo, hi = expandBitsGeneric(1)
		require.Equal(t, big.Word(1), lo)
		require.Zero(t, hi)

		lo, hi = expandBitsGeneric(0b_1111_1111)
		require.Equal(t, big.Word(0b_0101_0101_0101_0101), lo)
		require.Zero(t, hi)

		var evenMask uint64
		if WordBitSize == 64 {
			evenMask = 0x5555555555555555
		} else {
			evenMask = 0x55555555
		}
		lo, hi = expandBitsGeneric(WordMaxValue)
		require.Equal(t, evenMask, uint64(lo))
		require.Equal(t, evenMask, uint64(hi))
	})

	t.Run("single_bit", func(t *testing.T) {
		for b := range WordBitSize {
			x := big.Word(1) << b
			lo, hi := expandBitsGeneric(x)

			pos := 2 * b
			if pos < WordBitSize {
				require.Equal(t, big.Word(1)<<pos, lo, "b=%d", b)
				require.Zero(t, hi, "b=%d", b)
			} else {
				require.Zero(t, lo, "b=%d", b)
				require.Equal(t, big.Word(1)<<(pos-WordBitSize), hi, "b=%d", b)
			}
		}
	})

	t.Run("xor_homomorphism", func(t *testing.T) {
		var oddMask uint64
		if WordBitSize == 64 {
			oddMask = 0xAAAAAAAAAAAAAAAA
		} else {
			oddMask = 0xAAAAAAAA
		}

		for range 20_000 {
			a := randomWord(t)
			b := randomWord(t)

			alo, ahi := expandBitsGeneric(a)
			blo, bhi := expandBitsGeneric(b)
			xlo, xhi := expandBitsGeneric(a ^ b)

			require.Equal(t, alo^blo, xlo)
			require.Equal(t, ahi^bhi, xhi)

			require.Zero(t, uint64(xlo)&oddMask)
			require.Zero(t, uint64(xhi)&oddMask)

			pcIn := bits.OnesCount(uint(a))
			pcOut := bits.OnesCount(uint(alo)) + bits.OnesCount(uint(ahi))
			require.Equal(t, pcIn, pcOut)
		}
	})
}

func TestCompressBitsGeneric(t *testing.T) {
	t.Run("known_values", func(t *testing.T) {
		e, o := compressBitsGeneric(0)
		require.Zero(t, e)
		require.Zero(t, o)

		e, o = compressBitsGeneric(1)
		require.Equal(t, big.Word(1), e)
		require.Zero(t, o)

		e, o = compressBitsGeneric(2)
		require.Zero(t, e)
		require.Equal(t, big.Word(1), o)

		e, o = compressBitsGeneric(3)
		require.Equal(t, big.Word(1), e)
		require.Equal(t, big.Word(1), o)

		e, o = compressBitsGeneric(big.Word(uint64(0x5555555555555555) & uint64(WordMaxValue)))
		require.Equal(t, big.Word(0xFFFFFFFF), e)
		require.Zero(t, o)

		e, o = compressBitsGeneric(big.Word(uint64(0xAAAAAAAAAAAAAAAA) & uint64(WordMaxValue)))
		require.Zero(t, e)
		require.Equal(t, big.Word(0xFFFFFFFF), o)
	})

	t.Run("single_bit", func(t *testing.T) {
		for i := range WordBitSize {
			x := big.Word(1) << i
			e, o := compressBitsGeneric(x)

			targetBit := big.Word(1) << (i / 2)

			if i%2 == 0 {
				require.Equal(t, targetBit, e, "bit %d should map to even %d", i, i/2)
				require.Zero(t, o, "bit %d should not map to odd", i)
			} else {
				require.Zero(t, e, "bit %d should not map to even", i)
				require.Equal(t, targetBit, o, "bit %d should map to odd %d", i, i/2)
			}
		}
	})

	t.Run("round_trip_inverse", func(t *testing.T) {
		for range 20_000 {
			w := uint32(randomWord(t))

			expandedLo, _ := expandBitsGeneric(big.Word(w))
			expandedVal := uint64(expandedLo)

			e, o := compressBitsGeneric(big.Word(expandedVal))
			require.Equal(t, big.Word(w), e, "Compress(Expand(x)) should be x")
			require.Zero(t, o, "Odd part should be zero for expanded even bits")

			e, o = compressBitsGeneric(big.Word(expandedVal << 1))
			require.Zero(t, e, "Even part should be zero for expanded odd bits")
			require.Equal(t, big.Word(w), o, "Compress(Expand(x)<<1) should match x in odd part")
		}
	})

	t.Run("xor_homomorphism", func(t *testing.T) {
		for range 20_000 {
			a := uint64(randomWord(t)) | (uint64(randomWord(t)) << 32)
			b := uint64(randomWord(t)) | (uint64(randomWord(t)) << 32)

			ea, oa := compressBitsGeneric(big.Word(a))
			eb, ob := compressBitsGeneric(big.Word(b))
			ex, ox := compressBitsGeneric(big.Word(a ^ b))

			require.Equal(t, ea^eb, ex, "Even part linearity failed")
			require.Equal(t, oa^ob, ox, "Odd part linearity failed")

			pcWant := bits.OnesCount64(a ^ b)
			pcGot := bits.OnesCount64(uint64(ex)) + bits.OnesCount64(uint64(ox))
			require.Equal(t, pcWant, pcGot, "Population count mismatch")
		}
	})
}

func FuzzCLMULAssembly(f *testing.F) {
	if !isCLMULAsmMode {
		f.Skip("CLMUL is on generic mode")
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

func FuzzCLMULWordsAAssembly(f *testing.F) {
	if !isCLMULWordsAsmMode {
		f.Skip("CLMULWords is on generic mode")
	}

	f.Add(make([]byte, 2+Words*2*WordByteSize))

	toWords := func(b []byte) []big.Word {
		if len(b) == 0 {
			return nil
		}
		n := (len(b) + WordByteSize - 1) / WordByteSize

		lst := make([]big.Word, n)

		for idx := range b {
			wIdx := idx / WordByteSize
			shift := (idx % WordByteSize) * 8
			lst[wIdx] |= big.Word(b[idx]) << shift
		}
		return lst
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) < 2 {
			return
		}

		xLen := int(data[0]) % (Words + 1)
		yLen := int(data[1]) % (Words + 1)
		oLen := xLen + yLen
		if oLen == 0 {
			return
		}

		needed := 2 + (xLen+yLen)*WordByteSize
		if len(data) < needed {
			return
		}

		x := toWords(data[2 : 2+xLen*WordByteSize])
		y := toWords(data[2+xLen*WordByteSize:])

		zWant := make([]big.Word, oLen)
		zGot := make([]big.Word, oLen)

		clmulWordsGeneric(zWant, x, y)
		CLMULWords(zGot, x, y)

		require.Equal(t, zWant, zGot, "CLMULWords mismatch: x=%#v y=%#v", x, y)
	})
}

func FuzzExpandBitsAssembly(f *testing.F) {
	if !isExpandBitsAsmMode {
		f.Skip("ExpandBits is on generic mode")
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

func FuzzCompressBitsAssembly(f *testing.F) {
	if !isCompressBitsAsmMode {
		f.Skip("CompressBits is on generic mode")
	}

	f.Add(uint64(0))
	f.Add(uint64(WordMaxValue))

	f.Fuzz(func(t *testing.T, x uint64) {
		wa := big.Word(x)

		wantEven, wantOdd := compressBitsGeneric(wa)
		gotEven, otOdd := CompressBits(wa)

		require.Equal(t, wantEven, gotEven, "CompressBits Even mismatch: x=%#x", x)
		require.Equal(t, wantOdd, otOdd, "CompressBits Odd mismatch: x=%#x", x)
	})
}

func BenchmarkCLMUL(b *testing.B) {
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

	if isCLMULAsmMode {
		b.Run("Assembly", bench(CLMUL))
	}
}

func BenchmarkCLMULWords(b *testing.B) {
	bench := func(fn func(z, x, y []big.Word)) func(b *testing.B) {
		return func(b *testing.B) {
			x := make([][]big.Word, preallocSize)
			for idx := range x {
				x[idx] = randomWordSlice(b, Words)
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
	if isCLMULWordsAsmMode {
		b.Run("Assembly", bench(CLMULWords))
	}
}

func BenchmarkExpandBits(b *testing.B) {
	b.Run("Generic", benchBits(expandBitsGeneric))
	if isExpandBitsAsmMode {
		b.Run("Assembly", benchBits(ExpandBits))
	}
}

func BenchmarkCompressBits(b *testing.B) {
	b.Run("Generic", benchBits(compressBitsGeneric))
	if isCompressBitsAsmMode {
		b.Run("Assembly", benchBits(CompressBits))
	}
}

func benchBits(fn func(x big.Word) (big.Word, big.Word)) func(b *testing.B) {
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
