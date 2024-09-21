package weierstrassbinary

import (
	"math/big"

	"github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/RyuaNerin/elliptic2/internal/field"
)

type affineOp struct {
	*Curve
	t [2]field.GF2m
}

var _ curve.GF2mOperator = (*affineOp)(nil)

func SelectAffine(params *CurveParams) NewOpFunc {
	return NewAffine
}

func NewAffine(c *Curve) curve.GF2mOperator {
	op := new(affineOp)
	op.Curve = c
	for idx := range op.t {
		op.t[idx].SetModulus(c.Poly)
	}
	return op
}

func (op *affineOp) SetInfinity(coords *curve.GF2mCoordinate) {
	coords[0].SetInt64(0)
	coords[1].SetInt64(0)
}

func (op *affineOp) IsInfinity(coords *curve.GF2mCoordinate) bool {
	return coords[0].IsZero()
}

func (op *affineOp) ToCoordinate(coords *curve.GF2mCoordinate, x, y *big.Int) {
	coords[0].SetBigInt(x)
	coords[1].SetBigInt(y)
}

func (op *affineOp) ToAffinePoint(x, y *big.Int, coords *curve.GF2mCoordinate) {
	coords[0].ToBigInt(x)
	coords[1].ToBigInt(y)
}

func (op *affineOp) Add(dst, p1, p2 *curve.GF2mCoordinate) { op.add(dst, p1, p2) }
func (op *affineOp) Double(dst, p1 *curve.GF2mCoordinate)  { op.dbl(dst, p1) }
