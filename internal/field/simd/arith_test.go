package simd_test

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/RyuaNerin/elliptic2/internal/field/simd"
)

const testWords = 6

func randWords(rnd *rand.Rand, n int) []big.Word {
	w := make([]big.Word, n)
	for i := range w {
		w[i] = big.Word(rnd.Uint64())
	}
	return w
}

func toInt(w []big.Word) *big.Int {
	b := make([]big.Word, len(w))
	copy(b, w)
	return new(big.Int).SetBits(b)
}

func base(n int) *big.Int {
	return new(big.Int).Lsh(big.NewInt(1), uint(simd.WordBitSize*n))
}

func TestAddVV(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))

	for range 1000 {
		x := randWords(rnd, testWords)
		y := randWords(rnd, testWords)
		z := make([]big.Word, testWords)

		c := simd.AddVV(z, x, y)

		want := new(big.Int).Add(toInt(x), toInt(y))
		got := new(big.Int).Lsh(big.NewInt(int64(c)), uint(simd.WordBitSize*testWords))
		got.Add(got, toInt(z))

		require.Zero(t, want.Cmp(got))
	}
}

func TestSubVV(t *testing.T) {
	rnd := rand.New(rand.NewSource(2))

	for range 1000 {
		x := randWords(rnd, testWords)
		y := randWords(rnd, testWords)
		z := make([]big.Word, testWords)

		c := simd.SubVV(z, x, y)

		want := new(big.Int).Sub(toInt(x), toInt(y))
		got := toInt(z)
		if c != 0 {
			got.Sub(got, base(testWords))
		}

		require.Zero(t, want.Cmp(got))
	}
}

func TestAddVWSubVW(t *testing.T) {
	rnd := rand.New(rand.NewSource(3))

	for range 1000 {
		x := randWords(rnd, testWords)
		y := big.Word(rnd.Uint64())
		z := make([]big.Word, testWords)

		c := simd.AddVW(z, x, y)
		want := new(big.Int).Add(toInt(x), toInt([]big.Word{y}))
		got := new(big.Int).Lsh(big.NewInt(int64(c)), uint(simd.WordBitSize*testWords))
		got.Add(got, toInt(z))
		require.Zero(t, want.Cmp(got), "AddVW")

		c = simd.SubVW(z, x, y)
		want = new(big.Int).Sub(toInt(x), toInt([]big.Word{y}))
		got = toInt(z)
		if c != 0 {
			got.Sub(got, base(testWords))
		}
		require.Zero(t, want.Cmp(got), "SubVW")
	}
}

func TestShlVUShrVU(t *testing.T) {
	rnd := rand.New(rand.NewSource(4))

	for range 1000 {
		x := randWords(rnd, testWords)
		z := make([]big.Word, testWords)
		s := uint(rnd.Intn(simd.WordBitSize))

		c := simd.ShlVU(z, x, s)
		want := new(big.Int).Lsh(toInt(x), s)
		got := new(big.Int).Lsh(toInt([]big.Word{c}), uint(simd.WordBitSize*testWords))
		got.Add(got, toInt(z))
		require.Zero(t, want.Cmp(got), "ShlVU s=%d", s)

		simd.ShrVU(z, x, s)
		require.Zero(t, new(big.Int).Rsh(toInt(x), s).Cmp(toInt(z)), "ShrVU s=%d", s)
	}
}

func TestMulAddVWW(t *testing.T) {
	rnd := rand.New(rand.NewSource(5))

	for range 1000 {
		x := randWords(rnd, testWords)
		m := big.Word(rnd.Uint64())
		a := big.Word(rnd.Uint64())
		z := make([]big.Word, testWords)

		c := simd.MulAddVWW(z, x, m, a)

		want := new(big.Int).Mul(toInt(x), toInt([]big.Word{m}))
		want.Add(want, toInt([]big.Word{a}))

		got := new(big.Int).Mul(toInt([]big.Word{c}), base(testWords))
		got.Add(got, toInt(z))

		require.Zero(t, want.Cmp(got))
	}
}

func TestAddMulVVW(t *testing.T) {
	rnd := rand.New(rand.NewSource(6))

	for range 1000 {
		x := randWords(rnd, testWords)
		m := big.Word(rnd.Uint64())
		z := randWords(rnd, testWords)

		want := new(big.Int).Mul(toInt(x), toInt([]big.Word{m}))
		want.Add(want, toInt(z))

		c := simd.AddMulVVW(z, x, m)

		got := new(big.Int).Mul(toInt([]big.Word{c}), base(testWords))
		got.Add(got, toInt(z))

		require.Zero(t, want.Cmp(got))
	}
}

func TestMulWWDivWW(t *testing.T) {
	rnd := rand.New(rand.NewSource(7))

	for range 1000 {
		x := big.Word(rnd.Uint64())
		y := big.Word(rnd.Uint64())

		hi, lo := simd.MulWW(x, y)
		want := new(big.Int).Mul(toInt([]big.Word{x}), toInt([]big.Word{y}))
		require.Zero(t, want.Cmp(toInt([]big.Word{lo, hi})), "MulWW")

		// DivWW : hi < y 이어야 한다.
		if y == 0 {
			continue
		}
		hi %= y
		q, r := simd.DivWW(hi, lo, y)

		u := toInt([]big.Word{lo, hi})
		wantQ, wantR := new(big.Int).QuoRem(u, toInt([]big.Word{y}), new(big.Int))
		require.Zero(t, wantQ.Cmp(toInt([]big.Word{q})), "DivWW q")
		require.Zero(t, wantR.Cmp(toInt([]big.Word{r})), "DivWW r")
	}
}
