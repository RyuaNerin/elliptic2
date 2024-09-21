package twistededwards

import (
	"math/big"

	"github.com/RyuaNerin/elliptic2/internal/curve"
	"github.com/RyuaNerin/elliptic2/internal/field"
)

type (
	// X Y Z T
	extendedOp struct {
		*Curve
		t [2]field.GFp
	}
	ExtendedOp struct {
		extendedOp
	}
	ExtendedOpAm1 struct {
		extendedOp
	}
)

var (
	_ curve.GFpOperator = (*ExtendedOp)(nil)
	_ curve.GFpOperator = (*ExtendedOpAm1)(nil)
)

func SelectExtended(params *CurveParams) NewOpFunc {
	// (A + 1) % P == 0  <=>  A == -1 mod P
	var tmp field.GFp
	tmp.SetModulus(params.P)
	tmp.Add(params.A, field.GFpOne).Reduce()
	if tmp.Sign() == 0 {
		return NewExtendedAm1
	}
	return NewExtended
}

func NewExtended(c *Curve) curve.GFpOperator {
	op := new(ExtendedOp)
	op.Curve = c
	for idx := range op.t {
		op.t[idx].SetModulus(c.P)
	}
	return op
}

func NewExtendedAm1(c *Curve) curve.GFpOperator {
	op := new(ExtendedOpAm1)
	op.Curve = c
	for idx := range op.t {
		op.t[idx].SetModulus(c.P)
	}
	return op
}

func (op *extendedOp) SetInfinity(coords *curve.GFpCoordinate) {
	// X Y = 0 1 (affine)
	// -> X Y Z T = 0 1 1 0 (extended)
	coords[0].SetInt64(0)
	coords[1].SetInt64(1)
	coords[2].SetInt64(1)
	coords[3].SetInt64(0)
}

func (op *extendedOp) IsInfinity(coords *curve.GFpCoordinate) bool {
	// X == 0 && Y == 1
	return coords[0].Sign() == 0 && coords[1].Cmp(field.GFpOne) == 0
}

func (op *extendedOp) ToCoordinate(coords *curve.GFpCoordinate, x, y *big.Int) {
	coords[0].SetBigInt(x)                         // X = x
	coords[1].SetBigInt(y)                         // Y = y
	coords[2].SetInt64(1)                          // Z = 1
	coords[3].Mul(&coords[0], &coords[1]).Reduce() // T = X * Y
}

func (op *extendedOp) ToAffinePoint(x, y *big.Int, coords *curve.GFpCoordinate) {
	if op.IsInfinity(coords) {
		x.SetInt64(0)
		y.SetInt64(1)
		return
	}

	X1, Y1, Z1 := &coords[0], &coords[1], &coords[2]

	X3, Y3 := &op.t[0], &op.t[1]

	Y3.Inv(Z1)              // t0 = 1 / Z
	X3.Mul(X1, Y3).Reduce() // x  = X / Z
	Y3.Mul(Y1, Y3).Reduce() // y  = Y / Z

	X3.ToBigInt(x)
	Y3.ToBigInt(y)
}

func (op *extendedOp) ScaleZ(coords *curve.GFpCoordinate) {
	X1, Y1, Z1, T1 := &coords[0], &coords[1], &coords[2], &coords[3]

	t0 := op.t[0]

	t0.Inv(Z1)               //     1 / Z
	X1.Mul(X1, &t0).Reduce() // X = X / Z
	Y1.Mul(Y1, &t0).Reduce() // Y = Y / Z
	Z1.SetInt64(1)           // Z = 1
	T1.Mul(X1, Y1).Reduce()  // T = X * Y
}

func (op *extendedOp) Neg(dst, p1 *curve.GFpCoordinate) {
	// -P = (-x, y) in affine
	// T = X * Y
	// --> -P = (-X, Y, Z, -T) in extended
	dst[0].Neg(&p1[0]).Reduce() // X = -x
	dst[1].Set(&p1[1])          // Y =  y
	dst[2].Set(&p1[2])          // Z =  z
	dst[3].Neg(&p1[3]).Reduce() // T = -t
}

func (op *ExtendedOp) Add(dst, p1, p2 *curve.GFpCoordinate)  { op.add2008hwcd2(dst, p1, p2) }
func (op *ExtendedOp) Double(dst, p1 *curve.GFpCoordinate)   { op.dbl2008hwcd(dst, p1) }
func (op *ExtendedOp) Madd(dst, p1, p2 *curve.GFpCoordinate) { op.madd2008hwcd2(dst, p1, p2) }

func (op *ExtendedOpAm1) Add(dst, p1, p2 *curve.GFpCoordinate)  { op.add2008hwcd2(dst, p1, p2) }
func (op *ExtendedOpAm1) Double(dst, p1 *curve.GFpCoordinate)   { op.dbl2008hwcd(dst, p1) }
func (op *ExtendedOpAm1) Madd(dst, p1, p2 *curve.GFpCoordinate) { op.madd2008hwcd4(dst, p1, p2) }
