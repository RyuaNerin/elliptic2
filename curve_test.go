package elliptic2_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"runtime"
	"testing"

	"github.com/RyuaNerin/elliptic2"
	"github.com/RyuaNerin/elliptic2/internal"
	. "github.com/RyuaNerin/elliptic2/internal/curvetesting"
	"github.com/RyuaNerin/elliptic2/nist"
	"github.com/stretchr/testify/require"
)

func TestWithStdCurves(t *testing.T) {
	t.Run("P-224", test(elliptic.P224(), nist.P224()))
	t.Run("P-256", test(elliptic.P256(), nist.P256()))
	t.Run("P-384", test(elliptic.P384(), nist.P384()))
	t.Run("P-521", test(elliptic.P521(), nist.P521()))
}

func BenchmarkStd(b *testing.B) {
	b.Run("P-224", bench(elliptic.P224()))
	b.Run("P-256", bench(elliptic.P256()))
	b.Run("P-384", bench(elliptic.P384()))
	b.Run("P-521", bench(elliptic.P521()))
}

func BenchmarkLib(b *testing.B) {
	b.Run("P-192", bench(nist.P192()))
	b.Run("P-224", bench(nist.P224()))
	b.Run("P-256", bench(nist.P256()))
	b.Run("P-384", bench(nist.P384()))
	b.Run("P-521", bench(nist.P521()))

	b.Run("B-163", bench(nist.B163()))
	b.Run("B-233", bench(nist.B233()))
	b.Run("B-283", bench(nist.B283()))
	b.Run("B-409", bench(nist.B409()))
	b.Run("B-571", bench(nist.B571()))

	b.Run("K-163", bench(nist.K163()))
	b.Run("K-233", bench(nist.K233()))
	b.Run("K-283", bench(nist.K283()))
	b.Run("K-409", bench(nist.K409()))
	b.Run("K-571", bench(nist.K571()))
}

func test(std elliptic.Curve, lib elliptic.Curve) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Add", func(t *testing.T) {
			// disable deprecation warning of ScalarBaseMult
			std, lib := std.(elliptic2.Curve), lib.(elliptic2.Curve)

			for range 100 {
				x1, y1 := std.ScalarBaseMult(GetRandomK(t, std))
				x2, y2 := std.ScalarBaseMult(GetRandomK(t, std))

				RequireIsOnCurve(t, std, x1, y1, "input point 1 (std)")
				RequireIsOnCurve(t, std, x2, y2, "input point 2 (std)")
				RequireIsOnCurve(t, lib, x1, y1, "input point 1 (lib)")
				RequireIsOnCurve(t, lib, x2, y2, "input point 2 (lib)")

				xStd, yStd := std.Add(x1, y1, x2, y2)
				xLib, yLib := lib.Add(x1, y1, x2, y2)

				RequireIsOnCurve(t, std, xStd, yStd, "Add result (std on std)")
				RequireIsOnCurve(t, std, xLib, yLib, "Add result (lib on std)")
				RequireIsOnCurve(t, lib, xStd, yStd, "Add result (std on lib)")
				RequireIsOnCurve(t, lib, xLib, yLib, "Add result (lib on lib)")

				RequireXYEquals(t, &Point{X: xStd, Y: yStd}, &Point{X: xLib, Y: yLib}, "Add result")
			}
		})
		t.Run("Double", func(t *testing.T) {
			// disable deprecation warning of ScalarBaseMult
			std, lib := std.(elliptic2.Curve), lib.(elliptic2.Curve)

			for range 100 {
				x1, y1 := std.ScalarBaseMult(GetRandomK(t, std))

				RequireIsOnCurve(t, std, x1, y1, "input point (std)")
				RequireIsOnCurve(t, lib, x1, y1, "input point (lib)")

				xStd, yStd := std.Double(x1, y1)
				xLib, yLib := lib.Double(x1, y1)

				RequireIsOnCurve(t, std, xStd, yStd, "Double result (std on std)")
				RequireIsOnCurve(t, std, xLib, yLib, "Double result (lib on std)")
				RequireIsOnCurve(t, lib, xStd, yStd, "Double result (std on lib)")
				RequireIsOnCurve(t, lib, xLib, yLib, "Double result (lib on lib)")

				RequireXYEquals(t, &Point{X: xStd, Y: yStd}, &Point{X: xLib, Y: yLib}, "Double result")
			}
		})
		t.Run("ScalarMult", func(t *testing.T) {
			// disable deprecation warning of ScalarBaseMult
			std, lib := std.(elliptic2.Curve), lib.(elliptic2.Curve)

			for range 100 {
				k := GetRandomK(t, std)
				x1, y1 := std.ScalarBaseMult(GetRandomK(t, std))

				RequireIsOnCurve(t, std, x1, y1, "input point (std)")
				RequireIsOnCurve(t, lib, x1, y1, "input point (lib)")

				xStd, yStd := std.ScalarMult(x1, y1, k)
				xLib, yLib := lib.ScalarMult(x1, y1, k)

				RequireIsOnCurve(t, std, xStd, yStd, "ScalarMult result (std on std)")
				RequireIsOnCurve(t, std, xLib, yLib, "ScalarMult result (lib on std)")
				RequireIsOnCurve(t, lib, xStd, yStd, "ScalarMult result (std on lib)")
				RequireIsOnCurve(t, lib, xLib, yLib, "ScalarMult result (lib on lib)")

				RequireXYEquals(t, &Point{X: xStd, Y: yStd}, &Point{X: xLib, Y: yLib}, "ScalarMult result")
			}
		})
		t.Run("ScalarBaseMult", func(t *testing.T) {
			// disable deprecation warning of ScalarBaseMult
			std, lib := std.(elliptic2.Curve), lib.(elliptic2.Curve)

			for range 100 {
				k := GetRandomK(t, std)

				xStd, yStd := std.ScalarBaseMult(k)
				xLib, yLib := lib.ScalarBaseMult(k)

				RequireIsOnCurve(t, std, xStd, yStd, "ScalarBaseMult result (std on std)")
				RequireIsOnCurve(t, std, xLib, yLib, "ScalarBaseMult result (lib on std)")
				RequireIsOnCurve(t, lib, xStd, yStd, "ScalarBaseMult result (std on lib)")
				RequireIsOnCurve(t, lib, xLib, yLib, "ScalarBaseMult result (lib on lib)")

				RequireXYEquals(t, &Point{X: xStd, Y: yStd}, &Point{X: xLib, Y: yLib}, "ScalarBaseMult result")
			}
		})

		t.Run("ECDSA/SignAndVerify", func(t *testing.T) {
			data := []byte("Hello, World!")

			privStd, err := ecdsa.GenerateKey(std, internal.Random)
			require.NoError(t, err)

			privLib := *privStd
			privLib.Curve = lib

			for range 10 {
				sigStd, err := ecdsa.SignASN1(internal.Random, privStd, data)
				require.NoError(t, err)

				sigLib, err := ecdsa.SignASN1(internal.Random, &privLib, data)
				require.NoError(t, err)

				require.True(t, ecdsa.VerifyASN1(&privStd.PublicKey, data, sigStd), "std: verify failed")
				require.True(t, ecdsa.VerifyASN1(&privLib.PublicKey, data, sigLib), "lib: verify failed")

				// Modify data to change the signature next iteration
				data[0] ^= sigStd[0]
			}
		})
	}
}

func bench(c elliptic.Curve) func(b *testing.B) {
	return func(b *testing.B) {
		b.Run("Add", func(b *testing.B) { BAdd(b, c) })
		b.Run("Double", func(b *testing.B) { BDouble(b, c) })
		b.Run("ScalarMult", func(b *testing.B) { BMult(b, c) })
		b.Run("ScalarBaseMult", func(b *testing.B) { BBaseMult(b, c) })

		b.Run("ECDSA/GenerateKey", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				_, err := ecdsa.GenerateKey(c, internal.Random)
				require.NoError(b, err)
			}
		})
		b.Run("ECDSA/Sign", func(b *testing.B) {
			msg := []byte("Hello, World!")

			priv, err := ecdsa.GenerateKey(c, internal.Random)
			require.NoError(b, err)

			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				sig, err := ecdsa.SignASN1(internal.Random, priv, msg)
				require.NoError(b, err)
				msg[0] = sig[0]
			}
		})
		b.Run("ECDSA/Verify", func(b *testing.B) {
			msg := []byte("Hello, World!")

			priv, err := ecdsa.GenerateKey(c, internal.Random)
			require.NoError(b, err)

			sig, err := ecdsa.SignASN1(internal.Random, priv, msg)
			require.NoError(b, err)

			b.ReportAllocs()
			b.ResetTimer()

			var ok bool
			for range b.N {
				ok = ecdsa.VerifyASN1(&priv.PublicKey, msg, sig)
			}
			runtime.KeepAlive(ok)
		})
	}
}
