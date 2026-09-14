package field_test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/RyuaNerin/elliptic2/internal"
	. "github.com/RyuaNerin/elliptic2/internal/field"
)

func gfp(_ *GFp) {}

func TestGFpInvValidation(t *testing.T) { testInvValidation(t, gfp, gfpModulus) }
func TestGFpSqrValidation(t *testing.T) { testSqrValidation(t, gfp, gfpModulus) }

func BenchmarkGFpAdd(b *testing.B) { bArg(b, gfp, gfpModulus, add) }
func BenchmarkGFpMul(b *testing.B) { bArg(b, gfp, gfpModulus, mul) }
func BenchmarkGFpSqr(b *testing.B) { bArg(b, gfp, gfpModulus, sqr) }
func BenchmarkGFpInv(b *testing.B) { bArg(b, gfp, gfpModulus, inv) }

var gfpModulus = []*Modulus{
	NewGFpModulus(HI(`fffffffffffffffffffffffffffffffeffffffffffffffff`)),                                                                                    // p % 4 == 3
	NewGFpModulus(HI(`ffffffffffffffffffffffffffffffff000000000000000000000001`)),                                                                            // p % 4 == 1
	NewGFpModulus(HI(`ffffffff00000001000000000000000000000000ffffffffffffffffffffffff`)),                                                                    // p % 4 == 3
	NewGFpModulus(HI(`fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffeffffffff0000000000000000ffffffff`)),                                    // p % 4 == 3
	NewGFpModulus(HI(`1ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff`)), // p % 4 == 3
	NewGFpModulus(HI(`a9fb57dba1eea9bc3e660a909d838d726e3bf623d52620282013481d1f6e5377`)),                                                                    // brainpoolP256r1, p % 4 == 1
}

// randGFp 는 [0, p) 범위의 임의의 값을 반환한다.
func randGFp(t testing.TB, m *Modulus) (*GFp, *big.Int) {
	p := m.ToBigInt(nil)

	v, err := rand.Int(rand.Reader, p)
	require.NoError(t, err)

	return m.NewGFpFromBigInt(v), v
}

func TestGFpArithmetic(t *testing.T) {
	for _, m := range gfpModulus {
		p := m.ToBigInt(nil)

		z := m.NewGFp()
		want := new(big.Int)

		for range 100 {
			x, xi := randGFp(t, m)
			y, yi := randGFp(t, m)

			// Add
			z.Add(x, y)
			want.Add(xi, yi).Mod(want, p)
			requireEqualGFp(t, want, z, "Add")

			// Sub
			z.Sub(x, y)
			want.Sub(xi, yi).Mod(want, p)
			requireEqualGFp(t, want, z, "Sub")

			// Sub (reversed)
			z.Sub(y, x)
			want.Sub(yi, xi).Mod(want, p)
			requireEqualGFp(t, want, z, "Sub reversed")

			// Mul
			z.Mul(x, y)
			want.Mul(xi, yi).Mod(want, p)
			requireEqualGFp(t, want, z, "Mul")

			// Sqr
			z.Sqr(x)
			want.Mul(xi, xi).Mod(want, p)
			requireEqualGFp(t, want, z, "Sqr")

			// Neg
			z.Neg(x)
			want.Neg(xi).Mod(want, p)
			requireEqualGFp(t, want, z, "Neg")

			// Inv
			z.Inv(x)
			want.ModInverse(xi, p)
			requireEqualGFp(t, want, z, "Inv")

			// Lsh / Rsh
			for _, s := range []uint{0, 1, 3, 63, 64, 65, 130} {
				z.Lsh(x, s)
				want.Lsh(xi, s).Mod(want, p)
				requireEqualGFp(t, want, z, "Lsh")

				z.Rsh(x, s)
				want.Rsh(xi, s)
				requireEqualGFp(t, want, z, "Rsh")
			}
		}

		// 자기 자신을 destination 으로 사용하는 경우
		x, xi := randGFp(t, m)
		x.Mul(x, x)
		want.Mul(xi, xi).Mod(want, p)
		requireEqualGFp(t, want, x, "Mul aliased")

		x, xi = randGFp(t, m)
		x.Add(x, x)
		want.Add(xi, xi).Mod(want, p)
		requireEqualGFp(t, want, x, "Add aliased")

		x, xi = randGFp(t, m)
		x.Sqr(x)
		want.Mul(xi, xi).Mod(want, p)
		requireEqualGFp(t, want, x, "Sqr aliased")
	}
}

func TestGFpSqrt(t *testing.T) {
	for _, m := range gfpModulus {
		p := m.ToBigInt(nil)

		sqr := m.NewGFp()
		sqrt := m.NewGFp()
		want := new(big.Int)

		for range 20 {
			x, xi := randGFp(t, m)

			sqr.Sqr(x)
			require.NotNil(t, sqrt.Sqrt(sqr), "sqrt(x^2) must exist")

			// sqrt(x^2) 는 x 또는 p - x 이다.
			want.Mul(xi, xi).Mod(want, p)
			got := sqrt.ToBigInt(nil)
			require.Zero(
				t,
				new(big.Int).Mod(new(big.Int).Mul(got, got), p).Cmp(want),
				"sqrt(x)^2 != x, p: %s, x: %s, sqrt: %s", p.Text(16), xi.Text(16), got.Text(16),
			)
		}

		// 이차 비잉여는 nil 을 반환해야 한다.
		nonResidue := m.NewGFp()
		for i := int64(2); ; i++ {
			v := big.NewInt(i)
			if big.Jacobi(v, p) == -1 {
				nonResidue.SetBigInt(v)
				require.Nil(t, sqrt.Sqrt(nonResidue), "sqrt of non-residue must be nil")
				break
			}
		}
	}
}

// TestGFpCompatibleWithGFp 는 기존 GFp 구현과 결과가 (mod p 기준으로) 동일한지 확인한다.
func TestGFpCompatibleWithGFp(t *testing.T) {
	for _, m := range gfpModulus {
		v1 := m.NewGFp()
		v2 := m.NewGFp()

		x1, y1 := m.NewGFp(), m.NewGFp()
		x2, y2 := m.NewGFp(), m.NewGFp()

		for range 50 {
			_, xi := randGFp(t, m)
			_, yi := randGFp(t, m)

			x1.SetBigInt(xi)
			y1.SetBigInt(yi)
			x2.SetBigInt(xi)
			y2.SetBigInt(yi)

			requireSameGFp(t, v1.Add(x1, y1).Reduce(), v2.Add(x2, y2), "Add")
			requireSameGFp(t, v1.Sub(x1, y1).Reduce(), v2.Sub(x2, y2), "Sub")
			requireSameGFp(t, v1.Mul(x1, y1).Reduce(), v2.Mul(x2, y2), "Mul")
			requireSameGFp(t, v1.Sqr(x1).Reduce(), v2.Sqr(x2), "Sqr")
			requireSameGFp(t, v1.Neg(x1).Reduce(), v2.Neg(x2), "Neg")
			requireSameGFp(t, v1.Inv(x1).Reduce(), v2.Inv(x2), "Inv")
			requireSameGFp(t, v1.Lsh(x1, 3).Reduce(), v2.Lsh(x2, 3), "Lsh")
		}
	}
}

func TestGFpSetters(t *testing.T) {
	m := gfpModulus[2]
	p := m.ToBigInt(nil)

	z := m.NewGFp()
	want := new(big.Int)

	// p 보다 큰 값
	big1 := new(big.Int).Mul(p, big.NewInt(7))
	big1.Add(big1, big.NewInt(12345))
	requireEqualGFp(t, want.Mod(big1, p), z.SetBigInt(big1), "SetBigInt (large)")
	requireEqualGFp(t, want.Mod(big1, p), z.SetBytes(big1.Bytes()), "SetBytes (large)")

	// 음수
	neg := new(big.Int).Neg(big.NewInt(0x1234567890abcdef))
	requireEqualGFp(t, want.Mod(neg, p), z.SetBigInt(neg), "SetBigInt (negative)")
	requireEqualGFp(t, want.Mod(big.NewInt(-3), p), z.SetInt64(-3), "SetInt64 (negative)")

	// 작은 값
	requireEqualGFp(t, big.NewInt(0), z.SetInt64(0), "SetInt64 (0)")
	require.True(t, z.IsZero())
	require.Equal(t, 0, z.Sign())

	requireEqualGFp(t, big.NewInt(1), z.SetUint64(1), "SetUint64 (1)")
	require.Equal(t, 1, z.Sign())
	require.Equal(t, 1, z.BitLen())
	require.Equal(t, uint(1), z.Bit(0))

	// Reduce / ReduceSoft 는 정규화된 값을 바꾸지 않는다.
	x, xi := randGFp(t, m)
	requireEqualGFp(t, xi, x.Reduce(), "Reduce")
	requireEqualGFp(t, xi, x.ReduceSoft(), "ReduceSoft")
	require.Equal(t, xi.BitLen(), x.BitLen())
	require.Equal(t, xi.String(), x.String())
	require.Equal(t, xi.Text(16), x.Text(16))
}

func requireEqualGFp(t testing.TB, want *big.Int, got *GFp, msg string) {
	t.Helper()

	require.NotNil(t, got, msg)
	require.Zero(
		t,
		want.Cmp(got.ToBigInt(nil)),
		"%s\n  want: %s\n  got:  %s", msg, want.Text(16), got.Text(16),
	)
}

func requireSameGFp(t testing.TB, want *GFp, got *GFp, msg string) {
	t.Helper()

	require.NotNil(t, want, msg)
	require.NotNil(t, got, msg)
	require.Zero(
		t,
		want.ToBigInt(nil).Cmp(got.ToBigInt(nil)),
		"%s\n  GFp:   %s\n  GFp: %s", msg, want.Text(16), got.Text(16),
	)
}

func BenchmarkCompare(b *testing.B) {
	for _, m := range gfpModulus {
		p := m.ToBigInt(nil)

		xi, _ := rand.Int(rand.Reader, p)
		yi, _ := rand.Int(rand.Reader, p)

		// big.Int
		var bz, bx, by big.Int
		bx.Set(xi)
		by.Set(yi)

		// GFp (big.Int 기반)
		v1z, v1x, v1y := m.NewGFp(), m.NewGFpFromBigInt(xi), m.NewGFpFromBigInt(yi)

		// GFp (고정 크기)
		v2z, v2x, v2y := m.NewGFp(), m.NewGFpFromBigInt(xi), m.NewGFpFromBigInt(yi)

		impls := []struct {
			name string
			add  func()
			mul  func()
			sqr  func()
			inv  func()
		}{
			{
				name: "big.Int",
				add:  func() { bz.Add(&bx, &by); bz.Mod(&bz, p) },
				mul:  func() { bz.Mul(&bx, &by); bz.Mod(&bz, p) },
				sqr:  func() { bz.Mul(&bx, &bx); bz.Mod(&bz, p) },
				inv:  func() { bz.ModInverse(&bx, p) },
			},
			{
				name: "GFp",
				add:  func() { v1z.Add(v1x, v1y).Reduce() },
				mul:  func() { v1z.Mul(v1x, v1y).Reduce() },
				sqr:  func() { v1z.Sqr(v1x).Reduce() },
				inv:  func() { v1z.Inv(v1x) },
			},
			{
				name: "GFp",
				add:  func() { v2z.Add(v2x, v2y) },
				mul:  func() { v2z.Mul(v2x, v2y) },
				sqr:  func() { v2z.Sqr(v2x) },
				inv:  func() { v2z.Inv(v2x) },
			},
		}

		b.Run(fmt.Sprintf("bits=%d", m.BitLen()), func(b *testing.B) {
			for _, impl := range impls {
				b.Run(impl.name, func(b *testing.B) {
					for _, op := range []struct {
						name string
						fn   func()
					}{
						{"Add", impl.add},
						{"Mul", impl.mul},
						{"Sqr", impl.sqr},
						{"Inv", impl.inv},
					} {
						b.Run(op.name, func(b *testing.B) {
							b.ReportAllocs()
							for range b.N {
								op.fn()
							}
						})
					}
				})
			}
		})
	}
}
