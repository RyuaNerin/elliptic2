package twistededwards

import (
	"crypto/elliptic"
	"math/big"
	"strings"

	"github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/RyuaNerin/elliptic2/internal/curve/twistededwards"
	"github.com/RyuaNerin/elliptic2/internal/field"
)

type CurveParams struct {
	Name    string
	BitSize int
	P       *big.Int // prime
	A, D    *big.Int // a, d
	Gx, Gy  *big.Int // generator
	N       *big.Int // order
	OID     string
}

func NewCurve(params CurveParams) elliptic.Curve {
	// BitSize
	if params.BitSize <= 0 {
		panic("twistededwards: BitSize must be positive")
	}

	// P
	if params.P == nil || params.P.Sign() <= 0 {
		panic("twistededwards: P must be non-zero (singular curve)")
	}
	if params.P.BitLen() != params.BitSize {
		panic("twistededwards: BitSize does not match P")
	}

	// A
	if params.A == nil {
		panic("twistededwards: A is nil")
	}
	if params.A.Sign() < 0 {
		panic("twistededwards: A is negative")
	}
	if params.A.Cmp(params.P) >= 0 {
		panic("twistededwards: A is larger than or equal to P")
	}

	// D
	if params.D == nil {
		panic("twistededwards: D is nil")
	}
	if params.D.Sign() < 0 {
		panic("twistededwards: D is negative")
	}
	if params.D.Cmp(params.P) >= 0 {
		panic("twistededwards: D is larger than or equal to P")
	}

	// A != D
	if params.A.Cmp(params.D) == 0 {
		panic("twistededwards: A and D must be different")
	}

	// A * D * (A - D) != 0
	if params.A.Sign() == 0 {
		panic("twistededwards: A must be non-zero")
	}
	if params.D.Sign() == 0 {
		panic("twistededwards: D must be non-zero")
	}
	var t big.Int
	t.Sub(params.A, params.D)
	t.Mul(&t, params.A)
	t.Mul(&t, params.D)
	t.Mod(&t, params.P)
	if t.Sign() == 0 {
		panic("twistededwards: A * D * (A - D) must be non-zero")
	}

	// N
	if params.N == nil || params.N.Sign() <= 0 {
		panic("twistededwards: N must be non-zero (singular curve)")
	}
	if params.N.BitLen() > params.BitSize {
		panic("twistededwards: N is too large")
	}

	// Gx, Gy
	if (params.Gx == nil) != (params.Gy == nil) {
		panic("twistededwards: Gx and Gy must both be nil or both be non-nil")
	}
	if params.Gx != nil {
		if params.Gx.Sign() < 0 {
			panic("twistededwards: Gx is negative")
		}
		if params.Gx.Cmp(params.P) >= 0 {
			panic("twistededwards: Gx is larger than or equal to P")
		}
	}
	if params.Gy != nil {
		if params.Gy.Sign() < 0 {
			panic("twistededwards: Gy is negative")
		}
		if params.Gy.Cmp(params.P) >= 0 {
			panic("twistededwards: Gy is larger than or equal to P")
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
		twistededwards.Build(
			&twistededwards.CurveParams{
				Name:    strings.Clone(params.Name),
				BitSize: params.BitSize,
				P:       field.NewGFpModulus(params.P),
				A:       field.NewGFpInt(params.A),
				D:       field.NewGFpInt(params.D),
				N:       new(big.Int).Set(params.N),
				Gx:      gx,
				Gy:      gy,
			},
		),
		curve.WithOID(params.OID),
	)
}
