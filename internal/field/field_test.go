package field_test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/RyuaNerin/elliptic2/internal/curvetesting"
	"github.com/RyuaNerin/elliptic2/internal/field"
	"github.com/RyuaNerin/elliptic2/internal/field/simd"
	"github.com/stretchr/testify/require"
)

func testInvValidation[E any, GF field.GF[E]](t *testing.T, _ func(GF), modulus []*field.Modulus) {
	x := (GF)(new(E))
	inv := (GF)(new(E))
	one := (GF)(new(E))

	one.SetInt64(1)

	for _, m := range modulus {
		x.SetModulus(m)
		inv.SetModulus(m)
		one.SetModulus(m)

		for range 10 {
			fillWords(t, x)

			inv.Mul(x, inv.Inv(x))
			inv.Reduce()

			require.Equal(t, one, inv, "x * x^(-1) != 1")
			require.Zero(t,
				inv.Cmp(one),
				"x * x^(-1) != 1\n  x:   %s\n  inv: %s\n",
				x.String(), inv.String(),
			)
		}
	}
}

func testSqrValidation[E any, GF field.GF[E]](t *testing.T, _ func(GF), modulus []*field.Modulus) {
	x := (GF)(new(E))
	mul := (GF)(new(E))
	sqr := (GF)(new(E))

	for _, m := range modulus {
		x.SetModulus(m)
		mul.SetModulus(m)
		sqr.SetModulus(m)

		for range 10 {
			fillWords(t, x)

			mul.Mul(x, x)
			sqr.Sqr(x)

			require.Zero(t,
				mul.Cmp(sqr),
				"x * x != x^2, x: %s\nmul: %s\nsqr: %s",
				x.String(), mul.String(), sqr.String(),
			)
		}
	}
}

func fillWords[E any, EP field.GF[E]](t testing.TB, z EP) {
	bitSize := z.Modulus().BitLen()

	wordLen := ((bitSize + simd.WordBitSize - 1) / simd.WordBitSize)
	buf := make([]byte, wordLen*simd.WordByteSize)

	for {
		_, err := rand.Read(buf)
		require.NoError(t, err)

		z.SetBytes(buf)
		z.Reduce()

		if !z.IsZero() {
			return
		}
	}
}

func addWant(tc testCase, argIdx int) *big.Int { return tc.add[argIdx] }
func mulWant(tc testCase, argIdx int) *big.Int { return tc.mul[argIdx] }
func sqrWant(tc testCase, argIdx int) *big.Int { return tc.sqr[argIdx] }
func invWant(tc testCase, argIdx int) *big.Int { return tc.inv[argIdx] }

func add[E any, GF field.GF[E]](z GF, arg ...GF) {
	z.Add(arg[0], arg[1])
	for idx := 2; idx < len(arg); idx++ {
		z.Add(z, arg[idx])
	}
}

func mul[E any, GF field.GF[E]](z GF, arg ...GF) {
	z.Mul(arg[0], arg[1])
	for idx := 2; idx < len(arg); idx++ {
		z.Mul(z, arg[idx])
	}
}

func sqr[E any, GF field.GF[E]](dst GF, arg ...GF) {
	dst.Sqr(arg[0])
}

func inv[E any, GF field.GF[E]](dst GF, arg ...GF) {
	dst.Inv(arg[0])
}

func tArg[E any, GF field.GF[E]](
	t *testing.T,
	_ func(GF),
	testCases []testCase,
	indexes [][]int,
	fnOperation func(z GF, x ...GF),
	fnGetWant func(tc testCase, argIdx int) *big.Int,
) {
	argBackup := make([]big.Int, 3)
	fields := make([]GF, 5)
	for i := range fields {
		fields[i] = (GF)(new(E))
	}
	dst := (GF)(new(E))

	var dstInt *big.Int

	testName := runtime.FuncForPC(reflect.ValueOf(fnOperation).Pointer()).Name()
	if idx := strings.LastIndex(testName, "."); idx >= 0 {
		testName = testName[idx+1:]
	}

	for idx, tc := range testCases {
		for idx, value := range tc.arg {
			argBackup[idx].Set(value)
		}

		for idx := range fields {
			fields[idx].SetModulus(tc.modulus)
		}
		dst.SetModulus(tc.modulus)

		for inputIdx, args := range indexes {
			want := fnGetWant(tc, inputIdx)
			for argIdx, arg := range args {
				fields[argIdx].SetBigInt(tc.arg[arg])
			}

			fnOperation(dst, fields[:len(args)]...)

			dstInt = dst.ToBigInt(dstInt)
			curvetesting.RequireEqual(t,
				want, dstInt,
				"%d: %s[%d]: invalid result",
				idx, testName, inputIdx,
			)
		}

		for aidx := range 3 {
			curvetesting.RequireUnmodified(t, &argBackup[aidx], tc.arg[aidx], "%d: arg[%d]", idx, aidx)
		}
	}
}

func bArg[E any, GF field.GF[E]](
	b *testing.B,
	_ func(GF),
	modulus []*field.Modulus,
	fnOperation func(z GF, x ...GF),
) {
	for _, m := range modulus {
		b.Run(fmt.Sprintf("BitSize=%d", m.BitLen()), func(b *testing.B) {
			dst := (GF)(new(E))
			x := (GF)(new(E))
			y := (GF)(new(E))

			dst.SetModulus(m)
			x.SetModulus(m)
			y.SetModulus(m)

			fillWords(b, x)
			fillWords(b, y)

			b.ReportAllocs()
			b.ResetTimer()

			for range b.N {
				fnOperation(dst, x, y)
			}
		})
	}
}

var tcArgIndexes = struct {
	add [][]int
	mul [][]int
	sqr [][]int
	inv [][]int
}{
	add: [][]int{{0, 0}, {1, 1}, {2, 2}, {0, 1}, {0, 2}, {1, 2}, {0, 1, 2}},
	mul: [][]int{{0, 0}, {1, 1}, {2, 2}, {0, 1}, {0, 2}, {1, 2}, {0, 1, 2}},
	sqr: [][]int{{0}, {1}, {2}},
	inv: [][]int{{0}, {1}, {2}},
}

type testCase struct {
	modulus *field.Modulus
	arg     [3]*big.Int
	add     [7]*big.Int // (a + b)
	mul     [7]*big.Int // (a * b)
	sqr     [3]*big.Int // (a ^ 2)
	inv     [3]*big.Int // (a ^ -1)
}
