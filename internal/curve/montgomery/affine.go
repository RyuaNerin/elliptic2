package montgomery

import (
	"math/big"

	"github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/RyuaNerin/elliptic2/internal/field"
)

type affineOp struct {
	*Curve
	t [2]field.GFp
}

var _ curve.GFpOperator = (*affineOp)(nil)

func SelectAffine(*CurveParams) NewOpFunc {
	return NewAffine
}

func NewAffine(c *Curve) curve.GFpOperator {
	op := new(affineOp)
	op.Curve = c
	for idx := range op.t {
		op.t[idx].SetModulus(c.P)
	}
	return op
}

func (affineOp) SetInfinity(coords *curve.GFpCoordinate) {
	coords[0].SetInt64(0)
	coords[1].SetInt64(1)
}

func (affineOp) IsInfinity(coords *curve.GFpCoordinate) bool {
	return coords[0].Sign() == 0 && coords[1].Cmp(field.GFpOne) == 0
}

func (op *affineOp) ToCoordinate(coords *curve.GFpCoordinate, x, y *big.Int) {
	if x.Sign() == 0 && y.Sign() == 0 {
		op.SetInfinity(coords)
		return
	}
	coords[0].SetBigInt(x)
	coords[1].SetBigInt(y)
}

func (op *affineOp) ToAffinePoint(x, y *big.Int, coords *curve.GFpCoordinate) {
	if op.IsInfinity(coords) {
		x.SetInt64(0)
		y.SetInt64(0)
		return
	}
	coords[0].ToBigInt(x)
	coords[1].ToBigInt(y)
}

func (op *affineOp) Add(dst, p1, p2 *curve.GFpCoordinate) {
	switch {
	case p1[0].Cmp(&p2[0]) != 0:
		op.add(dst, p1, p2)

	case p1[1].Cmp(&p2[1]) == 0:
		op.Double(dst, p1)

	default:
		op.SetInfinity(dst)
	}
}

func (op *affineOp) Double(dst, p1 *curve.GFpCoordinate) {
	if p1[1].IsZero() {
		op.SetInfinity(dst)
		return
	}
	op.dbl(dst, p1)
}
