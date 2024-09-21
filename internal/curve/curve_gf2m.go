package curve

import (
	"crypto/elliptic"
	"fmt"
	"math/big"

	"github.com/RyuaNerin/elliptic2"
)

type (
	GF2mCurveArithmetic CurveArithmetic[GF2mCoordinate, *GF2mCoordinate, GF2mOperator]
	GF2mOperator        Operator[GF2mCoordinate, *GF2mCoordinate]
	GF2mMaddOperator    MaddOperator[GF2mCoordinate, *GF2mCoordinate]

	curveGF2m struct {
		curveBase
		base GF2mCurveArithmetic
	}
)

var (
	_ elliptic2.CurveExtended = (*curveGF2m)(nil)
	_ elliptic.Curve          = (*curveGF2m)(nil)
)

func (c *curveGF2m) Params() *elliptic.CurveParams   { return c.base.Params() }
func (c *curveGF2m) Params2() *elliptic2.CurveParams { return c.base.Params2() }

func (c *curveGF2m) IsOnCurve(x, y *big.Int) bool { return c.base.IsOnCurve(x, y) }

func (c *curveGF2m) HasGenerator() bool {
	_, _, ok := c.base.Generator()
	return ok
}

func (c *curveGF2m) ComputeY(x *big.Int, largeY bool) *big.Int { return c.base.ComputeY(x, largeY) }

func (c *curveGF2m) Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	c.panicIfNotOnCurve(x1, y1)
	c.panicIfNotOnCurve(x2, y2)

	if x1.Sign() == 0 && y1.Sign() == 0 {
		return new(big.Int).Set(x2), new(big.Int).Set(y2)
	}
	if x2.Sign() == 0 && y2.Sign() == 0 {
		return new(big.Int).Set(x1), new(big.Int).Set(y1)
	}

	op := c.base.NewOperator()

	var dst, p1 GF2mCoordinate
	dst.SetModulus(c.base.Modulus())
	p1.SetModulus(c.base.Modulus())

	op.ToCoordinate(&p1, x1, y1)

	if x1.Cmp(x2) == 0 && y1.Cmp(y2) == 0 {
		op.Double(&dst, &p1)
	} else {
		var p2 GF2mCoordinate
		p2.SetModulus(c.base.Modulus())

		op.ToCoordinate(&p2, x2, y2)
		op.Add(&dst, &p1, &p2)
	}

	x, y = new(big.Int), new(big.Int)
	op.ToAffinePoint(x, y, &dst)
	return x, y
}

func (c *curveGF2m) Double(x1, y1 *big.Int) (x, y *big.Int) {
	c.panicIfNotOnCurve(x1, y1)

	if x1.Sign() == 0 {
		return new(big.Int), new(big.Int)
	}

	op := c.base.NewOperator()

	var dst, p1 GF2mCoordinate
	dst.SetModulus(c.base.Modulus())
	p1.SetModulus(c.base.Modulus())

	op.ToCoordinate(&p1, x1, y1)

	op.Double(&dst, &p1)

	x, y = new(big.Int), new(big.Int)
	op.ToAffinePoint(x, y, &dst)
	return x, y
}

func (c *curveGF2m) ScalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	c.panicIfNotOnCurve(x1, y1)

	return c.scalarMult(x1, y1, k)
}

func (c *curveGF2m) ScalarBaseMult(k []byte) (x, y *big.Int) {
	x, y, ok := c.base.Generator()
	if !ok {
		panic("elliptic2: curve has no generator")
	}
	return c.scalarMult(x, y, k)
}

func (c *curveGF2m) scalarMult(x1, y1 *big.Int, k []byte) (x, y *big.Int) {
	if len(k) == 0 {
		return new(big.Int), new(big.Int)
	}

	var num big.Int
	num.SetBytes(k)
	if num.Sign() == 0 {
		return new(big.Int), new(big.Int)
	}

	op := c.base.NewOperator()

	var r0Value, r1Value, tmValue GF2mCoordinate
	r0, r1, tm := &r0Value, &r1Value, &tmValue

	r0.SetModulus(c.base.Modulus())
	r1.SetModulus(c.base.Modulus())
	tm.SetModulus(c.base.Modulus())

	op.ToCoordinate(r0, x1, y1)

	// Montgomery Ladder
	op.Double(r1, r0)
	for i := num.BitLen() - 2; i >= 0; i-- {
		if num.Bit(i) == 1 {
			op.Add(tm, r0, r1)
			r0, tm = tm, r0

			op.Double(tm, r1)
			r1, tm = tm, r1
		} else {
			op.Add(tm, r1, r0)
			r1, tm = tm, r1

			op.Double(tm, r0)
			r0, tm = tm, r0
		}
	}

	x, y = new(big.Int), new(big.Int)
	op.ToAffinePoint(x, y, r0)
	return x, y
}

func (c *curveGF2m) panicIfNotOnCurve(x, y *big.Int) {
	if x.Sign() == 0 && y.Sign() == 0 {
		return
	}

	if !c.base.IsOnCurve(x, y) {
		panic(fmt.Sprintf("elliptic2: point (%s, %s) is not on curve %s", x.Text(16), y.Text(16), c.base.Params().Name))
	}
}
