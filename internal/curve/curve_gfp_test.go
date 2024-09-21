package curve_test

import (
	"fmt"
	"testing"

	. "github.com/RyuaNerin/elliptic2/internal"
	. "github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/RyuaNerin/elliptic2/internal/curve/edwards"
	"github.com/RyuaNerin/elliptic2/internal/curve/montgomery"
	"github.com/RyuaNerin/elliptic2/internal/curve/twistededwards"
	"github.com/RyuaNerin/elliptic2/internal/curve/weierstrassprime"
	. "github.com/RyuaNerin/elliptic2/internal/curvetesting"
	. "github.com/RyuaNerin/elliptic2/internal/field"
)

func TestGFpCurveMadd(t *testing.T) { TestCurveMadd(t, gfpCurves...) }

func BenchmarkGFp(b *testing.B) {
	for _, gfpCurve := range gfpCurves {
		c := NewCurveSimple(gfpCurve)
		b.Run(
			fmt.Sprintf("%s/%s", GetCurveType(c), GetName(c)),
			func(b *testing.B) { BMult(b, c) },
		)
	}
}

func BenchmarkGFpMadd(b *testing.B) {
	for _, gfpCurve := range gfpCurves {
		c := NewCurveMadd(gfpCurve)
		b.Run(
			fmt.Sprintf("%s/%s", GetCurveType(c), GetName(c)),
			func(b *testing.B) { BMult(b, c) },
		)
	}
}

var gfpCurves = []CurveArithmeticBase{
	weierstrassprime.Build(&weierstrassprime.CurveParams{
		Name:    `P-192`,
		BitSize: 192,
		P:       NewGFpModulusFromHex(`fffffffffffffffffffffffffffffffeffffffffffffffff`),
		A:       ParseGFpHex(`fffffffffffffffffffffffffffffffefffffffffffffffc`),
		B:       ParseGFpHex(`64210519e59c80e70fa7e9ab72243049feb8deecc146b9b1`),
		N:       ParseGFpHex(`ffffffffffffffffffffffff99def836146bc9b1b4d22831`),
		Gx:      HI(`188da80eb03090f67cbf20eb43a18800f4ff0afd82ff1012`),
		Gy:      HI(`7192b95ffc8da78631011ed6b24cdd573f977a11e794811`),
	}),
	edwards.Build(&edwards.CurveParams{
		Name:    `E-222`,
		BitSize: 222,
		P:       NewGFpModulusFromHex(`3fffffffffffffffffffffffffffffffffffffffffffffffffffff8b`),
		C:       ParseGFpHex(`1`),
		D:       ParseGFpHex(`27166`),
		N:       ParseGFpHex(`ffffffffffffffffffffffffffff70cbc95e932f802f31423598cbf`),
		Gx:      HI(`19b12bb156a389e55c9768c303316d07c23adab3736eb2bc3eb54e51`),
		Gy:      HI(`1c`),
	}),
	montgomery.Build(&montgomery.CurveParams{
		Name:    `Curve25519`,
		BitSize: 255,
		P:       NewGFpModulusFromHex(`7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffed`),
		A:       ParseGFpHex(`76d06`),
		B:       ParseGFpHex(`1`),
		N:       ParseGFpHex(`1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3ed`),
		Gx:      HI(`9`),
		Gy:      HI(`20ae19a1b8a086b4e01edd2c7748d14c923d4d7e6d7c61b229e9c5a27eced3d9`),
	}),
	twistededwards.Build(&twistededwards.CurveParams{
		Name:    `Ed25519`,
		BitSize: 255,
		P:       NewGFpModulusFromHex(`7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffed`),
		A:       ParseGFpHex(`7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffec`),
		D:       ParseGFpHex(`52036cee2b6ffe738cc740797779e89800700a4d4141d8ab75eb4dca135978a3`),
		N:       ParseGFpHex(`1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3ed`),
		Gx:      HI(`216936d3cd6e53fec0a4e231fdd6dc5c692cc7609525a7b2c9562d608f25d51a`),
		Gy:      HI(`6666666666666666666666666666666666666666666666666666666666666658`),
	}),
}
