package weierstrassbinary

import (
	"math/big"

	"github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/RyuaNerin/elliptic2/internal/field"
)

type jacobianOp struct {
	*Curve
	t [4]field.GF2m
}

var (
	_ curve.GF2mOperator     = (*jacobianOp)(nil)
	_ curve.GF2mMaddOperator = (*jacobianOp)(nil)
)

func SelectJacobianOp(*CurveParams) NewOpFunc {
	return NewJacobian
}

func NewJacobian(c *Curve) curve.GF2mOperator {
	op := new(jacobianOp)
	op.Curve = c
	for idx := range op.t {
		op.t[idx].SetModulus(c.Poly)
	}
	return op
}

func (op *jacobianOp) SetInfinity(coords *curve.GF2mCoordinate) {
	// X Y Z = 1 1 0
	coords[0].SetInt64(1)
	coords[1].SetInt64(1)
	coords[2].SetInt64(0)
}

func (op *jacobianOp) IsInfinity(coords *curve.GF2mCoordinate) bool {
	return coords[2].IsZero() // Z == 0
}

func (op *jacobianOp) ToCoordinate(coords *curve.GF2mCoordinate, x, y *big.Int) {
	coords[0].SetBigInt(x)
	coords[1].SetBigInt(y)
	coords[2].SetInt64(1)
}

func (op *jacobianOp) ToAffinePoint(x, y *big.Int, coords *curve.GF2mCoordinate) {
	if op.IsInfinity(coords) {
		x.SetInt64(0)
		y.SetInt64(0)
		return
	}

	X1, Y1, Z1 := &coords[0], &coords[1], &coords[2]

	X3, Y3, t0 := &op.t[0], &op.t[1], &op.t[2]

	Y3.Inv(Z1)     //     1 /   Z
	t0.Sqr(Y3)     //     1 /  ZZ
	Y3.Mul(t0, Y3) //     1 / ZZZ
	X3.Mul(X1, t0) // x = X /  ZZ
	Y3.Mul(Y1, Y3) // y = Y / ZZZ

	X3.ToBigInt(x)
	Y3.ToBigInt(y)
}

func (op *jacobianOp) ScaleZ(p *curve.GF2mCoordinate) {
	X1, Y1, Z1 := &p[0], &p[1], &p[2]

	t0, t1 := &op.t[0], &op.t[1]

	t1.Inv(Z1)     //     1 /   Z
	t0.Sqr(t1)     //     1 /  ZZ
	t1.Mul(t0, t1) //     1 / ZZZ
	X1.Mul(X1, t0) // X = X /  ZZ
	Y1.Mul(Y1, t1) // Y = Y / ZZZ
	Z1.SetInt64(1) // Z = 1
}

func (op *jacobianOp) Neg(dst, coords *curve.GF2mCoordinate) {
	// -P = (X, X+Y, Z)
	dst[0].Set(&coords[0])
	dst[1].Add(&coords[0], &coords[1])
	dst[2].Set(&coords[2])
}

func (op *jacobianOp) Add(dst, p1, p2 *curve.GF2mCoordinate)  { op.add2005dl(dst, p1, p2) }
func (op *jacobianOp) Double(dst, p1 *curve.GF2mCoordinate)   { op.dbl2005dl(dst, p1) }
func (op *jacobianOp) Madd(dst, p1, p2 *curve.GF2mCoordinate) { op.madd2008bl(dst, p1, p2) }
