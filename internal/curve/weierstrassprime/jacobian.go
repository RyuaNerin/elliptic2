package weierstrassprime

import (
	"math/big"

	"github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/RyuaNerin/elliptic2/internal/field"
)

type (
	jacobianOp struct {
		*Curve
		t [4]field.GFp
	}
	JacobianOp    struct{ jacobianOp }
	JacobianOpA0  struct{ jacobianOp }
	JacobianOpAm3 struct{ jacobianOp }
)

var (
	_ curve.GFpOperator     = (*JacobianOp)(nil)
	_ curve.GFpOperator     = (*JacobianOpA0)(nil)
	_ curve.GFpOperator     = (*JacobianOpAm3)(nil)
	_ curve.GFpMaddOperator = (*JacobianOp)(nil)
	_ curve.GFpMaddOperator = (*JacobianOpA0)(nil)
	_ curve.GFpMaddOperator = (*JacobianOpAm3)(nil)
)

func SelectJacobian(params *CurveParams) NewOpFunc {
	if params.A.Cmp(field.GFpZero) == 0 {
		return NewJacobianA0
	}

	// (A + 3) % P == 0  <=>  A == -3 mod P
	var tmp field.GFp
	tmp.SetModulus(params.P).Add(params.A, field.GFpThree).Reduce()
	if tmp.IsZero() {
		return NewJacobianAm3
	}

	return NewJacobian
}

// general a
func NewJacobian(c *Curve) curve.GFpOperator {
	op := new(JacobianOp)
	op.Curve = c
	for idx := range op.t {
		op.t[idx].SetModulus(c.P)
	}
	return op
}

// a = 0
func NewJacobianA0(c *Curve) curve.GFpOperator {
	op := new(JacobianOpA0)
	op.Curve = c
	for idx := range op.t {
		op.t[idx].SetModulus(c.P)
	}
	return op
}

// a = -3
func NewJacobianAm3(c *Curve) curve.GFpOperator {
	op := new(JacobianOpAm3)
	op.Curve = c
	for idx := range op.t {
		op.t[idx].SetModulus(c.P)
	}
	return op
}

func (jacobianOp) SetInfinity(coords *curve.GFpCoordinate) {
	// X Y Z = 0 1 0
	coords[0].SetInt64(0) // X == 0
	coords[1].SetInt64(1) // Y == 1
	coords[2].SetInt64(0) // Z == 0
}

func (jacobianOp) IsInfinity(coords *curve.GFpCoordinate) bool {
	return coords[2].Sign() == 0 // Z == 0
}

func (jacobianOp) ToCoordinate(coords *curve.GFpCoordinate, x, y *big.Int) {
	coords[0].SetBigInt(x)
	coords[1].SetBigInt(y)
	coords[2].SetInt64(1)
}

func (op *jacobianOp) ToAffinePoint(x, y *big.Int, coords *curve.GFpCoordinate) {
	if op.IsInfinity(coords) {
		x.SetInt64(0)
		y.SetInt64(0)
		return
	}

	X, Y, Z := &coords[0], &coords[1], &coords[2]
	dstX, dstY, t0, t1 := &op.t[0], &op.t[1], &op.t[2], &op.t[3]

	t1.Inv(Z)                   //     1 /   Z
	t0.Sqr(t1).ReduceSoft()     //     1 /  ZZ
	t1.Mul(t1, t0).ReduceSoft() //     1 / ZZZ
	dstX.Mul(X, t0).Reduce()    // x = X /  ZZ
	dstY.Mul(Y, t1).Reduce()    // y = Y / ZZZ

	dstX.ToBigInt(x)
	dstY.ToBigInt(y)
}

func (op *jacobianOp) ScaleZ(p *curve.GFpCoordinate) {
	X1, Y1, Z1 := &p[0], &p[1], &p[2]
	t0, t1 := &op.t[0], &op.t[1]

	t1.Inv(Z1)                  //     1 /   Z
	t0.Sqr(t1).ReduceSoft()     //     1 /  ZZ
	t1.Mul(t1, t0).ReduceSoft() //     1 / ZZZ
	X1.Mul(X1, t0).Reduce()     // x = X /  ZZ
	Y1.Mul(Y1, t1).Reduce()     // y = Y / ZZZ
	Z1.SetInt64(1)              // z = 1
}

func (op *jacobianOp) Neg(dst, p *curve.GFpCoordinate) {
	dst[0].Set(&p[0])
	dst[1].Neg(&p[1]).Reduce()
	dst[2].Set(&p[2])
}

func (op *JacobianOp) Add(dst, p1, p2 *curve.GFpCoordinate)  { op.add2007bl(dst, p1, p2) }
func (op *JacobianOp) Double(dst, p1 *curve.GFpCoordinate)   { op.dbl2007bl(dst, p1) }
func (op *JacobianOp) Madd(dst, p1, p2 *curve.GFpCoordinate) { op.madd2007bl(dst, p1, p2) }

func (op *JacobianOpA0) Add(dst, p1, p2 *curve.GFpCoordinate)  { op.add2007bl(dst, p1, p2) }
func (op *JacobianOpA0) Double(dst, p1 *curve.GFpCoordinate)   { op.dbl2009l(dst, p1) }
func (op *JacobianOpA0) Madd(dst, p1, p2 *curve.GFpCoordinate) { op.madd2007bl(dst, p1, p2) }

func (op *JacobianOpAm3) Add(dst, p1, p2 *curve.GFpCoordinate)  { op.add2007bl(dst, p1, p2) }
func (op *JacobianOpAm3) Double(dst, p1 *curve.GFpCoordinate)   { op.dbl2001b(dst, p1) }
func (op *JacobianOpAm3) Madd(dst, p1, p2 *curve.GFpCoordinate) { op.madd2007bl(dst, p1, p2) }
