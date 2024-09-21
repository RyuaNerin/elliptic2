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

type CurveParams struct {
	Type       CurveType
	Name       string   // the name of the curve
	BitSize    int      // the size of the underlying field
	A, B, C, D *big.Int // curve parameters for prime curves
	A2, A6     *big.Int // curve parameters for binary curves
	P          *big.Int // Prime or polynomial defining the field
	N          *big.Int // order of the base point
	Gx, Gy     *big.Int // coordinates of the generator point
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

	elliptic.Curve
}
