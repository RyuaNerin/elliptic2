package weierstrassbinary

import (
	"math/big"

	"github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/RyuaNerin/elliptic2/internal/field"
)

type (
	lambdaOp struct {
		*Curve
		t [3]field.GF2m
	}
)

var _ curve.GF2mOperator = (*lambdaOp)(nil)

func SelectLambdaOp(params *CurveParams) NewOpFunc {
	return NewLambda
}

// general a
func NewLambda(c *Curve) curve.GF2mOperator {
	op := new(lambdaOp)
	op.Curve = c
	for idx := range op.t {
		op.t[idx].SetModulus(c.Poly)
	}
	return op
}

func (op *lambdaOp) SetInfinity(coords *curve.GF2mCoordinate) {
	// X L Z = 0 1 0
	coords[0].SetInt64(0)
	coords[1].SetInt64(1)
	coords[2].SetInt64(0)
}

func (op *lambdaOp) IsInfinity(coords *curve.GF2mCoordinate) bool {
	return coords[2].IsZero() // Z == 0
}

func (op *lambdaOp) ToCoordinate(coords *curve.GF2mCoordinate, xp, yp *big.Int) {
	x, y := &op.t[0], &op.t[1]
	x.SetBigInt(xp)
	y.SetBigInt(yp)

	X3, L3, Z3 := &coords[0], &coords[1], &coords[2]

	X3.Set(x)                       // X = x
	L3.Add(x, L3.Mul(y, L3.Inv(x))) // L = x + y/x
	Z3.SetInt64(1)                  // Z = 1
}

func (op *lambdaOp) ToAffinePoint(x, y *big.Int, coords *curve.GF2mCoordinate) {
	/**
	x=X/Z
	y/x=(L-X)/Z

	T = 1 / Z
	x = X * T
	y = x * (L - X) * T
	*/
	if op.IsInfinity(coords) {
		x.SetInt64(0)
		y.SetInt64(0)
		return
	}

	X1, L1, Z1 := &coords[0], &coords[1], &coords[2]
	X3, Y3, t0 := &op.t[0], &op.t[1], &op.t[2]

	t0.Inv(Z1)                             //               1 / Z
	X3.Mul(X1, t0)                         // x = X           / Z
	Y3.Mul(Y3.Mul(X3, Y3.Sub(L1, X1)), t0) // y = x * (L - X) / Z

	X3.ToBigInt(x)
	Y3.ToBigInt(y)
}

func (op *lambdaOp) Add(dst, p1, p2 *curve.GF2mCoordinate) { op.add2013olar(dst, p1, p2) }
func (op *lambdaOp) Double(dst, p1 *curve.GF2mCoordinate)  { op.dbl2013olar(dst, p1) }
