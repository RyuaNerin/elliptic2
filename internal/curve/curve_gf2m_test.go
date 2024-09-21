package curve_test

import (
	"fmt"
	"testing"

	. "github.com/RyuaNerin/elliptic2/internal"
	. "github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/RyuaNerin/elliptic2/internal/curve/weierstrassbinary"
	. "github.com/RyuaNerin/elliptic2/internal/curvetesting"
	. "github.com/RyuaNerin/elliptic2/internal/field"
)

func TestGF2mCurveMadd(t *testing.T) { TestCurveMadd(t, gf2mCurves...) }

func BenchmarkGF2m(b *testing.B) {
	for _, gf2mCurve := range gf2mCurves {
		c := NewCurveSimple(gf2mCurve)
		b.Run(
			fmt.Sprintf("%s/%s", GetCurveType(c), GetName(c)),
			func(b *testing.B) { BMult(b, c) },
		)
	}
}

func BenchmarkGF2mMadd(b *testing.B) {
	for _, gf2mCurve := range gf2mCurves {
		c := NewCurveMadd(gf2mCurve)
		b.Run(
			fmt.Sprintf("%s/%s", GetCurveType(c), GetName(c)),
			func(b *testing.B) { BMult(b, c) },
		)
	}
}

var gf2mCurves = []CurveArithmeticBase{
	weierstrassbinary.Build(&weierstrassbinary.CurveParams{
		Name:    "B-163",
		BitSize: 163,
		Poly:    NewGF2mModulusFromPolynomials(163, 7, 6, 3, 0),
		A2:      ParseGF2mHex("1"),
		A6:      ParseGF2mHex("20a601907b8c953ca1481eb10512f78744a3205fd"),
		N:       ParseGF2mHex("40000000000000000000292fe77e70c12a4234c33"),
		Gx:      HI("3f0eba16286a2d57ea0991168d4994637e8343e36"),
		Gy:      HI("d51fbc6c71a0094fa2cdd545b11c5c0c797324f1"),
	}),
}
