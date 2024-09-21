package curvetesting

import (
	"crypto/elliptic"
	"testing"

	"github.com/RyuaNerin/elliptic2"
	"github.com/RyuaNerin/elliptic2/internal/curve"
)

func B(b *testing.B, f func(*testing.B, elliptic2.Curve), curves ...elliptic.Curve) {
	for _, c := range curves {
		c := c
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

// ScalarMult
func BMult(b *testing.B, curve elliptic2.Curve) {
	RequireGenerator(b, curve)

	priv := GetRandomK(curve)

	x, y := curve.ScalarBaseMult(priv)
	RequireIsOnCurve(b, curve, x, y)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x, y = curve.ScalarMult(x, y, priv)
		priv[i%len(priv)] ^= byte(x.Bit(0)) << (i % 8)
	}
}

// Add
func BAdd(b *testing.B, curve elliptic2.Curve) {
	RequireGenerator(b, curve)

	x1, y1 := curve.ScalarBaseMult(GetRandomK(curve))
	RequireIsOnCurve(b, curve, x1, y1)
	x2, y2 := curve.ScalarBaseMult(GetRandomK(curve))
	RequireIsOnCurve(b, curve, x2, y2)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x, y := curve.Add(x1, y1, x2, y2)
		x2, y2 = x1, y1
		x1, y1 = x, y
	}
}

// Double
func BDouble(b *testing.B, curve elliptic2.Curve) {
	RequireGenerator(b, curve)

	x, y := curve.ScalarBaseMult(GetRandomK(curve))
	RequireIsOnCurve(b, curve, x, y)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x, y = curve.Double(x, y)
	}
}

// ScalarBaseMult
func BBaseMult(b *testing.B, curve elliptic2.Curve) {
	RequireGenerator(b, curve)

	k := GetRandomK(curve)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		x, _ := curve.ScalarBaseMult(k)
		k[i%len(k)] ^= byte(x.Bit(0)) << (i % 8)
	}
}
