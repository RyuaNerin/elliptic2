package curvetesting

import (
	"crypto/elliptic"
	"math/big"
	"runtime"
	"testing"

	"github.com/RyuaNerin/elliptic2"
	"github.com/RyuaNerin/elliptic2/internal/curve"
)

const benchKeyCount = 64

func B(b *testing.B, f func(*testing.B, elliptic2.Curve), curves ...elliptic.Curve) {
	for _, c := range curves {
		b.Run(GetName(c), func(b *testing.B) {
			f(b, c)
		})
	}
}

func BOp[
	TCurveParams any,
	TCurveArithmetic curve.CurveArithmeticBase,
	TOperator any,
](
	b *testing.B,
	fnBuild func(params *TCurveParams, fnNewOp func(c TCurveArithmetic) TOperator) TCurveArithmetic,
	fnNewOp func(c TCurveArithmetic) TOperator,
	fnBench func(b *testing.B, curve elliptic2.Curve),
	curves ...elliptic2.Curve,
) {
	for _, c := range curves {
		base := curve.GetBase(c)
		c := curve.NewCurve(fnBuild(base.RawParams().(*TCurveParams), fnNewOp))

		b.Run(GetName(c), func(b *testing.B) { fnBench(b, c) })
	}
}

func generateKeys(b *testing.B, c elliptic2.Curve, n int) [][]byte {
	b.Helper()
	keys := make([][]byte, n)
	for idx := range keys {
		keys[idx] = GetRandomK(b, c)
	}
	return keys
}

func generatePoints(b *testing.B, c elliptic2.Curve, keys [][]byte) (xs, ys []*big.Int) {
	b.Helper()
	xs = make([]*big.Int, len(keys))
	ys = make([]*big.Int, len(keys))
	for idx, key := range keys {
		xs[idx], ys[idx] = c.ScalarBaseMult(key)
		RequireIsOnCurve(b, c, xs[idx], ys[idx], "generated point")
	}
	return
}

// ScalarMult
func BMult(b *testing.B, c elliptic2.Curve) {
	RequireGenerator(b, c)

	keys := generateKeys(b, c, benchKeyCount)
	xs, ys := generatePoints(b, c, keys)

	b.ReportAllocs()
	b.ResetTimer()

	var x, y *big.Int
	for iter := range b.N {
		idxKey := iter % benchKeyCount
		x, y = c.ScalarMult(xs[idxKey], ys[idxKey], keys[idxKey])
	}
	runtime.KeepAlive(x)
	runtime.KeepAlive(y)
}

// Add
func BAdd(b *testing.B, c elliptic2.Curve) {
	RequireGenerator(b, c)

	keys := generateKeys(b, c, benchKeyCount*2)
	xs, ys := generatePoints(b, c, keys)

	b.ReportAllocs()
	b.ResetTimer()

	var x, y *big.Int
	for iter := range b.N {
		idxKey := (iter * 2) % (benchKeyCount * 2)
		x, y = c.Add(xs[idxKey], ys[idxKey], xs[idxKey+1], ys[idxKey+1])
	}
	runtime.KeepAlive(x)
	runtime.KeepAlive(y)
}

// Double
func BDouble(b *testing.B, c elliptic2.Curve) {
	RequireGenerator(b, c)

	keys := generateKeys(b, c, benchKeyCount)
	xs, ys := generatePoints(b, c, keys)

	b.ReportAllocs()
	b.ResetTimer()

	var x, y *big.Int
	for iter := range b.N {
		idxKey := iter % benchKeyCount
		x, y = c.Double(xs[idxKey], ys[idxKey])
	}
	runtime.KeepAlive(x)
	runtime.KeepAlive(y)
}

// ScalarBaseMult
func BBaseMult(b *testing.B, c elliptic2.Curve) {
	RequireGenerator(b, c)

	keys := generateKeys(b, c, benchKeyCount)

	b.ReportAllocs()
	b.ResetTimer()

	var x, y *big.Int
	for iter := range b.N {
		x, y = c.ScalarBaseMult(keys[iter%benchKeyCount])
	}
	runtime.KeepAlive(x)
	runtime.KeepAlive(y)
}
