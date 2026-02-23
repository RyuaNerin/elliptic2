package weierstrassprime

import (
	"crypto/elliptic"
	"math/big"
	"strings"

	"github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/RyuaNerin/elliptic2/internal/curve/weierstrassprime"
	"github.com/RyuaNerin/elliptic2/internal/field"
)

type CurveParams struct {
	Name    string
	BitSize int
	P       *big.Int // prime
	A, B    *big.Int // a, b
	N       *big.Int // order
	Gx, Gy  *big.Int // generator
	OID     string
}

func NewCurve(params CurveParams) elliptic.Curve {
	// BitSize
	if params.BitSize <= 0 {
		panic("weierstrassprime: BitSize is not positive")
	}

	// P
	if params.P == nil || params.P.Sign() <= 0 {
		panic("weierstrassprime: P must be non-zero (singular curve)")
	}
	if params.P.BitLen() != params.BitSize {
		panic("weierstrassprime: BitSize does not match P")
	}

	// A
	if params.A == nil {
		panic("weierstrassprime: A is nil")
	}
	if params.A.Sign() < 0 {
		panic("weierstrassprime: A is negative")
	}
	if params.A.Cmp(params.P) >= 0 {
		panic("weierstrassprime: A is larger than or equal to P")
	}

	// B
	if params.B == nil {
		panic("weierstrassprime: B is nil")
	}
	if params.B.Sign() < 0 {
		panic("weierstrassprime: B is negative")
	}
	if params.B.Cmp(params.P) >= 0 {
		panic("weierstrassprime: B is larger than or equal to P")
	}

	// N
	if params.N == nil || params.N.Sign() <= 0 {
		panic("weierstrassprime: N must be non-zero (singular curve)")
	}
	if params.N.BitLen() > params.BitSize {
		panic("weierstrassprime: N is too large")
	}

	// Gx, Gy
	if (params.Gx == nil) != (params.Gy == nil) {
		panic("weierstrassprime: Gx and Gy must both be nil or both be non-nil")
	}
	if params.Gx != nil {
		if params.Gx.Sign() < 0 {
			panic("weierstrassprime: Gx is negative")
		}
		if params.Gx.Cmp(params.P) >= 0 {
			panic("weierstrassprime: Gx is larger than or equal to P")
		}
	}
	if params.Gy != nil {
		if params.Gy.Sign() < 0 {
			panic("weierstrassprime: Gy is negative")
		}
		if params.Gy.Cmp(params.P) >= 0 {
			panic("weierstrassprime: Gy is larger than or equal to P")
		}
	}

	// 4aaa + 27bb != 0
	var a3, b2, disc big.Int
	a3.Exp(params.A, big.NewInt(3), params.P) // aaa
	a3.Lsh(&a3, 2)                            // 4aaa
	a3.Mod(&a3, params.P)                     // 4aaa mod p
	b2.Exp(params.B, big.NewInt(2), params.P) // bb
	b2.Mul(&b2, big.NewInt(27))               // 27bb
	b2.Mod(&b2, params.P)                     // 27bb mod p
	disc.Add(&a3, &b2)                        // 4aaa + 27bb
	disc.Mod(&disc, params.P)                 // (4aaa + 27bb) mod p
	if disc.Sign() == 0 {
		panic("weierstrassprime: singular curve (discriminant is zero)")
	}

	var gx, gy *big.Int
	if params.Gx != nil {
		gx = new(big.Int).Set(params.Gx)
	}
	if params.Gy != nil {
		gy = new(big.Int).Set(params.Gy)
	}

	return curve.NewCurve(
		weierstrassprime.Build(
			&weierstrassprime.CurveParams{
				Name:    strings.Clone(params.Name),
				BitSize: params.BitSize,
				P:       field.NewGFpModulus(params.P),
				A:       field.NewGFpInt(params.A),
				B:       field.NewGFpInt(params.B),
				N:       new(big.Int).Set(params.N),
				Gx:      gx,
				Gy:      gy,
			},
		),
		curve.WithOID(params.OID),
	)
}
