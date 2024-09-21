package weierstrassbinary

import (
	"crypto/elliptic"
	"math/big"
	"strings"

	"github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/RyuaNerin/elliptic2/internal/curve/weierstrassbinary"
	"github.com/RyuaNerin/elliptic2/internal/field"
)

type CurveParams struct {
	Name    string
	BitSize int
	Poly    *big.Int // polynomial
	A, B    *big.Int // a, b
	Gx, Gy  *big.Int // generator
	N       *big.Int // order
	OID     string
}

func NewCurve(params CurveParams) elliptic.Curve {
	// BitSize
	if params.BitSize <= 0 {
		panic("weierstrassbinary: BitSize is not positive")
	}

	// Poly
	if params.Poly == nil || params.Poly.Sign() <= 0 {
		panic("weierstrassbinary: Poly must be non-zero (singular curve)")
	}
	if params.Poly.BitLen() != params.BitSize+1 {
		panic("weierstrassbinary: BitSize does not match Poly degree")
	}

	// A
	if params.A == nil {
		panic("weierstrassbinary: A is nil")
	}
	if params.A.Sign() < 0 {
		panic("weierstrassbinary: A is negative")
	}
	if params.A.BitLen() > params.BitSize {
		panic("weierstrassbinary: A is too large")
	}

	// B
	if params.B == nil {
		panic("weierstrassbinary: B is nil")
	}
	if params.B.Sign() <= 0 {
		panic("weierstrassbinary: B must be non-zero (singular curve)")
	}
	if params.B.BitLen() > params.BitSize {
		panic("weierstrassbinary: B is too large")
	}

	// N
	if params.N == nil || params.N.Sign() <= 0 {
		panic("weierstrassbinary: N must be non-zero (singular curve)")
	}
	if params.N.BitLen() > params.BitSize+1 {
		panic("weierstrassbinary: N is too large")
	}

	// Gx, Gy
	if (params.Gx == nil) != (params.Gy == nil) {
		panic("weierstrassbinary: Gx and Gy must both be nil or both be non-nil")
	}
	if params.Gx != nil {
		if params.Gx.Sign() < 0 {
			panic("weierstrassbinary: Gx is negative")
		}
		if params.Gx.BitLen() > params.BitSize {
			panic("weierstrassbinary: Gx is too large")
		}
	}
	if params.Gy != nil {
		if params.Gy.Sign() < 0 {
			panic("weierstrassbinary: Gy is negative")
		}
		if params.Gy.BitLen() > params.BitSize {
			panic("weierstrassbinary: Gy is too large")
		}
	}

	var gx, gy *big.Int
	if params.Gx != nil {
		gx = new(big.Int).Set(params.Gx)
	}
	if params.Gy != nil {
		gy = new(big.Int).Set(params.Gy)
	}

	modulus := field.NewGF2mModulus(params.Poly)

	return curve.NewCurve(
		weierstrassbinary.Build(
			&weierstrassbinary.CurveParams{
				Name:    strings.Clone(params.Name),
				BitSize: params.BitSize,
				Poly:    modulus,
				A2:      modulus.NewGF2mFromBigInt(params.A),
				A6:      modulus.NewGF2mFromBigInt(params.B),
				N:       modulus.NewGF2mFromBigInt(params.N),
				Gx:      gx,
				Gy:      gy,
			},
		),
		curve.WithOID(params.OID),
	)
}

// Polynomials to big.Int
// f = [233, 74, 0] when x^233 +  x^74 + 1
// f = [283, 12, 7, 5, 0] when x^283 +  x^12 +  x^7 +  x^5 + 1
func Polynomials(f ...int) *big.Int       { return field.GF2mPolynomials(f...) }
func ExtractPolynomials(f *big.Int) []int { return field.ParseGF2mPolynomials(f) }
