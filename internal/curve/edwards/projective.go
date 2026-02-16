package edwards

import (
	"math/big"

	"github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/RyuaNerin/elliptic2/internal/field"
)

type projectiveOp struct {
	*Curve
	t [4]field.GFp
}

var _ curve.GFpOperator = (*projectiveOp)(nil)

func SelectProjectiveOp(*CurveParams) func(c *Curve) curve.GFpOperator {
	return NewProjective
}

func NewProjective(c *Curve) curve.GFpOperator {
	op := new(projectiveOp)
	op.Curve = c
	for idx := range op.t {
		op.t[idx].SetModulus(c.P)
	}
	return op
}

func (op *projectiveOp) SetInfinity(coords *curve.GFpCoordinate) {
	// X Y = 0 1 (affine)
	// -> X Y Z = 0 1 1 (projective)
	coords[0].SetInt64(0)
	coords[1].SetInt64(1)
	coords[2].SetInt64(1)
}

func (op *projectiveOp) IsInfinity(coords *curve.GFpCoordinate) bool {
	// X == 0 && Y == 1
	return coords[0].Sign() == 0 && coords[1].Cmp(field.GFpOne) == 0
}

func (op *projectiveOp) ToCoordinate(coords *curve.GFpCoordinate, x, y *big.Int) {
	// affine (x, y) -> projective (X, Y, Z=1)
	coords[0].SetBigInt(x)
	coords[1].SetBigInt(y)
	coords[2].SetInt64(1)
}

func (op *projectiveOp) ToAffinePoint(x, y *big.Int, coords *curve.GFpCoordinate) {
	if op.IsInfinity(coords) {
		x.SetInt64(0)
		y.SetInt64(1)
		return
	}

	X, Y, Z := &coords[0], &coords[1], &coords[2]

	dstX, dstY := &op.t[0], &op.t[1]

	if dstY.Inv(Z) == nil {
		x.SetInt64(0)
		y.SetInt64(1)
		return
	} // dstY = 1 / Z
	dstX.Mul(X, dstY).Reduce() // dstX = X / Z
	dstY.Mul(Y, dstY).Reduce() // dstY = Y / Z

	dstX.ToBigInt(x)
	dstY.ToBigInt(y)
}

func (op *projectiveOp) ScaleZ(coords *curve.GFpCoordinate) {
	/**
	(x, y, z) -> (x/z, y/z, 1)
	*/

	X, Y, Z := &coords[0], &coords[1], &coords[2]

	zInv := op.t[0]

	if zInv.Inv(Z) == nil {
		op.SetInfinity(coords)
		return
	}
	X.Mul(X, &zInv).Reduce() // X1 = X / Z
	Y.Mul(Y, &zInv).Reduce() // Y1 = Y / Z
	Z.SetInt64(1)            // Z1 = 1
}

func (op *projectiveOp) Neg(dst, p1 *curve.GFpCoordinate) {
	// -P = (-x, y)
	// --> -P = (-X, Y, Z) in projective
	dst[0].Neg(&p1[0]).Reduce()
	dst[1].Set(&p1[1])
	dst[2].Set(&p1[2])
}

func (op *projectiveOp) Add(dst, p1, p2 *curve.GFpCoordinate)  { op.add2007bl2(dst, p1, p2) }
func (op *projectiveOp) Double(dst, p1 *curve.GFpCoordinate)   { op.dbl2007bl2(dst, p1) }
func (op *projectiveOp) Madd(dst, p1, p2 *curve.GFpCoordinate) { op.madd2007bl2(dst, p1, p2) }
