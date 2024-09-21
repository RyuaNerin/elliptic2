package twistededwards

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

func SelectAffine(params *CurveParams) NewOpFunc {
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

func (op *affineOp) SetInfinity(coords *curve.GFpCoordinate) {
	coords[0].SetInt64(0)
	coords[1].SetInt64(0)
}

func (op *affineOp) IsInfinity(coords *curve.GFpCoordinate) bool {
	return coords[0].Sign() == 0
}

func (op *affineOp) ToCoordinate(coords *curve.GFpCoordinate, x, y *big.Int) {
	coords[0].SetBigInt(x)
	coords[1].SetBigInt(y)
}

func (op *affineOp) ToAffinePoint(x, y *big.Int, coords *curve.GFpCoordinate) {
	coords[0].ToBigInt(x)
	coords[1].ToBigInt(y)
}

func (op *affineOp) Add(dst, p1, p2 *curve.GFpCoordinate) { op.add(dst, p1, p2) }
func (op *affineOp) Double(dst, p1 *curve.GFpCoordinate)  { op.dbl(dst, p1) }
