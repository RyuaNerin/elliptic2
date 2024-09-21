package montgomery

import (
	"crypto/elliptic"
	"math/big"
	"strings"

	"github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/RyuaNerin/elliptic2/internal/curve/montgomery"
	"github.com/RyuaNerin/elliptic2/internal/field"
)

type CurveParams struct {
	Name    string
	BitSize int
	P       *big.Int // prime
	A, B    *big.Int // a, b
	Gx, Gy  *big.Int // generator
	N       *big.Int // order
	OID     string
}

func NewCurve(params CurveParams) elliptic.Curve {
	// BitSize
	if params.BitSize <= 0 {
		panic("montgomery: BitSize must be positive")
	}

	// P
	if params.P == nil || params.P.Sign() <= 0 {
		panic("montgomery: P must be non-zero (singular curve)")
	}
	if params.P.BitLen() != params.BitSize {
		panic("montgomery: BitSize does not match P")
	}

	// A
	if params.A == nil {
		panic("montgomery: A is nil")
	}
	if params.A.Sign() < 0 {
		panic("montgomery: A is negative")
	}
	if params.A.Cmp(params.P) >= 0 {
		panic("montgomery: A is larger than or equal to P")
	}

	// B
	if params.B == nil {
		panic("montgomery: B is nil")
	}
	if params.B.Sign() <= 0 {
		panic("montgomery: B must be positive")
	}
	if params.B.Cmp(params.P) >= 0 {
		panic("montgomery: B is larger than or equal to P")
	}

	// B * (A^2 - 4) != 0
	var t big.Int
	t.Mul(params.A, params.A)
	t.Sub(&t, big.NewInt(4))
	t.Mul(&t, params.B)
	t.Mod(&t, params.P)
	if t.Sign() == 0 {
		panic("montgomery: B * (A^2 - 4) must be non-zero")
	}

	// N
	if params.N == nil || params.N.Sign() <= 0 {
		panic("montgomery: N must be non-zero (singular curve)")
	}
	if params.N.BitLen() > params.BitSize {
		panic("montgomery: N is too large")
	}

	// Gx, Gy
	if (params.Gx == nil) != (params.Gy == nil) {
		panic("montgomery: Gx and Gy must both be nil or both be non-nil")
	}
	if params.Gx != nil {
		if params.Gx.Sign() < 0 {
			panic("montgomery: Gx is negative")
		}
		if params.Gx.Cmp(params.P) >= 0 {
			panic("montgomery: Gx is larger than or equal to P")
		}
	}
	if params.Gy != nil {
		if params.Gy.Sign() < 0 {
			panic("montgomery: Gy is negative")
		}
		if params.Gy.Cmp(params.P) >= 0 {
			panic("montgomery: Gy is larger than or equal to P")
		}
	}

	var gx, gy *big.Int
	if params.Gx != nil {
		gx = new(big.Int).Set(params.Gx)
	}
	if params.Gy != nil {
		gy = new(big.Int).Set(params.Gy)
	}

	return curve.NewCurve(
		montgomery.Build(
			&montgomery.CurveParams{
				Name:    strings.Clone(params.Name),
				BitSize: params.BitSize,
				P:       field.NewGFpModulus(params.P),
				A:       field.NewGFpInt(params.A),
				B:       field.NewGFpInt(params.B),
				N:       field.NewGFpInt(params.N),
				Gx:      gx,
				Gy:      gy,
			},
		),
		curve.WithOID(params.OID),
	)
}
