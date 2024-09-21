package edwards

import (
	"crypto/elliptic"
	"math/big"
	"strings"

	"github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/RyuaNerin/elliptic2/internal/curve/edwards"
	"github.com/RyuaNerin/elliptic2/internal/field"
)

type CurveParams struct {
	Name    string
	BitSize int
	P       *big.Int // prime
	C, D    *big.Int // c, d
	Gx, Gy  *big.Int // generator
	N       *big.Int // order
	OID     string
}

func NewCurve(params CurveParams) elliptic.Curve {
	// BitSize
	if params.BitSize <= 0 {
		panic("edwards: BitSize must be positive")
	}

	// P
	if params.P == nil || params.P.Sign() <= 0 {
		panic("edwards: P must be non-zero (singular curve)")
	}
	if params.P.BitLen() != params.BitSize {
		panic("edwards: BitSize does not match P")
	}

	// C
	if params.C == nil {
		panic("edwards: C is nil")
	}
	if params.C.Sign() <= 0 {
		panic("edwards: C must be positive")
	}
	if params.C.Cmp(params.P) >= 0 {
		panic("edwards: C is larger than or equal to P")
	}

	// D
	if params.D == nil {
		panic("edwards: D is nil")
	}
	if params.D.Sign() < 0 {
		panic("edwards: D is negative")
	}
	if params.D.Cmp(params.P) >= 0 {
		panic("edwards: D is larger than or equal to P")
	}

	// C * D * (1 - C^4 * D) != 0
	if params.D.Sign() == 0 {
		panic("edwards: D must be non-zero")
	}
	var t big.Int
	t.Exp(params.C, big.NewInt(4), params.P) //              c^4
	t.Mul(&t, params.D)                      //              c^4 * d
	t.Sub(big.NewInt(1), &t)                 //          1 - c^4 * d
	t.Mul(&t, params.D)                      //     D * (1 - c^4 * d)
	t.Mul(&t, params.C)                      // C * D * (1 - c^4 * d)
	t.Mod(&t, params.P)
	if t.Sign() == 0 {
		panic("edwards: C * D * (1 - C^4 * D) must be non-zero")
	}

	// N
	if params.N == nil || params.N.Sign() <= 0 {
		panic("edwards: N must be non-zero (singular curve)")
	}
	if params.N.BitLen() > params.BitSize {
		panic("edwards: N is too large")
	}

	// Gx, Gy
	if (params.Gx == nil) != (params.Gy == nil) {
		panic("edwards: Gx and Gy must both be nil or both be non-nil")
	}
	if params.Gx != nil {
		if params.Gx.Sign() < 0 {
			panic("edwards: Gx is negative")
		}
		if params.Gx.Cmp(params.P) >= 0 {
			panic("edwards: Gx is larger than or equal to P")
		}
	}
	if params.Gy != nil {
		if params.Gy.Sign() < 0 {
			panic("edwards: Gy is negative")
		}
		if params.Gy.Cmp(params.P) >= 0 {
			panic("edwards: Gy is larger than or equal to P")
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
		edwards.Build(
			&edwards.CurveParams{
				Name:    strings.Clone(params.Name),
				BitSize: params.BitSize,
				P:       field.NewGFpModulus(params.P),
				C:       field.NewGFpInt(params.C),
				D:       field.NewGFpInt(params.D),
				N:       field.NewGFpInt(params.N),
				Gx:      gx,
				Gy:      gy,
			},
		),
		curve.WithOID(params.OID),
	)
}
