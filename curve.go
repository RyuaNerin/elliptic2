package elliptic2

import (
	"crypto/elliptic"
	"encoding/asn1"
	"math/big"
)

type CurveType int

const (
	_ CurveType = iota
	CurveTypeWeierstrassPrime
	CurveTypeWeierstrassBinary
	CurveTypeMontgomery
	CurveTypeEdwards
	CurveTypeTwistedEdwards
	curveTypeEnd
)

var curveTypeNames = []string{
	CurveTypeWeierstrassPrime:  "Weierstrass (prime field)",
	CurveTypeWeierstrassBinary: "Weierstrass (binary field)",
	CurveTypeMontgomery:        "Montgomery",
	CurveTypeEdwards:           "Edwards",
	CurveTypeTwistedEdwards:    "Twisted Edwards",
}

func init() {
	for i := 1; i < int(curveTypeEnd); i++ {
		if i >= len(curveTypeNames) {
			panic("missing curve type name")
		}
		if curveTypeNames[i] == "" {
			panic("missing curve type name")
		}
	}
}

func (ct CurveType) String() string {
	if ct <= 0 || curveTypeEnd <= ct {
		return "Unknown"
	}
	return curveTypeNames[ct]
}

type Curve interface {
	// IsOnCurve reports whether the given (x,y) lies on the curve.
	IsOnCurve(x, y *big.Int) bool

	// Add returns the sum of (x1,y1) and (x2,y2).
	Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int)

	// Double returns 2*(x,y).
	Double(x1, y1 *big.Int) (x, y *big.Int)

	// ScalarMult returns k*(x,y) where k is an integer in big-endian form.
	ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int)

	// ScalarBaseMult returns k*G, where G is the base point of the group
	// and k is an integer in big-endian form.
	ScalarBaseMult(k []byte) (x, y *big.Int)
}

// CurveParams holds the domain parameters of an elliptic curve.
//
// The struct is a union across several curve families; Type selects which
// fields are populated and how the remaining ones are interpreted:
//
//	Weierstrass (prime):  y² = x³ + Ax + B          over GF(P)
//	    uses A, B, P, N, Gx, Gy, InfX, InfY
//	Weierstrass (binary): y² + xy = x³ + A2x² + A6  over GF(2^m)
//	    uses A2, A6, Poly, N, Gx, Gy, InfX, InfY
//	Montgomery:           By² = x³ + Ax² + x        over GF(P)
//	    uses A, B, P, N, Gx, Gy
//	Edwards:              x² + y² = C²(1 + Dx²y²)   over GF(P)
//	    uses C, D, P, N, Gx, Gy
//	Twisted Edwards:      Ax² + y² = 1 + Dx²y²      over GF(P)
//	    uses A, D, P, N, Gx, Gy
//
// P is the field characteristic and is unused by binary curves, which are
// defined instead by Poly, the irreducible polynomial of GF(2^m) given as
// the ascending list of exponents with non-zero coefficients (e.g.
// {0, 1, 2, 7, 163} for x¹⁶³ + x⁷ + x² + x + 1). BitSize is the bit length
// of the underlying field: P for prime fields, m for binary fields.
//
// (Gx, Gy) is the base point, whose order is N. InfX and InfY carry the
// encoding of the point at infinity for curve forms that need one; Edwards
// and Montgomery curves have an affine identity element and leave them nil.
//
// Fields are treated as read-only once the curve is constructed. Callers
// must not mutate the big.Int values, as they may be shared across curves
// and derived parameter sets.
type CurveParams struct {
	Type       CurveType
	Name       string   // the name of the curve
	BitSize    int      // the size of the underlying field
	A, B, C, D *big.Int // curve parameters for prime curves
	A2, A6     *big.Int // curve parameters for binary curves
	P          *big.Int // Prime defining the field
	Poly       []int    // Irreducible polynomial for binary fields, in ascending order of degree
	N          *big.Int // order of the base point
	Gx, Gy     *big.Int // coordinates of the generator point
	InfX, InfY *big.Int // coordinates of the point at infinity
}

// Params converts the curve parameters into a standard library
// *elliptic.CurveParams for compatibility with code that expects the
// elliptic.Curve interface.
//
// The conversion is lossy: only Name, BitSize, P, N, B, Gx and Gy are
// carried over. Parameters that have no counterpart in the standard
// library (A, C, D, A2, A6, Poly, InfX, InfY) are dropped.
//
// Because elliptic.CurveParams assumes the short Weierstrass form
// y² = x³ - 3x + B over a prime field, the returned value is only
// meaningful for curves matching that form. For other curves — notably
// binary curves — the result should be treated as metadata only, and its
// arithmetic methods (Add, Double, ScalarMult, IsOnCurve) must not be used.
//
// The returned value is a new struct, but the big.Int fields are shared
// with p and must not be mutated by the caller.
func (p *CurveParams) Params() *elliptic.CurveParams {
	return &elliptic.CurveParams{
		Name:    p.Name,
		BitSize: p.BitSize,
		P:       p.P,
		N:       p.N,
		B:       p.B,
		Gx:      p.Gx,
		Gy:      p.Gy,
	}
}

type CurveExtended interface {
	Curve

	// Params2 returns the parameters of the curve.
	// It is not compatible with [elliptic.Curve.Params].
	// - The returned value is a safe copy and can be modified by the caller.
	Params2() *CurveParams

	// HasGenerator reports whether the curve has a defined generator point.
	// If not, ScalarMult and ScalarBaseMult cannot be used.
	HasGenerator() bool

	// OID returns the object identifier of the curve.
	// The returned value is a safe copy and can be modified by the caller.
	// If the curve has no OID, nil is returned.
	OID() asn1.ObjectIdentifier

	// ComputeY computes the y coordinate for the given x coordinate.
	// If there is no valid y for the given x, nil is returned.
	// If largeY is true, the larger one is returned instead of the smaller one.
	ComputeY(x *big.Int, largeY bool) *big.Int

	// IsEquivalentCurve reports whether the given curve has equivalent mathematical
	// curve parameters. Curve names and OIDs are intentionally ignored.
	// This method is not compatible with [elliptic.Curve].
	IsEquivalentCurve(other Curve) bool

	// IsEquivalentCurveWithOID reports whether the given curve has equivalent
	// mathematical curve parameters and the same object identifier.
	// This method is not compatible with [elliptic.Curve].
	IsEquivalentCurveWithOID(other Curve) bool

	elliptic.Curve
}
